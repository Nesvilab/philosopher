package mzn

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/prvst/philosopher/lib/msg"

	"github.com/prvst/philosopher/lib/psi"
	"github.com/sirupsen/logrus"
)

// MsData top struct
type MsData struct {
	FileName string
	//	RefSpectra sync.Map
	Spectra Spectra
}

// Spectra struct
type Spectra []Spectrum

// Spectrum tag
type Spectrum struct {
	Index         string
	Scan          string
	Level         string
	SpectrumName  string
	ScanStartTime float64
	Precursor     Precursor
	Mz            Mz
	Intensity     Intensity
	IonMobility   IonMobility
}

// Precursor struct
type Precursor struct {
	ParentIndex                string
	ParentScan                 string
	ChargeState                int
	SelectedIon                float64
	TargetIon                  float64
	PeakIntensity              float64
	IsolationWindowLowerOffset float64
	IsolationWindowUpperOffset float64
}

// Mz struct
type Mz struct {
	Stream        []byte
	DecodedStream []float64
	Precision     string
	Compression   string
}

// Intensity struct
type Intensity struct {
	Stream        []byte
	DecodedStream []float64
	Precision     string
	Compression   string
}

// IonMobility struct
type IonMobility struct {
	Stream        []byte
	DecodedStream []float64
	Precision     string
	Compression   string
}

func (a Spectra) Len() int           { return len(a) }
func (a Spectra) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Spectra) Less(i, j int) bool { return a[i].Index < a[j].Index }

// Read is the main function for parsing mzML data
func (p *MsData) Read(f string, skipMS1, skipMS2, skipMS3 bool) {

	var xml psi.IndexedMzML
	xml.Parse(f)

	p.FileName = f

	var spectra Spectra
	sl := xml.MzML.Run.SpectrumList

	for _, i := range sl.Spectrum {

		var level string
		for _, j := range i.CVParam {
			if string(j.Accession) == "MS:1000511" {
				level = j.Value
			}
		}

		if skipMS1 == true && level == "1" {
			continue
		} else if skipMS2 == true && level == "2" {
			continue
		} else if skipMS3 == true && level == "3" {
			continue
		}

		spectrum := processSpectrum(i)

		spectra = append(spectra, spectrum)
	}

	if len(spectra) == 0 {
		msg.NoSpectraFound(errors.New(""), "fatal")
	}

	p.Spectra = spectra

	return
}

