package rep

import (
	"reflect"
	"testing"
)

func TestRemoveNonRazorProteinsNotPresentInPsms(t *testing.T) {

	evi := Evidence{
		PSM: PSMEvidenceList{
			{Protein: "sp|P1|A"},
			{Protein: "sp|P7|G", MappedProteins: map[string]string{"sp|P5|E": ""}},
		},
		Proteins: ProteinEvidenceList{
			{PartHeader: "sp|P1|A"},
			{PartHeader: "sp|P2|B"},
			{PartHeader: "sp|P3|C"},
			{PartHeader: "rev_sp|P4|D"},
			{PartHeader: "sp|P5|E"},
			{PartHeader: "sp|P6|F"},
			{PartHeader: "sp|P7|G"},
			{PartHeader: "sp|P8|H"},
		},
	}

	evi = evi.RemoveNonRazorProteinsNotPresentInPsms("rev_")

	var got []string
	for _, p := range evi.Proteins {
		got = append(got, p.PartHeader)
	}

	want := []string{"sp|P1|A", "rev_sp|P4|D", "sp|P5|E", "sp|P7|G"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("proteins = %q, want %q", got, want)
	}
}
