package id

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"philosopher/lib/uti"

	"philosopher/lib/msg"

	"philosopher/lib/mod"
	"philosopher/lib/spc"
	"philosopher/lib/sys"

	"github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack"
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
	SearchParameters      []spc.Parameter
	Database              string
	Prophet               string
	Modifications         mod.Modifications
	Models                []spc.DistributionPoint
	PeptideIdentification PepIDList
}

// PeptideIdentification struct
type PeptideIdentification struct {
	Index                            uint32
	Spectrum                         string
	SpectrumFile                     string
	Scan                             int
	Peptide                          string
	Protein                          string
	ModifiedPeptide                  string
	AlternativeProteins              []string
	AlternativeProteinsIndexed       map[string]int
	AssumedCharge                    uint8
	PrevAA                           string
	NextAA                           string
	HitRank                          uint8
	MissedCleavages                  uint8
	NumberTolTerm                    uint8
	NumberOfEnzymaticTermini         uint8
	NumberTotalProteins              uint16
	TotalNumberIons                  uint16
	NumberMatchedIons                uint16
	NumberofMissedCleavages          int
	UncalibratedPrecursorNeutralMass float64
	PrecursorNeutralMass             float64
	PrecursorExpMass                 float64
	RetentionTime                    float64
	CalcNeutralPepMass               float64
	Massdiff                         float64
	LocalizedPTMSites                map[string]int
	LocalizedPTMMassDiff             map[string]string
	Probability                      float64
	IsoMassD                         int
	Expectation                      float64
	Xcorr                            float64
	DeltaCN                          float64
	DeltaCNStar                      float64
	SPScore                          float64
	SPRank                           float64
	Hyperscore                       float64
	Nextscore                        float64
	DiscriminantValue                float64
	Intensity                        float64
	IonMobility                      float64
	IsRejected                       uint8
	Modifications                    mod.Modifications
}

// PepIDList is a list of PeptideSpectrumMatch
type PepIDList []PeptideIdentification

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

