package quan

import (
	"fmt"
	"math"
	"strings"

	"github.com/prvst/philosopher/lib/rep"
	"github.com/prvst/philosopher/lib/tmt"
	"github.com/prvst/philosopher/lib/utils"
	"github.com/prvst/philosopher/lib/xml"
)

const (
	mzDeltaWindow float64 = 0.5
)

// calculateIonPurity verifies how much interference there is on the precursor scans for each fragment
func calculateIonPurity(d, f string, mzData xml.Raw, evi []rep.PSMEvidence) ([]rep.PSMEvidence, error) {

	// index spectra in a dictionary
	var indexedMz = make(map[string]xml.Spectrum)
	for _, i := range mzData.Spectra {

		if i.Level == "2" && i.Precursor.IsolationWindowLowerOffset == 0 && i.Precursor.IsolationWindowUpperOffset == 0 {
			i.Precursor.IsolationWindowLowerOffset = mzDeltaWindow
			i.Precursor.IsolationWindowUpperOffset = mzDeltaWindow
		}

		// left-pad the spectrum index
		paddedIndex := fmt.Sprintf("%05s", i.Index)

		// left-pad the spectrum scan
		paddedScan := fmt.Sprintf("%05s", i.Scan)

		// left-pad the precursor spectrum index
		paddedPI := fmt.Sprintf("%05s", i.Precursor.ParentIndex)

		// left-pad the precursor spectrum scan
		paddedPS := fmt.Sprintf("%05s", i.Precursor.ParentScan)

		i.Index = paddedIndex
		i.Scan = paddedScan
		i.Precursor.ParentIndex = paddedPI
		i.Precursor.ParentScan = paddedPS

		indexedMz[paddedScan] = i
	}

	for i := range evi {

		// get spectrum index
		split := strings.Split(evi[i].Spectrum, ".")

		ms2, ok := indexedMz[split[1]]
		if ok {

			ms1 := indexedMz[ms2.Precursor.ParentScan]

			var ions = make(map[float64]float64)
			var isolationWindowSummedInt float64
			for k := range ms1.Peaks {
				if ms1.Peaks[k] >= (ms2.Precursor.TargetIon-ms2.Precursor.IsolationWindowUpperOffset) && ms1.Peaks[k] <= (ms2.Precursor.TargetIon+ms2.Precursor.IsolationWindowUpperOffset) {
					ions[ms1.Peaks[k]] = ms1.Intensities[k]
					isolationWindowSummedInt += ms1.Intensities[k]
				}
			}

			// create the list of mz differences for each peak
			var mzRatio []float64
			for k := 1; k <= 6; k++ {
				r := float64(k) * (float64(1) / float64(ms2.Precursor.ChargeState))
				mzRatio = append(mzRatio, utils.ToFixed(r, 2))
			}

			var isotopePackage = make(map[float64]float64)

			isotopePackage[ms2.Precursor.TargetIon] = ms2.Precursor.PeakIntensity
			isotopesInt := ms2.Precursor.PeakIntensity

			for k, v := range ions {
				for _, m := range mzRatio {
					if math.Abs(ms2.Precursor.TargetIon-k) <= (m+0.02) && math.Abs(ms2.Precursor.TargetIon-k) >= (m-0.02) {
						isotopePackage[k] = v
						isotopesInt += v
						break
					}
				}
			}

			if isotopesInt == 0 {
				evi[i].Purity = 0
			} else {
				evi[i].Purity = utils.Round((isotopesInt / isolationWindowSummedInt), 5, 2)
			}

			if evi[i].Purity > 1 {
				evi[i].Purity = 1
			}

		}

	}

	return evi, nil
}

// prepareLabelStructure instantiates the Label objects and maps them against the fragment scans in order to get the channel intensities
func prepareLabelStructure(dir, format, plex string, tol float64, mzData xml.Raw) (map[string]tmt.Labels, error) {

	// get all spectra names from PSMs and create the label list
	var labels = make(map[string]tmt.Labels)
	ppmPrecision := tol / math.Pow(10, 6)

	for _, i := range mzData.Spectra {
		if i.Level == "2" {

			tmt, err := tmt.New(plex)
			if err != nil {
				return labels, err
			}

			// left-pad the spectrum scan
			paddedScan := fmt.Sprintf("%05s", i.Scan)

			tmt.Index = i.Index
			tmt.Scan = paddedScan
			tmt.ChargeState = i.Precursor.ChargeState

			for j := range i.Peaks {

				if i.Peaks[j] <= (tmt.Channel1.Mz+(ppmPrecision*tmt.Channel1.Mz)) && i.Peaks[j] >= (tmt.Channel1.Mz-(ppmPrecision*tmt.Channel1.Mz)) {
					if i.Intensities[j] > tmt.Channel1.Intensity {
						tmt.Channel1.Intensity = i.Intensities[j]
					}
				}

				if i.Peaks[j] <= (tmt.Channel2.Mz+(ppmPrecision*tmt.Channel2.Mz)) && i.Peaks[j] >= (tmt.Channel2.Mz-(ppmPrecision*tmt.Channel2.Mz)) {
					if i.Intensities[j] > tmt.Channel2.Intensity {
						tmt.Channel2.Intensity = i.Intensities[j]
					}
				}

				if i.Peaks[j] <= (tmt.Channel3.Mz+(ppmPrecision*tmt.Channel3.Mz)) && i.Peaks[j] >= (tmt.Channel3.Mz-(ppmPrecision*tmt.Channel3.Mz)) {
					if i.Intensities[j] > tmt.Channel3.Intensity {
						tmt.Channel3.Intensity = i.Intensities[j]
					}
				}

				if i.Peaks[j] <= (tmt.Channel4.Mz+(ppmPrecision*tmt.Channel4.Mz)) && i.Peaks[j] >= (tmt.Channel4.Mz-(ppmPrecision*tmt.Channel4.Mz)) {
					if i.Intensities[j] > tmt.Channel4.Intensity {
						tmt.Channel4.Intensity = i.Intensities[j]
					}
				}

				if i.Peaks[j] <= (tmt.Channel5.Mz+(ppmPrecision*tmt.Channel5.Mz)) && i.Peaks[j] >= (tmt.Channel5.Mz-(ppmPrecision*tmt.Channel5.Mz)) {
					if i.Intensities[j] > tmt.Channel5.Intensity {
						tmt.Channel5.Intensity = i.Intensities[j]
					}
				}

				if i.Peaks[j] <= (tmt.Channel6.Mz+(ppmPrecision*tmt.Channel6.Mz)) && i.Peaks[j] >= (tmt.Channel6.Mz-(ppmPrecision*tmt.Channel6.Mz)) {
					if i.Intensities[j] > tmt.Channel6.Intensity {
						tmt.Channel6.Intensity = i.Intensities[j]
					}
				}

				if i.Peaks[j] <= (tmt.Channel7.Mz+(ppmPrecision*tmt.Channel7.Mz)) && i.Peaks[j] >= (tmt.Channel7.Mz-(ppmPrecision*tmt.Channel7.Mz)) {
					if i.Intensities[j] > tmt.Channel7.Intensity {
						tmt.Channel7.Intensity = i.Intensities[j]
					}
				}

				if i.Peaks[j] <= (tmt.Channel8.Mz+(ppmPrecision*tmt.Channel8.Mz)) && i.Peaks[j] >= (tmt.Channel8.Mz-(ppmPrecision*tmt.Channel8.Mz)) {
					if i.Intensities[j] > tmt.Channel8.Intensity {
						tmt.Channel8.Intensity = i.Intensities[j]
					}
				}

				if i.Peaks[j] <= (tmt.Channel9.Mz+(ppmPrecision*tmt.Channel9.Mz)) && i.Peaks[j] >= (tmt.Channel9.Mz-(ppmPrecision*tmt.Channel9.Mz)) {
					if i.Intensities[j] > tmt.Channel9.Intensity {
						tmt.Channel9.Intensity = i.Intensities[j]
					}
				}

				if i.Peaks[j] <= (tmt.Channel10.Mz+(ppmPrecision*tmt.Channel10.Mz)) && i.Peaks[j] >= (tmt.Channel10.Mz-(ppmPrecision*tmt.Channel10.Mz)) {
					if i.Intensities[j] > tmt.Channel10.Intensity {
						tmt.Channel10.Intensity = i.Intensities[j]
					}
				}

				if i.Peaks[j] > 135 {
					break
				}

			}

			labels[paddedScan] = tmt

		}
	}

	return labels, nil
}

// mapLabeledSpectra maps all labeled spectra to ions
func mapLabeledSpectra(labels map[string]tmt.Labels, purity float64, evi []rep.PSMEvidence) ([]rep.PSMEvidence, error) {

	for i := range evi {

		split := strings.Split(evi[i].Spectrum, ".")

		// referenced by scan number
		v, ok := labels[split[2]]
		if ok {

			evi[i].Labels.Spectrum = v.Spectrum
			evi[i].Labels.Index = v.Index
			evi[i].Labels.Scan = v.Scan
			evi[i].Labels.Channel1.Intensity = v.Channel1.Intensity
			evi[i].Labels.Channel2.Intensity = v.Channel2.Intensity
			evi[i].Labels.Channel3.Intensity = v.Channel3.Intensity
			evi[i].Labels.Channel4.Intensity = v.Channel4.Intensity
			evi[i].Labels.Channel5.Intensity = v.Channel5.Intensity
			evi[i].Labels.Channel6.Intensity = v.Channel6.Intensity
			evi[i].Labels.Channel7.Intensity = v.Channel7.Intensity
			evi[i].Labels.Channel8.Intensity = v.Channel8.Intensity
			evi[i].Labels.Channel9.Intensity = v.Channel9.Intensity
			evi[i].Labels.Channel10.Intensity = v.Channel10.Intensity

		}
	}

	return evi, nil
}

