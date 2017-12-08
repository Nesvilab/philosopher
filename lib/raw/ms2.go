package raw

// MS2 top struct
type MS2 struct {
	Ms2Scan []Ms2Scan
}

// Ms2Scan tag
type Ms2Scan struct {
	Index         string
	Scan          string
	SpectrumName  string
	ScanStartTime float64
	Precursor     Precursor
	Spectrum      Spectrum
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

// // Ms2Spectrum tag
// type Ms2Spectrum []Ms2Peak
//
// // Ms2Peak tag
// type Ms2Peak struct {
// 	Mz        float64
// 	Intensity float64
// }
//
// func (a Ms2Spectrum) Len() int           { return len(a) }
// func (a Ms2Spectrum) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
// func (a Ms2Spectrum) Less(i, j int) bool { return a[i].Mz < a[j].Mz }

// // ReadMzML parses only MS2 data from mzML
// func (m *MS2) ReadMzML(f string) error {
//
// 	m.Ms2Scan = make(map[string]Ms2Scan)
//
// 	var mz mzm.IndexedMzML
// 	err := mz.Parse(f)
// 	if err != nil {
// 		return err
// 	}
//
// 	// get the clean file name for spectra name
// 	ext := filepath.Ext(f)
// 	name := filepath.Base(f)
// 	cleanFileName := name[0 : len(name)-len(ext)]
//
// 	for _, i := range mz.MzML.Run.SpectrumList.Spectrum {
// 		for _, j := range i.CVParam {
// 			if j.Name == "ms level" && j.Value == "2" {
//
// 				var ms2 Ms2Scan
// 				ms2.Index = i.Index
//
// 				// parse and format scan ID for the SpectrumName
// 				idSplit := strings.Split(i.ID, "scan=")
// 				iID, _ := strconv.Atoi(idSplit[1])
// 				cusID := fmt.Sprintf("%05d", iID)
// 				ms2.SpectrumName = fmt.Sprintf("%s.%s.%s", cleanFileName, cusID, cusID)
// 				ms2.Scan = cusID
//
// 				for _, k := range i.ScanList.Scan {
// 					for _, l := range k.CVParam {
// 						if l.Name == "scan start time" {
// 							stt, _ := strconv.ParseFloat(l.Value, 64)
// 							minrt := (stt * 60)
// 							ms2.ScanStartTime = uti.ToFixed(minrt, 3)
// 						}
// 					}
// 				}
//
// 				for _, j := range i.PrecursorList.Precursor {
//
// 					var prec Precursor
//
// 					split := strings.Split(j.SpectrumRef, " ")
// 					scan := strings.Split(split[2], "=")
//
// 					ind, err := strconv.Atoi(scan[1])
// 					if err != nil {
// 						return err
// 					}
//
// 					adjInd := ind - 1
// 					prec.ParentIndex = strconv.Itoa(adjInd)
// 					prec.ParentScan = scan[1]
//
// 					for _, k := range j.IsolationWindow.CVParam {
// 						if k.Accession == "MS:1000827" {
// 							val, err := strconv.ParseFloat(k.Value, 64)
// 							if err != nil {
// 								return nil
// 							}
// 							prec.TargetIon = val
//
// 						} else if k.Accession == "MS:1000828" {
// 							val, err := strconv.ParseFloat(k.Value, 64)
// 							if err != nil {
// 								return err
// 							}
// 							prec.IsolationWindowLowerOffset = val
//
// 						} else if k.Accession == "MS:1000829" {
// 							val, err := strconv.ParseFloat(k.Value, 64)
// 							if err != nil {
// 								return err
// 							}
// 							prec.IsolationWindowUpperOffset = val
//
// 						}
// 					}
//
// 					for _, k := range j.SelectedIonList.SelectedIon[0].CVParam {
// 						if k.Accession == "MS:1000744" {
// 							val, err := strconv.ParseFloat(k.Value, 64)
// 							if err != nil {
// 								return nil
// 							}
// 							prec.SelectedIon = val
//
// 						} else if k.Accession == "MS:1000041" {
// 							val, err := strconv.Atoi(k.Value)
// 							if err != nil {
// 								return err
// 							}
// 							prec.ChargeState = val
//
// 						} else if k.Accession == "MS:1000042" {
// 							val, err := strconv.ParseFloat(k.Value, 64)
// 							if err != nil {
// 								return nil
// 							}
// 							prec.PeakIntensity = val
//
// 						}
// 					}
//
// 					ms2.Precursor = prec
// 				}
//
// 				mzPeaks, err := mzm.Decode("mz", i.BinaryDataArrayList.BinaryDataArray[0])
// 				if err != nil {
// 					return err
// 				}
//
// 				mzIntensities, err := mzm.Decode("int", i.BinaryDataArrayList.BinaryDataArray[1])
// 				if err != nil {
// 					return err
// 				}
//
// 				var ms2Peaks Ms2Spectrum
//
// 				for m := range mzPeaks {
// 					var peak Ms2Peak
// 					peak.Mz = mzPeaks[m]
// 					peak.Intensity = mzIntensities[m]
// 					ms2Peaks = append(ms2Peaks, peak)
// 				}
//
// 				ms2.Spectrum = ms2Peaks
// 				var fullSpecName = fmt.Sprintf("%s.%s.%s.%d", cleanFileName, cusID, cusID, ms2.Precursor.ChargeState)
// 				m.Ms2Scan[fullSpecName] = ms2
//
// 			}
// 		}
// 	}
//
// 	return nil
// }
//
// // GetMzMLMS2Spectra parses only MS2 data from mzML
// func GetMzMLMS2Spectra(mz mzm.IndexedMzML, cleanFileName string) (MS2, error) {
//
// 	var m MS2
// 	m.Ms2Scan = make(map[string]Ms2Scan)
//
// 	for _, i := range mz.MzML.Run.SpectrumList.Spectrum {
// 		for _, j := range i.CVParam {
// 			if j.Name == "ms level" && j.Value == "2" {
//
// 				var ms2 Ms2Scan
// 				ms2.Index = i.Index
//
// 				// parse and format scan ID for the SpectrumName
// 				idSplit := strings.Split(i.ID, "scan=")
// 				iID, _ := strconv.Atoi(idSplit[1])
// 				cusID := fmt.Sprintf("%05d", iID)
// 				ms2.SpectrumName = fmt.Sprintf("%s.%s.%s", cleanFileName, cusID, cusID)
// 				ms2.Scan = cusID
//
// 				for _, k := range i.ScanList.Scan {
// 					for _, l := range k.CVParam {
// 						if l.Name == "scan start time" {
// 							stt, _ := strconv.ParseFloat(l.Value, 64)
// 							minrt := (stt * 60)
// 							ms2.ScanStartTime = uti.ToFixed(minrt, 3)
// 						}
// 					}
// 				}
//
// 				for _, j := range i.PrecursorList.Precursor {
//
// 					var prec Precursor
// 					split := strings.Split(j.SpectrumRef, " ")
// 					scan := strings.Split(split[2], "=")
//
// 					ind, err := strconv.Atoi(scan[1])
// 					if err != nil {
// 						return m, err
// 					}
//
// 					adjInd := ind - 1
// 					prec.ParentIndex = strconv.Itoa(adjInd)
// 					prec.ParentScan = scan[1]
//
// 					for _, k := range j.SelectedIonList.SelectedIon[0].CVParam {
// 						if k.Accession == "MS:1000744" {
// 							val, err := strconv.ParseFloat(k.Value, 64)
// 							if err != nil {
// 								return m, err
// 							}
// 							prec.SelectedIon = val
// 						} else if k.Accession == "MS:1000041" {
// 							val, err := strconv.Atoi(k.Value)
// 							if err != nil {
// 								return m, err
// 							}
// 							prec.ChargeState = val
// 						} else if k.Accession == "MS:1000042" {
// 							val, err := strconv.ParseFloat(k.Value, 64)
// 							if err != nil {
// 								return m, err
// 							}
// 							prec.PeakIntensity = val
// 						}
// 					}
//
// 					ms2.Precursor = prec
// 				}
//
// 				mzPeaks, err := mzm.Decode("mz", i.BinaryDataArrayList.BinaryDataArray[0])
// 				if err != nil {
// 					return m, err
// 				}
//
// 				mzIntensities, err := mzm.Decode("int", i.BinaryDataArrayList.BinaryDataArray[1])
// 				if err != nil {
// 					return m, err
// 				}
//
// 				var ms2Peaks Ms2Spectrum
//
// 				for m := range mzPeaks {
// 					var peak Ms2Peak
// 					peak.Mz = mzPeaks[m]
// 					peak.Intensity = mzIntensities[m]
// 					ms2Peaks = append(ms2Peaks, peak)
// 				}
//
// 				ms2.Spectrum = ms2Peaks
// 				var fullSpecName = fmt.Sprintf("%s.%s.%s.%d", cleanFileName, cusID, cusID, ms2.Precursor.ChargeState)
// 				m.Ms2Scan[fullSpecName] = ms2
// 				//m.Ms2Scan = append(m.Ms2Scan, ms2)
//
// 			}
// 		}
// 	}
//
// 	return m, nil
// }
//
// // ReadMzXML parses only MS2 data from mzXML
// func (m *MS2) ReadMzXML(f string) error {
//
// 	xmlFile, err := os.Open(f)
//
// 	if err != nil {
// 		logrus.Fatal("Error trying to read mzXML file")
// 	}
// 	defer xmlFile.Close()
//
// 	// get the clean file name for spectra name
// 	ext := filepath.Ext(f)
// 	name := filepath.Base(f)
// 	cleanFileName := name[0 : len(name)-len(ext)]
//
// 	decoder := xml.NewDecoder(xmlFile)
// 	decoder.CharsetReader = charset.NewReaderLabel
//
// 	var mz mzx.MzXML
// 	if err = decoder.Decode(&mz); err != nil {
// 		return err
// 	}
//
// 	for _, i := range mz.MSRun.Scan {
//
// 		if i.MSLevel == 2 {
//
// 			var ms2 Ms2Scan
// 			ms2.Index = string(i.Num)
// 			ms2.SpectrumName = fmt.Sprintf("%s.%d.%d", cleanFileName, i.Num, i.Num)
//
// 			preRT := i.RetentionTime
// 			preRT = strings.Replace(preRT, "PT", "", -1)
// 			preRT = strings.Replace(preRT, "S", "", -1)
// 			stt, _ := strconv.ParseFloat(preRT, 64)
// 			ms2.ScanStartTime = uti.Round((stt / 60), 5, 4)
//
// 			var prec Precursor
// 			prec.ParentIndex = string(i.Precursor.PrecursorScanNum)
// 			prec.ParentScan = string(i.Precursor.PrecursorScanNum)
// 			prec.ChargeState = i.Precursor.PrecursorCharge
// 			prec.PeakIntensity = i.Precursor.PrecursorIntensity
// 			ms2.Precursor = prec
//
// 			var ms2Peaks Ms2Spectrum
// 			var peaklist []float64
// 			var intlist []float64
//
// 			covertedPeaks, err := mzx.Decode(i.Peaks)
// 			if err != nil {
// 				return err
// 			}
//
// 			var counter int
// 			for m := range covertedPeaks {
// 				if counter%2 == 0 {
// 					peaklist = append(peaklist, covertedPeaks[m])
// 				} else {
// 					intlist = append(intlist, covertedPeaks[m])
// 				}
// 				counter++
// 			}
//
// 			if len(peaklist) == len(intlist) {
// 				for n := range peaklist {
// 					var peak Ms2Peak
// 					peak.Mz = peaklist[n]
// 					peak.Intensity = intlist[n]
// 					ms2Peaks = append(ms2Peaks, peak)
// 				}
// 			}
//
// 			ms2.Spectrum = ms2Peaks
// 			//m.Ms2Scan = append(m.Ms2Scan, ms2)
// 			var fullSpecName = fmt.Sprintf("%s.%d.%d.%d", cleanFileName, i.Num, i.Num, ms2.Precursor.ChargeState)
// 			m.Ms2Scan[fullSpecName] = ms2
// 		}
//
// 	}
//
// 	return nil
// }
//
// // GetMzXMLMS2Spectra parses only MS2 data from mzXML
// func GetMzXMLMS2Spectra(mz mzx.MzXML, cleanFileName string) (MS2, error) {
//
// 	var m MS2
//
// 	for _, i := range mz.MSRun.Scan {
//
// 		if i.MSLevel == 2 {
//
// 			var ms2 Ms2Scan
// 			ms2.Index = string(i.Num)
// 			ms2.SpectrumName = fmt.Sprintf("%s.%d.%d", cleanFileName, i.Num, i.Num)
//
// 			preRT := i.RetentionTime
// 			preRT = strings.Replace(preRT, "PT", "", -1)
// 			preRT = strings.Replace(preRT, "S", "", -1)
// 			stt, _ := strconv.ParseFloat(preRT, 64)
// 			ms2.ScanStartTime = uti.Round((stt / 60), 5, 4)
//
// 			var prec Precursor
// 			prec.ParentIndex = string(i.Precursor.PrecursorScanNum)
// 			prec.ParentScan = string(i.Precursor.PrecursorScanNum)
// 			prec.ChargeState = i.Precursor.PrecursorCharge
// 			prec.PeakIntensity = i.Precursor.PrecursorIntensity
// 			ms2.Precursor = prec
//
// 			var ms2Peaks Ms2Spectrum
// 			var peaklist []float64
// 			var intlist []float64
//
// 			covertedPeaks, err := mzx.Decode(i.Peaks)
// 			if err != nil {
// 				return m, err
// 			}
//
// 			var counter int
// 			for m := range covertedPeaks {
// 				if counter%2 == 0 {
// 					peaklist = append(peaklist, covertedPeaks[m])
// 				} else {
// 					intlist = append(intlist, covertedPeaks[m])
// 				}
// 				counter++
// 			}
//
// 			if len(peaklist) == len(intlist) {
// 				for n := range peaklist {
// 					var peak Ms2Peak
// 					peak.Mz = peaklist[n]
// 					peak.Intensity = intlist[n]
// 					ms2Peaks = append(ms2Peaks, peak)
// 				}
// 			}
//
// 			ms2.Spectrum = ms2Peaks
// 			var fullSpecName = fmt.Sprintf("%s.%d.%d.%d", cleanFileName, i.Num, i.Num, ms2.Precursor.ChargeState)
// 			m.Ms2Scan[fullSpecName] = ms2
//
// 			//m.Ms2Scan = append(m.Ms2Scan, ms2)
// 		}
//
// 	}
//
// 	return m, nil
// }
