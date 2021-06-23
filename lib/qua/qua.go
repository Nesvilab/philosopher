package qua

import (
	"errors"
	"fmt"
	"math"
	"path/filepath"
	"sort"
	"strings"

	"philosopher/lib/ext/rawfilereader"
	"philosopher/lib/iso"
	"philosopher/lib/met"
	"philosopher/lib/msg"
	"philosopher/lib/mzn"
	"philosopher/lib/rep"
	"philosopher/lib/tmt"
	"philosopher/lib/trq"
	"philosopher/lib/uti"

	"github.com/sirupsen/logrus"
)

// Pair ...
type Pair struct {
	Key   string
	Value float64
}

// PairList ...
type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// RunLabelFreeQuantification is the top function for label free quantification
func RunLabelFreeQuantification(p met.Quantify) {

	// This parameter is hardcoded now because of the changes in the latest msconvert version 3.20.
	p.Isolated = true

	var evi rep.Evidence
	evi.RestoreGranular()

	evi = peakIntensity(evi, p.Dir, p.Format, p.RTWin, p.PTWin, p.Tol, p.Isolated, p.Raw)

	evi = calculateIntensities(evi)

	evi.SerializeGranular()

}

// RunIsobaricLabelQuantification is the top function for label quantification
func RunIsobaricLabelQuantification(p met.Quantify, mods bool) met.Quantify {

	var psmMap = make(map[string]rep.PSMEvidence)
	var sourceMap = make(map[string][]rep.PSMEvidence)
	var sourceList []string

	if p.Brand == "" {
		msg.NoParametersFound(errors.New("you need to specify a brand type (tmt or itraq)"), "fatal")
	}

	var evi rep.Evidence
	evi.RestoreGranular()

	// removed all calculated defined values from before
	evi = cleanPreviousData(evi, p.Brand, p.Plex)

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
		p.LabelNames = uti.GetLabelNames(p.Annot)
	}

	logrus.Info("Calculating intensities and ion interference")

	for i := range sourceList {

		var mz mzn.MsData
		var fileName string

		logrus.Info("Processing ", sourceList[i])

		if p.Raw {

			fileName = fmt.Sprintf("%s%s%s.raw", p.Dir, string(filepath.Separator), sourceList[i])
			stream := rawfilereader.Run(fileName, "")
			mz.ReadRaw(fileName, stream)

		} else {

			fileName = fmt.Sprintf("%s%s%s.mzML", p.Dir, string(filepath.Separator), sourceList[i])
			mz.Read(fileName)

			for i := range mz.Spectra {
				mz.Spectra[i].Decode()
			}
		}

		mappedPurity := calculateIonPurity(p.Dir, p.Format, mz, sourceMap[sourceList[i]])

		var labels map[string]iso.Labels
		if p.Level == 3 {
			labels = prepareLabelStructureWithMS3(p.Dir, p.Format, p.Brand, p.Plex, p.Tol, mz)

		} else {
			labels = prepareLabelStructureWithMS2(p.Dir, p.Format, p.Brand, p.Plex, p.Tol, mz)
		}

		labels = assignLabelNames(labels, p.LabelNames, p.Brand, p.Plex)

		mappedPSM := mapLabeledSpectra(labels, p.Purity, sourceMap[sourceList[i]])

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
	//psmMap = nil

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
	//rep.SerializeEVPSM(&evi)
	rep.SerializePSM(&evi.PSM)

	// create Ion
	rep.SerializeIon(&evi.Ions)

	// create Peptides
	rep.SerializePeptides(&evi.Peptides)

	// create Ion
	rep.SerializeProteins(&evi.Proteins)

	return p
}

// RunBioQuantification is the top function for functional-based quantification
func RunBioQuantification(c met.Data) {

	// run cdhit, create cluster file
	logrus.Info("Clustering")
	clusterFile, clusterFasta := execute(c.BioQuant.Level)

	// parse the cluster file
	logrus.Info("Parsing clusters")
	clusters := parseClusterFile(clusterFile, clusterFasta)

	// maps all proteins from the db against the clusters
	logrus.Info("Mapping proteins to clusters")
	mappedClust := mapProtXML2Clusters(clusters)

	logrus.Info("Retrieving Proteome data")
	//mappedClust = retrieveInfoFromUniProtDB(mappedClust)

	// mapping to functional annotation and save to disk
	savetoDisk(mappedClust, c.Temp, c.BioQuant.UID)

}