// Read is the main function for parsing pepxml data
func (p *PepXML) Read(f string) {

	var xml spc.PepXML
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

			_, ok := p.Modifications.Index[key]
			if !ok {

				m := mod.Modification{
					Index:            key,
					Type:             "Assigned",
					MonoIsotopicMass: i.Mass,
					MassDiff:         i.MassDiff,
					Variable:         string(i.Variable),
					AminoAcid:        string(i.AminoAcid),
					IsobaricMods:     make(map[string]float64),
				}

				p.Modifications.Index[key] = m
			}
		}

		// map terminal modifications from file
		for _, i := range mpa.MsmsRunSummary.SearchSummary.TerminalModifications {

			key := fmt.Sprintf("%s-term#%.4f", strings.ToUpper(string(i.Terminus)), i.Mass)

			_, ok := p.Modifications.Index[key]
			if !ok {

				m := mod.Modification{
					Index:             key,
					Type:              "Assigned",
					MonoIsotopicMass:  i.Mass,
					MassDiff:          i.MassDiff,
					Variable:          string(i.Variable),
					AminoAcid:         fmt.Sprintf("%s-term", i.Terminus),
					IsProteinTerminus: string(i.ProteinTerminus),
					Terminus:          strings.ToLower(string(i.Terminus)),
					IsobaricMods:      make(map[string]float64),
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

		massDeviation := getMassDeviation(mpa.MsmsRunSummary.SpectrumQuery)

		// start processing spectra queries
		var psmlist PepIDList
		sq := mpa.MsmsRunSummary.SpectrumQuery
		for _, i := range sq {
			psm := processSpectrumQuery(i, massDeviation, p.Modifications, p.DecoyTag, p.FileName)
			psmlist = append(psmlist, psm)
		}

		p.PeptideIdentification = psmlist
		p.Prophet = string(mpa.AnalysisSummary[0].Analysis)
		p.Models = models

		// p.adjustMassDeviation()

		if len(psmlist) == 0 {
			msg.NoPSMFound(errors.New(f), "warning")
		}

	}

	return
}

// ReadPepXMLInput reads one or more fies and organize the data into PSM list
func ReadPepXMLInput(xmlFile, decoyTag, temp string, models bool) (PepIDList, string) {

	var files = make(map[string]uint8)
	var fileCheckList []string
	var fileCheckFlag bool
	var pepIdent PepIDList
	var mods []mod.Modification
	var params []spc.Parameter
	var modsIndex = make(map[string]mod.Modification)
	var searchEngine string

	if strings.Contains(xmlFile, "pep.xml") || strings.Contains(xmlFile, "pepXML") {
		fileCheckList = append(fileCheckList, xmlFile)
		files[xmlFile] = 0
	} else {
		glob := fmt.Sprintf("%s%s*pep.xml", xmlFile, string(filepath.Separator))
		list, _ := filepath.Glob(glob)

		if len(list) == 0 {
			msg.NoParametersFound(errors.New("missing pepXML files"), "fatal")
		}

		for _, i := range list {
			absPath, _ := filepath.Abs(i)
			files[absPath] = 0
			fileCheckList = append(fileCheckList, absPath)

			if strings.Contains(i, "mod") {
				fileCheckFlag = true
			}
		}

	}

	//spew.Dump(files)

	// verify if the we have interact and interact.mod files for parsing.
	// To avoid reading both files, we keep the mod one and discard the other.
	if fileCheckFlag == true {
		for _, i := range fileCheckList {
			i = strings.Replace(i, "mod.", "", 1)
			_, ok := files[i]
			if ok {
				delete(files, i)
			}
		}
	}

	for i := range files {
		var p PepXML
		p.DecoyTag = decoyTag
		p.Read(i)

		params = p.SearchParameters

		// print models
		if models == true {
			if strings.EqualFold(p.Prophet, "interprophet") {
				logrus.Error("Cannot print models for interprophet files")
			} else {
				logrus.Info("Printing models")
				go p.ReportModels(temp, filepath.Base(i))
				time.Sleep(time.Second * 3)
			}
		}

		pepIdent = append(pepIdent, p.PeptideIdentification...)

		for _, k := range p.Modifications.Index {
			_, ok := modsIndex[k.Index]
			if !ok {
				mods = append(mods, k)
				modsIndex[k.Index] = k
			}
		}

		searchEngine = p.SearchEngine
	}

	// create a "fake" global pepXML comprising all data
	var pepXML PepXML
	pepXML.DecoyTag = decoyTag
	pepXML.SearchParameters = params
	pepXML.PeptideIdentification = pepIdent
	pepXML.Modifications.Index = modsIndex

	// promoting Spectra that matches to both decoys and targets to TRUE hits
	pepXML.PromoteProteinIDs()

	// serialize all pep files
	sort.Sort(pepXML.PeptideIdentification)
	pepXML.Serialize()

	return pepIdent, searchEngine
}

func processSpectrumQuery(sq spc.SpectrumQuery, massDeviation float64, mods mod.Modifications, decoyTag, FileName string) PeptideIdentification {

	var psm PeptideIdentification
	psm.Modifications.Index = make(map[string]mod.Modification)
	psm.AlternativeProteinsIndexed = make(map[string]int)

	psm.Index = sq.Index
	psm.SpectrumFile = FileName
	psm.Spectrum = string(sq.Spectrum)
	psm.Scan = sq.StartScan
	psm.AssumedCharge = sq.AssumedCharge
	psm.RetentionTime = sq.RetentionTimeSec
	psm.IonMobility = sq.IonMobility

	if sq.UncalibratedPrecursorNeutralMass > 0 {
		psm.PrecursorNeutralMass = sq.PrecursorNeutralMass
		psm.UncalibratedPrecursorNeutralMass = sq.UncalibratedPrecursorNeutralMass
	} else {
		psm.PrecursorNeutralMass = sq.PrecursorNeutralMass
		psm.UncalibratedPrecursorNeutralMass = sq.PrecursorNeutralMass
	}

	for _, i := range sq.SearchResult.SearchHit {

		psm.HitRank = i.HitRank
		psm.PrevAA = string(i.PrevAA)
		psm.NextAA = string(i.NextAA)
		psm.MissedCleavages = i.MissedCleavages
		psm.NumberTolTerm = i.TotalTerm
		psm.NumberTotalProteins = i.TotalProteins
		psm.TotalNumberIons = i.TotalIons
		psm.NumberMatchedIons = i.MatchedIons
		psm.IsRejected = i.IsRejected

		psm.Peptide = string(i.Peptide)
		psm.Protein = string(i.Protein)
		psm.CalcNeutralPepMass = i.CalcNeutralPepMass
		psm.Massdiff = uti.ToFixed((i.Massdiff - massDeviation), 4)

		psm.NumberofMissedCleavages = int(i.MissedCleavages)
		psm.NumberOfEnzymaticTermini = i.TotalTerm

		for _, j := range i.AnalysisResult {

			if string(j.Analysis) == "peptideprophet" {

				psm.Probability = j.PeptideProphetResult.Probability

				for _, k := range j.PeptideProphetResult.SearchScoreSummary.Parameter {

					if k.Name == "massd" {
						psm.IsoMassD, _ = strconv.Atoi(k.Value)
					}

					// if k.Name == "ntt" {
					// 	psm.NumberOfEnzymaticTermini, _ = strconv.Atoi(k.Value)
					// }

					// if k.Name == "nmc" {
					// 	psm.NMC, _ = strconv.Atoi(k.Value)
					// }
				}
			}

			if string(j.Analysis) == "interprophet" {
				psm.Probability = j.InterProphetResult.Probability
			}

			if string(j.Analysis) == "ptmprophet" {
				psm.LocalizedPTMSites = make(map[string]int)
				psm.LocalizedPTMMassDiff = make(map[string]string)
				for _, k := range j.PTMProphetResult {
					psm.LocalizedPTMSites[string(k.PTM)] = len(k.ModAminoAcidProbability)
					psm.LocalizedPTMMassDiff[string(k.PTM)] = string(k.PTMPeptide)
				}
			}
		}

		for _, j := range i.AlternativeProteins {
			psm.AlternativeProteins = append(psm.AlternativeProteins, string(j.Protein))
			psm.AlternativeProteinsIndexed[string(j.Protein)]++
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
				value, _ := strconv.ParseFloat(j.Value, 64)
				psm.DeltaCNStar = value
			} else if string(j.Name) == "spscore" {
				value, _ := strconv.ParseFloat(j.Value, 64)
				psm.SPScore = value
			} else if string(j.Name) == "sprank" {
				value, _ := strconv.ParseFloat(j.Value, 64)
				psm.SPRank = value
			} else if string(j.Name) == "hyperscore" {
				value, _ := strconv.ParseFloat(j.Value, 64)
				psm.Hyperscore = value
			} else if string(j.Name) == "nextscore" {
				value, _ := strconv.ParseFloat(j.Value, 64)
				psm.Nextscore = value
			}
		}

		// to be able to accept multiple entries with the same spectrum name, we fuse the
		// file name to the spectrum name. This is going to be used as an identifiable attribute
		// Before reporting the filtered PSMs, the file name is removed from the spectrum name.
		psm.Spectrum = fmt.Sprintf("%s#%s", psm.Spectrum, FileName)

		psm.mapModsFromPepXML(i.ModificationInfo, mods)
	}

	return psm
}

