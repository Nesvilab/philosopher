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
	"github.com/prvst/philosopher/lib/data"
	"github.com/prvst/philosopher/lib/fil"
	"github.com/prvst/philosopher/lib/meta"
	"github.com/prvst/philosopher/lib/rep"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/prvst/philosopher/lib/xml"
)

// Abacus main structure
type Abacus struct {
	meta.Data
	ProtProb float64
	PepProb  float64
	Comb     string
	Razor    bool
	Tag      string
}

// ExperimentalData ...
type ExperimentalData struct {
	Name        string
	PeptideIons map[string]int
}

// ExperimentalDataList ...
type ExperimentalDataList []ExperimentalData

// New constructor
func New() Abacus {

	var o Abacus
	var m meta.Data
	m.Restore(sys.Meta())

	o.UUID = m.UUID
	o.Distro = m.Distro
	o.Home = m.Home
	o.MetaFile = m.MetaFile
	o.MetaDir = m.MetaDir
	o.DB = m.DB
	o.Temp = m.Temp
	o.TimeStamp = m.TimeStamp
	o.OS = m.OS
	o.Arch = m.Arch

	return o
}

// Run abacus
func (a *Abacus) Run(args []string) error {

	var names []string
	var xmlFiles []string
	var database data.Base
	var datasets = make(map[string]rep.Evidence)

	// restore database
	database = data.Base{}
	database.RestoreWithPath(args[0])

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
				interactFile := fmt.Sprintf("%s%s", i, f.Name())
				absPath, _ := filepath.Abs(interactFile)
				xmlFiles = append(xmlFiles, absPath)
			}
		}

		// collect project names
		prjName := i
		if strings.Contains(prjName, string(filepath.Separator)) {
			prjName = strings.Replace(prjName, string(filepath.Separator), "", -1)
		}

		// unique list and map of datasets
		datasets[prjName] = e
		names = append(names, prjName)
	}

	sort.Strings(names)

	logrus.Info("Processing combined file")
	evidences, err := a.processCombinedFile(a.Comb, a.Tag, a.PepProb, a.ProtProb, database)
	if err != nil {
		return err
	}

	evidences = getSpectralCounts(evidences, datasets)

	evidences = sumIntensities(evidences, datasets)

	saveCompareResults(a.Temp, evidences, names)

	return nil
}

