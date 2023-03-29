package fil

import (
	"reflect"
	"testing"

	"github.com/Nesvilab/philosopher/lib/id"
	"github.com/Nesvilab/philosopher/lib/sys"
	"github.com/Nesvilab/philosopher/lib/tes"
	"github.com/Nesvilab/philosopher/lib/uti"
)

func Test_readPepXMLInput(t *testing.T) {

	tes.SetupTestEnv()
	var pepIDList id.PepIDListPtrs

	type args1 struct {
		xmlFile        string
		decoyTag       string
		temp           string
		models         bool
		calibratedMass int
	}
	test1 := []struct {
		name  string
		args  args1
		want  int
		want1 string
	}{
		{
			name:  "Testting pepXML reading and formating for the filter",
			args:  args1{xmlFile: "interact.pep.xml", decoyTag: "rev_", temp: sys.GetTemp(), models: false, calibratedMass: 0},
			want:  12517,
			want1: "MSFragger",
		},
	}

	for _, tt := range test1 {

		t.Run(tt.name, func(t *testing.T) {

			got, got1 := id.ReadPepXMLInput(tt.args.xmlFile, tt.args.decoyTag, tt.args.temp, tt.args.models)
			pepIDList = got

			if !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("readPepXMLInput() got = %v, want %v", len(got), tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("readPepXMLInput() got1 = %v, want %v", got1, tt.want1)
			}

			if got[0].Index != uint32(12107) {
				t.Errorf("Index is incorrect, got %d, want %d", got[0].Index, uint32(12107))
			}

			if got[0].SpectrumFileName().Str() != "z04397_tc-o238g-setB_MS3.26373.26373.2#interact.pep.xml" {
				t.Errorf("Spectrum is incorrect, got %s, want %s", got[0].SpectrumFileName().Str(), "z04397_tc-o238g-setB_MS3.55183.55183.3#interact.pep.xml")
			}

			if uti.ToFixed(got[0].PrecursorNeutralMass, 4) != 1618.8370 {
				t.Errorf("PrecursorNeutralMass is incorrect, got %f, want %f", got[0].PrecursorNeutralMass, 1618.8370)
			}

			if uti.ToFixed(got[0].RetentionTime, 2) != 4044.46 {
				t.Errorf("RetentionTime is incorrect, got %f, want %f", got[0].RetentionTime, 4044.46)
			}

			if got[0].Peptide != "QSLEASLAETEGR" {
				t.Errorf("Peptide is incorrect, got %s, want %s", got[0].Peptide, "QSLEASLAETEGR")
			}

			if uti.ToFixed(got[0].Massdiff, 2) != 0.00 {
				t.Errorf("Massdiff is incorrect, got %.2f, want %.2f", uti.ToFixed(got[0].Massdiff, 2), 0.00)
			}

			if uti.ToFixed(got[0].CalcNeutralPepMass, 4) != 1618.8364 {
				t.Errorf("CalcNeutralPepMass is incorrect, got %f, want %f", got[0].CalcNeutralPepMass, 1618.8364)
			}

			if got[0].NumberofMissedCleavages != 0 {
				t.Errorf("NumberofMissedCleavages is incorrect, got %d, want %d", got[0].NumberofMissedCleavages, 0)
			}

			if got[0].Protein != "contam_sp|P13645|K1C10_HUMAN" {
				t.Errorf("Protein is incorrect, got %s, want %s", got[0].Protein, "contam_sp|P13645|K1C10_HUMAN")
			}

			if got[0].Probability != 1.0000 {
				t.Errorf("Probability is incorrect, got %f, want %f", got[0].Probability, 1.0000)
			}

		})
	}

	type args2 struct {
		decoyTag string
		psm      float64
		peptide  float64
		ion      float64
	}
	test2 := []struct {
		name  string
		args  args2
		want  float64
		want1 float64
		want2 float64
	}{
		{
			name:  "Testting pepXML reading and formating for the filter",
			args:  args2{decoyTag: "rev_", psm: 0.01, peptide: 0.01, ion: 0.01},
			want:  0.80566,
			want1: 0.820338,
			want2: 0.809078,
		},
	}

	for _, tt := range test2 {

		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := processPeptideIdentifications(pepIDList, tt.args.decoyTag, "", tt.args.psm, tt.args.peptide, tt.args.ion, false, false)
			if got != tt.want {
				t.Errorf("processPeptideIdentifications(psm) got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("processPeptideIdentifications(peptide) got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("processPeptideIdentifications(ion) got2 = %v, want %v", got2, tt.want2)
			}
		})
	}

	type args3 struct {
		charge   uint8
		decoyTag string
	}
	test3 := []struct {
		name  string
		args  args3
		wantT int
		wantD int
	}{
		{
			name:  "Testing charge state 1 profile",
			args:  args3{charge: uint8(1), decoyTag: "rev_"},
			wantT: 0,
			wantD: 0,
		},
		{
			name:  "Testing charge state 2 profile",
			args:  args3{charge: uint8(2), decoyTag: "rev_"},
			wantT: 5748,
			wantD: 203,
		},
		{
			name:  "Testing charge state 3 profile",
			args:  args3{charge: uint8(3), decoyTag: "rev_"},
			wantT: 5482,
			wantD: 178,
		},
		{
			name:  "Testing charge state 4 profile",
			args:  args3{charge: uint8(4), decoyTag: "rev_"},
			wantT: 808,
			wantD: 48,
		},
		{
			name:  "Testing charge state 5 profile",
			args:  args3{charge: uint8(5), decoyTag: "rev_"},
			wantT: 37,
			wantD: 5,
		},
	}

	for _, tt := range test3 {
		t.Run(tt.name, func(t *testing.T) {
			gotT, gotD := chargeProfile(pepIDList, tt.args.charge, tt.args.decoyTag)
			if gotT != tt.wantT {
				t.Errorf("chargeProfile() gotT = %v, want %v", gotT, tt.wantT)
			}
			if gotD != tt.wantD {
				t.Errorf("chargeProfile() gotD = %v, want %v", gotD, tt.wantD)
			}
		})
	}

	test4 := []struct {
		name string
		want int
	}{
		{
			name: "Testing the generation of Unique PSMs",
			want: 12517,
		},
	}

	for _, tt := range test4 {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUniquePSMs(pepIDList); !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("GetUniquePSMs() = %v, want %v", len(got), tt.want)
			}
		})
	}

	test5 := []struct {
		name string
		want int
	}{
		{
			name: "Testing the generation of Unique Ions",
			want: 11685,
		},
	}

	for _, tt := range test5 {
		t.Run(tt.name, func(t *testing.T) {
			if got := getUniquePeptideIons(pepIDList); !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("getUniquePeptideIons() = %v, want %v", len(got), tt.want)
			}
		})
	}

	test6 := []struct {
		name string
		want int
	}{
		{
			name: "Testing the generation of Unique Peptides",
			want: 10794,
		},
	}

	for _, tt := range test6 {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUniquePeptides(pepIDList); !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("GetUniquePeptides() = %v, want %v", len(got), tt.want)
			}
		})
	}

	test7 := []struct {
		name string
		want int
	}{
		{
			name: "Testing the Ion extraction from PSM",
			want: 11685,
		},
	}

	for _, tt := range test7 {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractIonsFromPSMs(pepIDList); !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("ExtractIonsFromPSMs() = %v, want %v", len(got), tt.want)
			}
		})
	}

	//tes.ShutDowTestEnv()
}

