package scl

import (
	"github.com/Nesvilab/philosopher/lib/iso"
)

// New builds a new Labelled spectra object
func New(plex string) iso.Labels {

	var o iso.Labels

	o.Channel1.Name = "300"
	o.Channel2.Name = "301N"
	o.Channel3.Name = "301C"
	o.Channel4.Name = "302N"
	o.Channel5.Name = "302O"
	o.Channel6.Name = "302C"

	o.Channel1.Mz = 300.1918
	o.Channel2.Mz = 301.1888
	o.Channel3.Mz = 301.1951
	o.Channel4.Mz = 302.1922
	o.Channel5.Mz = 302.1960
	o.Channel6.Mz = 302.1985

	return o
}