// mapModsFromPepXML receives a pepXML struct with modifications and adds them to the given struct
func (p *PeptideIdentification) mapModsFromPepXML(m spc.ModificationInfo, mods mod.Modifications) {

	p.ModifiedPeptide = string(m.ModifiedPeptide)

	for _, i := range m.ModAminoacidMass {
		aa := strings.Split(p.Peptide, "")
		key := fmt.Sprintf("%s#%.4f", aa[i.Position-1], i.Mass)
		v, ok := mods.Index[key]
		if ok {
			m := v
			newKey := fmt.Sprintf("%s#%d#%.4f", aa[i.Position-1], i.Position, i.Mass)
			m.Index = newKey
			m.Position = strconv.Itoa(i.Position)
			m.IsobaricMods = make(map[string]float64)
			p.Modifications.Index[newKey] = m
		}
	}

	// n-terminal modifications
	if m.ModNTermMass != 0 {
		key := fmt.Sprintf("N-term#%.4f", m.ModNTermMass)
		v, ok := mods.Index[key]
		if ok {
			m := v
			m.AminoAcid = "N-term"
			m.IsobaricMods = make(map[string]float64)
			p.Modifications.Index[key] = m
		}
	}

	// c-terminal modifications
	if m.ModCTermMass != 0 {
		key := fmt.Sprintf("C-term#%.4f", m.ModCTermMass)
		v, ok := mods.Index[key]
		if ok {
			m := v
			m.AminoAcid = "C-term"
			m.IsobaricMods = make(map[string]float64)
			p.Modifications.Index[key] = m
		}
	}

	// if isotopicCorr >= 0.036386 || isotopicCorr <= -0.036386 {
	key := fmt.Sprintf("%.4f", p.Massdiff)
	_, ok := p.Modifications.Index[key]
	if !ok {
		m := mod.Modification{
			Index:        key,
			Name:         "Unknown",
			Type:         "Observed",
			MassDiff:     p.Massdiff,
			IsobaricMods: make(map[string]float64),
		}
		p.Modifications.Index[key] = m
	}

	//}

	return
}

