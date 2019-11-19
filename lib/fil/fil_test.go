package fil

import (
	"os"
	"philosopher/lib/sys"
	"philosopher/lib/wrk"
	"reflect"
	"testing"
)

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
			args:  args{xmlFile: "/workspace/philosopher/test/wrksp/interact.pep.xml", decoyTag: "rev_", temp: sys.GetTemp(), models: false, calibratedMass: 0},
			want:  64406,
			want1: "MSFragger",
		},
	}
	for _, tt := range tests {

		os.Chdir("../../test/")
		wrk.Init("0000", "0000")

		t.Run(tt.name, func(t *testing.T) {

			got, got1 := readPepXMLInput(tt.args.xmlFile, tt.args.decoyTag, tt.args.temp, tt.args.models, tt.args.calibratedMass)

			if !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("readPepXMLInput() got = %v, want %v", len(got), tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("readPepXMLInput() got1 = %v, want %v", got1, tt.want1)
			}

			if got[0].Index != uint32(18992) {
				t.Errorf("Meta path or name is incorrect, got %d, want %d", got[0].Index, uint32(18992))
			}

			if got[0].Spectrum != "b1906_293T_proteinID_01A_QE3_122212.60782.60782.2" {
				t.Errorf("Meta path or name is incorrect, got %s, want %s", got[0].Spectrum, "b1906_293T_proteinID_01A_QE3_122212.60782.60782.2")
			}

			if got[0].Scan != 60782 {
				t.Errorf("Meta path or name is incorrect, got %d, want %d", got[0].Scan, 60782)
			}

		})
	}
}
