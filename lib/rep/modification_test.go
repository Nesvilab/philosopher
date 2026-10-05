package rep

import (
	"testing"

	"github.com/Nesvilab/philosopher/lib/mod"
	"github.com/Nesvilab/philosopher/lib/wrk"
)

// A mass loss must be annotated like a mass gain: the ppm tolerance is taken on the absolute mass.
// MapMods loads the UniMod ontology through the workspace, so the test runs inside a fresh one.
func TestMapMods_NegativeDeltaMass(t *testing.T) {

	t.Chdir(t.TempDir())
	wrk.Init("test", "test", t.TempDir())

	evi := Evidence{PSM: PSMEvidenceList{{
		Peptide: "PEPTIDEK",
		Modifications: mod.ModificationsSlice{IndexSlice: []mod.Modification{
			{Index: "-0.9840", Name: "Unknown", Type: mod.Observed, MassDiff: -0.984016},
			{Index: "-0.5000", Name: "Unknown", Type: mod.Observed, MassDiff: -0.5},
		}},
	}}}

	evi.MapMods()

	names := make(map[string]string)
	for _, m := range evi.PSM[0].Modifications.IndexSlice {
		names[m.Index] = m.Name
	}

	if names["-0.9840"] != "Amidated" {
		t.Errorf("-0.984016 annotated as %q, want Amidated", names["-0.9840"])
	}
	if names["-0.5000"] != "Unknown" {
		t.Errorf("-0.5 annotated as %q, want Unknown", names["-0.5000"])
	}
}
