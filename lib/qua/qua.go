package qua

import (
	"errors"
	"fmt"
	"github.com/Nesvilab/philosopher/lib/ibt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Nesvilab/philosopher/lib/id"
	"github.com/Nesvilab/philosopher/lib/iso"
	"github.com/Nesvilab/philosopher/lib/met"
	"github.com/Nesvilab/philosopher/lib/msg"
	"github.com/Nesvilab/philosopher/lib/mzn"
	"github.com/Nesvilab/philosopher/lib/rep"
	"github.com/Nesvilab/philosopher/lib/scl"
	"github.com/Nesvilab/philosopher/lib/tmt"
	"github.com/Nesvilab/philosopher/lib/trq"
	"github.com/Nesvilab/philosopher/lib/uti"
	"github.com/Nesvilab/philosopher/lib/xta"
	"github.com/Nesvilab/philosopher/lib/xta2"

	"github.com/sirupsen/logrus"
)

// Pair ...
type Pair struct {
	Key   id.SpectrumType
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

	if len(evi.PSM) < 1 || len(evi.Ions) < 1 {
		msg.QuantifyingData(errors.New("the PSM list is empty."), "warning")
		os.Exit(0)
	}

	evi = peakIntensity(evi, p.Dir, p.Format, p.RTWin, p.PTWin, p.Tol, p.Isolated, p.Raw, p.Faims)

	evi = calculateIntensities(evi)

	evi.SerializeGranular()

}

