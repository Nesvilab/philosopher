package mod

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rogpeppe/go-charset/charset"
	// test
	_ "github.com/rogpeppe/go-charset/data"
)

// XML is the head struct
type XML struct {
	Name   string
	UniMOD UniMOD
}

// UniMOD ...
type UniMOD struct {
	XMLName       xml.Name     `xml:"unimod"`
	Modifications Modification `xml:"modifications"`
}

// Modification struct
type Modification struct {
	XMLName xml.Name `xml:"modifications"`
	Mods    []Mod    `xml:"mod"`
}

// Mod struct
type Mod struct {
	XMLName       xml.Name      `xml:"mod"`
	Title         string        `xml:"title,attr"`
	FullName      string        `xml:"full_name,attr"`
	Posted        string        `xml:"date_time_posted,attr"`
	Updated       string        `xml:"date_time_modified,attr"`
	RecordID      int           `xml:"record_id,attr"`
	Specificities []Specificity `xml:"specificity"`
	Xrefs         []Xref        `xml:"xref"`
	Delta         Delta         `xml:"delta"`
}

// Specificity struct
type Specificity struct {
	XMLName        xml.Name `xml:"specificity"`
	Site           string   `xml:"site,attr"`
	Position       string   `xml:"position,attr"`
	Classification string   `xml:"classification,attr"`
	Notes          Note     `xml:"misc_notes"`
}

// Delta struct
type Delta struct {
	XMLName     xml.Name  `xml:"delta"`
	MonoMass    float64   `xml:"mono_mass,attr"`
	AvgMass     float64   `xml:"avge_mass,attr"`
	Composition string    `xml:"composition,attr"`
	Elements    []Element `xml:"element"`
}

// Element struct
type Element struct {
	XMLName xml.Name `xml:"element"`
	Symbol  string   `xml:"symbol,attr"`
	Number  float64  `xml:"number,attr"`
}

// Note struct
type Note struct {
	XMLName xml.Name `xml:"misc_notes"`
	Value   string   `xml:",chardata"`
}

// Xref struct
type Xref struct {
	XMLName xml.Name `xml:"xref"`
	Text    Text     `xml:"text"`
	Source  Source   `xml:"source"`
	URL     URL      `xml:"url"`
}

// Text struct
type Text struct {
	XMLName xml.Name `xml:"text"`
	Value   string   `xml:",chardata"`
}

// Source struct
type Source struct {
	XMLName xml.Name `xml:"source"`
	Value   string   `xml:",chardata"`
}

// URL struct
type URL struct {
	XMLName xml.Name `xml:"url"`
	Value   string   `xml:",chardata"`
}

// Parse is the main function for parsing pepxml data
func (p *XML) Parse(f string) error {

	xmlFile, err := os.Open(f)
	if err != nil {
		msg := fmt.Sprintf("Error opening file: %s", err)
		return errors.New(msg)
	}
	defer xmlFile.Close()
	b, _ := ioutil.ReadAll(xmlFile)

	var uni UniMOD

	reader := bytes.NewReader(b)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReader

	if err = decoder.Decode(&uni); err != nil {
		msg := fmt.Sprintf("Unable to parse XML: %s", err)
		return errors.New(msg)
	}

	if len(uni.Modifications.Mods) < 1 {
		msg := fmt.Sprintf("Skipping file %s, no peptide validation found", filepath.Base(f))
		return errors.New(msg)
	}

	p.UniMOD = uni
	p.Name = filepath.Base(f)

	return nil
}
