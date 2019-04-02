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
	"github.com/prvst/philosopher/lib/psi/mzml"
	"github.com/prvst/philosopher/lib/rep"
	"github.com/prvst/philosopher/lib/tmt"
	"github.com/prvst/philosopher/lib/uti"
	"github.com/sirupsen/logrus"
)

type Pair struct {
	Key   string
	Value float64
}

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

	evi, e = peakIntensity(evi, p.Dir, p.Format, p.RTWin, p.PTWin, p.Tol, p.Isolated)
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

	// removed all calculated defined values from before
	evi, e = cleanPreviousData(evi, p.Plex)
	if e != nil {
		return p, e
	}

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

		var mz mzml.MsData

		logrus.Info("Processing ", sourceList[i])
		fileName := fmt.Sprintf("%s%s%s.mzML", p.Dir, string(filepath.Separator), sourceList[i])

		e = mz.Read(fileName, false, false, false)
		if e != nil {
			return p, e
		}

		// mz, e := getSpectra(p.Dir, p.Format, p.Level, sourceList[i])
		// if e != nil {
		// 	return p, e
		// }

		for i := range mz.Spectra {
			mz.Spectra[i].Decode()
		}

		mappedPurity, _ := calculateIonPurity(p.Dir, p.Format, mz, sourceMap[sourceList[i]])
		if e != nil {
			return p, e
		}

		//ms1 = raw.MS1{}

		var labels = make(map[string]tmt.Labels)
		if p.Level == 3 {
			var labE error
			labels, labE = prepareLabelStructureWithMS3(p.Dir, p.Format, p.Plex, p.Tol, mz)
			if labE != nil {
				return p, labE
			}
		} else {
			var labE error
			labels, labE = prepareLabelStructureWithMS2(p.Dir, p.Format, p.Plex, p.Tol, mz)
			if labE != nil {
				return p, labE
			}
		}

		//ms2 = raw.MS2{}
		//ms3 = raw.MS3{}

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

	// classification and filtering based on quality filters
	logrus.Info("Filtering spectra for label quantification")
	spectrumMap, phosphoSpectrumMap := classification(evi, mods, p.BestPSM, p.RemoveLow, p.Purity, p.MinProb)

	// assignment happens only for general PSMs
	evi = assignUsage(evi, spectrumMap)

	// forces psms with no label to have 0 intensities
	evi = correctUnlabelledSpectra(evi)

	evi = rollUpPeptides(evi, spectrumMap, phosphoSpectrumMap)

	evi = rollUpPeptideIons(evi, spectrumMap, phosphoSpectrumMap)

	evi = rollUpProteins(evi, spectrumMap, phosphoSpectrumMap)

	// normalize to the total protein levels
	logrus.Info("Calculating normalized protein levels")
	evi = NormToTotalProteins(evi)

	logrus.Info("Saving")

	// create EV PSM
	e = rep.SerializeEVPSM(&evi)
	if e != nil {
		return p, e
	}

	// create EV Ion
	e = rep.SerializeEVIon(&evi)
	if e != nil {
		return p, e
	}

	// create EV Peptides
	e = rep.SerializeEVPeptides(&evi)
	if e != nil {
		return p, e
	}
	// create EV Ion
	e = rep.SerializeEVProteins(&evi)
	if e != nil {
		return p, e
	}

	return p, nil
}

