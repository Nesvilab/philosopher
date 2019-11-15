package qua

import (
	"errors"

	"philosopher/lib/msg"
	"philosopher/lib/rep"
)

// CalculateSpectralCounts add Spc to ions and proteins
func CalculateSpectralCounts(e rep.Evidence) rep.Evidence {

	if len(e.PSM) < 1 && len(e.Ions) < 1 {
		msg.QuantifyingData(errors.New("The PSM list is enpty"), "fatal")
	}

	var uniqueIonPSM = make(map[string]string)
	var razorIonPSM = make(map[string]string)
	for _, i := range e.PSM {
		if i.IsUnique {
			uniqueIonPSM[i.Spectrum] = i.ProteinID
		}
		if i.IsURazor {
			razorIonPSM[i.Spectrum] = i.ProteinID
		}
	}

	for i := range e.Proteins {

		e.Proteins[i].TotalSpC = len(e.Proteins[i].SupportingSpectra)

		for _, j := range e.Proteins[i].TotalPeptideIons {

			if j.IsUnique == true {
				e.Proteins[i].UniqueSpC += len(j.Spectra)
			}

			if j.IsURazor == true {
				e.Proteins[i].URazorSpC += len(j.Spectra)
			}

		}

	}

	return e
}
