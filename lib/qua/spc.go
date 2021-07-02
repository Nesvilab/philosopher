package qua

import (
	"philosopher/lib/rep"
)

// CalculateSpectralCounts add Spc to ions and proteins
func CalculateSpectralCounts(e rep.Evidence) rep.Evidence {

	var uniqueIonPSM = make(map[string]string)
	var razorIonPSM = make(map[string]string)

	var sequences = make(map[string]int)

	for _, i := range e.PSM {

		sequences[i.Peptide]++

		if i.IsUnique {
			uniqueIonPSM[i.Spectrum] = i.ProteinID
		}
		if i.IsURazor {
			razorIonPSM[i.Spectrum] = i.ProteinID
		}
	}

	for i := range e.Peptides {
		v, ok := sequences[e.Peptides[i].Sequence]
		if ok {
			e.Peptides[i].Spc += v
		}
	}

	for i := range e.Proteins {

		e.Proteins[i].TotalSpC = len(e.Proteins[i].SupportingSpectra)

		for _, j := range e.Proteins[i].TotalPeptideIons {

			if j.IsUnique {
				e.Proteins[i].UniqueSpC += len(j.Spectra)
			}

			if j.IsURazor {
				e.Proteins[i].URazorSpC += len(j.Spectra)
			}

		}

	}

	return e
}
