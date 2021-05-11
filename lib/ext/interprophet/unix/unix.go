package interprophet

import (
	"errors"
	"io/ioutil"

	"philosopher/lib/msg"

	"philosopher/lib/sys"
)

// UnixInterProphetParser accessor
func UnixInterProphetParser(s string) {

	bin, e1 := Asset("InterProphetParser")
	if e1 != nil {
		msg.DeployAsset(errors.New("InterProphetParser"), "Cannot read InterProphetParser bin")
	}

	e2 := ioutil.WriteFile(s, bin, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("InterProphetParser"), "Cannot deploy InterProphetParser")
	}
}
