package inf

import (
	"fmt"
	"strings"

	"github.com/nesvilab/philosopher/lib/rep"
)

// Peptide ...
type Peptide struct {
	IonForm                  string
	Sequence                 map[string]uint16
	Protein                  string
	Charge                   uint8
	CalcNeutralPepMass       float64
	Probability              float64
	Weight                   float64
	Spectra                  map[string]uint16
	MappedProteins           map[string]uint16
	MappedProteinsWithDecoys map[string]uint16
}

// ProteinInference ...
func ProteinInference(psm rep.PSMEvidenceList) {

	var peptideList []Peptide
	var exclusionList = make(map[string]uint16)
	var peptideIndex = make(map[string]Peptide)
	var peptideCount = make(map[string]uint16)

	var proteinTNP = make(map[string]uint16)
	var probMap = make(map[string]map[string]float64)

	// build the peptide index
	for _, i := range psm {

		ionForm := fmt.Sprintf("%s#%d#%.4f", i.Peptide, i.AssumedCharge, i.CalcNeutralPepMass)

		_, ok := exclusionList[ionForm]
		if !ok {
			var p Peptide

			p.IonForm = ionForm
			p.Sequence = make(map[string]uint16)
			p.Spectra = make(map[string]uint16)
			p.MappedProteins = make(map[string]uint16)
			p.MappedProteinsWithDecoys = make(map[string]uint16)

			p.Sequence[i.Peptide]++
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

			obj.Sequence[i.Peptide] = peptideCount[i.Peptide]
			obj.Spectra[i.Spectrum]++

			obj.MappedProteins[i.Protein]++
			obj.MappedProteinsWithDecoys[i.Protein]++

			for j := range i.MappedProteins {
				if !strings.Contains(j, "rev_") && i.Protein != j {
					obj.MappedProteins[j]++
				}
				obj.MappedProteinsWithDecoys[j]++
			}

			// total number of peptides per protein
			proteinTNP[i.Protein] += obj.Spectra[i.Spectrum]

			peptideIndex[i.Peptide] = obj
		}
	}

	// update weight
	for i := range peptideList {
		if len(peptideList[i].MappedProteins) > 0 {
			peptideList[i].Weight = (float64(1.0) / float64(len(peptideList[i].MappedProteins)))
		}

	}

	//spew.Dump(probMap)
	for i := range probMap {
		for j := range i {
			fmt.Println(i, "\t", j)
		}
	}

	return
}
