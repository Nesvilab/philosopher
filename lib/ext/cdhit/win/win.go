package cdhit

import (
	"errors"
	"os"

	"github.com/Nesvilab/philosopher/lib/msg"

	"github.com/Nesvilab/philosopher/lib/sys"
)

// Win64 CD-HIT Deploy
func Win64(win64 string) {

	bin, e1 := Asset("cd-hit.exe")
	if e1 != nil {
		msg.DeployAsset(errors.New("CD-HIT"), "Cannot read CD-HIT 64-bit bin")
	}

	e2 := os.WriteFile(win64, bin, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("CD-HIT"), "Cannot deploy CD-HIT 64-bit")
	}
}
