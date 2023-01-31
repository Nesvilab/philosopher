package proteinprophet

import (
	"errors"
	"io/ioutil"

	"github.com/Nesvilab/philosopher/lib/msg"
	"github.com/Nesvilab/philosopher/lib/sys"
)

// WinBatchCoverage deploys batchcoverage
func WinBatchCoverage(s string) {

	bin, e1 := Asset("batchcoverage.exe")
	if e1 != nil {
		msg.DeployAsset(errors.New("batchcoverage"), "Cannot read batchcoverage bin")
	}

	e2 := ioutil.WriteFile(s, bin, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("batchcoverage"), "Cannot deploy batchcoverage")
	}

	return
}

// WinDatabaseParser deploys DatabaseParser
func WinDatabaseParser(s string) {

	bin, e1 := Asset("DatabaseParser.exe")
	if e1 != nil {
		msg.DeployAsset(errors.New("DatabaseParser"), "Cannot read DatabaseParser bin")
	}

	e2 := ioutil.WriteFile(s, bin, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("DatabaseParser"), "Cannot deploy DatabaseParser")
	}

	return
}

// WinProteinProphet deploys ProteinProphet.exe
func WinProteinProphet(s string) {

	bin, e1 := Asset("ProteinProphet.exe")
	if e1 != nil {
		msg.DeployAsset(errors.New("ProteinProphet"), "Cannot read ProteinProphet bin")
	}

	e2 := ioutil.WriteFile(s, bin, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("ProteinProphet"), "Cannot deploy ProteinProphet")
	}

	return
}

// LibgccDLL deploys libgcc_s_dw2.dll
func LibgccDLL(s string) {

	bin, e1 := Asset("libgcc_s_dw2-1.dll")
	if e1 != nil {
		msg.DeployAsset(errors.New("libgcc_s_dw2"), "Cannot read LibgccDLL bin")
	}

	e2 := ioutil.WriteFile(s, bin, 0755)
	if e2 != nil {
		msg.DeployAsset(errors.New("libgcc_s_dw2"), "Cannot deploy LibgccDLL")
	}

	return
}

// Zlib1DLL deploys zlib1.dll
func Zlib1DLL(s string) {

	bin, e1 := Asset("zlib1.dll")
	if e1 != nil {
		msg.DeployAsset(errors.New("Zlib1DLL"), "Cannot read Zlib1DLL bin")
	}

	e2 := ioutil.WriteFile(s, bin, 0755)
	if e2 != nil {
		msg.DeployAsset(errors.New("Zlib1DLL"), "Cannot deploy Zlib1DLL")
	}

	return
}
