//go:build windows
// +build windows

package peptideprophet

import (
	wPeP "github.com/Nesvilab/philosopher/lib/ext/peptideprophet/win"
)

// Deploy PeptideProphet binaries on binary directory
func (p *PeptideProphet) Deploy(distro string) {
	wPeP.WinInteractParser(p.WinInteractParser)
	p.DefaultInteractParser = p.WinInteractParser
	wPeP.WinRefreshParser(p.WinRefreshParser)
	p.DefaultRefreshParser = p.WinRefreshParser
	wPeP.WinPeptideProphetParser(p.WinPeptideProphetParser)
	p.DefaultPeptideProphetParser = p.WinPeptideProphetParser
	wPeP.LibgccDLL(p.LibgccDLL)
	wPeP.Zlib1DLL(p.Zlib1DLL)
	wPeP.Mv(p.Mv)
}
