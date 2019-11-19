package fil

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"philosopher/lib/cla"
	"philosopher/lib/dat"
	"philosopher/lib/id"
	"philosopher/lib/inf"
	"philosopher/lib/met"
	"philosopher/lib/mod"
	"philosopher/lib/msg"
	"philosopher/lib/qua"
	"philosopher/lib/rep"
	"philosopher/lib/spc"
	"philosopher/lib/sys"
)

// Run executes the Filter processing
func Run(f met.Data) met.Data {

	e := rep.New()
	var pepxml id.PepXML
	var pep id.PepIDList
	var pro id.ProtIDList

	// get the database tag from database command
	if len(f.Filter.Tag) == 0 {
		f.Filter.Tag = f.Database.Tag
	}

	logrus.Info("Processing peptide identification files")

	// if no method is selected, force the 2D to be default
	if len(f.Filter.Pox) > 0 && f.Filter.TwoD == false && f.Filter.Seq == false {
		f.Filter.TwoD = true
	}

	pepid, searchEngine := readPepXMLInput(f.Filter.Pex, f.Filter.Tag, f.Temp, f.Filter.Model, f.MSFragger.CalibrateMass)

	f.SearchEngine = searchEngine

	psmT, pepT, ionT := processPeptideIdentifications(pepid, f.Filter.Tag, f.Filter.PsmFDR, f.Filter.PepFDR, f.Filter.IonFDR)
	_ = psmT
	_ = pepT
	_ = ionT

	if len(f.Filter.Pox) > 0 {

		protXML := readProtXMLInput(sys.MetaDir(), f.Filter.Pox, f.Filter.Tag, f.Filter.Weight)
		processProteinIdentifications(protXML, f.Filter.PtFDR, f.Filter.PepFDR, f.Filter.ProtProb, f.Filter.Picked, f.Filter.Razor, f.Filter.Fo, f.Filter.Tag)

	} else {

		if f.Filter.Inference == true {
			pepid = inf.ProteinInference(pepid)

			pepid.Serialize("psm")
			pepid.Serialize("pep")
			pepid.Serialize("ion")

			processProteinInferenceIdentifications(pepid, f.Filter.PtFDR, f.Filter.PepFDR, f.Filter.ProtProb, f.Filter.Picked, f.Filter.Tag)
		}

	}

	if f.Filter.Seq == true {

		// sequential analysis
		// filtered psm list and filtered prot list
		pep.Restore("psm")
		pro.Restore()
		sequentialFDRControl(pep, pro, f.Filter.PsmFDR, f.Filter.PepFDR, f.Filter.IonFDR, f.Filter.Tag)
		pep = nil
		pro = nil

	} else if f.Filter.TwoD == true {

		// two-dimensional analysis
		// complete pep list and filtered mirror-image prot list
		pepxml.Restore()
		pro.Restore()
		twoDFDRFilter(pepxml.PeptideIdentification, pro, f.Filter.PsmFDR, f.Filter.PepFDR, f.Filter.IonFDR, f.Filter.Tag)
		pepxml = id.PepXML{}
		pro = nil

	}

	var dtb dat.Base
	dtb.Restore()
	if len(dtb.Records) < 1 {
		msg.Custom(errors.New("Database data not available, interrupting processing"), "fatal")
	}

	logrus.Info("Post processing identifications")

	// restoring for the modifications
	var pxml id.PepXML
	pxml.Restore()

	e.Mods = pxml.Modifications
	e.AssembleSearchParameters(pxml.SearchParameters)
	pxml = id.PepXML{}

	var psm id.PepIDList
	psm.Restore("psm")
	e.AssemblePSMReport(psm, f.Filter.Tag)
	psm = nil

	var ion id.PepIDList
	ion.Restore("ion")
	e.AssembleIonReport(ion, f.Filter.Tag)
	ion = nil

	// evaluate modifications in data set
	if f.Filter.Mapmods == true {
		logrus.Info("Mapping modifications")
		e.MapMods()

		logrus.Info("Processing modifications")
		e.AssembleModificationReport()
	}

	var pept id.PepIDList
	pept.Restore("pep")
	e.AssemblePeptideReport(pept, f.Filter.Tag)
	pept = nil

	// evaluate modifications in data set
	if f.Filter.Mapmods == true {
		e.UpdateIonModCount()
		e.UpdatePeptideModCount()
	}

	if len(f.Filter.Pox) > 0 {

		logrus.Info("Processing protein inference")
		pro.Restore()
		e.AssembleProteinReport(pro, f.Filter.Weight, f.Filter.Tag)
		pro = nil

		// Pushes the new ion status from the protein inferece to the other layers, the gene and protein ID
		// assignment gets corrected in the next function call (UpdateLayerswithDatabase)
		e.UpdateIonStatus(f.Filter.Tag)
	}

	logrus.Info("Assigning protein identifications to layers")
	e.UpdateLayerswithDatabase(f.Filter.Tag)

	// reorganizes the selected proteins and the alternative proteins list
	logrus.Info("Updating razor PSM assignment to proteins")
	if f.Filter.Razor == true {
		e.UpdateSupportingSpectra()
	}

	logrus.Info("Calculating spectral counts")
	e = qua.CalculateSpectralCounts(e)

	logrus.Info("Saving")
	e.SerializeGranular()

	return f
}

