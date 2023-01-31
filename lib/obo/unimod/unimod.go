package obo

import (
	"errors"
	"io/ioutil"

	"github.com/Nesvilab/philosopher/lib/msg"
	"github.com/Nesvilab/philosopher/lib/sys"
)

// Deploy deploys the Unimod OBO file
func Deploy(f string) {

	asset, e1 := Asset("unimod.obo")
	if e1 != nil {
		msg.DeployAsset(errors.New("Unimod"), "Cannot read unimod obo")
	}

	e2 := ioutil.WriteFile(f, asset, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("Unimod"), "Cannot deploy Unimod obo")
	}

}
