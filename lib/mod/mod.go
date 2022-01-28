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
	Index      string
	ID         string
	Name       string
	Definition string
	//MonoIsotopicMass  float64
	//AverageMass       float64
	AminoAcid string
	//IsProteinTerminus string
	//Terminus          string
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

// Serialize saves to disk a msgpack version of the Isobaric data structure
// func (m *Modifications) Serialize() {

// 	b, e := msgpack.Marshal(&m)
// 	if e != nil {
// 		msg.MarshalFile(e, "fatal")
// 	}

// 	e = ioutil.WriteFile(sys.EvModificationsBin(), b, sys.FilePermission())
// 	if e != nil {
// 		msg.SerializeFile(e, "fatal")
// 	}
// }

// Restore reads philosopher results files and restore the data sctructure
// func (m *Modifications) Restore() {

// 	b, e := ioutil.ReadFile(sys.EvModificationsBin())
// 	if e != nil {
// 		msg.MarshalFile(e, "warning")
// 	}

// 	e = msgpack.Unmarshal(b, &m)
// 	if e != nil {
// 		msg.SerializeFile(e, "warning")
// 	}
// }
