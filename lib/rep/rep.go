package rep

import (
	"fmt"
	"strconv"

	"philosopher/lib/id"
	"philosopher/lib/iso"
	"philosopher/lib/met"
	"philosopher/lib/mod"

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
	UseAllModsInFirstSearch            string
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
	CheckSpectralFiles                 string
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

func (p PSMEvidence) SpectrumFileName() id.SpectrumType {
	return id.SpectrumType{Spectrum: p.Spectrum, FileName: p.SpectrumFile}
}

// PSMEvidence struct
type PSMEvidence struct {
	Source                           string
	Spectrum                         string
	SpectrumFile                     string
	Peptide                          string
	Protein                          string
	ProteinDescription               string
	ProteinID                        string
	EntryName                        string
	GeneName                         string
	ModifiedPeptide                  string
	CompensationVoltage              string
	ProteinStart                     int
	ProteinEnd                       int
	NumberOfEnzymaticTermini         uint8
	NumberOfMissedCleavages          uint8
	AssumedCharge                    uint8
	HitRank                          uint8
	Index                            uint32
	UncalibratedPrecursorNeutralMass float64
	PrecursorNeutralMass             float64
	RetentionTime                    float64
	CalcNeutralPepMass               float64
	RawMassdiff                      float64
	Massdiff                         float64
	Probability                      float64
	Expectation                      float64
	Xcorr                            float64
	DeltaCN                          float64
	DeltaCNStar                      float64
	SPScore                          float64
	SPRank                           float64
	Hyperscore                       float64
	Nextscore                        float64
	SpectralSim                      float64
	Rtscore                          float64
	Intensity                        float64
	IonMobility                      float64
	Purity                           float64
	PrevAA                           byte
	NextAA                           byte
	IsDecoy                          bool
	IsUnique                         bool
	IsURazor                         bool
	PTM                              *id.PTM
	MSFraggerLoc                     *id.MSFraggerLoc
	Labels                           *iso.Labels
	Modifications                    mod.ModificationsSlice
	MappedProteins                   map[string]int
	MappedGenes                      map[string]struct{}
}

func (e PSMEvidence) IonForm() id.IonFormType {
	t, err := strconv.ParseFloat(fmt.Sprintf("%.4f", e.CalcNeutralPepMass), 32)
	if err != nil {
		panic(err)
	}
	return id.IonFormType{e.Peptide, float32(t), e.AssumedCharge}
}

func (e IonEvidence) IonForm() id.IonFormType {
	t, err := strconv.ParseFloat(fmt.Sprintf("%.4f", e.PeptideMass), 32)
	if err != nil {
		panic(err)
	}
	return id.IonFormType{e.Sequence, float32(t), e.ChargeState}
}

// PSMEvidenceList ...
type PSMEvidenceList []PSMEvidence

func (a PSMEvidenceList) Len() int      { return len(a) }
func (a PSMEvidenceList) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a PSMEvidenceList) Less(i, j int) bool {
	return a[i].SpectrumFileName().Str() < a[j].SpectrumFileName().Str()
}

