package comet

import (
	"errors"
	"io/ioutil"

	"github.com/prvst/philosopher/lib/err"

	"github.com/prvst/philosopher/lib/sys"
)

// UnixParameterFile ...
func UnixParameterFile(unixParam string) {

	param, e := Asset("comet.params")
	if e != nil {
		err.DeployAsset(errors.New("Comet parameter file"))
	}

	e = ioutil.WriteFile(unixParam, param, sys.FilePermission())
	if e != nil {
		err.DeployAsset(errors.New("Comet parameter file"))
	}

	return
}

// Unix64 ...
func Unix64(unix64 string) {

	bin, e := Asset("comet.2018014.linux.exe")
	if e != nil {
		err.DeployAsset(errors.New("Comet parameter file"))
	}

	e = ioutil.WriteFile(unix64, bin, sys.FilePermission())
	if e != nil {
		err.DeployAsset(errors.New("Comet parameter file"))
	}

	return
}
