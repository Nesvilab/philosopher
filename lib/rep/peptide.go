package rep

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/Nesvilab/philosopher/lib/cla"
	"github.com/Nesvilab/philosopher/lib/id"
	"github.com/Nesvilab/philosopher/lib/mod"
	"github.com/Nesvilab/philosopher/lib/msg"
)

// AssemblePeptideReport reports consist on ion reporting
func (evi *Evidence) AssemblePeptideReport(pep id.PepIDList, decoyTag string) {

	var pepSeqMap = make(map[string]bool) //is this a decoy
	var pepCSMap = make(map[string][]uint8)
	var pepProt = make(map[string]string)
	var mappedGenes = make(map[string][]string)
	var mappedProts = make(map[string][]string)
	var pepInt = make(map[string]float64)
	var bestProb = make(map[string]float64)
	var prevAA = make(map[string]string)
	var nextAA = make(map[string]string)
	var spectra = make(map[string][]id.SpectrumType)
	var pepMods = make(map[string][]mod.Modification)

	for _, i := range pep {
		pepSeqMap[i.Peptide] = cla.IsDecoyPSM(i, decoyTag)
	}

	for _, i := range evi.PSM {

		if _, ok := pepSeqMap[i.Peptide]; ok {

			pepCSMap[i.Peptide] = append(pepCSMap[i.Peptide], i.AssumedCharge)
			spectra[i.Peptide] = append(spectra[i.Peptide], i.SpectrumFileName())
			pepProt[i.Peptide] = i.Protein
			prevAA[i.Peptide] = i.PrevAA
			nextAA[i.Peptide] = i.NextAA

			if i.Intensity > pepInt[i.Peptide] {
				pepInt[i.Peptide] = i.Intensity
			}

			for j := range i.MappedProteins {
				mappedProts[i.Peptide] = append(mappedProts[i.Peptide], j)
			}

			for j := range i.MappedGenes {
				mappedGenes[i.Peptide] = append(mappedGenes[i.Peptide], j)
			}

			pepMods[i.Peptide] = append(pepMods[i.Peptide], i.Modifications.IndexSlice...)

			// for _, j := range i.Modifications.IndexSlice {
			// 	pepMods[i.Peptide] = append(pepMods[i.Peptide], j)
			// }

		}

		if i.Probability > bestProb[i.Peptide] {
			bestProb[i.Peptide] = i.Probability
		}

	}

	evi.Peptides = make(PeptideEvidenceList, len(pepSeqMap))
	idx := 0
	for k, v := range pepSeqMap {

		pep := &evi.Peptides[idx]
		idx++

		pep.Spectra = make(map[id.SpectrumType]uint8)
		pep.ChargeState = make(map[uint8]uint8)
		pep.MappedGenes = make(map[string]struct{})
		pep.MappedProteins = make(map[string]int)

		pep.Sequence = k

		pep.Probability = bestProb[k]

		pep.PrevAA = prevAA[k]
		pep.NextAA = nextAA[k]

		for _, i := range spectra[k] {
			pep.Spectra[i] = 0
		}

		for _, i := range pepCSMap[k] {
			pep.ChargeState[i] = 0
		}

		for _, i := range mappedGenes[k] {
			pep.MappedGenes[i] = struct{}{}
		}

		for _, i := range mappedProts[k] {
			pep.MappedProteins[i] = 0
		}

		d, ok := pepProt[k]
		if ok {
			pep.Protein = d
		}

		pepModificationsIndex := make(map[string]mod.Modification)
		if mods, ok := pepMods[pep.Sequence]; ok {
			for _, j := range mods {
				_, okMod := pepModificationsIndex[j.Index]
				if !okMod {
					pepModificationsIndex[j.Index] = j
				}
			}
		}
		if len(pepModificationsIndex) != 0 {
			pep.Modifications = mod.Modifications{Index: pepModificationsIndex}.ToSlice()
		}
		// is this a decoy ?
		pep.IsDecoy = v

	}

	sort.Sort(evi.Peptides)

}

