package wrk

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/gth"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/sys"
)

// Run is the workspace main entry point
func Run(Version, Build string, b, c, i bool) *err.Error {

	gth.UpdateChecker(Version, Build)

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

	da.Version = version
	da.Build = build

	os.Mkdir(da.MetaDir, 0755)
	if _, e := os.Stat(sys.MetaDir()); os.IsNotExist(e) {
		return &err.Error{Type: err.CannotCreateDirectory, Class: err.FATA, Argument: "Can't create meta directory; check folder permissions"}
	}

	os.Mkdir(da.Temp, 0755)
	if _, e := os.Stat(da.Temp); os.IsNotExist(e) {
		return &err.Error{Type: err.CannotCreateDirectory, Class: err.FATA, Argument: "Can't find temporary directory; check folder permissions"}
	}

	da.Home = fmt.Sprintf("%s", da.Home)
	da.MetaDir = fmt.Sprintf("%s", da.MetaDir)

	da.Serialize()

	return nil
}