func Test_readProtXMLInput(t *testing.T) {

	tes.SetupTestEnv()
	var proXML id.ProtXML

	type args1 struct {
		//meta     string
		xmlFile  string
		decoyTag string
		weight   float64
	}
	test1 := []struct {
		name string
		args args1
		want int
	}{
		{
			name: "Testting protXML reading and formating for the filter",
			args: args1{xmlFile: "interact.prot.xml", decoyTag: "rev_", weight: 1.00},
			want: 2358,
		},
	}

	for _, tt := range test1 {
		t.Run(tt.name, func(t *testing.T) {

			got := ReadProtXMLInput(tt.args.xmlFile, tt.args.decoyTag, tt.args.weight)
			proXML = got

			if len(got.Groups) != tt.want {
				t.Errorf("readProtXMLInput() = %v, want %v", len(got.Groups), tt.want)
			}
		})
	}

	test2 := []struct {
		name  string
		wantT int
		wantD int
	}{
		{
			name:  "Testing Protein Profile",
			wantT: 1937,
			wantD: 421,
		},
	}
	for _, tt := range test2 {
		t.Run(tt.name, func(t *testing.T) {
			gotT, gotD := proteinProfile(proXML)
			if gotT != tt.wantT {
				t.Errorf("proteinProfile() gotT = %v, want %v", gotT, tt.wantT)
			}
			if gotD != tt.wantD {
				t.Errorf("proteinProfile() gotD = %v, want %v", gotD, tt.wantD)
			}
		})
	}

	type args3 struct {
		ptFDR    float64
		pepProb  float64
		protProb float64
		isPicked bool
		isRazor  bool
		fo       bool
		decoyTag string
	}
	test3 := []struct {
		name string
		args args3
	}{
		{
			name: "Testing Protein Identifications",
			args: args3{ptFDR: 0.01, pepProb: 0.7, protProb: 0.5, isPicked: false, isRazor: true, fo: false, decoyTag: "rev_"},
		},
	}
	for _, tt := range test3 {
		t.Run(tt.name, func(t *testing.T) {
			ProcessProteinIdentifications(proXML, tt.args.ptFDR, tt.args.pepProb, tt.args.protProb, tt.args.isPicked, tt.args.isRazor, false, tt.args.decoyTag)
		})
	}
}
