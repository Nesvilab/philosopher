package qua

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/prvst/philosopher/lib/bio"
	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/raw"
	"github.com/prvst/philosopher/lib/rep"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/sirupsen/logrus"
)

// // peakIntensity collects PSM intensities from the apex peak
func peakIntensity(evi rep.Evidence, dir, format string, rTWin, pTWin, tol float64, isIso bool) (rep.Evidence, *err.Error) {

	logrus.Info("Indexing PSM information")
	var sourceMap = make(map[string]uint8)
	var spectra = make(map[string][]string)
	var ppmPrecision = make(map[string]float64)
	var mz = make(map[string]float64)
	var minRT = make(map[string]float64)
	var maxRT = make(map[string]float64)
	var retentionTime = make(map[string]float64)
	var intensity = make(map[string]float64)

	var isolatedMz = make(map[string]float64)
	var isolatedUW = make(map[string]float64)
	var isolatedLW = make(map[string]float64)

	for _, i := range evi.PSM {

		partName := strings.Split(i.Spectrum, ".")
		sourceMap[partName[0]] = 0
		spectra[partName[0]] = append(spectra[partName[0]], i.Spectrum)

		ppmPrecision[i.Spectrum] = tol / math.Pow(10, 6)
		mz[i.Spectrum] = ((i.PrecursorNeutralMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge))
		minRT[i.Spectrum] = (i.RetentionTime / 60) - rTWin
		maxRT[i.Spectrum] = (i.RetentionTime / 60) + rTWin
		retentionTime[i.Spectrum] = i.RetentionTime
	}

	logrus.Info("Reading spectra and tracing peaks")
	for source := range sourceMap {

		logrus.Info("Processing ", source)
		var ms1 raw.MS1
		var ms2 raw.MS2

		// get the clean name, remove the extension
		var extension = filepath.Ext(filepath.Base(source))
		var name = source[0 : len(source)-len(extension)]
		input := fmt.Sprintf("%s%s%s.bin", sys.MetaDir(), string(filepath.Separator), name)

		// get all MS1 spectra (and MS2 if necessary)
		if _, e := os.Stat(input); e == nil {

			spec, e := raw.Restore(source)
			if e != nil {
				return evi, &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: "error restoring indexed mz"}
			}

			ms1 = raw.GetMS1(spec)

			if isIso == true {
				ms2 = raw.GetMS2(spec)
			}

		} else {

			spec, rer := raw.RestoreFromFile(dir, source, format)
			if rer != nil {
				return evi, &err.Error{Type: err.CannotParseXML, Class: err.FATA, Argument: "cant read mz file"}
			}

			ms1 = raw.GetMS1(spec)

			if isIso == true {
				ms2 = raw.GetMS2(spec)
			}

		}

		for _, i := range ms2.Ms2Scan {
			fmt.Println(source, "\t", i.Index, "\t", i.Scan, "\t", i.Precursor.ParentIndex, "\t", i.Precursor.ParentScan, "\t", i.Precursor.SelectedIon, "\t", i.Precursor.TargetIon)
		}

		if isIso == true {
			for _, i := range ms2.Ms2Scan {
				_, ok := isolatedMz[i.Precursor.ParentScan]
				if !ok {
					isolatedMz[i.Precursor.ParentScan] = i.Precursor.TargetIon
					isolatedUW[i.Precursor.ParentScan] = i.Precursor.IsolationWindowUpperOffset
					isolatedLW[i.Precursor.ParentScan] = i.Precursor.IsolationWindowLowerOffset
				}
			}
		}

		v, ok := spectra[source]
		if ok {
			for _, j := range v {

				var measured = make(map[float64]float64)
				var retrieved bool

				measured, retrieved = xic(ms1.Ms1Scan, isolatedMz, isolatedUW, isolatedLW, minRT[j], maxRT[j], ppmPrecision[j], mz[j], isIso)

				if retrieved == true {
					var timeW = retentionTime[j] / 60
					var topI = 0.0

					for k, v := range measured {
						if k > (timeW-pTWin) && k < (timeW+pTWin) {
							if v > topI {
								topI = v
							}
						}
					}

					intensity[j] = topI
				}
			}
		}

	}

	for i := range evi.PSM {
		partName := strings.Split(evi.PSM[i].Spectrum, ".")
		_, ok := spectra[partName[0]]
		if ok {
			evi.PSM[i].Intensity = intensity[evi.PSM[i].Spectrum]
		}
	}

	return evi, nil
}

// xic extract ion chomatograms
func xic(ms1 []raw.Ms1Scan, isolatedIMz, isolatedUW, isolatedLW map[string]float64, minRT, maxRT, ppmPrecision, mz float64, isIso bool) (map[float64]float64, bool) {

	var list = make(map[float64]float64)

	for j := range ms1 {

		if isIso == true {
			mz = isolatedIMz[ms1[j].Scan]
		}

		if ms1[j].ScanStartTime >= minRT && ms1[j].ScanStartTime <= maxRT {

			lowi := sort.Search(len(ms1[j].Spectrum), func(i int) bool { return ms1[j].Spectrum[i].Mz >= mz-ppmPrecision*mz })
			highi := sort.Search(len(ms1[j].Spectrum), func(i int) bool { return ms1[j].Spectrum[i].Mz >= mz+ppmPrecision*mz })

			var maxI = 0.0

			for _, k := range ms1[j].Spectrum[lowi:highi] {
				if k.Intensity > maxI {
					maxI = k.Intensity
				}
			}

			if maxI > 0 {
				list[ms1[j].ScanStartTime] = maxI
			}

		}
	}

	if len(list) >= 5 {
		return list, true
	}

	return list, false
}

func calculateIntensities(e rep.Evidence) (rep.Evidence, *err.Error) {

	logrus.Info("Assigning intensities to data layers")

	if len(e.PSM) < 1 || len(e.Ions) < 1 {
		return e, &err.Error{Type: err.CannotFindPSMData, Class: err.FATA, Argument: "no PSMs or Ions found, cannot attribute intensity calculations"}
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

				if k.IsUnique == true {
					uniqueInt = append(uniqueInt, v)
				}

				if k.IsURazor == true {
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

	return e, nil
}
