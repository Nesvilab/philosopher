package aba

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/prvst/philosopher/lib/dat"
	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/fil"
	"github.com/prvst/philosopher/lib/id"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/rep"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/prvst/philosopher/lib/tmt"
	"github.com/sirupsen/logrus"
)

// Create protein combined report
func proteinLevelAbacus(a met.Abacus, temp string, args []string) error {

	var names []string
	var xmlFiles []string
	var database dat.Base
	var datasets = make(map[string]rep.Evidence)

	var labelList []DataSetLabelNames

	// restore database
	database = dat.Base{}
	database.RestoreWithPath(args[0])

	// restoring combined file
	logrus.Info("Processing combined file")
	evidences, err := processProteinCombinedFile(a, database)
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

	// If the name starts with CONTROL_  or Control_ then we put CONTROL (regardless of what follows after first '_')
	// If the name starts with something else, then we first determine, for each experiment, if the annotation
	// follows GENE_condition_replicate format (meaning there are two '_' in the name) or just GENE_replicate
	// format (meaning there is only one '_')

	// If two '_', then we put in the second row GENE_condition (i.e. remove the second _ and what follows)
	// If only one '_', then we put in the second row just GENE (i.e. remove the first _ and what follows after)

	var reprintLabels []string
	if a.Reprint == true {
		for _, i := range names {
			if strings.Contains(strings.ToUpper(i), "CONTROL") {
				reprintLabels = append(reprintLabels, "CONTROL")
			} else {
				parts := strings.Split(i, "_")
				if len(parts) == 3 {
					label := fmt.Sprintf("%s_%s", parts[0], parts[1])
					reprintLabels = append(reprintLabels, label)
				} else if len(parts) == 2 {
					label := fmt.Sprintf("%s", parts[0])
					reprintLabels = append(reprintLabels, label)
				}
			}
		}
	}

	sort.Strings(names)
	sort.Strings(reprintLabels)

	logrus.Info("Processing spectral counts")
	evidences = getProteinSpectralCounts(evidences, datasets)

	logrus.Info("Processing intensities")
	evidences = sumProteinIntensities(evidences, datasets)

	// collect TMT labels
	if a.Labels == true {
		evidences = getProteinLabelIntensities(evidences, datasets)
	}

	if a.Labels == true {
		saveProteinAbacusResult(temp, evidences, datasets, names, a.Unique, true, labelList)
	} else {
		saveProteinAbacusResult(temp, evidences, datasets, names, a.Unique, false, labelList)
	}

	if a.Reprint == true {
		logrus.Info("Creating Reprint report")
		saveReprintResults(temp, evidences, datasets, names, reprintLabels, a.Unique, false, labelList)
	}

	return nil
}

