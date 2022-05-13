// Package dat (Database)
package dat

import (
	"errors"
	"regexp"
	"strings"

	"philosopher/lib/msg"
)

// Record is the root of all database parsers
type Record struct {
	ID               string
	OriginalHeader   string
	PartHeader       string
	EntryName        string
	ProteinName      string
	Organism         string
	GeneNames        string
	ProteinExistence string
	SequenceVersion  string
	Description      string
	Sequence         string
	Length           int
	IsDecoy          bool
	IsContaminant    bool
}

// ProcessENSEMBL parses ENSEMBL like FASTA records
func ProcessENSEMBL(k, v, decoyTag string) Record {

	var e Record

	// this version accepts ENSEMBL headers from other souces like transcript and gene
	idReg1 := regexp.MustCompile(`(ENS\w+\.?\d{1,})`)
	idReg2 := regexp.MustCompile(`(CONTAM\w+_?:?\w+)`)
	desReg := regexp.MustCompile(`ENS\w+(.*)`)
	geneReg := regexp.MustCompile(`ENS\w+\.?\d{0,2}?\|\w+\.?\d{0,2}?\|\w+\.?\d{0,2}?\|\w+\.?\d{0,2}?\|.+?\|.+?\|(.+?)\|`)
	ensgReg := regexp.MustCompile(`(ENSG\w+\.?\d{1,})`)

	e.OriginalHeader = k

	part := strings.Split(k, " ")
	e.PartHeader = part[0]

	// ID and version
	idm := idReg1.FindStringSubmatch(k)
	if idm == nil {
		idm = idReg2.FindStringSubmatch(k)
		e.ID = idm[1]
	} else {
		e.ID = idm[1]
	}

	// Description
	desc := desReg.FindStringSubmatch(k)
	if desc == nil {
		e.Description = ""
	} else {
		e.Description = desc[1]
	}

	gene := geneReg.FindStringSubmatch(k)
	if gene == nil {
		e.GeneNames = ""
	} else {
		e.GeneNames = gene[1]
	}

	ensGene := ensgReg.FindStringSubmatch(k)
	if ensGene == nil {
		e.GeneNames = ""
	} else {
		e.GeneNames = ensGene[1]
	}

	e.EntryName = ""
	e.ProteinName = ""
	e.Organism = ""
	e.ProteinExistence = ""
	e.SequenceVersion = ""
	e.Description = ""

	// Sequence
	e.Sequence = v

	// Length
	e.Length = len(v)

	if strings.HasPrefix(k, decoyTag) {
		e.IsDecoy = true
	} else {
		e.IsDecoy = false
	}

	return e
}

// ProcessNCBI parses UniProt like FASTA records
func ProcessNCBI(k, v, decoyTag string) Record {

	var e Record

	idReg1 := regexp.MustCompile(`(\w{2}_\d{1,10}\.(\d{1,2}))`)
	idReg2 := regexp.MustCompile(`(\w{2}_\d{1,10}\.?(\d{1,2})?)`)
	pnReg := regexp.MustCompile(`\w{2}_\d{1,10}\.?\d{1,2}?\s(.+)GN`)
	genReg1 := regexp.MustCompile(`\sGN=(\w+)\s`)
	genReg2 := regexp.MustCompile(`\sGN=(.+)\s\[`)
	orReg := regexp.MustCompile(`\[(.+)\]`)
	desReg1 := regexp.MustCompile(`\w{2}_\d{1,10}\.?\d{1,2}?\s(.+)\sGN?\[?`)
	desReg2 := regexp.MustCompile(`\w{2}_\d{1,10}\.?\d{1,2}?(.+)\[.?`)

	e.OriginalHeader = k

	part := strings.Split(k, " ")
	e.PartHeader = part[0]

	// ID and version
	idm := idReg1.FindStringSubmatch(k)
	if idm == nil {
		idm = idReg2.FindStringSubmatch(k)
		e.ID = idm[1]
		e.SequenceVersion = ""
		e.EntryName = idm[1]
	} else {
		i := strings.Split(idm[1], ".")
		e.ID = i[0]
		e.SequenceVersion = idm[2]
		e.EntryName = idm[1]
	}

	// Protein Existence
	e.ProteinExistence = ""

	// Protein Name
	pnm := pnReg.FindStringSubmatch(k)
	if pnm == nil {
		e.ProteinName = ""
	} else {
		e.ProteinName = pnm[1]
	}

	// Gene Names
	genn1 := genReg1.FindStringSubmatch(k)

	if genn1 == nil {
		genn2 := genReg2.FindStringSubmatch(k)
		if genn2 == nil {
			e.GeneNames = ""
		} else {
			e.GeneNames = genn2[1]
		}
	} else {
		e.GeneNames = genn1[1]
	}

	// Description
	desc := desReg1.FindStringSubmatch(k)
	if desc == nil {
		desc = desReg2.FindStringSubmatch(k)
		if desc == nil {
			e.Description = ""
		} else {
			e.Description = desc[1]
		}
	} else {
		e.Description = desc[1]
	}

	// Organism Name
	orgn := orReg.FindStringSubmatch(k)
	if orgn == nil {
		e.Organism = ""
	} else {
		e.Organism = orgn[1]
	}

	// Sequence
	e.Sequence = v

	// Length
	e.Length = len(v)

	if strings.HasPrefix(k, decoyTag) {
		e.IsDecoy = true
	} else {
		e.IsDecoy = false
	}

	return e
}

