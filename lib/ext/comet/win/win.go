package comet

import (
	"errors"
	"io/ioutil"

	"github.com/prvst/philosopher/lib/sys"
)

// WinParameterFile ...
func WinParameterFile(winParam string) error {

	param, err := Asset("comet.params.txt")
	if err != nil {
		return errors.New("Cannot deploy Comet parameter file")
	}

	err = ioutil.WriteFile(winParam, param, sys.FilePermission())
	if err != nil {
		return errors.New("Cannot deploy Comet parameter file")
	}

	return nil
}

// Win32 ...
func Win32(win32 string) error {

	bin, err := Asset("comet.2018014.win32.exe")
	if err != nil {
		return errors.New("Cannot deploy Comet parameter file")
	}

	err = ioutil.WriteFile(win32, bin, sys.FilePermission())
	if err != nil {
		return errors.New("Cannot deploy Comet parameter file")
	}

	return nil
}

// Win64 ...
func Win64(win64 string) error {

	bin, err := Asset("comet.2018014.win64.exe")
	err = ioutil.WriteFile(win64, bin, sys.FilePermission())

	if err != nil {
		return errors.New("Cannot deploy Comet parameter file")
	}

	return nil
}
