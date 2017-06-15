package wrk

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pierrre/archivefile/zip"
	"github.com/prvst/cmsl-source/err"
	"github.com/prvst/philosopher-source/lib/meta"
	"github.com/prvst/philosopher-source/lib/sys"
)

// Init creates a new workspace
func Init() *err.Error {

	var m meta.Data
	m.Restore(sys.Meta())

	if len(m.UUID) > 1 && len(m.Home) > 1 {
		return &err.Error{Type: err.CannotOverwriteMeta, Class: err.WARN}
	}

	dir, e := os.Getwd()
	if e != nil {
		return &err.Error{Type: err.CannotStatLocalDirectory, Class: err.FATA, Argument: "check folder permissions"}
	}

	da := meta.New(dir)

	os.Mkdir(da.MetaDir, 0755)
	os.Mkdir(da.Temp, 0755)

	da.Home = fmt.Sprintf("%s", da.Home)
	da.MetaDir = fmt.Sprintf("%s", da.MetaDir)

	if _, e := os.Stat(sys.MetaDir()); os.IsNotExist(e) {
		return &err.Error{Type: err.CannotCreateDirectory, Class: err.FATA, Argument: "check folder permissions"}
	}

	da.Serialize()

	return nil
}

// Backup collects all binary files from the mea folder and zips them
func Backup() error {

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
		//fmt.Println(archivePath)
	}

	e := zip.ArchiveFile(sys.MetaDir(), outFilePath, progress)
	if e != nil {
		return &err.Error{Type: err.CannotZipMetaDirectory, Class: err.FATA}
	}

	return nil
}

// Clean deletes all meta data and the directory itself
func Clean() error {

	e := os.RemoveAll(sys.MetaDir())
	if e != nil {
		return &err.Error{Type: err.CannotDeleteMetaDirectory, Class: err.FATA, Argument: e.Error()}
	}

	return nil
}
