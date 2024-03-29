package wrk_test

import (
	"os"
	"testing"

	"github.com/Nesvilab/philosopher/lib/tes"
	. "github.com/Nesvilab/philosopher/lib/wrk"
)

func TestInit(t *testing.T) {

	tes.SetupTestEnv()

	type args struct {
		version string
		build   string
		temp    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Testing workspace initialization",
			args: args{version: "0000", build: "0000", temp: ""},
		},
	}

	os.Chdir("../../test/wrksp")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init(tt.args.version, tt.args.build, tt.args.temp)
		})
	}
}
