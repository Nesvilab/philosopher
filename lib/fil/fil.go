package fil

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/prvst/philosopher/lib/cla"
	"github.com/prvst/philosopher/lib/dat"
	"github.com/prvst/philosopher/lib/id"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/qua"
	"github.com/prvst/philosopher/lib/rep"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/prvst/philosopher/lib/uti"
	"github.com/sirupsen/logrus"
)

// Run executes the Filter processing
//func (f *Filter) Run(psmFDR, pepFDR, ionFDR, ptFDR, pepProb, protProb float64, isPicked, isRazor, mapmod bool) error {
func Run(f met.Data) (met.Data, error) {

	e := rep.New()
	var pepxml id.PepXML
	var pep id.PepIDList
	var pro id.ProtIDList
	var err error

	logrus.Info("Processing peptide identification files")

	pepid, searchEngine, err := readPepXMLInput(f.Filter.Pex, f.Filter.Tag, f.Filter.Model)
	if err != nil {
		return f, err
	}

	f.SearchEngine = searchEngine

	err = processPeptideIdentifications(pepid, f.Filter.Tag, f.Filter.PsmFDR, f.Filter.PepFDR, f.Filter.IonFDR)
	if err != nil {
		return f, err
	}

	if len(f.Filter.Pox) > 0 {

		protXML, proerr := readProtXMLInput(sys.MetaDir(), f.Filter.Pox, f.Filter.Tag, f.Filter.Weight)
		if proerr != nil {
			return f, proerr
		}

		err = processProteinIdentifications(protXML, f.Filter.PtFDR, f.Filter.PepFDR, f.Filter.ProtProb, f.Filter.Picked, f.Filter.Razor)
		if err != nil {
			return f, err
		}
		//protXML = id.ProtXML{}

		if f.Filter.Seq == true {

			// sequential analysis
			// filtered psm list and filtered prot list
			pep.Restore("psm")
			pro.Restore()
			err = sequentialFDRControl(pep, pro, f.Filter.PsmFDR, f.Filter.PepFDR, f.Filter.IonFDR, f.Filter.Tag)
			if err != nil {
				return f, err
			}
			pep = nil
			pro = nil

		} else {

			// two-dimensional analysis
			// complete pep list and filtered mirror-image prot list
			pepxml.Restore()
			pro.Restore()
			err = twoDFDRFilter(pepxml.PeptideIdentification, pro, f.Filter.PsmFDR, f.Filter.PepFDR, f.Filter.IonFDR, f.Filter.Tag)
			if err != nil {
				return f, err
			}
			pepxml = id.PepXML{}
			pro = nil

		}

	}

	var dtb dat.Base
	dtb.Restore()
	if len(dtb.Records) < 1 {
		return f, errors.New("Database data not available, interrupting processing")
	}

	logrus.Info("Post processing identifications")

	// restoring for the modifications
	var pxml id.PepXML
	pxml.Restore()
	e.Mods.DefinedModAminoAcid = pxml.DefinedModAminoAcid
	e.Mods.DefinedModMassDiff = pxml.DefinedModMassDiff
	pxml = id.PepXML{}

	var psm id.PepIDList
	psm.Restore("psm")
	e.AssemblePSMReport(psm, f.Filter.Tag)
	psm = nil

	// evaluate modifications in data set
	if f.Filter.Mapmods == true {
		logrus.Info("Mapping modifications")
		e.MapMassDiffToUniMod()

		logrus.Info("Processing modifications")
		e.AssembleModificationReport()
	}

	var ion id.PepIDList
	ion.Restore("ion")
	e.AssembleIonReport(ion, f.Filter.Tag)
	ion = nil

	var pept id.PepIDList
	pept.Restore("pep")
	e.AssemblePeptideReport(pept, f.Filter.Tag)
	pept = nil

	// evaluate modifications in data set
	if f.Filter.Mapmods == true {
		e.UpdateIonModCount()
		e.UpdatePeptideModCount()
	}

	logrus.Info("Processing Protein Inference")
	pro.Restore()
	err = e.AssembleProteinReport(pro, f.Filter.Tag)
	if err != nil {
		return f, err
	}
	pro = nil

	logrus.Info("Correcting PSM to Protein mappings")
	e.UpdateMappedProteins()

	// ADD ERROR CASES
	logrus.Info("Mapping Ion status to PSMs")
	e.UpdateIonStatus()

	// reorganizes the selected proteins and the alternative proteins list
	if f.Filter.Razor == true {
		logrus.Info("Updating razor PSM assingment to Proteins")
		e.UpdateProteinStatus()
		e.UpdateGeneNames()
		e.UpdateSupportingSpectra()
	}

	logrus.Info("Calculating Spectral Counts")
	e, cerr := qua.CalculateSpectralCounts(e)
	if cerr != nil {
		return f, cerr
	}

	cerr = e.SerializeGranular()
	if cerr != nil {
		return f, cerr
	}

	return f, nil
}

// readPepXMLInput reads one or more fies and organize the data into PSM list
func readPepXMLInput(xmlFile, decoyTag string, models bool) (id.PepIDList, string, error) {

	var files []string
	var pepIdent id.PepIDList
	var definedModMassDiff = make(map[float64]float64)
	var definedModAminoAcid = make(map[float64]string)
	var searchEngine string

	if strings.Contains(xmlFile, "pep.xml") || strings.Contains(xmlFile, "pepXML") {
		files = append(files, xmlFile)
	} else {
		glob := fmt.Sprintf("%s/*pep.xml", xmlFile)
		list, _ := filepath.Glob(glob)

		if len(list) == 0 {
			return pepIdent, "", errors.New("No pepXML files found, check your files and try again")
		}

		for _, i := range list {
			absPath, _ := filepath.Abs(i)
			files = append(files, absPath)
		}

	}

	for _, i := range files {
		var p id.PepXML
		p.DecoyTag = decoyTag
		e := p.Read(i)
		if e != nil {
			return nil, "", e
		}

		// print models
		if models == true {
			if strings.EqualFold(p.Prophet, "interprophet") {
				logrus.Error("Cannot print models for interprophet files")
			} else {
				logrus.Info("Printing models")
				temp, _ := sys.GetTemp()
				go p.ReportModels(temp, filepath.Base(i))
				time.Sleep(time.Second * 3)
			}
		}

		pepIdent = append(pepIdent, p.PeptideIdentification...)

		for k, v := range p.DefinedModAminoAcid {
			definedModAminoAcid[k] = v
		}

		for k, v := range p.DefinedModMassDiff {
			definedModMassDiff[k] = v
		}

		searchEngine = p.SearchEngine
	}

	// create a "fake" global pepXML comprising all data
	var pepXML id.PepXML
	pepXML.DecoyTag = decoyTag
	pepXML.PeptideIdentification = pepIdent
	pepXML.DefinedModAminoAcid = definedModAminoAcid
	pepXML.DefinedModMassDiff = definedModMassDiff

	// promoting Spectra that matches to both decoys and targets to TRUE hits
	pepXML.PromoteProteinIDs()

	// serialize all pep files
	sort.Sort(pepXML.PeptideIdentification)
	pepXML.Serialize()

	return pepIdent, searchEngine, nil
}

