package peptideprophet

import (
	"errors"
	"io/ioutil"

	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/sys"
)

// WinInteractParser deploys InteractParser.exe
func WinInteractParser(s string) {

	bin, e := Asset("InteractParser.exe")
	e = ioutil.WriteFile(s, bin, sys.FilePermission())

	if e != nil {
		err.DeployAsset(errors.New("InteractParser"), "trace")
	}

	return
}

// WinRefreshParser deploys Refreshparser.exe
func WinRefreshParser(s string) {

	bin, e := Asset("RefreshParser.exe")
	e = ioutil.WriteFile(s, bin, sys.FilePermission())

	if e != nil {
		err.DeployAsset(errors.New("RefreshParser"), "trace")
	}

	return
}

// WinPeptideProphetParser deploys Windows PeptideProphetParser
func WinPeptideProphetParser(s string) {

	bin, e := Asset("PeptideProphetParser.exe")
	e = ioutil.WriteFile(s, bin, sys.FilePermission())

	if e != nil {
		err.DeployAsset(errors.New("PeptideProphetParser"), "trace")
	}

	return
}

// Mv deploys mv.exe
func Mv(s string) {

	bin, e := Asset("mv.exe")
	e = ioutil.WriteFile(s, bin, sys.FilePermission())

	if e != nil {
		err.DeployAsset(errors.New("mv.exe"), "trace")
	}

	return
}

// LibgccDLL deploys libgcc_s_dw2.dll
func LibgccDLL(s string) {

	bin, e := Asset("libgcc_s_dw2-1.dll")
	e = ioutil.WriteFile(s, bin, sys.FilePermission())

	if e != nil {
		err.DeployAsset(errors.New("libgcc_s_dw2"), "trace")
	}

	return
}

// Zlib1DLL deploys zlib1.dll
func Zlib1DLL(s string) {

	bin, e := Asset("zlib1.dll")
	e = ioutil.WriteFile(s, bin, sys.FilePermission())

	if e != nil {
		err.DeployAsset(errors.New("Zlib1DLL"), "trace")
	}

	return
}
