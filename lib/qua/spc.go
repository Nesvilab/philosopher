package qua

import (
	"philosopher/lib/rep"
)

// CalculateSpectralCounts add Spc to ions and proteins
func CalculateSpectralCounts(e rep.Evidence) rep.Evidence {

	var totalIonPSM = make(map[string][]string)
	var uniqueIonPSM = make(map[string][]string)
	var razorIonPSM = make(map[string][]string)

	var sequences = make(map[string]int)

	for _, i := range e.PSM {

		sequences[i.Peptide]++

		totalIonPSM[i.ProteinID] = append(totalIonPSM[i.ProteinID], i.Spectrum)

		if i.IsUnique {
			uniqueIonPSM[i.ProteinID] = append(uniqueIonPSM[i.ProteinID], i.Spectrum)
		}
		if i.IsURazor {
			razorIonPSM[i.ProteinID] = append(razorIonPSM[i.ProteinID], i.Spectrum)
		}
	}

	for i := range e.Peptides {
		v, ok := sequences[e.Peptides[i].Sequence]
		if ok {
			e.Peptides[i].Spc += v
		}
	}

	for i := range e.Proteins {

		vT, okT := totalIonPSM[e.Proteins[i].ProteinID]
		if okT {
			e.Proteins[i].TotalSpC += len(vT)
		}

		vU, okU := uniqueIonPSM[e.Proteins[i].ProteinID]
		if okU {
			e.Proteins[i].UniqueSpC += len(vU)
		}

		vUR, okR := razorIonPSM[e.Proteins[i].ProteinID]
		if okR {
			e.Proteins[i].URazorSpC += len(vUR)
		}

	}

	return e
}

// func CalculateSpectralCounts(e rep.Evidence) rep.Evidence {

// 	var uniqueIonPSM = make(map[string]string)
// 	var razorIonPSM = make(map[string]string)

// 	var sequences = make(map[string]int)

// 	for _, i := range e.PSM {

// 		sequences[i.Peptide]++

// 		if i.IsUnique {
// 			uniqueIonPSM[i.Spectrum] = i.ProteinID
// 		}
// 		if i.IsURazor {
// 			razorIonPSM[i.Spectrum] = i.ProteinID
// 		}
// 	}

// 	for i := range e.Peptides {
// 		v, ok := sequences[e.Peptides[i].Sequence]
// 		if ok {
// 			e.Peptides[i].Spc += v
// 		}
// 	}

// 	for i := range e.Proteins {

// 		for _, j := range e.Proteins[i].TotalPeptideIons {

// 			e.Proteins[i].TotalSpC += len(j.Spectra)

// 			if j.IsUnique {
// 				e.Proteins[i].UniqueSpC += len(j.Spectra)
// 			}

// 			if j.IsURazor {
// 				e.Proteins[i].URazorSpC += len(j.Spectra)
// 			}

// 		}

// 	}

// 	return e
// }
