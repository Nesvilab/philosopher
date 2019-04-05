package obo

// Terms is a collection of Term
type Terms []Term

// Term refers to an atomic ontology definition
type Term struct {
	ID               string
	Name             string
	Definition       string
	DateTimePosted   string
	DateTimeModified string
	Comments         string
	Synonyms         []string
	IsA              string
	Mod              Mod
}

// Mod contains UniMod term definition
type Mod struct {
	MonoIsotopicMass float64
	AverageMass      float64
	Composition      string
}
