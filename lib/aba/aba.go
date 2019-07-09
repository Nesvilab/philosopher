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
func Run(a met.Abacus, temp string, args []string) error {

	if a.Peptide == false && a.Protein == false {
		logrus.Fatal("You need to specify a peptide or protein combined file for the Abacus analysis")
	}

	if a.Peptide == true {
		e := peptideLevelAbacus(a, temp, args)
		if e != nil {
			return e
		}
	}

	if a.Protein == true {
		e := proteinLevelAbacus(a, temp, args)
		if e != nil {
			return e
		}
	}

	return nil
}
