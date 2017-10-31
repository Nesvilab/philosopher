package xml

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/cmsl/data/mz/mzml"
	"github.com/prvst/cmsl/err"
	"github.com/prvst/philosopher/lib/sys"
)

// Raw struct
type Raw struct {
	FileName string
	Spectra  Spectra
}

// Spectra is a list of Spetrum
type Spectra []Spectrum

// Spectrum struct
type Spectrum struct {
	Index       string
	Scan        string
	Level       string
	StartTime   float64
	Precursor   Precursor
	Peaks       []float64
	Intensities []float64
}

// Precursor struct
type Precursor struct {
	ParentIndex                string
	ParentScan                 string
	ChargeState                int
	SelectedIon                float64
	TargetIon                  float64
	IsolationWindowLowerOffset float64
	IsolationWindowUpperOffset float64
	PeakIntensity              float64
}

// Len function for Sort
func (s Spectra) Len() int {
	return len(s)
}

// Less function for Sort
func (s Spectra) Less(i, j int) bool {
	return s[i].Index > s[j].Index
}

// Swap function for Sort
func (s Spectra) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *Raw) Read(f, t string) error {

	if t == "mzML" {

		var mz mzml.IndexedMzML
		e := mz.Parse(f)
		if e != nil {
			return &err.Error{Type: err.CannotParseXML, Class: err.FATA, Argument: filepath.Base(f)}
		}

		s.FileName = f
		var spectra Spectra

		for _, i := range mz.MzML.Run.SpectrumList.Spectrum {

			var spec Spectrum
			spec.Index = i.Index

			split := strings.Split(i.ID, " ")
			scan := strings.Split(split[2], "=")

			ind, _ := strconv.Atoi(scan[1])

			adjInd := ind - 1
			spec.Index = strconv.Itoa(adjInd)
			spec.Scan = scan[1]

			for _, j := range i.CVParam {
				if j.Accession == "MS:1000511" {
					spec.Level = j.Value
				}
			}

			for _, j := range i.ScanList.Scan {
				for _, k := range j.CVParam {
					if k.Accession == "MS:1000016" {
						val, e := strconv.ParseFloat(k.Value, 64)
						if e != nil {
							return &err.Error{Type: err.CannotConvertFloatToString, Class: err.FATA, Argument: filepath.Base(f)}
						}
						spec.StartTime = val
					}
				}
			}

			if spec.Level == "2" {
				for _, j := range i.PrecursorList.Precursor {

					var prec Precursor

					split := strings.Split(j.SpectrumRef, " ")
					scan := strings.Split(split[2], "=")

					ind, err := strconv.Atoi(scan[1])
					if err != nil {
						return err
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
						}

						if k.Accession == "MS:1000828" {
							val, err := strconv.ParseFloat(k.Value, 64)
							if err != nil {
								return err
							}
							prec.IsolationWindowLowerOffset = val
						}

						if k.Accession == "MS:1000829" {
							val, err := strconv.ParseFloat(k.Value, 64)
							if err != nil {
								return nil
							}
							prec.IsolationWindowUpperOffset = val
						}

						if k.Accession == "MS:1000042" {
							val, err := strconv.ParseFloat(k.Value, 64)
							if err != nil {
								return nil
							}
							prec.PeakIntensity = val
						}

					}

					for _, k := range j.SelectedIonList.SelectedIon {
						for _, l := range k.CVParam {

							if l.Accession == "MS:1000744" {
								val, err := strconv.ParseFloat(l.Value, 64)
								if err != nil {
									return nil
								}
								prec.SelectedIon = val
							}

							if l.Accession == "MS:1000041" {
								val, err := strconv.Atoi(l.Value)
								if err != nil {
									return err
								}
								prec.ChargeState = val
							}

							if l.Accession == "MS:1000042" {
								val, err := strconv.ParseFloat(l.Value, 64)
								if err != nil {
									return nil
								}
								prec.PeakIntensity = val
							}

						}
					}

					spec.Precursor = prec
				}
			}

			for m := 0; m <= len(i.Peaks)-1; m++ {
				spec.Peaks = append(spec.Peaks, i.Peaks[m])
				spec.Intensities = append(spec.Intensities, i.Intensities[m])
			}

			//spec.Peaks = i.BinaryDataArrayList.BinaryDataArray[0].Binary.Value
			//spec.Intensities = i.BinaryDataArrayList.BinaryDataArray[1].Binary.Value

			// spec.Peaks = i.Peaks
			// spec.Intensities = i.Intensities

			spectra = append(spectra, spec)
		}

		s.Spectra = spectra

	} else if t == "mzXML" {
		fmt.Println("Reader Not implement yet")
		return nil
	}

	return nil
}

// Serialize converts the whle structure to a gob file
func (s *Raw) Serialize() error {

	var err error

	// create a file
	dataFile, err := os.Create(sys.RawBin())
	if err != nil {
		return err
	}

	dataEncoder := gob.NewEncoder(dataFile)
	goberr := dataEncoder.Encode(s)
	if goberr != nil {
		logrus.Fatal("Cannot save results, Bad format", goberr)
	}
	dataFile.Close()

	return nil
}

// Restore reads philosopher results files and restore the data sctructure
func (s *Raw) Restore() error {

	file, _ := os.Open(sys.RawBin())

	dec := gob.NewDecoder(file)
	err := dec.Decode(&s)
	if err != nil {
		return errors.New("Could not restore Philosopher result. Please check file path")
	}

	return nil
}
