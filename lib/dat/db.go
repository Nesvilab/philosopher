// Package dat (Database)
package dat

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/Nesvilab/philosopher/lib/msg"
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
	Sequence         string
	IsDecoy          bool
}

// ProcessHeader parses FASTA records looking for individial elements
func ProcessHeader(k, v string, class dbtype, tag string, verb bool) Record {

	var r Record

	r.OriginalHeader = k
	idx := strings.Index(k, " ")
	if idx == -1 {
		r.PartHeader = k
	} else {
		r.PartHeader = k[:idx]
	}

	r.Sequence = v

	if strings.HasPrefix(k, tag) {
		r.IsDecoy = true
	}

	// if strings.Contains(k, "contam_") {
	// 	//r.IsContaminant = true
	// }

	r.ID = getID(k, class, verb)
	r.EntryName = getEntryName(k, class, verb)
	r.ProteinName = getProteinName(k, class, verb)
	r.Organism = getOrganism(k, class, verb)
	r.GeneNames = getGeneName(k, class, verb)
	r.ProteinExistence = getProteinExistence(k, class, verb)

	return r
}

var getProteinExistence_uniprot = regexp.MustCompile(`PE\=(.+)\s[\[|OX\=|GN\=|SV\=|$]`)

func getProteinExistence(header string, class dbtype, verb bool) (match string) {

	var r *regexp.Regexp

	switch class {
	case ensembl:
		return ""
	case cptac_ensembl:
		return ""
	case ncbi:
		return ""
	case uniprot:
		r = getProteinExistence_uniprot
	case uniref:
		return ""
	case tair:
		return ""
	case nextprot:
		return ""
	case generic:
		return ""
	default:
		return ""
	}

	reg := r.FindStringSubmatchIndex(header)

	if reg == nil || len(reg) != 4 {
		if verb {
			m := fmt.Sprintf("[protein existence]\n%s", header)
			msg.ParsingFASTAHeader(errors.New(m), "info")
		}

		return ""
	} else {
		match = header[reg[2]:reg[3]]
	}

	switch match {
	case "1":
		return "1:Experimental evidence at protein level"
	case "2":
		return "2:Experimental evidence at transcript level"
	case "3":
		return "3:Protein inferred from homology"
	case "4":
		return "4:Protein predicted"
	case "5":
		return "5:Protein uncertain"
	default:
		return ""
	}
}

var getGeneNameEnsembl = regexp.MustCompile(`(ENSG\d{1,11}\.?\d?\d?)`)
var getGeneNameCptacEnsembl = regexp.MustCompile(`(ENSG\d{1,11}\.?\d?\d?)`)
var getGeneNameNcbi = regexp.MustCompile(`GN\=(.+)\s[\[|OX\=|GN\=|PE\=|$]`)
var getGeneNameUniprot = regexp.MustCompile(`GN\=([[:graph:]]+)`)
var getGeneNameTair = regexp.MustCompile(`\|\sSymbols\:(.+?)\s\|`)

func getGeneName(header string, class dbtype, verb bool) (match string) {

	var r *regexp.Regexp

	switch class {
	case ensembl:
		r = getGeneNameEnsembl
	case cptac_ensembl:
		r = getGeneNameCptacEnsembl
	case ncbi:
		r = getGeneNameNcbi
	case uniprot:
		r = getGeneNameUniprot
	case uniref:
		return ""
	case tair:
		r = getGeneNameTair
	case nextprot:
		s := strings.Split(header, "|")
		s[2] = strings.TrimLeft(s[2], " ")
		s[2] = strings.TrimRight(s[2], " ")
		return s[2]
	case generic:
		return ""
	default:
		return ""
	}

	reg := r.FindStringSubmatchIndex(header)

	if reg == nil || len(reg) != 4 {

		if verb {
			m := fmt.Sprintf("[gene name]\n%s", header)
			msg.ParsingFASTAHeader(errors.New(m), "info")
		}

		return ""

	} else {
		match = header[reg[2]:reg[3]]
	}

	return match
}

