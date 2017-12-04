package raw

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/mzm"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/prvst/philosopher/lib/uti"
	"github.com/vmihailenco/msgpack"
)

// Spectra is a collection of MS1 and MS2
type Spectra struct {
	FileName string
	MS1      []Ms1Scan
	MS2      map[string]Ms2Scan
}

// Ms1Scan tag
type Ms1Scan struct {
	Scan          string
	Index         string
	SpectrumName  string
	ScanStartTime float64
	Spectrum      Stream
}

// Ms2Scan tag
type Ms2Scan struct {
	Scan          string
	Index         string
	SpectrumName  string
	ScanStartTime float64
	Precursor     Precursor
	Spectrum      Stream
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

// Stream tag
type Stream []Peak

// Peak tag
type Peak struct {
	Mz        float64
	Intensity float64
}

func (a Stream) Len() int           { return len(a) }
func (a Stream) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Stream) Less(i, j int) bool { return a[i].Mz < a[j].Mz }

// ReadMzML builds the spectra MS1 and MS2 indexes
func (m *Spectra) ReadMzML(f string) *err.Error {

	xmlFile, e := os.Open(f)

	if e != nil {
		return &err.Error{Type: err.CannotParseXML, Class: err.FATA, Argument: "mzML"}
	}
	defer xmlFile.Close()

	// get the clean file name for spectra name
	ext := filepath.Ext(f)
	name := filepath.Base(f)
	cleanFileName := name[0 : len(name)-len(ext)]

	var mz mzm.IndexedMzML
	e = mz.Parse(f)
	if e != nil {
		return &err.Error{Type: err.CannotParseXML, Class: err.FATA, Argument: "mzML"}
	}

	for _, i := range mz.MzML.Run.SpectrumList.Spectrum {
		for _, j := range i.CVParam {

			if strings.EqualFold(j.Name, "ms level") && strings.EqualFold(j.Value, "1") {

				var ms1 Ms1Scan
				ms1.Index = string(i.Index)

				// parse and format scan ID for the SpectrumName
				idSplit := strings.Split(i.ID, "scan=")
				iID, _ := strconv.Atoi(idSplit[1])
				cusID := fmt.Sprintf("%05d", iID)
				ms1.SpectrumName = fmt.Sprintf("%s.%s.%s", cleanFileName, cusID, cusID)

				for _, k := range i.ScanList.Scan {
					for _, l := range k.CVParam {
						if strings.EqualFold(l.Name, "scan start time") {
							stt, _ := strconv.ParseFloat(l.Value, 64)
							ms1.ScanStartTime = uti.Round(stt, 5, 4)
						}
					}
				}

				var stream Stream

				for m := 0; m <= len(i.Peaks)-1; m++ {
					var peak Peak
					peak.Mz = i.Peaks[m]
					peak.Intensity = i.Intensities[m]
					stream = append(stream, peak)
				}

				ms1.Spectrum = stream
				m.MS1 = append(m.MS1, ms1)

			} else if j.Name == "ms level" && j.Value == "2" {

				var ms2 Ms2Scan
				ms2.Index = i.Index

				// parse and format scan ID for the SpectrumName
				idSplit := strings.Split(i.ID, "scan=")
				iID, _ := strconv.Atoi(idSplit[1])
				cusID := fmt.Sprintf("%05d", iID)
				ms2.SpectrumName = fmt.Sprintf("%s.%s.%s", cleanFileName, cusID, cusID)
				ms2.Scan = cusID

				for _, k := range i.ScanList.Scan {
					for _, l := range k.CVParam {
						if l.Name == "scan start time" {
							stt, _ := strconv.ParseFloat(l.Value, 64)
							minrt := (stt * 60)
							ms2.ScanStartTime = uti.ToFixed(minrt, 3)
						}
					}
				}

				for _, j := range i.PrecursorList.Precursor {

					var prec Precursor

					split := strings.Split(j.SpectrumRef, " ")
					scan := strings.Split(split[2], "=")

					ind, e := strconv.Atoi(scan[1])
					if e != nil {
						return nil
					}

					adjInd := ind - 1
					prec.ParentIndex = strconv.Itoa(adjInd)
					prec.ParentScan = scan[1]

					for _, k := range j.IsolationWindow.CVParam {
						if k.Accession == "MS:1000827" {
							val, err := strconv.ParseFloat(k.Value, 64)
							if err != nil {
								return nil
							}
							prec.TargetIon = val

						} else if k.Accession == "MS:1000828" {
							val, err := strconv.ParseFloat(k.Value, 64)
							if err != nil {
								return nil
							}
							prec.IsolationWindowLowerOffset = val

						} else if k.Accession == "MS:1000829" {
							val, err := strconv.ParseFloat(k.Value, 64)
							if err != nil {
								return nil
							}
							prec.IsolationWindowUpperOffset = val

						}
					}

					for _, k := range j.SelectedIonList.SelectedIon[0].CVParam {
						if k.Accession == "MS:1000744" {
							val, err := strconv.ParseFloat(k.Value, 64)
							if err != nil {
								return nil
							}
							prec.SelectedIon = val

						} else if k.Accession == "MS:1000041" {
							val, err := strconv.Atoi(k.Value)
							if err != nil {
								return nil
							}
							prec.ChargeState = val

						} else if k.Accession == "MS:1000042" {
							val, err := strconv.ParseFloat(k.Value, 64)
							if err != nil {
								return nil
							}
							prec.PeakIntensity = val

						}
					}

					ms2.Precursor = prec
				}

				mzPeaks, err := mzm.Decode("mz", i.BinaryDataArrayList.BinaryDataArray[0])
				if err != nil {
					return nil
				}

				mzIntensities, err := mzm.Decode("int", i.BinaryDataArrayList.BinaryDataArray[1])
				if err != nil {
					return nil
				}

				var stream Stream

				for m := range mzPeaks {
					var peak Peak
					peak.Mz = mzPeaks[m]
					peak.Intensity = mzIntensities[m]
					stream = append(stream, peak)
				}

				ms2.Spectrum = stream
				var fullSpecName = fmt.Sprintf("%s.%s.%s.%d", cleanFileName, cusID, cusID, ms2.Precursor.ChargeState)
				m.MS2[fullSpecName] = ms2

			}

		}
	}

	return nil
}

