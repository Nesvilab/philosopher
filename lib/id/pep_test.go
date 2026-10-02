package id

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"sort"
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

	return processSpectrumQuery(sq, newDeclaredModifications(mod.Modifications{}), "rev_", "a.pepXML")
}

// searchSummaryMods is a header table as MSFragger declares it: fixed carbamidomethyl on C, a -43.0058
// mass offset on C reported as a variable modification, oxidation on M and N-terminal acetylation.
func searchSummaryMods() mod.Modifications {
	return mod.Modifications{Index: map[string]mod.Modification{
		"C#160.0307":     {Index: "C#160.0307", Type: mod.Assigned, AminoAcid: "C", MassDiff: 57.0214, Variable: false},
		"C#117.0249":     {Index: "C#117.0249", Type: mod.Assigned, AminoAcid: "C", MassDiff: -43.0058, Variable: true},
		"M#147.0354":     {Index: "M#147.0354", Type: mod.Assigned, AminoAcid: "M", MassDiff: 15.9949, Variable: true},
		"N-term#43.0184": {Index: "N-term#43.0184", Type: mod.Assigned, AminoAcid: "N-term", MassDiff: 42.0106, Variable: true},
	}}
}

// carbamidomethylOffsetMods declares carbamidomethyl twice on C, as the fixed modification and, as the
// FragPipe common-PTM offset list does, as a mass offset stacked on it with a different last decimal.
func carbamidomethylOffsetMods() mod.Modifications {
	return mod.Modifications{Index: map[string]mod.Modification{
		"C#160.0307": {Index: "C#160.0307", Type: mod.Assigned, AminoAcid: "C", MassDiff: 57.0214, Variable: false},
		"C#217.0521": {Index: "C#217.0521", Type: mod.Assigned, AminoAcid: "C", MassDiff: 57.0215, Variable: true},
	}}
}

// tmtSearchSummaryMods is a TMTpro header table with fixed labels on K and on the peptide N-terminus.
// The N-terminal key depends on how the label mass was typed: 304.20715 gives 305.2150, 304.2071 gives
// 305.2149.
func tmtSearchSummaryMods(nTermKey string) func() mod.Modifications {
	return func() mod.Modifications {
		return mod.Modifications{Index: map[string]mod.Modification{
			"K#432.3021": {Index: "K#432.3021", Type: mod.Assigned, AminoAcid: "K", MassDiff: 304.2071, Variable: false},
			nTermKey:     {Index: nTermKey, Type: mod.Assigned, AminoAcid: "N-term", MassDiff: 304.2071, Variable: false},
		}}
	}
}

// methylSearchSummaryMods declares a five-decimal methyl offset (14.01565) on K only, stored the way the
// header loader stores it: uti.ToFixed keeps 14.0156.
func methylSearchSummaryMods() mod.Modifications {
	return mod.Modifications{Index: map[string]mod.Modification{
		"K#142.1106": {Index: "K#142.1106", Type: mod.Assigned, AminoAcid: "K", MassDiff: 14.0156, Variable: true},
	}}
}

// assignedMods renders the assigned modifications of a PSM as "<site>:<massdiff>:<variable>" strings, sorted.
func assignedMods(p PeptideIdentification) []string {
	var out []string
	for _, m := range p.Modifications.IndexSlice {
		if m.Type != mod.Assigned {
			continue
		}
		site := m.AminoAcid
		if m.Position > 0 {
			site = fmt.Sprintf("%d%s", m.Position, m.AminoAcid)
		}
		out = append(out, fmt.Sprintf("%s:%.4f:%t", site, m.MassDiff, m.Variable))
	}
	sort.Strings(out)
	return out
}

