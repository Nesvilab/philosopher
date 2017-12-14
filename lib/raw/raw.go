package raw

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/raw/mz"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/vmihailenco/msgpack"
)

// Data represents parsed and processed MZ data from mz files
type Data struct {
	Raw *mz.Raw
}

// IndexMz receives a list of mz files and creates a binary index for each one
func IndexMz(f []string) *err.Error {

	for _, i := range f {

		var d Data

		if strings.Contains(i, "mzml") || strings.Contains(i, "mzML") {

			logrus.Info("Indexing ", i)

			raw := &mz.Raw{}
			raw.FileName = i

			// e := raw.Read(i)
			// if e != nil {
			// 	return e
			// }

			e := raw.ParRead(i)
			if e != nil {
				return e
			}

			s := make(map[interface{}]interface{})
			raw.RefSpectra.Range(func(k, v interface{}) bool {
				s[k] = v
				return true
			})

			for _, v := range s {

				var spectrum mz.Spectrum
				spectrum.Index = v.(mz.Spectrum).Index
				spectrum.Level = v.(mz.Spectrum).Level
				spectrum.Intensities = v.(mz.Spectrum).Intensities
				spectrum.Peaks = v.(mz.Spectrum).Peaks
				spectrum.Precursor = v.(mz.Spectrum).Precursor
				spectrum.Scan = v.(mz.Spectrum).Scan
				spectrum.StartTime = v.(mz.Spectrum).StartTime

				raw.Spectra = append(raw.Spectra, spectrum)
			}
			s = nil

			d.Raw = raw
			raw = nil

		} else if strings.Contains(i, "mzxml") || strings.Contains(i, "mzXML") {
			return &err.Error{Type: err.MethodNotImplemented, Class: err.FATA, Argument: "mzXML reader not implemented"}
		}

		d.Serialize()
		d = Data{}
	}

	return nil
}

// Serialize mz data structure to binary format
func (data *Data) Serialize() *err.Error {

	// remove the extension
	var extension = filepath.Ext(filepath.Base(data.Raw.FileName))
	var name = data.Raw.FileName[0 : len(data.Raw.FileName)-len(extension)]

	// overwrite raw file
	data.Raw.FileName = filepath.Base(name)

	output := fmt.Sprintf("%s%s%s.bin", sys.MetaDir(), string(filepath.Separator), filepath.Base(name))

	b, err := msgpack.Marshal(&data)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(output)
	if err != nil {
		return nil
	}
	defer f.Close()

	f.Write(b)
	_ = b

	return nil
}

// Restore reads philosopher results files and restore the data sctructure
func Restore(f string) (*Data, *err.Error) {

	var data Data

	// remove the extension
	var extension = filepath.Ext(filepath.Base(f))
	var name = f[0 : len(f)-len(extension)]

	input := fmt.Sprintf("%s%s%s.bin", sys.MetaDir(), string(filepath.Separator), name)

	file, _ := os.Open(input)

	dec := msgpack.NewDecoder(file)
	e := dec.Decode(&data)
	if e != nil {
		return &data, &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: e.Error()}
	}

	return &data, nil
}

// RestoreFromFile reads the mz information directly from a mz file, not from indxed binaries
func RestoreFromFile(dir, f, format string) (*Data, *err.Error) {

	var d Data

	fullPath := fmt.Sprintf("%s%s%s.%s", dir, string(filepath.Separator), f, format)

	if format == "mzML" {

		raw := &mz.Raw{}
		raw.FileName = f

		e := raw.ParRead(fullPath)
		if e != nil {
			return &d, e
		}

		s := make(map[interface{}]interface{})
		raw.RefSpectra.Range(func(k, v interface{}) bool {
			s[k] = v
			return true
		})

		for _, v := range s {

			var spectrum mz.Spectrum
			spectrum.Index = v.(mz.Spectrum).Index
			spectrum.Level = v.(mz.Spectrum).Level
			spectrum.Intensities = v.(mz.Spectrum).Intensities
			spectrum.Peaks = v.(mz.Spectrum).Peaks
			spectrum.Precursor = v.(mz.Spectrum).Precursor
			spectrum.Scan = v.(mz.Spectrum).Scan
			spectrum.StartTime = v.(mz.Spectrum).StartTime

			raw.Spectra = append(raw.Spectra, spectrum)
		}

		d.Raw = raw

		s = nil
		raw = nil

	} else if strings.Contains(f, "mzxml") || strings.Contains(f, "mzXML") {
		return &d, &err.Error{Type: err.MethodNotImplemented, Class: err.FATA, Argument: "mzXML reader not implemented"}
	}

	return &d, nil
}

// GetMS1 from Spectral Data
func GetMS1(d *Data) MS1 {

	var list MS1

	for _, i := range d.Raw.Spectra {
		if string(i.Level) == "1" {

			var scan Ms1Scan

			scan.Index = i.Index
			scan.Scan = i.Scan

			time := i.StartTime
			scan.ScanStartTime = time

			var stream Spectrum
			for m := 0; m <= len(i.Peaks.DecodedStream)-1; m++ {
				var peak Peak
				peak.Mz = i.Peaks.DecodedStream[m]
				peak.Intensity = i.Intensities.DecodedStream[m]
				stream = append(stream, peak)
			}

			scan.Spectrum = stream
			list.Ms1Scan = append(list.Ms1Scan, scan)
		}
	}

	return list
}

// GetMS2 from Spectral Data
func GetMS2(d *Data) MS2 {

	var list MS2

	for _, i := range d.Raw.Spectra {
		if string(i.Level) == "2" {

			var scan Ms2Scan

			scan.Index = i.Index
			scan.Scan = i.Scan

			time := i.StartTime
			scan.ScanStartTime = time

			scan.Precursor.ChargeState = i.Precursor.ChargeState
			scan.Precursor.IsolationWindowLowerOffset = i.Precursor.IsolationWindowLowerOffset
			scan.Precursor.IsolationWindowUpperOffset = i.Precursor.IsolationWindowUpperOffset
			scan.Precursor.ParentIndex = i.Precursor.ParentIndex
			scan.Precursor.ParentScan = i.Precursor.ParentScan
			scan.Precursor.PeakIntensity = i.Precursor.PeakIntensity
			scan.Precursor.SelectedIon = i.Precursor.SelectedIon
			scan.Precursor.TargetIon = i.Precursor.TargetIon

			var stream Spectrum
			for m := 0; m <= len(i.Peaks.DecodedStream)-1; m++ {
				var peak Peak
				peak.Mz = i.Peaks.DecodedStream[m]
				peak.Intensity = i.Intensities.DecodedStream[m]
				stream = append(stream, peak)
			}

			scan.Spectrum = stream
			list.Ms2Scan = append(list.Ms2Scan, scan)

		}
	}

	return list
}
