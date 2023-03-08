package comet

import (
	"errors"
	"os"

	"github.com/Nesvilab/philosopher/lib/msg"
	"github.com/Nesvilab/philosopher/lib/sys"
)

// WinParameterFile writes the parameter file to the disk
func WinParameterFile(winParam string) {

	param, e1 := Asset("comet.params.txt")
	if e1 != nil {
		msg.DeployAsset(errors.New("comet Parameter File"), "cannot read Comet parameter bin")
	}

	e2 := os.WriteFile(winParam, param, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("comet Parameter File"), "cannot deploy Comet parameter")
	}
}

// Win32 deploys win32 bits comet parameter file
func Win32(win32 string) {

	bin, e := Asset("comet.2019011.win32.exe")
	if e != nil {
		msg.DeployAsset(errors.New("comet Windows binary file"), "cannot read Comet bin")
	}

	e = os.WriteFile(win32, bin, sys.FilePermission())
	if e != nil {
		msg.DeployAsset(errors.New("comet Windows binary file"), "cannot deploy Comet")
	}
}

// Win64 deploys win64 bits comet parameter file
func Win64(win64 string) {

	bin, e := Asset("comet.2019011.win64.exe")
	if e != nil {
		msg.DeployAsset(errors.New("comet Windows binary file"), "cannot read Comet bin")
	}

	e = os.WriteFile(win64, bin, sys.FilePermission())
	if e != nil {
		msg.DeployAsset(errors.New("comet Windows binary file"), "cannot deploy Comet")
	}
}
