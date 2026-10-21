package fil

import (
	"os"
	"testing"

	"github.com/Nesvilab/philosopher/lib/id"
	"github.com/Nesvilab/philosopher/lib/sys"
)

// A decoy representative promoted to a target among its indistinguishable proteins keeps the razor peptides
// of its group. RazorFilter keeps its assignment in the workspace, so the test runs inside a fresh one.
func TestRazorFilter_PromotedProtein(t *testing.T) {

	t.Chdir(t.TempDir())
	os.Mkdir(sys.MetaDir(), sys.FilePermission())

	p := id.ProtXML{
		DecoyTag: "rev_",
		Groups: id.GroupList{
			{
				Proteins: id.ProtIDList{
					{
						ProteinName:              "rev_sp|P00001|X_HUMAN",
						IndistinguishableProtein: []string{"tr|Q00002|Y_HUMAN", "sp|P00003|Z_HUMAN"},
						TotalNumberPeptides:      2,
						GroupSiblingID:           "a",
						PeptideIons: []id.PeptideIonIdentification{
							{PeptideSequence: "ONLYHEREPEPK", Weight: 1, GroupWeight: 1},
							{PeptideSequence: "SHAREDPEPK", Weight: 0.5, GroupWeight: 0.5, PeptideParentProtein: []string{"sp|Q00006|Q_HUMAN"}},
						},
					},
				},
			},
			{
				Proteins: id.ProtIDList{
					{
						ProteinName:         "sp|Q00006|Q_HUMAN",
						TotalNumberPeptides: 1,
						GroupSiblingID:      "a",
						PeptideIons: []id.PeptideIonIdentification{
							{PeptideSequence: "SHAREDPEPK", Weight: 0.5, GroupWeight: 0.5, PeptideParentProtein: []string{"sp|P00003|Z_HUMAN"}},
						},
					},
				},
			},
		},
	}

	p.PromoteProteinIDs()
	p = RazorFilter(p)

	promoted := p.Groups[0].Proteins[0]
	for _, i := range promoted.PeptideIons {
		if i.Razor != 1 {
			t.Errorf("%s is not razor to %s", i.PeptideSequence, promoted.ProteinName)
		}
	}
	if !promoted.HasRazor {
		t.Errorf("%s has no razor peptide", promoted.ProteinName)
	}
}

// func TestPepXMLFDRFilter(t *testing.T) {

// 	tes.SetupTestEnv()

// 	pepID, _ := readPepXMLInput("interact.pep.xml", "rev_", sys.GetTemp(), false, 0)

// 	type args struct {
// 		input     map[string]id.PepIDList
// 		targetFDR float64
// 		level     string
// 		decoyTag  string
// 	}
// 	tests := []struct {
// 		name  string
// 		args  args
// 		want  int
// 		want1 float64
// 	}{
// 		{
// 			name:  "Testing PSM Filtering, 1st pass",
// 			args:  args{input: GetUniquePSMs(pepID), targetFDR: 0.01, level: "psm", decoyTag: "rev_"},
// 			want:  63387,
// 			want1: 0.1914,
// 		},
// 		{
// 			name:  "Testing Peptide Filtering, 1st pass",
// 			args:  args{input: GetUniquePeptides(pepID), targetFDR: 0.01, level: "peptide", decoyTag: "rev_"},
// 			want:  28284,
// 			want1: 0.723,
// 		},
// 		{
// 			name:  "Testing Ion Filtering, 1st pass",
// 			args:  args{input: getUniquePeptideIons(pepID), targetFDR: 0.01, level: "ion", decoyTag: "rev_"},
// 			want:  38151,
// 			want1: 0.5155,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, got1 := PepXMLFDRFilter(tt.args.input, tt.args.targetFDR, tt.args.level, tt.args.decoyTag)
// 			if !reflect.DeepEqual(len(got), tt.want) {
// 				t.Errorf("PepXMLFDRFilter() got = %v, want %v", len(got), tt.want)
// 			}
// 			if got1 != tt.want1 {
// 				t.Errorf("PepXMLFDRFilter() got1 = %v, want %v", got1, tt.want1)
// 			}
// 		})
// 	}

// 	//tes.ShutDowTestEnv()
// }

// func TestPickedFDR(t *testing.T) {

// 	tes.SetupTestEnv()

// 	proXML := readProtXMLInput("interact.prot.xml", "rev_", 1.00)

// 	type args struct {
// 		p id.ProtXML
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want int
// 	}{
// 		{
// 			name: "Testing PickedFDR Filter",
// 			args: args{p: proXML},
// 			want: 7926,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := PickedFDR(tt.args.p); !reflect.DeepEqual(len(got.Groups), tt.want) {
// 				t.Errorf("PickedFDR() = %v, want %v", len(got.Groups), tt.want)
// 			}
// 		})
// 	}

// 	//tes.ShutDowTestEnv()
// }

// func TestRazorFilter(t *testing.T) {

// 	tes.SetupTestEnv()

// 	proXML := readProtXMLInput("interact.prot.xml", "rev_", 1.00)

// 	type args struct {
// 		p id.ProtXML
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want int
// 	}{
// 		{
// 			name: "Testing Razor Filter",
// 			args: args{p: proXML},
// 			want: 7926,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := RazorFilter(tt.args.p); !reflect.DeepEqual(len(got.Groups), tt.want) {
// 				t.Errorf("RazorFilter() = %v, want %v", len(got.Groups), tt.want)
// 			}
// 		})
// 	}

// 	//tes.ShutDowTestEnv()
// }
