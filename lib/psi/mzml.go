package psi

import (
	"encoding/xml"
)

// IndexedMzML is the root level tag
type IndexedMzML struct {
	XMLName xml.Name `xml:"indexedmzML"`
	Name    string
	MzML    MzML `xml:"mzML"`
}

// MzML This is the root element for the Proteomics Standards Initiative (PSI) mzML schema, which is intended to
// capture the use of a mass spectrometer, the data generated, and the initial processing of that data
type MzML struct {
	XMLName                     xml.Name                    `xml:"mzML"`
	Accession                   string                      `xml:"accession,attr"`
	ID                          string                      `xml:"id,attr"`
	Version                     string                      `xml:"version,attr"`
	CvList                      CvList                      `xml:"cvList"`
	FileDescription             FileDescription             `ml:"fileDescription"`
	RefParamGroupList           RefParamGroupList           `xml:"referenceableParamGroupList"`
	SampleList                  SampleList                  `xml:"sampleList"`
	SoftwareList                SoftwareList                `xml:"softwareList"`
	ScanSettingsList            ScanSettingsList            `xml:"scanSettingsList"`
	InstrumentConfigurationList InstrumentConfigurationList `xml:"instrumentConfigurationList"`
	DataProcessingList          DataProcessingList          `xml:"dataProcessingList"`
	Run                         Run                         `xml:"run"`
}

// DataProcessingList is a list and descriptions of data processing applied to this data
type DataProcessingList struct {
	XMLName        xml.Name         `xml:"dataProcessingList"`
	Count          int              `xml:"count,attr,omitempty"`
	DataProcessing []DataProcessing `xml:"dataProcessing"`
}

// DataProcessing is a description of the way in which a particular software was used
type DataProcessing struct {
	XMLName          xml.Name           `xml:"dataProcessing"`
	ID               string             `xml:"id,attr,omitempty"`
	ProcessingMethod []ProcessingMethod `xml:"processingMethod"`
}

// ProcessingMethod is the description of the default peak processing method. This element describes the base method used in the generation of
// a particular mzML file. Variable methods should be described in the appropriate acquisition section - if no acquisition-specific details
// are found, then this information serves as the default
type ProcessingMethod struct {
	XMLName                    xml.Name                     `xml:"processingMethod"`
	Order                      int                          `xml:"order,attr,omitempty"`
	SoftwareRef                string                       `xml:"softwareRef,attr,omitempty"`
	ReferenceableParamGroupRef []ReferenceableParamGroupRef `xml:"referenceableParamGroupRef"`
	CVParam                    []CVParam                    `xml:"cvParam"`
	UserParam                  []UserParam                  `xml:"userParam"`
}

// InstrumentConfigurationList is a list and descriptions of instrument configurations. At least one instrument configuration MUST be specified,
// even if it is only to specify that the instrument is unknown. In that case, the "instrument model" term is used to indicate the unknown
// instrument in the instrumentConfiguration
type InstrumentConfigurationList struct {
	XMLName                 xml.Name                  `xml:"instrumentConfigurationList"`
	Count                   int                       `xml:"count,attr"`
	InstrumentConfiguration []InstrumentConfiguration `xml:"instrumentConfiguration"`
}

// InstrumentConfiguration tag
type InstrumentConfiguration struct {
	XMLName                    xml.Name                     `xml:"instrumentConfiguration"`
	ID                         string                       `xml:"id,att,omitempty"`
	ScanSettingsRef            string                       `xml:"scanSettingsRef,att,omitempty"`
	ReferenceableParamGroupRef []ReferenceableParamGroupRef `xml:"referenceableParamGroupRef"`
	CVParam                    []CVParam                    `xml:"cvParam"`
	UserParam                  []UserParam                  `xml:"userParam"`
	ComponentList              ComponentList                `xml:"componentList"`
	SoftwareRef                SoftwareRef                  `xml:"softwareRef"`
}

// ComponentList is a list with the different components used in the mass spectrometer. At least one source, one mass analyzer and one detector need to be specified
type ComponentList struct {
	XMLName  xml.Name `xml:"componentList"`
	Count    int      `xml:"count,attr,omitempty"`
	Source   Source   `xml:"source"`
	Analyzer Analyzer `xml:"analyzer"`
	Detector Detector `xml:"detector"`
}

// Source is a source component
type Source struct {
	XMLName                    xml.Name                     `xml:"source"`
	Order                      int                          `xml:"order,attr,omitempty"`
	ReferenceableParamGroupRef []ReferenceableParamGroupRef `xml:"referenceableParamGroupRef"`
	CVParam                    []CVParam                    `xml:"cvParam"`
	UserParam                  []UserParam                  `xml:"userParam"`
}

