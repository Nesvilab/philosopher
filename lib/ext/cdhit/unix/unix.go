package cdhit

import (
	"errors"
	"io/ioutil"
)

// Unix64 ...
func Unix64(unix64 string) error {

	bin, err := Asset("cd-hit")
	err = ioutil.WriteFile(unix64, bin, 0755)

	if err != nil {
		return errors.New("Cannot deploy CD-hit")
	}

	return nil
}
