package mz

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"io"
	"math"
	"os"
	"strconv"
	"sync"

	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/raw/mzml"
)

// Raw struct
type Raw struct {
	FileName   string
	RefSpectra sync.Map
	Spectra    Spectra
}

// Spectra is a list of Spetrum
type Spectra []Spectrum

// Spectrum struct
type Spectrum struct {
	Index       string
	Scan        string
	Level       string
	StartTime   string
	Precursor   Precursor
	Peaks       Peaks
	Intensities Intensities
}

// Peaks struct
type Peaks struct {
	Stream        []byte
	DecodedStream []float64
	Precision     string
	Compression   string
}

// Intensities struct
type Intensities struct {
	Stream        []byte
	DecodedStream []float64
	Precision     string
	Compression   string
}

// Precursor struct
type Precursor struct {
	ParentIndex                string
	ParentScan                 string
	ChargeState                string
	SelectedIon                string
	TargetIon                  string
	IsolationWindowLowerOffset string
	IsolationWindowUpperOffset string
	PeakIntensity              string
}

// Read is a simple mzML reader
func (r *Raw) Read(f string) *err.Error {

	xmlFile, e := os.Open(f)
	if e != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: e.Error()}
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder((bufio.NewReader(xmlFile)))

	var inElement string
	for {

		t, _ := decoder.Token()
		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:

			inElement = se.Name.Local

			if inElement == "spectrum" {
				var rawSpec mzml.Spectrum
				decoder.DecodeElement(&rawSpec, &se)
				_ = decoder.Decode(&se)

				var spec Spectrum
				spec.Index = string(rawSpec.Index)

				indexStr := string(rawSpec.Index)
				indexInt, _ := strconv.Atoi(indexStr)
				indexInt++
				spec.Scan = string(strconv.Itoa(indexInt))

				for _, j := range rawSpec.CVParam {
					if string(j.Accession) == "MS:1000511" {
						spec.Level = string(j.Value)
					}
				}

				for _, j := range rawSpec.ScanList.Scan[0].CVParam {
					if string(j.Accession) == "MS:1000016" {
						spec.StartTime = string(j.Value)
					}
				}

				spec.Precursor = Precursor{}
				if len(rawSpec.PrecursorList.Precursor) > 0 {
					for _, j := range rawSpec.PrecursorList.Precursor[0].IsolationWindow.CVParam {
						if string(j.Accession) == "MS:1000828" {
							spec.Precursor.IsolationWindowLowerOffset = string(j.Value)
						}

						if string(j.Accession) == "MS:1000829" {
							spec.Precursor.IsolationWindowUpperOffset = string(j.Value)
						}
					}

					for _, j := range rawSpec.PrecursorList.Precursor[0].SelectedIonList.SelectedIon[0].CVParam {
						if string(j.Accession) == "MS:1000744" {
							spec.Precursor.SelectedIon = string(j.Value)
						}

						if string(j.Accession) == "MS:1000041" {
							spec.Precursor.ChargeState = string(j.Value)
						}

						if string(j.Accession) == "MS:1000042" {
							spec.Precursor.PeakIntensity = string(j.Value)
						}
					}
				}

				//TODO TARGETION

				var binPeak Peaks
				binPeak.Stream = rawSpec.BinaryDataArrayList.BinaryDataArray[0].Binary.Value
				for _, j := range rawSpec.BinaryDataArrayList.BinaryDataArray[0].CVParam {
					if string(j.Accession) == "MS:1000523" {
						binPeak.Precision = "64"
					} else if string(j.Accession) == "MS:1000521" {
						binPeak.Precision = "32"
					}

					if string(j.Accession) == "MS:1000574" {
						binPeak.Compression = "1"
					} else if string(j.Accession) == "MS:1000576" {
						binPeak.Compression = "0"
					}
				}

				spec.Peaks = binPeak
				spec.Peaks.DecodedStream, _ = Decode("mz", rawSpec.BinaryDataArrayList.BinaryDataArray[0])
				spec.Peaks.Stream = nil

				var binInt Intensities
				binInt.Stream = rawSpec.BinaryDataArrayList.BinaryDataArray[1].Binary.Value
				for _, j := range rawSpec.BinaryDataArrayList.BinaryDataArray[1].CVParam {
					if string(j.Accession) == "MS:1000523" {
						binInt.Precision = "64"
					} else if string(j.Accession) == "MS:1000521" {
						binInt.Precision = "32"
					}

					if string(j.Accession) == "MS:1000574" {
						binInt.Compression = "1"
					} else if string(j.Accession) == "MS:1000576" {
						binInt.Compression = "0"
					}
				}

				spec.Intensities = binInt
				spec.Intensities.DecodedStream, _ = Decode("int", rawSpec.BinaryDataArrayList.BinaryDataArray[1])
				spec.Intensities.Stream = nil

				r.Spectra = append(r.Spectra, spec)

				//nil
				spec = Spectrum{}
				rawSpec = mzml.Spectrum{}
			}

		default:

		}

	}

	decoder = nil

	return nil
}

