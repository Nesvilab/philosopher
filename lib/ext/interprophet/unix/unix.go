package interprophet

import (
	"errors"
	"io/ioutil"

	"philosopher/lib/msg"

	"philosopher/lib/sys"
)

// UnixInterProphetParser accessor
func UnixInterProphetParser(s string) {

	bin, e := Asset("InterProphetParser")
	e = ioutil.WriteFile(s, bin, sys.FilePermission())

	if e != nil {
		msg.DeployAsset(errors.New("InterProphetParser"), "trace")
	}

	return
}
