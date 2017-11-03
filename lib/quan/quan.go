package quan

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/cmsl/err"
	"github.com/prvst/philosopher/lib/meta"
	"github.com/prvst/philosopher/lib/rep"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/prvst/philosopher/lib/tmt"
	"github.com/prvst/philosopher/lib/xml"
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

	var psmMap = make(map[string]rep.PSMEvidence)
	var sourceMap = make(map[string][]rep.PSMEvidence)
	var sourceList []string

	logrus.Info("Restoring Data")

	var evi rep.Evidence
	e := evi.RestoreGranular()
	if e != nil {
		return e
	}

	// removed all calculated defined bvalues from before
	cleanPreviousData(p.Plex)

	// collect all used source file names
	for _, i := range evi.PSM {
		specName := strings.Split(i.Spectrum, ".")
		source := fmt.Sprintf("%s.%s", specName[0], p.Format)
		sourceMap[source] = append(sourceMap[source], i)

		psmMap[i.Spectrum] = i
	}

	for i := range sourceMap {
		sourceList = append(sourceList, i)
	}

	sort.Strings(sourceList)

	logrus.Info("Calculating intensities and ion interference")

	for i := range sourceList {

		logrus.Info("Reading ", sourceList[i])

		mz, e := getSpectra(p.Dir, p.Format, sourceList[i])
		if e != nil {
			return e
		}

		mappedPurity, e := calculateIonPurity(p.Dir, p.Format, mz, sourceMap[sourceList[i]])
		if e != nil {
			return e
		}

		labels, err := prepareLabelStructure(p.Dir, p.Format, p.Plex, p.Tol, mz)
		if err != nil {
			return err
		}

		mappedPSM, err := mapLabeledSpectra(labels, p.Purity, sourceMap[sourceList[i]])
		if err != nil {
			return err
		}

		for _, j := range mappedPurity {
			v, ok := psmMap[j.Spectrum]
			if ok {
				psm := v
				psm.Purity = j.Purity
				psmMap[j.Spectrum] = psm
			}
		}

		for _, j := range mappedPSM {
			v, ok := psmMap[j.Spectrum]
			if ok {
				psm := v
				psm.Labels = j.Labels
				psmMap[j.Spectrum] = psm
			}
		}
	}

	for i := range evi.PSM {
		v, ok := psmMap[evi.PSM[i].Spectrum]
		if ok {
			evi.PSM[i].Purity = v.Purity
			evi.PSM[i].Labels = v.Labels
		}
	}

	evi = rollUpPeptides(evi, p.Purity)

	evi = rollUpPeptideIons(evi, p.Purity)

	evi = rollUpProteins(evi, p.Purity)

	e = evi.SerializeGranular()
	if e != nil {
		return e
	}

	return nil
}

func getSpectra(path, format string, spectra string) (xml.Raw, error) {

	var err error
	var mzData xml.Raw

	//ext := filepath.Ext(spectra)
	//clean := name[0 : len(name)-len(ext)]
	name := filepath.Base(spectra)
	fullpath, _ := filepath.Abs(path)
	name = fmt.Sprintf("%s%s%s", fullpath, string(filepath.Separator), name)

	if strings.Contains(spectra, "mzML") {

		err = mzData.Read(name, "mzML")
		if err != nil {
			return mzData, err
		}

	} else if strings.Contains(spectra, "mzXML") {

		fmt.Println("NULL")

	} else {
		logrus.Fatal("Cannot open file: ", name)
	}

	return mzData, nil
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

		evi.Proteins[i].URazorLabels, e = tmt.New(plex)
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