// getMassDeviation calculates the mass deviation for a pepXML file based on the 0 mass difference
func getMassDeviation(sq []spc.SpectrumQuery) float64 {

	var countZero int
	var massZero float64
	var adjustedMass float64

	for _, i := range sq {
		for _, j := range i.SearchResult.SearchHit {
			if math.Abs(j.Massdiff) >= -0.1 && math.Abs(j.Massdiff) <= 0.1 {
				countZero++
				massZero += j.Massdiff
			}
		}
	}

	adjustedMass = massZero / float64(countZero)

	return adjustedMass
}

// PromoteProteinIDs changes the identification in cases where the reference protein is a decoy and
// the alternative proteins contains target proteins.
func (p *PepXML) PromoteProteinIDs() {

	for i := range p.PeptideIdentification {

		var current string
		var alt string
		//var altNTT int
		var list = make(map[string]int)
		var isUniProt bool

		if strings.Contains(p.PeptideIdentification[i].Protein, p.DecoyTag) {

			current = p.PeptideIdentification[i].Protein

			for j := range p.PeptideIdentification[i].AlternativeProteins {

				if strings.Contains(p.PeptideIdentification[i].AlternativeProteins[j], "sp|") {
					isUniProt = true
				}

				if !strings.HasPrefix(p.PeptideIdentification[i].AlternativeProteins[j], p.DecoyTag) {
					list[p.PeptideIdentification[i].AlternativeProteins[j]] = j
				}
			}

		}

		if len(list) > 0 {

			// if a Uniprot database is used we give preference to SwissProt proteins
			if isUniProt == true {
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
				p.PeptideIdentification[i].AlternativeProteins[list[alt]] = p.PeptideIdentification[i].AlternativeProteins[len(p.PeptideIdentification[i].AlternativeProteins)-1]
				p.PeptideIdentification[i].AlternativeProteins[len(p.PeptideIdentification[i].AlternativeProteins)-1] = ""
				p.PeptideIdentification[i].AlternativeProteins = p.PeptideIdentification[i].AlternativeProteins[:len(p.PeptideIdentification[i].AlternativeProteins)-1]

				// add the replaces current to the list
				p.PeptideIdentification[i].AlternativeProteins = append(p.PeptideIdentification[i].AlternativeProteins, current)

			} else {
				for k := range list {
					alt = k
					break
				}
				p.PeptideIdentification[i].Protein = alt

				// remove the replaces protein from the alternative proteins list
				p.PeptideIdentification[i].AlternativeProteins[list[alt]] = p.PeptideIdentification[i].AlternativeProteins[len(p.PeptideIdentification[i].AlternativeProteins)-1]
				p.PeptideIdentification[i].AlternativeProteins[len(p.PeptideIdentification[i].AlternativeProteins)-1] = ""
				p.PeptideIdentification[i].AlternativeProteins = p.PeptideIdentification[i].AlternativeProteins[:len(p.PeptideIdentification[i].AlternativeProteins)-1]

				// add the replaces current to the list
				p.PeptideIdentification[i].AlternativeProteins = append(p.PeptideIdentification[i].AlternativeProteins, current)
			}

		}
	}

	return
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

	return
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

	return
}

