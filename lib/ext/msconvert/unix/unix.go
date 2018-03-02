package msconvert

import (
	"errors"
	"io/ioutil"
)

// Unix64 deploy
func Unix64(unix64 string) error {

	bin, err := Asset("msconvert")
	if err != nil {
		return errors.New("Cannot deploy msconvert")
	}

	err = ioutil.WriteFile(unix64, bin, 0755)
	if err != nil {
		return errors.New("Cannot deploy msconvert")
	}

	return nil
}
