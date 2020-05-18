package ptmprophet

import (
	"errors"
	"io/ioutil"

	"philosopher/lib/msg"
	"philosopher/lib/sys"
)

// UnixPTMProphetParser locates and extracts the PTMProphet binary
func UnixPTMProphetParser(s string) {

	bin, e1 := Asset("PTMProphetParser")
	if e1 != nil {
		msg.DeployAsset(errors.New("PTMProphetParser"), "Cannot read PTMProphetParser binary")
	}

	e2 := ioutil.WriteFile(s, bin, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("PTMProphetParser"), "Cannot deploy PTMProphetParser")
	}

	return
}
