package tmtintegrator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Nesvilab/philosopher/lib/met"
	"github.com/Nesvilab/philosopher/lib/msg"
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
func Run(m met.Data, args []string) met.Data {

	var tmti = New(m.Temp)

	// collect and store the mz files
	m.TMTIntegrator.Files = args

	if len(m.TMTIntegrator.Param) > 1 {
		// convert the param file to binary and store it in meta
		var binFile []byte
		paramAbs, _ := filepath.Abs(m.TMTIntegrator.Param)
		binFile, e := os.ReadFile(paramAbs)
		if e != nil {
			msg.ReadFile(e, "error")
		}
		m.TMTIntegrator.ParamFile = binFile
	}

	// run TMTIntegrator
	tmti.Execute(m.TMTIntegrator, args)

	return m
}

// Execute is the main function to execute TMTIntegrator
func (c *TMTIntegrator) Execute(params met.TMTIntegrator, cmdArgs []string) {

	cmd := appendParams(params)

	for _, i := range cmdArgs {
		file, _ := filepath.Abs(i)
		cmd.Args = append(cmd.Args, file)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	e := cmd.Start()
	if e != nil {
		msg.ExecutingBinary(e, "error")
	}

	_ = cmd.Wait()
}

func appendParams(params met.TMTIntegrator) *exec.Cmd {

	mem := fmt.Sprintf("-Xmx%dG", params.Memory)
	jarPath, _ := filepath.Abs(params.JarPath)

	args := exec.Command("java",
		mem,
		"-jar",
		jarPath,
		"philosopher.yml",
	)

	return args
}
