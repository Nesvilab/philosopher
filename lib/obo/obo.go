package obo

import (
	"io/ioutil"

	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/vmihailenco/msgpack"
)

// DataFormat defines different data type from PSI
type DataFormat interface {
	Parse(s string) error
	Serialize() *err.Error
	Restore() *err.Error
}

// Mod contains UniMod term definition
type Mod struct {
	met.Data
	MonoIsotopicMass float64
	AverageMass      float64
	Composition      string
	Term             []Term
}

// Term refers to an atomic ontology definition
type Term struct {
	ID               string
	Name             string
	Definition       string
	DateTimePosted   string
	DateTimeModified string
	Comments         string
	Synonyms         []string
	IsA              string
}

// New Onto constructor
func New(flavor string) DataFormat {

	var m met.Data
	m.Restore(sys.Meta())

	switch flavor {
	case "unimod":
		o := &Mod{}
		o.UUID = m.UUID
		o.Distro = m.Distro
		o.Home = m.Home
		o.MetaFile = m.MetaFile
		o.MetaDir = m.MetaDir
		o.DB = m.DB
		o.Temp = m.Temp
		o.TimeStamp = m.TimeStamp
		return o
	case "psi-ms":
	case "psi-mod":

	}

	return nil
}

// GetUniModTerms deploys, reads and assemble the unimod data into structs
// func GetUniModTerms(temp string) (Onto, error) {

// 	var e error
// 	var o Onto

// 	// deploys unimod database
// 	f, e := mod.DeployUniModObo(temp)
// 	if e != nil {
// 		return o, e
// 	}

// 	// process xml file and load structs
// 	e = o.Parse(f)
// 	if e != nil {
// 		return o, e
// 	}

// 	o.Serialize()

// 	return o, nil
// }

// Parse reads the unimod.obo file and creates the data structure
func (d *Mod) Parse(s string) error {

	return nil
}

// Serialize UniMod data structure
func (d *Mod) Serialize() *err.Error {

	b, er := msgpack.Marshal(&d)
	if er != nil {
		return &err.Error{Type: err.CannotCreateOutputFile, Class: err.FATA, Argument: er.Error()}
	}

	er = ioutil.WriteFile(sys.MODBin(), b, sys.FilePermission())
	if er != nil {
		return &err.Error{Type: err.CannotSerializeData, Class: err.FATA, Argument: er.Error()}
	}

	return nil
}

// Restore reads philosopher results files and restore the data sctructure
func (d *Mod) Restore() *err.Error {

	b, e := ioutil.ReadFile(sys.MODBin())
	if e != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: ": database data may be corrupted"}
	}

	e = msgpack.Unmarshal(b, &d)
	if e != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: ": database data may be corrupted"}
	}

	return nil
}
