package scl

import (
	"github.com/Nesvilab/philosopher/lib/iso"
)

// New builds a new Labelled spectra object
func New(plex string) iso.Labels {

	var o iso.Labels

	o.Channel1.Name = "sCLIP1"
	o.Channel2.Name = "sCLIP2"
	o.Channel3.Name = "sCLIP3"
	o.Channel4.Name = "sCLIP4"
	o.Channel5.Name = "sCLIP5"
	o.Channel6.Name = "sCLIP6"
	o.Channel7.Name = "sCLIP7"
	o.Channel8.Name = "sCLIP8"
	o.Channel9.Name = "sCLIP9"
	o.Channel10.Name = "sCLIP10"
	o.Channel11.Name = "sCLIP11"
	o.Channel12.Name = "sCLIP12"
	o.Channel13.Name = "sCLIP13"
	o.Channel14.Name = "sCLIP14"
	o.Channel15.Name = "sCLIP15"
	o.Channel16.Name = "sCLIP16"
	o.Channel17.Name = "sCLIP17"
	o.Channel18.Name = "sCLIP18"

	o.Channel1.Mz = 114.1277
	o.Channel2.Mz = 118.1528
	o.Channel3.Mz = 173.1284
	o.Channel4.Mz = 184.1076
	o.Channel5.Mz = 229.1910
	o.Channel6.Mz = 244.1292
	o.Channel7.Mz = 245.1325
	o.Channel8.Mz = 272.1612
	o.Channel9.Mz = 300.1918
	o.Channel10.Mz = 301.1888
	o.Channel11.Mz = 301.1951
	o.Channel12.Mz = 302.1922
	o.Channel13.Mz = 302.1960
	o.Channel14.Mz = 302.1985
	o.Channel15.Mz = 426.2823
	o.Channel16.Mz = 427.2856
	o.Channel17.Mz = 428.2890
	o.Channel18.Mz = 429.2923

	return o
}
