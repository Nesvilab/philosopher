package ptmprophet

import (
	"io/ioutil"

	"philosopher/lib/msg"
	"philosopher/lib/sys"
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
