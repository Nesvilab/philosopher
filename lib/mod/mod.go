package mod

// Modifications is a collection of modifications
type Modifications struct {
	Index map[string]Modification
}

// Modification is the basic attribute for each modification
type Modification struct {
	Index             string
	ID                string
	Name              string
	Definition        string
	Variable          string
	Position          string
	Type              string
	AminoAcid         string
	IsProteinTerminus string
	Terminus          string
	MonoIsotopicMass  float64
	AverageMass       float64
	MassDiff          float64
	IsobaricMods      map[string]float64
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
