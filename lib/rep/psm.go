package rep

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"philosopher/lib/msg"
	"philosopher/lib/uti"

	"philosopher/lib/bio"
	"philosopher/lib/cla"
	"philosopher/lib/dat"
	"philosopher/lib/id"
)

// AssemblePSMReport creates the PSM structure for reporting
func (evi *Evidence) AssemblePSMReport(pep id.PepIDList, decoyTag string) {

	var list PSMEvidenceList

	// collect database information
	var dtb dat.Base
	dtb.Restore()

	var genes = make(map[string]string)
	var ptid = make(map[string]string)
	for _, j := range dtb.Records {
		genes[j.PartHeader] = j.GeneNames
		ptid[j.PartHeader] = j.ID
	}

	for _, i := range pep {

		var p PSMEvidence

		source := strings.Split(i.Spectrum, ".")
		p.Source = source[0]
		p.Index = i.Index
		p.Spectrum = i.Spectrum
		p.SpectrumFile = i.SpectrumFile
		p.Scan = i.Scan
		p.PrevAA = ""
		p.NextAA = ""
		p.NumberOfEnzymaticTermini = int(i.NumberOfEnzymaticTermini)
		p.NumberOfMissedCleavages = i.NumberofMissedCleavages
		p.Peptide = i.Peptide
		p.IonForm = fmt.Sprintf("%s#%d#%.4f", i.Peptide, i.AssumedCharge, i.CalcNeutralPepMass)
		p.Protein = i.Protein
		p.ModifiedPeptide = i.ModifiedPeptide
		p.AssumedCharge = i.AssumedCharge
		p.HitRank = i.HitRank
		p.PrecursorExpMass = i.PrecursorExpMass
		p.RetentionTime = i.RetentionTime
		p.CalcNeutralPepMass = i.CalcNeutralPepMass
		p.Massdiff = i.Massdiff
		p.LocalizedPTMSites = i.LocalizedPTMSites
		p.LocalizedPTMMassDiff = i.LocalizedPTMMassDiff
		p.Probability = i.Probability
		p.Expectation = i.Expectation
		p.Xcorr = i.Xcorr
		p.DeltaCN = i.DeltaCN
		p.SPRank = i.SPRank
		p.Hyperscore = i.Hyperscore
		p.Nextscore = i.Nextscore
		p.DiscriminantValue = i.DiscriminantValue
		p.Intensity = i.Intensity
		p.IonMobility = i.IonMobility
		p.CompensationVoltage = i.CompesationVoltage
		p.MappedGenes = make(map[string]int)
		p.MappedProteins = make(map[string]int)
		p.Modifications = i.Modifications
		p.MSFragerLocalization = i.MSFragerLocalization
		p.MSFraggerLocalizationScoreWithPTM = i.MSFraggerLocalizationScoreWithPTM
		p.MSFraggerLocalizationScoreWithoutPTM = i.MSFraggerLocalizationScoreWithoutPTM

		if i.UncalibratedPrecursorNeutralMass > 0 {
			p.PrecursorNeutralMass = i.PrecursorNeutralMass
			p.UncalibratedPrecursorNeutralMass = i.UncalibratedPrecursorNeutralMass
		} else {
			p.PrecursorNeutralMass = i.PrecursorNeutralMass
			p.UncalibratedPrecursorNeutralMass = i.PrecursorNeutralMass
		}

		for j := range i.AlternativeProteins {
			p.MappedProteins[j]++
		}

		gn, ok := genes[i.Protein]
		if ok {
			p.GeneName = gn
		}

		id, ok := ptid[i.Protein]
		if ok {
			p.ProteinID = id
		}

		// is this bservation a decoy ?
		if cla.IsDecoyPSM(i, decoyTag) {
			p.IsDecoy = true
		}

		// the redudnancy check was introduced because of inconsistencies with
		// PeptideProphet. The Windows version is printing the same protein
		// as alternative when the peptide maps to the same protein multiple times
		var redudantMapping = 0
		if len(i.AlternativeProteins) == 0 {
			p.IsUnique = true
		} else {
			for k := range i.AlternativeProteins {
				if k == i.Protein {
					redudantMapping++
				}
			}
			p.IsUnique = false
		}

		if redudantMapping == len(i.AlternativeProteins) {
			p.IsUnique = true
		}

		list = append(list, p)
	}

	sort.Sort(list)
	evi.PSM = list
}

