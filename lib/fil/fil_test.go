package fil

import (
	"os"
	"philosopher/lib/id"
	"philosopher/lib/sys"
	"philosopher/lib/uti"
	"philosopher/lib/wrk"
	"reflect"
	"testing"
)

var pepID id.PepIDList
var proID id.ProtIDList

func Test_readPepXMLInput(t *testing.T) {
	type args struct {
		xmlFile        string
		decoyTag       string
		temp           string
		models         bool
		calibratedMass int
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 string
	}{
		{
			name:  "Testting pepXML reading and formating for the filter",
			args:  args{xmlFile: "interact.pep.xml", decoyTag: "rev_", temp: sys.GetTemp(), models: false, calibratedMass: 0},
			want:  64406,
			want1: "MSFragger",
		},
	}
	for _, tt := range tests {

		os.Chdir("../../test/wrksp/")
		wrk.Init("0000", "0000")

		t.Run(tt.name, func(t *testing.T) {

			got, got1 := readPepXMLInput(tt.args.xmlFile, tt.args.decoyTag, tt.args.temp, tt.args.models, tt.args.calibratedMass)
			pepID = got

			if !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("readPepXMLInput() got = %v, want %v", len(got), tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("readPepXMLInput() got1 = %v, want %v", got1, tt.want1)
			}

			if got[0].Index != uint32(18992) {
				t.Errorf("Index is incorrect, got %d, want %d", got[0].Index, uint32(18992))
			}

			if got[0].Spectrum != "b1906_293T_proteinID_01A_QE3_122212.60782.60782.2" {
				t.Errorf("Spectrum is incorrect, got %s, want %s", got[0].Spectrum, "b1906_293T_proteinID_01A_QE3_122212.60782.60782.2")
			}

			if got[0].Scan != 60782 {
				t.Errorf("Scan is incorrect, got %d, want %d", got[0].Scan, 60782)
			}

			if got[0].PrecursorNeutralMass != 1429.7663 {
				t.Errorf("PrecursorNeutralMass is incorrect, got %f, want %f", got[0].PrecursorNeutralMass, 1429.7663)
			}

			if got[0].RetentionTime != 11202.398 {
				t.Errorf("RetentionTime is incorrect, got %f, want %f", got[0].RetentionTime, 11202.398)
			}

			if got[0].Peptide != "LEESADNILSIVK" {
				t.Errorf("Peptide is incorrect, got %s, want %s", got[0].Peptide, "LEESADNILSIVK")
			}

			if uti.ToFixed(got[0].Massdiff, 2) != 0.00 {
				t.Errorf("Massdiff is incorrect, got %.2f, want %.2f", uti.ToFixed(got[0].Massdiff, 2), 0.00)
			}

			if got[0].CalcNeutralPepMass != 1429.7664 {
				t.Errorf("CalcNeutralPepMass is incorrect, got %.2f, want %.2f", got[0].CalcNeutralPepMass, 1429.7664)
			}

			if got[0].NextAA != "Q" {
				t.Errorf("NextAA is incorrect, got %s, want %s", got[0].NextAA, "Q")
			}

			if got[0].NumberofMissedCleavages != 0 {
				t.Errorf("NumberofMissedCleavages is incorrect, got %d, want %d", got[0].NumberofMissedCleavages, 0)
			}

			if got[0].Protein != "sp|O00287|RFXAP_HUMAN" {
				t.Errorf("Protein is incorrect, got %s, want %s", got[0].Protein, "sp|O00287|RFXAP_HUMAN")
			}

			if got[0].Probability != 1.0000 {
				t.Errorf("Probability is incorrect, got %f, want %f", got[0].Probability, 1.0000)
			}

		})
	}
}

func Test_processPeptideIdentifications(t *testing.T) {
	type args struct {
		p        id.PepIDList
		decoyTag string
		psm      float64
		peptide  float64
		ion      float64
	}
	tests := []struct {
		name  string
		args  args
		want  float64
		want1 float64
		want2 float64
	}{
		{
			name:  "Testting pepXML reading and formating for the filter",
			args:  args{p: pepID, decoyTag: "rev_", psm: 0.01, peptide: 0.01, ion: 0.01},
			want:  0.1914,
			want1: 0.723,
			want2: 0.5155,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := processPeptideIdentifications(tt.args.p, tt.args.decoyTag, tt.args.psm, tt.args.peptide, tt.args.ion)
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
}

func Test_chargeProfile(t *testing.T) {
	type args struct {
		p        id.PepIDList
		charge   uint8
		decoyTag string
	}
	tests := []struct {
		name  string
		args  args
		wantT int
		wantD int
	}{
		{
			name:  "Testing charge state 1 profile",
			args:  args{p: pepID, charge: uint8(1), decoyTag: "rev_"},
			wantT: 0,
			wantD: 0,
		},
		{
			name:  "Testing charge state 2 profile",
			args:  args{p: pepID, charge: uint8(2), decoyTag: "rev_"},
			wantT: 36174,
			wantD: 457,
		},
		{
			name:  "Testing charge state 3 profile",
			args:  args{p: pepID, charge: uint8(3), decoyTag: "rev_"},
			wantT: 22656,
			wantD: 317,
		},
		{
			name:  "Testing charge state 4 profile",
			args:  args{p: pepID, charge: uint8(4), decoyTag: "rev_"},
			wantT: 4272,
			wantD: 88,
		},
		{
			name:  "Testing charge state 5 profile",
			args:  args{p: pepID, charge: uint8(5), decoyTag: "rev_"},
			wantT: 432,
			wantD: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT, gotD := chargeProfile(tt.args.p, tt.args.charge, tt.args.decoyTag)
			if gotT != tt.wantT {
				t.Errorf("chargeProfile() gotT = %v, want %v", gotT, tt.wantT)
			}
			if gotD != tt.wantD {
				t.Errorf("chargeProfile() gotD = %v, want %v", gotD, tt.wantD)
			}
		})
	}
}

func TestGetUniquePSMs(t *testing.T) {
	type args struct {
		p id.PepIDList
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Testing the generation of Unique PSMs",
			args: args{pepID},
			want: 64406,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUniquePSMs(tt.args.p); !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("GetUniquePSMs() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func Test_getUniquePeptideIons(t *testing.T) {
	type args struct {
		p id.PepIDList
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Testing the generation of Unique Ions",
			args: args{pepID},
			want: 39716,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getUniquePeptideIons(tt.args.p); !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("getUniquePeptideIons() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestGetUniquePeptides(t *testing.T) {
	type args struct {
		p id.PepIDList
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Testing the generation of Unique Peptides",
			args: args{pepID},
			want: 30092,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUniquePeptides(tt.args.p); !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("GetUniquePeptides() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestExtractIonsFromPSMs(t *testing.T) {
	type args struct {
		p id.PepIDList
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Testing the Ion extraction from PSM",
			args: args{pepID},
			want: 39716,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractIonsFromPSMs(tt.args.p); !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("ExtractIonsFromPSMs() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func Test_readProtXMLInput(t *testing.T) {
	type args struct {
		meta     string
		xmlFile  string
		decoyTag string
		weight   float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Testting protXML reading and formating for the filter",
			args: args{xmlFile: "interact.prot.xml", decoyTag: "rev_", weight: 1.00},
			want: 7926,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readProtXMLInput(tt.args.meta, tt.args.xmlFile, tt.args.decoyTag, tt.args.weight); !reflect.DeepEqual(len(got.Groups), tt.want) {
				t.Errorf("readProtXMLInput() = %v, want %v", len(got.Groups), tt.want)
			}
		})
	}
}
