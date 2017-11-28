package raw

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/raw/mz"
	"github.com/prvst/philosopher/lib/raw/mzml"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/vmihailenco/msgpack"
)

// Data represents parsed and processed MZ data from mz files
type Data struct {
	Raw mz.Raw
}

// IndexMz receives a list of mz files and creates a binary index for each one
func IndexMz(f []string) *err.Error {

	for _, i := range f {

		var d Data

		if strings.Contains(i, "mzml") || strings.Contains(i, "mzML") {

			raw, e := mzml.Read(i)
			if e != nil {
				return e
			}

			d.Raw = raw

		} else if strings.Contains(i, "mzxml") || strings.Contains(i, "mzXML") {
			return &err.Error{Type: err.MethodNotImplemented, Class: err.FATA, Argument: "mzXML reader not implemented"}
		}

		d.Serialize()
	}

	return nil
}

// Serialize mz data structure to binary format
func (data *Data) Serialize() *err.Error {

	// remove the extension
	var extension = filepath.Ext(filepath.Base(data.Raw.FileName))
	var name = data.Raw.FileName[0 : len(data.Raw.FileName)-len(extension)]

	output := fmt.Sprintf("%s%s%s.bin", sys.MetaDir(), string(filepath.Separator), filepath.Base(name))

	// create a file
	dataFile, e := os.Create(output)
	if e != nil {
		return &err.Error{Type: err.CannotCreateOutputFile, Class: err.FATA, Argument: e.Error()}
	}

	dataEncoder := msgpack.NewEncoder(dataFile)
	goberr := dataEncoder.Encode(data)
	if goberr != nil {
		return &err.Error{Type: err.CannotSerializeData, Class: err.FATA, Argument: e.Error()}
	}
	dataFile.Close()

	return nil
}

// Restore reads philosopher results files and restore the data sctructure
func (data *Data) Restore() *err.Error {

	file, _ := os.Open(sys.RawBin())

	dec := msgpack.NewDecoder(file)
	e := dec.Decode(&data)
	if e != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: e.Error()}
	}

	return nil
}
