package tmt

import (
	"github.com/prvst/cmsl/err"
)

// Labels main struct
type Labels struct {
	Spectrum      string
	Index         string
	Scan          string
	RetentionTime float64
	ChargeState   int
	Channel1      Channel1
	Channel2      Channel2
	Channel3      Channel3
	Channel4      Channel4
	Channel5      Channel5
	Channel6      Channel6
	Channel7      Channel7
	Channel8      Channel8
	Channel9      Channel9
	Channel10     Channel10
	Channel11     Channel11
}

// LabeledSpectra is a list of spectra lables
type LabeledSpectra map[string]Labels

// Channel1 TMT
type Channel1 struct {
	Name      string
	Mz        float64
	Intensity float64
}

// Channel2 TMT
type Channel2 struct {
	Name      string
	Mz        float64
	Intensity float64
}

// Channel3 TMT
type Channel3 struct {
	Name      string
	Mz        float64
	Intensity float64
}

// Channel4 TMT
type Channel4 struct {
	Name      string
	Mz        float64
	Intensity float64
}

// Channel5 TMT
type Channel5 struct {
	Name      string
	Mz        float64
	Intensity float64
}

// Channel6 TMT
type Channel6 struct {
	Name      string
	Mz        float64
	Intensity float64
}

// Channel7 TMT
type Channel7 struct {
	Name      string
	Mz        float64
	Intensity float64
}

// Channel8 TMT
type Channel8 struct {
	Name      string
	Mz        float64
	Intensity float64
}

// Channel9 TMT
type Channel9 struct {
	Name      string
	Mz        float64
	Intensity float64
}

// Channel10 TMT
type Channel10 struct {
	Name      string
	Mz        float64
	Intensity float64
}

// Channel11 TMT
type Channel11 struct {
	Name      string
	Mz        float64
	Intensity float64
}

// New builds a new Labelled spectra object
func New(plex string) (Labels, *err.Error) {

	var o Labels

	o.Channel1.Name = "126"
	o.Channel2.Name = "127N"
	o.Channel3.Name = "127C"
	o.Channel4.Name = "128N"
	o.Channel5.Name = "128C"
	o.Channel6.Name = "129N"
	o.Channel7.Name = "129C"
	o.Channel8.Name = "130N"
	o.Channel9.Name = "130C"
	o.Channel10.Name = "131N"
	o.Channel11.Name = "131C"

	if plex == "6" {
		o.Channel1.Mz = 126.127726
		o.Channel2.Mz = 127.124761
		o.Channel3.Mz = 127.131081
		o.Channel4.Mz = 128.128116
		o.Channel5.Mz = 128.134436
		o.Channel6.Mz = 129.131471
	} else if plex == "10" {
		o.Channel1.Mz = 126.127726
		o.Channel2.Mz = 127.124761
		o.Channel3.Mz = 127.131081
		o.Channel4.Mz = 128.128116
		o.Channel5.Mz = 128.134436
		o.Channel6.Mz = 129.131471
		o.Channel7.Mz = 129.137790
		o.Channel8.Mz = 130.134825
		o.Channel9.Mz = 130.141145
		o.Channel10.Mz = 131.138180
	} else if plex == "11" {
		o.Channel1.Mz = 126.127726
		o.Channel2.Mz = 127.124761
		o.Channel3.Mz = 127.131081
		o.Channel4.Mz = 128.128116
		o.Channel5.Mz = 128.134436
		o.Channel6.Mz = 129.131471
		o.Channel7.Mz = 129.137790
		o.Channel8.Mz = 130.134825
		o.Channel9.Mz = 130.141145
		o.Channel10.Mz = 131.138180
		o.Channel11.Mz = 131.144499
	} else {
		return o, &err.Error{Type: err.UnknownMultiplex, Class: err.FATA}
	}

	return o, nil
}
