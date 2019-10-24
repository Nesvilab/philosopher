package rep

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nesvilab/philosopher/lib/id"
	"github.com/nesvilab/philosopher/lib/met"
	"github.com/nesvilab/philosopher/lib/mod"
	"github.com/nesvilab/philosopher/lib/msg"
	"github.com/nesvilab/philosopher/lib/tmt"
	"github.com/sirupsen/logrus"
)

// Evidence ...
type Evidence struct {
	Decoys          bool
	Parameters      SearchParametersEvidence
	PSM             PSMEvidenceList
	Ions            IonEvidenceList
	Peptides        PeptideEvidenceList
	Proteins        ProteinEvidenceList
	Mods            mod.Modifications
	Modifications   ModificationEvidence
	CombinedProtein CombinedProteinEvidenceList
	CombinedPeptide CombinedPeptideEvidenceList
}

// SearchParametersEvidence ...
type SearchParametersEvidence struct {
	MSFragger                          string
	DatabaseName                       string
	NumThreads                         string
	PrecursorMassLower                 string
	PrecursorMassUpper                 string
	PrecursorMassUnits                 string
	PrecursorTrueTolerance             string
	PrecursorTrueUnits                 string
	FragmentMassTolerance              string
	FragmentMassUnits                  string
	CalibrateMass                      string
	Ms1ToleranceMad                    string
	Ms2ToleranceMad                    string
	EvaluateMassCalibration            string
	IsotopeError                       string
	MassOffsets                        string
	PrecursorMassMode                  string
	ShiftedIons                        string
	ShiftedIonsExcludeRanges           string
	FragmentIonSeries                  string
	SearchEnzymeName                   string
	SearchEnzymeCutafter               string
	SearchEnzymeButnotafter            string
	NumEnzymeTermini                   string
	AllowedMissedCleavage              string
	ClipNTermM                         string
	AllowMultipleVariableModsOnResidue string
	MaxVariableModsPerMod              string
	MaxVariableModsCombinations        string
	OutputFileExtension                string
	OutputFormat                       string
	OutputReportTopN                   string
	OutputMaxExpect                    string
	ReportAlternativeProteins          string
	OverrideCharge                     string
	PrecursorCharge                    string
	DigestMinLength                    string
	DigestMaxLength                    string
	DigestMassRange                    string
	MaxFragmentCharge                  string
	TrackZeroTopN                      string
	ZeroBinAcceptExpect                string
	ZeroBinMultExpect                  string
	AddTopNComplementary               string
	MinimumPeaks                       string
	UseTopNPeaks                       string
	MinFragmentsModelling              string
	MinMatchedFragments                string
	MinimumRatio                       string
	ClearMzRange                       string
	VariableMod01                      string
	VariableMod02                      string
	Alanine                            string
	Cysteine                           string
	CTermPeptide                       string
	CTermProtein                       string
	AsparticAcid                       string
	GlutamicAcid                       string
	Phenylalanine                      string
	Glycine                            string
	Histidine                          string
	Isoleucine                         string
	Lysine                             string
	Leucine                            string
	Methionine                         string
	Asparagine                         string
	NTermPeptide                       string
	NTermProtein                       string
	Proline                            string
	Glutamine                          string
	Arginine                           string
	Serine                             string
	Threonine                          string
	Valine                             string
	Tryptophan                         string
	Tyrosine                           string
}

// PSMEvidence struct
type PSMEvidence struct {
	Source                           string
	Index                            uint32
	Spectrum                         string
	Scan                             int
	NTT                              int
	NMC                              int
	PrevAA                           string
	NextAA                           string
	Peptide                          string
	IonForm                          string
	Protein                          string
	ProteinDescription               string
	ProteinID                        string
	EntryName                        string
	GeneName                         string
	ModifiedPeptide                  string
	MappedProteins                   map[string]int
	MappedGenes                      map[string]int
	AssumedCharge                    uint8
	HitRank                          uint8
	UncalibratedPrecursorNeutralMass float64
	PrecursorNeutralMass             float64
	PrecursorExpMass                 float64
	RetentionTime                    float64
	CalcNeutralPepMass               float64
	RawMassdiff                      float64
	Massdiff                         float64
	LocalizedPTMSites                map[string]int
	LocalizedPTMMassDiff             map[string]string
	Probability                      float64
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
	Purity                           float64
	IsDecoy                          bool
	IsUnique                         bool
	IsURazor                         bool
	Labels                           tmt.Labels
	Modifications                    mod.Modifications
}

// PSMEvidenceList ...
type PSMEvidenceList []PSMEvidence

func (a PSMEvidenceList) Len() int           { return len(a) }
func (a PSMEvidenceList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a PSMEvidenceList) Less(i, j int) bool { return a[i].Spectrum < a[j].Spectrum }

