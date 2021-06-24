package qua

import (
	"errors"
	"fmt"
	"math"
	"path/filepath"
	"sort"
	"strings"

	"philosopher/lib/bio"
	"philosopher/lib/ext/rawfilereader"
	"philosopher/lib/msg"

	"philosopher/lib/mzn"
	"philosopher/lib/rep"

	"github.com/sirupsen/logrus"
)

// LFQ main structure
type LFQ struct {
	Intensities map[string]float64
}

// NewLFQ constructor
func NewLFQ() LFQ {

	var self LFQ

	self.Intensities = make(map[string]float64)

	return self
}

func peakIntensity(evi rep.Evidence, dir, format string, rTWin, pTWin, tol float64, isIso, isRaw, isFaims bool) rep.Evidence {

	logrus.Info("Indexing PSM information")

	var psmMap = make(map[string]rep.PSMEvidence)
	var sourceMap = make(map[string][]rep.PSMEvidence)
	var spectra = make(map[string][]string)
	var ppmPrecision = make(map[string]float64)
	var mzMap = make(map[string]float64)
	var mzCVMap = make(map[string]string)
	var minRT = make(map[string]float64)
	var maxRT = make(map[string]float64)
	var compVoltageMap = make(map[string]string)
	var retentionTime = make(map[string]float64)
	var intensity = make(map[string]float64)
	var instensityCV = make(map[string]float64)

	var charges = make(map[string]int)

	// collect attributes from PSM
	for _, i := range evi.PSM {
		partName := strings.Split(i.Spectrum, ".")
		sourceMap[partName[0]] = append(sourceMap[partName[0]], i)
		spectra[partName[0]] = append(spectra[partName[0]], i.Spectrum)

		ppmPrecision[i.Spectrum] = tol / math.Pow(10, 6)
		mzMap[i.Spectrum] = ((i.PrecursorNeutralMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge))
		minRT[i.Spectrum] = (i.RetentionTime / 60) - rTWin
		maxRT[i.Spectrum] = (i.RetentionTime / 60) + rTWin
		retentionTime[i.Spectrum] = i.RetentionTime
		compVoltageMap[i.Spectrum] = i.CompensationVoltage
		charges[i.Spectrum] = int(i.AssumedCharge)
		psmMap[i.Spectrum] = i
	}

	// get a sorted list of spectrum names
	var sourceList []string
	for i := range sourceMap {
		sourceList = append(sourceList, i)
	}

	sort.Strings(sourceList)

	logrus.Info("Reading spectra and tracing peaks")

	for _, s := range sourceList {

		logrus.Info("Processing ", s)
		var mz mzn.MsData
		var fileName string

		if isRaw {
			fileName = fmt.Sprintf("%s%s%s.raw", dir, string(filepath.Separator), s)
			stream := rawfilereader.Run(fileName, "")
			mz.ReadRaw(s, stream)
		} else {
			fileName = fmt.Sprintf("%s%s%s.mzML", dir, string(filepath.Separator), s)
			mz.Read(fileName)
		}

		for i := range mz.Spectra {

			spectrum := fmt.Sprintf("%s.%05s.%05s.%d", s, mz.Spectra[i].Scan, mz.Spectra[i].Scan, mz.Spectra[i].Precursor.ChargeState)

			if mz.Spectra[i].Level == "1" {
				if !isRaw {
					mz.Spectra[i].Decode()
				}

				if isFaims {
					mzCVMap[mz.Spectra[i].Scan] = mz.Spectra[i].CompensationVoltage
				}

			} else if mz.Spectra[i].Level == "2" {
				_, ok := mzMap[spectrum]
				if ok {
					mzMap[spectrum] = mz.Spectra[i].Precursor.TargetIon
				}
			}
		}

		mappedPurity := calculateIonPurity(dir, format, mz, sourceMap[s])

		for _, j := range mappedPurity {
			v, ok := psmMap[j.Spectrum]
			if ok {
				psm := v
				psm.Purity = j.Purity
				psmMap[j.Spectrum] = psm
			}
		}

		v, ok := spectra[s]
		if ok {
			for _, j := range v {

				measuredFaims, measured, retrieved := xic(mz.Spectra, minRT[j], maxRT[j], ppmPrecision[j], mzMap[j])

				if retrieved {

					var timeW = retentionTime[j] / 60
					var topI = 0.0
					var topCVI = 0.0
					var ms2CompensationVoltage = compVoltageMap[j]

					for k, v := range measured {

						if k > (timeW-pTWin) && k < (timeW+pTWin) {
							if v > topI {
								topI = v
							}
						}

						if isFaims {
							v1, ok := measuredFaims[ms2CompensationVoltage]
							if ok {
								if v1 > topCVI {
									topCVI = v1
								}
							}
						}
					}

					intensity[j] = topI
					instensityCV[j] = topCVI
				}
			}
		}
	}

	for i := range evi.PSM {
		partName := strings.Split(evi.PSM[i].Spectrum, ".")
		_, ok := spectra[partName[0]]
		if ok {
			evi.PSM[i].Intensity = intensity[evi.PSM[i].Spectrum]
			evi.PSM[i].IntensityCV = instensityCV[evi.PSM[i].Spectrum]
		}

		v, ok := psmMap[evi.PSM[i].Spectrum]
		if ok {
			evi.PSM[i].Purity = v.Purity
		}

	}

	return evi
}

