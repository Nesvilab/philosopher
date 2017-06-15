package quan

import (
	"errors"
	"fmt"
	"math"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/montanaflynn/stats"
	"github.com/prvst/cmsl/bio"
	"github.com/prvst/cmsl/data/mz"
	"github.com/prvst/cmsl/data/mz/mzml"
	"github.com/prvst/cmsl/data/mz/mzxml"
	"github.com/prvst/cmsl/utils"
	"github.com/prvst/philosopher/lib/meta"
	"github.com/prvst/philosopher/lib/rep"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/prvst/philosopher/lib/tmt"
)

// Quantify ...
type Quantify struct {
	meta.Data
	Phi      string
	Format   string
	Dir      string
	Brand    string
	Plex     string
	ChanNorm string
	RTWin    float64
	PTWin    float64
	Tol      float64
	Purity   float64
	IntNorm  bool
}

// Pair struct
type Pair struct {
	Key   rep.IonEvidence
	Value float64
}

// PairList struict
type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// New constructor
func New() Quantify {

	var o Quantify
	var m meta.Data
	m.Restore(sys.Meta())

	o.UUID = m.UUID
	o.Distro = m.Distro
	o.Home = m.Home
	o.MetaFile = m.MetaFile
	o.MetaDir = m.MetaDir
	o.DB = m.DB
	o.Temp = m.Temp
	o.TimeStamp = m.TimeStamp
	o.OS = m.OS
	o.Arch = m.Arch

	return o
}

// RunLabelFreeQuantification is the top function for label free quantification
func (p *Quantify) RunLabelFreeQuantification() error {

	var err error

	var evi rep.Evidence
	evi.Restore()

	if len(evi.Proteins) < 1 {
		logrus.Fatal("This restult file does not contains report data")
	}

	logrus.Info("Calculating Spectral Counts")
	evi, err = CalculateSpectralCounts(evi)
	if err != nil {
		return err
	}

	logrus.Info("Calculating MS1 Intensities")
	evi, err = peakIntensity(evi, p.Dir, p.Format, p.RTWin, p.PTWin, p.Tol)
	if err != nil {
		return err
	}

	evi, err = calculateIntensities(evi)
	if err != nil {
		return err
	}

	evi.Serialize()

	return nil
}

// RunLabeledQuantification is the top function for label quantification
func (p *Quantify) RunLabeledQuantification() error {

	var evi rep.Evidence
	evi.Restore()

	// removed all calculated defined bvalues from before
	cleanPreviousData(p.Plex)

	logrus.Info("Reading spectra files")
	ms1, ms2, err := getSpectra(p.Dir, p.Format, evi)
	if err != nil {
		return err
	}

	logrus.Info("Calculating ion purity")
	evi, err = calculateIonPurity(p.Dir, p.Format, ms1, ms2, evi)
	if err != nil {
		return err
	}

	logrus.Info("Calculating label intensities")
	labels, err := labeledPeakIntensity(p.Dir, p.Format, strings.ToLower(p.Brand), p.Plex, p.Tol, evi, ms2)
	if err != nil {
		return err
	}

	logrus.Info("Mapping quantification values")
	evi, err = mapLabeledSpectra(labels, p.Purity, evi)
	if err != nil {
		return err
	}

	// logrus.Info("Normalizing peptide channels")
	// evi, err = totalTop3LabelQuantification(evi)
	// if err != nil {
	// 	return err
	// }

	evi, err = labelQuantificationOnTotalIons(evi)
	if err != nil {
		return err
	}

	evi, err = labelQuantificationOnUniqueIons(evi)
	if err != nil {
		return err
	}

	evi, err = labelQuantificationOnURazors(evi)
	if err != nil {
		return err
	}

	if len(p.ChanNorm) > 0 {
		logrus.Info("Applying normalization to control channel ", p.ChanNorm)
		evi, err = ratioToControlChannel(evi, p.ChanNorm)
		if err != nil {
			return err
		}
	}

	if p.IntNorm == true {
		logrus.Info("Applying normalization to intensity means")
		evi, err = ratioToIntensityMean(evi)
		if err != nil {
			return err
		}
	}

	evi.Serialize()

	return nil
}

// cleanPreviousData cleans previous label quantifications
func cleanPreviousData(plex string) error {

	var err error
	var evi rep.Evidence
	evi.Restore()

	for i := range evi.PSM {
		evi.PSM[i].Labels, err = tmt.New(plex)
		if err != nil {
			return err
		}
	}

	for i := range evi.Ions {
		evi.Ions[i].Labels, err = tmt.New(plex)
		if err != nil {
			return err
		}
	}

	for i := range evi.Proteins {
		evi.Proteins[i].TotalLabels, err = tmt.New(plex)
		if err != nil {
			return err
		}
		evi.Proteins[i].UniqueLabels, err = tmt.New(plex)
		if err != nil {
			return err
		}
		evi.Proteins[i].RazorLabels, err = tmt.New(plex)
		if err != nil {
			return err
		}
	}

	evi.Serialize()

	return nil
}

// labeledPeakIntensity ...
func labeledPeakIntensity(dir, format, brand, plex string, tol float64, evi rep.Evidence, ms2 map[string]mz.MS2) (map[string]tmt.Labels, error) {

	// get all spectra names from PSMs and create the label list
	var spectra = make(map[string]tmt.Labels)

	for _, i := range evi.PSM {

		ls, err := tmt.New(plex)
		if err != nil {
			return spectra, err
		}

		// remove the charge state from the spectrum name key
		split := strings.Split(i.Spectrum, ".")
		name := fmt.Sprintf("%s.%s.%s", split[0], split[1], split[2])

		ls.Spectrum = i.Spectrum
		ls.RetentionTime = i.RetentionTime

		if format == "mzML" {
			index, err := strconv.Atoi(split[1])
			if err != nil {
				return spectra, err
			}
			ls.Index = (uint32(index) - 1)
		} else {
			index, err := strconv.Atoi(split[1])
			if err != nil {
				return spectra, err
			}
			ls.Index = uint32(index)
		}

		spectra[name] = ls
	}

	ppmPrecision := tol / math.Pow(10, 6)

	spectra = getLabels(spectra, ms2, ppmPrecision)

	return spectra, nil
}

// getLabels extract ion chomatograms
func getLabels(spec map[string]tmt.Labels, ms2 map[string]mz.MS2, ppmPrecision float64) map[string]tmt.Labels {

	// for each ms2 data from a different file
	for _, v := range ms2 {
		// for each ms2 spectra
		for _, i := range v.Ms2Scan {

			split := strings.Split(i.SpectrumName, ".")
			name := fmt.Sprintf("%s.%s.%s", split[0], split[1], split[2])

			sp, ok := spec[name]
			if ok {

				v := sp

				for _, k := range i.Spectrum {

					if k.Mz <= (v.Channel1.Mz+(ppmPrecision*v.Channel1.Mz)) && k.Mz >= (v.Channel1.Mz-(ppmPrecision*v.Channel1.Mz)) {
						if k.Intensity > v.Channel1.Intensity {
							v.Channel1.Intensity = k.Intensity
						}
					}

					if k.Mz <= (v.Channel2.Mz+(ppmPrecision*v.Channel2.Mz)) && k.Mz >= (v.Channel2.Mz-(ppmPrecision*v.Channel2.Mz)) {
						if k.Intensity > v.Channel2.Intensity {
							v.Channel2.Intensity = k.Intensity
						}
					}

					if k.Mz <= (v.Channel3.Mz+(ppmPrecision*v.Channel3.Mz)) && k.Mz >= (v.Channel3.Mz-(ppmPrecision*v.Channel3.Mz)) {
						if k.Intensity > v.Channel3.Intensity {
							v.Channel3.Intensity = k.Intensity
						}
					}

					if k.Mz <= (v.Channel4.Mz+(ppmPrecision*v.Channel4.Mz)) && k.Mz >= (v.Channel4.Mz-(ppmPrecision*v.Channel4.Mz)) {
						if k.Intensity > v.Channel4.Intensity {
							v.Channel4.Intensity = k.Intensity
						}
					}

					if k.Mz <= (v.Channel5.Mz+(ppmPrecision*v.Channel5.Mz)) && k.Mz >= (v.Channel5.Mz-(ppmPrecision*v.Channel5.Mz)) {
						if k.Intensity > v.Channel5.Intensity {
							v.Channel5.Intensity = k.Intensity
						}
					}

					if k.Mz <= (v.Channel6.Mz+(ppmPrecision*v.Channel6.Mz)) && k.Mz >= (v.Channel6.Mz-(ppmPrecision*v.Channel6.Mz)) {
						if k.Intensity > v.Channel6.Intensity {
							v.Channel6.Intensity = k.Intensity
						}
					}

					if k.Mz <= (v.Channel7.Mz+(ppmPrecision*v.Channel7.Mz)) && k.Mz >= (v.Channel7.Mz-(ppmPrecision*v.Channel7.Mz)) {
						if k.Intensity > v.Channel7.Intensity {
							v.Channel7.Intensity = k.Intensity
						}
					}

					if k.Mz <= (v.Channel8.Mz+(ppmPrecision*v.Channel8.Mz)) && k.Mz >= (v.Channel8.Mz-(ppmPrecision*v.Channel8.Mz)) {
						if k.Intensity > v.Channel8.Intensity {
							v.Channel8.Intensity = k.Intensity
						}
					}

					if k.Mz <= (v.Channel9.Mz+(ppmPrecision*v.Channel9.Mz)) && k.Mz >= (v.Channel9.Mz-(ppmPrecision*v.Channel9.Mz)) {
						if k.Intensity > v.Channel9.Intensity {
							v.Channel9.Intensity = k.Intensity
						}
					}

					if k.Mz <= (v.Channel10.Mz+(ppmPrecision*v.Channel10.Mz)) && k.Mz >= (v.Channel10.Mz-(ppmPrecision*v.Channel10.Mz)) {
						if k.Intensity > v.Channel10.Intensity {
							v.Channel10.Intensity = k.Intensity
						}
					}

					// if k.Mz <= (v.Channel1.Mz+(ppmPrecision*v.Channel1.Mz)) && k.Mz >= (v.Channel1.Mz-(ppmPrecision*v.Channel1.Mz)) {
					// 	v.Channel1.Intensity += k.Intensity
					// }
					//
					// if k.Mz <= (v.Channel2.Mz+(ppmPrecision*v.Channel2.Mz)) && k.Mz >= (v.Channel2.Mz-(ppmPrecision*v.Channel2.Mz)) {
					// 	v.Channel2.Intensity += k.Intensity
					// }
					//
					// if k.Mz <= (v.Channel3.Mz+(ppmPrecision*v.Channel2.Mz)) && k.Mz >= (v.Channel3.Mz-(ppmPrecision*v.Channel2.Mz)) {
					// 	v.Channel3.Intensity += k.Intensity
					// }
					//
					// if k.Mz <= (v.Channel4.Mz+(ppmPrecision*v.Channel4.Mz)) && k.Mz >= (v.Channel4.Mz-(ppmPrecision*v.Channel4.Mz)) {
					// 	v.Channel4.Intensity += k.Intensity
					// }
					//
					// if k.Mz <= (v.Channel5.Mz+(ppmPrecision*v.Channel5.Mz)) && k.Mz >= (v.Channel5.Mz-(ppmPrecision*v.Channel5.Mz)) {
					// 	v.Channel5.Intensity += k.Intensity
					// }
					//
					// if k.Mz <= (v.Channel6.Mz+(ppmPrecision*v.Channel6.Mz)) && k.Mz >= (v.Channel6.Mz-(ppmPrecision*v.Channel6.Mz)) {
					// 	v.Channel6.Intensity += k.Intensity
					// }
					//
					// if k.Mz <= (v.Channel7.Mz+(ppmPrecision*v.Channel7.Mz)) && k.Mz >= (v.Channel7.Mz-(ppmPrecision*v.Channel7.Mz)) {
					// 	v.Channel7.Intensity += k.Intensity
					// }
					//
					// if k.Mz <= (v.Channel8.Mz+(ppmPrecision*v.Channel8.Mz)) && k.Mz >= (v.Channel8.Mz-(ppmPrecision*v.Channel8.Mz)) {
					// 	v.Channel8.Intensity += k.Intensity
					// }
					//
					// if k.Mz <= (v.Channel9.Mz+(ppmPrecision*v.Channel9.Mz)) && k.Mz >= (v.Channel9.Mz-(ppmPrecision*v.Channel9.Mz)) {
					// 	v.Channel9.Intensity += k.Intensity
					// }
					//
					// if k.Mz <= (v.Channel10.Mz+(ppmPrecision*v.Channel10.Mz)) && k.Mz >= (v.Channel10.Mz-(ppmPrecision*v.Channel10.Mz)) {
					// 	v.Channel10.Intensity += k.Intensity
					// }

					if k.Mz > 150 {
						break
					}
				}

				spec[name] = v

			}
		}
	}

	return spec
}

