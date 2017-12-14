package ptmprophet

import (
	"io/ioutil"

	"github.com/prvst/philosopher/lib/err"
)

// UnixPTMProphetParser locates and extracts the PTMProphet binary
func UnixPTMProphetParser(s string) *err.Error {

	bin, e := Asset("PTMProphetParser")
	e = ioutil.WriteFile(s, bin, 0755)

	if e != nil {
		return &err.Error{Type: err.CannotExtractAsset, Class: err.FATA, Argument: "PTMProphetParser"}
	}

	return nil
}
