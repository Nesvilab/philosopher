package quan

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/cmsl/data/mz"
	"github.com/prvst/cmsl/data/mz/mzml"
	"github.com/prvst/cmsl/data/mz/mzxml"
	"github.com/prvst/cmsl/err"
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
func (p *Quantify) RunLabelFreeQuantification() *err.Error {

	var evi rep.Evidence
	e := evi.RestoreGranular()
	if e != nil {
		return e
	}

	if len(evi.Proteins) < 1 {
		logrus.Fatal("This result file does not contains report data")
	}

	// logrus.Info("Calculating Spectral Counts")
	// evi, e = CalculateSpectralCounts(evi)
	// if e != nil {
	// 	return e
	// }

	// logrus.Info("Calculating MS1 Intensities")
	// evi, e = traceApexPeaks(evi, p.Dir, p.Format, p.RTWin, p.PTWin, p.Tol)
	// if e != nil {
	// 	return e
	// }

	evi, e = peakIntensity(evi, p.Dir, p.Format, p.RTWin, p.PTWin, p.Tol)
	if e != nil {
		return e
	}

	evi, e = calculateIntensities(evi)
	if e != nil {
		return e
	}

	e = evi.SerializeGranular()
	if e != nil {
		return e
	}

	return nil
}

// RunTMTQuantification is the top function for label quantification
func (p *Quantify) RunTMTQuantification() error {

	var evi rep.Evidence
	e := evi.RestoreGranular()
	if e != nil {
		return e
	}

	// removed all calculated defined bvalues from before
	cleanPreviousData(p.Plex)

	var spectra []string
	var specGroup = make(map[string][]rep.PSMEvidence)
	var annotSpecGroup = make(map[string][]rep.PSMEvidence)
	var specMap = make(map[string]rep.PSMEvidence)

	// group ms2 spectra by source file
	for _, i := range evi.PSM {
		specName := strings.Split(i.Spectrum, ".")
		source := fmt.Sprintf("%s.%s", specName[0], p.Format)
		specGroup[source] = append(specGroup[source], i)
		// mapp all spectra for latter when we recover the quantifications
		specMap[i.Spectrum] = i
	}

	for i := range specGroup {
		spectra = append(spectra, i)
	}
	sort.Strings(spectra)

	logrus.Info("Calculating intensities and ion interference")
	for _, k := range spectra {

		logrus.Info("Processing ", k)

		ms1, ms2, e := getSpectra(p.Dir, p.Format, k)
		if e != nil {
			return e
		}

		annotSpecGroup[k], e = calculateIonPurity(p.Dir, p.Format, ms1, ms2, specGroup[k])
		if e != nil {
			return e
		}

		labels, err := labeledPeakIntensity(p.Dir, p.Format, strings.ToLower(p.Brand), p.Plex, p.Tol, specGroup[k], ms2)
		if err != nil {
			return err
		}

		annotSpecGroup[k], err = mapLabeledSpectra(labels, p.Purity, specGroup[k])
		if err != nil {
			return err
		}

	}

	logrus.Info("Mapping quantification values")

	// collect back the values from the annotate spectra to the original evidence structure
	for _, v := range annotSpecGroup {
		for _, i := range v {
			_, ok := specMap[i.Spectrum]
			if ok {
				specMap[i.Spectrum] = i
			}
		}
	}

	var annotEviList rep.PSMEvidenceList
	for _, i := range specMap {
		annotEviList = append(annotEviList, i)
	}

	evi.PSM = annotEviList

	e = evi.SerializeGranular()
	if e != nil {
		return e
	}

	return nil
}