// readPepXMLInput reads one or more fies and organize the data into PSM list
func readPepXMLInput(xmlFile, decoyTag, temp string, models bool, calibratedMass int) (id.PepIDList, string) {

	var files []string
	var pepIdent id.PepIDList
	var mods []mod.Modification
	var params []spc.Parameter
	var modsIndex = make(map[string]mod.Modification)
	var searchEngine string

	if strings.Contains(xmlFile, "pep.xml") || strings.Contains(xmlFile, "pepXML") {
		files = append(files, xmlFile)
	} else {
		glob := fmt.Sprintf("%s/*pep.xml", xmlFile)
		list, _ := filepath.Glob(glob)

		if len(list) == 0 {
			msg.NoParametersFound(errors.New("missing pepXML files"), "fatal")
		}

		for _, i := range list {
			absPath, _ := filepath.Abs(i)
			files = append(files, absPath)
		}

	}

	for _, i := range files {
		var p id.PepXML
		p.DecoyTag = decoyTag
		p.Read(i)

		params = p.SearchParameters

		// print models
		if models == true {
			if strings.EqualFold(p.Prophet, "interprophet") {
				logrus.Error("Cannot print models for interprophet files")
			} else {
				logrus.Info("Printing models")
				go p.ReportModels(temp, filepath.Base(i))
				time.Sleep(time.Second * 3)
			}
		}

		pepIdent = append(pepIdent, p.PeptideIdentification...)

		for _, k := range p.Modifications.Index {
			_, ok := modsIndex[k.Index]
			if !ok {
				mods = append(mods, k)
				modsIndex[k.Index] = k
			}
		}

		searchEngine = p.SearchEngine
	}

	// create a "fake" global pepXML comprising all data
	var pepXML id.PepXML
	pepXML.DecoyTag = decoyTag
	pepXML.SearchParameters = params
	pepXML.PeptideIdentification = pepIdent
	pepXML.Modifications.Index = modsIndex

	// promoting Spectra that matches to both decoys and targets to TRUE hits
	pepXML.PromoteProteinIDs()

	// serialize all pep files
	sort.Sort(pepXML.PeptideIdentification)
	pepXML.Serialize()

	return pepIdent, searchEngine
}

