package proteinprophet

import (
	"errors"
	"io/ioutil"

	"github.com/prvst/philosopher/lib/msg"
	"github.com/prvst/philosopher/lib/sys"
)

// UnixBatchCoverage deploys batchcoverage
func UnixBatchCoverage(s string) {

	bin, e := Asset("batchcoverage")
	e = ioutil.WriteFile(s, bin, sys.FilePermission())

	if e != nil {
		msg.DeployAsset(errors.New("batchcoverage"), "trace")
	}

	return
}

// UnixDatabaseParser deploys DatabaseParser
func UnixDatabaseParser(s string) {

	bin, e := Asset("DatabaseParser")
	e = ioutil.WriteFile(s, bin, sys.FilePermission())

	if e != nil {
		msg.DeployAsset(errors.New("DatabaseParser"), "trace")
	}

	return
}

// UnixProteinProphet deploys Proteinprophet
func UnixProteinProphet(s string) {

	bin, e := Asset("ProteinProphet")
	e = ioutil.WriteFile(s, bin, sys.FilePermission())

	if e != nil {
		msg.DeployAsset(errors.New("ProteinProphet"), "trace")
	}

	return
}
