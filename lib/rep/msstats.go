package rep

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/prvst/philosopher/lib/bio"
	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/sys"
)

// MSstatsReport report all psms from study that passed the FDR filter
func (evi *Evidence) MSstatsReport(decoyTag string, hasRazor bool) {

	output := fmt.Sprintf("%s%smsstats.csv", sys.MetaDir(), string(filepath.Separator))

	// create result file
	file, e := os.Create(output)
	if e != nil {
		err.WriteFile(e, "fatal")
	}
	defer file.Close()

	_, e = io.WriteString(file, "Spectrum.Name\tSpectrum.File\tPeptide.Sequence\tModified.Peptide.Sequence\tCharge\tCalculated.MZ\tPeptideProphet.Probability\tIntensity\tIs.Unique\tGene\tProtein.Accessions\tModifications\n")
	if e != nil {
		err.WriteToFile(errors.New("Cannot write to MSstats report"), "fatal")
	}

	// building the printing set tat may or not contain decoys
	var printSet PSMEvidenceList

	for _, i := range evi.PSM {
		if hasRazor == true {

			if i.IsURazor == true {
				if evi.Decoys == false {
					if i.IsDecoy == false && len(i.Protein) > 0 && !strings.HasPrefix(i.Protein, decoyTag) {
						printSet = append(printSet, i)
					}
				} else {
					printSet = append(printSet, i)
				}
			}

		} else {

			if evi.Decoys == false {
				if i.IsDecoy == false && len(i.Protein) > 0 && !strings.HasPrefix(i.Protein, decoyTag) {
					printSet = append(printSet, i)
				}
			} else {
				printSet = append(printSet, i)
			}

		}
	}

	for _, i := range printSet {

		var fileName string
		parts := strings.Split(i.Spectrum, ".")
		fileName = fmt.Sprintf("%s.raw", parts[0])

		line := fmt.Sprintf("%s\t%s\t%s\t%s\t%d\t%.4f\t%.4f\t%.4f\t%t\t%s\t%s\n",
			i.Spectrum,
			fileName,
			i.Peptide,
			i.ModifiedPeptide,
			i.AssumedCharge,
			((i.CalcNeutralPepMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge)),
			i.Probability,
			i.Intensity,
			i.IsUnique,
			i.GeneName,
			i.Protein,
		)
		_, e = io.WriteString(file, line)
		if e != nil {
			err.WriteToFile(errors.New("Cannot write to MSstats report"), "fatal")
		}
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}

// MSstatsTMTReport report all psms with TMT labels from study that passed the FDR filter
func (evi *Evidence) MSstatsTMTReport(labels map[string]string, decoyTag string, hasRazor bool) {

	output := fmt.Sprintf("%s%smsstats.csv", sys.MetaDir(), string(filepath.Separator))

	// create result file
	file, e := os.Create(output)
	if e != nil {
		err.WriteFile(errors.New("Cannot create report MSstats TMT file"), "fatal")
	}
	defer file.Close()

	header := "Spectrum.Name\tSpectrum.File\tPeptide.Sequence\tModified.Peptide.Sequence\tCharge.State\tCalculated.MZ\tPeptideProphet.Probability\tIntensity\tIs.Unique\tGene\tProtein\tPurity\t126.Abundance\t127N.Abundance\t127C.Abundance\t128N.Abundance\t128C.Abundance\t129N.Abundance\t129C.Abundance\t130N.Abundance\t130C.Abundance\t131N.Abundance\t131C.Abundance\n"

	if len(labels) > 0 {
		for k, v := range labels {
			header = strings.Replace(header, k, v, -1)
		}
	}

	_, e = io.WriteString(file, header)
	if e != nil {
		err.WriteToFile(errors.New("Cannot write to MSstats report"), "fatal")
	}

	// building the printing set tat may or not contain decoys
	var printSet PSMEvidenceList
	for _, i := range evi.PSM {
		if hasRazor == true {

			if i.IsURazor == true {
				if evi.Decoys == false {
					if i.IsDecoy == false && len(i.Protein) > 0 && !strings.HasPrefix(i.Protein, decoyTag) {
						printSet = append(printSet, i)
					}
				} else {
					printSet = append(printSet, i)
				}
			}

		} else {

			if evi.Decoys == false {
				if i.IsDecoy == false && len(i.Protein) > 0 && !strings.HasPrefix(i.Protein, decoyTag) {
					printSet = append(printSet, i)
				}
			} else {
				printSet = append(printSet, i)
			}

		}
	}

	for _, i := range printSet {

		var fileName string
		parts := strings.Split(i.Spectrum, ".")
		fileName = fmt.Sprintf("%s.raw", parts[0])

		line := fmt.Sprintf("%s\t%s\t%s\t%s\t%d\t%.4f\t%.4f\t%.4f\t%t\t%s\t%s\t%.2f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\n",
			i.Spectrum,
			fileName,
			i.Peptide,
			i.ModifiedPeptide,
			i.AssumedCharge,
			((i.CalcNeutralPepMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge)),
			i.Probability,
			i.Intensity,
			i.IsUnique,
			i.GeneName,
			i.Protein,
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
			err.WriteToFile(errors.New("Cannot write to MSstats report"), "fatal")
		}
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}
