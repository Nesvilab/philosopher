package mod

// Modifications is a collections of modifications
type Modifications struct {
	InternalIndex map[float64]uint8
	Internal      []InternalModification
	TerminalIndex map[float64]uint8
	Terminal      []TerminalModification
}

// Mod is the basic attribute for each modification
type Mod struct {
	ID               string
	Name             string
	Definition       string
	Variable         string
	Position         string
	Type             string
	MonoIsotopicMass float64
	AverageMass      float64
	MassDiff         float64
}

// InternalModification is a modification that happens inside the peptide structure
type InternalModification struct {
	Mod
	AminoAcid string
}

// TerminalModification is a list of assigned terminal modifications from the database search
type TerminalModification struct {
	Mod
	ProteinTerminus string
	Terminus        string
}
