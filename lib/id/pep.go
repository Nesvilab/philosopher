package id

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"philosopher/lib/uti"

	"philosopher/lib/msg"

	"philosopher/lib/mod"
	"philosopher/lib/spc"
	"philosopher/lib/sys"

	"github.com/sirupsen/logrus"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// PepXML data
type PepXML struct {
	FileName              string
	SpectraFile           string
	SearchEngine          string
	DecoyTag              string
	Database              string
	Prophet               string
	SearchParameters      []spc.Parameter
	Models                []spc.DistributionPoint
	Modifications         mod.Modifications
	PeptideIdentification PepIDList
}

type PepXML4Serialiazation struct {
	FileName              string
	SpectraFile           string
	SearchEngine          string
	DecoyTag              string
	Database              string
	Prophet               string
	SearchParameters      []spc.Parameter
	Models                []spc.DistributionPoint
	Modifications         mod.Modifications
	PeptideIdentification PepIDListPtrs
}

// PeptideIdentificationOld struct
// type PeptideIdentificationOld struct {
// 	Spectrum                             string
// 	SpectrumFile                         string
// 	Peptide                              string
// 	Protein                              string
// 	ModifiedPeptide                      string
// 	CompensationVoltage                  string
// 	PrevAA                               string
// 	NextAA                               string
// 	LocalizationRange                    string
// 	MSFragerLocalization                 string
// 	MSFraggerLocalizationScoreWithPTM    string
// 	MSFraggerLocalizationScoreWithoutPTM string
// 	Scan                                 int
// 	NumberofMissedCleavages              int
// 	IsoMassD                             int
// 	AssumedCharge                        uint8
// 	HitRank                              uint8
// 	MissedCleavages                      uint8
// 	NumberTolTerm                        uint8
// 	NumberOfEnzymaticTermini             uint8
// 	IsRejected                           uint8
// 	NumberTotalProteins                  uint16
// 	TotalNumberIons                      uint16
// 	NumberMatchedIons                    uint16
// 	Index                                uint32
// 	UncalibratedPrecursorNeutralMass     float64
// 	PrecursorNeutralMass                 float64
// 	PrecursorExpMass                     float64
// 	RetentionTime                        float64
// 	CalcNeutralPepMass                   float64
// 	Massdiff                             float64
// 	Probability                          float64
// 	Expectation                          float64
// 	Xcorr                                float64
// 	DeltaCN                              float64
// 	DeltaCNStar                          float64
// 	SPScore                              float64
// 	SPRank                               float64
// 	Hyperscore                           float64
// 	Nextscore                            float64
// 	DiscriminantValue                    float64
// 	Intensity                            float64
// 	IonMobility                          float64
// 	AlternativeProteins                  map[string]int
// 	LocalizedPTMSites                    map[string]int
// 	LocalizedPTMMassDiff                 map[string]string
// 	Modifications                        mod.Modifications
// }

type SpectrumType struct{ Spectrum, FileName string }

func (s SpectrumType) Str() string { return fmt.Sprintf("%s#%s", s.Spectrum, s.FileName) }
func (p PeptideIdentification) SpectrumFileName() SpectrumType {
	return SpectrumType{p.Spectrum, p.SpectrumFile}
}

// PeptideIdentification struct
type PeptideIdentification struct {
	Spectrum                         string
	SpectrumFile                     string
	Peptide                          string
	Protein                          string
	ModifiedPeptide                  string
	CompensationVoltage              string
	AssumedCharge                    uint8
	HitRank                          uint8
	NumberOfEnzymaticTermini         uint8
	NumberofMissedCleavages          uint8
	Index                            uint32
	UncalibratedPrecursorNeutralMass float64
	PrecursorNeutralMass             float64
	RetentionTime                    float64
	CalcNeutralPepMass               float64
	Massdiff                         float64
	Probability                      float64
	Expectation                      float64
	Xcorr                            float64
	DeltaCN                          float64
	SPRank                           float64
	Hyperscore                       float64
	Nextscore                        float64
	SpectralSim                      float64
	Rtscore                          float64
	IonMobility                      float64
	Intensity                        float64
	AlternativeProteins              map[string]int
	MSFragerLoc                      *MSFraggerLoc
	PTM                              *PTM
	Modifications                    mod.ModificationsSlice
}