// rollUpPeptides gathers PSM info and filters them before summing the instensities to the peptide level
func rollUpPeptides(evi rep.Evidence, spectrumMap map[string]tmt.Labels) rep.Evidence {

	for j := range evi.Peptides {
		for k := range evi.Peptides[j].Spectra {

			i, ok := spectrumMap[k]
			if ok {

				evi.Peptides[j].Labels.Channel1.Name = i.Channel1.Name
				evi.Peptides[j].Labels.Channel1.Mz = i.Channel1.Mz
				evi.Peptides[j].Labels.Channel1.Intensity += i.Channel1.Intensity

				evi.Peptides[j].Labels.Channel2.Name = i.Channel2.Name
				evi.Peptides[j].Labels.Channel2.Mz = i.Channel2.Mz
				evi.Peptides[j].Labels.Channel2.Intensity += i.Channel2.Intensity

				evi.Peptides[j].Labels.Channel3.Name = i.Channel3.Name
				evi.Peptides[j].Labels.Channel3.Mz = i.Channel3.Mz
				evi.Peptides[j].Labels.Channel3.Intensity += i.Channel3.Intensity

				evi.Peptides[j].Labels.Channel4.Name = i.Channel4.Name
				evi.Peptides[j].Labels.Channel4.Mz = i.Channel4.Mz
				evi.Peptides[j].Labels.Channel4.Intensity += i.Channel4.Intensity

				evi.Peptides[j].Labels.Channel5.Name = i.Channel5.Name
				evi.Peptides[j].Labels.Channel5.Mz = i.Channel5.Mz
				evi.Peptides[j].Labels.Channel5.Intensity += i.Channel5.Intensity

				evi.Peptides[j].Labels.Channel6.Name = i.Channel6.Name
				evi.Peptides[j].Labels.Channel6.Mz = i.Channel6.Mz
				evi.Peptides[j].Labels.Channel6.Intensity += i.Channel6.Intensity

				evi.Peptides[j].Labels.Channel7.Name = i.Channel7.Name
				evi.Peptides[j].Labels.Channel7.Mz = i.Channel7.Mz
				evi.Peptides[j].Labels.Channel7.Intensity += i.Channel7.Intensity

				evi.Peptides[j].Labels.Channel8.Name = i.Channel8.Name
				evi.Peptides[j].Labels.Channel8.Mz = i.Channel8.Mz
				evi.Peptides[j].Labels.Channel8.Intensity += i.Channel8.Intensity

				evi.Peptides[j].Labels.Channel9.Name = i.Channel9.Name
				evi.Peptides[j].Labels.Channel9.Mz = i.Channel9.Mz
				evi.Peptides[j].Labels.Channel9.Intensity += i.Channel9.Intensity

				evi.Peptides[j].Labels.Channel10.Name = i.Channel10.Name
				evi.Peptides[j].Labels.Channel10.Mz = i.Channel10.Mz
				evi.Peptides[j].Labels.Channel10.Intensity += i.Channel10.Intensity

				evi.Peptides[j].Labels.Channel11.Name = i.Channel11.Name
				evi.Peptides[j].Labels.Channel11.Mz = i.Channel11.Mz
				evi.Peptides[j].Labels.Channel11.Intensity += i.Channel11.Intensity
			}
		}
	}

	return evi
}

// rollUpPeptideIons gathers PSM info and filters them before summing the instensities to the peptide ION level
func rollUpPeptideIons(evi rep.Evidence, spectrumMap map[string]tmt.Labels) rep.Evidence {

	for j := range evi.Ions {
		for k := range evi.Ions[j].Spectra {

			i, ok := spectrumMap[k]
			if ok {

				evi.Ions[j].Labels.Channel1.Name = i.Channel1.Name
				evi.Ions[j].Labels.Channel1.Mz = i.Channel1.Mz
				evi.Ions[j].Labels.Channel1.Intensity += i.Channel1.Intensity

				evi.Ions[j].Labels.Channel2.Name = i.Channel2.Name
				evi.Ions[j].Labels.Channel2.Mz = i.Channel2.Mz
				evi.Ions[j].Labels.Channel2.Intensity += i.Channel2.Intensity

				evi.Ions[j].Labels.Channel3.Name = i.Channel3.Name
				evi.Ions[j].Labels.Channel3.Mz = i.Channel3.Mz
				evi.Ions[j].Labels.Channel3.Intensity += i.Channel3.Intensity

				evi.Ions[j].Labels.Channel4.Name = i.Channel4.Name
				evi.Ions[j].Labels.Channel4.Mz = i.Channel4.Mz
				evi.Ions[j].Labels.Channel4.Intensity += i.Channel4.Intensity

				evi.Ions[j].Labels.Channel5.Name = i.Channel5.Name
				evi.Ions[j].Labels.Channel5.Mz = i.Channel5.Mz
				evi.Ions[j].Labels.Channel5.Intensity += i.Channel5.Intensity

				evi.Ions[j].Labels.Channel6.Name = i.Channel6.Name
				evi.Ions[j].Labels.Channel6.Mz = i.Channel6.Mz
				evi.Ions[j].Labels.Channel6.Intensity += i.Channel6.Intensity

				evi.Ions[j].Labels.Channel7.Name = i.Channel7.Name
				evi.Ions[j].Labels.Channel7.Mz = i.Channel7.Mz
				evi.Ions[j].Labels.Channel7.Intensity += i.Channel7.Intensity

				evi.Ions[j].Labels.Channel8.Name = i.Channel8.Name
				evi.Ions[j].Labels.Channel8.Mz = i.Channel8.Mz
				evi.Ions[j].Labels.Channel8.Intensity += i.Channel8.Intensity

				evi.Ions[j].Labels.Channel9.Name = i.Channel9.Name
				evi.Ions[j].Labels.Channel9.Mz = i.Channel9.Mz
				evi.Ions[j].Labels.Channel9.Intensity += i.Channel9.Intensity

				evi.Ions[j].Labels.Channel10.Name = i.Channel10.Name
				evi.Ions[j].Labels.Channel10.Mz = i.Channel10.Mz
				evi.Ions[j].Labels.Channel10.Intensity += i.Channel10.Intensity

				evi.Ions[j].Labels.Channel11.Name = i.Channel11.Name
				evi.Ions[j].Labels.Channel11.Mz = i.Channel11.Mz
				evi.Ions[j].Labels.Channel11.Intensity += i.Channel11.Intensity
			}
		}
	}

	return evi
}

// rollUpProteins gathers PSM info and filters them before summing the instensities to the peptide ION level
func rollUpProteins(evi rep.Evidence, spectrumMap map[string]tmt.Labels) rep.Evidence {

	for j := range evi.Proteins {
		for _, k := range evi.Proteins[j].TotalPeptideIons {
			for l := range k.Spectra {

				i, ok := spectrumMap[l]
				if ok {
					evi.Proteins[j].TotalLabels.Channel1.Name = i.Channel1.Name
					evi.Proteins[j].TotalLabels.Channel1.Mz = i.Channel1.Mz
					evi.Proteins[j].TotalLabels.Channel1.Intensity += i.Channel1.Intensity

					evi.Proteins[j].TotalLabels.Channel2.Name = i.Channel2.Name
					evi.Proteins[j].TotalLabels.Channel2.Mz = i.Channel2.Mz
					evi.Proteins[j].TotalLabels.Channel2.Intensity += i.Channel2.Intensity

					evi.Proteins[j].TotalLabels.Channel3.Name = i.Channel3.Name
					evi.Proteins[j].TotalLabels.Channel3.Mz = i.Channel3.Mz
					evi.Proteins[j].TotalLabels.Channel3.Intensity += i.Channel3.Intensity

					evi.Proteins[j].TotalLabels.Channel4.Name = i.Channel4.Name
					evi.Proteins[j].TotalLabels.Channel4.Mz = i.Channel4.Mz
					evi.Proteins[j].TotalLabels.Channel4.Intensity += i.Channel4.Intensity

					evi.Proteins[j].TotalLabels.Channel5.Name = i.Channel5.Name
					evi.Proteins[j].TotalLabels.Channel5.Mz = i.Channel5.Mz
					evi.Proteins[j].TotalLabels.Channel5.Intensity += i.Channel5.Intensity

					evi.Proteins[j].TotalLabels.Channel6.Name = i.Channel6.Name
					evi.Proteins[j].TotalLabels.Channel6.Mz = i.Channel6.Mz
					evi.Proteins[j].TotalLabels.Channel6.Intensity += i.Channel6.Intensity

					evi.Proteins[j].TotalLabels.Channel7.Name = i.Channel7.Name
					evi.Proteins[j].TotalLabels.Channel7.Mz = i.Channel7.Mz
					evi.Proteins[j].TotalLabels.Channel7.Intensity += i.Channel7.Intensity

					evi.Proteins[j].TotalLabels.Channel8.Name = i.Channel8.Name
					evi.Proteins[j].TotalLabels.Channel8.Mz = i.Channel8.Mz
					evi.Proteins[j].TotalLabels.Channel8.Intensity += i.Channel8.Intensity

					evi.Proteins[j].TotalLabels.Channel9.Name = i.Channel9.Name
					evi.Proteins[j].TotalLabels.Channel9.Mz = i.Channel9.Mz
					evi.Proteins[j].TotalLabels.Channel9.Intensity += i.Channel9.Intensity

					evi.Proteins[j].TotalLabels.Channel10.Name = i.Channel10.Name
					evi.Proteins[j].TotalLabels.Channel10.Mz = i.Channel10.Mz
					evi.Proteins[j].TotalLabels.Channel10.Intensity += i.Channel10.Intensity

					evi.Proteins[j].TotalLabels.Channel11.Name = i.Channel11.Name
					evi.Proteins[j].TotalLabels.Channel11.Mz = i.Channel11.Mz
					evi.Proteins[j].TotalLabels.Channel11.Intensity += i.Channel11.Intensity

					if k.IsNondegenerateEvidence {
						evi.Proteins[j].UniqueLabels.Channel1.Name = i.Channel1.Name
						evi.Proteins[j].UniqueLabels.Channel1.Mz = i.Channel1.Mz
						evi.Proteins[j].UniqueLabels.Channel1.Intensity += i.Channel1.Intensity

						evi.Proteins[j].UniqueLabels.Channel2.Name = i.Channel2.Name
						evi.Proteins[j].UniqueLabels.Channel2.Mz = i.Channel2.Mz
						evi.Proteins[j].UniqueLabels.Channel2.Intensity += i.Channel2.Intensity

						evi.Proteins[j].UniqueLabels.Channel3.Name = i.Channel3.Name
						evi.Proteins[j].UniqueLabels.Channel3.Mz = i.Channel3.Mz
						evi.Proteins[j].UniqueLabels.Channel3.Intensity += i.Channel3.Intensity

						evi.Proteins[j].UniqueLabels.Channel4.Name = i.Channel4.Name
						evi.Proteins[j].UniqueLabels.Channel4.Mz = i.Channel4.Mz
						evi.Proteins[j].UniqueLabels.Channel4.Intensity += i.Channel4.Intensity

						evi.Proteins[j].UniqueLabels.Channel5.Name = i.Channel5.Name
						evi.Proteins[j].UniqueLabels.Channel5.Mz = i.Channel5.Mz
						evi.Proteins[j].UniqueLabels.Channel5.Intensity += i.Channel5.Intensity

						evi.Proteins[j].UniqueLabels.Channel6.Name = i.Channel6.Name
						evi.Proteins[j].UniqueLabels.Channel6.Mz = i.Channel6.Mz
						evi.Proteins[j].UniqueLabels.Channel6.Intensity += i.Channel6.Intensity

						evi.Proteins[j].UniqueLabels.Channel7.Name = i.Channel7.Name
						evi.Proteins[j].UniqueLabels.Channel7.Mz = i.Channel7.Mz
						evi.Proteins[j].UniqueLabels.Channel7.Intensity += i.Channel7.Intensity

						evi.Proteins[j].UniqueLabels.Channel8.Name = i.Channel8.Name
						evi.Proteins[j].UniqueLabels.Channel8.Mz = i.Channel8.Mz
						evi.Proteins[j].UniqueLabels.Channel8.Intensity += i.Channel8.Intensity

						evi.Proteins[j].UniqueLabels.Channel9.Name = i.Channel9.Name
						evi.Proteins[j].UniqueLabels.Channel9.Mz = i.Channel9.Mz
						evi.Proteins[j].UniqueLabels.Channel9.Intensity += i.Channel9.Intensity

						evi.Proteins[j].UniqueLabels.Channel10.Name = i.Channel10.Name
						evi.Proteins[j].UniqueLabels.Channel10.Mz = i.Channel10.Mz
						evi.Proteins[j].UniqueLabels.Channel10.Intensity += i.Channel10.Intensity

						evi.Proteins[j].UniqueLabels.Channel11.Name = i.Channel11.Name
						evi.Proteins[j].UniqueLabels.Channel11.Mz = i.Channel11.Mz
						evi.Proteins[j].UniqueLabels.Channel11.Intensity += i.Channel11.Intensity
					}

					if k.IsURazor {
						evi.Proteins[j].URazorLabels.Channel1.Name = i.Channel1.Name
						evi.Proteins[j].URazorLabels.Channel1.Mz = i.Channel1.Mz
						evi.Proteins[j].URazorLabels.Channel1.Intensity += i.Channel1.Intensity

						evi.Proteins[j].URazorLabels.Channel2.Name = i.Channel2.Name
						evi.Proteins[j].URazorLabels.Channel2.Mz = i.Channel2.Mz
						evi.Proteins[j].URazorLabels.Channel2.Intensity += i.Channel2.Intensity

						evi.Proteins[j].URazorLabels.Channel3.Name = i.Channel3.Name
						evi.Proteins[j].URazorLabels.Channel3.Mz = i.Channel3.Mz
						evi.Proteins[j].URazorLabels.Channel3.Intensity += i.Channel3.Intensity

						evi.Proteins[j].URazorLabels.Channel4.Name = i.Channel4.Name
						evi.Proteins[j].URazorLabels.Channel4.Mz = i.Channel4.Mz
						evi.Proteins[j].URazorLabels.Channel4.Intensity += i.Channel4.Intensity

						evi.Proteins[j].URazorLabels.Channel5.Name = i.Channel5.Name
						evi.Proteins[j].URazorLabels.Channel5.Mz = i.Channel5.Mz
						evi.Proteins[j].URazorLabels.Channel5.Intensity += i.Channel5.Intensity

						evi.Proteins[j].URazorLabels.Channel6.Name = i.Channel6.Name
						evi.Proteins[j].URazorLabels.Channel6.Mz = i.Channel6.Mz
						evi.Proteins[j].URazorLabels.Channel6.Intensity += i.Channel6.Intensity

						evi.Proteins[j].URazorLabels.Channel7.Name = i.Channel7.Name
						evi.Proteins[j].URazorLabels.Channel7.Mz = i.Channel7.Mz
						evi.Proteins[j].URazorLabels.Channel7.Intensity += i.Channel7.Intensity

						evi.Proteins[j].URazorLabels.Channel8.Name = i.Channel8.Name
						evi.Proteins[j].URazorLabels.Channel8.Mz = i.Channel8.Mz
						evi.Proteins[j].URazorLabels.Channel8.Intensity += i.Channel8.Intensity

						evi.Proteins[j].URazorLabels.Channel9.Name = i.Channel9.Name
						evi.Proteins[j].URazorLabels.Channel9.Mz = i.Channel9.Mz
						evi.Proteins[j].URazorLabels.Channel9.Intensity += i.Channel9.Intensity

						evi.Proteins[j].URazorLabels.Channel10.Name = i.Channel10.Name
						evi.Proteins[j].URazorLabels.Channel10.Mz = i.Channel10.Mz
						evi.Proteins[j].URazorLabels.Channel10.Intensity += i.Channel10.Intensity

						evi.Proteins[j].URazorLabels.Channel11.Name = i.Channel11.Name
						evi.Proteins[j].URazorLabels.Channel11.Mz = i.Channel11.Mz
						evi.Proteins[j].URazorLabels.Channel11.Intensity += i.Channel11.Intensity
					}

				}

			}
		}
	}

	return evi
}

