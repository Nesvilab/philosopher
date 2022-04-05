package met

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"time"

	"philosopher/lib/msg"

	"philosopher/lib/sys"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v5"
)

// Data is the global parameter container
type Data struct {
	UUID           string
	Home           string
	Temp           string
	MetaFile       string
	MetaDir        string
	DB             string
	OS             string
	Arch           string
	Distro         string
	TimeStamp      string
	Version        string
	Build          string
	ProjectName    string
	SearchEngine   string
	Msconvert      Msconvert
	Idconvert      Idconvert
	Database       Database
	MSFragger      MSFragger
	Comet          Comet
	PeptideProphet PeptideProphet
	InterProphet   InterProphet
	ProteinProphet ProteinProphet
	PTMProphet     PTMProphet
	Filter         Filter
	Quantify       Quantify
	BioQuant       BioQuant
	Abacus         Abacus
	Report         Report
	TMTIntegrator  TMTIntegrator
	Index          Index
	Pipeline       Pipeline
}

// Msconvert options and parameters
type Msconvert struct {
	Output                  string
	Format                  string
	MZBinaryEncoding        string
	IntensityBinaryEncoding string
	NoIndex                 bool
	Zlib                    bool
}

// Idconvert optioons and parameters
type Idconvert struct {
	Format string
}

// Database options and parameters
type Database struct {
	ID        string `yaml:"id"`
	Annot     string `yaml:"protein_database"`
	Enz       string `yaml:"enzyme"`
	Tag       string `yaml:"decoy_tag"`
	Add       string `yaml:"add"`
	Custom    string `yaml:"custom"`
	TimeStamp string `yaml:"timestamp"`
	Crap      bool   `yaml:"contam"`
	CrapTag   bool   `yaml:"contaminant_tag"`
	Rev       bool   `yaml:"reviewed"`
	Iso       bool   `yaml:"isoform"`
	NoD       bool   `yaml:"nodecoys"`
}

// Comet options and parameters
type Comet struct {
	Param        string `yaml:"param"`
	RawExtension string `yaml:"raw"`
	RawFiles     []string
	ParamFile    []byte
	Print        bool
	NoIndex      bool `yaml:"noindex"`
}