// IonEvidence groups all valid info about peptide ions for reports
type IonEvidence struct {
	Sequence             string
	IonForm              string
	ModifiedSequence     string
	RetentionTime        string
	ChargeState          uint8
	Spectra              map[string]int
	MappedProteins       map[string]int
	MappedGenes          map[string]int
	MZ                   float64
	PeptideMass          float64
	PrecursorNeutralMass float64
	Weight               float64
	GroupWeight          float64
	Intensity            float64
	Probability          float64
	Expectation          float64
	SummedLabelIntensity float64
	IsUnique             bool
	IsURazor             bool
	IsDecoy              bool
	Protein              string
	ProteinID            string
	GeneName             string
	EntryName            string
	ProteinDescription   string
	Labels               tmt.Labels
	PhosphoLabels        tmt.Labels
	Modifications        mod.Modifications
}

// IonEvidenceList ...
type IonEvidenceList []IonEvidence

func (a IonEvidenceList) Len() int           { return len(a) }
func (a IonEvidenceList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a IonEvidenceList) Less(i, j int) bool { return a[i].Sequence < a[j].Sequence }

// PeptideEvidence groups all valid info about peptide ions for reports
type PeptideEvidence struct {
	Sequence               string
	ChargeState            map[uint8]uint8
	Spectra                map[string]uint8
	Protein                string
	ProteinID              string
	GeneName               string
	EntryName              string
	ProteinDescription     string
	MappedProteins         map[string]int
	MappedGenes            map[string]int
	Spc                    int
	Intensity              float64
	Probability            float64
	ModifiedObservations   int
	UnModifiedObservations int
	IsDecoy                bool
	Labels                 tmt.Labels
	PhosphoLabels          tmt.Labels
	Modifications          mod.Modifications
}

// PeptideEvidenceList ...
type PeptideEvidenceList []PeptideEvidence

func (a PeptideEvidenceList) Len() int           { return len(a) }
func (a PeptideEvidenceList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a PeptideEvidenceList) Less(i, j int) bool { return a[i].Sequence < a[j].Sequence }

// ProteinEvidence ...
type ProteinEvidence struct {
	OriginalHeader         string
	PartHeader             string
	ProteinName            string
	ProteinGroup           uint32
	ProteinSubGroup        string
	ProteinID              string
	EntryName              string
	Description            string
	Organism               string
	Length                 int
	Coverage               float32
	GeneNames              string
	ProteinExistence       string
	Sequence               string
	SupportingSpectra      map[string]int
	IndiProtein            map[string]uint8
	UniqueStrippedPeptides int
	TotalPeptideIons       map[string]IonEvidence
	TotalSpC               int
	UniqueSpC              int
	URazorSpC              int // Unique + razor
	TotalIntensity         float64
	UniqueIntensity        float64
	URazorIntensity        float64 // Unique + razor
	Probability            float64
	TopPepProb             float64
	IsDecoy                bool
	IsContaminant          bool
	TotalLabels            tmt.Labels
	UniqueLabels           tmt.Labels
	URazorLabels           tmt.Labels // Unique + razor
	PhosphoTotalLabels     tmt.Labels
	PhosphoUniqueLabels    tmt.Labels
	PhosphoURazorLabels    tmt.Labels // Unique + razor
	Modifications          mod.Modifications
}

// ProteinEvidenceList list
type ProteinEvidenceList []ProteinEvidence

func (a ProteinEvidenceList) Len() int           { return len(a) }
func (a ProteinEvidenceList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ProteinEvidenceList) Less(i, j int) bool { return a[i].ProteinGroup < a[j].ProteinGroup }

// CombinedProteinEvidence represents all combined proteins detected
type CombinedProteinEvidence struct {
	GroupNumber            uint32
	SiblingID              string
	ProteinName            string
	ProteinID              string
	IndiProtein            []string
	EntryName              string
	Organism               string
	Length                 int
	Coverage               float32
	GeneNames              string
	ProteinExistence       string
	Description            string
	Names                  []string
	UniqueStrippedPeptides int
	SupportingSpectra      map[string]string
	ProteinProbability     float64
	TopPepProb             float64
	PeptideIons            []id.PeptideIonIdentification
	TotalSpc               map[string]int
	UniqueSpc              map[string]int
	UrazorSpc              map[string]int
	TotalIntensity         map[string]float64
	UniqueIntensity        map[string]float64
	UrazorIntensity        map[string]float64
	TotalLabels            map[string]tmt.Labels
	UniqueLabels           map[string]tmt.Labels
	URazorLabels           map[string]tmt.Labels // Unique + razor
}

// CombinedProteinEvidenceList is a list of Combined Protein Evidences
type CombinedProteinEvidenceList []CombinedProteinEvidence

func (a CombinedProteinEvidenceList) Len() int           { return len(a) }
func (a CombinedProteinEvidenceList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CombinedProteinEvidenceList) Less(i, j int) bool { return a[i].GroupNumber < a[j].GroupNumber }

