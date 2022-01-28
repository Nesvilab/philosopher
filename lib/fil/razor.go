package fil

import (
	"io/ioutil"
	"philosopher/lib/msg"
	"philosopher/lib/sys"

	"github.com/vmihailenco/msgpack"
)

// RazorCandidate is a peptide sequence to be evaluated as a razor
type RazorCandidate struct {
	Sequence          string
	MappedProtein     string
	MappedproteinsSID map[string]string
	MappedProteinsTNP map[string]int
	MappedProteinsW   map[string]float64
	MappedProteinsGW  map[string]float64
}

// a Map fo Razor candidates
type RazorMap map[string]RazorCandidate

// Serialize converts the razor structure to a gob file
func (p *RazorMap) Serialize() {

	b, e := msgpack.Marshal(&p)
	if e != nil {
		msg.MarshalFile(e, "fatal")
	}

	e = ioutil.WriteFile(sys.RazorBin(), b, sys.FilePermission())
	if e != nil {
		msg.WriteFile(e, "fatal")
	}

}

// Restore reads razor bin files and restore the data sctructure
func (p *RazorMap) Restore(silent bool) {

	b, e := ioutil.ReadFile(sys.RazorBin())
	if e != nil && !silent {
		msg.ReadFile(e, "warning")
	}

	e = msgpack.Unmarshal(b, &p)
	if e != nil && !silent {
		msg.DecodeMsgPck(e, "warning")
	}

}
