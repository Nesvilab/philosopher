package rep

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"philosopher/lib/bio"
	"philosopher/lib/msg"
)

// MetaMSstatsReport report all psms from study that passed the FDR filter
func (evi Evidence) MetaMSstatsReport(workspace, brand string, channels int, hasDecoys, hasPrefix bool) {

	if evi.PSM == nil {
		RestorePSM(&evi.PSM)
	}

	var header string
	var output string

	if hasPrefix {
		output = fmt.Sprintf("%s%s%s_msstats.tsv", workspace, string(filepath.Separator), path.Base(workspace))
	} else {
		output = fmt.Sprintf("%s%smsstats.csv", workspace, string(filepath.Separator))
	}

	// create result file
	file, e := os.Create(output)
	if e != nil {
		msg.WriteFile(errors.New("cannot create MSstats report"), "fatal")
	}
	defer file.Close()

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

	header = "Spectrum.Name,Spectrum.File,Peptide.Sequence,Modified.Peptide.Sequence,Charge,Calculated.MZ,PeptideProphet.Probability,Intensity,Is.Unique,Gene,Protein.Accessions,Modifications"

	if brand == "tmt" {
		switch channels {
		case 6:
			header += ",Purity,Channel 126,Channel 127N,Channel 128C,Channel 129N,Channel 130C,Channel 131"
		case 10:
			header += ",Purity,Channel 126,Channel 127N,Channel 127C,Channel 128N,Channel 128C,Channel 129N,Channel 129C,Channel 130N,Channel 130C,Channel 131N"
		case 11:
			header += ",Purity,Channel 126,Channel 127N,Channel 127C,Channel 128N,Channel 128C,Channel 129N,Channel 129C,Channel 130N,Channel 130C,Channel 131N,Channel 131C"
		case 16:
			header += ",Purity,Channel 126,Channel 127N,Channel 127C,Channel 128N,Channel 128C,Channel 129N,Channel 129C,Channel 130N,Channel 130C,Channel 131N,Channel 131C,Channel 132N,Channel 132C,Channel 133N,Channel 133C,Channel 134N"
		case 18:
			header += ",Purity,Channel 126,Channel 127N,Channel 127C,Channel 128N,Channel 128C,Channel 129N,Channel 129C,Channel 130N,Channel 130C,Channel 131N,Channel 131C,Channel 132N,Channel 132C,Channel 133N,Channel 133C,Channel 134N,Channel 134C,Channel 135N"
		default:
			header += ""
		}
	} else if brand == "itraq" {
		switch channels {
		case 4:
			header += ",Purity,Channel 114,Channel 115,Channel 116,Channel 117"
		case 8:
			header += ",Purity,Channel 113,Channel 114,Channel 115,Channel 116,Channel 117,Channel 118,Channel 119,Channel 121"
		default:
			header += ""
		}
	} else if brand == "xtag" {
		header += ",Purity,Channel xTag1,Channel xTag2,Channel xTag3,Channel xTag4,Channel xTag5,Channel xTag6,Channel xTag7,Channel xTag8,Channel xTag9,Channel xTag10,Channel xTag11,Channel xTag12,Channel xTag13,Channel xTag14,Channel xTag15,Channel xTag16,Channel xTag17,Channel xTag18"
	}

	header += "\n"

	_, e = io.WriteString(file, header)
	if e != nil {
		msg.WriteToFile(errors.New("cannot print PSM to file"), "error")
	}

	for _, i := range printSet {

		var fileName string
		parts := strings.Split(i.Spectrum, ".")
		fileName = fmt.Sprintf("%s.raw", parts[0])

		line := fmt.Sprintf("%s,%s,%s,%s,%d,%.4f,%.4f,%.4f,%t,%s,%s,%s",
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
			"",
		)

		if brand == "tmt" {
			switch channels {
			case 6:
				line = fmt.Sprintf("%s,%.2f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f",
					line,
					i.Purity,
					i.Labels.Channel1.Intensity,
					i.Labels.Channel2.Intensity,
					i.Labels.Channel5.Intensity,
					i.Labels.Channel6.Intensity,
					i.Labels.Channel9.Intensity,
					i.Labels.Channel10.Intensity,
				)
			case 10:
				line = fmt.Sprintf("%s,%.2f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f",
					line,
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
				)
			case 11:
				line = fmt.Sprintf("%s,%.2f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f",
					line,
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
			case 16:
				line = fmt.Sprintf("%s,%.2f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f",
					line,
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
					i.Labels.Channel12.Intensity,
					i.Labels.Channel13.Intensity,
					i.Labels.Channel14.Intensity,
					i.Labels.Channel15.Intensity,
					i.Labels.Channel16.Intensity,
				)
			case 18:
				line = fmt.Sprintf("%s,%.2f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f",
					line,
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
				line = fmt.Sprintf("%s,%.2f,%.4f,%.4f,%.4f,%.4f",
					line,
					i.Purity,
					i.Labels.Channel1.Intensity,
					i.Labels.Channel2.Intensity,
					i.Labels.Channel3.Intensity,
					i.Labels.Channel4.Intensity,
				)
			case 8:
				line = fmt.Sprintf("%s,%.2f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f",
					line,
					i.Purity,
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
			line = fmt.Sprintf("%s,%.2f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f",
				line,
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

		_, e = io.WriteString(file, line)
		if e != nil {
			msg.WriteToFile(errors.New("cannot write to MSstats report"), "error")
		}
	}
}
