package inf

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"philosopher/lib/dat"
	"philosopher/lib/id"
	"philosopher/lib/uti"
)

// Peptide ...
type Peptide struct {
	IonForm                  string
	Sequence                 string
	Protein                  string
	Charge                   uint8
	CalcNeutralPepMass       float64
	Probability              float64
	Weight                   float64
	Spectra                  map[id.SpectrumType]int
	MappedProteins           map[string]int
	MappedProteinsWithDecoys map[string]int
}

// ProteinInference ...
func ProteinInference(psm id.PepIDList) (id.PepIDList, map[string]string, map[string]float64) {

	var peptideList []Peptide
	var exclusionList = make(map[string]int)
	var peptideIndex = make(map[string]Peptide)
	var proteinTNP = make(map[string]int)
	var probMap = make(map[string]map[string]float64)
	var proteinPepSeqMap = make(map[string][]string)

	// collect database information
	var db dat.Base
	db.Restore()

	// build the peptide index
	for _, i := range psm {

		ionForm := fmt.Sprintf("%s#%d#%.4f", i.Peptide, i.AssumedCharge, i.CalcNeutralPepMass)

		_, ok := exclusionList[ionForm]
		if !ok {
			var p Peptide

			p.IonForm = ionForm
			p.Spectra = make(map[id.SpectrumType]int)
			p.MappedProteins = make(map[string]int)
			p.MappedProteinsWithDecoys = make(map[string]int)

			p.Sequence = i.Peptide
			p.Charge = i.AssumedCharge
			p.CalcNeutralPepMass = i.CalcNeutralPepMass
			p.Probability = i.Probability
			p.Weight = -1.0

			peptideList = append(peptideList, p)
			peptideIndex[ionForm] = p

			// build the peptide to protein prob map
			v, okPM := probMap[i.Peptide]
			if okPM {
				inner := v
				if i.Probability > inner[i.Protein] {
					inner[i.Protein] = i.Probability
				}
				probMap[i.Peptide] = inner
			} else {
				var obj = make(map[string]float64)
				obj[i.Protein] = i.Probability
				probMap[i.Peptide] = obj
			}

			exclusionList[ionForm] = 0
		}

		// total number of peptides per protein
		proteinTNP[i.Protein]++

		for j := range i.AlternativeProteins {
			if j != i.Protein {
				proteinTNP[j]++
			}
		}

		proteinPepSeqMap[i.Protein] = append(proteinPepSeqMap[i.Protein], i.Peptide)
	}

	for _, i := range psm {

		ionForm := fmt.Sprintf("%s#%d#%.4f", i.Peptide, i.AssumedCharge, i.CalcNeutralPepMass)

		v, ok := peptideIndex[ionForm]
		if ok {
			obj := v

			obj.Sequence = i.Peptide
			obj.Spectra[i.SpectrumFileName()]++

			obj.MappedProteins[i.Protein] = proteinTNP[i.Protein]
			obj.MappedProteinsWithDecoys[i.Protein] = proteinTNP[i.Protein]

			for j := range i.AlternativeProteins {
				obj.MappedProteins[j] = -1
				obj.MappedProteinsWithDecoys[j] = -1
			}

			// assign razor for absolute mappings
			if len(i.AlternativeProteins) == 1 {
				obj.Protein = i.Protein
			}

			// get the highest probability
			if i.Probability > obj.Probability {
				obj.Probability = i.Probability
			}

			peptideIndex[ionForm] = obj
		}
	}

	// update Total Number of Peptides per Protein
	for i := range peptideList {

		if len(peptideList[i].MappedProteins) > 0 {
			peptideList[i].Weight = (float64(1.0) / float64(len(peptideList[i].MappedProteins)))
		}

		for k := range peptideList[i].MappedProteins {
			_, ok := proteinTNP[k]
			if ok {
				peptideList[i].MappedProteins[k] = proteinTNP[k]
			}
		}

		for k := range peptideList[i].MappedProteinsWithDecoys {
			_, ok := proteinTNP[k]
			if ok {
				peptideList[i].MappedProteinsWithDecoys[k] = proteinTNP[k]
			}
		}

	}

	proteinCoverageMap := calculateProteinCoverage(proteinPepSeqMap, db)

	// assign razor
	var razorMap = make(map[string]string)
	for i := range peptideList {

		var protein string
		var candidateProteins []string
		var tnp int
		var coverage float64

		for k := range peptideList[i].MappedProteins {
			candidateProteins = append(candidateProteins, k)
		}

		sort.Strings(candidateProteins)

		for _, j := range candidateProteins {

			if peptideList[i].MappedProteins[j] > tnp {
				tnp = peptideList[i].MappedProteins[j]
				protein = j
			}
		}

		for _, j := range candidateProteins {

			if peptideList[i].MappedProteins[j] == tnp && proteinCoverageMap[j] > coverage {
				coverage = proteinCoverageMap[j]
				protein = j
			}
		}

		if len(protein) > 0 {
			peptideList[i].Protein = protein
		}

		razorMap[peptideList[i].Sequence] = peptideList[i].Protein
	}

	//spew.Dump(peptideList)
	//spew.Dump(proteinCoverageMap)

	// update PSMs
	for i := range psm {
		pt, ok := razorMap[psm[i].Peptide]
		if ok {

			if pt != psm[i].Protein {

				psm[i].AlternativeProteins[psm[i].Protein]++

				var toRemove string
				for j := range psm[i].AlternativeProteins {
					if j == pt {
						toRemove = j
						break
					}
				}

				psm[i].AlternativeProteins[psm[i].Protein]++
				delete(psm[i].AlternativeProteins, toRemove)

				psm[i].Protein = pt
			}

		}
	}

	return psm, razorMap, proteinCoverageMap
}

// calculateProteinCoverage returns a percentage of coverage based on a set of peptides
func calculateProteinCoverage(proteinPepSeqMap map[string][]string, db dat.Base) map[string]float64 {

	var coverage = make(map[string]float64)
	var protSeq = make(map[string]string)

	for _, i := range db.Records {
		protSeq[i.PartHeader] = i.Sequence
	}

	for k, v := range proteinPepSeqMap {

		var peptides = make(map[string]uint8)
		for _, i := range v {
			peptides[i]++
		}

		CoveredSequence := protSeq[k]
		for pep := range peptides {

			re := regexp.MustCompile(pep)
			loc := re.FindStringIndex(protSeq[k])

			if len(loc) > 0 {
				CoveredSequence = strings.ReplaceAll(protSeq[k], protSeq[k][loc[0]:loc[1]], strings.Repeat("X", len(pep)))
			}
		}

		cov := regexp.MustCompile("X")
		matches := cov.FindAllStringIndex(CoveredSequence, -1)
		coverage[k] = uti.Round(float64(len(matches))/float64(len(CoveredSequence))*100, 5, 2)
	}

	return coverage
}
