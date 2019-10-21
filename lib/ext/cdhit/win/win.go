package cdhit

import (
	"errors"
	"io/ioutil"

	"github.com/nesvilab/philosopher/lib/msg"

	"github.com/nesvilab/philosopher/lib/sys"
)

// Win64 ...
func Win64(win64 string) {

	bin, e := Asset("cd-hit.exe")
	e = ioutil.WriteFile(win64, bin, sys.FilePermission())

	if e != nil {
		msg.ExecutingBinary(errors.New("CD-hit"), "trace")
	}

	return
}
