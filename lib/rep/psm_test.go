package rep

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func indexOf(fields []string, name string) int {
	for i, f := range fields {
		if f == name {
			return i
		}
	}
	return -1
}

func TestPSMReport_MSFraggerScoreColumns(t *testing.T) {

	workspace := t.TempDir()

	evi := PSMEvidenceList{
		{
			Spectrum:                         "a.00001.00001.2",
			SpectrumFile:                     "a.pepXML",
			Peptide:                          "PEPTIDEK",
			ExtendedPeptide:                  "K.PEPTIDEK.A",
			PrevAA:                           "K",
			NextAA:                           "A",
			AssumedCharge:                    2,
			RetentionTime:                    60,
			UncalibratedPrecursorNeutralMass: 1000.5,
			PrecursorNeutralMass:             1000.5,
			CalcNeutralPepMass:               1000.5,
			Expectation:                      0.0015,
			Hyperscore:                       20.5,
			Nextscore:                        10.25,
			BCS:                              5,
			FINterm:                          0.1234,
			Probability:                      0.99,
			Qvalue:                           0.001,
			Protein:                          "sp|P1|TEST",
			ProteinID:                        "P1",
			EntryName:                        "TEST",
			GeneName:                         "T1",
		},
	}

	evi.PSMReport(workspace, "", "rev_", 0, false, false, false, false, false, false, false)

	content, err := os.ReadFile(filepath.Join(workspace, "psm.tsv"))
	if err != nil {
		t.Fatalf("read psm.tsv: %v", err)
	}

	lines := strings.Split(strings.TrimRight(string(content), "\r\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("psm.tsv has %d lines, want header plus one PSM", len(lines))
	}

	header := strings.Split(lines[0], "\t")
	row := strings.Split(lines[1], "\t")

	if len(header) != len(row) {
		t.Fatalf("header has %d columns but the PSM row has %d fields", len(header), len(row))
	}

	idx := indexOf(header, "fI_nterm")
	if idx < 0 {
		t.Fatalf("psm.tsv header lacks fI_nterm: %q", lines[0])
	}

	if header[idx-1] != "BCS" || header[idx+1] != "Probability" {
		t.Errorf("fI_nterm is between %q and %q, want between BCS and Probability", header[idx-1], header[idx+1])
	}

	if row[idx-1] != "5" {
		t.Errorf("BCS column = %q, want 5", row[idx-1])
	}

	if row[idx] != "0.1234" {
		t.Errorf("fI_nterm column = %q, want 0.1234", row[idx])
	}

	if row[idx+1] != "0.9900" {
		t.Errorf("Probability column = %q, want 0.9900", row[idx+1])
	}
}