// cleanPreviousData cleans previous label quantifications
func cleanPreviousData(evi rep.Evidence, brand, plex string) rep.Evidence {

	for i := range evi.PSM {
		if brand == "tmt" {
			evi.PSM[i].Labels = tmt.New(plex)
		} else if brand == "itraq" {
			evi.PSM[i].Labels = trq.New(plex)
		}
	}

	for i := range evi.Ions {
		if brand == "tmt" {
			evi.Ions[i].Labels = tmt.New(plex)
		} else if brand == "itraq" {
			evi.Ions[i].Labels = trq.New(plex)
		}
	}

	for i := range evi.Proteins {
		if brand == "tmt" {
			evi.Proteins[i].TotalLabels = tmt.New(plex)
			evi.Proteins[i].UniqueLabels = tmt.New(plex)
			evi.Proteins[i].URazorLabels = tmt.New(plex)
		} else if brand == "itraq" {
			evi.Proteins[i].TotalLabels = trq.New(plex)
			evi.Proteins[i].UniqueLabels = trq.New(plex)
			evi.Proteins[i].URazorLabels = trq.New(plex)
		}
	}

	return evi
}

// checks for custom names and assign the normal channel or the custom name to the CustomName
func assignLabelNames(labels map[string]iso.Labels, labelNames map[string]string, brand, plex string) map[string]iso.Labels {

	for k, v := range labels {
		v2 := v

		if brand == "tmt" {

			if len(labelNames["126"]) < 1 {
				v2.Channel1.CustomName = "126"
			} else {
				v2.Channel1.CustomName = labelNames["126"]
			}

			if len(labelNames["127N"]) < 1 {
				v2.Channel2.CustomName = "127N"
			} else {
				v2.Channel2.CustomName = labelNames["127N"]
			}

			if len(labelNames["127C"]) < 1 {
				v2.Channel3.CustomName = "127C"
			} else {
				v2.Channel3.CustomName = labelNames["127C"]
			}

			if len(labelNames["128N"]) < 1 {
				v2.Channel4.CustomName = "128N"
			} else {
				v2.Channel4.CustomName = labelNames["128N"]
			}

			if len(labelNames["128C"]) < 1 {
				v2.Channel5.CustomName = "128C"
			} else {
				v2.Channel5.CustomName = labelNames["128C"]
			}

			if len(labelNames["129N"]) < 1 {
				v2.Channel6.CustomName = "129N"
			} else {
				v2.Channel6.CustomName = labelNames["129N"]
			}

			if len(labelNames["129C"]) < 1 {
				v2.Channel7.CustomName = "129C"
			} else {
				v2.Channel7.CustomName = labelNames["129C"]
			}

			if len(labelNames["130N"]) < 1 {
				v2.Channel8.CustomName = "130N"
			} else {
				v2.Channel8.CustomName = labelNames["130N"]
			}

			if len(labelNames["130C"]) < 1 {
				v2.Channel9.CustomName = "130C"
			} else {
				v2.Channel9.CustomName = labelNames["130C"]
			}

			if len(labelNames["131N"]) < 1 {
				v2.Channel10.CustomName = "131N"
			} else {
				v2.Channel10.CustomName = labelNames["131N"]
			}

			if len(labelNames["131C"]) < 1 {
				v2.Channel11.CustomName = "131C"
			} else {
				v2.Channel11.CustomName = labelNames["131C"]
			}

			if len(labelNames["132N"]) < 1 {
				v2.Channel12.CustomName = "132N"
			} else {
				v2.Channel12.CustomName = labelNames["132N"]
			}

			if len(labelNames["132C"]) < 1 {
				v2.Channel13.CustomName = "132C"
			} else {
				v2.Channel13.CustomName = labelNames["132C"]
			}

			if len(labelNames["133N"]) < 1 {
				v2.Channel14.CustomName = "133N"
			} else {
				v2.Channel14.CustomName = labelNames["133N"]
			}

			if len(labelNames["133C"]) < 1 {
				v2.Channel15.CustomName = "133C"
			} else {
				v2.Channel15.CustomName = labelNames["133C"]
			}

			if len(labelNames["134N"]) < 1 {
				v2.Channel16.CustomName = "134N"
			} else {
				v2.Channel16.CustomName = labelNames["134N"]
			}

		} else if brand == "itraq" && plex == "4" {

			if len(labelNames["114"]) < 1 {
				v2.Channel1.CustomName = "114"
			} else {
				v2.Channel1.CustomName = labelNames["114"]
			}

			if len(labelNames["115"]) < 1 {
				v2.Channel2.CustomName = "115"
			} else {
				v2.Channel2.CustomName = labelNames["115"]
			}

			if len(labelNames["116"]) < 1 {
				v2.Channel3.CustomName = "116"
			} else {
				v2.Channel3.CustomName = labelNames["116"]
			}

			if len(labelNames["117"]) < 1 {
				v2.Channel4.CustomName = "117"
			} else {
				v2.Channel4.CustomName = labelNames["117"]
			}

		} else if brand == "itraq" && plex == "8" {
			if len(labelNames["113"]) < 1 {
				v2.Channel1.CustomName = "113"
			} else {
				v2.Channel1.CustomName = labelNames["113"]
			}

			if len(labelNames["114"]) < 1 {
				v2.Channel2.CustomName = "114"
			} else {
				v2.Channel2.CustomName = labelNames["114"]
			}

			if len(labelNames["115"]) < 1 {
				v2.Channel3.CustomName = "115"
			} else {
				v2.Channel3.CustomName = labelNames["115"]
			}

			if len(labelNames["116"]) < 1 {
				v2.Channel4.CustomName = "116"
			} else {
				v2.Channel4.CustomName = labelNames["116"]
			}

			if len(labelNames["117"]) < 1 {
				v2.Channel5.CustomName = "117"
			} else {
				v2.Channel5.CustomName = labelNames["117"]
			}

			if len(labelNames["118"]) < 1 {
				v2.Channel6.CustomName = "118"
			} else {
				v2.Channel6.CustomName = labelNames["118"]
			}

			if len(labelNames["119"]) < 1 {
				v2.Channel7.CustomName = "119"
			} else {
				v2.Channel7.CustomName = labelNames["119"]
			}

			if len(labelNames["121"]) < 1 {
				v2.Channel8.CustomName = "121"
			} else {
				v2.Channel8.CustomName = labelNames["121"]
			}

		}

		labels[k] = v2
	}

	return labels
}

