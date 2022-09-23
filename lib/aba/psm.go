package aba

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"philosopher/lib/met"
	"philosopher/lib/msg"
	"philosopher/lib/rep"
	"philosopher/lib/sys"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

func psmLevelAbacus(m met.Data, args []string) {

	var names []string

	var labelList []DataSetLabelNames

	// recover all files
	logrus.Info("Restoring PSM results")

	var evidences rep.CombinedPSMEvidenceList

	for _, i := range args {

		// restoring the database
		var e rep.Evidence
		e.RestoreGranularWithPath(i)

		var labels DataSetLabelNames
		labels.LabelName = make(map[string]string)

		// collect interact full file names
		files, _ := ioutil.ReadDir(i)
		for _, f := range files {
			if strings.Contains(f.Name(), "annotation") {
				var annot = fmt.Sprintf("%s%s%s", i, string(filepath.Separator), f.Name())
				labels.Name = annot

				if len(m.Quantify.Annot) > 0 {
					labels.LabelName = getLabelNames(annot)
				}
			}
		}

		// collect project names
		prjName := i
		if strings.Contains(prjName, string(filepath.Separator)) {
			prjName = strings.Replace(filepath.Base(prjName), string(filepath.Separator), "", -1)
		}

		labelList = append(labelList, labels)

		// unique list and map of datasets
		names = append(names, prjName)
		sort.Strings(names)

		for _, j := range e.PSM {

			var psm rep.CombinedPSMEvidence

			psm.DataSet = prjName
			psm.Source = j.Source
			psm.Spectrum = j.Spectrum
			psm.SpectrumFile = j.SpectrumFile
			psm.Peptide = j.Peptide
			psm.ModifiedPeptide = j.ModifiedPeptide
			psm.Protein = j.Protein
			psm.ProteinDescription = j.ProteinDescription
			psm.ProteinID = j.ProteinID
			psm.EntryName = j.EntryName
			psm.GeneName = j.GeneName
			psm.AssumedCharge = j.AssumedCharge

			evidences = append(evidences, psm)
		}
	}

	if m.Abacus.Labels {
		savePSMAbacusResult(m.Temp, evidences, names, m.Abacus.Unique, true, m.Abacus.Full, labelList)
	} else {
		savePSMAbacusResult(m.Temp, evidences, names, m.Abacus.Unique, false, m.Abacus.Full, labelList)
	}

}

// savePSMAbacusResult creates a single report using 1 or more philosopher result files
func savePSMAbacusResult(session string, evidences rep.CombinedPSMEvidenceList, namesList []string, uniqueOnly, hasTMT, full bool, labelsList []DataSetLabelNames) {

	// create result file
	output := fmt.Sprintf("%s%scombined_psm.tsv", session, string(filepath.Separator))

	// create result file
	file, e := os.Create(output)
	if e != nil {
		msg.WriteFile(e, "error")
	}
	defer file.Close()

	header := "Source\tSpectrum\tSpectrumFile\tPeptide\tModified Peptide\tCharge\tProtein\tProtein ID\tEntry Name\tProtein Description\tGene"

	// Add Unique+Razor Intensity
	for _, i := range namesList {
		header += fmt.Sprintf("\t%s Intensity", i)
	}

	if hasTMT {
		for _, i := range namesList {
			header += fmt.Sprintf("\t%s 126 Abundance", i)
			header += fmt.Sprintf("\t%s 127N Abundance", i)
			header += fmt.Sprintf("\t%s 127C Abundance", i)
			header += fmt.Sprintf("\t%s 128N Abundance", i)
			header += fmt.Sprintf("\t%s 128C Abundance", i)
			header += fmt.Sprintf("\t%s 129N Abundance", i)
			header += fmt.Sprintf("\t%s 129C Abundance", i)
			header += fmt.Sprintf("\t%s 130N Abundance", i)
			header += fmt.Sprintf("\t%s 130C Abundance", i)
			header += fmt.Sprintf("\t%s 131N Abundance", i)

			for _, j := range labelsList {
				if j.Name == i {
					for k, v := range j.LabelName {
						before := fmt.Sprintf("%s %s Abundance", i, k)
						after := fmt.Sprintf("%s Abundance", v)
						header = strings.Replace(header, before, after, -1)
					}
				}
			}
		}
	}

	header += "\n"
	_, e = io.WriteString(file, header)
	if e != nil {
		msg.WriteToFile(e, "fatal")
	}

	for _, i := range evidences {
		var line string

		line += fmt.Sprintf("%s\t", i.Source)

		line += fmt.Sprintf("%s\t", i.Spectrum)

		line += fmt.Sprintf("%s\t", i.SpectrumFile)

		line += fmt.Sprintf("%s\t", i.Peptide)

		line += fmt.Sprintf("%s\t", i.ModifiedPeptide)

		line += fmt.Sprintf("%d\t", i.AssumedCharge)

		line += fmt.Sprintf("%s\t", i.Protein)

		line += fmt.Sprintf("%s\t", i.ProteinID)

		line += fmt.Sprintf("%s\t", i.EntryName)

		line += fmt.Sprintf("%s\t", i.ProteinDescription)

		line += fmt.Sprintf("%s\t", i.GeneName)

		line += "\n"
		_, e := io.WriteString(file, line)
		if e != nil {
			msg.WriteToFile(e, "fatal")
		}
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))
}
