package ptmprophet

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"philosopher/lib/met"
	"philosopher/lib/msg"
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

// Run PTMProphet
func Run(m met.Data, args []string) met.Data {

	var ptm = New(m.Temp)

	// deploy the binaries
	ptm.Deploy(m.Distro)

	// run
	ptm.Execute(m.PTMProphet, args)

	m.PTMProphet.InputFiles = args

	return m
}

// Execute PTMProphet
func (p *PTMProphet) Execute(params met.PTMProphet, args []string) []string {

	// get the execution commands
	bin := p.DefaultPTMProphetParser
	cmd := exec.Command(bin)

	// append pepxml files
	for i := range args {
		file, _ := filepath.Abs(args[i])
		cmd.Args = append(cmd.Args, args[i])
		cmd.Dir = filepath.Dir(file)
	}

	cmd = p.appendParams(params, cmd)

	// append output file
	var output string
	if params.KeepOld {

		//if len(params.Output) > 0 {
		if args[0] == "interact.pep.xml" {
			output = "interact.mod.pep.xml"
		} else {
			output = strings.Replace(args[0], "pep.xml", "mod.pep.xml", 1)
		}

		cmd.Args = append(cmd.Args, output)
		cmd.Dir = filepath.Dir(output)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	e := cmd.Start()
	if e != nil {
		msg.ExecutingBinary(e, "fatal")
	}
	_ = cmd.Wait()

	// if cmd.ProcessState.ExitCode() != 0 {
	// 	fmt.Println(cmd.Stdout, cmd.Stderr)
	// 	msg.ExecutingBinary(errors.New("there was an error with PTMProphet, please check your parameters and input files"), "fatal")
	// }

	// collect all resulting files
	var customOutput []string
	if params.KeepOld {
		for _, i := range cmd.Args {
			if strings.Contains(i, output) || i == params.Output {
				customOutput = append(customOutput, i)
			}
		}
	}

	if params.KeepOld {
		return customOutput
	}
	return args
}

func (p PTMProphet) appendParams(params met.PTMProphet, cmd *exec.Cmd) *exec.Cmd {

	if params.NoUpdate {
		cmd.Args = append(cmd.Args, "NOUPDATE")
	}

	if params.KeepOld {
		cmd.Args = append(cmd.Args, "KEEPOLD")
	}

	if params.Verbose {
		cmd.Args = append(cmd.Args, "VERBOSE")
	}

	if params.Lability {
		cmd.Args = append(cmd.Args, "LABILITY")
	}

	if params.Ifrags {
		cmd.Args = append(cmd.Args, "IFRAGS")
	}

	if params.Autodirect {
		cmd.Args = append(cmd.Args, "AUTORIDECT")
	}

	if params.MassDiffMode {
		cmd.Args = append(cmd.Args, "MASSDIFFMODE")
	}

	if params.NoMinoFactor {
		cmd.Args = append(cmd.Args, "NOMINOFACTOR")
	}

	if params.Static {
		cmd.Args = append(cmd.Args, "STATIC")
	}

	if params.EM != 2 {
		v := fmt.Sprintf("EM=%d", params.EM)
		cmd.Args = append(cmd.Args, v)
	}

	if params.FragPPMTol != 15 {
		v := fmt.Sprintf("FRAGPPMTOL=%d", params.FragPPMTol)
		cmd.Args = append(cmd.Args, v)
	}

	if params.MaxThreads != 1 {
		v := fmt.Sprintf("MAXTHREADS=%d", params.MaxThreads)
		cmd.Args = append(cmd.Args, v)
	}

	if params.MaxFragZ != 0 {
		v := fmt.Sprintf("MAXFRAGZ=%d", params.MaxFragZ)
		cmd.Args = append(cmd.Args, v)
	}

	if params.Mino != 0 {
		v := fmt.Sprintf("MINO=%d", params.Mino)
		cmd.Args = append(cmd.Args, v)
	}

	if params.MassOffset != 0 {
		v := fmt.Sprintf("MASSOFFSET=%d", params.MassOffset)
		cmd.Args = append(cmd.Args, v)
	}

	if params.PPMTol != 1 {
		v := fmt.Sprintf("PPMTOL=%.4f", params.PPMTol)
		cmd.Args = append(cmd.Args, v)
	}

	if params.MinProb != 0.9 {
		v := fmt.Sprintf("MINPROB=%.4f", params.MinProb)
		cmd.Args = append(cmd.Args, v)
	}

	if params.ExcludeMassDiffMin != 0 {
		v := fmt.Sprintf("EXCLUDEMASSDIFFMIN=%.2f", params.ExcludeMassDiffMin)
		cmd.Args = append(cmd.Args, v)
	}

	if params.ExcludeMassDiffMax != 0 {
		v := fmt.Sprintf("EXCLUDEMASSDIFFMAX=%.2f", params.ExcludeMassDiffMax)
		cmd.Args = append(cmd.Args, v)
	}

	if len(params.NIons) > 0 {
		v := fmt.Sprintf("NIONS=%s", params.NIons)
		cmd.Args = append(cmd.Args, v)
	}

	if len(params.CIons) > 0 {
		v := fmt.Sprintf("CIONS=%s", params.CIons)
		cmd.Args = append(cmd.Args, v)
	}

	if len(params.Mods) > 0 {
		cmd.Args = append(cmd.Args, params.Mods)
	}

	return cmd
}
