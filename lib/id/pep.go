package id

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/prvst/philosopher/lib/err"

	"github.com/prvst/philosopher/lib/bio"

	"github.com/prvst/philosopher/lib/mod"
	"github.com/prvst/philosopher/lib/spc"
	"github.com/prvst/philosopher/lib/sys"
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
	Index                uint32
	Spectrum             string
	Scan                 int
	Peptide              string
	Protein              string
	ModifiedPeptide      string
	AlternativeProteins  []string
	AssumedCharge        uint8
	PrevAA               string
	NextAA               string
	HitRank              uint8
	MissedCleavages      uint8
	NumberTolTerm        uint8
	NumberTotalProteins  uint16
	TotalNumberIons      uint16
	NumberMatchedIons    uint16
	PrecursorNeutralMass float64
	PrecursorExpMass     float64
	RetentionTime        float64
	CalcNeutralPepMass   float64
	RawMassDiff          float64
	Massdiff             float64
	LocalizedPTMSites    map[string]int
	LocalizedPTMMassDiff map[string]string
	Probability          float64
	IsoMassD             int
	Expectation          float64
	Xcorr                float64
	DeltaCN              float64
	DeltaCNStar          float64
	SPScore              float64
	SPRank               float64
	Hyperscore           float64
	Nextscore            float64
	DiscriminantValue    float64
	Intensity            float64
	IsRejected           uint8
	Modifications        mod.Modifications
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
		p.FileName = f
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
					IsobaricMods:     make(map[string]uint8),
				}

				p.Modifications.Index[key] = m
			}
		}

		// map terminal modifications from file
		for _, i := range mpa.MsmsRunSummary.SearchSummary.TerminalModifications {

			key := fmt.Sprintf("%s-term#%.4f", i.Terminus, i.Mass)

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
					IsobaricMods:      make(map[string]uint8),
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

		// start processing spectra queries
		var psmlist PepIDList
		sq := mpa.MsmsRunSummary.SpectrumQuery
		for _, i := range sq {
			psm := processSpectrumQuery(i, p.Modifications, p.DecoyTag)
			psmlist = append(psmlist, psm)
		}

		if len(psmlist) == 0 {
			err.NoPSMFound()
		}

		p.PeptideIdentification = psmlist
		p.Prophet = string(mpa.AnalysisSummary[0].Analysis)
		p.Models = models

		p.adjustMassDeviation()

		if len(psmlist) == 0 {
			err.NoPSMFound()
		}

	}

	return
}