// processPeptideIdentifications reads and process pepXML
func processPeptideIdentifications(p id.PepIDList, decoyTag string, psm, peptide, ion float64) error {

	var err error

	// report charge profile
	var t, d int

	t, d, _ = chargeProfile(p, 1, decoyTag)
	logrus.WithFields(logrus.Fields{
		"target": t,
		"decoy":  d,
	}).Info("1+ Charge profile")

	t, d, _ = chargeProfile(p, 2, decoyTag)
	logrus.WithFields(logrus.Fields{
		"target": t,
		"decoy":  d,
	}).Info("2+ Charge profile")

	t, d, _ = chargeProfile(p, 3, decoyTag)
	logrus.WithFields(logrus.Fields{
		"target": t,
		"decoy":  d,
	}).Info("3+ Charge profile")

	t, d, _ = chargeProfile(p, 4, decoyTag)
	logrus.WithFields(logrus.Fields{
		"target": t,
		"decoy":  d,
	}).Info("4+ Charge profile")

	t, d, _ = chargeProfile(p, 5, decoyTag)
	logrus.WithFields(logrus.Fields{
		"target": t,
		"decoy":  d,
	}).Info("5+ Charge profile")

	t, d, _ = chargeProfile(p, 6, decoyTag)
	logrus.WithFields(logrus.Fields{
		"target": t,
		"decoy":  d,
	}).Info("6+ Charge profile")

	uniqPsms := getUniquePSMs(p)
	uniqPeps := getUniquePeptides(p)
	uniqIons := getUniquePeptideIons(p)

	logrus.WithFields(logrus.Fields{
		"psms":     len(p),
		"peptides": len(uniqPeps),
		"ions":     len(uniqIons),
	}).Info("Database search results")

	filteredPSM, err := pepXMLFDRFilter(uniqPsms, psm, "PSM", decoyTag)
	if err != nil {
		logrus.Fatal(err)
	}
	filteredPSM.Serialize("psm")

	filteredPeptides, err := pepXMLFDRFilter(uniqPeps, peptide, "Peptide", decoyTag)
	if err != nil {
		logrus.Fatal(err)
	}
	filteredPeptides.Serialize("pep")

	filteredIons, err := pepXMLFDRFilter(uniqIons, ion, "Ion", decoyTag)
	if err != nil {
		logrus.Fatal(err)
	}
	filteredIons.Serialize("ion")

	return nil
}

// chargeProfile ...
func chargeProfile(p id.PepIDList, charge uint8, decoyTag string) (t, d int, err error) {

	for _, i := range p {
		if i.AssumedCharge == charge {
			if strings.Contains(i.Protein, decoyTag) {
				d++
			} else {
				t++
			}
		}
	}

	if t < 1 || d < 1 {
		err = errors.New("Invalid charge state count")
	}

	return t, d, err
}

//getUniquePSMs selects only unique pepetide ions for the given data stucture
func getUniquePSMs(p id.PepIDList) map[string]id.PepIDList {

	uniqMap := make(map[string]id.PepIDList)

	for _, i := range p {
		uniqMap[i.Spectrum] = append(uniqMap[i.Spectrum], i)
	}

	return uniqMap
}

//getUniquePeptideIons selects only unique pepetide ions for the given data stucture
func getUniquePeptideIons(p id.PepIDList) map[string]id.PepIDList {

	uniqMap := ExtractIonsFromPSMs(p)

	return uniqMap
}

// ExtractIonsFromPSMs takes a pepidlist and transforms into an ion map
func ExtractIonsFromPSMs(p id.PepIDList) map[string]id.PepIDList {

	uniqMap := make(map[string]id.PepIDList)

	for _, i := range p {
		ion := fmt.Sprintf("%s#%d#%.4f", i.Peptide, i.AssumedCharge, i.CalcNeutralPepMass)
		uniqMap[ion] = append(uniqMap[ion], i)
	}

	// organize id list by score
	for _, v := range uniqMap {
		sort.Sort(v)
	}

	return uniqMap
}

// getUniquePeptides selects only unique pepetide for the given data stucture
func getUniquePeptides(p id.PepIDList) map[string]id.PepIDList {

	uniqMap := make(map[string]id.PepIDList)

	for _, i := range p {
		uniqMap[string(i.Peptide)] = append(uniqMap[string(i.Peptide)], i)
	}

	// organize id list by score
	for _, v := range uniqMap {
		sort.Sort(v)
	}

	return uniqMap
}

// pepXMLFDRFilter applies FDR filtering at the PSM level
func pepXMLFDRFilter(input map[string]id.PepIDList, targetFDR float64, level, decoyTag string) (id.PepIDList, error) {

	//var msg string
	var targets float64
	var decoys float64
	var calcFDR float64
	var list id.PepIDList
	var peplist id.PepIDList
	var minProb float64 = 10
	var err error

	if strings.EqualFold(level, "PSM") {

		// move all entries to list and count the number of targets and decoys
		for _, i := range input {
			for _, j := range i {
				if cla.IsDecoyPSM(j, decoyTag) {
					decoys++
				} else {
					targets++
				}
				list = append(list, j)
			}
		}

	} else if strings.EqualFold(level, "Peptide") {

		// 0 index means the one with highest score
		for _, i := range input {
			peplist = append(peplist, i[0])
		}

		for i := range peplist {
			if cla.IsDecoyPSM(peplist[i], decoyTag) {
				decoys++
			} else {
				targets++
			}
			list = append(list, peplist[i])
		}

	} else if strings.EqualFold(level, "Ion") {

		// 0 index means the one with highest score
		for _, i := range input {
			peplist = append(peplist, i[0])
		}

		for i := range peplist {
			if cla.IsDecoyPSM(peplist[i], decoyTag) {
				decoys++
			} else {
				targets++
			}
			list = append(list, peplist[i])
		}

	} else {
		err = errors.New("Error applying FDR score; unknown level")
	}

	sort.Sort(list)

	var scoreMap = make(map[float64]float64)
	limit := (len(list) - 1)
	for j := limit; j >= 0; j-- {
		_, ok := scoreMap[list[j].Probability]
		if !ok {
			scoreMap[list[j].Probability] = (decoys / targets)
		}
		if cla.IsDecoyPSM(list[j], decoyTag) {
			decoys--
		} else {
			targets--
		}
	}

	var keys []float64
	for k := range scoreMap {
		keys = append(keys, k)
	}

	sort.Sort(sort.Reverse(sort.Float64Slice(keys)))

	var probList = make(map[float64]uint8)
	for i := range keys {

		//f := fmt.Sprintf("%.2f", scoreMap[keys[i]]*100)
		//f := uti.Round(scoreMap[keys[i]]*100, 5, 2)
		//fmt.Println(keys[i], "\t", scoreMap[keys[i]], "\t", uti.ToFixed(scoreMap[keys[i]], 4), "\t", f, "\t", targetFDR)

		if uti.ToFixed(scoreMap[keys[i]], 4) <= targetFDR {
			probList[keys[i]] = 0
			minProb = keys[i]
			calcFDR = uti.ToFixed(scoreMap[keys[i]], 4)
		}

	}

	var cleanlist id.PepIDList
	decoys = 0
	targets = 0

	for i := range list {
		_, ok := probList[list[i].Probability]
		if ok {
			cleanlist = append(cleanlist, list[i])
			if cla.IsDecoyPSM(list[i], decoyTag) {
				decoys++
			} else {
				targets++
			}
		}
	}

	msg := fmt.Sprintf("Converged to %.2f %% FDR with %0.f %ss", (calcFDR * 100), targets, level)
	logrus.WithFields(logrus.Fields{
		"decoy":     decoys,
		"total":     (targets + decoys),
		"threshold": minProb,
	}).Info(msg)

	return cleanlist, err
}

