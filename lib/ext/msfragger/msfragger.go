package msfragger

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/met"
	"github.com/sirupsen/logrus"
)

// MSFragger represents the tool configuration
type MSFragger struct {
	DefaultBin   string
	DefaultParam string
}

// New constructor
func New(temp string) MSFragger {

	var self MSFragger

	self.DefaultBin = ""
	self.DefaultParam = ""

	return self
}

// Run is the Fragger main entry point
func Run(m met.Data, args []string) (met.Data, *err.Error) {

	var frg = New(m.Temp)

	// if len(m.MSFragger.Param) < 1 {
	// 	return m, &err.Error{Type: err.CannotRunMSFragger, Class: err.WARN, Argument: "No parameter file found, using values defined via command line"}
	// }

	// collect and store the mz files
	m.MSFragger.RawFiles = args

	if len(m.MSFragger.Param) > 1 {
		// convert the param file to binary and store it in meta
		var binFile []byte
		paramAbs, _ := filepath.Abs(m.MSFragger.Param)
		binFile, e := ioutil.ReadFile(paramAbs)
		if e != nil {
			logrus.Fatal(e)
		}
		m.MSFragger.ParamFile = binFile
	}

	// run comet
	e := frg.Execute(args, m.MSFragger)
	if e != nil {
		//logrus.Fatal(e)
	}

	return m, nil
}

// Execute is the main fucntion to execute MSFragger
func (c *MSFragger) Execute(cmdArgs []string, m met.MSFragger) *err.Error {

	mem := fmt.Sprintf("-Xmx%sG", m.Memmory)
	cmd := exec.Command("java", "-jar", mem, m.JarPath, m.Param)

	for _, i := range cmdArgs {
		file, _ := filepath.Abs(i)
		cmd.Args = append(cmd.Args, file)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	e := cmd.Start()
	if e != nil {
		return nil
	}

	_ = cmd.Wait()

	return nil
}
