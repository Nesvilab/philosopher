package aba

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"philosopher/lib/iso"
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

	var labels = make(map[string]string)

	// recover all files
	logrus.Info("Restoring PSM results")

	var evidences rep.CombinedPSMEvidenceList

	for _, i := range args {

		// restoring the database
		var e rep.Evidence
		e.RestoreGranularWithPath(i)

		// collect interact full file names
		files, _ := ioutil.ReadDir(i)
		for _, f := range files {
			if strings.Contains(f.Name(), "annotation") {
				var annot = fmt.Sprintf("%s%s%s", i, string(filepath.Separator), f.Name())

				file, e := os.Open(annot)
				if e != nil {
					msg.ReadFile(errors.New("cannot open annotation file"), "error")
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					names := strings.Fields(scanner.Text())

					name := i + " " + names[0]

					labels[name] = names[1]
				}

				if e = scanner.Err(); e != nil {
					msg.Custom(errors.New("the annotation file looks to be empty"), "fatal")
				}
			}
		}

		// collect project names
		prjName := i
		if strings.Contains(prjName, string(filepath.Separator)) {
			prjName = strings.Replace(filepath.Base(prjName), string(filepath.Separator), "", -1)
		}

		// unique list and map of datasets
		names = append(names, prjName)
		sort.Strings(names)

		for _, j := range e.PSM {

			var psm rep.CombinedPSMEvidence
			psm.Intensity = make(map[string]float64)
			psm.Labels = make(map[string]iso.Labels)

			psm.DataSet = prjName
			psm.Source = j.Source
			psm.Spectrum = j.Spectrum
			psm.Peptide = j.Peptide
			psm.ModifiedPeptide = j.ModifiedPeptide
			psm.Protein = j.Protein
			psm.ProteinDescription = j.ProteinDescription
			psm.ProteinID = j.ProteinID
			psm.EntryName = j.EntryName
			psm.GeneName = j.GeneName
			psm.AssumedCharge = j.AssumedCharge
			psm.IsUnique = j.IsUnique
			psm.Purity = j.Purity

			psm.Intensity[prjName] = j.Intensity

			if j.Labels != nil {
				psm.Labels[prjName] = *j.Labels
				if j.Labels.IsUsed {
					psm.IsUsed = true
				}
			}

			evidences = append(evidences, psm)
		}
	}

	if m.Abacus.Labels {
		savePSMAbacusResult(m.Temp, evidences, names, m.Abacus.Unique, true, m.Abacus.Full, labels)
	} else {
		savePSMAbacusResult(m.Temp, evidences, names, m.Abacus.Unique, false, m.Abacus.Full, labels)
	}

}

// savePSMAbacusResult creates a single report using 1 or more philosopher result files
func savePSMAbacusResult(session string, evidences rep.CombinedPSMEvidenceList, namesList []string, uniqueOnly, hasLabels, full bool, labelsList map[string]string) {

	// create result file
	output := fmt.Sprintf("%s%scombined_psm.tsv", session, string(filepath.Separator))

	// create result file
	file, e := os.Create(output)
	if e != nil {
		msg.WriteFile(e, "error")
	}
	defer file.Close()

	header := "Spectrum\tSpectrum File\tPeptide\tModified Peptide\tCharge\tGene\tProtein\tProtein ID\tEntry Name\tProtein Description\tIs Unique\tQuan Usage\tPurity"

	// Add Unique+Razor Intensity
	for _, i := range namesList {
		header += fmt.Sprintf("\t%s", i)
	}

	chs := []string{"126", "127N", "127C", "128N", "128C", "129N", "129C", "130N", "130C", "131N", "131C", "132N", "132C", "133N", "133C", "134N", "134C", "135N", "135C"}

	if hasLabels {
		for _, i := range namesList {
			for _, j := range chs {
				l := fmt.Sprintf("%s %s", i, j)
				v, ok := labelsList[l]
				if ok {
					header += fmt.Sprintf("\t%s", v)
				} else {
					header += fmt.Sprintf("\t%s %s", i, j)
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

		line += fmt.Sprintf("%s\t", i.Spectrum)

		line += fmt.Sprintf("%s.raw\t", i.Source)

		line += fmt.Sprintf("%s\t", i.Peptide)

		line += fmt.Sprintf("%s\t", i.ModifiedPeptide)

		line += fmt.Sprintf("%d\t", i.AssumedCharge)

		line += fmt.Sprintf("%s\t", i.GeneName)

		line += fmt.Sprintf("%s\t", i.Protein)

		line += fmt.Sprintf("%s\t", i.ProteinID)

		line += fmt.Sprintf("%s\t", i.EntryName)

		line += fmt.Sprintf("%s\t", i.ProteinDescription)

		line += fmt.Sprintf("%t\t", i.IsUnique)

		line += fmt.Sprintf("%t\t", i.IsUsed)

		line += fmt.Sprintf("%.2f\t", i.Purity)

		for _, j := range namesList {
			line += fmt.Sprintf("%6.f\t", i.Intensity[j])
		}

		if hasLabels {
			for _, j := range namesList {
				line += fmt.Sprintf("%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
					i.Labels[j].Channel1.Intensity,
					i.Labels[j].Channel2.Intensity,
					i.Labels[j].Channel3.Intensity,
					i.Labels[j].Channel4.Intensity,
					i.Labels[j].Channel5.Intensity,
					i.Labels[j].Channel6.Intensity,
					i.Labels[j].Channel7.Intensity,
					i.Labels[j].Channel8.Intensity,
					i.Labels[j].Channel9.Intensity,
					i.Labels[j].Channel10.Intensity,
					i.Labels[j].Channel11.Intensity,
					i.Labels[j].Channel12.Intensity,
					i.Labels[j].Channel13.Intensity,
					i.Labels[j].Channel14.Intensity,
					i.Labels[j].Channel15.Intensity,
					i.Labels[j].Channel16.Intensity,
					i.Labels[j].Channel17.Intensity,
					i.Labels[j].Channel18.Intensity,
				)
			}
		}

		line += "\n"
		_, e := io.WriteString(file, line)
		if e != nil {
			msg.WriteToFile(e, "fatal")
		}
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))
}
