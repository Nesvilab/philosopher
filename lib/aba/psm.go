package aba

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Nesvilab/philosopher/lib/id"
	"github.com/Nesvilab/philosopher/lib/iso"
	"github.com/Nesvilab/philosopher/lib/met"
	"github.com/Nesvilab/philosopher/lib/msg"
	"github.com/Nesvilab/philosopher/lib/rep"
	"github.com/Nesvilab/philosopher/lib/sys"

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
		files, _ := os.ReadDir(i)
		for _, f := range files {
			if strings.Contains(f.Name(), "annotation") {
				var annot = fmt.Sprintf("%s%s%s", i, string(filepath.Separator), f.Name())

				file, e := os.Open(annot)
				if e != nil {
					msg.ReadFile(errors.New("cannot open annotation file"), "fatal")
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {

					if len(scanner.Text()) > 3 {
						names := strings.Fields(scanner.Text())

						if len(names) <= 1 {
							msg.Custom(errors.New("the annotation file looks to be empty"), "error")
						}

						name := i + " " + names[0]
						labels[name] = names[1]
					}
				}

				if e = scanner.Err(); e != nil {
					msg.Custom(errors.New("the annotation file looks to be empty"), "error")
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
			psm.NamedIntensity = make(map[string]float64)
			psm.NamedLabels = make(map[string]iso.Labels)

			psm.DataSet = prjName
			psm.Source = j.Source
			psm.Spectrum = j.Spectrum
			psm.Peptide = j.Peptide
			psm.ModifiedPeptide = j.ModifiedPeptide
			psm.Probability = j.Probability
			psm.Protein = j.Protein
			psm.ProteinStart = j.ProteinStart
			psm.ProteinEnd = j.ProteinEnd
			psm.ProteinDescription = strings.Replace(j.ProteinDescription, ",", " ", -1)
			psm.ProteinID = j.ProteinID
			psm.EntryName = j.EntryName
			psm.GeneName = j.GeneName
			psm.AssumedCharge = j.AssumedCharge
			psm.IsUnique = j.IsUnique
			psm.Purity = j.Purity

			psm.Intensity = j.Intensity
			psm.NamedIntensity[prjName] = j.Intensity

			psm.MappedProteins = j.MappedProteins
			psm.MappedGenes = j.MappedGenes

			if j.PTM == nil {
				psm.PTM = id.PTM{LocalizedPTMSites: map[string]int{}, LocalizedPTMMassDiff: map[string]string{}}
			} else {
				psm.PTM = *j.PTM
			}

			if j.Labels != nil {
				psm.Labels = *j.Labels
				psm.NamedLabels[prjName] = *j.Labels
				if j.Labels.IsUsed {
					psm.IsUsed = true
				}
			}

			evidences = append(evidences, psm)
		}
	}

	if m.Abacus.Labels {
		//savePSMAbacusResult(m.Temp, m.Abacus.Plex, evidences, names, m.Abacus.Unique, true, m.Abacus.Full, labels)
		saveMSstatsResult(m.Temp, m.Abacus.Plex, evidences, true)
	} else {
		//savePSMAbacusResult(m.Temp, m.Abacus.Plex, evidences, names, m.Abacus.Unique, false, m.Abacus.Full, labels)
		saveMSstatsResult(m.Temp, m.Abacus.Plex, evidences, false)
	}

}

// saveMSstatsResult creates a msstats report using 1 or more philosopher result files
func saveMSstatsResult(session, plex string, evidences rep.CombinedPSMEvidenceList, hasLabels bool) {

	var modMap = make(map[string]string)
	var modList []string

	// create result file
	output := fmt.Sprintf("%s%smsstats.csv", session, string(filepath.Separator))

	// create result file
	file, e := os.Create(output)
	if e != nil {
		msg.WriteFile(e, "error")
	}
	defer file.Close()

	for _, i := range evidences {
		for k := range i.PTM.LocalizedPTMMassDiff {
			_, ok := modMap[k]
			if !ok {
				modMap[k] = ""
			} else {
				modMap[k] = ""
			}
		}
	}

	for k := range modMap {
		modList = append(modList, k)
	}

	sort.Strings(modList)

	header := "Spectrum.Name,Spectrum.File,Peptide.Sequence,Modified.Peptide.Sequence,Probability,Charge,Protein.Start,Protein.End,Gene,Mapped.Genes,Protein,Protein.ID,Mapped.Proteins,Protein.Description,Is.Unique,Purity,Intensity"

	if len(modList) > 0 {
		for _, i := range modList {
			header += "," + i
		}
	}

	if hasLabels {
		if plex == "10" {
			header = fmt.Sprintf("%s,Channel 126,Channel 127N,Channel 127C,Channel 128N,Channel 128C,Channel 129N,Channel 129C,Channel 130N,Channel 130C,Channel 131N", header)
		} else if plex == "11" {
			header = fmt.Sprintf("%s,Channel 126,Channel 127N,Channel 127C,Channel 128N,Channel 128C,Channel 129N,Channel 129C,Channel 130N,Channel 130C,Channel 131N,Channel 131C", header)
		} else if plex == "16" {
			header = fmt.Sprintf("%s,Channel 126,Channel 127N,Channel 127C,Channel 128N,Channel 128C,Channel 129N,Channel 129C,Channel 130N,Channel 130C,Channel 131N,Channel 131C,Channel 132N,Channel 132C,Channel 133N,Channel 133C,Channel 134N", header)
		} else if plex == "18" {
			header = fmt.Sprintf("%s,Channel 126,Channel 127N,Channel 127C,Channel 128N,Channel 128C,Channel 129N,Channel 129C,Channel 130N,Channel 130C,Channel 131N,Channel 131C,Channel 32N,Channel 132C,Channel 133N,Channel 133C,Channel 134N,Channel 134C,Channel 135N", header)
		} else {
			msg.Custom(errors.New("unsupported number of labels"), "error")
		}
	}

	header += "\n"
	_, e = io.WriteString(file, header)
	if e != nil {
		msg.WriteToFile(e, "error")
	}

	for _, i := range evidences {
		var line string
		var mods string

		var mappedGenes []string
		for j := range i.MappedGenes {
			if j != i.GeneName && len(j) > 0 {
				mappedGenes = append(mappedGenes, j)
			}
		}
		sort.Strings(mappedGenes)

		var mappedProteins []string
		for j := range i.MappedProteins {
			if j != i.Protein && len(j) > 0 {
				mappedProteins = append(mappedProteins, j)
			}
		}
		sort.Strings(mappedProteins)

		line += fmt.Sprintf("%s,", i.Spectrum)

		line += fmt.Sprintf("%s.mzML,", i.Source)

		line += fmt.Sprintf("%s,", i.Peptide)

		line += fmt.Sprintf("%s,", i.ModifiedPeptide)

		line += fmt.Sprintf("%.4f,", i.Probability)

		line += fmt.Sprintf("%d,", i.AssumedCharge)

		line += fmt.Sprintf("%d,", i.ProteinStart)

		line += fmt.Sprintf("%d,", i.ProteinEnd)

		line += fmt.Sprintf("%s,", i.GeneName)

		line += fmt.Sprintf("%s,", strings.Join(mappedGenes, ";"))

		line += fmt.Sprintf("%s,", i.Protein)

		line += fmt.Sprintf("%s,", i.ProteinID)

		line += fmt.Sprintf("%s,", strings.Join(mappedProteins, ";"))

		line += fmt.Sprintf("%s,", i.ProteinDescription)

		line += fmt.Sprintf("%t,", i.IsUnique)

		line += fmt.Sprintf("%.2f,", i.Purity)

		line += fmt.Sprintf("%6.f,", i.Intensity)

		if len(modList) > 0 {
			for _, j := range modList {
				mods += fmt.Sprintf("%s,", i.PTM.LocalizedPTMMassDiff[j])
			}
		}

		line += mods

		if hasLabels {
			switch plex {
			case "10":
				line += fmt.Sprintf("%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f",
					i.Labels.Channel1.Intensity,
					i.Labels.Channel2.Intensity,
					i.Labels.Channel3.Intensity,
					i.Labels.Channel4.Intensity,
					i.Labels.Channel5.Intensity,
					i.Labels.Channel6.Intensity,
					i.Labels.Channel7.Intensity,
					i.Labels.Channel8.Intensity,
					i.Labels.Channel9.Intensity,
					i.Labels.Channel10.Intensity,
				)
			case "11":
				line += fmt.Sprintf("%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f",
					i.Labels.Channel1.Intensity,
					i.Labels.Channel2.Intensity,
					i.Labels.Channel3.Intensity,
					i.Labels.Channel4.Intensity,
					i.Labels.Channel5.Intensity,
					i.Labels.Channel6.Intensity,
					i.Labels.Channel7.Intensity,
					i.Labels.Channel8.Intensity,
					i.Labels.Channel9.Intensity,
					i.Labels.Channel10.Intensity,
					i.Labels.Channel11.Intensity,
				)
			case "16":
				line += fmt.Sprintf("%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f",
					i.Labels.Channel1.Intensity,
					i.Labels.Channel2.Intensity,
					i.Labels.Channel3.Intensity,
					i.Labels.Channel4.Intensity,
					i.Labels.Channel5.Intensity,
					i.Labels.Channel6.Intensity,
					i.Labels.Channel7.Intensity,
					i.Labels.Channel8.Intensity,
					i.Labels.Channel9.Intensity,
					i.Labels.Channel10.Intensity,
					i.Labels.Channel11.Intensity,
					i.Labels.Channel12.Intensity,
					i.Labels.Channel13.Intensity,
					i.Labels.Channel14.Intensity,
					i.Labels.Channel15.Intensity,
					i.Labels.Channel16.Intensity,
				)
			case "18":
				line += fmt.Sprintf("%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f",
					i.Labels.Channel1.Intensity,
					i.Labels.Channel2.Intensity,
					i.Labels.Channel3.Intensity,
					i.Labels.Channel4.Intensity,
					i.Labels.Channel5.Intensity,
					i.Labels.Channel6.Intensity,
					i.Labels.Channel7.Intensity,
					i.Labels.Channel8.Intensity,
					i.Labels.Channel9.Intensity,
					i.Labels.Channel10.Intensity,
					i.Labels.Channel11.Intensity,
					i.Labels.Channel12.Intensity,
					i.Labels.Channel13.Intensity,
					i.Labels.Channel14.Intensity,
					i.Labels.Channel15.Intensity,
					i.Labels.Channel16.Intensity,
					i.Labels.Channel17.Intensity,
					i.Labels.Channel18.Intensity,
				)
			}
		}

		line += "\n"
		_, e := io.WriteString(file, line)
		if e != nil {
			msg.WriteToFile(e, "error")
		}
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))
}

