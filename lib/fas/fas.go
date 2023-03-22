package fas

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/Nesvilab/philosopher/lib/msg"
)

// ParseFile a fasta file and returns a map with the header as key and sequence as value
func ParseFile(filename string) map[string]string {

	var fastaMap = make(map[string]string)
	fastaSlice := ParseFile2(filename)
	for _, e := range fastaSlice {
		fastaMap[e.Header] = e.Seq
	}
	return fastaMap
}

type FastaEntry struct {
	Header string
	Seq    string
}

func ParseFile2(filename string) []FastaEntry {

	f, e := os.Open(filename)
	if filename == "" || e != nil {
		msg.ReadFile(errors.New("cannot open the database file"), "fatal")
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		panic(err)
	}

	fastaSlice := make([]FastaEntry, 0, stat.Size()/450)

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if len(scanner.Bytes()) > 0 && scanner.Bytes()[0] == '>' {
			line := scanner.Bytes()[1:]
			for i, e := range line {
				if e == '\t' {
					line[i] = ' '
				}
			}
			fastaSlice = append(fastaSlice, FastaEntry{Header: string(line), Seq: ""})
		} else {
			fastaSlice[len(fastaSlice)-1].Seq += scanner.Text()
		}
	}

	return fastaSlice
}

// ParseUniProtDescriptionMap parses a UniProt FASTA file and returns a map with ID and DESC
func ParseUniProtDescriptionMap(database string) (fastaMap map[string]string) {

	fastaMap = make(map[string]string)

	// parse fasta file
	file := ParseFile(database)
	faseq, _ := regexp.Compile(`\w+\|(.*?)\|(.*?)\s(.*)`)

	// get protein name and description and add them to fastaMap
	for k := range file {
		reg := faseq.FindStringSubmatch(k)
		desc := strings.Split(reg[3], "OS=")
		fastaMap[strings.TrimSpace(reg[1])] = strings.TrimSpace(desc[0])
	}

	return
}

// ParseUniProtSequencenMap parses a UniProt FASTA file and returns a map with ID and DESC
func ParseUniProtSequencenMap(database string) (fastaMap map[string]string) {

	fastaMap = make(map[string]string)

	// parse fasta file
	file := ParseFile(database)
	faseq, _ := regexp.Compile(`\w+\|(.*?)\|(.*?)\s(.*)`)

	// get protein name and description and add them to fastaMap
	for k, v := range file {
		reg := faseq.FindStringSubmatch(k)
		fastaMap[strings.TrimSpace(reg[1])] = v
	}

	return
}

// ParseFastaDescription a fasta file and returns a map with the header as key and sequence as value
func ParseFastaDescription(filename string) map[string][]string {

	f, e := os.Open(filename)
	if filename == "" || e != nil {
		msg.ReadFile(errors.New("cannot open FASTA file"), "error")
	}
	defer f.Close()

	reHeader, _ := regexp.Compile(`>\w+\|(.*?)\|(.*?)\s(.*)`)
	scanner := bufio.NewScanner(f)

	fastaMap := make(map[string][]string)

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), ">") {
			header := reHeader.FindStringSubmatch(scanner.Text())
			shortHeader := strings.Split(header[3], "OS=")
			var list []string
			list = append(list, shortHeader[0])
			fastaMap[header[1]] = list
		}
	}

	return fastaMap
}

// CleanDatabase removes decoys and contaminants
func CleanDatabase(db map[string]string, decoytag, contag string) (cleanMap map[string]string) {

	cleanMap = make(map[string]string)

	for k, v := range db {
		if !strings.Contains(k, decoytag) && !strings.Contains(k, contag) {
			cleanMap[k] = v
		}
	}

	return cleanMap
}
