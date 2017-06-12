package cdhit

import (
	"errors"
	"io/ioutil"
)

// Win64 ...
func Win64(win64 string) error {

	bin, err := Asset("cd-hit.exe")
	err = ioutil.WriteFile(win64, bin, 0755)

	if err != nil {
		return errors.New("Cannot deploy CD-hit")
	}

	return nil
}
