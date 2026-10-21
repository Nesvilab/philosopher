package id

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Nesvilab/philosopher/lib/tes"
)

func TestProtXML_PromoteProteinIDs(t *testing.T) {

	p := ProtXML{
		DecoyTag: "rev_",
		Groups: GroupList{
			{
				Proteins: ProtIDList{
					{
						OriginalHeader:           "rev_sp|P00001|X_HUMAN decoy of X",
						ProteinName:              "rev_sp|P00001|X_HUMAN",
						Description:              "decoy of X",
						IndistinguishableProtein: []string{"tr|Q00002|Y_HUMAN", "sp|P00003|Z_HUMAN"},
						Length:                   120,
					},
					{
						OriginalHeader:           "sp|P00004|W_HUMAN Protein W",
						ProteinName:              "sp|P00004|W_HUMAN",
						Description:              "Protein W",
						IndistinguishableProtein: []string{"rev_sp|P00005|V_HUMAN"},
						Length:                   80,
					},
				},
			},
		},
	}

	p.PromoteProteinIDs()

	promoted := p.Groups[0].Proteins[0]

	if promoted.ProteinName != "sp|P00003|Z_HUMAN" {
		t.Errorf("ProteinName = %q, want %q", promoted.ProteinName, "sp|P00003|Z_HUMAN")
	}

	// the decoy takes the promoted protein's place among the indistinguishable proteins
	wantIndi := []string{"tr|Q00002|Y_HUMAN", "rev_sp|P00001|X_HUMAN"}
	if !reflect.DeepEqual(promoted.IndistinguishableProtein, wantIndi) {
		t.Errorf("IndistinguishableProtein = %q, want %q", promoted.IndistinguishableProtein, wantIndi)
	}

	// the decoy's header, description and length are dropped; the header is the one Read builds for a
	// protein without annotation (TestAssembleProteinReport_PromotedProtein resolves it in lib/rep)
	if promoted.OriginalHeader != "sp|P00003|Z_HUMAN " || promoted.Description != "" || promoted.Length != 0 {
		t.Errorf("OriginalHeader, Description, Length = %q, %q, %d, want %q, %q, 0",
			promoted.OriginalHeader, promoted.Description, promoted.Length, "sp|P00003|Z_HUMAN ", "")
	}

	// a target representative is left as it is
	target := p.Groups[0].Proteins[1]
	if target.ProteinName != "sp|P00004|W_HUMAN" || target.OriginalHeader != "sp|P00004|W_HUMAN Protein W" ||
		target.Description != "Protein W" || target.Length != 80 ||
		!reflect.DeepEqual(target.IndistinguishableProtein, []string{"rev_sp|P00005|V_HUMAN"}) {
		t.Errorf("target representative changed: %+v", target)
	}
}

func TestProtXML_MarkUniquePeptides(t *testing.T) {

	tes.SetupTestEnv()

	var p ProtXML

	p.Read("interact.prot.xml")
	p.DecoyTag = "rev_"

	var unique int
	var flag bool

	type fields struct {
		FileName   string
		DecoyTag   string
		Groups     GroupList
		RunOptions string
	}
	type args struct {
		w float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Testing protein Unique marking.",
			args: args{w: 1.00},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			p.MarkUniquePeptides(tt.args.w)

			p.PromoteProteinIDs()

			for _, i := range p.Groups {
				for _, j := range i.Proteins {

					if strings.HasPrefix(string(j.ProteinName), p.DecoyTag) {
						for _, k := range j.IndistinguishableProtein {
							if !strings.HasPrefix(string(k), p.DecoyTag) {
								flag = true
							}
						}
					}

					for _, k := range j.PeptideIons {

						if k.IsUnique == true {
							unique++
						}

					}
				}
			}

			if unique != 11153 {
				t.Errorf("Number of Unque ions in ProtXML is wrong, got %v, want %v", unique, 11153)
			}

			if flag == true {
				t.Errorf("Protein Promotion is no working properly")
			}

		})
	}

}
