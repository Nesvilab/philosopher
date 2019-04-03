package psi

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/prvst/philosopher/lib/err"
	"github.com/rogpeppe/go-charset/charset"
	// anon charset
	_ "github.com/rogpeppe/go-charset/data"
)

// DataFormat defines different data type from PSI
type DataFormat interface {
	Parse() error
}

// CvList tag
type CvList struct {
	XMLName xml.Name `xml:"cvList"`
	Count   int      `xml:"count,attr"`
	CV      []CV     `xml:"cv"`
}

// CV tag
type CV struct {
	XMLName  xml.Name `xml:"cv"`
	ID       string   `xml:"id,attr"`
	Version  string   `xml:"version,attr"`
	URI      string   `xml:"URI,attr"`
	FullName string   `xml:"fullName,attr"`
}

// CVParam tag
type CVParam struct {
	XMLName       xml.Name `xml:"cvParam"`
	CVRef         string   `xml:"cvRef,attr"`
	Accession     string   `xml:"accession,attr"`
	Name          string   `xml:"name,attr"`
	Value         string   `xml:"value,attr"`
	UnitCvRef     string   `xml:"unitCvRef,attr"`
	UnitAccession string   `xml:"unitAccession,attr"`
	UnitName      string   `xml:"unitName,attr"`
}

// UserParam tag
type UserParam struct {
	XMLName       xml.Name `xml:"userParam"`
	Name          string   `xml:"name,attr"`
	Type          string   `xml:"type,attr"`
	UnitAccession string   `xml:"unitAccession,attr"`
	UnitCvRef     string   `xml:"unitCvRef,attr"`
	UnitName      string   `xml:"UnitName,attr"`
	Value         string   `xml:"value,attr"`
}

// Parse is the main function for parsing pepxml data
func (p *IndexedMzML) Parse(f string) error {

	xmlFile, e := os.Open(f)
	if e != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: filepath.Base(f)}
	}
	defer xmlFile.Close()
	b, _ := ioutil.ReadAll(xmlFile)

	var mzml IndexedMzML

	reader := bytes.NewReader(b)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReader

	if e = decoder.Decode(&mzml); e != nil {
		return &err.Error{Type: err.CannotParseXML, Class: err.FATA, Argument: e.Error()}
	}

	p.MzML = mzml.MzML
	p.Name = filepath.Base(f)

	return nil
}

// Parse is the main function for parsing pepxml data
func (p *MzIdentML) Parse(f string) error {

	xmlFile, e := os.Open(f)
	if e != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: filepath.Base(f)}
	}
	defer xmlFile.Close()
	b, _ := ioutil.ReadAll(xmlFile)

	var mzid MzIdentML

	reader := bytes.NewReader(b)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReader

	if e = decoder.Decode(&mzid); e != nil {
		return &err.Error{Type: err.CannotParseXML, Class: err.FATA, Argument: e.Error()}
	}

	p = &mzid

	return nil
}

// Parse is the main function for parsing pepxml data
func (p *MzIdentML) Write(f string) error {
	//
	// xmlFile, e := os.Open(f)
	// if e != nil {
	// 	return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: filepath.Base(f)}
	// }
	// defer xmlFile.Close()
	// b, _ := ioutil.ReadAll(xmlFile)
	//
	// var mzid MzIdentML
	//
	// reader := bytes.NewReader(b)
	// decoder := xml.NewDecoder(reader)
	// decoder.CharsetReader = charset.NewReader
	//
	// if e = decoder.Decode(&mzid); e != nil {
	// 	return &err.Error{Type: err.CannotParseXML, Class: err.FATA, Argument: e.Error()}
	// }
	//
	// t := time.Now()
	//
	// p = &mzid
	// p.ID = "Philosopher"
	// p.Name = filepath.Base(f)
	// p.Version = "1.2.0"
	// p.CreationDate = t.Format(time.ANSIC)

	return nil
}
