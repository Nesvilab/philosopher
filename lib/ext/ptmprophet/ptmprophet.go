package ptmprophet

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	unix "github.com/prvst/philosopher/lib/ext/ptmprophet/unix"
	wPeP "github.com/prvst/philosopher/lib/ext/ptmprophet/win"

	"github.com/prvst/cmsl/err"
	"github.com/prvst/philosopher/lib/meta"
	"github.com/prvst/philosopher/lib/sys"
)

// PTMProphet is the main tool data configuration structure
type PTMProphet struct {
	meta.Data
	Output                  string
	EM                      int
	MzTol                   float64
	PPMTol                  float64
	MinProb                 float64
	NoUpdate                bool
	KeepOld                 bool
	Verbose                 bool
	MassDiffMode            bool
	DefaultPTMProphetParser string
	WinPTMProphetParser     string
	UnixPTMProphetParser    string
}

// New constructor
func New() PTMProphet {

	var o PTMProphet
	var m meta.Data
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

	o.UnixPTMProphetParser = o.Temp + string(filepath.Separator) + "PTMProphetParser"
	o.WinPTMProphetParser = o.Temp + string(filepath.Separator) + "PTMProphetParser.exe"

	return o
}

// Deploy PTMProphet binaries on binary directory
func (c *PTMProphet) Deploy() *err.Error {

	if c.OS == sys.Windows() {
		wPeP.WinPTMProphetParser(c.WinPTMProphetParser)
		c.DefaultPTMProphetParser = c.WinPTMProphetParser
	} else {
		if strings.EqualFold(c.Distro, sys.Debian()) {
			unix.UnixPTMProphetParser(c.UnixPTMProphetParser)
			c.DefaultPTMProphetParser = c.UnixPTMProphetParser
		} else if strings.EqualFold(c.Distro, sys.Redhat()) {
			unix.UnixPTMProphetParser(c.UnixPTMProphetParser)
			c.DefaultPTMProphetParser = c.UnixPTMProphetParser
		} else {
			return &err.Error{Type: err.UnsupportedDistribution, Class: err.FATA, Argument: "PTMProphetParser"}
		}
	}

	return nil
}

// Run PTMProphet
func (c *PTMProphet) Run(args []string) *err.Error {

	// get the execution commands
	bin := c.DefaultPTMProphetParser
	cmd := exec.Command(bin)

	// append pepxml files
	for i := range args {
		file, _ := filepath.Abs(args[i])
		//cmd.Args = append(cmd.Args, file)
		cmd.Args = append(cmd.Args, args[i])
		cmd.Dir = filepath.Dir(file)
	}

	cmd = c.appendParams(cmd)

	// append output file
	var output string
	output = "interact.mod.pep.xml"
	if len(c.Output) > 0 {
		output = fmt.Sprintf("%s.pep.xml", c.Output)
	}
	cmd.Args = append(cmd.Args, output)
	cmd.Dir = filepath.Dir(output)
	// var output string
	// if len(c.Output) > 0 {
	// 	output = fmt.Sprintf("%s%s%s.mod.pep.xml", c.Temp, string(filepath.Separator), c.Output)
	// 	output, _ = filepath.Abs(output)
	// 	cmd.Args = append(cmd.Args, output)
	// 	cmd.Dir = filepath.Dir(output)
	// }

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	e := cmd.Start()
	if e != nil {
		return &err.Error{Type: err.CannotExecuteBinary, Class: err.FATA, Argument: "PTMprophet"}
	}
	_ = cmd.Wait()

	// var baseDir string
	// baseDir = filepath.Dir(args[0])
	//
	// // copy to work directory
	// if len(c.Output) > 0 {
	// 	dest := fmt.Sprintf("%s%s%s", baseDir, string(filepath.Separator), filepath.Base(output))
	// 	e = sys.CopyFile(output, dest)
	// 	if e != nil {
	// 		return &err.Error{Type: err.CannotCopyFile, Class: err.FATA, Argument: "PTMProphet results"}
	// 	}
	// }

	return nil
}

func (c *PTMProphet) appendParams(cmd *exec.Cmd) *exec.Cmd {

	if c.NoUpdate == true {
		cmd.Args = append(cmd.Args, "NOUPDATE")
	}

	if c.KeepOld == true {
		cmd.Args = append(cmd.Args, "KEEPOLD")
	}

	if c.Verbose == true {
		cmd.Args = append(cmd.Args, "VERBOSE")
	}

	if c.MassDiffMode == true {
		cmd.Args = append(cmd.Args, "MASSDIFFMODE")
	}

	if c.EM != 1 {
		v := fmt.Sprintf("EM=%d", c.EM)
		cmd.Args = append(cmd.Args, v)
	}

	if c.MzTol != 0.1 {
		v := fmt.Sprintf("MZTOL=%.4f", c.MzTol)
		cmd.Args = append(cmd.Args, v)
	}

	if c.PPMTol != 1 {
		v := fmt.Sprintf("PPMTOL=%.4f", c.PPMTol)
		cmd.Args = append(cmd.Args, v)
	}

	if c.MinProb != 0 {
		v := fmt.Sprintf("MINPROB=%.4f", c.MinProb)
		cmd.Args = append(cmd.Args, v)
	}

	return cmd
}
