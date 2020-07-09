package qua

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"philosopher/lib/dat"
	"philosopher/lib/id"
	"philosopher/lib/iso"
	"philosopher/lib/met"
	"philosopher/lib/msg"
	"philosopher/lib/mzn"
	"philosopher/lib/rep"
	"philosopher/lib/sys"
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
// This function can be used for both pre and post filtering quantification
func RunLabelFreeQuantification(p met.Quantify) {

	var lfq = NewLFQ()

	// collect database information
	var db dat.Base
	db.Restore()

	psm, _ := id.ReadPepXMLInput(".", db.Prefix, sys.GetTemp(), false)

	//evi = peakIntensity(evi, p.Dir, p.Format, p.RTWin, p.PTWin, p.Tol, p.Isolated)
	psm = peakIntensity(psm, p.Dir, p.Format, p.RTWin, p.PTWin, p.Tol, p.Isolated)

	for _, i := range psm {
		lfq.Intensities[i.Spectrum] = i.Intensity
	}

	lfq.Serialize()

	// checks if the Evidence structure exists, if so, update it.
	if _, err := os.Stat(sys.EvPSMBin()); err == nil {

		var evi rep.Evidence
		evi.RestoreGranular()

		evi = PropagateIntensities(evi, lfq)
		evi.SerializeGranular()

	}

	return
}

