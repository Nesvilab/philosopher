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

	if len(a.CombPep) == 0 && len(a.CombPro) == 0 {
		logrus.Fatal("You need to specify a peptide or protein combined file for the Abacus analysis")
	}

	if len(a.CombPep) > 0 {
		e := peptideLevelAbacus(a, temp, args)
		if e != nil {
			return e
		}
	}

	if len(a.CombPro) > 0 {
		e := proteinLevelAbacus(a, temp, args)
		if e != nil {
			return e
		}
	}

	return nil
}
