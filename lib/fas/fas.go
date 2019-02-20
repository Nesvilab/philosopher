package fas

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/prvst/philosopher/lib/err"
)

// ParseFile a fasta file and returns a map with the header as key and sequence as value
func ParseFile(filename string) (map[string]string, *err.Error) {

	var fastaHeader string
	var fastaSeq string
	var fastaMap = make(map[string]string)

	f, e := os.Open(filename)
	if filename == "" || e != nil {
		return fastaMap, &err.Error{Type: err.CannotParseFastaFile, Class: err.FATA, Argument: e.Error()}
	}
	defer f.Close()

	reHeader, _ := regexp.Compile("^>(.*)")
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), ">") {
			header := reHeader.FindStringSubmatch(scanner.Text())
			fastaHeader = header[1]
			fastaMap[fastaHeader] = ""
		} else {
			fastaSeq = fastaMap[fastaHeader]
			fastaSeq += scanner.Text()
			fastaMap[fastaHeader] = fastaSeq
		}
	}

	return fastaMap, nil
}

// ParseUniProtDescriptionMap parses a UniProt FASTA file and returns a map with ID and DESC
func ParseUniProtDescriptionMap(database string) (fastaMap map[string]string) {

	fastaMap = make(map[string]string)

	// parse fasta file
	file, _ := ParseFile(database)
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
	file, _ := ParseFile(database)
	faseq, _ := regexp.Compile(`\w+\|(.*?)\|(.*?)\s(.*)`)

	// get protein name and description and add them to fastaMap
	for k, v := range file {
		reg := faseq.FindStringSubmatch(k)
		fastaMap[strings.TrimSpace(reg[1])] = v
	}

	return
}

// ParseFastaDescription a fasta file and returns a map with the header as key and sequence as value
func ParseFastaDescription(filename string) (map[string][]string, *err.Error) {

	f, e := os.Open(filename)
	if filename == "" || e != nil {
		return nil, &err.Error{Type: err.CannotParseFastaFile, Class: err.FATA, Argument: e.Error()}
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

	return fastaMap, nil
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
