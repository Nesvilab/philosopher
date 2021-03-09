package rawfilereader

import (
	"errors"
	"io/ioutil"

	"philosopher/lib/msg"

	"philosopher/lib/sys"
)

// Reh64 deploys RawfileReader for Red Hat
func Reh64(unix64 string) {

	bin, e1 := Asset("rawFileReaderReH")
	if e1 != nil {
		msg.DeployAsset(errors.New("rawFileReaderReH"), "Cannot read rawFileReaderDeb obo")
	}

	e2 := ioutil.WriteFile(unix64, bin, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("rawFileReaderReH"), "Cannot deploy rawFileReaderDeb 64-bit")
	}

	return
}
