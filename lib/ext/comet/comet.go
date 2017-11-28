package comet

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/prvst/philosopher/lib/err"
	ucomet "github.com/prvst/philosopher/lib/ext/comet/unix"
	wcomet "github.com/prvst/philosopher/lib/ext/comet/win"
	"github.com/prvst/philosopher/lib/sys"
)

// Comet represents the tool configuration
type Comet struct {
	DefaultBin   string
	DefaultParam string
	Win32        string
	Win64        string
	Unix64       string
	WinParam     string
	UnixParam    string
}

// New constructor
func New() Comet {

	var self Comet

	temp, _ := sys.GetTemp()

	self.DefaultBin = ""
	self.DefaultParam = ""
	self.Win32 = temp + string(filepath.Separator) + "comet.2017012.win32.exe"
	self.Win64 = temp + string(filepath.Separator) + "comet.2017012.win64.exe"
	self.Unix64 = temp + string(filepath.Separator) + "comet.2017012.linux.exe"
	self.WinParam = temp + string(filepath.Separator) + "comet.params.txt"
	self.UnixParam = temp + string(filepath.Separator) + "comet.params"

	return self
}

// Deploy generates comet binary on workdir bin directory
func (c *Comet) Deploy(os, arch string) {

	if os == sys.Windows() {

		// deploy comet param file
		wcomet.WinParameterFile(c.WinParam)
		c.DefaultParam = c.WinParam

		if arch == sys.Arch386() {
			wcomet.Win32(c.Win32)
			c.DefaultBin = c.Win32

		} else {
			wcomet.Win64(c.Win64)
			c.DefaultBin = c.Win64
		}

	} else {

		// deploy comet param file
		ucomet.UnixParameterFile(c.UnixParam)
		c.DefaultParam = c.UnixParam

		// deploy comet
		ucomet.Unix64(c.Unix64)
		c.DefaultBin = c.Unix64

	}

	return
}

// Run is the main fucntion to execute Comet
func (c *Comet) Run(cmdArgs []string, param string) *err.Error {

	par := fmt.Sprintf("-P%s", param)
	args := []string{par}

	for _, i := range cmdArgs {
		file, _ := filepath.Abs(i)
		args = append(args, file)
	}

	run := exec.Command(c.DefaultBin, args...)
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr
	e := run.Start()
	if e != nil {
		//return &err.Error{Type: err.CannotRunComet, Class: err.FATA}
		return nil
	}
	_ = run.Wait()

	return nil
}
