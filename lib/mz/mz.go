package mz

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

// IndexedMzML is the root level tag
type IndexedMzML struct {
	Name    string
	XMLName xml.Name `xml:"indexedmzML"`
	MzML    ML       `xml:"mzML"`
}

// ML is the root level tag
type ML struct {
	XMLName           xml.Name          `xml:"mzML"`
	Accession         string            `xml:"accession,attr"`
	Version           string            `xml:"version,attr"`
	FileDescription   FileDescription   `ml:"fileDescription"`
	RefParamGroupList RefParamGroupList `xml:"referenceableParamGroupList"`
	SoftwareList      SoftwareList      `xml:"softwareList"`
	Run               Run               `xml:"run"`
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
	XMLName xml.Name `xml:"userParam"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:"value,attr"`
	Type    string   `xml:"type,attr"`
}

// FileDescription tag
type FileDescription struct {
	XMLName        xml.Name       `xml:"fileDescription"`
	FileContent    FileContent    `xml:"fileContent"`
	SourceFileList SourceFileList `xml:"sourceFileList"`
}

// FileContent tag
type FileContent struct {
	XMLName xml.Name `xml:"fileContent"`
	CVParam CVParam  `xml:"cvParam"`
}

// SourceFileList tag
type SourceFileList struct {
	XMLName    xml.Name     `xml:"sourceFileList"`
	Count      int          `xml:"count,attr"`
	SourceFile []SourceFile `xml:"sourceFile"`
}

// SourceFile tag
type SourceFile struct {
	XMLName  xml.Name  `xml:"sourceFile"`
	ID       string    `xml:"id,attr"`
	Name     string    `xml:"name,attr"`
	Location string    `xml:"location,attr"`
	CVParam  []CVParam `xml:"cvParam"`
}

// RefParamGroupList tag
type RefParamGroupList struct {
	XMLName       xml.Name        `xml:"referenceableParamGroupList"`
	Count         int             `xml:"count,attr"`
	RefParamGroup []RefParamGroup `xml:"referenceableParamGroup"`
}

// RefParamGroup tag
type RefParamGroup struct {
	XMLName   xml.Name    `xml:"referenceableParamGroup"`
	ID        string      `xml:"id,attr"`
	CVParam   []CVParam   `xml:"cvParam"`
	UserParam []UserParam `xml:"userParam"`
}

// SoftwareList tag
type SoftwareList struct {
	XMLName  xml.Name   `xml:"softwareList"`
	Count    int        `xml:"count,attr"`
	Software []Software `xml:"software"`
}

// Software tag
type Software struct {
	XMLName       xml.Name        `xml:"software"`
	ID            string          `xml:"id,attr"`
	Version       string          `xml:"version,attr"`
	CVParam       []CVParam       `xml:"cvParam"`
	UserParam     []UserParam     `xml:"userParam"`
	RefParamGroup []RefParamGroup `xml:"referenceableParamGroup"`
}

// Run tag
type Run struct {
	XMLName                           xml.Name         `xml:"run"`
	ID                                string           `xml:"id,attr"`
	DefaultInstrumentConfigurationRef string           `xml:"defaultInstrumentConfigurationRef,attr"`
	StartTimeStamp                    string           `xml:"startTimeStamp,attr"`
	SpectrumList                      SpectrumList     `xml:"spectrumList"`
	ChromatogramList                  ChromatogramList `xml:"chromatogramList"`
}

// SpectrumList tag
type SpectrumList struct {
	XMLName                  xml.Name   `xml:"spectrumList"`
	Count                    int        `xml:"count,attr"`
	DefaultDataProcessingRef string     `xml:"defaultDataProcessingRef,attr"`
	Spectrum                 []Spectrum `xml:"spectrum"`
}

// Spectrum tag
type Spectrum struct {
	XMLName             xml.Name            `xml:"spectrum"`
	Index               string              `xml:"index,attr"`
	ID                  string              `xml:"id,attr"`
	DefaultArrayLength  float64             `xml:"defaultArrayLength,attr"`
	DataProcessingRef   string              `xml:"dataProcessingRef,att"`
	CVParam             []CVParam           `xml:"cvParam"`
	ScanList            ScanList            `xml:"scanList"`
	PrecursorList       PrecursorList       `xml:"precursorList"`
	BinaryDataArrayList BinaryDataArrayList `xml:"binaryDataArrayList"`
	Peaks               []float64
	Intensities         []float64
}

