package comet

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/prvst/philosopher/lib/err"
	ucomet "github.com/prvst/philosopher/lib/ext/comet/unix"
	wcomet "github.com/prvst/philosopher/lib/ext/comet/win"
	"github.com/prvst/philosopher/lib/meta"
	"github.com/prvst/philosopher/lib/sys"
)

// Comet represents the tool configuration
type Comet struct {
	meta.Data
	OS           string
	Arch         string
	DefaultBin   string
	DefaultParam string
	Win32        string
	Win64        string
	Unix64       string
	WinParam     string
	UnixParam    string
	Param        string
	Print        bool
}

// New constructor
func New() Comet {

	var o Comet
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

	o.Win32 = o.Temp + string(filepath.Separator) + "comet.2016012.win32.exe"
	o.Win64 = o.Temp + string(filepath.Separator) + "comet.2016012.win64.exe"
	o.WinParam = o.Temp + string(filepath.Separator) + "comet.params.txt"
	o.Unix64 = o.Temp + string(filepath.Separator) + "comet.2016012.linux.exe"
	o.UnixParam = o.Temp + string(filepath.Separator) + "comet.params"

	return o
}

// Deploy generates comet binary on workdir bin directory
func (c *Comet) Deploy() {

	if c.OS == sys.Windows() {

		// deploy comet param file
		wcomet.WinParameterFile(c.WinParam)
		c.DefaultParam = c.WinParam

		if c.Arch == sys.Arch386() {
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

		// // deploy comet
		ucomet.Unix64(c.Unix64)
		c.DefaultBin = c.Unix64

	}

	return
}

// Run is the main fucntion to execute Comet
func (c *Comet) Run(cmdArgs []string) *err.Error {

	param := fmt.Sprintf("-P%s", c.Param)
	args := []string{param}

	for _, i := range cmdArgs {
		file, _ := filepath.Abs(i)
		args = append(args, file)
	}

	run := exec.Command(c.DefaultBin, args...)
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr
	e := run.Start()
	if e != nil {
		return &err.Error{Type: err.CannotRunComet, Class: err.FATA}
	}
	_ = run.Wait()

	return nil
}
