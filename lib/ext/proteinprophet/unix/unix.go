package proteinprophet

import (
	"errors"
	"io/ioutil"

	"github.com/prvst/philosopher/lib/sys"
)

// UnixBatchCoverage ...
func UnixBatchCoverage(s string) error {

	bin, err := Asset("batchcoverage")
	err = ioutil.WriteFile(s, bin, sys.FilePermission())

	if err != nil {
		return errors.New("Cannot deploy batchcoverage")
	}

	return nil
}

// UnixDatabaseParser ...
func UnixDatabaseParser(s string) error {

	bin, err := Asset("DatabaseParser")
	err = ioutil.WriteFile(s, bin, sys.FilePermission())

	if err != nil {
		return errors.New("Cannot deploy DatabaseParser")
	}

	return nil
}

// UnixProteinProphet ...
func UnixProteinProphet(s string) error {

	bin, err := Asset("ProteinProphet")
	err = ioutil.WriteFile(s, bin, sys.FilePermission())

	if err != nil {
		return errors.New("Cannot deploy ProteinProphet")
	}

	return nil
}
