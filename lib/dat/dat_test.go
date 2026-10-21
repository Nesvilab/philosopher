package dat_test

import (
	"github.com/Nesvilab/philosopher/lib/fas"
	"math"
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

func TestProcessHeader_ProteinExistence(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   string
	}{
		{
			name:   "PE before SV",
			header: "sp|P04637|P53_HUMAN Cellular tumor antigen p53 OS=Homo sapiens OX=9606 GN=TP53 PE=1 SV=4",
			want:   "1:Experimental evidence at protein level",
		},
		{
			name:   "PE at the end",
			header: "sp|fragnovo1|FRAGNOVO_1 peptide 1 from FragNovo PE=2",
			want:   "2:Experimental evidence at transcript level",
		},
		{
			name:   "PE followed by OX and GN",
			header: "tr|A0A024R161|A0A024R161_HUMAN Guanine nucleotide-binding protein subunit gamma OS=Homo sapiens PE=3 OX=9606 GN=DNAJC25-GNG10",
			want:   "3:Protein inferred from homology",
		},
		{
			name:   "TYPE lookalike",
			header: "sp|fragnovo1|FRAGNOVO_1 peptide 1 from FragNovo TYPE=3 SV=1",
			want:   "",
		},
		{
			name:   "no PE",
			header: "sp|P04637|P53_HUMAN Cellular tumor antigen p53 OS=Homo sapiens OX=9606 GN=TP53 SV=4",
			want:   "",
		},
		{
			name:   "unknown PE level",
			header: "sp|P04637|P53_HUMAN Cellular tumor antigen p53 OS=Homo sapiens OX=9606 GN=TP53 PE=9 SV=4",
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ProcessHeader(tt.header, "", Classify(tt.header, "rev_"), "rev_", false)
			if r.ProteinExistence != tt.want {
				t.Errorf("ProteinExistence = %q, want %q", r.ProteinExistence, tt.want)
			}
		})
	}
}
