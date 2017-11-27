package mzxml

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher/lib/fin"
	"github.com/rogpeppe/go-charset/charset"
	// anon charset
	_ "github.com/rogpeppe/go-charset/data"
)

// MzXML tag
type MzXML struct {
	XMLName xml.Name `xml:"mzXML"`
	Xmlns   []byte   `xml:"xmlns,attr"`
	MSRun   MSRun    `xml:"msRun"`
}

// MSRun tag
type MSRun struct {
	XMLName        xml.Name       `xml:"msRun"`
	ScanCount      uint64         `xml:"scanCount,attr"`
	StartTime      []byte         `xml:"startTime,attr"`
	EndTime        []byte         `xml:"endTime,attr"`
	ParentFile     ParentFile     `xml:"parentFile"`
	MSInstrument   MSInstrument   `xml:"msInstrument"`
	DataProcessing DataProcessing `xml:"dataProcessing"`
	Scan           []Scan         `xml:"scan"`
}

// ParentFile tag
type ParentFile struct {
	XMLName  xml.Name `xml:"parentFile"`
	FileName []byte   `xml:"fileName,attr"`
	FileType []byte   `xml:"fileType,attr"`
	FileSha1 []byte   `xml:"fileSha1,attr"`
}

// MSInstrument tag
type MSInstrument struct {
	XMLName        xml.Name       `xml:"msInstrument"`
	MSInstrumentID uint8          `xml:"msInstrumentID,attr"`
	MSManufacturer MSManufacturer `xml:"msManufacturer"`
	MSModel        MSModel        `xml:"msModel"`
	MSIonisation   MSIonisation   `xml:"msIonisation"`
	MSMassAnalyzer MSMassAnalyzer `xml:"msMassAnalyzer"`
	MSDetector     MSDetector     `xml:"msDetector"`
	Software       MSSoftware     `xml:"software"`
}

//MSManufacturer tag
type MSManufacturer struct {
	XMLName  xml.Name `xml:"msManufacturer"`
	Category []byte   `xml:"category,attr"`
	Value    []byte   `xml:"value,attr"`
}

//MSModel tag
type MSModel struct {
	XMLName  xml.Name `xml:"msModel"`
	Category []byte   `xml:"category,attr"`
	Value    []byte   `xml:"value,attr"`
}

//MSIonisation tag
type MSIonisation struct {
	XMLName  xml.Name `xml:"msIonisation"`
	Category []byte   `xml:"category,attr"`
	Value    []byte   `xml:"value,attr"`
}

//MSMassAnalyzer tag
type MSMassAnalyzer struct {
	XMLName  xml.Name `xml:"msMassAnalyzer"`
	Category []byte   `xml:"category,attr"`
	Value    []byte   `xml:"value,attr"`
}

//MSDetector tag
type MSDetector struct {
	XMLName  xml.Name `xml:"msDetector"`
	Category []byte   `xml:"category,attr"`
	Value    []byte   `xml:"value,attr"`
}

//MSSoftware tag
type MSSoftware struct {
	XMLName xml.Name `xml:"software"`
	Type    []byte   `xml:"type,attr"`
	Name    []byte   `xml:"name,attr"`
	Version []byte   `xml:"version,attr"`
}

//DataProcessing tag
type DataProcessing struct {
	XMLName             xml.Name            `xml:"dataProcessing"`
	Centroided          uint8               `xml:"centroided,attr"`
	Software            MSSoftware          `xml:"software"`
	ProcessingOperation ProcessingOperation `xml:"processingOperation"`
}

//ProcessingOperation tag
type ProcessingOperation struct {
	XMLName xml.Name `xml:"processingOperation"`
	Name    []byte   `xml:"name,attr"`
}

// Scan tag
type Scan struct {
	XMLName           xml.Name    `xml:"scan"`
	Num               int         `xml:"num,attr"`
	ScanType          []byte      `xml:"scanType,attr"`
	MSLevel           uint8       `xml:"msLevel,attr"`
	Centroided        []byte      `xml:"centroided,attr"`
	PeaksCount        int         `xml:"peaksCount,attr"`
	Polarity          []byte      `xml:"polarity,attr"`
	RetentionTime     string      `xml:"retentionTime,attr"`
	CollisionEnergy   float64     `xml:"collisionEnergy,attr"`
	LowMz             float64     `xml:"lowMz,attr"`
	HighMz            float64     `xml:"highMz,attr"`
	BasePeakMz        float64     `xml:"basePeakMz,attr"`
	BasePeakIntensity float64     `xml:"basePeakIntensity,attr"`
	TotIonCurrent     float64     `xml:"totIonCurrent,attr"`
	Precursor         PrecursorMz `xml:"precursorMz"`
	Peaks             Peaks       `xml:"peaks"`
	ConvertedPeaks    []float64
}

