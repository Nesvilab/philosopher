package wmm

import (
	"fmt"
	"os"
	"philosopher/lib/met"
)

// Run executes the Filter processing
func Run(m met.Data) met.Data {

	var text string

	// Write about the Database
	text = writeDatabase(m.Database, text)

	if len(m.MSFragger.JarPath) > 0 {
		text = writeSearch(m.MSFragger, text)
	}

	f, _ := os.Create("test.txt")
	defer f.Close()

	f.WriteString(text)

	return m
}

func writeDatabase(d met.Database, text string) string {

	text = fmt.Sprintf("A ORGANISM-NAME database file was downloaded from UniProt [CITATION] using the proteome ID %s on DATE-STAMP\n", d.ID)

	return text
}

func writeSearch(d met.MSFragger, text string) string {

	var searchText string

	// var precursorUnits string
	// var fragmentUnits string

	// if d.PrecursorMassMode == "1" {
	// 	precursorUnits = "ppm"
	// } else if d.PrecursorMassMode == "0" {
	// 	precursorUnits = "Da"
	// }

	// if d.FragmentMassUnits == 1 {
	// 	fragmentUnits = "ppm"
	// } else if d.FragmentMassUnits == 0 {
	// 	fragmentUnits = "Da"
	// }

	searchText = fmt.Sprintf("Database searching was performed on %s files with MSFragger [CITATION] using a precursor tolerance of", d.RawExtension)

	// if d.PrecursorMassUpper > 75 {
	// 	searchText = fmt.Sprintf("An open database search on %s files was performed with MSFragger [CITATION] using precursor tolerance set from %d to %d Da and fragment tolerance of %d %s", d.RawExtension, d.PrecursorMassLower, d.PrecursorMassUpper, d.FragmentMassTolerance, fragmentUnits)
	// } else {
	// 	searchText = fmt.Sprintf("Database searching was performed on %s files with MSFragger [CITATION] using a precursor tolerance of %s %s, fragment tolerance of %s %s", d.RawExtension, d.PrecursorMassLower, precursorUnits, d.FragmentMassTolerance, fragmentUnits)
	// }

	text = text + searchText

	return text
}