// NormToTotalProteins calculates the protein level normalization based on total proteins
func NormToTotalProteins(evi rep.Evidence) rep.Evidence {

	var topValue float64
	var channelSum = [10]float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	var normFactors = [10]float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	// sum TMT singal for each column
	for _, i := range evi.Proteins {
		channelSum[0] += i.URazorLabels.Channel1.Intensity
		channelSum[1] += i.URazorLabels.Channel2.Intensity
		channelSum[2] += i.URazorLabels.Channel3.Intensity
		channelSum[3] += i.URazorLabels.Channel4.Intensity
		channelSum[4] += i.URazorLabels.Channel5.Intensity
		channelSum[5] += i.URazorLabels.Channel6.Intensity
		channelSum[6] += i.URazorLabels.Channel7.Intensity
		channelSum[7] += i.URazorLabels.Channel8.Intensity
		channelSum[8] += i.URazorLabels.Channel9.Intensity
		channelSum[9] += i.URazorLabels.Channel10.Intensity
	}

	// find the higest value amongst channels
	for _, i := range channelSum {
		if i > topValue {
			topValue = i
		}
	}

	// calculate normalizing factors
	for i := range channelSum {
		normFactors[i] = channelSum[i] / topValue
	}

	// multiply each protein TMT set by the factors to get normalized values
	for _, i := range evi.Proteins {
		i.URazorLabels.Channel1.Intensity *= normFactors[0]
		i.URazorLabels.Channel2.Intensity *= normFactors[1]
		i.URazorLabels.Channel3.Intensity *= normFactors[2]
		i.URazorLabels.Channel4.Intensity *= normFactors[3]
		i.URazorLabels.Channel5.Intensity *= normFactors[4]
		i.URazorLabels.Channel6.Intensity *= normFactors[5]
		i.URazorLabels.Channel7.Intensity *= normFactors[6]
		i.URazorLabels.Channel8.Intensity *= normFactors[7]
		i.URazorLabels.Channel9.Intensity *= normFactors[8]
		i.URazorLabels.Channel10.Intensity *= normFactors[9]
	}

	return evi
}

// func totalTop3LabelQuantification(evi rep.Evidence) (rep.Evidence, error) {
//
// 	for i := range evi.Proteins {
//
// 		p := make(PairList, len(evi.Proteins[i].TotalPeptideIons))
//
// 		j := 0
// 		for k, v := range evi.Proteins[i].TotalPeptideIons {
// 			p[j] = Pair{evi.Proteins[i].TotalPeptideIons[k], v.SummedLabelIntensity}
// 			j++
// 		}
//
// 		sort.Sort(sort.Reverse(p))
//
// 		var selectedIons []rep.IonEvidence
//
// 		var limit = 0
// 		if len(p) >= 3 {
// 			limit = 3
// 		} else if len(p) == 2 {
// 			limit = 2
// 		} else if len(p) == 1 {
// 			limit = 1
// 		}
//
// 		var counter = 0
// 		for _, j := range p {
// 			counter++
// 			if counter > limit {
// 				break
// 			}
// 			selectedIons = append(selectedIons, j.Key)
// 		}
//
// 		var c1Data float64
// 		var c2Data float64
// 		var c3Data float64
// 		var c4Data float64
// 		var c5Data float64
// 		var c6Data float64
// 		var c7Data float64
// 		var c8Data float64
// 		var c9Data float64
// 		var c10Data float64
//
// 		for _, j := range selectedIons {
// 			c1Data += j.Labels.Channel1.NormIntensity
// 			c2Data += j.Labels.Channel2.NormIntensity
// 			c3Data += j.Labels.Channel3.NormIntensity
// 			c4Data += j.Labels.Channel4.NormIntensity
// 			c5Data += j.Labels.Channel5.NormIntensity
// 			c6Data += j.Labels.Channel6.NormIntensity
// 			c7Data += j.Labels.Channel7.NormIntensity
// 			c8Data += j.Labels.Channel8.NormIntensity
// 			c9Data += j.Labels.Channel9.NormIntensity
// 			c10Data += j.Labels.Channel10.NormIntensity
// 		}
//
// 		evi.Proteins[i].TotalLabels.Channel1.TopIntensity = (c1Data / float64(limit))
// 		evi.Proteins[i].TotalLabels.Channel2.TopIntensity = (c2Data / float64(limit))
// 		evi.Proteins[i].TotalLabels.Channel3.TopIntensity = (c3Data / float64(limit))
// 		evi.Proteins[i].TotalLabels.Channel4.TopIntensity = (c4Data / float64(limit))
// 		evi.Proteins[i].TotalLabels.Channel5.TopIntensity = (c5Data / float64(limit))
// 		evi.Proteins[i].TotalLabels.Channel6.TopIntensity = (c6Data / float64(limit))
// 		evi.Proteins[i].TotalLabels.Channel7.TopIntensity = (c7Data / float64(limit))
// 		evi.Proteins[i].TotalLabels.Channel8.TopIntensity = (c8Data / float64(limit))
// 		evi.Proteins[i].TotalLabels.Channel9.TopIntensity = (c9Data / float64(limit))
// 		evi.Proteins[i].TotalLabels.Channel10.TopIntensity = (c10Data / float64(limit))
//
// 	}
//
// 	return evi, nil
// }

