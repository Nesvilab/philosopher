package tmt

import (
	"github.com/Nesvilab/philosopher/lib/iso"
)

// New builds a new Labelled spectra object
func New(plex string) iso.Labels {

	var o iso.Labels

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
	o.Channel12.Name = "132N"
	o.Channel13.Name = "132C"
	o.Channel14.Name = "133N"
	o.Channel15.Name = "133C"
	o.Channel16.Name = "134N"
	o.Channel17.Name = "134C"
	o.Channel18.Name = "135N"

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
	o.Channel11.Mz = 131.144500
	o.Channel12.Mz = 132.141535
	o.Channel13.Mz = 132.147855
	o.Channel14.Mz = 133.144890
	o.Channel15.Mz = 133.151210
	o.Channel16.Mz = 134.148245
	o.Channel17.Mz = 134.154565
	o.Channel18.Mz = 135.151600

	return o
}
