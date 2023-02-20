package peptideprophet

import (
	"errors"
	"os"

	"github.com/Nesvilab/philosopher/lib/msg"
	"github.com/Nesvilab/philosopher/lib/sys"
)

// WinInteractParser deploys InteractParser.exe
func WinInteractParser(s string) {

	bin, e1 := Asset("InteractParser.exe")
	if e1 != nil {
		msg.DeployAsset(errors.New("InteractParser"), "Cannot read InteractParser bin")
	} else {

		e2 := os.WriteFile(s, bin, sys.FilePermission())
		if e2 != nil {
			msg.DeployAsset(errors.New("InteractParser"), "Cannot deploy InteractParser")
		}
	}

	return
}

// WinRefreshParser deploys Refreshparser.exe
func WinRefreshParser(s string) {

	bin, e1 := Asset("RefreshParser.exe")
	if e1 != nil {
		msg.DeployAsset(errors.New("RefreshParser"), "Cannot read RefreshParser bin")
	}

	e2 := os.WriteFile(s, bin, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("RefreshParser"), "Cannot deploy RefreshParser")
	}

	return
}

// WinPeptideProphetParser deploys Windows PeptideProphetParser
func WinPeptideProphetParser(s string) {

	bin, e1 := Asset("PeptideProphetParser.exe")
	if e1 != nil {
		msg.DeployAsset(errors.New("PeptideProphetParser"), "Cannot read PeptideProphet bin")
	}

	e2 := os.WriteFile(s, bin, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("PeptideProphetParser"), "Cannot deploy PeptideProphet")
	}

	return
}

// Mv deploys mv.exe
func Mv(s string) {

	bin, e1 := Asset("mv.exe")
	if e1 != nil {
		msg.DeployAsset(errors.New("mv.exe"), "Cannot read mv.exe bin")
	}

	e2 := os.WriteFile(s, bin, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("mv.exe"), "Cannot deploy mv.exe")
	}

	return
}

// LibgccDLL deploys libgcc_s_dw2.dll
func LibgccDLL(s string) {

	bin, e1 := Asset("libgcc_s_dw2-1.dll")
	if e1 != nil {
		msg.DeployAsset(errors.New("LibgccDLL"), "Cannot read LibgccDLL bin")
	}

	e2 := os.WriteFile(s, bin, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("LibgccDLL"), "Cannot deploy LibgccDLL")
	}

	return
}

// Zlib1DLL deploys zlib1.dll
func Zlib1DLL(s string) {

	bin, e1 := Asset("zlib1.dll")
	if e1 != nil {
		msg.DeployAsset(errors.New("Zlib1DLL"), "Cannot read Zlib1DLL bin")
	}

	e2 := os.WriteFile(s, bin, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("Zlib1DLL"), "Cannot deploy Zlib1DLL")
	}

	return
}