// CombinedPeptideEvidence represents all combined peptides detected
type CombinedPeptideEvidence struct {
	Key                string
	BestPSM            float64
	Sequence           string
	Protein            string
	Gene               string
	ProteinDescription string
	ChargeStates       []string
	AssignedMassDiffs  []string
	Spc                map[string]int
	Intensity          map[string]float64
}

// CombinedPeptideEvidenceList is a list of Combined Peptide Evidences
type CombinedPeptideEvidenceList []CombinedPeptideEvidence

func (a CombinedPeptideEvidenceList) Len() int           { return len(a) }
func (a CombinedPeptideEvidenceList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CombinedPeptideEvidenceList) Less(i, j int) bool { return a[i].Sequence < a[j].Sequence }

// ModificationEvidence represents the list of modifications and the mod bins
type ModificationEvidence struct {
	MassBins []MassBin
}

// MassBin represents each bin from the mass distribution
type MassBin struct {
	LowerMass     float64
	HigherRight   float64
	MassCenter    float64
	AverageMass   float64
	CorrectedMass float64
	Modifications []string
	AssignedMods  PSMEvidenceList
	ObservedMods  PSMEvidenceList
}

// New constructor
func New() Evidence {

	var self Evidence

	return self
}

// Run is the main entry poit for Report
func Run(m met.Data) {

	var repo = New()
	repo.RestoreGranular()

	var isComet bool
	var hasLoc bool
	var isoBrand string
	var isoChannels int
	var labels = make(map[string]string)

	if len(m.Comet.Param) > 0 {
		isComet = true
	}

	if m.PTMProphet.InputFiles != nil || len(m.PTMProphet.InputFiles) > 0 {
		hasLoc = true
	}

	if m.Quantify.Brand == "tmt" {
		isoBrand = "tmt"
	} else if m.Quantify.Brand == "itraq" {
		isoBrand = "itraq"
	}

	if len(m.Quantify.Plex) > 0 {
		isoChannels, _ = strconv.Atoi(m.Quantify.Plex)
	}

	// get the labels from the annotation file
	if len(m.Quantify.Annot) > 0 {
		if len(m.Quantify.Annot) > 0 {
			annotfile := fmt.Sprintf(".%sannotation.txt", string(filepath.Separator))
			annotfile, _ = filepath.Abs(annotfile)
			labels = getLabelNames(annotfile)
		}
	}

	logrus.Info("Creating reports")

	// PSM
	repo.MetaPSMReport(labels, isoBrand, isoChannels, m.Report.Decoys, isComet, hasLoc)

	// Ion
	repo.MetaIonReport(labels, isoBrand, isoChannels, m.Report.Decoys)

	// Peptide
	repo.MetaPeptideReport(labels, isoBrand, isoChannels, m.Report.Decoys)

	// Protein
	if len(m.Filter.Pox) > 0 {
		repo.MetaProteinReport(labels, isoBrand, isoChannels, m.Report.Decoys, m.Quantify.Unique)
	}

	// Modifications
	if len(repo.Modifications.MassBins) > 0 {
		repo.ModificationReport()

		if m.PTMProphet.InputFiles != nil || len(m.PTMProphet.InputFiles) > 0 {
			repo.PSMLocalizationReport(m.Filter.Tag, m.Filter.Razor, m.Report.Decoys)
		}

		repo.PlotMassHist()
	}

	// MSstats
	if m.Report.MSstats == true {
		repo.MetaMSstatsReport(labels, isoBrand, isoChannels, m.Report.Decoys)
	}

	// MzID
	if m.Report.MZID == true {
		repo.MzIdentMLReport(m.Version, m.Database.Annot)
	}

	return
}

// addCustomNames adds to the label structures user-defined names to be used on the TMT labels
func getLabelNames(annot string) map[string]string {

	var labels = make(map[string]string)

	file, e := os.Open(annot)
	if e != nil {
		msg.ReadFile(e, "fatal")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		names := strings.Split(scanner.Text(), " ")
		labels[names[0]] = names[1]
	}

	if e = scanner.Err(); e != nil {
		msg.Custom(errors.New("Annotation file seems to be empty"), "error")
	}

	return labels
}

// prepares the list of modifications to be printed by the report functions
func getModsList(m map[string]mod.Modification) ([]string, []string) {

	var a []string
	var o []string

	for _, i := range m {
		if i.Type == "Assigned" && i.Name != "Unknown" {
			a = append(a, fmt.Sprintf("%s%s(%.4f)", i.Position, i.AminoAcid, i.MassDiff))
		}
		if i.Type == "Observed" && i.Name != "Unknown" {
			for k, v := range i.IsobaricMods {
				o = append(o, fmt.Sprintf("%s(%f)", k, v))
			}
		}
	}

	return a, o
}
