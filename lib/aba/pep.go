// Package aba (Abacus), peptide level
package aba

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"philosopher/lib/fil"
	"philosopher/lib/id"
	"philosopher/lib/met"
	"philosopher/lib/msg"
	"philosopher/lib/rep"
	"philosopher/lib/sys"
)

// Create peptide combined report
func peptideLevelAbacus(m met.Data, args []string) {

	var names []string
	var xmlFiles []string
	var datasets = make(map[string]rep.Evidence)
	var labelList []DataSetLabelNames

	// restoring combined file
	logrus.Info("Processing combined file")
	seqMap, chargeMap := processPeptideCombinedFile(m.Abacus)

	// recover all files
	logrus.Info("Restoring peptide results")

	for _, i := range args {

		// restoring the database
		var e rep.Evidence
		//e.RestoreGranularWithPath(i)
		rep.RestoreEVPSMWithPath(&e, i)

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

		if len(m.Quantify.Annot) > 0 {
			labels.LabelName = getLabelNames(annot)
		}

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
	evidences := collectPeptideDatafromExperiments(datasets, seqMap, chargeMap)

	logrus.Info("Processing spectral counts")
	evidences = getPeptideSpectralCounts(evidences, datasets)

	logrus.Info("Processing intensities")
	evidences = getIntensities(evidences, datasets)

	savePeptideAbacusResult(m.Temp, evidences, datasets, names, m.Abacus.Unique, false, labelList)

	return
}

// processPeptideCombinedFile reads and filter the combined peptide report
func processPeptideCombinedFile(a met.Abacus) (map[string]int8, map[string][]string) {

	//var list rep.CombinedPeptideEvidenceList
	var seqMap = make(map[string]int8)
	var chargeMap = make(map[string][]string)

	if _, e := os.Stat("combined.pep.xml"); os.IsNotExist(e) {

		msg.NoParametersFound(errors.New("Cannot find the combined.pep.xml file"), "fatal")

	} else {

		var pep id.PepXML
		var pepID id.PepIDList
		pep.Read("combined.pep.xml")
		pep.DecoyTag = a.Tag

		for _, i := range pep.PeptideIdentification {
			pepID = append(pepID, i)
		}

		uniqPsms := fil.GetUniquePSMs(pepID)
		uniqPeps := fil.GetUniquePeptides(pepID)

		filteredPSMs, _ := fil.PepXMLFDRFilter(uniqPsms, 0.01, "PSM", a.Tag)
		filteredPeptides, _ := fil.PepXMLFDRFilter(uniqPeps, 0.01, "Peptide", a.Tag)

		// get all peptide sequences from combined file and collapse them
		for _, i := range filteredPeptides {
			if !strings.HasPrefix(i.Protein, a.Tag) {
				seqMap[i.Peptide] = 0
			}
		}

		// get all charge states
		for _, i := range filteredPSMs {
			if !strings.HasPrefix(i.Protein, a.Tag) {
				chargeMap[i.Peptide] = append(chargeMap[i.Peptide], strconv.Itoa(int(i.AssumedCharge)))
			}
		}

	}

	return seqMap, chargeMap
}

// collectPeptideDatafromExperiments reads each individual data set peptide output and collects the quantification data to the combined report
func collectPeptideDatafromExperiments(datasets map[string]rep.Evidence, seqMap map[string]int8, chargeMap map[string][]string) rep.CombinedPeptideEvidenceList {

	var evidences rep.CombinedPeptideEvidenceList
	var uniqPeptides = make(map[string]uint8)
	var proteinMap = make(map[string]string)
	var proteinIDMap = make(map[string]string)
	var proteinDesc = make(map[string]string)
	var geneMap = make(map[string]string)
	var probMap = make(map[string]float64)

	for _, v := range datasets {
		for _, i := range v.PSM {

			_, ok := seqMap[i.Peptide]
			if ok {

				var keys []string
				keys = append(keys, i.Peptide)

				if i.Probability > probMap[i.Peptide] {
					probMap[i.Peptide] = i.Probability
				}

				var uniqMds = make(map[string]uint8)

				for _, j := range i.Modifications.Index {
					if j.Type == "Assigned" {
						mass := strconv.FormatFloat(j.MassDiff, 'f', 6, 64)
						uniqMds[mass] = 0
					}
				}

				// this forces the unmodified peps to collapse with peps containing +16
				if len(uniqMds) == 0 || uniqMds["0"] == 0 {
					delete(uniqMds, "0")
					delete(uniqMds, "0.000000")
				}
				delete(uniqMds, "15.994900")

				// if len(uniqMds) == 0 || uniqMds["0"] == 0 {
				// 	uniqMds["15.994900"] = 0
				// 	delete(uniqMds, "0")
				// 	delete(uniqMds, "0.000000")
				// }

				for j := range uniqMds {
					keys = append(keys, j)
				}

				sort.Strings(keys[1:])

				key := strings.Join(keys, "#")
				uniqPeptides[key] = 0

				proteinMap[i.Peptide] = i.Protein
				proteinIDMap[i.Peptide] = i.ProteinID
				proteinDesc[i.Peptide] = i.ProteinDescription
				geneMap[i.Peptide] = i.GeneName

			}
		}
	}

	for k := range uniqPeptides {

		var e rep.CombinedPeptideEvidence
		e.Spc = make(map[string]int)
		e.Intensity = make(map[string]float64)

		parts := strings.Split(k, "#")

		e.Key = k
		e.Sequence = parts[0]
		e.BestPSM = probMap[parts[0]]

		sort.Strings(parts[1:])
		e.AssignedMassDiffs = parts[1:]

		charges, ok := chargeMap[parts[0]]
		if ok {

			var uniqCharges = make(map[string]uint8)
			for _, ch := range charges {
				uniqCharges[ch] = 0
			}

			for ch := range uniqCharges {
				e.ChargeStates = append(e.ChargeStates, ch)
			}
			sort.Strings(e.ChargeStates)
		}

		v, ok := proteinMap[parts[0]]
		if ok {
			e.Protein = v
			e.ProteinID = proteinIDMap[parts[0]]
			e.ProteinDescription = proteinDesc[parts[0]]
			e.Gene = geneMap[parts[0]]
		}

		evidences = append(evidences, e)

	}

	return evidences
}

// getPeptideSpectralCounts collects spectral counts from the individual data sets for the combined peptide report
func getPeptideSpectralCounts(combined rep.CombinedPeptideEvidenceList, datasets map[string]rep.Evidence) rep.CombinedPeptideEvidenceList {

	for k, v := range datasets {

		var keyMaps = make(map[string]int)

		for _, j := range v.PSM {

			var keys []string
			keys = append(keys, j.Peptide)

			var uniqMds = make(map[string]uint8)

			for _, k := range j.Modifications.Index {
				if k.Type == "Assigned" {
					mass := strconv.FormatFloat(k.MassDiff, 'f', 6, 64)
					uniqMds[mass] = 0
				}
			}

			// this forces the unmodified pepes to collapse with peps containing +16
			// if len(uniqMds) == 0 || uniqMds["0"] == 0 {
			// 	uniqMds["15.994900"] = 0
			// 	delete(uniqMds, "0")
			// 	delete(uniqMds, "0.000000")
			// }
			if len(uniqMds) == 0 || uniqMds["0"] == 0 {
				delete(uniqMds, "0")
				delete(uniqMds, "0.000000")
			}
			delete(uniqMds, "15.994900")

			for k := range uniqMds {
				keys = append(keys, k)
			}

			sort.Strings(keys[1:])

			key := strings.Join(keys, "#")

			keyMaps[key]++
		}

		for i := range combined {
			count, ok := keyMaps[combined[i].Key]
			if ok {
				combined[i].Spc[k] = count
			}
		}

	}

	return combined
}

// getIntensities collects intensities from the individual data sets for the combined peptide report
func getIntensities(combined rep.CombinedPeptideEvidenceList, datasets map[string]rep.Evidence) rep.CombinedPeptideEvidenceList {

	for k, v := range datasets {

		var keyMaps = make(map[string]float64)

		for _, j := range v.PSM {

			var keys []string
			keys = append(keys, j.Peptide)

			var uniqMds = make(map[string]uint8)

			for _, k := range j.Modifications.Index {
				if k.Type == "Assigned" {
					mass := strconv.FormatFloat(k.MassDiff, 'f', 6, 64)
					uniqMds[mass] = 0
				}
			}

			// this forces the unmodified pepes to collapse with peps containing +16
			// if len(uniqMds) == 0 || uniqMds["0"] == 0 {
			// 	uniqMds["15.994900"] = 0
			// 	delete(uniqMds, "0")
			// 	delete(uniqMds, "0.000000")
			// }
			if len(uniqMds) == 0 || uniqMds["0"] == 0 {
				delete(uniqMds, "0")
				delete(uniqMds, "0.000000")
			}
			delete(uniqMds, "15.994900")

			for k := range uniqMds {
				keys = append(keys, k)
			}

			sort.Strings(keys[1:])

			key := strings.Join(keys, "#")

			v, ok := keyMaps[key]
			if !ok {
				keyMaps[key] = j.Intensity
			} else {
				if j.Intensity > v {
					keyMaps[key] = j.Intensity
				}
			}

		}

		for i := range combined {
			int, ok := keyMaps[combined[i].Key]
			if ok {
				combined[i].Intensity[k] = int
			}
		}

	}

	return combined
}

// savePeptideAbacusResult creates a single report using 1 or more philosopher result files
func savePeptideAbacusResult(session string, evidences rep.CombinedPeptideEvidenceList, datasets map[string]rep.Evidence, namesList []string, uniqueOnly, hasTMT bool, labelsList []DataSetLabelNames) {

	// create result file
	output := fmt.Sprintf("%s%scombined_peptide.tsv", session, string(filepath.Separator))

	// create result file
	file, e := os.Create(output)
	if e != nil {
		msg.WriteFile(e, "error")
	}
	defer file.Close()

	line := "Sequence\tCharge States\tProbability\tAssigned Modifications\tGene\tProtein\tProtein ID\tProtein Description\t"

	for _, i := range namesList {
		line += fmt.Sprintf("%s Spectral Count\t", i)
		line += fmt.Sprintf("%s Intensity\t", i)
	}

	line += "\n"
	_, e = io.WriteString(file, line)
	if e != nil {
		msg.WriteToFile(e, "fatal")
	}

	// organize by group number
	sort.Sort(evidences)

	for _, i := range evidences {

		var line string

		line += fmt.Sprintf("%s\t", i.Sequence)

		line += fmt.Sprintf("%v\t", strings.Join(i.ChargeStates, ","))

		line += fmt.Sprintf("%f\t", i.BestPSM)

		line += fmt.Sprintf("%v\t", strings.Join(i.AssignedMassDiffs, ","))

		line += fmt.Sprintf("%s\t", i.Gene)

		line += fmt.Sprintf("%s\t", i.Protein)

		line += fmt.Sprintf("%s\t", i.ProteinID)

		line += fmt.Sprintf("%s\t", i.ProteinDescription)

		for _, j := range namesList {
			line += fmt.Sprintf("%d\t%.4f\t", i.Spc[j], i.Intensity[j])
		}

		line += "\n"
		_, e = io.WriteString(file, line)
		if e != nil {
			msg.WriteToFile(e, "fatal")
		}

	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}
