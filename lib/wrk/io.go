package wrk

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pierrre/archivefile/zip"
	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/meta"
	"github.com/prvst/philosopher/lib/sys"
)

// Backup collects all binary files from the workspace and zips them
func Backup() *err.Error {

	var m meta.Data
	m.Restore(sys.Meta())

	if len(m.UUID) < 1 && len(m.Home) < 1 {
		return &err.Error{Type: err.CannotFindMetaDirectory, Class: err.FATA}
	}

	var name string
	if m.OS == sys.Windows() {
		name = "backup.zip"
	} else {
		t := time.Now()
		timestamp := t.Format(time.RFC3339)
		name = fmt.Sprintf("%s-%s.zip", m.ProjectName, timestamp)
	}
	outFilePath := filepath.Join(m.Home, name)

	progress := func(archivePath string) {
	}

	e := zip.ArchiveFile(sys.MetaDir(), outFilePath, progress)
	if e != nil {
		return &err.Error{Type: err.CannotZipMetaDirectory, Class: err.FATA}
	}

	return nil
}

// Clean deletes all meta data and the workspace directory
func Clean() *err.Error {

	e := os.RemoveAll(sys.MetaDir())
	if e != nil {
		return &err.Error{Type: err.CannotDeleteMetaDirectory, Class: err.FATA, Argument: e.Error()}
	}

	return nil
}
