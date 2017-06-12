package tmt

import "errors"

// Labels main struct
type Labels struct {
	Spectrum      string
	Index         uint32
	RetentionTime float64
	Channel1      Channel1
	Channel2      Channel2
	Channel3      Channel3
	Channel4      Channel4
	Channel5      Channel5
	Channel6      Channel6
	Channel7      Channel7
	Channel8      Channel8
	Channel9      Channel9
	Channel10     Channel10
}

// LabeledSpectra is a list of spectra lables
type LabeledSpectra map[string]Labels

// Channel1 TMT
type Channel1 struct {
	Mz             float64
	Intensity      float64
	NormIntensity  float64
	RatioIntensity float64
	TopIntensity   float64
	Mean           float64
	StDev          float64
}

// Channel2 TMT
type Channel2 struct {
	Mz             float64
	Intensity      float64
	NormIntensity  float64
	RatioIntensity float64
	TopIntensity   float64
	Mean           float64
	StDev          float64
}

// Channel3 TMT
type Channel3 struct {
	Mz             float64
	Intensity      float64
	NormIntensity  float64
	RatioIntensity float64
	TopIntensity   float64
	Mean           float64
	StDev          float64
}

// Channel4 TMT
type Channel4 struct {
	Mz             float64
	Intensity      float64
	NormIntensity  float64
	RatioIntensity float64
	TopIntensity   float64
	Mean           float64
	StDev          float64
}

// Channel5 TMT
type Channel5 struct {
	Mz             float64
	Intensity      float64
	NormIntensity  float64
	RatioIntensity float64
	TopIntensity   float64
	Mean           float64
	StDev          float64
}

// Channel6 TMT
type Channel6 struct {
	Mz             float64
	Intensity      float64
	NormIntensity  float64
	RatioIntensity float64
	TopIntensity   float64
	Mean           float64
	StDev          float64
}

// Channel7 TMT
type Channel7 struct {
	Mz             float64
	Intensity      float64
	NormIntensity  float64
	RatioIntensity float64
	TopIntensity   float64
	Mean           float64
	StDev          float64
}

// Channel8 TMT
type Channel8 struct {
	Mz             float64
	Intensity      float64
	NormIntensity  float64
	RatioIntensity float64
	TopIntensity   float64
	Mean           float64
	StDev          float64
}

// Channel9 TMT
type Channel9 struct {
	Mz             float64
	Intensity      float64
	NormIntensity  float64
	RatioIntensity float64
	TopIntensity   float64
	Mean           float64
	StDev          float64
}

// Channel10 TMT
type Channel10 struct {
	Mz             float64
	Intensity      float64
	NormIntensity  float64
	RatioIntensity float64
	TopIntensity   float64
	Mean           float64
	StDev          float64
}

// New builds a new Labelled spectra object
func New(plex string) (Labels, error) {

	var o Labels

	if plex == "6" {
		o.Channel1.Mz = 126.1277
		o.Channel2.Mz = 127.1248
		o.Channel3.Mz = 128.1344
		o.Channel4.Mz = 129.1315
		o.Channel5.Mz = 130.1411
		o.Channel6.Mz = 131.1382
	} else if plex == "10" {
		o.Channel1.Mz = 126.127725
		o.Channel2.Mz = 127.124760
		o.Channel3.Mz = 127.131079
		o.Channel4.Mz = 128.128114
		o.Channel5.Mz = 128.134433
		o.Channel6.Mz = 129.131468
		o.Channel7.Mz = 129.137787
		o.Channel8.Mz = 130.134822
		o.Channel9.Mz = 130.141141
		o.Channel10.Mz = 131.138176
	} else {
		return o, errors.New("Unknown multiplex value")
	}

	return o, nil
}
