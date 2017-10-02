package quan

import (
	"fmt"
	"math"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/cmsl/bio"
	"github.com/prvst/cmsl/data/mz"
	"github.com/prvst/cmsl/err"
	"github.com/prvst/philosopher/lib/rep"
)

// peakIntensity ...
func peakIntensity(e rep.Evidence, dir, format string, rTWin, pTWin, tol float64) (rep.Evidence, *err.Error) {

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
		//mz := utils.Round(((e.PSM[i].PrecursorNeutralMass + (float64(e.PSM[i].AssumedCharge) * bio.Proton)) / float64(e.PSM[i].AssumedCharge)), 5, 4)
		mz := ((e.PSM[i].PrecursorNeutralMass + (float64(e.PSM[i].AssumedCharge) * bio.Proton)) / float64(e.PSM[i].AssumedCharge))
		minRT := (e.PSM[i].RetentionTime / 60) - rTWin
		maxRT := (e.PSM[i].RetentionTime / 60) + rTWin

		var measured = make(map[float64]float64)
		var retrieved bool

		// XIC on MS1 level
		for k, j := range ms1Map {
			if strings.Contains(e.PSM[i].Spectrum, k) {
				measured, retrieved = xic(j, minRT, maxRT, ppmPrecision, mz)
			}
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

			// if topI > e.PSM[i].Intensity {
			// 	e.PSM[i].Intensity = topI
			// }

			e.PSM[i].Intensity = topI
		}

	}

	return e, nil
}

