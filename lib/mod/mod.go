package mod

// Modifications is a collection of modifications
type Modifications struct {
	Index map[string]Modification
}

type ModTypeType uint8

const (
	Assigned ModTypeType = iota
	Observed
)

// Modification is the basic attribute for each modification
type Modification struct {
	Index        string
	ID           string
	Name         string
	Definition   string
	AminoAcid    string
	IsobaricMods map[string]float64
	MassDiff     float64
	Position     int
	Type         ModTypeType
	Variable     bool
}

// Modifications is a collection of modifications
type ModificationsSlice struct {
	IndexSlice []Modification
}

func (m Modifications) ToSlice() ModificationsSlice {
	IndexSlice := make([]Modification, 0, len(m.Index))
	for k, v := range m.Index {
		IndexSlice = append(IndexSlice, v)
		if v.Index != k {
			panic(nil)
		}
	}
	return ModificationsSlice{IndexSlice: IndexSlice}
}
func (m ModificationsSlice) ToMap() Modifications {
	Index := make(map[string]Modification, len(m.IndexSlice))
	for _, e := range m.IndexSlice {
		Index[e.Index] = e
	}
	return Modifications{Index: Index}
}
