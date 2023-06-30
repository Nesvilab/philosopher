package rep

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Nesvilab/philosopher/lib/id"
	"github.com/Nesvilab/philosopher/lib/msg"
)

// MetaMSstatsReport report all psms from study that passed the FDR filter
func (evi Evidence) MetaMSstatsReport(workspace, brand string, channels int, hasDecoys, hasPrefix bool) {

	if evi.PSM == nil {
		RestorePSM(&evi.PSM)
	}

	var output string
	var modMap = make(map[string]string)
	var modList []string

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

		if i.PTM != nil {
			for k := range i.PTM.LocalizedPTMMassDiff {
				_, ok := modMap[k]
				if !ok {
					modMap[k] = ""
				} else {
					modMap[k] = ""
				}
			}
		}
	}

	for k := range modMap {
		modList = append(modList, k)
	}

	sort.Strings(modList)

	header := "Spectrum.Name,Spectrum.File,Peptide.Sequence,Modified.Peptide.Sequence,Probability,Charge,Protein.Start,Protein.End,Gene,Mapped.Genes,Protein,Protein.ID,Mapped.Proteins,Protein.Description,Is.Unique,Purity,Intensity"

	if len(modList) > 0 {
		for _, i := range modList {
			if strings.Contains(i, "STY:79.966331") {
				i = "STY:79.9663"
			}
			header += "," + i
		}
	}

	if brand == "tmt" {
		switch channels {
		case 6:
			header += ",Channel 126,Channel 127N,Channel 128C,Channel 129N,Channel 130C,Channel 131"
		case 10:
			header += ",Channel 126,Channel 127N,Channel 127C,Channel 128N,Channel 128C,Channel 129N,Channel 129C,Channel 130N,Channel 130C,Channel 131N"
		case 11:
			header += ",Channel 126,Channel 127N,Channel 127C,Channel 128N,Channel 128C,Channel 129N,Channel 129C,Channel 130N,Channel 130C,Channel 131N,Channel 131C"
		case 16:
			header += ",Channel 126,Channel 127N,Channel 127C,Channel 128N,Channel 128C,Channel 129N,Channel 129C,Channel 130N,Channel 130C,Channel 131N,Channel 131C,Channel 132N,Channel 132C,Channel 133N,Channel 133C,Channel 134N"
		case 18:
			header += ",Channel 126,Channel 127N,Channel 127C,Channel 128N,Channel 128C,Channel 129N,Channel 129C,Channel 130N,Channel 130C,Channel 131N,Channel 131C,Channel 132N,Channel 132C,Channel 133N,Channel 133C,Channel 134N,Channel 134C,Channel 135N"
		default:
			header += ""
		}
	} else if brand == "itraq" {
		switch channels {
		case 4:
			header += ",Channel 114,Channel 115,Channel 116,Channel 117"
		case 8:
			header += ",Channel 113,Channel 114,Channel 115,Channel 116,Channel 117,Channel 118,Channel 119,Channel 121"
		default:
			header += ""
		}
	} else if brand == "sclip" {
		header += ",Channel sCLIP1,Channel sCLIP2,Channel sCLIP3,Channel sCLIPv4,Channel sCLIP5,Channel sCLIP6"
	} else if brand == "xtag" {
		header += ",Channel xTag1,Channel xTag2,Channel xTag3,Channel xTag4,Channel xTag5,Channel xTag6,Channel xTag7,Channel xTag8,Channel xTag9,Channel xTag10,Channel xTag11,Channel xTag12,Channel xTag13,Channel xTag14,Channel xTag15,Channel xTag16,Channel xTag17,Channel xTag19,Channel xTag20,Channel xTag21,Channel xTag22,Channel xTag23,Channel xTag24,Channel xTag25,Channel xTag26,Channel xTag27,Channel xTag28,Channel xTag29,Channel xTag30,Channel xTag31,Channel xTag32"
	}

	header += "\n"

	_, e = io.WriteString(file, header)
	if e != nil {
		msg.WriteToFile(errors.New("cannot print PSM to file"), "error")
	}

	for _, i := range printSet {

		var fileName string
		parts := strings.Split(i.Spectrum, ".")
		fileName = fmt.Sprintf("%s.mzML", parts[0])

		var mappedGenes []string
		for j := range i.MappedGenes {
			mappedGenes = append(mappedGenes, strings.Replace(j, ",", ";", -1))
		}

		var mappedProteins []string
		for j := range i.MappedProteins {
			mappedProteins = append(mappedProteins, j)
		}

		line := fmt.Sprintf("%s,%s,%s,%s,%.4f,%d,%d,%d,%s,%s,%s,%s,%s,%s,%t,%.4f,%.6f",
			i.Spectrum,
			fileName,
			i.Peptide,
			i.ModifiedPeptide,
			i.Probability,
			i.AssumedCharge,
			i.ProteinStart,
			i.ProteinEnd,
			strings.Replace(i.GeneName, ",", "-", -1),
			strings.Join(mappedGenes, ";"),
			i.Protein,
			i.ProteinID,
			strings.Join(mappedProteins, ";"),
			strings.Replace(i.ProteinDescription, ",", "-", -1),
			i.IsUnique,
			i.Purity,
			i.Intensity,
		)

		if len(modList) > 0 {
			for _, j := range modList {

				PTM := i.PTM
				if PTM == nil {
					PTM = &id.PTM{LocalizedPTMSites: map[string]int{}, LocalizedPTMMassDiff: map[string]string{}}
				}

				line = fmt.Sprintf("%s,%s",
					line,
					PTM.LocalizedPTMMassDiff[j],
				)
			}
		}

		// if len(modList) > 0 {

		// 	//line += ","

		// 	for _, j := range modList {
		// 		if i.PTM != nil {
		// 			mods += fmt.Sprintf("%s,", i.PTM.LocalizedPTMMassDiff[j])
		// 		} else {
		// 			line += ","
		// 		}
		// 	}
		// }

		// line += mods

		if brand == "tmt" {

			switch channels {
			case 6:
				line = fmt.Sprintf("%s,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f",
					line,
					i.Labels.Channel1.Intensity,
					i.Labels.Channel2.Intensity,
					i.Labels.Channel5.Intensity,
					i.Labels.Channel6.Intensity,
					i.Labels.Channel9.Intensity,
					i.Labels.Channel10.Intensity,
				)
			case 10:
				line = fmt.Sprintf("%s,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f",
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
				)
			case 11:
				line = fmt.Sprintf("%s,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f",
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
				)
			case 16:
				line = fmt.Sprintf("%s,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f",
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
				)
			case 18:
				line = fmt.Sprintf("%s,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f",
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
			default:
				header += ""
			}
		} else if brand == "itraq" {
			switch channels {
			case 4:
				line = fmt.Sprintf("%s,%.4f,%.4f,%.4f,%.4f",
					line,
					i.Labels.Channel1.Intensity,
					i.Labels.Channel2.Intensity,
					i.Labels.Channel3.Intensity,
					i.Labels.Channel4.Intensity,
				)
			case 8:
				line = fmt.Sprintf("%s,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f",
					line,
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
			line = fmt.Sprintf("%s,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f",
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

		_, e = io.WriteString(file, line)
		if e != nil {
			msg.WriteToFile(errors.New("cannot write to MSstats report"), "error")
		}
	}
}