func calculateIonPurity(d, f string, ms1 map[string]mz.MS1, ms2 map[string]mz.MS2, evi rep.Evidence) (rep.Evidence, error) {

	// organize them by index
	var indexedMs1 = make(map[int]mz.Ms1Scan)

	for _, v := range ms1 {
		for _, i := range v.Ms1Scan {
			indexedMs1[i.Index] = i
		}
	}

	ms1 = nil

	// range over IDs and spectra searching for a match
	for i := range evi.PSM {
		// locate the corresponding mz file for this identification
		name := strings.Split(evi.PSM[i].Spectrum, ".")

		s2, ok2 := ms2[name[0]]
		if ok2 {

			S2spec, S2ok := s2.Ms2Scan[evi.PSM[i].Spectrum]
			if S2ok {

				//for _, j := range s2.Ms2Scan {
				// search for the given spectra correspondent to the identification
				//FullSpectrumName := fmt.Sprintf("%s.%d", j.SpectrumName, j.Precursor.ChargeState)

				//if FullSpectrumName == evi.PSM[i].Spectrum {

				// recover the matching ms1 structure based on index number
				s1, ok1 := indexedMs1[S2spec.Precursor.ParentIndex]
				if ok1 {

					lowerDelta := S2spec.Precursor.SelectedIon - 1.5
					higherDelta := S2spec.Precursor.SelectedIon + 1.5
					var ions []mz.Ms1Peak

					for _, k := range s1.Spectrum {
						if k.Mz <= higherDelta && k.Mz >= lowerDelta {
							ions = append(ions, k)
						}
					}

					// create the list of mz differences for each peak
					var mzRatio []float64
					for k := 1; k <= 6; k++ {
						r := float64(k) * (float64(1) / float64(S2spec.Precursor.ChargeState))
						mzRatio = append(mzRatio, utils.ToFixed(r, 2))
					}

					var ionPackage []mz.Ms1Peak
					var summedInt float64
					for _, l := range ions {

						summedInt += l.Intensity

						for _, m := range mzRatio {
							if math.Abs(S2spec.Precursor.SelectedIon-l.Mz) <= (m+0.05) && math.Abs(S2spec.Precursor.SelectedIon-l.Mz) >= (m-0.05) {
								ionPackage = append(ionPackage, l)
								break
							}
						}
					}
					summedInt += S2spec.Precursor.PeakIntensity

					// calculate the total inensity for the selected ions from the ion package
					var summedPackageInt float64
					for _, k := range ionPackage {
						summedPackageInt += k.Intensity
					}
					summedPackageInt += S2spec.Precursor.PeakIntensity

					if summedInt == 0 {
						evi.PSM[i].Purity = 0
					} else {
						evi.PSM[i].Purity = utils.Round((summedPackageInt / summedInt), 5, 2)
					}

				}
			}

		}
	}

	return evi, nil
}

