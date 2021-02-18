package dat_test

import (
	. "philosopher/lib/dat"
	"philosopher/lib/sys"
	"testing"
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
			d.Fetch(tt.args.id, tt.args.temp, tt.args.iso, tt.args.rev)
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
			d := &Base{
				UniProtDB: tt.fields.UniProtDB,
				CrapDB:    tt.fields.CrapDB,
				TaDeDB:    tt.fields.TaDeDB,
				Records:   tt.fields.Records,
			}
			d.ProcessDB(tt.args.file, tt.args.decoyTag)

			if len(d.Records) != 20381 {
				t.Errorf("Number of FASTA entries is incorrect, got %d, want %d", len(d.Records), 20381)
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
