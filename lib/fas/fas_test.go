package fas_test

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	. "github.com/Nesvilab/philosopher/lib/fas"
	"github.com/Nesvilab/philosopher/lib/tes"
)

func TestParseFile(t *testing.T) {

	tes.SetupTestEnv()

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)

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
			args: args{filename: "../db/uniprot/2019-02-05-td-hsa-reviewed-2019-02-04.fasta"},
			want: 40896,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			path, err := os.Getwd()
			if err != nil {
				log.Println(err)
			}
			fmt.Println(path)

			if got := ParseFile(tt.args.filename); !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("ParseFile() = %d, want %d", len(got), tt.want)
			}
		})
	}

	//tes.ShutDowTestEnv()
}