// mapLabeledSpectra maps all labeled spectra to ions
func mapLabeledSpectra(spectra map[string]tmt.Labels, purity float64, evi rep.Evidence) (rep.Evidence, error) {

	var ionLabelMap = make(map[string]tmt.Labels)
	var purityMap = make(map[string]float64)

	for i := range evi.PSM {
		split := strings.Split(evi.PSM[i].Spectrum, ".")
		name := fmt.Sprintf("%s.%s.%s", split[0], split[1], split[2])
		v, ok := spectra[name]
		if ok {
			evi.PSM[i].Labels.Spectrum = v.Spectrum
			evi.PSM[i].Labels.Index = v.Index
			evi.PSM[i].Labels.Channel1.Intensity = v.Channel1.Intensity
			evi.PSM[i].Labels.Channel2.Intensity = v.Channel2.Intensity
			evi.PSM[i].Labels.Channel3.Intensity = v.Channel3.Intensity
			evi.PSM[i].Labels.Channel4.Intensity = v.Channel4.Intensity
			evi.PSM[i].Labels.Channel5.Intensity = v.Channel5.Intensity
			evi.PSM[i].Labels.Channel6.Intensity = v.Channel6.Intensity
			evi.PSM[i].Labels.Channel7.Intensity = v.Channel7.Intensity
			evi.PSM[i].Labels.Channel8.Intensity = v.Channel8.Intensity
			evi.PSM[i].Labels.Channel9.Intensity = v.Channel9.Intensity
			evi.PSM[i].Labels.Channel10.Intensity = v.Channel10.Intensity

			// create a purity map for later use from ions and proteins
			if evi.PSM[i].Purity >= purity && evi.PSM[i].Probability >= 0.9 {
				purityMap[name] = evi.PSM[i].Purity
			}

		}
	}

	for i := range evi.Ions {
		for k := range evi.Ions[i].Spectra {

			v, ok := spectra[k]
			_, pure := purityMap[k]

			if ok && pure {
				evi.Ions[i].Labels.Channel1.Intensity += v.Channel1.Intensity
				evi.Ions[i].Labels.Channel2.Intensity += v.Channel2.Intensity
				evi.Ions[i].Labels.Channel3.Intensity += v.Channel3.Intensity
				evi.Ions[i].Labels.Channel4.Intensity += v.Channel4.Intensity
				evi.Ions[i].Labels.Channel5.Intensity += v.Channel5.Intensity
				evi.Ions[i].Labels.Channel6.Intensity += v.Channel6.Intensity
				evi.Ions[i].Labels.Channel7.Intensity += v.Channel7.Intensity
				evi.Ions[i].Labels.Channel8.Intensity += v.Channel8.Intensity
				evi.Ions[i].Labels.Channel9.Intensity += v.Channel9.Intensity
				evi.Ions[i].Labels.Channel10.Intensity += v.Channel10.Intensity

				var ion string
				if len(evi.Ions[i].ModifiedSequence) > 0 {
					ion = fmt.Sprintf("%s#%d", evi.Ions[i].ModifiedSequence, evi.Ions[i].ChargeState)
				} else {
					ion = fmt.Sprintf("%s#%d", evi.Ions[i].Sequence, evi.Ions[i].ChargeState)
				}

				ionLabelMap[ion] = evi.Ions[i].Labels

			}
		}
	}

	// global normalization
	var sumMap = make(map[uint8]float64)
	var allSum []float64

	for _, l := range ionLabelMap {
		sumMap[1] += l.Channel1.Intensity
		sumMap[2] += l.Channel2.Intensity
		sumMap[3] += l.Channel3.Intensity
		sumMap[4] += l.Channel4.Intensity
		sumMap[5] += l.Channel5.Intensity
		sumMap[6] += l.Channel6.Intensity
		sumMap[7] += l.Channel7.Intensity
		sumMap[8] += l.Channel8.Intensity
		sumMap[9] += l.Channel9.Intensity
		sumMap[10] += l.Channel10.Intensity
	}

	for _, v := range sumMap {
		allSum = append(allSum, v)
	}

	max, err := stats.Max(allSum)
	if err != nil {
		return evi, err
	}

	for i, l := range ionLabelMap {
		for j := range evi.Proteins {

			v, ok := evi.Proteins[j].TotalPeptideIons[i]
			if ok {
				value := v

				// Protein raw Intensity summation of all ions
				evi.Proteins[j].TotalLabels.Channel1.Intensity += l.Channel1.Intensity
				evi.Proteins[j].TotalLabels.Channel2.Intensity += l.Channel2.Intensity
				evi.Proteins[j].TotalLabels.Channel3.Intensity += l.Channel3.Intensity
				evi.Proteins[j].TotalLabels.Channel4.Intensity += l.Channel4.Intensity
				evi.Proteins[j].TotalLabels.Channel5.Intensity += l.Channel5.Intensity
				evi.Proteins[j].TotalLabels.Channel6.Intensity += l.Channel6.Intensity
				evi.Proteins[j].TotalLabels.Channel7.Intensity += l.Channel7.Intensity
				evi.Proteins[j].TotalLabels.Channel8.Intensity += l.Channel8.Intensity
				evi.Proteins[j].TotalLabels.Channel9.Intensity += l.Channel9.Intensity
				evi.Proteins[j].TotalLabels.Channel10.Intensity += l.Channel10.Intensity

				// Protein Normalized Intensity of all Ions
				evi.Proteins[j].TotalLabels.Channel1.NormIntensity += (l.Channel1.Intensity * (max / sumMap[1]))
				evi.Proteins[j].TotalLabels.Channel2.NormIntensity += (l.Channel2.Intensity * (max / sumMap[2]))
				evi.Proteins[j].TotalLabels.Channel3.NormIntensity += (l.Channel3.Intensity * (max / sumMap[3]))
				evi.Proteins[j].TotalLabels.Channel4.NormIntensity += (l.Channel4.Intensity * (max / sumMap[4]))
				evi.Proteins[j].TotalLabels.Channel5.NormIntensity += (l.Channel5.Intensity * (max / sumMap[5]))
				evi.Proteins[j].TotalLabels.Channel6.NormIntensity += (l.Channel6.Intensity * (max / sumMap[6]))
				evi.Proteins[j].TotalLabels.Channel7.NormIntensity += (l.Channel7.Intensity * (max / sumMap[7]))
				evi.Proteins[j].TotalLabels.Channel8.NormIntensity += (l.Channel8.Intensity * (max / sumMap[8]))
				evi.Proteins[j].TotalLabels.Channel9.NormIntensity += (l.Channel9.Intensity * (max / sumMap[9]))
				evi.Proteins[j].TotalLabels.Channel10.NormIntensity += (l.Channel10.Intensity * (max / sumMap[10]))

				// Peptide Ion Raw Intensity
				value.Labels.Channel1.Intensity = l.Channel1.Intensity
				value.Labels.Channel2.Intensity = l.Channel2.Intensity
				value.Labels.Channel3.Intensity = l.Channel3.Intensity
				value.Labels.Channel4.Intensity = l.Channel4.Intensity
				value.Labels.Channel5.Intensity = l.Channel5.Intensity
				value.Labels.Channel6.Intensity = l.Channel6.Intensity
				value.Labels.Channel7.Intensity = l.Channel7.Intensity
				value.Labels.Channel8.Intensity = l.Channel8.Intensity
				value.Labels.Channel9.Intensity = l.Channel9.Intensity
				value.Labels.Channel10.Intensity = l.Channel10.Intensity

				// Peptide Ion Normalized Intensity
				value.Labels.Channel1.NormIntensity = (l.Channel1.Intensity * (max / sumMap[1]))
				value.Labels.Channel2.NormIntensity = (l.Channel2.Intensity * (max / sumMap[2]))
				value.Labels.Channel3.NormIntensity = (l.Channel3.Intensity * (max / sumMap[3]))
				value.Labels.Channel4.NormIntensity = (l.Channel4.Intensity * (max / sumMap[4]))
				value.Labels.Channel5.NormIntensity = (l.Channel5.Intensity * (max / sumMap[5]))
				value.Labels.Channel6.NormIntensity = (l.Channel6.Intensity * (max / sumMap[6]))
				value.Labels.Channel7.NormIntensity = (l.Channel7.Intensity * (max / sumMap[7]))
				value.Labels.Channel8.NormIntensity = (l.Channel8.Intensity * (max / sumMap[8]))
				value.Labels.Channel9.NormIntensity = (l.Channel9.Intensity * (max / sumMap[9]))
				value.Labels.Channel10.NormIntensity = (l.Channel10.Intensity * (max / sumMap[10]))

				// Peptide Ion Summed Intensity
				// value.SummedLabelIntensity += (l.Channel1.Intensity * (max / sumMap[1]))
				// value.SummedLabelIntensity += (l.Channel2.Intensity * (max / sumMap[2]))
				// value.SummedLabelIntensity += (l.Channel3.Intensity * (max / sumMap[3]))
				// value.SummedLabelIntensity += (l.Channel4.Intensity * (max / sumMap[4]))
				// value.SummedLabelIntensity += (l.Channel5.Intensity * (max / sumMap[5]))
				// value.SummedLabelIntensity += (l.Channel6.Intensity * (max / sumMap[6]))
				// value.SummedLabelIntensity += (l.Channel7.Intensity * (max / sumMap[7]))
				// value.SummedLabelIntensity += (l.Channel8.Intensity * (max / sumMap[8]))
				// value.SummedLabelIntensity += (l.Channel9.Intensity * (max / sumMap[9]))
				// value.SummedLabelIntensity += (l.Channel10.Intensity * (max / sumMap[10]))

				evi.Proteins[j].TotalPeptideIons[i] = value
			}

			v, ok = evi.Proteins[j].UniquePeptideIons[i]
			if ok {
				value := v

				// update the protein intensity to keep the raw sum of all ions
				evi.Proteins[j].UniqueLabels.Channel1.Intensity += l.Channel1.Intensity
				evi.Proteins[j].UniqueLabels.Channel2.Intensity += l.Channel2.Intensity
				evi.Proteins[j].UniqueLabels.Channel3.Intensity += l.Channel3.Intensity
				evi.Proteins[j].UniqueLabels.Channel4.Intensity += l.Channel4.Intensity
				evi.Proteins[j].UniqueLabels.Channel5.Intensity += l.Channel5.Intensity
				evi.Proteins[j].UniqueLabels.Channel6.Intensity += l.Channel6.Intensity
				evi.Proteins[j].UniqueLabels.Channel7.Intensity += l.Channel7.Intensity
				evi.Proteins[j].UniqueLabels.Channel8.Intensity += l.Channel8.Intensity
				evi.Proteins[j].UniqueLabels.Channel9.Intensity += l.Channel9.Intensity
				evi.Proteins[j].UniqueLabels.Channel10.Intensity += l.Channel10.Intensity

				// update the protein inormintensities
				evi.Proteins[j].UniqueLabels.Channel1.NormIntensity += (l.Channel1.Intensity * (max / sumMap[1]))
				evi.Proteins[j].UniqueLabels.Channel2.NormIntensity += (l.Channel2.Intensity * (max / sumMap[2]))
				evi.Proteins[j].UniqueLabels.Channel3.NormIntensity += (l.Channel3.Intensity * (max / sumMap[3]))
				evi.Proteins[j].UniqueLabels.Channel4.NormIntensity += (l.Channel4.Intensity * (max / sumMap[4]))
				evi.Proteins[j].UniqueLabels.Channel5.NormIntensity += (l.Channel5.Intensity * (max / sumMap[5]))
				evi.Proteins[j].UniqueLabels.Channel6.NormIntensity += (l.Channel6.Intensity * (max / sumMap[6]))
				evi.Proteins[j].UniqueLabels.Channel7.NormIntensity += (l.Channel7.Intensity * (max / sumMap[7]))
				evi.Proteins[j].UniqueLabels.Channel8.NormIntensity += (l.Channel8.Intensity * (max / sumMap[8]))
				evi.Proteins[j].UniqueLabels.Channel9.NormIntensity += (l.Channel9.Intensity * (max / sumMap[9]))
				evi.Proteins[j].UniqueLabels.Channel10.NormIntensity += (l.Channel10.Intensity * (max / sumMap[10]))

				// Peptide Ion Raw Intensity
				value.Labels.Channel1.Intensity = l.Channel1.Intensity
				value.Labels.Channel2.Intensity = l.Channel2.Intensity
				value.Labels.Channel3.Intensity = l.Channel3.Intensity
				value.Labels.Channel4.Intensity = l.Channel4.Intensity
				value.Labels.Channel5.Intensity = l.Channel5.Intensity
				value.Labels.Channel6.Intensity = l.Channel6.Intensity
				value.Labels.Channel7.Intensity = l.Channel7.Intensity
				value.Labels.Channel8.Intensity = l.Channel8.Intensity
				value.Labels.Channel9.Intensity = l.Channel9.Intensity
				value.Labels.Channel10.Intensity = l.Channel10.Intensity

				// Peptide Ion Normalized Intensity
				value.Labels.Channel1.NormIntensity = (l.Channel1.Intensity * (max / sumMap[1]))
				value.Labels.Channel2.NormIntensity = (l.Channel2.Intensity * (max / sumMap[2]))
				value.Labels.Channel3.NormIntensity = (l.Channel3.Intensity * (max / sumMap[3]))
				value.Labels.Channel4.NormIntensity = (l.Channel4.Intensity * (max / sumMap[4]))
				value.Labels.Channel5.NormIntensity = (l.Channel5.Intensity * (max / sumMap[5]))
				value.Labels.Channel6.NormIntensity = (l.Channel6.Intensity * (max / sumMap[6]))
				value.Labels.Channel7.NormIntensity = (l.Channel7.Intensity * (max / sumMap[7]))
				value.Labels.Channel8.NormIntensity = (l.Channel8.Intensity * (max / sumMap[8]))
				value.Labels.Channel9.NormIntensity = (l.Channel9.Intensity * (max / sumMap[9]))
				value.Labels.Channel10.NormIntensity = (l.Channel10.Intensity * (max / sumMap[10]))

				// Peptide Ion Summed Intensity
				// value.SummedLabelIntensity += (l.Channel1.Intensity * (max / sumMap[1]))
				// value.SummedLabelIntensity += (l.Channel2.Intensity * (max / sumMap[2]))
				// value.SummedLabelIntensity += (l.Channel3.Intensity * (max / sumMap[3]))
				// value.SummedLabelIntensity += (l.Channel4.Intensity * (max / sumMap[4]))
				// value.SummedLabelIntensity += (l.Channel5.Intensity * (max / sumMap[5]))
				// value.SummedLabelIntensity += (l.Channel6.Intensity * (max / sumMap[6]))
				// value.SummedLabelIntensity += (l.Channel7.Intensity * (max / sumMap[7]))
				// value.SummedLabelIntensity += (l.Channel8.Intensity * (max / sumMap[8]))
				// value.SummedLabelIntensity += (l.Channel9.Intensity * (max / sumMap[9]))
				// value.SummedLabelIntensity += (l.Channel10.Intensity * (max / sumMap[10]))

				evi.Proteins[j].UniquePeptideIons[i] = value
			}

			v, ok = evi.Proteins[j].URazorPeptideIons[i]
			if ok {
				value := v

				// update the protein intensity to keep the raw sum of all ions
				evi.Proteins[j].RazorLabels.Channel1.Intensity += l.Channel1.Intensity
				evi.Proteins[j].RazorLabels.Channel2.Intensity += l.Channel2.Intensity
				evi.Proteins[j].RazorLabels.Channel3.Intensity += l.Channel3.Intensity
				evi.Proteins[j].RazorLabels.Channel4.Intensity += l.Channel4.Intensity
				evi.Proteins[j].RazorLabels.Channel5.Intensity += l.Channel5.Intensity
				evi.Proteins[j].RazorLabels.Channel6.Intensity += l.Channel6.Intensity
				evi.Proteins[j].RazorLabels.Channel7.Intensity += l.Channel7.Intensity
				evi.Proteins[j].RazorLabels.Channel8.Intensity += l.Channel8.Intensity
				evi.Proteins[j].RazorLabels.Channel9.Intensity += l.Channel9.Intensity
				evi.Proteins[j].RazorLabels.Channel10.Intensity += l.Channel10.Intensity

				// update the protein inormintensities
				evi.Proteins[j].RazorLabels.Channel1.NormIntensity += (l.Channel1.Intensity * (max / sumMap[1]))
				evi.Proteins[j].RazorLabels.Channel2.NormIntensity += (l.Channel2.Intensity * (max / sumMap[2]))
				evi.Proteins[j].RazorLabels.Channel3.NormIntensity += (l.Channel3.Intensity * (max / sumMap[3]))
				evi.Proteins[j].RazorLabels.Channel4.NormIntensity += (l.Channel4.Intensity * (max / sumMap[4]))
				evi.Proteins[j].RazorLabels.Channel5.NormIntensity += (l.Channel5.Intensity * (max / sumMap[5]))
				evi.Proteins[j].RazorLabels.Channel6.NormIntensity += (l.Channel6.Intensity * (max / sumMap[6]))
				evi.Proteins[j].RazorLabels.Channel7.NormIntensity += (l.Channel7.Intensity * (max / sumMap[7]))
				evi.Proteins[j].RazorLabels.Channel8.NormIntensity += (l.Channel8.Intensity * (max / sumMap[8]))
				evi.Proteins[j].RazorLabels.Channel9.NormIntensity += (l.Channel9.Intensity * (max / sumMap[9]))
				evi.Proteins[j].RazorLabels.Channel10.NormIntensity += (l.Channel10.Intensity * (max / sumMap[10]))

				// Peptide Ion Raw Intensity
				value.Labels.Channel1.Intensity = l.Channel1.Intensity
				value.Labels.Channel2.Intensity = l.Channel2.Intensity
				value.Labels.Channel3.Intensity = l.Channel3.Intensity
				value.Labels.Channel4.Intensity = l.Channel4.Intensity
				value.Labels.Channel5.Intensity = l.Channel5.Intensity
				value.Labels.Channel6.Intensity = l.Channel6.Intensity
				value.Labels.Channel7.Intensity = l.Channel7.Intensity
				value.Labels.Channel8.Intensity = l.Channel8.Intensity
				value.Labels.Channel9.Intensity = l.Channel9.Intensity
				value.Labels.Channel10.Intensity = l.Channel10.Intensity

				// Peptide Ion Normalized Intensity
				value.Labels.Channel1.NormIntensity = (l.Channel1.Intensity * (max / sumMap[1]))
				value.Labels.Channel2.NormIntensity = (l.Channel2.Intensity * (max / sumMap[2]))
				value.Labels.Channel3.NormIntensity = (l.Channel3.Intensity * (max / sumMap[3]))
				value.Labels.Channel4.NormIntensity = (l.Channel4.Intensity * (max / sumMap[4]))
				value.Labels.Channel5.NormIntensity = (l.Channel5.Intensity * (max / sumMap[5]))
				value.Labels.Channel6.NormIntensity = (l.Channel6.Intensity * (max / sumMap[6]))
				value.Labels.Channel7.NormIntensity = (l.Channel7.Intensity * (max / sumMap[7]))
				value.Labels.Channel8.NormIntensity = (l.Channel8.Intensity * (max / sumMap[8]))
				value.Labels.Channel9.NormIntensity = (l.Channel9.Intensity * (max / sumMap[9]))
				value.Labels.Channel10.NormIntensity = (l.Channel10.Intensity * (max / sumMap[10]))

				// Peptide Ion Summed Intensity
				// value.SummedLabelIntensity += (l.Channel1.Intensity * (max / sumMap[1]))
				// value.SummedLabelIntensity += (l.Channel2.Intensity * (max / sumMap[2]))
				// value.SummedLabelIntensity += (l.Channel3.Intensity * (max / sumMap[3]))
				// value.SummedLabelIntensity += (l.Channel4.Intensity * (max / sumMap[4]))
				// value.SummedLabelIntensity += (l.Channel5.Intensity * (max / sumMap[5]))
				// value.SummedLabelIntensity += (l.Channel6.Intensity * (max / sumMap[6]))
				// value.SummedLabelIntensity += (l.Channel7.Intensity * (max / sumMap[7]))
				// value.SummedLabelIntensity += (l.Channel8.Intensity * (max / sumMap[8]))
				// value.SummedLabelIntensity += (l.Channel9.Intensity * (max / sumMap[9]))
				// value.SummedLabelIntensity += (l.Channel10.Intensity * (max / sumMap[10]))

				evi.Proteins[j].URazorPeptideIons[i] = value
			}

		}
	}

	for i := range evi.Proteins {
		for _, v := range evi.Proteins[i].TotalPeptideIons {
			v.SummedLabelIntensity = (v.Labels.Channel1.NormIntensity +
				v.Labels.Channel2.NormIntensity +
				v.Labels.Channel3.NormIntensity +
				v.Labels.Channel4.NormIntensity +
				v.Labels.Channel5.NormIntensity +
				v.Labels.Channel6.NormIntensity +
				v.Labels.Channel7.NormIntensity +
				v.Labels.Channel8.NormIntensity +
				v.Labels.Channel9.NormIntensity +
				v.Labels.Channel10.NormIntensity)
		}

		for _, v := range evi.Proteins[i].UniquePeptideIons {
			v.SummedLabelIntensity = (v.Labels.Channel1.NormIntensity +
				v.Labels.Channel2.NormIntensity +
				v.Labels.Channel3.NormIntensity +
				v.Labels.Channel4.NormIntensity +
				v.Labels.Channel5.NormIntensity +
				v.Labels.Channel6.NormIntensity +
				v.Labels.Channel7.NormIntensity +
				v.Labels.Channel8.NormIntensity +
				v.Labels.Channel9.NormIntensity +
				v.Labels.Channel10.NormIntensity)
		}

		for _, v := range evi.Proteins[i].URazorPeptideIons {
			v.SummedLabelIntensity = (v.Labels.Channel1.NormIntensity +
				v.Labels.Channel2.NormIntensity +
				v.Labels.Channel3.NormIntensity +
				v.Labels.Channel4.NormIntensity +
				v.Labels.Channel5.NormIntensity +
				v.Labels.Channel6.NormIntensity +
				v.Labels.Channel7.NormIntensity +
				v.Labels.Channel8.NormIntensity +
				v.Labels.Channel9.NormIntensity +
				v.Labels.Channel10.NormIntensity)
		}
	}

	return evi, nil
}