// getMS1Spectra gets MS1 infor from spectra files
func getMS1Spectra(path, format string, pep rep.PSMEvidenceList) (map[string][]mz.Ms1Scan, *err.Error) {

	// get the name of all raw files used in the experiment from pepxml
	var spec = make(map[string][]mz.Ms1Scan)
	var mzs = make(map[string]int)

	// collects all mz file names from identified spectra
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
			err := ms1Reader.ReadMzML(name)
			if err != nil {
				return spec, err
			}
		} else if strings.Contains(k, "mzXML") {
			err := ms1Reader.ReadMzXML(name)
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

		if v[j].ScanStartTime >= minRT && v[j].ScanStartTime <= maxRT {
			//if v[j].ScanStartTime >= minRT && v[j].ScanStartTime < maxRT {

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

// calculateIntensities calculates the protein intensity
func calculateIntensities(e rep.Evidence) (rep.Evidence, *err.Error) {

	var intPepMap = make(map[string]float64)
	var intIonMap = make(map[string]float64)

	if len(e.PSM) < 1 || len(e.Ions) < 1 {
		return e, &err.Error{Type: err.CannotFindPSMData, Class: err.FATA, Argument: "cannot attribute intensity calculations"}
	}

	for i := range e.PSM {

		var key string
		if len(e.PSM[i].ModifiedPeptide) > 0 {
			key = fmt.Sprintf("%s#%d", e.PSM[i].ModifiedPeptide, e.PSM[i].AssumedCharge)
		} else {
			key = fmt.Sprintf("%s#%d", e.PSM[i].Peptide, e.PSM[i].AssumedCharge)
		}

		// global intensity map for Peptides, getting the most intense
		_, okPep := intPepMap[e.PSM[i].Peptide]
		if okPep {
			if e.PSM[i].Intensity > intPepMap[e.PSM[i].Peptide] {
				intPepMap[e.PSM[i].Peptide] = e.PSM[i].Intensity
			}
		} else {
			intPepMap[e.PSM[i].Peptide] = e.PSM[i].Intensity
		}

		// global intensity map for Ions, getting the most intense
		_, okIon := intIonMap[key]
		if okIon {
			if e.PSM[i].Intensity > intIonMap[key] {
				intIonMap[key] = e.PSM[i].Intensity
			}
		} else {
			intIonMap[key] = e.PSM[i].Intensity
		}

	}

	// attribute intensities to peptide evidences
	for i := range e.Peptides {
		v, ok := intPepMap[e.Peptides[i].Sequence]
		if ok {
			e.Peptides[i].Intensity = v
		}
	}

	// attribute intensities to ion evidences
	for i := range e.Ions {

		var key string

		if len(e.Ions[i].ModifiedSequence) > 0 {
			key = fmt.Sprintf("%s#%d", e.Ions[i].ModifiedSequence, e.Ions[i].ChargeState)
		} else {
			key = fmt.Sprintf("%s#%d", e.Ions[i].Sequence, e.Ions[i].ChargeState)
		}

		v, ok := intIonMap[key]
		if ok {
			e.Ions[i].Intensity = v
		}

	}

	// attribute intensities to protein evidences
	for i := range e.Proteins {

		var totalInt []float64
		var uniqueInt []float64
		var razorInt []float64

		// for unique ions
		for k := range e.Proteins[i].UniquePeptideIons {
			v, ok := intIonMap[k]
			if ok {
				uniqueInt = append(uniqueInt, v)
			}
		}

		// for razor ions
		for k := range e.Proteins[i].URazorPeptideIons {
			v, ok := intIonMap[k]
			if ok {
				razorInt = append(razorInt, v)
			}
		}

		// for total ions
		for k := range e.Proteins[i].TotalPeptideIons {
			v, ok := intIonMap[k]
			if ok {
				totalInt = append(totalInt, v)
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

// // calculateIntensities calculates the protein intensity
// func calculateIntensities(e rep.Evidence) (rep.Evidence, *err.Error) {
//
// 	var intMap = make(map[string]float64)
// 	var intRefMap = make(map[string]float64)
// 	var intPepMap = make(map[string]float64)
// 	var intIonMap = make(map[string]float64)
//
// 	if len(e.PSM) < 1 || len(e.Ions) < 1 {
// 		return e, &err.Error{Type: err.CannotFindPSMData, Class: err.FATA, Argument: "cannot attribute intensity calculations"}
// 	}
//
// 	for i := range e.PSM {
//
// 		var key string
// 		if len(e.PSM[i].ModifiedPeptide) > 0 {
// 			key = fmt.Sprintf("%s#%d", e.PSM[i].ModifiedPeptide, e.PSM[i].AssumedCharge)
// 		} else {
// 			key = fmt.Sprintf("%s#%d", e.PSM[i].Peptide, e.PSM[i].AssumedCharge)
// 		}
//
// 		// global intensity map for spectra
// 		_, ok := intMap[key]
// 		if ok {
// 			if e.PSM[i].Intensity > intMap[key] {
// 				intMap[key] = e.PSM[i].Intensity
// 			}
// 		} else {
// 			intMap[key] = e.PSM[i].Intensity
// 		}
//
// 		// global intensity map for Peptides
// 		_, okPep := intPepMap[e.PSM[i].Peptide]
// 		if okPep {
// 			if e.PSM[i].Intensity > intPepMap[e.PSM[i].Peptide] {
// 				intPepMap[e.PSM[i].Peptide] = e.PSM[i].Intensity
// 			}
// 		} else {
// 			intPepMap[e.PSM[i].Peptide] = e.PSM[i].Intensity
// 		}
//
// 		// global intensity map for Ions
// 		_, okIon := intIonMap[key]
// 		if okIon {
// 			if e.PSM[i].Intensity > intIonMap[key] {
// 				intIonMap[key] = e.PSM[i].Intensity
// 			}
// 		} else {
// 			intIonMap[key] = e.PSM[i].Intensity
// 		}
//
// 	}
//
// 	// peptides get the higest intense from the matching sequences
// 	for i := range e.Peptides {
// 		v, ok := intMap[e.Peptides[i].Sequence]
// 		if ok {
// 			e.Peptides[i].Intensity = v
// 		}
// 	}
//
// 	// ions get the higest intense from the matching sequences
// 	for i := range e.Ions {
//
// 		var key string
// 		if len(e.Ions[i].ModifiedSequence) > 0 {
// 			key = fmt.Sprintf("%s#%d", e.Ions[i].ModifiedSequence, e.Ions[i].ChargeState)
// 		} else {
// 			key = fmt.Sprintf("%s#%d", e.Ions[i].Sequence, e.Ions[i].ChargeState)
// 		}
//
// 		v, ok := intMap[key]
// 		if ok {
// 			e.Ions[i].Intensity = v
// 			intRefMap[key] = v
// 		}
//
// 	}
//
// 	for i := range e.Proteins {
//
// 		var totalInt []float64
// 		var uniqueInt []float64
// 		var razorInt []float64
//
// 		// make a reference for razor peptides
// 		var uniqIons = make(map[string]uint8)
//
// 		for k := range e.Proteins[i].UniquePeptideIons {
// 			v, ok := intRefMap[k]
// 			if ok {
// 				uniqueInt = append(uniqueInt, v)
// 				razorInt = append(razorInt, v)
// 				uniqIons[k] = 0
// 			}
// 		}
//
// 		for k, j := range e.Proteins[i].TotalPeptideIons {
// 			v, ok := intRefMap[k]
// 			if ok {
// 				totalInt = append(totalInt, v)
// 				if j.IsRazor {
// 					_, ok := uniqIons[k]
// 					if !ok {
// 						razorInt = append(razorInt, v)
// 					}
// 				}
// 			}
// 		}
//
// 		sort.Float64s(uniqueInt)
// 		sort.Float64s(totalInt)
// 		sort.Float64s(razorInt)
//
// 		if len(uniqueInt) >= 3 {
// 			e.Proteins[i].UniqueIntensity = (uniqueInt[len(uniqueInt)-1] + uniqueInt[len(uniqueInt)-2] + uniqueInt[len(uniqueInt)-3])
// 		} else if len(uniqueInt) == 2 {
// 			e.Proteins[i].UniqueIntensity = (uniqueInt[len(uniqueInt)-1] + uniqueInt[len(uniqueInt)-2])
// 		} else if len(uniqueInt) == 1 {
// 			e.Proteins[i].UniqueIntensity = (uniqueInt[len(uniqueInt)-1])
// 		}
//
// 		if len(totalInt) >= 3 {
// 			e.Proteins[i].TotalIntensity = (totalInt[len(totalInt)-1] + totalInt[len(totalInt)-2] + totalInt[len(totalInt)-3])
// 		} else if len(totalInt) == 2 {
// 			e.Proteins[i].TotalIntensity = (totalInt[len(totalInt)-1] + totalInt[len(totalInt)-2])
// 		} else if len(totalInt) == 1 {
// 			e.Proteins[i].TotalIntensity = (totalInt[len(totalInt)-1])
// 		}
//
// 		if len(razorInt) >= 3 {
// 			e.Proteins[i].RazorIntensity = (razorInt[len(razorInt)-1] + razorInt[len(razorInt)-2] + razorInt[len(razorInt)-3])
// 		} else if len(razorInt) == 2 {
// 			e.Proteins[i].RazorIntensity = (razorInt[len(razorInt)-1] + razorInt[len(razorInt)-2])
// 		} else if len(razorInt) == 1 {
// 			e.Proteins[i].RazorIntensity = (razorInt[len(razorInt)-1])
// 		}
//
// 	}
//
// 	return e, nil
// }
