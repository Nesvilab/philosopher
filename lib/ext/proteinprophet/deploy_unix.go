//go:build linux
// +build linux

package proteinprophet

import (
	"errors"
	unix "philosopher/lib/ext/proteinprophet/unix"
	"philosopher/lib/msg"
	"philosopher/lib/sys"
	"strings"
)

// Deploy generates comet binary on workdir bin directory
func (p *ProteinProphet) Deploy(distro string) {

	if strings.EqualFold(distro, sys.Debian()) {
		unix.UnixBatchCoverage(p.UnixBatchCoverage)
		p.DefaultBatchCoverage = p.UnixBatchCoverage
		unix.UnixDatabaseParser(p.UnixDatabaseParser)
		p.DefaultDatabaseParser = p.UnixDatabaseParser
		unix.UnixProteinProphet(p.UnixProteinProphet)
		p.DefaultProteinProphet = p.UnixProteinProphet
	} else if strings.EqualFold(distro, sys.Redhat()) {
		unix.UnixBatchCoverage(p.UnixBatchCoverage)
		p.DefaultBatchCoverage = p.UnixBatchCoverage
		unix.UnixDatabaseParser(p.UnixDatabaseParser)
		p.DefaultDatabaseParser = p.UnixDatabaseParser
		unix.UnixProteinProphet(p.UnixProteinProphet)
		p.DefaultProteinProphet = p.UnixProteinProphet
	} else {
		msg.UnsupportedDistribution(errors.New(""), "error")
	}
}
