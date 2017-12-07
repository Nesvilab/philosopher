package mzml

import (
	"encoding/xml"

	// anon charset
	_ "github.com/rogpeppe/go-charset/data"
)

// IndexedMzML is the root level tag
type IndexedMzML struct {
	XMLName xml.Name `xml:"indexedmzML"`
	MzML    MzML     `xml:"mzML"`
}

// MzML is the root level tag
type MzML struct {
	XMLName           xml.Name          `xml:"mzML"`
	Accession         []byte            `xml:"accession,attr"`
	Version           []byte            `xml:"version,attr"`
	FileDescription   FileDescription   `ml:"fileDescription"`
	RefParamGroupList RefParamGroupList `xml:"referenceableParamGroupList"`
	SoftwareList      SoftwareList      `xml:"softwareList"`
	Run               Run               `xml:"run"`
}

// CvList tag
type CvList struct {
	XMLName xml.Name `xml:"cvList"`
	Count   []byte   `xml:"count,attr"`
	CV      []CV     `xml:"cv"`
}

// CV tag
type CV struct {
	XMLName  xml.Name `xml:"cv"`
	ID       []byte   `xml:"id,attr"`
	Version  []byte   `xml:"version,attr"`
	URI      []byte   `xml:"URI,attr"`
	FullName []byte   `xml:"fullName,attr"`
}

// CVParam tag
type CVParam struct {
	XMLName       xml.Name `xml:"cvParam"`
	CVRef         []byte   `xml:"cvRef,attr"`
	Accession     []byte   `xml:"accession,attr"`
	Name          []byte   `xml:"name,attr"`
	Value         []byte   `xml:"value,attr"`
	UnitCvRef     []byte   `xml:"unitCvRef,attr"`
	UnitAccession []byte   `xml:"unitAccession,attr"`
	UnitName      []byte   `xml:"unitName,attr"`
}

// UserParam tag
type UserParam struct {
	XMLName xml.Name `xml:"userParam"`
	Name    []byte   `xml:"name,attr"`
	Value   []byte   `xml:"value,attr"`
	Type    []byte   `xml:"type,attr"`
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
	Count      []byte       `xml:"count,attr"`
	SourceFile []SourceFile `xml:"sourceFile"`
}

// SourceFile tag
type SourceFile struct {
	XMLName  xml.Name  `xml:"sourceFile"`
	ID       []byte    `xml:"id,attr"`
	Name     []byte    `xml:"name,attr"`
	Location []byte    `xml:"location,attr"`
	CVParam  []CVParam `xml:"cvParam"`
}

// RefParamGroupList tag
type RefParamGroupList struct {
	XMLName       xml.Name        `xml:"referenceableParamGroupList"`
	Count         []byte          `xml:"count,attr"`
	RefParamGroup []RefParamGroup `xml:"referenceableParamGroup"`
}

// RefParamGroup tag
type RefParamGroup struct {
	XMLName   xml.Name    `xml:"referenceableParamGroup"`
	ID        []byte      `xml:"id,attr"`
	CVParam   []CVParam   `xml:"cvParam"`
	UserParam []UserParam `xml:"userParam"`
}

// SoftwareList tag
type SoftwareList struct {
	XMLName  xml.Name   `xml:"softwareList"`
	Count    []byte     `xml:"count,attr"`
	Software []Software `xml:"software"`
}

// Software tag
type Software struct {
	XMLName       xml.Name        `xml:"software"`
	ID            []byte          `xml:"id,attr"`
	Version       []byte          `xml:"version,attr"`
	CVParam       []CVParam       `xml:"cvParam"`
	UserParam     []UserParam     `xml:"userParam"`
	RefParamGroup []RefParamGroup `xml:"referenceableParamGroup"`
}