func getSpectra(path, format string, spectra string) (map[string]mz.MS1, map[string]mz.MS2, error) {

	var err error
	var ms1 = make(map[string]mz.MS1)
	var ms2 = make(map[string]mz.MS2)

	ext := filepath.Ext(spectra)
	name := filepath.Base(spectra)
	clean := name[0 : len(name)-len(ext)]
	fullpath, _ := filepath.Abs(path)
	name = fmt.Sprintf("%s%s%s", fullpath, string(filepath.Separator), name)

	if strings.Contains(spectra, "mzML") {

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

	} else if strings.Contains(spectra, "mzXML") {

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

	return ms1, ms2, nil
}

// cleanPreviousData cleans previous label quantifications
func cleanPreviousData(plex string) *err.Error {

	var evi rep.Evidence
	e := evi.RestoreGranular()
	if e != nil {
		return e
	}

	for i := range evi.PSM {
		evi.PSM[i].Labels, e = tmt.New(plex)
		if e != nil {
			return e
		}
	}

	for i := range evi.Ions {
		evi.Ions[i].Labels, e = tmt.New(plex)
		if e != nil {
			return e
		}
	}

	for i := range evi.Proteins {
		evi.Proteins[i].TotalLabels, e = tmt.New(plex)
		if e != nil {
			return e
		}

		evi.Proteins[i].UniqueLabels, e = tmt.New(plex)
		if e != nil {
			return e
		}

		evi.Proteins[i].RazorLabels, e = tmt.New(plex)
		if e != nil {
			return e
		}

	}

	e = evi.SerializeGranular()
	if e != nil {
		return e
	}

	return nil
}

// func getSpectra(path, format string, evi rep.Evidence) (map[string]mz.MS1, map[string]mz.MS2, error) {
//
// 	var err error
// 	var ms1 = make(map[string]mz.MS1)
// 	var ms2 = make(map[string]mz.MS2)
// 	var mzs = make(map[string]int)
//
// 	// get specta file names from identifications
// 	for _, i := range evi.PSM {
// 		specName := strings.Split(i.Spectrum, ".")
// 		source := fmt.Sprintf("%s.%s", specName[0], format)
// 		mzs[source]++
// 	}
//
// 	for k := range mzs {
//
// 		ext := filepath.Ext(k)
// 		name := filepath.Base(k)
// 		clean := name[0 : len(name)-len(ext)]
// 		fullpath, _ := filepath.Abs(path)
// 		name = fmt.Sprintf("%s%s%s", fullpath, string(filepath.Separator), name)
//
// 		if strings.Contains(k, "mzML") {
//
// 			var ms mzml.IndexedMzML
// 			err = ms.Parse(name)
// 			if err != nil {
// 				return ms1, ms2, err
// 			}
//
// 			ms1Scans, ms1err := mz.GetMzMLSpectra(ms, clean)
// 			if ms1err != nil {
// 				return nil, nil, ms1err
// 			}
// 			ms1[clean] = ms1Scans
//
// 			ms2Scans, ms2err := mz.GetMzMLMS2Spectra(ms, clean)
// 			if ms2err != nil {
// 				return nil, nil, ms2err
// 			}
// 			ms2[clean] = ms2Scans
//
// 		} else if strings.Contains(k, "mzXML") {
//
// 			var ms mzxml.MzXML
// 			err = ms.Parse(name)
// 			if err != nil {
// 				return nil, nil, err
// 			}
//
// 			ms1Scans, err := mz.GetMzXMLSpectra(ms, clean)
// 			if err != nil {
// 				return nil, nil, err
// 			}
// 			ms1[clean] = ms1Scans
//
// 			ms2Scans, err := mz.GetMzXMLMS2Spectra(ms, clean)
// 			if err != nil {
// 				return nil, nil, err
// 			}
// 			ms2[clean] = ms2Scans
//
// 		} else {
// 			logrus.Fatal("Cannot open file: ", name)
// 		}
//
// 	}
//
// 	return ms1, ms2, nil
// }

// RunLabeledQuantification is the top function for label quantification
// func (p *Quantify) RunLabeledQuantification() error {
//
// 	var evi rep.Evidence
// 	e := evi.RestoreGranular()
// 	if e != nil {
// 		return e
// 	}
//
// 	// removed all calculated defined bvalues from before
// 	cleanPreviousData(p.Plex)
//
// 	logrus.Info("Reading spectra files")
// 	ms1, ms2, err := getSpectra(p.Dir, p.Format, evi)
// 	if err != nil {
// 		return err
// 	}
//
// 	logrus.Info("Calculating ion purity")
// 	evi, err = calculateIonPurity(p.Dir, p.Format, ms1, ms2, evi)
// 	if err != nil {
// 		return err
// 	}
//
// 	logrus.Info("Calculating label intensities")
// 	labels, err := labeledPeakIntensity(p.Dir, p.Format, strings.ToLower(p.Brand), p.Plex, p.Tol, evi, ms2)
// 	if err != nil {
// 		return err
// 	}
//
// 	logrus.Info("Mapping quantification values")
// 	evi, err = mapLabeledSpectra(labels, p.Purity, evi)
// 	if err != nil {
// 		return err
// 	}
//
// 	// optional
// 	// logrus.Info("Normalizing peptide channels")
// 	// evi, err = totalTop3LabelQuantification(evi)
// 	// if err != nil {
// 	// 	return err
// 	// }
//
// 	// evi, err = labelQuantificationOnTotalIons(evi)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	//
// 	// evi, err = labelQuantificationOnUniqueIons(evi)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	//
// 	// evi, err = labelQuantificationOnURazors(evi)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	//
// 	// if len(p.ChanNorm) > 0 {
// 	// 	logrus.Info("Applying normalization to control channel ", p.ChanNorm)
// 	// 	evi, err = ratioToControlChannel(evi, p.ChanNorm)
// 	// 	if err != nil {
// 	// 		return err
// 	// 	}
// 	// }
// 	//
// 	// if p.IntNorm == true {
// 	// 	logrus.Info("Applying normalization to intensity means")
// 	// 	evi, err = ratioToIntensityMean(evi)
// 	// 	if err != nil {
// 	// 		return err
// 	// 	}
// 	// }
//
// 	e = evi.SerializeGranular()
// 	if e != nil {
// 		fmt.Println("Err")
// 		return e
// 	}
//
// 	return nil
// }
