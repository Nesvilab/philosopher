package bio

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want AminoAcid
	}{
		{
			name: "Testing Alanine",
			args: args{"Alanine"},
			want: AminoAcid{Code: "A", ShortName: "Ala", Name: "Alanine", MonoIsotopeMass: 71.037113805, AverageMass: 71.0779},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
