package qua

import (
	"philosopher/lib/rep"
	"philosopher/lib/uti"
)

// CalculatePeptideCounts counts peptides for proteins
func CalculatePeptideCounts(e rep.Evidence) rep.Evidence {

	var total = make(map[string][]string)
	var unique = make(map[string][]string)
	var razor = make(map[string][]string)

	for _, i := range e.PSM {

		total[i.Protein] = append(total[i.Protein], i.Peptide)
		for j := range i.MappedProteins {
			total[j] = append(total[j], i.Peptide)
		}

		if i.IsUnique {
			unique[i.Protein] = append(unique[i.Protein], i.Peptide)
		}

		if i.IsURazor {
			razor[i.Protein] = append(razor[i.Protein], i.Peptide)
		}
	}

	for k, v := range total {
		total[k] = uti.RemoveDuplicateStrings(v)
	}

	for k, v := range unique {
		unique[k] = uti.RemoveDuplicateStrings(v)
	}

	for k, v := range razor {
		razor[k] = uti.RemoveDuplicateStrings(v)
	}

	for i := range e.Proteins {

		vTP, okTP := total[e.Proteins[i].PartHeader]
		if okTP {
			for _, j := range vTP {
				e.Proteins[i].TotalPeptides[j]++
			}
		}

		vuP, okuP := unique[e.Proteins[i].PartHeader]
		if okuP {
			for _, j := range vuP {
				e.Proteins[i].UniquePeptides[j]++
			}
		}

		vRP, okRP := razor[e.Proteins[i].PartHeader]
		if okRP {
			for _, j := range vRP {
				e.Proteins[i].URazorPeptides[j]++
			}
		}

	}

	return e
}
