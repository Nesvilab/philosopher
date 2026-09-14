package id

import (
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/Nesvilab/philosopher/lib/mod"
	"github.com/Nesvilab/philosopher/lib/spc"
)

// spectrumQueryTemplate is a minimal MSFragger spectrum_query. The %s slot receives the
// search_score elements under test.
const spectrumQueryTemplate = `<spectrum_query spectrum="a.00001.00001.2" start_scan="1" end_scan="1" precursor_neutral_mass="1000.5" assumed_charge="2" index="1" retention_time_sec="60.0">
  <search_result>
    <search_hit hit_rank="1" peptide="PEPTIDEK" peptide_prev_aa="K" peptide_next_aa="A" protein="sp|P1|TEST" num_tot_proteins="1" num_matched_ions="10" tot_num_ions="14" calc_neutral_pep_mass="1000.5" massdiff="0.0" num_tol_term="2" num_missed_cleavages="0" num_matched_peptides="3">
      %s
    </search_hit>
  </search_result>
</spectrum_query>`

func parseSpectrumQuery(t *testing.T, scores string) PeptideIdentification {
	t.Helper()

	var sq spc.SpectrumQuery
	if err := xml.Unmarshal([]byte(fmt.Sprintf(spectrumQueryTemplate, scores)), &sq); err != nil {
		t.Fatalf("unmarshal spectrum_query: %v", err)
	}

	return processSpectrumQuery(sq, mod.Modifications{}, "rev_", "a.pepXML")
}

func TestProcessSpectrumQuery_MSFraggerSearchScores(t *testing.T) {

	tests := []struct {
		name           string
		scores         string
		wantHyperscore float64
		wantNextscore  float64
		wantBCS        int
		wantFINterm    float64
	}{
		{
			name: "bcs and fI_nterm are read from the search scores",
			scores: `<search_score name="hyperscore" value="20.5"/>
      <search_score name="nextscore" value="10.25"/>
      <search_score name="expect" value="1.500000e-03"/>
      <search_score name="bcs" value="5"/>
      <search_score name="fI_nterm" value="0.1234"/>`,
			wantHyperscore: 20.5,
			wantNextscore:  10.25,
			wantBCS:        5,
			wantFINterm:    0.1234,
		},
		{
			name: "fI_nterm of zero is kept as zero",
			scores: `<search_score name="hyperscore" value="20.5"/>
      <search_score name="nextscore" value="10.25"/>
      <search_score name="bcs" value="0"/>
      <search_score name="fI_nterm" value="0.0000"/>`,
			wantHyperscore: 20.5,
			wantNextscore:  10.25,
			wantBCS:        0,
			wantFINterm:    0,
		},
		{
			name: "older pepXML without bcs or fI_nterm defaults both to zero",
			scores: `<search_score name="hyperscore" value="20.5"/>
      <search_score name="nextscore" value="10.25"/>
      <search_score name="expect" value="1.500000e-03"/>`,
			wantHyperscore: 20.5,
			wantNextscore:  10.25,
			wantBCS:        0,
			wantFINterm:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			psm := parseSpectrumQuery(t, tt.scores)

			if psm.Hyperscore != tt.wantHyperscore {
				t.Errorf("Hyperscore = %v, want %v", psm.Hyperscore, tt.wantHyperscore)
			}
			if psm.Nextscore != tt.wantNextscore {
				t.Errorf("Nextscore = %v, want %v", psm.Nextscore, tt.wantNextscore)
			}
			if psm.BCS != tt.wantBCS {
				t.Errorf("BCS = %v, want %v", psm.BCS, tt.wantBCS)
			}
			if psm.FINterm != tt.wantFINterm {
				t.Errorf("FINterm = %v, want %v", psm.FINterm, tt.wantFINterm)
			}
		})
	}
}
