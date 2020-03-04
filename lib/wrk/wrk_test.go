package wrk_test

import (
	"os"
	"philosopher/lib/tes"
	. "philosopher/lib/wrk"
	"testing"
)

func TestInit(t *testing.T) {

	tes.SetupTestEnv()

	type args struct {
		version string
		build   string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Testing workspace initialization",
			args: args{version: "0000", build: "0000"},
		},
	}

	os.Chdir("../../test/wrksp")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init(tt.args.version, tt.args.build)
		})
	}

	tes.ShutDowTestEnv()
}
