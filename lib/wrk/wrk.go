package wrk

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/pierrre/archivefile/zip"
	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/gth"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/sirupsen/logrus"
)

// Run is the workspace main entry point
func Run(Version, Build string, b, c, i, n bool) *err.Error {

	if n == false {
		gth.UpdateChecker(Version, Build)
	}

	if (i == true && b == true && c == true) || (i == true && b == true) || (i == true && c == true) || (c == true && b == true) {
		logrus.Fatal("this command accepts only one parameter")
	}

	if i == true {

		logrus.Info("Creating workspace")
		e := Init(Version, Build)
		if e != nil {
			if e.Class == "warning" {
				logrus.Warn(e.Error())
			}
		}
		return e

	} else if b == true {

		logrus.Info("Creating backup")
		e := Backup()
		if e != nil {
			logrus.Warn(e.Error())
		}
		return e

	} else if c == true {

		logrus.Info("Removing workspace")
		e := Clean()
		if e != nil {
			logrus.Warn(e.Error())
		}
		return e

	}

	return nil
}

// Init creates a new workspace
func Init(version, build string) *err.Error {

	var m met.Data
	m.Restore(sys.Meta())

	if len(m.UUID) > 1 && len(m.Home) > 1 {
		return &err.Error{Type: err.CannotOverwriteMeta, Class: err.WARN}
	}

	dir, e := os.Getwd()
	if e != nil {
		return &err.Error{Type: err.CannotStatLocalDirectory, Class: err.FATA, Argument: "check folder permissions"}
	}

	da := met.New(dir)
	// if e != nil {
	// 	logrus.Fatal(e.Error())
	// }

	da.Version = version
	da.Build = build

	os.Mkdir(da.MetaDir, sys.FilePermission())
	if _, e := os.Stat(sys.MetaDir()); os.IsNotExist(e) {
		return &err.Error{Type: err.CannotCreateDirectory, Class: err.FATA, Argument: "Can't create meta directory; check folder permissions"}
	}

	if runtime.GOOS == sys.Windows() {
		e = HideFile(sys.MetaDir())
		if _, e := os.Stat(sys.MetaDir()); os.IsNotExist(e) {
			return &err.Error{Type: err.CannotCreateDirectory, Class: err.FATA, Argument: "Can't create meta directory; check folder permissions"}
		}
	}

	os.Mkdir(da.Temp, sys.FilePermission())
	if _, e := os.Stat(da.Temp); os.IsNotExist(e) {
		return &err.Error{Type: err.CannotCreateDirectory, Class: err.FATA, Argument: "Can't find temporary directory; check folder permissions"}
	}

	da.Home = fmt.Sprintf("%s", da.Home)
	da.MetaDir = fmt.Sprintf("%s", da.MetaDir)

	da.Serialize()

	return nil
}

// Backup collects all binary files from the workspace and zips them
func Backup() *err.Error {

	var m met.Data
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
