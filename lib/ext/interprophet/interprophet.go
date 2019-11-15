package interprophet

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	unix "philosopher/lib/ext/interprophet/unix"
	wiPr "philosopher/lib/ext/interprophet/win"
	"philosopher/lib/met"
	"philosopher/lib/msg"
	"philosopher/lib/sys"
)

// InterProphet represents the tool configuration
type InterProphet struct {
	DefaultInterProphetParser string
	WinInterProphetParser     string
	UnixInterProphetParser    string
	LibgccDLL                 string
	Zlib1DLL                  string
}

// New constructor
func New(temp string) InterProphet {

	var self InterProphet

	self.UnixInterProphetParser = temp + string(filepath.Separator) + "InterProphetParser"
	self.WinInterProphetParser = temp + string(filepath.Separator) + "InterProphetParser.exe"
	self.LibgccDLL = temp + string(filepath.Separator) + "libgcc_s_dw2-1.dll"
	self.Zlib1DLL = temp + string(filepath.Separator) + "zlib1.dll"

	return self
}

// Run is the main entry point for InterProphet
func Run(m met.Data, args []string) met.Data {

	var itp = New(m.Temp)
	m.InterProphet.InputFiles = args

	if len(args) < 1 {
		msg.NoParametersFound(errors.New("IProphet input files"), "fatal")
	}

	// deploy the binaries
	itp.Deploy(m.OS, m.Distro)

	// run InterProphet
	itp.Execute(m.InterProphet, m.Home, m.Temp, args)

	m.InterProphet.InputFiles = args

	return m
}

// Deploy generates comet binary on workdir bin directory
func (i *InterProphet) Deploy(os, distro string) {

	if os == sys.Windows() {
		wiPr.WinInterProphetParser(i.WinInterProphetParser)
		i.DefaultInterProphetParser = i.WinInterProphetParser
		wiPr.LibgccDLL(i.LibgccDLL)
		wiPr.Zlib1DLL(i.Zlib1DLL)
	} else {
		if strings.EqualFold(distro, sys.Debian()) {
			unix.UnixInterProphetParser(i.UnixInterProphetParser)
			i.DefaultInterProphetParser = i.UnixInterProphetParser
		} else if strings.EqualFold(distro, sys.Redhat()) {
			unix.UnixInterProphetParser(i.UnixInterProphetParser)
			i.DefaultInterProphetParser = i.UnixInterProphetParser
		} else {
			msg.UnsupportedDistribution(errors.New(""), "fatal")
		}
	}

	return
}

// Execute IProphet
func (i InterProphet) Execute(params met.InterProphet, home, temp string, args []string) []string {

	// run
	bin := i.DefaultInterProphetParser
	cmd := exec.Command(bin)

	for i := 0; i <= len(args)-1; i++ {
		file, _ := filepath.Abs(args[i])
		cmd.Args = append(cmd.Args, file)
	}

	// append output file
	output := fmt.Sprintf("%s%s%s.pep.xml", temp, string(filepath.Separator), params.Output)
	output, _ = filepath.Abs(output)

	cmd = i.appendParams(params, cmd)
	cmd.Args = append(cmd.Args, output)
	cmd.Dir = filepath.Dir(output)

	env := os.Environ()
	env = append(env, fmt.Sprintf("XML_ONLY=%d", 1))
	env = append(env, fmt.Sprintf("WEBSERVER_ROOT=%s", temp))
	for i := range env {
		if strings.HasPrefix(strings.ToUpper(env[i]), "PATH=") {
			env[i] = env[i] + ";" + temp
		}
	}
	cmd.Env = env

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	e := cmd.Start()
	if e != nil {
		msg.ExecutingBinary(e, "fatal")
	}
	_ = cmd.Wait()

	if cmd.ProcessState.ExitCode() != 0 {
		msg.ExecutingBinary(errors.New("There was an error with iProphet, please check your parameters and input files"), "fatal")
	}

	// copy to work directory
	dest := fmt.Sprintf("%s%s%s", home, string(filepath.Separator), filepath.Base(output))
	sys.CopyFile(output, dest)

	// collect all resulting files
	var processedOutput []string
	for _, i := range cmd.Args {
		if strings.Contains(i, params.Output) || i == params.Output {
			processedOutput = append(processedOutput, i)
		}
	}

	return processedOutput
}

func (i InterProphet) appendParams(params met.InterProphet, cmd *exec.Cmd) *exec.Cmd {

	if params.Length == true {
		cmd.Args = append(cmd.Args, "LENGTH")
	}

	if params.Nofpkm == true {
		cmd.Args = append(cmd.Args, "NOFPKM")
	}

	if params.Nonss == true {
		cmd.Args = append(cmd.Args, "NONSS")
	}

	if params.Nonse == true {
		cmd.Args = append(cmd.Args, "NONSE")
	}

	if params.Nonrs == true {
		cmd.Args = append(cmd.Args, "NONRS")
	}

	if params.Nonsm == true {
		cmd.Args = append(cmd.Args, "NONSM")
	}

	if params.Nonsp == true {
		cmd.Args = append(cmd.Args, "NONSP")
	}

	if params.Sharpnse == true {
		cmd.Args = append(cmd.Args, "SHARPNSE")
	}

	if params.Nonsi == true {
		cmd.Args = append(cmd.Args, "NONSI")
	}

	if params.Threads != 1 {
		v := fmt.Sprintf("THREADS=%d", params.Threads)
		cmd.Args = append(cmd.Args, v)
	}

	if len(params.Decoy) > 0 {
		v := fmt.Sprintf("DECOY=%s", params.Decoy)
		cmd.Args = append(cmd.Args, v)
	}

	if len(params.Cat) > 0 {
		v := fmt.Sprintf("CAT=%s", params.Cat)
		cmd.Args = append(cmd.Args, v)
	}

	if params.MinProb != 0 {
		v := fmt.Sprintf("MINPROB=%.4f", params.MinProb)
		cmd.Args = append(cmd.Args, v)
	}

	return cmd
}
