package idconvert

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/prvst/philosopher/lib/err"
	didc "github.com/prvst/philosopher/lib/ext/idconvert/darwin"
	uidcv "github.com/prvst/philosopher/lib/ext/idconvert/unix"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/sys"
)

// Idconvert represents the tool configuration
type Idconvert struct {
	DefaultBin string
	Unix64     string
	Darwinx64  string
}

// New constructor
func New(temp string) Idconvert {

	var self Idconvert

	self.DefaultBin = ""
	self.Unix64 = temp + string(filepath.Separator) + "idconvert"
	self.Darwinx64 = temp + string(filepath.Separator) + "idconvert"

	return self
}

// Run is the Msconvert main entry point
func Run(m met.Data, args []string) (met.Data, *err.Error) {

	var msc = New(m.Temp)

	// TODO create an error class for msconvert
	if len(args) < 1 {
		return m, &err.Error{Type: err.CannotRunProgram, Class: err.FATA, Argument: "[idconvert] Missing files for conversion"}
	}

	// deploy msconvert
	msc.Deploy(m.OS, m.Arch)

	// run msconvert
	e := msc.Execute(m.Idconvert, m.Home, m.Temp, args)
	if e != nil {
		//logrus.Fatal(e)
	}

	return m, nil
}

// Deploy generates comet binary on workdir bin directory
func (c *Idconvert) Deploy(os, arch string) {

	if os == sys.Darwin() {

		didc.Darwinx64(c.Darwinx64)
		c.DefaultBin = c.Darwinx64

	} else if os == sys.Linux() {

		// deploy msconvert
		uidcv.Unix64(c.Unix64)
		c.DefaultBin = c.Unix64

	} else {

	}

	return
}

// Execute is the main fucntion to execute Msconvert
func (c *Idconvert) Execute(params met.Idconvert, home, temp string, args []string) *err.Error {

	bin := c.DefaultBin
	cmd := exec.Command(bin)

	// append files
	for i := range args {
		file, _ := filepath.Abs(args[i])
		cmd.Args = append(cmd.Args, file)
	}

	cmd = c.appendParams(params, cmd)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	e := cmd.Start()
	if e != nil {
		return &err.Error{Type: err.CannotExecuteBinary, Class: err.FATA, Argument: "Idconvert"}
	}
	_ = cmd.Wait()

	return nil
}

func (c Idconvert) appendParams(params met.Idconvert, cmd *exec.Cmd) *exec.Cmd {

	if params.Format == "pepXML" || params.Format == "pepxml" {
		cmd.Args = append(cmd.Args, "--pepXML")
	}

	if params.Format == "mzIdentML" || params.Format == "mzidentml" {
		cmd.Args = append(cmd.Args, "--mzIdentML")
	}

	if params.Format == "text" {
		cmd.Args = append(cmd.Args, "--text")
	}

	return cmd
}
