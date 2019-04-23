// +build !windows

package wrk

import (
	"os"
	"path/filepath"
	"strings"
)

// HideFile makes the .meta folder hidden on Windows
func HideFile(filename string) error {
	if !strings.HasPrefix(filepath.Base(filename), ".") {
		err := os.Rename(filename, "."+filename)
		if err != nil {
			return err
		}
	}
	return nil
}