func totalTop3LabelQuantification(evi rep.Evidence) (rep.Evidence, error) {

	for i := range evi.Proteins {

		p := make(PairList, len(evi.Proteins[i].TotalPeptideIons))

		j := 0
		for k, v := range evi.Proteins[i].TotalPeptideIons {
			p[j] = Pair{evi.Proteins[i].TotalPeptideIons[k], v.SummedLabelIntensity}
			j++
		}

		sort.Sort(sort.Reverse(p))

		var selectedIons []rep.IonEvidence

		var limit = 0
		if len(p) >= 3 {
			limit = 3
		} else if len(p) == 2 {
			limit = 2
		} else if len(p) == 1 {
			limit = 1
		}

		var counter = 0
		for _, j := range p {
			counter++
			if counter > limit {
				break
			}
			selectedIons = append(selectedIons, j.Key)
		}

		var c1Data float64
		var c2Data float64
		var c3Data float64
		var c4Data float64
		var c5Data float64
		var c6Data float64
		var c7Data float64
		var c8Data float64
		var c9Data float64
		var c10Data float64

		for _, j := range selectedIons {
			c1Data += j.Labels.Channel1.NormIntensity
			c2Data += j.Labels.Channel2.NormIntensity
			c3Data += j.Labels.Channel3.NormIntensity
			c4Data += j.Labels.Channel4.NormIntensity
			c5Data += j.Labels.Channel5.NormIntensity
			c6Data += j.Labels.Channel6.NormIntensity
			c7Data += j.Labels.Channel7.NormIntensity
			c8Data += j.Labels.Channel8.NormIntensity
			c9Data += j.Labels.Channel9.NormIntensity
			c10Data += j.Labels.Channel10.NormIntensity
		}

		evi.Proteins[i].TotalLabels.Channel1.TopIntensity = (c1Data / float64(limit))
		evi.Proteins[i].TotalLabels.Channel2.TopIntensity = (c2Data / float64(limit))
		evi.Proteins[i].TotalLabels.Channel3.TopIntensity = (c3Data / float64(limit))
		evi.Proteins[i].TotalLabels.Channel4.TopIntensity = (c4Data / float64(limit))
		evi.Proteins[i].TotalLabels.Channel5.TopIntensity = (c5Data / float64(limit))
		evi.Proteins[i].TotalLabels.Channel6.TopIntensity = (c6Data / float64(limit))
		evi.Proteins[i].TotalLabels.Channel7.TopIntensity = (c7Data / float64(limit))
		evi.Proteins[i].TotalLabels.Channel8.TopIntensity = (c8Data / float64(limit))
		evi.Proteins[i].TotalLabels.Channel9.TopIntensity = (c9Data / float64(limit))
		evi.Proteins[i].TotalLabels.Channel10.TopIntensity = (c10Data / float64(limit))

	}

	return evi, nil
}

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

// labelQuantificationOnTotalIons applies normalization to lable intensities
func labelQuantificationOnTotalIons(evi rep.Evidence) (rep.Evidence, error) {

	for i := range evi.Proteins {

		var totalIons []rep.IonEvidence
		for _, v := range evi.Proteins[i].TotalPeptideIons {
			totalIons = append(totalIons, v)
		}

		var c1Data []float64
		var c2Data []float64
		var c3Data []float64
		var c4Data []float64
		var c5Data []float64
		var c6Data []float64
		var c7Data []float64
		var c8Data []float64
		var c9Data []float64
		var c10Data []float64

		// determine the mean and the standard deviation of the mean
		for j := range totalIons {
			c1Data = append(c1Data, totalIons[j].Labels.Channel1.NormIntensity)
			c2Data = append(c2Data, totalIons[j].Labels.Channel2.NormIntensity)
			c3Data = append(c3Data, totalIons[j].Labels.Channel3.NormIntensity)
			c4Data = append(c4Data, totalIons[j].Labels.Channel4.NormIntensity)
			c5Data = append(c5Data, totalIons[j].Labels.Channel5.NormIntensity)
			c6Data = append(c6Data, totalIons[j].Labels.Channel6.NormIntensity)
			c7Data = append(c7Data, totalIons[j].Labels.Channel7.NormIntensity)
			c8Data = append(c8Data, totalIons[j].Labels.Channel8.NormIntensity)
			c9Data = append(c9Data, totalIons[j].Labels.Channel9.NormIntensity)
			c10Data = append(c10Data, totalIons[j].Labels.Channel10.NormIntensity)
		}

		c1Mean, _ := stats.Mean(c1Data)
		c2Mean, _ := stats.Mean(c2Data)
		c3Mean, _ := stats.Mean(c3Data)
		c4Mean, _ := stats.Mean(c4Data)
		c5Mean, _ := stats.Mean(c5Data)
		c6Mean, _ := stats.Mean(c6Data)
		c7Mean, _ := stats.Mean(c7Data)
		c8Mean, _ := stats.Mean(c8Data)
		c9Mean, _ := stats.Mean(c9Data)
		c10Mean, _ := stats.Mean(c10Data)
		// if err != nil {
		// 	fmt.Println("AQUI")
		// 	return err
		// }

		c1StDev, _ := stats.StandardDeviationPopulation(c1Data)
		c2StDev, _ := stats.StandardDeviationPopulation(c2Data)
		c3StDev, _ := stats.StandardDeviationPopulation(c3Data)
		c4StDev, _ := stats.StandardDeviationPopulation(c4Data)
		c5StDev, _ := stats.StandardDeviationPopulation(c5Data)
		c6StDev, _ := stats.StandardDeviationPopulation(c6Data)
		c7StDev, _ := stats.StandardDeviationPopulation(c7Data)
		c8StDev, _ := stats.StandardDeviationPopulation(c8Data)
		c9StDev, _ := stats.StandardDeviationPopulation(c9Data)
		c10StDev, _ := stats.StandardDeviationPopulation(c10Data)
		// if err != nil {
		// 	return err
		// }

		// remov those that deviate from the mean by more than 2 sigma
		loC1Sigma := (c1Mean - 2*(c1StDev))
		hiC1Sigma := (c1Mean + 2*(c1StDev))

		loC2Sigma := (c2Mean - 2*(c2StDev))
		hiC2Sigma := (c2Mean + 2*(c2StDev))

		loC3Sigma := (c3Mean - 2*(c3StDev))
		hiC3Sigma := (c3Mean + 2*(c3StDev))

		loC4Sigma := (c4Mean - 2*(c4StDev))
		hiC4Sigma := (c4Mean + 2*(c4StDev))

		loC5Sigma := (c5Mean - 2*(c5StDev))
		hiC5Sigma := (c5Mean + 2*(c5StDev))

		loC6Sigma := (c6Mean - 2*(c6StDev))
		hiC6Sigma := (c6Mean + 2*(c6StDev))

		loC7Sigma := (c7Mean - 2*(c7StDev))
		hiC7Sigma := (c7Mean + 2*(c7StDev))

		loC8Sigma := (c8Mean - 2*(c8StDev))
		hiC8Sigma := (c8Mean + 2*(c8StDev))

		loC9Sigma := (c9Mean - 2*(c9StDev))
		hiC9Sigma := (c9Mean + 2*(c9StDev))

		loC10Sigma := (c10Mean - 2*(c10StDev))
		hiC10Sigma := (c10Mean + 2*(c10StDev))

		var normIons = make(map[int][]rep.IonEvidence)

		for i := range totalIons {

			if totalIons[i].Labels.Channel1.NormIntensity > 0 && totalIons[i].Labels.Channel1.NormIntensity >= loC1Sigma && totalIons[i].Labels.Channel1.NormIntensity <= hiC1Sigma {
				normIons[1] = append(normIons[1], totalIons[i])
			}

			if totalIons[i].Labels.Channel2.NormIntensity > 0 && totalIons[i].Labels.Channel2.NormIntensity >= loC2Sigma && totalIons[i].Labels.Channel2.NormIntensity <= hiC2Sigma {
				normIons[2] = append(normIons[2], totalIons[i])
			}

			if totalIons[i].Labels.Channel3.NormIntensity > 0 && totalIons[i].Labels.Channel3.NormIntensity >= loC3Sigma && totalIons[i].Labels.Channel3.NormIntensity <= hiC3Sigma {
				normIons[3] = append(normIons[3], totalIons[i])
			}

			if totalIons[i].Labels.Channel4.NormIntensity > 0 && totalIons[i].Labels.Channel4.NormIntensity >= loC4Sigma && totalIons[i].Labels.Channel4.NormIntensity <= hiC4Sigma {
				normIons[4] = append(normIons[4], totalIons[i])
			}

			if totalIons[i].Labels.Channel5.NormIntensity > 0 && totalIons[i].Labels.Channel5.NormIntensity >= loC5Sigma && totalIons[i].Labels.Channel5.NormIntensity <= hiC5Sigma {
				normIons[5] = append(normIons[5], totalIons[i])
			}

			if totalIons[i].Labels.Channel6.NormIntensity > 0 && totalIons[i].Labels.Channel6.NormIntensity >= loC6Sigma && totalIons[i].Labels.Channel6.NormIntensity <= hiC6Sigma {
				normIons[6] = append(normIons[6], totalIons[i])
			}

			if totalIons[i].Labels.Channel7.NormIntensity > 0 && totalIons[i].Labels.Channel7.NormIntensity >= loC7Sigma && totalIons[i].Labels.Channel7.NormIntensity <= hiC7Sigma {
				normIons[7] = append(normIons[7], totalIons[i])
			}

			if totalIons[i].Labels.Channel8.NormIntensity > 0 && totalIons[i].Labels.Channel8.NormIntensity >= loC8Sigma && totalIons[i].Labels.Channel8.NormIntensity <= hiC8Sigma {
				normIons[8] = append(normIons[8], totalIons[i])
			}

			if totalIons[i].Labels.Channel9.NormIntensity > 0 && totalIons[i].Labels.Channel9.NormIntensity >= loC9Sigma && totalIons[i].Labels.Channel9.NormIntensity <= hiC9Sigma {
				normIons[9] = append(normIons[9], totalIons[i])
			}

			if totalIons[i].Labels.Channel10.NormIntensity > 0 && totalIons[i].Labels.Channel10.NormIntensity >= loC10Sigma && totalIons[i].Labels.Channel10.NormIntensity <= hiC10Sigma {
				normIons[10] = append(normIons[10], totalIons[i])
			}

		}

		// recalculate the mean and standard deviation
		c1Data = nil
		c2Data = nil
		c3Data = nil
		c4Data = nil
		c5Data = nil
		c6Data = nil
		c7Data = nil
		c8Data = nil
		c9Data = nil
		c10Data = nil

		for _, v := range normIons[1] {
			c1Data = append(c1Data, v.Labels.Channel1.NormIntensity)
		}

		for _, v := range normIons[2] {
			c2Data = append(c2Data, v.Labels.Channel2.NormIntensity)
		}

		for _, v := range normIons[3] {
			c3Data = append(c3Data, v.Labels.Channel3.NormIntensity)
		}

		for _, v := range normIons[4] {
			c4Data = append(c4Data, v.Labels.Channel4.NormIntensity)
		}

		for _, v := range normIons[5] {
			c5Data = append(c5Data, v.Labels.Channel5.NormIntensity)
		}

		for _, v := range normIons[6] {
			c6Data = append(c6Data, v.Labels.Channel6.NormIntensity)
		}

		for _, v := range normIons[7] {
			c7Data = append(c7Data, v.Labels.Channel7.NormIntensity)
		}

		for _, v := range normIons[8] {
			c8Data = append(c8Data, v.Labels.Channel8.NormIntensity)
		}

		for _, v := range normIons[9] {
			c9Data = append(c9Data, v.Labels.Channel9.NormIntensity)
		}

		for _, v := range normIons[10] {
			c10Data = append(c10Data, v.Labels.Channel10.NormIntensity)
		}

		evi.Proteins[i].TotalLabels.Channel1.Mean, _ = stats.Mean(c1Data)
		evi.Proteins[i].TotalLabels.Channel2.Mean, _ = stats.Mean(c2Data)
		evi.Proteins[i].TotalLabels.Channel3.Mean, _ = stats.Mean(c3Data)
		evi.Proteins[i].TotalLabels.Channel4.Mean, _ = stats.Mean(c4Data)
		evi.Proteins[i].TotalLabels.Channel5.Mean, _ = stats.Mean(c5Data)
		evi.Proteins[i].TotalLabels.Channel6.Mean, _ = stats.Mean(c6Data)
		evi.Proteins[i].TotalLabels.Channel7.Mean, _ = stats.Mean(c7Data)
		evi.Proteins[i].TotalLabels.Channel8.Mean, _ = stats.Mean(c8Data)
		evi.Proteins[i].TotalLabels.Channel9.Mean, _ = stats.Mean(c9Data)
		evi.Proteins[i].TotalLabels.Channel10.Mean, _ = stats.Mean(c10Data)

		evi.Proteins[i].TotalLabels.Channel1.StDev, _ = stats.StandardDeviationPopulation(c1Data)
		evi.Proteins[i].TotalLabels.Channel2.StDev, _ = stats.StandardDeviationPopulation(c2Data)
		evi.Proteins[i].TotalLabels.Channel3.StDev, _ = stats.StandardDeviationPopulation(c3Data)
		evi.Proteins[i].TotalLabels.Channel4.StDev, _ = stats.StandardDeviationPopulation(c4Data)
		evi.Proteins[i].TotalLabels.Channel5.StDev, _ = stats.StandardDeviationPopulation(c5Data)
		evi.Proteins[i].TotalLabels.Channel6.StDev, _ = stats.StandardDeviationPopulation(c6Data)
		evi.Proteins[i].TotalLabels.Channel7.StDev, _ = stats.StandardDeviationPopulation(c7Data)
		evi.Proteins[i].TotalLabels.Channel8.StDev, _ = stats.StandardDeviationPopulation(c8Data)
		evi.Proteins[i].TotalLabels.Channel9.StDev, _ = stats.StandardDeviationPopulation(c9Data)
		evi.Proteins[i].TotalLabels.Channel10.StDev, _ = stats.StandardDeviationPopulation(c10Data)

	}

	return evi, nil
}