func processSpectrumQuery(sq spc.SpectrumQuery, mods mod.Modifications, decoyTag string) PeptideIdentification {

	var psm PeptideIdentification
	psm.Modifications.Index = make(map[string]mod.Modification)

	psm.Index = sq.Index
	psm.Spectrum = string(sq.Spectrum)
	psm.Scan = sq.StartScan
	psm.PrecursorNeutralMass = sq.PrecursorNeutralMass
	psm.AssumedCharge = sq.AssumedCharge
	psm.RetentionTime = sq.RetentionTimeSec

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
		psm.Massdiff = i.Massdiff

		for _, j := range i.AnalysisResult {
			if string(j.Analysis) == "peptideprophet" {
				psm.Probability = j.PeptideProphetResult.Probability
				for _, k := range j.PeptideProphetResult.SearchScoreSummary.Parameter {
					if k.Name == "isomassd" {
						psm.IsoMassD, _ = strconv.Atoi(k.Value)
					}
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
		}

		for _, j := range i.Score {
			if string(j.Name) == "expect" {
				psm.Expectation = j.Value
			} else if string(j.Name) == "xcorr" {
				psm.Xcorr = j.Value
			} else if string(j.Name) == "deltacn" {
				psm.DeltaCN = j.Value
			} else if string(j.Name) == "deltacnstar" {
				psm.DeltaCNStar = j.Value
			} else if string(j.Name) == "spscore" {
				psm.SPScore = j.Value
			} else if string(j.Name) == "sprank" {
				psm.SPRank = j.Value
			} else if string(j.Name) == "hyperscore" {
				psm.Hyperscore = j.Value
			} else if string(j.Name) == "nextscore" {
				psm.Nextscore = j.Value
			}
		}

		psm.mapModsFromPepXML(i.ModificationInfo, mods)
	}

	return psm
}

// mapModsFromPepXML receives a pepXML struct with modifications and adds them to
// the given struct
func (p *PeptideIdentification) mapModsFromPepXML(m spc.ModificationInfo, mods mod.Modifications) {

	p.ModifiedPeptide = string(m.ModifiedPeptide)

	var isotopicCorr float64
	if p.IsoMassD != 0 {
		isotopicCorr = p.Massdiff - (bio.Proton * float64(p.IsoMassD))
	} else {
		isotopicCorr = p.Massdiff
	}

	for _, i := range m.ModAminoacidMass {
		aa := strings.Split(p.Peptide, "")
		key := fmt.Sprintf("%s#%.4f", aa[i.Position-1], i.Mass)
		v, ok := mods.Index[key]
		if ok {
			m := v
			newKey := fmt.Sprintf("%s#%d#%.4f", aa[i.Position-1], i.Position, i.Mass)
			m.Index = newKey
			m.Position = strconv.Itoa(i.Position)
			m.IsobaricMods = make(map[string]uint8)
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
			m.IsobaricMods = make(map[string]uint8)
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
			m.IsobaricMods = make(map[string]uint8)
			p.Modifications.Index[key] = m
		}
	}

	if isotopicCorr >= 0.036386 || isotopicCorr <= -0.036386 {
		key := fmt.Sprintf("%.4f", isotopicCorr)
		_, ok := p.Modifications.Index[key]
		if !ok {
			m := mod.Modification{
				Index:        key,
				Name:         "Unknown",
				Type:         "Observed",
				MassDiff:     isotopicCorr,
				IsobaricMods: make(map[string]uint8),
			}
			p.Modifications.Index[key] = m
		}

	}

	return
}

// adjustMassDeviation calculates the mass deviation for a pepXML file based on the 0 mass difference
func (p *PepXML) adjustMassDeviation() {

	var countZero int
	var massZero float64
	var adjustedMass float64

	for _, i := range p.PeptideIdentification {
		if math.Abs(i.Massdiff) >= -0.1 && math.Abs(i.Massdiff) <= 0.1 {
			countZero++
			massZero += i.Massdiff
		}
	}

	adjustedMass = massZero / float64(countZero)

	// keep the original massdiff on the raw variable just in case
	for i := range p.PeptideIdentification {
		p.PeptideIdentification[i].RawMassDiff = p.PeptideIdentification[i].Massdiff
		p.PeptideIdentification[i].Massdiff = (p.PeptideIdentification[i].Massdiff - adjustedMass)
	}

	return
}

// PromoteProteinIDs promotes protein identifications where the reference protein
// is indistinguishable to other target proteins.
func (p *PepXML) PromoteProteinIDs() {

	for i := range p.PeptideIdentification {

		var list []string
		var ref string

		if strings.Contains(p.PeptideIdentification[i].Protein, p.DecoyTag) {
			for j := range p.PeptideIdentification[i].AlternativeProteins {
				if !strings.HasPrefix(p.PeptideIdentification[i].AlternativeProteins[j], p.DecoyTag) {
					list = append(list, p.PeptideIdentification[i].AlternativeProteins[j])
				}
			}
		}

		if len(list) > 0 {
			for i := range list {
				if strings.HasPrefix(list[i], "sp|") {
					ref = list[i]
					break
				} else {
					ref = list[i]
				}
			}
			p.PeptideIdentification[i].Protein = ref
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

		err.Plotter(e)

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
	var class = true

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
		err.MarshalFile(e)
	}

	e = ioutil.WriteFile(sys.PepxmlBin(), b, sys.FilePermission())
	if e != nil {
		err.WriteFile(e)
	}

	return
}

// Restore reads philosopher results files and restore the data sctructure
func (p *PepXML) Restore() {

	b, e := ioutil.ReadFile(sys.PepxmlBin())
	if e != nil {
		err.ReadFile(e)
	}

	e = msgpack.Unmarshal(b, &p)
	if e != nil {
		err.DecodeMsgPck(e)
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
		err.WarnCustom(errors.New("Cannot determine binary data class"))
	}

	b, e := msgpack.Marshal(&p)
	if e != nil {
		err.MarshalFile(e)
	}

	e = ioutil.WriteFile(dest, b, sys.FilePermission())
	if e != nil {
		err.WriteFile(e)
	}

	return
}

// Restore reads philosopher results files and restore the data sctructure
func (p *PepIDList) Restore(level string) error {

	var dest string

	if level == "psm" {
		dest = sys.PsmBin()
	} else if level == "pep" {
		dest = sys.PepBin()
	} else if level == "ion" {
		dest = sys.IonBin()
	} else {
		err.WarnCustom(errors.New("Cannot determine binary data class"))
	}

	b, e := ioutil.ReadFile(dest)
	if e != nil {
		err.ReadFile(e)
	}

	e = msgpack.Unmarshal(b, &p)
	if e != nil {
		err.DecodeMsgPck(e)
	}

	return nil
}