// savePSMAbacusResult creates a single report using 1 or more philosopher result files
func savePSMAbacusResult(session, plex string, evidences rep.CombinedPSMEvidenceList, namesList []string, hasLabels, full bool, labelsList map[string]string) {

	// create result file
	output := fmt.Sprintf("%s%scombined_psm.tsv", session, string(filepath.Separator))

	// create result file
	file, e := os.Create(output)
	if e != nil {
		msg.WriteFile(e, "fatal")
	}
	defer file.Close()

	header := "Spectrum Name\tSpectrum File\tPeptide\tModified Peptide\tCharge\tGene\tProtein\tProtein ID\tEntry Name\tProtein Description\tIs Unique\tQuan Usage\tPurity"

	// Add Unique+Razor Intensity
	for _, i := range namesList {
		header += fmt.Sprintf("\t%s", i)
	}

	var chs []string

	if plex == "10" {
		chs = append(chs, "126", "127N", "127C", "128N", "128C", "129N", "129C", "130N", "130C", "131N")
	} else if plex == "11" {
		chs = append(chs, "126", "127N", "127C", "128N", "128C", "129N", "129C", "130N", "130C", "131N", "131C")
	} else if plex == "16" {
		chs = append(chs, "126", "127N", "127C", "128N", "128C", "129N", "129C", "130N", "130C", "131N", "131C", "132N", "132C", "133N", "133C", "134N")
	} else if plex == "18" {
		chs = append(chs, "126", "127N", "127C", "128N", "128C", "129N", "129C", "130N", "130C", "131N", "131C", "132N", "132C", "133N", "133C", "134N", "134C", "135N")
	} else {
		msg.Custom(errors.New("unsupported number of labels"), "error")
	}

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
		msg.WriteToFile(e, "error")
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
			line += fmt.Sprintf("%6.f\t", i.NamedIntensity[j])
		}

		if hasLabels {
			switch plex {
			case "10":
				for _, j := range namesList {
					line += fmt.Sprintf("%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
						i.NamedLabels[j].Channel1.Intensity,
						i.NamedLabels[j].Channel2.Intensity,
						i.NamedLabels[j].Channel3.Intensity,
						i.NamedLabels[j].Channel4.Intensity,
						i.NamedLabels[j].Channel5.Intensity,
						i.NamedLabels[j].Channel6.Intensity,
						i.NamedLabels[j].Channel7.Intensity,
						i.NamedLabels[j].Channel8.Intensity,
						i.NamedLabels[j].Channel9.Intensity,
						i.NamedLabels[j].Channel10.Intensity,
					)
				}
			case "16":
				for _, j := range namesList {
					line += fmt.Sprintf("%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
						i.NamedLabels[j].Channel1.Intensity,
						i.NamedLabels[j].Channel2.Intensity,
						i.NamedLabels[j].Channel3.Intensity,
						i.NamedLabels[j].Channel4.Intensity,
						i.NamedLabels[j].Channel5.Intensity,
						i.NamedLabels[j].Channel6.Intensity,
						i.NamedLabels[j].Channel7.Intensity,
						i.NamedLabels[j].Channel8.Intensity,
						i.NamedLabels[j].Channel9.Intensity,
						i.NamedLabels[j].Channel10.Intensity,
						i.NamedLabels[j].Channel11.Intensity,
						i.NamedLabels[j].Channel12.Intensity,
						i.NamedLabels[j].Channel13.Intensity,
						i.NamedLabels[j].Channel14.Intensity,
						i.NamedLabels[j].Channel15.Intensity,
						i.NamedLabels[j].Channel16.Intensity,
					)
				}
			case "18":
				for _, j := range namesList {
					line += fmt.Sprintf("%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
						i.NamedLabels[j].Channel1.Intensity,
						i.NamedLabels[j].Channel2.Intensity,
						i.NamedLabels[j].Channel3.Intensity,
						i.NamedLabels[j].Channel4.Intensity,
						i.NamedLabels[j].Channel5.Intensity,
						i.NamedLabels[j].Channel6.Intensity,
						i.NamedLabels[j].Channel7.Intensity,
						i.NamedLabels[j].Channel8.Intensity,
						i.NamedLabels[j].Channel9.Intensity,
						i.NamedLabels[j].Channel10.Intensity,
						i.NamedLabels[j].Channel11.Intensity,
						i.NamedLabels[j].Channel12.Intensity,
						i.NamedLabels[j].Channel13.Intensity,
						i.NamedLabels[j].Channel14.Intensity,
						i.NamedLabels[j].Channel15.Intensity,
						i.NamedLabels[j].Channel16.Intensity,
						i.NamedLabels[j].Channel17.Intensity,
						i.NamedLabels[j].Channel18.Intensity,
					)
				}
			}
		}

		line += "\n"
		_, e := io.WriteString(file, line)
		if e != nil {
			msg.WriteToFile(e, "error")
		}
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))
}
