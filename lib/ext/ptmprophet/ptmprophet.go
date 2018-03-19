package ptmprophet

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/prvst/philosopher/lib/err"
	unix "github.com/prvst/philosopher/lib/ext/ptmprophet/unix"
	wPeP "github.com/prvst/philosopher/lib/ext/ptmprophet/win"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/sys"
)

// PTMProphet is the main tool data configuration structure
type PTMProphet struct {
	DefaultPTMProphetParser string
	WinPTMProphetParser     string
	UnixPTMProphetParser    string
}

// New constructor
func New(temp string) PTMProphet {

	var self PTMProphet

	//temp, _ := sys.GetTemp()

	self.UnixPTMProphetParser = temp + string(filepath.Separator) + "PTMProphetParser"
	self.WinPTMProphetParser = temp + string(filepath.Separator) + "PTMProphetParser.exe"

	return self
}

func Run(m met.Data, args []string) met.Data {

	var ptm = New(m.Temp)

	// deploy the binaries
	e := ptm.Deploy(m.OS, m.Distro)
	if e != nil {
		fmt.Println(e.Message)
	}

	// run
	xml, e := ptm.Execute(m.PTMProphet, args)
	if e != nil {
		fmt.Println(e.Message)
	}
	_ = xml

	m.PTMProphet.InputFiles = args

	return m
}

// Deploy PTMProphet binaries on binary directory
func (p *PTMProphet) Deploy(os, distro string) *err.Error {

	if os == sys.Windows() {
		wPeP.WinPTMProphetParser(p.WinPTMProphetParser)
		p.DefaultPTMProphetParser = p.WinPTMProphetParser
	} else {
		if strings.EqualFold(distro, sys.Debian()) {
			unix.UnixPTMProphetParser(p.UnixPTMProphetParser)
			p.DefaultPTMProphetParser = p.UnixPTMProphetParser
		} else if strings.EqualFold(distro, sys.Redhat()) {
			unix.UnixPTMProphetParser(p.UnixPTMProphetParser)
			p.DefaultPTMProphetParser = p.UnixPTMProphetParser
		} else {
			return &err.Error{Type: err.UnsupportedDistribution, Class: err.FATA, Argument: "PTMProphetParser"}
		}
	}

	return nil
}

// Execute PTMProphet
func (p *PTMProphet) Execute(params met.PTMProphet, args []string) ([]string, *err.Error) {

	// get the execution commands
	bin := p.DefaultPTMProphetParser
	cmd := exec.Command(bin)

	// append pepxml files
	for i := range args {
		file, _ := filepath.Abs(args[i])
		//cmd.Args = append(cmd.Args, file)
		cmd.Args = append(cmd.Args, args[i])
		cmd.Dir = filepath.Dir(file)
	}

	cmd = p.appendParams(params, cmd)

	// append output file
	var output string
	output = "interact.mod.pep.xml"
	if len(params.Output) > 0 {
		output = fmt.Sprintf("%s.pep.xml", params.Output)
	}
	cmd.Args = append(cmd.Args, output)
	cmd.Dir = filepath.Dir(output)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	e := cmd.Start()
	if e != nil {
		return nil, &err.Error{Type: err.CannotExecuteBinary, Class: err.FATA, Argument: "PTMprophet"}
	}
	_ = cmd.Wait()

	// collect all resulting files
	var processedOutput []string
	for _, i := range cmd.Args {
		if strings.Contains(i, output) || i == params.Output {
			processedOutput = append(processedOutput, i)
		}
	}

	return processedOutput, nil
}

func (p PTMProphet) appendParams(params met.PTMProphet, cmd *exec.Cmd) *exec.Cmd {

	if params.NoUpdate == true {
		cmd.Args = append(cmd.Args, "NOUPDATE")
	}

	if params.KeepOld == true {
		cmd.Args = append(cmd.Args, "KEEPOLD")
	}

	if params.Verbose == true {
		cmd.Args = append(cmd.Args, "VERBOSE")
	}

	if params.MassDiffMode == true {
		cmd.Args = append(cmd.Args, "MASSDIFFMODE")
	}

	if params.EM != 1 {
		v := fmt.Sprintf("EM=%d", params.EM)
		cmd.Args = append(cmd.Args, v)
	}

	if params.MzTol != 0.1 {
		v := fmt.Sprintf("MZTOL=%.4f", params.MzTol)
		cmd.Args = append(cmd.Args, v)
	}

	if params.PPMTol != 1 {
		v := fmt.Sprintf("PPMTOL=%.4f", params.PPMTol)
		cmd.Args = append(cmd.Args, v)
	}

	if params.MinProb != 0 {
		v := fmt.Sprintf("MINPROB=%.4f", params.MinProb)
		cmd.Args = append(cmd.Args, v)
	}

	if len(params.Mods) > 0 {
		cmd.Args = append(cmd.Args, params.Mods)
	}

	return cmd
}