// ScanList tag
type ScanList struct {
	XMLName xml.Name  `xml:"scanList"`
	Count   int       `xml:"count,attr"`
	CVParam []CVParam `xml:"cvParam"`
	Scan    []Scan    `xml:"scan"`
}

// PrecursorList tag
type PrecursorList struct {
	XMLName   xml.Name    `xml:"precursorList"`
	Count     int         `xml:"count,attr"`
	Precursor []Precursor `xml:"precursor"`
}

// Precursor tag
type Precursor struct {
	XMLName         xml.Name        `xml:"precursor"`
	SpectrumRef     string          `xml:"spectrumRef,attr"`
	IsolationWindow IsolationWindow `xml:"isolationWindow"`
	SelectedIonList SelectedIonList `xml:"selectedIonList"`
	Activation      Activation      `xml:"activation"`
}

// IsolationWindow tag
type IsolationWindow struct {
	XMLName xml.Name  `xml:"isolationWindow"`
	CVParam []CVParam `xml:"cvParam"`
}

// SelectedIonList tag
type SelectedIonList struct {
	XMLName     xml.Name      `xml:"selectedIonList"`
	Count       int           `xml:"count,attr"`
	SelectedIon []SelectedIon `xml:"selectedIon"`
}

// SelectedIon tag
type SelectedIon struct {
	XMLName xml.Name  `xml:"selectedIon"`
	CVParam []CVParam `xml:"cvParam"`
}

// Scan tag
type Scan struct {
	XMLName              xml.Name       `xml:"scan"`
	InstConfigurationRef string         `xml:"instrumentConfigurationRef,attr"`
	CVParam              []CVParam      `xml:"cvParam"`
	UserParam            []UserParam    `xml:"userParam"`
	ScanWindowList       ScanWindowList `xml:"scanWindowList"`
}

// ScanWindowList tag
type ScanWindowList struct {
	XMLName    xml.Name     `xml:"scanWindowList"`
	Count      int          `xml:"count,attr"`
	ScanWindow []ScanWindow `xml:"scanWindow"`
}

// ScanWindow tag
type ScanWindow struct {
	XMLName xml.Name  `xml:"scanWindow"`
	CVParam []CVParam `xml:"cvParam"`
}

// Activation tag
type Activation struct {
	XMLName xml.Name  `xml:"activation"`
	CVParam []CVParam `xml:"cvParam"`
}

// ChromatogramList tag
type ChromatogramList struct {
	XMLName                  xml.Name       `xml:"chromatogramList"`
	Count                    int            `xml:"count,attr"`
	DefaultDataProcessingRef string         `xml:"defaultDataProcessingRef,attr"`
	Chromatogram             []Chromatogram `xml:"chromatogram"`
}

// Chromatogram tag
type Chromatogram struct {
	XMLName             xml.Name            `xml:"chromatogram"`
	Index               int                 `xml:"index,attr"`
	ID                  string              `xml:"id,attr"`
	DefaultArrayLength  float64             `xml:"defaultArrayLength,attr"`
	CVParam             []CVParam           `xml:"cvParam"`
	BinaryDataArrayList BinaryDataArrayList `xml:"binaryDataArrayList"`
}

// BinaryDataArrayList tag
type BinaryDataArrayList struct {
	XMLName         xml.Name          `xml:"binaryDataArrayList"`
	Count           int               `xml:"count,attr"`
	BinaryDataArray []BinaryDataArray `xml:"binaryDataArray"`
}

// BinaryDataArray tag
type BinaryDataArray struct {
	XMLName       xml.Name  `xml:"binaryDataArray"`
	EncodedLength float64   `xml:"encodedLength,attr"`
	CVParam       []CVParam `xml:"cvParam"`
	Binary        Binary    `xml:"binary"`
}

// Binary tag
type Binary struct {
	XMLName xml.Name `xml:"binary"`
	Value   []byte   `xml:",chardata"`
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
