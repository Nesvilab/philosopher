package meta

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
	UUID        string
	Home        string
	Temp        string
	MetaFile    string
	MetaDir     string
	DB          string
	OS          string
	Arch        string
	Distro      string
	TimeStamp   string
	ProjectName string
	Database    Database
	Comet       Comet
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
	Print     bool
}

// // Experimental data
// type Experimental struct {
// 	ProjectName string
// 	DecoyTag    string
// 	ConTag      string
// 	PsmFDR      float64
// 	PepFDR      float64
// 	IonFDR      float64
// 	PrtFDR      float64
// 	PepProb     float64
// 	PrtProb     float64
// 	topPepProb  float64
// 	CometParam  []byte
// }

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