// // labelQuantificationOnTotalIons applies normalization to lable intensities
// func labelQuantificationOnTotalIons(evi rep.Evidence) (rep.Evidence, error) {
//
// 	for i := range evi.Proteins {
//
// 		var totalIons []rep.IonEvidence
// 		for _, v := range evi.Proteins[i].TotalPeptideIons {
// 			totalIons = append(totalIons, v)
// 		}
//
// 		var c1Data []float64
// 		var c2Data []float64
// 		var c3Data []float64
// 		var c4Data []float64
// 		var c5Data []float64
// 		var c6Data []float64
// 		var c7Data []float64
// 		var c8Data []float64
// 		var c9Data []float64
// 		var c10Data []float64
//
// 		// determine the mean and the standard deviation of the mean
// 		for j := range totalIons {
// 			c1Data = append(c1Data, totalIons[j].Labels.Channel1.NormIntensity)
// 			c2Data = append(c2Data, totalIons[j].Labels.Channel2.NormIntensity)
// 			c3Data = append(c3Data, totalIons[j].Labels.Channel3.NormIntensity)
// 			c4Data = append(c4Data, totalIons[j].Labels.Channel4.NormIntensity)
// 			c5Data = append(c5Data, totalIons[j].Labels.Channel5.NormIntensity)
// 			c6Data = append(c6Data, totalIons[j].Labels.Channel6.NormIntensity)
// 			c7Data = append(c7Data, totalIons[j].Labels.Channel7.NormIntensity)
// 			c8Data = append(c8Data, totalIons[j].Labels.Channel8.NormIntensity)
// 			c9Data = append(c9Data, totalIons[j].Labels.Channel9.NormIntensity)
// 			c10Data = append(c10Data, totalIons[j].Labels.Channel10.NormIntensity)
// 		}
//
// 		c1Mean, _ := stats.Mean(c1Data)
// 		c2Mean, _ := stats.Mean(c2Data)
// 		c3Mean, _ := stats.Mean(c3Data)
// 		c4Mean, _ := stats.Mean(c4Data)
// 		c5Mean, _ := stats.Mean(c5Data)
// 		c6Mean, _ := stats.Mean(c6Data)
// 		c7Mean, _ := stats.Mean(c7Data)
// 		c8Mean, _ := stats.Mean(c8Data)
// 		c9Mean, _ := stats.Mean(c9Data)
// 		c10Mean, _ := stats.Mean(c10Data)
// 		// if err != nil {
// 		// 	fmt.Println("AQUI")
// 		// 	return err
// 		// }
//
// 		c1StDev, _ := stats.StandardDeviationPopulation(c1Data)
// 		c2StDev, _ := stats.StandardDeviationPopulation(c2Data)
// 		c3StDev, _ := stats.StandardDeviationPopulation(c3Data)
// 		c4StDev, _ := stats.StandardDeviationPopulation(c4Data)
// 		c5StDev, _ := stats.StandardDeviationPopulation(c5Data)
// 		c6StDev, _ := stats.StandardDeviationPopulation(c6Data)
// 		c7StDev, _ := stats.StandardDeviationPopulation(c7Data)
// 		c8StDev, _ := stats.StandardDeviationPopulation(c8Data)
// 		c9StDev, _ := stats.StandardDeviationPopulation(c9Data)
// 		c10StDev, _ := stats.StandardDeviationPopulation(c10Data)
// 		// if err != nil {
// 		// 	return err
// 		// }
//
// 		// remov those that deviate from the mean by more than 2 sigma
// 		loC1Sigma := (c1Mean - 2*(c1StDev))
// 		hiC1Sigma := (c1Mean + 2*(c1StDev))
//
// 		loC2Sigma := (c2Mean - 2*(c2StDev))
// 		hiC2Sigma := (c2Mean + 2*(c2StDev))
//
// 		loC3Sigma := (c3Mean - 2*(c3StDev))
// 		hiC3Sigma := (c3Mean + 2*(c3StDev))
//
// 		loC4Sigma := (c4Mean - 2*(c4StDev))
// 		hiC4Sigma := (c4Mean + 2*(c4StDev))
//
// 		loC5Sigma := (c5Mean - 2*(c5StDev))
// 		hiC5Sigma := (c5Mean + 2*(c5StDev))
//
// 		loC6Sigma := (c6Mean - 2*(c6StDev))
// 		hiC6Sigma := (c6Mean + 2*(c6StDev))
//
// 		loC7Sigma := (c7Mean - 2*(c7StDev))
// 		hiC7Sigma := (c7Mean + 2*(c7StDev))
//
// 		loC8Sigma := (c8Mean - 2*(c8StDev))
// 		hiC8Sigma := (c8Mean + 2*(c8StDev))
//
// 		loC9Sigma := (c9Mean - 2*(c9StDev))
// 		hiC9Sigma := (c9Mean + 2*(c9StDev))
//
// 		loC10Sigma := (c10Mean - 2*(c10StDev))
// 		hiC10Sigma := (c10Mean + 2*(c10StDev))
//
// 		var normIons = make(map[int][]rep.IonEvidence)
//
// 		for i := range totalIons {
//
// 			if totalIons[i].Labels.Channel1.NormIntensity > 0 && totalIons[i].Labels.Channel1.NormIntensity >= loC1Sigma && totalIons[i].Labels.Channel1.NormIntensity <= hiC1Sigma {
// 				normIons[1] = append(normIons[1], totalIons[i])
// 			}
//
// 			if totalIons[i].Labels.Channel2.NormIntensity > 0 && totalIons[i].Labels.Channel2.NormIntensity >= loC2Sigma && totalIons[i].Labels.Channel2.NormIntensity <= hiC2Sigma {
// 				normIons[2] = append(normIons[2], totalIons[i])
// 			}
//
// 			if totalIons[i].Labels.Channel3.NormIntensity > 0 && totalIons[i].Labels.Channel3.NormIntensity >= loC3Sigma && totalIons[i].Labels.Channel3.NormIntensity <= hiC3Sigma {
// 				normIons[3] = append(normIons[3], totalIons[i])
// 			}
//
// 			if totalIons[i].Labels.Channel4.NormIntensity > 0 && totalIons[i].Labels.Channel4.NormIntensity >= loC4Sigma && totalIons[i].Labels.Channel4.NormIntensity <= hiC4Sigma {
// 				normIons[4] = append(normIons[4], totalIons[i])
// 			}
//
// 			if totalIons[i].Labels.Channel5.NormIntensity > 0 && totalIons[i].Labels.Channel5.NormIntensity >= loC5Sigma && totalIons[i].Labels.Channel5.NormIntensity <= hiC5Sigma {
// 				normIons[5] = append(normIons[5], totalIons[i])
// 			}
//
// 			if totalIons[i].Labels.Channel6.NormIntensity > 0 && totalIons[i].Labels.Channel6.NormIntensity >= loC6Sigma && totalIons[i].Labels.Channel6.NormIntensity <= hiC6Sigma {
// 				normIons[6] = append(normIons[6], totalIons[i])
// 			}
//
// 			if totalIons[i].Labels.Channel7.NormIntensity > 0 && totalIons[i].Labels.Channel7.NormIntensity >= loC7Sigma && totalIons[i].Labels.Channel7.NormIntensity <= hiC7Sigma {
// 				normIons[7] = append(normIons[7], totalIons[i])
// 			}
//
// 			if totalIons[i].Labels.Channel8.NormIntensity > 0 && totalIons[i].Labels.Channel8.NormIntensity >= loC8Sigma && totalIons[i].Labels.Channel8.NormIntensity <= hiC8Sigma {
// 				normIons[8] = append(normIons[8], totalIons[i])
// 			}
//
// 			if totalIons[i].Labels.Channel9.NormIntensity > 0 && totalIons[i].Labels.Channel9.NormIntensity >= loC9Sigma && totalIons[i].Labels.Channel9.NormIntensity <= hiC9Sigma {
// 				normIons[9] = append(normIons[9], totalIons[i])
// 			}
//
// 			if totalIons[i].Labels.Channel10.NormIntensity > 0 && totalIons[i].Labels.Channel10.NormIntensity >= loC10Sigma && totalIons[i].Labels.Channel10.NormIntensity <= hiC10Sigma {
// 				normIons[10] = append(normIons[10], totalIons[i])
// 			}
//
// 		}
//
// 		// recalculate the mean and standard deviation
// 		c1Data = nil
// 		c2Data = nil
// 		c3Data = nil
// 		c4Data = nil
// 		c5Data = nil
// 		c6Data = nil
// 		c7Data = nil
// 		c8Data = nil
// 		c9Data = nil
// 		c10Data = nil
//
// 		for _, v := range normIons[1] {
// 			c1Data = append(c1Data, v.Labels.Channel1.NormIntensity)
// 		}
//
// 		for _, v := range normIons[2] {
// 			c2Data = append(c2Data, v.Labels.Channel2.NormIntensity)
// 		}
//
// 		for _, v := range normIons[3] {
// 			c3Data = append(c3Data, v.Labels.Channel3.NormIntensity)
// 		}
//
// 		for _, v := range normIons[4] {
// 			c4Data = append(c4Data, v.Labels.Channel4.NormIntensity)
// 		}
//
// 		for _, v := range normIons[5] {
// 			c5Data = append(c5Data, v.Labels.Channel5.NormIntensity)
// 		}
//
// 		for _, v := range normIons[6] {
// 			c6Data = append(c6Data, v.Labels.Channel6.NormIntensity)
// 		}
//
// 		for _, v := range normIons[7] {
// 			c7Data = append(c7Data, v.Labels.Channel7.NormIntensity)
// 		}
//
// 		for _, v := range normIons[8] {
// 			c8Data = append(c8Data, v.Labels.Channel8.NormIntensity)
// 		}
//
// 		for _, v := range normIons[9] {
// 			c9Data = append(c9Data, v.Labels.Channel9.NormIntensity)
// 		}
//
// 		for _, v := range normIons[10] {
// 			c10Data = append(c10Data, v.Labels.Channel10.NormIntensity)
// 		}
//
// 		evi.Proteins[i].TotalLabels.Channel1.Mean, _ = stats.Mean(c1Data)
// 		evi.Proteins[i].TotalLabels.Channel2.Mean, _ = stats.Mean(c2Data)
// 		evi.Proteins[i].TotalLabels.Channel3.Mean, _ = stats.Mean(c3Data)
// 		evi.Proteins[i].TotalLabels.Channel4.Mean, _ = stats.Mean(c4Data)
// 		evi.Proteins[i].TotalLabels.Channel5.Mean, _ = stats.Mean(c5Data)
// 		evi.Proteins[i].TotalLabels.Channel6.Mean, _ = stats.Mean(c6Data)
// 		evi.Proteins[i].TotalLabels.Channel7.Mean, _ = stats.Mean(c7Data)
// 		evi.Proteins[i].TotalLabels.Channel8.Mean, _ = stats.Mean(c8Data)
// 		evi.Proteins[i].TotalLabels.Channel9.Mean, _ = stats.Mean(c9Data)
// 		evi.Proteins[i].TotalLabels.Channel10.Mean, _ = stats.Mean(c10Data)
//
// 		evi.Proteins[i].TotalLabels.Channel1.StDev, _ = stats.StandardDeviationPopulation(c1Data)
// 		evi.Proteins[i].TotalLabels.Channel2.StDev, _ = stats.StandardDeviationPopulation(c2Data)
// 		evi.Proteins[i].TotalLabels.Channel3.StDev, _ = stats.StandardDeviationPopulation(c3Data)
// 		evi.Proteins[i].TotalLabels.Channel4.StDev, _ = stats.StandardDeviationPopulation(c4Data)
// 		evi.Proteins[i].TotalLabels.Channel5.StDev, _ = stats.StandardDeviationPopulation(c5Data)
// 		evi.Proteins[i].TotalLabels.Channel6.StDev, _ = stats.StandardDeviationPopulation(c6Data)
// 		evi.Proteins[i].TotalLabels.Channel7.StDev, _ = stats.StandardDeviationPopulation(c7Data)
// 		evi.Proteins[i].TotalLabels.Channel8.StDev, _ = stats.StandardDeviationPopulation(c8Data)
// 		evi.Proteins[i].TotalLabels.Channel9.StDev, _ = stats.StandardDeviationPopulation(c9Data)
// 		evi.Proteins[i].TotalLabels.Channel10.StDev, _ = stats.StandardDeviationPopulation(c10Data)
//
// 	}
//
// 	return evi, nil
// }