// readProtXMLInput reads one or more fies and organize the data into PSM list
func readProtXMLInput(meta, xmlFile, decoyTag string, weight float64) (id.ProtXML, error) {

	var protXML id.ProtXML

	err := protXML.Read(xmlFile)
	if err != nil {
		return protXML, err
	}

	protXML.DecoyTag = decoyTag

	protXML.MarkUniquePeptides(weight)

	protXML.PromoteProteinIDs()

	protXML.Serialize()

	return protXML, nil
}

// processProteinIdentifications checks if pickedFDR ar razor options should be applied to given data set, if they do,
// the inputed protXML data is processed before filtered.
func processProteinIdentifications(p id.ProtXML, ptFDR, pepProb, protProb float64, isPicked, isRazor bool) error {

	var err error
	var pid id.ProtIDList

	// tagget / decoy / threshold
	t, d, _ := proteinProfile(p)
	logrus.WithFields(logrus.Fields{
		"target": t,
		"decoy":  d,
	}).Info("Protein inference results")

	// applies pickedFDR algorithm
	if isPicked == true {
		p = PickedFDR(p)
	}

	// applies razor algorithm
	if isRazor == true {
		p, err = RazorFilter(p)
		if err != nil {
			return err
		}
	}

	// run the FDR filter for proteins
	pid, err = ProtXMLFilter(p, ptFDR, pepProb, protProb, isPicked, isRazor)
	if err != nil {
		return err
	}

	// save results on meta folder
	pid.Serialize()

	return nil
}

// proteinProfile ...
func proteinProfile(p id.ProtXML) (t, d int, err error) {

	for _, i := range p.Groups {
		for _, j := range i.Proteins {
			if cla.IsDecoyProtein(j, p.DecoyTag) {
				d++
			} else {
				t++
			}
		}
	}

	return t, d, err
}

// PickedFDR employs the picked FDR strategy
func PickedFDR(p id.ProtXML) id.ProtXML {

	// var appMap = make(map[string]int)
	var targetMap = make(map[string]float64)
	var decoyMap = make(map[string]float64)
	var recordMap = make(map[string]int)

	// collect all proteins from every group
	for _, i := range p.Groups {
		for _, j := range i.Proteins {
			if cla.IsDecoyProtein(j, p.DecoyTag) {
				decoyMap[string(j.ProteinName)] = j.PeptideIons[0].InitialProbability
			} else {
				targetMap[string(j.ProteinName)] = j.PeptideIons[0].InitialProbability
			}
		}
	}

	// check unique targets
	for k := range targetMap {
		iKey := fmt.Sprintf("%s%s", p.DecoyTag, k)
		_, ok := decoyMap[iKey]
		if !ok {
			recordMap[k] = 1
		}
	}

	// check unique decoys
	for k := range decoyMap {
		iKey := strings.Replace(k, p.DecoyTag, "", -1)
		_, ok := targetMap[iKey]
		if !ok {
			recordMap[k] = 1
		}
	}

	// check paired observations
	for k, v := range targetMap {
		iKey := fmt.Sprintf("%s%s", p.DecoyTag, k)
		vok, ok := decoyMap[iKey]
		if ok {
			if vok > v {
				recordMap[k] = 0
				recordMap[iKey] = 1
			} else if v > vok {
				recordMap[k] = 1
				recordMap[iKey] = 0
			} else {
				recordMap[k] = 1
				recordMap[iKey] = 1
			}
		}
	}

	// collect all proteins from every group
	for i := range p.Groups {
		for j := range p.Groups[i].Proteins {
			v, ok := recordMap[string(p.Groups[i].Proteins[j].ProteinName)]
			if ok {
				p.Groups[i].Proteins[j].Picked = v
			}
		}
	}

	return p
}

