package fil

import (
	"philosopher/lib/sys"
)

// RazorCandidate is a peptide sequence to be evaluated as a razor
type RazorCandidate struct {
	Sequence          string
	MappedProteinsW   map[string]float64
	MappedProteinsGW  map[string]float64
	MappedProteinsTNP map[string]int
	MappedproteinsSID map[string]string
	MappedProtein     string
}

// a Map fo Razor candidates
type RazorMap map[string]RazorCandidate

// Serialize converts the razor structure to a gob file
func (p *RazorMap) Serialize() {
	sys.Serialize(p, sys.RazorBin())
}

// Restore reads razor bin files and restore the data sctructure
func (p *RazorMap) Restore(silent bool) {
	sys.Restore(p, sys.RazorBin(), silent)
}
