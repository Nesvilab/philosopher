package dat_test

import (
	"philosopher/lib/dat"
	"philosopher/lib/sys"
	"reflect"
	"testing"
)

func TestDat(t *testing.T) {

}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want dat.Base
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dat.New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBase_Fetch(t *testing.T) {
	type fields struct {
		UniProtDB string
		CrapDB    string
		TaDeDB    map[string]string
		Records   []dat.Record
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
			d := &dat.Base{
				UniProtDB: tt.fields.UniProtDB,
				CrapDB:    tt.fields.CrapDB,
				TaDeDB:    tt.fields.TaDeDB,
				Records:   tt.fields.Records,
			}
			d.Fetch(tt.args.id, tt.args.temp, tt.args.iso, tt.args.rev)
		})
	}
}
