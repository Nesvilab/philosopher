package data

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/cmsl/bio"
	"github.com/prvst/cmsl/data/fas"
	"github.com/prvst/cmsl/db"
	"github.com/prvst/philosopher/lib/meta"
	"github.com/prvst/philosopher/lib/sys"
)

// Base main structure
type Base struct {
	meta.Data
	ID        string
	Enz       string
	Tag       string
	UniProtDB string
	CrapDB    string
	Add       string
	Custom    string
	Annot     string
	Crap      bool
	Iso       bool
	Rev       bool
	TaDeDB    map[string]string
	Records   []db.Record
}

// New constructor
func New() Base {

	var o Base
	var m meta.Data
	m.Restore(sys.Meta())

	o.UUID = m.UUID
	o.Distro = m.Distro
	o.Home = m.Home
	o.MetaFile = m.MetaFile
	o.MetaDir = m.MetaDir
	o.DB = m.DB
	o.Temp = m.Temp
	o.TimeStamp = m.TimeStamp

	return o
}

// ProcessDB ...
func (d *Base) ProcessDB(file, decoyTag string) error {

	fastaMap, err := fas.ParseFile(file)
	if err != nil {
		return errors.New("Error parsing FASTA database")
	}

	for k, v := range fastaMap {

		class, err := db.Classify(k, decoyTag)
		if err != nil {
			return err
		}

		if class == "uniprot" {

			e, err := db.ProcessUniProtKB(k, v, decoyTag)
			if err != nil {
				return err
			}
			d.Records = append(d.Records, e)

		} else if class == "ncbi" {

			e, err := db.ProcessNCBI(k, v, decoyTag)
			if err != nil {
				return err
			}
			d.Records = append(d.Records, e)
		} else {
			return errors.New("Unknown database class")
		}

	}

	return nil
}

// Fetch downloads a database file from UniProt
func (d *Base) Fetch() error {

	var query string
	d.UniProtDB = fmt.Sprintf("%s%s%s.fas", d.Temp, string(filepath.Separator), d.ID)

	if d.Rev == true {
		query = fmt.Sprintf("%s%s%s", "http://www.uniprot.org/uniprot/?query=reviewed:yes+AND+proteome:", d.ID, "&format=fasta")
	} else {
		query = fmt.Sprintf("%s%s%s", "http://www.uniprot.org/uniprot/?query=proteome:", d.ID, "&format=fasta")
	}

	if d.Iso == true {
		query = fmt.Sprintf("%s&include=yes", query)
	} else {
		query = fmt.Sprintf("%s&include=no", query)
	}

	// tries to create an output file
	output, err := os.Create(d.UniProtDB)
	if err != nil {
		msg := fmt.Sprintf("Cannot create output file %s - %s", query, err)
		return errors.New(msg)
	}
	defer output.Close()

	// Tries to query data from Uniprot
	response, err := http.Get(query)
	if err != nil {
		msg := fmt.Sprintf("Cannot find annotation file %s", err)
		return errors.New(msg)
	}
	defer response.Body.Close()

	// if len(response.Body.) > 10 {
	// 	return errors.New("Cannot connect to UniProt, hceck your internet connection")
	// }

	// Tries to download data from Uniprot
	n, err := io.Copy(output, response.Body)
	if err != nil {
		msg := fmt.Sprintf("Cannot download annotation file %d - %s", n, err)
		return errors.New(msg)
	}

	return nil
}

