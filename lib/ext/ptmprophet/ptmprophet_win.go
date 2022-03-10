//go:build windows
// +build windows

package ptmprophet

import (
	wPeP "philosopher/lib/ext/ptmprophet/win"
)

// Deploy PTMProphet binaries on binary directory
func (p *PTMProphet) Deploy(distro string) {
	wPeP.WinPTMProphetParser(p.WinPTMProphetParser)
	p.DefaultPTMProphetParser = p.WinPTMProphetParser
}
