//go:build !windows
// +build !windows

package wrk

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"philosopher/lib/msg"
)

// HideFile makes the .meta folder hidden on Windows
func HideFile(filename string) {
	if !strings.HasPrefix(filepath.Base(filename), ".") {
		e := os.Rename(filename, "."+filename)
		if e != nil {
			msg.Custom(errors.New("cannot hide file"), "fatal")
		}
	}
}
