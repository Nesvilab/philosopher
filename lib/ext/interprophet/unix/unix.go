package interprophet

import (
	"errors"
	"io/ioutil"

	"github.com/prvst/philosopher/lib/err"

	"github.com/prvst/philosopher/lib/sys"
)

// UnixInterProphetParser accessor
func UnixInterProphetParser(s string) {

	bin, e := Asset("InterProphetParser")
	e = ioutil.WriteFile(s, bin, sys.FilePermission())

	if e != nil {
		err.DeployAsset(errors.New("InterProphetParser"))
	}

	return
}
