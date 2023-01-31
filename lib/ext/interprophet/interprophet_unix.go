//go:build linux
// +build linux

package interprophet

import (
	"errors"
	"strings"

	unix "github.com/Nesvilab/philosopher/lib/ext/interprophet/unix"
	"github.com/Nesvilab/philosopher/lib/msg"
	"github.com/Nesvilab/philosopher/lib/sys"
)

// Deploy generates comet binary on workdir bin directory
func (i *InterProphet) Deploy(distro string) {

	if strings.EqualFold(distro, sys.Debian()) {
		unix.UnixInterProphetParser(i.UnixInterProphetParser)
		i.DefaultInterProphetParser = i.UnixInterProphetParser
	} else if strings.EqualFold(distro, sys.Redhat()) {
		unix.UnixInterProphetParser(i.UnixInterProphetParser)
		i.DefaultInterProphetParser = i.UnixInterProphetParser
	} else {
		msg.UnsupportedDistribution(errors.New(""), "error")
	}
}
