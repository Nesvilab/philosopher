package fil

import (
	"philosopher/lib/id"
	"philosopher/lib/sys"
	"philosopher/lib/tes"
	"reflect"
	"testing"
)

func TestPepXMLFDRFilter(t *testing.T) {

	tes.SetupTestEnv()

	pepID, _ := readPepXMLInput("interact.pep.xml", "rev_", sys.GetTemp(), false, 0)

	type args struct {
		input     map[string]id.PepIDList
		targetFDR float64
		level     string
		decoyTag  string
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 float64
	}{
		{
			name:  "Testing PSM Filtering, 1st pass",
			args:  args{input: GetUniquePSMs(pepID), targetFDR: 0.01, level: "psm", decoyTag: "rev_"},
			want:  63387,
			want1: 0.1914,
		},
		{
			name:  "Testing Peptide Filtering, 1st pass",
			args:  args{input: GetUniquePeptides(pepID), targetFDR: 0.01, level: "peptide", decoyTag: "rev_"},
			want:  28284,
			want1: 0.723,
		},
		{
			name:  "Testing Ion Filtering, 1st pass",
			args:  args{input: getUniquePeptideIons(pepID), targetFDR: 0.01, level: "ion", decoyTag: "rev_"},
			want:  38151,
			want1: 0.5155,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := PepXMLFDRFilter(tt.args.input, tt.args.targetFDR, tt.args.level, tt.args.decoyTag)
			if !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("PepXMLFDRFilter() got = %v, want %v", len(got), tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("PepXMLFDRFilter() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	tes.ShutDowTestEnv()
}

func TestPickedFDR(t *testing.T) {

	tes.SetupTestEnv()

	proXML := readProtXMLInput("interact.prot.xml", "rev_", 1.00)

	type args struct {
		p id.ProtXML
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Testing PickedFDR Filter",
			args: args{p: proXML},
			want: 7926,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PickedFDR(tt.args.p); !reflect.DeepEqual(len(got.Groups), tt.want) {
				t.Errorf("PickedFDR() = %v, want %v", len(got.Groups), tt.want)
			}
		})
	}

	//tes.ShutDowTestEnv()
}

func TestRazorFilter(t *testing.T) {

	tes.SetupTestEnv()

	proXML := readProtXMLInput("interact.prot.xml", "rev_", 1.00)

	type args struct {
		p id.ProtXML
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Testing Razor Filter",
			args: args{p: proXML},
			want: 7926,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RazorFilter(tt.args.p); !reflect.DeepEqual(len(got.Groups), tt.want) {
				t.Errorf("RazorFilter() = %v, want %v", len(got.Groups), tt.want)
			}
		})
	}

	//tes.ShutDowTestEnv()
}
