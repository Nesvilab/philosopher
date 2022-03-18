package k2s

import (
	"philosopher/lib/iso"
)

// New builds a new Labelled spectra object
func New(plex string) iso.Labels {

	var o iso.Labels

	o.Channel1.Name = "284"
	o.Channel2.Name = "290"

	o.Channel1.Mz = 284.1427
	o.Channel2.Mz = 290.1804

	return o
}
