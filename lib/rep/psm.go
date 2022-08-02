package rep

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
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

	var genes = make(map[string]string)
	var ptid = make(map[string]string)
	{
		// collect database information
		var dtb dat.Base
		dtb.Restore()

		for _, j := range dtb.Records {
			genes[j.PartHeader] = j.GeneNames
			ptid[j.PartHeader] = j.ID
		}
	}
	evi.PSM = make(PSMEvidenceList, len(pep))
	for idx, i := range pep {

		p := &evi.PSM[idx]

		source := strings.Split(i.Spectrum, ".")
		p.Source = source[0]
		p.Index = i.Index
		p.Spectrum = i.Spectrum
		p.SpectrumFile = i.SpectrumFile
		p.NumberOfEnzymaticTermini = i.NumberOfEnzymaticTermini
		p.NumberOfMissedCleavages = i.NumberofMissedCleavages
		p.Peptide = i.Peptide
		p.Protein = i.Protein
		p.ModifiedPeptide = i.ModifiedPeptide
		p.AssumedCharge = i.AssumedCharge
		p.HitRank = i.HitRank
		p.RetentionTime = i.RetentionTime
		p.CalcNeutralPepMass = i.CalcNeutralPepMass
		p.Massdiff = i.Massdiff
		p.PTM = i.PTM
		p.Probability = i.Probability
		p.Expectation = i.Expectation
		p.Xcorr = i.Xcorr
		p.DeltaCN = i.DeltaCN
		p.SPRank = i.SPRank
		p.Hyperscore = i.Hyperscore
		p.Nextscore = i.Nextscore
		p.SpectralSim = i.SpectralSim
		p.Rtscore = i.Rtscore
		p.Intensity = i.Intensity
		p.IonMobility = i.IonMobility
		p.CompensationVoltage = i.CompensationVoltage
		p.MappedGenes = make(map[string]struct{})
		p.MappedProteins = make(map[string]int)
		p.Modifications = i.Modifications
		p.MSFraggerLoc = i.MSFragerLoc
		if i.UncalibratedPrecursorNeutralMass > 0 {
			p.PrecursorNeutralMass = float64(i.PrecursorNeutralMass)
			p.UncalibratedPrecursorNeutralMass = float64(i.UncalibratedPrecursorNeutralMass)
		} else {
			p.PrecursorNeutralMass = float64(i.PrecursorNeutralMass)
			p.UncalibratedPrecursorNeutralMass = float64(i.PrecursorNeutralMass)
		}

		// Forcing the modified peptide string to be empty in case no mods are present
		if len(p.Modifications.IndexSlice) == 0 {
			p.ModifiedPeptide = ""
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

	}

	sort.Sort(evi.PSM)
}