// ParRead is a parallel reader implementing sync.Map
func (r *Raw) ParRead(f string) *err.Error {

	xmlFile, e := os.Open(f)
	if e != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: e.Error()}
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder((bufio.NewReader(xmlFile)))

	var inElement string
	for {

		t, _ := decoder.Token()
		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:

			inElement = se.Name.Local

			if inElement == "spectrum" {
				var rawSpec mzml.Spectrum
				decoder.DecodeElement(&rawSpec, &se)
				//				_ = decoder.Decode(&se)
				go procSpectra(r, rawSpec)
			}

		default:

		}

	}

	decoder = nil

	return nil
}

// func (r *Raw) ParRead(f string) *err.Error {
//
// 	xmlFile, e := os.Open(f)
// 	if e != nil {
// 		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: e.Error()}
// 	}
// 	defer xmlFile.Close()
//
// 	decoder := xml.NewDecoder((bufio.NewReader(xmlFile)))
//
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
// 				var rawSpec mzml.Spectrum
// 				decoder.DecodeElement(&rawSpec, &se)
//
// 				go procSpectra(r, rawSpec)
//
// 				//_ = decoder.Decode(&se)
// 				// var spec Spectrum
// 				// spec.Index = string(rawSpec.Index)
// 				//
// 				// indexStr := string(rawSpec.Index)
// 				// indexInt, _ := strconv.Atoi(indexStr)
// 				// indexInt++
// 				// spec.Scan = string(strconv.Itoa(indexInt))
// 				//
// 				// for _, j := range rawSpec.CVParam {
// 				// 	if string(j.Accession) == "MS:1000511" {
// 				// 		spec.Level = string(j.Value)
// 				// 	}
// 				// }
// 				//
// 				// for _, j := range rawSpec.ScanList.Scan[0].CVParam {
// 				// 	if string(j.Accession) == "MS:1000016" {
// 				// 		spec.StartTime = string(j.Value)
// 				// 	}
// 				// }
// 				//
// 				// spec.Precursor = Precursor{}
// 				// if len(rawSpec.PrecursorList.Precursor) > 0 {
// 				// 	for _, j := range rawSpec.PrecursorList.Precursor[0].IsolationWindow.CVParam {
// 				// 		if string(j.Accession) == "MS:1000828" {
// 				// 			spec.Precursor.IsolationWindowLowerOffset = string(j.Value)
// 				// 		}
// 				//
// 				// 		if string(j.Accession) == "MS:1000829" {
// 				// 			spec.Precursor.IsolationWindowUpperOffset = string(j.Value)
// 				// 		}
// 				// 	}
// 				//
// 				// 	for _, j := range rawSpec.PrecursorList.Precursor[0].SelectedIonList.SelectedIon[0].CVParam {
// 				// 		if string(j.Accession) == "MS:1000744" {
// 				// 			spec.Precursor.SelectedIon = string(j.Value)
// 				// 		}
// 				//
// 				// 		if string(j.Accession) == "MS:1000041" {
// 				// 			spec.Precursor.ChargeState = string(j.Value)
// 				// 		}
// 				//
// 				// 		if string(j.Accession) == "MS:1000042" {
// 				// 			spec.Precursor.PeakIntensity = string(j.Value)
// 				// 		}
// 				// 	}
// 				// }
// 				//
// 				// //TODO TARGETION
// 				//
// 				// var binPeak Peaks
// 				// binPeak.Stream = rawSpec.BinaryDataArrayList.BinaryDataArray[0].Binary.Value
// 				// for _, j := range rawSpec.BinaryDataArrayList.BinaryDataArray[0].CVParam {
// 				// 	if string(j.Accession) == "MS:1000523" {
// 				// 		binPeak.Precision = "64"
// 				// 	} else if string(j.Accession) == "MS:1000521" {
// 				// 		binPeak.Precision = "32"
// 				// 	}
// 				//
// 				// 	if string(j.Accession) == "MS:1000574" {
// 				// 		binPeak.Compression = "1"
// 				// 	} else if string(j.Accession) == "MS:1000576" {
// 				// 		binPeak.Compression = "0"
// 				// 	}
// 				// }
// 				//
// 				// spec.Peaks = binPeak
// 				// spec.Peaks.DecodedStream, _ = Decode("mz", rawSpec.BinaryDataArrayList.BinaryDataArray[0])
// 				// spec.Peaks.Stream = nil
// 				//
// 				// var binInt Intensities
// 				// binInt.Stream = rawSpec.BinaryDataArrayList.BinaryDataArray[1].Binary.Value
// 				// for _, j := range rawSpec.BinaryDataArrayList.BinaryDataArray[1].CVParam {
// 				// 	if string(j.Accession) == "MS:1000523" {
// 				// 		binInt.Precision = "64"
// 				// 	} else if string(j.Accession) == "MS:1000521" {
// 				// 		binInt.Precision = "32"
// 				// 	}
// 				//
// 				// 	if string(j.Accession) == "MS:1000574" {
// 				// 		binInt.Compression = "1"
// 				// 	} else if string(j.Accession) == "MS:1000576" {
// 				// 		binInt.Compression = "0"
// 				// 	}
// 				// }
// 				//
// 				// spec.Intensities = binInt
// 				// spec.Intensities.DecodedStream, _ = Decode("int", rawSpec.BinaryDataArrayList.BinaryDataArray[1])
// 				// spec.Intensities.Stream = nil
// 				//
// 				// r.RefSpectra.Store(spec.Scan, spec)
// 				//
// 				// //nil
// 				// spec = Spectrum{}
// 				// rawSpec = mzml.Spectrum{}
// 			}
//
// 		default:
//
// 		}
//
// 	}
//
// 	decoder = nil
//
// 	return nil
// }

