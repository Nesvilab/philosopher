package wrk

import (
	"fmt"
	"os"

	"github.com/prvst/cmsl/err"
	"github.com/prvst/philosopher/lib/meta"
	"github.com/prvst/philosopher/lib/sys"
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
