// +build windows

package wrk

import (
	"syscall"

	"github.com/prvst/philosopher/lib/err"
)

// HideFile makes the .meta folder hidden on Windows
func HideFile(filename string) {
	filenameW, e := syscall.UTF16PtrFromString(filename)
	if e != nil {
		err.FatalCustom(e)
	}
	e = syscall.SetFileAttributes(filenameW, syscall.FILE_ATTRIBUTE_HIDDEN)
	if e != nil {
		err.ErrorCustom(e)
	}
	return
}