// // labelQuantificationOnUniqueIons applies normalization to lable intensities
// func labelQuantificationOnUniqueIons(evi rep.Evidence) (rep.Evidence, error) {
//
// 	for i := range evi.Proteins {
//
// 		var ions []rep.IonEvidence
// 		for _, v := range evi.Proteins[i].TotalPeptideIons {
// 			if v.IsNondegenerateEvidence == true {
// 				ions = append(ions, v)
// 			}
// 		}
//
// 		var c1Data []float64
// 		var c2Data []float64
// 		var c3Data []float64
// 		var c4Data []float64
// 		var c5Data []float64
// 		var c6Data []float64
// 		var c7Data []float64
// 		var c8Data []float64
// 		var c9Data []float64
// 		var c10Data []float64
//
// 		// determine the mean and the standard deviation of the mean
// 		for i := range ions {
// 			c1Data = append(c1Data, ions[i].Labels.Channel1.NormIntensity)
// 			c2Data = append(c2Data, ions[i].Labels.Channel2.NormIntensity)
// 			c3Data = append(c3Data, ions[i].Labels.Channel3.NormIntensity)
// 			c4Data = append(c4Data, ions[i].Labels.Channel4.NormIntensity)
// 			c5Data = append(c5Data, ions[i].Labels.Channel5.NormIntensity)
// 			c6Data = append(c6Data, ions[i].Labels.Channel6.NormIntensity)
// 			c7Data = append(c7Data, ions[i].Labels.Channel7.NormIntensity)
// 			c8Data = append(c8Data, ions[i].Labels.Channel8.NormIntensity)
// 			c9Data = append(c9Data, ions[i].Labels.Channel9.NormIntensity)
// 			c10Data = append(c10Data, ions[i].Labels.Channel10.NormIntensity)
// 		}
//
// 		c1Mean, _ := stats.Mean(c1Data)
// 		c2Mean, _ := stats.Mean(c2Data)
// 		c3Mean, _ := stats.Mean(c3Data)
// 		c4Mean, _ := stats.Mean(c4Data)
// 		c5Mean, _ := stats.Mean(c5Data)
// 		c6Mean, _ := stats.Mean(c6Data)
// 		c7Mean, _ := stats.Mean(c7Data)
// 		c8Mean, _ := stats.Mean(c8Data)
// 		c9Mean, _ := stats.Mean(c9Data)
// 		c10Mean, _ := stats.Mean(c10Data)
// 		// if err != nil {
// 		// 	return err
// 		// }
//
// 		c1StDev, _ := stats.StandardDeviationPopulation(c1Data)
// 		c2StDev, _ := stats.StandardDeviationPopulation(c2Data)
// 		c3StDev, _ := stats.StandardDeviationPopulation(c3Data)
// 		c4StDev, _ := stats.StandardDeviationPopulation(c4Data)
// 		c5StDev, _ := stats.StandardDeviationPopulation(c5Data)
// 		c6StDev, _ := stats.StandardDeviationPopulation(c6Data)
// 		c7StDev, _ := stats.StandardDeviationPopulation(c7Data)
// 		c8StDev, _ := stats.StandardDeviationPopulation(c8Data)
// 		c9StDev, _ := stats.StandardDeviationPopulation(c9Data)
// 		c10StDev, _ := stats.StandardDeviationPopulation(c10Data)
// 		// if err != nil {
// 		// 	return err
// 		// }
//
// 		// remov those that deviate from the mean by more than 2 sigma
// 		loC1Sigma := (c1Mean - 2*(c1StDev))
// 		hiC1Sigma := (c1Mean + 2*(c1StDev))
//
// 		loC2Sigma := (c2Mean - 2*(c2StDev))
// 		hiC2Sigma := (c2Mean + 2*(c2StDev))
//
// 		loC3Sigma := (c3Mean - 2*(c3StDev))
// 		hiC3Sigma := (c3Mean + 2*(c3StDev))
//
// 		loC4Sigma := (c4Mean - 2*(c4StDev))
// 		hiC4Sigma := (c4Mean + 2*(c4StDev))
//
// 		loC5Sigma := (c5Mean - 2*(c5StDev))
// 		hiC5Sigma := (c5Mean + 2*(c5StDev))
//
// 		loC6Sigma := (c6Mean - 2*(c6StDev))
// 		hiC6Sigma := (c6Mean + 2*(c6StDev))
//
// 		loC7Sigma := (c7Mean - 2*(c7StDev))
// 		hiC7Sigma := (c7Mean + 2*(c7StDev))
//
// 		loC8Sigma := (c8Mean - 2*(c8StDev))
// 		hiC8Sigma := (c8Mean + 2*(c8StDev))
//
// 		loC9Sigma := (c9Mean - 2*(c9StDev))
// 		hiC9Sigma := (c9Mean + 2*(c9StDev))
//
// 		loC10Sigma := (c10Mean - 2*(c10StDev))
// 		hiC10Sigma := (c10Mean + 2*(c10StDev))
//
// 		var normIons = make(map[int][]rep.IonEvidence)
//
// 		for i := range ions {
//
// 			if ions[i].Labels.Channel1.NormIntensity > 0 && ions[i].Labels.Channel1.NormIntensity >= loC1Sigma && ions[i].Labels.Channel1.NormIntensity <= hiC1Sigma {
// 				normIons[1] = append(normIons[1], ions[i])
// 			}
//
// 			if ions[i].Labels.Channel2.NormIntensity > 0 && ions[i].Labels.Channel2.NormIntensity >= loC2Sigma && ions[i].Labels.Channel2.NormIntensity <= hiC2Sigma {
// 				normIons[2] = append(normIons[2], ions[i])
// 			}
//
// 			if ions[i].Labels.Channel3.NormIntensity > 0 && ions[i].Labels.Channel3.NormIntensity >= loC3Sigma && ions[i].Labels.Channel3.NormIntensity <= hiC3Sigma {
// 				normIons[3] = append(normIons[3], ions[i])
// 			}
//
// 			if ions[i].Labels.Channel4.NormIntensity > 0 && ions[i].Labels.Channel4.NormIntensity >= loC4Sigma && ions[i].Labels.Channel4.NormIntensity <= hiC4Sigma {
// 				normIons[4] = append(normIons[4], ions[i])
// 			}
//
// 			if ions[i].Labels.Channel5.NormIntensity > 0 && ions[i].Labels.Channel5.NormIntensity >= loC5Sigma && ions[i].Labels.Channel5.NormIntensity <= hiC5Sigma {
// 				normIons[5] = append(normIons[5], ions[i])
// 			}
//
// 			if ions[i].Labels.Channel6.NormIntensity > 0 && ions[i].Labels.Channel6.NormIntensity >= loC6Sigma && ions[i].Labels.Channel6.NormIntensity <= hiC6Sigma {
// 				normIons[6] = append(normIons[6], ions[i])
// 			}
//
// 			if ions[i].Labels.Channel7.NormIntensity > 0 && ions[i].Labels.Channel7.NormIntensity >= loC7Sigma && ions[i].Labels.Channel7.NormIntensity <= hiC7Sigma {
// 				normIons[7] = append(normIons[7], ions[i])
// 			}
//
// 			if ions[i].Labels.Channel8.NormIntensity > 0 && ions[i].Labels.Channel8.NormIntensity >= loC8Sigma && ions[i].Labels.Channel8.NormIntensity <= hiC8Sigma {
// 				normIons[8] = append(normIons[8], ions[i])
// 			}
//
// 			if ions[i].Labels.Channel9.NormIntensity > 0 && ions[i].Labels.Channel9.NormIntensity >= loC9Sigma && ions[i].Labels.Channel9.NormIntensity <= hiC9Sigma {
// 				normIons[9] = append(normIons[9], ions[i])
// 			}
//
// 			if ions[i].Labels.Channel10.NormIntensity > 0 && ions[i].Labels.Channel10.NormIntensity >= loC10Sigma && ions[i].Labels.Channel10.NormIntensity <= hiC10Sigma {
// 				normIons[10] = append(normIons[10], ions[i])
// 			}
//
// 		}
//
// 		// recalculate the mean and standard deviation
// 		c1Data = nil
// 		c2Data = nil
// 		c3Data = nil
// 		c4Data = nil
// 		c5Data = nil
// 		c6Data = nil
// 		c7Data = nil
// 		c8Data = nil
// 		c9Data = nil
// 		c10Data = nil
//
// 		for _, v := range normIons[1] {
// 			c1Data = append(c1Data, v.Labels.Channel1.NormIntensity)
// 		}
//
// 		for _, v := range normIons[2] {
// 			c2Data = append(c2Data, v.Labels.Channel2.NormIntensity)
// 		}
//
// 		for _, v := range normIons[3] {
// 			c3Data = append(c3Data, v.Labels.Channel3.NormIntensity)
// 		}
//
// 		for _, v := range normIons[4] {
// 			c4Data = append(c4Data, v.Labels.Channel4.NormIntensity)
// 		}
//
// 		for _, v := range normIons[5] {
// 			c5Data = append(c5Data, v.Labels.Channel5.NormIntensity)
// 		}
//
// 		for _, v := range normIons[6] {
// 			c6Data = append(c6Data, v.Labels.Channel6.NormIntensity)
// 		}
//
// 		for _, v := range normIons[7] {
// 			c7Data = append(c7Data, v.Labels.Channel7.NormIntensity)
// 		}
//
// 		for _, v := range normIons[8] {
// 			c8Data = append(c8Data, v.Labels.Channel8.NormIntensity)
// 		}
//
// 		for _, v := range normIons[9] {
// 			c9Data = append(c9Data, v.Labels.Channel9.NormIntensity)
// 		}
//
// 		for _, v := range normIons[10] {
// 			c10Data = append(c10Data, v.Labels.Channel10.NormIntensity)
// 		}
//
// 		evi.Proteins[i].UniqueLabels.Channel1.Mean, _ = stats.Mean(c1Data)
// 		evi.Proteins[i].UniqueLabels.Channel2.Mean, _ = stats.Mean(c2Data)
// 		evi.Proteins[i].UniqueLabels.Channel3.Mean, _ = stats.Mean(c3Data)
// 		evi.Proteins[i].UniqueLabels.Channel4.Mean, _ = stats.Mean(c4Data)
// 		evi.Proteins[i].UniqueLabels.Channel5.Mean, _ = stats.Mean(c5Data)
// 		evi.Proteins[i].UniqueLabels.Channel6.Mean, _ = stats.Mean(c6Data)
// 		evi.Proteins[i].UniqueLabels.Channel7.Mean, _ = stats.Mean(c7Data)
// 		evi.Proteins[i].UniqueLabels.Channel8.Mean, _ = stats.Mean(c8Data)
// 		evi.Proteins[i].UniqueLabels.Channel9.Mean, _ = stats.Mean(c9Data)
// 		evi.Proteins[i].UniqueLabels.Channel10.Mean, _ = stats.Mean(c10Data)
//
// 		evi.Proteins[i].UniqueLabels.Channel1.StDev, _ = stats.StandardDeviationPopulation(c1Data)
// 		evi.Proteins[i].UniqueLabels.Channel2.StDev, _ = stats.StandardDeviationPopulation(c2Data)
// 		evi.Proteins[i].UniqueLabels.Channel3.StDev, _ = stats.StandardDeviationPopulation(c3Data)
// 		evi.Proteins[i].UniqueLabels.Channel4.StDev, _ = stats.StandardDeviationPopulation(c4Data)
// 		evi.Proteins[i].UniqueLabels.Channel5.StDev, _ = stats.StandardDeviationPopulation(c5Data)
// 		evi.Proteins[i].UniqueLabels.Channel6.StDev, _ = stats.StandardDeviationPopulation(c6Data)
// 		evi.Proteins[i].UniqueLabels.Channel7.StDev, _ = stats.StandardDeviationPopulation(c7Data)
// 		evi.Proteins[i].UniqueLabels.Channel8.StDev, _ = stats.StandardDeviationPopulation(c8Data)
// 		evi.Proteins[i].UniqueLabels.Channel9.StDev, _ = stats.StandardDeviationPopulation(c9Data)
// 		evi.Proteins[i].UniqueLabels.Channel10.StDev, _ = stats.StandardDeviationPopulation(c10Data)
//
// 	}
//
// 	return evi, nil
// }

