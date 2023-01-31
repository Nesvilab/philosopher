package ptmprophet

import (
	"errors"
	"io/ioutil"

	"github.com/Nesvilab/philosopher/lib/msg"
	"github.com/Nesvilab/philosopher/lib/sys"
)

// WinPTMProphetParser locates and extracts the PTMProphet binary
func WinPTMProphetParser(s string) {

	bin, e1 := Asset("PTMProphetParser.exe")
	if e1 != nil {
		msg.DeployAsset(errors.New("PTMProphetParser"), "Cannot read PTMProphet bin")
	}

	e2 := ioutil.WriteFile(s, bin, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("PTMProphetParser"), "Cannot deploy PTMProphet")
	}
}
