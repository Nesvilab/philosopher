package qua

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"path/filepath"
	"sort"
	"strings"

	"philosopher/lib/bio"
	"philosopher/lib/id"
	"philosopher/lib/msg"
	"philosopher/lib/sys"
	"philosopher/lib/uti"

	"philosopher/lib/mzn"
	"philosopher/lib/rep"

	"github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack"
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

func peakIntensity(psm id.PepIDList, dir, format string, rTWin, pTWin, tol float64, isIso bool) id.PepIDList {

	logrus.Info("Indexing PSM information")

	var sourceMap = make(map[string]uint8)
	var spectra = make(map[string][]string)
	var ppmPrecision = make(map[string]float64)
	var mzMap = make(map[string]float64)
	var minRT = make(map[string]float64)
	var maxRT = make(map[string]float64)
	var retentionTime = make(map[string]float64)
	var intensity = make(map[string]float64)

	var charges = make(map[string]int)

	// collect attributes from PSM
	for _, i := range psm {
		partName := strings.Split(i.Spectrum, ".")
		sourceMap[partName[0]] = 0
		spectra[partName[0]] = append(spectra[partName[0]], i.Spectrum)

		ppmPrecision[i.Spectrum] = tol / math.Pow(10, 6)
		mzMap[i.Spectrum] = ((i.PrecursorNeutralMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge))
		minRT[i.Spectrum] = (i.RetentionTime / 60) - rTWin
		maxRT[i.Spectrum] = (i.RetentionTime / 60) + rTWin
		retentionTime[i.Spectrum] = i.RetentionTime

		charges[i.Spectrum] = int(i.AssumedCharge)
	}

	// get a sorted list of spectrum names
	var sourceMapList []string
	for source := range sourceMap {
		sourceMapList = append(sourceMapList, source)
	}
	sort.Strings(sourceMapList)

	logrus.Info("Reading spectra and tracing peaks")
	for _, s := range sourceMapList {

		logrus.Info("Processing ", s)
		var mz mzn.MsData

		fileName := fmt.Sprintf("%s%s%s.mzML", dir, string(filepath.Separator), s)

		// load MS1, ignore MS2 and MS3
		mz.Read(fileName, false, false, true)

		for i := range mz.Spectra {
			if mz.Spectra[i].Level == "1" {
				mz.Spectra[i].Decode()
			} else if mz.Spectra[i].Level == "2" {
				spectrum := fmt.Sprintf("%s.%05s.%05s.%d", s, mz.Spectra[i].Scan, mz.Spectra[i].Scan, mz.Spectra[i].Precursor.ChargeState)
				_, ok := mzMap[spectrum]
				if ok {
					// update the MZ with the desired Precursor value from mzML
					if isIso == true {
						mzMap[spectrum] = mz.Spectra[i].Precursor.TargetIon
					} else {
						mzMap[spectrum] = mz.Spectra[i].Precursor.SelectedIon
					}
				}
			}
		}

		v, ok := spectra[s]
		if ok {
			for _, j := range v {

				measured, retrieved := xic(mz.Spectra, minRT[j], maxRT[j], ppmPrecision[j], mzMap[j])

				// if j == "20180209_03_TP_1A.03130.03130.2#interact.pep.xml" {
				// 	fmt.Println(measured)
				// }

				if retrieved == true {

					// create the list of mz differences for each peak
					var mzRatio []float64
					for k := 1; k <= 6; k++ {
						r := float64(k) * (float64(1) / float64(charges[j]))
						mzRatio = append(mzRatio, uti.ToFixed(r, 2))
					}

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

	for i := range psm {
		partName := strings.Split(psm[i].Spectrum, ".")
		_, ok := spectra[partName[0]]
		if ok {
			psm[i].Intensity = intensity[psm[i].Spectrum]
		}
	}

	return psm
}

// xic extract ion chomatograms
func xic(mz mzn.Spectra, minRT, maxRT, ppmPrecision, mzValue float64) (map[float64]float64, bool) {

	var list = make(map[float64]float64)

	for j := range mz {
		if mz[j].Level == "1" {

			if mz[j].ScanStartTime >= minRT && mz[j].ScanStartTime <= maxRT {

				lowi := sort.Search(len(mz[j].Mz.DecodedStream), func(i int) bool { return mz[j].Mz.DecodedStream[i] >= mzValue-ppmPrecision*mzValue })
				highi := sort.Search(len(mz[j].Mz.DecodedStream), func(i int) bool { return mz[j].Mz.DecodedStream[i] >= mzValue+ppmPrecision*mzValue })

				var maxI = 0.0

				// if mz[j].Index == "3106" {
				// 	spew.Dump(mzValue, mz[j].Intensity.DecodedStream[lowi:highi])
				// }

				for _, k := range mz[j].Intensity.DecodedStream[lowi:highi] {
					if k > maxI {
						maxI = k
					}
				}

				if maxI > 0 {
					list[mz[j].ScanStartTime] = maxI
				}

			}
		}
	}

	if len(list) >= 5 {
		return list, true
	}

	return list, false
}

// PropagateIntensities propagates the PSM intensities to the other layers
func PropagateIntensities(e rep.Evidence, lfq LFQ) rep.Evidence {

	logrus.Info("Assigning intensities to data layers")

	if len(e.PSM) < 1 || len(e.Ions) < 1 {
		msg.QuantifyingData(errors.New("The PSM list is enpty"), "fatal")
	}

	var peptideIntMap = make(map[string]float64)
	var ionIntMap = make(map[string]float64)

	for i := range e.PSM {

		v, ok1 := lfq.Intensities[e.PSM[i].Spectrum]
		if ok1 {
			e.PSM[i].Intensity = v
		}

		// peptide intensity : sum of all
		_, ok := peptideIntMap[e.PSM[i].Peptide]
		if ok {
			peptideIntMap[e.PSM[i].Peptide] += e.PSM[i].Intensity
		} else {
			peptideIntMap[e.PSM[i].Peptide] += e.PSM[i].Intensity
		}

		// ion intensity : most intense ion
		ionV, ok := ionIntMap[e.PSM[i].IonForm]
		if ok {
			if e.PSM[i].Intensity > ionV {
				ionIntMap[e.PSM[i].IonForm] = e.PSM[i].Intensity
			}
		} else {
			ionIntMap[e.PSM[i].IonForm] = e.PSM[i].Intensity
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

	return e
}

// Serialize saves to disk a msgpack version of the LFQ data structure
func (d *LFQ) Serialize() {

	b, e := msgpack.Marshal(&d)
	if e != nil {
		msg.MarshalFile(e, "fatal")
	}

	e = ioutil.WriteFile(sys.LFQBin(), b, sys.FilePermission())
	if e != nil {
		msg.SerializeFile(e, "fatal")
	}

	return
}

// Restore reads philosopher results files and restore the data sctructure
func (d *LFQ) Restore() {

	b, e := ioutil.ReadFile(sys.LFQBin())
	if e != nil {
		msg.MarshalFile(e, "warning")
	}

	e = msgpack.Unmarshal(b, &d)
	if e != nil {
		msg.SerializeFile(e, "warning")
	}

	return
}