// processCombinedFile reads the combined protXML and creates a unique protein list as a reference fo all counts
func processProteinCombinedFile(a met.Abacus, database dat.Base) (rep.CombinedProteinEvidenceList, error) {

	var list rep.CombinedProteinEvidenceList

	if _, err := os.Stat(a.CombPro); os.IsNotExist(err) {
		logrus.Fatal("Cannot find combined.prot.xml file")
	} else {

		var protxml id.ProtXML
		protxml.Read(a.CombPro)
		protxml.DecoyTag = a.Tag

		// promote decoy proteins with indistinguishable target proteins
		protxml.PromoteProteinIDs()

		// applies pickedFDR algorithm
		if a.Picked == true {
			protxml = fil.PickedFDR(protxml)
		}

		// applies razor algorithm
		if a.Razor == true {
			protxml, err = fil.RazorFilter(protxml)
			if err != nil {
				return list, err
			}
		}

		proid, err := fil.ProtXMLFilter(protxml, 0.01, a.PepProb, a.ProtProb, a.Picked, a.Razor)
		if err != nil {
			return nil, err
		}

		for _, j := range proid {

			if !strings.Contains(j.ProteinName, a.Tag) {

				var ce rep.CombinedProteinEvidence

				ce.TotalSpc = make(map[string]int)
				ce.UniqueSpc = make(map[string]int)
				ce.UrazorSpc = make(map[string]int)

				ce.TotalIntensity = make(map[string]float64)
				ce.UniqueIntensity = make(map[string]float64)
				ce.UrazorIntensity = make(map[string]float64)

				ce.TotalLabels = make(map[string]tmt.Labels)
				ce.UniqueLabels = make(map[string]tmt.Labels)
				ce.URazorLabels = make(map[string]tmt.Labels)

				ce.SupportingSpectra = make(map[string]string)
				ce.ProteinName = j.ProteinName
				ce.Length = j.Length
				ce.Coverage = j.PercentCoverage
				ce.GroupNumber = j.GroupNumber
				ce.SiblingID = j.GroupSiblingID
				ce.IndiProtein = j.IndistinguishableProtein
				ce.UniqueStrippedPeptides = len(j.UniqueStrippedPeptides)
				ce.PeptideIons = j.PeptideIons
				ce.ProteinProbability = j.Probability
				ce.TopPepProb = j.TopPepProb

				list = append(list, ce)
			}
		}

	}

	for i := range list {
		for _, j := range database.Records {
			if strings.Contains(j.OriginalHeader, list[i].ProteinName) {
				list[i].ProteinID = j.ID
				list[i].EntryName = j.EntryName
				list[i].GeneNames = j.GeneNames
				list[i].Organism = j.Organism
				list[i].Description = j.Description
				list[i].ProteinExistence = j.ProteinExistence
				break
			}
		}
	}

	return list, nil
}

func getProteinSpectralCounts(combined rep.CombinedProteinEvidenceList, datasets map[string]rep.Evidence) rep.CombinedProteinEvidenceList {

	for k, v := range datasets {

		for i := range combined {
			for _, j := range v.Proteins {
				if combined[i].ProteinID == j.ProteinID {
					combined[i].UniqueSpc[k] = j.UniqueSpC
					combined[i].TotalSpc[k] = j.TotalSpC
					combined[i].UrazorSpc[k] = j.URazorSpC
					break
				}
			}
		}

	}

	return combined
}

func getProteinLabelIntensities(combined rep.CombinedProteinEvidenceList, datasets map[string]rep.Evidence) rep.CombinedProteinEvidenceList {

	for k, v := range datasets {

		for i := range combined {
			for _, j := range v.Proteins {
				if combined[i].ProteinID == j.ProteinID {
					combined[i].TotalLabels[k] = j.TotalLabels
					combined[i].UniqueLabels[k] = j.UniqueLabels
					combined[i].URazorLabels[k] = j.URazorLabels
					break
				}
			}
		}

	}

	return combined
}

// sumIntensities calculates the protein intensity
func sumProteinIntensities(combined rep.CombinedProteinEvidenceList, datasets map[string]rep.Evidence) rep.CombinedProteinEvidenceList {

	for k, v := range datasets {

		var ions = make(map[string]float64)
		for _, i := range v.Ions {
			ions[i.IonForm] = i.Intensity
		}

		for _, i := range combined {
			for j := range v.Proteins {
				if i.ProteinID == v.Proteins[j].ProteinID {
					i.TotalIntensity[k] = v.Proteins[j].TotalIntensity
					i.UniqueIntensity[k] = v.Proteins[j].UniqueIntensity
					i.UrazorIntensity[k] = v.Proteins[j].URazorIntensity
					break
				}
			}
		}

	}

	return combined
}

