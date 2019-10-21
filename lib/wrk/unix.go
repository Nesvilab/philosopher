// +build !windows

package wrk

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/nesvilab/philosopher/lib/msg"
)

// HideFile makes the .meta folder hidden on Windows
func HideFile(filename string) {
	if !strings.HasPrefix(filepath.Base(filename), ".") {
		e := os.Rename(filename, "."+filename)
		if e != nil {
			msg.Custom(errors.New("Cannot hide file"), "error")
		}
	}
	return
}
