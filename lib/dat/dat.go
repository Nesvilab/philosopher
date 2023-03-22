// Package dat (Database)
package dat

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/Nesvilab/philosopher/lib/fas"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Nesvilab/philosopher/lib/msg"

	"github.com/Nesvilab/philosopher/lib/met"
	"github.com/Nesvilab/philosopher/lib/sys"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v5"
)

// Base main structure
type Base struct {
	FileName        string
	UniProtDB       string
	CrapDB          string
	Prefix          string
	Proteomes       string
	DownloadedFiles []string
	Records         []Record
	RecordsLen      int
	NParts          uint
	PartsLen        []int
	TaDeDB          map[string]string
}

// New constructor
func New() Base {

	var self Base

	self.TaDeDB = make(map[string]string)
	self.Records = []Record{}

	return self
}

// Run is the main entry point for the databse command
func Run(m met.Data) met.Data {

	var db = New()

	if len(m.Database.ID) == 0 && (len(m.Database.Annot) == 0 || m.Database.Annot == "--contam" || m.Database.Annot == "--prefix") && (len(m.Database.Custom) == 0 || m.Database.Custom == "--contam" || m.Database.Custom == "--prefix") {
		msg.InputNotFound(errors.New("provide a protein FASTA file or Proteome ID"), "error")
	}

	if len(m.Database.Annot) > 0 {

		logrus.Info("Annotating the database")

		m.DB = m.Database.Annot

		db.ProcessDB_and_serialize(m.Database.Annot, m.Database.Tag, m.Database.Verbose)

		db.Serialize()

		return m
	}

	if len(m.Database.ID) < 1 && len(m.Database.Custom) < 1 {
		msg.InputNotFound(errors.New("you need to provide a taxon ID or a custom FASTA file"), "error")
	}

	if !m.Database.Crap {
		msg.Custom(errors.New("contaminants are not going to be added to database"), "warning")
	}

	// bool variable will control the adition fo contaminant tags to contam proteins from the same organism.
	var ids = make(map[string]string)

	if len(m.Database.Custom) < 1 {

		m.DB = m.Database.Custom

		dbs := strings.Split(m.Database.ID, ",")
		for _, i := range dbs {

			organism, proteomeID := GetOrganismID(sys.GetTemp(), i)

			logrus.Info("Fetching ", organism, " database ", i)

			currentTime := time.Now()
			m.Database.TimeStamp = currentTime.Format("2006.01.02 15:04:05")

			db.Fetch(i, proteomeID, m.Temp, m.Database.Iso, m.Database.Rev)

			ids[proteomeID] = organism
		}

	} else {
		dbPath, _ := filepath.Abs(m.Database.Custom)
		db.UniProtDB = dbPath
		db.DownloadedFiles = append(db.DownloadedFiles, dbPath)
	}

	logrus.Info("Generating the target-decoy database")
	db.Create(m.Temp, m.Database.Add, m.Database.Enz, m.Database.Tag, m.Database.Crap, m.Database.NoD, m.Database.CrapTag, ids)

	logrus.Info("Creating file")
	customDB := db.Save(m.Home, m.Temp, m.Database.ID, m.Database.Tag, m.Database.Rev, m.Database.Iso, m.Database.NoD, m.Database.Crap)

	db.ProcessDB_and_serialize(customDB, m.Database.Tag, m.Database.Verbose)

	logrus.Info("Processing decoys")
	db.Create(m.Temp, m.Database.Add, m.Database.Enz, m.Database.Tag, m.Database.Crap, m.Database.NoD, m.Database.CrapTag, ids)

	logrus.Info("Creating file")
	db.Save(m.Home, m.Temp, m.Database.ID, m.Database.Tag, m.Database.Rev, m.Database.Iso, m.Database.NoD, m.Database.Crap)

	db.Prefix = m.Database.Tag

	db.Serialize()

	return m
}

