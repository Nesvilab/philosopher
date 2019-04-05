package obo

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// GetUniModTerms deploys, reads and assemble the unimod data into structs
func GetUniModTerms(temp string) (Terms, error) {

	// deploys unimod database
	f, err := DeployUniModObo(temp)
	if err != nil {
		return err
	}

	// process xml file and load structs
	err = read(f)
	if err != nil {
		return err
	}

	serialize()

	return nil
}

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

func read(s string) error {

	return nil
}

func serialize() {

	return
}
