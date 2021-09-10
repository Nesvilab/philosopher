package fil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"philosopher/lib/cla"
	"philosopher/lib/dat"
	"philosopher/lib/id"
	"philosopher/lib/inf"
	"philosopher/lib/met"
	"philosopher/lib/mod"
	"philosopher/lib/msg"
	"philosopher/lib/rep"
	"philosopher/lib/sys"

	"github.com/sirupsen/logrus"
)

// Run executes the Filter processing
func Run(f met.Data) met.Data {

	e := rep.New()
	var pepxml id.PepXML
	var pep id.PepIDList
	var pro id.ProtIDList

	if len(f.Filter.RazorBin) > 0 {

		f.Filter.Razor = true

		if _, err := os.Stat(f.Filter.RazorBin); os.IsNotExist(err) {
			logrus.Warn("razor peptides not found: ", f.Filter.RazorBin, ". Skipping razor assignment")
			f.Filter.RazorBin = ""
		} else {
			rdest := fmt.Sprintf("%s%s.meta%srazor.bin", f.Home, string(filepath.Separator), string(filepath.Separator))
			sys.CopyFile(f.Filter.RazorBin, rdest)
		}
	}

	// get the database tag from database command
	if len(f.Filter.Tag) == 0 {
		f.Filter.Tag = f.Database.Tag
	}

	logrus.Info("Processing peptide identification files")

	// if no method is selected, force the 2D to be default
	if len(f.Filter.Pox) > 0 && !f.Filter.TwoD && !f.Filter.Seq {
		f.Filter.TwoD = true
	}

	pepid, searchEngine := id.ReadPepXMLInput(f.Filter.Pex, f.Filter.Tag, f.Temp, f.Filter.Model)

	f.SearchEngine = searchEngine

	psmT, pepT, ionT := processPeptideIdentifications(pepid, f.Filter.Tag, f.Filter.Mods, f.Filter.PsmFDR, f.Filter.PepFDR, f.Filter.IonFDR)
	_ = psmT
	_ = pepT
	_ = ionT

	if _, err := os.Stat(sys.ProBin()); err == nil {

		pro.Restore()

	} else if len(f.Filter.Pox) > 0 && !strings.EqualFold(f.Filter.Pox, "combined") {

		protXML := ReadProtXMLInput(f.Filter.Pox, f.Filter.Tag, f.Filter.Weight)
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

	pepxml.Restore()

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

	var dtb dat.Base
	dtb.Restore()
	if len(dtb.Records) < 1 {
		msg.Custom(errors.New("database annotation not found, interrupting the processing"), "fatal")
	}

	if f.Filter.TwoD || f.Filter.Razor {
		var psm id.PepIDList
		psm.Restore("psm")
		psm = correctRazorAssignment(psm)
		psm.Serialize("psm")
		psm = nil

		var pep id.PepIDList
		pep.Restore("pep")
		pep = correctRazorAssignment(pep)
		pep.Serialize("pep")
		pep = nil

		var ion id.PepIDList
		ion.Restore("ion")
		ion = correctRazorAssignment(ion)
		ion.Serialize("ion")
		ion = nil
	}

	logrus.Info("Post processing identifications")

	// restoring for the modifications
	e.Mods = pepxml.Modifications
	e.AssembleSearchParameters(pepxml.SearchParameters)
	pepxml = id.PepXML{}
	os.RemoveAll(sys.PepxmlBin())

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

	logrus.Info("Assigning protein identifications to layers")
	e.UpdateLayerswithDatabase(f.Filter.Tag)

	// evaluate modifications in data set
	if f.Filter.Mapmods {
		e.UpdateIonModCount()
		e.UpdatePeptideModCount()
	}

	if f.Filter.Razor {

		var razor RazorMap = make(map[string]RazorCandidate)
		razor.Restore()

		for i := range e.PSM {

			for j := range e.PSM[i].MappedProteins {
				if strings.Contains(j, f.Filter.Tag) {
					delete(e.PSM[i].MappedProteins, j)
				}
			}

			v, ok := razor[e.PSM[i].Peptide]

			if ok {
				if len(v.MappedProtein) > 0 {
					if e.PSM[i].Protein != v.MappedProtein {
						e.PSM[i].MappedProteins[e.PSM[i].Protein]++
						delete(e.PSM[i].MappedProteins, v.MappedProtein)
						e.PSM[i].Protein = v.MappedProtein
					}
					delete(e.PSM[i].MappedProteins, v.MappedProtein)
				}

				e.PSM[i].IsURazor = true
			}

			if e.PSM[i].IsUnique {
				e.PSM[i].IsURazor = true
			}

			if len(e.PSM[i].MappedProteins) == 0 {
				e.PSM[i].IsURazor = true
				e.PSM[i].IsUnique = true

				e.PSM[i].MappedGenes = make(map[string]int)
			}

		}

		razor = nil
	}

	if len(f.Filter.Pox) > 0 || f.Filter.Inference {

		logrus.Info("Processing protein inference")
		pro.Restore()
		e.AssembleProteinReport(pro, f.Filter.Weight, f.Filter.Tag)
		pro = nil

		// Pushes the new ion status from the protein inferece to the other layers, the gene and protein ID
		// assignment gets corrected in the next function call (UpdateLayerswithDatabase)
		e.UpdateIonStatus(f.Filter.Tag)

		logrus.Info("Synchronizing PSMs and proteins")

		e = e.SyncPSMToProteins(f.Filter.Tag)

		e.UpdateNumberOfEnzymaticTermini()
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
	}).Info("Total report numbers after FDR filtering, and post-processing")

	logrus.Info("Saving")
	e.SerializeGranular()

	return f
}

