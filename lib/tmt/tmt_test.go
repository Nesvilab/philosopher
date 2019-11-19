package tmt

import (
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
		want Labels
	}{
		{
			name: "Testting 10 plex",
			args: args{plex: "10"},
			want: Labels{
				Channel1: Channel1{
					Name: "126",
					Mz:   126.127726,
				},
				Channel2: Channel2{
					Name: "127N",
					Mz:   127.124761,
				},
				Channel3: Channel3{
					Name: "127C",
					Mz:   127.131081,
				},
				Channel4: Channel4{
					Name: "128N",
					Mz:   128.128116,
				},
				Channel5: Channel5{
					Name: "128C",
					Mz:   128.134436,
				},
				Channel6: Channel6{
					Name: "129N",
					Mz:   129.131471,
				},
				Channel7: Channel7{
					Name: "129C",
					Mz:   129.137790,
				},
				Channel8: Channel8{
					Name: "130N",
					Mz:   130.134825,
				},
				Channel9: Channel9{
					Name: "130C",
					Mz:   130.141145,
				},
				Channel10: Channel10{
					Name: "131N",
					Mz:   131.138180,
				},
			},
		},
		{
			name: "Testting 11 plex",
			args: args{plex: "11"},
			want: Labels{
				Channel1: Channel1{
					Name: "126",
					Mz:   126.127726,
				},
				Channel2: Channel2{
					Name: "127N",
					Mz:   127.124761,
				},
				Channel3: Channel3{
					Name: "127C",
					Mz:   127.131081,
				},
				Channel4: Channel4{
					Name: "128N",
					Mz:   128.128116,
				},
				Channel5: Channel5{
					Name: "128C",
					Mz:   128.134436,
				},
				Channel6: Channel6{
					Name: "129N",
					Mz:   129.131471,
				},
				Channel7: Channel7{
					Name: "129C",
					Mz:   129.137790,
				},
				Channel8: Channel8{
					Name: "130N",
					Mz:   130.134825,
				},
				Channel9: Channel9{
					Name: "130C",
					Mz:   130.141145,
				},
				Channel10: Channel10{
					Name: "131N",
					Mz:   131.138180,
				},
				Channel11: Channel11{
					Name: "131C",
					Mz:   131.144499,
				},
			},
		},
		{
			name: "Testting 16 plex",
			args: args{plex: "16"},
			want: Labels{
				Channel1: Channel1{
					Name: "126",
					Mz:   126.127726,
				},
				Channel2: Channel2{
					Name: "127N",
					Mz:   127.124761,
				},
				Channel3: Channel3{
					Name: "127C",
					Mz:   127.131081,
				},
				Channel4: Channel4{
					Name: "128N",
					Mz:   128.128116,
				},
				Channel5: Channel5{
					Name: "128C",
					Mz:   128.134436,
				},
				Channel6: Channel6{
					Name: "129N",
					Mz:   129.131471,
				},
				Channel7: Channel7{
					Name: "129C",
					Mz:   129.137790,
				},
				Channel8: Channel8{
					Name: "130N",
					Mz:   130.134825,
				},
				Channel9: Channel9{
					Name: "130C",
					Mz:   130.141145,
				},
				Channel10: Channel10{
					Name: "131N",
					Mz:   131.138180,
				},

				Channel11: Channel11{
					Name: "131C",
					Mz:   131.144499,
				},

				Channel12: Channel12{
					Name: "132N",
					Mz:   132.141535,
				},

				Channel13: Channel13{
					Name: "132C",
					Mz:   132.147855,
				},

				Channel14: Channel14{
					Name: "133N",
					Mz:   133.144890,
				},

				Channel15: Channel15{
					Name: "133C",
					Mz:   133.151210,
				},

				Channel16: Channel16{
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
