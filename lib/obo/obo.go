package obo

// // OBO is the top level struct
// type OBO struct {
// 	FormatVersion string
// 	Date          string
// 	Terms         []Term
// }

// // Term represents the ontology terms
// type Term struct {
// 	ID           string
// 	Name         string
// 	Definition   string
// 	isA          string
// 	Relationship []string
// 	Synonym      []string
// 	Comment      string
// 	isObsolete   bool
// 	XRefs        map[string]string
// }

// // Parse is the function that will read an OBO file and return the filled structs
// func (o *OBO) Parse(f string) *err.Error {

// 	oboFile, e := os.Open(f)
// 	if e != nil {
// 		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: filepath.Base(f)}
// 	}
// 	defer oboFile.Close()

// 	scanner := bufio.NewScanner(oboFile)
// 	for scanner.Scan() {
// 		fmt.Println(scanner.Text())
// 	}

// 	if err := scanner.Err(); err != nil {
// 		log.Fatal(err)
// 	}

// 	return nil
// }
