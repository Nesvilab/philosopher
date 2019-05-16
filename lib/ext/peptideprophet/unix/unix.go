package peptideprophet

import (
	"errors"
	"io/ioutil"
)

// UnixInteractParser ...
func UnixInteractParser(s string) error {

	bin, err := Asset("InteractParser")
	err = ioutil.WriteFile(s, bin, sys.FilePermission())

	if err != nil {
		return errors.New("Cannot deploy InteractParser")
	}

	return nil
}

// UnixRefreshParser ...
func UnixRefreshParser(s string) error {

	bin, err := Asset("RefreshParser")
	err = ioutil.WriteFile(s, bin, sys.FilePermission())

	if err != nil {
		return errors.New("Cannot deploy RefreshParser")
	}

	return nil
}

// UnixPeptideProphetParser ...
func UnixPeptideProphetParser(s string) error {

	bin, err := Asset("PeptideProphetParser")
	err = ioutil.WriteFile(s, bin, sys.FilePermission())

	if err != nil {
		return errors.New("Cannot deploy PeptideProphetParser")
	}

	return nil
}
