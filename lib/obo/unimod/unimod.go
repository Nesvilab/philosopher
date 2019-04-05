package obo

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// DeployUniModObo deploys the OBO file to the temp folder
func DeployUniModObo(temp string) (string, error) {

	oboFile := fmt.Sprintf("%s%sunimod.obo", temp, string(filepath.Separator))

	param, err := Asset("unimod.obo")
	err = ioutil.WriteFile(oboFile, param, 0644)

	if err != nil {
		msg := fmt.Sprintf("Could not deploy UniMOD database %s", err)
		return oboFile, errors.New(msg)
	}

	return oboFile, nil
}
