package cdhit

import (
	"errors"
	"io/ioutil"

	"github.com/prvst/philosopher/lib/sys"
)

// Win64 ...
func Win64(win64 string) error {

	bin, err := Asset("cd-hit.exe")
	err = ioutil.WriteFile(win64, bin, sys.FilePermission())

	if err != nil {
		return errors.New("Cannot deploy CD-hit")
	}

	return nil
}
