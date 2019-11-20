package bio_test

import (
	. "philosopher/lib/bio"
	"philosopher/test"
	"testing"
)

func TestAminoAcids(t *testing.T) {

	test.SetupTestEnv()

	a := New("Alanine")
	if a.Name != "Alanine" {
		t.Errorf("Aminoacid name is incorrect, got %s, want %s", a.Name, "Alanine")
	}

	a = New("Arginine")
	if a.Name != "Arginine" {
		t.Errorf("Aminoacid name is incorrect, got %s, want %s", a.Name, "Arginine")
	}

	a = New("Asparagine")
	if a.Name != "Asparagine" {
		t.Errorf("Aminoacid name is incorrect, got %s, want %s", a.Name, "Asparagine")
	}

	a = New("Aspartic Acid")
	if a.Name != "Aspartic Acid" {
		t.Errorf("Aminoacid name is incorrect, got %s, want %s", a.Name, "Aspartic Acid")
	}

	a = New("Cysteine")
	if a.Name != "Cysteine" {
		t.Errorf("Aminoacid name is incorrect, got %s, want %s", a.Name, "Cysteine")
	}

	a = New("Glutamine")
	if a.Name != "Glutamine" {
		t.Errorf("Aminoacid name is incorrect, got %s, want %s", a.Name, "Glutamine")
	}

	a = New("Glutamic Acid")
	if a.Name != "Glutamic Acid" {
		t.Errorf("Aminoacid name is incorrect, got %s, want %s", a.Name, "Glutamic Acid")
	}

	a = New("Glycine")
	if a.Name != "Glycine" {
		t.Errorf("Aminoacid name is incorrect, got %s, want %s", a.Name, "Glycine")
	}

	a = New("Histidine")
	if a.Name != "Histidine" {
		t.Errorf("Aminoacid name is incorrect, got %s, want %s", a.Name, "Histidine")
	}

	a = New("Isoleucine")
	if a.Name != "Isoleucine" {
		t.Errorf("Aminoacid name is incorrect, got %s, want %s", a.Name, "Isoleucine")
	}

	a = New("Leucine")
	if a.Name != "Leucine" {
		t.Errorf("Aminoacid name is incorrect, got %s, want %s", a.Name, "Leucine")
	}

	a = New("Lysine")
	if a.Name != "Lysine" {
		t.Errorf("Aminoacid name is incorrect, got %s, want %s", a.Name, "Lysine")
	}

	a = New("Methionine")
	if a.Name != "Methionine" {
		t.Errorf("Aminoacid name is incorrect, got %s, want %s", a.Name, "Methionine")
	}

	a = New("Phenylalanine")
	if a.Name != "Phenylalanine" {
		t.Errorf("Aminoacid name is incorrect, got %s, want %s", a.Name, "Phenylalanine")
	}

	a = New("Proline")
	if a.Name != "Proline" {
		t.Errorf("Aminoacid name is incorrect, got %s, want %s", a.Name, "Proline")
	}

	a = New("Serine")
	if a.Name != "Serine" {
		t.Errorf("Aminoacid name is incorrect, got %s, want %s", a.Name, "Serine")
	}

	a = New("Threonine")
	if a.Name != "Threonine" {
		t.Errorf("Aminoacid name is incorrect, got %s, want %s", a.Name, "Threonine")
	}

	a = New("Tryptophan")
	if a.Name != "Tryptophan" {
		t.Errorf("Aminoacid name is incorrect, got %s, want %s", a.Name, "Tryptophan")
	}

	a = New("Tyrosine")
	if a.Name != "Tyrosine" {
		t.Errorf("Aminoacid name is incorrect, got %s, want %s", a.Name, "Tyrosine")
	}

	a = New("Valine")
	if a.Name != "Valine" {
		t.Errorf("Aminoacid name is incorrect, got %s, want %s", a.Name, "Valine")
	}

	test.ShutDowTestEnv()
}

// func TestProtonMass(t *testing.T) {
// 	p := Proton
// 	if p != float64(1.007276) {
// 		t.Errorf("Proton mass is incorrect, got %f, want %f", p, 1.007276)
// 	}
// }

func TestEnzymes(t *testing.T) {

	test.SetupTestEnv()

	var e Enzyme

	e.Synth("Trypsin")
	if e.Name != "trypsin" {
		t.Errorf("Enzyme is incorrect, got %s, want %s", e.Name, "trypsin")
	}

	e.Synth("Lys_c")
	if e.Name != "lys_c" {
		t.Errorf("Enzyme is incorrect, got %s, want %s", e.Name, "lys_c")
	}

	e.Synth("Chymotrypsin")
	if e.Name != "chymotrypsin" {
		t.Errorf("Enzyme is incorrect, got %s, want %s", e.Name, "chymotrypsin")
	}

	e.Synth("Glu_c")
	if e.Name != "glu_c" {
		t.Errorf("Enzyme is incorrect, got %s, want %s", e.Name, "glu_c")
	}

	test.ShutDowTestEnv()
}
