package rawfilereader

import (
	"os/exec"
	"path/filepath"
	"philosopher/lib/met"
	"philosopher/lib/msg"
	"philosopher/lib/sys"
	"runtime"

	uDeb "philosopher/lib/ext/rawfilereader/deb64"
	uRH "philosopher/lib/ext/rawfilereader/reh64"
	wRaw "philosopher/lib/ext/rawfilereader/win"
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

// Deploy generates binaries on workdir
func (c *RawFileReader) Deploy() {

	if c.OS == sys.Windows() {

		// deploy windows binary
		wRaw.Win(c.WinBin)
		wRaw.ThermoFisherCommonCoreDataDLL(c.ThermoFisherCommonCoreDataDLL)
		wRaw.ThermoFisherCommonCoreRawFileReaderDLL(c.ThermoFisherCommonCoreRawFileReaderDLL)
		c.DefaultBin = c.WinBin

	} else if c.OS == "linux" && c.Distro == sys.Debian() {

		// deploy debian binary
		uDeb.Deb64(c.Deb64Bin)
		c.DefaultBin = c.Deb64Bin

	} else {

		// deploy red hat binary
		uRH.Reh64(c.ReH64Bin)
		c.DefaultBin = c.ReH64Bin

	}

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
