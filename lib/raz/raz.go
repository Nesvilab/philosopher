package raz

import (
	"github.com/Nesvilab/philosopher/lib/sys"
)

// RazorCandidate is a peptide sequence to be evaluated as a razor
type RazorCandidate struct {
	Sequence       string
	MappedProtein  string
	MappedProteins map[string]MappedProtein
}
type MappedProtein struct {
	SID string
	W   float64
	GW  float64
	TNP int
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
