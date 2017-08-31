package quan

import (
	"fmt"

	"github.com/prvst/cmsl/err"
	"github.com/prvst/philosopher/lib/rep"
)

// CalculateSpectralCounts add Spc to ions and proteins
func CalculateSpectralCounts(e rep.Evidence) (rep.Evidence, *err.Error) {

	var spcMap = make(map[string]int)
	var ionRefMap = make(map[string]int)

	if len(e.PSM) < 1 && len(e.Ions) < 1 {
		return e, &err.Error{Type: err.WorkspaceNotFound, Class: err.FATA}
	}

	for _, i := range e.PSM {
		var key string
		if len(i.ModifiedPeptide) > 0 {
			key = fmt.Sprintf("%s#%d", i.ModifiedPeptide, i.AssumedCharge)
		} else {
			key = fmt.Sprintf("%s#%d", i.Peptide, i.AssumedCharge)
		}
		spcMap[key]++
	}

	for i := range e.Ions {
		var key string
		if len(e.Ions[i].ModifiedSequence) > 0 {
			key = fmt.Sprintf("%s#%d", e.Ions[i].ModifiedSequence, e.Ions[i].ChargeState)
		} else {
			key = fmt.Sprintf("%s#%d", e.Ions[i].Sequence, e.Ions[i].ChargeState)
		}
		v1, ok := spcMap[key]
		if ok {
			e.Ions[i].Spc = v1
			ionRefMap[key] = v1
		}
	}

	for i := range e.Proteins {

		// make a reference for razor peptides
		var uniqIons = make(map[string]uint8)

		for k := range e.Proteins[i].URazorPeptideIons {
			v, ok := ionRefMap[k]
			if ok {
				e.Proteins[i].UniqueSpC += v
				e.Proteins[i].RazorSpC += v
				uniqIons[k] = 0
			}
		}

		// also checks if peptides are razor, then add them to the uniquerazor field
		for k, j := range e.Proteins[i].TotalPeptideIons {
			v, ok := ionRefMap[k]
			if ok {
				e.Proteins[i].TotalSpC += v
				if j.IsRazor {
					_, ok := uniqIons[k]
					if !ok {
						e.Proteins[i].RazorSpC += v
					}
				}
			}
		}

	}

	return e, nil
}
