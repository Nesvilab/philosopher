package interprophet

import (
	"errors"
	"io/ioutil"

	"github.com/nesvilab/philosopher/lib/msg"

	"github.com/nesvilab/philosopher/lib/sys"
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
