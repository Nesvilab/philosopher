package comet

import (
	"errors"
	"io/ioutil"
)

// UnixParameterFile ...
func UnixParameterFile(unixParam string) error {

	param, err := Asset("comet.params")
	err = ioutil.WriteFile(unixParam, param, 0644)

	if err != nil {
		return errors.New("Cannot deploy Comet parameter file")
	}

	return nil
}

// Unix64 ...
func Unix64(unix64 string) error {

	bin, err := Asset("comet.2016012.linux.exe")
	err = ioutil.WriteFile(unix64, bin, 0755)

	if err != nil {
		return errors.New("Cannot deploy Comet parameter file")
	}

	return nil
}
