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
	"github.com/prvst/philosopher-source/lib/data"
	"github.com/prvst/philosopher-source/lib/ext/proteinprophet"
	"github.com/prvst/philosopher-source/lib/fil"
	"github.com/prvst/philosopher-source/lib/meta"
	"github.com/prvst/philosopher-source/lib/rep"
	"github.com/prvst/philosopher-source/lib/sys"
	"github.com/prvst/philosopher-source/lib/xml"
)

// Abacus main structure
type Abacus struct {
	meta.Data
	ProtProb float64
	PepProb  float64
	Comb     string
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

	var xmlFiles []string
	var database data.Base
	var globalPepMap = make(map[string]int)
	var PsmMap = make(map[string]rep.PSMEvidenceList)
	var names []string

	// initialize proteinprophet
	var e rep.Evidence
	e.RestoreWithPath(args[0])
	pop := initializeProteinProphet(e)

	// recover all files
	logrus.Info("Restoring Philospher results")

	for _, i := range args {

		var e rep.Evidence
		e.RestoreWithPath(i)

		files, _ := ioutil.ReadDir(i)
		for _, f := range files {
			if strings.Contains(f.Name(), "pep.xml") {
				absPath, _ := filepath.Abs(f.Name())
				xmlFiles = append(xmlFiles, absPath)
			}
		}

		// if len(e.Proteins) < 1 {
		// 	logrus.Fatal("result files does not contains protein inference information. Run the filter option with a protXML file in order to have the results.")
		// }

		//xmlFiles = append(xmlFiles, r.PepXML.Files...)
		database = data.Base{}
		database.RestoreWithPath(i)
		names = append(names, e.ProjectName)

		for _, j := range e.PSM {
			var ion string
			if j.Probability >= a.PepProb {
				if len(j.ModifiedPeptide) > 0 {
					ion = fmt.Sprintf("%s#%d", j.ModifiedPeptide, j.AssumedCharge)
				} else {
					ion = fmt.Sprintf("%s#%d", j.Peptide, j.AssumedCharge)
				}
				key := fmt.Sprintf("%s@%s", e.ProjectName, ion)
				globalPepMap[key]++
			}
		}
		PsmMap[e.ProjectName] = e.PSM
	}

	sort.Strings(names)

	var combinedFile string

	if len(a.Comb) < 1 {
		logrus.Info("Creating the combined protXML file")

		// set the output name
		//pop.Combine = true
		// deploy the binaries
		pop.Deploy()

		// run ProteinProphet
		err := pop.Run(xmlFiles)
		if err != nil {
			return err
		}

		combinedFile = "combined.prot.xml"
	} else {
		combinedFile = a.Comb
	}

	var m meta.Data
	metaPath := fmt.Sprintf("%s%s%s", args[0], string(filepath.Separator), sys.Meta())
	m.Restore(metaPath)

	var decoyTag = m.DecoyTag
	//var conTag = data.ConTag()

	evidences, err := a.processCombinedFile(combinedFile, decoyTag, a.PepProb, a.ProtProb, database)
	if err != nil {
		return err
	}

	// build map list with all centroids and quantifications
	// one report to rule them all
	// Assuming that the same database was used for everyone
	saveCompareResults(m.Temp, evidences, globalPepMap, PsmMap, names)

	return nil
}

