package msconvert

import (
	"errors"
	"io/ioutil"

	"github.com/prvst/philosopher/lib/sys"
)

// Darwinx64 deploy
func Darwinx64(unix64 string) error {

	bin, err := Asset("msconvert")
	if err != nil {
		return errors.New("Cannot deploy msconvert")
	}

	err = ioutil.WriteFile(unix64, bin, sys.FilePermission())
	if err != nil {
		return errors.New("Cannot deploy msconvert")
	}

	return nil
}
