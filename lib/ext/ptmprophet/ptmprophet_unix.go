//go:build linux
// +build linux

package ptmprophet

import (
	"errors"
	"strings"

	unix "philosopher/lib/ext/ptmprophet/unix"
	"philosopher/lib/msg"
	"philosopher/lib/sys"
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
