package wrk_test

import (
	"os"
	. "philosopher/lib/wrk"
	"testing"
)

func TestInit(t *testing.T) {
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
	for _, tt := range tests {
		os.Chdir("../../test/wrksp")
		t.Run(tt.name, func(t *testing.T) {
			Init(tt.args.version, tt.args.build)
			Clean()
		})
	}
}
