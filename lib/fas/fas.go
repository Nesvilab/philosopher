package fas

import (
	"bufio"
	"errors"
	"log"
	"os"

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

	// increase buffer size to handle long FASTA headers
	const maxLineSize = 1024 * 1024 // 1MB
	buf := make([]byte, maxLineSize)
	scanner.Buffer(buf, maxLineSize)

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

	// check for scann error after loop in case the actual header length exceeds the maxLineSize
	if err := scanner.Err(); err != nil {
		log.Fatalf("Scanner error: %v", err)
	}

	return fastaSlice
}
