package aba

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/prvst/philosopher/lib/id"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/rep"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/sirupsen/logrus"
)

// Create peptide combined report
func peptideLevelAbacus(a met.Abacus, temp string, args []string) error {

	var names []string
	var xmlFiles []string
	var datasets = make(map[string]rep.Evidence)

	var labelList []DataSetLabelNames

	// TODO: create a combined pepXML file by running Interprophet on all data sets

	// restoring combined file
	logrus.Info("Processing combined file")
	seqMap, chargeMap, err := processPeptideCombinedFile(a)
	if err != nil {
		return err
	}

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

	logrus.Info("Collecting data from individual experiments")
	evidences, err := collectPeptideDatafromExperiments(datasets, seqMap, chargeMap)

	logrus.Info("Processing spectral counts")
	evidences = getPeptideSpectralCounts(evidences, datasets)

	savePeptideAbacusResult(temp, evidences, datasets, names, a.Unique, false, labelList)

	return nil
}

func processPeptideCombinedFile(a met.Abacus) (map[string]int8, map[string][]string, error) {

	//var list rep.CombinedPeptideEvidenceList
	var seqMap = make(map[string]int8)
	var chargeMap = make(map[string][]string)

	if _, err := os.Stat(a.CombPep); os.IsNotExist(err) {
		logrus.Fatal("Cannot find combined.pep.xml file")
	} else {

		var pep id.PepXML
		pep.Read(a.CombPep)
		pep.DecoyTag = a.Tag

		// get all peptide sequences from combined file and collapse them
		for _, i := range pep.PeptideIdentification {
			if !strings.Contains(i.Protein, a.Tag) {
				seqMap[i.Peptide] = 0
				chargeMap[i.Peptide] = append(chargeMap[i.Peptide], strconv.Itoa(int(i.AssumedCharge)))
			}
		}

	}

	return seqMap, chargeMap, nil
}

func collectPeptideDatafromExperiments(datasets map[string]rep.Evidence, seqMap map[string]int8, chargeMap map[string][]string) (rep.CombinedPeptideEvidenceList, error) {

	var evidences rep.CombinedPeptideEvidenceList
	var uniqPeptides = make(map[string]uint8)
	var proteinMap = make(map[string]string)
	var proteinDesc = make(map[string]string)

	for _, v := range datasets {
		for _, i := range v.PSM {

			_, ok := seqMap[i.Peptide]
			if ok {

				var keys []string
				keys = append(keys, i.Peptide)

				var uniqMds = make(map[string]uint8)

				for _, j := range i.AssignedMassDiffs {
					mass := strconv.FormatFloat(j, 'f', 6, 64)
					uniqMds[mass] = 0
				}

				// this forces the unmodified pepes to collapse with peps containing +16
				if len(uniqMds) == 0 || uniqMds["0"] == 0 {
					uniqMds["15.994900"] = 0
					delete(uniqMds, "0")
					delete(uniqMds, "0.000000")
				}

				for j, _ := range uniqMds {
					keys = append(keys, j)
				}

				sort.Strings(keys[1:])

				key := strings.Join(keys, "#")
				uniqPeptides[key] = 0

				proteinMap[i.Peptide] = i.Protein
				proteinDesc[i.Peptide] = i.ProteinDescription

			}
		}
	}

	for k, _ := range uniqPeptides {

		var e rep.CombinedPeptideEvidence
		e.Spc = make(map[string]int)
		e.Intensity = make(map[string]float64)

		parts := strings.Split(k, "#")

		e.Key = k
		e.Sequence = parts[0]

		sort.Strings(parts[1:])
		e.AssignedMassDiffs = parts[1:]

		charges, ok := chargeMap[parts[0]]
		if ok {

			var uniqCharges = make(map[string]uint8)
			for _, ch := range charges {
				uniqCharges[ch] = 0
			}

			for ch, _ := range uniqCharges {
				e.ChargeStates = append(e.ChargeStates, ch)
			}
			sort.Strings(e.ChargeStates)
		}

		v, ok := proteinMap[parts[0]]
		if ok {
			e.Protein = v
			e.ProteinDescription = proteinDesc[parts[0]]
		}

		evidences = append(evidences, e)

	}

	return evidences, nil
}

func getPeptideSpectralCounts(combined rep.CombinedPeptideEvidenceList, datasets map[string]rep.Evidence) rep.CombinedPeptideEvidenceList {

	for k, v := range datasets {

		var keyMaps = make(map[string]int)

		for _, j := range v.PSM {

			var keys []string
			keys = append(keys, j.Peptide)

			var uniqMds = make(map[string]uint8)

			for _, k := range j.AssignedMassDiffs {
				mass := strconv.FormatFloat(k, 'f', 6, 64)
				uniqMds[mass] = 0
			}

			// this forces the unmodified pepes to collapse with peps containing +16
			if len(uniqMds) == 0 || uniqMds["0"] == 0 {
				uniqMds["15.994900"] = 0
				delete(uniqMds, "0")
				delete(uniqMds, "0.000000")
			}

			for k, _ := range uniqMds {
				keys = append(keys, k)
			}

			sort.Strings(keys[1:])

			key := strings.Join(keys, "#")

			keyMaps[key]++
		}

		for i, _ := range combined {
			count, ok := keyMaps[combined[i].Key]
			if ok {
				combined[i].Spc[k] = count
			}
		}

		// for i := range combined {
		// 	for _, j := range v.Peptides {
		// 		if combined[i].Sequence == j.Sequence {
		// 			combined[i].Spc[k] = j.Spc
		// 			break
		// 		}
		// 	}
		// }

	}

	return combined
}

