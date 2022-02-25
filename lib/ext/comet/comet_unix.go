//go:build linux
// +build linux

package comet

import (
	ucomet "philosopher/lib/ext/comet/unix"
)

// Deploy generates comet binary on workdir bin directory
func (c *Comet) Deploy(os, arch string) {
	// deploy comet param file
	ucomet.UnixParameterFile(c.UnixParam)
	c.DefaultParam = c.UnixParam

	// deploy comet
	ucomet.Unix64(c.Unix64)
	c.DefaultBin = c.Unix64
}
