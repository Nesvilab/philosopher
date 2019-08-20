package obo

import (
	"fmt"
	"io/ioutil"

	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/sys"
)

// Deploy deploys the OBO file to the temp folder
func Deploy(f string) {

	asset, e := Asset("unimod.obo")
	if e != nil {
		err.DeployAsset(e, "trace")
	}

	e = ioutil.WriteFile(f, asset, sys.FilePermission())
	if e != nil {
		fmt.Println(e.Error())
		err.DeployAsset(e, "trace")
	}

	return
}