// savePeptideAbacusResult creates a single report using 1 or more philosopher result files
func savePeptideAbacusResult(session string, evidences rep.CombinedPeptideEvidenceList, datasets map[string]rep.Evidence, namesList []string, uniqueOnly, hasTMT bool, labelsList []DataSetLabelNames) {

	// create result file
	output := fmt.Sprintf("%s%scombined_peptide.csv", session, string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create report file:", err)
	}
	defer file.Close()

	line := "Sequence\tCharge States\tAssigned Massdiffs\tProtein\tProtein Description\t"

	for _, i := range namesList {
		line += fmt.Sprintf("%s Spectral Count\t", i)
	}

	line += "\n"
	n, err := io.WriteString(file, line)
	if err != nil {
		logrus.Fatal(n, err)
	}

	// organize by group number
	sort.Sort(evidences)

	for _, i := range evidences {

		var line string

		line += fmt.Sprintf("%s\t", i.Sequence)

		line += fmt.Sprintf("%v\t", strings.Join(i.ChargeStates, ","))

		line += fmt.Sprintf("%v\t", strings.Join(i.AssignedMassDiffs, ","))

		line += fmt.Sprintf("%s\t", i.Protein)

		line += fmt.Sprintf("%s\t", i.ProteinDescription)

		for _, j := range namesList {
			line += fmt.Sprintf("%d\t", i.Spc[j])
		}

		line += "\n"
		n, err := io.WriteString(file, line)
		if err != nil {
			logrus.Fatal(n, err)
		}

	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}

// Create peptide combined report
// func peptideLevelAbacus(a met.Abacus, temp string, args []string) error {
//
// 	var names []string
// 	var xmlFiles []string
// 	var database dat.Base
// 	var datasets = make(map[string]rep.Evidence)
//
// 	var labelList []DataSetLabelNames
//
// 	// restore database
// 	database = dat.Base{}
// 	database.RestoreWithPath(args[0])
//
// 	// recover all files
// 	logrus.Info("Restoring results")
//
// 	for _, i := range args {
//
// 		// restoring the database
// 		var e rep.Evidence
// 		e.RestoreGranularWithPath(i)
//
// 		var labels DataSetLabelNames
// 		labels.LabelName = make(map[string]string)
//
// 		// collect interact full file names
// 		files, _ := ioutil.ReadDir(i)
// 		for _, f := range files {
// 			if strings.Contains(f.Name(), "pep.xml") {
// 				interactFile := fmt.Sprintf("%s%s%s", i, string(filepath.Separator), f.Name())
// 				absPath, _ := filepath.Abs(interactFile)
// 				xmlFiles = append(xmlFiles, absPath)
// 			}
// 		}
//
// 		var annot = fmt.Sprintf("%s%sannotation.txt", i, string(filepath.Separator))
// 		if strings.Contains(i, string(filepath.Separator)) {
// 			i = strings.Replace(i, string(filepath.Separator), "", -1)
// 			labels.Name = i
// 		} else {
// 			labels.Name = i
// 		}
// 		labels.LabelName, _ = getLabelNames(annot)
//
// 		// collect project names
// 		prjName := i
// 		if strings.Contains(prjName, string(filepath.Separator)) {
// 			prjName = strings.Replace(filepath.Base(prjName), string(filepath.Separator), "", -1)
// 		}
//
// 		labelList = append(labelList, labels)
//
// 		// unique list and map of datasets
// 		datasets[prjName] = e
// 		names = append(names, prjName)
// 	}
//
// 	sort.Strings(names)
//
// 	// logrus.Info("Processing spectral counts")
// 	// evidences = getProteinSpectralCounts(evidences, datasets)
// 	//
// 	// logrus.Info("Processing intensities")
// 	// evidences = sumProteinIntensities(evidences, datasets)
// 	//
// 	// // collect TMT labels
// 	// if a.Labels == true {
// 	// 	evidences = getProteinLabelIntensities(evidences, datasets)
// 	// }
// 	//
// 	// if a.Labels == true {
// 	// 	saveProteinAbacusResult(temp, evidences, datasets, names, a.Unique, true, labelList)
// 	// } else {
// 	// 	saveProteinAbacusResult(temp, evidences, datasets, names, a.Unique, false, labelList)
// 	// }
//
// 	return nil
// }
//
