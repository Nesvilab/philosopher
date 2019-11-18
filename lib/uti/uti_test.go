package uti_test

import (
	"philosopher/lib/uti"
	"testing"
)

func TestUti(t *testing.T) {

	x := uti.Round(5.3557876867, 5, 2)
	if x != 5.35 {
		t.Errorf("Aminoacid name is incorrect, got %f, want %f", x, 5.35)
	}

	y := uti.ToFixed(5.3557876867, 3)
	if y != 5.355 {
		t.Errorf("Aminoacid name is incorrect, got %f, want %f", y, 5.3557876867)
	}

}
