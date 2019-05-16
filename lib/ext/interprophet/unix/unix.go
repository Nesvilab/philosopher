package interprophet

import (
	"errors"
	"io/ioutil"

	"github.com/prvst/philosopher/lib/sys"
)

// UnixInterProphetParser accessor
func UnixInterProphetParser(s string) error {

	bin, err := Asset("InterProphetParser")
	err = ioutil.WriteFile(s, bin, sys.FilePermission())

	if err != nil {
		return errors.New("Cannot deploy InterProphetParser")
	}

	return nil
}
