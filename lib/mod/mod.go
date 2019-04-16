package mod

// Modifications is a collections of modification
type Modifications struct {
	Mods  []Modification
	Index map[string]Modification
}

// Modification is the basic attribute for each modification
type Modification struct {
	Index            string
	ID               string
	Name             string
	Definition       string
	Variable         string
	Position         string
	Type             string
	MonoIsotopicMass float64
	AverageMass      float64
	MassDiff         float64
	Internal         InternalModification
	Terminal         TerminalModification
}

// InternalModification is a modification that happens inside the peptide structure
type InternalModification struct {
	AminoAcid string
}

// TerminalModification is a list of assigned terminal modifications from the database search
type TerminalModification struct {
	IsProteinTerminus string
	Terminus          string
}
