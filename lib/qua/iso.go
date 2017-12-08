package qua

import (
	"fmt"
	"math"
	"strings"

	"github.com/prvst/philosopher/lib/raw"
	"github.com/prvst/philosopher/lib/rep"
	"github.com/prvst/philosopher/lib/uti"
)

const (
	mzDeltaWindow float64 = 0.5
)

// calculateIonPurity verifies how much interference there is on the precursor scans for each fragment
func calculateIonPurity(d, f string, ms1 raw.MS1, ms2 raw.MS2, evi []rep.PSMEvidence) ([]rep.PSMEvidence, error) {

	// index MS1 spectra in a dictionary
	var indexedMS1 = make(map[string]raw.Ms1Scan)
	for i := range ms1.Ms1Scan {
		// left-pad the spectrum index
		paddedIndex := fmt.Sprintf("%05s", ms1.Ms1Scan[i].Index)

		// left-pad the spectrum scan
		paddedScan := fmt.Sprintf("%05s", ms1.Ms1Scan[i].Scan)

		ms1.Ms1Scan[i].Index = paddedIndex
		ms1.Ms1Scan[i].Scan = paddedScan

		indexedMS1[paddedScan] = ms1.Ms1Scan[i]
	}

	// index MS2 spectra in a dictionary
	var indexedMS2 = make(map[string]raw.Ms2Scan)
	for i := range ms2.Ms2Scan {

		if ms2.Ms2Scan[i].Precursor.IsolationWindowLowerOffset == 0 && ms2.Ms2Scan[i].Precursor.IsolationWindowUpperOffset == 0 {
			ms2.Ms2Scan[i].Precursor.IsolationWindowLowerOffset = mzDeltaWindow
			ms2.Ms2Scan[i].Precursor.IsolationWindowUpperOffset = mzDeltaWindow
		}

		// left-pad the spectrum index
		paddedIndex := fmt.Sprintf("%05s", ms2.Ms2Scan[i].Index)

		// left-pad the spectrum scan
		paddedScan := fmt.Sprintf("%05s", ms2.Ms2Scan[i].Scan)

		// left-pad the precursor spectrum index
		paddedPI := fmt.Sprintf("%05s", ms2.Ms2Scan[i].Precursor.ParentIndex)

		// left-pad the precursor spectrum scan
		paddedPS := fmt.Sprintf("%05s", ms2.Ms2Scan[i].Precursor.ParentScan)

		ms2.Ms2Scan[i].Index = paddedIndex
		ms2.Ms2Scan[i].Scan = paddedScan
		ms2.Ms2Scan[i].Precursor.ParentIndex = paddedPI
		ms2.Ms2Scan[i].Precursor.ParentScan = paddedPS

		indexedMS2[paddedScan] = ms2.Ms2Scan[i]
	}

	for i := range evi {

		// get spectrum index
		split := strings.Split(evi[i].Spectrum, ".")

		v2, ok := indexedMS2[split[1]]
		if ok {

			v1 := indexedMS1[v2.Precursor.ParentScan]

			var ions = make(map[float64]float64)
			var isolationWindowSummedInt float64

			for k := range v1.Spectrum {
				if v1.Spectrum[k].Mz >= (v2.Precursor.TargetIon-v2.Precursor.IsolationWindowUpperOffset) && v1.Spectrum[k].Mz <= (v2.Precursor.TargetIon+v2.Precursor.IsolationWindowUpperOffset) {
					ions[v1.Spectrum[k].Mz] = v1.Spectrum[k].Intensity
					isolationWindowSummedInt += v1.Spectrum[k].Intensity
				}
			}

			if evi[i].Spectrum == "20170314_LC_TMTB4_prot_F24_01.02219.02219.2" {
				fmt.Println(ions)
				fmt.Println(isolationWindowSummedInt)
			}

			// create the list of mz differences for each peak
			var mzRatio []float64
			for k := 1; k <= 6; k++ {
				r := float64(k) * (float64(1) / float64(v2.Precursor.ChargeState))
				mzRatio = append(mzRatio, uti.ToFixed(r, 2))
			}

			var isotopePackage = make(map[float64]float64)

			isotopePackage[v2.Precursor.TargetIon] = v2.Precursor.PeakIntensity
			isotopesInt := v2.Precursor.PeakIntensity

			for k, v := range ions {
				for _, m := range mzRatio {
					if math.Abs(v2.Precursor.TargetIon-k) <= (m+0.02) && math.Abs(v2.Precursor.TargetIon-k) >= (m-0.02) {
						isotopePackage[k] = v
						isotopesInt += v
						break
					}
				}
			}

			if isotopesInt == 0 {
				evi[i].Purity = 0
			} else {
				evi[i].Purity = uti.Round((isotopesInt / isolationWindowSummedInt), 5, 2)
			}

			if evi[i].Purity > 1 {
				evi[i].Purity = 1
			}

		}

	}

	return evi, nil
}

