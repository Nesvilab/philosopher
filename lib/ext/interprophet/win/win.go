package interprophet

import (
	"errors"
	"io/ioutil"

	"github.com/prvst/philosopher/lib/msg"

	"github.com/prvst/philosopher/lib/sys"
)

// WinInterProphetParser accessor
func WinInterProphetParser(s string) {

	bin, e := Asset("InterProphetParser.exe")
	e = ioutil.WriteFile(s, bin, sys.FilePermission())

	if e != nil {
		msg.DeployAsset(errors.New("InterProphetParser"), "trace")
	}

	return
}

// LibgccDLL accessor
func LibgccDLL(s string) {

	bin, e := Asset("libgcc_s_dw2-1.dll")
	e = ioutil.WriteFile(s, bin, sys.FilePermission())

	if e != nil {
		msg.DeployAsset(errors.New("libgcc_s_dw2"), "trace")
	}

	return
}

// Zlib1DLL accessor
func Zlib1DLL(s string) {

	bin, e := Asset("zlib1.dll")
	e = ioutil.WriteFile(s, bin, sys.FilePermission())

	if e != nil {
		msg.DeployAsset(errors.New("Zlib1DLL"), "trace")
	}

	return
}
