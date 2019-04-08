package obo

import (
	"fmt"
	"io/ioutil"

	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/sys"
)

// Deploy deploys the OBO file to the temp folder
func Deploy(f string) *err.Error {

	asset, e := Asset("unimod.obo")
	if e != nil {
		return &err.Error{Type: err.CannotDeployAsset, Class: err.FATA, Argument: "UniMod Obo not found"}
	}

	e = ioutil.WriteFile(f, asset, sys.FilePermission())
	if e != nil {
		fmt.Println(e.Error())
		return &err.Error{Type: err.CannotDeployAsset, Class: err.FATA, Argument: "Could not deploy UniMod obo"}
	}

	return nil
}
