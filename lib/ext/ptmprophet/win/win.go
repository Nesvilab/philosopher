package ptmprophet

import (
	"io/ioutil"

	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/sys"
)

// WinPTMProphetParser locates and extracts the PTMProphet binary
func WinPTMProphetParser(s string) {

	bin, e := Asset("PTMProphetParser.exe")
	e = ioutil.WriteFile(s, bin, sys.FilePermission())

	if e != nil {
		err.ExecutingBinary(e, "trace")
	}

	return
}
