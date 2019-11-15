// Package aba (Abacus)
package aba

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"philosopher/lib/met"
	"philosopher/lib/msg"
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

// addCustomNames adds to the label structures user-defined names to be used on the TMT labels
func getLabelNames(annot string) map[string]string {

	var labels = make(map[string]string)

	file, e := os.Open(annot)
	if e != nil {
		msg.ReadFile(errors.New("Cannot open annotation file"), "error")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		names := strings.Split(scanner.Text(), " ")
		labels[names[0]] = names[1]
	}

	if e = scanner.Err(); e != nil {
		msg.Custom(errors.New("The annotation file looks to be empty"), "fatal")
	}

	return labels
}
