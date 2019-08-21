package spc

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/prvst/philosopher/lib/msg"

	"github.com/rogpeppe/go-charset/charset"

	// anon charset
	_ "github.com/rogpeppe/go-charset/data"
)

// DataFormat defines different types of data formats from the SPC
type DataFormat interface {
	Parse()
}

// Parameter tag
type Parameter struct {
	XMLName xml.Name `xml:"parameter"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:"value,attr"`
}

// Annotation tag
type Annotation struct {
	XMLName            xml.Name `xml:"annotation"`
	ProteinDescription []byte   `xml:"protein_description,attr"`
}

// ModificationInfo tag
type ModificationInfo struct {
	XMLName          xml.Name           `xml:"modification_info"`
	ModNTermMass     float64            `xml:"mod_nterm_mass,attr"`
	ModCTermMass     float64            `xml:"mod_cterm_mass,attr"`
	ModifiedPeptide  []byte             `xml:"modified_peptide,attr"`
	ModAminoacidMass []ModAminoacidMass `xml:"mod_aminoacid_mass"`
}

// ModAminoacidMass tag
type ModAminoacidMass struct {
	XMLName  xml.Name `xml:"mod_aminoacid_mass"`
	Position int      `xml:"position,attr"`
	Mass     float64  `xml:"mass,attr"`
}

// Parse is the main function for parsing pepxml data
func (p *PepXML) Parse(f string) {

	xmlFile, e := os.Open(f)
	if e != nil {
		msg.ReadFile(e, "fatal")
	}
	defer xmlFile.Close()
	b, _ := ioutil.ReadAll(xmlFile)

	var mpa MsmsPipelineAnalysis

	reader := bytes.NewReader(b)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReader

	if e = decoder.Decode(&mpa); e != nil {
		msg.DecodeMsgPck(e, "fatal")
	}

	p.MsmsPipelineAnalysis = mpa
	p.Name = filepath.Base(f)

	return
}

// Parse is the main function for parsing pepxml data
func (p *ProtXML) Parse(f string) {

	xmlFile, e := os.Open(f)
	if e != nil {
		msg.ReadFile(e, "fatal")
	}
	defer xmlFile.Close()
	b, _ := ioutil.ReadAll(xmlFile)

	var ps ProteinSummary

	reader := bytes.NewReader(b)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReader

	if e = decoder.Decode(&ps); e != nil {
		msg.DecodeMsgPck(e, "fatal")
	}

	p.ProteinSummary = ps
	p.Name = filepath.Base(f)

	return
}