// RunIsobaricLabelQuantification is the top function for label quantification
func RunIsobaricLabelQuantification(p met.Quantify, mods bool) met.Quantify {

	var psmMap = make(map[id.SpectrumType]rep.PSMEvidence)
	var sourceMap = make(map[string][]rep.PSMEvidence)
	var sourceList []string

	if p.Brand == "" {
		msg.NoParametersFound(errors.New("you need to specify a brand type (tmt or itraq)"), "error")
	}

	var evi rep.Evidence
	evi.RestoreGranular()

	if len(evi.PSM) < 1 || len(evi.Ions) < 1 {
		msg.QuantifyingData(errors.New("the PSM list is empty."), "warning")
		os.Exit(0)
	}

	// removed all calculated defined values from before
	evi = cleanPreviousData(evi, p.Brand, p.Plex)

	// collect all used source file names
	for _, i := range evi.PSM {
		specName := strings.Split(i.Spectrum, ".")
		sourceMap[specName[0]] = append(sourceMap[specName[0]], i)
		psmMap[i.SpectrumFileName()] = i
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
			//fileName = fmt.Sprintf("%s%s%s.raw", p.Dir, string(filepath.Separator), sourceList[i])
			//stream := rawfilereader.Run(fileName, "")
			//mz.ReadRaw(fileName, stream)

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
			v, ok := psmMap[j.SpectrumFileName()]
			if ok {
				psm := v
				psm.Purity = j.Purity
				psmMap[j.SpectrumFileName()] = psm
			}
		}

		for _, j := range mappedPSM {
			v, ok := psmMap[j.SpectrumFileName()]
			if ok {
				psm := v
				psm.Labels = j.Labels
				psmMap[j.SpectrumFileName()] = psm
			}
		}

	}

	for i := range evi.PSM {
		v, ok := psmMap[evi.PSM[i].SpectrumFileName()]
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
	//logrus.Info("Calculating normalized protein levels")
	//evi = NormToTotalProteins(evi)

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
			evi.PSM[i].Labels = &iso.Labels{}
			*evi.PSM[i].Labels = tmt.New(plex)
		} else if brand == "itraq" {
			evi.PSM[i].Labels = &iso.Labels{}
			*evi.PSM[i].Labels = trq.New(plex)
		} else if brand == "sclip" {
			evi.PSM[i].Labels = &iso.Labels{}
			*evi.PSM[i].Labels = scl.New(plex)
		} else if brand == "ibt" {
			evi.PSM[i].Labels = &iso.Labels{}
			*evi.PSM[i].Labels = ibt.New(plex)
		} else if brand == "xtag" {
			evi.PSM[i].Labels = &iso.Labels{}
			*evi.PSM[i].Labels = xta.New(plex)
		} else if brand == "xtag2" {
			evi.PSM[i].Labels = &iso.Labels{}
			*evi.PSM[i].Labels = xta2.New(plex)
		}
	}

	for i := range evi.Ions {
		if brand == "tmt" {
			evi.Ions[i].Labels = &iso.Labels{}
			*evi.Ions[i].Labels = tmt.New(plex)
		} else if brand == "itraq" {
			evi.Ions[i].Labels = &iso.Labels{}
			*evi.Ions[i].Labels = trq.New(plex)
		} else if brand == "sclip" {
			evi.Ions[i].Labels = &iso.Labels{}
			*evi.Ions[i].Labels = scl.New(plex)
		} else if brand == "ibt" {
			evi.Ions[i].Labels = &iso.Labels{}
			*evi.Ions[i].Labels = ibt.New(plex)
		} else if brand == "xtag" {
			evi.Ions[i].Labels = &iso.Labels{}
			*evi.Ions[i].Labels = xta.New(plex)
		} else if brand == "xtag2" {
			evi.Ions[i].Labels = &iso.Labels{}
			*evi.Ions[i].Labels = xta2.New(plex)
		}
	}

	for i := range evi.Proteins {
		if brand == "tmt" {
			evi.Proteins[i].TotalLabels = &iso.Labels{}
			evi.Proteins[i].UniqueLabels = &iso.Labels{}
			evi.Proteins[i].URazorLabels = &iso.Labels{}
			*evi.Proteins[i].TotalLabels = tmt.New(plex)
			*evi.Proteins[i].UniqueLabels = tmt.New(plex)
			*evi.Proteins[i].URazorLabels = tmt.New(plex)
		} else if brand == "itraq" {
			evi.Proteins[i].TotalLabels = &iso.Labels{}
			evi.Proteins[i].UniqueLabels = &iso.Labels{}
			evi.Proteins[i].URazorLabels = &iso.Labels{}
			*evi.Proteins[i].TotalLabels = trq.New(plex)
			*evi.Proteins[i].UniqueLabels = trq.New(plex)
			*evi.Proteins[i].URazorLabels = trq.New(plex)
		} else if brand == "sclip" {
			evi.Proteins[i].TotalLabels = &iso.Labels{}
			evi.Proteins[i].UniqueLabels = &iso.Labels{}
			evi.Proteins[i].URazorLabels = &iso.Labels{}
			*evi.Proteins[i].TotalLabels = scl.New(plex)
			*evi.Proteins[i].UniqueLabels = scl.New(plex)
			*evi.Proteins[i].URazorLabels = scl.New(plex)
		} else if brand == "ibt" {
			evi.Proteins[i].TotalLabels = &iso.Labels{}
			evi.Proteins[i].UniqueLabels = &iso.Labels{}
			evi.Proteins[i].URazorLabels = &iso.Labels{}
			*evi.Proteins[i].TotalLabels = ibt.New(plex)
			*evi.Proteins[i].UniqueLabels = ibt.New(plex)
			*evi.Proteins[i].URazorLabels = ibt.New(plex)
		} else if brand == "xtag" {
			evi.Proteins[i].TotalLabels = &iso.Labels{}
			evi.Proteins[i].UniqueLabels = &iso.Labels{}
			evi.Proteins[i].URazorLabels = &iso.Labels{}
			*evi.Proteins[i].TotalLabels = xta.New(plex)
			*evi.Proteins[i].UniqueLabels = xta.New(plex)
			*evi.Proteins[i].URazorLabels = xta.New(plex)
		} else if brand == "xtag2" {
			evi.Proteins[i].TotalLabels = &iso.Labels{}
			evi.Proteins[i].UniqueLabels = &iso.Labels{}
			evi.Proteins[i].URazorLabels = &iso.Labels{}
			*evi.Proteins[i].TotalLabels = xta2.New(plex)
			*evi.Proteins[i].UniqueLabels = xta2.New(plex)
			*evi.Proteins[i].URazorLabels = xta2.New(plex)
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

			if len(labelNames["134C"]) < 1 {
				v2.Channel17.CustomName = "134C"
			} else {
				v2.Channel17.CustomName = labelNames["134C"]
			}

			if len(labelNames["135N"]) < 1 {
				v2.Channel18.CustomName = "135N"
			} else {
				v2.Channel18.CustomName = labelNames["135N"]
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

		} else if brand == "sclip" {

			if len(labelNames["sCLIP1"]) < 1 {
				v2.Channel1.CustomName = "sCLIP1"
			} else {
				v2.Channel1.CustomName = labelNames["sCLIP1"]
			}

			if len(labelNames["sCLIP2"]) < 1 {
				v2.Channel2.CustomName = "sCLIP2"
			} else {
				v2.Channel2.CustomName = labelNames["sCLIP2"]
			}

			if len(labelNames["sCLIP3"]) < 1 {
				v2.Channel3.CustomName = "sCLIP3"
			} else {
				v2.Channel3.CustomName = labelNames["sCLIP3"]
			}

			if len(labelNames["sCLIP4"]) < 1 {
				v2.Channel4.CustomName = "sCLIP4"
			} else {
				v2.Channel4.CustomName = labelNames["sCLIP4"]
			}

			if len(labelNames["sCLIP5"]) < 1 {
				v2.Channel5.CustomName = "sCLIP5"
			} else {
				v2.Channel5.CustomName = labelNames["sCLIP5"]
			}

			if len(labelNames["sCLIP6"]) < 1 {
				v2.Channel6.CustomName = "sCLIP6"
			} else {
				v2.Channel6.CustomName = labelNames["sCLIP6"]
			}

		} else if brand == "ibt" {

			if len(labelNames["114"]) < 1 {
				v2.Channel1.CustomName = "114"
			} else {
				v2.Channel1.CustomName = labelNames["114"]
			}

			if len(labelNames["115N"]) < 1 {
				v2.Channel2.CustomName = "115N"
			} else {
				v2.Channel2.CustomName = labelNames["115N"]
			}

			if len(labelNames["115C"]) < 1 {
				v2.Channel3.CustomName = "115C"
			} else {
				v2.Channel3.CustomName = labelNames["115C"]
			}

			if len(labelNames["116N"]) < 1 {
				v2.Channel4.CustomName = "116N"
			} else {
				v2.Channel4.CustomName = labelNames["116N"]
			}

			if len(labelNames["116C"]) < 1 {
				v2.Channel5.CustomName = "116C"
			} else {
				v2.Channel5.CustomName = labelNames["116C"]
			}

			if len(labelNames["117N"]) < 1 {
				v2.Channel6.CustomName = "117N"
			} else {
				v2.Channel6.CustomName = labelNames["117N"]
			}

			if len(labelNames["117C"]) < 1 {
				v2.Channel7.CustomName = "117C"
			} else {
				v2.Channel7.CustomName = labelNames["117C"]
			}

			if len(labelNames["118N"]) < 1 {
				v2.Channel8.CustomName = "118N"
			} else {
				v2.Channel8.CustomName = labelNames["118N"]
			}

			if len(labelNames["118C"]) < 1 {
				v2.Channel9.CustomName = "118C"
			} else {
				v2.Channel9.CustomName = labelNames["118C"]
			}

			if len(labelNames["119N"]) < 1 {
				v2.Channel10.CustomName = "119N"
			} else {
				v2.Channel10.CustomName = labelNames["119N"]
			}

			if len(labelNames["119C"]) < 1 {
				v2.Channel11.CustomName = "119C"
			} else {
				v2.Channel11.CustomName = labelNames["119C"]
			}

			if len(labelNames["120N"]) < 1 {
				v2.Channel12.CustomName = "120N"
			} else {
				v2.Channel12.CustomName = labelNames["120N"]
			}

			if len(labelNames["120C"]) < 1 {
				v2.Channel13.CustomName = "120C"
			} else {
				v2.Channel13.CustomName = labelNames["120C"]
			}

			if len(labelNames["121N"]) < 1 {
				v2.Channel14.CustomName = "121N"
			} else {
				v2.Channel14.CustomName = labelNames["121N"]
			}

			if len(labelNames["121C"]) < 1 {
				v2.Channel15.CustomName = "121C"
			} else {
				v2.Channel15.CustomName = labelNames["121C"]
			}

			if len(labelNames["122"]) < 1 {
				v2.Channel16.CustomName = "122"
			} else {
				v2.Channel16.CustomName = labelNames["122"]
			}

		} else if brand == "xtag" {

			if len(labelNames["xTag1"]) < 1 {
				v2.Channel1.CustomName = "xTag1"
			} else {
				v2.Channel1.CustomName = labelNames["xTag1"]
			}

			if len(labelNames["xTag2"]) < 1 {
				v2.Channel2.CustomName = "xTag2"
			} else {
				v2.Channel2.CustomName = labelNames["xTag2"]
			}

			if len(labelNames["xTag3"]) < 1 {
				v2.Channel3.CustomName = "xTag3"
			} else {
				v2.Channel3.CustomName = labelNames["xTag3"]
			}

			if len(labelNames["xTag4"]) < 1 {
				v2.Channel4.CustomName = "xTag4"
			} else {
				v2.Channel4.CustomName = labelNames["xTag4"]
			}

			if len(labelNames["xTag5"]) < 1 {
				v2.Channel5.CustomName = "xTag5"
			} else {
				v2.Channel5.CustomName = labelNames["xTag5"]
			}

			if len(labelNames["xTag6"]) < 1 {
				v2.Channel6.CustomName = "xTag6"
			} else {
				v2.Channel6.CustomName = labelNames["xTag6"]
			}

			if len(labelNames["xTag7"]) < 1 {
				v2.Channel7.CustomName = "xTag7"
			} else {
				v2.Channel7.CustomName = labelNames["xTag7"]
			}

			if len(labelNames["xTag8"]) < 1 {
				v2.Channel8.CustomName = "xTag8"
			} else {
				v2.Channel8.CustomName = labelNames["xTag8"]
			}

			if len(labelNames["xTag9"]) < 1 {
				v2.Channel9.CustomName = "xTag9"
			} else {
				v2.Channel9.CustomName = labelNames["xTag9"]
			}

			if len(labelNames["xTag10"]) < 1 {
				v2.Channel10.CustomName = "xTag10"
			} else {
				v2.Channel10.CustomName = labelNames["xTag10"]
			}

			if len(labelNames["xTag11"]) < 1 {
				v2.Channel11.CustomName = "xTag11"
			} else {
				v2.Channel11.CustomName = labelNames["xTag11"]
			}

			if len(labelNames["xTag12"]) < 1 {
				v2.Channel12.CustomName = "xTag12"
			} else {
				v2.Channel12.CustomName = labelNames["xTag12"]
			}

			if len(labelNames["xTag13"]) < 1 {
				v2.Channel13.CustomName = "xTag13"
			} else {
				v2.Channel13.CustomName = labelNames["xTag13"]
			}

			if len(labelNames["xTag14"]) < 1 {
				v2.Channel14.CustomName = "xTag14"
			} else {
				v2.Channel14.CustomName = labelNames["xTag14"]
			}

			if len(labelNames["xTag15"]) < 1 {
				v2.Channel15.CustomName = "xTag15"
			} else {
				v2.Channel15.CustomName = labelNames["xTag15"]
			}

			if len(labelNames["xTag16"]) < 1 {
				v2.Channel16.CustomName = "xTag16"
			} else {
				v2.Channel16.CustomName = labelNames["xTag16"]
			}

			if len(labelNames["xTag17"]) < 1 {
				v2.Channel17.CustomName = "xTag17"
			} else {
				v2.Channel17.CustomName = labelNames["xTag17"]
			}

			if len(labelNames["xTag18"]) < 1 {
				v2.Channel18.CustomName = "xTag18"
			} else {
				v2.Channel18.CustomName = labelNames["xTag18"]
			}

			if len(labelNames["xTag19"]) < 1 {
				v2.Channel19.CustomName = "xTag19"
			} else {
				v2.Channel19.CustomName = labelNames["xTag19"]
			}

			if len(labelNames["xTag20"]) < 1 {
				v2.Channel20.CustomName = "xTag20"
			} else {
				v2.Channel20.CustomName = labelNames["xTag20"]
			}

			if len(labelNames["xTag21"]) < 1 {
				v2.Channel21.CustomName = "xTag21"
			} else {
				v2.Channel21.CustomName = labelNames["xTag21"]
			}

			if len(labelNames["xTag22"]) < 1 {
				v2.Channel22.CustomName = "xTag22"
			} else {
				v2.Channel22.CustomName = labelNames["xTag22"]
			}

			if len(labelNames["xTag23"]) < 1 {
				v2.Channel23.CustomName = "xTag23"
			} else {
				v2.Channel23.CustomName = labelNames["xTag23"]
			}

			if len(labelNames["xTag24"]) < 1 {
				v2.Channel24.CustomName = "xTag24"
			} else {
				v2.Channel24.CustomName = labelNames["xTag24"]
			}

			if len(labelNames["xTag25"]) < 1 {
				v2.Channel25.CustomName = "xTag25"
			} else {
				v2.Channel25.CustomName = labelNames["xTag25"]
			}

			if len(labelNames["xTag26"]) < 1 {
				v2.Channel26.CustomName = "xTag26"
			} else {
				v2.Channel26.CustomName = labelNames["xTag26"]
			}

			if len(labelNames["xTag27"]) < 1 {
				v2.Channel27.CustomName = "xTag27"
			} else {
				v2.Channel27.CustomName = labelNames["xTag27"]
			}

			if len(labelNames["xTag28"]) < 1 {
				v2.Channel28.CustomName = "xTag28"
			} else {
				v2.Channel28.CustomName = labelNames["xTag28"]
			}

			if len(labelNames["xTag29"]) < 1 {
				v2.Channel29.CustomName = "xTag29"
			} else {
				v2.Channel29.CustomName = labelNames["xTag29"]
			}

			if len(labelNames["xTag30"]) < 1 {
				v2.Channel30.CustomName = "xTag30"
			} else {
				v2.Channel30.CustomName = labelNames["xTag30"]
			}

			if len(labelNames["xTag31"]) < 1 {
				v2.Channel31.CustomName = "xTag31"
			} else {
				v2.Channel31.CustomName = labelNames["xTag31"]
			}

			if len(labelNames["xTag32"]) < 1 {
				v2.Channel32.CustomName = "xTag32"
			} else {
				v2.Channel32.CustomName = labelNames["xTag32"]
			}

		} else if brand == "xtag2" {

			if len(labelNames["114"]) < 1 {
				v2.Channel1.CustomName = "114"
			} else {
				v2.Channel1.CustomName = labelNames["114"]
			}
			if len(labelNames["115a"]) < 1 {
				v2.Channel2.CustomName = "115a"
			} else {
				v2.Channel2.CustomName = labelNames["115a"]
			}
			if len(labelNames["115b"]) < 1 {
				v2.Channel3.CustomName = "115b"
			} else {
				v2.Channel3.CustomName = labelNames["115b"]
			}
			if len(labelNames["115c"]) < 1 {
				v2.Channel4.CustomName = "115c"
			} else {
				v2.Channel4.CustomName = labelNames["115c"]
			}
			if len(labelNames["116a"]) < 1 {
				v2.Channel5.CustomName = "116a"
			} else {
				v2.Channel5.CustomName = labelNames["116a"]
			}
			if len(labelNames["116b"]) < 1 {
				v2.Channel6.CustomName = "116b"
			} else {
				v2.Channel6.CustomName = labelNames["116b"]
			}
			if len(labelNames["116c"]) < 1 {
				v2.Channel7.CustomName = "116c"
			} else {
				v2.Channel7.CustomName = labelNames["116c"]
			}
			if len(labelNames["116d"]) < 1 {
				v2.Channel8.CustomName = "116d"
			} else {
				v2.Channel8.CustomName = labelNames["116d"]
			}
			if len(labelNames["116e"]) < 1 {
				v2.Channel9.CustomName = "116e"
			} else {
				v2.Channel9.CustomName = labelNames["116e"]
			}
			if len(labelNames["117a"]) < 1 {
				v2.Channel10.CustomName = "117a"
			} else {
				v2.Channel10.CustomName = labelNames["117a"]
			}
			if len(labelNames["117b"]) < 1 {
				v2.Channel11.CustomName = "117b"
			} else {
				v2.Channel11.CustomName = labelNames["117b"]
			}
			if len(labelNames["117c"]) < 1 {
				v2.Channel12.CustomName = "117c"
			} else {
				v2.Channel12.CustomName = labelNames["117c"]
			}
			if len(labelNames["117d"]) < 1 {
				v2.Channel13.CustomName = "117d"
			} else {
				v2.Channel13.CustomName = labelNames["117d"]
			}
			if len(labelNames["117e"]) < 1 {
				v2.Channel14.CustomName = "117e"
			} else {
				v2.Channel14.CustomName = labelNames["117e"]
			}
			if len(labelNames["117f"]) < 1 {
				v2.Channel15.CustomName = "117f"
			} else {
				v2.Channel15.CustomName = labelNames["117f"]
			}
			if len(labelNames["118a"]) < 1 {
				v2.Channel16.CustomName = "118a"
			} else {
				v2.Channel16.CustomName = labelNames["118a"]
			}
			if len(labelNames["118b"]) < 1 {
				v2.Channel17.CustomName = "118b"
			} else {
				v2.Channel17.CustomName = labelNames["118b"]
			}
			if len(labelNames["118c"]) < 1 {
				v2.Channel18.CustomName = "118c"
			} else {
				v2.Channel18.CustomName = labelNames["118c"]
			}
			if len(labelNames["118d"]) < 1 {
				v2.Channel19.CustomName = "118d"
			} else {
				v2.Channel19.CustomName = labelNames["118d"]
			}
			if len(labelNames["118e"]) < 1 {
				v2.Channel20.CustomName = "118e"
			} else {
				v2.Channel20.CustomName = labelNames["118e"]
			}
			if len(labelNames["118f"]) < 1 {
				v2.Channel21.CustomName = "118f"
			} else {
				v2.Channel21.CustomName = labelNames["118f"]
			}
			if len(labelNames["118g"]) < 1 {
				v2.Channel22.CustomName = "118g"
			} else {
				v2.Channel22.CustomName = labelNames["118g"]
			}
			if len(labelNames["119a"]) < 1 {
				v2.Channel23.CustomName = "119a"
			} else {
				v2.Channel23.CustomName = labelNames["119a"]
			}
			if len(labelNames["119b"]) < 1 {
				v2.Channel24.CustomName = "119b"
			} else {
				v2.Channel24.CustomName = labelNames["119b"]
			}
			if len(labelNames["119c"]) < 1 {
				v2.Channel25.CustomName = "119c"
			} else {
				v2.Channel25.CustomName = labelNames["119c"]
			}
			if len(labelNames["119d"]) < 1 {
				v2.Channel26.CustomName = "119d"
			} else {
				v2.Channel26.CustomName = labelNames["119d"]
			}
			if len(labelNames["119e"]) < 1 {
				v2.Channel27.CustomName = "119e"
			} else {
				v2.Channel27.CustomName = labelNames["119e"]
			}
			if len(labelNames["119f"]) < 1 {
				v2.Channel28.CustomName = "119f"
			} else {
				v2.Channel28.CustomName = labelNames["119f"]
			}
			if len(labelNames["119g"]) < 1 {
				v2.Channel29.CustomName = "119g"
			} else {
				v2.Channel29.CustomName = labelNames["119g"]
			}

		}

		labels[k] = v2
	}

	return labels
}

