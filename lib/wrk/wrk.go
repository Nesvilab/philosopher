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
func Run(Version, Build string, b, c, i, n bool) {

	if n == false {
		gth.UpdateChecker(Version, Build)
	}

	if (i == true && b == true && c == true) || (i == true && b == true) || (i == true && c == true) || (c == true && b == true) {
		logrus.Fatal("this command accepts only one parameter")
	}

	if i == true {

		logrus.Info("Creating workspace")
		Init(Version, Build)

	} else if b == true {

		logrus.Info("Creating backup")
		Backup()

	} else if c == true {

		logrus.Info("Removing workspace")
		Clean()
	}

	return
}

// Init creates a new workspace
func Init(version, build string) {

	var m met.Data
	m.Restore(sys.Meta())

	if len(m.UUID) > 1 && len(m.Home) > 1 {
		err.OverwrittingMeta()
	}

	dir, e := os.Getwd()
	if e != nil {
		err.GettingLocalDir(e)
	}

	da := met.New(dir)

	da.Version = version
	da.Build = build

	os.Mkdir(da.MetaDir, sys.FilePermission())
	if _, e := os.Stat(sys.MetaDir()); os.IsNotExist(e) {
		err.CreatingMetaDirectory(e)
	}

	if runtime.GOOS == sys.Windows() {
		HideFile(sys.MetaDir())
	}

	os.Mkdir(da.Temp, sys.FilePermission())
	if _, e := os.Stat(da.Temp); os.IsNotExist(e) {
		err.LocatingTemDirecotry(e)
	}

	da.Home = fmt.Sprintf("%s", da.Home)
	da.MetaDir = fmt.Sprintf("%s", da.MetaDir)

	da.Serialize()

	return
}

// Backup collects all binary files from the workspace and zips them
func Backup() {

	var m met.Data
	m.Restore(sys.Meta())

	if len(m.UUID) < 1 && len(m.Home) < 1 {
		err.LocatingMetaDirecotry()
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
		err.ArchivingMetaDirecotry(e)
	}

	return
}

// Clean deletes all meta data and the workspace directory
func Clean() {

	var d met.Data
	d.Restore(sys.Meta())

	e := os.RemoveAll(sys.MetaDir())
	if e != nil {
		err.DeletingMetaDirecotry(e)
	}

	if len(d.Temp) > 0 {
		e := os.RemoveAll(d.Temp)
		if e != nil {
			err.DeletingMetaDirecotry(e)
		}
	}

	return
}