// Analyzer is a analyzer component
type Analyzer struct {
	XMLName                    xml.Name                     `xml:"analyzer"`
	Order                      int                          `xml:"order,attr,omitempty"`
	ReferenceableParamGroupRef []ReferenceableParamGroupRef `xml:"referenceableParamGroupRef"`
	CVParam                    []CVParam                    `xml:"cvParam"`
	UserParam                  []UserParam                  `xml:"userParam"`
}

// Detector is a detector component
type Detector struct {
	XMLName                    xml.Name                     `xml:"detector"`
	Order                      int                          `xml:"order,attr,omitempty"`
	ReferenceableParamGroupRef []ReferenceableParamGroupRef `xml:"referenceableParamGroupRef"`
	CVParam                    []CVParam                    `xml:"cvParam"`
	UserParam                  []UserParam                  `xml:"userParam"`
}

// FileDescription contains the information pertaining to the entire mzML file (i.e. not specific to any part of the data set) is stored here
type FileDescription struct {
	XMLName        xml.Name       `xml:"fileDescription"`
	FileContent    FileContent    `xml:"fileContent"`
	SourceFileList SourceFileList `xml:"sourceFileList"`
}

// FileContent tag
type FileContent struct {
	XMLName                    xml.Name                     `xml:"fileContent"`
	ReferenceableParamGroupRef []ReferenceableParamGroupRef `xml:"referenceableParamGroupRef"`
	CVParam                    []CVParam                    `xml:"cvParam"`
	UserParam                  []UserParam                  `xml:"userParam"`
}

// SourceFileList tag
type SourceFileList struct {
	XMLName    xml.Name         `xml:"sourceFileList"`
	Count      int              `xml:"count,attr"`
	SourceFile []MzMLSourceFile `xml:"sourceFile"`
}

// SourceFileRefList is a list with the source files containing the acquisition settings
type SourceFileRefList struct {
	XMLName       xml.Name        `xml:"sourceFileRefList"`
	Count         int             `xml:"count,attr"`
	SourceFileRef []SourceFileRef `xml:"sourceFileRef"`
}

// SourceFileRef is a file from which this instance was created
type SourceFileRef struct {
	XMLName xml.Name `xml:"SourceFile"`
	Ref     string   `xml:"ref,attr,omitempty"`
}

// MzMLSourceFile is a file from which this instance was created
type MzMLSourceFile struct {
	XMLName                     xml.Name                    `xml:"sourceFile"`
	ID                          string                      `xml:"id,attr"`
	Location                    string                      `xml:"location,attr"`
	Name                        string                      `xml:"name,attr"`
	ExternalFormatDocumentation ExternalFormatDocumentation `xml:"ExternalFormatDocumentation"`
	FileFormat                  FileFormat                  `xml:"FileFormat"`
	CVParam                     []CVParam                   `xml:"cvParam"`
	UserParam                   []UserParam                 `xml:"userParam"`
}

// RefParamGroupList is the container for a list of referenceableParamGroups
type RefParamGroupList struct {
	XMLName                 xml.Name                  `xml:"referenceableParamGroupList"`
	Count                   int                       `xml:"count,attr"`
	ReferenceableParamGroup []ReferenceableParamGroup `xml:"referenceableParamGroup"`
}

// ReferenceableParamGroupRef is a reference to a previously defined ParamGroup, which is a reusable container of one or more cvParams
type ReferenceableParamGroupRef struct {
	XMLName xml.Name `xml:"referenceableParamGroupRef"`
	Ref     string   `xml:"id,ref"`
}

// ReferenceableParamGroup is a collection of CVParam and UserParam elements that can be referenced from elsewhere in this mzML
// document by using the 'paramGroupRef' element in that location to reference the 'id' attribute value of this element
type ReferenceableParamGroup struct {
	XMLName   xml.Name    `xml:"referenceableParamGroup"`
	ID        string      `xml:"id,attr"`
	CVParam   []CVParam   `xml:"cvParam"`
	UserParam []UserParam `xml:"userParam"`
}

// ScanSettingsList is a list with the descriptions of the acquisition settings applied prior to the start of data acquisition
type ScanSettingsList struct {
	XMLName      xml.Name       `xml:"scanSettingsList"`
	Count        int            `xml:"count,attr"`
	ScanSettings []ScanSettings `xml:"scanSettings"`
}

