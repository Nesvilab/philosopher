package ktag2

import (
	"philosopher/lib/iso"
)

// New builds a new Labelled spectra object
func New(plex string) iso.Labels {

	var o iso.Labels

	o.Channel1.Name = "286"
	o.Channel2.Name = "290"

	o.Channel1.Mz = 286.1761
	o.Channel2.Mz = 290.1832

	return o
}