func procSpectra(r *Raw, rawSpec mzml.Spectrum) {

	var spec Spectrum
	spec.Index = string(rawSpec.Index)

	indexStr := string(rawSpec.Index)
	indexInt, _ := strconv.Atoi(indexStr)
	indexInt++
	spec.Scan = string(strconv.Itoa(indexInt))

	for _, j := range rawSpec.CVParam {
		if string(j.Accession) == "MS:1000511" {
			spec.Level = string(j.Value)
		}
	}

	for _, j := range rawSpec.ScanList.Scan[0].CVParam {
		if string(j.Accession) == "MS:1000016" {
			spec.StartTime = string(j.Value)
		}
	}

	spec.Precursor = Precursor{}
	if len(rawSpec.PrecursorList.Precursor) > 0 {
		for _, j := range rawSpec.PrecursorList.Precursor[0].IsolationWindow.CVParam {
			if string(j.Accession) == "MS:1000828" {
				spec.Precursor.IsolationWindowLowerOffset = string(j.Value)
			}

			if string(j.Accession) == "MS:1000829" {
				spec.Precursor.IsolationWindowUpperOffset = string(j.Value)
			}
		}

		for _, j := range rawSpec.PrecursorList.Precursor[0].SelectedIonList.SelectedIon[0].CVParam {
			if string(j.Accession) == "MS:1000744" {
				spec.Precursor.SelectedIon = string(j.Value)
			}

			if string(j.Accession) == "MS:1000041" {
				spec.Precursor.ChargeState = string(j.Value)
			}

			if string(j.Accession) == "MS:1000042" {
				spec.Precursor.PeakIntensity = string(j.Value)
			}
		}
	}

	//TODO TARGETION

	var binPeak Peaks
	binPeak.Stream = rawSpec.BinaryDataArrayList.BinaryDataArray[0].Binary.Value
	for _, j := range rawSpec.BinaryDataArrayList.BinaryDataArray[0].CVParam {
		if string(j.Accession) == "MS:1000523" {
			binPeak.Precision = "64"
		} else if string(j.Accession) == "MS:1000521" {
			binPeak.Precision = "32"
		}

		if string(j.Accession) == "MS:1000574" {
			binPeak.Compression = "1"
		} else if string(j.Accession) == "MS:1000576" {
			binPeak.Compression = "0"
		}
	}

	spec.Peaks = binPeak
	spec.Peaks.DecodedStream, _ = Decode("mz", rawSpec.BinaryDataArrayList.BinaryDataArray[0])
	spec.Peaks.Stream = nil

	var binInt Intensities
	binInt.Stream = rawSpec.BinaryDataArrayList.BinaryDataArray[1].Binary.Value
	for _, j := range rawSpec.BinaryDataArrayList.BinaryDataArray[1].CVParam {
		if string(j.Accession) == "MS:1000523" {
			binInt.Precision = "64"
		} else if string(j.Accession) == "MS:1000521" {
			binInt.Precision = "32"
		}

		if string(j.Accession) == "MS:1000574" {
			binInt.Compression = "1"
		} else if string(j.Accession) == "MS:1000576" {
			binInt.Compression = "0"
		}
	}

	spec.Intensities = binInt
	spec.Intensities.DecodedStream, _ = Decode("int", rawSpec.BinaryDataArrayList.BinaryDataArray[1])
	spec.Intensities.Stream = nil

	r.RefSpectra.Store(spec.Scan, spec)

	//nil
	spec = Spectrum{}
	rawSpec = mzml.Spectrum{}

	return
}

