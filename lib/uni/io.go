package uni

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/vmihailenco/msgpack"
)

// Serialize UniMod data structure
func (u *MOD) Serialize() *err.Error {

	b, er := msgpack.Marshal(&u)
	if er != nil {
		return &err.Error{Type: err.CannotCreateOutputFile, Class: err.FATA, Argument: er.Error()}
	}

	er = ioutil.WriteFile(sys.MODBin(), b, 0644)
	if er != nil {
		return &err.Error{Type: err.CannotSerializeData, Class: err.FATA, Argument: er.Error()}
	}

	return nil
}

// Restore reads philosopher results files and restore the data sctructure
func (u *MOD) Restore() error {

	file, _ := os.Open(sys.MODBin())

	dec := msgpack.NewDecoder(file)
	err := dec.Decode(&u)
	if err != nil {
		return errors.New("Could not restore Philosopher result. Please check file path")
	}

	return nil
}
