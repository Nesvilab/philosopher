package rep

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"philosopher/lib/msg"

	"philosopher/lib/bio"
	"philosopher/lib/cla"
	"philosopher/lib/id"
	"philosopher/lib/mod"
	"philosopher/lib/uti"
)

// AssembleIonReport reports consist on ion reporting
func (evi *Evidence) AssembleIonReport(ion id.PepIDList, decoyTag string) {

	var psmPtMap = make(map[id.IonFormType][]string)
	var psmIonMap = make(map[id.IonFormType][]id.SpectrumType)
	var bestProb = make(map[id.IonFormType]float64)

	var ionMods = make(map[id.IonFormType][]mod.Modification)

	// collapse all psm to protein based on Peptide-level identifications
	for _, i := range evi.PSM {

		psmIonMap[i.IonForm()] = append(psmIonMap[i.IonForm()], i.SpectrumFileName())
		psmPtMap[i.IonForm()] = append(psmPtMap[i.IonForm()], i.Protein)

		if i.Probability > bestProb[i.IonForm()] {
			bestProb[i.IonForm()] = i.Probability
		}

		for j := range i.MappedProteins {
			psmPtMap[i.IonForm()] = append(psmPtMap[i.IonForm()], j)
		}

		for _, j := range i.Modifications.IndexSlice {
			ionMods[i.IonForm()] = append(ionMods[i.IonForm()], j)
		}

	}

	evi.Ions = make(IonEvidenceList, len(ion))
	for idx, i := range ion {
		pr := &evi.Ions[idx]

		//pr.IonForm() = fmt.Sprintf("%s#%d#%.4f", i.Peptide, i.AssumedCharge, i.CalcNeutralPepMass)

		pr.Spectra = make(map[id.SpectrumType]int)
		pr.MappedGenes = make(map[string]struct{})
		pr.MappedProteins = make(map[string]int)
		//pr.Modifications.Index = make(map[string]mod.Modification)

		pr.Sequence = i.Peptide
		pr.ModifiedSequence = i.ModifiedPeptide
		pr.MZ = uti.Round(((i.CalcNeutralPepMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge)), 5, 4)
		pr.ChargeState = i.AssumedCharge
		pr.PeptideMass = i.CalcNeutralPepMass
		if v, ok := psmIonMap[pr.IonForm()]; ok {
			for _, j := range v {
				pr.Spectra[j]++
			}
		}
		pr.PrecursorNeutralMass = i.PrecursorNeutralMass
		pr.Expectation = i.Expectation
		pr.NumberOfEnzymaticTermini = i.NumberOfEnzymaticTermini
		pr.Protein = i.Protein
		pr.MappedProteins[i.Protein] = 0
		pr.Modifications = i.Modifications
		pr.Probability = bestProb[pr.IonForm()]

		// get the mapped proteins
		for _, j := range psmPtMap[pr.IonForm()] {
			pr.MappedProteins[j] = 0
		}
		prModifications := pr.Modifications.ToMap()
		if mods, ok := ionMods[pr.IonForm()]; ok {
			for _, j := range mods {
				_, okMod := prModifications.Index[j.Index]
				if !okMod {
					prModifications.Index[j.Index] = j
				}
			}
		}
		pr.Modifications = prModifications.ToSlice()
		// is this bservation a decoy ?
		if cla.IsDecoyPSM(i, decoyTag) {
			pr.IsDecoy = true
		}

	}

	sort.Sort(evi.Ions)
}

