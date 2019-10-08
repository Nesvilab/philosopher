package rep

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/prvst/philosopher/lib/msg"

	"github.com/prvst/philosopher/lib/bio"
	"github.com/prvst/philosopher/lib/cla"
	"github.com/prvst/philosopher/lib/dat"
	"github.com/prvst/philosopher/lib/id"
	"github.com/prvst/philosopher/lib/sys"
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
		p.Scan = i.Scan
		p.PrevAA = i.PrevAA
		p.NextAA = i.NextAA
		p.Peptide = i.Peptide
		p.IonForm = fmt.Sprintf("%s#%d#%.4f", i.Peptide, i.AssumedCharge, i.CalcNeutralPepMass)
		p.Protein = i.Protein
		p.ModifiedPeptide = i.ModifiedPeptide
		p.AssumedCharge = i.AssumedCharge
		p.HitRank = i.HitRank
		p.PrecursorNeutralMass = i.PrecursorNeutralMass
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
		p.MappedProteins = make(map[string]int)
		p.Modifications = i.Modifications

		for _, j := range i.AlternativeProteins {
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

		list = append(list, p)
	}

	sort.Sort(list)
	evi.PSM = list

	return
}

// PSMReport report all psms from study that passed the FDR filter
func (evi *Evidence) PSMReport(decoyTag string, hasRazor, hasDecoys bool) {

	output := fmt.Sprintf("%s%spsm.tsv", sys.MetaDir(), string(filepath.Separator))

	// create result file
	file, e := os.Create(output)
	if e != nil {
		msg.WriteFile(errors.New("Cannot create report file"), "fatal")
	}
	defer file.Close()

	_, e = io.WriteString(file, "Spectrum\tPeptide\tModified Peptide\tCharge\tRetention\tCalculated M/Z\tObserved M/Z\tDelta Mass\tExperimental Mass\tPeptide Mass\tXCorr\tDeltaCN\tDeltaCNStar\tSPScore\tSPRank\tExpectation\tHyperscore\tNextscore\tPeptideProphet Probability\tIntensity\tAssigned Modifications\tObserved Modifications\tNumber of Phospho Sites\tPhospho Site Localization\tIs Unique\tProtein\tProtein ID\tEntry Name\tGene\tProtein Description\tMapped Proteins\n")
	if e != nil {
		msg.WriteToFile(errors.New("Cannot print PSM to file"), "fatal")
	}

	// building the printing set tat may or not contain decoys
	var printSet PSMEvidenceList
	for _, i := range evi.PSM {
		if hasDecoys == false {
			if i.IsDecoy == false {
				printSet = append(printSet, i)
			}
		} else {
			printSet = append(printSet, i)
		}
	}

	for _, i := range printSet {

		assL, obs := getModsList(i.Modifications.Index)

		var mappedProteins []string
		for j := range i.MappedProteins {
			if j != i.Protein {
				mappedProteins = append(mappedProteins, j)
			}
		}

		sort.Strings(mappedProteins)
		sort.Strings(assL)
		sort.Strings(obs)

		line := fmt.Sprintf("%s\t%s\t%s\t%d\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%e\t%.4f\t%.4f\t%.4f\t%.4f\t%s\t%s\t%d\t%s\t%t\t%s\t%s\t%s\t%s\t%s\t%s\n",
			i.Spectrum,
			i.Peptide,
			i.ModifiedPeptide,
			i.AssumedCharge,
			i.RetentionTime,
			((i.CalcNeutralPepMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge)),
			((i.PrecursorNeutralMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge)),
			i.Massdiff,
			i.PrecursorNeutralMass,
			i.CalcNeutralPepMass,
			i.Xcorr,
			i.DeltaCN,
			i.DeltaCNStar,
			i.SPScore,
			i.SPRank,
			i.Expectation,
			i.Hyperscore,
			i.Nextscore,
			i.Probability,
			i.Intensity,
			strings.Join(assL, ", "),
			strings.Join(obs, ", "),
			i.LocalizedPTMSites["STY:79.966331"],
			i.LocalizedPTMMassDiff["STY:79.966331"],
			i.IsUnique,
			i.Protein,
			i.ProteinID,
			i.EntryName,
			i.GeneName,
			i.ProteinDescription,
			strings.Join(mappedProteins, ", "),
		)
		_, e = io.WriteString(file, line)
		if e != nil {
			msg.WriteToFile(e, "fatal")
		}
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}

// PSMTMTReport report all psms with TMT labels from study that passed the FDR filter
func (evi *Evidence) PSMTMTReport(labels map[string]string, decoyTag string, hasRazor, hasDecoys bool) {

	output := fmt.Sprintf("%s%spsm.tsv", sys.MetaDir(), string(filepath.Separator))

	// create result file
	file, e := os.Create(output)
	if e != nil {
		msg.WriteFile(e, "fatal")
	}
	defer file.Close()

	header := "Spectrum\tPeptide\tModified Peptide\tCharge\tRetention\tCalculated M/Z\tObserved M/Z\tDelta Mass\tExperimental Mass\tPeptide Mass\tXCorr\tDeltaCN\tDeltaCNStar\tSPScore\tSPRank\tExpectation\tHyperscore\tNextscore\tPeptideProphet Probability\tIntensity\tIs Unique\tIs Used\tAssigned Modifications\tObserved Modifications\tNumber of Phospho Sites\tPhospho Site Localization\tProtein\tProtein ID\tEntry Name\tGene\tProtein Description\tMapped Proteins\tPurity\t126 Abundance\t127N Abundance\t127C Abundance\t128N Abundance\t128C Abundance\t129N Abundance\t129C Abundance\t130N Abundance\t130C Abundance\t131N Abundance\t131C Abundance\n"

	if len(labels) > 0 {
		for k, v := range labels {
			header = strings.Replace(header, k, v, -1)
		}
	}

	_, e = io.WriteString(file, header)
	if e != nil {
		msg.WriteToFile(e, "fatal")
	}

	// building the printing set tat may or not contain decoys
	var printSet PSMEvidenceList
	for _, i := range evi.PSM {
		if hasDecoys == false {
			if i.IsDecoy == false {
				printSet = append(printSet, i)
			}
		} else {
			printSet = append(printSet, i)
		}
	}

	for _, i := range printSet {

		assL, obs := getModsList(i.Modifications.Index)

		var mappedProteins []string
		for j := range i.MappedProteins {
			if j != i.Protein {
				mappedProteins = append(mappedProteins, j)
			}
		}

		sort.Strings(mappedProteins)
		sort.Strings(assL)
		sort.Strings(obs)

		line := fmt.Sprintf("%s\t%s\t%s\t%d\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%e\t%.4f\t%.4f\t%.4f\t%.4f\t%t\t%t\t%s\t%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%.2f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\n",
			i.Spectrum,
			i.Peptide,
			i.ModifiedPeptide,
			i.AssumedCharge,
			i.RetentionTime,
			((i.CalcNeutralPepMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge)),
			((i.PrecursorNeutralMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge)),
			i.Massdiff,
			i.PrecursorNeutralMass,
			i.CalcNeutralPepMass,
			i.Xcorr,
			i.DeltaCN,
			i.DeltaCNStar,
			i.SPScore,
			i.SPRank,
			i.Expectation,
			i.Hyperscore,
			i.Nextscore,
			i.Probability,
			i.Intensity,
			i.IsUnique,
			i.Labels.IsUsed,
			strings.Join(assL, ", "),
			strings.Join(obs, ", "),
			i.LocalizedPTMSites["STY:79.966331"],
			i.LocalizedPTMMassDiff["STY:79.966331"],
			i.Protein,
			i.ProteinID,
			i.EntryName,
			i.GeneName,
			i.ProteinDescription,
			strings.Join(mappedProteins, ", "),
			i.Purity,
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
		_, e = io.WriteString(file, line)
		if e != nil {
			msg.WriteToFile(e, "fatal")
		}
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}

// PSMFraggerReport report all psms from study that passed the FDR filter
func (evi *Evidence) PSMFraggerReport(decoyTag string, hasRazor, hasDecoys bool) {

	output := fmt.Sprintf("%s%spsm.tsv", sys.MetaDir(), string(filepath.Separator))

	// create result file
	file, e := os.Create(output)
	if e != nil {
		msg.WriteFile(e, "fatal")
	}
	defer file.Close()

	_, e = io.WriteString(file, "Spectrum\tPeptide\tModified Peptide\tCharge\tRetention\tCalculated M/Z\tObserved M/Z\tDelta Mass\tExperimental Mass\tPeptide Mass\tExpectation\tHyperscore\tNextscore\tPeptideProphet Probability\tIntensity\tAssigned Modifications\tObserved Modifications\tNumber of Phospho Sites\tPhospho Site Localization\tIs Unique\tProtein\tProtein ID\tEntry Name\tGene\tProtein Description\tMapped Proteins\n")
	if e != nil {
		msg.WriteToFile(e, "fatal")
	}

	// building the printing set tat may or not contain decoys
	var printSet PSMEvidenceList
	for _, i := range evi.PSM {
		if hasDecoys == false {
			if i.IsDecoy == false {
				printSet = append(printSet, i)
			}
		} else {
			printSet = append(printSet, i)
		}
	}

	for _, i := range printSet {

		assL, obs := getModsList(i.Modifications.Index)

		var mappedProteins []string
		for j := range i.MappedProteins {
			if j != i.Protein {
				mappedProteins = append(mappedProteins, j)
			}
		}

		sort.Strings(mappedProteins)
		sort.Strings(assL)
		sort.Strings(obs)

		line := fmt.Sprintf("%s\t%s\t%s\t%d\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%e\t%.4f\t%.4f\t%.4f\t%.4f\t%s\t%s\t%d\t%s\t%t\t%s\t%s\t%s\t%s\t%s\t%s\n",
			i.Spectrum,
			i.Peptide,
			i.ModifiedPeptide,
			i.AssumedCharge,
			i.RetentionTime,
			((i.CalcNeutralPepMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge)),
			((i.PrecursorNeutralMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge)),
			i.Massdiff,
			i.PrecursorNeutralMass,
			i.CalcNeutralPepMass,
			i.Expectation,
			i.Hyperscore,
			i.Nextscore,
			i.Probability,
			i.Intensity,
			strings.Join(assL, ", "),
			strings.Join(obs, ", "),
			i.LocalizedPTMSites["STY:79.966331"],
			i.LocalizedPTMMassDiff["STY:79.966331"],
			i.IsUnique,
			i.Protein,
			i.ProteinID,
			i.EntryName,
			i.GeneName,
			i.ProteinDescription,
			strings.Join(mappedProteins, ", "),
		)
		_, e = io.WriteString(file, line)
		if e != nil {
			msg.WriteToFile(e, "fatal")
		}
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}

// PSMTMTFraggerReport report all psms with TMT labels from study that passed the FDR filter
func (evi *Evidence) PSMTMTFraggerReport(labels map[string]string, decoyTag string, hasRazor, hasDecoys bool) {

	output := fmt.Sprintf("%s%spsm.tsv", sys.MetaDir(), string(filepath.Separator))

	// create result file
	file, e := os.Create(output)
	if e != nil {
		msg.WriteFile(e, "fatal")
	}
	defer file.Close()

	header := "Spectrum\tPeptide\tModified Peptide\tCharge\tRetention\tCalculated M/Z\tObserved M/Z\tDelta Mass\tExperimental Mass\tPeptide Mass\tExpectation\tHyperscore\tNextscore\tPeptideProphet Probability\tIntensity\tIs Unique\tIs Used\tAssigned Modifications\tObserved Modifications\tNumber of Phospho Sites\tPhospho Site Localization\tProtein\tProtein ID\tEntry Name\tGene\tProtein Description\tMapped Proteins\tPurity\t126 Abundance\t127N Abundance\t127C Abundance\t128N Abundance\t128C Abundance\t129N Abundance\t129C Abundance\t130N Abundance\t130C Abundance\t131N Abundance\t131C Abundance\n"

	if len(labels) > 0 {
		for k, v := range labels {
			header = strings.Replace(header, k, v, -1)
		}
	}

	_, e = io.WriteString(file, header)
	if e != nil {
		msg.WriteToFile(e, "fatal")
	}

	// building the printing set tat may or not contain decoys
	var printSet PSMEvidenceList
	for _, i := range evi.PSM {
		if hasDecoys == false {
			if i.IsDecoy == false {
				printSet = append(printSet, i)
			}
		} else {
			printSet = append(printSet, i)
		}
	}

	for _, i := range printSet {

		assL, obs := getModsList(i.Modifications.Index)

		var mappedProteins []string
		for j := range i.MappedProteins {
			if j != i.Protein {
				mappedProteins = append(mappedProteins, j)
			}
		}

		sort.Strings(mappedProteins)
		sort.Strings(assL)
		sort.Strings(obs)

		line := fmt.Sprintf("%s\t%s\t%s\t%d\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%e\t%.4f\t%.4f\t%.4f\t%.4f\t%t\t%t\t%s\t%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%.2f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\n",
			i.Spectrum,
			i.Peptide,
			i.ModifiedPeptide,
			i.AssumedCharge,
			i.RetentionTime,
			((i.CalcNeutralPepMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge)),
			((i.PrecursorNeutralMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge)),
			i.Massdiff,
			i.PrecursorNeutralMass,
			i.CalcNeutralPepMass,
			i.Expectation,
			i.Hyperscore,
			i.Nextscore,
			i.Probability,
			i.Intensity,
			i.IsUnique,
			i.Labels.IsUsed,
			strings.Join(assL, ", "),
			strings.Join(obs, ", "),
			i.LocalizedPTMSites["STY:79.966331"],
			i.LocalizedPTMMassDiff["STY:79.966331"],
			i.Protein,
			i.ProteinID,
			i.EntryName,
			i.GeneName,
			i.ProteinDescription,
			strings.Join(mappedProteins, ", "),
			i.Purity,
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
		_, e = io.WriteString(file, line)
		if e != nil {
			msg.WriteToFile(e, "fatal")
		}
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}

// PSMLocalizationReport report ptm localization based on PTMProphet outputs
func (evi *Evidence) PSMLocalizationReport(decoyTag string, hasRazor, hasDecoys bool) {

	output := fmt.Sprintf("%s%slocalization.tsv", sys.MetaDir(), string(filepath.Separator))

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
		if hasDecoys == false {
			if i.IsDecoy == false {
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

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}