// labelQuantificationOnUniqueIons applies normalization to lable intensities
func labelQuantificationOnUniqueIons(evi rep.Evidence) (rep.Evidence, error) {

	for i := range evi.Proteins {

		var ions []rep.IonEvidence
		for _, v := range evi.Proteins[i].UniquePeptideIons {
			ions = append(ions, v)
		}

		var c1Data []float64
		var c2Data []float64
		var c3Data []float64
		var c4Data []float64
		var c5Data []float64
		var c6Data []float64
		var c7Data []float64
		var c8Data []float64
		var c9Data []float64
		var c10Data []float64

		// determine the mean and the standard deviation of the mean
		for i := range ions {
			c1Data = append(c1Data, ions[i].Labels.Channel1.NormIntensity)
			c2Data = append(c2Data, ions[i].Labels.Channel2.NormIntensity)
			c3Data = append(c3Data, ions[i].Labels.Channel3.NormIntensity)
			c4Data = append(c4Data, ions[i].Labels.Channel4.NormIntensity)
			c5Data = append(c5Data, ions[i].Labels.Channel5.NormIntensity)
			c6Data = append(c6Data, ions[i].Labels.Channel6.NormIntensity)
			c7Data = append(c7Data, ions[i].Labels.Channel7.NormIntensity)
			c8Data = append(c8Data, ions[i].Labels.Channel8.NormIntensity)
			c9Data = append(c9Data, ions[i].Labels.Channel9.NormIntensity)
			c10Data = append(c10Data, ions[i].Labels.Channel10.NormIntensity)
		}

		c1Mean, _ := stats.Mean(c1Data)
		c2Mean, _ := stats.Mean(c2Data)
		c3Mean, _ := stats.Mean(c3Data)
		c4Mean, _ := stats.Mean(c4Data)
		c5Mean, _ := stats.Mean(c5Data)
		c6Mean, _ := stats.Mean(c6Data)
		c7Mean, _ := stats.Mean(c7Data)
		c8Mean, _ := stats.Mean(c8Data)
		c9Mean, _ := stats.Mean(c9Data)
		c10Mean, _ := stats.Mean(c10Data)
		// if err != nil {
		// 	return err
		// }

		c1StDev, _ := stats.StandardDeviationPopulation(c1Data)
		c2StDev, _ := stats.StandardDeviationPopulation(c2Data)
		c3StDev, _ := stats.StandardDeviationPopulation(c3Data)
		c4StDev, _ := stats.StandardDeviationPopulation(c4Data)
		c5StDev, _ := stats.StandardDeviationPopulation(c5Data)
		c6StDev, _ := stats.StandardDeviationPopulation(c6Data)
		c7StDev, _ := stats.StandardDeviationPopulation(c7Data)
		c8StDev, _ := stats.StandardDeviationPopulation(c8Data)
		c9StDev, _ := stats.StandardDeviationPopulation(c9Data)
		c10StDev, _ := stats.StandardDeviationPopulation(c10Data)
		// if err != nil {
		// 	return err
		// }

		// remov those that deviate from the mean by more than 2 sigma
		loC1Sigma := (c1Mean - 2*(c1StDev))
		hiC1Sigma := (c1Mean + 2*(c1StDev))

		loC2Sigma := (c2Mean - 2*(c2StDev))
		hiC2Sigma := (c2Mean + 2*(c2StDev))

		loC3Sigma := (c3Mean - 2*(c3StDev))
		hiC3Sigma := (c3Mean + 2*(c3StDev))

		loC4Sigma := (c4Mean - 2*(c4StDev))
		hiC4Sigma := (c4Mean + 2*(c4StDev))

		loC5Sigma := (c5Mean - 2*(c5StDev))
		hiC5Sigma := (c5Mean + 2*(c5StDev))

		loC6Sigma := (c6Mean - 2*(c6StDev))
		hiC6Sigma := (c6Mean + 2*(c6StDev))

		loC7Sigma := (c7Mean - 2*(c7StDev))
		hiC7Sigma := (c7Mean + 2*(c7StDev))

		loC8Sigma := (c8Mean - 2*(c8StDev))
		hiC8Sigma := (c8Mean + 2*(c8StDev))

		loC9Sigma := (c9Mean - 2*(c9StDev))
		hiC9Sigma := (c9Mean + 2*(c9StDev))

		loC10Sigma := (c10Mean - 2*(c10StDev))
		hiC10Sigma := (c10Mean + 2*(c10StDev))

		var normIons = make(map[int][]rep.IonEvidence)

		for i := range ions {

			if ions[i].Labels.Channel1.NormIntensity > 0 && ions[i].Labels.Channel1.NormIntensity >= loC1Sigma && ions[i].Labels.Channel1.NormIntensity <= hiC1Sigma {
				normIons[1] = append(normIons[1], ions[i])
			}

			if ions[i].Labels.Channel2.NormIntensity > 0 && ions[i].Labels.Channel2.NormIntensity >= loC2Sigma && ions[i].Labels.Channel2.NormIntensity <= hiC2Sigma {
				normIons[2] = append(normIons[2], ions[i])
			}

			if ions[i].Labels.Channel3.NormIntensity > 0 && ions[i].Labels.Channel3.NormIntensity >= loC3Sigma && ions[i].Labels.Channel3.NormIntensity <= hiC3Sigma {
				normIons[3] = append(normIons[3], ions[i])
			}

			if ions[i].Labels.Channel4.NormIntensity > 0 && ions[i].Labels.Channel4.NormIntensity >= loC4Sigma && ions[i].Labels.Channel4.NormIntensity <= hiC4Sigma {
				normIons[4] = append(normIons[4], ions[i])
			}

			if ions[i].Labels.Channel5.NormIntensity > 0 && ions[i].Labels.Channel5.NormIntensity >= loC5Sigma && ions[i].Labels.Channel5.NormIntensity <= hiC5Sigma {
				normIons[5] = append(normIons[5], ions[i])
			}

			if ions[i].Labels.Channel6.NormIntensity > 0 && ions[i].Labels.Channel6.NormIntensity >= loC6Sigma && ions[i].Labels.Channel6.NormIntensity <= hiC6Sigma {
				normIons[6] = append(normIons[6], ions[i])
			}

			if ions[i].Labels.Channel7.NormIntensity > 0 && ions[i].Labels.Channel7.NormIntensity >= loC7Sigma && ions[i].Labels.Channel7.NormIntensity <= hiC7Sigma {
				normIons[7] = append(normIons[7], ions[i])
			}

			if ions[i].Labels.Channel8.NormIntensity > 0 && ions[i].Labels.Channel8.NormIntensity >= loC8Sigma && ions[i].Labels.Channel8.NormIntensity <= hiC8Sigma {
				normIons[8] = append(normIons[8], ions[i])
			}

			if ions[i].Labels.Channel9.NormIntensity > 0 && ions[i].Labels.Channel9.NormIntensity >= loC9Sigma && ions[i].Labels.Channel9.NormIntensity <= hiC9Sigma {
				normIons[9] = append(normIons[9], ions[i])
			}

			if ions[i].Labels.Channel10.NormIntensity > 0 && ions[i].Labels.Channel10.NormIntensity >= loC10Sigma && ions[i].Labels.Channel10.NormIntensity <= hiC10Sigma {
				normIons[10] = append(normIons[10], ions[i])
			}

		}

		// recalculate the mean and standard deviation
		c1Data = nil
		c2Data = nil
		c3Data = nil
		c4Data = nil
		c5Data = nil
		c6Data = nil
		c7Data = nil
		c8Data = nil
		c9Data = nil
		c10Data = nil

		for _, v := range normIons[1] {
			c1Data = append(c1Data, v.Labels.Channel1.NormIntensity)
		}

		for _, v := range normIons[2] {
			c2Data = append(c2Data, v.Labels.Channel2.NormIntensity)
		}

		for _, v := range normIons[3] {
			c3Data = append(c3Data, v.Labels.Channel3.NormIntensity)
		}

		for _, v := range normIons[4] {
			c4Data = append(c4Data, v.Labels.Channel4.NormIntensity)
		}

		for _, v := range normIons[5] {
			c5Data = append(c5Data, v.Labels.Channel5.NormIntensity)
		}

		for _, v := range normIons[6] {
			c6Data = append(c6Data, v.Labels.Channel6.NormIntensity)
		}

		for _, v := range normIons[7] {
			c7Data = append(c7Data, v.Labels.Channel7.NormIntensity)
		}

		for _, v := range normIons[8] {
			c8Data = append(c8Data, v.Labels.Channel8.NormIntensity)
		}

		for _, v := range normIons[9] {
			c9Data = append(c9Data, v.Labels.Channel9.NormIntensity)
		}

		for _, v := range normIons[10] {
			c10Data = append(c10Data, v.Labels.Channel10.NormIntensity)
		}

		evi.Proteins[i].UniqueLabels.Channel1.Mean, _ = stats.Mean(c1Data)
		evi.Proteins[i].UniqueLabels.Channel2.Mean, _ = stats.Mean(c2Data)
		evi.Proteins[i].UniqueLabels.Channel3.Mean, _ = stats.Mean(c3Data)
		evi.Proteins[i].UniqueLabels.Channel4.Mean, _ = stats.Mean(c4Data)
		evi.Proteins[i].UniqueLabels.Channel5.Mean, _ = stats.Mean(c5Data)
		evi.Proteins[i].UniqueLabels.Channel6.Mean, _ = stats.Mean(c6Data)
		evi.Proteins[i].UniqueLabels.Channel7.Mean, _ = stats.Mean(c7Data)
		evi.Proteins[i].UniqueLabels.Channel8.Mean, _ = stats.Mean(c8Data)
		evi.Proteins[i].UniqueLabels.Channel9.Mean, _ = stats.Mean(c9Data)
		evi.Proteins[i].UniqueLabels.Channel10.Mean, _ = stats.Mean(c10Data)

		evi.Proteins[i].UniqueLabels.Channel1.StDev, _ = stats.StandardDeviationPopulation(c1Data)
		evi.Proteins[i].UniqueLabels.Channel2.StDev, _ = stats.StandardDeviationPopulation(c2Data)
		evi.Proteins[i].UniqueLabels.Channel3.StDev, _ = stats.StandardDeviationPopulation(c3Data)
		evi.Proteins[i].UniqueLabels.Channel4.StDev, _ = stats.StandardDeviationPopulation(c4Data)
		evi.Proteins[i].UniqueLabels.Channel5.StDev, _ = stats.StandardDeviationPopulation(c5Data)
		evi.Proteins[i].UniqueLabels.Channel6.StDev, _ = stats.StandardDeviationPopulation(c6Data)
		evi.Proteins[i].UniqueLabels.Channel7.StDev, _ = stats.StandardDeviationPopulation(c7Data)
		evi.Proteins[i].UniqueLabels.Channel8.StDev, _ = stats.StandardDeviationPopulation(c8Data)
		evi.Proteins[i].UniqueLabels.Channel9.StDev, _ = stats.StandardDeviationPopulation(c9Data)
		evi.Proteins[i].UniqueLabels.Channel10.StDev, _ = stats.StandardDeviationPopulation(c10Data)

	}

	return evi, nil
}

