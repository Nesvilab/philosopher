package k2s

import (
	"philosopher/lib/iso"
)

// New builds a new Labelled spectra object
func New(plex string) iso.Labels {

	var o iso.Labels

	o.Channel1.Name = "284"
	o.Channel2.Name = "290"
	o.Channel3.Name = "301"
	o.Channel4.Name = "307"
	o.Channel5.Name = "327"
	o.Channel6.Name = "333"

	o.Channel1.Mz = 284.1427
	o.Channel2.Mz = 290.1804
	o.Channel3.Mz = 301.1693
	o.Channel4.Mz = 307.2069
	o.Channel5.Mz = 327.1849
	o.Channel6.Mz = 333.2226

	return o
}
