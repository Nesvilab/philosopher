package aba

import (
	"github.com/prvst/philosopher/lib/met"
)

// DataSetLabelNames maps all custom names to each TMT tags
type DataSetLabelNames struct {
	Name      string
	LabelName map[string]string
}

// TODO update error methos on the abacus function
// Run abacus
func Run(a met.Abacus, temp string, args []string) error {

	// e := peptideLevelAbacus(a, temp, args)
	// if e != nil {
	// 	return e
	// }

	e := proteinLevelAbacus(a, temp, args)
	if e != nil {
		return e
	}

	return nil
}