func TestMapModsFromPepXML(t *testing.T) {

	tests := []struct {
		name    string
		peptide string
		mods    func() mod.Modifications
		info    spc.ModificationInfo
		want    []string
	}{
		{
			// MSFragger writes the offset as a second mod_aminoacid_mass at the carbamidomethyl site whose
			// total mass (C + 57.02146 - 32.0085 = 128.02216) no search summary header declares.
			name:    "an undeclared mass offset stacked on a fixed modification is derived from the declared entry",
			peptide: "DDDIAALVVDNGSGMCK",
			info: spc.ModificationInfo{
				ModNTermMass:    43.018425,
				ModifiedPeptide: []byte("n[43]DDDIAALVVDNGSGM[147]C[71]K"),
				ModAminoacidMass: []spc.ModAminoacidMass{
					{Position: 15, Mass: 147.0354},
					{Position: 16, Mass: 160.03065},
					{Position: 16, Mass: 128.02216},
				},
			},
			want: []string{"15M:15.9949:true", "16C:-32.0085:true", "16C:57.0214:false", "N-term:42.0106:true"},
		},
		{
			name:    "a declared offset stacked on a fixed modification keeps both tokens",
			peptide: "HQGVMVGMGQKDCYVGDEAQSK",
			info: spc.ModificationInfo{
				ModAminoacidMass: []spc.ModAminoacidMass{
					{Position: 13, Mass: 160.03065},
					{Position: 13, Mass: 117.02486},
				},
			},
			want: []string{"13C:-43.0058:true", "13C:57.0214:false"},
		},
		{
			name:    "an undeclared offset on a residue without fixed modifications is derived from the residue mass",
			peptide: "RSAAEMYGSVTEHPSPSPLLSSSFDLDYDFQR",
			info: spc.ModificationInfo{
				ModAminoacidMass: []spc.ModAminoacidMass{{Position: 14, Mass: 177.01906}},
			},
			want: []string{"14P:79.9663:true"},
		},
		{
			name:    "an undeclared terminal mass is derived from the terminal group",
			peptide: "PEPTIDEK",
			info: spc.ModificationInfo{
				ModNTermMass: 29.039125, // H + dimethyl 28.0313
				ModCTermMass: 16.018724, // OH - 0.984016 (amidation)
			},
			want: []string{"C-term:-0.9840:true", "N-term:28.0313:true"},
		},
		{
			name:    "a mass one 4-decimal bin away from the declared one still maps",
			peptide: "PEPCK",
			info: spc.ModificationInfo{
				ModAminoacidMass: []spc.ModAminoacidMass{{Position: 4, Mass: 160.0308}},
			},
			want: []string{"4C:57.0214:false"},
		},
		{
			// PTMProphet rewrites masses with fewer decimals; that is still the declared oxidation, not a
			// new 15.9945 modification.
			name:    "a mass several bins away from the declared one maps to it instead of a derived variant",
			peptide: "PEPMK",
			info: spc.ModificationInfo{
				ModAminoacidMass: []spc.ModAminoacidMass{{Position: 4, Mass: 147.035}},
			},
			want: []string{"4M:15.9949:true"},
		},
		{
			name:    "a mass 0.0015 Da from the declared one still maps",
			peptide: "PEPMK",
			info: spc.ModificationInfo{
				ModAminoacidMass: []spc.ModAminoacidMass{{Position: 4, Mass: 147.0339}},
			},
			want: []string{"4M:15.9949:true"},
		},
		{
			name:    "a mass beyond the tolerance is a modification of its own",
			peptide: "PEPMK",
			info: spc.ModificationInfo{
				ModAminoacidMass: []spc.ModAminoacidMass{{Position: 4, Mass: 147.0380}},
			},
			want: []string{"4M:15.9975:true"},
		},
		{
			// 117.025 is the declared -43.0058 offset on carbamidomethyl C (stated mass 117.0249), so the
			// stated mass, not the residue plus the delta, decides the match.
			name:    "a declared offset on a fixed-modified residue is found by its stated mass",
			peptide: "PEPCK",
			info: spc.ModificationInfo{
				ModAminoacidMass: []spc.ModAminoacidMass{{Position: 4, Mass: 117.025}},
			},
			want: []string{"4C:-43.0058:true"},
		},
		{
			name:    "a drifted fixed modification maps to the fixed entry, not to a stacked offset of the same delta",
			peptide: "PEPCK",
			mods:    carbamidomethylOffsetMods,
			info: spc.ModificationInfo{
				ModAminoacidMass: []spc.ModAminoacidMass{{Position: 4, Mass: 160.0308}},
			},
			want: []string{"4C:57.0214:false"},
		},
		{
			// 0.0023 Da above carbamidomethyl C: too far for the tolerance, far too small to be a modification.
			name:    "a drift larger than the tolerance but smaller than any modification maps to the declared entry",
			peptide: "PEPCK",
			info: spc.ModificationInfo{
				ModAminoacidMass: []spc.ModAminoacidMass{{Position: 4, Mass: 160.0330}},
			},
			want: []string{"4C:57.0214:false"},
		},
		{
			name:    "a terminal mass within tolerance of the declared one maps to it",
			peptide: "PEPTIDEK",
			info:    spc.ModificationInfo{ModNTermMass: 43.0181},
			want:    []string{"N-term:42.0106:true"},
		},
		{
			// S + 14.01565 = 101.04768. The header stores the K offset as 14.0156, so the S token must read
			// 14.0156 as well, not the independently rounded 14.0157.
			name:    "an undeclared offset is written like the same offset declared on another residue",
			peptide: "PEPSK",
			mods:    methylSearchSummaryMods,
			info: spc.ModificationInfo{
				ModAminoacidMass: []spc.ModAminoacidMass{{Position: 4, Mass: 101.04768}},
			},
			want: []string{"4S:14.0156:true"},
		},
		{
			name:    "a terminus carrying only its hydrogen adds nothing",
			peptide: "PEPTIDEK",
			info:    spc.ModificationInfo{ModNTermMass: 1.00782503207},
			want:    nil,
		},
		{
			name:    "a mass on a non-standard residue is dropped",
			peptide: "PEUCK",
			info: spc.ModificationInfo{
				ModAminoacidMass: []spc.ModAminoacidMass{{Position: 3, Mass: 200.0}},
			},
			want: nil,
		},
		{
			name:    "a position outside the peptide is dropped instead of crashing",
			peptide: "PEPCK",
			info: spc.ModificationInfo{
				ModAminoacidMass: []spc.ModAminoacidMass{{Position: 9, Mass: 160.03065}, {Position: 0, Mass: 160.03065}},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mods := searchSummaryMods
			if tt.mods != nil {
				mods = tt.mods
			}
			p := PeptideIdentification{Peptide: tt.peptide, Massdiff: 22.9714}
			p.mapModsFromPepXML(tt.info, newDeclaredModifications(mods()))

			if got := assignedMods(p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("assigned modifications = %v, want %v", got, tt.want)
			}
		})
	}
}

