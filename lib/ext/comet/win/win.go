package comet

import (
	"errors"
	"io/ioutil"

	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/sys"
)

// WinParameterFile writes the parameter file to the disk
func WinParameterFile(winParam string) {

	param, e := Asset("comet.params.txt")
	if e != nil {
		err.DeployAsset(errors.New("Comet parameter file"), "fatal")
	}

	e = ioutil.WriteFile(winParam, param, sys.FilePermission())
	if e != nil {
		err.DeployAsset(errors.New("Comet parameter file"), "fatal")
	}

	return
}

// Win32 deploys win32 bits comt parameter file
func Win32(win32 string) {

	bin, e := Asset("comet.2018014.win32.exe")
	if e != nil {
		err.DeployAsset(errors.New("Comet parameter file"), "fatal")
	}

	e = ioutil.WriteFile(win32, bin, sys.FilePermission())
	if e != nil {
		err.DeployAsset(errors.New("Comet parameter file"), "fatal")
	}

	return
}

// Win64 deploys win64 bits comt parameter file
func Win64(win64 string) {

	bin, e := Asset("comet.2018014.win64.exe")
	e = ioutil.WriteFile(win64, bin, sys.FilePermission())

	if e != nil {
		err.DeployAsset(errors.New("Comet parameter file"), "fatal")
	}

	return
}