// PrecursorMz ...
type PrecursorMz struct {
	XMLName            xml.Name `xml:"precursorMz"`
	PrecursorScanNum   int      `xml:"precursorScanNum,attr"`
	PrecursorIntensity float64  `xml:"precursorIntensity,attr"`
	PrecursorCharge    int      `xml:"precursorCharge,attr"`
	ActivationMethod   []byte   `xml:"activationMethod,attr"`
	WindowWideness     float64  `xml:"windowWideness,attr"`
	Value              []byte   `xml:",chardata"`
}

// Peaks tag
type Peaks struct {
	XMLName         xml.Name `xml:"peaks"`
	CompressionType []byte   `xml:"compressionType,attr"`
	CompressionLen  uint32   `xml:"compressionLen,attr"`
	Precision       []byte   `xml:"precision,attr"`
	ByteOrder       []byte   `xml:"byteOrder,attr"`
	ContentType     []byte   `xml:"contentType,attr"`
	Value           []byte   `xml:",chardata"`
}

// Parse is the main function for parsing data
func (m *MzXML) Parse(f string) error {

	xmlFile, err := os.Open(f)
	if err != nil {
		return errors.New("Error trying to read mzXML file")
	}
	defer xmlFile.Close()
	b, _ := ioutil.ReadAll(xmlFile)

	var mz MzXML

	reader := bytes.NewReader(b)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReader

	if err = decoder.Decode(&mz); err != nil {
		msg := fmt.Sprintf("Unable to parse XML: %s", err)
		return errors.New(msg)
	}

	for i := range mz.MSRun.Scan {
		mz.MSRun.Scan[i].ConvertedPeaks, err = Decode(mz.MSRun.Scan[i].Peaks)
		if err != nil {
			return err
		}

	}

	m.MSRun = mz.MSRun

	return nil
}

// Decode processes the binary data
func Decode(peaks Peaks) ([]float64, error) {

	var compression bool
	var precision string
	var err error

	if string(peaks.Precision) == "64" {
		precision = "64"
	} else if string(peaks.Precision) == "32" {
		precision = "32"
	}

	if string(peaks.CompressionType) == "zlib" {
		compression = true
	} else {
		compression = false
	}

	f, err := readEncoded(peaks, precision, compression)
	if err != nil {
		return f, err
	}

	return f, nil
}

// readEncoded transforms the binary data into float64 values
func readEncoded(peaks Peaks, precision string, isCompressed bool) ([]float64, error) {

	var stream []uint8
	var floatArray []float64

	b := bytes.NewReader(peaks.Value)
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
				bits := binary.BigEndian.Uint32(stream)
				converted := math.Float32frombits(bits)
				//floatArray = append(floatArray, uti.Round(float64(converted), 5, 2))
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
				bits := binary.BigEndian.Uint64(stream)
				converted := math.Float64frombits(bits)
				//floatArray = append(floatArray, uti.Round(converted, 5, 2))
				floatArray = append(floatArray, converted)
				stream = nil
				counter = 0
			}
		}
	} else {
		return floatArray, errors.New("Undefined binary precision")
	}

	return floatArray, nil
}

// // ReadEncoded function converts the base64 encrypt info from mzxml files
// // and converts to a float64 array
// func ReadEncoded(peaks Peaks) []float64 {
//
// 	var dataArray []uint8
// 	var floatArray []float64
// 	var stream []uint8
//
// 	if strings.EqualFold(string(peaks.CompressionType), "zlib") {
//
// 		b := bytes.NewReader(peaks.Value)
// 		b64 := base64.NewDecoder(base64.StdEncoding, b)
// 		r, err := zlib.NewReader(b64)
//
// 		if err != nil {
// 			panic(err)
// 		}
//
// 		var bytestream bytes.Buffer
// 		io.Copy(&bytestream, r)
// 		dataArray = bytestream.Bytes()
//
// 	} else {
// 		if strings.EqualFold(string(peaks.Precision), "32") {
// 			dataArray, _ = base32.StdEncoding.DecodeString(string(peaks.Value))
// 		} else {
// 			dataArray, _ = base64.StdEncoding.DecodeString(string(peaks.Value))
// 		}
// 	}
//
// 	var counter int
// 	for i := range dataArray {
// 		counter++
// 		stream = append(stream, dataArray[i])
// 		if counter == 8 {
// 			bits := binary.BigEndian.Uint64(stream)
// 			converted := math.Float64frombits(bits)
// 			floatArray = append(floatArray, converted)
// 			stream = nil
// 			counter = 0
// 		}
// 		//i++
// 	}
//
// 	return floatArray
// }