// labelQuantificationOnUniqueIons applies normalization to lable intensities
func labelQuantificationOnURazors(evi rep.Evidence) (rep.Evidence, error) {

	for i := range evi.Proteins {

		var ions []rep.IonEvidence
		for _, v := range evi.Proteins[i].URazorPeptideIons {
			ions = append(ions, v)
		}

		var c1Data []float64
		var c2Data []float64
		var c3Data []float64
		var c4Data []float64
		var c5Data []float64
		var c6Data []float64
		var c7Data []float64
		var c8Data []float64
		var c9Data []float64
		var c10Data []float64

		// determine the mean and the standard deviation of the mean
		for i := range ions {
			c1Data = append(c1Data, ions[i].Labels.Channel1.NormIntensity)
			c2Data = append(c2Data, ions[i].Labels.Channel2.NormIntensity)
			c3Data = append(c3Data, ions[i].Labels.Channel3.NormIntensity)
			c4Data = append(c4Data, ions[i].Labels.Channel4.NormIntensity)
			c5Data = append(c5Data, ions[i].Labels.Channel5.NormIntensity)
			c6Data = append(c6Data, ions[i].Labels.Channel6.NormIntensity)
			c7Data = append(c7Data, ions[i].Labels.Channel7.NormIntensity)
			c8Data = append(c8Data, ions[i].Labels.Channel8.NormIntensity)
			c9Data = append(c9Data, ions[i].Labels.Channel9.NormIntensity)
			c10Data = append(c10Data, ions[i].Labels.Channel10.NormIntensity)
		}

		c1Mean, _ := stats.Mean(c1Data)
		c2Mean, _ := stats.Mean(c2Data)
		c3Mean, _ := stats.Mean(c3Data)
		c4Mean, _ := stats.Mean(c4Data)
		c5Mean, _ := stats.Mean(c5Data)
		c6Mean, _ := stats.Mean(c6Data)
		c7Mean, _ := stats.Mean(c7Data)
		c8Mean, _ := stats.Mean(c8Data)
		c9Mean, _ := stats.Mean(c9Data)
		c10Mean, _ := stats.Mean(c10Data)
		// if err != nil {
		// 	return err
		// }

		c1StDev, _ := stats.StandardDeviationPopulation(c1Data)
		c2StDev, _ := stats.StandardDeviationPopulation(c2Data)
		c3StDev, _ := stats.StandardDeviationPopulation(c3Data)
		c4StDev, _ := stats.StandardDeviationPopulation(c4Data)
		c5StDev, _ := stats.StandardDeviationPopulation(c5Data)
		c6StDev, _ := stats.StandardDeviationPopulation(c6Data)
		c7StDev, _ := stats.StandardDeviationPopulation(c7Data)
		c8StDev, _ := stats.StandardDeviationPopulation(c8Data)
		c9StDev, _ := stats.StandardDeviationPopulation(c9Data)
		c10StDev, _ := stats.StandardDeviationPopulation(c10Data)
		// if err != nil {
		// 	return err
		// }

		// remov those that deviate from the mean by more than 2 sigma
		loC1Sigma := (c1Mean - 2*(c1StDev))
		hiC1Sigma := (c1Mean + 2*(c1StDev))

		loC2Sigma := (c2Mean - 2*(c2StDev))
		hiC2Sigma := (c2Mean + 2*(c2StDev))

		loC3Sigma := (c3Mean - 2*(c3StDev))
		hiC3Sigma := (c3Mean + 2*(c3StDev))

		loC4Sigma := (c4Mean - 2*(c4StDev))
		hiC4Sigma := (c4Mean + 2*(c4StDev))

		loC5Sigma := (c5Mean - 2*(c5StDev))
		hiC5Sigma := (c5Mean + 2*(c5StDev))

		loC6Sigma := (c6Mean - 2*(c6StDev))
		hiC6Sigma := (c6Mean + 2*(c6StDev))

		loC7Sigma := (c7Mean - 2*(c7StDev))
		hiC7Sigma := (c7Mean + 2*(c7StDev))

		loC8Sigma := (c8Mean - 2*(c8StDev))
		hiC8Sigma := (c8Mean + 2*(c8StDev))

		loC9Sigma := (c9Mean - 2*(c9StDev))
		hiC9Sigma := (c9Mean + 2*(c9StDev))

		loC10Sigma := (c10Mean - 2*(c10StDev))
		hiC10Sigma := (c10Mean + 2*(c10StDev))

		var normIons = make(map[int][]rep.IonEvidence)

		for i := range ions {

			if ions[i].Labels.Channel1.NormIntensity > 0 && ions[i].Labels.Channel1.NormIntensity >= loC1Sigma && ions[i].Labels.Channel1.NormIntensity <= hiC1Sigma {
				normIons[1] = append(normIons[1], ions[i])
			}

			if ions[i].Labels.Channel2.NormIntensity > 0 && ions[i].Labels.Channel2.NormIntensity >= loC2Sigma && ions[i].Labels.Channel2.NormIntensity <= hiC2Sigma {
				normIons[2] = append(normIons[2], ions[i])
			}

			if ions[i].Labels.Channel3.NormIntensity > 0 && ions[i].Labels.Channel3.NormIntensity >= loC3Sigma && ions[i].Labels.Channel3.NormIntensity <= hiC3Sigma {
				normIons[3] = append(normIons[3], ions[i])
			}

			if ions[i].Labels.Channel4.NormIntensity > 0 && ions[i].Labels.Channel4.NormIntensity >= loC4Sigma && ions[i].Labels.Channel4.NormIntensity <= hiC4Sigma {
				normIons[4] = append(normIons[4], ions[i])
			}

			if ions[i].Labels.Channel5.NormIntensity > 0 && ions[i].Labels.Channel5.NormIntensity >= loC5Sigma && ions[i].Labels.Channel5.NormIntensity <= hiC5Sigma {
				normIons[5] = append(normIons[5], ions[i])
			}

			if ions[i].Labels.Channel6.NormIntensity > 0 && ions[i].Labels.Channel6.NormIntensity >= loC6Sigma && ions[i].Labels.Channel6.NormIntensity <= hiC6Sigma {
				normIons[6] = append(normIons[6], ions[i])
			}

			if ions[i].Labels.Channel7.NormIntensity > 0 && ions[i].Labels.Channel7.NormIntensity >= loC7Sigma && ions[i].Labels.Channel7.NormIntensity <= hiC7Sigma {
				normIons[7] = append(normIons[7], ions[i])
			}

			if ions[i].Labels.Channel8.NormIntensity > 0 && ions[i].Labels.Channel8.NormIntensity >= loC8Sigma && ions[i].Labels.Channel8.NormIntensity <= hiC8Sigma {
				normIons[8] = append(normIons[8], ions[i])
			}

			if ions[i].Labels.Channel9.NormIntensity > 0 && ions[i].Labels.Channel9.NormIntensity >= loC9Sigma && ions[i].Labels.Channel9.NormIntensity <= hiC9Sigma {
				normIons[9] = append(normIons[9], ions[i])
			}

			if ions[i].Labels.Channel10.NormIntensity > 0 && ions[i].Labels.Channel10.NormIntensity >= loC10Sigma && ions[i].Labels.Channel10.NormIntensity <= hiC10Sigma {
				normIons[10] = append(normIons[10], ions[i])
			}

		}

		// recalculate the mean and standard deviation
		c1Data = nil
		c2Data = nil
		c3Data = nil
		c4Data = nil
		c5Data = nil
		c6Data = nil
		c7Data = nil
		c8Data = nil
		c9Data = nil
		c10Data = nil

		for _, v := range normIons[1] {
			c1Data = append(c1Data, v.Labels.Channel1.NormIntensity)
		}

		for _, v := range normIons[2] {
			c2Data = append(c2Data, v.Labels.Channel2.NormIntensity)
		}

		for _, v := range normIons[3] {
			c3Data = append(c3Data, v.Labels.Channel3.NormIntensity)
		}

		for _, v := range normIons[4] {
			c4Data = append(c4Data, v.Labels.Channel4.NormIntensity)
		}

		for _, v := range normIons[5] {
			c5Data = append(c5Data, v.Labels.Channel5.NormIntensity)
		}

		for _, v := range normIons[6] {
			c6Data = append(c6Data, v.Labels.Channel6.NormIntensity)
		}

		for _, v := range normIons[7] {
			c7Data = append(c7Data, v.Labels.Channel7.NormIntensity)
		}

		for _, v := range normIons[8] {
			c8Data = append(c8Data, v.Labels.Channel8.NormIntensity)
		}

		for _, v := range normIons[9] {
			c9Data = append(c9Data, v.Labels.Channel9.NormIntensity)
		}

		for _, v := range normIons[10] {
			c10Data = append(c10Data, v.Labels.Channel10.NormIntensity)
		}

		evi.Proteins[i].RazorLabels.Channel1.Mean, _ = stats.Mean(c1Data)
		evi.Proteins[i].RazorLabels.Channel2.Mean, _ = stats.Mean(c2Data)
		evi.Proteins[i].RazorLabels.Channel3.Mean, _ = stats.Mean(c3Data)
		evi.Proteins[i].RazorLabels.Channel4.Mean, _ = stats.Mean(c4Data)
		evi.Proteins[i].RazorLabels.Channel5.Mean, _ = stats.Mean(c5Data)
		evi.Proteins[i].RazorLabels.Channel6.Mean, _ = stats.Mean(c6Data)
		evi.Proteins[i].RazorLabels.Channel7.Mean, _ = stats.Mean(c7Data)
		evi.Proteins[i].RazorLabels.Channel8.Mean, _ = stats.Mean(c8Data)
		evi.Proteins[i].RazorLabels.Channel9.Mean, _ = stats.Mean(c9Data)
		evi.Proteins[i].RazorLabels.Channel10.Mean, _ = stats.Mean(c10Data)

		evi.Proteins[i].RazorLabels.Channel1.StDev, _ = stats.StandardDeviationPopulation(c1Data)
		evi.Proteins[i].RazorLabels.Channel2.StDev, _ = stats.StandardDeviationPopulation(c2Data)
		evi.Proteins[i].RazorLabels.Channel3.StDev, _ = stats.StandardDeviationPopulation(c3Data)
		evi.Proteins[i].RazorLabels.Channel4.StDev, _ = stats.StandardDeviationPopulation(c4Data)
		evi.Proteins[i].RazorLabels.Channel5.StDev, _ = stats.StandardDeviationPopulation(c5Data)
		evi.Proteins[i].RazorLabels.Channel6.StDev, _ = stats.StandardDeviationPopulation(c6Data)
		evi.Proteins[i].RazorLabels.Channel7.StDev, _ = stats.StandardDeviationPopulation(c7Data)
		evi.Proteins[i].RazorLabels.Channel8.StDev, _ = stats.StandardDeviationPopulation(c8Data)
		evi.Proteins[i].RazorLabels.Channel9.StDev, _ = stats.StandardDeviationPopulation(c9Data)
		evi.Proteins[i].RazorLabels.Channel10.StDev, _ = stats.StandardDeviationPopulation(c10Data)

	}

	return evi, nil
}

