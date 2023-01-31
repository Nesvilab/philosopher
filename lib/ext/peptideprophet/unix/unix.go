package peptideprophet

import (
	"errors"
	"io/ioutil"

	"github.com/Nesvilab/philosopher/lib/msg"
	"github.com/Nesvilab/philosopher/lib/sys"
)

// UnixInteractParser deploys InteractParser
func UnixInteractParser(s string) {

	bin, e1 := Asset("InteractParser")
	if e1 != nil {
		msg.DeployAsset(errors.New("InteractParser"), "Cannot read InteractParser bin")
	}

	e2 := ioutil.WriteFile(s, bin, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("InteractParser"), "Cannot deploy InteractParser")
	}

	return
}

// UnixRefreshParser deploys RefreshParser
func UnixRefreshParser(s string) {

	bin, e1 := Asset("RefreshParser")
	if e1 != nil {
		msg.DeployAsset(errors.New("RefreshParser"), "Cannot read RefreshParser bin")
	}

	e2 := ioutil.WriteFile(s, bin, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("RefreshParser"), "Cannot deploy RefreshParser")
	}

	return
}

// UnixPeptideProphetParser deployes PeptideProphetParser
func UnixPeptideProphetParser(s string) {

	bin, e1 := Asset("PeptideProphetParser")
	if e1 != nil {
		msg.DeployAsset(errors.New("PeptideProphetParser"), "Cannot read PeptideProphetParser bin")
	}

	e2 := ioutil.WriteFile(s, bin, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("PeptideProphetParser"), "Cannot read PeptideProphetParser bin")
	}

	return
}
