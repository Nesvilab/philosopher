package msconvert

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/prvst/philosopher/lib/err"
	dmsc "github.com/prvst/philosopher/lib/ext/msconvert/darwin"
	umsc "github.com/prvst/philosopher/lib/ext/msconvert/unix"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/sys"
)

// Msconvert represents the tool configuration
type Msconvert struct {
	DefaultBin string
	Unix64     string
	Darwinx64  string
}

// New constructor
func New(temp string) Msconvert {

	var self Msconvert

	self.DefaultBin = ""
	self.Unix64 = temp + string(filepath.Separator) + "msconvert"
	self.Darwinx64 = temp + string(filepath.Separator) + "msconvert"

	return self
}

// Run is the Msconvert main entry point
func Run(m met.Data, args []string) (met.Data, *err.Error) {

	var msc = New(m.Temp)

	if len(args) < 1 {
		return m, &err.Error{Type: err.CannotRunComet, Class: err.FATA, Argument: "Missing parameter file or data file for analysis"}
	}

	// deploy msconvert
	msc.Deploy(m.OS, m.Arch)

	// run msconvert
	e := msc.Execute(args, m.Msconvert)
	if e != nil {
		//logrus.Fatal(e)
	}

	return m, nil
}

// Deploy generates comet binary on workdir bin directory
func (c *Msconvert) Deploy(os, arch string) {

	if os == sys.Darwin() {

		dmsc.Darwinx64(c.Darwinx64)
		c.DefaultBin = c.Darwinx64

	} else if os == sys.Linux() {

		// deploy msconvert
		umsc.Unix64(c.Unix64)
		c.DefaultBin = c.Unix64

	} else {

	}

	return
}

// Execute is the main fucntion to execute Msconvert
func (c *Msconvert) Execute(cmdArgs []string, param met.Msconvert) *err.Error {

	run := exec.Command(c.DefaultBin)
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr
	e := run.Start()
	if e != nil {
		//return &err.Error{Type: err.CannotRunComet, Class: err.FATA}
		return nil
	}
	_ = run.Wait()

	return nil
}
