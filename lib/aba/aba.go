package aba

import (
	"github.com/prvst/philosopher/lib/met"
	"github.com/sirupsen/logrus"
)

// DataSetLabelNames maps all custom names to each TMT tags
type DataSetLabelNames struct {
	Name      string
	LabelName map[string]string
}

// Run abacus
// TODO update error methos on the abacus function
func Run(m met.Data, args []string) error {

	if m.Abacus.Peptide == false && m.Abacus.Protein == false {
		logrus.Fatal("You need to specify a peptide or protein combined file for the Abacus analysis")
	}

	if m.Abacus.Peptide == true {
		e := peptideLevelAbacus(m, args)
		if e != nil {
			return e
		}
	}

	if m.Abacus.Protein == true {
		e := proteinLevelAbacus(m, args)
		if e != nil {
			return e
		}
	}

	return nil
}
