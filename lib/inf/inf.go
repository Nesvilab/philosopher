package inf

import (
	"fmt"
	"strings"

	"github.com/nesvilab/philosopher/lib/id"
)

// Peptide ...
type Peptide struct {
	Sequence                      map[string]uint16
	Charge                        uint8
	CalcNeutralPepMass            float64
	Probability                   float64
	Weight                        float64
	Spectra                       map[string]uint16
	Protein                       map[string]uint16
	AlternativeProteins           map[string]uint16
	AlternativeProteinsWithDecoys map[string]uint16
}

// ProteinInference ...
func ProteinInference(psm id.PepIDList) {

	var peptideList []Peptide
	var exclusionList = make(map[string]uint16)
	var peptideIndex = make(map[string]Peptide)
	var peptideCount = make(map[string]uint16)
	var proteinTNP = make(map[string]uint16)

	// build the peptide index
	for _, i := range psm {

		ionForm := fmt.Sprintf("%s#%d#%.4f", i.Peptide, i.AssumedCharge, i.CalcNeutralPepMass)

		_, ok := exclusionList[ionForm]
		if !ok {
			var p Peptide

			p.Sequence = make(map[string]uint16)
			p.Spectra = make(map[string]uint16)
			p.Protein = make(map[string]uint16)
			p.AlternativeProteins = make(map[string]uint16)
			p.AlternativeProteinsWithDecoys = make(map[string]uint16)

			p.Sequence[i.Peptide]++
			p.Charge = i.AssumedCharge
			p.CalcNeutralPepMass = i.CalcNeutralPepMass
			p.Probability = i.Probability
			p.Weight = 1.0

			exclusionList[ionForm] = 0
			peptideCount[i.Peptide]++

			peptideList = append(peptideList, p)
			peptideIndex[ionForm] = p
		}
	}

	for _, i := range psm {

		ionForm := fmt.Sprintf("%s#%d#%.4f", i.Peptide, i.AssumedCharge, i.CalcNeutralPepMass)

		v, ok := peptideIndex[ionForm]
		if ok {
			obj := v

			obj.Sequence[i.Peptide] = peptideCount[i.Peptide]
			obj.Spectra[i.Spectrum]++
			obj.Protein[i.Protein]++

			proteinTNP[i.Protein] += obj.Spectra[i.Spectrum]

			for _, j := range i.AlternativeProteins {
				if !strings.Contains(j, "rev_") && i.Protein != j {
					obj.AlternativeProteins[j]++
				}
				obj.AlternativeProteinsWithDecoys[j]++
			}

			peptideIndex[i.Peptide] = obj
		}
	}

	// update weight
	for i := range peptideList {
		if len(peptideList[i].AlternativeProteins) > 0 {
			peptideList[i].Weight = (float64(1.0) / float64(len(peptideList[i].AlternativeProteins)))
		}
	}

	//spew.Dump(proteinTNP)
	// for _, i := range peptideList {
	// 	//if i.Weight >= 0.9 {
	// 	fmt.Println(i.Sequence, "\t", i.Protein)
	// 	//}
	// }

	return
}
