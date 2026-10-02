package rep

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Nesvilab/philosopher/lib/id"
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

const (
	isfColumn                  = "Is In-Source Fragmented Peptide"
	isfParentPeptideColumn     = "Parent Peptide"
	isfTestParentPeptide       = "n[42.0106]PEPTM[15.9949]IDEKAR"
	isfTestFragmentPeptide     = "PEPTIDEK"
	isfTestNonFragmentPeptide  = "ELVISLIVESK"
	isfTestFragmentSpectrum    = "a.00001.00001.2"
	isfTestNonFragmentSpectrum = "a.00002.00002.2"
	isfTestPeptideColumn       = "Peptide"
	isfTestContaminantColumn   = "Is Contaminant"
	isfTestUniqueColumn        = "Is Unique"
)

// isfTestPSM returns a minimal target PSM for psm.tsv tests.
func isfTestPSM(spectrum, peptide, parentPeptide string) PSMEvidence {
	return PSMEvidence{
		Spectrum:                         spectrum,
		SpectrumFile:                     "a.pepXML",
		Peptide:                          peptide,
		PrevAA:                           "K",
		NextAA:                           "A",
		AssumedCharge:                    2,
		RetentionTime:                    60,
		UncalibratedPrecursorNeutralMass: 1000.5,
		PrecursorNeutralMass:             1000.5,
		CalcNeutralPepMass:               1000.5,
		Probability:                      0.99,
		Protein:                          "sp|P1|TEST",
		ProteinID:                        "P1",
		EntryName:                        "TEST",
		GeneName:                         "T1",
		IsUnique:                         true,
		ISFParentPeptide:                 parentPeptide,
	}
}

// writeAndReadPSMReport writes psm.tsv for evi and returns its header and the rows keyed by peptide.
func writeAndReadPSMReport(t *testing.T, evi PSMEvidenceList) ([]string, map[string][]string) {
	t.Helper()

	workspace := t.TempDir()
	evi.PSMReport(workspace, "", "rev_", 0, false, false, false, false, false, false, false)

	content, err := os.ReadFile(filepath.Join(workspace, "psm.tsv"))
	if err != nil {
		t.Fatalf("read psm.tsv: %v", err)
	}

	lines := strings.Split(strings.TrimRight(string(content), "\r\n"), "\n")
	header := strings.Split(lines[0], "\t")
	peptideIdx := indexOf(header, isfTestPeptideColumn)
	if peptideIdx < 0 {
		t.Fatalf("psm.tsv header lacks %s: %q", isfTestPeptideColumn, lines[0])
	}

	rows := make(map[string][]string, len(lines)-1)
	for n, line := range lines[1:] {
		row := strings.Split(line, "\t")
		if len(row) != len(header) {
			t.Fatalf("row %d has %d fields, header has %d", n+1, len(row), len(header))
		}
		rows[row[peptideIdx]] = row
	}

	return header, rows
}

func TestPSMReport_ISFColumns(t *testing.T) {

	evi := PSMEvidenceList{
		isfTestPSM(isfTestFragmentSpectrum, isfTestFragmentPeptide, isfTestParentPeptide),
		isfTestPSM(isfTestNonFragmentSpectrum, isfTestNonFragmentPeptide, ""),
	}

	header, rows := writeAndReadPSMReport(t, evi)

	idx := indexOf(header, isfColumn)
	if idx < 0 {
		t.Fatalf("psm.tsv header lacks %q: %q", isfColumn, header)
	}

	wantOrder := []string{isfTestContaminantColumn, isfColumn, isfParentPeptideColumn, isfTestUniqueColumn}
	gotOrder := header[idx-1 : idx+3]
	if strings.Join(gotOrder, "|") != strings.Join(wantOrder, "|") {
		t.Errorf("ISF columns are in %q, want %q", gotOrder, wantOrder)
	}

	tests := []struct {
		name              string
		peptide           string
		wantIsISF         string
		wantParentPeptide string
	}{
		{
			name:              "in-source fragment PSM reports its parent modified peptide",
			peptide:           isfTestFragmentPeptide,
			wantIsISF:         "true",
			wantParentPeptide: isfTestParentPeptide,
		},
		{
			name:              "other PSM has an empty parent column",
			peptide:           isfTestNonFragmentPeptide,
			wantIsISF:         "false",
			wantParentPeptide: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row, ok := rows[tt.peptide]
			if !ok {
				t.Fatalf("psm.tsv lacks the PSM of %s", tt.peptide)
			}
			if row[idx] != tt.wantIsISF {
				t.Errorf("%s = %q, want %q", isfColumn, row[idx], tt.wantIsISF)
			}
			if row[idx+1] != tt.wantParentPeptide {
				t.Errorf("%s = %q, want %q", isfParentPeptideColumn, row[idx+1], tt.wantParentPeptide)
			}
			if row[idx-1] != "false" || row[idx+2] != "true" {
				t.Errorf("neighbouring columns are %q and %q, want Is Contaminant false and Is Unique true", row[idx-1], row[idx+2])
			}
		})
	}
}

func TestPSMReport_NoISFColumnsWithoutISF(t *testing.T) {

	evi := PSMEvidenceList{
		isfTestPSM(isfTestNonFragmentSpectrum, isfTestNonFragmentPeptide, ""),
	}

	header, _ := writeAndReadPSMReport(t, evi)

	for _, column := range []string{isfColumn, isfParentPeptideColumn} {
		if indexOf(header, column) >= 0 {
			t.Errorf("psm.tsv header has %q although no PSM is an in-source fragment", column)
		}
	}

	idx := indexOf(header, isfTestContaminantColumn)
	if idx < 0 || header[idx+1] != isfTestUniqueColumn {
		t.Errorf("Is Contaminant is not directly followed by Is Unique: %q", header)
	}
}

func TestAssemblePSMReport_CopiesISFParent(t *testing.T) {

	pep := id.PepIDList{
		{
			Spectrum:         isfTestFragmentSpectrum,
			SpectrumFile:     "a.pepXML",
			Peptide:          isfTestFragmentPeptide,
			Protein:          "sp|P1|TEST",
			AssumedCharge:    2,
			ISFParentPeptide: isfTestParentPeptide,
		},
		{
			Spectrum:      isfTestNonFragmentSpectrum,
			SpectrumFile:  "a.pepXML",
			Peptide:       isfTestNonFragmentPeptide,
			Protein:       "sp|P1|TEST",
			AssumedCharge: 2,
		},
	}

	var evi Evidence
	evi.AssemblePSMReport(pep, "rev_")

	got := make(map[string]PSMEvidence, len(evi.PSM))
	for _, p := range evi.PSM {
		got[p.Peptide] = p
	}

	fragment := got[isfTestFragmentPeptide]
	if fragment.ISFParentPeptide != isfTestParentPeptide {
		t.Errorf("fragment PSM parent = %q, want %q", fragment.ISFParentPeptide, isfTestParentPeptide)
	}

	other := got[isfTestNonFragmentPeptide]
	if other.ISFParentPeptide != "" {
		t.Errorf("non-fragment PSM parent = %q, want empty", other.ISFParentPeptide)
	}
}