// MetaPSMReport report all psms from study that passed the FDR filter
func (evi Evidence) MetaPSMReport(workspace, brand string, channels int, hasDecoys, isComet, hasLoc, hasLabels bool) {

	var header string
	var modMap = make(map[string]string)
	var modList []string
	var hasCompVolt bool
	var hasIonMob bool
	var hasPurity bool

	output := fmt.Sprintf("%s%spsm.tsv", workspace, string(filepath.Separator))

	// create result file
	file, e := os.Create(output)
	if e != nil {
		msg.WriteFile(errors.New("cannot create report file"), "fatal")
	}
	defer file.Close()

	// building the printing set tat may or not contain decoys
	var printSet PSMEvidenceList
	for i := range evi.PSM {

		compositeName := strings.Split(evi.PSM[i].Spectrum, "#")
		evi.PSM[i].Spectrum = compositeName[0]

		if !hasDecoys {
			if !evi.PSM[i].IsDecoy {
				printSet = append(printSet, evi.PSM[i])
			}
		} else {
			printSet = append(printSet, evi.PSM[i])
		}

		for k := range evi.PSM[i].LocalizedPTMMassDiff {
			_, ok := modMap[k]
			if !ok {
				modMap[k] = ""
			} else {
				modMap[k] = ""
			}
		}

		if len(evi.PSM[i].CompensationVoltage) > 0 {
			hasCompVolt = true
		}

		if evi.PSM[i].IonMobility > 0 {
			hasIonMob = true
		}

		if evi.PSM[i].Purity > 0 {
			hasPurity = true
		}

		if len(evi.PSM[i].MSFragerLocalization) > 0 {
			hasLoc = true
		}

	}

	for k := range modMap {
		modList = append(modList, k)
	}

	sort.Strings(modList)

	header = "Spectrum\tSpectrum File\tPeptide\tModified Peptide\tPrev AA\tNext AA\tPeptide Length\tCharge\tRetention\tObserved Mass\tCalibrated Observed Mass\tObserved M/Z\tCalibrated Observed M/Z\tCalculated Peptide Mass\tCalculated M/Z\tDelta Mass"

	if isComet {
		header += "\tXCorr\tDeltaCN\tDeltaCNStar\tSPScore\tSPRank"
	}

	header += "\tExpectation\tHyperscore\tNextscore\tPeptideProphet Probability\tNumber of Enzymatic Termini\tNumber of Missed Cleavages\tProtein Start\tProtein End\tIntensity\tAssigned Modifications\tObserved Modifications"

	if len(modList) > 0 {
		for _, i := range modList {
			if strings.Contains(i, "STY:79.966331") {
				i = "STY:79.9663"
			}
			header += "\t" + i + "\t" + i + " Best Localization"
		}
	}

	if hasLoc {
		header += "\tMSFragger Localization\tBest Score with Delta Mass\tBest Score without Delta Mass"
	}

	if hasIonMob {
		header += "\tIon Mobility"
	}

	if hasCompVolt {
		header += "\tCompensation Voltage"
	}

	if hasPurity {
		header += "\tPurity"
	}

	header += "\tIs Unique\tProtein\tProtein ID\tEntry Name\tGene\tProtein Description\tMapped Genes\tMapped Proteins"

	if brand == "tmt" {
		switch channels {
		case 6:
			header += "\tQuan Usage\tChannel 126\tChannel 127N\tChannel 128C\tChannel 129N\tChannel 130C\tChannel 131N"
		case 10:
			header += "\tQuan Usage\tChannel 126\tChannel 127N\tChannel 127C\tChannel 128N\tChannel 128C\tChannel 129N\tChannel 129C\tChannel 130N\tChannel 130C\tChannel 131N"
		case 11:
			header += "\tQuan Usage\tChannel 126\tChannel 127N\tChannel 127C\tChannel 128N\tChannel 128C\tChannel 129N\tChannel 129C\tChannel 130N\tChannel 130C\tChannel 131N\tChannel 131C"
		case 16:
			header += "\tQuan Usage\tChannel 126\tChannel 127N\tChannel 127C\tChannel 128N\tChannel 128C\tChannel 129N\tChannel 129C\tChannel 130N\tChannel 130C\tChannel 131N\tChannel 131C\tChannel 132N\tChannel 132C\tChannel 133N\tChannel 133C\tChannel 134N"
		case 18:
			header += "\tChannel 126\tChannel 127N\tChannel 127C\tChannel 128N\tChannel 128C\tChannel 129N\tChannel 129C\tChannel 130N\tChannel 130C\tChannel 131N\tChannel 131C\tChannel 132N\tChannel 132C\tChannel 133N\tChannel 133C\tChannel 134N\tChannel 134C\tChannel 135N"
		default:
			header += ""
		}
	} else if brand == "itraq" {
		switch channels {
		case 4:
			header += "\tQuan Usage\tChannel 114\tChannel 115\tChannel 116\tChannel 117"
		case 8:
			header += "\tQuan Usage\tChannel 113\tChannel 114\tChannel 115\tChannel 116\tChannel 117\tChannel 118\tChannel 119\tChannel 121"
		default:
			header += ""
		}
	}

	header += "\n"

	// verify if the structure has labels, if so, replace the original channel names by them.
	if hasLabels {

		var c1, c2, c3, c4, c5, c6, c7, c8, c9, c10, c11, c12, c13, c14, c15, c16, c17, c18 string

		for _, i := range printSet {
			if len(i.Labels.Channel1.CustomName) >= 1 {
				c1 = i.Labels.Channel1.CustomName
				c2 = i.Labels.Channel2.CustomName
				c3 = i.Labels.Channel3.CustomName
				c4 = i.Labels.Channel4.CustomName
				c5 = i.Labels.Channel5.CustomName
				c6 = i.Labels.Channel6.CustomName
				c7 = i.Labels.Channel7.CustomName
				c8 = i.Labels.Channel8.CustomName
				c9 = i.Labels.Channel9.CustomName
				c10 = i.Labels.Channel10.CustomName
				c11 = i.Labels.Channel11.CustomName
				c12 = i.Labels.Channel12.CustomName
				c13 = i.Labels.Channel13.CustomName
				c14 = i.Labels.Channel14.CustomName
				c15 = i.Labels.Channel15.CustomName
				c16 = i.Labels.Channel16.CustomName
				c17 = i.Labels.Channel17.CustomName
				c18 = i.Labels.Channel18.CustomName
				break
			}
		}

		header = strings.Replace(header, "Channel "+printSet[10].Labels.Channel1.Name, c1, -1)
		header = strings.Replace(header, "Channel "+printSet[10].Labels.Channel2.Name, c2, -1)
		header = strings.Replace(header, "Channel "+printSet[10].Labels.Channel3.Name, c3, -1)
		header = strings.Replace(header, "Channel "+printSet[10].Labels.Channel4.Name, c4, -1)
		header = strings.Replace(header, "Channel "+printSet[10].Labels.Channel5.Name, c5, -1)
		header = strings.Replace(header, "Channel "+printSet[10].Labels.Channel6.Name, c6, -1)
		header = strings.Replace(header, "Channel "+printSet[10].Labels.Channel7.Name, c7, -1)
		header = strings.Replace(header, "Channel "+printSet[10].Labels.Channel8.Name, c8, -1)
		header = strings.Replace(header, "Channel "+printSet[10].Labels.Channel9.Name, c9, -1)
		header = strings.Replace(header, "Channel "+printSet[10].Labels.Channel10.Name, c10, -1)
		header = strings.Replace(header, "Channel "+printSet[10].Labels.Channel11.Name, c11, -1)
		header = strings.Replace(header, "Channel "+printSet[10].Labels.Channel12.Name, c12, -1)
		header = strings.Replace(header, "Channel "+printSet[10].Labels.Channel13.Name, c13, -1)
		header = strings.Replace(header, "Channel "+printSet[10].Labels.Channel14.Name, c14, -1)
		header = strings.Replace(header, "Channel "+printSet[10].Labels.Channel15.Name, c15, -1)
		header = strings.Replace(header, "Channel "+printSet[10].Labels.Channel16.Name, c16, -1)
		header = strings.Replace(header, "Channel "+printSet[10].Labels.Channel17.Name, c17, -1)
		header = strings.Replace(header, "Channel "+printSet[10].Labels.Channel18.Name, c18, -1)
	}

	_, e = io.WriteString(file, header)
	if e != nil {
		msg.WriteToFile(errors.New("cannot print PSM to file"), "fatal")
	}

	for _, i := range printSet {

		assL, obs := getModsList(i.Modifications.Index)

		var mappedProteins []string
		for j := range i.MappedProteins {
			if j != i.Protein && len(j) > 0 {
				mappedProteins = append(mappedProteins, j)
			}
		}

		var mappedGenes []string
		for j := range i.MappedGenes {
			if j != i.GeneName && len(j) > 0 {
				mappedGenes = append(mappedGenes, j)
			}
		}

		sort.Strings(mappedGenes)
		sort.Strings(mappedProteins)
		sort.Strings(assL)
		sort.Strings(obs)

		line := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%d\t%d\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
			i.Spectrum,
			i.SpectrumFile,
			i.Peptide,
			i.ModifiedPeptide,
			i.PrevAA,
			i.NextAA,
			len(i.Peptide),
			i.AssumedCharge,
			i.RetentionTime,
			i.UncalibratedPrecursorNeutralMass,
			i.PrecursorNeutralMass,
			((i.UncalibratedPrecursorNeutralMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge)),
			((i.PrecursorNeutralMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge)),
			i.CalcNeutralPepMass,
			((i.CalcNeutralPepMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge)),
			i.Massdiff,
		)

		if isComet {
			line = fmt.Sprintf("%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
				line,
				i.Xcorr,
				i.DeltaCN,
				i.DeltaCNStar,
				i.SPScore,
				i.SPRank,
			)
		}

		line = fmt.Sprintf("%s\t%.14f\t%.4f\t%.4f\t%.4f\t%d\t%d\t%d\t%d\t%.4f\t%s\t%s",
			line,
			i.Expectation,
			i.Hyperscore,
			i.Nextscore,
			i.Probability,
			i.NumberOfEnzymaticTermini,
			i.NumberOfMissedCleavages,
			i.ProteinStart,
			i.ProteinEnd,
			i.Intensity,
			strings.Join(assL, ", "),
			strings.Join(obs, ", "),
		)

		if len(modList) > 0 {
			for _, j := range modList {

				r := regexp.MustCompile(`\d\.\d{3}`)
				matches := r.FindAllString(i.LocalizedPTMMassDiff[j], -1)
				max := uti.GetMaxNumber(matches)

				line = fmt.Sprintf("%s\t%s\t%s",
					line,
					i.LocalizedPTMMassDiff[j],
					max,
				)
			}
		}

		if hasLoc {
			line = fmt.Sprintf("%s\t%s\t%s\t%s",
				line,
				i.MSFragerLocalization,
				i.MSFraggerLocalizationScoreWithPTM,
				i.MSFraggerLocalizationScoreWithoutPTM,
			)
		}

		if hasIonMob {
			line = fmt.Sprintf("%s\t%.4f",
				line,
				i.IonMobility,
			)
		}

		if hasCompVolt {
			line = fmt.Sprintf("%s\t%s",
				line,
				i.CompensationVoltage,
			)
		}

		if hasPurity {
			line = fmt.Sprintf("%s\t%.2f",
				line,
				i.Purity,
			)
		}

		line = fmt.Sprintf("%s\t%t\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
			line,
			i.IsUnique,
			i.Protein,
			i.ProteinID,
			i.EntryName,
			i.GeneName,
			i.ProteinDescription,
			strings.Join(mappedGenes, ", "),
			strings.Join(mappedProteins, ", "),
		)

		if brand == "tmt" {
			switch channels {
			case 6:
				line = fmt.Sprintf("%s\t%t\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
					line,
					i.Labels.IsUsed,
					i.Labels.Channel1.Intensity,
					i.Labels.Channel2.Intensity,
					i.Labels.Channel5.Intensity,
					i.Labels.Channel6.Intensity,
					i.Labels.Channel9.Intensity,
					i.Labels.Channel10.Intensity,
				)
			case 10:
				line = fmt.Sprintf("%s\t%t\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
					line,
					i.Labels.IsUsed,
					i.Labels.Channel1.Intensity,
					i.Labels.Channel2.Intensity,
					i.Labels.Channel3.Intensity,
					i.Labels.Channel4.Intensity,
					i.Labels.Channel5.Intensity,
					i.Labels.Channel6.Intensity,
					i.Labels.Channel7.Intensity,
					i.Labels.Channel8.Intensity,
					i.Labels.Channel9.Intensity,
					i.Labels.Channel10.Intensity,
				)
			case 11:
				line = fmt.Sprintf("%s\t%t\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
					line,
					i.Labels.IsUsed,
					i.Labels.Channel1.Intensity,
					i.Labels.Channel2.Intensity,
					i.Labels.Channel3.Intensity,
					i.Labels.Channel4.Intensity,
					i.Labels.Channel5.Intensity,
					i.Labels.Channel6.Intensity,
					i.Labels.Channel7.Intensity,
					i.Labels.Channel8.Intensity,
					i.Labels.Channel9.Intensity,
					i.Labels.Channel10.Intensity,
					i.Labels.Channel11.Intensity,
				)
			case 16:
				line = fmt.Sprintf("%s\t%t\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
					line,
					i.Labels.IsUsed,
					i.Labels.Channel1.Intensity,
					i.Labels.Channel2.Intensity,
					i.Labels.Channel3.Intensity,
					i.Labels.Channel4.Intensity,
					i.Labels.Channel5.Intensity,
					i.Labels.Channel6.Intensity,
					i.Labels.Channel7.Intensity,
					i.Labels.Channel8.Intensity,
					i.Labels.Channel9.Intensity,
					i.Labels.Channel10.Intensity,
					i.Labels.Channel11.Intensity,
					i.Labels.Channel12.Intensity,
					i.Labels.Channel13.Intensity,
					i.Labels.Channel14.Intensity,
					i.Labels.Channel15.Intensity,
					i.Labels.Channel16.Intensity,
				)
			case 18:
				line = fmt.Sprintf("%s\t%t\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
					line,
					i.Labels.IsUsed,
					i.Labels.Channel1.Intensity,
					i.Labels.Channel2.Intensity,
					i.Labels.Channel3.Intensity,
					i.Labels.Channel4.Intensity,
					i.Labels.Channel5.Intensity,
					i.Labels.Channel6.Intensity,
					i.Labels.Channel7.Intensity,
					i.Labels.Channel8.Intensity,
					i.Labels.Channel9.Intensity,
					i.Labels.Channel10.Intensity,
					i.Labels.Channel11.Intensity,
					i.Labels.Channel12.Intensity,
					i.Labels.Channel13.Intensity,
					i.Labels.Channel14.Intensity,
					i.Labels.Channel15.Intensity,
					i.Labels.Channel16.Intensity,
					i.Labels.Channel17.Intensity,
					i.Labels.Channel18.Intensity,
				)
			default:
				header += ""
			}
		} else if brand == "itraq" {
			switch channels {
			case 4:
				line = fmt.Sprintf("%s\t%t\t%.4f\t%.4f\t%.4f\t%.4f",
					line,
					i.Labels.IsUsed,
					i.Labels.Channel1.Intensity,
					i.Labels.Channel2.Intensity,
					i.Labels.Channel3.Intensity,
					i.Labels.Channel4.Intensity,
				)
			case 8:
				line = fmt.Sprintf("%s\t%t\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
					line,
					i.Labels.IsUsed,
					i.Labels.Channel1.Intensity,
					i.Labels.Channel2.Intensity,
					i.Labels.Channel3.Intensity,
					i.Labels.Channel4.Intensity,
					i.Labels.Channel5.Intensity,
					i.Labels.Channel6.Intensity,
					i.Labels.Channel7.Intensity,
					i.Labels.Channel8.Intensity,
				)
			default:
				header += ""
			}
		}

		line += "\n"

		_, e = io.WriteString(file, line)
		if e != nil {
			msg.WriteToFile(e, "fatal")
		}
	}
}