// processCombinedFile reads the combined protXML and creates a unique protein list as a reference fo all counts
func (a *Abacus) processCombinedFile(combinedFile, decoyTag string, pepProb, protProb float64, database data.Base) (rep.CombinedEvidenceList, error) {

	var list rep.CombinedEvidenceList

	if _, err := os.Stat(combinedFile); os.IsNotExist(err) {
		logrus.Fatal("Cannot find combined.prot.xml file")
	} else {

		var protxml xml.ProtXML
		protxml.Read(combinedFile)
		protxml.DecoyTag = decoyTag

		// promote decoy proteins with indistinguishable target proteins
		protxml.PromoteProteinIDs()

		// applies razor algorithm
		if a.Razor == true {
			protxml, err = fil.RazorFilter(protxml)
			if err != nil {
				return list, err
			}
		}

		proid, err := fil.ProtXMLFilter(protxml, 0.01, pepProb, protProb, false, a.Razor)
		if err != nil {
			return nil, err
		}

		for _, j := range proid {

			var ce rep.CombinedEvidence

			ce.TotalSpc = make(map[string]int)
			ce.UniqueSpc = make(map[string]int)
			ce.UrazorSpc = make(map[string]int)

			ce.TotalIntensity = make(map[string]float64)
			ce.UniqueIntensity = make(map[string]float64)
			ce.UrazorIntensity = make(map[string]float64)

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

		var ions = make(map[string]int)

		for _, i := range v.PSM {
			ions[i.IonForm]++
		}

		for _, i := range v.Proteins {
			for j := range combined {
				if i.ProteinID == combined[j].ProteinID {
					combined[j].UniqueSpc[k] = i.UniqueSpC
					combined[j].TotalSpc[k] = i.TotalSpC
					break
				}
			}
		}

		for i := range combined {
			for _, j := range combined[i].PeptideIons {

				ion := fmt.Sprintf("%s#%d#%.4f", j.PeptideSequence, j.Charge, j.CalcNeutralPepMass)

				sum, ok := ions[ion]
				if ok {

					//combined[i].TotalSpc[k] += sum

					// if j.IsUnique == true {
					// 	combined[i].UniqueSpc[k] += sum
					// }

					if j.Razor == 1 {
						combined[i].UrazorSpc[k] += sum
					}

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
// 		var ions = make(map[string]int)
//
// 		for _, i := range v.PSM {
// 			ions[i.IonForm]++
// 		}
//
// 		for i := range combined {
// 			for _, j := range combined[i].PeptideIons {
//
// 				ion := fmt.Sprintf("%s#%d#%.4f", j.PeptideSequence, j.Charge, j.CalcNeutralPepMass)
//
// 				sum, ok := ions[ion]
// 				if ok {
//
// 					combined[i].TotalSpc[k] += sum
//
// 					if j.IsUnique == true {
// 						combined[i].UniqueSpc[k] += sum
// 					}
//
// 					if j.Razor == 1 {
// 						combined[i].UrazorSpc[k] += sum
// 					}
//
// 				}
//
// 			}
// 		}
//
// 	}
//
// 	return combined
// }

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

// func sumIntensities(combined rep.CombinedEvidenceList, datasets map[string]rep.Evidence) rep.CombinedEvidenceList {
//
// 	//var totalInt []float64
// 	//var uniqueInt []float64
// 	//var urazorInt []float64
//
// 	for k, v := range datasets {
//
// 		var ions = make(map[string][]float64)
//
// 		// for _, i := range v.PSM {
// 		// 	ions[i.IonForm] = append(ions[i.IonForm], i.Intensity)
// 		// }
//
// 		var intIonMap = make(map[string]float64)
// 		for _, i := range v.PSM {
// 			// global intensity map for Ions, getting the most intense
// 			_, ok := intIonMap[i.IonForm]
// 			if ok {
// 				if i.Intensity > intIonMap[i.IonForm] {
// 					intIonMap[i.IonForm] = i.Intensity
// 				}
// 			} else {
// 				intIonMap[i.IonForm] = i.Intensity
// 			}
// 		}
//
// 		for _, i := range v.Proteins {
// 			for j := range combined {
// 				if i.ProteinID == combined[j].ProteinID {
// 					combined[j].TotalIntensity[k] = i.TotalIntensity
// 					combined[j].UniqueIntensity[k] = i.UniqueIntensity
// 					break
// 				}
// 			}
// 		}
//
// 		for i := range combined {
//
// 			var urazorInt []float64
//
// 			for _, j := range combined[i].PeptideIons {
//
// 				ion := fmt.Sprintf("%s#%d#%.4f", j.PeptideSequence, j.Charge, j.CalcNeutralPepMass)
//
// 				intList, ok := ions[ion]
// 				if ok {
//
// 					// totalInt = intList
// 					// sort.Float64s(totalInt)
// 					//
// 					// if len(totalInt) >= 3 {
// 					// 	combined[i].TotalIntensity[k] = (totalInt[len(totalInt)-1] + totalInt[len(totalInt)-2] + totalInt[len(totalInt)-3])
// 					// } else if len(totalInt) >= 2 {
// 					// 	combined[i].TotalIntensity[k] = (totalInt[len(totalInt)-1] + totalInt[len(totalInt)-2])
// 					// } else if len(totalInt) == 1 {
// 					// 	combined[i].TotalIntensity[k] = (totalInt[len(totalInt)-1])
// 					// }
//
// 					// if j.IsUnique == true {
// 					// 	uniqueInt = intList
// 					// 	sort.Float64s(uniqueInt)
// 					//
// 					// 	if len(uniqueInt) >= 3 {
// 					// 		combined[i].UniqueIntensity[k] = (uniqueInt[len(uniqueInt)-1] + uniqueInt[len(uniqueInt)-2] + uniqueInt[len(uniqueInt)-3])
// 					// 	} else if len(uniqueInt) >= 2 {
// 					// 		combined[i].UniqueIntensity[k] = (uniqueInt[len(uniqueInt)-1] + uniqueInt[len(uniqueInt)-2])
// 					// 	} else if len(uniqueInt) == 1 {
// 					// 		combined[i].UniqueIntensity[k] = (uniqueInt[len(uniqueInt)-1])
// 					// 	}
// 					//
// 					// }
//
// 					if j.Razor == 1 {
// 						urazorInt = intList
// 						sort.Float64s(urazorInt)
//
// 						if len(urazorInt) >= 3 {
// 							combined[i].UrazorIntensity[k] = (urazorInt[len(urazorInt)-1] + urazorInt[len(urazorInt)-2] + urazorInt[len(urazorInt)-3])
// 						} else if len(urazorInt) >= 2 {
// 							combined[i].UrazorIntensity[k] = (urazorInt[len(urazorInt)-1] + urazorInt[len(urazorInt)-2])
// 						} else if len(urazorInt) == 1 {
// 							combined[i].UrazorIntensity[k] = (urazorInt[len(urazorInt)-1])
// 						}
// 					}
//
// 				}
//
// 			}
// 		}
//
// 	}
//
// 	return combined
// }

// saveCompareResults creates a single report using 1 or more philosopher result files
func saveCompareResults(session string, evidences rep.CombinedEvidenceList, namesList []string) {

	// create result file
	output := fmt.Sprintf("%s%scombined.tsv", session, string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create report file:", err)
	}
	defer file.Close()

	//line := "Protein Group\tSubGroup\tProtein ID\tEntry Name\tGene Names\tDescription\tProtein Length\tProtein Probability\tTop Peptide Probability\tUnique Stripped Peptides\tTotal Peptide Ions\tUnique Peptide Ions\tRazor Peptide Ions\tIndistinguishable Proteins\t"
	line := "Protein Group\tSubGroup\tProtein ID\tEntry Name\tGene Names\tProtein Length\tProtein Probability\tTop Peptide Probability\tUnique Stripped Peptides\tTotal Peptide Ions\t"

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

		line += fmt.Sprintf("%d\t", len(i.PeptideIons))

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
