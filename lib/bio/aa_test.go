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
		{
			name: "Testing Glutamine",
			args: args{"Glutamine"},
			want: AminoAcid{Code: "Q", ShortName: "Gln", Name: "Glutamine", MonoIsotopeMass: 128.058577540, AverageMass: 128.12922},
		},
		{
			name: "Testing Glutamic Acid",
			args: args{"Glutamic Acid"},
			want: AminoAcid{Code: "E", ShortName: "Glu", Name: "Glutamic Acid", MonoIsotopeMass: 129.042593135, AverageMass: 129.11398},
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

func TestMonoIsotopeMassByCode(t *testing.T) {

	for _, code := range "ACDEFGHIKLMNPQRSTVWY" {
		if _, ok := MonoIsotopeMassByCode(string(code)); !ok {
			t.Errorf("MonoIsotopeMassByCode(%q) is unknown, want a standard residue", string(code))
		}
	}

	known := map[string]float64{"Q": 128.058577540, "E": 129.042593135, "C": 103.009184505}
	for code, want := range known {
		if got, _ := MonoIsotopeMassByCode(code); got != want {
			t.Errorf("MonoIsotopeMassByCode(%q) = %v, want %v", code, got, want)
		}
	}

	for _, code := range []string{"U", "X", "B", "n", ""} {
		if _, ok := MonoIsotopeMassByCode(code); ok {
			t.Errorf("MonoIsotopeMassByCode(%q) is known, want unknown", code)
		}
	}
}
