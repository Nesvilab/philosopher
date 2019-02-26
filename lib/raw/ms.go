package raw

// Spectrum tag
type Spectrum []Peak

// Peak tag
type Peak struct {
	Mz          float64
	Intensity   float64
	IonMobility float64
}

func (a Spectrum) Len() int           { return len(a) }
func (a Spectrum) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Spectrum) Less(i, j int) bool { return a[i].Mz < a[j].Mz }
