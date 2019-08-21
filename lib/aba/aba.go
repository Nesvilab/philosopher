// Package aba (Abacus)
package aba

import (
	"errors"

	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/msg"
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
		msg.Custom(errors.New("You need to specify a peptide or protein combined file for the Abacus analysis"), "fatal")
	}

	if m.Abacus.Peptide == true {
		peptideLevelAbacus(m, args)
	}

	if m.Abacus.Protein == true {
		proteinLevelAbacus(m, args)
	}

	return
}