// RazorFilter classifies peptides as razor
// func RazorFilter(p id.ProtXML) (id.ProtXML, error) {
//
// 	//var razorMap = make(map[string]string)
// 	var razorMap = make(map[string][]string)
// 	var refMap = make(map[string][]string)
// 	var pepMap = make(map[string]string)
// 	var checkRazor = make(map[string]uint8)
//
// 	// create reference entries by collapsing into a single string all necessary information about the peptide-to-protein assignment
// 	for i := range p.Groups {
// 		for j := range p.Groups[i].Proteins {
// 			for k := range p.Groups[i].Proteins[j].PeptideIons {
//
// 				ref := fmt.Sprintf("%d#%s#%s#%s#%f#%f#%f#%d#%.4f",
// 					p.Groups[i].GroupNumber,
// 					string(p.Groups[i].Proteins[j].GroupSiblingID),
// 					string(p.Groups[i].Proteins[j].ProteinName),
// 					string(p.Groups[i].Proteins[j].PeptideIons[k].PeptideSequence),
// 					p.Groups[i].Proteins[j].PeptideIons[k].InitialProbability,
// 					p.Groups[i].Proteins[j].PeptideIons[k].Weight,
// 					p.Groups[i].Proteins[j].PeptideIons[k].GroupWeight,
// 					p.Groups[i].Proteins[j].PeptideIons[k].Charge,
// 					p.Groups[i].Proteins[j].PeptideIons[k].CalcNeutralPepMass,
// 				)
//
// 				refMap[string(p.Groups[i].Proteins[j].PeptideIons[k].PeptideSequence)] = append(refMap[string(p.Groups[i].Proteins[j].PeptideIons[k].PeptideSequence)], ref)
// 				pepMap[string(p.Groups[i].Proteins[j].PeptideIons[k].PeptideSequence)] = ""
//
// 			}
// 		}
// 	}
//
// 	// for each unique peptide sequence
// 	for k := range pepMap {
//
// 		var gw float64
// 		var w float64
// 		var mgw []string
//
// 		// retrieve the list of references based on the peptide sequence
// 		v, ok := refMap[k]
// 		if ok {
//
// 			// for each reference in the list
// 			for i := range v {
//
// 				pep := strings.Split(v[i], "#")
//
// 				weight, err := strconv.ParseFloat(pep[5], 64)
// 				if err != nil {
// 					return p, err
// 				}
//
// 				groupWeight, err := strconv.ParseFloat(pep[6], 64)
// 				if err != nil {
// 					return p, err
// 				}
//
// 				cnpm, err := strconv.ParseFloat(pep[8], 64)
// 				if err != nil {
// 					return p, err
// 				}
//
// 				pepCheck := fmt.Sprintf("%s#%s#%.4f", pep[3], pep[7], cnpm)
//
// 				// references with weight > 0.5 are easy cases, and clearly assinged as razor
// 				if weight > 0.5 {
//
// 					mgw = append(mgw, v[i])
// 					_, ok := checkRazor[pepCheck]
// 					if !ok {
// 						mgw = append(mgw, v[i])
// 						checkRazor[pepCheck] = 1
// 					}
//
// 				} else {
//
// 					if groupWeight > gw {
// 						gw = groupWeight
// 						w = weight
// 						_, ok := checkRazor[pepCheck]
// 						if !ok {
// 							mgw = append(mgw, v[i])
// 							checkRazor[pepCheck] = 1
// 						}
//
// 					} else if groupWeight == gw {
// 						if weight > w {
// 							w = weight
//
// 							_, ok := checkRazor[pepCheck]
// 							if !ok {
// 								mgw = append(mgw, v[i])
// 								checkRazor[pepCheck] = 1
// 							}
//
// 						} else if weight == w {
// 							w = weight
//
// 							_, ok := checkRazor[pepCheck]
// 							if !ok {
// 								mgw = append(mgw, v[i])
// 								checkRazor[pepCheck] = 1
// 							}
//
// 						}
// 					}
//
// 				}
//
// 			}
// 		}
// 		razorMap[k] = mgw
// 	}
//
// 	for i := range p.Groups {
// 		for j := range p.Groups[i].Proteins {
// 			for k := range p.Groups[i].Proteins[j].PeptideIons {
//
// 				ref := fmt.Sprintf("%d#%s#%s#%s#%f#%f#%f#%d#%.4f",
// 					p.Groups[i].GroupNumber,
// 					string(p.Groups[i].Proteins[j].GroupSiblingID),
// 					string(p.Groups[i].Proteins[j].ProteinName),
// 					string(p.Groups[i].Proteins[j].PeptideIons[k].PeptideSequence),
// 					p.Groups[i].Proteins[j].PeptideIons[k].InitialProbability,
// 					p.Groups[i].Proteins[j].PeptideIons[k].Weight,
// 					p.Groups[i].Proteins[j].PeptideIons[k].GroupWeight,
// 					p.Groups[i].Proteins[j].PeptideIons[k].Charge,
// 					p.Groups[i].Proteins[j].PeptideIons[k].CalcNeutralPepMass)
//
// 				v, ok := razorMap[string(p.Groups[i].Proteins[j].PeptideIons[k].PeptideSequence)]
// 				if ok {
//
// 					for _, l := range v {
// 						if strings.EqualFold(ref, l) {
// 							p.Groups[i].Proteins[j].PeptideIons[k].Razor = 1
// 							p.Groups[i].Proteins[j].HasRazor = true
// 							break
// 						}
// 					}
// 				}
//
// 				// for debuging
// 				// if strings.Contains(string(p.Groups[i].Proteins[j].ProteinName), "") {
// 				// 	fmt.Println("PING")
// 				// 	fmt.Println(p.Groups[i].Proteins[j].PeptideIons[k].Razor)
// 				// 	fmt.Println(p.Groups[i].Proteins[j].HasRazor)
// 				// 	os.Exit(1)
// 				// }
//
// 			}
// 		}
// 	}
//
// 	// mark as razor all peptides in the reference map
// 	for i := range p.Groups {
// 		for j := range p.Groups[i].Proteins {
// 			var r float64
// 			for k := range p.Groups[i].Proteins[j].PeptideIons {
// 				if p.Groups[i].Proteins[j].PeptideIons[k].Razor == 1 || p.Groups[i].Proteins[j].PeptideIons[k].IsUnique {
// 					if p.Groups[i].Proteins[j].PeptideIons[k].InitialProbability > r {
// 						r = p.Groups[i].Proteins[j].PeptideIons[k].InitialProbability
// 					}
// 				}
// 			}
// 			p.Groups[i].Proteins[j].TopPepProb = r
// 		}
// 	}
//
// 	return p, nil
// }

// RazorCandidate is a peptide sequence to be evaluated as a razor
type RazorCandidate struct {
	Sequence          string
	MappedProteinsW   map[string]float64
	MappedProteinsGW  map[string]float64
	MappedProteinsTNP map[string]int
	MappedProtein     string
}

// RazorCandidateMap is a list of razor candidates
type RazorCandidateMap map[string]RazorCandidate

