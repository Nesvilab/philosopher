package peptideprophet

import (
	"errors"
	"io/ioutil"
)

// WinInteractParser ...
func WinInteractParser(s string) error {

	bin, err := Asset("InteractParser.exe")
	err = ioutil.WriteFile(s, bin, 0755)

	if err != nil {
		return errors.New("Cannot deploy InteractParser")
	}

	return nil
}

// WinRefreshParser ...
func WinRefreshParser(s string) error {

	bin, err := Asset("RefreshParser.exe")
	err = ioutil.WriteFile(s, bin, 0755)

	if err != nil {
		return errors.New("Cannot deploy RefreshParser")
	}

	return nil
}

// WinPeptideProphetParser ...
func WinPeptideProphetParser(s string) error {

	bin, err := Asset("PeptideProphetParser.exe")
	err = ioutil.WriteFile(s, bin, 0755)

	if err != nil {
		return errors.New("Cannot deploy PeptideProphetParser")
	}

	return nil
}

// Mv ...
func Mv(s string) error {

	bin, err := Asset("mv.exe")
	err = ioutil.WriteFile(s, bin, 0755)

	if err != nil {
		return errors.New("Cannot deploy mv.exe")
	}

	return nil
}

// LibgccDLL ...
func LibgccDLL(s string) error {

	bin, err := Asset("libgcc_s_dw2-1.dll")
	err = ioutil.WriteFile(s, bin, 0755)

	if err != nil {
		return errors.New("Cannot deploy libgcc_s_dw2")
	}

	return nil
}

// Zlib1DLL ...
func Zlib1DLL(s string) error {

	bin, err := Asset("zlib1.dll")
	err = ioutil.WriteFile(s, bin, 0755)

	if err != nil {
		return errors.New("Cannot deploy Zlib1DLL")
	}

	return nil
}
