package ptmprophet

import (
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

	return nil
}
