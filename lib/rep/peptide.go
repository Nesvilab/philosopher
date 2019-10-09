package rep

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/prvst/philosopher/lib/cla"
	"github.com/prvst/philosopher/lib/id"
	"github.com/prvst/philosopher/lib/mod"
	"github.com/prvst/philosopher/lib/msg"
	"github.com/prvst/philosopher/lib/sys"
)

// AssemblePeptideReport reports consist on ion reporting
func (evi *Evidence) AssemblePeptideReport(pep id.PepIDList, decoyTag string) {

	var list PeptideEvidenceList
	var pepSeqMap = make(map[string]bool) //is this a decoy
	var pepCSMap = make(map[string][]uint8)
	var pepInt = make(map[string]float64)
	var pepProt = make(map[string]string)
	var spectra = make(map[string][]string)
	var mappedProts = make(map[string][]string)
	var bestProb = make(map[string]float64)
	var pepMods = make(map[string][]mod.Modification)

	for _, i := range pep {
		if !cla.IsDecoyPSM(i, decoyTag) {
			pepSeqMap[i.Peptide] = false
		} else {
			pepSeqMap[i.Peptide] = true
		}
	}

	for _, i := range evi.PSM {

		_, ok := pepSeqMap[i.Peptide]
		if ok {

			pepCSMap[i.Peptide] = append(pepCSMap[i.Peptide], i.AssumedCharge)
			spectra[i.Peptide] = append(spectra[i.Peptide], i.Spectrum)
			pepProt[i.Peptide] = i.Protein

			if i.Intensity > pepInt[i.Peptide] {
				pepInt[i.Peptide] = i.Intensity
			}

			for j := range i.MappedProteins {
				mappedProts[i.Peptide] = append(mappedProts[i.Peptide], j)
			}

			for _, j := range i.Modifications.Index {
				pepMods[i.Peptide] = append(pepMods[i.Peptide], j)
			}

		}

		if i.Probability > bestProb[i.Peptide] {
			bestProb[i.Peptide] = i.Probability
		}

	}

	for k, v := range pepSeqMap {

		var pep PeptideEvidence
		pep.Spectra = make(map[string]uint8)
		pep.ChargeState = make(map[uint8]uint8)
		pep.MappedProteins = make(map[string]int)
		pep.Modifications.Index = make(map[string]mod.Modification)

		pep.Sequence = k

		pep.Probability = bestProb[k]

		for _, i := range spectra[k] {
			pep.Spectra[i] = 0
		}

		for _, i := range pepCSMap[k] {
			pep.ChargeState[i] = 0
		}

		for _, i := range mappedProts[k] {
			pep.MappedProteins[i] = 0
		}

		d, ok := pepProt[k]
		if ok {
			pep.Protein = d
		}

		mods, ok := pepMods[pep.Sequence]
		if ok {
			for _, j := range mods {
				_, okMod := pep.Modifications.Index[j.Index]
				if !okMod {
					pep.Modifications.Index[j.Index] = j
				}
			}
		}

		pep.Spc = len(spectra[k])
		pep.Intensity = pepInt[k]

		// is this a decoy ?
		pep.IsDecoy = v

		list = append(list, pep)
	}

	sort.Sort(list)
	evi.Peptides = list

	return
}

// MetaPeptideReport report consist on ion reporting
func (evi Evidence) MetaPeptideReport(labels map[string]string, brand string, channels int, hasDecoys bool) {

	var header string
	output := fmt.Sprintf("%s%speptide.tsv", sys.MetaDir(), string(filepath.Separator))

	file, e := os.Create(output)
	if e != nil {
		msg.WriteFile(errors.New("peptide output file"), "fatal")
	}
	defer file.Close()

	// building the printing set tat may or not contain decoys
	var printSet PeptideEvidenceList
	for _, i := range evi.Peptides {
		if hasDecoys == false {
			if i.IsDecoy == false {
				printSet = append(printSet, i)
			}
		} else {
			printSet = append(printSet, i)
		}
	}

	header = "Peptide\tCharges\tProbability\tSpectral Count\tIntensity\tAssigned Modifications\tObserved Modifications\tProtein\tProtein ID\tEntry Name\tGene\tProtein Description\tMapped Proteins"

	if brand == "tmt" {
		switch channels {
		case 10:
			header += "\tChannel 126\tChannel 127N\tChannel 127C\tChannel 128N\tChannel 128C\tChannel 129N\tChannel 129C\tChannel 130N\tChannel 130C\tChannel 131N"
		case 11:
			header += "\tChannel 126\tChannel 127N\tChannel 127C\tChannel 128N\tChannel 128C\tChannel 129N\tChannel 129C\tChannel 130N\tChannel 130C\tChannel 131N\tChannel 131C"
		case 16:
			header += "\tChannel 126\tChannel 127N\tChannel 127C\tChannel 128N\tChannel 128C\tChannel 129N\tChannel 129C\tChannel 130N\tChannel 130C\tChannel 131N\tChannel 131C\tChannel 132N\tChannel 132C\tChannel 133N\tChannel 133C\tChannel 134N"
		default:
			header += ""
		}
	}

	header += "\n"

	if len(labels) > 0 {
		for k, v := range labels {
			k = fmt.Sprintf("Channel %s", k)
			header = strings.Replace(header, k, v, -1)
		}
	}

	_, e = io.WriteString(file, header)
	if e != nil {
		msg.WriteToFile(errors.New("Cannot print PSM to file"), "fatal")
	}

	for _, i := range printSet {

		assL, obs := getModsList(i.Modifications.Index)

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

		sort.Strings(mappedProteins)
		sort.Strings(assL)
		sort.Strings(obs)
		sort.Strings(cs)

		line := fmt.Sprintf("%s\t%s\t%.4f\t%d\t%f\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
			i.Sequence,
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
			strings.Join(mappedProteins, ", "),
		)

		if brand == "tmt" {
			switch channels {
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
					//i.Labels.Channel12.Intensity,
					//i.Labels.Channel13.Intensity,
					//i.Labels.Channel14.Intensity,
					//i.Labels.Channel15.Intensity,
					//i.Labels.Channel16.Intensity,
				)
			default:
				header += ""
			}
		}

		line += "\n"

		_, e = io.WriteString(file, line)
		if e != nil {
			msg.WriteToFile(errors.New("Cannot print Peptides to file"), "fatal")
		}
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}