// processPeptideIdentifications reads and process pepXML
func processPeptideIdentifications(p id.PepIDList, decoyTag, mods string, psm, peptide, ion float64) (float64, float64, float64) {

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

	filteredPSM, psmThreshold := PepXMLFDRFilter(uniqPsms, psm, "PSM", decoyTag)
	filteredPSM.Serialize("psm")

	filteredPeptides, peptideThreshold := PepXMLFDRFilter(uniqPeps, peptide, "Peptide", decoyTag)
	filteredPeptides.Serialize("pep")

	filteredIons, ionThreshold := PepXMLFDRFilter(uniqIons, ion, "Ion", decoyTag)
	filteredIons.Serialize("ion")

	// sug-group FDR filtering
	if len(mods) > 0 {
		ptmBasedPSMFiltering(uniqPsms, psm, decoyTag, mods)
	}

	return psmThreshold, peptideThreshold, ionThreshold
}

func ptmBasedPSMFiltering(uniqPsms map[string]id.PepIDList, targetFDR float64, decoyTag, mods string) {

	// unmodified = no ptms
	// defined = only the ptms defined, nothing else
	// remaining or all = one or more ptms that might include the combination of the defined + something else

	logrus.Info("Separating PSMs based on the modification profile")

	unModPSMs := make(map[string]id.PepIDList)
	definedModPSMs := make(map[string]id.PepIDList)
	restModPSMs := make(map[string]id.PepIDList)

	modsMap := make(map[string]string)

	modsList := strings.Split(mods, ",")
	for _, i := range modsList {
		m := strings.Split(i, ":")
		modsMap[i] = m[0]
	}

	exclusionList := make(map[string]uint8)
	psmsWithOtherPTMs := make(map[string]id.PepIDList)

	for k, v := range uniqPsms {

		if !strings.Contains(v[0].ModifiedPeptide, "[") || len(v[0].ModifiedPeptide) == 0 {

			unModPSMs[k] = v
			exclusionList[v[0].Spectrum] = 0

		} else {

			// if PSM contains other mods than the ones defined by the flag, mark them to be ignored
			for _, i := range v[0].Modifications.Index {
				if i.Variable == "Y" {
					m := fmt.Sprintf("%s:%.4f", i.AminoAcid, i.MassDiff)
					_, ok := modsMap[m]
					if !ok {
						psmsWithOtherPTMs[v[0].Spectrum] = v
					}
				}
			}

			// if PSM contains only the defined mod and the correct amino acid, teh add to defined category
			// and mark it for being excluded from rest
			for _, i := range v[0].Modifications.Index {
				m := fmt.Sprintf("%s:%.4f", i.AminoAcid, i.MassDiff)
				aa, ok1 := modsMap[m]
				_, ok2 := psmsWithOtherPTMs[v[0].Spectrum]

				if ok1 && !ok2 && aa == i.AminoAcid {
					definedModPSMs[k] = v
					exclusionList[v[0].Spectrum] = 0
				}
			}

		}

		_, ok := exclusionList[v[0].Spectrum]
		if !ok {
			restModPSMs[k] = v
		}
	}

	logrus.Info("Filtering unmodified PSMs")
	filteredUnmodPSM, _ := PepXMLFDRFilter(unModPSMs, targetFDR, "PSM", decoyTag)

	logrus.Info("Filtering defined modified PSMs")
	filteredDefinedPSM, _ := PepXMLFDRFilter(definedModPSMs, targetFDR, "PSM", decoyTag)

	logrus.Info("Filtering all modified PSMs")
	filteredAllPSM, _ := PepXMLFDRFilter(restModPSMs, targetFDR, "PSM", decoyTag)

	var combinedFiltered id.PepIDList

	combinedFiltered = append(combinedFiltered, filteredUnmodPSM...)

	combinedFiltered = append(combinedFiltered, filteredDefinedPSM...)

	combinedFiltered = append(combinedFiltered, filteredAllPSM...)

	combinedFiltered.Serialize("psm")

}