// MSFragger options and parameters
type MSFragger struct {
	JarPath                            string `yaml:"path"`
	Extension                          string `yaml:"extension"`
	DatabaseName                       string `yaml:"database_name"`
	DecoyPrefix                        string `yaml:"decoy_prefix"`
	IsotopeError                       string `yaml:"isotope_error"`
	MassOffsets                        string `yaml:"mass_offsets"`
	PrecursorMassMode                  string `yaml:"precursor_mass_mode"`
	DeltaMassExcludeRanges             string `yaml:"delta_mass_exclude_ranges"`
	FragmentIonSeries                  string `yaml:"fragment_ion_series"`
	IonSeriesDefinitions               string `yaml:"ion_series_definitions"`
	SearchEnzymeName1                  string `yaml:"search_enzyme_name_1"`
	SearchEnzymeCut1                   string `yaml:"search_enzyme_cut_1"`
	SearchEnzymeNocut1                 string `yaml:"search_enzyme_nocut_1"`
	SearchEnzymeSense1                 string `yaml:"search_enzyme_sense_1"`
	SearchEnzymeName2                  string `yaml:"search_enzyme_name_2"`
	SearchEnzymeCut2                   string `yaml:"search_enzyme_cut_2"`
	SearchEnzymeNocut2                 string `yaml:"search_enzyme_nocut_2"`
	SearchEnzymeSense2                 string `yaml:"search_enzyme_sense_2"`
	OutputFormat                       string `yaml:"output_format"`
	PrecursorCharge                    string `yaml:"precursor_charge"`
	DigestMassRange                    string `yaml:"digest_mass_range"`
	ClearMzRange                       string `yaml:"clear_mz_range"`
	RemovePrecursorRange               string `yaml:"remove_precursor_range"`
	LabileSearchMode                   string `yaml:"labile_search_mode"`
	RestrictDeltaMassTo                string `yaml:"restrict_deltamass_to"`
	DiagnosticFragments                string `yaml:"diagnostic_fragments"`
	YTypeMasses                        string `yaml:"Y_type_masses"`
	VariableMod01                      string `yaml:"variable_mod_01"`
	VariableMod02                      string `yaml:"variable_mod_02"`
	VariableMod03                      string `yaml:"variable_mod_03"`
	VariableMod04                      string `yaml:"variable_mod_04"`
	VariableMod05                      string `yaml:"variable_mod_05"`
	VariableMod06                      string `yaml:"variable_mod_06"`
	VariableMod07                      string `yaml:"variable_mod_07"`
	RawFiles                           []string
	Memory                             int     `yaml:"memory"`
	Threads                            int     `yaml:"num_threads"`
	DataType                           int     `yaml:"data_type"`
	PrecursorMassLower                 int     `yaml:"precursor_mass_lower"`
	PrecursorMassUpper                 int     `yaml:"precursor_mass_upper"`
	PrecursorMassUnits                 int     `yaml:"precursor_mass_units"`
	PrecursorTrueTolerance             int     `yaml:"precursor_true_tolerance"`
	PrecursorTrueUnits                 int     `yaml:"precursor_true_units"`
	FragmentMassUnits                  int     `yaml:"fragment_mass_units"`
	CalibrateMass                      int     `yaml:"calibrate_mass"`
	UseAllModsInFirstSearch            int     `yaml:"use_all_mods_in_first_search"`
	WriteCalibratedMGF                 int     `yaml:"write_calibrated_mgf"`
	EvaluateMassCalibration            int     `yaml:"evaluate_mass_calibration"`
	Deisotope                          int     `yaml:"deisotope"`
	Deneutralloss                      int     `yaml:"deneutralloss"`
	LocalizeDeltaMass                  int     `yaml:"localize_delta_mass"`
	AllowedMissedCleavage1             int     `yaml:"allowed_missed_cleavage_1"`
	AllowedMissedCleavage2             int     `yaml:"allowed_missed_cleavage_2"`
	NumEnzymeTermini                   int     `yaml:"num_enzyme_termini"`
	ClipNTermM                         int     `yaml:"clip_nTerm_M"`
	AllowMultipleVariableModsOnResidue int     `yaml:"allow_multiple_variable_mods_on_residue"`
	MaxVariableModsPerPeptide          int     `yaml:"max_variable_mods_per_peptide"`
	MaxVariableModsCombinations        int     `yaml:"max_variable_mods_combinations"`
	OutputReportTopN                   int     `yaml:"output_report_topN"`
	OutputMaxExpect                    int     `yaml:"output_max_expect"`
	ReportAlternativeProteins          int     `yaml:"report_alternative_proteins"`
	OverrideCharge                     int     `yaml:"override_charge"`
	DigestMinLength                    int     `yaml:"digest_min_length"`
	DigestMaxLength                    int     `yaml:"digest_max_length"`
	MaxFragmentCharge                  int     `yaml:"max_fragment_charge"`
	TrackZeroTopN                      int     `yaml:"track_zero_topN"`
	ZeroBinAcceptExpect                int     `yaml:"zero_bin_accept_expect"`
	ZeroBinMultExpect                  int     `yaml:"zero_bin_mult_expect"`
	AddTopNComplementary               int     `yaml:"add_topN_complementary"`
	CheckSpectralFiles                 int     `yaml:"check_spectral_files"`
	MinimumPeaks                       int     `yaml:"minimum_peaks"`
	UseTopNPeaks                       int     `yaml:"use_topN_peaks"`
	MinFragmentsModelling              int     `yaml:"min_fragments_modelling"`
	MinMatchedFragments                int     `yaml:"min_matched_fragments"`
	RemovePrecursorPeak                int     `yaml:"remove_precursor_peak"`
	IntensityTransform                 int     `yaml:"intensity_transform"`
	MassDiffToVariableMod              int     `yaml:"mass_diff_to_variable_mod"`
	DiagnosticIntensityFilter          int     `yaml:"diagnostic_intensity_filter"`
	MinimumRatio                       float64 `yaml:"minimum_ratio"`
	FragmentMassTolerance              float64 `yaml:"fragment_mass_tolerance"`
	AddCtermPeptide                    float64 `yaml:"add_Cterm_peptide"`
	AddCtermProtein                    float64 `yaml:"add_Cterm_protein"`
	AddNTermPeptide                    float64 `yaml:"add_Nterm_peptide"`
	AddNtermProteine                   float64 `yaml:"add_Nterm_protein"`
	AddAlanine                         float64 `yaml:"add_A_alanine"`
	AddCysteine                        float64 `yaml:"add_C_cysteine"`
	AddAsparticAcid                    float64 `yaml:"add_D_aspartic_acid"`
	AddGlutamicAcid                    float64 `yaml:"add_E_glutamic_acid"`
	AddPhenylAlnine                    float64 `yaml:"add_F_phenylalanine"`
	AddGlycine                         float64 `yaml:"add_G_glycine"`
	AddHistidine                       float64 `yaml:"add_H_histidine"`
	AddIsoleucine                      float64 `yaml:"add_I_isoleucine"`
	AddLysine                          float64 `yaml:"add_K_lysine"`
	AddLeucine                         float64 `yaml:"add_L_leucine"`
	AddMethionine                      float64 `yaml:"add_M_methionine"`
	AddAsparagine                      float64 `yaml:"add_N_asparagine"`
	AddProline                         float64 `yaml:"add_P_proline"`
	AddGlutamine                       float64 `yaml:"add_Q_glutamine"`
	AddArginine                        float64 `yaml:"add_R_arginine"`
	AddSerine                          float64 `yaml:"add_S_serine"`
	AddThreonine                       float64 `yaml:"add_T_threonine"`
	AddValine                          float64 `yaml:"add_V_valine"`
	AddTryptophan                      float64 `yaml:"add_W_tryptophan"`
	AddTyrosine                        float64 `yaml:"add_Y_tyrosine"`
	Param                              string  `yaml:"param"`
	ParamFile                          []byte
	//SearchEnzymeName                   string  `yaml:"search_enzyme_name"`
	//SearchEnzymeCutafter               string  `yaml:"search_enzyme_cutafter"`
	//SearchEnzymeButNotAfter            string  `yaml:"search_enzyme_butnotafter"`
	//AllowedMissedCleavage              int     `yaml:"allowed_missed_cleavage"`
}