// PTMProphet rewrites mod_nterm_mass, so the PSM mass can sit in a neighbouring 4-decimal bin of the
// declared TMTpro mass or further away. Every variant must yield exactly one N-terminal label token,
// whichever way the header key was rounded.
func TestMapModsFromPepXML_TMTNTerm(t *testing.T) {

	want := []string{"8K:304.2071:false", "N-term:304.2071:false"}

	for _, key := range []string{"N-term#305.2150", "N-term#305.2149"} {
		declared := newDeclaredModifications(tmtSearchSummaryMods(key)())

		for _, mass := range []float64{305.2150, 305.2149, 305.2151, 305.2140, 305.2290} {
			p := PeptideIdentification{Peptide: "PEPTIDEK"}
			p.mapModsFromPepXML(spc.ModificationInfo{
				ModNTermMass:     mass,
				ModAminoacidMass: []spc.ModAminoacidMass{{Position: 8, Mass: 432.30206}},
			}, declared)

			if got := assignedMods(p); !reflect.DeepEqual(got, want) {
				t.Errorf("header %s, mod_nterm_mass %.4f: assigned modifications = %v, want %v", key, mass, got, want)
			}
		}
	}
}

func TestMapModsFromPepXML_ObservedMassDiff(t *testing.T) {

	p := PeptideIdentification{Peptide: "PEPTIDEK", Massdiff: 22.9714}
	p.mapModsFromPepXML(spc.ModificationInfo{}, newDeclaredModifications(searchSummaryMods()))

	var observed []string
	for _, m := range p.Modifications.IndexSlice {
		if m.Type == mod.Observed {
			observed = append(observed, fmt.Sprintf("%s:%.4f:%s", m.Index, m.MassDiff, m.Name))
		}
	}

	if want := []string{"22.9714:22.9714:Unknown"}; !reflect.DeepEqual(observed, want) {
		t.Errorf("observed modifications = %v, want %v", observed, want)
	}
}

