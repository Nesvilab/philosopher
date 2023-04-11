package fil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/Nesvilab/philosopher/lib/cla"
	"github.com/Nesvilab/philosopher/lib/id"
	"github.com/Nesvilab/philosopher/lib/inf"
	"github.com/Nesvilab/philosopher/lib/met"
	"github.com/Nesvilab/philosopher/lib/mod"
	"github.com/Nesvilab/philosopher/lib/rep"
	"github.com/Nesvilab/philosopher/lib/sys"

	"github.com/sirupsen/logrus"
)

// Run executes the Filter processing
func Run(f met.Data) met.Data {

	e := rep.New()
	var pep id.PepIDList
	var pro id.ProtIDList
	var dbBin string

	if len(f.Filter.ProBin) > 0 {

		f.Filter.Razor = true

		if _, err := os.Stat(f.Filter.ProBin); err == nil {

			p := fmt.Sprintf("%s%s.meta/protxml.bin", f.Filter.ProBin, string(filepath.Separator))
			r := fmt.Sprintf("%s%s.meta/razor.bin", f.Filter.ProBin, string(filepath.Separator))
			dbBin = fmt.Sprintf("%s%s", f.Filter.ProBin, string(filepath.Separator))

			logrus.Info("Fetching protein inference from ", f.Filter.ProBin)

			rdest := fmt.Sprintf("%s%s.meta%sprotxml.bin", f.Home, string(filepath.Separator), string(filepath.Separator))
			sys.CopyFile(p, rdest)

			rdest = fmt.Sprintf("%s%s.meta%srazor.bin", f.Home, string(filepath.Separator), string(filepath.Separator))
			sys.CopyFile(r, rdest)

		} else if errors.Is(err, os.ErrNotExist) {

			logrus.Warn("protein inference not found: ", f.Filter.ProBin)

			f.Filter.ProBin = ""
		}

	}

	logrus.Info("Processing peptide identification files")

	// if no method is selected, force the 2D to be default
	if len(f.Filter.Pox) > 0 && !f.Filter.TwoD && !f.Filter.Seq {
		f.Filter.TwoD = true
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	var protXML id.ProtXML
	go func() {
		defer wg.Done()
		if len(f.Filter.Pox) > 0 {
			protXML = ReadProtXMLInput(f.Filter.Pox, f.Filter.Tag, f.Filter.Weight)
		}
	}()
	pepid, searchEngine := id.ReadPepXMLInput(f.Filter.Pex, f.Filter.Tag, f.Temp, f.Filter.Model)
	wg.Wait()

	f.SearchEngine = searchEngine

	psmT, pepT, ionT := processPeptideIdentifications(pepid, f.Filter.Tag, f.Filter.Mods, f.Filter.PsmFDR, f.Filter.PepFDR, f.Filter.IonFDR, f.Filter.Delta, f.Filter.Group)
	_ = psmT
	_ = pepT
	_ = ionT
	if _, err := os.Stat(sys.ProBin()); err == nil {

		pro.Restore()

	} else if len(f.Filter.Pox) > 0 && !strings.EqualFold(f.Filter.Pox, "combined") {

		//protXML := ReadProtXMLInput(f.Filter.Pox, f.Filter.Tag, f.Filter.Weight)

		ProcessProteinIdentifications(protXML, f.Filter.PtFDR, f.Filter.PepFDR, f.Filter.ProtProb, f.Filter.Picked, f.Filter.Razor, false, f.Filter.Tag)
		pro.Restore()

	} else {

		if f.Filter.Inference {

			var filteredPSM id.PepIDList
			filteredPSM.Restore("psm")

			pepid, razorMap, coverMap := inf.ProteinInference(filteredPSM)
			filteredPSM = nil

			pepid.Serialize("psm")
			pepid.Serialize("pep")
			pepid.Serialize("ion")

			processProteinInferenceIdentifications(pepid, razorMap, coverMap, f.Filter.PtFDR, f.Filter.PepFDR, f.Filter.ProtProb, f.Filter.Picked, f.Filter.Tag)
		}
	}

	var pepxml id.PepXML
	pepxml.Restore()

	// restoring for the modifications
	e.Mods = pepxml.Modifications
	e.AssembleSearchParameters(pepxml.SearchParameters)

	if f.Filter.Seq {
		// sequential analysis
		// filtered psm list and filtered prot list
		pep.Restore("psm")
		sequentialFDRControl(pep, pro, f.Filter.PsmFDR, f.Filter.PepFDR, f.Filter.IonFDR, f.Filter.Tag)
		pep = nil

	} else if f.Filter.TwoD {

		// two-dimensional analysis
		// complete pep list and filtered mirror-image prot list
		twoDFDRFilter(pepxml.PeptideIdentification, pro, f.Filter.PsmFDR, f.Filter.PepFDR, f.Filter.IonFDR, f.Filter.Tag)

	}

	os.RemoveAll(sys.PepxmlBin())

	logrus.Info("Post processing identifications")

	var psm id.PepIDList
	psm.Restore("psm")
	e.AssemblePSMReport(psm, f.Filter.Tag)
	psm = nil

	var ion id.PepIDList
	ion.Restore("ion")
	e.AssembleIonReport(ion, f.Filter.Tag)
	ion = nil

	// evaluate modifications in data set
	if f.Filter.Mapmods {
		logrus.Info("Mapping modifications")
		//should include observed mods into mapping?
		e.MapMods()
	}

	var pept id.PepIDList
	pept.Restore("pep")
	e.AssemblePeptideReport(pept, f.Filter.Tag)
	pept = nil

	// evaluate modifications in data set
	if f.Filter.Mapmods {
		e.UpdatePeptideModCount()
	}

	// Apply the razor assignment to all data
	if f.Filter.Razor || len(f.Filter.ProBin) > 0 {
		e.ApplyRazorAssignment(f.Filter.Tag)
	}

	logrus.Info("Assigning protein identifications to layers")

	// object d is for reuising databasee paths
	e.UpdateLayerswithDatabase(dbBin, f.Filter.Tag)

	if len(f.Filter.Pox) > 0 || f.Filter.Inference {

		logrus.Info("Processing protein inference")
		pro.Restore()

		e.AssembleProteinReport(pro, f.Filter.Weight, dbBin, f.Filter.Tag)
		pro = nil

		logrus.Info("Synchronizing PSMs and proteins")

		e.SyncPSMToProteins(f.Filter.Tag)

		e.UpdateNumberOfEnzymaticTermini(f.Filter.Tag)

		e.CalculateProteinCoverage()
	}

	e = e.SyncPSMToPeptides(f.Filter.Tag)
	e = e.SyncPSMToPeptideIons(f.Filter.Tag)

	var countPSM, countPep, countIon, coutProtein int
	for _, i := range e.PSM {
		if !i.IsDecoy {
			countPSM++
		}
	}

	for _, i := range e.Peptides {
		if !i.IsDecoy {
			countPep++
		}
	}

	for _, i := range e.Ions {
		if !i.IsDecoy {
			countIon++
		}
	}

	for _, i := range e.Proteins {
		if !i.IsDecoy {
			coutProtein++
		}
	}

	logrus.WithFields(logrus.Fields{
		"psms":     countPSM,
		"peptides": countPep,
		"ions":     countIon,
		"proteins": coutProtein,
	}).Info("Final report numbers after FDR filtering, and post-processing")
	logrus.Info("Saving")
	e.SerializeGranular()

	return f
}

// processPeptideIdentifications reads and process pepXML
func processPeptideIdentifications(p id.PepIDListPtrs, decoyTag, mods string, psm, peptide, ion float64, delta, class bool) (float64, float64, float64) {

	// report charge profile
	var t, d int

	t, d = chargeProfile(p, 1, decoyTag)
	logrus.WithFields(logrus.Fields{
		"target": t,
		"decoy":  d,
	}).Info("1+ Charge profile")

	t, d = chargeProfile(p, 2, decoyTag)
	logrus.WithFields(logrus.Fields{
		"target": t,
		"decoy":  d,
	}).Info("2+ Charge profile")

	t, d = chargeProfile(p, 3, decoyTag)
	logrus.WithFields(logrus.Fields{
		"target": t,
		"decoy":  d,
	}).Info("3+ Charge profile")

	t, d = chargeProfile(p, 4, decoyTag)
	logrus.WithFields(logrus.Fields{
		"target": t,
		"decoy":  d,
	}).Info("4+ Charge profile")

	t, d = chargeProfile(p, 5, decoyTag)
	logrus.WithFields(logrus.Fields{
		"target": t,
		"decoy":  d,
	}).Info("5+ Charge profile")

	t, d = chargeProfile(p, 6, decoyTag)
	logrus.WithFields(logrus.Fields{
		"target": t,
		"decoy":  d,
	}).Info("6+ Charge profile")

	uniqPsms := GetUniquePSMs(p)
	uniqPeps := GetUniquePeptides(p)
	uniqIons := getUniquePeptideIons(p)

	logrus.WithFields(logrus.Fields{
		"psms":     len(p),
		"peptides": len(uniqPeps),
		"ions":     len(uniqIons),
	}).Info("Database search results")

	filteredPSM, psmThreshold := PepXMLFDRFilter(uniqPsms, psm, "PSM", decoyTag, "")
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() { defer wg.Done(); filteredPSM.Serialize("psm") }()

	filteredPeptides, peptideThreshold := PepXMLFDRFilter(uniqPeps, peptide, "Peptide", decoyTag, "")
	go func() { defer wg.Done(); filteredPeptides.Serialize("pep") }()

	filteredIons, ionThreshold := PepXMLFDRFilter(uniqIons, ion, "Ion", decoyTag, "")
	go func() { defer wg.Done(); filteredIons.Serialize("ion") }()
	wg.Wait()

	// sug-group FDR filtering
	if len(mods) > 0 {
		ptmBasedPSMFiltering(uniqPsms, psm, decoyTag, mods)
	}

	if delta {
		deltaMassBasedPSMFiltering(uniqPsms, psm, decoyTag)
	}

	if class {
		classBasedPSMFiltering(uniqPsms, psm, decoyTag)
	}

	return psmThreshold, peptideThreshold, ionThreshold
}

// classBasedPSMFiltering applies FDR filtering on PSMs based on the class
func classBasedPSMFiltering(uniqPsms map[string]id.PepIDListPtrs, targetFDR float64, decoyTag string) {

	logrus.Info("Separating PSMs based on group")

	var classes []string
	classMap := make(map[string][]id.PepIDListPtrs)

	for _, v := range uniqPsms {
		classMap[v[0].Class] = append(classMap[v[0].Class], v)
	}

	uniqueKeys := make(map[string]bool)
	for key := range classMap {
		if !uniqueKeys[key] {
			uniqueKeys[key] = true
			classes = append(classes, key)
		}
	}

	sort.Strings(classes)

	var combinedFiltered id.PepIDListPtrs

	for i := 0; i < len(classes); i++ {

		psms := make(map[string]id.PepIDListPtrs)

		for _, v := range classMap[classes[i]] {
			psms[v[0].Spectrum] = v
		}

		logrus.Info("Filtering group ", classes[i])
		filteredPSMs, _ := PepXMLFDRFilter(psms, targetFDR, "PSM", decoyTag, "")

		combinedFiltered = append(combinedFiltered, filteredPSMs...)

	}

	combinedFiltered.Serialize("psm")
}

// func classBasedPSMFiltering(uniqPsms map[string]id.PepIDListPtrs, targetFDR float64, decoyTag string) {

// 	logrus.Info("Separating PSMs based on class")

// 	var classes []string
// 	classMap := make(map[string][]id.PepIDListPtrs)

// 	for _, v := range uniqPsms {
// 		classMap[v[0].Class] = append(classMap[v[0].Class], v)
// 	}

// 	uniqueKeys := make(map[string]bool)
// 	for key := range classMap {
// 		if !uniqueKeys[key] {
// 			uniqueKeys[key] = true
// 			classes = append(classes, key)
// 		}
// 	}

// 	sort.Strings(classes)

// 	var combinedFiltered *id.PepIDListPtrs = new(id.PepIDListPtrs)

// 	var wg sync.WaitGroup
// 	results := make(chan id.PepIDListPtrs, len(classes))

// 	for i := 0; i < len(classes); i++ {
// 		wg.Add(1)

// 		go func(class string) {
// 			defer wg.Done()

// 			psms := make(map[string]id.PepIDListPtrs)

// 			for _, v := range classMap[class] {
// 				psms[v[0].Spectrum] = v
// 			}

// 			logrus.Info("Filtering class ", class)
// 			filteredPSMs, _ := PepXMLFDRFilter(psms, targetFDR, "PSM", decoyTag, "")

// 			results <- filteredPSMs
// 		}(classes[i])
// 	}

// 	go func() {
// 		wg.Wait()
// 		close(results)
// 	}()

// 	for filteredPSMs := range results {
// 		*combinedFiltered = append(*combinedFiltered, filteredPSMs...)
// 	}

// 	combinedFiltered.Serialize("psm")
// }

func deltaMassBasedPSMFiltering(uniqPsms map[string]id.PepIDListPtrs, targetFDR float64, decoyTag string) {

	logrus.Info("Separating PSMs based on the delta mass profile")

	// unmodified: delta mass < 1000
	unModPSMs := make(map[string]id.PepIDListPtrs)

	// common: only the most common PTMs; > 1000 & < 100000
	commonModPSMs := make(map[string]id.PepIDListPtrs)

	// glyco: all the remaining PTMs including glyco
	glycoModPSMs := make(map[string]id.PepIDListPtrs)

	for k, v := range uniqPsms {

		var glyco, common bool

		if v[0].Massdiff > 145 {
			glyco = true
		} else if v[0].Massdiff >= 3.5 && v[0].Massdiff <= 145 {
			common = true
		}

		if glyco && !common {
			glycoModPSMs[k] = v
		} else if !glyco && common {
			commonModPSMs[k] = v
		} else {
			unModPSMs[k] = v
		}

	}

	logrus.Info("Filtering unmodified PSMs")
	filteredUnmodPSM, _ := PepXMLFDRFilter(unModPSMs, targetFDR, "PSM", decoyTag, "")

	logrus.Info("Filtering commonly modified PSMs")
	filteredDefinedPSM, _ := PepXMLFDRFilter(commonModPSMs, targetFDR, "PSM", decoyTag, "")

	logrus.Info("Filtering glyco-modified PSMs")
	filteredAllPSM, _ := PepXMLFDRFilter(glycoModPSMs, targetFDR, "PSM", decoyTag, "")

	var combinedFiltered id.PepIDListPtrs

	combinedFiltered = append(combinedFiltered, filteredUnmodPSM...)

	combinedFiltered = append(combinedFiltered, filteredDefinedPSM...)

	combinedFiltered = append(combinedFiltered, filteredAllPSM...)

	combinedFiltered.Serialize("psm")

}

func ptmBasedPSMFiltering(uniqPsms map[string]id.PepIDListPtrs, targetFDR float64, decoyTag, mods string) {

	logrus.Info("Separating PSMs based on the modification profile")

	// unmodified: no ptms
	unModPSMs := make(map[string]id.PepIDListPtrs)

	// defined: only the ptms defined, nothing else
	definedModPSMs := make(map[string]id.PepIDListPtrs)

	// other: one or more ptms that might include the combination of the defined + something else
	restModPSMs := make(map[string]id.PepIDListPtrs)

	modsMap := make(map[string]string)

	modsList := strings.Split(mods, ",")
	for _, i := range modsList {
		m := strings.Split(i, ":")
		modsMap[i] = m[0]
	}

	for k, v := range uniqPsms {

		var other, defined bool

		for _, i := range v[0].Modifications.IndexSlice {

			if i.Variable {

				var m string
				if i.AminoAcid == "N-term" {
					m = fmt.Sprintf("%s:%.4f", "n", i.MassDiff)
				} else {
					m = fmt.Sprintf("%s:%.4f", i.AminoAcid, i.MassDiff)
				}

				_, ok := modsMap[m]
				if ok {
					defined = true
				} else {
					other = true
				}

			}
		}

		if other && defined {
			restModPSMs[k] = v
		} else if other && !defined {
			restModPSMs[k] = v
		} else if !other && defined {
			definedModPSMs[k] = v
		} else {
			unModPSMs[k] = v
		}

	}

	logrus.Info("Filtering unmodified PSMs")
	filteredUnmodPSM, _ := PepXMLFDRFilter(unModPSMs, targetFDR, "PSM", decoyTag, "")

	logrus.Info("Filtering defined modified PSMs")
	filteredDefinedPSM, _ := PepXMLFDRFilter(definedModPSMs, targetFDR, "PSM", decoyTag, "")

	logrus.Info("Filtering all other PSMs")
	filteredAllPSM, _ := PepXMLFDRFilter(restModPSMs, targetFDR, "PSM", decoyTag, "X")

	var combinedFiltered id.PepIDListPtrs

	combinedFiltered = append(combinedFiltered, filteredUnmodPSM...)

	combinedFiltered = append(combinedFiltered, filteredDefinedPSM...)

	combinedFiltered = append(combinedFiltered, filteredAllPSM...)

	combinedFiltered.Serialize("psm")

}

// chargeProfile ...
func chargeProfile(p id.PepIDListPtrs, charge uint8, decoyTag string) (t, d int) {

	for _, i := range p {
		if i.AssumedCharge == charge {
			if strings.HasPrefix(i.Protein, decoyTag) {
				d++
			} else {
				t++
			}
		}
	}

	return t, d
}

// GetUniquePSMs selects only unique pepetide ions for the given data structure
func GetUniquePSMs(p id.PepIDListPtrs) map[string]id.PepIDListPtrs {
	uniqMap := make(map[string]id.PepIDListPtrs)

	for _, i := range p {
		uniqMap[i.SpectrumFileName().Str()] = append(uniqMap[i.SpectrumFileName().Str()], i)
	}
	return uniqMap
}

// getUniquePeptideIons selects only unique pepetide ions for the given data structure
func getUniquePeptideIons(p id.PepIDListPtrs) map[string]id.PepIDListPtrs {

	uniqMap := ExtractIonsFromPSMs(p)

	return uniqMap
}

// ExtractIonsFromPSMs takes a pepidlist and transforms into an ion map
func ExtractIonsFromPSMs(p id.PepIDListPtrs) map[string]id.PepIDListPtrs {

	uniqMap := make(map[string]id.PepIDListPtrs)

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

// GetUniquePeptides selects only unique pepetide for the given data structure
func GetUniquePeptides(p id.PepIDListPtrs) map[string]id.PepIDListPtrs {

	uniqMap := make(map[string]id.PepIDListPtrs)

	for _, i := range p {
		uniqMap[i.Peptide] = append(uniqMap[i.Peptide], i)
	}

	// organize id list by score
	for _, v := range uniqMap {
		sort.Sort(v)
	}

	return uniqMap
}

// ReadProtXMLInput reads one or more fies and organize the data into PSM list
func ReadProtXMLInput(xmlFile, decoyTag string, weight float64) id.ProtXML {

	var protXML id.ProtXML

	protXML.Read(xmlFile)

	protXML.DecoyTag = decoyTag

	protXML.MarkUniquePeptides(weight)

	protXML.PromoteProteinIDs()

	protXML.Serialize()

	return protXML
}

// ProcessProteinIdentifications checks if pickedFDR ar razor options should be applied to given data set, if they do,
// the inputed protXML data is processed before filtered.
func ProcessProteinIdentifications(p id.ProtXML, ptFDR, pepProb, protProb float64, isPicked, isRazor, isCombined bool, decoyTag string) string {

	var pid id.ProtIDList

	// tagget / decoy / threshold
	t, d := proteinProfile(p)
	logrus.WithFields(logrus.Fields{
		"target": t,
		"decoy":  d,
	}).Info("Protein inference results")

	// applies pickedFDR algorithm
	if isPicked {
		p = PickedFDR(p)
	}

	// applies razor algorithm
	if isRazor {
		p = RazorFilter(p)
	}

	// run the FDR filter for proteins
	pid = ProtXMLFilter(p, ptFDR, pepProb, protProb, isPicked, isRazor, decoyTag)

	// save results on meta folder
	if isCombined {
		proBin := pid.SerializeToTemp()
		return proBin
	} else {
		pid.Serialize()
	}

	return ""
}

// processProteinInferenceIdentifications checks if pickedFDR ar razor options should be applied to given data set, if they do,
// the inputed Philosopher inference data is processed before filtered.
func processProteinInferenceIdentifications(psm id.PepIDList, razorMap map[string]string, coverMap map[string]float64, ptFDR, pepProb, protProb float64, isPicked bool, decoyTag string) {

	var t int
	var d int
	var proXML id.ProtXML
	var proGrps id.GroupList
	var proteinList = make(map[string]id.ProteinIdentification)

	// build the ProtXML strct
	grpID := id.GroupIdentification{
		GroupNumber: 0,
		Probability: 1.00,
	}

	proGrps = append(proGrps, grpID)

	proXML.DecoyTag = decoyTag
	proXML.Groups = proGrps

	for _, i := range psm {
		_, ok := proteinList[i.Protein]
		if !ok {

			p := id.ProteinIdentification{
				GroupNumber:    0,
				GroupSiblingID: "a",
				ProteinName:    i.Protein,
				Picked:         0,
				HasRazor:       false,
			}

			proteinList[i.Protein] = p
		}
	}

	for i := range proteinList {
		if strings.HasPrefix(i, decoyTag) {
			d++
		} else {
			t++
		}
	}

	// add the razor / non-razor marked proteins
	var razorMarked = make(map[string]uint8)
	for _, i := range psm {

		pro := proteinList[i.Protein]
		razorProtein, ok := razorMap[i.Peptide]

		if ok && pro.ProteinName == razorProtein {

			pro.Length = 0
			pro.PercentCoverage = float32(coverMap[pro.ProteinName])
			pro.HasRazor = true

			if i.Probability > pro.Probability {
				pro.Probability = i.Probability
				pro.TopPepProb = i.Probability
			}

			razorMarked[pro.ProteinName] = 0

		} else {

			_, ok := razorMarked[pro.ProteinName]
			if !ok {
				pro.Length = 0
				pro.PercentCoverage = float32(coverMap[pro.ProteinName])
				pro.HasRazor = false

				if i.Probability > pro.Probability {
					pro.Probability = i.Probability
					pro.TopPepProb = i.Probability
				}

			}
		}

		proteinList[i.Protein] = pro
	}

	// add the ions
	//var addedIon = make(map[string]uint8)
	for _, i := range psm {

		pro := proteinList[i.Protein]
		razorProtein, ok := razorMap[i.Peptide]

		if ok && pro.ProteinName == razorProtein {

			pro.UniqueStrippedPeptides = append(pro.UniqueStrippedPeptides, i.Peptide)
			pro.TotalNumberPeptides++

			pep := id.PeptideIonIdentification{
				PeptideSequence:    i.Peptide,
				ModifiedPeptide:    i.ModifiedPeptide,
				Charge:             i.AssumedCharge,
				Weight:             1,
				GroupWeight:        0,
				CalcNeutralPepMass: i.CalcNeutralPepMass,
				Razor:              1,
			}

			for j := range i.AlternativeProteins {
				pep.PeptideParentProtein = append(pep.PeptideParentProtein, j)
				pro.IndistinguishableProtein = append(pro.IndistinguishableProtein, j)
			}

			if i.Probability > pep.InitialProbability {
				pep.InitialProbability = i.Probability
			}

			if len(i.AlternativeProteins) < 2 {
				pep.IsUnique = true
			} else {
				pep.IsUnique = false
			}

			pep.Modifications.Index = make(map[string]mod.Modification)
			for k, v := range i.Modifications.ToMap().Index {
				pep.Modifications.Index[k] = v
			}

			pro.HasRazor = true
			pro.PeptideIons = append(pro.PeptideIons, pep)

			proteinList[i.Protein] = pro

		} else {

			pro.UniqueStrippedPeptides = append(pro.UniqueStrippedPeptides, i.Peptide)
			pro.TotalNumberPeptides++

			pep := id.PeptideIonIdentification{
				PeptideSequence:    i.Peptide,
				ModifiedPeptide:    i.ModifiedPeptide,
				Charge:             i.AssumedCharge,
				Weight:             0,
				GroupWeight:        0,
				CalcNeutralPepMass: i.CalcNeutralPepMass,
				Razor:              0,
			}

			if i.Probability > pep.InitialProbability {
				pep.InitialProbability = i.Probability
			}

			if len(i.AlternativeProteins) < 2 {
				pep.IsUnique = true
			} else {
				pep.IsUnique = false
			}

			pep.Modifications.Index = make(map[string]mod.Modification)
			for k, v := range i.Modifications.ToMap().Index {
				pep.Modifications.Index[k] = v
			}

			pro.PeptideIons = append(pro.PeptideIons, pep)

			proteinList[i.Protein] = pro

		}
	}

	for _, i := range proteinList {
		proXML.Groups[0].Proteins = append(proXML.Groups[0].Proteins, i)
	}

	// tagget / decoy / threshold
	logrus.WithFields(logrus.Fields{
		"target": t,
		"decoy":  d,
	}).Info("Protein inference results")

	// run the FDR filter for proteins
	pid := ProtXMLFilter(proXML, ptFDR, pepProb, protProb, false, true, decoyTag)

	// save results on meta folder
	pid.Serialize()

}

// proteinProfile ...
func proteinProfile(p id.ProtXML) (t, d int) {

	for _, i := range p.Groups {
		for _, j := range i.Proteins {
			if cla.IsDecoyProtein(j, p.DecoyTag) {
				d++
			} else {
				t++
			}
		}
	}

	return t, d
}

// extractPSMfromPepXML retrieves all psm from protxml that maps into pepxml files
// using protein names from <protein> and <alternative_proteins> tags
func extractPSMfromPepXML(filter string, peplist id.PepIDList, pro id.ProtIDList) id.PepIDListPtrs {

	var protmap = make(map[string]struct{})
	var filterMap = make(map[id.SpectrumType]*id.PeptideIdentification)
	var output id.PepIDListPtrs

	if filter == "sequential" {

		// get all protein and peptide pairs from protxml
		for _, i := range pro {
			for _, j := range i.UniqueStrippedPeptides {
				key := fmt.Sprintf("%s#%s", i.ProteinName, j)
				protmap[key] = struct{}{}
			}
		}

		for idx, i := range peplist {

			key := fmt.Sprintf("%s#%s", i.Protein, i.Peptide)

			_, ok := protmap[key]
			if ok {
				filterMap[i.SpectrumFileName()] = &peplist[idx]
			} else {

				for j := range i.AlternativeProteins {
					key := fmt.Sprintf("%s#%s", j, i.Peptide)
					_, ap := protmap[key]
					if ap {
						filterMap[i.SpectrumFileName()] = &peplist[idx]
					}
				}

			}

		}

	} else if filter == "2d" {

		// get all protein names from protxml
		for _, i := range pro {
			protmap[string(i.ProteinName)] = struct{}{}
		}

		for idx, i := range peplist {
			_, ok := protmap[string(i.Protein)]
			if ok {
				filterMap[i.SpectrumFileName()] = &peplist[idx]
			} else {
				for j := range i.AlternativeProteins {
					_, ap := protmap[j]
					if ap {
						filterMap[i.SpectrumFileName()] = &peplist[idx]
					}
				}
			}
		}

	}

	for _, v := range filterMap {
		output = append(output, v)
	}

	return output
}

// proteinProfileWithList
func proteinProfileWithList(list []id.ProteinIdentification, decoyTag string) (t, d int) {

	for i := range list {
		if cla.IsDecoyProtein(list[i], decoyTag) {
			d++
		} else {
			t++
		}
	}
	return t, d
}
