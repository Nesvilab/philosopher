package raw

// MS1 top struct
type MS1 struct {
	Ms1Scan []Ms1Scan
}

// Ms1Scan tag
type Ms1Scan struct {
	Index         string
	Scan          string
	SpectrumName  string
	ScanStartTime float64
	Spectrum      Spectrum
}
