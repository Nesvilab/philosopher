package met

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher/lib/sys"
	uuid "github.com/satori/go.uuid"
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
	ProjectName    string
	Database       Database
	Comet          Comet
	PeptideProphet PeptideProphet
	InterProphet   InterProphet
	ProteinProphet ProteinProphet
	PTMProphet     PTMProphet
	Filter         Filter
	Quantify       Quantify
}

// Database options and parameters
type Database struct {
	ID     string
	Annot  string
	Enz    string
	Tag    string
	Add    string
	Custom string
	Crap   bool
	Rev    bool
	Iso    bool
}

// Comet options and parameters
type Comet struct {
	Param     string
	ParamFile []byte
	RawFiles  []string
	Print     bool
}

// PeptideProphet options and parameters
type PeptideProphet struct {
	InputFiles   []string
	Output       string
	Database     string
	Rtcat        string
	Decoy        string
	Minpiprob    float64
	Minrtprob    float64
	Minprob      float64
	Masswidth    float64
	MinPepLen    int
	Clevel       int
	Minpintt     int
	Ignorechg    int
	Minrtntt     int
	Combine      bool
	Exclude      bool
	Leave        bool
	Perfectlib   bool
	Icat         bool
	Noicat       bool
	Zero         bool
	Accmass      bool
	Ppm          bool
	Nomass       bool
	Pi           bool
	Rt           bool
	Glyc         bool
	Phospho      bool
	Maldi        bool
	Instrwarn    bool
	Decoyprobs   bool
	Nontt        bool
	Nonmc        bool
	Expectscore  bool
	Nonparam     bool
	Neggamma     bool
	Forcedistr   bool
	Optimizefval bool
}

// InterProphet options and parameters
type InterProphet struct {
	InputFiles []string
	Threads    int
	Decoy      string
	Cat        string
	MinProb    float64
	Output     string
	Length     bool
	Nofpkm     bool
	Nonss      bool
	Nonse      bool
	Nonrs      bool
	Nonsm      bool
	Nonsp      bool
	Sharpnse   bool
	Nonsi      bool
}

// ProteinProphet options and parameters
type ProteinProphet struct {
	InputFiles  []string
	Minprob     float64
	Minindep    int
	Mufactor    int
	Output      string
	Maxppmdiff  int
	ExcludeZ    bool
	Noplot      bool
	Nooccam     bool
	Softoccam   bool
	Icat        bool
	Glyc        bool
	Nogroupwts  bool
	NonSP       bool
	Accuracy    bool
	Asap        bool
	Refresh     bool
	Normprotlen bool
	Logprobs    bool
	Confem      bool
	Allpeps     bool
	Unmapped    bool
	Noprotlen   bool
	Instances   bool
	Fpkm        bool
	Protmw      bool
	Iprophet    bool
	Asapprophet bool
	Delude      bool
	Excludemods bool
}

// PTMProphet options and parameters
type PTMProphet struct {
	InputFiles   []string
	Output       string
	EM           int
	MzTol        float64
	PPMTol       float64
	MinProb      float64
	NoUpdate     bool
	KeepOld      bool
	Verbose      bool
	MassDiffMode bool
}

// Filter options and parameters
type Filter struct {
	Phi      string
	Pex      string
	Pox      string
	Tag      string
	Con      string
	Ptconf   string
	RepProt  string
	Save     string
	Database string
	PsmFDR   float64
	PepFDR   float64
	IonFDR   float64
	PtFDR    float64
	ProtProb float64
	PepProb  float64
	TopPep   bool
	Model    bool
	RepPSM   bool
	Razor    bool
	Picked   bool
	Seq      bool
	Mapmods  bool
}

// Quantify options and parameters
type Quantify struct {
	Phi      string
	Format   string
	Dir      string
	Brand    string
	Plex     string
	ChanNorm string
	RTWin    float64
	PTWin    float64
	Tol      float64
	Purity   float64
	IntNorm  bool
}

var err error

// New initializes the structure with the system information needed
// to run all the follwing commands
func New(h string) Data {

	var d Data

	var fmtuuid = uuid.NewV4()
	var uuid = fmt.Sprintf("%s", fmtuuid)
	d.UUID = uuid

	d.OS = runtime.GOOS
	d.Arch = runtime.GOARCH

	d.Distro, err = sys.GetLinuxFlavor()
	if err != nil {
		logrus.Fatal(err)
	}

	d.Home = h
	d.ProjectName = string(filepath.Base(h))

	d.MetaFile = d.Home + string(filepath.Separator) + sys.Meta()
	d.MetaDir = d.Home + string(filepath.Separator) + sys.MetaDir()

	d.DB = d.MetaDir + string(filepath.Separator) + sys.DBBin()

	d.Temp, err = sys.GetTemp()
	d.Temp += string(filepath.Separator) + uuid

	t := time.Now()
	d.TimeStamp = t.Format(time.RFC3339)

	return d
}

// Serialize converts the whole structure to a gob file
func (d *Data) Serialize() error {

	output := fmt.Sprintf("%s", sys.Meta())

	// create a file
	dataFile, err := os.Create(output)
	if err != nil {
		return err
	}

	dataEncoder := msgpack.NewEncoder(dataFile)
	err = dataEncoder.Encode(d)
	if err != nil {
		return err
	}
	dataFile.Close()

	return nil
}

// Restore reads philosopher results files and restore the data sctructure
func (d *Data) Restore(f string) error {

	file, _ := os.Open(f)

	dec := msgpack.NewDecoder(file)
	err := dec.Decode(&d)
	if err != nil {
		return errors.New("Could not restore meta data")
	}

	if len(d.UUID) < 1 {
		return errors.New("Could not restore meta data")
	}

	if _, err := os.Stat(d.Temp); os.IsNotExist(err) {
		os.Mkdir(d.Temp, 0755)
	}

	return nil
}