// xic extract ion chomatograms
func xic(mz mzn.Spectra, minRT, maxRT, ppmPrecision, mzValue float64) (map[string]float64, map[float64]float64, bool) {

	var list = make(map[float64]float64)
	var ms1CompensationVoltage = make(map[string]float64)

	for j := range mz {
		if mz[j].Level == "1" {

			if mz[j].ScanStartTime >= minRT && mz[j].ScanStartTime <= maxRT {

				lowi := sort.Search(len(mz[j].Mz.DecodedStream), func(i int) bool { return mz[j].Mz.DecodedStream[i] >= mzValue-ppmPrecision*mzValue })
				highi := sort.Search(len(mz[j].Mz.DecodedStream), func(i int) bool { return mz[j].Mz.DecodedStream[i] >= mzValue+ppmPrecision*mzValue })

				var maxI = 0.0

				for _, k := range mz[j].Intensity.DecodedStream[lowi:highi] {
					if k > maxI {
						maxI = k
					}
				}

				if maxI > 0 {
					list[mz[j].ScanStartTime] = maxI
					ms1CompensationVoltage[mz[j].CompensationVoltage] = maxI
				}

			}
		}
	}

	if len(list) >= 5 {
		return ms1CompensationVoltage, list, true
	}

	return ms1CompensationVoltage, list, false
}

func calculateIntensities(e rep.Evidence) rep.Evidence {

	logrus.Info("Assigning intensities to data layers")

	if len(e.PSM) < 1 || len(e.Ions) < 1 {
		msg.QuantifyingData(errors.New("the PSM list is enpty"), "fatal")
	}

	var peptideIntMap = make(map[string]float64)
	var ionIntMap = make(map[string]float64)

	for _, i := range e.PSM {

		// peptide intensity : sum of all
		_, ok := peptideIntMap[i.Peptide]
		if ok {
			peptideIntMap[i.Peptide] += i.Intensity
		} else {
			peptideIntMap[i.Peptide] += i.Intensity
		}

		// ion intensity : most intense ion
		ionV, ok := ionIntMap[i.IonForm]
		if ok {
			if i.Intensity > ionV {
				ionIntMap[i.IonForm] = i.Intensity
			}
		} else {
			ionIntMap[i.IonForm] = i.Intensity
		}

	}

	for i := range e.Peptides {
		v, ok := peptideIntMap[e.Peptides[i].Sequence]
		if ok {
			e.Peptides[i].Intensity = v
		}
	}

	for i := range e.Ions {
		v, ok := ionIntMap[e.Ions[i].IonForm]
		if ok {
			e.Ions[i].Intensity = v
		}
	}

	// protein intensities : top 3 most intense ions
	for i := range e.Proteins {

		var totalInt []float64
		var uniqueInt []float64
		var razorInt []float64

		for _, k := range e.Proteins[i].TotalPeptideIons {
			v, ok := ionIntMap[k.IonForm]
			if ok {

				totalInt = append(totalInt, v)

				if k.IsUnique {
					uniqueInt = append(uniqueInt, v)
				}

				if k.IsURazor {
					razorInt = append(razorInt, v)
				}

			}
		}

		sort.Float64s(totalInt)
		sort.Float64s(uniqueInt)
		sort.Float64s(razorInt)

		if len(totalInt) >= 3 {
			e.Proteins[i].TotalIntensity = (totalInt[len(totalInt)-1] + totalInt[len(totalInt)-2] + totalInt[len(totalInt)-3])
		} else if len(totalInt) == 2 {
			e.Proteins[i].TotalIntensity = (totalInt[len(totalInt)-1] + totalInt[len(totalInt)-2])
		} else if len(totalInt) == 1 {
			e.Proteins[i].TotalIntensity = (totalInt[len(totalInt)-1])
		}

		if len(uniqueInt) >= 3 {
			e.Proteins[i].UniqueIntensity = (uniqueInt[len(uniqueInt)-1] + uniqueInt[len(uniqueInt)-2] + uniqueInt[len(uniqueInt)-3])
		} else if len(uniqueInt) == 2 {
			e.Proteins[i].UniqueIntensity = (uniqueInt[len(uniqueInt)-1] + uniqueInt[len(uniqueInt)-2])
		} else if len(uniqueInt) == 1 {
			e.Proteins[i].UniqueIntensity = (uniqueInt[len(uniqueInt)-1])
		}

		if len(razorInt) >= 3 {
			e.Proteins[i].URazorIntensity = (razorInt[len(razorInt)-1] + razorInt[len(razorInt)-2] + razorInt[len(razorInt)-3])
		} else if len(razorInt) == 2 {
			e.Proteins[i].URazorIntensity = (razorInt[len(razorInt)-1] + razorInt[len(razorInt)-2])
		} else if len(razorInt) == 1 {
			e.Proteins[i].URazorIntensity = (razorInt[len(razorInt)-1])
		}

	}

	return e
}