// // prepareLabelStructure instantiates the Label objects and maps them against the fragment scans in order to get the channel intensities
// func prepareLabelStructure(dir, format, plex string, tol float64, mzData mz.Raw) (map[string]tmt.Labels, error) {
//
// 	// get all spectra names from PSMs and create the label list
// 	var labels = make(map[string]tmt.Labels)
// 	ppmPrecision := tol / math.Pow(10, 6)
//
// 	for _, i := range mzData.Spectra {
// 		if i.Level == "2" {
//
// 			tmt, err := tmt.New(plex)
// 			if err != nil {
// 				return labels, err
// 			}
//
// 			// left-pad the spectrum scan
// 			paddedScan := fmt.Sprintf("%05s", i.Scan)
//
// 			tmt.Index = i.Index
// 			tmt.Scan = paddedScan
// 			tmt.ChargeState = i.Precursor.ChargeState
//
// 			for j := range i.Peaks.DecodedStream {
//
// 				if i.Peaks.DecodedStream[j] <= (tmt.Channel1.Mz+(ppmPrecision*tmt.Channel1.Mz)) && i.Peaks.DecodedStream[j] >= (tmt.Channel1.Mz-(ppmPrecision*tmt.Channel1.Mz)) {
// 					if i.Intensities.DecodedStream[j] > tmt.Channel1.Intensity {
// 						tmt.Channel1.Intensity = i.Intensities.DecodedStream[j]
// 					}
// 				}
//
// 				if i.Peaks.DecodedStream[j] <= (tmt.Channel2.Mz+(ppmPrecision*tmt.Channel2.Mz)) && i.Peaks.DecodedStream[j] >= (tmt.Channel2.Mz-(ppmPrecision*tmt.Channel2.Mz)) {
// 					if i.Intensities.DecodedStream[j] > tmt.Channel2.Intensity {
// 						tmt.Channel2.Intensity = i.Intensities.DecodedStream[j]
// 					}
// 				}
//
// 				if i.Peaks.DecodedStream[j] <= (tmt.Channel3.Mz+(ppmPrecision*tmt.Channel3.Mz)) && i.Peaks.DecodedStream[j] >= (tmt.Channel3.Mz-(ppmPrecision*tmt.Channel3.Mz)) {
// 					if i.Intensities.DecodedStream[j] > tmt.Channel3.Intensity {
// 						tmt.Channel3.Intensity = i.Intensities.DecodedStream[j]
// 					}
// 				}
//
// 				if i.Peaks.DecodedStream[j] <= (tmt.Channel4.Mz+(ppmPrecision*tmt.Channel4.Mz)) && i.Peaks.DecodedStream[j] >= (tmt.Channel4.Mz-(ppmPrecision*tmt.Channel4.Mz)) {
// 					if i.Intensities.DecodedStream[j] > tmt.Channel4.Intensity {
// 						tmt.Channel4.Intensity = i.Intensities.DecodedStream[j]
// 					}
// 				}
//
// 				if i.Peaks.DecodedStream[j] <= (tmt.Channel5.Mz+(ppmPrecision*tmt.Channel5.Mz)) && i.Peaks.DecodedStream[j] >= (tmt.Channel5.Mz-(ppmPrecision*tmt.Channel5.Mz)) {
// 					if i.Intensities.DecodedStream[j] > tmt.Channel5.Intensity {
// 						tmt.Channel5.Intensity = i.Intensities.DecodedStream[j]
// 					}
// 				}
//
// 				if i.Peaks.DecodedStream[j] <= (tmt.Channel6.Mz+(ppmPrecision*tmt.Channel6.Mz)) && i.Peaks.DecodedStream[j] >= (tmt.Channel6.Mz-(ppmPrecision*tmt.Channel6.Mz)) {
// 					if i.Intensities.DecodedStream[j] > tmt.Channel6.Intensity {
// 						tmt.Channel6.Intensity = i.Intensities.DecodedStream[j]
// 					}
// 				}
//
// 				if i.Peaks.DecodedStream[j] <= (tmt.Channel7.Mz+(ppmPrecision*tmt.Channel7.Mz)) && i.Peaks.DecodedStream[j] >= (tmt.Channel7.Mz-(ppmPrecision*tmt.Channel7.Mz)) {
// 					if i.Intensities.DecodedStream[j] > tmt.Channel7.Intensity {
// 						tmt.Channel7.Intensity = i.Intensities.DecodedStream[j]
// 					}
// 				}
//
// 				if i.Peaks.DecodedStream[j] <= (tmt.Channel8.Mz+(ppmPrecision*tmt.Channel8.Mz)) && i.Peaks.DecodedStream[j] >= (tmt.Channel8.Mz-(ppmPrecision*tmt.Channel8.Mz)) {
// 					if i.Intensities.DecodedStream[j] > tmt.Channel8.Intensity {
// 						tmt.Channel8.Intensity = i.Intensities.DecodedStream[j]
// 					}
// 				}
//
// 				if i.Peaks.DecodedStream[j] <= (tmt.Channel9.Mz+(ppmPrecision*tmt.Channel9.Mz)) && i.Peaks.DecodedStream[j] >= (tmt.Channel9.Mz-(ppmPrecision*tmt.Channel9.Mz)) {
// 					if i.Intensities.DecodedStream[j] > tmt.Channel9.Intensity {
// 						tmt.Channel9.Intensity = i.Intensities.DecodedStream[j]
// 					}
// 				}
//
// 				if i.Peaks.DecodedStream[j] <= (tmt.Channel10.Mz+(ppmPrecision*tmt.Channel10.Mz)) && i.Peaks.DecodedStream[j] >= (tmt.Channel10.Mz-(ppmPrecision*tmt.Channel10.Mz)) {
// 					if i.Intensities.DecodedStream[j] > tmt.Channel10.Intensity {
// 						tmt.Channel10.Intensity = i.Intensities.DecodedStream[j]
// 					}
// 				}
//
// 				if i.Peaks.DecodedStream[j] > 135 {
// 					break
// 				}
//
// 			}
//
// 			labels[paddedScan] = tmt
//
// 		}
// 	}
//
// 	return labels, nil
// }
//
// // mapLabeledSpectra maps all labeled spectra to ions
// func mapLabeledSpectra(labels map[string]tmt.Labels, purity float64, evi []rep.PSMEvidence) ([]rep.PSMEvidence, error) {
//
// 	for i := range evi {
//
// 		split := strings.Split(evi[i].Spectrum, ".")
//
// 		// referenced by scan number
// 		v, ok := labels[split[2]]
// 		if ok {
//
// 			evi[i].Labels.Spectrum = v.Spectrum
// 			evi[i].Labels.Index = v.Index
// 			evi[i].Labels.Scan = v.Scan
// 			evi[i].Labels.Channel1.Intensity = v.Channel1.Intensity
// 			evi[i].Labels.Channel2.Intensity = v.Channel2.Intensity
// 			evi[i].Labels.Channel3.Intensity = v.Channel3.Intensity
// 			evi[i].Labels.Channel4.Intensity = v.Channel4.Intensity
// 			evi[i].Labels.Channel5.Intensity = v.Channel5.Intensity
// 			evi[i].Labels.Channel6.Intensity = v.Channel6.Intensity
// 			evi[i].Labels.Channel7.Intensity = v.Channel7.Intensity
// 			evi[i].Labels.Channel8.Intensity = v.Channel8.Intensity
// 			evi[i].Labels.Channel9.Intensity = v.Channel9.Intensity
// 			evi[i].Labels.Channel10.Intensity = v.Channel10.Intensity
//
// 		}
// 	}
//
// 	return evi, nil
// }
//
// // rollUpPeptides gathers PSM info and filters them before summing the instensities to the peptide level
// func rollUpPeptides(evi rep.Evidence, spectrumMap map[string]tmt.Labels) rep.Evidence {
//
// 	for j := range evi.Peptides {
// 		for k := range evi.Peptides[j].Spectra {
//
// 			i, ok := spectrumMap[k]
// 			if ok {
//
// 				evi.Peptides[j].Labels.Channel1.Name = i.Channel1.Name
// 				evi.Peptides[j].Labels.Channel1.Mz = i.Channel1.Mz
// 				evi.Peptides[j].Labels.Channel1.Intensity += i.Channel1.Intensity
//
// 				evi.Peptides[j].Labels.Channel2.Name = i.Channel2.Name
// 				evi.Peptides[j].Labels.Channel2.Mz = i.Channel2.Mz
// 				evi.Peptides[j].Labels.Channel2.Intensity += i.Channel2.Intensity
//
// 				evi.Peptides[j].Labels.Channel3.Name = i.Channel3.Name
// 				evi.Peptides[j].Labels.Channel3.Mz = i.Channel3.Mz
// 				evi.Peptides[j].Labels.Channel3.Intensity += i.Channel3.Intensity
//
// 				evi.Peptides[j].Labels.Channel4.Name = i.Channel4.Name
// 				evi.Peptides[j].Labels.Channel4.Mz = i.Channel4.Mz
// 				evi.Peptides[j].Labels.Channel4.Intensity += i.Channel4.Intensity
//
// 				evi.Peptides[j].Labels.Channel5.Name = i.Channel5.Name
// 				evi.Peptides[j].Labels.Channel5.Mz = i.Channel5.Mz
// 				evi.Peptides[j].Labels.Channel5.Intensity += i.Channel5.Intensity
//
// 				evi.Peptides[j].Labels.Channel6.Name = i.Channel6.Name
// 				evi.Peptides[j].Labels.Channel6.Mz = i.Channel6.Mz
// 				evi.Peptides[j].Labels.Channel6.Intensity += i.Channel6.Intensity
//
// 				evi.Peptides[j].Labels.Channel7.Name = i.Channel7.Name
// 				evi.Peptides[j].Labels.Channel7.Mz = i.Channel7.Mz
// 				evi.Peptides[j].Labels.Channel7.Intensity += i.Channel7.Intensity
//
// 				evi.Peptides[j].Labels.Channel8.Name = i.Channel8.Name
// 				evi.Peptides[j].Labels.Channel8.Mz = i.Channel8.Mz
// 				evi.Peptides[j].Labels.Channel8.Intensity += i.Channel8.Intensity
//
// 				evi.Peptides[j].Labels.Channel9.Name = i.Channel9.Name
// 				evi.Peptides[j].Labels.Channel9.Mz = i.Channel9.Mz
// 				evi.Peptides[j].Labels.Channel9.Intensity += i.Channel9.Intensity
//
// 				evi.Peptides[j].Labels.Channel10.Name = i.Channel10.Name
// 				evi.Peptides[j].Labels.Channel10.Mz = i.Channel10.Mz
// 				evi.Peptides[j].Labels.Channel10.Intensity += i.Channel10.Intensity
//
// 				evi.Peptides[j].Labels.Channel11.Name = i.Channel11.Name
// 				evi.Peptides[j].Labels.Channel11.Mz = i.Channel11.Mz
// 				evi.Peptides[j].Labels.Channel11.Intensity += i.Channel11.Intensity
// 			}
// 		}
// 	}
//
// 	return evi
// }
//
// // rollUpPeptideIons gathers PSM info and filters them before summing the instensities to the peptide ION level
// func rollUpPeptideIons(evi rep.Evidence, spectrumMap map[string]tmt.Labels) rep.Evidence {
//
// 	for j := range evi.Ions {
// 		for k := range evi.Ions[j].Spectra {
//
// 			i, ok := spectrumMap[k]
// 			if ok {
//
// 				evi.Ions[j].Labels.Channel1.Name = i.Channel1.Name
// 				evi.Ions[j].Labels.Channel1.Mz = i.Channel1.Mz
// 				evi.Ions[j].Labels.Channel1.Intensity += i.Channel1.Intensity
//
// 				evi.Ions[j].Labels.Channel2.Name = i.Channel2.Name
// 				evi.Ions[j].Labels.Channel2.Mz = i.Channel2.Mz
// 				evi.Ions[j].Labels.Channel2.Intensity += i.Channel2.Intensity
//
// 				evi.Ions[j].Labels.Channel3.Name = i.Channel3.Name
// 				evi.Ions[j].Labels.Channel3.Mz = i.Channel3.Mz
// 				evi.Ions[j].Labels.Channel3.Intensity += i.Channel3.Intensity
//
// 				evi.Ions[j].Labels.Channel4.Name = i.Channel4.Name
// 				evi.Ions[j].Labels.Channel4.Mz = i.Channel4.Mz
// 				evi.Ions[j].Labels.Channel4.Intensity += i.Channel4.Intensity
//
// 				evi.Ions[j].Labels.Channel5.Name = i.Channel5.Name
// 				evi.Ions[j].Labels.Channel5.Mz = i.Channel5.Mz
// 				evi.Ions[j].Labels.Channel5.Intensity += i.Channel5.Intensity
//
// 				evi.Ions[j].Labels.Channel6.Name = i.Channel6.Name
// 				evi.Ions[j].Labels.Channel6.Mz = i.Channel6.Mz
// 				evi.Ions[j].Labels.Channel6.Intensity += i.Channel6.Intensity
//
// 				evi.Ions[j].Labels.Channel7.Name = i.Channel7.Name
// 				evi.Ions[j].Labels.Channel7.Mz = i.Channel7.Mz
// 				evi.Ions[j].Labels.Channel7.Intensity += i.Channel7.Intensity
//
// 				evi.Ions[j].Labels.Channel8.Name = i.Channel8.Name
// 				evi.Ions[j].Labels.Channel8.Mz = i.Channel8.Mz
// 				evi.Ions[j].Labels.Channel8.Intensity += i.Channel8.Intensity
//
// 				evi.Ions[j].Labels.Channel9.Name = i.Channel9.Name
// 				evi.Ions[j].Labels.Channel9.Mz = i.Channel9.Mz
// 				evi.Ions[j].Labels.Channel9.Intensity += i.Channel9.Intensity
//
// 				evi.Ions[j].Labels.Channel10.Name = i.Channel10.Name
// 				evi.Ions[j].Labels.Channel10.Mz = i.Channel10.Mz
// 				evi.Ions[j].Labels.Channel10.Intensity += i.Channel10.Intensity
//
// 				evi.Ions[j].Labels.Channel11.Name = i.Channel11.Name
// 				evi.Ions[j].Labels.Channel11.Mz = i.Channel11.Mz
// 				evi.Ions[j].Labels.Channel11.Intensity += i.Channel11.Intensity
// 			}
// 		}
// 	}
//
// 	return evi
// }
//
// // rollUpProteins gathers PSM info and filters them before summing the instensities to the peptide ION level
// func rollUpProteins(evi rep.Evidence, spectrumMap map[string]tmt.Labels) rep.Evidence {
//
// 	for j := range evi.Proteins {
// 		for _, k := range evi.Proteins[j].TotalPeptideIons {
// 			for l := range k.Spectra {
//
// 				i, ok := spectrumMap[l]
// 				if ok {
// 					evi.Proteins[j].TotalLabels.Channel1.Name = i.Channel1.Name
// 					evi.Proteins[j].TotalLabels.Channel1.Mz = i.Channel1.Mz
// 					evi.Proteins[j].TotalLabels.Channel1.Intensity += i.Channel1.Intensity
//
// 					evi.Proteins[j].TotalLabels.Channel2.Name = i.Channel2.Name
// 					evi.Proteins[j].TotalLabels.Channel2.Mz = i.Channel2.Mz
// 					evi.Proteins[j].TotalLabels.Channel2.Intensity += i.Channel2.Intensity
//
// 					evi.Proteins[j].TotalLabels.Channel3.Name = i.Channel3.Name
// 					evi.Proteins[j].TotalLabels.Channel3.Mz = i.Channel3.Mz
// 					evi.Proteins[j].TotalLabels.Channel3.Intensity += i.Channel3.Intensity
//
// 					evi.Proteins[j].TotalLabels.Channel4.Name = i.Channel4.Name
// 					evi.Proteins[j].TotalLabels.Channel4.Mz = i.Channel4.Mz
// 					evi.Proteins[j].TotalLabels.Channel4.Intensity += i.Channel4.Intensity
//
// 					evi.Proteins[j].TotalLabels.Channel5.Name = i.Channel5.Name
// 					evi.Proteins[j].TotalLabels.Channel5.Mz = i.Channel5.Mz
// 					evi.Proteins[j].TotalLabels.Channel5.Intensity += i.Channel5.Intensity
//
// 					evi.Proteins[j].TotalLabels.Channel6.Name = i.Channel6.Name
// 					evi.Proteins[j].TotalLabels.Channel6.Mz = i.Channel6.Mz
// 					evi.Proteins[j].TotalLabels.Channel6.Intensity += i.Channel6.Intensity
//
// 					evi.Proteins[j].TotalLabels.Channel7.Name = i.Channel7.Name
// 					evi.Proteins[j].TotalLabels.Channel7.Mz = i.Channel7.Mz
// 					evi.Proteins[j].TotalLabels.Channel7.Intensity += i.Channel7.Intensity
//
// 					evi.Proteins[j].TotalLabels.Channel8.Name = i.Channel8.Name
// 					evi.Proteins[j].TotalLabels.Channel8.Mz = i.Channel8.Mz
// 					evi.Proteins[j].TotalLabels.Channel8.Intensity += i.Channel8.Intensity
//
// 					evi.Proteins[j].TotalLabels.Channel9.Name = i.Channel9.Name
// 					evi.Proteins[j].TotalLabels.Channel9.Mz = i.Channel9.Mz
// 					evi.Proteins[j].TotalLabels.Channel9.Intensity += i.Channel9.Intensity
//
// 					evi.Proteins[j].TotalLabels.Channel10.Name = i.Channel10.Name
// 					evi.Proteins[j].TotalLabels.Channel10.Mz = i.Channel10.Mz
// 					evi.Proteins[j].TotalLabels.Channel10.Intensity += i.Channel10.Intensity
//
// 					evi.Proteins[j].TotalLabels.Channel11.Name = i.Channel11.Name
// 					evi.Proteins[j].TotalLabels.Channel11.Mz = i.Channel11.Mz
// 					evi.Proteins[j].TotalLabels.Channel11.Intensity += i.Channel11.Intensity
//
// 					if k.IsNondegenerateEvidence {
// 						evi.Proteins[j].UniqueLabels.Channel1.Name = i.Channel1.Name
// 						evi.Proteins[j].UniqueLabels.Channel1.Mz = i.Channel1.Mz
// 						evi.Proteins[j].UniqueLabels.Channel1.Intensity += i.Channel1.Intensity
//
// 						evi.Proteins[j].UniqueLabels.Channel2.Name = i.Channel2.Name
// 						evi.Proteins[j].UniqueLabels.Channel2.Mz = i.Channel2.Mz
// 						evi.Proteins[j].UniqueLabels.Channel2.Intensity += i.Channel2.Intensity
//
// 						evi.Proteins[j].UniqueLabels.Channel3.Name = i.Channel3.Name
// 						evi.Proteins[j].UniqueLabels.Channel3.Mz = i.Channel3.Mz
// 						evi.Proteins[j].UniqueLabels.Channel3.Intensity += i.Channel3.Intensity
//
// 						evi.Proteins[j].UniqueLabels.Channel4.Name = i.Channel4.Name
// 						evi.Proteins[j].UniqueLabels.Channel4.Mz = i.Channel4.Mz
// 						evi.Proteins[j].UniqueLabels.Channel4.Intensity += i.Channel4.Intensity
//
// 						evi.Proteins[j].UniqueLabels.Channel5.Name = i.Channel5.Name
// 						evi.Proteins[j].UniqueLabels.Channel5.Mz = i.Channel5.Mz
// 						evi.Proteins[j].UniqueLabels.Channel5.Intensity += i.Channel5.Intensity
//
// 						evi.Proteins[j].UniqueLabels.Channel6.Name = i.Channel6.Name
// 						evi.Proteins[j].UniqueLabels.Channel6.Mz = i.Channel6.Mz
// 						evi.Proteins[j].UniqueLabels.Channel6.Intensity += i.Channel6.Intensity
//
// 						evi.Proteins[j].UniqueLabels.Channel7.Name = i.Channel7.Name
// 						evi.Proteins[j].UniqueLabels.Channel7.Mz = i.Channel7.Mz
// 						evi.Proteins[j].UniqueLabels.Channel7.Intensity += i.Channel7.Intensity
//
// 						evi.Proteins[j].UniqueLabels.Channel8.Name = i.Channel8.Name
// 						evi.Proteins[j].UniqueLabels.Channel8.Mz = i.Channel8.Mz
// 						evi.Proteins[j].UniqueLabels.Channel8.Intensity += i.Channel8.Intensity
//
// 						evi.Proteins[j].UniqueLabels.Channel9.Name = i.Channel9.Name
// 						evi.Proteins[j].UniqueLabels.Channel9.Mz = i.Channel9.Mz
// 						evi.Proteins[j].UniqueLabels.Channel9.Intensity += i.Channel9.Intensity
//
// 						evi.Proteins[j].UniqueLabels.Channel10.Name = i.Channel10.Name
// 						evi.Proteins[j].UniqueLabels.Channel10.Mz = i.Channel10.Mz
// 						evi.Proteins[j].UniqueLabels.Channel10.Intensity += i.Channel10.Intensity
//
// 						evi.Proteins[j].UniqueLabels.Channel11.Name = i.Channel11.Name
// 						evi.Proteins[j].UniqueLabels.Channel11.Mz = i.Channel11.Mz
// 						evi.Proteins[j].UniqueLabels.Channel11.Intensity += i.Channel11.Intensity
// 					}
//
// 					if k.IsURazor {
// 						evi.Proteins[j].URazorLabels.Channel1.Name = i.Channel1.Name
// 						evi.Proteins[j].URazorLabels.Channel1.Mz = i.Channel1.Mz
// 						evi.Proteins[j].URazorLabels.Channel1.Intensity += i.Channel1.Intensity
//
// 						evi.Proteins[j].URazorLabels.Channel2.Name = i.Channel2.Name
// 						evi.Proteins[j].URazorLabels.Channel2.Mz = i.Channel2.Mz
// 						evi.Proteins[j].URazorLabels.Channel2.Intensity += i.Channel2.Intensity
//
// 						evi.Proteins[j].URazorLabels.Channel3.Name = i.Channel3.Name
// 						evi.Proteins[j].URazorLabels.Channel3.Mz = i.Channel3.Mz
// 						evi.Proteins[j].URazorLabels.Channel3.Intensity += i.Channel3.Intensity
//
// 						evi.Proteins[j].URazorLabels.Channel4.Name = i.Channel4.Name
// 						evi.Proteins[j].URazorLabels.Channel4.Mz = i.Channel4.Mz
// 						evi.Proteins[j].URazorLabels.Channel4.Intensity += i.Channel4.Intensity
//
// 						evi.Proteins[j].URazorLabels.Channel5.Name = i.Channel5.Name
// 						evi.Proteins[j].URazorLabels.Channel5.Mz = i.Channel5.Mz
// 						evi.Proteins[j].URazorLabels.Channel5.Intensity += i.Channel5.Intensity
//
// 						evi.Proteins[j].URazorLabels.Channel6.Name = i.Channel6.Name
// 						evi.Proteins[j].URazorLabels.Channel6.Mz = i.Channel6.Mz
// 						evi.Proteins[j].URazorLabels.Channel6.Intensity += i.Channel6.Intensity
//
// 						evi.Proteins[j].URazorLabels.Channel7.Name = i.Channel7.Name
// 						evi.Proteins[j].URazorLabels.Channel7.Mz = i.Channel7.Mz
// 						evi.Proteins[j].URazorLabels.Channel7.Intensity += i.Channel7.Intensity
//
// 						evi.Proteins[j].URazorLabels.Channel8.Name = i.Channel8.Name
// 						evi.Proteins[j].URazorLabels.Channel8.Mz = i.Channel8.Mz
// 						evi.Proteins[j].URazorLabels.Channel8.Intensity += i.Channel8.Intensity
//
// 						evi.Proteins[j].URazorLabels.Channel9.Name = i.Channel9.Name
// 						evi.Proteins[j].URazorLabels.Channel9.Mz = i.Channel9.Mz
// 						evi.Proteins[j].URazorLabels.Channel9.Intensity += i.Channel9.Intensity
//
// 						evi.Proteins[j].URazorLabels.Channel10.Name = i.Channel10.Name
// 						evi.Proteins[j].URazorLabels.Channel10.Mz = i.Channel10.Mz
// 						evi.Proteins[j].URazorLabels.Channel10.Intensity += i.Channel10.Intensity
//
// 						evi.Proteins[j].URazorLabels.Channel11.Name = i.Channel11.Name
// 						evi.Proteins[j].URazorLabels.Channel11.Mz = i.Channel11.Mz
// 						evi.Proteins[j].URazorLabels.Channel11.Intensity += i.Channel11.Intensity
// 					}
//
// 				}
//
// 			}
// 		}
// 	}
//
// 	return evi
// }
//
// // NormToTotalProteins calculates the protein level normalization based on total proteins
// func NormToTotalProteins(evi rep.Evidence) rep.Evidence {
//
// 	var topValue float64
// 	var channelSum = [10]float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
// 	var normFactors = [10]float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
//
// 	// sum TMT singal for each column
// 	for _, i := range evi.Proteins {
// 		channelSum[0] += i.URazorLabels.Channel1.Intensity
// 		channelSum[1] += i.URazorLabels.Channel2.Intensity
// 		channelSum[2] += i.URazorLabels.Channel3.Intensity
// 		channelSum[3] += i.URazorLabels.Channel4.Intensity
// 		channelSum[4] += i.URazorLabels.Channel5.Intensity
// 		channelSum[5] += i.URazorLabels.Channel6.Intensity
// 		channelSum[6] += i.URazorLabels.Channel7.Intensity
// 		channelSum[7] += i.URazorLabels.Channel8.Intensity
// 		channelSum[8] += i.URazorLabels.Channel9.Intensity
// 		channelSum[9] += i.URazorLabels.Channel10.Intensity
// 	}
//
// 	// find the higest value amongst channels
// 	for _, i := range channelSum {
// 		if i > topValue {
// 			topValue = i
// 		}
// 	}
//
// 	// calculate normalizing factors
// 	for i := range channelSum {
// 		normFactors[i] = channelSum[i] / topValue
// 	}
//
// 	// multiply each protein TMT set by the factors to get normalized values
// 	for _, i := range evi.Proteins {
// 		i.URazorLabels.Channel1.Intensity *= normFactors[0]
// 		i.URazorLabels.Channel2.Intensity *= normFactors[1]
// 		i.URazorLabels.Channel3.Intensity *= normFactors[2]
// 		i.URazorLabels.Channel4.Intensity *= normFactors[3]
// 		i.URazorLabels.Channel5.Intensity *= normFactors[4]
// 		i.URazorLabels.Channel6.Intensity *= normFactors[5]
// 		i.URazorLabels.Channel7.Intensity *= normFactors[6]
// 		i.URazorLabels.Channel8.Intensity *= normFactors[7]
// 		i.URazorLabels.Channel9.Intensity *= normFactors[8]
// 		i.URazorLabels.Channel10.Intensity *= normFactors[9]
// 	}
//
// 	return evi
// }
