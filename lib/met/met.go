package met

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/sys"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack"
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
	Cluster        Cluster
	Abacus         Abacus
	Report         Report
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
	ID     string `yaml:"id"`
	Annot  string `yaml:"annotate"`
	Enz    string `yaml:"enzyme"`
	Tag    string `yaml:"prefix"`
	Add    string `yaml:"add"`
	Custom string `yaml:"custom"`
	Crap   bool   `yaml:"contam"`
	Rev    bool   `yaml:"reviewed"`
	Iso    bool   `yaml:"isoform"`
}

// Comet options and parameters
type Comet struct {
	Param        string `yaml:"param"`
	ParamFile    []byte
	RawExtension string `yaml:"raw"`
	RawFiles     []string
	Print        bool
	NoIndex      bool `yaml:"noindex"`
}

// MSFragger options and parameters
type MSFragger struct {
	JarPath      string `yaml:"path"`
	Memmory      string `yaml:"memmory"`
	Param        string `yaml:"param"`
	RawExtension string `yaml:"raw"`
	ParamFile    []byte
	RawFiles     []string
}

// PeptideProphet options and parameters
type PeptideProphet struct {
	InputFiles    []string
	FileExtension string  `yaml:"extension"`
	Output        string  `yaml:"output"`
	Database      string  `yaml:"database"`
	Rtcat         string  `yaml:"rtcat"`
	Decoy         string  `yaml:"decoy"`
	Enzyme        string  `yaml:"enzyme"`
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
	InputFiles   []string
	Output       string  `yaml:"output"`
	Mods         string  `yaml:"mods"`
	NIons        string  `yaml:"nions"`
	CIons        string  `yaml:"cions"`
	EM           int     `yaml:"em"`
	FragPPMTol   int     `yaml:"fragppmtol"`
	MaxThreads   int     `yaml:"maxthreads"`
	MaxFragZ     int     `yaml:"maxfragz"`
	Mino         int     `yaml:"mino"`
	MassOffset   int     `yaml:"massoffset"`
	MzTol        float64 `yaml:"mztol"`
	PPMTol       float64 `yaml:"ppmtol"`
	MinProb      float64 `yaml:"minprob"`
	NoUpdate     bool    `yaml:"noupdate"`
	KeepOld      bool    `yaml:"keepold"`
	Verbose      bool    `yaml:"verbose"`
	MassDiffMode bool    `yaml:"massdiffmode"`
	Lability     bool    `yaml:"lability"`
	Direct       bool    `yaml:"direct"`
	Ifrags       bool    `yaml:"ifrags"`
	Autodirect   bool    `yaml:"autodirect"`
	NoMinoFactor bool    `yaml:"nominofactor"`
}

// Filter options and parameters
type Filter struct {
	Pex      string  `yaml:"pepxml"`
	Pox      string  `yaml:"protxml"`
	Tag      string  `yaml:"tag"`
	PsmFDR   float64 `yaml:"psmFDR"`
	PepFDR   float64 `yaml:"peptideFDR"`
	IonFDR   float64 `yaml:"ionFDR"`
	PtFDR    float64 `yaml:"proteinFDR"`
	ProtProb float64 `yaml:"proteinProbability"`
	PepProb  float64 `yaml:"peptideProbability"`
	Weight   float64 `yaml:"peptideWeight"`
	Model    bool    `yaml:"models"`
	Razor    bool    `yaml:"razor"`
	Picked   bool    `yaml:"picked"`
	Seq      bool    `yaml:"sequential"`
	Cap      bool    `yaml:"cappedsequential"`
	Mapmods  bool    `yaml:"mapMods"`
	Fo       bool
}

// Quantify options and parameters
type Quantify struct {
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
	LabelNames map[string]string
}

// Abacus options ad parameters
type Abacus struct {
	CombPro  string  `yaml:"protein"`
	CombPep  string  `yaml:"peptide"`
	Tag      string  `yaml:"tag"`
	Annot    string  `yaml:"annotation"`
	ProtProb float64 `yaml:"proteinProbability"`
	PepProb  float64 `yaml:"peptideProbability"`
	Razor    bool    `yaml:"razor"`
	Picked   bool    `yaml:"picked"`
	Labels   bool    `yaml:"labels"`
	Unique   bool    `yaml:"uniqueOnly"`
	Reprint  bool    `yaml:"reprint"`
}

// Cluster options and parameters
type Cluster struct {
	UID   string  `yaml:"organismUniProtID"`
	Level float64 `yaml:"level"`
}

// Report options and parameters
type Report struct {
	Decoys  bool `yaml:"withDecoys"`
	MSstats bool `yaml:"msstats"`
}

// Index options and parameters
type Index struct {
	Spectra string
}

// Pipeline options and parameters
type Pipeline struct {
	Directives string
	Print      bool
	//Parallel   bool
	//Dataset    string
}

// New initializes the structure with the system information needed
// to run all the follwing commands
func New(h string) Data {

	var d Data

	var fmtuuid, _ = uuid.NewV4()
	var uuid = fmt.Sprintf("%s", fmtuuid)
	d.UUID = uuid

	d.OS = runtime.GOOS
	d.Arch = runtime.GOARCH

	distro, e := sys.GetLinuxFlavor()
	if e != nil {
		logrus.Fatal(e)
	}
	d.Distro = distro

	d.Home = h
	d.ProjectName = string(filepath.Base(h))

	d.MetaFile = d.Home + string(filepath.Separator) + sys.Meta()
	d.MetaDir = d.Home + string(filepath.Separator) + sys.MetaDir()

	d.DB = d.Home + string(filepath.Separator) + sys.DBBin()

	temp, e := sys.GetTemp()
	temp += string(filepath.Separator) + uuid
	d.Temp = temp

	t := time.Now()
	d.TimeStamp = t.Format(time.RFC3339)

	return d
}

// CleanTemp removes all files from the given temp directory
func CleanTemp(dir string) error {

	//os.RemoveAll(tmp)

	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return err
	}
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
	}

	return nil
}

// Serialize converts the whole structure to a gob file
func (d *Data) Serialize() *err.Error {

	b, e := msgpack.Marshal(&d)
	if e != nil {
		return &err.Error{Type: err.CannotSerializeData, Class: err.FATA}
	}

	e = ioutil.WriteFile(sys.Meta(), b, 0644)
	if e != nil {
		return &err.Error{Type: err.CannotSerializeData, Class: err.FATA}
	}

	return nil
}

// Restore reads philosopher results files and restore the data sctructure
func (d *Data) Restore(f string) *err.Error {

	b, e := ioutil.ReadFile(f)
	if e != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA}
	}

	e = msgpack.Unmarshal(b, &d)
	if e != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA}
	}

	if len(d.UUID) < 1 {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA}
	}

	if _, err := os.Stat(d.Temp); os.IsNotExist(err) {
		os.Mkdir(d.Temp, sys.FilePermission())
		//0755
	}

	return nil
}

// FunctionInitCheckUp does initilization checkup and verification if meta and temp folders are up.
// In case not, meta troews an error and folder is created.
func (d Data) FunctionInitCheckUp() *err.Error {

	if len(d.UUID) < 1 && len(d.Home) < 1 {
		return &err.Error{Type: err.WorkspaceNotFound, Class: err.FATA}
	}

	if _, e := os.Stat(d.Temp); os.IsNotExist(e) && len(d.UUID) > 0 {
		os.Mkdir(d.Temp, sys.FilePermission())
		return &err.Error{Type: err.CannotCreateDirectory, Class: err.FATA, Argument: "Can't create temporary directory; check folder permissions"}
	}

	return nil
}
