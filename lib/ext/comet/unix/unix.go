package comet

import (
	"errors"
	"io/ioutil"

	"github.com/prvst/philosopher/lib/sys"
)

// UnixParameterFile ...
func UnixParameterFile(unixParam string) error {

	param, err := Asset("comet.params")
	if err != nil {
		return errors.New("Cannot deploy Comet parameter file")
	}

	err = ioutil.WriteFile(unixParam, param, sys.FilePermission())
	if err != nil {
		return errors.New("Cannot deploy Comet parameter file")
	}

	return nil
}

// Unix64 ...
func Unix64(unix64 string) error {

	bin, err := Asset("comet.2018014.linux.exe")
	if err != nil {
		return errors.New("Cannot deploy Comet parameter file")
	}

	err = ioutil.WriteFile(unix64, bin, sys.FilePermission())
	if err != nil {
		return errors.New("Cannot deploy Comet parameter file")
	}

	return nil
}
