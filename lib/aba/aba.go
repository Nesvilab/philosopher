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
	//var globalPepMap = make(map[string]int)
	var PsmMap = make(map[string]rep.PSMEvidenceList)

	var totalIons = make(map[string][]string)
	var uniqueIons = make(map[string][]string)
	var razorIons = make(map[string][]string)

	// recover all files
	logrus.Info("Restoring Philospher results")

	for _, i := range args {

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

		// restore database
		database = data.Base{}
		database.RestoreWithPath(i)

		// collect project names
		prjName := i
		if strings.Contains(prjName, string(filepath.Separator)) {
			prjName = strings.Replace(prjName, string(filepath.Separator), "", -1)
		}
		names = append(names, prjName)
		sort.Strings(names)

		for _, j := range e.PSM {
			var ion string

			if len(j.ModifiedPeptide) > 0 {
				ion = fmt.Sprintf("%s#%d", j.ModifiedPeptide, j.AssumedCharge)
			} else {
				ion = fmt.Sprintf("%s#%d", j.Peptide, j.AssumedCharge)
			}

			key := fmt.Sprintf("%s#%s", prjName, j.Protein)

			totalIons[key] = append(totalIons[key], ion)

			if j.IsUnique == true {
				uniqueIons[key] = append(uniqueIons[key], ion)
			}

			if j.IsURazor == true {
				razorIons[key] = append(razorIons[key], ion)
			}

		}
		PsmMap[prjName] = e.PSM
	}

	logrus.Info("Processing combined file")
	evidences, err := a.processCombinedFile(a.Comb, a.Tag, a.PepProb, a.ProtProb, database)
	if err != nil {
		return err
	}

	// build map list with all centroids and quantifications
	// one report to rule them all
	// Assuming that the same database was used for everyone
	saveCompareResults(a.Temp, evidences, totalIons, uniqueIons, razorIons, PsmMap, names)

	return nil
}

// processCombinedFile ...
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

		//	p.Filter(0.01, pepProb, protProb, false, false, false)
		proid, err := fil.ProtXMLFilter(protxml, 0.01, pepProb, protProb, false, a.Razor)
		if err != nil {
			return nil, err
		}

		for _, j := range proid {

			var ce rep.CombinedEvidence
			ce.TotalPeptideIonStrings = make(map[string]int)
			ce.UniquePeptideIonStrings = make(map[string]int)
			ce.RazorPeptideIonStrings = make(map[string]int)

			ce.ProteinName = j.ProteinName
			ce.Length = j.Length
			ce.GroupNumber = j.GroupNumber
			ce.SiblingID = j.GroupSiblingID
			ce.IndiProtein = j.IndistinguishableProtein

			for _, k := range j.PeptideIons {

				var ion string
				if len(k.ModifiedPeptide) > 0 {
					ion = fmt.Sprintf("%s#%d", k.ModifiedPeptide, k.Charge)
				} else {
					ion = fmt.Sprintf("%s#%d", k.PeptideSequence, k.Charge)
				}

				ce.TotalPeptideIonStrings[ion] = 0

				if k.IsUnique == true {
					ce.UniquePeptideIonStrings[ion] = 0
				}

				if k.Razor == 1 {
					ce.RazorPeptideIonStrings[ion] = 0
				}

			}

			ce.UniqueStrippedPeptides = len(j.UniqueStrippedPeptides)
			ce.TotalPeptideIons = len(j.PeptideIons)

			for _, j := range j.PeptideIons {
				if j.IsUnique == false {
					ce.SharedPeptideIons++
				} else {
					ce.UniquePeptideIons++
				}
				if j.Razor == 1 {
					ce.RazorPeptideIons++
				}
			}

			ce.ProteinProbability = j.Probability
			ce.TopPepProb = j.TopPepProb

			list = append(list, ce)
		}
		//}

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

