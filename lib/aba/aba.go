// Package aba (Abacus)
package aba

import (
	"github.com/Nesvilab/philosopher/lib/met"
)

// DataSetLabelNames maps all custom names to each TMT tags
type DataSetLabelNames struct {
	Name      string
	LabelName map[string]string
}

// Run abacus
func Run(m met.Data, args []string) {

	psmLevelAbacus(m, args)

	if m.Abacus.Peptide {
		peptideLevelAbacus(m, args)
	}

	if m.Abacus.Protein {
		proteinLevelAbacus(m, args)
	}
}
