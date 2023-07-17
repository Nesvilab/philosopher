package xta2

import (
	"github.com/Nesvilab/philosopher/lib/iso"
)

// New builds a new Labelled spectra object
func New(plex string) iso.Labels {

	var o iso.Labels

	o.Channel1.Name = "114"
	o.Channel2.Name = "115a"
	o.Channel3.Name = "115b"
	o.Channel4.Name = "115c"
	o.Channel5.Name = "116a"
	o.Channel6.Name = "116b"
	o.Channel7.Name = "116c"
	o.Channel8.Name = "116d"
	o.Channel9.Name = "116e"
	o.Channel10.Name = "117a"
	o.Channel11.Name = "117b"
	o.Channel12.Name = "117c"
	o.Channel13.Name = "117d"
	o.Channel14.Name = "117e"
	o.Channel15.Name = "117f"
	o.Channel16.Name = "118a"
	o.Channel17.Name = "118b"
	o.Channel18.Name = "118c"
	o.Channel19.Name = "118d"
	o.Channel20.Name = "118e"
	o.Channel21.Name = "118f"
	o.Channel22.Name = "118g"
	o.Channel23.Name = "119a"
	o.Channel24.Name = "119b"
	o.Channel25.Name = "119c"
	o.Channel26.Name = "119d"
	o.Channel27.Name = "119e"
	o.Channel28.Name = "119f"
	o.Channel29.Name = "119g"

	o.Channel1.Mz = 114.12773
	o.Channel2.Mz = 115.12476
	o.Channel3.Mz = 115.13108
	o.Channel4.Mz = 115.134
	o.Channel5.Mz = 116.12812
	o.Channel6.Mz = 116.13104
	o.Channel7.Mz = 116.13444
	o.Channel8.Mz = 116.13736
	o.Channel9.Mz = 116.14028
	o.Channel10.Mz = 117.13147
	o.Channel11.Mz = 117.13439
	o.Channel12.Mz = 117.13779
	o.Channel13.Mz = 117.14071
	o.Channel14.Mz = 117.14363
	o.Channel15.Mz = 117.14656
	o.Channel16.Mz = 118.13483
	o.Channel17.Mz = 118.13775
	o.Channel18.Mz = 118.14067
	o.Channel19.Mz = 118.14407
	o.Channel20.Mz = 118.14699
	o.Channel21.Mz = 118.14991
	o.Channel22.Mz = 118.15283
	o.Channel23.Mz = 119.1411
	o.Channel24.Mz = 119.14402
	o.Channel25.Mz = 119.14695
	o.Channel26.Mz = 119.14987
	o.Channel27.Mz = 119.15327
	o.Channel28.Mz = 119.15619
	o.Channel29.Mz = 119.15911

	return o
}
