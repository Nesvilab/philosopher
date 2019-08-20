package proteinprophet

import (
	"errors"
	"io/ioutil"

	"github.com/prvst/philosopher/lib/err"
)

// WinBatchCoverage deploys batchcoverage
func WinBatchCoverage(s string) {

	bin, e := Asset("batchcoverage.exe")
	e = ioutil.WriteFile(s, bin, 0755)

	if e != nil {
		err.DeployAsset(errors.New("batchcoverage"), "trace")
	}

	return
}

// WinDatabaseParser deploys DatabaseParser
func WinDatabaseParser(s string) {

	bin, e := Asset("DatabaseParser.exe")
	e = ioutil.WriteFile(s, bin, 0755)

	if e != nil {
		err.DeployAsset(errors.New("DatabaseParser"), "trace")
	}

	return
}

// WinProteinProphet deploys ProteinProphet.exe
func WinProteinProphet(s string) {

	bin, e := Asset("ProteinProphet.exe")
	e = ioutil.WriteFile(s, bin, 0755)

	if e != nil {
		err.DeployAsset(errors.New("ProteinProphet"), "trace")
	}

	return
}

// LibgccDLL deploys libgcc_s_dw2.dll
func LibgccDLL(s string) {

	bin, e := Asset("libgcc_s_dw2-1.dll")
	e = ioutil.WriteFile(s, bin, 0755)

	if e != nil {
		err.DeployAsset(errors.New("libgcc_s_dw2"), "trace")
	}

	return
}

// Zlib1DLL deploys zlib1.dll
func Zlib1DLL(s string) {

	bin, e := Asset("zlib1.dll")
	e = ioutil.WriteFile(s, bin, 0755)

	if e != nil {
		err.DeployAsset(errors.New("zlib1"), "trace")
	}

	return
}