// cleanPreviousData cleans previous label quantifications
func cleanPreviousData(evi rep.Evidence, plex string) (rep.Evidence, *err.Error) {

	var e *err.Error

	for i := range evi.PSM {
		evi.PSM[i].Labels, e = tmt.New(plex)
		if e != nil {
			return evi, e
		}
	}

	for i := range evi.Ions {
		evi.Ions[i].Labels, e = tmt.New(plex)
		if e != nil {
			return evi, e
		}
	}

	for i := range evi.Proteins {
		evi.Proteins[i].TotalLabels, e = tmt.New(plex)
		if e != nil {
			return evi, e
		}

		evi.Proteins[i].UniqueLabels, e = tmt.New(plex)
		if e != nil {
			return evi, e
		}

		evi.Proteins[i].URazorLabels, e = tmt.New(plex)
		if e != nil {
			return evi, e
		}

	}

	return evi, nil
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

func classification(evi rep.Evidence, mods, best bool, remove, purity, probability float64) (map[string]tmt.Labels, map[string]tmt.Labels) {

	var spectrumMap = make(map[string]tmt.Labels)
	var phosphoSpectrumMap = make(map[string]tmt.Labels)

	var bestMap = make(map[string]uint8)

	var psmLabelSumList PairList

	// 1st check: Purity the score and the Probability levels
	for _, i := range evi.PSM {
		if i.Probability >= probability && i.Purity >= purity {

			spectrumMap[i.Spectrum] = i.Labels
			bestMap[i.Spectrum] = 0

			if mods == true {
				_, ok := i.LocalizedPTMSites["PTMProphet_STY79.9663"]
				if ok {
					phosphoSpectrumMap[i.Spectrum] = i.Labels
				}
			}

		}

		if remove != 0 {
			sum := i.Labels.Channel1.Intensity +
				i.Labels.Channel2.Intensity +
				i.Labels.Channel3.Intensity +
				i.Labels.Channel4.Intensity +
				i.Labels.Channel5.Intensity +
				i.Labels.Channel6.Intensity +
				i.Labels.Channel7.Intensity +
				i.Labels.Channel8.Intensity +
				i.Labels.Channel9.Intensity +
				i.Labels.Channel10.Intensity +
				i.Labels.Channel11.Intensity
			psmLabelSumList = append(psmLabelSumList, Pair{i.Spectrum, sum})
		}
	}

	// 2nd check: best PSM
	// collect all ion-related spectra from the each fraction/file
	// var bestMap = make(map[string]uint8)
	if best == true {
		var groupedPSMMap = make(map[string][]rep.PSMEvidence)
		for _, i := range evi.PSM {
			specName := strings.Split(i.Spectrum, ".")
			fqn := fmt.Sprintf("%s#%s", specName[0], i.IonForm)
			groupedPSMMap[fqn] = append(groupedPSMMap[fqn], i)
		}

		for _, v := range groupedPSMMap {
			if len(v) == 1 {
				bestMap[v[0].Spectrum] = 0
			} else {

				var bestPSM string
				var bestPSMInt float64
				for _, i := range v {
					tmtSum := i.Labels.Channel1.Intensity +
						i.Labels.Channel2.Intensity +
						i.Labels.Channel3.Intensity +
						i.Labels.Channel4.Intensity +
						i.Labels.Channel5.Intensity +
						i.Labels.Channel6.Intensity +
						i.Labels.Channel7.Intensity +
						i.Labels.Channel8.Intensity +
						i.Labels.Channel9.Intensity +
						i.Labels.Channel10.Intensity +
						i.Labels.Channel11.Intensity

					if tmtSum > bestPSMInt {
						bestPSM = i.Spectrum
						bestPSMInt = tmtSum
					}

				}

				bestMap[bestPSM] = 0

			}
		}
	}

	var toDelete = make(map[string]uint8)
	var toDeletePhospho = make(map[string]uint8)

	// 3rd check: remove the lower 3%
	// Ignore all PSMs that fall under the lower 3% based on their summed TMT labels
	if remove != 0 {
		sort.Sort(psmLabelSumList)
		lowerFive := float64(len(psmLabelSumList)) * remove
		lowerFiveInt := int(uti.Round(lowerFive, 5, 0))

		for i := 0; i <= lowerFiveInt; i++ {
			toDelete[psmLabelSumList[i].Key] = 0
			toDeletePhospho[psmLabelSumList[i].Key] = 0
		}
	}

	for k := range spectrumMap {
		_, ok := bestMap[k]
		if !ok {
			toDelete[k] = 0
		}
	}

	for i := range toDelete {
		delete(spectrumMap, i)
	}

	for k := range phosphoSpectrumMap {
		_, ok := bestMap[k]
		if !ok {
			toDeletePhospho[k] = 0
		}
	}

	logrus.Info("Removing ", len(toDelete), " PSMs from isobaric quantification")
	for i := range toDeletePhospho {
		delete(phosphoSpectrumMap, i)
	}

	return spectrumMap, phosphoSpectrumMap
}