// RemovePSMByIndex perfomrs a re-slicing by removing an element from a list
func RemovePSMByIndex(s []PSMEvidence, i int) []PSMEvidence {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// IonEvidence groups all valid info about peptide ions for reports
type IonEvidence struct {
	Sequence                 string
	ModifiedSequence         string
	Protein                  string
	ProteinID                string
	GeneName                 string
	EntryName                string
	ProteinDescription       string
	ChargeState              uint8
	NumberOfEnzymaticTermini uint8
	PrevAA                   byte
	NextAA                   byte
	MZ                       float64
	PeptideMass              float64
	PrecursorNeutralMass     float64
	Weight                   float64
	GroupWeight              float64
	Intensity                float64
	Probability              float64
	Expectation              float64
	SummedLabelIntensity     float64
	IsUnique                 bool
	IsURazor                 bool
	IsDecoy                  bool
	Labels                   *iso.Labels
	PhosphoLabels            *iso.Labels
	Modifications            mod.ModificationsSlice
	Spectra                  map[id.SpectrumType]int
	MappedProteins           map[string]int
	MappedGenes              map[string]struct{}
	//IonForm()              string
	//RetentionTime          string

}

// IonEvidenceList ...
type IonEvidenceList []IonEvidence

func (a IonEvidenceList) Len() int           { return len(a) }
func (a IonEvidenceList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a IonEvidenceList) Less(i, j int) bool { return a[i].Sequence < a[j].Sequence }

// RemoveIonsByIndex perfomrs a re-slicing by removing an element from a list
func RemoveIonsByIndex(s []IonEvidence, i int) []IonEvidence {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// PeptideEvidence groups all valid info about peptide ions for reports
type PeptideEvidence struct {
	Sequence               string
	Protein                string
	ProteinID              string
	GeneName               string
	EntryName              string
	ProteinDescription     string
	Spc                    int
	ModifiedObservations   int
	UnModifiedObservations int
	Intensity              float64
	Probability            float64
	PrevAA                 byte
	NextAA                 byte
	IsUnique               bool
	IsURazor               bool
	IsDecoy                bool
	ChargeState            map[uint8]uint8
	Spectra                map[id.SpectrumType]uint8
	MappedProteins         map[string]int
	MappedGenes            map[string]struct{}
	Labels                 *iso.Labels
	PhosphoLabels          *iso.Labels
	Modifications          mod.ModificationsSlice
}

// PeptideEvidenceList ...
type PeptideEvidenceList []PeptideEvidence

func (a PeptideEvidenceList) Len() int           { return len(a) }
func (a PeptideEvidenceList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a PeptideEvidenceList) Less(i, j int) bool { return a[i].Sequence < a[j].Sequence }

// RemovePeptidesByIndex perfomrs a re-slicing by removing an element from a list
func RemovePeptidesByIndex(s []PeptideEvidence, i int) []PeptideEvidence {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// ProteinEvidence ...
type ProteinEvidence struct {
	OriginalHeader         string
	PartHeader             string
	ProteinName            string
	ProteinSubGroup        string
	ProteinID              string
	EntryName              string
	Description            string
	Organism               string
	GeneNames              string
	ProteinExistence       string
	Sequence               string
	Length                 int
	TotalSpC               int
	UniqueSpC              int
	URazorSpC              int // Unique + razor
	UniqueStrippedPeptides int
	ProteinGroup           uint32
	Coverage               float32
	TotalIntensity         float64
	UniqueIntensity        float64
	URazorIntensity        float64 // Unique + razor
	Probability            float64
	TopPepProb             float64
	IsDecoy                bool
	IsContaminant          bool
	SupportingSpectra      map[id.SpectrumType]int
	IndiProtein            map[string]struct{}
	TotalPeptideIons       map[id.IonFormType]IonEvidence
	TotalPeptides          map[string]int
	UniquePeptides         map[string]int
	URazorPeptides         map[string]int // Unique + razor
	TotalLabels            *iso.Labels
	UniqueLabels           *iso.Labels
	URazorLabels           *iso.Labels // Unique + razor
	PhosphoTotalLabels     *iso.Labels
	PhosphoUniqueLabels    *iso.Labels
	PhosphoURazorLabels    *iso.Labels // Unique + razor
	Modifications          mod.ModificationsSlice
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
	EntryName              string
	Organism               string
	GeneNames              string
	ProteinExistence       string
	Description            string
	IndiProtein            []string
	Names                  []string
	Length                 int
	UniqueStrippedPeptides int
	Coverage               float32
	ProteinProbability     float64
	TopPepProb             float64
	SupportingSpectra      map[string]string
	TotalSpc               map[string]int
	UniqueSpc              map[string]int
	UrazorSpc              map[string]int
	TotalPeptides          map[string]map[string]bool
	UniquePeptides         map[string]map[string]bool
	UrazorPeptides         map[string]map[string]bool
	TotalIntensity         map[string]float64
	UniqueIntensity        map[string]float64
	UrazorIntensity        map[string]float64
	TotalLabels            map[string]iso.Labels
	UniqueLabels           map[string]iso.Labels
	URazorLabels           map[string]iso.Labels // Unique + razor
	PeptideIons            []id.PeptideIonIdentification
}

// CombinedProteinEvidenceList is a list of Combined Protein Evidences
type CombinedProteinEvidenceList []CombinedProteinEvidence

func (a CombinedProteinEvidenceList) Len() int           { return len(a) }
func (a CombinedProteinEvidenceList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CombinedProteinEvidenceList) Less(i, j int) bool { return a[i].GroupNumber < a[j].GroupNumber }

// CombinedPeptideEvidence represents all combined peptides detected
type CombinedPeptideEvidence struct {
	Sequence           string
	Protein            string
	ProteinID          string
	EntryName          string
	Gene               string
	ProteinDescription string
	BestPSM            float64
	ChargeStates       map[uint8]uint8
	AssignedMassDiffs  map[string]uint8
	Spc                map[string]int
	Intensity          map[string]float64
}

// CombinedPeptideEvidenceList is a list of Combined Peptide Evidences
type CombinedPeptideEvidenceList []CombinedPeptideEvidence

func (a CombinedPeptideEvidenceList) Len() int           { return len(a) }
func (a CombinedPeptideEvidenceList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CombinedPeptideEvidenceList) Less(i, j int) bool { return a[i].Sequence < a[j].Sequence }

// CombinedPSMEvidence represents all combined PSMs detected
type CombinedPSMEvidence struct {
	DataSet            string
	Source             string
	Spectrum           string
	SpectrumFile       string
	Peptide            string
	ModifiedPeptide    string
	Protein            string
	ProteinDescription string
	ProteinID          string
	EntryName          string
	GeneName           string
	AssumedCharge      uint8
	IsUnique           bool
	IsUsed             bool
	Intensity          map[string]float64
	Labels             map[string]iso.Labels
}

// CombinedPSMEvidenceList is a list of Combined PSM Evidences
type CombinedPSMEvidenceList []CombinedPSMEvidence

func (a CombinedPSMEvidenceList) Len() int           { return len(a) }
func (a CombinedPSMEvidenceList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CombinedPSMEvidenceList) Less(i, j int) bool { return a[i].Spectrum < a[j].Spectrum }

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

	var isComet bool
	var hasLoc bool
	var hasLabels bool
	var isoBrand string
	var isoChannels int

	if len(m.Comet.Param) > 0 {
		isComet = true
	}

	if m.MSFragger.LocalizeDeltaMass == 1 {
		hasLoc = true
	}

	if m.Quantify.Brand == "tmt" {
		isoBrand = "tmt"
	} else if m.Quantify.Brand == "itraq" {
		isoBrand = "itraq"
	} else if m.Quantify.Brand == "xtag" {
		isoBrand = "xtag"
	}

	if len(m.Quantify.Plex) > 0 {
		isoChannels, _ = strconv.Atoi(m.Quantify.Plex)
	}

	if len(m.Quantify.Annot) > 0 {
		hasLabels = true
	}

	logrus.Info("Creating reports")
	{
		var repoPSM PSMEvidenceList
		RestorePSM(&repoPSM)
		// PSM
		repoPSM.MetaPSMReport(m.Home, isoBrand, m.Database.Tag, isoChannels, m.Report.Decoys, isComet, hasLoc, m.Report.IonMob, hasLabels, m.Report.Prefix)
	}
	{
		var repoIons IonEvidenceList
		RestoreIon(&repoIons)
		// Ion
		repoIons.MetaIonReport(m.Home, isoBrand, m.Database.Tag, isoChannels, m.Report.Decoys, hasLabels, m.Report.Prefix)
	}
	{
		// Peptide
		var repoPeptides PeptideEvidenceList
		RestorePeptide(&repoPeptides)
		repoPeptides.MetaPeptideReport(m.Home, isoBrand, m.Database.Tag, isoChannels, m.Report.Decoys, hasLabels, m.Report.Prefix)
	}
	// Protein
	if len(m.Filter.Pox) > 0 || m.Filter.Inference {
		var repoProteins ProteinEvidenceList
		RestoreProtein(&repoProteins)
		repoProteins.MetaProteinReport(m.Home, isoBrand, m.Database.Tag, isoChannels, m.Report.Decoys, m.Filter.Razor, m.Quantify.Unique, hasLabels, m.Report.Prefix)
		repoProteins.ProteinFastaReport(m.Home, m.Report.Decoys)
	}

	// Modifications
	repo := New()
	if len(repo.Modifications.MassBins) > 0 {
		if repo.PSM == nil {
			RestorePSM(&repo.PSM)
		}
		repo.ModificationReport(m.Home, m.Report.Prefix)

		if m.PTMProphet.InputFiles != nil || len(m.PTMProphet.InputFiles) > 0 {
			repo.PSMLocalizationReport(m.Home, m.Filter.Tag, m.Filter.Razor, m.Report.Decoys)
		}

		repo.PlotMassHist()
	}

	// MSstats
	if m.Report.MSstats {
		repo.MetaMSstatsReport(m.Home, isoBrand, isoChannels, m.Report.Decoys, m.Report.Prefix)
	}

	// MzID
	if m.Report.MZID {
		repo.RestoreGranular()
		repo.MzIdentMLReport(m.Version, m.Database.Annot)
	}

}

// prepares the list of modifications to be printed by the report functions
func getModsList(m map[string]mod.Modification) ([]string, []string) {

	var a []string
	var o []string

	for _, i := range m {
		if i.Type == mod.Assigned && i.Name != "Unknown" {
			pos := ""
			if len(i.AminoAcid) == 1 {
				pos = strconv.Itoa(i.Position)
			}
			a = append(a, fmt.Sprintf("%s%s(%.4f)", pos, i.AminoAcid, i.MassDiff))
		}
		if i.Type == mod.Observed && i.Name != "Unknown" {
			for k, v := range i.IsobaricMods {
				o = append(o, fmt.Sprintf("%s(%f)", k, v))
			}
		}
	}

	return a, o
}
