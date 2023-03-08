//go:build linux
// +build linux

package ptmprophet

import (
	"errors"
	"strings"

	unix "github.com/Nesvilab/philosopher/lib/ext/ptmprophet/unix"
	"github.com/Nesvilab/philosopher/lib/msg"
	"github.com/Nesvilab/philosopher/lib/sys"
)

// Deploy PTMProphet binaries on binary directory
func (p *PTMProphet) Deploy(distro string) {
	if strings.EqualFold(distro, sys.Debian()) {
		unix.UnixPTMProphetParser(p.UnixPTMProphetParser)
		p.DefaultPTMProphetParser = p.UnixPTMProphetParser
	} else if strings.EqualFold(distro, sys.Redhat()) {
		unix.UnixPTMProphetParser(p.UnixPTMProphetParser)
		p.DefaultPTMProphetParser = p.UnixPTMProphetParser
	} else {
		msg.UnsupportedDistribution(errors.New(""), "error")
	}
}
