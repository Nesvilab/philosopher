package rawfilereader

import (
	"os/exec"
	"path/filepath"
	"philosopher/lib/met"
	"philosopher/lib/msg"
	"philosopher/lib/sys"
	"runtime"
)

// RawFileReader represents the tool configuration
type RawFileReader struct {
	met.Data
	OS                                     string
	Arch                                   string
	Deb64Bin                               string
	ReH64Bin                               string
	WinBin                                 string
	DefaultBin                             string
	ThermoFisherCommonCoreDataDLL          string
	ThermoFisherCommonCoreRawFileReaderDLL string
}

// New constructor
func New() RawFileReader {

	var self RawFileReader
	var m met.Data
	m.Restore(sys.Meta())

	self.UUID = m.UUID
	self.Distro = m.Distro
	self.Home = m.Home
	self.MetaFile = m.MetaFile
	self.MetaDir = m.MetaDir
	self.DB = m.DB
	self.Temp = m.Temp
	self.TimeStamp = m.TimeStamp
	self.OS = m.OS
	self.Arch = m.Arch

	self.OS = runtime.GOOS
	self.Arch = runtime.GOARCH
	self.Deb64Bin = m.Temp + string(filepath.Separator) + "rawFileReaderDeb"
	self.ReH64Bin = m.Temp + string(filepath.Separator) + "rawFileReaderReH"
	self.WinBin = m.Temp + string(filepath.Separator) + "RawFileReader.exe"
	self.ThermoFisherCommonCoreDataDLL = m.Temp + string(filepath.Separator) + "ThermoFisher.CommonCore.Data.dll"
	self.ThermoFisherCommonCoreRawFileReaderDLL = m.Temp + string(filepath.Separator) + "ThermoFisher.CommonCore.RawFileReader.dll"

	return self
}

// Run is the main entry point for rawfilereader
func Run(rawFileName, scanQuery string) string {

	var reader = New()

	// deploy the binaries
	reader.Deploy()

	// run
	stream := reader.Execute(rawFileName, scanQuery)

	return stream
}

// Execute is the main function to execute RawFileReader
func (c *RawFileReader) Execute(rawFileName, scanQuery string) string {

	bin := c.DefaultBin
	cmd := exec.Command(bin)

	file, _ := filepath.Abs(rawFileName)
	cmd.Args = append(cmd.Args, file)

	if len(scanQuery) > 0 {
		cmd.Args = append(cmd.Args, scanQuery)
	}

	out, e := cmd.CombinedOutput()
	if e != nil {
		msg.ExecutingBinary(e, "fatal")
	}
	_ = cmd.Wait()

	return string(out)
}
