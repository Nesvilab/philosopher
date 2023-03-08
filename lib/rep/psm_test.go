package rep

import (
	"testing"

	"github.com/Nesvilab/philosopher/lib/tes"
)

func TestPSMEvidenceList_PSMReport(t *testing.T) {

	tes.SetupTestEnv()

	//var repoPSM PSMEvidenceList
	//RestorePSM(&repoPSM)

	type args struct {
		workspace    string
		brand        string
		decoyTag     string
		channels     int
		hasDecoys    bool
		isComet      bool
		hasLoc       bool
		hasIonMob    bool
		hasLabels    bool
		hasPrefix    bool
		removeContam bool
	}
	tests := []struct {
		name string
		evi  PSMEvidenceList
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		_ = tt
		// t.Run(tt.name, func(t *testing.T) {
		// 	tt.evi.PSMReport(tt.args.workspace, tt.args.brand, tt.args.decoyTag, tt.args.channels, tt.args.hasDecoys, tt.args.isComet, tt.args.hasLoc, tt.args.hasIonMob, tt.args.hasLabels, tt.args.hasPrefix, tt.args.removeContam)
		// })
	}
}
