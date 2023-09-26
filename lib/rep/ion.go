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
	"strings"

	"github.com/Nesvilab/philosopher/lib/msg"

	"github.com/Nesvilab/philosopher/lib/bio"
	"github.com/Nesvilab/philosopher/lib/cla"
	"github.com/Nesvilab/philosopher/lib/id"
	"github.com/Nesvilab/philosopher/lib/mod"
	"github.com/Nesvilab/philosopher/lib/uti"
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

		ionMods[i.IonForm()] = append(ionMods[i.IonForm()], i.Modifications.IndexSlice...)

		// for _, j := range i.Modifications.IndexSlice {
		// 	ionMods[i.IonForm()] = append(ionMods[i.IonForm()], j)
		// }

	}

	evi.Ions = make(IonEvidenceList, len(ion))
	for idx, i := range ion {
		pr := &evi.Ions[idx]

		pr.Spectra = make(map[id.SpectrumType]int)
		pr.MappedGenes = make(map[string]struct{})
		pr.MappedProteins = make(map[string]int)
		pr.Sequence = i.Peptide
		pr.ModifiedSequence = i.ModifiedPeptide
		pr.MZ = uti.Round(((i.CalcNeutralPepMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge)), 5, 4)
		pr.ChargeState = i.AssumedCharge
		pr.PeptideMass = i.CalcNeutralPepMass

		pr.PrevAA = string(i.PrevAA)
		pr.NextAA = string(i.NextAA)

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

// IonReport reports consist on ion reporting
func (evi IonEvidenceList) IonReport(workspace, brand, decoyTag string, channels int, hasDecoys, hasLabels, hasPrefix, removeContam bool) {

	var header string
	var output string

	if hasPrefix {
		output = fmt.Sprintf("%s%s%s_ion.tsv", workspace, string(filepath.Separator), path.Base(workspace))
	} else {
		output = fmt.Sprintf("%s%sion.tsv", workspace, string(filepath.Separator))
	}

	file, e := os.Create(output)
	bw := bufio.NewWriter(file)
	if e != nil {
		msg.WriteFile(errors.New("peptide ion output file"), "error")
	}
	defer file.Close()
	defer bw.Flush()

	// building the printing set tat may or not contain decoys
	var printSet []*IonEvidence
	for idx, i := range evi {

		if removeContam && (strings.HasPrefix(i.Protein, "contam_") || strings.HasPrefix(i.Protein, "Cont_")) {
			continue
		}

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

	header = "Peptide Sequence\tModified Sequence\tPrev AA\tNext AA\tPeptide Length\tProtein Start\tProtein End\tM/Z\tCharge\tObserved Mass\tProbability\tExpectation\tSpectral Count\tIntensity\tAssigned Modifications\tObserved Modifications\tProtein\tProtein ID\tEntry Name\tGene\tProtein Description\tMapped Genes\tMapped Proteins"

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
		header = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s",
			header,
			printSet[headerIndex].Labels.Channel1.CustomName,
			printSet[headerIndex].Labels.Channel2.CustomName,
			printSet[headerIndex].Labels.Channel3.CustomName,
			printSet[headerIndex].Labels.Channel4.CustomName,
			printSet[headerIndex].Labels.Channel5.CustomName,
			printSet[headerIndex].Labels.Channel6.CustomName,
		)
	} else if brand == "ibt" {
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
	} else if brand == "xtag" {
		header = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
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
			printSet[headerIndex].Labels.Channel19.CustomName,
			printSet[headerIndex].Labels.Channel20.CustomName,
			printSet[headerIndex].Labels.Channel21.CustomName,
			printSet[headerIndex].Labels.Channel22.CustomName,
			printSet[headerIndex].Labels.Channel23.CustomName,
			printSet[headerIndex].Labels.Channel24.CustomName,
			printSet[headerIndex].Labels.Channel25.CustomName,
			printSet[headerIndex].Labels.Channel26.CustomName,
			printSet[headerIndex].Labels.Channel27.CustomName,
			printSet[headerIndex].Labels.Channel28.CustomName,
			printSet[headerIndex].Labels.Channel29.CustomName,
			printSet[headerIndex].Labels.Channel30.CustomName,
			printSet[headerIndex].Labels.Channel31.CustomName,
			printSet[headerIndex].Labels.Channel32.CustomName,
		)
	} else if brand == "xtag2" {

		header += "\tQuan Usage"

		header = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
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
			printSet[headerIndex].Labels.Channel19.CustomName,
			printSet[headerIndex].Labels.Channel20.CustomName,
			printSet[headerIndex].Labels.Channel21.CustomName,
			printSet[headerIndex].Labels.Channel22.CustomName,
			printSet[headerIndex].Labels.Channel23.CustomName,
			printSet[headerIndex].Labels.Channel24.CustomName,
			printSet[headerIndex].Labels.Channel25.CustomName,
			printSet[headerIndex].Labels.Channel26.CustomName,
			printSet[headerIndex].Labels.Channel27.CustomName,
			printSet[headerIndex].Labels.Channel28.CustomName,
			printSet[headerIndex].Labels.Channel29.CustomName,
		)
	}

	header += "\n"

	_, e = io.WriteString(bw, header)
	if e != nil {
		msg.WriteToFile(errors.New("cannot print Ion to file"), "error")
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
		if i.IsDecoy || strings.HasPrefix(i.Protein, decoyTag) {
			i.ProteinID = decoyTag + i.ProteinID
			i.GeneName = decoyTag + i.GeneName
			i.EntryName = decoyTag + i.EntryName
		}

		line := fmt.Sprintf("%s\t%s\t%s\t%s\t%d\t%d\t%d\t%.4f\t%d\t%.4f\t%.4f\t%.14f\t%d\t%.4f\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
			i.Sequence,
			i.ModifiedSequence,
			string(i.PrevAA),
			string(i.NextAA),
			len(i.Sequence),
			i.ProteinStart,
			i.ProteinEnd,
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
			line = fmt.Sprintf("%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
				line,
				i.Labels.Channel1.Intensity,
				i.Labels.Channel2.Intensity,
				i.Labels.Channel3.Intensity,
				i.Labels.Channel4.Intensity,
				i.Labels.Channel5.Intensity,
				i.Labels.Channel6.Intensity,
			)
		} else if brand == "ibt" {
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
		} else if brand == "xtag" {
			line = fmt.Sprintf("%s\t%t\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
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
				i.Labels.Channel19.Intensity,
				i.Labels.Channel20.Intensity,
				i.Labels.Channel21.Intensity,
				i.Labels.Channel22.Intensity,
				i.Labels.Channel23.Intensity,
				i.Labels.Channel24.Intensity,
				i.Labels.Channel25.Intensity,
				i.Labels.Channel26.Intensity,
				i.Labels.Channel27.Intensity,
				i.Labels.Channel28.Intensity,
				i.Labels.Channel29.Intensity,
				i.Labels.Channel30.Intensity,
				i.Labels.Channel31.Intensity,
				i.Labels.Channel32.Intensity,
			)
		} else if brand == "xtag2" {
			line = fmt.Sprintf("%s\t%t\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
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
				i.Labels.Channel19.Intensity,
				i.Labels.Channel20.Intensity,
				i.Labels.Channel21.Intensity,
				i.Labels.Channel22.Intensity,
				i.Labels.Channel23.Intensity,
				i.Labels.Channel24.Intensity,
				i.Labels.Channel25.Intensity,
				i.Labels.Channel26.Intensity,
				i.Labels.Channel27.Intensity,
				i.Labels.Channel28.Intensity,
				i.Labels.Channel29.Intensity,
			)
		}
		line += "\n"

		_, e = io.WriteString(bw, line)
		if e != nil {
			msg.WriteToFile(errors.New("cannot print Ions to file"), "error")
		}
	}

}