// MetaIonReport reports consist on ion reporting
func (evi IonEvidenceList) MetaIonReport(workspace, brand, decoyTag string, channels int, hasDecoys, hasLabels bool) {

	var header string
	output := fmt.Sprintf("%s%sion.tsv", workspace, string(filepath.Separator))

	file, e := os.Create(output)
	bw := bufio.NewWriter(file)
	if e != nil {
		msg.WriteFile(errors.New("peptide ion output file"), "fatal")
	}
	defer file.Close()
	defer bw.Flush()

	// building the printing set tat may or not contain decoys
	var printSet []*IonEvidence
	for idx, i := range evi {
		// This inclusion is necessary to avoid unexistent observations from being included after using the filter --mods options
		if i.Probability > 0 {
			if !hasDecoys {
				if !i.IsDecoy {
					printSet = append(printSet, &evi[idx])
				}
			} else {
				printSet = append(printSet, &evi[idx])
			}
		}
	}

	header = "Peptide Sequence\tModified Sequence\tPrev AA\tNext AA\tPeptide Length\tM/Z\tCharge\tObserved Mass\tProbability\tExpectation\tSpectral Count\tIntensity\tAssigned Modifications\tObserved Modifications\tProtein\tProtein ID\tEntry Name\tGene\tProtein Description\tMapped Genes\tMapped Proteins"

	if brand == "tmt" {
		switch channels {
		case 6:
			header += "\tChannel 126\tChannel 127N\tChannel 128C\tChannel 129N\tChannel 130C\tChannel 131"
		case 10:
			header += "\tChannel 126\tChannel 127N\tChannel 127C\tChannel 128N\tChannel 128C\tChannel 129N\tChannel 129C\tChannel 130N\tChannel 130C\tChannel 131N"
		case 11:
			header += "\tChannel 126\tChannel 127N\tChannel 127C\tChannel 128N\tChannel 128C\tChannel 129N\tChannel 129C\tChannel 130N\tChannel 130C\tChannel 131N\tChannel 131C"
		case 16:
			header += "\tChannel 126\tChannel 127N\tChannel 127C\tChannel 128N\tChannel 128C\tChannel 129N\tChannel 129C\tChannel 130N\tChannel 130C\tChannel 131N\tChannel 131C\tChannel 132N\tChannel 132C\tChannel 133N\tChannel 133C\tChannel 134N"
		case 18:
			header += "\tChannel 126\tChannel 127N\tChannel 127C\tChannel 128N\tChannel 128C\tChannel 129N\tChannel 129C\tChannel 130N\tChannel 130C\tChannel 131N\tChannel 131C\tChannel 132N\tChannel 132C\tChannel 133N\tChannel 133C\tChannel 134N\tChannel 134C\tChannel 135N"
		default:
			header += ""
		}
	} else if brand == "itraq" {
		switch channels {
		case 4:
			header += "\tChannel 114\tChannel 115\tChannel 116\tChannel 117"
		case 8:
			header += "\tChannel 113\tChannel 114\tChannel 115\tChannel 116\tChannel 117\tChannel 118\tChannel 119\tChannel 121"
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

	_, e = io.WriteString(bw, header)
	if e != nil {
		msg.WriteToFile(errors.New("cannot print Ion to file"), "fatal")
	}

	for _, i := range printSet {

		assL, obs := getModsList(i.Modifications.ToMap().Index)

		var mappedProteins []string
		for j := range i.MappedProteins {
			if j != i.Protein {
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
		}

		line := fmt.Sprintf("%s\t%s\t%s\t%s\t%d\t%.4f\t%d\t%.4f\t%.4f\t%.14f\t%d\t%.4f\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
			i.Sequence,
			i.ModifiedSequence,
			string(i.PrevAA),
			string(i.NextAA),
			len(i.Sequence),
			i.MZ,
			i.ChargeState,
			i.PeptideMass,
			i.Probability,
			i.Expectation,
			len(i.Spectra),
			i.Intensity,
			strings.Join(assL, ", "),
			strings.Join(obs, ", "),
			i.Protein,
			i.ProteinID,
			i.EntryName,
			i.GeneName,
			i.ProteinDescription,
			strings.Join(mappedGenes, ","),
			strings.Join(mappedProteins, ","),
		)

		switch channels {
		case 4:
			line = fmt.Sprintf("%s\t%.4f\t%.4f\t%.4f\t%.4f",
				line,
				i.Labels.Channel1.Intensity,
				i.Labels.Channel2.Intensity,
				i.Labels.Channel3.Intensity,
				i.Labels.Channel4.Intensity,
			)
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

		line += "\n"

		_, e = io.WriteString(bw, line)
		if e != nil {
			msg.WriteToFile(errors.New("cannot print Ions to file"), "fatal")
		}
	}

}
