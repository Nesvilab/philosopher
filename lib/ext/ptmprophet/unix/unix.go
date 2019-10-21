package ptmprophet

import (
	"io/ioutil"

	"github.com/nesvilab/philosopher/lib/msg"
	"github.com/nesvilab/philosopher/lib/sys"
)

// UnixPTMProphetParser locates and extracts the PTMProphet binary
func UnixPTMProphetParser(s string) {

	bin, e := Asset("PTMProphetParser")
	e = ioutil.WriteFile(s, bin, sys.FilePermission())

	if e != nil {
		msg.ExecutingBinary(e, "trace")
	}

	return
}
