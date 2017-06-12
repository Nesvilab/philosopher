package proteinprophet

import (
	"errors"
	"io/ioutil"
)

// WinBatchCoverage ...
func WinBatchCoverage(s string) error {

	bin, err := Asset("batchcoverage.exe")
	err = ioutil.WriteFile(s, bin, 0755)

	if err != nil {
		return errors.New("Cannot deploy batchcoverage")
	}

	return nil
}

// WinDatabaseParser ...
func WinDatabaseParser(s string) error {

	bin, err := Asset("DatabaseParser.exe")
	err = ioutil.WriteFile(s, bin, 0755)

	if err != nil {
		return errors.New("Cannot deploy DatabaseParser")
	}

	return nil
}

// WinProteinProphet ...
func WinProteinProphet(s string) error {

	bin, err := Asset("ProteinProphet.exe")
	err = ioutil.WriteFile(s, bin, 0755)

	if err != nil {
		return errors.New("Cannot deploy ProteinProphet")
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