// tdclassifier identifies a PSM as target or Decoy based on the
// presence of the TAG string on <protein> and <alternative_proteins>
func tdclassifier(p PeptideIdentification, tag string) bool {

	// default for TRUE ( DECOY)
	//var class = true
	var class bool

	if strings.HasPrefix(string(p.Protein), tag) {
		class = true
	} else {
		class = false
	}

	for i := range p.AlternativeProteins {
		if !strings.HasPrefix(p.AlternativeProteins[i], tag) {
			class = false
		}
		break
	}

	return class
}

// Serialize converts the whle structure to a gob file
func (p *PepXML) Serialize() {

	b, e := msgpack.Marshal(&p)
	if e != nil {
		msg.MarshalFile(e, "fatal")
	}

	e = ioutil.WriteFile(sys.PepxmlBin(), b, sys.FilePermission())
	if e != nil {
		msg.WriteFile(e, "fatal")
	}

	return
}

// Restore reads philosopher results files and restore the data sctructure
func (p *PepXML) Restore() {

	b, e := ioutil.ReadFile(sys.PepxmlBin())
	if e != nil {
		msg.ReadFile(e, "warning")
	}

	e = msgpack.Unmarshal(b, &p)
	if e != nil {
		msg.DecodeMsgPck(e, "warning")
	}

	return
}

// Serialize converts the whle structure to a gob file
func (p *PepIDList) Serialize(level string) {

	var dest string

	if level == "psm" {
		dest = sys.PsmBin()
	} else if level == "pep" {
		dest = sys.PepBin()
	} else if level == "ion" {
		dest = sys.IonBin()
	} else {
		msg.Custom(errors.New("Cannot determine binary data class"), "fatal")
	}

	b, e := msgpack.Marshal(&p)
	if e != nil {
		msg.MarshalFile(e, "fatal")
	}

	e = ioutil.WriteFile(dest, b, sys.FilePermission())
	if e != nil {
		msg.WriteFile(e, "fatal")
	}

	return
}

// Restore reads philosopher results files and restore the data sctructure
func (p *PepIDList) Restore(level string) {

	var dest string

	if level == "psm" {
		dest = sys.PsmBin()
	} else if level == "pep" {
		dest = sys.PepBin()
	} else if level == "ion" {
		dest = sys.IonBin()
	} else {
		msg.Custom(errors.New("Cannot determine binary data class"), "fatal")
	}

	b, e := ioutil.ReadFile(dest)
	if e != nil {
		msg.ReadFile(e, "fatal")
	}

	e = msgpack.Unmarshal(b, &p)
	if e != nil {
		msg.DecodeMsgPck(e, "fatal")
	}

	return
}