// ProcessDB_and_serialize determines the type of sequence and sends it to the appropriate parsing function
func (d *Base) ProcessDB_and_serialize(filename, decoyTag string, verbose bool) {
	nproc := runtime.GOMAXPROCS(0)
	entriesChunk := make(chan []fas.FastaEntry, nproc)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		d.RecordsLen = ParseFile(filename, entriesChunk)
	}()
	var entriesChunkRecv <-chan []fas.FastaEntry = entriesChunk

	wg.Add(nproc)
	d.NParts = uint(nproc)
	d.PartsLen = make([]int, nproc)
	for i := 0; i < nproc; i++ {
		go func(i int) {
			defer wg.Done()
			output, e := os.OpenFile(fmt.Sprintf("%s__%d", sys.DBBin(), i), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, sys.FilePermission())
			defer func(output *os.File) {
				err := output.Close()
				if err != nil {
					msg.MarshalFile(err, "error")
					panic(err)
				}
			}(output)
			if e != nil {
				msg.WriteFile(e, "error")
				panic(e)
			}
			bo := bufio.NewWriter(output)
			defer func(bo *bufio.Writer) {
				err := bo.Flush()
				if err != nil {
					msg.MarshalFile(err, "error")
					panic(err)
				}
			}(bo)
			enc := msgpack.NewEncoder(bo)
			enc.UseInternedStrings(false)
			enc.UseArrayEncodedStructs(true)
			for fastaSlice := range entriesChunkRecv {
				d.PartsLen[i] += len(fastaSlice)
				for _, e := range fastaSlice {
					class := Classify(e.Header, decoyTag)
					enc.Encode(ProcessHeader(e.Header, e.Seq, class, decoyTag, verbose))
				}
			}
		}(i)
	}
	wg.Wait()
}
func ParseFile(filename string, entriesChunk chan<- []fas.FastaEntry) int {

	f, e := os.Open(filename)
	if filename == "" || e != nil {
		msg.ReadFile(errors.New("cannot open the database file"), "fatal")
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	n_entries := 0
	const chunk_size = 1 << 14
	fastaSlice := make([]fas.FastaEntry, 0, chunk_size)
	for scanner.Scan() {
		if len(scanner.Bytes()) > 0 && scanner.Bytes()[0] == '>' {
			line := scanner.Bytes()[1:]
			for i, e := range line {
				if e == '\t' {
					line[i] = ' '
				}
			}
			n_entries++
			if len(fastaSlice) == chunk_size {
				entriesChunk <- fastaSlice
				fastaSlice = make([]fas.FastaEntry, 0, chunk_size)
			}
			fastaSlice = append(fastaSlice, fas.FastaEntry{Header: string(line), Seq: ""})
		} else {
			fastaSlice[len(fastaSlice)-1].Seq += scanner.Text()
		}
	}
	entriesChunk <- fastaSlice
	close(entriesChunk)
	return n_entries
}

// Fetch downloads a database file from UniProt
func (d *Base) Fetch(uniprotID, proteomeID, temp string, iso, rev bool) {

	d.UniProtDB = fmt.Sprintf("%s%s%s.fas", temp, string(filepath.Separator), uniprotID)

	base := "https://rest.uniprot.org/uniprotkb/"

	// add the parameters
	query := base + "stream?compressed=false&format=fasta&"

	// add isoforms?
	if iso {
		query = query + "includeIsoform=true&"
	} else {
		query = query + "includeIsoform=false&"
	}

	// add the proteome parameter
	query = fmt.Sprintf("%squery=(proteome:%s)", query, uniprotID)

	// is reviewed?
	if rev {
		query = query + "+AND+(reviewed:true)"
	}

	client := resty.New()

	// HTTP response gets saved into file, similar to curl -o flag
	f := d.UniProtDB + ".gz"
	_, e := client.R().
		SetOutput(f).
		SetHeader("Accept-Encoding", "gzip,deflate").
		SetHeader("Content-Encoding", "gzip,deflate").
		SetHeader("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36").
		Get(query)

	if e != nil {
		msg.DatabaseNotFound(errors.New("cannot reach out UniProt database, are you connected to the Internet?"), "error")
	}

	file, err := os.Open(f)

	if err != nil {
		msg.ReadFile(errors.New("cannot open FASTA file, are you connected to the Internet?"), "error")
	}

	gz, err := gzip.NewReader(file)

	if err != nil && rev {
		msg.ReadFile(errors.New("please check if your parameters are correct, including the Uniprot ID and if the organism has reviewed sequences"), "error")
	} else if err != nil {
		msg.ReadFile(errors.New("please check if your parameters are correct, including the Uniprot ID"), "error")
	}

	defer file.Close()
	defer gz.Close()

	// tries to create an output file
	output, e := os.Create(d.UniProtDB)
	if e != nil {
		msg.WriteFile(errors.New("cannot create a local database file"), "error")
	}
	defer output.Close()

	// tries to download data from Uniprot
	_, e = io.Copy(output, gz)
	if e != nil {
		msg.Custom(errors.New("UniProt download failed, please check your connection"), "error")
	}

	d.DownloadedFiles = append(d.DownloadedFiles, d.UniProtDB)
}

// Create processes the given fasta file and add decoy sequences
func (d *Base) Create(temp, add, enz, tag string, crap, noD, cTag bool, ids map[string]string) {

	d.TaDeDB = make(map[string]string)

	for _, i := range d.DownloadedFiles {

		dbfile, _ := filepath.Abs(i)
		db := fas.ParseFile(dbfile)

		if len(add) > 0 {
			add := fas.ParseFile(add)

			for k, v := range add {
				db[k] = v
			}
		}

		// adding contaminants to database before reversion
		// repeated entries are removed and substituted by contaminants
		if crap {

			d.Deploy(temp)

			crapMap := fas.ParseFile(d.CrapDB)

			for k, v := range crapMap {
				for key := range ids {

					if cTag {
						if strings.Contains(k, key) {
							// Do not add contaminant tags to contam. proteins from the same organism
						} else {
							k = "contam_" + k
						}
					}

				}

				split := strings.Split(k, "|")
				for i := range db {
					if strings.Contains(i, split[1]) {
						delete(db, i)
					}
				}
				db[k] = v
			}

		}

		for h, s := range db {

			th := ">" + h
			d.TaDeDB[th] = s

			if !noD {
				dh := ">" + tag + h
				d.TaDeDB[dh] = reverseSeq(s)
			}

		}

	}

}

// Deploy crap file to session folder
func (d *Base) Deploy(temp string) {

	d.CrapDB = fmt.Sprintf("%s%scrap-gpmdb.fas", temp, string(filepath.Separator))

	param, e1 := Asset("crap-gpmdb.fas")
	if e1 != nil {
		msg.WriteFile(e1, "error")
	}

	e2 := os.WriteFile(d.CrapDB, param, sys.FilePermission())
	if e2 != nil {
		msg.WriteFile(e2, "error")
	}

}

// GetOrganismID maps the UniprotID to organismID
func GetOrganismID(temp string, uniprotID string) (string, string) {

	var proteomes = make(map[string]string)
	var organisms = make(map[string]string)
	proteomeFile := fmt.Sprintf("%s%sproteomes.csv", temp, string(filepath.Separator))

	param, e1 := Asset("proteomes.csv")
	if e1 != nil {
		msg.WriteFile(e1, "error")
	}

	e2 := os.WriteFile(proteomeFile, param, sys.FilePermission())
	if e2 != nil {
		msg.WriteFile(e2, "error")
	}

	f, e := os.Open(proteomeFile)
	if e != nil {
		log.Fatal(e)
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		organisms[parts[0]] = parts[1]
		proteomes[parts[0]] = parts[2]
	}

	return organisms[uniprotID], proteomes[uniprotID]
}

// Save fasta file to disk
func (d *Base) Save(home, temp, ids, tag string, isRev, hasIso, noD, Crap bool) string {

	var base string

	if len(ids) > 0 {
		base = strings.Replace(ids, ",", "-", -1)
	} else {
		base = filepath.Base(d.UniProtDB)
	}

	t := time.Now()
	stamp := t.Format("2006-01-02")

	baseName := fmt.Sprintf("%s%s", string(filepath.Separator), stamp)

	if !noD {
		baseName = baseName + "-decoys"
	}

	if isRev {
		baseName = baseName + "-reviewed"
	}

	if hasIso {
		baseName = baseName + "-isoforms"
	}

	if Crap {
		baseName = baseName + "-contam"
	}

	workfile := fmt.Sprintf("%s%s-%s.fas", temp, baseName, base)
	outfile := fmt.Sprintf("%s%s-%s.fas", home, baseName, base)

	// create db file
	file, e := os.Create(workfile)
	if e != nil {
		msg.ReadFile(errors.New("cannot open the database file"), "error")
	}
	defer file.Close()

	var headers []string
	for k := range d.TaDeDB {
		headers = append(headers, k)
	}

	sort.Strings(headers)

	for _, i := range headers {
		line := i + "\n" + d.TaDeDB[i] + "\n"
		_, e = io.WriteString(file, line)
		if e != nil {
			msg.WriteFile(e, "error")
		}
	}

	sys.CopyFile(workfile, outfile)

	return outfile
}

// Serialize saves to disk a msgpack version of the database data structure
func (d *Base) Serialize() {
	output, e := os.OpenFile(sys.DBBin(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, sys.FilePermission())
	defer func(output *os.File) {
		err := output.Close()
		if err != nil {
			msg.MarshalFile(err, "error")
			panic(err)
		}
	}(output)
	if e != nil {
		msg.WriteFile(e, "error")
		panic(e)
	}
	bo := bufio.NewWriter(output)
	defer func(bo *bufio.Writer) {
		err := bo.Flush()
		if err != nil {
			msg.MarshalFile(err, "error")
			panic(err)
		}
	}(bo)
	enc := msgpack.NewEncoder(bo)
	enc.UseInternedStrings(false)
	enc.UseArrayEncodedStructs(true)
	err := enc.Encode(&d)
	bo.Flush()
	if err != nil {
		msg.MarshalFile(err, "error")
		panic(err)
	}
}

// Restore reads philosopher results files and restore the data sctructure
func (d *Base) Restore() {
	d.restoreImpl(sys.DBBin())
}

func (d *Base) restoreImpl(filename string) {
	input, e := os.Open(filename)
	if e != nil {
		msg.ReadFile(e, "error")
		panic(e)
	}
	bi := bufio.NewReader(input)
	dec := msgpack.NewDecoder(bi)
	dec.UseInternedStrings(false)
	err := dec.Decode(&d)
	errClose := input.Close()
	if errClose != nil {
		panic(errClose)
	}
	if err != nil {
		msg.DecodeMsgPck(err, "error")
		panic(err)
	}
	d.Records = make([]Record, d.RecordsLen)
	var wg sync.WaitGroup
	wg.Add(int(d.NParts))
	partsLen_cumsum := make([]int, len(d.PartsLen)+1)
	for i, e := range d.PartsLen {
		partsLen_cumsum[i+1] = e + partsLen_cumsum[i]
	}
	for i := uint(0); i < d.NParts; i++ {
		go func(i uint) {
			defer wg.Done()
			fn := fmt.Sprintf("%s__%d", filename, i)
			input, e := os.OpenFile(fn, os.O_RDONLY, sys.FilePermission())
			if e != nil {
				msg.ReadFile(e, "error")
				panic(e)
			}
			bi := bufio.NewReader(input)
			dec := msgpack.NewDecoder(bi)
			start, end := partsLen_cumsum[i], partsLen_cumsum[i+1]
			for idx := start; idx < end; idx++ {
				err := dec.Decode(&d.Records[idx])
				if err != nil {
					panic(err)
				}
			}
		}(i)
	}
	wg.Wait()
}

// RestoreWithPath reads philosopher results files and restore the data sctructure
func (d *Base) RestoreWithPath(p string) {

	path := fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.DBBin())
	path, _ = filepath.Abs(path)
	d.restoreImpl(path)
}

// reverseSeq returns its argument string reversed rune-wise left to right.
func reverseSeq(s string) string {

	var index = 0
	r := []rune(s)

	if strings.HasPrefix(s, "M") {
		index = 1
	}

	for i, j := index, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
