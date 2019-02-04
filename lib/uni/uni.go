package uni

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/mod"
	"github.com/prvst/philosopher/lib/sys"
)

// MOD main structure
type MOD struct {
	met.Data
	XMLFile       string
	Modifications ModList
}

// MODElement struct
type MODElement struct {
	RecordID    int
	Title       string
	Description string
	Posted      string
	Updated     string
	MonoMass    float64
	AvgMass     float64
	Composition string
	Specificity []Specificity
	Xref        []Xref
	Elements    []Element
}

// Specificity struct
type Specificity struct {
	Site           string
	Position       string
	Classification string
	Note           string
}

// Xref struct
type Xref struct {
	Text   string
	Source string
	URL    string
}

// Element struct
type Element struct {
	Symbol string
	Number float64
}

// ModList is a list of UniMOD modifications
type ModList []MODElement

// Len function for Sort
func (p ModList) Len() int {
	return len(p)
}

// Less function for Sort
func (p ModList) Less(i, j int) bool {
	return p[i].RecordID > p[j].RecordID
}

// Swap function for Sort
func (p ModList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// New UniMOD constructor
func New() MOD {

	var o MOD
	var m met.Data
	m.Restore(sys.Meta())

	o.UUID = m.UUID
	o.Distro = m.Distro
	o.Home = m.Home
	o.MetaFile = m.MetaFile
	o.MetaDir = m.MetaDir
	o.DB = m.DB
	o.Temp = m.Temp
	o.TimeStamp = m.TimeStamp

	return o
}

// ProcessUniMOD deploys, reads and assemble the unimod data into structs
func (u *MOD) ProcessUniMOD() error {

	// deploys unimod database
	f, err := u.Deploy()
	if err != nil {
		return err
	}

	// process xml file and load structs
	err = u.Read(f)
	if err != nil {
		return err
	}

	u.Serialize()

	return nil
}

// Deploy unimod xml file to session folder
func (u *MOD) Deploy() (string, error) {

	u.XMLFile = fmt.Sprintf("%s%sunimod.xml", u.Temp, string(filepath.Separator))

	param, err := Asset("unimod.xml")
	err = ioutil.WriteFile(u.XMLFile, param, 0644)

	if err != nil {
		msg := fmt.Sprintf("Could not deploy UniMOD database %s", err)
		return u.XMLFile, errors.New(msg)
	}

	return u.XMLFile, nil
}

// Read is the main function for parsing UniMOD data
func (u *MOD) Read(f string) error {

	var xml mod.XML
	err := xml.Parse(f)
	if err != nil {
		return err
	}

	var list ModList

	for _, i := range xml.UniMOD.Modifications.Mods {

		var u MODElement

		u.Title = i.Title
		u.Description = i.FullName
		u.Posted = i.Posted
		u.Updated = i.Updated
		u.MonoMass = i.Delta.MonoMass
		u.AvgMass = i.Delta.AvgMass
		u.Composition = i.Delta.Composition

		for _, j := range i.Specificities {
			var spec Specificity
			spec.Site = j.Site
			spec.Position = j.Position
			spec.Classification = j.Classification
			spec.Note = j.Notes.Value
			u.Specificity = append(u.Specificity, spec)
		}

		for _, j := range i.Xrefs {
			var x Xref
			x.Text = j.Text.Value
			x.Source = j.Source.Value
			x.URL = j.URL.Value
			u.Xref = append(u.Xref, x)
		}

		for _, j := range i.Delta.Elements {
			var e Element
			e.Symbol = j.Symbol
			e.Number = j.Number
			u.Elements = append(u.Elements, e)
		}

		list = append(list, u)
	}

	u.Modifications = list
	u.XMLFile = filepath.Base(f)

	return nil
}