// PeptideProphet options and parameters
type PeptideProphet struct {
	FileExtension string `yaml:"extension"`
	Output        string `yaml:"output"`
	Database      string `yaml:"database"`
	Rtcat         string `yaml:"rtcat"`
	Decoy         string `yaml:"decoy"`
	Enzyme        string `yaml:"enzyme"`
	Ignorechg     string `yaml:"ignorechg"`
	InputFiles    []string
	Minpiprob     float64 `yaml:"minpiprob"`
	Minrtprob     float64 `yaml:"minrtprob"`
	Minprob       float64 `yaml:"minprob"`
	Masswidth     float64 `yaml:"masswidth"`
	MinPepLen     int     `yaml:"minpeplen"`
	Clevel        int     `yaml:"clevel"`
	Minpintt      int     `yaml:"minpintt"`
	Minrtntt      int     `yaml:"minrtntt"`
	Combine       bool    `yaml:"combine"`
	Exclude       bool    `yaml:"exclude"`
	Leave         bool    `yaml:"leave"`
	Perfectlib    bool    `yaml:"perfectlib"`
	Icat          bool    `yaml:"icat"`
	Noicat        bool    `yaml:"noicat"`
	Zero          bool    `yaml:"zero"`
	Accmass       bool    `yaml:"accmass"`
	Ppm           bool    `yaml:"ppm"`
	Nomass        bool    `yaml:"nomass"`
	Pi            bool    `yaml:"pi"`
	Rt            bool    `yaml:"rt"`
	Glyc          bool    `yaml:"glyc"`
	Phospho       bool    `yaml:"phospho"`
	Maldi         bool    `yaml:"maldi"`
	Instrwarn     bool    `yaml:"instrwarn"`
	Decoyprobs    bool    `yaml:"decoyprobs"`
	Nontt         bool    `yaml:"nontt"`
	Nonmc         bool    `yaml:"nonmc"`
	Expectscore   bool    `yaml:"expectscore"`
	Nonparam      bool    `yaml:"nonparam"`
	Neggamma      bool    `yaml:"neggamma"`
	Forcedistr    bool    `yaml:"forcedistr"`
	Optimizefval  bool    `yaml:"optimizefval"`
	Concurrent    bool    `yaml:"concurrent"`
}

