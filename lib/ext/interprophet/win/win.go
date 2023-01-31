package interprophet

import (
	"errors"
	"io/ioutil"

	"github.com/Nesvilab/philosopher/lib/msg"

	"github.com/Nesvilab/philosopher/lib/sys"
)

// WinInterProphetParser accessor
func WinInterProphetParser(s string) {

	bin, e1 := Asset("InterProphetParser.exe")
	if e1 != nil {
		msg.DeployAsset(errors.New("InterProphetParser.exe"), "Cannot read InterProphetParser.exe bin")
	}

	e2 := ioutil.WriteFile(s, bin, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("InterProphetParser.exe"), "Cannot deploy InterProphetParser.exe")
	}

	return
}

// LibgccDLL accessor
func LibgccDLL(s string) {

	bin, e1 := Asset("libgcc_s_dw2-1.dll")
	if e1 != nil {
		msg.DeployAsset(errors.New("LibgccDLL"), "Cannot read LibgccDLL bin")
	}

	e2 := ioutil.WriteFile(s, bin, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("LibgccDLL"), "Cannot deploy LibgccDLL")
	}

	return
}

// Zlib1DLL accessor
func Zlib1DLL(s string) {

	bin, e1 := Asset("zlib1.dll")
	if e1 != nil {
		msg.DeployAsset(errors.New("Zlib1DLL"), "Cannot read Zlib1DLL bin")
	}

	e2 := ioutil.WriteFile(s, bin, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("Zlib1DLL"), "Cannot deploy Zlib1DLL")
	}

	return
}
