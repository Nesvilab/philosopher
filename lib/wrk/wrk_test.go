package wrk_test

import (
	"os"
	"philosopher/lib/sys"
	"philosopher/lib/wrk"
	"testing"
)

func TestWorkspace(t *testing.T) {

	e := os.Chdir("../../test/wrksp")
	if e != nil {
		t.Errorf("Path is incorrect, got %s", e)
	}

	wrk.Init("0000", "0000")

	if _, e := os.Stat(sys.MetaDir()); os.IsNotExist(e) {
		if e != nil {
			t.Errorf("Meta Dir path is incorrect, got %s", e)
		}
	}

	if _, e := os.Stat(sys.Meta()); os.IsNotExist(e) {
		if e != nil {
			t.Errorf("Meta file path is incorrect, got %s", e)
		}
	}

	wrk.Clean()
}