// processPeptideIdentifications reads and process pepXML
func processPeptideIdentifications(p id.PepIDList, decoyTag string, psm, peptide, ion float64) (float64, float64, float64) {

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

	return psmThreshold, peptideThreshold, ionThreshold
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

//GetUniquePSMs selects only unique pepetide ions for the given data stucture
func GetUniquePSMs(p id.PepIDList) map[string]id.PepIDList {

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

// GetUniquePeptides selects only unique pepetide for the given data stucture
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

// readProtXMLInput reads one or more fies and organize the data into PSM list
func readProtXMLInput(meta, xmlFile, decoyTag string, weight float64) id.ProtXML {

	var protXML id.ProtXML

	protXML.Read(xmlFile)

	protXML.DecoyTag = decoyTag

	protXML.MarkUniquePeptides(weight)

	protXML.PromoteProteinIDs()

	protXML.Serialize()

	return protXML
}

// processProteinIdentifications checks if pickedFDR ar razor options should be applied to given data set, if they do,
// the inputed protXML data is processed before filtered.
func processProteinIdentifications(p id.ProtXML, ptFDR, pepProb, protProb float64, isPicked, isRazor, fo bool, decoyTag string) {

	var pid id.ProtIDList

	// tagget / decoy / threshold
	t, d := proteinProfile(p)
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
		p = RazorFilter(p)
	}

	// run the FDR filter for proteins
	pid = ProtXMLFilter(p, ptFDR, pepProb, protProb, isPicked, isRazor, decoyTag)

	if fo == true {
		output := fmt.Sprintf("%s%spep_pro_mappings.tsv", sys.MetaDir(), string(filepath.Separator))

		file, e := os.Create(output)
		if e != nil {
			msg.WriteFile(e, "fatal")
		}
		defer file.Close()

		for _, i := range pid {
			if !strings.HasPrefix(i.ProteinName, decoyTag) {

				var line []string

				line = append(line, i.ProteinName)

				for _, j := range i.PeptideIons {
					if j.Razor == 1 {
						line = append(line, j.PeptideSequence)
					}
				}

				mapping := strings.Join(line, "\t")
				_, e = io.WriteString(file, mapping)
				if e != nil {
					msg.WriteToFile(e, "fatal")
				}

			}
		}
	}

	// save results on meta folder
	pid.Serialize()

	return
}

// processProteinInferenceIdentifications checks if pickedFDR ar razor options should be applied to given data set, if they do,
// the inputed Philospher inference data is processed before filtered.
func processProteinInferenceIdentifications(psm id.PepIDList, ptFDR, pepProb, protProb float64, isPicked bool, decoyTag string) {

	var t int
	var d int
	var proteinIndex = make(map[string]uint)
	//var pid id.ProtIDList

	for _, i := range psm {

		proteinIndex[i.Protein]++
		for j := range i.AlternativeProteinsIndexed {
			proteinIndex[j]++
		}

		// var p id.ProteinIdentification
		// p.GroupNumber = 0
		// p.GroupSiblingID = "a"
		// p.ProteinName = i.Protein
		// p.Probability = i.Probability
		// p.TopPepProb = i.Probability
		// p.HasRazor = true
		// p.Picked = 0
	}

	for i := range proteinIndex {
		if strings.HasPrefix(i, decoyTag) {
			d++
		} else {
			t++
		}
	}

	// tagget / decoy / threshold
	logrus.WithFields(logrus.Fields{
		"target": t,
		"decoy":  d,
	}).Info("Protein inference results")

	// // run the FDR filter for proteins
	// pid = ProtXMLFilter(p, ptFDR, pepProb, protProb, isPicked, isRazor, decoyTag)

	// // save results on meta folder
	// pid.Serialize()

	return
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
func extractPSMfromPepXML(peplist id.PepIDList, pro id.ProtIDList) id.PepIDList {

	var protmap = make(map[string]uint16)
	var filterMap = make(map[string]id.PeptideIdentification)
	var output id.PepIDList

	// get all protein names from protxml
	for _, i := range pro {
		protmap[string(i.ProteinName)] = 0
	}

	for _, i := range peplist {
		_, ok := protmap[string(i.Protein)]
		if ok {
			filterMap[string(i.Spectrum)] = i
		} else {
			for _, j := range i.AlternativeProteins {
				_, ap := protmap[string(j)]
				if ap {
					filterMap[string(i.Spectrum)] = i
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
