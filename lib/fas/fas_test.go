package fas_test

import (
	"os"
	. "philosopher/lib/fas"
	"philosopher/lib/wrk"
	"reflect"
	"testing"
)

func TestParseFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Testing Fasta file parsing",
			args: args{filename: "uniprot/2019-02-05-td-hsa-reviewed-2019-02-04.fasta"},
			want: 40896,
		},
	}
	for _, tt := range tests {

		os.Chdir("../../test/db/")
		wrk.Init("0000", "0000")

		t.Run(tt.name, func(t *testing.T) {
			if got := ParseFile(tt.args.filename); !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("ParseFile() = %d, want %d", len(got), tt.want)
			}
		})

		wrk.Clean()
	}
}