// processCombinedFile ...
func (a *Abacus) processCombinedFile(combinedFile, decoyTag string, pepProb, protProb float64, database data.Base) (rep.CombinedEvidenceList, error) {

	var list rep.CombinedEvidenceList

	if _, err := os.Stat(combinedFile); os.IsNotExist(err) {
		logrus.Fatal("Cannot find combined.prot.xml file")
	} else {

		var protxml xml.ProtXML
		// protxml.Read(combinedFile)

		protxml.Read(combinedFile)
		protxml.DecoyTag = decoyTag
		//protxml.ConTag = data.ConTag()

		// promote decoy proteins with indistinguishable target proteins
		protxml.PromoteProteinIDs()

		//	p.Filter(0.01, pepProb, protProb, false, false, false)
		proid, err := fil.ProtXMLFilter(protxml, 0.01, pepProb, protProb, false, false)
		if err != nil {
			return nil, err
		}

		for _, j := range proid {

			var ce rep.CombinedEvidence
			ce.TotalPeptideIonStrings = make(map[string]int)
			ce.UniquePeptideIonStrings = make(map[string]int)

			ce.ProteinName = j.ProteinName
			ce.Length = j.Length
			ce.GroupNumber = j.GroupNumber
			ce.SiblingID = j.GroupSiblingID

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

			}

			ce.UniqueStrippedPeptides = len(j.UniqueStrippedPeptides)
			ce.TotalPeptideIons = len(j.PeptideIons)

			for _, j := range j.PeptideIons {
				if j.IsUnique == false {
					ce.SharedPeptideIons++
				} else {
					ce.UniquePeptideIons++
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
func saveCompareResults(session string, evidences rep.CombinedEvidenceList, globalPepMap map[string]int, psmMap map[string]rep.PSMEvidenceList, namesList []string) {

	// create result file
	output := fmt.Sprintf("%s%scombined.tsv", session, string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create report file:", err)
	}
	defer file.Close()

	line := "Protein Group\tProtein ID\tEntry Name\tGene Names\tDescription\tProtein Length\tProtein Probability\tTop Peptide Probability\tUnique Stripped Peptides\tTotal Peptide Ions\tUnique Peptide Ions\t"
	for _, i := range namesList {
		line += fmt.Sprintf("%s Total Spectral Count\t", i)
		line += fmt.Sprintf("%s Unique Spectral Count\t", i)
		line += fmt.Sprintf("%s Total Intensity\t", i)
		line += fmt.Sprintf("%s Unique Intensity\t", i)
	}

	line += "\n"
	n, err := io.WriteString(file, line)
	if err != nil {
		logrus.Fatal(n, err)
	}

	for _, i := range evidences {

		var line string

		line += fmt.Sprintf("%d-%s\t", i.GroupNumber, i.SiblingID)

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

		var tIons []string
		for j := range i.TotalPeptideIonStrings {
			tIons = append(tIons, j)
		}

		var uIons []string
		for j := range i.UniquePeptideIonStrings {
			uIons = append(uIons, j)
		}

		for _, j := range namesList {
			totalSpC, uniqueSpC := getSpectralCounts(tIons, uIons, globalPepMap, j)
			totalInt, uniqueInt := sumIntensities(tIons, uIons, psmMap[j], j)
			line += fmt.Sprintf("%d\t%d\t%6.f\t%6.f\t", totalSpC, uniqueSpC, totalInt, uniqueInt)
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

func getSpectralCounts(tIons, uIons []string, globalPepMap map[string]int, name string) (int, int) {

	var totalSpc int
	var uniqueSpc int

	for _, i := range tIons {
		key := fmt.Sprintf("%s@%s", name, i)
		v, okT := globalPepMap[key]
		if okT {
			totalSpc += v
		}
	}

	for _, i := range uIons {
		key := fmt.Sprintf("%s@%s", name, i)
		v, okU := globalPepMap[key]
		if okU {
			uniqueSpc += v
		}
	}

	// for _, i := range tIons {
	// 	key := fmt.Sprintf("%s@%s", name, i)
	// 	_, okT := globalPepMap[key]
	// 	if okT {
	// 		totalSpc++
	// 	}
	// }
	//
	// for _, i := range uIons {
	// 	key := fmt.Sprintf("%s@%s", name, i)
	// 	_, okU := globalPepMap[key]
	// 	if okU {
	// 		uniqueSpc++
	// 	}
	// }

	return totalSpc, uniqueSpc
}

// sumIntensities calculates the protein intensity
func sumIntensities(tIons, uIons []string, pep rep.PSMEvidenceList, name string) (float64, float64) {

	var totalInt []float64
	var uniqueInt []float64

	var totalQuantInt float64
	var uniqueQuantInt float64

	var totalMap = make(map[string]int)
	for _, i := range tIons {
		totalMap[i]++
	}

	var uniqueMap = make(map[string]int)
	for _, i := range uIons {
		uniqueMap[i]++
	}

	for _, i := range pep {

		var ion string
		if len(i.ModifiedPeptide) > 0 {
			ion = fmt.Sprintf("%s#%d", i.ModifiedPeptide, i.AssumedCharge)
		} else {
			ion = fmt.Sprintf("%s#%d", i.Peptide, i.AssumedCharge)
		}

		_, okT := totalMap[ion]
		if okT {
			totalInt = append(totalInt, i.Intensity)
		}

		_, okU := uniqueMap[ion]
		if okU {
			uniqueInt = append(uniqueInt, i.Intensity)
		}

		sort.Float64s(uniqueInt)
		sort.Float64s(totalInt)

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

	return totalQuantInt, uniqueQuantInt
}

func initializeProteinProphet(e rep.Evidence) proteinprophet.ProteinProphet {

	var pop proteinprophet.ProteinProphet

	pop.UUID = e.UUID
	pop.Distro = e.Distro
	pop.Home = e.Home
	pop.MetaFile = e.MetaFile
	pop.MetaDir = e.MetaDir
	pop.DB = e.DB
	pop.Temp = e.Temp
	pop.TimeStamp = e.TimeStamp
	pop.OS = e.OS
	pop.Arch = e.Arch

	pop.UnixBatchCoverage = pop.Temp + string(filepath.Separator) + "batchcoverage"
	pop.UnixDatabaseParser = pop.Temp + string(filepath.Separator) + "DatabaseParser"
	pop.UnixProteinProphet = pop.Temp + string(filepath.Separator) + "ProteinProphet"
	pop.WinBatchCoverage = pop.Temp + string(filepath.Separator) + "batchcoverage.exe"
	pop.WinDatabaseParser = pop.Temp + string(filepath.Separator) + "DatabaseParser.exe"
	pop.WinProteinProphet = pop.Temp + string(filepath.Separator) + "ProteinProphet.exe"
	pop.LibgccDLL = pop.Temp + string(filepath.Separator) + "libgcc_s_dw2-1.dll"
	pop.Zlib1DLL = pop.Temp + string(filepath.Separator) + "zlib1.dll"

	return pop
}
