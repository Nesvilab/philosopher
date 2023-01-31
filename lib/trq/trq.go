package trq

import (
	"errors"

	"github.com/Nesvilab/philosopher/lib/iso"
	"github.com/Nesvilab/philosopher/lib/msg"
)

// New builds a new Labelled spectra object
func New(plex string) iso.Labels {

	var o iso.Labels

	if plex == "4" {

		o.Channel1.Name = "114"
		o.Channel2.Name = "115"
		o.Channel3.Name = "116"
		o.Channel4.Name = "117"

		o.Channel1.Mz = 114.1112
		o.Channel2.Mz = 115.1083
		o.Channel3.Mz = 116.1116
		o.Channel4.Mz = 117.1150

	} else if plex == "8" {

		o.Channel1.Name = "113"
		o.Channel2.Name = "114"
		o.Channel3.Name = "115"
		o.Channel4.Name = "116"
		o.Channel5.Name = "117"
		o.Channel6.Name = "118"
		o.Channel7.Name = "119"
		o.Channel8.Name = "121"

		o.Channel1.Mz = 113.1078
		o.Channel2.Mz = 114.1112
		o.Channel3.Mz = 115.1082
		o.Channel4.Mz = 116.1116
		o.Channel5.Mz = 117.1149
		o.Channel6.Mz = 118.1120
		o.Channel7.Mz = 119.1153
		o.Channel8.Mz = 121.1220

	} else {
		msg.Custom(errors.New("unknown multiplex setting, please define the plex number used in your experiment"), "fatal")
	}

	return o
}