// MetaPSMReport report all psms from study that passed the FDR filter
func (evi PSMEvidenceList) MetaPSMReport(workspace, brand, decoyTag string, channels int, hasDecoys, isComet, hasLoc, hasIonMob, hasLabels, hasPrefix bool) {

	var header string
	var output string
	var modMap = make(map[string]string)
	var modList []string
	var hasCompVolt bool
	var hasPurity bool
	var hasSpectralSim bool
	var hasRtScore bool

	if hasPrefix {
		output = fmt.Sprintf("%s%s%s_psm.tsv", workspace, string(filepath.Separator), path.Base(workspace))
	} else {
		output = fmt.Sprintf("%s%spsm.tsv", workspace, string(filepath.Separator))
	}

	// create result file
	file, e := os.Create(output)
	bw := bufio.NewWriter(file)
	if e != nil {
		msg.WriteFile(errors.New("cannot create report file, "+e.Error()), "fatal")
	}
	defer file.Close()
	defer bw.Flush()

	// building the printing set tat may or not contain decoys
	//var printSet PSMEvidenceList
	var printSet []*PSMEvidence
	for i := range evi {

		if !hasDecoys {
			if !evi[i].IsDecoy {
				printSet = append(printSet, &evi[i])
			}
		} else {
			printSet = append(printSet, &evi[i])
		}

		if evi[i].PTM != nil {
			for k := range evi[i].PTM.LocalizedPTMMassDiff {
				_, ok := modMap[k]
				if !ok {
					modMap[k] = ""
				} else {
					modMap[k] = ""
				}
			}
		}

		if len(evi[i].CompensationVoltage) > 0 {
			hasCompVolt = true
		}

		if !hasIonMob && evi[i].IonMobility > 0 {
			hasIonMob = true
		}

		if evi[i].Purity > 0 {
			hasPurity = true
		}

		if evi[i].MSFraggerLoc != nil && len(evi[i].MSFraggerLoc.MSFragerLocalization) > 0 {
			hasLoc = true
		}

		if evi[i].SpectralSim != 0 {
			hasSpectralSim = true
		}

		if evi[i].Rtscore != 0 {
			hasRtScore = true
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

	if hasSpectralSim {
		header += "\tSpectralSim"
	}

	if hasRtScore {
		header += "\tRTScore"
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
			header = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s",
				header,
				printSet[0].Labels.Channel1.CustomName,
				printSet[0].Labels.Channel2.CustomName,
				printSet[0].Labels.Channel5.CustomName,
				printSet[0].Labels.Channel6.CustomName,
				printSet[0].Labels.Channel9.CustomName,
				printSet[0].Labels.Channel10.CustomName,
			)
		case 10:
			header = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
				header,
				printSet[0].Labels.Channel1.CustomName,
				printSet[0].Labels.Channel2.CustomName,
				printSet[0].Labels.Channel3.CustomName,
				printSet[0].Labels.Channel4.CustomName,
				printSet[0].Labels.Channel5.CustomName,
				printSet[0].Labels.Channel6.CustomName,
				printSet[0].Labels.Channel7.CustomName,
				printSet[0].Labels.Channel8.CustomName,
				printSet[0].Labels.Channel9.CustomName,
				printSet[0].Labels.Channel10.CustomName,
			)
		case 11:
			header = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
				header,
				printSet[0].Labels.Channel1.CustomName,
				printSet[0].Labels.Channel2.CustomName,
				printSet[0].Labels.Channel3.CustomName,
				printSet[0].Labels.Channel4.CustomName,
				printSet[0].Labels.Channel5.CustomName,
				printSet[0].Labels.Channel6.CustomName,
				printSet[0].Labels.Channel7.CustomName,
				printSet[0].Labels.Channel8.CustomName,
				printSet[0].Labels.Channel9.CustomName,
				printSet[0].Labels.Channel10.CustomName,
				printSet[0].Labels.Channel11.CustomName,
			)
		case 16:
			header = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
				header,
				printSet[0].Labels.Channel1.CustomName,
				printSet[0].Labels.Channel2.CustomName,
				printSet[0].Labels.Channel3.CustomName,
				printSet[0].Labels.Channel4.CustomName,
				printSet[0].Labels.Channel5.CustomName,
				printSet[0].Labels.Channel6.CustomName,
				printSet[0].Labels.Channel7.CustomName,
				printSet[0].Labels.Channel8.CustomName,
				printSet[0].Labels.Channel9.CustomName,
				printSet[0].Labels.Channel10.CustomName,
				printSet[0].Labels.Channel11.CustomName,
				printSet[0].Labels.Channel12.CustomName,
				printSet[0].Labels.Channel13.CustomName,
				printSet[0].Labels.Channel14.CustomName,
				printSet[0].Labels.Channel15.CustomName,
				printSet[0].Labels.Channel16.CustomName,
			)
		case 18:
			header = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
				header,
				printSet[0].Labels.Channel1.CustomName,
				printSet[0].Labels.Channel2.CustomName,
				printSet[0].Labels.Channel3.CustomName,
				printSet[0].Labels.Channel4.CustomName,
				printSet[0].Labels.Channel5.CustomName,
				printSet[0].Labels.Channel6.CustomName,
				printSet[0].Labels.Channel7.CustomName,
				printSet[0].Labels.Channel8.CustomName,
				printSet[0].Labels.Channel9.CustomName,
				printSet[0].Labels.Channel10.CustomName,
				printSet[0].Labels.Channel11.CustomName,
				printSet[0].Labels.Channel12.CustomName,
				printSet[0].Labels.Channel13.CustomName,
				printSet[0].Labels.Channel14.CustomName,
				printSet[0].Labels.Channel15.CustomName,
				printSet[0].Labels.Channel16.CustomName,
				printSet[0].Labels.Channel17.CustomName,
				printSet[0].Labels.Channel18.CustomName,
			)
		default:
			header += ""
		}
	} else if brand == "itraq" {
		switch channels {
		case 4:
			header = fmt.Sprintf("%s\t%s\t%s\t%s\t%s",
				header,
				printSet[0].Labels.Channel1.CustomName,
				printSet[0].Labels.Channel2.CustomName,
				printSet[0].Labels.Channel3.CustomName,
				printSet[0].Labels.Channel4.CustomName,
			)
		case 8:
			header = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
				header,
				printSet[0].Labels.Channel1.CustomName,
				printSet[0].Labels.Channel2.CustomName,
				printSet[0].Labels.Channel3.CustomName,
				printSet[0].Labels.Channel4.CustomName,
				printSet[0].Labels.Channel5.CustomName,
				printSet[0].Labels.Channel6.CustomName,
				printSet[0].Labels.Channel7.CustomName,
				printSet[0].Labels.Channel8.CustomName,
			)
		default:
			header += ""
		}
	} else if brand == "xtag" {
		header = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
			header,
			printSet[0].Labels.Channel1.CustomName,
			printSet[0].Labels.Channel2.CustomName,
			printSet[0].Labels.Channel3.CustomName,
			printSet[0].Labels.Channel4.CustomName,
			printSet[0].Labels.Channel5.CustomName,
			printSet[0].Labels.Channel6.CustomName,
			printSet[0].Labels.Channel7.CustomName,
			printSet[0].Labels.Channel8.CustomName,
			printSet[0].Labels.Channel9.CustomName,
			printSet[0].Labels.Channel10.CustomName,
			printSet[0].Labels.Channel11.CustomName,
			printSet[0].Labels.Channel12.CustomName,
			printSet[0].Labels.Channel13.CustomName,
			printSet[0].Labels.Channel14.CustomName,
			printSet[0].Labels.Channel15.CustomName,
			printSet[0].Labels.Channel16.CustomName,
			printSet[0].Labels.Channel17.CustomName,
			printSet[0].Labels.Channel18.CustomName,
		)
	}

	header += "\n"

	_, e = io.WriteString(bw, header)
	if e != nil {
		msg.WriteToFile(errors.New("cannot print PSM to file"), "fatal")
	}

	for _, i := range printSet {

		assL, obs := getModsList(i.Modifications.ToMap().Index)

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

		// append decoy tags on the gene and proteinID names
		if i.IsDecoy {
			i.ProteinID = decoyTag + i.ProteinID
			i.GeneName = decoyTag + i.GeneName
			i.EntryName = decoyTag + i.EntryName
		}

		line := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%d\t%d\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
			i.Spectrum,
			i.SpectrumFile,
			i.Peptide,
			i.ModifiedPeptide,
			string(i.PrevAA),
			string(i.NextAA),
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

		if hasSpectralSim {
			line = fmt.Sprintf("%s\t%.4f",
				line,
				i.SpectralSim,
			)
		}

		if hasRtScore {
			line = fmt.Sprintf("%s\t%.4f",
				line,
				i.Rtscore,
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
				PTM := i.PTM
				if PTM == nil {
					PTM = &id.PTM{LocalizedPTMSites: map[string]int{}, LocalizedPTMMassDiff: map[string]string{}}
				}
				matches := r.FindAllString(PTM.LocalizedPTMMassDiff[j], -1)
				max := uti.GetMaxNumber(matches)

				line = fmt.Sprintf("%s\t%s\t%s",
					line,
					PTM.LocalizedPTMMassDiff[j],
					max,
				)
			}
		}

		if hasLoc {
			MSFraggerLoc := i.MSFraggerLoc
			if MSFraggerLoc == nil {
				MSFraggerLoc = &id.MSFraggerLoc{}
			}
			line = fmt.Sprintf("%s\t%s\t%s\t%s",
				line,
				MSFraggerLoc.MSFragerLocalization,
				MSFraggerLoc.MSFraggerLocalizationScoreWithPTM,
				MSFraggerLoc.MSFraggerLocalizationScoreWithoutPTM)
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
		} else if brand == "xtag" {
			line = fmt.Sprintf("%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
				line,
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
		}
		line += "\n"

		_, e = io.WriteString(bw, line)
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
		if i.PTM != nil {
			for j := range i.PTM.LocalizedPTMMassDiff {
				line := fmt.Sprintf("%s\t%s\t%s\t%d\t%.4f\t%s\t%d\t%s\n",
					i.SpectrumFileName().Str(),
					i.Peptide,
					i.ModifiedPeptide,
					i.AssumedCharge,
					i.RetentionTime,
					j,
					i.PTM.LocalizedPTMSites[j],
					i.PTM.LocalizedPTMMassDiff[j],
				)
				_, e = io.WriteString(file, line)
				if e != nil {
					msg.WriteToFile(e, "fatal")
				}
			}
		}
	}
}
