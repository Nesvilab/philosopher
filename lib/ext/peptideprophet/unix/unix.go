package peptideprophet

import (
	"errors"
	"io/ioutil"

	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/sys"
)

// UnixInteractParser deploys InteractParser
func UnixInteractParser(s string) {

	bin, e := Asset("InteractParser")
	e = ioutil.WriteFile(s, bin, sys.FilePermission())

	if e != nil {
		err.DeployAsset(errors.New("InteractParser"))
	}

	return
}

// UnixRefreshParser ...
func UnixRefreshParser(s string) {

	bin, e := Asset("RefreshParser")
	e = ioutil.WriteFile(s, bin, sys.FilePermission())

	if e != nil {
		err.DeployAsset(errors.New("RefreshParser"))
	}

	return
}

// UnixPeptideProphetParser ...
func UnixPeptideProphetParser(s string) {

	bin, e := Asset("PeptideProphetParser")
	e = ioutil.WriteFile(s, bin, sys.FilePermission())

	if e != nil {
		err.DeployAsset(errors.New("PeptideProphetParser"))
	}

	return
}