// PeptideReport report consist on ion reporting
func (evi PeptideEvidenceList) PeptideReport(workspace, brand, decoyTag string, channels int, hasDecoys, hasLabels, hasPrefix, removeContam bool) {

	var header string
	var output string

	if hasPrefix {
		output = fmt.Sprintf("%s%s%s_peptide.tsv", workspace, string(filepath.Separator), path.Base(workspace))
	} else {
		output = fmt.Sprintf("%s%speptide.tsv", workspace, string(filepath.Separator))
	}

	file, e := os.Create(output)
	bw := bufio.NewWriter(file)
	if e != nil {
		msg.WriteFile(errors.New("peptide output file"), "error")
	}
	defer file.Close()
	defer bw.Flush()

	// building the printing set tat may or not contain decoys
	var printSet []*PeptideEvidence
	for idx, i := range evi {

		if removeContam && (strings.HasPrefix(i.Protein, "contam_") || strings.HasPrefix(i.Protein, "Cont_")) {
			continue
		}

		if !hasDecoys {
			if !i.IsDecoy {
				printSet = append(printSet, &evi[idx])
			}
		} else {
			printSet = append(printSet, &evi[idx])
		}
	}

	header = "Peptide\tPrev AA\tNext AA\tPeptide Length\tProtein Start\tProtein End\tCharges\tProbability\tSpectral Count\tIntensity\tAssigned Modifications\tObserved Modifications\tProtein\tProtein ID\tEntry Name\tGene\tProtein Description\tMapped Genes\tMapped Proteins"

	var headerIndex int
	for i := range printSet {
		if printSet[i].Labels != nil && len(printSet[i].Labels.Channel1.CustomName) > 0 {
			headerIndex = i
			break
		}
	}

	if brand == "tmt" {
		switch channels {
		case 6:
			header = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s",
				header,
				printSet[headerIndex].Labels.Channel1.CustomName,
				printSet[headerIndex].Labels.Channel2.CustomName,
				printSet[headerIndex].Labels.Channel5.CustomName,
				printSet[headerIndex].Labels.Channel6.CustomName,
				printSet[headerIndex].Labels.Channel9.CustomName,
				printSet[headerIndex].Labels.Channel10.CustomName,
			)
		case 10:
			header = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
				header,
				printSet[headerIndex].Labels.Channel1.CustomName,
				printSet[headerIndex].Labels.Channel2.CustomName,
				printSet[headerIndex].Labels.Channel3.CustomName,
				printSet[headerIndex].Labels.Channel4.CustomName,
				printSet[headerIndex].Labels.Channel5.CustomName,
				printSet[headerIndex].Labels.Channel6.CustomName,
				printSet[headerIndex].Labels.Channel7.CustomName,
				printSet[headerIndex].Labels.Channel8.CustomName,
				printSet[headerIndex].Labels.Channel9.CustomName,
				printSet[headerIndex].Labels.Channel10.CustomName,
			)
		case 11:
			header = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
				header,
				printSet[headerIndex].Labels.Channel1.CustomName,
				printSet[headerIndex].Labels.Channel2.CustomName,
				printSet[headerIndex].Labels.Channel3.CustomName,
				printSet[headerIndex].Labels.Channel4.CustomName,
				printSet[headerIndex].Labels.Channel5.CustomName,
				printSet[headerIndex].Labels.Channel6.CustomName,
				printSet[headerIndex].Labels.Channel7.CustomName,
				printSet[headerIndex].Labels.Channel8.CustomName,
				printSet[headerIndex].Labels.Channel9.CustomName,
				printSet[headerIndex].Labels.Channel10.CustomName,
				printSet[headerIndex].Labels.Channel11.CustomName,
			)
		case 16:
			header = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
				header,
				printSet[headerIndex].Labels.Channel1.CustomName,
				printSet[headerIndex].Labels.Channel2.CustomName,
				printSet[headerIndex].Labels.Channel3.CustomName,
				printSet[headerIndex].Labels.Channel4.CustomName,
				printSet[headerIndex].Labels.Channel5.CustomName,
				printSet[headerIndex].Labels.Channel6.CustomName,
				printSet[headerIndex].Labels.Channel7.CustomName,
				printSet[headerIndex].Labels.Channel8.CustomName,
				printSet[headerIndex].Labels.Channel9.CustomName,
				printSet[headerIndex].Labels.Channel10.CustomName,
				printSet[headerIndex].Labels.Channel11.CustomName,
				printSet[headerIndex].Labels.Channel12.CustomName,
				printSet[headerIndex].Labels.Channel13.CustomName,
				printSet[headerIndex].Labels.Channel14.CustomName,
				printSet[headerIndex].Labels.Channel15.CustomName,
				printSet[headerIndex].Labels.Channel16.CustomName,
			)
		case 18:
			header = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
				header,
				printSet[headerIndex].Labels.Channel1.CustomName,
				printSet[headerIndex].Labels.Channel2.CustomName,
				printSet[headerIndex].Labels.Channel3.CustomName,
				printSet[headerIndex].Labels.Channel4.CustomName,
				printSet[headerIndex].Labels.Channel5.CustomName,
				printSet[headerIndex].Labels.Channel6.CustomName,
				printSet[headerIndex].Labels.Channel7.CustomName,
				printSet[headerIndex].Labels.Channel8.CustomName,
				printSet[headerIndex].Labels.Channel9.CustomName,
				printSet[headerIndex].Labels.Channel10.CustomName,
				printSet[headerIndex].Labels.Channel11.CustomName,
				printSet[headerIndex].Labels.Channel12.CustomName,
				printSet[headerIndex].Labels.Channel13.CustomName,
				printSet[headerIndex].Labels.Channel14.CustomName,
				printSet[headerIndex].Labels.Channel15.CustomName,
				printSet[headerIndex].Labels.Channel16.CustomName,
				printSet[headerIndex].Labels.Channel17.CustomName,
				printSet[headerIndex].Labels.Channel18.CustomName,
			)
		default:
			header += ""
		}
	} else if brand == "itraq" {
		switch channels {
		case 4:
			header = fmt.Sprintf("%s\t%s\t%s\t%s\t%s",
				header,
				printSet[headerIndex].Labels.Channel1.CustomName,
				printSet[headerIndex].Labels.Channel2.CustomName,
				printSet[headerIndex].Labels.Channel3.CustomName,
				printSet[headerIndex].Labels.Channel4.CustomName,
			)
		case 8:
			header = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
				header,
				printSet[headerIndex].Labels.Channel1.CustomName,
				printSet[headerIndex].Labels.Channel2.CustomName,
				printSet[headerIndex].Labels.Channel3.CustomName,
				printSet[headerIndex].Labels.Channel4.CustomName,
				printSet[headerIndex].Labels.Channel5.CustomName,
				printSet[headerIndex].Labels.Channel6.CustomName,
				printSet[headerIndex].Labels.Channel7.CustomName,
				printSet[headerIndex].Labels.Channel8.CustomName,
			)
		default:
			header += ""
		}
	} else if brand == "sclip" {
		header = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
			header,
			printSet[headerIndex].Labels.Channel1.CustomName,
			printSet[headerIndex].Labels.Channel2.CustomName,
			printSet[headerIndex].Labels.Channel3.CustomName,
			printSet[headerIndex].Labels.Channel4.CustomName,
			printSet[headerIndex].Labels.Channel5.CustomName,
			printSet[headerIndex].Labels.Channel6.CustomName,
			printSet[headerIndex].Labels.Channel7.CustomName,
			printSet[headerIndex].Labels.Channel8.CustomName,
			printSet[headerIndex].Labels.Channel9.CustomName,
			printSet[headerIndex].Labels.Channel10.CustomName,
			printSet[headerIndex].Labels.Channel11.CustomName,
			printSet[headerIndex].Labels.Channel12.CustomName,
			printSet[headerIndex].Labels.Channel13.CustomName,
			printSet[headerIndex].Labels.Channel14.CustomName,
			printSet[headerIndex].Labels.Channel15.CustomName,
			printSet[headerIndex].Labels.Channel16.CustomName,
			printSet[headerIndex].Labels.Channel17.CustomName,
			printSet[headerIndex].Labels.Channel18.CustomName,
		)
	} else if brand == "xtag" {
		header = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
			header,
			printSet[headerIndex].Labels.Channel1.CustomName,
			printSet[headerIndex].Labels.Channel2.CustomName,
			printSet[headerIndex].Labels.Channel3.CustomName,
			printSet[headerIndex].Labels.Channel4.CustomName,
			printSet[headerIndex].Labels.Channel5.CustomName,
			printSet[headerIndex].Labels.Channel6.CustomName,
			printSet[headerIndex].Labels.Channel7.CustomName,
			printSet[headerIndex].Labels.Channel8.CustomName,
			printSet[headerIndex].Labels.Channel9.CustomName,
			printSet[headerIndex].Labels.Channel10.CustomName,
			printSet[headerIndex].Labels.Channel11.CustomName,
			printSet[headerIndex].Labels.Channel12.CustomName,
			printSet[headerIndex].Labels.Channel13.CustomName,
			printSet[headerIndex].Labels.Channel14.CustomName,
			printSet[headerIndex].Labels.Channel15.CustomName,
			printSet[headerIndex].Labels.Channel16.CustomName,
			printSet[headerIndex].Labels.Channel17.CustomName,
			printSet[headerIndex].Labels.Channel18.CustomName,
		)
	}

	header += "\n"

	//_, e = io.WriteString(file, header)
	_, e = io.WriteString(bw, header)
	if e != nil {
		msg.WriteToFile(errors.New("cannot print PSM to file"), "error")
	}

	for _, i := range printSet {

		assL, obs := getModsList(i.Modifications.ToMap().Index)

		var mappedProteins []string
		for j := range i.MappedProteins {
			if j != i.Protein {
				mappedProteins = append(mappedProteins, j)
			}
		}

		var cs []string
		for j := range i.ChargeState {
			cs = append(cs, strconv.Itoa(int(j)))
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
		sort.Strings(cs)

		// append decoy tags on the gene and proteinID names
		if i.IsDecoy || strings.HasPrefix(i.Protein, decoyTag) {
			i.ProteinID = decoyTag + i.ProteinID
			i.GeneName = decoyTag + i.GeneName
			i.EntryName = decoyTag + i.EntryName
		}

		line := fmt.Sprintf("%s\t%s\t%s\t%d\t%d\t%d\t%s\t%.4f\t%d\t%f\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
			i.Sequence,
			string(i.PrevAA),
			string(i.NextAA),
			len(i.Sequence),
			i.ProteinStart,
			i.ProteinEnd,
			strings.Join(cs, ", "),
			i.Probability,
			i.Spc,
			i.Intensity,
			strings.Join(assL, ", "),
			strings.Join(obs, ", "),
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
				line = fmt.Sprintf("%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
					line,
					i.Labels.Channel1.Intensity,
					i.Labels.Channel2.Intensity,
					i.Labels.Channel5.Intensity,
					i.Labels.Channel6.Intensity,
					i.Labels.Channel9.Intensity,
					i.Labels.Channel10.Intensity,
				)
			case 10:
				line = fmt.Sprintf("%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
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
				line = fmt.Sprintf("%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
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
				line = fmt.Sprintf("%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
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
			default:
				header += ""
			}
		} else if brand == "itraq" {
			switch channels {
			case 4:
				line = fmt.Sprintf("%s\t%.4f\t%.4f\t%.4f\t%.4f",
					line,
					i.Labels.Channel1.Intensity,
					i.Labels.Channel2.Intensity,
					i.Labels.Channel3.Intensity,
					i.Labels.Channel4.Intensity,
				)
			case 8:
				line = fmt.Sprintf("%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
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
		} else if brand == "sclip" {
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
			msg.WriteToFile(errors.New("cannot print Peptides to file"), "error")
		}
	}
}
