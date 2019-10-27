package inf

import (
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/nesvilab/philosopher/lib/rep"
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
	Spectra                  map[string]int
	MappedProteins           map[string]int
	MappedProteinsWithDecoys map[string]int
}

// ProteinInference ...
func ProteinInference(psm rep.PSMEvidenceList) {

	var peptideList []Peptide
	var exclusionList = make(map[string]int)
	var peptideIndex = make(map[string]Peptide)
	var peptideCount = make(map[string]int)

	var proteinTNP = make(map[string]int)
	var probMap = make(map[string]map[string]float64)

	// build the peptide index
	for _, i := range psm {

		ionForm := fmt.Sprintf("%s#%d#%.4f", i.Peptide, i.AssumedCharge, i.CalcNeutralPepMass)

		_, ok := exclusionList[ionForm]
		if !ok {
			var p Peptide

			p.IonForm = ionForm
			p.Spectra = make(map[string]int)
			p.MappedProteins = make(map[string]int)
			p.MappedProteinsWithDecoys = make(map[string]int)

			p.Sequence = i.Peptide
			p.Charge = i.AssumedCharge
			p.CalcNeutralPepMass = i.CalcNeutralPepMass
			p.Probability = -1.0
			p.Weight = 1.0

			exclusionList[ionForm] = 0
			peptideCount[i.Peptide]++

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
				probMap[i.Peptide] = map[string]float64{i.Protein: i.Probability}
			}

		}
	}

	for _, i := range psm {

		ionForm := fmt.Sprintf("%s#%d#%.4f", i.Peptide, i.AssumedCharge, i.CalcNeutralPepMass)

		v, ok := peptideIndex[ionForm]
		if ok {
			obj := v

			obj.Sequence = i.Peptide
			obj.Spectra[i.Spectrum]++

			obj.MappedProteins[i.Protein]++
			obj.MappedProteinsWithDecoys[i.Protein]++

			for j := range i.MappedProteins {
				if !strings.Contains(j, "rev_") && i.Protein != j {
					obj.MappedProteins[j]++
				}
				obj.MappedProteinsWithDecoys[j]++
			}

			// assign razor for absolute mappings
			if len(i.MappedProteins) == 1 {
				obj.Protein = i.Protein
				obj.Probability = i.Probability
			}

			// total number of peptides per protein
			proteinTNP[i.Protein] += int(obj.Spectra[i.Spectrum])

			peptideIndex[i.Peptide] = obj
		}
	}

	// update weight and assign razor
	for i := range peptideList {

		if len(peptideList[i].MappedProteins) > 0 {
			peptideList[i].Weight = (float64(1.0) / float64(len(peptideList[i].MappedProteins)))
		}

		var protein string
		var tnp int

		for k, v := range peptideList[i].MappedProteins {
			if v > tnp {
				tnp = v
				protein = k
			}
		}

		pm := probMap[peptideList[i].Sequence]

		peptideList[i].Protein = protein
		peptideList[i].Probability = pm[protein]

	}

	spew.Dump(peptideList)
	// for k, v := range probMap {
	// 	for i := range v {
	// 		fmt.Println(k, "\t", i)
	// 	}
	// }

	return
}
