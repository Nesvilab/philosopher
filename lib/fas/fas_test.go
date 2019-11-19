package fas_test

import (
	"os"
	. "philosopher/lib/fas"
	"reflect"
	"testing"
)

// 		It("Parsing UniProt Description", func() {
// 			f := ParseUniProtDescriptionMap("db/uniprot/2019-02-05-td-hsa-reviewed-2019-02-04.fasta")
// 			Expect(len(f)).To(Equal(20448))
// 		})

// 		It("Parsing UniProt Sequence", func() {
// 			f := ParseUniProtSequencenMap("db/uniprot/2019-02-05-td-hsa-reviewed-2019-02-04.fasta")
// 			Expect(len(f)).To(Equal(20448))
// 		})

// 		It("Parsing FASTA Description", func() {
// 			f := ParseFastaDescription("db/uniprot/2019-02-05-td-hsa-reviewed-2019-02-04.fasta")
// 			Expect(len(f)).To(Equal(20448))
// 		})

// 		It("Parsing FASTA File", func() {
// 			f := ParseFile("db/uniprot/2019-02-05-td-hsa-reviewed-2019-02-04.fasta")

// 			f = CleanDatabase(f, "rev_", "cont_")
// 			Expect(len(f)).To(Equal(20448))
// 		})

// 	})

// })

func TestParseFile(t *testing.T) {

	os.Chdir("../../test/db/")

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
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseFile(tt.args.filename); !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("ParseFile() = %d, want %d", len(got), tt.want)
			}
		})
	}
}
