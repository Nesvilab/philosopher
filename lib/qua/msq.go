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

// peakIntensity collects PSM intensities from the apex peak
func peakIntensity(evi rep.Evidence, dir, format string, rTWin, pTWin, tol float64) (rep.Evidence, *err.Error) {

	// collect all source file names present on the PSM list
	var sourceMap = make(map[string]uint8)
	for _, i := range evi.PSM {
		specName := strings.Split(i.Spectrum, ".")
		sourceMap[specName[0]] = 0
	}

	logrus.Info("Reading spectra and tracing peaks")
	for k := range sourceMap {

		var ms1 raw.MS1

		// get the clean name, remove the extension
		var extension = filepath.Ext(filepath.Base(k))
		var name = k[0 : len(k)-len(extension)]
		input := fmt.Sprintf("%s%s%s.bin", sys.MetaDir(), string(filepath.Separator), name)

		// get all MS1 spectra
		if _, e := os.Stat(input); e == nil {

			spec, e := raw.Restore(k)
			if e != nil {
				return evi, &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: "error restoring indexed mz"}
			}

			ms1 = raw.GetMS1(spec)

		} else {

			spec, rer := raw.RestoreFromFile(dir, k, format)
			if rer != nil {
				return evi, &err.Error{Type: err.CannotParseXML, Class: err.FATA, Argument: "cant read mz file"}
			}

			ms1 = raw.GetMS1(spec)
		}

		for i := range evi.PSM {

			// process pepXML information using the experimental mass to calculate the mz
			ppmPrecision := tol / math.Pow(10, 6)
			mz := ((evi.PSM[i].PrecursorNeutralMass + (float64(evi.PSM[i].AssumedCharge) * bio.Proton)) / float64(evi.PSM[i].AssumedCharge))
			minRT := (evi.PSM[i].RetentionTime / 60) - rTWin
			maxRT := (evi.PSM[i].RetentionTime / 60) + rTWin

			var measured = make(map[float64]float64)
			var retrieved bool

			// XIC on MS1 level
			if strings.Contains(evi.PSM[i].Spectrum, k) {
				measured, retrieved = xic(ms1.Ms1Scan, minRT, maxRT, ppmPrecision, mz)
			}

			if retrieved == true {
				var timeW = evi.PSM[i].RetentionTime / 60
				var topI = 0.0

				for k, v := range measured {
					if k > (timeW-pTWin) && k < (timeW+pTWin) {
						if v > topI {
							topI = v
						}
					}
				}

				evi.PSM[i].Intensity = topI
			}

		}

	}

	return evi, nil
}

// xic extract ion chomatograms
func xic(v []raw.Ms1Scan, minRT, maxRT, ppmPrecision, mz float64) (map[float64]float64, bool) {

	var list = make(map[float64]float64)

	for j := range v {

		if v[j].ScanStartTime >= minRT && v[j].ScanStartTime <= maxRT {

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

func calculateIntensities(e rep.Evidence) (rep.Evidence, *err.Error) {

	var intPepKeyMap = make(map[string]float64)
	var intPepMap = make(map[string]float64)
	var intIonMap = make(map[string]float64)

	if len(e.PSM) < 1 || len(e.Ions) < 1 {
		return e, &err.Error{Type: err.CannotFindPSMData, Class: err.FATA, Argument: "cannot attribute intensity calculations"}
	}

	for i := range e.PSM {

		// global intensity map for Peptides, getting the sum of all intensities
		// pepetide intensity is calculated by grouping PSM by sequence and calculted MZ

		pepKey := fmt.Sprintf("%s#%f", e.PSM[i].Peptide, e.PSM[i].CalcNeutralPepMass)
		_, okPep := intPepKeyMap[pepKey]
		if okPep {
			if e.PSM[i].Intensity > intPepKeyMap[pepKey] {
				intPepKeyMap[pepKey] = e.PSM[i].Intensity
			}
		} else {
			intPepKeyMap[pepKey] = e.PSM[i].Intensity
		}

		for k, v := range intPepKeyMap {
			key := strings.Split(k, "#")

			_, okPep := intPepMap[key[0]]
			if okPep {
				if e.PSM[i].Intensity > intPepMap[key[0]] {
					intPepMap[key[0]] = v
				}
			} else {
				intPepMap[key[0]] = v
			}

		}

		// global intensity map for Ions, getting the most intense ion
		_, okIon := intIonMap[e.PSM[i].IonForm]
		if okIon {
			if e.PSM[i].Intensity > intIonMap[e.PSM[i].IonForm] {
				intIonMap[e.PSM[i].IonForm] = e.PSM[i].Intensity
			}
		} else {
			intIonMap[e.PSM[i].IonForm] = e.PSM[i].Intensity
		}

	}

	// attribute intensities to peptide evidences
	for i := range e.Peptides {
		v, ok := intPepMap[e.Peptides[i].Sequence]
		if ok {
			e.Peptides[i].Intensity += v
		}
	}

	// attribute intensities to ion evidences
	for i := range e.Ions {
		v, ok := intIonMap[e.Ions[i].IonForm]
		if ok {
			e.Ions[i].Intensity = v
		}
	}

	// attribute intensities to protein evidences: getting the top 3 most intense ions
	for i := range e.Proteins {

		var totalInt []float64
		var uniqueInt []float64
		var razorInt []float64

		// for unique ions
		for _, k := range e.Proteins[i].TotalPeptideIons {
			v, ok := intIonMap[k.IonForm]
			if ok {
				//if k.IsNondegenerateEvidence == true {
				if k.IsUnique == true {
					uniqueInt = append(uniqueInt, v)
				}
			}
		}

		// for razor ions
		for _, k := range e.Proteins[i].TotalPeptideIons {
			v, ok := intIonMap[k.IonForm]
			if ok {
				if k.IsURazor == true {
					razorInt = append(razorInt, v)
				}
			}
		}

		// for total ions
		for k := range e.Proteins[i].TotalPeptideIons {
			v, ok := intIonMap[k]
			if ok {
				totalInt = append(totalInt, v)
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
