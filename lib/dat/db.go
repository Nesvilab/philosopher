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
	Class            string
	Length           int
	IsDecoy          bool
	IsContaminant    bool
}

// ProcessHeader parses FASTA records looking for individial elements
func ProcessHeader(k, v, class, tag string, verb bool) Record {

	var r Record

	r.Class = class
	r.OriginalHeader = k
	r.PartHeader = strings.Split(k, " ")[0]
	r.Length = len(v)
	r.Sequence = v

	if strings.HasPrefix(k, tag) {
		r.IsDecoy = true
	}

	if strings.Contains(k, "contam_") {
		r.IsContaminant = true
	}

	r.ID = getID(k, class, verb)
	r.EntryName = getEntryName(k, class, verb)
	r.ProteinName = getProteinName(k, class, verb)
	r.Organism = getOrganism(k, class, verb)
	r.GeneNames = getGeneName(k, class, verb)
	r.ProteinExistence = getProteinExistence(k, class, verb)

	return r
}

func getProteinExistence(header, class string, verb bool) (match string) {

	var r *regexp.Regexp
	var reg []string

	switch class {
	case "ensembl":
		return ""
	case "cptac-ensembl":
		return ""
	case "ncbi":
		return ""
	case "uniprot":
		r = regexp.MustCompile(`PE\=(.+)\s[\[|OX\=|GN\=|SV\=|$]`)
	case "uniref":
		return ""
	case "tair":
		return ""
	case "nextprot":
		return ""
	case "generic":
		return ""
	default:
		return ""
	}

	reg = r.FindStringSubmatch(header)

	if reg == nil || len(reg) <= 1 {
		if verb {
			m := fmt.Sprintf("[protein existence]\n%s", header)
			msg.ParsingFASTAHeader(errors.New(m), "info")
		}

		return ""
	} else {
		match = reg[1]
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

func getGeneName(header, class string, verb bool) (match string) {

	var r *regexp.Regexp
	var reg []string

	switch class {
	case "ensembl":
		r = regexp.MustCompile(`(ENSG\d{1,11}\.?\d?\d?)`)
	case "cptac-ensembl":
		r = regexp.MustCompile(`(ENSG\d{1,11}\.?\d?\d?)`)
	case "ncbi":
		return ""
	case "uniprot":
		r = regexp.MustCompile(`GN\=([[:alnum:]]+)`)
	case "uniref":
		return ""
	case "tair":
		r = regexp.MustCompile(`\|\sSymbols\:(.+?)\s\|`)
	case "nextprot":
		s := strings.Split(header, "|")
		s[2] = strings.TrimLeft(s[2], " ")
		s[2] = strings.TrimRight(s[2], " ")
		return s[2]
	case "generic":
		return ""
	default:
		return ""
	}

	reg = r.FindStringSubmatch(header)

	if reg == nil || len(reg) <= 1 {

		if verb {
			m := fmt.Sprintf("[gene name]\n%s", header)
			msg.ParsingFASTAHeader(errors.New(m), "info")
		}

		return ""

	} else {
		match = reg[1]
	}

	return match
}

func getOrganism(header, class string, verb bool) (match string) {

	var r *regexp.Regexp
	var reg []string

	switch class {
	case "ensembl":
		return ""
	case "cptac-ensembl":
		return ""
	case "ncbi":
		r = regexp.MustCompile(`\[(.+)\]$`)
	case "uniprot":
		r = regexp.MustCompile(`OS\=(.+?)\s?[OX\=|GN\=|PE\=|$?]`)
	case "uniref":
		return ""
	case "tair":
		return ""
	case "nextprot":
		return " Homo sapiens"
	case "generic":
		return ""
	default:
		return ""
	}

	reg = r.FindStringSubmatch(header)

	if reg == nil || len(reg) <= 1 {

		if verb {
			m := fmt.Sprintf("[organism name]\n%s", header)
			msg.ParsingFASTAHeader(errors.New(m), "info")
		}

		return ""

	} else {
		match = reg[1]
	}

	return match
}

func getProteinName(header, class string, verb bool) (match string) {

	var r *regexp.Regexp
	var reg []string

	switch class {
	case "ensembl":
		r = regexp.MustCompile(`description\:(.+)\s?$`)
	case "cptac-ensembl":
		r = regexp.MustCompile(`ENS[P|T|G]\d{1,11}\|ENS[P|T|G]\d{1,11}\|ENS[P|T|G]\d{1,11}\|(.+)$`)
	case "ncbi":
		r = regexp.MustCompile(`\s(.+)\sGN?\[?`)
	case "uniprot":
		r = regexp.MustCompile(`[[:alnum:]]+\_[[:alnum:]]+\s(.+?)\s[[:upper:]][[:upper:]]\=.+`)
	case "uniref":
		r = regexp.MustCompile(`(UniRef\w+)`)
	case "tair":
		s := strings.Split(header, "|")
		s[2] = strings.TrimLeft(s[2], " ")
		s[2] = strings.TrimRight(s[2], " ")
		return s[2]
	case "nextprot":
		s := strings.Split(header, "|")
		s[3] = strings.TrimLeft(s[3], " ")
		s[3] = strings.TrimRight(s[3], " ")
		return s[3]
	case "generic":
		return header
	default:
		return ""
	}

	reg = r.FindStringSubmatch(header)

	if reg == nil || len(reg) <= 1 {

		if verb {
			m := fmt.Sprintf("[protein name]\n%s", header)
			msg.ParsingFASTAHeader(errors.New(m), "info")
		}

		return ""

	} else {
		match = reg[1]
	}

	return match
}

func getEntryName(header, class string, verb bool) (match string) {

	var r *regexp.Regexp
	var reg []string

	switch class {
	case "ensembl":
		r = regexp.MustCompile(`(ENSP\w+\.?\d{1,})`)
	case "cptac-ensembl":
		r = regexp.MustCompile(`(ENSP\w+\.?\d{1,})`)
	case "ncbi":
		r = regexp.MustCompile(`(\w{2}_\d{1,10}\.?(\d{1,2})?)`)
	case "uniprot":
		r = regexp.MustCompile(`\w+\|.+?\|(.+?)\s`)
	case "uniref":
		r = regexp.MustCompile(`(UniRef\w+)`)
	case "tair":
		r = regexp.MustCompile(`^(AT.+)\s\|\sSymbols`)
	case "nextprot":
		r = regexp.MustCompile(`nxp\|(.+?)\|`)
	case "generic":
		return header
	default:
		return ""
	}

	reg = r.FindStringSubmatch(header)

	if reg == nil || len(reg) < 1 {

		if verb {
			m := fmt.Sprintf("[entry name]\n%s", header)
			msg.ParsingFASTAHeader(errors.New(m), "info")
		}

		return ""

	} else {
		match = reg[1]
	}

	return match
}

func getID(header, class string, verb bool) (match string) {

	var r *regexp.Regexp
	var reg []string

	switch class {
	case "ensembl":
		r = regexp.MustCompile(`(ENSP\w+\.?\d{1,})`)
	case "cptac-ensembl":
		r = regexp.MustCompile(`(ENSP\w+\.?\d{1,})`)
	case "ncbi":
		r = regexp.MustCompile(`(\w{2}_\d{1,10}\.?(\d{1,2})?)`)
	case "uniprot":
		r = regexp.MustCompile(`[sp|tr]\|(.+?)\|`)
	case "uniref":
		r = regexp.MustCompile(`(UniRef\w+)`)
	case "tair":
		r = regexp.MustCompile(`^(AT.+)\s\|\sSymbols`)
	case "nextprot":
		r = regexp.MustCompile(`nxp\|(.+?)\|`)
	case "generic":
		return header
	default:
		return ""
	}

	reg = r.FindStringSubmatch(header)

	if reg == nil || len(reg) <= 1 {

		if verb {
			m := fmt.Sprintf("[protein ID]\n%s", header)
			msg.ParsingFASTAHeader(errors.New(m), "warning")
		}

		return ""

	} else {
		match = reg[1]
	}

	return match
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
	} else if strings.Contains(seq, "ENSP") && strings.Contains(seq, "|ENST") && strings.Contains(seq, "|ENSG") {
		return "cptac-ensembl"
	} else if strings.HasPrefix(seq, "ENSP") {
		return "ensembl"
	} else if strings.HasPrefix(seq, "UniRef") {
		return "uniref"
	} else if strings.HasPrefix(seq, "AT") {
		return "tair"
	} else if strings.HasPrefix(seq, "nxp") {
		return "nextprot"
	}

	return "generic"
}
