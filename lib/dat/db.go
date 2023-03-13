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
	//Class            string
	//Length           int
	IsDecoy bool
	//IsContaminant    bool
}

// ProcessHeader parses FASTA records looking for individial elements
func ProcessHeader(k, v string, class dbtype, tag string, verb bool) Record {

	var r Record

	//r.Class = class
	r.OriginalHeader = k
	idx := strings.Index(k, " ")
	if idx == -1 {
		r.PartHeader = k
	} else {
		r.PartHeader = k[:idx]
	}
	//r.Length = len(v)
	r.Sequence = v

	if strings.HasPrefix(k, tag) {
		r.IsDecoy = true
	}

	if strings.Contains(k, "contam_") {
		//r.IsContaminant = true
	}

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

var getGeneName_ensembl = regexp.MustCompile(`(ENSG\d{1,11}\.?\d?\d?)`)
var getGeneName_cptac_ensembl = regexp.MustCompile(`(ENSG\d{1,11}\.?\d?\d?)`)
var getGeneName_ncbi = regexp.MustCompile(`GN\=(.+)\s[\[|OX\=|GN\=|PE\=|$]`)
var getGeneName_uniprot = regexp.MustCompile(`GN\=(.+)\s[\[|OX\=|GN\=|PE\=|$]`)
var getGeneName_tair = regexp.MustCompile(`\|\sSymbols\:(.+?)\s\|`)

func getGeneName(header string, class dbtype, verb bool) (match string) {

	var r *regexp.Regexp

	switch class {
	case ensembl:
		r = getGeneName_ensembl
	case cptac_ensembl:
		r = getGeneName_cptac_ensembl
	case ncbi:
		r = getGeneName_ncbi
	case uniprot:
		r = getGeneName_uniprot
	case uniref:
		return ""
	case tair:
		r = getGeneName_tair
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

var getOrganism_ncbi = regexp.MustCompile(`\[(.+)\]$`)
var getOrganism_uniprot = regexp.MustCompile(`OS\=(.+?)[OX\=|GN\=|PE\=|$]`)

func getOrganism(header string, class dbtype, verb bool) (match string) {

	var r *regexp.Regexp

	switch class {
	case ensembl:
		return ""
	case cptac_ensembl:
		return ""
	case ncbi:
		r = getOrganism_ncbi
	case uniprot:
		r = getOrganism_uniprot
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

var getProteinName_ensembl = regexp.MustCompile(`description\:(.+)\s?$`)
var getProteinName_cptac_ensembl = regexp.MustCompile(`ENS[P|T|G]\d{1,11}\|ENS[P|T|G]\d{1,11}\|ENS[P|T|G]\d{1,11}\|(.+)$`)
var getProteinName_ncbi = regexp.MustCompile(`\s(.+)\sGN?\[?`)
var getProteinName_uniprot = regexp.MustCompile(`\w+\|.+?\|.+?\s(.+?)\s[?OS|?(|?OX|?GN|?PE|?SV]`)
var getProteinName_uniref = regexp.MustCompile(`(UniRef\w+)`)

func getProteinName(header string, class dbtype, verb bool) (match string) {

	var r *regexp.Regexp
	//var reg []string

	switch class {
	case ensembl:
		r = getProteinName_ensembl
	case cptac_ensembl:
		r = getProteinName_cptac_ensembl
	case ncbi:
		r = getProteinName_ncbi
	case uniprot:
		r = getProteinName_uniprot
	case uniref:
		r = getProteinName_uniref
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

	//reg := r.FindStringSubmatch(header)
	reg := r.FindStringSubmatchIndex(header)

	//if reg == nil || len(reg) <= 1 {
	if reg == nil || len(reg) != 4 {

		if verb {
			m := fmt.Sprintf("[protein name]\n%s", header)
			msg.ParsingFASTAHeader(errors.New(m), "info")
		}

		return ""

	} else {
		//match = reg[1]
		match = header[reg[2]:reg[3]]
	}

	return match
}

var getEntryName_ensembl = regexp.MustCompile(`(ENSP\w+\.?\d{1,})`)
var getEntryName_cptac_ensembl = regexp.MustCompile(`(ENSP\w+\.?\d{1,})`)
var getEntryName_ncbi = regexp.MustCompile(`(\w{2}_\d{1,10}\.?(\d{1,2})?)`)
var getEntryName_uniprot = regexp.MustCompile(`\w+\|.+?\|(.+?)\s`)
var getEntryName_uniref = regexp.MustCompile(`(UniRef\w+)`)
var getEntryName_tair = regexp.MustCompile(`^(AT.+)\s\|\sSymbols`)
var getEntryName_nextprot = regexp.MustCompile(`nxp\|(.+?)\|`)

func getEntryName(header string, class dbtype, verb bool) (match string) {

	var r *regexp.Regexp

	switch class {
	case ensembl:
		r = getEntryName_ensembl
	case cptac_ensembl:
		r = getEntryName_cptac_ensembl
	case ncbi:
		r = getEntryName_ncbi
	case uniprot:
		r = getEntryName_uniprot
	case uniref:
		r = getEntryName_uniref
	case tair:
		r = getEntryName_tair
	case nextprot:
		r = getEntryName_nextprot
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

var getID_ensembl = regexp.MustCompile(`(ENSP\w+\.?\d{1,})`)
var getID_cptac_ensembl = regexp.MustCompile(`(ENSP\w+\.?\d{1,})`)
var getID_ncbi = regexp.MustCompile(`(\w{2}_\d{1,10}\.?(\d{1,2})?)`)
var getID_uniprot = regexp.MustCompile(`[sp|tr]\|(.+?)\|`)
var getID_uniref = regexp.MustCompile(`(UniRef\w+)`)
var getID_tair = regexp.MustCompile(`^(AT.+)\s\|\sSymbols`)
var getID_nextprot = regexp.MustCompile(`nxp\|(.+?)\|`)

func getID(header string, class dbtype, verb bool) (match string) {

	var r *regexp.Regexp

	switch class {
	case ensembl:
		r = getID_ensembl
	case cptac_ensembl:
		r = getID_cptac_ensembl
	case ncbi:
		r = getID_ncbi
	case uniprot:
		r = getID_uniprot
	case uniref:
		r = getID_uniref
	case tair:
		r = getID_tair
	case nextprot:
		r = getID_nextprot
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