// saveProteinAbacusResult creates a single report using 1 or more philosopher result files
func saveProteinAbacusResult(session string, evidences rep.CombinedProteinEvidenceList, datasets map[string]rep.Evidence, namesList []string, uniqueOnly, hasTMT bool, labelsList []DataSetLabelNames) {

	// create result file
	output := fmt.Sprintf("%s%scombined_protein.tsv", session, string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create report file:", err)
	}
	defer file.Close()

	line := "Protein Group\tSubGroup\tProtein ID\tEntry Name\tGene Names\tProtein Length\tCoverage\tOrganism\tProtein Existence\tDescription\tProtein Probability\tTop Peptide Probability\tUnique Stripped Peptides\tSummarized Total Spectral Count\tSummarized Unique Spectral Count\tSummarized Razor Spectral Count\t"

	for _, i := range namesList {
		line += fmt.Sprintf("%s Total Spectral Count\t", i)
		line += fmt.Sprintf("%s Unique Spectral Count\t", i)
		line += fmt.Sprintf("%s Razor Spectral Count\t", i)
		line += fmt.Sprintf("%s Total Intensity\t", i)
		line += fmt.Sprintf("%s Unique Intensity\t", i)
		line += fmt.Sprintf("%s Razor Intensity\t", i)
	}

	if hasTMT == true {
		for _, i := range namesList {
			line += fmt.Sprintf("%s 126 Abundance\t", i)
			line += fmt.Sprintf("%s 127N Abundance\t", i)
			line += fmt.Sprintf("%s 127C Abundance\t", i)
			line += fmt.Sprintf("%s 128N Abundance\t", i)
			line += fmt.Sprintf("%s 128C Abundance\t", i)
			line += fmt.Sprintf("%s 129N Abundance\t", i)
			line += fmt.Sprintf("%s 129C Abundance\t", i)
			line += fmt.Sprintf("%s 130N Abundance\t", i)
			line += fmt.Sprintf("%s 130C Abundance\t", i)
			line += fmt.Sprintf("%s 131N Abundance\t", i)

			for _, j := range labelsList {
				if j.Name == i {
					for k, v := range j.LabelName {
						before := fmt.Sprintf("%s %s Abundance", i, k)
						after := fmt.Sprintf("%s Abundance", v)
						line = strings.Replace(line, before, after, -1)
					}
				}
			}
		}

	}

	line += "Indistinguishable Proteins\t"

	line += "\n"
	n, err := io.WriteString(file, line)
	if err != nil {
		logrus.Fatal(n, err)
	}

	// organize by group number
	sort.Sort(evidences)

	var summTotalSpC = make(map[string]int)
	var summUniqueSpC = make(map[string]int)
	var summURazorSpC = make(map[string]int)

	// collect and sum all evidences from all data sets for each protein
	for _, i := range evidences {
		for _, j := range namesList {
			summTotalSpC[i.ProteinID] += i.TotalSpc[j]
			summUniqueSpC[i.ProteinID] += i.UniqueSpc[j]
			summURazorSpC[i.ProteinID] += i.UrazorSpc[j]
		}
	}

	for _, i := range evidences {

		var line string

		line += fmt.Sprintf("%d\t", i.GroupNumber)

		line += fmt.Sprintf("%s\t", i.SiblingID)

		line += fmt.Sprintf("%s\t", i.ProteinID)

		line += fmt.Sprintf("%s\t", i.EntryName)

		line += fmt.Sprintf("%s\t", i.GeneNames)

		line += fmt.Sprintf("%d\t", i.Length)

		line += fmt.Sprintf("%d\t", int(i.Coverage))

		line += fmt.Sprintf("%s\t", i.Organism)

		line += fmt.Sprintf("%s\t", i.ProteinExistence)

		line += fmt.Sprintf("%s\t", i.Description)

		line += fmt.Sprintf("%.4f\t", i.ProteinProbability)

		line += fmt.Sprintf("%.4f\t", i.TopPepProb)

		line += fmt.Sprintf("%d\t", i.UniqueStrippedPeptides)

		line += fmt.Sprintf("%d\t", summTotalSpC[i.ProteinID])

		line += fmt.Sprintf("%d\t", summUniqueSpC[i.ProteinID])

		line += fmt.Sprintf("%d\t", summURazorSpC[i.ProteinID])

		for _, j := range namesList {
			line += fmt.Sprintf("%d\t%d\t%d\t%6.f\t%6.f\t%6.f\t", i.TotalSpc[j], i.UniqueSpc[j], i.UrazorSpc[j], i.TotalIntensity[j], i.UniqueIntensity[j], i.UrazorIntensity[j])
		}

		if hasTMT == true {
			if uniqueOnly == true {
				for _, j := range namesList {
					line += fmt.Sprintf("%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t",
						i.UniqueLabels[j].Channel1.Intensity,
						i.UniqueLabels[j].Channel2.Intensity,
						i.UniqueLabels[j].Channel3.Intensity,
						i.UniqueLabels[j].Channel4.Intensity,
						i.UniqueLabels[j].Channel5.Intensity,
						i.UniqueLabels[j].Channel6.Intensity,
						i.UniqueLabels[j].Channel7.Intensity,
						i.UniqueLabels[j].Channel8.Intensity,
						i.UniqueLabels[j].Channel9.Intensity,
						i.UniqueLabels[j].Channel10.Intensity,
					)
				}
			} else {
				for _, j := range namesList {
					line += fmt.Sprintf("%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t",
						i.URazorLabels[j].Channel1.Intensity,
						i.URazorLabels[j].Channel2.Intensity,
						i.URazorLabels[j].Channel3.Intensity,
						i.URazorLabels[j].Channel4.Intensity,
						i.URazorLabels[j].Channel5.Intensity,
						i.URazorLabels[j].Channel6.Intensity,
						i.URazorLabels[j].Channel7.Intensity,
						i.URazorLabels[j].Channel8.Intensity,
						i.URazorLabels[j].Channel9.Intensity,
						i.URazorLabels[j].Channel10.Intensity,
					)
				}
			}
		}

		ip := strings.Join(i.IndiProtein, ", ")
		line += fmt.Sprintf("%s\t", ip)

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

// saveReprintResults creates a single report using 1 or more philosopher result files using the Reprint format
func saveReprintResults(session string, evidences rep.CombinedProteinEvidenceList, datasets map[string]rep.Evidence, namesList, labelList []string, uniqueOnly, hasTMT bool, labelsList []DataSetLabelNames) {

	// create result file
	output := fmt.Sprintf("%s%sreprint.tsv", session, string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create report file:", err)
	}
	defer file.Close()

	line := "PROTID\t"

	for _, i := range namesList {
		line += fmt.Sprintf("%s_SPC\t", i)
	}

	line += "\n"
	line += "na\t"

	for _, i := range labelList {
		line += fmt.Sprintf("%s\t", i)
	}

	line += "\n"

	n, err := io.WriteString(file, line)
	if err != nil {
		logrus.Fatal(n, err)
	}

	// organize by group number
	sort.Sort(evidences)

	var summTotalSpC = make(map[string]int)
	var summUniqueSpC = make(map[string]int)
	var summURazorSpC = make(map[string]int)

	// collect and sum all evidences from all data sets for each protein
	for _, i := range evidences {
		for _, j := range namesList {
			summTotalSpC[i.ProteinID] += i.TotalSpc[j]
			summUniqueSpC[i.ProteinID] += i.UniqueSpc[j]
			summURazorSpC[i.ProteinID] += i.UrazorSpc[j]
		}
	}

	for _, i := range evidences {

		var line string

		line += fmt.Sprintf("%s\t", i.ProteinID)

		for _, j := range namesList {
			line += fmt.Sprintf("%d\t", i.UrazorSpc[j])
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

// addCustomNames adds to the label structures user-defined names to be used on the TMT labels
func getLabelNames(annot string) (map[string]string, *err.Error) {

	var labels = make(map[string]string)

	file, e := os.Open(annot)
	if e != nil {
		return labels, &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: e.Error()}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		names := strings.Split(scanner.Text(), " ")
		labels[names[0]] = names[1]
	}

	if e = scanner.Err(); e != nil {
		return labels, &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: e.Error()}
	}

	return labels, nil
}
