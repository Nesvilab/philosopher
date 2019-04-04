package psi

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/rogpeppe/go-charset/charset"
	"github.com/sirupsen/logrus"
	// anon charset
	_ "github.com/rogpeppe/go-charset/data"
)

// DataFormat defines different data type from PSI
type DataFormat interface {
	Parse() error
}

// SourceFile is a file from which this instance was created
type SourceFile struct {
	XMLName                     xml.Name                    `xml:"SourceFile"`
	ID                          string                      `xml:"id,attr,omitempty"`
	Location                    string                      `xml:"location,attr,omitempty"`
	Name                        string                      `xml:"name,attr,omitempty"`
	ExternalFormatDocumentation ExternalFormatDocumentation `xml:"ExternalFormatDocumentation"`
	FileFormat                  FileFormat                  `xml:"FileFormat"`
	CVParam                     []CVParam                   `xml:"cvParam"`
	UserParam                   []UserParam                 `xml:"userParam"`
}

// CvList is the list of controlled vocabularies used in the file
type CvList struct {
	XMLName xml.Name `xml:"cvList"`
	Count   int      `xml:"count,attr,omitempty"`
	CV      []CV     `xml:"cv"`
}

// CV is a ource controlled vocabulary from which cvParams will be obtained
type CV struct {
	XMLName  xml.Name `xml:"cv"`
	ID       string   `xml:"id,attr,omitempty"`
	Version  string   `xml:"version,attr,omitempty,omitempty"`
	URI      string   `xml:"URI,attr,omitempty"`
	FullName string   `xml:"fullName,attr,omitempty"`
}

// CVParam is single entry from an ontology or a controlled vocabulary
type CVParam struct {
	XMLName       xml.Name `xml:"cvParam"`
	Accession     string   `xml:"accession,attr,omitempty"`
	CVRef         string   `xml:"cvRef,attr,omitempty"`
	Name          string   `xml:"name,attr,omitempty"`
	UnitAccession string   `xml:"unitAccession,attr,omitempty"`
	UnitCvRef     string   `xml:"unitCvRef,attr,omitempty"`
	UnitName      string   `xml:"unitName,attr,omitempty"`
	Value         string   `xml:"value,attr,omitempty"`
}

// UserParam In case more information about the ions annotation has to be
// conveyed, that has no fit in FragmentArray. Note: It is suggested that the
// value attribute takes the form of a list of the same size as FragmentArray
// values. However, there is no formal encoding and it cannot be expeceted that
// other software will process or impart that information properly
type UserParam struct {
	XMLName       xml.Name `xml:"userParam"`
	Name          string   `xml:"name,attr,omitempty"`
	Type          string   `xml:"type,attr,omitempty"`
	UnitAccession string   `xml:"unitAccession,attr,omitempty"`
	UnitCvRef     string   `xml:"unitCvRef,attr,omitempty"`
	UnitName      string   `xml:"UnitName,attr,omitempty"`
	Value         string   `xml:"value,attr,omitempty"`
}

// Parse is the main function for parsing IndexedMzML data
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

// Parse is the main function for parsing MzIdentML data
func (p *MzIdentML) Parse(f string) error {

	xmlFile, e := os.Open(f)
	if e != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: filepath.Base(f)}
	}
	defer xmlFile.Close()
	b, _ := ioutil.ReadAll(xmlFile)

	reader := bytes.NewReader(b)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReader

	if e = decoder.Decode(p); e != nil {
		return &err.Error{Type: err.CannotParseXML, Class: err.FATA, Argument: e.Error()}
	}

	return nil
}

// Parse is the main function for parsing pepxml data
func (p *MzIdentML) Write() error {

	output := fmt.Sprintf("%s%sreport.mzid", sys.MetaDir(), string(filepath.Separator))

	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create MzId report:", err)
	}
	defer file.Close()

	file.WriteString(xml.Header)

	enc := xml.NewEncoder(file)
	enc.Indent("  ", "    ")

	if err := enc.Encode(p); err != nil {
		fmt.Printf("error: %v\n", err)
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return nil
}
