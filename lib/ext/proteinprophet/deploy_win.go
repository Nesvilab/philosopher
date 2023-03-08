//go:build windows
// +build windows

package proteinprophet

import (
	wPoP "github.com/Nesvilab/philosopher/lib/ext/proteinprophet/win"
)

// Deploy generates comet binary on workdir bin directory
func (p *ProteinProphet) Deploy(distro string) {

	wPoP.WinBatchCoverage(p.WinBatchCoverage)
	p.DefaultBatchCoverage = p.WinBatchCoverage
	wPoP.WinDatabaseParser(p.WinDatabaseParser)
	p.DefaultDatabaseParser = p.WinDatabaseParser
	wPoP.WinProteinProphet(p.WinProteinProphet)
	p.DefaultProteinProphet = p.WinProteinProphet
	wPoP.LibgccDLL(p.LibgccDLL)
	wPoP.Zlib1DLL(p.Zlib1DLL)
}