// InterProphet options and parameters
type InterProphet struct {
	InputFiles []string
	Output     string  `yaml:"output"`
	Decoy      string  `yaml:"decoy"`
	Cat        string  `yaml:"cat"`
	Threads    int     `yaml:"threads"`
	MinProb    float64 `yaml:"minprob"`
	Length     bool    `yaml:"length"`
	Nofpkm     bool    `yaml:"nofpkm"`
	Nonss      bool    `yaml:"nonss"`
	Nonse      bool    `yaml:"nonse"`
	Nonrs      bool    `yaml:"nonrs"`
	Nonsm      bool    `yaml:"nonsm"`
	Nonsp      bool    `yaml:"nonsp"`
	Sharpnse   bool    `yaml:"sharpnse"`
	Nonsi      bool    `yaml:"nonsi"`
}

// ProteinProphet options and parameters
type ProteinProphet struct {
	InputFiles  []string
	Output      string  `yaml:"output"`
	Minindep    int     `yaml:"minidep"`
	Mufactor    int     `yaml:"mufactor"`
	Maxppmdiff  int     `yaml:"maxppmdiff"`
	Minprob     float64 `yaml:"minprob"`
	ExcludeZ    bool    `yaml:"excludez"`
	Noplot      bool    `yaml:"noplot"`
	Nooccam     bool    `yaml:"noocam"`
	Softoccam   bool    `yaml:"softocam"`
	Icat        bool    `yaml:"icat"`
	Glyc        bool    `yaml:"glyc"`
	Nogroupwts  bool    `yaml:"nogroupwts"`
	NonSP       bool    `yaml:"nonsp"`
	Subgroups   bool    `yaml:"subgroups"`
	Accuracy    bool    `yaml:"accuracy"`
	Asap        bool    `yaml:"asap"`
	Refresh     bool    `yaml:"refresh"`
	Normprotlen bool    `yaml:"normprotlen"`
	Logprobs    bool    `yaml:"logprobs"`
	Confem      bool    `yaml:"confem"`
	Allpeps     bool    `yaml:"allpeps"`
	Unmapped    bool    `yaml:"unmapped"`
	Noprotlen   bool    `yaml:"noprotlen"`
	Instances   bool    `yaml:"instances"`
	Fpkm        bool    `yaml:"fpkm"`
	Protmw      bool    `yaml:"protmw"`
	Iprophet    bool    `yaml:"iprophet"`
	Asapprophet bool    `yaml:"asapprophet"`
	Delude      bool    `yaml:"delude"`
	Excludemods bool    `yaml:"excludemods"`
}

// PTMProphet options and parameters
type PTMProphet struct {
	InputFiles         []string
	Output             string  `yaml:"output"`
	Mods               string  `yaml:"mods"`
	NIons              string  `yaml:"nions"`
	CIons              string  `yaml:"cions"`
	EM                 int     `yaml:"em"`
	FragPPMTol         int     `yaml:"fragppmtol"`
	MaxThreads         int     `yaml:"maxthreads"`
	MaxFragZ           int     `yaml:"maxfragz"`
	Mino               int     `yaml:"mino"`
	MassOffset         int     `yaml:"massoffset"`
	PPMTol             float64 `yaml:"ppmtol"`
	MinProb            float64 `yaml:"minprob"`
	ExcludeMassDiffMin float64 `yaml:"excludemassdiffmin"`
	ExcludeMassDiffMax float64 `yaml:"excludemassdiffmax"`
	Static             bool    `yaml:"static"`
	NoUpdate           bool    `yaml:"noupdate"`
	KeepOld            bool    `yaml:"keepold"`
	Verbose            bool    `yaml:"verbose"`
	MassDiffMode       bool    `yaml:"massdiffmode"`
	Lability           bool    `yaml:"lability"`
	Direct             bool    `yaml:"direct"`
	Ifrags             bool    `yaml:"ifrags"`
	Autodirect         bool    `yaml:"autodirect"`
	NoMinoFactor       bool    `yaml:"nominofactor"`
}

