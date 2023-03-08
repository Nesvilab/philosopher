package met_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/Nesvilab/philosopher/lib/met"
	"github.com/Nesvilab/philosopher/lib/sys"
)

var dir string
var d met.Data

func TestMetaData(t *testing.T) {

	var e error

	dir, e = os.Getwd()
	d = met.New(dir)

	if e != nil {
		t.Errorf("Path is incorrect, got %s", e)
	}

	if len(d.UUID) == 0 {
		t.Errorf("UUID is incorrect, got %s", d.UUID)
	}

	if d.OS != runtime.GOOS {
		t.Errorf("OS name is incorrect, got %s, want %s", d.OS, runtime.GOOS)
	}

	if d.Arch != runtime.GOARCH {
		t.Errorf("Architecture name is incorrect, got %s, want %s", d.Arch, runtime.GOARCH)
	}

	if d.Home != dir {
		t.Errorf("Home name is incorrect, got %s, want %s", d.Home, dir)
	}

	if d.DB != d.Home+string(filepath.Separator)+sys.DBBin() {
		t.Errorf("Database name is incorrect, got %s, want %s", d.DB, d.Home+string(filepath.Separator)+sys.DBBin())
	}

	temp := sys.GetTemp()
	temp += string(filepath.Separator) + d.UUID
	if d.Temp != temp {
		t.Errorf("Temp folder is incorrect, got %s, want %s", d.Temp, temp)
	}

}
