package cdhit

import (
	"errors"
	"io/ioutil"

	"github.com/prvst/philosopher/lib/sys"
)

// Unix64 ...
func Unix64(unix64 string) error {

	bin, err := Asset("cd-hit")
	err = ioutil.WriteFile(unix64, bin, sys.FilePermission())

	if err != nil {
		return errors.New("Cannot deploy CD-hit")
	}

	return nil
}