var getOrganismNcbi = regexp.MustCompile(`\[(.+)\]$`)
var getOrganismUniprot = regexp.MustCompile(`OS\=(.+?)\s?[OX\=|GN\=|PE\=|$?]`)

func getOrganism(header string, class dbtype, verb bool) (match string) {

	var r *regexp.Regexp

	switch class {
	case ensembl:
		return ""
	case cptac_ensembl:
		return ""
	case ncbi:
		r = getOrganismNcbi
	case uniprot:
		r = getOrganismUniprot
	case uniref:
		return ""
	case tair:
		return ""
	case nextprot:
		return " Homo sapiens"
	case generic:
		return ""
	default:
		return ""
	}

	reg := r.FindStringSubmatchIndex(header)

	if reg == nil || len(reg) != 4 {

		if verb {
			m := fmt.Sprintf("[organism name]\n%s", header)
			msg.ParsingFASTAHeader(errors.New(m), "info")
		}

		return ""

	} else {
		match = header[reg[2]:reg[3]]
	}

	return match
}

var getProteinNameEnsembl = regexp.MustCompile(`description\:(.+)\s?$`)
var getProteinNameCptacEnsembl = regexp.MustCompile(`ENS[P|T|G]\d{1,11}\|ENS[P|T|G]\d{1,11}\|ENS[P|T|G]\d{1,11}\|(.+)$`)
var getProteinNameNcbi = regexp.MustCompile(`\s(.+)\sGN?\[?`)
var getProteinNameUniprot = regexp.MustCompile(`[[:alnum:]]+\_[[:alnum:]]+\s(.+?)\s[[:upper:]][[:upper:]]\=.+`)
var getProteinNameUniref = regexp.MustCompile(`(UniRef\w+)`)

func getProteinName(header string, class dbtype, verb bool) (match string) {

	var r *regexp.Regexp

	switch class {
	case ensembl:
		r = getProteinNameEnsembl
	case cptac_ensembl:
		r = getProteinNameCptacEnsembl
	case ncbi:
		r = getProteinNameNcbi
	case uniprot:
		r = getProteinNameUniprot
	case uniref:
		r = getProteinNameUniref
	case tair:
		s := strings.Split(header, "|")
		s[2] = strings.TrimLeft(s[2], " ")
		s[2] = strings.TrimRight(s[2], " ")
		return s[2]
	case nextprot:
		s := strings.Split(header, "|")
		s[3] = strings.TrimLeft(s[3], " ")
		s[3] = strings.TrimRight(s[3], " ")
		return s[3]
	case generic:
		return header
	default:
		return ""
	}

	reg := r.FindStringSubmatchIndex(header)

	if reg == nil || len(reg) != 4 {

		if verb {
			m := fmt.Sprintf("[protein name]\n%s", header)
			msg.ParsingFASTAHeader(errors.New(m), "info")
		}

		return ""

	} else {
		match = header[reg[2]:reg[3]]
	}

	return match
}

var getEntryNameEnsembl = regexp.MustCompile(`(ENSP\w+\.?\d{1,})`)
var getEntryNameCptacEnsembl = regexp.MustCompile(`(ENSP\w+\.?\d{1,})`)
var getEntryNameNcbi = regexp.MustCompile(`(\w{2}_\d{1,10}\.?(\d{1,2})?)`)
var getEntryNameUniprot = regexp.MustCompile(`\w+\|.+?\|(.+?)\s`)
var getEntryNameUniref = regexp.MustCompile(`(UniRef\w+)`)
var getEntryNameTair = regexp.MustCompile(`^(AT.+)\s\|\sSymbols`)
var getEntryNameNextprot = regexp.MustCompile(`nxp\|(.+?)\|`)