// Decode processes the binary data
func Decode(class string, bin mzml.BinaryDataArray) ([]float64, error) {

	var compression bool
	var precision string
	var err error

	for i := range bin.CVParam {

		if string(bin.CVParam[i].Accession) == "MS:1000523" {
			precision = "64"
		} else if string(bin.CVParam[i].Accession) == "MS:1000521" {
			precision = "32"
		}

		if string(bin.CVParam[i].Accession) == "MS:1000574" {
			compression = true
		} else if string(bin.CVParam[i].Accession) == "MS:1000576" {
			compression = false
		}

	}

	f, err := readEncoded(class, bin, precision, compression)
	if err != nil {
		return f, err
	}

	return f, nil
}

// readEncoded transforms the binary data into float64 values
func readEncoded(class string, bin mzml.BinaryDataArray, precision string, isCompressed bool) ([]float64, error) {

	var stream []uint8
	var floatArray []float64

	b := bytes.NewReader(bin.Binary.Value)
	b64 := base64.NewDecoder(base64.StdEncoding, b)

	var bytestream bytes.Buffer
	if isCompressed == true {
		r, err := zlib.NewReader(b64)
		if err != nil {
			return floatArray, err
		}
		io.Copy(&bytestream, r)
	} else {
		io.Copy(&bytestream, b64)
	}

	dataArray := bytestream.Bytes()

	var counter int

	if precision == "32" {
		for i := range dataArray {
			counter++
			stream = append(stream, dataArray[i])
			if counter == 4 {
				bits := binary.LittleEndian.Uint32(stream)
				converted := math.Float32frombits(bits)

				if class == "mz" {
					//floatArray = append(floatArray, utils.Round(float64(converted), 5, 6))
					floatArray = append(floatArray, float64(converted))
				} else if class == "int" {
					//floatArray = append(floatArray, utils.Round(float64(converted), 5, 6))
					floatArray = append(floatArray, float64(converted))
				}

				stream = nil
				counter = 0
			}
		}
	} else if precision == "64" {
		for i := range dataArray {
			counter++
			stream = append(stream, dataArray[i])
			if counter == 8 {
				bits := binary.LittleEndian.Uint64(stream)
				converted := math.Float64frombits(bits)

				if class == "mz" {
					//floatArray = append(floatArray, utils.Round(float64(converted), 5, 6))
					floatArray = append(floatArray, float64(converted))
				} else if class == "int" {
					//floatArray = append(floatArray, utils.Round(float64(converted), 5, 6))
					floatArray = append(floatArray, float64(converted))
				}

				stream = nil
				counter = 0
			}
		}
	} else {
		return floatArray, errors.New("Undefined binary precision")
	}

	return floatArray, nil
}
