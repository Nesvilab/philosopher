package mz

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/mzm"
	"github.com/prvst/philosopher/lib/mzx"
	"github.com/prvst/philosopher/lib/uti"
)

// MS1 top struct
type MS1 struct {
	Ms1Scan []Ms1Scan
}

// Ms1Scan tag
type Ms1Scan struct {
	Index         string
	SpectrumName  string
	ScanStartTime float64
	Spectrum      Ms1Spectrum
}

// Ms1Spectrum tag
type Ms1Spectrum []Ms1Peak

// Ms1Peak tag
type Ms1Peak struct {
	Mz        float64
	Intensity float64
}

func (a Ms1Spectrum) Len() int           { return len(a) }
func (a Ms1Spectrum) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Ms1Spectrum) Less(i, j int) bool { return a[i].Mz < a[j].Mz }

// ReadMzML parses only MS1 data from mzML
func (m *MS1) ReadMzML(f string) *err.Error {

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

				var ms1Peaks Ms1Spectrum

				for m := 0; m <= len(i.Peaks)-1; m++ {
					var peak Ms1Peak
					peak.Mz = i.Peaks[m]
					peak.Intensity = i.Intensities[m]
					ms1Peaks = append(ms1Peaks, peak)
				}

				ms1.Spectrum = ms1Peaks
				m.Ms1Scan = append(m.Ms1Scan, ms1)

			}
		}
	}

	return nil
}

// GetMzMLSpectra parses only MS1 data from mzML
func GetMzMLSpectra(mz mzm.IndexedMzML, cleanFileName string) (MS1, error) {

	var m MS1

	for _, i := range mz.MzML.Run.SpectrumList.Spectrum {
		for _, j := range i.CVParam {
			if strings.EqualFold(j.Name, "ms level") && strings.EqualFold(j.Value, "1") {

				var ms1 Ms1Scan
				ms1.Index = i.Index

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

				var ms1Peaks Ms1Spectrum

				for m := 0; m <= len(i.Peaks)-1; m++ {
					var peak Ms1Peak
					peak.Mz = i.Peaks[m]
					peak.Intensity = i.Intensities[m]
					ms1Peaks = append(ms1Peaks, peak)
				}

				ms1.Spectrum = ms1Peaks
				m.Ms1Scan = append(m.Ms1Scan, ms1)

			}
		}
	}

	return m, nil
}

// ReadMzXML parses only MS1 data from mzXML
func (m *MS1) ReadMzXML(f string) *err.Error {

	xmlFile, e := os.Open(f)

	if e != nil {
		return &err.Error{Type: err.CannotParseXML, Class: err.FATA, Argument: "mzXML"}
	}
	defer xmlFile.Close()

	// get the clean file name for spectra name
	ext := filepath.Ext(f)
	name := filepath.Base(f)
	cleanFileName := name[0 : len(name)-len(ext)]

	var mz mzx.MzXML
	e = mz.Parse(f)
	if e != nil {
		return &err.Error{Type: err.CannotParseXML, Class: err.FATA, Argument: "mzXML"}
	}

	for _, i := range mz.MSRun.Scan {

		if i.MSLevel == 1 {

			var ms1 Ms1Scan
			var ms1Peaks Ms1Spectrum
			var peaklist []float64
			var intlist []float64

			ms1.Index = string(i.Num)
			ms1.SpectrumName = fmt.Sprintf("%s.%d.%d", cleanFileName, i.Num, i.Num)

			preRT := i.RetentionTime
			preRT = strings.Replace(preRT, "PT", "", -1)
			preRT = strings.Replace(preRT, "S", "", -1)
			stt, _ := strconv.ParseFloat(preRT, 64)
			ms1.ScanStartTime = uti.Round((stt / 60), 5, 4)

			// convertedPeaks, err := mzxml.Decode(i.Peaks)
			// if err != nil {
			// 	return err
			// }

			var counter int
			for m := range i.ConvertedPeaks {
				if counter%2 == 0 {
					peaklist = append(peaklist, i.ConvertedPeaks[m])
				} else {
					intlist = append(intlist, i.ConvertedPeaks[m])
				}
				counter++
			}

			// var counter int
			// for m := range convertedPeaks {
			// 	if counter%2 == 0 {
			// 		peaklist = append(peaklist, convertedPeaks[m])
			// 	} else {
			// 		intlist = append(intlist, convertedPeaks[m])
			// 	}
			// 	counter++
			// }

			for n := range peaklist {
				var peak Ms1Peak
				peak.Mz = peaklist[n]
				peak.Intensity = intlist[n]
				ms1Peaks = append(ms1Peaks, peak)
			}

			ms1.Spectrum = ms1Peaks
			m.Ms1Scan = append(m.Ms1Scan, ms1)
		}

	}

	return nil
}

// GetMzXMLSpectra parses only MS1 data from mzXML
func GetMzXMLSpectra(mz mzx.MzXML, cleanFileName string) (MS1, error) {

	var m MS1

	for _, i := range mz.MSRun.Scan {

		if i.MSLevel == 1 {

			var ms1 Ms1Scan
			var ms1Peaks Ms1Spectrum
			var peaklist []float64
			var intlist []float64

			ms1.Index = string(i.Num)
			ms1.SpectrumName = fmt.Sprintf("%s.%d.%d", cleanFileName, i.Num, i.Num)

			preRT := i.RetentionTime
			preRT = strings.Replace(preRT, "PT", "", -1)
			preRT = strings.Replace(preRT, "S", "", -1)
			stt, _ := strconv.ParseFloat(preRT, 64)
			ms1.ScanStartTime = uti.Round((stt / 60), 5, 4)

			var counter int
			for m := range i.ConvertedPeaks {
				if counter%2 == 0 {
					peaklist = append(peaklist, i.ConvertedPeaks[m])
				} else {
					intlist = append(intlist, i.ConvertedPeaks[m])
				}
				counter++
			}

			for n := range peaklist {
				var peak Ms1Peak
				peak.Mz = peaklist[n]
				peak.Intensity = intlist[n]
				ms1Peaks = append(ms1Peaks, peak)
			}

			ms1.Spectrum = ms1Peaks
			m.Ms1Scan = append(m.Ms1Scan, ms1)
		}

	}

	if len(m.Ms1Scan) == 0 {
		return m, nil
	}

	return m, nil
}