// // labelQuantificationOnUniqueIons applies normalization to lable intensities
// func labelQuantificationOnURazors(evi rep.Evidence) (rep.Evidence, error) {
//
// 	for i := range evi.Proteins {
//
// 		var ions []rep.IonEvidence
// 		for _, v := range evi.Proteins[i].TotalPeptideIons {
// 			if v.IsURazor == true {
// 				ions = append(ions, v)
// 			}
// 		}
//
// 		var c1Data []float64
// 		var c2Data []float64
// 		var c3Data []float64
// 		var c4Data []float64
// 		var c5Data []float64
// 		var c6Data []float64
// 		var c7Data []float64
// 		var c8Data []float64
// 		var c9Data []float64
// 		var c10Data []float64
//
// 		// determine the mean and the standard deviation of the mean
// 		for i := range ions {
// 			c1Data = append(c1Data, ions[i].Labels.Channel1.NormIntensity)
// 			c2Data = append(c2Data, ions[i].Labels.Channel2.NormIntensity)
// 			c3Data = append(c3Data, ions[i].Labels.Channel3.NormIntensity)
// 			c4Data = append(c4Data, ions[i].Labels.Channel4.NormIntensity)
// 			c5Data = append(c5Data, ions[i].Labels.Channel5.NormIntensity)
// 			c6Data = append(c6Data, ions[i].Labels.Channel6.NormIntensity)
// 			c7Data = append(c7Data, ions[i].Labels.Channel7.NormIntensity)
// 			c8Data = append(c8Data, ions[i].Labels.Channel8.NormIntensity)
// 			c9Data = append(c9Data, ions[i].Labels.Channel9.NormIntensity)
// 			c10Data = append(c10Data, ions[i].Labels.Channel10.NormIntensity)
// 		}
//
// 		c1Mean, _ := stats.Mean(c1Data)
// 		c2Mean, _ := stats.Mean(c2Data)
// 		c3Mean, _ := stats.Mean(c3Data)
// 		c4Mean, _ := stats.Mean(c4Data)
// 		c5Mean, _ := stats.Mean(c5Data)
// 		c6Mean, _ := stats.Mean(c6Data)
// 		c7Mean, _ := stats.Mean(c7Data)
// 		c8Mean, _ := stats.Mean(c8Data)
// 		c9Mean, _ := stats.Mean(c9Data)
// 		c10Mean, _ := stats.Mean(c10Data)
// 		// if err != nil {
// 		// 	return err
// 		// }
//
// 		c1StDev, _ := stats.StandardDeviationPopulation(c1Data)
// 		c2StDev, _ := stats.StandardDeviationPopulation(c2Data)
// 		c3StDev, _ := stats.StandardDeviationPopulation(c3Data)
// 		c4StDev, _ := stats.StandardDeviationPopulation(c4Data)
// 		c5StDev, _ := stats.StandardDeviationPopulation(c5Data)
// 		c6StDev, _ := stats.StandardDeviationPopulation(c6Data)
// 		c7StDev, _ := stats.StandardDeviationPopulation(c7Data)
// 		c8StDev, _ := stats.StandardDeviationPopulation(c8Data)
// 		c9StDev, _ := stats.StandardDeviationPopulation(c9Data)
// 		c10StDev, _ := stats.StandardDeviationPopulation(c10Data)
// 		// if err != nil {
// 		// 	return err
// 		// }
//
// 		// remov those that deviate from the mean by more than 2 sigma
// 		loC1Sigma := (c1Mean - 2*(c1StDev))
// 		hiC1Sigma := (c1Mean + 2*(c1StDev))
//
// 		loC2Sigma := (c2Mean - 2*(c2StDev))
// 		hiC2Sigma := (c2Mean + 2*(c2StDev))
//
// 		loC3Sigma := (c3Mean - 2*(c3StDev))
// 		hiC3Sigma := (c3Mean + 2*(c3StDev))
//
// 		loC4Sigma := (c4Mean - 2*(c4StDev))
// 		hiC4Sigma := (c4Mean + 2*(c4StDev))
//
// 		loC5Sigma := (c5Mean - 2*(c5StDev))
// 		hiC5Sigma := (c5Mean + 2*(c5StDev))
//
// 		loC6Sigma := (c6Mean - 2*(c6StDev))
// 		hiC6Sigma := (c6Mean + 2*(c6StDev))
//
// 		loC7Sigma := (c7Mean - 2*(c7StDev))
// 		hiC7Sigma := (c7Mean + 2*(c7StDev))
//
// 		loC8Sigma := (c8Mean - 2*(c8StDev))
// 		hiC8Sigma := (c8Mean + 2*(c8StDev))
//
// 		loC9Sigma := (c9Mean - 2*(c9StDev))
// 		hiC9Sigma := (c9Mean + 2*(c9StDev))
//
// 		loC10Sigma := (c10Mean - 2*(c10StDev))
// 		hiC10Sigma := (c10Mean + 2*(c10StDev))
//
// 		var normIons = make(map[int][]rep.IonEvidence)
//
// 		for i := range ions {
//
// 			if ions[i].Labels.Channel1.NormIntensity > 0 && ions[i].Labels.Channel1.NormIntensity >= loC1Sigma && ions[i].Labels.Channel1.NormIntensity <= hiC1Sigma {
// 				normIons[1] = append(normIons[1], ions[i])
// 			}
//
// 			if ions[i].Labels.Channel2.NormIntensity > 0 && ions[i].Labels.Channel2.NormIntensity >= loC2Sigma && ions[i].Labels.Channel2.NormIntensity <= hiC2Sigma {
// 				normIons[2] = append(normIons[2], ions[i])
// 			}
//
// 			if ions[i].Labels.Channel3.NormIntensity > 0 && ions[i].Labels.Channel3.NormIntensity >= loC3Sigma && ions[i].Labels.Channel3.NormIntensity <= hiC3Sigma {
// 				normIons[3] = append(normIons[3], ions[i])
// 			}
//
// 			if ions[i].Labels.Channel4.NormIntensity > 0 && ions[i].Labels.Channel4.NormIntensity >= loC4Sigma && ions[i].Labels.Channel4.NormIntensity <= hiC4Sigma {
// 				normIons[4] = append(normIons[4], ions[i])
// 			}
//
// 			if ions[i].Labels.Channel5.NormIntensity > 0 && ions[i].Labels.Channel5.NormIntensity >= loC5Sigma && ions[i].Labels.Channel5.NormIntensity <= hiC5Sigma {
// 				normIons[5] = append(normIons[5], ions[i])
// 			}
//
// 			if ions[i].Labels.Channel6.NormIntensity > 0 && ions[i].Labels.Channel6.NormIntensity >= loC6Sigma && ions[i].Labels.Channel6.NormIntensity <= hiC6Sigma {
// 				normIons[6] = append(normIons[6], ions[i])
// 			}
//
// 			if ions[i].Labels.Channel7.NormIntensity > 0 && ions[i].Labels.Channel7.NormIntensity >= loC7Sigma && ions[i].Labels.Channel7.NormIntensity <= hiC7Sigma {
// 				normIons[7] = append(normIons[7], ions[i])
// 			}
//
// 			if ions[i].Labels.Channel8.NormIntensity > 0 && ions[i].Labels.Channel8.NormIntensity >= loC8Sigma && ions[i].Labels.Channel8.NormIntensity <= hiC8Sigma {
// 				normIons[8] = append(normIons[8], ions[i])
// 			}
//
// 			if ions[i].Labels.Channel9.NormIntensity > 0 && ions[i].Labels.Channel9.NormIntensity >= loC9Sigma && ions[i].Labels.Channel9.NormIntensity <= hiC9Sigma {
// 				normIons[9] = append(normIons[9], ions[i])
// 			}
//
// 			if ions[i].Labels.Channel10.NormIntensity > 0 && ions[i].Labels.Channel10.NormIntensity >= loC10Sigma && ions[i].Labels.Channel10.NormIntensity <= hiC10Sigma {
// 				normIons[10] = append(normIons[10], ions[i])
// 			}
//
// 		}
//
// 		// recalculate the mean and standard deviation
// 		c1Data = nil
// 		c2Data = nil
// 		c3Data = nil
// 		c4Data = nil
// 		c5Data = nil
// 		c6Data = nil
// 		c7Data = nil
// 		c8Data = nil
// 		c9Data = nil
// 		c10Data = nil
//
// 		for _, v := range normIons[1] {
// 			c1Data = append(c1Data, v.Labels.Channel1.NormIntensity)
// 		}
//
// 		for _, v := range normIons[2] {
// 			c2Data = append(c2Data, v.Labels.Channel2.NormIntensity)
// 		}
//
// 		for _, v := range normIons[3] {
// 			c3Data = append(c3Data, v.Labels.Channel3.NormIntensity)
// 		}
//
// 		for _, v := range normIons[4] {
// 			c4Data = append(c4Data, v.Labels.Channel4.NormIntensity)
// 		}
//
// 		for _, v := range normIons[5] {
// 			c5Data = append(c5Data, v.Labels.Channel5.NormIntensity)
// 		}
//
// 		for _, v := range normIons[6] {
// 			c6Data = append(c6Data, v.Labels.Channel6.NormIntensity)
// 		}
//
// 		for _, v := range normIons[7] {
// 			c7Data = append(c7Data, v.Labels.Channel7.NormIntensity)
// 		}
//
// 		for _, v := range normIons[8] {
// 			c8Data = append(c8Data, v.Labels.Channel8.NormIntensity)
// 		}
//
// 		for _, v := range normIons[9] {
// 			c9Data = append(c9Data, v.Labels.Channel9.NormIntensity)
// 		}
//
// 		for _, v := range normIons[10] {
// 			c10Data = append(c10Data, v.Labels.Channel10.NormIntensity)
// 		}
//
// 		evi.Proteins[i].URazorLabels.Channel1.Mean, _ = stats.Mean(c1Data)
// 		evi.Proteins[i].URazorLabels.Channel2.Mean, _ = stats.Mean(c2Data)
// 		evi.Proteins[i].URazorLabels.Channel3.Mean, _ = stats.Mean(c3Data)
// 		evi.Proteins[i].URazorLabels.Channel4.Mean, _ = stats.Mean(c4Data)
// 		evi.Proteins[i].URazorLabels.Channel5.Mean, _ = stats.Mean(c5Data)
// 		evi.Proteins[i].URazorLabels.Channel6.Mean, _ = stats.Mean(c6Data)
// 		evi.Proteins[i].URazorLabels.Channel7.Mean, _ = stats.Mean(c7Data)
// 		evi.Proteins[i].URazorLabels.Channel8.Mean, _ = stats.Mean(c8Data)
// 		evi.Proteins[i].URazorLabels.Channel9.Mean, _ = stats.Mean(c9Data)
// 		evi.Proteins[i].URazorLabels.Channel10.Mean, _ = stats.Mean(c10Data)
//
// 		evi.Proteins[i].URazorLabels.Channel1.StDev, _ = stats.StandardDeviationPopulation(c1Data)
// 		evi.Proteins[i].URazorLabels.Channel2.StDev, _ = stats.StandardDeviationPopulation(c2Data)
// 		evi.Proteins[i].URazorLabels.Channel3.StDev, _ = stats.StandardDeviationPopulation(c3Data)
// 		evi.Proteins[i].URazorLabels.Channel4.StDev, _ = stats.StandardDeviationPopulation(c4Data)
// 		evi.Proteins[i].URazorLabels.Channel5.StDev, _ = stats.StandardDeviationPopulation(c5Data)
// 		evi.Proteins[i].URazorLabels.Channel6.StDev, _ = stats.StandardDeviationPopulation(c6Data)
// 		evi.Proteins[i].URazorLabels.Channel7.StDev, _ = stats.StandardDeviationPopulation(c7Data)
// 		evi.Proteins[i].URazorLabels.Channel8.StDev, _ = stats.StandardDeviationPopulation(c8Data)
// 		evi.Proteins[i].URazorLabels.Channel9.StDev, _ = stats.StandardDeviationPopulation(c9Data)
// 		evi.Proteins[i].URazorLabels.Channel10.StDev, _ = stats.StandardDeviationPopulation(c10Data)
//
// 	}
//
// 	return evi, nil
// }