// The stated mass of a header line, not the residue plus the delta, is what a pepXML mass is compared
// with: for a modification stacked on a fixed one the two differ by the fixed delta.
func TestDeclaredModifications_FindDeclared(t *testing.T) {

	declared := newDeclaredModifications(searchSummaryMods())

	tests := []struct {
		site string
		mass float64
		want string
		ok   bool
	}{
		{"C", 160.03065, "C#160.0307", true},
		{"C", 117.025, "C#117.0249", true},
		{"C", 60.0034, "", false}, // C minus 43.0058 without the fixed carbamidomethyl is no declared mass
		{"c", 160.0306, "C#160.0307", true},
		{"n-term", 43.0183, "N-term#43.0184", true},
		{"M", 147.0380, "", false},
	}

	for _, tt := range tests {
		got, ok := declared.findDeclared(tt.site, tt.mass)
		if ok != tt.ok || got.Index != tt.want {
			t.Errorf("findDeclared(%q, %.4f) = %q, %t; want %q, %t", tt.site, tt.mass, got.Index, ok, tt.want, tt.ok)
		}
	}
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

// twoHitSpectrumQueryTemplate is an MSFragger spectrum_query with two search hits. The two %s
// slots receive the search_score elements of hit 1 and hit 2.
const twoHitSpectrumQueryTemplate = `<spectrum_query spectrum="a.00002.00002.2" start_scan="2" end_scan="2" precursor_neutral_mass="800.4" assumed_charge="2" index="2" retention_time_sec="61.0">
  <search_result>
    <search_hit hit_rank="1" peptide="PEPTIDE" peptide_prev_aa="K" peptide_next_aa="K" protein="sp|P1|TEST" num_tot_proteins="1" num_matched_ions="8" tot_num_ions="12" calc_neutral_pep_mass="800.4" massdiff="0.0" num_tol_term="1" num_missed_cleavages="0" num_matched_peptides="3">
      %s
    </search_hit>
    <search_hit hit_rank="2" peptide="TIDEKA" peptide_prev_aa="P" peptide_next_aa="K" protein="sp|P1|TEST" num_tot_proteins="1" num_matched_ions="6" tot_num_ions="10" calc_neutral_pep_mass="800.4" massdiff="0.0" num_tol_term="1" num_missed_cleavages="0" num_matched_peptides="3">
      %s
    </search_hit>
  </search_result>
</spectrum_query>`

func TestProcessSpectrumQuery_ISFParentScores(t *testing.T) {

	tests := []struct {
		name                 string
		scores               string
		wantISFParentPeptide string
	}{
		{
			name: "isf parent peptide is the modified peptide read from the search scores",
			scores: `<search_score name="hyperscore" value="20.5"/>
      <search_score name="bcs" value="5"/>
      <search_score name="fI_nterm" value="0.1234"/>
      <search_score name="isf_parent_peptide" value="n[42.0106]PEPTM[15.9949]IDEK"/>
      <search_score name="isf_parent_charge" value="2"/>
      <search_score name="isf_apex_rt_delta" value="0.0123"/>`,
			wantISFParentPeptide: "n[42.0106]PEPTM[15.9949]IDEK",
		},
		{
			name: "isf parent peptide without modifications is read as is",
			scores: `<search_score name="hyperscore" value="20.5"/>
      <search_score name="isf_parent_peptide" value="PEPTIDEK"/>
      <search_score name="isf_parent_charge" value="3"/>
      <search_score name="isf_apex_rt_delta" value="-0.0500"/>`,
			wantISFParentPeptide: "PEPTIDEK",
		},
		{
			name: "stale isf_parent_modified_peptide score alone is ignored",
			scores: `<search_score name="hyperscore" value="20.5"/>
      <search_score name="isf_parent_modified_peptide" value="n[42.0106]PEPTM[15.9949]IDEK"/>`,
			wantISFParentPeptide: "",
		},
		{
			name: "hit without isf scores has an empty parent",
			scores: `<search_score name="hyperscore" value="20.5"/>
      <search_score name="bcs" value="5"/>`,
			wantISFParentPeptide: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			psm := parseSpectrumQuery(t, tt.scores)

			if psm.ISFParentPeptide != tt.wantISFParentPeptide {
				t.Errorf("ISFParentPeptide = %q, want %q", psm.ISFParentPeptide, tt.wantISFParentPeptide)
			}
		})
	}
}

func TestProcessSpectrumQuery_ISFParentScoresAreResetPerHit(t *testing.T) {

	isfScores := `<search_score name="hyperscore" value="20.5"/>
      <search_score name="isf_parent_peptide" value="PEPTM[15.9949]IDEK"/>
      <search_score name="isf_parent_charge" value="2"/>
      <search_score name="isf_apex_rt_delta" value="0.0123"/>`
	plainScores := `<search_score name="hyperscore" value="10.5"/>`

	tests := []struct {
		name                 string
		firstHitScores       string
		lastHitScores        string
		wantISFParentPeptide string
	}{
		{
			name:                 "isf parent of an earlier hit does not leak into a later non-isf hit",
			firstHitScores:       isfScores,
			lastHitScores:        plainScores,
			wantISFParentPeptide: "",
		},
		{
			name:                 "isf parent of the last hit is kept",
			firstHitScores:       plainScores,
			lastHitScores:        isfScores,
			wantISFParentPeptide: "PEPTM[15.9949]IDEK",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sq spc.SpectrumQuery
			raw := fmt.Sprintf(twoHitSpectrumQueryTemplate, tt.firstHitScores, tt.lastHitScores)
			if err := xml.Unmarshal([]byte(raw), &sq); err != nil {
				t.Fatalf("unmarshal spectrum_query: %v", err)
			}

			psm := processSpectrumQuery(sq, mod.Modifications{}, "rev_", "a.pepXML")

			if psm.ISFParentPeptide != tt.wantISFParentPeptide {
				t.Errorf("ISFParentPeptide = %q, want %q", psm.ISFParentPeptide, tt.wantISFParentPeptide)
			}
		})
	}
}