// // RazorFilter classifies peptides as razor
// func RazorFilter(p id.ProtXML) (id.ProtXML, error) {
//
// 	var r = make(map[string]RazorCandidate)
//
// 	// for each peptide sequence, collapse oll parent protein peptides from ions originated from the same sequence
// 	for _, i := range p.Groups {
// 		for _, j := range i.Proteins {
// 			for _, k := range j.PeptideIons {
//
// 				v, ok := r[k.PeptideSequence]
// 				if !ok {
//
// 					var rc RazorCandidate
// 					rc.Sequence = k.PeptideSequence
// 					rc.MappedProteinsW = make(map[string]float64)
// 					rc.MappedProteinsGW = make(map[string]float64)
// 					rc.MappedProteinsTNP = make(map[string]int)
//
// 					rc.MappedProteinsW[j.ProteinName] = k.Weight
// 					rc.MappedProteinsGW[j.ProteinName] = k.GroupWeight
// 					rc.MappedProteinsTNP[j.ProteinName] = j.TotalNumberPeptides
//
// 					for _, i := range j.IndistinguishableProtein {
// 						rc.MappedProteinsW[i] = -1
// 						rc.MappedProteinsGW[i] = -1
// 						rc.MappedProteinsTNP[i] = -1
// 					}
//
// 					for _, i := range k.PeptideParentProtein {
// 						rc.MappedProteinsW[i] = -1
// 						rc.MappedProteinsGW[i] = -1
// 						rc.MappedProteinsTNP[i] = -1
// 					}
//
// 					r[k.PeptideSequence] = rc
//
// 				} else {
// 					var c = v
//
// 					if c.MappedProteinsW[j.ProteinName] == -1 {
// 						c.MappedProteinsW[j.ProteinName] = k.Weight
// 					}
//
// 					if c.MappedProteinsGW[j.ProteinName] == -1 {
// 						c.MappedProteinsGW[j.ProteinName] = k.GroupWeight
// 					}
//
// 					if c.MappedProteinsTNP[j.ProteinName] == -1 {
// 						c.MappedProteinsTNP[j.ProteinName] = j.TotalNumberPeptides
// 					}
//
// 					r[k.PeptideSequence] = c
//
// 				}
// 			}
// 		}
// 	}
//
// 	var razorPair = make(map[string]string)
//
// 	// get the best protein candidate for each pepetide sequence and make the razor pair
// 	for k, v := range r {
// 		// 1st pass: mark all cases with weight > 0.5
// 		for pt, w := range v.MappedProteinsW {
// 			if w > 0.5 {
// 				razorPair[k] = pt
// 			}
// 		}
// 	}
//
// 	for xk, xv := range razorPair {
// 		fmt.Println(xk, "\t", xv)
// 	}
//
// 	// 2nd pass: mark all cases with highest group weight in the list
// 	for k, v := range r {
//
// 		_, ok := razorPair[k]
// 		if !ok {
//
// 			var topPT string
// 			var topCount int
// 			var topGW float64
// 			var topGWMap = make(map[float64]uint8)
//
// 			if len(v.MappedProteinsGW) == 1 {
//
// 				for pt := range v.MappedProteinsGW {
// 					razorPair[k] = pt
// 				}
//
// 			} else if len(v.MappedProteinsGW) > 1 {
//
// 				for pt, tnp := range v.MappedProteinsGW {
// 					if tnp >= topGW {
// 						topGW = tnp
// 						topPT = pt
// 						topGWMap[topGW]++
// 					}
// 				}
//
// 				var tie bool
// 				if topGWMap[topGW] >= 2 {
// 					tie = true
// 				}
//
// 				if tie == false {
// 					razorPair[k] = topPT
//
// 				} else {
//
// 					for pt, tnp := range v.MappedProteinsTNP {
// 						if tnp >= topCount {
// 							topCount = tnp
// 							topPT = pt
// 						}
// 					}
//
// 					razorPair[k] = topPT
//
// 				}
// 			} else {
//
// 				for pt, tnp := range v.MappedProteinsTNP {
// 					if tnp >= topCount {
// 						topCount = tnp
// 						topPT = pt
// 					}
// 				}
//
// 				razorPair[k] = topPT
//
// 			}
// 		}
//
// 	}
//
// 	for k, v := range r {
// 		pt, ok := razorPair[k]
// 		if ok {
// 			razor := v
// 			razor.MappedProtein = pt
// 			r[k] = razor
// 		}
// 	}
//
// 	for i := range p.Groups {
// 		for j := range p.Groups[i].Proteins {
// 			for k := range p.Groups[i].Proteins[j].PeptideIons {
// 				v, ok := r[string(p.Groups[i].Proteins[j].PeptideIons[k].PeptideSequence)]
// 				if ok {
// 					if p.Groups[i].Proteins[j].ProteinName == v.MappedProtein {
// 						p.Groups[i].Proteins[j].PeptideIons[k].Razor = 1
// 						p.Groups[i].Proteins[j].HasRazor = true
// 					}
// 				}
// 			}
// 		}
// 	}
//
// 	// 	// mark as razor all peptides in the reference map
// 	for i := range p.Groups {
// 		for j := range p.Groups[i].Proteins {
// 			var r float64
// 			for k := range p.Groups[i].Proteins[j].PeptideIons {
// 				if p.Groups[i].Proteins[j].PeptideIons[k].Razor == 1 || p.Groups[i].Proteins[j].PeptideIons[k].IsUnique {
// 					if p.Groups[i].Proteins[j].PeptideIons[k].InitialProbability > r {
// 						r = p.Groups[i].Proteins[j].PeptideIons[k].InitialProbability
// 					}
// 				}
// 			}
// 			p.Groups[i].Proteins[j].TopPepProb = r
// 		}
// 	}
//
// 	return p, nil
// }

