package interprophet

import (
	"errors"
	"io/ioutil"

	"github.com/prvst/philosopher/lib/sys"
)

// WinInterProphetParser accessor
func WinInterProphetParser(s string) error {

	bin, err := Asset("InterProphetParser.exe")
	err = ioutil.WriteFile(s, bin, sys.FilePermission())

	if err != nil {
		return errors.New("Cannot deploy InterProphetParser")
	}

	return nil
}

// LibgccDLL accessor
func LibgccDLL(s string) error {

	bin, err := Asset("libgcc_s_dw2-1.dll")
	err = ioutil.WriteFile(s, bin, sys.FilePermission())

	if err != nil {
		return errors.New("Cannot deploy libgcc_s_dw2")
	}

	return nil
}

// Zlib1DLL accessor
func Zlib1DLL(s string) error {

	bin, err := Asset("zlib1.dll")
	err = ioutil.WriteFile(s, bin, sys.FilePermission())

	if err != nil {
		return errors.New("Cannot deploy Zlib1DLL")
	}

	return nil
}