func ratioToIntensityMean(evi rep.Evidence) (rep.Evidence, error) {

	for i := range evi.Proteins {

		var totalRef float64
		var uniqRef float64
		var razorRef float64

		totalRef += evi.Proteins[i].TotalLabels.Channel1.Mean
		totalRef += evi.Proteins[i].TotalLabels.Channel2.Mean
		totalRef += evi.Proteins[i].TotalLabels.Channel3.Mean
		totalRef += evi.Proteins[i].TotalLabels.Channel4.Mean
		totalRef += evi.Proteins[i].TotalLabels.Channel5.Mean
		totalRef += evi.Proteins[i].TotalLabels.Channel6.Mean
		totalRef += evi.Proteins[i].TotalLabels.Channel7.Mean
		totalRef += evi.Proteins[i].TotalLabels.Channel8.Mean
		totalRef += evi.Proteins[i].TotalLabels.Channel9.Mean
		totalRef += evi.Proteins[i].TotalLabels.Channel10.Mean

		uniqRef += evi.Proteins[i].UniqueLabels.Channel1.Mean
		uniqRef += evi.Proteins[i].UniqueLabels.Channel2.Mean
		uniqRef += evi.Proteins[i].UniqueLabels.Channel3.Mean
		uniqRef += evi.Proteins[i].UniqueLabels.Channel4.Mean
		uniqRef += evi.Proteins[i].UniqueLabels.Channel5.Mean
		uniqRef += evi.Proteins[i].UniqueLabels.Channel6.Mean
		uniqRef += evi.Proteins[i].UniqueLabels.Channel7.Mean
		uniqRef += evi.Proteins[i].UniqueLabels.Channel8.Mean
		uniqRef += evi.Proteins[i].UniqueLabels.Channel9.Mean
		uniqRef += evi.Proteins[i].UniqueLabels.Channel10.Mean

		razorRef += evi.Proteins[i].RazorLabels.Channel1.Mean
		razorRef += evi.Proteins[i].RazorLabels.Channel2.Mean
		razorRef += evi.Proteins[i].RazorLabels.Channel3.Mean
		razorRef += evi.Proteins[i].RazorLabels.Channel4.Mean
		razorRef += evi.Proteins[i].RazorLabels.Channel5.Mean
		razorRef += evi.Proteins[i].RazorLabels.Channel6.Mean
		razorRef += evi.Proteins[i].RazorLabels.Channel7.Mean
		razorRef += evi.Proteins[i].RazorLabels.Channel8.Mean
		razorRef += evi.Proteins[i].RazorLabels.Channel9.Mean
		razorRef += evi.Proteins[i].RazorLabels.Channel10.Mean

		evi.Proteins[i].TotalLabels.Channel1.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel1.Mean/totalRef), 4, 5) * 100)
		evi.Proteins[i].TotalLabels.Channel2.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel2.Mean/totalRef), 4, 5) * 100)
		evi.Proteins[i].TotalLabels.Channel3.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel3.Mean/totalRef), 4, 5) * 100)
		evi.Proteins[i].TotalLabels.Channel4.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel4.Mean/totalRef), 4, 5) * 100)
		evi.Proteins[i].TotalLabels.Channel5.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel5.Mean/totalRef), 4, 5) * 100)
		evi.Proteins[i].TotalLabels.Channel6.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel6.Mean/totalRef), 4, 5) * 100)
		evi.Proteins[i].TotalLabels.Channel7.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel7.Mean/totalRef), 4, 5) * 100)
		evi.Proteins[i].TotalLabels.Channel8.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel8.Mean/totalRef), 4, 5) * 100)
		evi.Proteins[i].TotalLabels.Channel9.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel9.Mean/totalRef), 4, 5) * 100)
		evi.Proteins[i].TotalLabels.Channel10.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel10.Mean/totalRef), 4, 5) * 100)

		evi.Proteins[i].UniqueLabels.Channel1.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel1.Mean/uniqRef), 4, 5) * 100)
		evi.Proteins[i].UniqueLabels.Channel2.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel2.Mean/uniqRef), 4, 5) * 100)
		evi.Proteins[i].UniqueLabels.Channel3.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel3.Mean/uniqRef), 4, 5) * 100)
		evi.Proteins[i].UniqueLabels.Channel4.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel4.Mean/uniqRef), 4, 5) * 100)
		evi.Proteins[i].UniqueLabels.Channel5.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel5.Mean/uniqRef), 4, 5) * 100)
		evi.Proteins[i].UniqueLabels.Channel6.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel6.Mean/uniqRef), 4, 5) * 100)
		evi.Proteins[i].UniqueLabels.Channel7.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel7.Mean/uniqRef), 4, 5) * 100)
		evi.Proteins[i].UniqueLabels.Channel8.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel8.Mean/uniqRef), 4, 5) * 100)
		evi.Proteins[i].UniqueLabels.Channel9.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel9.Mean/uniqRef), 4, 5) * 100)
		evi.Proteins[i].UniqueLabels.Channel10.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel10.Mean/uniqRef), 4, 5) * 100)

		evi.Proteins[i].RazorLabels.Channel1.RatioIntensity = (utils.Round((evi.Proteins[i].RazorLabels.Channel1.Mean/razorRef), 4, 5) * 100)
		evi.Proteins[i].RazorLabels.Channel2.RatioIntensity = (utils.Round((evi.Proteins[i].RazorLabels.Channel2.Mean/razorRef), 4, 5) * 100)
		evi.Proteins[i].RazorLabels.Channel3.RatioIntensity = (utils.Round((evi.Proteins[i].RazorLabels.Channel3.Mean/razorRef), 4, 5) * 100)
		evi.Proteins[i].RazorLabels.Channel4.RatioIntensity = (utils.Round((evi.Proteins[i].RazorLabels.Channel4.Mean/razorRef), 4, 5) * 100)
		evi.Proteins[i].RazorLabels.Channel5.RatioIntensity = (utils.Round((evi.Proteins[i].RazorLabels.Channel5.Mean/razorRef), 4, 5) * 100)
		evi.Proteins[i].RazorLabels.Channel6.RatioIntensity = (utils.Round((evi.Proteins[i].RazorLabels.Channel6.Mean/razorRef), 4, 5) * 100)
		evi.Proteins[i].RazorLabels.Channel7.RatioIntensity = (utils.Round((evi.Proteins[i].RazorLabels.Channel7.Mean/razorRef), 4, 5) * 100)
		evi.Proteins[i].RazorLabels.Channel8.RatioIntensity = (utils.Round((evi.Proteins[i].RazorLabels.Channel8.Mean/razorRef), 4, 5) * 100)
		evi.Proteins[i].RazorLabels.Channel9.RatioIntensity = (utils.Round((evi.Proteins[i].RazorLabels.Channel9.Mean/razorRef), 4, 5) * 100)
		evi.Proteins[i].RazorLabels.Channel10.RatioIntensity = (utils.Round((evi.Proteins[i].RazorLabels.Channel10.Mean/razorRef), 4, 5) * 100)

	}

	return evi, nil
}

func ratioToControlChannel(evi rep.Evidence, control string) (rep.Evidence, error) {

	for i := range evi.Proteins {

		var totalRef float64
		var uniqRef float64
		var razorRef float64

		switch control {
		case "1":
			totalRef = evi.Proteins[i].TotalLabels.Channel1.Mean
			uniqRef = evi.Proteins[i].UniqueLabels.Channel1.Mean
			razorRef = evi.Proteins[i].RazorLabels.Channel1.Mean
		case "2":
			totalRef = evi.Proteins[i].TotalLabels.Channel2.Mean
			uniqRef = evi.Proteins[i].UniqueLabels.Channel2.Mean
			razorRef = evi.Proteins[i].RazorLabels.Channel2.Mean
		case "3":
			totalRef = evi.Proteins[i].TotalLabels.Channel3.Mean
			uniqRef = evi.Proteins[i].UniqueLabels.Channel3.Mean
			razorRef = evi.Proteins[i].RazorLabels.Channel3.Mean
		case "4":
			totalRef = evi.Proteins[i].TotalLabels.Channel4.Mean
			uniqRef = evi.Proteins[i].UniqueLabels.Channel4.Mean
			razorRef = evi.Proteins[i].RazorLabels.Channel4.Mean
		case "5":
			totalRef = evi.Proteins[i].TotalLabels.Channel5.Mean
			uniqRef = evi.Proteins[i].UniqueLabels.Channel5.Mean
			razorRef = evi.Proteins[i].RazorLabels.Channel5.Mean
		case "6":
			totalRef = evi.Proteins[i].TotalLabels.Channel6.Mean
			uniqRef = evi.Proteins[i].UniqueLabels.Channel6.Mean
			razorRef = evi.Proteins[i].RazorLabels.Channel6.Mean
		case "7":
			totalRef = evi.Proteins[i].TotalLabels.Channel7.Mean
			uniqRef = evi.Proteins[i].UniqueLabels.Channel7.Mean
			razorRef = evi.Proteins[i].RazorLabels.Channel7.Mean
		case "8":
			totalRef = evi.Proteins[i].TotalLabels.Channel8.Mean
			uniqRef = evi.Proteins[i].UniqueLabels.Channel8.Mean
			razorRef = evi.Proteins[i].RazorLabels.Channel8.Mean
		case "9":
			totalRef = evi.Proteins[i].TotalLabels.Channel9.Mean
			uniqRef = evi.Proteins[i].UniqueLabels.Channel9.Mean
			razorRef = evi.Proteins[i].RazorLabels.Channel9.Mean
		case "10":
			totalRef = evi.Proteins[i].TotalLabels.Channel10.Mean
			uniqRef = evi.Proteins[i].UniqueLabels.Channel10.Mean
			razorRef = evi.Proteins[i].RazorLabels.Channel10.Mean
		default:
			return evi, errors.New("Cant find the given channel for normalization")
		}

		evi.Proteins[i].TotalLabels.Channel1.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel1.Mean/totalRef), 4, 5) * 100)
		evi.Proteins[i].TotalLabels.Channel2.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel2.Mean/totalRef), 4, 5) * 100)
		evi.Proteins[i].TotalLabels.Channel3.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel3.Mean/totalRef), 4, 5) * 100)
		evi.Proteins[i].TotalLabels.Channel4.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel4.Mean/totalRef), 4, 5) * 100)
		evi.Proteins[i].TotalLabels.Channel5.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel5.Mean/totalRef), 4, 5) * 100)
		evi.Proteins[i].TotalLabels.Channel6.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel6.Mean/totalRef), 4, 5) * 100)
		evi.Proteins[i].TotalLabels.Channel7.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel7.Mean/totalRef), 4, 5) * 100)
		evi.Proteins[i].TotalLabels.Channel8.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel8.Mean/totalRef), 4, 5) * 100)
		evi.Proteins[i].TotalLabels.Channel9.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel9.Mean/totalRef), 4, 5) * 100)
		evi.Proteins[i].TotalLabels.Channel10.RatioIntensity = (utils.Round((evi.Proteins[i].TotalLabels.Channel10.Mean/totalRef), 4, 5) * 100)

		evi.Proteins[i].UniqueLabels.Channel1.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel1.Mean/uniqRef), 4, 5) * 100)
		evi.Proteins[i].UniqueLabels.Channel2.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel2.Mean/uniqRef), 4, 5) * 100)
		evi.Proteins[i].UniqueLabels.Channel3.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel3.Mean/uniqRef), 4, 5) * 100)
		evi.Proteins[i].UniqueLabels.Channel4.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel4.Mean/uniqRef), 4, 5) * 100)
		evi.Proteins[i].UniqueLabels.Channel5.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel5.Mean/uniqRef), 4, 5) * 100)
		evi.Proteins[i].UniqueLabels.Channel6.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel6.Mean/uniqRef), 4, 5) * 100)
		evi.Proteins[i].UniqueLabels.Channel7.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel7.Mean/uniqRef), 4, 5) * 100)
		evi.Proteins[i].UniqueLabels.Channel8.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel8.Mean/uniqRef), 4, 5) * 100)
		evi.Proteins[i].UniqueLabels.Channel9.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel9.Mean/uniqRef), 4, 5) * 100)
		evi.Proteins[i].UniqueLabels.Channel10.RatioIntensity = (utils.Round((evi.Proteins[i].UniqueLabels.Channel10.Mean/uniqRef), 4, 5) * 100)

		evi.Proteins[i].RazorLabels.Channel1.RatioIntensity = (utils.Round((evi.Proteins[i].RazorLabels.Channel1.Mean/razorRef), 4, 5) * 100)
		evi.Proteins[i].RazorLabels.Channel2.RatioIntensity = (utils.Round((evi.Proteins[i].RazorLabels.Channel2.Mean/razorRef), 4, 5) * 100)
		evi.Proteins[i].RazorLabels.Channel3.RatioIntensity = (utils.Round((evi.Proteins[i].RazorLabels.Channel3.Mean/razorRef), 4, 5) * 100)
		evi.Proteins[i].RazorLabels.Channel4.RatioIntensity = (utils.Round((evi.Proteins[i].RazorLabels.Channel4.Mean/razorRef), 4, 5) * 100)
		evi.Proteins[i].RazorLabels.Channel5.RatioIntensity = (utils.Round((evi.Proteins[i].RazorLabels.Channel5.Mean/razorRef), 4, 5) * 100)
		evi.Proteins[i].RazorLabels.Channel6.RatioIntensity = (utils.Round((evi.Proteins[i].RazorLabels.Channel6.Mean/razorRef), 4, 5) * 100)
		evi.Proteins[i].RazorLabels.Channel7.RatioIntensity = (utils.Round((evi.Proteins[i].RazorLabels.Channel7.Mean/razorRef), 4, 5) * 100)
		evi.Proteins[i].RazorLabels.Channel8.RatioIntensity = (utils.Round((evi.Proteins[i].RazorLabels.Channel8.Mean/razorRef), 4, 5) * 100)
		evi.Proteins[i].RazorLabels.Channel9.RatioIntensity = (utils.Round((evi.Proteins[i].RazorLabels.Channel9.Mean/razorRef), 4, 5) * 100)
		evi.Proteins[i].RazorLabels.Channel10.RatioIntensity = (utils.Round((evi.Proteins[i].RazorLabels.Channel10.Mean/razorRef), 4, 5) * 100)

	}

	return evi, nil
}

