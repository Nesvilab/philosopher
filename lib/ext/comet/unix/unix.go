package comet

import (
	"errors"
	"os"

	"github.com/Nesvilab/philosopher/lib/msg"
	"github.com/Nesvilab/philosopher/lib/sys"
)

// UnixParameterFile deploys Comet parameter file
func UnixParameterFile(unixParam string) {

	param, e1 := Asset("comet.params")
	if e1 != nil {
		msg.DeployAsset(errors.New("comet Parameter File"), "Cannot read Comet parameter bin")
	}

	e2 := os.WriteFile(unixParam, param, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("comet Parameter File"), "Cannot deploy Comet parameter")
	}

}

// Unix64 deploys Comet binary
func Unix64(unix64 string) {

	bin, e1 := Asset("comet.2019011.linux.exe")
	if e1 != nil {
		msg.DeployAsset(errors.New("comet Linux binary"), "Cannot read Comet parameter bin")
	} else {

		e2 := os.WriteFile(unix64, bin, sys.FilePermission())
		if e2 != nil {
			msg.DeployAsset(errors.New("comet Linux binary"), "Cannot deploy Comet binary")
		}
	}

}