func getEntryName(header string, class dbtype, verb bool) (match string) {

	var r *regexp.Regexp

	switch class {
	case ensembl:
		r = getEntryNameEnsembl
	case cptac_ensembl:
		r = getEntryNameCptacEnsembl
	case ncbi:
		r = getEntryNameNcbi
	case uniprot:
		r = getEntryNameUniprot
	case uniref:
		r = getEntryNameUniref
	case tair:
		r = getEntryNameTair
	case nextprot:
		r = getEntryNameNextprot
	case generic:
		return header
	default:
		return ""
	}

	reg := r.FindStringSubmatchIndex(header)

	if reg == nil || len(reg) != 4 {

		if verb {
			m := fmt.Sprintf("[entry name]\n%s", header)
			msg.ParsingFASTAHeader(errors.New(m), "info")
		}

		return ""

	} else {
		match = header[reg[2]:reg[3]]
	}

	return match
}

var getIDEnsembl = regexp.MustCompile(`(ENSP\w+\.?\d{1,})`)
var getIDCptacEnsembl = regexp.MustCompile(`(ENSP\w+\.?\d{1,})`)
var getIDNcbi = regexp.MustCompile(`(\w{2}_\d{1,10}\.?(\d{1,2})?)`)
var getIDUniprot = regexp.MustCompile(`[sp|tr]\|(.+?)\|`)
var getIDUniref = regexp.MustCompile(`(UniRef\w+)`)
var getIDTair = regexp.MustCompile(`^(AT.+)\s\|\sSymbols`)
var getIDNextprot = regexp.MustCompile(`nxp\|(.+?)\|`)

func getID(header string, class dbtype, verb bool) (match string) {

	var r *regexp.Regexp

	switch class {
	case ensembl:
		r = getIDEnsembl
	case cptac_ensembl:
		r = getIDCptacEnsembl
	case ncbi:
		r = getIDNcbi
	case uniprot:
		r = getIDUniprot
	case uniref:
		r = getIDUniref
	case tair:
		r = getIDTair
	case nextprot:
		r = getIDNextprot
	case generic:
		return header
	default:
		return ""
	}

	reg := r.FindStringSubmatchIndex(header)

	if reg == nil || len(reg) != 4 {

		if verb {
			m := fmt.Sprintf("[protein ID]\n%s", header)
			msg.ParsingFASTAHeader(errors.New(m), "warning")
		}

		return ""

	} else {
		match = header[reg[2]:reg[3]]
	}

	return match
}

type dbtype uint8

const (
	uniprot dbtype = iota
	ncbi
	cptac_ensembl
	ensembl
	uniref
	tair
	nextprot
	generic
)

// Classify determines what kind of database originated the given sequence
func Classify(s, decoyTag string) dbtype {

	// remove the decoy and contamintant tags so we can see better the seq header
	seq := s
	if strings.HasPrefix(seq, decoyTag) {
		seq = seq[len(decoyTag):]
	}
	if strings.HasPrefix(seq, "contam_") {
		seq = seq[len("contam_"):]
	}
	if strings.HasPrefix(seq, decoyTag) {
		seq = seq[len(decoyTag):]
	}

	if strings.HasPrefix(seq, "sp|") || strings.HasPrefix(seq, "tr|") || strings.HasPrefix(seq, "db|") {
		return uniprot
	} else if strings.HasPrefix(seq, "AP_") || strings.HasPrefix(seq, "NP_") || strings.HasPrefix(seq, "YP_") || strings.HasPrefix(seq, "XP_") || strings.HasPrefix(seq, "ZP") || strings.HasPrefix(seq, "WP_") {
		return ncbi
	} else if strings.Contains(seq, "ENSP") && strings.Contains(seq, "|ENST") && strings.Contains(seq, "|ENSG") {
		return cptac_ensembl
	} else if strings.HasPrefix(seq, "ENSP") {
		return ensembl
	} else if strings.HasPrefix(seq, "UniRef") {
		return uniref
	} else if strings.HasPrefix(seq, "AT") {
		return tair
	} else if strings.HasPrefix(seq, "nxp") {
		return nextprot
	}

	return generic
}
