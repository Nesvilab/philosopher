// Package aba (Abacus)
package aba

import (
	"philosopher/lib/met"
)

// DataSetLabelNames maps all custom names to each TMT tags
type DataSetLabelNames struct {
	Name      string
	LabelName map[string]string
}

// Run abacus
// TODO update error methos on the abacus function
func Run(m met.Data, args []string) {

	// if !m.Abacus.Peptide && !m.Abacus.Protein {
	// 	msg.Custom(errors.New("you need to specify a peptide or protein combined file for the Abacus analysis"), "fatal")
	// }

	psmLevelAbacus(m, args)

	if m.Abacus.Peptide {
		peptideLevelAbacus(m, args)
	}

	if m.Abacus.Protein {
		proteinLevelAbacus(m, args)
	}
}

// addCustomNames adds to the label structures user-defined names to be used on the TMT labels
// func getLabelNames(dataSet, annot string) map[string]string {

// 	var labels = make(map[string]string)

// 	file, e := os.Open(annot)
// 	if e != nil {
// 		msg.ReadFile(errors.New("cannot open annotation file"), "error")
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		names := strings.Fields(scanner.Text())

// 		name := dataSet + " " + names[0]

// 		labels[name] = names[1]
// 	}

// 	if e = scanner.Err(); e != nil {
// 		msg.Custom(errors.New("the annotation file looks to be empty"), "fatal")
// 	}

// 	return labels
// }
