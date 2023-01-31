package cdhit

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/Nesvilab/philosopher/lib/msg"

	"github.com/Nesvilab/philosopher/lib/met"
	"github.com/Nesvilab/philosopher/lib/sys"
)

// CDhit represents the tool configuration
type CDhit struct {
	met.Data
	OS           string
	Arch         string
	UnixBin      string
	WinBin       string
	DefaultBin   string
	FileName     string
	FastaDB      string
	ClusterFile  string
	ClusterFasta string
}

// New constructor
func New() CDhit {

	var o CDhit
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

	o.OS = runtime.GOOS
	o.Arch = runtime.GOARCH
	o.UnixBin = m.Temp + string(filepath.Separator) + "cd-hit"
	o.WinBin = m.Temp + string(filepath.Separator) + "cd-hit.exe"
	o.FastaDB = m.Temp + string(filepath.Separator) + o.UUID + ".fasta"
	o.FileName = m.Temp + string(filepath.Separator) + "cdhit"

	return o
}

// Run runs the cdhit binary with user's information
func (c *CDhit) Run(level float64) {

	l := strconv.FormatFloat(level, 'E', -1, 64)

	cmd := c.DefaultBin
	args := []string{"-i", c.DB, "-o", c.ClusterFasta, "-c", l}

	run := exec.Command(cmd, args...)
	e := run.Start()
	_ = run.Wait()

	if e != nil {
		msg.ExecutingBinary(e, "fatal")
	}

}
