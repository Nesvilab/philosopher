package interprophet

import (
	"errors"
	"io/ioutil"
)

// UnixInterProphetParser accessor
func UnixInterProphetParser(s string) error {

	bin, err := Asset("InterProphetParser")
	err = ioutil.WriteFile(s, bin, 0755)

	if err != nil {
		return errors.New("Cannot deploy InterProphetParser")
	}

	return nil
}
