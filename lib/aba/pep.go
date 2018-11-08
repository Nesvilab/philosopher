package aba

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"github.com/prvst/philosopher/lib/dat"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/rep"
	"github.com/sirupsen/logrus"
)

// Create peptide combined report
func peptideLevelAbacus(a met.Abacus, temp string, args []string) error {

	var names []string
	var xmlFiles []string
	var database dat.Base
	var datasets = make(map[string]rep.Evidence)

	var labelList []DataSetLabelNames

	// restore database
	database = dat.Base{}
	database.RestoreWithPath(args[0])

	// recover all files
	logrus.Info("Restoring results")

	for _, i := range args {

		// restoring the database
		var e rep.Evidence
		e.RestoreGranularWithPath(i)

		var labels DataSetLabelNames
		labels.LabelName = make(map[string]string)

		// collect interact full file names
		files, _ := ioutil.ReadDir(i)
		for _, f := range files {
			if strings.Contains(f.Name(), "pep.xml") {
				interactFile := fmt.Sprintf("%s%s%s", i, string(filepath.Separator), f.Name())
				absPath, _ := filepath.Abs(interactFile)
				xmlFiles = append(xmlFiles, absPath)
			}
		}

		var annot = fmt.Sprintf("%s%sannotation.txt", i, string(filepath.Separator))
		if strings.Contains(i, string(filepath.Separator)) {
			i = strings.Replace(i, string(filepath.Separator), "", -1)
			labels.Name = i
		} else {
			labels.Name = i
		}
		labels.LabelName, _ = getLabelNames(annot)

		// collect project names
		prjName := i
		if strings.Contains(prjName, string(filepath.Separator)) {
			prjName = strings.Replace(filepath.Base(prjName), string(filepath.Separator), "", -1)
		}

		labelList = append(labelList, labels)

		// unique list and map of datasets
		datasets[prjName] = e
		names = append(names, prjName)
	}

	sort.Strings(names)

	// logrus.Info("Processing spectral counts")
	// evidences = getProteinSpectralCounts(evidences, datasets)
	//
	// logrus.Info("Processing intensities")
	// evidences = sumProteinIntensities(evidences, datasets)
	//
	// // collect TMT labels
	// if a.Labels == true {
	// 	evidences = getProteinLabelIntensities(evidences, datasets)
	// }
	//
	// if a.Labels == true {
	// 	saveProteinAbacusResult(temp, evidences, datasets, names, a.Unique, true, labelList)
	// } else {
	// 	saveProteinAbacusResult(temp, evidences, datasets, names, a.Unique, false, labelList)
	// }

	return nil
}
