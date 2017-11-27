package interprophet

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	unix "github.com/prvst/philosopher/lib/ext/interprophet/unix"
	wiPr "github.com/prvst/philosopher/lib/ext/interprophet/win"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/sys"
)

// InterProphet represents the tool configuration
type InterProphet struct {
	met.Data
	DefaultInterProphetParser string
	WinInterProphetParser     string
	UnixInterProphetParser    string
	LibgccDLL                 string
	Zlib1DLL                  string
	Threads                   int
	Decoy                     string
	Cat                       string
	MinProb                   float64
	Output                    string
	Length                    bool
	Nofpkm                    bool
	Nonss                     bool
	Nonse                     bool
	Nonrs                     bool
	Nonsm                     bool
	Nonsp                     bool
	Sharpnse                  bool
	Nonsi                     bool
}

// New constructor
func New() InterProphet {

	var o InterProphet
	var m met.Data
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

	o.UnixInterProphetParser = o.Temp + string(filepath.Separator) + "InterProphetParser"
	o.WinInterProphetParser = o.Temp + string(filepath.Separator) + "InterProphetParser.exe"
	o.LibgccDLL = o.Temp + string(filepath.Separator) + "libgcc_s_dw2-1.dll"
	o.Zlib1DLL = o.Temp + string(filepath.Separator) + "zlib1.dll"

	return o
}

// Deploy generates comet binary on workdir bin directory
func (c *InterProphet) Deploy() error {

	if c.OS == sys.Windows() {
		wiPr.WinInterProphetParser(c.WinInterProphetParser)
		c.DefaultInterProphetParser = c.WinInterProphetParser
		wiPr.LibgccDLL(c.LibgccDLL)
		wiPr.Zlib1DLL(c.Zlib1DLL)
	} else {
		if strings.EqualFold(c.Distro, sys.Debian()) {
			unix.UnixInterProphetParser(c.UnixInterProphetParser)
			c.DefaultInterProphetParser = c.UnixInterProphetParser
		} else if strings.EqualFold(c.Distro, sys.Redhat()) {
			unix.UnixInterProphetParser(c.UnixInterProphetParser)
			c.DefaultInterProphetParser = c.UnixInterProphetParser
		} else {
			return errors.New("Unsupported distribution for InterProphet")
		}
	}

	return nil
}

// Run IProphet ...
func (c *InterProphet) Run(args []string) error {

	// run
	bin := c.DefaultInterProphetParser
	cmd := exec.Command(bin)

	if len(args) < 1 {
		return errors.New("You need to provide a pepXML file")
	}

	for i := 0; i <= len(args)-1; i++ {
		file, _ := filepath.Abs(args[i])
		cmd.Args = append(cmd.Args, file)
	}

	if c.Length == true {
		cmd.Args = append(cmd.Args, "LENGTH")
	}

	if c.Nofpkm == true {
		cmd.Args = append(cmd.Args, "NOFPKM")
	}

	if c.Nonss == true {
		cmd.Args = append(cmd.Args, "NONSS")
	}

	if c.Nonse == true {
		cmd.Args = append(cmd.Args, "NONSE")
	}

	if c.Nonrs == true {
		cmd.Args = append(cmd.Args, "NONRS")
	}

	if c.Nonsm == true {
		cmd.Args = append(cmd.Args, "NONSM")
	}

	if c.Nonsp == true {
		cmd.Args = append(cmd.Args, "NONSP")
	}

	if c.Sharpnse == true {
		cmd.Args = append(cmd.Args, "SHARPNSE")
	}

	if c.Nonsi == true {
		cmd.Args = append(cmd.Args, "NONSI")
	}

	if c.Threads != 1 {
		v := fmt.Sprintf("THREADS=%d", c.Threads)
		cmd.Args = append(cmd.Args, v)
	}

	if len(c.Decoy) > 0 {
		v := fmt.Sprintf("DECOY=%s", c.Decoy)
		cmd.Args = append(cmd.Args, v)
	}

	if len(c.Cat) > 0 {
		v := fmt.Sprintf("CAT=%s", c.Cat)
		cmd.Args = append(cmd.Args, v)
	}

	if c.MinProb != 0 {
		v := fmt.Sprintf("MINPROB=%.4f", c.MinProb)
		cmd.Args = append(cmd.Args, v)
	}

	cmd.Args = append(cmd.Args, c.Output)

	cmd.Dir = filepath.Dir(args[0])

	env := os.Environ()
	env = append(env, fmt.Sprintf("XML_ONLY=%d", 1))
	env = append(env, fmt.Sprintf("WEBSERVER_ROOT=%s", c.Temp))
	for i := range env {
		if strings.HasPrefix(strings.ToUpper(env[i]), "PATH=") {
			env[i] = env[i] + ";" + c.Temp
		}
	}
	cmd.Env = env

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return errors.New("Cannot run iProphet")
	}
	_ = cmd.Wait()

	return nil
}
