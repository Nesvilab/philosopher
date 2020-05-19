package comet

import (
	"errors"
	"io/ioutil"

	"philosopher/lib/msg"
	"philosopher/lib/sys"
)

// UnixParameterFile deploys Comet parameter file
func UnixParameterFile(unixParam string) {

	param, e1 := Asset("comet.params")
	if e1 != nil {
		msg.DeployAsset(errors.New("Comet Parameter File"), "Cannot read Comet parameter bin")
	}

	e2 := ioutil.WriteFile(unixParam, param, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("Comet Parameter File"), "Cannot deploy Comet parameter")
	}

	return
}

// Unix64 deploys Comet binary
func Unix64(unix64 string) {

	bin, e1 := Asset("comet.2019011.linux.exe")
	if e1 != nil {
		msg.DeployAsset(errors.New("Comet Linux binary"), "Cannot read Comet parameter bin")
	} else {

		e2 := ioutil.WriteFile(unix64, bin, sys.FilePermission())
		if e2 != nil {
			msg.DeployAsset(errors.New("Comet Linux binary"), "Cannot deploy Comet binary")
		}
	}

	return
}
