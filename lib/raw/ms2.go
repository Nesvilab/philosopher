package raw

// MS2 top struct
type MS2 struct {
	Ms2Scan []Ms2Scan
}

// Ms2Scan tag
type Ms2Scan struct {
	Index         string
	Scan          string
	SpectrumName  string
	ScanStartTime float64
	Precursor     Precursor
	Spectrum      Spectrum
}
