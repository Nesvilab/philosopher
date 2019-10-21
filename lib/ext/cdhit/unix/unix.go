package cdhit

import (
	"errors"
	"io/ioutil"

	"github.com/nesvilab/philosopher/lib/msg"

	"github.com/nesvilab/philosopher/lib/sys"
)

// Unix64 ...
func Unix64(unix64 string) {

	bin, e := Asset("cd-hit")
	e = ioutil.WriteFile(unix64, bin, sys.FilePermission())

	if e != nil {
		msg.DeployAsset(errors.New("CD-HIT"), "trace")
	}

	return
}