// // RazorFilter classifies peptides as razor
// func RazorFilter(p id.ProtXML) (id.ProtXML, error) {
//
// 	var r = make(map[string]RazorCandidate)
// 	var rList []string
//
// 	// for each peptide sequence, collapse all parent protein peptides from ions originated from the same sequence
// 	for _, i := range p.Groups {
// 		for _, j := range i.Proteins {
// 			for _, k := range j.PeptideIons {
//
// 				v, ok := r[k.PeptideSequence]
// 				if !ok {
//
// 					var rc RazorCandidate
// 					rc.Sequence = k.PeptideSequence
// 					rc.MappedProteinsW = make(map[string]float64)
// 					rc.MappedProteinsGW = make(map[string]float64)
// 					rc.MappedProteinsTNP = make(map[string]int)
//
// 					rc.MappedProteinsW[j.ProteinName] = k.Weight
// 					rc.MappedProteinsGW[j.ProteinName] = k.GroupWeight
// 					rc.MappedProteinsTNP[j.ProteinName] = j.TotalNumberPeptides
//
// 					for _, i := range j.IndistinguishableProtein {
// 						rc.MappedProteinsW[i] = -1
// 						rc.MappedProteinsGW[i] = -1
// 						rc.MappedProteinsTNP[i] = -1
// 					}
//
// 					for _, i := range k.PeptideParentProtein {
// 						rc.MappedProteinsW[i] = -1
// 						rc.MappedProteinsGW[i] = -1
// 						rc.MappedProteinsTNP[i] = -1
// 					}
//
// 					r[k.PeptideSequence] = rc
//
// 				} else {
// 					var c = v
//
// 					// doing like this will allow proteins that map to shared peptidesto be considered
// 					c.MappedProteinsW[j.ProteinName] = k.Weight
// 					c.MappedProteinsGW[j.ProteinName] = k.GroupWeight
// 					c.MappedProteinsTNP[j.ProteinName] = j.TotalNumberPeptides
//
// 					// if c.MappedProteinsW[j.ProteinName] == -1 {
// 					// 	c.MappedProteinsW[j.ProteinName] = k.Weight
// 					// }
// 					//
// 					// if c.MappedProteinsGW[j.ProteinName] == -1 {
// 					// 	c.MappedProteinsGW[j.ProteinName] = k.GroupWeight
// 					// }
// 					//
// 					// if c.MappedProteinsTNP[j.ProteinName] == -1 {
// 					// 	c.MappedProteinsTNP[j.ProteinName] = j.TotalNumberPeptides
// 					// }
//
// 					r[k.PeptideSequence] = c
//
// 				}
// 			}
// 		}
// 	}
//
// 	// spew.Dump(r)
// 	// os.Exit(1)
//
// 	// this will make the assignment more deterministic
// 	for k := range r {
// 		rList = append(rList, k)
// 	}
// 	sort.Strings(rList)
//
// 	var razorPair = make(map[string]string)
//
// 	// get the best protein candidate for each pepetide sequence and make the razor pair
// 	for _, k := range rList {
// 		// 1st pass: mark all cases with weight > 0.5
// 		for pt, w := range r[k].MappedProteinsW {
// 			if w > 0.5 {
// 				razorPair[k] = pt
// 			}
// 		}
// 	}
//
// 	// 2nd pass: mark all cases with highest group weight in the list
// 	for _, k := range rList {
//
// 		_, ok := razorPair[k]
// 		if !ok {
//
// 			var topPT string
// 			var topCount int
// 			var topGW float64
// 			var topGWMap = make(map[float64]uint8)
//
// 			if len(r[k].MappedProteinsGW) == 1 {
//
// 				for pt := range r[k].MappedProteinsGW {
// 					razorPair[k] = pt
// 				}
//
// 			} else if len(r[k].MappedProteinsGW) > 1 {
//
// 				for pt, tnp := range r[k].MappedProteinsGW {
// 					if tnp >= topGW {
// 						topGW = tnp
// 						topPT = pt
// 						topGWMap[topGW]++
// 					}
// 				}
//
// 				var tie bool
// 				if topGWMap[topGW] >= 2 {
// 					tie = true
// 				}
//
// 				if tie == false {
// 					razorPair[k] = topPT
//
// 				} else {
//
// 					var mplist []string
// 					for pt, _ := range r[k].MappedProteinsTNP {
// 						mplist = append(mplist, pt)
// 					}
// 					sort.Strings(mplist)
//
// 					for _, pt := range mplist {
// 						if r[k].MappedProteinsTNP[pt] >= topCount {
// 							topCount = r[k].MappedProteinsTNP[pt]
// 							topPT = pt
// 							break
// 						}
// 					}
//
// 					razorPair[k] = topPT
// 				}
//
// 			} else {
//
// 				var mplist []string
// 				for pt, _ := range r[k].MappedProteinsTNP {
// 					mplist = append(mplist, pt)
// 				}
// 				sort.Strings(mplist)
//
// 				for _, pt := range mplist {
// 					if r[k].MappedProteinsTNP[pt] >= topCount {
// 						topCount = r[k].MappedProteinsTNP[pt]
// 						topPT = pt
// 						//break
// 					}
// 				}
//
// 				razorPair[k] = topPT
//
// 			}
// 		}
//
// 	}
//
// 	for _, k := range rList {
// 		pt, ok := razorPair[k]
// 		if ok {
// 			razor := r[k]
// 			razor.MappedProtein = pt
// 			r[k] = razor
// 		}
// 	}
//
// 	spew.Dump(r)
// 	os.Exit(1)
//
// 	for i := range p.Groups {
// 		for j := range p.Groups[i].Proteins {
// 			for k := range p.Groups[i].Proteins[j].PeptideIons {
// 				v, ok := r[string(p.Groups[i].Proteins[j].PeptideIons[k].PeptideSequence)]
// 				if ok {
// 					if p.Groups[i].Proteins[j].ProteinName == v.MappedProtein {
// 						p.Groups[i].Proteins[j].PeptideIons[k].Razor = 1
// 						p.Groups[i].Proteins[j].HasRazor = true
// 					}
// 				}
// 			}
// 		}
// 	}
//
// 	// 	// mark as razor all peptides in the reference map
// 	for i := range p.Groups {
// 		for j := range p.Groups[i].Proteins {
// 			var r float64
// 			for k := range p.Groups[i].Proteins[j].PeptideIons {
// 				if p.Groups[i].Proteins[j].PeptideIons[k].Razor == 1 || p.Groups[i].Proteins[j].PeptideIons[k].IsUnique {
// 					if p.Groups[i].Proteins[j].PeptideIons[k].InitialProbability > r {
// 						r = p.Groups[i].Proteins[j].PeptideIons[k].InitialProbability
// 					}
// 				}
// 			}
// 			p.Groups[i].Proteins[j].TopPepProb = r
// 		}
// 	}
//
// 	return p, nil
// }

// RazorFilter classifies peptides as razor
func RazorFilter(p id.ProtXML) (id.ProtXML, error) {

	var r = make(map[string]RazorCandidate)
	var rList []string

	// for each peptide sequence, collapse all parent protein peptides from ions originated from the same sequence
	for _, i := range p.Groups {
		for _, j := range i.Proteins {
			for _, k := range j.PeptideIons {

				v, ok := r[k.PeptideSequence]
				if !ok {

					var rc RazorCandidate
					rc.Sequence = k.PeptideSequence
					rc.MappedProteinsW = make(map[string]float64)
					rc.MappedProteinsGW = make(map[string]float64)
					rc.MappedProteinsTNP = make(map[string]int)

					rc.MappedProteinsW[j.ProteinName] = k.Weight
					rc.MappedProteinsGW[j.ProteinName] = k.GroupWeight
					rc.MappedProteinsTNP[j.ProteinName] = j.TotalNumberPeptides

					for _, i := range j.IndistinguishableProtein {
						rc.MappedProteinsW[i] = -1
						rc.MappedProteinsGW[i] = -1
						rc.MappedProteinsTNP[i] = -1
					}

					for _, i := range k.PeptideParentProtein {
						rc.MappedProteinsW[i] = -1
						rc.MappedProteinsGW[i] = -1
						rc.MappedProteinsTNP[i] = -1
					}

					r[k.PeptideSequence] = rc

				} else {
					var c = v

					// doing like this will allow proteins that map to shared peptidesto be considered
					c.MappedProteinsW[j.ProteinName] = k.Weight
					c.MappedProteinsGW[j.ProteinName] = k.GroupWeight
					c.MappedProteinsTNP[j.ProteinName] = j.TotalNumberPeptides

					// if c.MappedProteinsW[j.ProteinName] == -1 {
					// 	c.MappedProteinsW[j.ProteinName] = k.Weight
					// }
					//
					// if c.MappedProteinsGW[j.ProteinName] == -1 {
					// 	c.MappedProteinsGW[j.ProteinName] = k.GroupWeight
					// }
					//
					// if c.MappedProteinsTNP[j.ProteinName] == -1 {
					// 	c.MappedProteinsTNP[j.ProteinName] = j.TotalNumberPeptides
					// }

					r[k.PeptideSequence] = c

				}
			}
		}
	}

	// spew.Dump(r)
	// os.Exit(1)

	// this will make the assignment more deterministic
	for k := range r {
		rList = append(rList, k)
	}
	sort.Strings(rList)

	var razorPair = make(map[string]string)

	// get the best protein candidate for each pepetide sequence and make the razor pair
	for _, k := range rList {
		// 1st pass: mark all cases with weight > 0.5
		for pt, w := range r[k].MappedProteinsW {
			if w > 0.5 {
				razorPair[k] = pt
			}
		}
	}

	// 2nd pass: mark all cases with highest group weight in the list
	for _, k := range rList {

		_, ok := razorPair[k]
		if !ok {

			var topPT string
			var topCount int
			var topGW float64
			var topGWMap = make(map[float64]uint8)

			if len(r[k].MappedProteinsGW) == 1 {

				for pt := range r[k].MappedProteinsGW {
					razorPair[k] = pt
				}

			} else if len(r[k].MappedProteinsGW) > 1 {

				for pt, tnp := range r[k].MappedProteinsGW {
					if tnp >= topGW {
						topGW = tnp
						topPT = pt
						topGWMap[topGW]++
					}
				}

				var tie bool
				if topGWMap[topGW] >= 2 {
					tie = true
				}

				if tie == false {
					razorPair[k] = topPT

				} else {

					var mplist []string
					for pt, _ := range r[k].MappedProteinsTNP {
						mplist = append(mplist, pt)
					}
					sort.Strings(mplist)

					for _, pt := range mplist {
						if r[k].MappedProteinsTNP[pt] >= topCount {
							topCount = r[k].MappedProteinsTNP[pt]
							topPT = pt
							//break
						}
					}

					razorPair[k] = topPT
				}

			}
		}
	}

	for _, k := range rList {
		pt, ok := razorPair[k]
		if ok {
			razor := r[k]
			razor.MappedProtein = pt
			r[k] = razor
		}
	}

	// spew.Dump(r)
	// os.Exit(1)

	for i := range p.Groups {
		for j := range p.Groups[i].Proteins {
			for k := range p.Groups[i].Proteins[j].PeptideIons {
				v, ok := r[string(p.Groups[i].Proteins[j].PeptideIons[k].PeptideSequence)]
				if ok {
					if p.Groups[i].Proteins[j].ProteinName == v.MappedProtein {
						p.Groups[i].Proteins[j].PeptideIons[k].Razor = 1
						p.Groups[i].Proteins[j].HasRazor = true
					}
				}
			}
		}
	}

	// 	// mark as razor all peptides in the reference map
	for i := range p.Groups {
		for j := range p.Groups[i].Proteins {
			var r float64
			for k := range p.Groups[i].Proteins[j].PeptideIons {
				if p.Groups[i].Proteins[j].PeptideIons[k].Razor == 1 || p.Groups[i].Proteins[j].PeptideIons[k].IsUnique {
					if p.Groups[i].Proteins[j].PeptideIons[k].InitialProbability > r {
						r = p.Groups[i].Proteins[j].PeptideIons[k].InitialProbability
					}
				}
			}
			p.Groups[i].Proteins[j].TopPepProb = r
		}
	}

	return p, nil
}

