package id

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"path/filepath"
	"strings"

	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/mod"
	"github.com/prvst/philosopher/lib/spc"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/prvst/philosopher/lib/uti"
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
	Database              string
	Prophet               string
	DefinedModMassDiff    map[float64]float64
	DefinedModAminoAcid   map[float64]string
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
	ModPositions         []string
	AssignedModMasses    []float64
	AssignedMassDiffs    []float64
	AssignedAminoAcid    []string
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
func (p *PepXML) Read(f string) error {

	var xml spc.PepXML
	e := xml.Parse(f)
	if e != nil {
		return e
	}

	var mpa = xml.MsmsPipelineAnalysis

	if len(mpa.AnalysisSummary) > 0 {
		p.FileName = f
		p.Database = string(mpa.MsmsRunSummary.SearchSummary.SearchDatabase.LocalPath)
		p.SpectraFile = fmt.Sprintf("%s%s", mpa.MsmsRunSummary.BaseName, mpa.MsmsRunSummary.RawData)

		var models []spc.DistributionPoint
		pps := mpa.AnalysisSummary[0].PeptideprophetSummary

		// collect distribution points from meta
		for _, i := range pps.DistributionPoint {
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

		// collect variable and fixed modifications
		if len(p.DefinedModMassDiff) == 0 {
			p.DefinedModMassDiff = make(map[float64]float64)
			p.DefinedModAminoAcid = make(map[float64]string)
		}

		// get the search engine
		p.SearchEngine = string(mpa.MsmsRunSummary.SearchSummary.SearchEngine)
		if strings.Contains(string(mpa.MsmsRunSummary.SearchSummary.SearchEngineVersion), "MSFragger") {
			p.SearchEngine = "MSFragger"
		}

		// internal modifications (variable only)
		for _, i := range mpa.MsmsRunSummary.SearchSummary.AminoAcidModifications {
			p.DefinedModMassDiff[uti.Round(i.Mass, 5, 2)] = i.MassDiff
			p.DefinedModAminoAcid[uti.Round(i.Mass, 5, 2)] = string(i.AminoAcid)
		}

		// termini modifications
		for _, i := range mpa.MsmsRunSummary.SearchSummary.TerminalModifications {
			p.DefinedModMassDiff[uti.Round(i.Mass, 5, 2)] = i.MassDiff
			if string(i.Terminus) == "N" {
				p.DefinedModAminoAcid[uti.Round(i.Mass, 5, 2)] = "n"
			} else if string(i.Terminus) == "C" {
				p.DefinedModAminoAcid[uti.Round(i.Mass, 5, 2)] = "c"
			}
		}

		// start processing spectra queries
		var psmlist PepIDList
		sq := mpa.MsmsRunSummary.SpectrumQuery
		for _, i := range sq {
			psm := processSpectrumQuery(i, p.DefinedModMassDiff, p.DefinedModAminoAcid, p.DecoyTag)
			psmlist = append(psmlist, psm)
		}

		if len(psmlist) == 0 {
			return &err.Error{Type: err.NoPSMFound, Class: err.FATA}
		}

		p.PeptideIdentification = psmlist
		p.Prophet = string(mpa.AnalysisSummary[0].Analysis)
		p.Models = models

		p.adjustMassDeviation()

		if len(psmlist) == 0 {
			logrus.Error("No PSM detected, check your files and try agains")
		}

	}

	return nil
}

func processSpectrumQuery(sq spc.SpectrumQuery, definedModMassDiff map[float64]float64, definedModAminoAcid map[float64]string, decoyTag string) PeptideIdentification {

	var psm PeptideIdentification
	psm.Modifications.MassDiffIndex = make(map[float64]float64)
	psm.Modifications.AminoAcidIndex = make(map[float64]string)

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

		if len(string(i.ModificationInfo.ModifiedPeptide)) > 0 {

			psm.ModifiedPeptide = string(i.ModificationInfo.ModifiedPeptide)

			for _, j := range i.ModificationInfo.ModAminoacidMass {
				pos := fmt.Sprintf("%d", j.Position)
				psm.ModPositions = append(psm.ModPositions, pos)
				psm.AssignedModMasses = append(psm.AssignedModMasses, j.Mass)
				psm.AssignedMassDiffs = append(psm.AssignedMassDiffs, definedModMassDiff[uti.Round(j.Mass, 5, 2)])
				psm.AssignedAminoAcid = append(psm.AssignedAminoAcid, definedModAminoAcid[uti.Round(j.Mass, 5, 2)])
			}

			// n-temrinal modifications
			if i.ModificationInfo.ModNTermMass != 0 {
				psm.ModPositions = append(psm.ModPositions, "n")
				psm.AssignedModMasses = append(psm.AssignedModMasses, i.ModificationInfo.ModNTermMass)
				psm.AssignedMassDiffs = append(psm.AssignedMassDiffs, definedModMassDiff[uti.Round(i.ModificationInfo.ModNTermMass, 5, 2)])
				psm.AssignedAminoAcid = append(psm.AssignedAminoAcid, definedModAminoAcid[uti.Round(i.ModificationInfo.ModNTermMass, 5, 2)])
			}

			// c-terminal modifications
			if i.ModificationInfo.ModCTermMass != 0 {
				psm.ModPositions = append(psm.ModPositions, "c")
				psm.AssignedModMasses = append(psm.AssignedModMasses, i.ModificationInfo.ModCTermMass)
				psm.AssignedMassDiffs = append(psm.AssignedMassDiffs, definedModMassDiff[uti.Round(i.ModificationInfo.ModCTermMass, 5, 2)])
				psm.AssignedAminoAcid = append(psm.AssignedAminoAcid, definedModAminoAcid[uti.Round(i.ModificationInfo.ModCTermMass, 5, 2)])
			}

		}

	}

	return psm
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
				if !strings.Contains(p.PeptideIdentification[i].AlternativeProteins[j], p.DecoyTag) {
					list = append(list, p.PeptideIdentification[i].AlternativeProteins[j])
				}
			}
		}

		if len(list) > 0 {
			for i := range list {
				if strings.Contains(list[i], "sp|") {
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
func (p *PepXML) ReportModels(session, name string) (err error) {

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

	return err
}

func printModel(v, path string, xAxis, obs, pos, neg []float64) error {

	p, e := plot.New()
	if e != nil {
		return &err.Error{Type: err.CannotInstantiateStruct, Class: err.FATA, Argument: "plotter"}
	}

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

	return nil
}

// tdclassifier identifies a PSM as target or Decoy based on the
// presence of the TAG string on <protein> and <alternative_proteins>
func tdclassifier(p PeptideIdentification, tag string) bool {

	// default for TRUE ( DECOY)
	var class = true

	if strings.Contains(string(p.Protein), tag) {
		class = true
	} else {
		class = false
	}

	for i := range p.AlternativeProteins {
		if !strings.Contains(p.AlternativeProteins[i], tag) {
			class = false
		}
		break
	}

	return class
}

// Serialize converts the whle structure to a gob file
func (p *PepXML) Serialize() *err.Error {

	b, e := msgpack.Marshal(&p)
	if e != nil {
		return &err.Error{Type: err.CannotCreateOutputFile, Class: err.FATA, Argument: e.Error()}
	}

	e = ioutil.WriteFile(sys.PepxmlBin(), b, 0644)
	if e != nil {
		return &err.Error{Type: err.CannotSerializeData, Class: err.FATA, Argument: e.Error()}
	}

	return nil
}

// Restore reads philosopher results files and restore the data sctructure
func (p *PepXML) Restore() error {

	b, e := ioutil.ReadFile(sys.PepxmlBin())
	if e != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: ": Could not restore Philosopher result"}
	}

	e = msgpack.Unmarshal(b, &p)
	if e != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: ": Could not restore Philosopher result"}
	}

	return nil
}

// Serialize converts the whle structure to a gob file
func (p *PepIDList) Serialize(level string) *err.Error {

	var dest string

	if level == "psm" {
		dest = sys.PsmBin()
	} else if level == "pep" {
		dest = sys.PepBin()
	} else if level == "ion" {
		dest = sys.IonBin()
	} else {
		return &err.Error{Type: err.CannotSerializeData, Class: err.FATA}
	}

	b, er := msgpack.Marshal(&p)
	if er != nil {
		return &err.Error{Type: err.CannotCreateOutputFile, Class: err.FATA, Argument: er.Error()}
	}

	er = ioutil.WriteFile(dest, b, 0644)
	if er != nil {
		return &err.Error{Type: err.CannotSerializeData, Class: err.FATA, Argument: er.Error()}
	}

	return nil
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
		return errors.New("Cannot seralize Peptide identifications")
	}

	b, e := ioutil.ReadFile(dest)
	if e != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: ": Cannot seralize Peptide identifications"}
	}

	e = msgpack.Unmarshal(b, &p)
	if e != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: ": Cannot seralize Peptide identifications"}
	}

	return nil
}
