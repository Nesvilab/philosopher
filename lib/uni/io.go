package uni

import (
	"io/ioutil"

	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/vmihailenco/msgpack"
)

// Serialize UniMod data structure
func (d *MOD) Serialize() *err.Error {

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
func (d *MOD) Restore() *err.Error {

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
