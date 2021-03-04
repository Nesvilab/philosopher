package rawfilereader

import (
	"errors"
	"io/ioutil"

	"philosopher/lib/msg"

	"philosopher/lib/sys"
)

// Deb64 deploys RawfileReader for Debian
func Deb64(unix64 string) {

	bin, e1 := Asset("rawFileReaderDeb")
	if e1 != nil {
		msg.DeployAsset(errors.New("rawFileReaderDeb"), "Cannot read rawFileReaderDeb obo")
	}

	e2 := ioutil.WriteFile(unix64, bin, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("rawFileReaderDeb"), "Cannot deploy rawFileReaderDeb 64-bit")
	}

	return
}