// ProcessUniProtKB parses UniProt like FASTA records
func ProcessUniProtKB(k, v, decoyTag string) Record {

	var e Record

	idReg := regexp.MustCompile(`\w+\|(.+?)\|`)
	enReg := regexp.MustCompile(`\w+\|.+?\|(.+?)\s`)
	smEnR := regexp.MustCompile(`\w+\|.+?\|(.+)`)
	pnReg := regexp.MustCompile(`\w+\|.+?\|.+?\s(.+?)OS`)
	orReg1 := regexp.MustCompile(`OS=(.+?)(\sGN.+|\sPE.+|\sSV.+)`)
	orReg2 := regexp.MustCompile(`OS=(.+)(\sGN.+|\sPE.+|\sSV.+)?`)

	part := strings.Split(k, " ")
	e.PartHeader = part[0]

	// ID
	idm := idReg.FindStringSubmatch(k)
	e.ID = idm[1]

	// Entry Name
	enm := enReg.FindStringSubmatch(k)
	if enm == nil {
		smEnm := smEnR.FindStringSubmatch(k)

		if smEnm == nil {
			e.EntryName = ""
		} else {
			e.EntryName = smEnm[1]
		}

	} else {
		e.EntryName = enm[1]
	}

	// Protein Name
	pnm := pnReg.FindStringSubmatch(k)
	if pnm == nil {
		e.ProteinName = ""
		e.ProteinName = ""
		e.Description = ""
	} else {
		e.ProteinName = pnm[1]
		e.ProteinName = pnm[1]
		e.Description = pnm[1]
	}

	// Organism
	var orn []string
	if strings.Contains(k, "GN=") || strings.Contains(k, "PE=") || strings.Contains(k, "SV=") {
		orn = orReg1.FindStringSubmatch(k)
	} else {
		orn = orReg2.FindStringSubmatch(k)
	}

	if pnm == nil {
		e.Organism = ""
	} else {
		e.Organism = orn[1]
	}

	// Gene Names
	var gnm []string
	if strings.Contains(k, "GN=") && (strings.Contains(k, "PE=") || strings.Contains(k, "SV=")) {

		if len(orn) < 2 {
			msg.ParsingFASTA(errors.New(""), "fatal")
		}

		gnReg := regexp.MustCompile(`GN=(.+?)(\s.+)`)
		gnm = gnReg.FindStringSubmatch(orn[2])
	} else if strings.Contains(k, "GN=") {
		gnReg := regexp.MustCompile(`GN=(.+)$?\s?`)
		gnm = gnReg.FindStringSubmatch(orn[2])
	}

	if gnm != nil {
		e.GeneNames = gnm[1]
	} else {
		e.GeneNames = ""
	}

	var pem []string
	if strings.Contains(k, "PE=") && strings.Contains(k, "SV=") {
		gnReg := regexp.MustCompile(`PE=(.+?)(\s.+)`)
		pem = gnReg.FindStringSubmatch(orn[2])
	} else if strings.Contains(k, "PE=") {
		gnReg := regexp.MustCompile(`PE=(.+)$?\s?`)
		pem = gnReg.FindStringSubmatch(orn[2])
	}

	if pem != nil {
		switch pem[1] {
		case "1":
			e.ProteinExistence = "1:Experimental evidence at protein level"
		case "2":
			e.ProteinExistence = "2:Experimental evidence at transcript level"
		case "3":
			e.ProteinExistence = "3:Protein inferred from homology"
		case "4":
			e.ProteinExistence = "4:Protein predicted"
		case "5":
			e.ProteinExistence = "5:Protein uncertain"
		}
	} else {
		if len(pem) > 0 {
			e.ProteinExistence = pem[0]
		} else {
			e.ProteinExistence = ""
		}
	}

	var svm []string
	if strings.Contains(k, "PE=") {
		svReg := regexp.MustCompile(`SV=(.+)$?\s?`)
		svm = svReg.FindStringSubmatch(orn[2])
	}

	if svm != nil {
		e.SequenceVersion = svm[1]
	} else {
		e.SequenceVersion = ""
	}

	e.Sequence = v
	e.Length = len(v)

	if strings.HasPrefix(k, decoyTag) {
		e.IsDecoy = true
	} else {
		e.IsDecoy = false
	}

	if strings.Contains(k, "contam_") {
		e.IsContaminant = true
	} else {
		e.IsContaminant = false
	}

	e.OriginalHeader = k

	return e
}

