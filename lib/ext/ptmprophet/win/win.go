package ptmprophet

import (
	"io/ioutil"

	"philosopher/lib/msg"
	"philosopher/lib/sys"
)

// WinPTMProphetParser locates and extracts the PTMProphet binary
func WinPTMProphetParser(s string) {

	bin, e := Asset("PTMProphetParser.exe")
	e = ioutil.WriteFile(s, bin, sys.FilePermission())

	if e != nil {
		msg.ExecutingBinary(e, "trace")
	}

	return
}
