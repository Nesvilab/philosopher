package aba

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher/lib/dat"
	"github.com/prvst/philosopher/lib/fil"
	"github.com/prvst/philosopher/lib/id"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/rep"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/prvst/philosopher/lib/tmt"
)

// ExperimentalData ...
type ExperimentalData struct {
	Name        string
	PeptideIons map[string]int
}

// ExperimentalDataList ...
type ExperimentalDataList []ExperimentalData

// Run abacus
func Run(a met.Abacus, temp string, args []string) error {

	var names []string
	var xmlFiles []string
	var database dat.Base
	var datasets = make(map[string]rep.Evidence)

	// restore database
	database = dat.Base{}
	database.RestoreWithPath(args[0])

	// restoring combined file
	logrus.Info("Processing combined file")
	evidences, err := processCombinedFile(a, database)
	if err != nil {
		return err
	}

	// recover all files
	logrus.Info("Restoring results")

	for _, i := range args {

		// restoring the database
		var e rep.Evidence
		e.RestoreGranularWithPath(i)

		// collect interact full file names
		files, _ := ioutil.ReadDir(i)
		for _, f := range files {
			if strings.Contains(f.Name(), "pep.xml") {
				interactFile := fmt.Sprintf("%s%s%s", i, string(filepath.Separator), f.Name())
				absPath, _ := filepath.Abs(interactFile)
				xmlFiles = append(xmlFiles, absPath)
			}
		}

		// collect project names
		prjName := i
		if strings.Contains(prjName, string(filepath.Separator)) {
			prjName = strings.Replace(filepath.Base(prjName), string(filepath.Separator), "", -1)
		}

		// unique list and map of datasets
		datasets[prjName] = e
		names = append(names, prjName)
	}

	sort.Strings(names)

	logrus.Info("Processing spectral counts")
	evidences = getSpectralCounts(evidences, datasets)

	logrus.Info("Processing intensities")
	evidences = sumIntensities(evidences, datasets)

	// collect TMT labels
	if a.Labels == true {
		evidences = getLabelIntensities(evidences, datasets)
	}

	// collect TMT labels
	if a.Labels == true {
		saveCompareTMTResults(temp, evidences, datasets, names)
	} else {
		saveCompareResults(temp, evidences, datasets, names)
	}

	return nil
}

// processCombinedFile reads the combined protXML and creates a unique protein list as a reference fo all counts
func processCombinedFile(a met.Abacus, database dat.Base) (rep.CombinedEvidenceList, error) {

	var list rep.CombinedEvidenceList

	if _, err := os.Stat(a.Comb); os.IsNotExist(err) {
		logrus.Fatal("Cannot find combined.prot.xml file")
	} else {

		var protxml id.ProtXML
		protxml.Read(a.Comb)
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

				var ce rep.CombinedEvidence

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
				break
			}
		}
	}

	return list, nil
}

