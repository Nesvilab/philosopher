package aba

import (
	"errors"

	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/met"
)

// DataSetLabelNames maps all custom names to each TMT tags
type DataSetLabelNames struct {
	Name      string
	LabelName map[string]string
}

// Run abacus
// TODO update error methos on the abacus function
func Run(m met.Data, args []string) {

	if m.Abacus.Peptide == false && m.Abacus.Protein == false {
		err.FatalCustom(errors.New("You need to specify a peptide or protein combined file for the Abacus analysis"))
	}

	if m.Abacus.Peptide == true {
		peptideLevelAbacus(m, args)
	}

	if m.Abacus.Protein == true {
		proteinLevelAbacus(m, args)
	}

	return
}
