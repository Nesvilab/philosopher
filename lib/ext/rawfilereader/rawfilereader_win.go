//go:build windows
// +build windows

package rawfilereader

import (
	wRaw "philosopher/lib/ext/rawfilereader/win"
)

// Deploy generates binaries on workdir
func (c *RawFileReader) Deploy() {
	// deploy windows binary
	wRaw.Win(c.WinBin)
	wRaw.ThermoFisherCommonCoreDataDLL(c.ThermoFisherCommonCoreDataDLL)
	wRaw.ThermoFisherCommonCoreRawFileReaderDLL(c.ThermoFisherCommonCoreRawFileReaderDLL)
	c.DefaultBin = c.WinBin
}
