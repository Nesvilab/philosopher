package proteinprophet

import (
	"errors"
	"io/ioutil"
)

// UnixBatchCoverage ...
func UnixBatchCoverage(s string) error {

	bin, err := Asset("batchcoverage")
	err = ioutil.WriteFile(s, bin, 0755)

	if err != nil {
		return errors.New("Cannot deploy batchcoverage")
	}

	return nil
}

// UnixDatabaseParser ...
func UnixDatabaseParser(s string) error {

	bin, err := Asset("DatabaseParser")
	err = ioutil.WriteFile(s, bin, 0755)

	if err != nil {
		return errors.New("Cannot deploy DatabaseParser")
	}

	return nil
}

// UnixProteinProphet ...
func UnixProteinProphet(s string) error {

	bin, err := Asset("ProteinProphet")
	err = ioutil.WriteFile(s, bin, 0755)

	if err != nil {
		return errors.New("Cannot deploy ProteinProphet")
	}

	return nil
}
