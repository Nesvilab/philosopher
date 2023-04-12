package xta

import (
	"github.com/Nesvilab/philosopher/lib/iso"
)

// New builds a new Labelled spectra object
func New(plex string) iso.Labels {

	var o iso.Labels

	o.Channel1.Name = "xTag1"
	o.Channel2.Name = "xTag2"
	o.Channel3.Name = "xTag3"
	o.Channel4.Name = "xTag4"
	o.Channel5.Name = "xTag5"
	o.Channel6.Name = "xTag6"
	o.Channel7.Name = "xTag7"
	o.Channel8.Name = "xTag8"
	o.Channel9.Name = "xTag9"
	o.Channel10.Name = "xTag10"
	o.Channel11.Name = "xTag11"
	o.Channel12.Name = "xTag12"
	o.Channel13.Name = "xTag13"
	o.Channel14.Name = "xTag14"
	o.Channel15.Name = "xTag15"
	o.Channel16.Name = "xTag16"
	o.Channel17.Name = "xTag17"
	o.Channel18.Name = "xTag18"

	o.Channel1.Mz = 114.1277
	o.Channel2.Mz = 118.1528
	o.Channel3.Mz = 173.1284
	o.Channel4.Mz = 184.1076
	o.Channel5.Mz = 229.1910
	o.Channel6.Mz = 244.1292
	o.Channel7.Mz = 245.1325
	o.Channel8.Mz = 272.1612
	o.Channel9.Mz = 426.28233
	o.Channel10.Mz = 427.28568
	o.Channel11.Mz = 428.28904
	o.Channel12.Mz = 429.29239
	o.Channel13.Mz = 0
	o.Channel14.Mz = 0
	o.Channel15.Mz = 0
	o.Channel16.Mz = 0
	o.Channel17.Mz = 0
	o.Channel18.Mz = 0

	return o
}
