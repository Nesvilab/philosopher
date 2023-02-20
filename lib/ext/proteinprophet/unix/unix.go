package proteinprophet

import (
	"errors"
	"os"

	"github.com/Nesvilab/philosopher/lib/msg"
	"github.com/Nesvilab/philosopher/lib/sys"
)

// UnixBatchCoverage deploys batchcoverage
func UnixBatchCoverage(s string) {

	bin, e1 := Asset("batchcoverage")
	if e1 != nil {
		msg.DeployAsset(errors.New("batchcoverage"), "Cannot read batchcoverage bin")
	}

	e2 := os.WriteFile(s, bin, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("batchcoverage"), "Cannot deploy batchcoverage binary")
	}

	return
}

// UnixDatabaseParser deploys DatabaseParser
func UnixDatabaseParser(s string) {

	bin, e1 := Asset("DatabaseParser")
	if e1 != nil {
		msg.DeployAsset(errors.New("DatabaseParser"), "Cannot read batchcoverage bin")
	}

	e2 := os.WriteFile(s, bin, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("DatabaseParser"), "Cannot deploy batchcoverage binary")
	}

	return
}

// UnixProteinProphet deploys Proteinprophet
func UnixProteinProphet(s string) {

	bin, e1 := Asset("ProteinProphet")
	if e1 != nil {
		msg.DeployAsset(errors.New("ProteinProphet"), "Cannot read ProteinProphet bin")
	}

	e2 := os.WriteFile(s, bin, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("ProteinProphet"), "Cannot deploy ProteinProphet binary")
	}

	return
}