// ProtXMLFilter filters the protein list under a specific fdr
func ProtXMLFilter(p id.ProtXML, targetFDR, pepProb, protProb float64, isPicked, isRazor bool) (id.ProtIDList, error) {

	//var proteinIDs ProtIDList
	var list id.ProtIDList
	var targets float64
	var decoys float64
	var calcFDR float64
	var minProb float64 = 10
	var err error

	// collect all proteins from every group
	for i := range p.Groups {
		for j := range p.Groups[i].Proteins {

			if isRazor == true {

				if isPicked == true {
					if p.Groups[i].Proteins[j].Picked == 1 && p.Groups[i].Proteins[j].HasRazor == true {
						list = append(list, p.Groups[i].Proteins[j])
					}
				} else {
					if p.Groups[i].Proteins[j].HasRazor == true {
						list = append(list, p.Groups[i].Proteins[j])
					}
				}

			} else {

				if isPicked == true {
					if p.Groups[i].Proteins[j].Probability >= protProb && p.Groups[i].Proteins[j].Picked == 1 {
						list = append(list, p.Groups[i].Proteins[j])
					}

				} else {
					if p.Groups[i].Proteins[j].TopPepProb >= pepProb && p.Groups[i].Proteins[j].Probability >= protProb {
						list = append(list, p.Groups[i].Proteins[j])
					}
				}

			}

		}
	}

	for i := range list {
		if cla.IsDecoyProtein(list[i], p.DecoyTag) {
			decoys++
		} else {
			targets++
		}
	}

	sort.Sort(&list)

	// from botttom to top, classify every protein block with a given fdr score
	// the score is only calculates to the first (last) protein in each block
	// proteins with the same score, get the same fdr value.
	var scoreMap = make(map[float64]float64)
	for j := (len(list) - 1); j >= 0; j-- {
		_, ok := scoreMap[list[j].TopPepProb]
		if !ok {
			scoreMap[list[j].TopPepProb] = (decoys / targets)
		}

		if cla.IsDecoyProtein(list[j], p.DecoyTag) {
			decoys--
		} else {
			targets--
		}
	}

	var keys []float64
	for k := range scoreMap {
		keys = append(keys, k)
	}

	sort.Sort(sort.Reverse(sort.Float64Slice(keys)))

	var curProb = 10.0
	var curScore = 0.0
	var probArray []float64
	var probList = make(map[float64]uint8)

	for i := range keys {

		// for inspections
		//f := uti.Round(scoreMap[keys[i]]*100, 5, 2)
		//fmt.Println(keys[i], "\t", scoreMap[keys[i]], "\t", uti.ToFixed(scoreMap[keys[i]], 4), "\t", f)
		//fmt.Println(keys[i], "\t", scoreMap[keys[i]], "\t", uti.ToFixed(scoreMap[keys[i]], 4), "\t", f, "\t", targetFDR)

		probArray = append(probArray, keys[i])

		if uti.ToFixed(scoreMap[keys[i]], 4) <= targetFDR {
			probList[keys[i]] = 0
			minProb = keys[i]
			calcFDR = scoreMap[keys[i]]
			if keys[i] < curProb {
				curProb = keys[i]
			}
			if scoreMap[keys[i]] > curScore {
				curScore = scoreMap[keys[i]]
			}
		}

	}

	if curProb == 10 {
		msgProb := fmt.Sprintf("The protein FDR filter didn't reached the desired threshold of %.4f, try a higher threshold using the --prot parameter", targetFDR)
		err = errors.New(msgProb)
	}

	fmtScore := uti.ToFixed(curScore, 4)

	// for inspections
	//fmt.Println("curscore:", curScore, "\t", "fmtScore:", fmtScore, "\t", "targetfdr:", targetFDR)

	if curScore < targetFDR && fmtScore != targetFDR && probArray[len(probArray)-1] != curProb {

		for i := 0; i <= len(probArray); i++ {

			if probArray[i] == curProb {
				probList[probArray[i+1]] = 0
				minProb = probArray[i+1]
				calcFDR = scoreMap[probArray[i+1]]
				if probArray[i+1] < curProb {
					curProb = probArray[i+1]
				}
				if scoreMap[probArray[i+1]] > curScore {
					curScore = scoreMap[probArray[i+1]]
				}
				break
			}

		}

	}

	// for inspections
	//fmt.Println("curscore:", curScore, "\t", "fmtScore:", fmtScore, "\t", "targetfdr:", targetFDR)

	var cleanlist id.ProtIDList
	for i := range list {
		_, ok := probList[list[i].TopPepProb]
		if ok {
			cleanlist = append(cleanlist, list[i])
			if cla.IsDecoyProtein(list[i], p.DecoyTag) {
				decoys++
			} else {
				targets++
			}
		}
	}

	msg := fmt.Sprintf("Converged to %.2f %% FDR with %0.f Proteins", (calcFDR * 100), targets)
	logrus.WithFields(logrus.Fields{
		"decoy":     decoys,
		"total":     (targets + decoys),
		"threshold": minProb,
	}).Info(msg)

	return cleanlist, err
}