// PSMLocalizationReport report ptm localization based on PTMProphet outputs
func (evi *Evidence) PSMLocalizationReport(workspace, decoyTag string, hasRazor, hasDecoys bool) {

	output := fmt.Sprintf("%s%slocalization.tsv", workspace, string(filepath.Separator))

	// create result file
	file, e := os.Create(output)
	if e != nil {
		msg.WriteFile(e, "fatal")
	}
	defer file.Close()

	_, e = io.WriteString(file, "Spectrum\tPeptide\tModified Peptide\tCharge\tRetention\tModification\tNumber of Sites\tObserved Mass Localization\n")
	if e != nil {
		msg.WriteToFile(e, "fatal")
	}

	// building the printing set tat may or not contain decoys
	var printSet PSMEvidenceList
	for _, i := range evi.PSM {
		if !hasDecoys {
			if !i.IsDecoy {
				printSet = append(printSet, i)
			}
		} else {
			printSet = append(printSet, i)
		}
	}

	for _, i := range printSet {
		for j := range i.LocalizedPTMMassDiff {
			line := fmt.Sprintf("%s\t%s\t%s\t%d\t%.4f\t%s\t%d\t%s\n",
				i.Spectrum,
				i.Peptide,
				i.ModifiedPeptide,
				i.AssumedCharge,
				i.RetentionTime,
				j,
				i.LocalizedPTMSites[j],
				i.LocalizedPTMMassDiff[j],
			)
			_, e = io.WriteString(file, line)
			if e != nil {
				msg.WriteToFile(e, "fatal")
			}
		}
	}
}