// Filter options and parameters
type Filter struct {
	Pex       string  `yaml:"pepxml"`
	Pox       string  `yaml:"protxml"`
	Tag       string  `yaml:"tag"`
	Mods      string  `yaml:"mods"`
	RazorBin  string  `yaml:"razorbin"`
	PsmFDR    float64 `yaml:"psmFDR"`
	PepFDR    float64 `yaml:"peptideFDR"`
	IonFDR    float64 `yaml:"ionFDR"`
	PtFDR     float64 `yaml:"proteinFDR"`
	ProtProb  float64 `yaml:"proteinProbability"`
	PepProb   float64 `yaml:"peptideProbability"`
	Weight    float64 `yaml:"peptideWeight"`
	Model     bool    `yaml:"models"`
	Razor     bool    `yaml:"razor"`
	Picked    bool    `yaml:"picked"`
	Seq       bool    `yaml:"sequential"`
	TwoD      bool    `yaml:"two-dimensional"`
	Mapmods   bool    `yaml:"mapMods"`
	Inference bool
}

// Quantify options and parameters
type Quantify struct {
	Pex        string  `yaml:"pepxml"`
	Tag        string  `yaml:"tag"`
	Format     string  `yaml:"format"`
	Dir        string  `yaml:"dir"`
	Brand      string  `yaml:"brand"`
	Plex       string  `yaml:"plex"`
	ChanNorm   string  `yaml:"chanNorm"`
	Annot      string  `yaml:"annotation"`
	Level      int     `yaml:"level"`
	RTWin      float64 `yaml:"retentionTimeWindow"`
	PTWin      float64 `yaml:"peakTimeWindow"`
	Tol        float64 `yaml:"tolerance"`
	Purity     float64 `yaml:"purity"`
	MinProb    float64 `yaml:"minprob"`
	RemoveLow  float64 `yaml:"removeLow"`
	Isolated   bool    `yaml:"isolated"`
	IntNorm    bool    `yaml:"intNorm"`
	Unique     bool    `yaml:"uniqueOnly"`
	BestPSM    bool    `yaml:"bestPSM"`
	Raw        bool    `yaml:"raw"`
	Faims      bool    `yaml:"faims"`
	LabelNames map[string]string
}

// Abacus options ad parameters
type Abacus struct {
	Tag      string  `yaml:"tag"`
	ProtProb float64 `yaml:"proteinProbability"`
	PepProb  float64 `yaml:"peptideProbability"`
	Peptide  bool    `yaml:"peptide"`
	Protein  bool    `yaml:"protein"`
	Razor    bool    `yaml:"razor"`
	Picked   bool    `yaml:"picked"`
	Labels   bool    `yaml:"labels"`
	Unique   bool    `yaml:"uniqueOnly"`
	Reprint  bool    `yaml:"reprint"`
	Full     bool    `yaml:"full"`
}

// BioQuant options and parameters
type BioQuant struct {
	UID   string  `yaml:"organismUniProtID"`
	Level float64 `yaml:"level"`
}

// Report options and parameters
type Report struct {
	Decoys  bool `yaml:"withDecoys"`
	MSstats bool `yaml:"msstats"`
	MZID    bool `yaml:"mzID"`
	IonMob  bool `yaml:"ionmobility"`
}