// peakIntensity ...
func peakIntensity(e rep.Evidence, dir, format string, rTWin, pTWin, tol float64) (rep.Evidence, error) {

	// get all spectra in centralized structure
	logrus.Info("Reading spectra")
	ms1Map, err := getMS1Spectra(dir, format, e.PSM)
	if err != nil {
		return e, err
	}

	logrus.Info("Tracing Peaks")
	for i := range e.PSM {

		// process pepXML information
		ppmPrecision := tol / math.Pow(10, 6)
		mz := utils.Round(((e.PSM[i].PrecursorNeutralMass + (float64(e.PSM[i].AssumedCharge) * bio.Proton)) / float64(e.PSM[i].AssumedCharge)), 5, 4)
		minRT := (e.PSM[i].RetentionTime / 60) - rTWin
		maxRT := (e.PSM[i].RetentionTime / 60) + rTWin

		var measured = make(map[float64]float64)
		var retrieved bool

		// XIC on MS1 level
		for _, j := range ms1Map {
			measured, retrieved = xic(j, minRT, maxRT, ppmPrecision, mz)
		}

		if retrieved == true {
			var timeW = e.PSM[i].RetentionTime / 60
			var topI = 0.0

			for k, v := range measured {
				if k > (timeW-pTWin) && k < (timeW+pTWin) {
					if v > topI {
						topI = v
					}
				}
			}
			e.PSM[i].Intensity = topI
		}

	}

	return e, nil
}

func getSpectra(path, format string, evi rep.Evidence) (map[string]mz.MS1, map[string]mz.MS2, error) {

	var err error
	var ms1 = make(map[string]mz.MS1)
	var ms2 = make(map[string]mz.MS2)
	var mzs = make(map[string]int)

	// get specta file names from identifications
	for _, i := range evi.PSM {
		specName := strings.Split(i.Spectrum, ".")
		source := fmt.Sprintf("%s.%s", specName[0], format)
		mzs[source]++
	}

	for k := range mzs {

		ext := filepath.Ext(k)
		name := filepath.Base(k)
		clean := name[0 : len(name)-len(ext)]
		fullpath, _ := filepath.Abs(path)
		name = fmt.Sprintf("%s%s%s", fullpath, string(filepath.Separator), name)

		if strings.Contains(k, "mzML") {

			var ms mzml.IndexedMzML
			err = ms.Parse(name)
			if err != nil {
				return ms1, ms2, err
			}

			ms1Scans, ms1err := mz.GetMzMLSpectra(ms, clean)
			if ms1err != nil {
				return nil, nil, ms1err
			}
			ms1[clean] = ms1Scans

			ms2Scans, ms2err := mz.GetMzMLMS2Spectra(ms, clean)
			if ms2err != nil {
				return nil, nil, ms2err
			}
			ms2[clean] = ms2Scans

		} else if strings.Contains(k, "mzXML") {

			var ms mzxml.MzXML
			err = ms.Parse(name)
			if err != nil {
				return nil, nil, err
			}

			ms1Scans, err := mz.GetMzXMLSpectra(ms, clean)
			if err != nil {
				return nil, nil, err
			}
			ms1[clean] = ms1Scans

			ms2Scans, err := mz.GetMzXMLMS2Spectra(ms, clean)
			if err != nil {
				return nil, nil, err
			}
			ms2[clean] = ms2Scans

		} else {
			logrus.Fatal("Cannot open file: ", name)
		}

	}

	return ms1, ms2, nil
}

// getMS1Spectra gets MS1 infor from spectra files
func getMS1Spectra(path, format string, pep rep.PSMEvidenceList) (map[string][]mz.Ms1Scan, error) {

	var err error
	// get the name of all raw files used in the experiment from pepxml
	var spec = make(map[string][]mz.Ms1Scan)
	var mzs = make(map[string]int)

	for _, i := range pep {
		specName := strings.Split(i.Spectrum, ".")
		source := fmt.Sprintf("%s.%s", specName[0], format)
		mzs[source]++
	}

	for k := range mzs {

		ext := filepath.Ext(k)
		name := filepath.Base(k)
		clean := name[0 : len(name)-len(ext)]
		fullpath, _ := filepath.Abs(path)
		name = fmt.Sprintf("%s%s%s", fullpath, string(filepath.Separator), name)

		var ms1Reader mz.MS1

		if strings.Contains(k, "mzML") {
			err = ms1Reader.ReadMzML(name)
			if err != nil {
				return spec, err
			}
		} else if strings.Contains(k, "mzXML") {
			err = ms1Reader.ReadMzXML(name)
			if err != nil {
				return spec, err
			}
		} else {
			logrus.Fatal("Cannot open file: ", name)
		}

		spec[clean] = ms1Reader.Ms1Scan

	}

	return spec, nil
}

// getMS2Spectra gets MS1 infor from spectra files
func getMS2Spectra(path, format string, pep rep.PSMEvidenceList) (map[string]map[string]mz.Ms2Scan, error) {

	// get the name of all raw files used in the experiment from pepxml
	spec := make(map[string]map[string]mz.Ms2Scan)
	var mzs = make(map[string]int)

	for _, i := range pep {
		specName := strings.Split(i.Spectrum, ".")
		source := fmt.Sprintf("%s.%s", specName[0], format)
		mzs[source]++
	}

	for k := range mzs {

		ext := filepath.Ext(k)
		name := filepath.Base(k)
		clean := name[0 : len(name)-len(ext)]
		fullpath, _ := filepath.Abs(path)
		name = fmt.Sprintf("%s%s%s", fullpath, string(filepath.Separator), name)

		if strings.Contains(k, "mzML") {

			var ms2Reader mz.MS2
			err := ms2Reader.ReadMzML(name)
			if err != nil {
				return spec, err
			}
			spec[clean] = ms2Reader.Ms2Scan

		} else if strings.Contains(k, "mzXML") {

			var ms2Reader mz.MS2
			ms2Reader.ReadMzXML(name)
			spec[clean] = ms2Reader.Ms2Scan

		} else {
			logrus.Fatal("Cannot open file: ", name)
		}
	}

	return spec, nil
}

// xic extract ion chomatograms
func xic(v []mz.Ms1Scan, minRT, maxRT, ppmPrecision, mz float64) (map[float64]float64, bool) {

	var list = make(map[float64]float64)

	for j := range v {

		if v[j].ScanStartTime > minRT && v[j].ScanStartTime < maxRT {

			lowi := sort.Search(len(v[j].Spectrum), func(i int) bool { return v[j].Spectrum[i].Mz >= mz-ppmPrecision*mz })
			highi := sort.Search(len(v[j].Spectrum), func(i int) bool { return v[j].Spectrum[i].Mz >= mz+ppmPrecision*mz })

			var maxI = 0.0

			for _, k := range v[j].Spectrum[lowi:highi] {
				if k.Intensity > maxI {
					maxI = k.Intensity
				}
			}

			if maxI > 0 {
				list[v[j].ScanStartTime] = maxI
			}

		}
	}

	if len(list) >= 5 {
		return list, true
	}

	return list, false
}

// CalculateSpectralCounts add Spc to ions and proteins
func CalculateSpectralCounts(e rep.Evidence) (rep.Evidence, error) {

	var spcMap = make(map[string]int)
	var ionRefMap = make(map[string]int)

	if len(e.PSM) < 1 || len(e.Ions) < 1 {
		return e, errors.New("No peptide identification found")
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

// calculateIntensities calculates the protein intensity
func calculateIntensities(e rep.Evidence) (rep.Evidence, error) {

	var intMap = make(map[string]float64)
	var intRefMap = make(map[string]float64)
	var intPepMap = make(map[string]float64)

	if len(e.PSM) < 1 || len(e.Ions) < 1 {
		return e, errors.New("No peptide identification found")
	}

	for i := range e.PSM {

		var key string
		if len(e.PSM[i].ModifiedPeptide) > 0 {
			key = fmt.Sprintf("%s#%d", e.PSM[i].ModifiedPeptide, e.PSM[i].AssumedCharge)
		} else {
			key = fmt.Sprintf("%s#%d", e.PSM[i].Peptide, e.PSM[i].AssumedCharge)
		}

		_, ok := intMap[key]
		if ok {
			if e.PSM[i].Intensity > intMap[key] {
				intMap[key] = e.PSM[i].Intensity
			}
		} else {
			intMap[key] = e.PSM[i].Intensity
		}

		_, okPep := intPepMap[e.PSM[i].Peptide]
		if okPep {
			if e.PSM[i].Intensity > intPepMap[e.PSM[i].Peptide] {
				intPepMap[e.PSM[i].Peptide] = e.PSM[i].Intensity
			}
		} else {
			intPepMap[e.PSM[i].Peptide] = e.PSM[i].Intensity
		}

	}

	// peptides get the higest intense from the matching sequences
	for i := range e.Peptides {
		v, ok := intMap[e.Peptides[i].Sequence]
		if ok {
			e.Peptides[i].Intensity = v
		}
	}

	// ions get the higest intense from the matching sequences
	for i := range e.Ions {

		var key string
		if len(e.Ions[i].ModifiedSequence) > 0 {
			key = fmt.Sprintf("%s#%d", e.Ions[i].ModifiedSequence, e.Ions[i].ChargeState)
		} else {
			key = fmt.Sprintf("%s#%d", e.Ions[i].Sequence, e.Ions[i].ChargeState)
		}

		v, ok := intMap[key]
		if ok {
			e.Ions[i].Intensity = v
			intRefMap[key] = v
		}

	}

	for i := range e.Proteins {

		var totalInt []float64
		var uniqueInt []float64
		var razorInt []float64

		// make a reference for razor peptides
		var uniqIons = make(map[string]uint8)

		for k := range e.Proteins[i].UniquePeptideIons {
			v, ok := intRefMap[k]
			if ok {
				uniqueInt = append(uniqueInt, v)
				razorInt = append(razorInt, v)
				uniqIons[k] = 0
			}
		}

		for k, j := range e.Proteins[i].TotalPeptideIons {
			v, ok := intRefMap[k]
			if ok {
				totalInt = append(totalInt, v)
				if j.IsRazor {
					_, ok := uniqIons[k]
					if !ok {
						razorInt = append(razorInt, v)
					}
				}
			}
		}

		sort.Float64s(uniqueInt)
		sort.Float64s(totalInt)
		sort.Float64s(razorInt)

		if len(uniqueInt) >= 3 {
			e.Proteins[i].UniqueIntensity = (uniqueInt[len(uniqueInt)-1] + uniqueInt[len(uniqueInt)-2] + uniqueInt[len(uniqueInt)-3])
		} else if len(uniqueInt) == 2 {
			e.Proteins[i].UniqueIntensity = (uniqueInt[len(uniqueInt)-1] + uniqueInt[len(uniqueInt)-2])
		} else if len(uniqueInt) == 1 {
			e.Proteins[i].UniqueIntensity = (uniqueInt[len(uniqueInt)-1])
		}

		if len(totalInt) >= 3 {
			e.Proteins[i].TotalIntensity = (totalInt[len(totalInt)-1] + totalInt[len(totalInt)-2] + totalInt[len(totalInt)-3])
		} else if len(totalInt) == 2 {
			e.Proteins[i].TotalIntensity = (totalInt[len(totalInt)-1] + totalInt[len(totalInt)-2])
		} else if len(totalInt) == 1 {
			e.Proteins[i].TotalIntensity = (totalInt[len(totalInt)-1])
		}

		if len(razorInt) >= 3 {
			e.Proteins[i].RazorIntensity = (razorInt[len(razorInt)-1] + razorInt[len(razorInt)-2] + razorInt[len(razorInt)-3])
		} else if len(razorInt) == 2 {
			e.Proteins[i].RazorIntensity = (razorInt[len(razorInt)-1] + razorInt[len(razorInt)-2])
		} else if len(razorInt) == 1 {
			e.Proteins[i].RazorIntensity = (razorInt[len(razorInt)-1])
		}

	}

	return e, nil
}
