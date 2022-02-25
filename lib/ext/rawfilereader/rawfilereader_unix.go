//go:build linux
// +build linux

package rawfilereader

import (
	uDeb "philosopher/lib/ext/rawfilereader/deb64"
	uRH "philosopher/lib/ext/rawfilereader/reh64"
	"philosopher/lib/sys"
)

// Deploy generates binaries on workdir
func (c *RawFileReader) Deploy() {
	if c.OS == "linux" && c.Distro == sys.Debian() {

		// deploy debian binary
		uDeb.Deb64(c.Deb64Bin)
		c.DefaultBin = c.Deb64Bin

	} else {

		// deploy red hat binary
		uRH.Reh64(c.ReH64Bin)
		c.DefaultBin = c.ReH64Bin

	}
}