// Create processes the given fasta file and add decoy sequences
func (d *Base) Create() error {

	d.TaDeDB = make(map[string]string)
	var crapMap = make(map[string]string)

	dbfile, _ := filepath.Abs(d.UniProtDB)
	db, err := fas.ParseFile(dbfile)
	if err != nil {
		return err
	}

	if len(d.Add) > 0 {
		add, adderr := fas.ParseFile(d.Add)
		if adderr != nil {
			return adderr
		}

		for k, v := range add {
			db[k] = v
		}
	}

	// adding contaminants to database before reversion
	// repeated entries are removed and substituted by contaminants
	if d.Crap == true {
		d.Deploy()
		//crapFile = Deploy(p)
		crapMap, err = fas.ParseFile(d.CrapDB)
		for k, v := range crapMap {
			split := strings.Split(k, "|")

			for i := range db {
				if strings.Contains(i, split[1]) {
					delete(db, i)
				}
			}

			db[k] = v

		}
	}

	var e bio.Enzyme
	e.Synth(d.Enz)
	reg := regexp.MustCompile(e.Pattern)

	for h, s := range db {

		th := ">" + h
		d.TaDeDB[th] = s

		var revPeptides []string
		split := reg.Split(s, -1)
		if len(split) > 1 {
			for i := range split {
				r := reverseSeq(split[i])
				revPeptides = append(revPeptides, r)
			}
		} else {
			r := reverseSeq(s)
			revPeptides = append(revPeptides, r)
		}

		rev := strings.Join(revPeptides, e.Join)
		dh := ">" + d.Tag + h
		d.TaDeDB[dh] = rev
	}

	return nil
}

// Deploy crap file to session folder
func (d *Base) Deploy() error {

	d.CrapDB = fmt.Sprintf("%s%scrap.fas", d.Temp, string(filepath.Separator))

	param, err := Asset("crap.fas")
	err = ioutil.WriteFile(d.CrapDB, param, 0644)

	if err != nil {
		msg := fmt.Sprintf("Could not deploy Crap fasta file %s", err)
		return errors.New(msg)
	}

	return nil
}

// Save fasta file to disk
func (d *Base) Save() error {

	base := filepath.Base(d.UniProtDB)

	t := time.Now()
	stamp := fmt.Sprintf(t.Format("2006-01-02"))

	workfile := fmt.Sprintf("%s%s%s-td-%s", d.Temp, string(filepath.Separator), stamp, base)
	outfile := fmt.Sprintf("%s%s%s-td-%s", d.Home, string(filepath.Separator), stamp, base)

	// create decoy db file
	file, err := os.Create(workfile)
	if err != nil {
		msg := fmt.Sprintf("Cannot create decoy database %s", err)
		return errors.New(msg)
	}
	defer file.Close()

	var headers []string
	for k := range d.TaDeDB {
		headers = append(headers, k)
	}

	sort.Strings(headers)

	for _, i := range headers {
		line := i + "\n" + d.TaDeDB[i] + "\n"
		_, err = io.WriteString(file, line)
		if err != nil {
			return errors.New("Cannot write database file")
		}
	}

	sys.CopyFile(workfile, outfile)

	d.ProcessDB(outfile, d.Tag)

	err = d.Serialize()
	if err != nil {
		logrus.Fatal(err)
	}

	return nil
}

// Serialize ...
func (d *Base) Serialize() error {

	var err error

	// create a file
	dataFile, err := os.Create(sys.DBBin())
	if err != nil {
		return err
	}

	dataEncoder := gob.NewEncoder(dataFile)
	goberr := dataEncoder.Encode(d)
	if goberr != nil {
		logrus.Fatal("Cannot save results, Bad format", goberr)
	}
	dataFile.Close()

	return nil
}

// Restore reads philosopher results files and restore the data sctructure
func (d *Base) Restore() error {

	file, _ := os.Open(sys.DBBin())

	dec := gob.NewDecoder(file)
	err := dec.Decode(&d)
	if err != nil {
		return errors.New("Could not restore Philosopher result. Please check file path")
	}

	return nil
}

// RestoreWithPath reads philosopher results files and restore the data sctructure
func (d *Base) RestoreWithPath(p string) error {

	var path string

	if strings.Contains(p, string(filepath.Separator)) {
		path = fmt.Sprintf("%s%s", p, sys.DBBin())
	} else {
		path = fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.DBBin())
	}

	file, _ := os.Open(path)

	dec := gob.NewDecoder(file)
	err := dec.Decode(&d)
	if err != nil {
		return errors.New("Could not restore Philosopher result. Please check file path")
	}

	return nil
}

// reverseSeq returns its argument string reversed rune-wise left to right.
func reverseSeq(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