// Run tag
type Run struct {
	XMLName                           xml.Name         `xml:"run"`
	ID                                []byte           `xml:"id,attr"`
	DefaultInstrumentConfigurationRef []byte           `xml:"defaultInstrumentConfigurationRef,attr"`
	StartTimeStamp                    []byte           `xml:"startTimeStamp,attr"`
	SpectrumList                      SpectrumList     `xml:"spectrumList"`
	ChromatogramList                  ChromatogramList `xml:"chromatogramList"`
}

// SpectrumList tag
type SpectrumList struct {
	XMLName                  xml.Name   `xml:"spectrumList"`
	Count                    []byte     `xml:"count,attr"`
	DefaultDataProcessingRef []byte     `xml:"defaultDataProcessingRef,attr"`
	Spectrum                 []Spectrum `xml:"spectrum"`
}

// Spectrum tag
type Spectrum struct {
	XMLName             xml.Name            `xml:"spectrum"`
	Index               []byte              `xml:"index,attr"`
	ID                  []byte              `xml:"id,attr"`
	DefaultArrayLength  []byte              `xml:"defaultArrayLength,attr"`
	CVParam             []CVParam           `xml:"cvParam"`
	ScanList            ScanList            `xml:"scanList"`
	PrecursorList       PrecursorList       `xml:"precursorList"`
	BinaryDataArrayList BinaryDataArrayList `xml:"binaryDataArrayList"`
	Peaks               []byte
	Intensities         []byte
}

// ScanList tag
type ScanList struct {
	XMLName xml.Name  `xml:"scanList"`
	Count   []byte    `xml:"count,attr"`
	CVParam []CVParam `xml:"cvParam"`
	Scan    []Scan    `xml:"scan"`
}

// PrecursorList tag
type PrecursorList struct {
	XMLName   xml.Name    `xml:"precursorList"`
	Count     []byte      `xml:"count,attr"`
	Precursor []Precursor `xml:"precursor"`
}

// Precursor tag
type Precursor struct {
	XMLName         xml.Name        `xml:"precursor"`
	SpectrumRef     []byte          `xml:"spectrumRef,attr"`
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
	Count       []byte        `xml:"count,attr"`
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
	InstConfigurationRef []byte         `xml:"instrumentConfigurationRef,attr"`
	CVParam              []CVParam      `xml:"cvParam"`
	UserParam            []UserParam    `xml:"userParam"`
	ScanWindowList       ScanWindowList `xml:"scanWindowList"`
}

// ScanWindowList tag
type ScanWindowList struct {
	XMLName    xml.Name     `xml:"scanWindowList"`
	Count      []byte       `xml:"count,attr"`
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
	Count                    []byte         `xml:"count,attr"`
	DefaultDataProcessingRef []byte         `xml:"defaultDataProcessingRef,attr"`
	Chromatogram             []Chromatogram `xml:"chromatogram"`
}

// Chromatogram tag
type Chromatogram struct {
	XMLName             xml.Name            `xml:"chromatogram"`
	Index               []byte              `xml:"index,attr"`
	ID                  []byte              `xml:"id,attr"`
	DefaultArrayLength  []byte              `xml:"defaultArrayLength,attr"`
	CVParam             []CVParam           `xml:"cvParam"`
	BinaryDataArrayList BinaryDataArrayList `xml:"binaryDataArrayList"`
}

// BinaryDataArrayList tag
type BinaryDataArrayList struct {
	XMLName         xml.Name          `xml:"binaryDataArrayList"`
	Count           []byte            `xml:"count,attr"`
	BinaryDataArray []BinaryDataArray `xml:"binaryDataArray"`
}

// BinaryDataArray tag
type BinaryDataArray struct {
	XMLName       xml.Name  `xml:"binaryDataArray"`
	EncodedLength []byte    `xml:"encodedLength,attr"`
	CVParam       []CVParam `xml:"cvParam"`
	Binary        Binary    `xml:"binary"`
	//ConvertedBinary []float64
}

// Binary tag
type Binary struct {
	XMLName xml.Name `xml:"binary"`
	Value   []byte   `xml:",chardata"`
}

