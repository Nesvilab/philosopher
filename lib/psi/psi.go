package psi

import "encoding/xml"

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