func classification(evi rep.Evidence, mods, best bool, remove, purity, probability float64) (map[id.SpectrumType]iso.Labels, map[id.SpectrumType]iso.Labels) {

	var spectrumMap = make(map[id.SpectrumType]iso.Labels)
	var phosphoSpectrumMap = make(map[id.SpectrumType]iso.Labels)
	var bestMap = make(map[id.SpectrumType]uint8)
	var psmLabelSumList PairList
	var quantCheckUp bool

	// 1st check: Purity the score and the Probability levels
	for _, i := range evi.PSM {
		if i.Probability >= probability && i.Purity >= purity && !(i.IsDecoy) {

			spectrumMap[i.SpectrumFileName()] = *i.Labels
			bestMap[i.SpectrumFileName()] = 0

			if mods && i.PTM != nil {
				_, ok1 := i.PTM.LocalizedPTMSites["PTMProphet_STY79.9663"]
				_, ok2 := i.PTM.LocalizedPTMSites["PTMProphet_STY79.96633"]
				_, ok3 := i.PTM.LocalizedPTMSites["PTMProphet_STY79.966331"]
				if ok1 || ok2 || ok3 {
					phosphoSpectrumMap[i.SpectrumFileName()] = *i.Labels
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
				i.Labels.Channel16.Intensity +
				i.Labels.Channel17.Intensity +
				i.Labels.Channel18.Intensity +
				i.Labels.Channel19.Intensity +
				i.Labels.Channel20.Intensity +
				i.Labels.Channel21.Intensity +
				i.Labels.Channel22.Intensity +
				i.Labels.Channel23.Intensity +
				i.Labels.Channel24.Intensity +
				i.Labels.Channel25.Intensity +
				i.Labels.Channel26.Intensity +
				i.Labels.Channel27.Intensity +
				i.Labels.Channel28.Intensity +
				i.Labels.Channel29.Intensity +
				i.Labels.Channel30.Intensity +
				i.Labels.Channel31.Intensity +
				i.Labels.Channel32.Intensity
			psmLabelSumList = append(psmLabelSumList, Pair{i.SpectrumFileName(), sum})

			if sum > 0 {
				quantCheckUp = true
			}
		}
	}

	if remove != 0 && !quantCheckUp {
		remove = 0
		msg.Custom(errors.New("there are no non-zero intensities. Set 'removelow' to 0"), "warning")
	}

	// 2nd check: best PSM
	// collect all ion-related spectra from the each fraction/file
	// var bestMap = make(map[string]uint8)
	if best {
		var groupedPSMMap = make(map[string][]rep.PSMEvidence)
		for _, i := range evi.PSM {
			specName := strings.Split(i.Spectrum, ".")
			fqn := fmt.Sprintf("%s#%s", specName[0], i.IonForm().Str())
			groupedPSMMap[fqn] = append(groupedPSMMap[fqn], i)
		}

		for _, v := range groupedPSMMap {
			if len(v) == 1 {
				bestMap[v[0].SpectrumFileName()] = 0
			} else {

				var bestPSM id.SpectrumType
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
						i.Labels.Channel16.Intensity +
						i.Labels.Channel17.Intensity +
						i.Labels.Channel18.Intensity +
						i.Labels.Channel19.Intensity +
						i.Labels.Channel20.Intensity +
						i.Labels.Channel21.Intensity +
						i.Labels.Channel22.Intensity +
						i.Labels.Channel23.Intensity +
						i.Labels.Channel24.Intensity +
						i.Labels.Channel25.Intensity +
						i.Labels.Channel26.Intensity +
						i.Labels.Channel27.Intensity +
						i.Labels.Channel28.Intensity +
						i.Labels.Channel29.Intensity +
						i.Labels.Channel30.Intensity +
						i.Labels.Channel31.Intensity +
						i.Labels.Channel32.Intensity

					if tmtSum > bestPSMInt {
						bestPSM = i.SpectrumFileName()
						bestPSMInt = tmtSum
					}

				}

				bestMap[bestPSM] = 0

			}
		}
	}

	var toDelete = make(map[id.SpectrumType]uint8)
	var toDeletePhospho = make(map[id.SpectrumType]uint8)

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

	logrus.Info("Removing ", len(evi.PSM)-len(spectrumMap), " PSMs from isobaric quantification")
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
				if v1.Mz.DecodedStream[k] >= (v2.Precursor.TargetIon-v2.Precursor.IsolationWindowLowerOffset) && v1.Mz.DecodedStream[k] <= (v2.Precursor.TargetIon+v2.Precursor.IsolationWindowUpperOffset) {
					ions[v1.Mz.DecodedStream[k]] = v1.Intensity.DecodedStream[k]
					isolationWindowSummedInt += v1.Intensity.DecodedStream[k]
				}
			}

			// create the list of mz differences for each peak
			var mzRatio []float64
			for k := 1; k <= 6; k++ {
				r := float64(k) * 1.0033548378 / float64(v2.Precursor.ChargeState)
				mzRatio = append(mzRatio, uti.Round(r, 5, 2))
			}

			isotopesInt := v2.Precursor.TargetIonIntensity

			for k, v := range ions {
				for _, m := range mzRatio {
					if math.Abs(v2.Precursor.TargetIon-k) <= (m+0.025) && math.Abs(v2.Precursor.TargetIon-k) >= (m-0.025) {
						if v != v2.Precursor.TargetIonIntensity {
							isotopesInt += v
						}
						break
					}
				}
			}

			if isolationWindowSummedInt < 0 {
				msg.Custom(errors.New("summed intensity within isolation window is negative, should not happen"), "warning")
			}
			if isotopesInt < 0 {
				msg.Custom(errors.New("isotopes summed intensity is negative, should not happen"), "warning")
			}

			if isolationWindowSummedInt <= 0 || isotopesInt <= 0 {
				evi[i].Purity = 0
			} else {
				evi[i].Purity = uti.Round((isotopesInt / isolationWindowSummedInt), 5, 2)
			}

		}
	}

	return evi
}