// Write ...
func Write(raw fin.RawData) error {

	var mz MzXML

	mz.Xmlns = []byte("http://sashimi.sourceforge.net/schema_revision/mzXML_3.1")
	mz.MSRun.ScanCount = raw.ScanCount

	startTime := fmt.Sprintf("PT%.4fS", 60*raw.StartTime)
	endTime := fmt.Sprintf("PT%.4fS", 60*raw.EndTime)
	mz.MSRun.StartTime = []byte(startTime)
	mz.MSRun.EndTime = []byte(endTime)

	absFile, _ := filepath.Abs(raw.FileName)
	mz.MSRun.ParentFile.FileName = []byte(absFile)
	mz.MSRun.ParentFile.FileType = []byte("RAWData")

	sha1, err := calcSha1(absFile)
	if err != nil {
		return err
	}
	mz.MSRun.ParentFile.FileSha1 = []byte(sha1)

	mz.MSRun.MSInstrument.MSInstrumentID = 1

	mz.MSRun.MSInstrument.MSManufacturer.Category = []byte("msManufacturer")
	mz.MSRun.MSInstrument.MSManufacturer.Value = []byte(raw.Manufacturer)

	mz.MSRun.MSInstrument.MSModel.Category = []byte("msModel")
	mz.MSRun.MSInstrument.MSModel.Value = []byte(raw.Model)

	mz.MSRun.MSInstrument.MSIonisation.Category = []byte("msIonisation")
	mz.MSRun.MSInstrument.MSIonisation.Value = []byte(raw.Ionization[0])

	mz.MSRun.MSInstrument.MSMassAnalyzer.Category = []byte("msMassAnalyzer")
	mz.MSRun.MSInstrument.MSMassAnalyzer.Value = []byte(raw.Analyzer[0])

	mz.MSRun.MSInstrument.MSDetector.Category = []byte("msDetector")
	mz.MSRun.MSInstrument.MSDetector.Value = []byte(raw.Detector[0])

	mz.MSRun.MSInstrument.Software.Type = []byte("aquisition")
	mz.MSRun.MSInstrument.Software.Name = []byte(raw.SoftwareName)
	mz.MSRun.MSInstrument.Software.Version = []byte(raw.SoftwareVersion)

	mz.MSRun.DataProcessing.Software.Type = []byte("conversion")
	mz.MSRun.DataProcessing.Software.Name = []byte("OpenConverter")
	mz.MSRun.DataProcessing.Software.Version = []byte("0.1")
	mz.MSRun.DataProcessing.ProcessingOperation.Name = []byte("Conversion to mzXML")

	for i := 1; i <= raw.NScans(); i++ {
		var scan Scan

		rawScan := raw.Scan(i)
		scan.Num = i
		scan.ScanType = rawScan.Type
		scan.MSLevel = rawScan.MSLevel
		scan.Centroided = rawScan.Mode
		scan.PeaksCount = len(rawScan.Spectrum())
		scan.Polarity = rawScan.Polarity
		scan.RetentionTime = fmt.Sprintf("PT%.4fS", 60*rawScan.Time)
		scan.LowMz = rawScan.LowMz
		scan.HighMz = rawScan.HighMz
		scan.BasePeakMz = rawScan.BaseMz
		scan.BasePeakIntensity = rawScan.BaseIntensity
		scan.TotIonCurrent = rawScan.TotalCurrent

		if scan.MSLevel > 1 {
			var p PrecursorMz
			scan.CollisionEnergy = rawScan.Fragment[0].ColisionEnergy[0]
			stringMz := strconv.FormatFloat(rawScan.Fragment[0].PrecursorMzs[0], 'E', -1, 64)
			p.Value = []byte(stringMz)
		}

		for _, peak := range rawScan.Spectrum() {
			var peaks Peaks
			var value string

			peaks.Precision = []byte("64")
			peaks.ByteOrder = []byte("network")
			peaks.ContentType = []byte("m/z-int")

			value = fmt.Sprintf("%f %f ", peak.Mz, peak.I)
			byteValue := []byte(value)
			peaks.Value = append(scan.Peaks.Value, byteValue...)
			scan.Peaks = peaks
		}

		// base64 encoding
		sEnc := base64.StdEncoding.EncodeToString(scan.Peaks.Value)
		scan.Peaks.Value = []byte(sEnc)

		// zlib
		var zip = false
		if zip == true {
			var b bytes.Buffer
			w := zlib.NewWriter(&b)
			w.Write(scan.Peaks.Value)
			w.Close()
		}

		mz.MSRun.Scan = append(mz.MSRun.Scan, scan)

	}

	marh, err := xml.MarshalIndent(mz, " ", "  ")
	if err != nil {
		logrus.Fatal(err)
	}

	// to file
	stream := []byte(xml.Header + string(marh))
	err = ioutil.WriteFile("test.mzXML", stream, 0644)
	if err != nil {
		return errors.New("Cannot write output file")
	}

	return nil
}

func calcSha1(f string) (string, error) {

	hasher := sha256.New()
	s, err := ioutil.ReadFile(f)
	hasher.Write(s)
	if err != nil {
		return hex.EncodeToString(hasher.Sum(nil)), err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