func getSpectralCounts(combined rep.CombinedEvidenceList, datasets map[string]rep.Evidence) rep.CombinedEvidenceList {

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

// func getSpectralCounts(combined rep.CombinedEvidenceList, datasets map[string]rep.Evidence) rep.CombinedEvidenceList {
//
// 	for k, v := range datasets {
//
// 		//var ions = make(map[string]int)
// 		//var exclusion = make(map[string]uint8)
//
// 		// for _, i := range v.PSM {
// 		// 	ions[i.IonForm]++
// 		// }
//
// 		for i := range combined {
// 			for _, j := range v.Proteins {
// 				if combined[i].ProteinID == j.ProteinID {
// 					combined[i].UniqueSpc[k] = j.UniqueSpC
// 					combined[i].TotalSpc[k] = j.TotalSpC
// 					combined[i].UrazorSpc[k] = j.URazorSpC
// 					break
// 				}
// 			}
// 		}
//
// 		// for _, i := range v.Proteins {
// 		// 	for j := range combined {
// 		// 		if i.ProteinID == combined[j].ProteinID {
// 		// 			combined[j].UniqueSpc[k] = i.UniqueSpC
// 		// 			combined[j].TotalSpc[k] = i.TotalSpC
// 		// 			combined[j].UrazorSpc[k] = i.URazorSpC
// 		//
// 		// 			if combined[j].ProteinID == "NP_078966" {
// 		// 				fmt.Println("PING combined")
// 		// 				os.Exit(1)
// 		// 			}
// 		//
// 		// 			if i.ProteinID == "NP_078966" && combined[j].ProteinID == "NP_078966" {
// 		// 				fmt.Println("PING")
// 		// 				os.Exit(1)
// 		// 			}
// 		//
// 		// 			break
// 		// 		}
// 		// 	}
// 		// }
//
// 		// for i := range combined {
// 		// 	for _, j := range combined[i].PeptideIons {
// 		// 		ion := fmt.Sprintf("%s#%d#%.4f", j.PeptideSequence, j.Charge, j.CalcNeutralPepMass)
// 		// 		sum, ok := ions[ion]
// 		// 		if ok {
// 		// 			_, excl := exclusion[ion]
// 		// 			if !excl {
// 		// 				if j.Razor == 1 {
// 		// 					combined[i].UrazorSpc[k] += sum
// 		// 					exclusion[ion] = 0
// 		// 				}
// 		// 			}
// 		// 		}
// 		// 	}
// 		// }
//
// 	}
//
// 	// for _, i := range combined {
// 	// 	if i.ProteinID == "" {
// 	// 		litter.Dump(i.TotalSpc)
// 	// 		litter.Dump(i.UniqueSpc)
// 	// 		litter.Dump(i.UrazorSpc)
// 	// 	}
// 	// }
//
// 	fmt.Println(len(combined))
//
// 	return combined
// }

func getLabelIntensities(combined rep.CombinedEvidenceList, datasets map[string]rep.Evidence) rep.CombinedEvidenceList {

	for k, v := range datasets {

		for _, i := range v.Proteins {
			for j := range combined {
				if i.ProteinID == combined[j].ProteinID {
					combined[j].TotalLabels[k] = i.TotalLabels
					combined[j].UniqueLabels[k] = i.UniqueLabels
					combined[j].URazorLabels[k] = i.URazorLabels
					break
				}
			}
		}

	}

	return combined
}

// sumIntensities calculates the protein intensity
func sumIntensities(combined rep.CombinedEvidenceList, datasets map[string]rep.Evidence) rep.CombinedEvidenceList {

	for k, v := range datasets {

		var ions = make(map[string]float64)
		for _, i := range v.Ions {
			ions[i.IonForm] = i.Intensity
		}

		for _, i := range v.Proteins {
			for j := range combined {
				if i.ProteinID == combined[j].ProteinID {
					combined[j].TotalIntensity[k] = i.TotalIntensity
					combined[j].UniqueIntensity[k] = i.UniqueIntensity
					break
				}
			}
		}

		for i := range combined {

			var urazorInt []float64

			for _, j := range combined[i].PeptideIons {

				ion := fmt.Sprintf("%s#%d#%.4f", j.PeptideSequence, j.Charge, j.CalcNeutralPepMass)

				intList, ok := ions[ion]
				if ok {

					if j.Razor == 1 {
						urazorInt = append(urazorInt, intList)
						sort.Float64s(urazorInt)

						if len(urazorInt) >= 3 {
							combined[i].UrazorIntensity[k] = (urazorInt[len(urazorInt)-1] + urazorInt[len(urazorInt)-2] + urazorInt[len(urazorInt)-3])
						} else if len(urazorInt) >= 2 {
							combined[i].UrazorIntensity[k] = (urazorInt[len(urazorInt)-1] + urazorInt[len(urazorInt)-2])
						} else if len(urazorInt) == 1 {
							combined[i].UrazorIntensity[k] = (urazorInt[len(urazorInt)-1])
						}
					}

				}

			}
		}

	}

	return combined
}

// saveCompareResults creates a single report using 1 or more philosopher result files
func saveCompareResults(session string, evidences rep.CombinedEvidenceList, datasets map[string]rep.Evidence, namesList []string) {

	// create result file
	output := fmt.Sprintf("%s%scombined.tsv", session, string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create report file:", err)
	}
	defer file.Close()

	//line := "Protein Group\tSubGroup\tProtein ID\tEntry Name\tGene Names\tDescription\tProtein Length\tProtein Probability\tTop Peptide Probability\tUnique Stripped Peptides\tTotal Peptide Ions\tUnique Peptide Ions\tRazor Peptide Ions\tIndistinguishable Proteins\t"
	//line := "Protein Group\tSubGroup\tProtein ID\tEntry Name\tGene Names\tProtein Length\tProtein Probability\tTop Peptide Probability\tUnique Stripped Peptides\tTotal Peptide Ions\t"
	line := "Protein Group\tSubGroup\tProtein ID\tEntry Name\tGene Names\tProtein Length\tProtein Probability\tTop Peptide Probability\tUnique Stripped Peptides\t"

	for _, i := range namesList {
		line += fmt.Sprintf("%s Total Spectral Count\t", i)
		line += fmt.Sprintf("%s Unique Spectral Count\t", i)
		line += fmt.Sprintf("%s Razor Spectral Count\t", i)
		line += fmt.Sprintf("%s Total Intensity\t", i)
		line += fmt.Sprintf("%s Unique Intensity\t", i)
		line += fmt.Sprintf("%s Razor Intensity\t", i)
	}

	line += "Indistinguishable Proteins\t"

	line += "\n"
	n, err := io.WriteString(file, line)
	if err != nil {
		logrus.Fatal(n, err)
	}

	// organize by group number
	sort.Sort(evidences)

	for _, i := range evidences {

		var line string

		line += fmt.Sprintf("%d\t", i.GroupNumber)

		line += fmt.Sprintf("%s\t", i.SiblingID)

		line += fmt.Sprintf("%s\t", i.ProteinID)

		line += fmt.Sprintf("%s\t", i.EntryName)

		line += fmt.Sprintf("%s\t", i.GeneNames)

		line += fmt.Sprintf("%d\t", i.Length)

		line += fmt.Sprintf("%.4f\t", i.ProteinProbability)

		line += fmt.Sprintf("%.4f\t", i.TopPepProb)

		line += fmt.Sprintf("%d\t", i.UniqueStrippedPeptides)

		//line += fmt.Sprintf("%d\t", len(i.PeptideIons))

		for _, j := range namesList {
			line += fmt.Sprintf("%d\t%d\t%d\t%6.f\t%6.f\t%6.f\t", i.TotalSpc[j], i.UniqueSpc[j], i.UrazorSpc[j], i.TotalIntensity[j], i.UniqueIntensity[j], i.UrazorIntensity[j])
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

// saveCompareTMTResults creates a single report using 1 or more philosopher result files
func saveCompareTMTResults(session string, evidences rep.CombinedEvidenceList, datasets map[string]rep.Evidence, namesList []string) {

	// create result file
	output := fmt.Sprintf("%s%scombined.tsv", session, string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create report file:", err)
	}
	defer file.Close()

	//line := "Protein Group\tSubGroup\tProtein ID\tEntry Name\tGene Names\tDescription\tProtein Length\tProtein Probability\tTop Peptide Probability\tUnique Stripped Peptides\tTotal Peptide Ions\tUnique Peptide Ions\tRazor Peptide Ions\tIndistinguishable Proteins\t"
	//line := "Protein Group\tSubGroup\tProtein ID\tEntry Name\tGene Names\tProtein Length\tProtein Probability\tTop Peptide Probability\tUnique Stripped Peptides\tTotal Peptide Ions\t"
	line := "Protein Group\tSubGroup\tProtein ID\tEntry Name\tGene Names\tProtein Length\tProtein Probability\tTop Peptide Probability\tUnique Stripped Peptides\t"

	for _, i := range namesList {
		line += fmt.Sprintf("%s Total Spectral Count\t", i)
		line += fmt.Sprintf("%s Unique Spectral Count\t", i)
		line += fmt.Sprintf("%s Razor Spectral Count\t", i)
		line += fmt.Sprintf("%s Total Intensity\t", i)
		line += fmt.Sprintf("%s Unique Intensity\t", i)
		line += fmt.Sprintf("%s Razor Intensity\t", i)
	}

	line += "Indistinguishable Proteins\t"

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

		line += fmt.Sprintf("%d\t", i.GroupNumber)

		line += fmt.Sprintf("%s\t", i.SiblingID)

		line += fmt.Sprintf("%s\t", i.ProteinID)

		line += fmt.Sprintf("%s\t", i.EntryName)

		line += fmt.Sprintf("%s\t", i.GeneNames)

		line += fmt.Sprintf("%d\t", i.Length)

		line += fmt.Sprintf("%.4f\t", i.ProteinProbability)

		line += fmt.Sprintf("%.4f\t", i.TopPepProb)

		line += fmt.Sprintf("%d\t", i.UniqueStrippedPeptides)

		//line += fmt.Sprintf("%d\t", len(i.PeptideIons))

		for _, j := range namesList {
			line += fmt.Sprintf("%d\t%d\t%d\t%6.f\t%6.f\t%6.f\t", i.TotalSpc[j], i.UniqueSpc[j], i.UrazorSpc[j], i.TotalIntensity[j], i.UniqueIntensity[j], i.UrazorIntensity[j])
		}

		ip := strings.Join(i.IndiProtein, ", ")
		line += fmt.Sprintf("%s\t", ip)

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