// RunIsobaricLabelQuantification is the top function for label quantification
func RunIsobaricLabelQuantification(p met.Quantify, mods bool) met.Quantify {

	// collect database information
	var db dat.Base

	var labels = iso.NewIsoLabels()
	var psmMap = make(map[string]id.PeptideIdentification)
	var sourceMap = make(map[string][]id.PeptideIdentification)
	var sourceList []string

	if p.Brand == "" {
		msg.NoParametersFound(errors.New("You need to specify a brand type (tmt or itraq)"), "fatal")
	}

	var input string
	if len(p.Pex) > 0 {
		input = p.Pex
	} else {
		input = "."
	}

	psm, _ := id.ReadPepXMLInput(input, db.Prefix, sys.GetTemp(), false)

	if len(psm) < 1 {
		msg.NoPSMFound(errors.New("PSMs not found in data set"), "fatal")
	}

	// collect all used source file names
	for _, i := range psm {
		specName := strings.Split(i.Spectrum, ".")
		sourceMap[specName[0]] = append(sourceMap[specName[0]], i)
		psmMap[i.Spectrum] = i

		// remove the pep.xml file name from the spectrum name
		spectrumName := strings.Split(i.Spectrum, "#")

		// left-pad the spectrum scan
		paddedScan := fmt.Sprintf("%05d", i.Scan)

		// left-pad the spectrum index
		//paddedIndex := fmt.Sprintf("%05d", i.Index)

		var l iso.Labels
		if p.Brand == "tmt" {
			l = tmt.New(p.Plex)
		} else if p.Brand == "itraq" {
			l = trq.New(p.Plex)
		}

		l.Spectrum = i.Spectrum
		//l.Index = paddedIndex
		l.Scan = paddedScan
		l.RetentionTime = i.RetentionTime
		l.ChargeState = i.AssumedCharge

		labels.LabeledSpectra[spectrumName[0]] = l
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

	var purities []string
	for i := range sourceList {

		var mz mzn.MsData

		logrus.Info("Processing ", sourceList[i])
		fileName := fmt.Sprintf("%s%s%s.mzML", p.Dir, string(filepath.Separator), sourceList[i])

		mz.Read(fileName, false, false, false)

		for i := range mz.Spectra {
			mz.Spectra[i].Decode()
		}

		if p.Level == 3 {
			labels = prepareLabelStructureWithMS3(labels, p.Dir, p.Format, p.Brand, p.Plex, p.Tol, mz)

		} else {
			labels = prepareLabelStructureWithMS2(labels, p.Dir, p.Format, p.Brand, p.Plex, p.Tol, mz)
		}

		purities = append(purities, calculateIonPurity(labels, p.Dir, p.Format, mz)...)
	}

	for _, i := range purities {
		s := strings.Split(i, "#")
		v, ok := labels.LabeledSpectra[s[0]]
		if ok {
			v2 := v
			f, _ := strconv.ParseFloat(s[1], 64)
			v2.Purity = f
			labels.LabeledSpectra[s[0]] = v2
		}
	}

	labels = assignLabelNames(labels, p.LabelNames, p.Brand, p.Plex)

	labels.Serialize()

	// checks if the Evidence structure exists, if so, update it.
	if _, err := os.Stat(sys.EvPSMBin()); err == nil {

		var evi rep.Evidence
		evi.RestoreGranular()

		for i := range evi.PSM {
			s := strings.Split(evi.PSM[i].Spectrum, "#")
			v, ok := labels.LabeledSpectra[s[0]]
			if ok {
				evi.PSM[i].Labels = v
			}
		}

		// classification and filtering based on quality filters
		logrus.Info("Filtering spectra for label quantification")
		evi = Classification(evi, mods, p.BestPSM, p.RemoveLow, p.Purity, p.MinProb)

		// forces psms with no label to have 0 intensities
		evi = CorrectUnlabelledSpectra(evi)

		evi = RollUpPeptides(evi)

		evi = RollUpPeptideIons(evi)

		evi = RollUpProteins(evi)

		// normalize to the total protein levels
		logrus.Info("Calculating normalized protein levels")
		evi = NormToTotalProteins(evi)

		logrus.Info("Saving")

		// create EV PSM
		rep.SerializeEVPSM(&evi)

		// create EV Ion
		rep.SerializeEVIon(&evi)

		// create EV Peptides
		rep.SerializeEVPeptides(&evi)

		// create EV Ion
		rep.SerializeEVProteins(&evi)

	}

	return p
}

// RunBioQuantification is the top function for functional-based quantification
func RunBioQuantification(c met.Data) {

	// create clean reference db for clustering
	//clusterFasta := createCleanDataBaseReference(c.UUID, c.Temp)

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

	return
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
func assignLabelNames(labels iso.Tag, labelNames map[string]string, brand, plex string) iso.Tag {

	for k, v := range labels.LabeledSpectra {
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

		labels.LabeledSpectra[k] = v2
	}

	return labels
}

// Classification determines if PSMs should be used or not for isobaric tag quantification rollup
func Classification(evi rep.Evidence, mods, best bool, remove, purity, probability float64) rep.Evidence {

	var approvedPSMs = make(map[string]uint8)
	var bestMap = make(map[string]uint8)
	var psmLabelSumList PairList

	// 1st check: Purity the score and the Probability levels
	for _, i := range evi.PSM {

		if i.Probability >= probability && i.Labels.Purity >= purity {
			approvedPSMs[i.Spectrum] = 0
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

	// 3rd check: remove the lower 3%
	// Ignore all PSMs that fall under the lower 3% based on their summed TMT labels
	var toDelete = make(map[string]uint8)

	if remove != 0 {
		sort.Sort(psmLabelSumList)
		lowerFive := float64(len(psmLabelSumList)) * remove
		lowerFiveInt := int(uti.Round(lowerFive, 5, 0))

		for i := 0; i <= lowerFiveInt; i++ {
			toDelete[psmLabelSumList[i].Key] = 0
		}
	}

	logrus.Info("Removing ", len(toDelete), " PSMs from isobaric quantification")

	for _, i := range evi.PSM {
		_, ok := toDelete[i.Spectrum]
		if !ok {
			approvedPSMs[i.Spectrum] = 0
		}
	}

	for i := range evi.PSM {

		if mods == true {
			_, ok := evi.PSM[i].LocalizedPTMSites["PTMProphet_STY79.9663"]
			if ok {
				evi.PSM[i].Labels.HasPhospho = true
			}
		}

		_, ok1 := approvedPSMs[evi.PSM[i].Spectrum]
		_, ok2 := bestMap[evi.PSM[i].Spectrum]

		if best == true {
			if ok2 {
				evi.PSM[i].Labels.IsUsed = true
			}
		} else {
			if ok1 {
				evi.PSM[i].Labels.IsUsed = true
			}
		}
	}

	return evi
}
