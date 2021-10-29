package iso

// Labels main struct
type Labels struct {
	Spectrum      string
	Index         string
	Scan          string
	RetentionTime float64
	ChargeState   int
	IsUsed        bool
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
	Channel11     Channel11
	Channel12     Channel12
	Channel13     Channel13
	Channel14     Channel14
	Channel15     Channel15
	Channel16     Channel16
	Channel17     Channel17
	Channel18     Channel18
}

// LabeledSpectra is a list of spectra lables
type LabeledSpectra map[string]Labels

// Channel1 TMT
type Channel1 struct {
	Name       string
	CustomName string
	Mz         float64
	Intensity  float64
}

// Channel2 TMT
type Channel2 struct {
	Name       string
	CustomName string
	Mz         float64
	Intensity  float64
}

// Channel3 TMT
type Channel3 struct {
	Name       string
	CustomName string
	Mz         float64
	Intensity  float64
}

// Channel4 TMT
type Channel4 struct {
	Name       string
	CustomName string
	Mz         float64
	Intensity  float64
}

// Channel5 TMT
type Channel5 struct {
	Name       string
	CustomName string
	Mz         float64
	Intensity  float64
}

// Channel6 TMT
type Channel6 struct {
	Name       string
	CustomName string
	Mz         float64
	Intensity  float64
}

// Channel7 TMT
type Channel7 struct {
	Name       string
	CustomName string
	Mz         float64
	Intensity  float64
}

// Channel8 TMT
type Channel8 struct {
	Name       string
	CustomName string
	Mz         float64
	Intensity  float64
}

// Channel9 TMT
type Channel9 struct {
	Name       string
	CustomName string
	Mz         float64
	Intensity  float64
}

// Channel10 TMT
type Channel10 struct {
	Name       string
	CustomName string
	Mz         float64
	Intensity  float64
}

// Channel11 TMT
type Channel11 struct {
	Name       string
	CustomName string
	Mz         float64
	Intensity  float64
}

// Channel12 TMT
type Channel12 struct {
	Name       string
	CustomName string
	Mz         float64
	Intensity  float64
}

// Channel13 TMT
type Channel13 struct {
	Name       string
	CustomName string
	Mz         float64
	Intensity  float64
}

// Channel14 TMT
type Channel14 struct {
	Name       string
	CustomName string
	Mz         float64
	Intensity  float64
}

// Channel15 TMT
type Channel15 struct {
	Name       string
	CustomName string
	Mz         float64
	Intensity  float64
}

// Channel16 TMT
type Channel16 struct {
	Name       string
	CustomName string
	Mz         float64
	Intensity  float64
}

// Channel17 TMT
type Channel17 struct {
	Name       string
	CustomName string
	Mz         float64
	Intensity  float64
}

// Channel18 TMT
type Channel18 struct {
	Name       string
	CustomName string
	Mz         float64
	Intensity  float64
}
