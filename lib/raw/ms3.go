package raw

// MS3 top struct
type MS3 struct {
	Ms3Scan []Ms3Scan
}

// Ms3Scan tag
type Ms3Scan struct {
	Index         string
	Scan          string
	SpectrumName  string
	ScanStartTime float64
	Precursor     Precursor
	Spectrum      Spectrum
}
