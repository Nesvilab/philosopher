package tmt

import (
	"philosopher/lib/iso"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		plex string
	}
	tests := []struct {
		name string
		args args
		want iso.Labels
	}{
		{
			name: "Testting 16 plex",
			args: args{plex: "16"},
			want: iso.Labels{
				Channel1: iso.Channel1{
					Name: "126",
					Mz:   126.127726,
				},
				Channel2: iso.Channel2{
					Name: "127N",
					Mz:   127.124761,
				},
				Channel3: iso.Channel3{
					Name: "127C",
					Mz:   127.131081,
				},
				Channel4: iso.Channel4{
					Name: "128N",
					Mz:   128.128116,
				},
				Channel5: iso.Channel5{
					Name: "128C",
					Mz:   128.134436,
				},
				Channel6: iso.Channel6{
					Name: "129N",
					Mz:   129.131471,
				},
				Channel7: iso.Channel7{
					Name: "129C",
					Mz:   129.137790,
				},
				Channel8: iso.Channel8{
					Name: "130N",
					Mz:   130.134825,
				},
				Channel9: iso.Channel9{
					Name: "130C",
					Mz:   130.141145,
				},
				Channel10: iso.Channel10{
					Name: "131N",
					Mz:   131.138180,
				},

				Channel11: iso.Channel11{
					Name: "131C",
					Mz:   131.144499,
				},

				Channel12: iso.Channel12{
					Name: "132N",
					Mz:   132.141535,
				},

				Channel13: iso.Channel13{
					Name: "132C",
					Mz:   132.147855,
				},

				Channel14: iso.Channel14{
					Name: "133N",
					Mz:   133.144890,
				},

				Channel15: iso.Channel15{
					Name: "133C",
					Mz:   133.151210,
				},

				Channel16: iso.Channel16{
					Name: "134N",
					Mz:   134.148245,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.plex); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