// sequentialFDRControl estimates FDR levels by applying a second filter where all
// proteins from the protein filtered list are matched against filtered PSMs
func sequentialFDRControl(pep id.PepIDList, pro id.ProtIDList, psm, peptide, ion float64, decoyTag string) error {

	extPep := extractPSMfromPepXML(pep, pro)

	// organize enties by score (probability or expectation)
	sort.Sort(extPep)

	uniqPsms := getUniquePSMs(extPep)
	uniqPeps := getUniquePeptides(extPep)
	uniqIons := getUniquePeptideIons(extPep)

	logrus.WithFields(logrus.Fields{
		"psms":     len(uniqPsms),
		"peptides": len(uniqPeps),
		"ions":     len(uniqIons),
	}).Info("Applying sequential FDR estimation")

	filteredPSM, err := pepXMLFDRFilter(uniqPsms, psm, "PSM", decoyTag)
	if err != nil {
		logrus.Fatal(err)
	}
	filteredPSM.Serialize("psm")

	filteredPeptides, err := pepXMLFDRFilter(uniqPeps, peptide, "Peptide", decoyTag)
	if err != nil {
		logrus.Fatal(err)
	}
	filteredPeptides.Serialize("pep")

	filteredIons, err := pepXMLFDRFilter(uniqIons, ion, "Ion", decoyTag)
	if err != nil {
		logrus.Fatal(err)
	}
	filteredIons.Serialize("ion")

	return nil
}

// twoDFDRFilter estimates FDR levels by applying a second filter by regenerating
// a protein list with decoys from protXML and pepXML.
func twoDFDRFilter(pep id.PepIDList, pro id.ProtIDList, psm, peptide, ion float64, decoyTag string) error {

	// filter protein list at given FDR level and regenerate protein list by adding pairing decoys
	//logrus.Info("Creating mirror image from filtered protein list")
	mirrorProteinList := mirrorProteinList(pro, decoyTag)

	// get new protein list profile
	//logrus.Info(protxml.ProteinProfileWithList(mirrorProteinList, pa.Tag, pa.Con))
	t, d, _ := proteinProfileWithList(mirrorProteinList, decoyTag)
	logrus.WithFields(logrus.Fields{
		"target": t,
		"decoy":  d,
	}).Info("2D FDR estimation: Protein mirror image")

	// get PSM from the original pepXML using protein REGENERATED protein list, using protein names
	extPep := extractPSMfromPepXML(pep, mirrorProteinList)

	// organize enties by score (probability or expectation)
	sort.Sort(extPep)

	uniqPsms := getUniquePSMs(extPep)
	uniqPeps := getUniquePeptides(extPep)
	uniqIons := getUniquePeptideIons(extPep)

	logrus.WithFields(logrus.Fields{
		"psms":     len(uniqPsms),
		"peptides": len(uniqPeps),
		"ions":     len(uniqIons),
	}).Info("Second filtering results")

	filteredPSM, err := pepXMLFDRFilter(uniqPsms, psm, "PSM", decoyTag)
	if err != nil {
		logrus.Fatal(err)
	}
	filteredPSM.Serialize("psm")

	filteredPeptides, err := pepXMLFDRFilter(uniqPeps, peptide, "Peptide", decoyTag)
	if err != nil {
		logrus.Fatal(err)
	}
	filteredPeptides.Serialize("pep")

	filteredIons, err := pepXMLFDRFilter(uniqIons, ion, "Ion", decoyTag)
	if err != nil {
		logrus.Fatal(err)
	}
	filteredIons.Serialize("ion")

	return nil
}

// extractPSMfromPepXML retrieves all psm from protxml that maps into pepxml files
// using protein names from <protein> and <alternative_proteins> tags
func extractPSMfromPepXML(peplist id.PepIDList, pro id.ProtIDList) id.PepIDList {

	var protmap = make(map[string]uint16)
	var filterMap = make(map[string]id.PeptideIdentification)
	var output id.PepIDList

	// get all protein names from protxml
	for _, i := range pro {
		protmap[string(i.ProteinName)] = 0
	}

	for _, i := range peplist {

		_, ptTag := protmap[string(i.Protein)]
		if ptTag {
			filterMap[string(i.Spectrum)] = i
			protmap[string(i.Protein)]++
		} else {
			for _, j := range i.AlternativeProteins {
				_, altTag := protmap[j]
				if altTag {
					filterMap[string(i.Spectrum)] = i
					protmap[string(j)]++
				}
			}
		}

	}

	// // match protein names to <protein> tag on pepxml
	// for j := range peplist {
	// 	_, ok := protmap[string(peplist[j].Protein)]
	// 	if ok {
	// 		filterMap[string(peplist[j].Spectrum)] = peplist[j]
	// 	}
	// }
	//
	// // match protein names to <alternative_proteins> tag on pepxml
	// for m := range peplist {
	// 	for n := range peplist[m].AlternativeProteins {
	// 		_, ok := protmap[peplist[m].AlternativeProteins[n]]
	// 		if ok {
	// 			filterMap[string(peplist[m].Spectrum)] = peplist[m]
	// 		}
	// 	}
	// }

	for _, v := range filterMap {
		output = append(output, v)
	}

	return output
}

// mirrorProteinList takes a filtered list and regenerate the correspondedn decoys
func mirrorProteinList(p id.ProtIDList, decoyTag string) id.ProtIDList {

	var targets = make(map[string]uint8)
	var decoys = make(map[string]uint8)

	// get filtered list
	var list id.ProtIDList
	for _, i := range p {
		if !cla.IsDecoyProtein(i, decoyTag) {
			list = append(list, i)
		}
	}

	// get the list of identified taget proteins
	for _, i := range p {
		if cla.IsDecoy(i.ProteinName, decoyTag) {
			decoys[i.ProteinName] = 0
		} else {
			targets[i.ProteinName] = 0
		}
	}

	// collect all original protein ids in case we need to put them on mirror list
	var refMap = make(map[string]id.ProteinIdentification)
	for _, i := range p {
		refMap[i.ProteinName] = i
	}

	// add decoys correspondent to the given targets.
	// first check if the oposite list doesn't have an entry already.
	// if not, search for the mirror entry on the original list, if found
	// move it to the mirror list, otherwise add fake entry.
	for _, k := range list {
		decoy := decoyTag + k.ProteinName
		v, ok := refMap[decoy]
		if ok {
			list = append(list, v)
		} else {
			var pt id.ProteinIdentification
			pt.ProteinName = decoy
			list = append(list, pt)
		}
	}

	return list
}

// proteinProfileWithList ...
func proteinProfileWithList(list []id.ProteinIdentification, decoyTag string) (t, d int, err error) {

	for i := range list {
		if cla.IsDecoyProtein(list[i], decoyTag) {
			d++
		} else {
			t++
		}
	}

	return t, d, err
}
