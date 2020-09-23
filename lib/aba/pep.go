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

	"philosopher/lib/fil"
	"philosopher/lib/id"
	"philosopher/lib/met"
	"philosopher/lib/msg"
	"philosopher/lib/rep"
	"philosopher/lib/sys"

	"github.com/sirupsen/logrus"
)

// Create peptide combined report
func peptideLevelAbacus(m met.Data, args []string) {

	var names []string
	var xmlFiles []string
	var datasets = make(map[string]rep.Evidence)
	var labelList []DataSetLabelNames

	// restoring combined file
	logrus.Info("Processing combined file")
	processPeptideCombinedFile(m.Abacus)

	// recover all files
	logrus.Info("Restoring peptide results")

	local, _ := os.Getwd()
	local, _ = filepath.Abs(local)

	for _, i := range args {

		os.Chdir(i)

		// restoring the PSMs
		var psm rep.Evidence
		rep.RestoreEVPSM(&psm)

		var labels DataSetLabelNames
		labels.LabelName = make(map[string]string)

		files, _ := ioutil.ReadDir(i)

		// collect interact full file names
		for _, f := range files {
			if strings.Contains(f.Name(), "pep.xml") {
				interactFile := fmt.Sprintf("%s%s%s", i, string(filepath.Separator), f.Name())
				absPath, _ := filepath.Abs(interactFile)
				xmlFiles = append(xmlFiles, absPath)
			}
			if strings.Contains(f.Name(), "annotation") {
				var annot = fmt.Sprintf("%s%s%s", i, string(filepath.Separator), f.Name())
				labels.Name = annot

				if len(m.Quantify.Annot) > 0 {
					labels.LabelName = getLabelNames(annot)
				}
			}

			os.Chdir(local)
		}

		// collect project names
		prjName := i
		if strings.Contains(prjName, string(filepath.Separator)) {
			prjName = strings.Replace(filepath.Base(prjName), string(filepath.Separator), "", -1)
		}

		labelList = append(labelList, labels)

		// unique list and map of datasets
		datasets[prjName] = psm
		names = append(names, prjName)
	}

	os.Chdir(local)

	sort.Strings(names)

	logrus.Info("Collecting data from individual experiments")
	evidences := collectPeptideDatafromExperiments(datasets, m.Abacus.Tag)

	logrus.Info("Summarizing the quantification")
	evidences = SummarizeAttributes(evidences, datasets, local)

	os.Chdir(local)

	savePeptideAbacusResult(m.Temp, evidences, datasets, names, m.Abacus.Unique, false, labelList)

	return
}

// processPeptideCombinedFile reads and filter the combined peptide report
func processPeptideCombinedFile(a met.Abacus) {

	var pepID id.PepIDList
	var filteredPeptides id.PepIDList

	if _, e := os.Stat("combined.pep.xml"); os.IsNotExist(e) {

		msg.NoParametersFound(errors.New("Cannot find the combined.pep.xml file"), "fatal")

	} else {

		var pep id.PepXML
		pep.DecoyTag = a.Tag

		pepID, _ = id.ReadPepXMLInput("combined.pep.xml", a.Tag, sys.GetTemp(), false)

		//uniqPsms := fil.GetUniquePSMs(pepID)
		uniqPeps := fil.GetUniquePeptides(pepID)

		//filteredPSMs, _ := fil.PepXMLFDRFilter(uniqPsms, 0.01, "PSM", a.Tag)
		filteredPeptides, _ = fil.PepXMLFDRFilter(uniqPeps, 0.01, "Peptide", a.Tag)
		filteredPeptides.Serialize("pep")

	}

	return
}

// collectPeptideDatafromExperiments reads each individual data set peptide output and collects the quantification data to the combined report
func collectPeptideDatafromExperiments(datasets map[string]rep.Evidence, decoyTag string) rep.CombinedPeptideEvidenceList {

	var pep id.PepIDList
	pep.Restore("pep")

	var evidences rep.CombinedPeptideEvidenceList

	for _, i := range pep {
		if !strings.HasPrefix(i.Protein, decoyTag) {
			var e rep.CombinedPeptideEvidence
			e.Spc = make(map[string]int)
			e.Intensity = make(map[string]float64)
			e.AssignedMassDiffs = make(map[string]uint8)
			e.ChargeStates = make(map[uint8]uint8)

			e.Sequence = i.Peptide
			e.Protein = i.Protein

			evidences = append(evidences, e)
		}
	}

	return evidences
}

// SummarizeAttributes collects spectral counts and intensities from the individual data sets for the combined peptide report
func SummarizeAttributes(evidences rep.CombinedPeptideEvidenceList, datasets map[string]rep.Evidence, local string) rep.CombinedPeptideEvidenceList {

	var chargeMap = make(map[string][]uint8)
	var bestPSM = make(map[string]float64)

	for k := range datasets {

		os.Chdir(k)

		var evi rep.Evidence
		evi.RestoreGranular()

		SpcMap := make(map[string]int)
		IntMap := make(map[string]float64)
		ModsMap := make(map[string][]string)

		protIDMap := make(map[string]string)
		protDescMap := make(map[string]string)
		GeneMap := make(map[string]string)

		for _, j := range evi.Peptides {

			for _, k := range j.Modifications.Index {
				if k.Type == "Assigned" {
					mass := strconv.FormatFloat(k.MassDiff, 'f', 6, 64)
					ModsMap[j.Sequence] = append(ModsMap[j.Sequence], mass)
				}
			}

			SpcMap[j.Sequence] = j.Spc
			IntMap[j.Sequence] = j.Intensity

			protIDMap[j.Sequence] = j.ProteinID
			protDescMap[j.Sequence] = j.ProteinDescription
			GeneMap[j.Sequence] = j.GeneName

			// get all charge states
			for l := range j.ChargeState {
				chargeMap[j.Sequence] = append(chargeMap[j.Sequence], l)
			}

			if j.Probability > bestPSM[j.Sequence] {
				bestPSM[j.Sequence] = j.Probability
			}

		}

		for i := range evidences {
			spc, ok := SpcMap[evidences[i].Sequence]
			if ok {
				evidences[i].Spc[k] = spc
			}
			it, ok := IntMap[evidences[i].Sequence]
			if ok {
				evidences[i].Intensity[k] = it
			}
			m, ok := ModsMap[evidences[i].Sequence]
			if ok {
				for _, l := range m {
					evidences[i].AssignedMassDiffs[l] = 0
				}
			}
			c, ok := chargeMap[evidences[i].Sequence]
			if ok {
				for _, l := range c {
					evidences[i].ChargeStates[l] = 0
				}
			}
			id, ok := protIDMap[evidences[i].Sequence]
			if ok {
				evidences[i].ProteinID = id
			}
			desc, ok := protDescMap[evidences[i].Sequence]
			if ok {
				evidences[i].ProteinDescription = desc
			}
			gene, ok := GeneMap[evidences[i].Sequence]
			if ok {
				evidences[i].Gene = gene
			}
			prob, ok := bestPSM[evidences[i].Sequence]
			if ok {
				evidences[i].BestPSM = prob
			}
		}

		os.Chdir(local)
	}

	os.Chdir(local)

	return evidences
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

		var c []string
		for j := range i.ChargeStates {
			c = append(c, strconv.Itoa(int(j)))
		}
		line += fmt.Sprintf("%s\t", strings.Join(c, ","))

		line += fmt.Sprintf("%f\t", i.BestPSM)

		var m []string
		for j := range i.AssignedMassDiffs {
			m = append(m, j)
		}
		sort.Strings(m)
		line += fmt.Sprintf("%v\t", strings.Join(m, ","))

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