// ScanSettings contains the description of the acquisition settings of the instrument prior to the start of the run
type ScanSettings struct {
	XMLName                    xml.Name                     `xml:"scanSettings"`
	ID                         string                       `xml:"id,attr"`
	ReferenceableParamGroupRef []ReferenceableParamGroupRef `xml:"referenceableParamGroupRef"`
	CVParam                    []CVParam                    `xml:"cvParam"`
	UserParam                  []UserParam                  `xml:"userParam"`
	SourceFileRefList          SourceFileRefList            `xml:"sourceFileRefList"`
	TargetList                 TargetList                   `xml:"targetList"`
}

// SampleList is a list and descriptions of samples
type SampleList struct {
	XMLName xml.Name `xml:"sampleList"`
	Count   int      `xml:"count,attr"`
	Sample  []Sample `xml:"sample"`
}

// TargetList (or 'inclusion list') configured prior to the run
type TargetList struct {
	XMLName xml.Name `xml:"targetList"`
	Count   int      `xml:"count,attr"`
	Target  []Target `xml:"target"`
}

// Target is a structure allowing the use of a controlled (cvParam) or uncontrolled vocabulary (userParam),
// or a reference to a predefined set of these in this mzML file (paramGroupRef)
type Target struct {
	XMLName                    xml.Name                     `xml:"target"`
	ReferenceableParamGroupRef []ReferenceableParamGroupRef `xml:"referenceableParamGroupRef"`
	CVParam                    []CVParam                    `xml:"cvParam"`
	UserParam                  []UserParam                  `xml:"userParam"`
}

// SoftwareList contains descriptions of software used to acquire and/or process the data in this mzML file
type SoftwareList struct {
	XMLName  xml.Name   `xml:"softwareList"`
	Count    int        `xml:"count,attr"`
	Software []Software `xml:"software"`
}

// Software ia a piece of software
type Software struct {
	XMLName                    xml.Name                     `xml:"software"`
	ID                         string                       `xml:"id,attr"`
	Version                    string                       `xml:"version,attr"`
	ReferenceableParamGroupRef []ReferenceableParamGroupRef `xml:"referenceableParamGroupRef"`
	CVParam                    []CVParam                    `xml:"cvParam"`
	UserParam                  []UserParam                  `xml:"userParam"`
}

// SoftwareRef is a reference to a previously defined software element
type SoftwareRef struct {
	XMLName xml.Name `xml:"softwareRef"`
	Ref     string   `xml:"ref,attr,omitempty"`
}

// Run tag
type Run struct {
	XMLName                           xml.Name                     `xml:"run"`
	DefaultInstrumentConfigurationRef string                       `xml:"defaultInstrumentConfigurationRef,attr,omitempy"`
	DefaultSourceFileRef              string                       `xml:"defaultSourceFileRef,attr,omitempy"`
	ID                                string                       `xml:"id,attr,omitempy"`
	SampleRef                         string                       `xml:"sampleRef,attr,omitempy"`
	StartTimeStamp                    string                       `xml:"startTimeStamp,att,omitempty"`
	ReferenceableParamGroupRef        []ReferenceableParamGroupRef `xml:"referenceableParamGroupRef"`
	CVParam                           []CVParam                    `xml:"cvParam"`
	UserParam                         []UserParam                  `xml:"userParam"`
	SpectrumList                      SpectrumList                 `xml:"spectrumList"`
	ChromatogramList                  ChromatogramList             `xml:"chromatogramList"`
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
	DataProcessingRef   string              `xml:"dataProcessingRef,att"`
	DefaultArrayLength  float64             `xml:"defaultArrayLength,attr"`
	ID                  string              `xml:"id,attr"`
	Index               string              `xml:"index,attr"`
	SourceFileRef       string              `xml:"sourceFileRef,attr"`
	SpotID              string              `xml:"spotID,attr"`
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
	InstConfigurationRef string      `xml:"isolationWindow,attr"`
	CVParam              []CVParam   `xml:"cvParam"`
	UserParam            []UserParam `xml:"userParam"`
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
	UserParam           []UserParam         `xml:"userParam"`
	Precursor           Precursor           `xml:"precursor"`
	Product             Product             `xml:"product"`
	BinaryDataArrayList BinaryDataArrayList `xml:"binaryDataArrayList"`
}

// Product is the method of product ion selection and activation in a precursor ion scan
type Product struct {
	XMLName         xml.Name        `xml:"product"`
	IsolationWindow IsolationWindow `xml:"isolationWindow"`
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
