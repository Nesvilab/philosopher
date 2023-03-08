//go:build linux
// +build linux

package cdhit

import (
	ucdhit "github.com/Nesvilab/philosopher/lib/ext/cdhit/unix"
)

// Deploy generates binaries on workdir
func (c *CDhit) Deploy() {
	// deploy cd-hit binary
	ucdhit.Unix64(c.UnixBin)
	c.DefaultBin = c.UnixBin
}