func classification(evi rep.Evidence, mods, best bool, remove, purity, probability float64) (map[string]iso.Labels, map[string]iso.Labels) {

	var spectrumMap = make(map[string]iso.Labels)
	var phosphoSpectrumMap = make(map[string]iso.Labels)
	var bestMap = make(map[string]uint8)
	var psmLabelSumList PairList
	var quantCheckUp bool

	// 1st check: Purity the score and the Probability levels
	for _, i := range evi.PSM {
		if i.Probability >= probability && i.Purity >= purity {

			spectrumMap[i.Spectrum] = i.Labels
			bestMap[i.Spectrum] = 0

			if mods {
				_, ok1 := i.LocalizedPTMSites["PTMProphet_STY79.9663"]
				_, ok2 := i.LocalizedPTMSites["PTMProphet_STY79.96633"]
				_, ok3 := i.LocalizedPTMSites["PTMProphet_STY79.966331"]
				if ok1 || ok2 || ok3 {
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
				i.Labels.Channel11.Intensity +
				i.Labels.Channel12.Intensity +
				i.Labels.Channel13.Intensity +
				i.Labels.Channel14.Intensity +
				i.Labels.Channel15.Intensity +
				i.Labels.Channel16.Intensity
			psmLabelSumList = append(psmLabelSumList, Pair{i.Spectrum, sum})

			if sum > 0 {
				quantCheckUp = true
			}
		}
	}

	if remove != 0 && !quantCheckUp {
		msg.NoParametersFound(errors.New("no reporter ions found. Check your MS level, or update msconvert"), "fatal")
	}

	// 2nd check: best PSM
	// collect all ion-related spectra from the each fraction/file
	// var bestMap = make(map[string]uint8)
	if best {
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
						i.Labels.Channel11.Intensity +
						i.Labels.Channel12.Intensity +
						i.Labels.Channel13.Intensity +
						i.Labels.Channel14.Intensity +
						i.Labels.Channel15.Intensity +
						i.Labels.Channel16.Intensity

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

// calculateIonPurity verifies how much interference there is on the precursor scans for each fragment
func calculateIonPurity(d, f string, mz mzn.MsData, evi []rep.PSMEvidence) []rep.PSMEvidence {

	// index MS1 and MS2 spectra in a dictionary
	var indexedMS1 = make(map[string]mzn.Spectrum)
	var indexedMS2 = make(map[string]mzn.Spectrum)

	var MS1Peaks = make(map[string][]float64)
	var MS1Int = make(map[string][]float64)

	for i := range mz.Spectra {

		if mz.Spectra[i].Level == "1" {

			// left-pad the spectrum index
			paddedIndex := fmt.Sprintf("%05s", mz.Spectra[i].Index)

			// left-pad the spectrum scan
			paddedScan := fmt.Sprintf("%05s", mz.Spectra[i].Scan)

			mz.Spectra[i].Index = paddedIndex
			mz.Spectra[i].Scan = paddedScan

			indexedMS1[paddedScan] = mz.Spectra[i]

			MS1Peaks[paddedScan] = mz.Spectra[i].Mz.DecodedStream
			MS1Int[paddedScan] = mz.Spectra[i].Intensity.DecodedStream

		} else if mz.Spectra[i].Level == "2" {

			if mz.Spectra[i].Precursor.IsolationWindowLowerOffset == 0 && mz.Spectra[i].Precursor.IsolationWindowUpperOffset == 0 {
				mz.Spectra[i].Precursor.IsolationWindowLowerOffset = mzDeltaWindow
				mz.Spectra[i].Precursor.IsolationWindowUpperOffset = mzDeltaWindow
			}

			// left-pad the spectrum index
			paddedIndex := fmt.Sprintf("%05s", mz.Spectra[i].Index)

			// left-pad the spectrum scan
			paddedScan := fmt.Sprintf("%05s", mz.Spectra[i].Scan)

			// left-pad the precursor spectrum index
			paddedPI := fmt.Sprintf("%05s", mz.Spectra[i].Precursor.ParentIndex)

			// left-pad the precursor spectrum scan
			paddedPS := fmt.Sprintf("%05s", mz.Spectra[i].Precursor.ParentScan)

			mz.Spectra[i].Index = paddedIndex
			mz.Spectra[i].Scan = paddedScan
			mz.Spectra[i].Precursor.ParentIndex = paddedPI
			mz.Spectra[i].Precursor.ParentScan = paddedPS

			stream := MS1Peaks[paddedPS]

			for j := range stream {
				if stream[j] >= (mz.Spectra[i].Precursor.TargetIon-mz.Spectra[i].Precursor.IsolationWindowLowerOffset) && stream[j] <= (mz.Spectra[i].Precursor.TargetIon+mz.Spectra[i].Precursor.IsolationWindowUpperOffset) {
					if MS1Int[mz.Spectra[i].Precursor.ParentScan][j] > mz.Spectra[i].Precursor.TargetIonIntensity {
						mz.Spectra[i].Precursor.TargetIonIntensity = MS1Int[mz.Spectra[i].Precursor.ParentScan][j]
					}
				}
			}

			indexedMS2[paddedScan] = mz.Spectra[i]
		}
	}

	for i := range evi {

		// get spectrum index
		split := strings.Split(evi[i].Spectrum, ".")

		v2, ok := indexedMS2[split[1]]
		if ok {

			v1 := indexedMS1[v2.Precursor.ParentScan]

			var ions = make(map[float64]float64)
			var isolationWindowSummedInt float64

			for k := range v1.Mz.DecodedStream {
				if v1.Mz.DecodedStream[k] >= (v2.Precursor.TargetIon-v2.Precursor.IsolationWindowUpperOffset) && v1.Mz.DecodedStream[k] <= (v2.Precursor.TargetIon+v2.Precursor.IsolationWindowUpperOffset) {
					ions[v1.Mz.DecodedStream[k]] = v1.Intensity.DecodedStream[k]
					isolationWindowSummedInt += v1.Intensity.DecodedStream[k]
				}
			}

			// create the list of mz differences for each peak
			var mzRatio []float64
			for k := 1; k <= 6; k++ {
				r := float64(k) * (float64(1) / float64(v2.Precursor.ChargeState))
				mzRatio = append(mzRatio, uti.ToFixed(r, 2))
			}

			var isotopePackage = make(map[float64]float64)
			isotopePackage[v2.Precursor.TargetIon] = v2.Precursor.TargetIonIntensity
			isotopesInt := v2.Precursor.TargetIonIntensity

			for k, v := range ions {
				for _, m := range mzRatio {
					if math.Abs(v2.Precursor.TargetIon-k) <= (m+0.025) && math.Abs(v2.Precursor.TargetIon-k) >= (m-0.025) {
						if v != v2.Precursor.TargetIonIntensity {
							isotopePackage[k] = v
							isotopesInt += v
						}
						break
					}
				}
			}

			if isotopesInt == 0 {
				evi[i].Purity = 0
			} else {
				evi[i].Purity = uti.Round((isotopesInt / isolationWindowSummedInt), 5, 2)
			}

		}
	}

	return evi
}
