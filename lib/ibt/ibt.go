package ibt

import (
	"github.com/Nesvilab/philosopher/lib/iso"
)

// New builds a new Labelled spectra object
func New(plex string) iso.Labels {

	var o iso.Labels

	o.Channel1.Name = "114"
	o.Channel2.Name = "115N"
	o.Channel3.Name = "115C"
	o.Channel4.Name = "116N"
	o.Channel5.Name = "116C"
	o.Channel6.Name = "117N"
	o.Channel7.Name = "117C"
	o.Channel8.Name = "118N"
	o.Channel9.Name = "118C"
	o.Channel10.Name = "119N"
	o.Channel11.Name = "119C"
	o.Channel12.Name = "120N"
	o.Channel13.Name = "120C"
	o.Channel14.Name = "121N"
	o.Channel15.Name = "121C"
	o.Channel16.Name = "122"

	o.Channel1.Mz = 114.1277
	o.Channel2.Mz = 115.1248
	o.Channel3.Mz = 115.1311
	o.Channel4.Mz = 116.1281
	o.Channel5.Mz = 116.1344
	o.Channel6.Mz = 117.1315
	o.Channel7.Mz = 117.1378
	o.Channel8.Mz = 118.1348
	o.Channel9.Mz = 118.1411
	o.Channel10.Mz = 119.1382
	o.Channel11.Mz = 119.1445
	o.Channel12.Mz = 120.1415
	o.Channel13.Mz = 120.1479
	o.Channel14.Mz = 121.1449
	o.Channel15.Mz = 121.1512
	o.Channel16.Mz = 122.1482

	return o
}
