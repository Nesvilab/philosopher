//go:build windows
// +build windows

package comet

import (
	wcomet "github.com/Nesvilab/philosopher/lib/ext/comet/win"
	"github.com/Nesvilab/philosopher/lib/sys"
)

// Deploy generates comet binary on workdir bin directory
func (c *Comet) Deploy(arch string) {
	// deploy comet param file
	wcomet.WinParameterFile(c.WinParam)
	c.DefaultParam = c.WinParam

	if arch == sys.Arch386() {
		wcomet.Win32(c.Win32)
		c.DefaultBin = c.Win32

	} else {
		wcomet.Win64(c.Win64)
		c.DefaultBin = c.Win64
	}
}