// chargeProfile ...
func chargeProfile(p id.PepIDList, charge uint8, decoyTag string) (t, d int) {

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

//GetUniquePSMs selects only unique pepetide ions for the given data structure
func GetUniquePSMs(p id.PepIDList) map[string]id.PepIDList {

	uniqMap := make(map[string]id.PepIDList)

	for _, i := range p {
		uniqMap[i.Spectrum] = append(uniqMap[i.Spectrum], i)
	}

	return uniqMap
}

//getUniquePeptideIons selects only unique pepetide ions for the given data structure
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

// GetUniquePeptides selects only unique pepetide for the given data structure
func GetUniquePeptides(p id.PepIDList) map[string]id.PepIDList {

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

// ReadProtXMLInput reads one or more fies and organize the data into PSM list
func ReadProtXMLInput(xmlFile, decoyTag string, weight float64) id.ProtXML {

	var protXML id.ProtXML

	protXML.Read(xmlFile)

	protXML.DecoyTag = decoyTag

	protXML.MarkUniquePeptides(weight)

	protXML.PromoteProteinIDs()

	//protXML.Serialize()

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

			pro.Length = "0"
			pro.PercentCoverage = float32(coverMap[pro.ProteinName])
			pro.PctSpectrumIDs = 0.0
			pro.GroupProbability = 1.00
			pro.Confidence = 1.00
			pro.HasRazor = true

			if i.Probability > pro.Probability {
				pro.Probability = i.Probability
				pro.TopPepProb = i.Probability
			}

			razorMarked[pro.ProteinName] = 0

		} else {

			_, ok := razorMarked[pro.ProteinName]
			if !ok {
				pro.Length = "0"
				pro.PercentCoverage = float32(coverMap[pro.ProteinName])
				pro.PctSpectrumIDs = 0.0
				pro.GroupProbability = 1.00
				pro.Confidence = 1.00
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

		//ionForm := fmt.Sprintf("%s#%d#%.4f", i.Peptide, i.AssumedCharge, i.CalcNeutralPepMass)
		pro := proteinList[i.Protein]
		razorProtein, ok := razorMap[i.Peptide]

		if ok && pro.ProteinName == razorProtein {

			pro.UniqueStrippedPeptides = append(pro.UniqueStrippedPeptides, i.Peptide)
			pro.TotalNumberPeptides++

			pep := id.PeptideIonIdentification{
				PeptideSequence:      i.Peptide,
				ModifiedPeptide:      i.ModifiedPeptide,
				Charge:               i.AssumedCharge,
				Weight:               1,
				GroupWeight:          0,
				CalcNeutralPepMass:   i.CalcNeutralPepMass,
				SharedParentProteins: len(i.AlternativeProteins),
				Razor:                1,
			}

			for j := range i.AlternativeProteins {
				pep.PeptideParentProtein = append(pep.PeptideParentProtein, j)
				pro.IndistinguishableProtein = append(pro.IndistinguishableProtein, j)
			}

			pep.NumberOfInstances++

			if i.Probability > pep.InitialProbability {
				pep.InitialProbability = i.Probability
			}

			if len(i.AlternativeProteins) < 2 {
				pep.IsNondegenerateEvidence = true
				pep.IsUnique = true
			} else {
				pep.IsNondegenerateEvidence = false
				pep.IsUnique = false
			}

			pep.Modifications.Index = make(map[string]mod.Modification)
			for k, v := range i.Modifications.Index {
				pep.Modifications.Index[k] = v
			}

			//pro.IndistinguishableProtein = i.AlternativeProteins
			pro.HasRazor = true
			pro.PeptideIons = append(pro.PeptideIons, pep)

			proteinList[i.Protein] = pro

		} else {

			pro.UniqueStrippedPeptides = append(pro.UniqueStrippedPeptides, i.Peptide)
			pro.TotalNumberPeptides++

			pep := id.PeptideIonIdentification{
				PeptideSequence:      i.Peptide,
				ModifiedPeptide:      i.ModifiedPeptide,
				Charge:               i.AssumedCharge,
				Weight:               0,
				GroupWeight:          0,
				CalcNeutralPepMass:   i.CalcNeutralPepMass,
				SharedParentProteins: len(i.AlternativeProteins),
				Razor:                0,
			}

			if i.Probability > pep.InitialProbability {
				pep.InitialProbability = i.Probability
			}

			//pep.PeptideParentProtein = i.AlternativeProteins

			pep.NumberOfInstances++

			if len(i.AlternativeProteins) < 2 {
				pep.IsNondegenerateEvidence = true
				pep.IsUnique = true
			} else {
				pep.IsNondegenerateEvidence = false
				pep.IsUnique = false
			}

			pep.Modifications.Index = make(map[string]mod.Modification)
			for k, v := range i.Modifications.Index {
				pep.Modifications.Index[k] = v
			}

			//pro.IndistinguishableProtein = i.AlternativeProteins
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
func extractPSMfromPepXML(filter string, peplist id.PepIDList, pro id.ProtIDList) id.PepIDList {

	var protmap = make(map[string]uint16)
	var filterMap = make(map[string]id.PeptideIdentification)
	var output id.PepIDList

	if filter == "sequential" {

		// get all protein and peptide pairs from protxml
		for _, i := range pro {
			for _, j := range i.UniqueStrippedPeptides {
				key := fmt.Sprintf("%s#%s", i.ProteinName, j)
				protmap[string(key)] = 0
			}
		}

		for _, i := range peplist {

			key := fmt.Sprintf("%s#%s", i.Protein, i.Peptide)

			_, ok := protmap[key]
			if ok {
				filterMap[string(i.Spectrum)] = i
			} else {

				for j := range i.AlternativeProteins {
					key := fmt.Sprintf("%s#%s", j, i.Peptide)
					_, ap := protmap[key]
					if ap {
						filterMap[string(i.Spectrum)] = i
					}
				}

			}

		}

	} else if filter == "2d" {

		// get all protein names from protxml
		for _, i := range pro {
			protmap[string(i.ProteinName)] = 0
		}

		for _, i := range peplist {
			_, ok := protmap[string(i.Protein)]
			if ok {
				filterMap[string(i.Spectrum)] = i
			} else {
				for j := range i.AlternativeProteins {
					_, ap := protmap[j]
					if ap {
						filterMap[string(i.Spectrum)] = i
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
