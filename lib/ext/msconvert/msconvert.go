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

	// TODO create an error class for msconvert
	if len(args) < 1 {
		return m, &err.Error{Type: err.CannotRunProgram, Class: err.FATA, Argument: "[msconvert] Missing files for conversion"}
	}

	// deploy msconvert
	msc.Deploy(m.OS, m.Arch)

	// run msconvert
	e := msc.Execute(m.Msconvert, m.Home, m.Temp, args)
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
func (c *Msconvert) Execute(params met.Msconvert, home, temp string, args []string) *err.Error {

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
		return &err.Error{Type: err.CannotExecuteBinary, Class: err.FATA, Argument: "ProteinProphet"}
	}
	_ = cmd.Wait()

	return nil
}

func (c Msconvert) appendParams(params met.Msconvert, cmd *exec.Cmd) *exec.Cmd {

	if params.NoIndex == true {
		cmd.Args = append(cmd.Args, "--noindex")
	}

	if params.Zlib == true {
		cmd.Args = append(cmd.Args, "--zlib")
	}

	if params.Format == "mzML" || params.Format == "mzml" {
		cmd.Args = append(cmd.Args, "--mzML")
	}

	if params.Format == "mzXML" || params.Format == "mzxml" {
		cmd.Args = append(cmd.Args, "--mzXML")
	}

	if params.Format == "mz5" {
		cmd.Args = append(cmd.Args, "--mz5")
	}

	if params.Format == "mgf" {
		cmd.Args = append(cmd.Args, "--mgf")
	}

	if params.Format == "text" {
		cmd.Args = append(cmd.Args, "--text")
	}

	if params.Format == "ms1" {
		cmd.Args = append(cmd.Args, "--ms1")
	}

	if params.Format == "cms1" {
		cmd.Args = append(cmd.Args, "--cms1")
	}

	if params.Format == "ms2" {
		cmd.Args = append(cmd.Args, "--ms2")
	}

	if params.Format == "cms2" {
		cmd.Args = append(cmd.Args, "--cms2")
	}

	if params.MZBinaryEncoding == "64" {
		cmd.Args = append(cmd.Args, "--64")
	}

	if params.MZBinaryEncoding == "32" {
		cmd.Args = append(cmd.Args, "--32")
	}

	if params.IntensityBinaryEncoding == "64" {
		cmd.Args = append(cmd.Args, "--inten64")
	}

	if params.IntensityBinaryEncoding == "32" {
		cmd.Args = append(cmd.Args, "--inten32")
	}

	// if len(params.Output) > 0 {
	// 	v := fmt.Sprintf("--outfile%s", params.Output)
	// 	cmd.Args = append(cmd.Args, v)
	// }

	return cmd
}
