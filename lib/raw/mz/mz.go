package mz

// Raw struct
type Raw struct {
	FileName string
	Spectra  Spectra
}

// Spectra is a list of Spetrum
type Spectra []Spectrum

// Spectrum struct
type Spectrum struct {
	Index       []byte // uint32
	Scan        []byte // uint32
	Level       []byte // uint16
	StartTime   []byte // float64
	Precursor   Precursor
	Peaks       Peaks
	Intensities Intensities
}

// Peaks struct
type Peaks struct {
	Stream      []byte
	Precision   []byte
	Compression []byte
}

// Intensities struct
type Intensities struct {
	Stream      []byte
	Precision   []byte
	Compression []byte
}

// Precursor struct
type Precursor struct {
	ParentIndex                []byte // string
	ParentScan                 []byte // string
	ChargeState                []byte // uint16
	SelectedIon                []byte // float64
	TargetIon                  []byte // float64
	IsolationWindowLowerOffset []byte // float64
	IsolationWindowUpperOffset []byte // float64
	PeakIntensity              []byte // float64
}
