//go:build linux
// +build linux

package peptideprophet

import (
	"errors"
	"strings"

	unix "philosopher/lib/ext/peptideprophet/unix"
	"philosopher/lib/msg"
	"philosopher/lib/sys"
)

// Deploy PeptideProphet binaries on binary directory
func (p *PeptideProphet) Deploy(os, distro string) {

	if strings.EqualFold(distro, sys.Debian()) {
		unix.UnixInteractParser(p.UnixInteractParser)
		p.DefaultInteractParser = p.UnixInteractParser
		unix.UnixRefreshParser(p.UnixRefreshParser)
		p.DefaultRefreshParser = p.UnixRefreshParser
		unix.UnixPeptideProphetParser(p.UnixPeptideProphetParser)
		p.DefaultPeptideProphetParser = p.UnixPeptideProphetParser
	} else if strings.EqualFold(distro, sys.Redhat()) {
		unix.UnixInteractParser(p.UnixInteractParser)
		p.DefaultInteractParser = p.UnixInteractParser
		unix.UnixRefreshParser(p.UnixRefreshParser)
		p.DefaultRefreshParser = p.UnixRefreshParser
		unix.UnixPeptideProphetParser(p.UnixPeptideProphetParser)
		p.DefaultPeptideProphetParser = p.UnixPeptideProphetParser
	} else {
		msg.UnsupportedDistribution(errors.New(""), "fatal")
	}
}
