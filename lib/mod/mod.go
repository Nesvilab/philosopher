package mod

// Modifications is a collections of modification
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
	MonoIsotopicMass  float64
	AverageMass       float64
	MassDiff          float64
	AminoAcid         string
	IsProteinTerminus string
	Terminus          string
	IsobaricMods      map[string]uint8
}
