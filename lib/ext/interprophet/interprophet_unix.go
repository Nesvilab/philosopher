//go:build linux
// +build linux

package interprophet

import (
	"errors"
	"strings"

	unix "philosopher/lib/ext/interprophet/unix"
	"philosopher/lib/msg"
	"philosopher/lib/sys"
)

// Deploy generates comet binary on workdir bin directory
func (i *InterProphet) Deploy(os, distro string) {

	if strings.EqualFold(distro, sys.Debian()) {
		unix.UnixInterProphetParser(i.UnixInterProphetParser)
		i.DefaultInterProphetParser = i.UnixInterProphetParser
	} else if strings.EqualFold(distro, sys.Redhat()) {
		unix.UnixInterProphetParser(i.UnixInterProphetParser)
		i.DefaultInterProphetParser = i.UnixInterProphetParser
	} else {
		msg.UnsupportedDistribution(errors.New(""), "fatal")
	}
}