// func totalTop3LabelQuantification(evi rep.Evidence) (rep.Evidence, error) {
//
// 	for i := range evi.Proteins {
//
// 		var pairlist PairList
//
// 		for _, v := range evi.Proteins[i].TotalPeptideIons {
// 			var pair Pair
// 			pair.Key = v
// 			pair.Value = v.SummedLabelIntensity
// 			pairlist = append(pairlist, pair)
// 		}
//
// 		sort.Sort(sort.Reverse(pairlist))
//
// 		var selectedIons []rep.IonEvidence
//
// 		var limit = 0
// 		if len(pairlist) >= 3 {
// 			limit = 3
// 		} else if len(pairlist) == 2 {
// 			limit = 2
// 		} else if len(pairlist) == 1 {
// 			limit = 1
// 		}
//
// 		var counter = 0
// 		for _, j := range pairlist {
// 			counter++
// 			if counter > limit {
// 				break
// 			}
// 			selectedIons = append(selectedIons, j.Key)
// 		}
//
// 		var c1Data float64
// 		var c2Data float64
// 		var c3Data float64
// 		var c4Data float64
// 		var c5Data float64
// 		var c6Data float64
// 		var c7Data float64
// 		var c8Data float64
// 		var c9Data float64
// 		var c10Data float64
//
// 		// determine the mean and the standard deviation of the mean
// 		for _, j := range selectedIons {
// 			c1Data += j.Labels.Channel1.Intensity
// 			c2Data += j.Labels.Channel2.Intensity
// 			c3Data += j.Labels.Channel3.Intensity
// 			c4Data += j.Labels.Channel4.Intensity
// 			c5Data += j.Labels.Channel5.Intensity
// 			c6Data += j.Labels.Channel6.Intensity
// 			c7Data += j.Labels.Channel7.Intensity
// 			c8Data += j.Labels.Channel8.Intensity
// 			c9Data += j.Labels.Channel9.Intensity
// 			c10Data += j.Labels.Channel10.Intensity
// 		}
//
// 		evi.Proteins[i].TotalLabels.Channel1.TopIntensity = (c1Data / float64(limit))
// 		evi.Proteins[i].TotalLabels.Channel2.TopIntensity = (c2Data / float64(limit))
// 		evi.Proteins[i].TotalLabels.Channel3.TopIntensity = (c3Data / float64(limit))
// 		evi.Proteins[i].TotalLabels.Channel4.TopIntensity = (c4Data / float64(limit))
// 		evi.Proteins[i].TotalLabels.Channel5.TopIntensity = (c5Data / float64(limit))
// 		evi.Proteins[i].TotalLabels.Channel6.TopIntensity = (c6Data / float64(limit))
// 		evi.Proteins[i].TotalLabels.Channel7.TopIntensity = (c7Data / float64(limit))
// 		evi.Proteins[i].TotalLabels.Channel8.TopIntensity = (c8Data / float64(limit))
// 		evi.Proteins[i].TotalLabels.Channel9.TopIntensity = (c9Data / float64(limit))
// 		evi.Proteins[i].TotalLabels.Channel10.TopIntensity = (c10Data / float64(limit))
//
// 	}
//
// 	return evi, nil
// }

// func ratioToIntensityMean(evi rep.Evidence) (rep.Evidence, error) {
//
// 	for i := range evi.Proteins {
//
// 		var totalRef float64
// 		var uniqRef float64
// 		var razorRef float64
//
// 		totalRef += evi.Proteins[i].TotalLabels.Channel1.Mean
// 		totalRef += evi.Proteins[i].TotalLabels.Channel2.Mean
// 		totalRef += evi.Proteins[i].TotalLabels.Channel3.Mean
// 		totalRef += evi.Proteins[i].TotalLabels.Channel4.Mean
// 		totalRef += evi.Proteins[i].TotalLabels.Channel5.Mean
// 		totalRef += evi.Proteins[i].TotalLabels.Channel6.Mean
// 		totalRef += evi.Proteins[i].TotalLabels.Channel7.Mean
// 		totalRef += evi.Proteins[i].TotalLabels.Channel8.Mean
// 		totalRef += evi.Proteins[i].TotalLabels.Channel9.Mean
// 		totalRef += evi.Proteins[i].TotalLabels.Channel10.Mean
//
// 		uniqRef += evi.Proteins[i].UniqueLabels.Channel1.Mean
// 		uniqRef += evi.Proteins[i].UniqueLabels.Channel2.Mean
// 		uniqRef += evi.Proteins[i].UniqueLabels.Channel3.Mean
// 		uniqRef += evi.Proteins[i].UniqueLabels.Channel4.Mean
// 		uniqRef += evi.Proteins[i].UniqueLabels.Channel5.Mean
// 		uniqRef += evi.Proteins[i].UniqueLabels.Channel6.Mean
// 		uniqRef += evi.Proteins[i].UniqueLabels.Channel7.Mean
// 		uniqRef += evi.Proteins[i].UniqueLabels.Channel8.Mean
// 		uniqRef += evi.Proteins[i].UniqueLabels.Channel9.Mean
// 		uniqRef += evi.Proteins[i].UniqueLabels.Channel10.Mean
//
// 		razorRef += evi.Proteins[i].URazorLabels.Channel1.Mean
// 		razorRef += evi.Proteins[i].URazorLabels.Channel2.Mean
// 		razorRef += evi.Proteins[i].URazorLabels.Channel3.Mean
// 		razorRef += evi.Proteins[i].URazorLabels.Channel4.Mean
// 		razorRef += evi.Proteins[i].URazorLabels.Channel5.Mean
// 		razorRef += evi.Proteins[i].URazorLabels.Channel6.Mean
// 		razorRef += evi.Proteins[i].URazorLabels.Channel7.Mean
// 		razorRef += evi.Proteins[i].URazorLabels.Channel8.Mean
// 		razorRef += evi.Proteins[i].URazorLabels.Channel9.Mean
// 		razorRef += evi.Proteins[i].URazorLabels.Channel10.Mean
//
// 		evi.Proteins[i].TotalLabels.Channel1.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel1.Mean/totalRef), 4, 5) * 100)
// 		evi.Proteins[i].TotalLabels.Channel2.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel2.Mean/totalRef), 4, 5) * 100)
// 		evi.Proteins[i].TotalLabels.Channel3.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel3.Mean/totalRef), 4, 5) * 100)
// 		evi.Proteins[i].TotalLabels.Channel4.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel4.Mean/totalRef), 4, 5) * 100)
// 		evi.Proteins[i].TotalLabels.Channel5.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel5.Mean/totalRef), 4, 5) * 100)
// 		evi.Proteins[i].TotalLabels.Channel6.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel6.Mean/totalRef), 4, 5) * 100)
// 		evi.Proteins[i].TotalLabels.Channel7.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel7.Mean/totalRef), 4, 5) * 100)
// 		evi.Proteins[i].TotalLabels.Channel8.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel8.Mean/totalRef), 4, 5) * 100)
// 		evi.Proteins[i].TotalLabels.Channel9.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel9.Mean/totalRef), 4, 5) * 100)
// 		evi.Proteins[i].TotalLabels.Channel10.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel10.Mean/totalRef), 4, 5) * 100)
//
// 		evi.Proteins[i].UniqueLabels.Channel1.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel1.Mean/uniqRef), 4, 5) * 100)
// 		evi.Proteins[i].UniqueLabels.Channel2.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel2.Mean/uniqRef), 4, 5) * 100)
// 		evi.Proteins[i].UniqueLabels.Channel3.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel3.Mean/uniqRef), 4, 5) * 100)
// 		evi.Proteins[i].UniqueLabels.Channel4.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel4.Mean/uniqRef), 4, 5) * 100)
// 		evi.Proteins[i].UniqueLabels.Channel5.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel5.Mean/uniqRef), 4, 5) * 100)
// 		evi.Proteins[i].UniqueLabels.Channel6.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel6.Mean/uniqRef), 4, 5) * 100)
// 		evi.Proteins[i].UniqueLabels.Channel7.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel7.Mean/uniqRef), 4, 5) * 100)
// 		evi.Proteins[i].UniqueLabels.Channel8.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel8.Mean/uniqRef), 4, 5) * 100)
// 		evi.Proteins[i].UniqueLabels.Channel9.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel9.Mean/uniqRef), 4, 5) * 100)
// 		evi.Proteins[i].UniqueLabels.Channel10.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel10.Mean/uniqRef), 4, 5) * 100)
//
// 		evi.Proteins[i].URazorLabels.Channel1.RatioIntensity = (utils.Round((evi.Proteins[i].URazorLabels.Channel1.Mean/razorRef), 4, 5) * 100)
// 		evi.Proteins[i].URazorLabels.Channel2.RatioIntensity = (utils.Round((evi.Proteins[i].URazorLabels.Channel2.Mean/razorRef), 4, 5) * 100)
// 		evi.Proteins[i].URazorLabels.Channel3.RatioIntensity = (utils.Round((evi.Proteins[i].URazorLabels.Channel3.Mean/razorRef), 4, 5) * 100)
// 		evi.Proteins[i].URazorLabels.Channel4.RatioIntensity = (utils.Round((evi.Proteins[i].URazorLabels.Channel4.Mean/razorRef), 4, 5) * 100)
// 		evi.Proteins[i].URazorLabels.Channel5.RatioIntensity = (utils.Round((evi.Proteins[i].URazorLabels.Channel5.Mean/razorRef), 4, 5) * 100)
// 		evi.Proteins[i].URazorLabels.Channel6.RatioIntensity = (utils.Round((evi.Proteins[i].URazorLabels.Channel6.Mean/razorRef), 4, 5) * 100)
// 		evi.Proteins[i].URazorLabels.Channel7.RatioIntensity = (utils.Round((evi.Proteins[i].URazorLabels.Channel7.Mean/razorRef), 4, 5) * 100)
// 		evi.Proteins[i].URazorLabels.Channel8.RatioIntensity = (utils.Round((evi.Proteins[i].URazorLabels.Channel8.Mean/razorRef), 4, 5) * 100)
// 		evi.Proteins[i].URazorLabels.Channel9.RatioIntensity = (utils.Round((evi.Proteins[i].URazorLabels.Channel9.Mean/razorRef), 4, 5) * 100)
// 		evi.Proteins[i].URazorLabels.Channel10.RatioIntensity = (utils.Round((evi.Proteins[i].URazorLabels.Channel10.Mean/razorRef), 4, 5) * 100)
//
// 	}
//
// 	return evi, nil
// }

