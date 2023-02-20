package interprophet

import (
	"errors"
	"os"

	"github.com/Nesvilab/philosopher/lib/msg"

	"github.com/Nesvilab/philosopher/lib/sys"
)

// UnixInterProphetParser accessor
func UnixInterProphetParser(s string) {

	bin, e1 := Asset("InterProphetParser")
	if e1 != nil {
		msg.DeployAsset(errors.New("InterProphetParser"), "Cannot read InterProphetParser bin")
	}

	e2 := os.WriteFile(s, bin, sys.FilePermission())
	if e2 != nil {
		msg.DeployAsset(errors.New("InterProphetParser"), "Cannot deploy InterProphetParser")
	}
}