// Index receives a list of mz files and creates a binary index for each one
func Index(f []string) *err.Error {

	for _, i := range f {

		var d Spectra

		if strings.Contains(i, "mzml") || strings.Contains(i, "mzML") {

			e := d.ReadMzML(i)
			if e != nil {
				return e
			}

			d.FileName = i

		} else if strings.Contains(i, "mzxml") || strings.Contains(i, "mzXML") {
			return &err.Error{Type: err.MethodNotImplemented, Class: err.FATA, Argument: "mzXML reader not implemented"}
		}

		d.Serialize()
	}

	return nil
}

// Serialize mz data structure to binary format
func (m *Spectra) Serialize() *err.Error {

	// remove the extension
	var extension = filepath.Ext(filepath.Base(m.FileName))
	var name = m.FileName[0 : len(m.FileName)-len(extension)]

	output := fmt.Sprintf("%s%s%s.bin", sys.MetaDir(), string(filepath.Separator), filepath.Base(name))

	// create a file
	dataFile, e := os.Create(output)
	if e != nil {
		return &err.Error{Type: err.CannotCreateOutputFile, Class: err.FATA, Argument: e.Error()}
	}

	dataEncoder := msgpack.NewEncoder(dataFile)
	goberr := dataEncoder.Encode(m)
	if goberr != nil {
		return &err.Error{Type: err.CannotSerializeData, Class: err.FATA, Argument: e.Error()}
	}
	dataFile.Close()

	return nil
}

// Restore reads philosopher results files and restore the data sctructure
func (m *Spectra) Restore() *err.Error {

	file, _ := os.Open(sys.RawBin())

	dec := msgpack.NewDecoder(file)
	e := dec.Decode(&m)
	if e != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: e.Error()}
	}

	return nil
}
