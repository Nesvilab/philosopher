package qua

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/raw"
	"github.com/prvst/philosopher/lib/rep"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/prvst/philosopher/lib/tmt"
	"github.com/sirupsen/logrus"
)

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

// RunLabelFreeQuantification is the top function for label free quantification
func RunLabelFreeQuantification(p met.Quantify) *err.Error {

	var evi rep.Evidence
	e := evi.RestoreGranular()
	if e != nil {
		return e
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
func RunTMTQuantification(p met.Quantify, mods bool) (met.Quantify, error) {

	var psmMap = make(map[string]rep.PSMEvidence)
	var sourceMap = make(map[string][]rep.PSMEvidence)
	var sourceList []string

	logrus.Info("Restoring data")

	var evi rep.Evidence
	e := evi.RestoreGranular()
	if e != nil {
		return p, e
	}

	// removed all calculated defined bvalues from before
	cleanPreviousData(p.Plex)

	// collect all used source file names
	for _, i := range evi.PSM {
		specName := strings.Split(i.Spectrum, ".")
		sourceMap[specName[0]] = append(sourceMap[specName[0]], i)
		psmMap[i.Spectrum] = i
	}

	for i := range sourceMap {
		sourceList = append(sourceList, i)
	}

	if len(sourceMap) > 1 {
		sort.Strings(sourceList)
	}

	// read the annotation file
	p.LabelNames = make(map[string]string)
	if len(p.Annot) > 0 {
		p.LabelNames, e = getLabelNames(p.Annot)
		if e != nil {
			return p, e
		}
	}

	logrus.Info("Calculating intensities and ion interference")

	for i := range sourceList {

		var ms1 raw.MS1
		var ms2 raw.MS2

		logrus.Info("Reading ", sourceList[i])

		ms1, ms2, e = getSpectra(p.Dir, p.Format, sourceList[i])
		if e != nil {
			return p, e
		}

		mappedPurity, _ := calculateIonPurity(p.Dir, p.Format, ms1, ms2, sourceMap[sourceList[i]])
		if e != nil {
			return p, e
		}

		ms1 = raw.MS1{}

		labels, err := prepareLabelStructure(p.Dir, p.Format, p.Plex, p.Tol, ms2)
		if err != nil {
			return p, err
		}

		ms2 = raw.MS2{}

		labels = assignLabelNames(labels, p.LabelNames)

		mappedPSM, err := mapLabeledSpectra(labels, p.Purity, sourceMap[sourceList[i]])
		if err != nil {
			return p, err
		}

		for _, j := range mappedPurity {
			v, ok := psmMap[j.Spectrum]
			if ok {
				psm := v
				psm.Purity = j.Purity
				psmMap[j.Spectrum] = psm
			}
		}
		mappedPurity = nil

		for _, j := range mappedPSM {
			v, ok := psmMap[j.Spectrum]
			if ok {
				psm := v
				psm.Labels = j.Labels
				psmMap[j.Spectrum] = psm
			}
		}
		mappedPSM = nil

	}

	for i := range evi.PSM {
		v, ok := psmMap[evi.PSM[i].Spectrum]
		if ok {
			evi.PSM[i].Purity = v.Purity
			evi.PSM[i].Labels = v.Labels
		}
	}
	psmMap = nil

	var spectrumMap = make(map[string]tmt.Labels)
	for _, i := range evi.PSM {
		if i.Purity >= p.Purity {
			spectrumMap[i.Spectrum] = i.Labels
		}
	}

	// forces psms with no label to have 0 intensities
	evi = correctUnlabelledSpectra(evi)

	evi = rollUpPeptides(evi, spectrumMap)

	evi = rollUpPeptideIons(evi, spectrumMap)

	evi = rollUpProteins(evi, spectrumMap)

	// normalize to the total protein levels
	logrus.Info("Calculating normalized protein levels")
	evi = NormToTotalProteins(evi)

	logrus.Info("Saving")
	e = evi.SerializeGranular()
	if e != nil {
		return p, e
	}

	return p, nil
}

func getSpectra(dir, format string, k string) (raw.MS1, raw.MS2, *err.Error) {

	var ms1 raw.MS1
	var ms2 raw.MS2

	// get the clean name, remove the extension
	var extension = filepath.Ext(filepath.Base(k))
	var name = k[0 : len(k)-len(extension)]
	input := fmt.Sprintf("%s%s%s.bin", sys.MetaDir(), string(filepath.Separator), name)

	// get all MS1 spectra
	if _, e := os.Stat(input); e == nil {

		spec, e := raw.Restore(k)
		if e != nil {
			return ms1, ms2, &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: "error restoring indexed mz"}
		}

		ms1 = raw.GetMS1(spec)
		ms2 = raw.GetMS2(spec)

	} else {

		spec, rer := raw.RestoreFromFile(dir, k, format)
		if rer != nil {
			return ms1, ms2, &err.Error{Type: err.CannotParseXML, Class: err.FATA, Argument: "cant read mz file"}
		}

		ms1 = raw.GetMS1(spec)
		ms2 = raw.GetMS2(spec)
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

// addCustomNames adds to the label structures user-defined names to be used on the TMT labels
func getLabelNames(annot string) (map[string]string, *err.Error) {

	var labels = make(map[string]string)

	file, e := os.Open(annot)
	if e != nil {
		return labels, &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: e.Error()}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		names := strings.Split(scanner.Text(), " ")
		labels[names[0]] = names[1]
	}

	if e = scanner.Err(); e != nil {
		return labels, &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: e.Error()}
	}

	return labels, nil
}

// checks for custom names and assign the normal channel or the custom name to the CustomName
func assignLabelNames(labels map[string]tmt.Labels, labelNames map[string]string) map[string]tmt.Labels {

	for _, i := range labels {

		switch chnl := i.Channel1.Name; chnl {
		case "126":
			i.Channel1.CustomName = labelNames["126"]

			if len(i.Channel1.CustomName) < 1 {
				i.Channel1.CustomName = "126"
			}

		case "127N":
			i.Channel2.CustomName = labelNames["127N"]
			if len(i.Channel2.CustomName) < 1 {
				i.Channel2.CustomName = "127N"
			}

		case "127C":
			i.Channel3.CustomName = labelNames["127C"]
			if len(i.Channel3.CustomName) < 1 {
				i.Channel3.CustomName = "127C"
			}

		case "128N":
			i.Channel4.CustomName = labelNames["128N"]
			if len(i.Channel4.CustomName) < 1 {
				i.Channel4.CustomName = "128N"
			}

		case "128C":
			i.Channel5.CustomName = labelNames["128C"]
			if len(i.Channel5.CustomName) < 1 {
				i.Channel5.CustomName = "128C"
			}

		case "129N":
			i.Channel6.CustomName = labelNames["129N"]
			if len(i.Channel6.CustomName) < 1 {
				i.Channel6.CustomName = "129N"
			}

		case "129C":
			i.Channel7.CustomName = labelNames["129C"]
			if len(i.Channel7.CustomName) < 1 {
				i.Channel7.CustomName = "129C"
			}

		case "130N":
			i.Channel8.CustomName = labelNames["130N"]
			if len(i.Channel8.CustomName) < 1 {
				i.Channel8.CustomName = "130N"
			}

		case "130C":
			i.Channel9.CustomName = labelNames["130C"]
			if len(i.Channel9.CustomName) < 1 {
				i.Channel9.CustomName = "130C"
			}

		case "131N":
			i.Channel10.CustomName = labelNames["131N"]
			if len(i.Channel10.CustomName) < 1 {
				i.Channel10.CustomName = "131N"
			}

		case "131C":
			i.Channel11.CustomName = labelNames["131C"]
			if len(i.Channel11.CustomName) < 1 {
				i.Channel11.CustomName = "131C"
			}

		default:

		}
	}

	return labels
}