func processSpectrum(mzSpec psi.Spectrum) Spectrum {

	var spec Spectrum

	spec.Index = string(mzSpec.Index)

	indexStr := string(mzSpec.Index)
	indexInt, _ := strconv.Atoi(indexStr)
	indexInt++
	spec.Scan = string(strconv.Itoa(indexInt))

	for _, j := range mzSpec.CVParam {
		if string(j.Accession) == "MS:1000511" {
			spec.Level = j.Value
		}
	}

	for _, j := range mzSpec.ScanList.Scan[0].CVParam {
		if string(j.Accession) == "MS:1000016" {
			val, e := strconv.ParseFloat(j.Value, 64)
			if e != nil {
				msg.CastFloatToString(e, "error")
			}
			spec.ScanStartTime = val
		}
	}

	spec.Precursor = Precursor{}
	if len(mzSpec.PrecursorList.Precursor) > 0 {

		// parent index and parent scan
		var ref []string
		var precRef []string

		if len(mzSpec.PrecursorList.Precursor[0].SpectrumRef) == 0 {

			precRef = append(precRef, "-1")
			precRef = append(precRef, "-1")

		} else {

			ref = strings.Split(mzSpec.PrecursorList.Precursor[0].SpectrumRef, " ")
			precRef = strings.Split(ref[2], "=")

			// ABSCIEX has a different way of reporting the prcursor reference spectrum
			if len(ref) < 1 || len(precRef) < 1 {
				precRef = strings.Split(mzSpec.PrecursorList.Precursor[0].SpectrumRef, "=")
			}

		}

		spec.Precursor.ParentScan = strings.TrimSpace(precRef[1])
		pi, _ := strconv.Atoi(precRef[1])
		pi = (pi - 1)
		spec.Precursor.ParentIndex = strconv.Itoa(pi)

		for _, j := range mzSpec.PrecursorList.Precursor[0].IsolationWindow.CVParam {

			if string(j.Accession) == "MS:1000827" {
				val, e := strconv.ParseFloat(j.Value, 64)
				if e != nil {
					msg.CastFloatToString(e, "fatal")
				}
				spec.Precursor.TargetIon = val
			}

			if string(j.Accession) == "MS:1000828" {
				val, e := strconv.ParseFloat(j.Value, 64)
				if e != nil {
					msg.CastFloatToString(e, "fatal")
				}
				spec.Precursor.IsolationWindowLowerOffset = val
			}

			if string(j.Accession) == "MS:1000829" {
				val, e := strconv.ParseFloat(j.Value, 64)
				if e != nil {
					msg.CastFloatToString(e, "fatal")
				}
				spec.Precursor.IsolationWindowUpperOffset = val
			}

		}

		for _, j := range mzSpec.PrecursorList.Precursor[0].SelectedIonList.SelectedIon[0].CVParam {
			if string(j.Accession) == "MS:1000744" {
				val, e := strconv.ParseFloat(j.Value, 64)
				if e != nil {
					msg.CastFloatToString(e, "fatal")
				}
				spec.Precursor.SelectedIon = val
			}

			if string(j.Accession) == "MS:1000041" {
				val, e := strconv.Atoi(j.Value)
				if e != nil {
					msg.CastFloatToString(e, "fatal")
				}
				spec.Precursor.ChargeState = val
			}

			if string(j.Accession) == "MS:1000042" {
				val, e := strconv.ParseFloat(j.Value, 64)
				if e != nil {
					msg.CastFloatToString(e, "fatal")
				}
				spec.Precursor.PeakIntensity = val
			}
		}
	}

	spec.Mz.Stream = mzSpec.BinaryDataArrayList.BinaryDataArray[0].Binary.Value
	for _, j := range mzSpec.BinaryDataArrayList.BinaryDataArray[0].CVParam {
		if string(j.Accession) == "MS:1000523" {
			spec.Mz.Precision = "64"
		} else if string(j.Accession) == "MS:1000521" {
			spec.Mz.Precision = "32"
		}

		if string(j.Accession) == "MS:1000574" {
			spec.Mz.Compression = "1"
		} else if string(j.Accession) == "MS:1000576" {
			spec.Mz.Compression = "0"
		}
	}

	spec.Intensity.Stream = mzSpec.BinaryDataArrayList.BinaryDataArray[1].Binary.Value
	for _, j := range mzSpec.BinaryDataArrayList.BinaryDataArray[1].CVParam {
		if string(j.Accession) == "MS:1000523" {
			spec.Intensity.Precision = "64"
		} else if string(j.Accession) == "MS:1000521" {
			spec.Intensity.Precision = "32"
		}

		if string(j.Accession) == "MS:1000574" {
			spec.Intensity.Compression = "1"
		} else if string(j.Accession) == "MS:1000576" {
			spec.Intensity.Compression = "0"
		}
	}

	if mzSpec.BinaryDataArrayList.Count == 3 {
		spec.IonMobility.Stream = mzSpec.BinaryDataArrayList.BinaryDataArray[2].Binary.Value
		for _, j := range mzSpec.BinaryDataArrayList.BinaryDataArray[2].CVParam {
			if string(j.Accession) == "MS:1000523" {
				spec.IonMobility.Precision = "64"
			} else if string(j.Accession) == "MS:1000521" {
				spec.IonMobility.Precision = "32"
			}

			if string(j.Accession) == "MS:1000574" {
				spec.IonMobility.Compression = "1"
			} else if string(j.Accession) == "MS:1000576" {
				spec.IonMobility.Compression = "0"
			}
		}
	}

	return spec
}

// Decode processes the binary data
func (s *Spectrum) Decode() {

	if len(s.Mz.Stream) > 0 && len(s.Intensity.Stream) > 0 {
		s.Mz.DecodedStream = readEncoded(s.Mz.Stream, s.Mz.Precision, s.Mz.Compression)
		s.Mz.Stream = nil

		s.Intensity.DecodedStream = readEncoded(s.Intensity.Stream, s.Intensity.Precision, s.Intensity.Compression)
		s.Intensity.Stream = nil
	}

	if len(s.IonMobility.Stream) > 0 {
		s.IonMobility.DecodedStream = readEncoded(s.IonMobility.Stream, s.IonMobility.Precision, s.IonMobility.Compression)
		s.IonMobility.Stream = nil
	}

	return
}

// readEncoded transforms the binary data into float64 values
func readEncoded(bin []byte, precision, isCompressed string) []float64 {

	var stream []uint8
	var floatArray []float64

	b := bytes.NewReader(bin)
	b64 := base64.NewDecoder(base64.StdEncoding, b)

	var bytestream bytes.Buffer
	if isCompressed == "1" {
		r, e := zlib.NewReader(b64)
		if e != nil {
			msg.ReadingMzMLZlib(e, "error")
			var emptyArray []float64
			emptyArray = append(emptyArray, 0.0)
			return emptyArray
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
				floatArray = append(floatArray, float64(converted))
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
				floatArray = append(floatArray, float64(converted))
				stream = nil
				counter = 0
			}
		}
	} else {
		logrus.Trace("Error trying to define mzML binary precision")
	}

	return floatArray
}