// Read and process mzML spectral data
// func Read(f string) (mz.Raw, *err.Error) {
//
// 	xmlFile, e := os.Open(f)
// 	if e != nil {
// 		return mz.Raw{}, &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: e.Error()}
// 	}
// 	defer xmlFile.Close()
//
// 	decoder := xml.NewDecoder((bufio.NewReader(xmlFile)))
//
// 	var list mz.Spectra
// 	var inElement string
// 	for {
//
// 		t, _ := decoder.Token()
// 		if t == nil {
// 			break
// 		}
//
// 		switch se := t.(type) {
// 		case xml.StartElement:
//
// 			inElement = se.Name.Local
//
// 			if inElement == "spectrum" {
// 				var rawSpec Spectrum
// 				decoder.DecodeElement(&rawSpec, &se)
//
// 				var spec mz.Spectrum
// 				spec.Index = string(rawSpec.Index)
//
// 				indexStr := string(rawSpec.Index)
// 				indexInt, _ := strconv.Atoi(indexStr)
// 				indexInt++
// 				spec.Scan = string(strconv.Itoa(indexInt))
//
// 				for _, j := range rawSpec.CVParam {
// 					if string(j.Accession) == "MS:1000511" {
// 						spec.Level = string(j.Value)
// 					}
// 				}
//
// 				for _, j := range rawSpec.ScanList.Scan[0].CVParam {
// 					if string(j.Accession) == "MS:1000016" {
// 						spec.StartTime = string(j.Value)
// 					}
// 				}
//
// 				spec.Precursor = mz.Precursor{}
// 				if len(rawSpec.PrecursorList.Precursor) > 0 {
// 					for _, j := range rawSpec.PrecursorList.Precursor[0].IsolationWindow.CVParam {
// 						if string(j.Accession) == "MS:1000828" {
// 							spec.Precursor.IsolationWindowLowerOffset = string(j.Value)
// 						}
//
// 						if string(j.Accession) == "MS:1000829" {
// 							spec.Precursor.IsolationWindowUpperOffset = string(j.Value)
// 						}
// 					}
//
// 					for _, j := range rawSpec.PrecursorList.Precursor[0].SelectedIonList.SelectedIon[0].CVParam {
// 						if string(j.Accession) == "MS:1000744" {
// 							spec.Precursor.SelectedIon = string(j.Value)
// 						}
//
// 						if string(j.Accession) == "MS:1000041" {
// 							spec.Precursor.ChargeState = string(j.Value)
// 						}
//
// 						if string(j.Accession) == "MS:1000042" {
// 							spec.Precursor.PeakIntensity = string(j.Value)
// 						}
// 					}
// 				}
//
// 				//TODO TARGETION
//
// 				var binPeak mz.Peaks
// 				binPeak.Stream = rawSpec.BinaryDataArrayList.BinaryDataArray[0].Binary.Value
// 				for _, j := range rawSpec.BinaryDataArrayList.BinaryDataArray[0].CVParam {
// 					if string(j.Accession) == "MS:1000523" {
// 						binPeak.Precision = "64"
// 					} else if string(j.Accession) == "MS:1000521" {
// 						binPeak.Precision = "32"
// 					}
//
// 					if string(j.Accession) == "MS:1000574" {
// 						binPeak.Compression = "1"
// 					} else if string(j.Accession) == "MS:1000576" {
// 						binPeak.Compression = "0"
// 					}
// 				}
//
// 				spec.Peaks = binPeak
// 				spec.Peaks.DecodedStream, _ = Decode("mz", rawSpec.BinaryDataArrayList.BinaryDataArray[0])
// 				spec.Peaks.Stream = nil
//
// 				var binInt mz.Intensities
// 				binInt.Stream = rawSpec.BinaryDataArrayList.BinaryDataArray[1].Binary.Value
// 				for _, j := range rawSpec.BinaryDataArrayList.BinaryDataArray[1].CVParam {
// 					if string(j.Accession) == "MS:1000523" {
// 						binInt.Precision = "64"
// 					} else if string(j.Accession) == "MS:1000521" {
// 						binInt.Precision = "32"
// 					}
//
// 					if string(j.Accession) == "MS:1000574" {
// 						binInt.Compression = "1"
// 					} else if string(j.Accession) == "MS:1000576" {
// 						binInt.Compression = "0"
// 					}
// 				}
//
// 				spec.Intensities = binInt
// 				spec.Intensities.DecodedStream, _ = Decode("int", rawSpec.BinaryDataArrayList.BinaryDataArray[1])
// 				spec.Intensities.Stream = nil
//
// 				list = append(list, spec)
//
// 				//nil
// 				spec = mz.Spectrum{}
//
// 			}
//
// 		default:
//
// 		}
//
// 	}
//
// 	var raw mz.Raw
// 	raw.FileName = f
// 	raw.Spectra = list
//
// 	return raw, nil
// }

