package dat_test

import (
	"github.com/Nesvilab/philosopher/lib/fas"
	"math"
	"sort"
	"testing"

	. "github.com/Nesvilab/philosopher/lib/dat"
	"github.com/Nesvilab/philosopher/lib/sys"
)

func TestBase_Fetch(t *testing.T) {
	type fields struct {
		UniProtDB string
		CrapDB    string
		TaDeDB    map[string]string
		Records   []Record
	}
	type args struct {
		id   string
		temp string
		iso  bool
		rev  bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Download human proteome",
			args: args{id: "UP000005640", temp: sys.GetTemp(), iso: false, rev: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Base{
				UniProtDB: tt.fields.UniProtDB,
				CrapDB:    tt.fields.CrapDB,
				TaDeDB:    tt.fields.TaDeDB,
				Records:   tt.fields.Records,
			}
			d.Fetch(tt.args.id, "9606", tt.args.temp, tt.args.iso, tt.args.rev)
		})
	}
}

func TestBase_ProcessDB(t *testing.T) {
	type fields struct {
		UniProtDB string
		CrapDB    string
		TaDeDB    map[string]string
		Records   []Record
	}
	type args struct {
		file     string
		decoyTag string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Testing Sequence Parsing - UniProt",
			args: args{file: "/tmp/UP000005640.fas", decoyTag: "rev_"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			len_records := ParseFile(tt.args.file, make(chan<- []fas.FastaEntry, 1024))
			if math.Abs(float64(len_records-20413)) > 100 {
				t.Errorf("Number of FASTA entries is incorrect, got %d, want around %d", len_records, 20413)
			}
		})
	}
}

func TestReverseKeepKR(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		// peptide-level database: one tryptic peptide per entry
		{name: "Tryptic peptide ending in K", args: "AAAAALQAK", want: "AQLAAAAAK"},
		{name: "Tryptic peptide ending in R", args: "LSAAEAQAVR", want: "VAQAEAASLR"},
		{name: "Leading M is kept in place", args: "MPEPTIDEK", want: "MEDITPEPK"},
		{name: "Odd length with leading M", args: "MASDEFGH", want: "MHGFEDSA"},
		{name: "Internal K and R stay put", args: "SAKLGVR", want: "ASKVGLR"},
		// protein-level database: digestion still yields K/R-terminated decoys
		{name: "Protein with a non-tryptic tail", args: "MPEPTIDEKSAMPLERAAA", want: "MEDITPEPKELPMASRAAA"},
		// edge cases
		{name: "Empty sequence", args: "", want: ""},
		{name: "No cleavage residue at all", args: "AAGGSS", want: "SSGGAA"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReverseKeepKR(tt.args); got != tt.want {
				t.Errorf("ReverseKeepKR(%q) = %q, want %q", tt.args, got, tt.want)
			}
		})
	}
}

func TestReverseKeepKR_PreservesTrypticCTerm(t *testing.T) {
	// The defect this mode fixes: on a peptide-level database the existing modes
	// move the C-terminal K/R to the front, so the decoys are no longer tryptic
	// and are trivially separable from the targets.
	peptides := []string{
		"AAAAALQAK", "LSAAEAQAVR", "MPEPTIDEK", "HSEAAASK",
		"EKETALTELR", "TGGRTEDELTR", "SPGDSPGGK", "EALEREAPGGK",
	}
	for _, p := range peptides {
		d := ReverseKeepKR(p)
		if len(d) != len(p) {
			t.Errorf("ReverseKeepKR(%q) changed the length: got %d, want %d", p, len(d), len(p))
		}
		if d[len(d)-1] != p[len(p)-1] {
			t.Errorf("ReverseKeepKR(%q) = %q, C-terminal residue changed from %c to %c",
				p, d, p[len(p)-1], d[len(d)-1])
		}
		if sortedResidues(d) != sortedResidues(p) {
			t.Errorf("ReverseKeepKR(%q) = %q, amino acid composition changed", p, d)
		}
	}
}

func sortedResidues(s string) string {
	r := []byte(s)
	sort.Slice(r, func(i, j int) bool { return r[i] < r[j] })
	return string(r)
}

func TestBase_Deploy(t *testing.T) {
	type fields struct {
		UniProtDB string
		CrapDB    string
		TaDeDB    map[string]string
		Records   []Record
	}
	type args struct {
		temp string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Testing Crapome deployment",
			args: args{temp: sys.GetTemp()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Base{
				UniProtDB: tt.fields.UniProtDB,
				CrapDB:    tt.fields.CrapDB,
				TaDeDB:    tt.fields.TaDeDB,
				Records:   tt.fields.Records,
			}
			d.Deploy(tt.args.temp)
		})
	}
}
