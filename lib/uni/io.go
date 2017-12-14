package uni

import (
	"errors"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/vmihailenco/msgpack"
)

// Serialize UniMod data structure
func (u *MOD) Serialize() error {

	var err error

	// create a file
	dataFile, err := os.Create(sys.MODBin())
	if err != nil {
		return err
	}

	dataEncoder := msgpack.NewEncoder(dataFile)
	goberr := dataEncoder.Encode(u)
	if goberr != nil {
		logrus.Fatal("Cannot save results, Bad format", goberr)
	}
	dataFile.Close()

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