type PTM struct {
	LocalizedPTMSites    map[string]int
	LocalizedPTMMassDiff map[string]string
}

type MSFraggerLoc struct {
	MSFragerLocalization                 string
	MSFraggerLocalizationScoreWithPTM    string
	MSFraggerLocalizationScoreWithoutPTM string
}

// PepIDList is a list of PeptideSpectrumMatch
type PepIDList []PeptideIdentification
type PepIDListPtrs []*PeptideIdentification

func ToPepIDListPtrs(p PepIDList) PepIDListPtrs {
	pptrs := make(PepIDListPtrs, len(p), len(p))
	for i := range p {
		pptrs[i] = &p[i]
	}
	return pptrs
}

// Len function for Sort
func (p PepIDList) Len() int {
	return len(p)
}

// Less function for Sort
func (p PepIDList) Less(i, j int) bool {
	return p[i].Probability > p[j].Probability
}

// Swap function for Sort
func (p PepIDList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Len function for Sort
func (p PepIDListPtrs) Len() int {
	return len(p)
}

// Less function for Sort
func (p PepIDListPtrs) Less(i, j int) bool {
	return p[i].Probability > p[j].Probability
}

// Swap function for Sort
func (p PepIDListPtrs) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Read is the main function for parsing pepxml data
func (p *PepXML) Read(f string) {

	var xml spc.PepXML
	logrus.Info("Parsing ", f)
	xml.Parse(f)

	var mpa = xml.MsmsPipelineAnalysis

	if len(mpa.AnalysisSummary) > 0 {
		p.FileName = path.Base(f)
		p.Database = string(mpa.MsmsRunSummary.SearchSummary.SearchDatabase.LocalPath)
		p.SpectraFile = fmt.Sprintf("%s%s", mpa.MsmsRunSummary.BaseName, mpa.MsmsRunSummary.RawData)

		var models []spc.DistributionPoint

		// collect distribution points from meta
		for _, i := range mpa.AnalysisSummary[0].PeptideprophetSummary.DistributionPoint {
			var m spc.DistributionPoint
			m.Fvalue = i.Fvalue
			m.Obs1Distr = i.Obs1Distr
			m.Model1PosDistr = i.Model1PosDistr
			m.Model1NegDistr = i.Model1NegDistr
			m.Obs2Distr = i.Obs2Distr
			m.Model2PosDistr = i.Model2PosDistr
			m.Model2NegDistr = i.Model2NegDistr
			m.Obs3Distr = i.Obs3Distr
			m.Model3PosDistr = i.Model3PosDistr
			m.Model3NegDistr = i.Model3NegDistr
			m.Obs4Distr = i.Obs4Distr
			m.Model4PosDistr = i.Model4PosDistr
			m.Model4NegDistr = i.Model4NegDistr
			m.Obs5Distr = i.Obs5Distr
			m.Model5PosDistr = i.Model5PosDistr
			m.Model5NegDistr = i.Model5NegDistr
			m.Obs6Distr = i.Obs6Distr
			m.Model6PosDistr = i.Model6PosDistr
			m.Model6NegDistr = i.Model6NegDistr
			m.Obs7Distr = i.Obs7Distr
			m.Model7PosDistr = i.Model7PosDistr
			m.Model7NegDistr = i.Model7NegDistr
			models = append(models, m)
		}

		p.Modifications.Index = make(map[string]mod.Modification)

		// get the search engine
		p.SearchEngine = string(mpa.MsmsRunSummary.SearchSummary.SearchEngine)
		if strings.Contains(string(mpa.MsmsRunSummary.SearchSummary.SearchEngineVersion), "MSFragger") {
			p.SearchEngine = "MSFragger"
		}

		// map internal modifications from file
		for _, i := range mpa.MsmsRunSummary.SearchSummary.AminoAcidModifications {

			key := fmt.Sprintf("%s#%.4f", i.AminoAcid, i.Mass)
			variableStr := string(i.Variable)
			if variableStr != "Y" && variableStr != "N" {
				panic(nil)
			}
			variable := variableStr == "Y"
			_, ok := p.Modifications.Index[key]
			if !ok {
				m := mod.Modification{
					Index: key,
					Type:  mod.Assigned,
					//MonoIsotopicMass: i.Mass,
					MassDiff:  uti.ToFixed(i.MassDiff, 4),
					Variable:  variable,
					AminoAcid: string(i.AminoAcid),
					//IsobaricMods: make(map[string]float64),
				}

				p.Modifications.Index[key] = m
			}
		}

		// map terminal modifications from file
		for _, i := range mpa.MsmsRunSummary.SearchSummary.TerminalModifications {

			key := fmt.Sprintf("%s-term#%.4f", strings.ToUpper(string(i.Terminus)), i.Mass)
			variableStr := string(i.Variable)
			if variableStr != "Y" && variableStr != "N" {
				panic(nil)
			}
			variable := variableStr == "Y"
			_, ok := p.Modifications.Index[key]
			if !ok {

				m := mod.Modification{
					Index: key,
					Type:  mod.Assigned,
					//MonoIsotopicMass:  i.Mass,
					MassDiff:  uti.ToFixed(i.MassDiff, 4),
					Variable:  variable,
					AminoAcid: fmt.Sprintf("%s-term", i.Terminus),
					//IsProteinTerminus: string(i.ProteinTerminus),
					//Terminus:          strings.ToLower(string(i.Terminus)),
					//IsobaricMods: make(map[string]float64),
				}

				p.Modifications.Index[key] = m
			}
		}

		for _, i := range xml.MsmsPipelineAnalysis.MsmsRunSummary.SearchSummary.Parameter {
			par := &spc.Parameter{
				Name:  i.Name,
				Value: i.Value,
			}
			p.SearchParameters = append(p.SearchParameters, *par)

		}

		//massDeviation := getMassDeviation(mpa.MsmsRunSummary.SpectrumQuery)

		// start processing spectra queries
		sq := mpa.MsmsRunSummary.SpectrumQuery
		p.PeptideIdentification = make(PepIDList, len(sq), len(sq))
		for idx, i := range sq {
			p.PeptideIdentification[idx] = processSpectrumQuery(i, p.Modifications, p.DecoyTag, p.FileName)
		}

		p.Prophet = string(mpa.AnalysisSummary[0].Analysis)
		p.Models = models

		// p.adjustMassDeviation()

		if len(p.PeptideIdentification) == 0 {
			msg.NoPSMFound(errors.New(f), "warning")
		}

	}
}

// ReadPepXMLInput reads one or more fies and organize the data into PSM list
func ReadPepXMLInput(xmlFile, decoyTag, temp string, models bool) (PepIDListPtrs, string) {

	var files = make(map[string]struct{})
	var params []spc.Parameter
	var modsIndex = make(map[string]mod.Modification)
	var searchEngine string

	if strings.Contains(xmlFile, "pep.xml") || strings.Contains(xmlFile, "pepXML") {
		files[xmlFile] = struct{}{}
	} else {

		list := uti.IOReadDir(xmlFile, "pep.xml")

		if len(list) == 0 {
			msg.NoParametersFound(errors.New("missing PeptideProphet pepXML files"), "fatal")
		}

		// in case both PeptideProphet and PTMProphet files are present, use
		// PTMProphet results and ignore peptide prophet.
		for _, i := range list {
			base := filepath.Base(i)
			if strings.Contains(base, ".mod.") {
				files[i] = struct{}{}
			}
		}

		// if no PptideProphet results are present, then use all PeptideProphet files.
		if len(files) == 0 {
			for _, i := range list {
				base := filepath.Base(i)
				if !strings.Contains(base, ".mod.") {
					files[i] = struct{}{}
				}
			}
		}

	}

	pepIdentList := make([]PepIDList, len(files))
	sortedFiles := make([]string, 0, len(files))
	for i := range files {
		sortedFiles = append(sortedFiles, i)
	}
	sort.Strings(sortedFiles)
	mu := sync.Mutex{}
	processSinglePepXML := func(idx int, i string) {
		var p PepXML
		p.DecoyTag = decoyTag
		p.Read(i)
		if idx == 0 {
			params = p.SearchParameters
		}

		// print models
		if models {
			if strings.EqualFold(p.Prophet, "interprophet") {
				logrus.Error("Cannot print models for interprophet files")
			} else {
				logrus.Info("Printing models")
				go p.ReportModels(temp, filepath.Base(i))
				time.Sleep(time.Second * 3)
			}
		}

		pepIdentList[idx] = p.PeptideIdentification
		mu.Lock()
		for _, k := range p.Modifications.Index {
			if _, ok := modsIndex[k.Index]; !ok {
				modsIndex[k.Index] = k
			}
		}
		mu.Unlock()
		if idx == 0 {
			searchEngine = p.SearchEngine
		}
	}
	wg := sync.WaitGroup{}
	parallelism := 6
	parallelismTokens := make(chan struct{}, parallelism)
	wg.Add(len(sortedFiles))
	for idx, i := range sortedFiles {
		go func(idx int, i string) {
			defer wg.Done()
			parallelismTokens <- struct{}{}
			processSinglePepXML(idx, i)
			<-parallelismTokens
		}(idx, i)
	}
	wg.Wait()
	// create a "fake" global pepXML comprising all data
	var pepXML PepXML4Serialiazation
	pepXML.DecoyTag = decoyTag
	pepXML.SearchParameters = params
	pepXML.PeptideIdentification = make(PepIDListPtrs, 0)
	for _, pepIdent := range pepIdentList {
		for idx := range pepIdent {
			pepXML.PeptideIdentification = append(pepXML.PeptideIdentification, &pepIdent[idx])
		}
	}
	pepXML.Modifications.Index = modsIndex

	// promoting Spectra that matches to both decoys and targets to TRUE hits
	pepXML.PromoteProteinIDs()

	// serialize all pep files
	sort.Sort(pepXML.PeptideIdentification)
	pepXML.Serialize()

	return pepXML.PeptideIdentification, searchEngine
}

func processSpectrumQuery(sq spc.SpectrumQuery, mods mod.Modifications, decoyTag, FileName string) PeptideIdentification {

	var psm PeptideIdentification
	//psm.Modifications.Index = make(map[string]mod.Modification)
	psm.AlternativeProteins = make(map[string]int)

	psm.Index = sq.Index
	psm.SpectrumFile = FileName
	//psm.Scan = sq.StartScan
	psm.AssumedCharge = sq.AssumedCharge
	psm.RetentionTime = sq.RetentionTimeSec
	psm.IonMobility = sq.IonMobility
	psm.CompensationVoltage = sq.CompensationVoltage

	if sq.UncalibratedPrecursorNeutralMass > 0 {
		psm.PrecursorNeutralMass = sq.PrecursorNeutralMass
		psm.UncalibratedPrecursorNeutralMass = sq.UncalibratedPrecursorNeutralMass
	} else {
		psm.PrecursorNeutralMass = sq.PrecursorNeutralMass
		psm.UncalibratedPrecursorNeutralMass = sq.PrecursorNeutralMass
	}

	for _, i := range sq.SearchResult.SearchHit {

		psm.HitRank = i.HitRank
		//psm.PrevAA = string(i.PrevAA)
		//psm.NextAA = string(i.NextAA)
		//psm.MissedCleavages = i.MissedCleavages
		//psm.NumberTolTerm = i.TotalTerm
		//psm.NumberTotalProteins = i.TotalProteins
		//psm.TotalNumberIons = i.TotalIons
		//psm.NumberMatchedIons = i.MatchedIons
		//psm.IsRejected = i.IsRejected

		psm.Peptide = string(i.Peptide)
		psm.Protein = string(i.Protein)
		psm.CalcNeutralPepMass = i.CalcNeutralPepMass

		//psm.Massdiff = uti.ToFixed((i.Massdiff - massDeviation), 4)
		psm.Massdiff = uti.ToFixed(i.Massdiff, 4)

		psm.NumberofMissedCleavages = i.MissedCleavages
		psm.NumberOfEnzymaticTermini = i.TotalTerm

		for _, j := range i.AnalysisResult {

			if string(j.Analysis) == "peptideprophet" {

				psm.Probability = j.PeptideProphetResult.Probability

				for _, k := range j.PeptideProphetResult.SearchScoreSummary.Parameter {

					if k.Name == "massd" {
						//psm.IsoMassD, _ = strconv.Atoi(k.Value)
					}
				}
			}

			if string(j.Analysis) == "interprophet" {
				psm.Probability = j.InterProphetResult.Probability
			}

			if string(j.Analysis) == "ptmprophet" {
				if len(j.PTMProphetResult) != 0 {
					psm.PTM = &PTM{LocalizedPTMSites: make(map[string]int), LocalizedPTMMassDiff: make(map[string]string)}
					for _, k := range j.PTMProphetResult {
						psm.PTM.LocalizedPTMSites[string(k.PTM)] = len(k.ModAminoAcidProbability)
						psm.PTM.LocalizedPTMMassDiff[string(k.PTM)] = string(k.PTMPeptide)
					}
				}
			}
		}

		for _, j := range i.AlternativeProteins {
			psm.AlternativeProteins[string(j.Protein)]++
		}

		for _, j := range i.Score {
			if string(j.Name) == "expect" {
				eValue, _ := uti.ParseFloat(j.Value)
				psm.Expectation = eValue
			} else if string(j.Name) == "xcorr" {
				value, _ := strconv.ParseFloat(j.Value, 64)
				psm.Xcorr = value
			} else if string(j.Name) == "deltacn" {
				value, _ := strconv.ParseFloat(j.Value, 64)
				psm.DeltaCN = value
			} else if string(j.Name) == "deltacnstar" {
				//value, _ := strconv.ParseFloat(j.Value, 64)
				//psm.DeltaCNStar = value
			} else if string(j.Name) == "spscore" {
				//value, _ := strconv.ParseFloat(j.Value, 64)
				//psm.SPScore = value
			} else if string(j.Name) == "sprank" {
				value, _ := strconv.ParseFloat(j.Value, 64)
				psm.SPRank = value
			} else if string(j.Name) == "hyperscore" {
				value, _ := strconv.ParseFloat(j.Value, 64)
				psm.Hyperscore = value
			} else if string(j.Name) == "nextscore" {
				value, _ := strconv.ParseFloat(j.Value, 64)
				psm.Nextscore = value
			} else if string(j.Name) == "spectralsim" {
				value, _ := strconv.ParseFloat(j.Value, 64)
				psm.SpectralSim = value
			} else if string(j.Name) == "rtscore" {
				value, _ := strconv.ParseFloat(j.Value, 64)
				psm.Rtscore = value
			}
		}

		//psm.LocalizationRange = i.PTMResult.LocalizationPeptide
		if len(i.PTMResult.LocalizationPeptide+i.PTMResult.BestScoreWithPTM+i.PTMResult.ScoreWithoutPTM) != 0 {
			psm.MSFragerLoc = &MSFraggerLoc{
				MSFragerLocalization:                 i.PTMResult.LocalizationPeptide,
				MSFraggerLocalizationScoreWithPTM:    i.PTMResult.BestScoreWithPTM,
				MSFraggerLocalizationScoreWithoutPTM: i.PTMResult.ScoreWithoutPTM}
		}

		// to be able to accept multiple entries with the same spectrum name, we fuse the
		// file name to the spectrum name. This is going to be used as an identifiable attribute
		// Before reporting the filtered PSMs, the file name is removed from the spectrum name.
		//psm.Spectrum = fmt.Sprintf("%s#%s", psm.Spectrum, FileName)
		psm.Spectrum = string(sq.Spectrum)

		psm.mapModsFromPepXML(i.ModificationInfo, mods)
	}

	return psm
}

// mapModsFromPepXML receives a pepXML struct with modifications and adds them to the given struct
func (p *PeptideIdentification) mapModsFromPepXML(m spc.ModificationInfo, mods mod.Modifications) {

	p.ModifiedPeptide = string(m.ModifiedPeptide)
	pModificationsIndex := make(map[string]mod.Modification)
	for _, i := range m.ModAminoacidMass {

		aa := strings.Split(p.Peptide, "")
		key := fmt.Sprintf("%s#%.4f", aa[i.Position-1], i.Mass)

		// This is related to a rounding issue that prevents the correct mapping between
		// PTMProphet and MSFragger masses
		keyPlus := fmt.Sprintf("%s#%.4f", aa[i.Position-1], i.Mass+0.0001)
		keyMinus := fmt.Sprintf("%s#%.4f", aa[i.Position-1], i.Mass-0.0001)

		v, ok := mods.Index[key]
		if ok {
			m := v
			newKey := fmt.Sprintf("%s#%d#%.4f", aa[i.Position-1], i.Position, i.Mass)
			m.Index = newKey
			m.Position = i.Position
			//m.IsobaricMods = make(map[string]float64)
			pModificationsIndex[newKey] = m
		} else {

			v, ok = mods.Index[keyPlus]
			if ok {
				m := v
				newKey := fmt.Sprintf("%s#%d#%.4f", aa[i.Position-1], i.Position, i.Mass)
				m.Index = newKey
				m.Position = i.Position
				//m.IsobaricMods = make(map[string]float64)
				pModificationsIndex[newKey] = m
			}

			v, ok = mods.Index[keyMinus]
			if ok {
				m := v
				newKey := fmt.Sprintf("%s#%d#%.4f", aa[i.Position-1], i.Position, i.Mass)
				m.Index = newKey
				m.Position = i.Position
				//m.IsobaricMods = make(map[string]float64)
				pModificationsIndex[newKey] = m
			}
		}
	}

	// n-terminal modifications
	if m.ModNTermMass != 0 {
		key := fmt.Sprintf("N-term#%.4f", m.ModNTermMass)
		v, ok := mods.Index[key]
		if ok {
			m := v
			m.AminoAcid = "N-term"
			//m.IsobaricMods = make(map[string]float64)
			pModificationsIndex[key] = m
		}

		// this rule was added because PTMProphet is changing the mod_nterm_mass
		// in the PSM to something that does not exists in the header table.
		if strings.Contains(key, "305") {
			key = "N-term#305.2150"
			v, ok := mods.Index[key]
			if ok {
				m := v
				m.AminoAcid = "N-term"
				//m.IsobaricMods = make(map[string]float64)
				pModificationsIndex[key] = m
			}
		}

	}

	// c-terminal modifications
	if m.ModCTermMass != 0 {
		key := fmt.Sprintf("C-term#%.4f", m.ModCTermMass)
		v, ok := mods.Index[key]
		if ok {
			m := v
			m.AminoAcid = "C-term"
			//m.IsobaricMods = make(map[string]float64)
			pModificationsIndex[key] = m
		}
	}

	// if isotopicCorr >= 0.036386 || isotopicCorr <= -0.036386 {
	key := fmt.Sprintf("%.4f", p.Massdiff)
	_, ok := pModificationsIndex[key]
	if !ok {
		m := mod.Modification{
			Index:    key,
			Name:     "Unknown",
			Type:     mod.Observed,
			MassDiff: p.Massdiff,
			//IsobaricMods: make(map[string]float64),
		}
		pModificationsIndex[key] = m
	}
	if len(pModificationsIndex) != 0 {
		p.Modifications = mod.Modifications{Index: pModificationsIndex}.ToSlice()
	}
}

// getMassDeviation calculates the mass deviation for a pepXML file based on the 0 mass difference
// func getMassDeviation(sq []spc.SpectrumQuery) float64 {

// 	var countZero int
// 	var massZero float64
// 	var adjustedMass float64

// 	for _, i := range sq {
// 		for _, j := range i.SearchResult.SearchHit {
// 			if math.Abs(j.Massdiff) >= -0.1 && math.Abs(j.Massdiff) <= 0.1 {
// 				countZero++
// 				massZero += j.Massdiff
// 			}
// 		}
// 	}

// 	adjustedMass = massZero / float64(countZero)

// 	return adjustedMass
// }

// PromoteProteinIDs changes the identification in cases where the reference protein is a decoy and
// the alternative proteins contains target proteins.
func (p *PepXML4Serialiazation) PromoteProteinIDs() {

	for i := range p.PeptideIdentification {

		var current string
		var alt string
		//var altNTT int
		var list = make(map[string]int)
		var isUniProt bool

		if strings.Contains(p.PeptideIdentification[i].Protein, p.DecoyTag) {

			current = p.PeptideIdentification[i].Protein

			for j := range p.PeptideIdentification[i].AlternativeProteins {

				if strings.Contains(j, "sp|") {
					isUniProt = true
				}

				if !strings.HasPrefix(j, p.DecoyTag) {
					list[j]++
				}
			}

		}

		if len(list) > 0 {

			// if a Uniprot database is used we give preference to SwissProt proteins
			if isUniProt {
				for k := range list {
					if strings.HasPrefix(k, "sp|") {
						alt = k
						break
					} else {
						alt = k
					}
				}
				p.PeptideIdentification[i].Protein = alt

				// remove the replaces protein from the alternative proteins list
				//p.PeptideIdentification[i].AlternativeProteins[list[alt]] = p.PeptideIdentification[i].AlternativeProteins[len(p.PeptideIdentification[i].AlternativeProteins)-1]
				//p.PeptideIdentification[i].AlternativeProteins[len(p.PeptideIdentification[i].AlternativeProteins)-1] = ""
				//p.PeptideIdentification[i].AlternativeProteins = p.PeptideIdentification[i].AlternativeProteins[:len(p.PeptideIdentification[i].AlternativeProteins)-1]

				// add the replaces current to the list
				p.PeptideIdentification[i].AlternativeProteins[current]++

			} else {
				for k := range list {
					alt = k
					break
				}
				p.PeptideIdentification[i].Protein = alt

				// remove the replaces protein from the alternative proteins list
				//p.PeptideIdentification[i].AlternativeProteins[list[alt]] = p.PeptideIdentification[i].AlternativeProteins[len(p.PeptideIdentification[i].AlternativeProteins)-1]
				//p.PeptideIdentification[i].AlternativeProteins[len(p.PeptideIdentification[i].AlternativeProteins)-1] = ""
				//p.PeptideIdentification[i].AlternativeProteins = p.PeptideIdentification[i].AlternativeProteins[:len(p.PeptideIdentification[i].AlternativeProteins)-1]

				// add the replaces current to the list
				p.PeptideIdentification[i].AlternativeProteins[current]++
			}

		}
	}

}

// ReportModels creates PNG images using the PeptideProphet TD score distribution
func (p *PepXML) ReportModels(session, name string) {

	var xAxis []float64

	for i := range p.Models {
		xAxis = append(xAxis, p.Models[i].Fvalue)
	}

	name = strings.Replace(name, ".pep", "", -1)
	name = strings.Replace(name, ".Pep", "", -1)
	name = strings.Replace(name, ".xml", "", -1)

	for i := 2; i < 8; i++ {

		var obs []float64
		var pos []float64
		var neg []float64

		if i == 2 {

			path := fmt.Sprintf("%s%s%s_2.png", session, string(filepath.Separator), name)
			for j := range p.Models {
				obs = append(obs, p.Models[j].Obs2Distr)
				pos = append(pos, p.Models[j].Model2PosDistr)
				neg = append(neg, p.Models[j].Model2NegDistr)
			}
			printModel("2", path, xAxis, obs, pos, neg)

		} else if i == 3 {

			path := fmt.Sprintf("%s%s%s_3.png", session, string(filepath.Separator), name)
			for j := range p.Models {
				obs = append(obs, p.Models[j].Obs3Distr)
				pos = append(pos, p.Models[j].Model3PosDistr)
				neg = append(neg, p.Models[j].Model3NegDistr)
			}
			printModel("3", path, xAxis, obs, pos, neg)

		} else if i == 4 {

			path := fmt.Sprintf("%s%s%s_4.png", session, string(filepath.Separator), name)
			for j := range p.Models {
				obs = append(obs, p.Models[j].Obs4Distr)
				pos = append(pos, p.Models[j].Model4PosDistr)
				neg = append(neg, p.Models[j].Model4NegDistr)
			}
			printModel("4", path, xAxis, obs, pos, neg)

		} else if i == 5 {

			path := fmt.Sprintf("%s%s%s_5.png", session, string(filepath.Separator), name)
			for j := range p.Models {
				obs = append(obs, p.Models[j].Obs5Distr)
				pos = append(pos, p.Models[j].Model5PosDistr)
				neg = append(neg, p.Models[j].Model5NegDistr)
			}
			printModel("5", path, xAxis, obs, pos, neg)

		} else if i == 6 {

			path := fmt.Sprintf("%s%s%s_6.png", session, string(filepath.Separator), name)
			for j := range p.Models {
				obs = append(obs, p.Models[j].Obs6Distr)
				pos = append(pos, p.Models[j].Model6PosDistr)
				neg = append(neg, p.Models[j].Model6NegDistr)
			}
			printModel("6", path, xAxis, obs, pos, neg)

		} else if i == 7 {

			path := fmt.Sprintf("%s%s%s_7.png", session, string(filepath.Separator), name)
			for j := range p.Models {
				obs = append(obs, p.Models[j].Obs7Distr)
				pos = append(pos, p.Models[j].Model7PosDistr)
				neg = append(neg, p.Models[j].Model7NegDistr)
			}
			printModel("7", path, xAxis, obs, pos, neg)

		}
	}

}

func printModel(v, path string, xAxis, obs, pos, neg []float64) {

	p, e := plot.New()

	if e != nil {

		msg.Plotter(e, "fatal")

	} else {

		p.Title.Text = "FVAL" + v
		p.X.Label.Text = "FVAL"
		p.Y.Label.Text = "Density"

		obsPts := make(plotter.XYs, len(xAxis))
		posPts := make(plotter.XYs, len(xAxis))
		negPts := make(plotter.XYs, len(xAxis))
		for i := range obs {
			obsPts[i].X = xAxis[i]
			obsPts[i].Y = obs[i]
			posPts[i].X = xAxis[i]
			posPts[i].Y = pos[i]
			negPts[i].X = xAxis[i]
			negPts[i].Y = neg[i]
		}

		e = plotutil.AddLinePoints(p, "Observed", obsPts, "Positive", posPts, "Negative", negPts)
		if e != nil {
			panic(e)
		}

		// Save the plot to a PNG file.
		if err := p.Save(8*vg.Inch, 6*vg.Inch, path); err != nil {
			panic(err)
		}

		// copy to work directory
		sys.CopyFile(path, filepath.Base(path))
	}

}

// Serialize converts the whole structure to a msgpack file
func (p *PepXML4Serialiazation) Serialize() {
	sys.Serialize(p, sys.PepxmlBin())
}

// Restore reads philosopher results files and restore the data sctructure
func (p *PepXML) Restore() {
	sys.Restore(p, sys.PepxmlBin(), false)
}

// Serialize converts the whle structure to a gob file
func (p *PepIDList) Serialize(level string) {

	var dest string

	if level == "psm" {
		dest = sys.PSMBin()
	} else if level == "pep" {
		dest = sys.PepBin()
	} else if level == "ion" {
		dest = sys.IonBin()
	} else {
		msg.Custom(errors.New("cannot determine binary data class"), "fatal")
	}
	sys.Serialize(p, dest)
}

// Serialize converts the whle structure to a gob file
func (p *PepIDListPtrs) Serialize(level string) {

	var dest string

	if level == "psm" {
		dest = sys.PSMBin()
	} else if level == "pep" {
		dest = sys.PepBin()
	} else if level == "ion" {
		dest = sys.IonBin()
	} else {
		msg.Custom(errors.New("cannot determine binary data class"), "fatal")
	}
	sys.Serialize(p, dest)
}

// Restore reads philosopher results files and restore the data sctructure
func (p *PepIDList) Restore(level string) {

	var dest string

	if level == "psm" {
		dest = sys.PSMBin()
	} else if level == "pep" {
		dest = sys.PepBin()
	} else if level == "ion" {
		dest = sys.IonBin()
	} else {
		msg.Custom(errors.New("cannot determine binary data class"), "fatal")
	}
	sys.Restore(p, dest, false)
}
