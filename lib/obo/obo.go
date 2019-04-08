package obo

import (
	"io/ioutil"

	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/met"
	unmd "github.com/prvst/philosopher/lib/obo/unimod"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/vmihailenco/msgpack"
)

// DataFormat defines different data type from PSI
type DataFormat interface {
	Serialize() *err.Error
	Restore() *err.Error
}

// Mod contains UniMod term definition
type Mod struct {
	met.Data
	Flavor           string
	OboFile          string
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
func New(flavor string) (DataFormat, *err.Error) {

	var m met.Data
	var e error
	m.Restore(sys.Meta())

	switch flavor {
	case "unimod":
		o := &Mod{}

		o.Flavor = "unimod.obo"
		o.UUID = m.UUID
		o.Distro = m.Distro
		o.Home = m.Home
		o.MetaFile = m.MetaFile
		o.MetaDir = m.MetaDir
		o.DB = m.DB
		o.Temp = m.Temp
		o.TimeStamp = m.TimeStamp

		// Deploy
		o.OboFile, e = unmd.Deploy(m.Temp)
		if e != nil {
			return nil, &err.Error{Type: err.CannotDeployAsset, Class: err.FATA}
		}

		// Read
		unmd.Parse(o.OboFile)

		// Serielize
		o.Serialize()

		return o, nil
	case "psi-ms":
	case "psi-mod":

	}

	return nil, nil
}

// Serialize UniMod data structure
func (m *Mod) Serialize() *err.Error {

	b, er := msgpack.Marshal(&m)
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
func (m *Mod) Restore() *err.Error {

	b, e := ioutil.ReadFile(sys.MODBin())
	if e != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: ": database data may be corrupted"}
	}

	e = msgpack.Unmarshal(b, &m)
	if e != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: ": database data may be corrupted"}
	}

	return nil
}
