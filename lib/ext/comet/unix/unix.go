package comet

import (
	"errors"
	"io/ioutil"

	"github.com/prvst/philosopher/lib/msg"

	"github.com/prvst/philosopher/lib/sys"
)

// UnixParameterFile ...
func UnixParameterFile(unixParam string) {

	param, e := Asset("comet.params")
	if e != nil {
		msg.DeployAsset(errors.New("Comet parameter file"), "fatal")
	}

	e = ioutil.WriteFile(unixParam, param, sys.FilePermission())
	if e != nil {
		msg.DeployAsset(errors.New("Comet parameter file"), "fatal")
	}

	return
}

// Unix64 ...
func Unix64(unix64 string) {

	bin, e := Asset("comet.2018014.linux.exe")
	if e != nil {
		msg.DeployAsset(errors.New("Comet parameter file"), "fatal")
	}

	e = ioutil.WriteFile(unix64, bin, sys.FilePermission())
	if e != nil {
		msg.DeployAsset(errors.New("Comet parameter file"), "fatal")
	}

	return
}
