package interprophet

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/prvst/philosopher/lib/err"
	unix "github.com/prvst/philosopher/lib/ext/interprophet/unix"
	wiPr "github.com/prvst/philosopher/lib/ext/interprophet/win"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/sys"
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
func New() InterProphet {

	var self InterProphet

	temp, _ := sys.GetTemp()

	self.UnixInterProphetParser = temp + string(filepath.Separator) + "InterProphetParser"
	self.WinInterProphetParser = temp + string(filepath.Separator) + "InterProphetParser.exe"
	self.LibgccDLL = temp + string(filepath.Separator) + "libgcc_s_dw2-1.dll"
	self.Zlib1DLL = temp + string(filepath.Separator) + "zlib1.dll"

	return self
}

// Deploy generates comet binary on workdir bin directory
func (i *InterProphet) Deploy(os, distro string) *err.Error {

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
			return &err.Error{Type: err.UnsupportedDistribution, Class: err.FATA, Argument: "dont know how to deploy InterProphet"}
		}
	}

	return nil
}

// Run IProphet ...
func (i InterProphet) Run(params met.InterProphet, home, temp string, args []string) ([]string, *err.Error) {

	// run
	bin := i.DefaultInterProphetParser
	cmd := exec.Command(bin)

	for i := 0; i <= len(args)-1; i++ {
		file, _ := filepath.Abs(args[i])
		cmd.Args = append(cmd.Args, file)
	}

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

	cmd.Args = append(cmd.Args, params.Output)

	cmd.Dir = filepath.Dir(args[0])

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
		return nil, &err.Error{Type: err.CannotExecuteBinary, Class: err.FATA, Argument: "InterProphet"}
	}
	_ = cmd.Wait()

	// collect all resulting files
	var output []string
	for _, i := range cmd.Args {
		if strings.Contains(i, "iproph") || i == params.Output {
			output = append(output, i)
		}
	}

	return output, nil
}
