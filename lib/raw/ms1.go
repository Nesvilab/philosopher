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
	Spectrum      Ms1Spectrum
}

// Ms1Spectrum tag
type Ms1Spectrum []Ms1Peak

// Ms1Peak tag
type Ms1Peak struct {
	Mz        float64
	Intensity float64
}

func (a Ms1Spectrum) Len() int           { return len(a) }
func (a Ms1Spectrum) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Ms1Spectrum) Less(i, j int) bool { return a[i].Mz < a[j].Mz }
