//go:build windows
// +build windows

package cdhit

import (
	wcdhit "github.com/Nesvilab/philosopher/lib/ext/cdhit/win"
)

// Deploy generates binaries on workdir
func (c *CDhit) Deploy() {
	// deploy cd-hit binary
	wcdhit.Win64(c.WinBin)
	c.DefaultBin = c.WinBin
}
