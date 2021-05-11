package rawfilereader

import (
	"errors"
	"io/ioutil"

	"philosopher/lib/msg"

	"philosopher/lib/sys"
)

// Win deploys RawfileReader for Red Hat
func Win(win string) {

	bin, e1 := Asset("RawFileReader.exe")
	if e1 != nil {
		msg.DeployAsset(errors.New("RawFileReader.exe"), "Cannot read rawFileReader.exe obo")
	}

	e2 := ioutil.WriteFile(win, bin, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("rawFileReader.exe"), "Cannot deploy rawFileReader.exe 64-bit")
	}
}

// ThermoFisherCommonCoreDataDLL deploys libgcc_s_dw2.dll
func ThermoFisherCommonCoreDataDLL(s string) {

	bin, e1 := Asset("ThermoFisher.CommonCore.Data.dll")
	if e1 != nil {
		msg.DeployAsset(errors.New("ThermoFisherCommonCoreData"), "Cannot read ThermoFisherCommonCoreData.dll bin")
	}

	e2 := ioutil.WriteFile(s, bin, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("ThermoFisherCommonCoreData"), "Cannot deploy ThermoFisherCommonCoreData.dll.dllcDLL")
	}
}

// ThermoFisherCommonCoreRawFileReaderDLL deploys libgcc_s_dw2.dll
func ThermoFisherCommonCoreRawFileReaderDLL(s string) {

	bin, e1 := Asset("ThermoFisher.CommonCore.RawFileReader.dll")
	if e1 != nil {
		msg.DeployAsset(errors.New("ThermoFisher.CommonCore.RawFileReader.dll"), "Cannot read ThermoFisher.CommonCore.RawFileReader.dll bin")
	}

	e2 := ioutil.WriteFile(s, bin, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("ThermoFisher.CommonCore.RawFileReader.dll"), "Cannot deploy ThermoFisher.CommonCore.RawFileReader.dll")
	}
}
