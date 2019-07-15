package tmtintegrator

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

// TMTIntegrator represents the tool configuration
type TMTIntegrator struct {
	DefaultBin   string
	DefaultParam string
}

// New constructor
func New(temp string) TMTIntegrator {

	var self TMTIntegrator

	self.DefaultBin = ""
	self.DefaultParam = ""

	return self
}

// Run is the TMTIntegrator main entry point
func Run(m met.Data, args []string) (met.Data, *err.Error) {

	var tmti = New(m.Temp)

	// collect and store the mz files
	m.TMTIntegrator.Files = args

	if len(m.TMTIntegrator.Param) > 1 {
		// convert the param file to binary and store it in meta
		var binFile []byte
		paramAbs, _ := filepath.Abs(m.TMTIntegrator.Param)
		binFile, e := ioutil.ReadFile(paramAbs)
		if e != nil {
			logrus.Fatal(e)
		}
		m.TMTIntegrator.ParamFile = binFile
	}

	// run TMTIntegrator
	e := tmti.Execute(m.TMTIntegrator, args)
	if e != nil {
		//logrus.Fatal(e)
	}

	return m, nil
}

// Execute is the main fucntion to execute TMTIntegrator
func (c *TMTIntegrator) Execute(params met.TMTIntegrator, cmdArgs []string) *err.Error {

	cmd := appendParams(params)

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

func appendParams(params met.TMTIntegrator) *exec.Cmd {

	mem := fmt.Sprintf("-Xmx%dG", params.Memmory)
	jarPath, _ := filepath.Abs(params.JarPath)

	args := exec.Command("java",
		mem,
		"-jar",
		jarPath,
		"philosopher.yml",
	)

	return args
}