// ProcessUniRef parses UniProt like FASTA records
func ProcessUniRef(k, v, decoyTag string) Record {

	var e Record

	pnReg := regexp.MustCompile(`UniRef.+?\s(.+?)n\=`)
	orReg := regexp.MustCompile(`UniRef.+?\s.+?n\=\d{1,}\sTax=(.+)TaxID`)

	// PartHeader
	part := strings.Split(k, " ")
	e.PartHeader = part[0]

	// ID
	e.ID = part[0]

	// Entry Name
	e.EntryName = part[0]

	// Protein Name
	pnm := pnReg.FindStringSubmatch(k)
	if pnm == nil {
		e.ProteinName = ""
	} else {
		e.ProteinName = pnm[1]
	}

	// Organism
	orn := orReg.FindStringSubmatch(k)
	if orn == nil {
		e.Organism = ""
	} else {
		e.Organism = orn[1]
	}

	// Gene Names
	e.GeneNames = ""

	e.Sequence = v
	e.Length = len(v)

	if strings.HasPrefix(k, decoyTag) {
		e.IsDecoy = true
	} else {
		e.IsDecoy = false
	}

	e.OriginalHeader = k

	return e
}

// ProcessGeneric parses generic and uknown database headers
func ProcessGeneric(k, v, decoyTag string) Record {

	var e Record

	idReg := regexp.MustCompile(`(.*)`)

	// ID
	idm := idReg.FindStringSubmatch(k)
	e.ID = idm[1]

	part := strings.Split(k, " ")

	if len(part) > 1 {
		e.Description = part[1]
	} else {
		e.Description = k
	}

	e.EntryName = k
	e.GeneNames = ""
	e.Organism = ""
	e.SequenceVersion = ""

	e.Sequence = v
	e.Length = len(v)
	e.OriginalHeader = k

	e.PartHeader = part[0]

	if strings.HasPrefix(k, decoyTag) {
		e.IsDecoy = true
	} else {
		e.IsDecoy = false
	}

	return e
}

// ProcessTair parses Arabidopsis DB database headers
func ProcessTair(k, v, decoyTag string) Record {

	var e Record

	parts := strings.Split(k, "|")

	// ID
	e.ID = parts[0]
	e.ID = strings.TrimLeft(e.ID, " ")
	e.ID = strings.TrimRight(e.ID, " ")

	// Gene
	if strings.Contains(parts[1], "no symbol available") {
		e.GeneNames = ""
	} else {
		geneReg := regexp.MustCompile(`Symbols\:\s(.+)\s`)
		gm := geneReg.FindStringSubmatch(parts[1])
		e.GeneNames = gm[1]
	}

	// Descripion
	e.Description = parts[2]
	e.Description = strings.TrimLeft(e.Description, " ")
	e.Description = strings.TrimRight(e.Description, " ")

	part := strings.Split(k, " ")

	e.EntryName = k
	e.Organism = ""
	e.SequenceVersion = ""

	e.Sequence = v
	e.Length = len(v)
	e.OriginalHeader = k

	e.PartHeader = part[0]

	if strings.HasPrefix(k, decoyTag) {
		e.IsDecoy = true
	} else {
		e.IsDecoy = false
	}

	return e
}

// Classify determines what kind of database originated the given sequence
func Classify(s, decoyTag string) string {

	// remove the decoy and contamintant tags so we can see better the seq header
	seq := strings.Replace(s, decoyTag, "", -1)
	seq = strings.Replace(seq, "contam_", "", -1)

	if strings.HasPrefix(seq, "sp|") || strings.HasPrefix(seq, "tr|") || strings.HasPrefix(seq, "db|") {
		return "uniprot"
	} else if strings.HasPrefix(seq, "AP_") || strings.HasPrefix(seq, "NP_") || strings.HasPrefix(seq, "YP_") || strings.HasPrefix(seq, "XP_") || strings.HasPrefix(seq, "ZP") || strings.HasPrefix(seq, "WP_") {
		return "ncbi"
	} else if strings.HasPrefix(seq, "ENSP") {
		return "ensembl"
	} else if strings.HasPrefix(seq, "UniRef") {
		return "uniref"
	} else if strings.HasPrefix(seq, "AT") {
		return "tair"
	}

	return "generic"
}