// func ratioToControlChannel(evi rep.Evidence, control string) (rep.Evidence, error) {
//
// 	for i := range evi.Proteins {
//
// 		var totalRef float64
// 		var uniqRef float64
// 		var razorRef float64
//
// 		switch control {
// 		case "1":
// 			totalRef = evi.Proteins[i].TotalLabels.Channel1.Mean
// 			uniqRef = evi.Proteins[i].UniqueLabels.Channel1.Mean
// 			razorRef = evi.Proteins[i].URazorLabels.Channel1.Mean
// 		case "2":
// 			totalRef = evi.Proteins[i].TotalLabels.Channel2.Mean
// 			uniqRef = evi.Proteins[i].UniqueLabels.Channel2.Mean
// 			razorRef = evi.Proteins[i].URazorLabels.Channel2.Mean
// 		case "3":
// 			totalRef = evi.Proteins[i].TotalLabels.Channel3.Mean
// 			uniqRef = evi.Proteins[i].UniqueLabels.Channel3.Mean
// 			razorRef = evi.Proteins[i].URazorLabels.Channel3.Mean
// 		case "4":
// 			totalRef = evi.Proteins[i].TotalLabels.Channel4.Mean
// 			uniqRef = evi.Proteins[i].UniqueLabels.Channel4.Mean
// 			razorRef = evi.Proteins[i].URazorLabels.Channel4.Mean
// 		case "5":
// 			totalRef = evi.Proteins[i].TotalLabels.Channel5.Mean
// 			uniqRef = evi.Proteins[i].UniqueLabels.Channel5.Mean
// 			razorRef = evi.Proteins[i].URazorLabels.Channel5.Mean
// 		case "6":
// 			totalRef = evi.Proteins[i].TotalLabels.Channel6.Mean
// 			uniqRef = evi.Proteins[i].UniqueLabels.Channel6.Mean
// 			razorRef = evi.Proteins[i].URazorLabels.Channel6.Mean
// 		case "7":
// 			totalRef = evi.Proteins[i].TotalLabels.Channel7.Mean
// 			uniqRef = evi.Proteins[i].UniqueLabels.Channel7.Mean
// 			razorRef = evi.Proteins[i].URazorLabels.Channel7.Mean
// 		case "8":
// 			totalRef = evi.Proteins[i].TotalLabels.Channel8.Mean
// 			uniqRef = evi.Proteins[i].UniqueLabels.Channel8.Mean
// 			razorRef = evi.Proteins[i].URazorLabels.Channel8.Mean
// 		case "9":
// 			totalRef = evi.Proteins[i].TotalLabels.Channel9.Mean
// 			uniqRef = evi.Proteins[i].UniqueLabels.Channel9.Mean
// 			razorRef = evi.Proteins[i].URazorLabels.Channel9.Mean
// 		case "10":
// 			totalRef = evi.Proteins[i].TotalLabels.Channel10.Mean
// 			uniqRef = evi.Proteins[i].UniqueLabels.Channel10.Mean
// 			razorRef = evi.Proteins[i].URazorLabels.Channel10.Mean
// 		default:
// 			return evi, errors.New("Cant find the given channel for normalization")
// 		}
//
// 		evi.Proteins[i].TotalLabels.Channel1.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel1.Mean/totalRef), 4, 5) * 100)
// 		evi.Proteins[i].TotalLabels.Channel2.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel2.Mean/totalRef), 4, 5) * 100)
// 		evi.Proteins[i].TotalLabels.Channel3.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel3.Mean/totalRef), 4, 5) * 100)
// 		evi.Proteins[i].TotalLabels.Channel4.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel4.Mean/totalRef), 4, 5) * 100)
// 		evi.Proteins[i].TotalLabels.Channel5.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel5.Mean/totalRef), 4, 5) * 100)
// 		evi.Proteins[i].TotalLabels.Channel6.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel6.Mean/totalRef), 4, 5) * 100)
// 		evi.Proteins[i].TotalLabels.Channel7.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel7.Mean/totalRef), 4, 5) * 100)
// 		evi.Proteins[i].TotalLabels.Channel8.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel8.Mean/totalRef), 4, 5) * 100)
// 		evi.Proteins[i].TotalLabels.Channel9.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel9.Mean/totalRef), 4, 5) * 100)
// 		evi.Proteins[i].TotalLabels.Channel10.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel10.Mean/totalRef), 4, 5) * 100)
//
// 		evi.Proteins[i].UniqueLabels.Channel1.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel1.Mean/uniqRef), 4, 5) * 100)
// 		evi.Proteins[i].UniqueLabels.Channel2.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel2.Mean/uniqRef), 4, 5) * 100)
// 		evi.Proteins[i].UniqueLabels.Channel3.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel3.Mean/uniqRef), 4, 5) * 100)
// 		evi.Proteins[i].UniqueLabels.Channel4.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel4.Mean/uniqRef), 4, 5) * 100)
// 		evi.Proteins[i].UniqueLabels.Channel5.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel5.Mean/uniqRef), 4, 5) * 100)
// 		evi.Proteins[i].UniqueLabels.Channel6.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel6.Mean/uniqRef), 4, 5) * 100)
// 		evi.Proteins[i].UniqueLabels.Channel7.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel7.Mean/uniqRef), 4, 5) * 100)
// 		evi.Proteins[i].UniqueLabels.Channel8.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel8.Mean/uniqRef), 4, 5) * 100)
// 		evi.Proteins[i].UniqueLabels.Channel9.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel9.Mean/uniqRef), 4, 5) * 100)
// 		evi.Proteins[i].UniqueLabels.Channel10.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel10.Mean/uniqRef), 4, 5) * 100)
//
// 		evi.Proteins[i].URazorLabels.Channel1.RatioIntensity = (utils.Round((evi.Proteins[i].URazorLabels.Channel1.Mean/razorRef), 4, 5) * 100)
// 		evi.Proteins[i].URazorLabels.Channel2.RatioIntensity = (utils.Round((evi.Proteins[i].URazorLabels.Channel2.Mean/razorRef), 4, 5) * 100)
// 		evi.Proteins[i].URazorLabels.Channel3.RatioIntensity = (utils.Round((evi.Proteins[i].URazorLabels.Channel3.Mean/razorRef), 4, 5) * 100)
// 		evi.Proteins[i].URazorLabels.Channel4.RatioIntensity = (utils.Round((evi.Proteins[i].URazorLabels.Channel4.Mean/razorRef), 4, 5) * 100)
// 		evi.Proteins[i].URazorLabels.Channel5.RatioIntensity = (utils.Round((evi.Proteins[i].URazorLabels.Channel5.Mean/razorRef), 4, 5) * 100)
// 		evi.Proteins[i].URazorLabels.Channel6.RatioIntensity = (utils.Round((evi.Proteins[i].URazorLabels.Channel6.Mean/razorRef), 4, 5) * 100)
// 		evi.Proteins[i].URazorLabels.Channel7.RatioIntensity = (utils.Round((evi.Proteins[i].URazorLabels.Channel7.Mean/razorRef), 4, 5) * 100)
// 		evi.Proteins[i].URazorLabels.Channel8.RatioIntensity = (utils.Round((evi.Proteins[i].URazorLabels.Channel8.Mean/razorRef), 4, 5) * 100)
// 		evi.Proteins[i].URazorLabels.Channel9.RatioIntensity = (utils.Round((evi.Proteins[i].URazorLabels.Channel9.Mean/razorRef), 4, 5) * 100)
// 		evi.Proteins[i].URazorLabels.Channel10.RatioIntensity = (utils.Round((evi.Proteins[i].URazorLabels.Channel10.Mean/razorRef), 4, 5) * 100)
//
// 	}
//
// 	return evi, nil
// }
