package psi

import (
	"encoding/xml"
)

// IndexedMzML is the root level tag
type IndexedMzML struct {
	Name    string
	XMLName xml.Name `xml:"indexedmzML"`
	MzML    MzML     `xml:"mzML"`
}

// MzML is the root level tag
type MzML struct {
	XMLName           xml.Name          `xml:"mzML"`
	Accession         string            `xml:"accession,attr"`
	Version           string            `xml:"version,attr"`
	FileDescription   FileDescription   `ml:"fileDescription"`
	RefParamGroupList RefParamGroupList `xml:"referenceableParamGroupList"`
	SoftwareList      SoftwareList      `xml:"softwareList"`
	Run               Run               `xml:"run"`
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
