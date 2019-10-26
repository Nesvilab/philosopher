package inf

import (
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/nesvilab/philosopher/lib/id"
)

// Peptide ...
type Peptide struct {
	Sequence            map[string]uint8
	Charge              uint8
	CalcNeutralPepMass  float64
	Probability         float64
	Spectra             map[string]uint8
	Protein             map[string]uint8
	AlternativeProteins map[string]uint8
	Weight              float64
}

// ProteinInference ...
func ProteinInference(psm id.PepIDList) {

	var peptideList []Peptide
	var exclusionList = make(map[string]uint8)
	var peptideIndex = make(map[string]Peptide)
	var peptideCount = make(map[string]uint8)

	// build the peptide index
	for _, i := range psm {

		ionForm := fmt.Sprintf("%s#%d#%.4f", i.Peptide, i.AssumedCharge, i.CalcNeutralPepMass)

		_, ok := exclusionList[ionForm]
		if !ok {
			var p Peptide

			p.Sequence = make(map[string]uint8)
			p.Spectra = make(map[string]uint8)
			p.Protein = make(map[string]uint8)
			p.AlternativeProteins = make(map[string]uint8)

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

			for _, j := range i.AlternativeProteins {
				if !strings.Contains(j, "rev_") && i.Protein != j {
					obj.AlternativeProteins[j]++
				}
			}
			//obj.AlternativeProteins[i.Protein]++

			peptideIndex[i.Peptide] = obj
		}
	}

	// update weight
	for i := range peptideList {
		if len(peptideList[i].AlternativeProteins) > 0 {
			peptideList[i].Weight = (float64(1.0) / float64(len(peptideList[i].AlternativeProteins)))
		}
	}

	spew.Dump(peptideList)

	return
}
