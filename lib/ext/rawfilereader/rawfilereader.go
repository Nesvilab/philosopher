package rawfilereader

import (
	"os/exec"
	"path/filepath"
	"philosopher/lib/met"
	"philosopher/lib/msg"
	"philosopher/lib/sys"
	"runtime"

	uDeb "philosopher/lib/ext/rawfilereader/deb64"
)

// RawFileReader represents the tool configuration
type RawFileReader struct {
	met.Data
	OS         string
	Arch       string
	Deb64Bin   string
	ReH64Bin   string
	WinBin     string
	DefaultBin string
}

// New constructor
func New() RawFileReader {

	var o RawFileReader
	var m met.Data
	m.Restore(sys.Meta())

	o.UUID = m.UUID
	o.Distro = m.Distro
	o.Home = m.Home
	o.MetaFile = m.MetaFile
	o.MetaDir = m.MetaDir
	o.DB = m.DB
	o.Temp = m.Temp
	o.TimeStamp = m.TimeStamp
	o.OS = m.OS
	o.Arch = m.Arch

	o.OS = runtime.GOOS
	o.Arch = runtime.GOARCH
	o.Deb64Bin = m.Temp + string(filepath.Separator) + "rawFileReaderDeb"
	o.ReH64Bin = m.Temp + string(filepath.Separator) + "rawFileReaderReH"
	//o.WinBin = m.Temp + string(filepath.Separator) + "rawFileReader"

	return o
}

// Run is the main entry point for rawfilereader
func Run(args string) string {

	var reader = New()

	// deploy the binaries
	reader.Deploy()

	// run
	stream := reader.Execute(args)

	return stream
}

// Deploy generates binaries on workdir
func (c *RawFileReader) Deploy() {

	if c.OS == sys.Windows() {

		// deploy cd-hit binary
		//wcdhit.Win64(c.WinBin)
		//c.DefaultBin = c.WinBin

	} else if c.OS == "linux" && c.Distro == sys.Debian() {

		// deploy cd-hit binary
		uDeb.Deb64(c.Deb64Bin)
		c.DefaultBin = c.Deb64Bin

	} else if c.OS == "linux" && c.Distro == sys.Centos() {

		// deploy cd-hit binary
		//uRH.Deb64(c.Deb64Bin)
		//c.DefaultBin = c.Deb64Bin

	}

	return
}

// Execute is the main function to execute RawFileReader
func (c *RawFileReader) Execute(args string) string {

	bin := c.DefaultBin
	cmd := exec.Command(bin)

	file, _ := filepath.Abs(args)
	cmd.Args = append(cmd.Args, file)

	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr

	out, e := cmd.CombinedOutput()
	if e != nil {
		msg.ExecutingBinary(e, "fatal")
	}
	_ = cmd.Wait()

	return string(out)
}