// TMTIntegrator options and parameters
type TMTIntegrator struct {
	JarPath   string `yaml:"path"`
	Memory    int    `yaml:"memory"`
	Param     string `yaml:"param"`
	Files     []string
	ParamFile []byte
}

// Index options and parameters
type Index struct {
	Spectra string
}

// Pipeline options and parameters
type Pipeline struct {
	Directives string
	Print      bool
	Verbose    bool
}

// New initializes the structure with the system information needed
// to run all the follwing commands
func New(h string) Data {

	var d Data

	var fmtuuid = uuid.NewV4()
	var uuid = fmtuuid.String()
	d.UUID = uuid

	d.OS = runtime.GOOS
	d.Arch = runtime.GOARCH

	distro := sys.GetLinuxFlavor()

	d.Distro = distro

	d.Home = h
	d.ProjectName = string(filepath.Base(h))

	d.MetaFile = d.Home + string(filepath.Separator) + sys.Meta()
	d.MetaDir = d.Home + string(filepath.Separator) + sys.MetaDir()

	d.DB = d.Home + string(filepath.Separator) + sys.DBBin()

	temp := sys.GetTemp()
	temp += string(filepath.Separator) + uuid
	d.Temp = temp

	t := time.Now()
	d.TimeStamp = t.Format(time.RFC3339)

	return d
}

// CleanTemp removes all files from the given temp directory
func CleanTemp(dir string) {

	e := os.RemoveAll(dir)
	if e != nil {
		msg.Custom(e, "error")
	}
}

// Serialize converts the whole structure to a gob file
func (d *Data) Serialize() {

	b, e := msgpack.Marshal(&d)
	if e != nil {
		msg.MarshalFile(e, "fatal")
	}

	e = ioutil.WriteFile(sys.Meta(), b, sys.FilePermission())
	if e != nil {
		msg.WriteFile(e, "fatal")
	}

}

// Restore reads philosopher results files and restore the data sctructure
func (d *Data) Restore(f string) {

	b, e1 := ioutil.ReadFile(f)

	e2 := msgpack.Unmarshal(b, &d)

	if e1 != nil && e2 != nil {
		msg.Custom(errors.New("workspace not detected"), "warning")
	} else if len(d.UUID) < 1 {
		msg.Custom(errors.New("the current Workspace is corrupted or was created with an older version. Please remove it and create a new one"), "warning")
	}

	// checks if the temp is still there, if not recreate it
	if _, err := os.Stat(d.Temp); os.IsNotExist(err) {
		os.Mkdir(d.Temp, sys.FilePermission())
	}

}

// FunctionInitCheckUp does initilization checkup and verification if meta and temp folders are up.
// In case not, meta trows an error and folder is created.
func (d Data) FunctionInitCheckUp() {

	if len(d.UUID) < 1 && len(d.Home) < 1 {
		msg.WorkspaceNotFound(errors.New("failed to checkup the initialization"), "fatal")
	}

	if _, e := os.Stat(d.Temp); os.IsNotExist(e) && len(d.UUID) > 0 {
		os.Mkdir(d.Temp, sys.FilePermission())
		msg.LocatingTemDirecotry(e, "warning")
	}

}

// ToCmdString converts the MSFragger struct into a CMD string
func (d MSFragger) ToCmdString() {

	var cmd = "CMD string: philosopher msfragger"

	v := reflect.ValueOf(d)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {

		if typeOfS.Field(i).Name == "Param" || typeOfS.Field(i).Name == "RawFiles" || typeOfS.Field(i).Name == "ParamFile" {
			continue
		}

		cmd = fmt.Sprintf("%s --%s %v", cmd, typeOfS.Field(i).Name, v.Field(i).Interface())
	}

	logrus.Info(cmd)

}

// ToCmdString converts the PeptideProphet struct into a CMD string
func (d PeptideProphet) ToCmdString() {

	var cmd = "CMD string: philosopher peptideprophet"

	v := reflect.ValueOf(d)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		cmd = fmt.Sprintf("%s --%s %v", cmd, typeOfS.Field(i).Name, v.Field(i).Interface())
	}

	logrus.Info(cmd)

}
