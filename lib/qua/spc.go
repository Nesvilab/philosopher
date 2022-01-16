package qua

import (
	"philosopher/lib/id"
	"philosopher/lib/rep"
)

// CalculateSpectralCounts add Spc to ions and proteins
func CalculateSpectralCounts(e rep.Evidence) rep.Evidence {

	var total = make(map[string][]id.SpectrumType)
	var unique = make(map[string][]id.SpectrumType)
	var razor = make(map[string][]id.SpectrumType)

	var sequences = make(map[string]int)

	for _, i := range e.PSM {

		sequences[i.Peptide]++

		total[i.Protein] = append(total[i.Protein], i.SpectrumFileName())
		for j := range i.MappedProteins {
			total[j] = append(total[j], i.SpectrumFileName())
		}

		if i.IsUnique {
			unique[i.Protein] = append(unique[i.Protein], i.SpectrumFileName())
		}

		if i.IsURazor {
			razor[i.Protein] = append(razor[i.Protein], i.SpectrumFileName())
		}
	}

	for i := range e.Peptides {
		v, ok := sequences[e.Peptides[i].Sequence]
		if ok {
			e.Peptides[i].Spc += v
		}
	}

	for i := range e.Proteins {

		//e.Proteins[i].TotalSpC = 0
		//e.Proteins[i].UniqueSpC = 0
		//e.Proteins[i].URazorSpC = 0

		vT, okT := total[e.Proteins[i].PartHeader]
		if okT {
			e.Proteins[i].TotalSpC += len(vT)
		}

		vU, okU := unique[e.Proteins[i].PartHeader]
		if okU {
			e.Proteins[i].UniqueSpC += len(vU)
		}

		vUR, okR := razor[e.Proteins[i].PartHeader]
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
// 			uniqueIonPSM[i.SpectrumFileName()] = i.ProteinID
// 		}
// 		if i.IsURazor {
// 			razorIonPSM[i.SpectrumFileName()] = i.ProteinID
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
