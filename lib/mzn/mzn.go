package mzn

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"io"
	"math"
	"regexp"
	"strconv"
	"strings"

	"philosopher/lib/msg"

	"philosopher/lib/psi"

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
	Index               string
	Scan                string
	Level               string
	SpectrumName        string
	CompensationVoltage string
	ScanStartTime       float64
	Precursor           Precursor
	Mz                  Mz
	Intensity           Intensity
	IonMobility         IonMobility
}

// Precursor struct
type Precursor struct {
	ParentIndex                string
	ParentScan                 string
	ChargeState                int
	SelectedIon                float64
	SelectedIonIntensity       float64
	TargetIon                  float64
	TargetIonIntensity         float64
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

// ReadRaw is the main function for parsing Thermo Raw data
func (p *MsData) ReadRaw(fileName, f string) {

	p.FileName = fileName
	var spectra Spectra

	lines := strings.Split(f, "%")

	for _, i := range lines {
		parts := strings.Split(i, "#")

		if len(parts) >= 8 {

			var spec Spectrum

			if strings.Contains(parts[0], "ms3") {
				spec.Level = "3"
			} else if strings.Contains(parts[0], "ms2") {
				spec.Level = "2"
			} else {
				spec.Level = "1"
			}

			spec.Scan = parts[1]
			indexInt, _ := strconv.Atoi(parts[1])
			indexInt--
			spec.Index = string(strconv.Itoa(indexInt))

			parts2, _ := strconv.Atoi(parts[2])
			spec.Precursor.ChargeState = parts2

			spec.Precursor.ParentScan = parts[3]
			parentIndexInt, _ := strconv.Atoi(parts[3])
			parentIndexInt--
			spec.Precursor.ParentIndex = string(strconv.Itoa(parentIndexInt))

			spec.Precursor.IsolationWindowLowerOffset = 0.6
			spec.Precursor.IsolationWindowUpperOffset = 0.6

			siVal1, e := strconv.ParseFloat(parts[5], 64)
			if e != nil {
				msg.CastFloatToString(e, "fatal")
			} else {
				spec.Precursor.SelectedIon = siVal1
				spec.Precursor.TargetIon = siVal1
			}

			siVal2, e := strconv.ParseFloat(parts[6], 64)
			if e != nil {
				msg.CastFloatToString(e, "fatal")
			} else {
				//spec.Precursor.SelectedIonIntensity = siVal2
				spec.Precursor.TargetIonIntensity = siVal2
			}

			if spec.Level == "1" {
				spec.Precursor.ParentScan = ""
				spec.Precursor.ParentIndex = ""
				spec.Precursor.ChargeState = 0
				spec.Precursor.SelectedIon = 0
				spec.Precursor.SelectedIonIntensity = 0
				spec.Precursor.TargetIon = 0
				spec.Precursor.TargetIonIntensity = 0
				spec.Precursor.IsolationWindowLowerOffset = 0.6
				spec.Precursor.IsolationWindowUpperOffset = 0.6
			}

			rtVal, _ := strconv.ParseFloat(parts[4], 64)
			spec.ScanStartTime = rtVal

			//spec.Mz.Stream = []byte(parts[7])
			spec.Mz.Precision = "64"
			spec.Mz.Compression = "1"

			//spec.Intensity.Stream = []byte(parts[8])
			spec.Intensity.Precision = "64"
			spec.Intensity.Compression = "1"

			peaks := strings.Split(parts[8], " ")
			for _, arg := range peaks {
				if n, e := strconv.ParseFloat(string(arg), 64); e == nil {
					spec.Mz.DecodedStream = append(spec.Mz.DecodedStream, n)
				}
			}

			intensities := strings.Split(parts[9], " ")
			for _, arg := range intensities {
				if n, e := strconv.ParseFloat(string(arg), 64); e == nil {
					spec.Intensity.DecodedStream = append(spec.Intensity.DecodedStream, n)
				}
			}

			spectra = append(spectra, spec)
		}
	}

	if len(spectra) == 0 {
		msg.NoSpectraFound(errors.New(""), "fatal")
	}

	p.Spectra = spectra
}

// Read is the main function for parsing mzML data
func (p *MsData) Read(f string) {

	var xml psi.IndexedMzML
	xml.Parse(f)

	if xml.MzML.SoftwareList.Software[0].ID == "pwiz" {
		version, _ := strconv.Atoi(strings.Replace(xml.MzML.SoftwareList.Software[0].Version, ".", "", -1))
		if version <= 3019127 {
			msg.Custom(errors.New("the pwiz version used to convert this file is not supported, or deprecated. Please update your pwiz and convert the raw files again"), "warning")
		}
	}

	p.FileName = f

	var spectra Spectra
	sl := xml.MzML.Run.SpectrumList

	for _, i := range sl.Spectrum {

		spectrum := processSpectrum(i)

		spectra = append(spectra, spectrum)
	}

	if len(spectra) == 0 {
		msg.NoSpectraFound(errors.New(""), "fatal")
	}

	p.Spectra = spectra

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

		if string(j.Accession) == "MS:1001581" {
			spec.CompensationVoltage = j.Value
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

		if len(mzSpec.PrecursorList.Precursor[0].SpectrumRef) > 0 {

			scanRG := regexp.MustCompile(`scan=(.+)`)
			match := scanRG.FindStringSubmatch(mzSpec.PrecursorList.Precursor[0].SpectrumRef)

			spec.Precursor.ParentScan = match[1]
			pi, _ := strconv.Atoi(match[1])
			pi = (pi - 1)
			spec.Precursor.ParentIndex = strconv.Itoa(pi)
		}

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
				spec.Precursor.SelectedIonIntensity = val
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
