package xml

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/cmsl/data/mz/mzml"
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
	Index       int
	Scan        int
	Level       int
	StartTime   float64
	EndTime     float64
	Precursor   Precursor
	Peaks       []byte
	Intensities []byte
	//Peaks       []float64
	//Intensities []float64
}

// Precursor struct
type Precursor struct {
	ParentIndex   int
	ParentScan    int
	SelectedIon   float64
	ChargeState   int
	PeakIntensity float64
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
		err := mz.Parse(f)
		if err != nil {
			return err
		}

		s.FileName = f
		var spectra Spectra

		for _, i := range mz.MzML.Run.SpectrumList.Spectrum {

			var spec Spectrum
			spec.Index = i.Index
			spec.Scan = i.Index + 1

			for _, j := range i.CVParam {
				if j.Accession == "MS:1000579" {
					spec.Level = 1
				} else if j.Accession == "MS:1000511" {
					spec.Level = 2
				}
			}

			for _, j := range i.ScanList.Scan {
				for _, k := range j.CVParam {
					if k.Accession == "MS:1000016" {
						val, err := strconv.ParseFloat(k.Value, 64)
						if err != nil {
							return nil
						}
						spec.StartTime = val
					}
				}
			}

			for _, j := range i.PrecursorList.Precursor {

				var prec Precursor
				split := strings.Split(j.SpectrumRef, " ")
				scan := strings.Split(split[2], "=")

				ind, err := strconv.Atoi(scan[1])
				if err != nil {
					return err
				}

				prec.ParentIndex = ind - 1
				prec.ParentScan = ind

				for _, k := range j.IsolationWindow.CVParam {
					if k.Accession == "MS:1000744" {
						val, err := strconv.ParseFloat(k.Value, 64)
						if err != nil {
							return nil
						}
						prec.SelectedIon = val
					} else if k.Accession == "MS:1000041" {
						val, err := strconv.Atoi(k.Value)
						if err != nil {
							return err
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

				spec.Precursor = prec
			}

			spec.Peaks = i.BinaryDataArrayList.BinaryDataArray[0].Binary.Value
			spec.Intensities = i.BinaryDataArrayList.BinaryDataArray[1].Binary.Value

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
