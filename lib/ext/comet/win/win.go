package comet

import (
	"errors"
	"io/ioutil"

	"philosopher/lib/msg"
	"philosopher/lib/sys"
)

// WinParameterFile writes the parameter file to the disk
func WinParameterFile(winParam string) {

	param, e1 := Asset("comet.params.txt")
	if e1 != nil {
		msg.DeployAsset(errors.New("Comet Parameter File"), "Cannot read Comet parameter bin")
	}

	e2 := ioutil.WriteFile(winParam, param, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("Comet Parameter File"), "Cannot deploy Comet parameter")
	}

	return
}

// Win32 deploys win32 bits comt parameter file
func Win32(win32 string) {

	bin, e := Asset("comet.2019011.win32.exe")
	if e != nil {
		msg.DeployAsset(errors.New("Comet Windows binary file"), "fatal")
	}

	e = ioutil.WriteFile(win32, bin, sys.FilePermission())
	if e != nil {
		msg.DeployAsset(errors.New("Comet Windows binary file"), "fatal")
	}

	return
}

// Win64 deploys win64 bits comt parameter file
func Win64(win64 string) {

	bin, e := Asset("comet.2019011.win64.exe")
	e = ioutil.WriteFile(win64, bin, sys.FilePermission())

	if e != nil {
		msg.DeployAsset(errors.New("Comet Windows binary file"), "fatal")
	}

	return
}
