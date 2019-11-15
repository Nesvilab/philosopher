package peptideprophet

import (
	"errors"
	"io/ioutil"

	"philosopher/lib/msg"
	"philosopher/lib/sys"
)

// UnixInteractParser deploys InteractParser
func UnixInteractParser(s string) {

	bin, e := Asset("InteractParser")
	e = ioutil.WriteFile(s, bin, sys.FilePermission())

	if e != nil {
		msg.DeployAsset(errors.New("InteractParser"), "trace")
	}

	return
}

// UnixRefreshParser ...
func UnixRefreshParser(s string) {

	bin, e := Asset("RefreshParser")
	e = ioutil.WriteFile(s, bin, sys.FilePermission())

	if e != nil {
		msg.DeployAsset(errors.New("RefreshParser"), "trace")
	}

	return
}

// UnixPeptideProphetParser ...
func UnixPeptideProphetParser(s string) {

	bin, e := Asset("PeptideProphetParser")
	e = ioutil.WriteFile(s, bin, sys.FilePermission())

	if e != nil {
		msg.DeployAsset(errors.New("PeptideProphetParser"), "trace")
	}

	return
}