// // Decode processes the binary data
// func Decode(class string, bin BinaryDataArray) ([]float64, error) {
//
// 	var compression bool
// 	var precision string
// 	var err error
//
// 	for i := range bin.CVParam {
//
// 		if string(bin.CVParam[i].Accession) == "MS:1000523" {
// 			precision = "64"
// 		} else if string(bin.CVParam[i].Accession) == "MS:1000521" {
// 			precision = "32"
// 		}
//
// 		if string(bin.CVParam[i].Accession) == "MS:1000574" {
// 			compression = true
// 		} else if string(bin.CVParam[i].Accession) == "MS:1000576" {
// 			compression = false
// 		}
//
// 	}
//
// 	f, err := readEncoded(class, bin, precision, compression)
// 	if err != nil {
// 		return f, err
// 	}
//
// 	return f, nil
// }
//
// // readEncoded transforms the binary data into float64 values
// func readEncoded(class string, bin BinaryDataArray, precision string, isCompressed bool) ([]float64, error) {
//
// 	var stream []uint8
// 	var floatArray []float64
//
// 	b := bytes.NewReader(bin.Binary.Value)
// 	b64 := base64.NewDecoder(base64.StdEncoding, b)
//
// 	var bytestream bytes.Buffer
// 	if isCompressed == true {
// 		r, err := zlib.NewReader(b64)
// 		if err != nil {
// 			return floatArray, err
// 		}
// 		io.Copy(&bytestream, r)
// 	} else {
// 		io.Copy(&bytestream, b64)
// 	}
//
// 	dataArray := bytestream.Bytes()
//
// 	var counter int
//
// 	if precision == "32" {
// 		for i := range dataArray {
// 			counter++
// 			stream = append(stream, dataArray[i])
// 			if counter == 4 {
// 				bits := binary.LittleEndian.Uint32(stream)
// 				converted := math.Float32frombits(bits)
//
// 				if class == "mz" {
// 					//floatArray = append(floatArray, utils.Round(float64(converted), 5, 6))
// 					floatArray = append(floatArray, float64(converted))
// 				} else if class == "int" {
// 					//floatArray = append(floatArray, utils.Round(float64(converted), 5, 6))
// 					floatArray = append(floatArray, float64(converted))
// 				}
//
// 				stream = nil
// 				counter = 0
// 			}
// 		}
// 	} else if precision == "64" {
// 		for i := range dataArray {
// 			counter++
// 			stream = append(stream, dataArray[i])
// 			if counter == 8 {
// 				bits := binary.LittleEndian.Uint64(stream)
// 				converted := math.Float64frombits(bits)
//
// 				if class == "mz" {
// 					//floatArray = append(floatArray, utils.Round(float64(converted), 5, 6))
// 					floatArray = append(floatArray, float64(converted))
// 				} else if class == "int" {
// 					//floatArray = append(floatArray, utils.Round(float64(converted), 5, 6))
// 					floatArray = append(floatArray, float64(converted))
// 				}
//
// 				stream = nil
// 				counter = 0
// 			}
// 		}
// 	} else {
// 		return floatArray, errors.New("Undefined binary precision")
// 	}
//
// 	return floatArray, nil
// }
