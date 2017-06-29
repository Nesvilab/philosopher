package ptmprophet

import (
	"io/ioutil"

	"github.com/prvst/cmsl/err"
)

// WinPTMProphetParser locates and extracts the PTMProphet binary
func WinPTMProphetParser(s string) *err.Error {

	bin, e := Asset("PTMProphetParser.exe")
	e = ioutil.WriteFile(s, bin, 0755)

	if e != nil {
		return &err.Error{Type: err.CannotExtractAsset, Class: err.FATA, Argument: "PTMProphetParser"}
	}

	return nil
}