// saveCompareResults creates a single report using 1 or more philosopher result files
func saveCompareResults(session string, evidences rep.CombinedEvidenceList, totalIons, uniqueIons, razorIons map[string][]string, psmMap map[string]rep.PSMEvidenceList, namesList []string) {

	// create result file
	output := fmt.Sprintf("%s%scombined.tsv", session, string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create report file:", err)
	}
	defer file.Close()

	line := "Protein Group\tSubGroup\tProtein ID\tEntry Name\tGene Names\tDescription\tProtein Length\tProtein Probability\tTop Peptide Probability\tUnique Stripped Peptides\tTotal Peptide Ions\tUnique Peptide Ions\tRazor Peptide Ions\tIndistinguishable Proteins\t"
	for _, i := range namesList {
		line += fmt.Sprintf("%s Total Spectral Count\t", i)
		line += fmt.Sprintf("%s Unique Spectral Count\t", i)
		line += fmt.Sprintf("%s Razor Spectral Count\t", i)
		line += fmt.Sprintf("%s Total Intensity\t", i)
		line += fmt.Sprintf("%s Unique Intensity\t", i)
		line += fmt.Sprintf("%s Razor Intensity\t", i)
	}

	line += "\n"
	n, err := io.WriteString(file, line)
	if err != nil {
		logrus.Fatal(n, err)
	}

	for _, i := range evidences {

		var line string

		line += fmt.Sprintf("%d\t", i.GroupNumber)

		line += fmt.Sprintf("%s\t", i.SiblingID)

		line += fmt.Sprintf("%s\t", i.ProteinID)

		line += fmt.Sprintf("%s\t", i.EntryName)

		line += fmt.Sprintf("%s\t", i.GeneNames)

		line += fmt.Sprintf("%s\t", i.ProteinName)

		line += fmt.Sprintf("%d\t", i.Length)

		line += fmt.Sprintf("%.4f\t", i.ProteinProbability)

		line += fmt.Sprintf("%.4f\t", i.TopPepProb)

		line += fmt.Sprintf("%d\t", i.UniqueStrippedPeptides)

		line += fmt.Sprintf("%d\t", i.TotalPeptideIons)

		line += fmt.Sprintf("%d\t", i.UniquePeptideIons)

		line += fmt.Sprintf("%d\t", i.RazorPeptideIons)

		ip := strings.Join(i.IndiProtein, ", ")
		line += fmt.Sprintf("%s\t", ip)

		for _, j := range namesList {
			totalSpC, uniqueSpC, razorSpC := getSpectralCounts(totalIons, uniqueIons, razorIons, i.ProteinName, j)
			totalInt, uniqueInt, razorInt := sumIntensities(totalIons, uniqueIons, razorIons, psmMap[j], i.ProteinName, j)
			line += fmt.Sprintf("%d\t%d\t%d\t%6.f\t%6.f\t%6.f\t", totalSpC, uniqueSpC, razorSpC, totalInt, uniqueInt, razorInt)
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

// // saveCompareResults creates a single report using 1 or more philosopher result files
// func saveCompareResults(session string, evidences rep.CombinedEvidenceList, globalPepMap map[string]int, psmMap map[string]rep.PSMEvidenceList, namesList []string) {
//
// 	// create result file
// 	output := fmt.Sprintf("%s%scombined.tsv", session, string(filepath.Separator))
//
// 	// create result file
// 	file, err := os.Create(output)
// 	if err != nil {
// 		logrus.Fatal("Cannot create report file:", err)
// 	}
// 	defer file.Close()
//
// 	line := "Protein Group\tSubGroup\tProtein ID\tEntry Name\tGene Names\tDescription\tProtein Length\tProtein Probability\tTop Peptide Probability\tUnique Stripped Peptides\tTotal Peptide Ions\tUnique Peptide Ions\tRazor Peptide Ions\tIndistinguishable Proteins\t"
// 	for _, i := range namesList {
// 		line += fmt.Sprintf("%s Total Spectral Count\t", i)
// 		line += fmt.Sprintf("%s Unique Spectral Count\t", i)
// 		line += fmt.Sprintf("%s Razor Spectral Count\t", i)
// 		line += fmt.Sprintf("%s Total Intensity\t", i)
// 		line += fmt.Sprintf("%s Unique Intensity\t", i)
// 		line += fmt.Sprintf("%s Razor Intensity\t", i)
// 	}
//
// 	line += "\n"
// 	n, err := io.WriteString(file, line)
// 	if err != nil {
// 		logrus.Fatal(n, err)
// 	}
//
// 	for _, i := range evidences {
//
// 		var line string
//
// 		line += fmt.Sprintf("%d\t", i.GroupNumber)
//
// 		line += fmt.Sprintf("%s\t", i.SiblingID)
//
// 		line += fmt.Sprintf("%s\t", i.ProteinID)
//
// 		line += fmt.Sprintf("%s\t", i.EntryName)
//
// 		line += fmt.Sprintf("%s\t", i.GeneNames)
//
// 		line += fmt.Sprintf("%s\t", i.ProteinName)
//
// 		line += fmt.Sprintf("%d\t", i.Length)
//
// 		line += fmt.Sprintf("%.4f\t", i.ProteinProbability)
//
// 		line += fmt.Sprintf("%.4f\t", i.TopPepProb)
//
// 		line += fmt.Sprintf("%d\t", i.UniqueStrippedPeptides)
//
// 		line += fmt.Sprintf("%d\t", i.TotalPeptideIons)
//
// 		line += fmt.Sprintf("%d\t", i.UniquePeptideIons)
//
// 		line += fmt.Sprintf("%d\t", i.RazorPeptideIons)
//
// 		ip := strings.Join(i.IndiProtein, ", ")
// 		line += fmt.Sprintf("%s\t", ip)
//
// 		for _, j := range namesList {
//
// 			totalSpC, uniqueSpC, razorSpC := getSpectralCounts(globalPepMap, i.ProteinName, j)
// 			//totalInt, uniqueInt, razorInt := sumIntensities(tIons, uIons, rIons, psmMap[j], j)
// 			//line += fmt.Sprintf("%d\t%d\t%d\t%6.f\t%6.f\t%6.f\t", totalSpC, uniqueSpC, razorSpC, totalInt, uniqueInt, razorInt)
//
// 			if j == "pcdna-E1" && i.EntryName == "CIAO1_HUMAN" {
// 				fmt.Println(totalSpC, uniqueSpC, razorSpC)
// 			}
//
// 		}
//
// 		line += "\n"
// 		n, err := io.WriteString(file, line)
// 		if err != nil {
// 			logrus.Fatal(n, err)
// 		}
//
// 	}
//
// 	// copy to work directory
// 	sys.CopyFile(output, filepath.Base(output))
//
// 	return
// }

func getSpectralCounts(totalIons, uniqueIons, razorIons map[string][]string, proteinName, name string) (int, int, int) {

	var totalSpc int
	var uniqueSpc int
	var razorSpc int

	key := fmt.Sprintf("%s#%s", name, proteinName)

	Tv, Tok := totalIons[key]
	if Tok {
		totalSpc = len(Tv)
	}

	Uv, Uok := uniqueIons[key]
	if Uok {
		uniqueSpc = len(Uv)
	}

	Rv, Rok := razorIons[key]
	if Rok {
		razorSpc = len(Rv)
	}

	return totalSpc, uniqueSpc, razorSpc
}

// sumIntensities calculates the protein intensity
func sumIntensities(totalIons, uniqueIons, razorIons map[string][]string, pep rep.PSMEvidenceList, proteinName, name string) (float64, float64, float64) {

	var totalInt []float64
	var uniqueInt []float64
	var razorInt []float64

	var totalQuantInt float64
	var uniqueQuantInt float64
	var razorQuantInt float64

	var totalIonsList = make(map[string]uint8)
	var uniqueIonsList = make(map[string]uint8)
	var razorIonsList = make(map[string]uint8)

	key := fmt.Sprintf("%s#%s", name, proteinName)

	Tv, Tok := totalIons[key]
	if Tok {
		for _, i := range Tv {
			totalIonsList[i] = 0
		}
	}

	Uv, Uok := uniqueIons[key]
	if Uok {
		for _, i := range Uv {
			uniqueIonsList[i] = 0
		}
	}

	Rv, Rok := razorIons[key]
	if Rok {
		for _, i := range Rv {
			razorIonsList[i] = 0
		}
	}

	for _, i := range pep {

		var ion string
		if len(i.ModifiedPeptide) > 0 {
			ion = fmt.Sprintf("%s#%d", i.ModifiedPeptide, i.AssumedCharge)
		} else {
			ion = fmt.Sprintf("%s#%d", i.Peptide, i.AssumedCharge)
		}

		_, okT := totalIonsList[ion]
		if okT {
			totalInt = append(totalInt, i.Intensity)
		}

		_, okU := uniqueIonsList[ion]
		if okU {
			uniqueInt = append(uniqueInt, i.Intensity)
		}

		_, okR := razorIonsList[ion]
		if okR {
			razorInt = append(razorInt, i.Intensity)
		}

		sort.Float64s(uniqueInt)
		sort.Float64s(totalInt)
		sort.Float64s(razorInt)

		if len(razorInt) >= 3 {
			razorQuantInt = (razorInt[len(razorInt)-1] + razorInt[len(razorInt)-2] + razorInt[len(razorInt)-3])
		} else if len(razorInt) >= 2 {
			razorQuantInt = (razorInt[len(razorInt)-1] + razorInt[len(razorInt)-2])
		} else if len(razorInt) == 1 {
			razorQuantInt = (razorInt[len(razorInt)-1])
		}

		if len(uniqueInt) >= 3 {
			uniqueQuantInt = (uniqueInt[len(uniqueInt)-1] + uniqueInt[len(uniqueInt)-2] + uniqueInt[len(uniqueInt)-3])
		} else if len(uniqueInt) >= 2 {
			uniqueQuantInt = (uniqueInt[len(uniqueInt)-1] + uniqueInt[len(uniqueInt)-2])
		} else if len(uniqueInt) == 1 {
			uniqueQuantInt = (uniqueInt[len(uniqueInt)-1])
		}

		if len(totalInt) >= 3 {
			totalQuantInt = (totalInt[len(totalInt)-1] + totalInt[len(totalInt)-2] + totalInt[len(totalInt)-3])
		} else if len(totalInt) >= 2 {
			totalQuantInt = (totalInt[len(totalInt)-1] + totalInt[len(totalInt)-2])
		} else if len(totalInt) == 1 {
			totalQuantInt = (totalInt[len(totalInt)-1])
		}

	}

	return totalQuantInt, uniqueQuantInt, razorQuantInt
}
