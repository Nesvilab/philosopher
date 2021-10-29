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

	"philosopher/lib/dat"
	"philosopher/lib/id"
	"philosopher/lib/mod"
	"philosopher/lib/msg"
)

// AssembleProteinReport creates the post processed protein strcuture
func (evi *Evidence) AssembleProteinReport(pro id.ProtIDList, weight float64, decoyTag string) {

	var list ProteinEvidenceList
	var protMods = make(map[string][]mod.Modification)
	var evidenceIons = make(map[string]IonEvidence)

	for _, i := range evi.Ions {
		evidenceIons[i.IonForm] = i
	}

	for _, i := range evi.PSM {
		for _, j := range i.Modifications.Index {
			protMods[i.IonForm] = append(protMods[i.IonForm], j)
		}
	}

	for _, i := range pro {

		var rep ProteinEvidence

		rep.SupportingSpectra = make(map[string]int)
		rep.TotalPeptideIons = make(map[string]IonEvidence)
		rep.IndiProtein = make(map[string]uint8)
		rep.Modifications.Index = make(map[string]mod.Modification)

		rep.ProteinName = i.ProteinName
		rep.ProteinGroup = i.GroupNumber
		rep.ProteinSubGroup = i.GroupSiblingID
		rep.Length, _ = strconv.Atoi(i.Length)
		rep.Coverage = i.PercentCoverage
		rep.UniqueStrippedPeptides = len(i.UniqueStrippedPeptides)
		rep.Probability = i.Probability
		rep.TopPepProb = i.TopPepProb

		rep.TotalPeptides = make(map[string]int)
		rep.UniquePeptides = make(map[string]int)
		rep.URazorPeptides = make(map[string]int)

		if strings.HasPrefix(i.ProteinName, decoyTag) {
			rep.IsDecoy = true
		} else {
			rep.IsDecoy = false
		}

		for j := range i.IndistinguishableProtein {
			rep.IndiProtein[i.IndistinguishableProtein[j]] = 0
		}

		for _, k := range i.PeptideIons {

			ion := fmt.Sprintf("%s#%d#%.4f", k.PeptideSequence, k.Charge, k.CalcNeutralPepMass)

			v, ok := evidenceIons[ion]
			if ok {

				for spec := range v.Spectra {
					rep.SupportingSpectra[spec]++
				}

				ref := v
				ref.Weight = k.Weight
				ref.GroupWeight = k.GroupWeight

				for _, l := range k.PeptideParentProtein {
					ref.MappedProteins[l] = 0
				}
				delete(ref.MappedProteins, i.ProteinName)

				ref.Modifications = k.Modifications

				if len(ref.MappedProteins) == 0 && ref.Weight >= weight {
					ref.IsUnique = true
				} else {
					ref.IsUnique = false
				}

				if k.Razor == 1 {
					ref.IsURazor = true
				}

				mods, ok := protMods[ion]
				if ok {
					for _, j := range mods {
						_, okMod := ref.Modifications.Index[j.Index]
						if !okMod && k.IsUnique {
							ref.Modifications.Index[j.Index] = j
							rep.Modifications.Index[j.Index] = j
						}

						if !okMod && k.Razor == 1 {
							ref.Modifications.Index[j.Index] = j
							rep.Modifications.Index[j.Index] = j
						}
					}
				}

				rep.TotalPeptideIons[ion] = ref

			} else {

				var ref IonEvidence
				ref.MappedProteins = make(map[string]int)
				ref.Spectra = make(map[string]int)

				ref.Protein = i.ProteinName

				ref.Sequence = k.PeptideSequence
				ref.IonForm = ion
				ref.ModifiedSequence = k.ModifiedPeptide
				ref.ChargeState = k.Charge
				ref.Probability = k.InitialProbability
				ref.Weight = k.Weight
				ref.GroupWeight = k.GroupWeight
				ref.NumberOfEnzymaticTermini = k.NumberOfEnzymaticTermini
				ref.Labels = k.Labels

				ref.MappedProteins = make(map[string]int)
				for _, l := range k.PeptideParentProtein {
					ref.MappedProteins[l] = 0
				}
				delete(ref.MappedProteins, i.ProteinName)

				ref.Modifications = k.Modifications

				if len(ref.MappedProteins) == 0 && ref.Weight >= weight {
					ref.IsUnique = true
				} else {
					ref.IsUnique = false
				}

				mods, ok := protMods[ion]
				if ok {
					for _, j := range mods {
						_, okMod := ref.Modifications.Index[j.Index]
						if !okMod && k.IsUnique {
							ref.Modifications.Index[j.Index] = j
							rep.Modifications.Index[j.Index] = j
						}

						if !okMod && k.Razor == 1 {
							ref.Modifications.Index[j.Index] = j
							rep.Modifications.Index[j.Index] = j
						}
					}
				}

				rep.TotalPeptideIons[ion] = ref
			}

		}

		// if strings.Contains(rep.ProteinName, "Q8WXG9") {
		// 	spew.Dump(rep)
		// }

		list = append(list, rep)
	}

	var dtb dat.Base
	dtb.Restore()

	if len(dtb.Records) < 1 {
		msg.DatabaseNotFound(errors.New(""), "fatal")
	}

	// fix the name sand headers and pull database information into protein report
	for i := range list {
		for _, j := range dtb.Records {
			//if strings.Contains(j.OriginalHeader, list[i].ProteinName) {
			if strings.HasPrefix(j.OriginalHeader, list[i].ProteinName) {

				//fmt.Println("A:", j.OriginalHeader, "\t", "B:", list[i].ProteinName)

				if (j.IsDecoy && list[i].IsDecoy) || (!j.IsDecoy && !list[i].IsDecoy) {

					list[i].OriginalHeader = j.OriginalHeader
					list[i].PartHeader = j.PartHeader
					list[i].ProteinID = j.ID
					list[i].EntryName = j.EntryName
					list[i].ProteinExistence = j.ProteinExistence
					list[i].GeneNames = j.GeneNames
					list[i].Sequence = j.Sequence
					list[i].ProteinName = j.ProteinName
					list[i].Organism = j.Organism

					// uniprot entries have the description on ProteinName
					if len(j.Description) < 1 {
						list[i].Description = j.ProteinName
					} else {
						list[i].Description = j.Description
					}

					// updating the protein ions
					for _, k := range list[i].TotalPeptideIons {
						k.Protein = j.PartHeader
						k.ProteinID = j.ID
						k.GeneName = j.GeneNames
					}

					break
				}
			}
		}
	}

	sort.Sort(list)
	evi.Proteins = list

}

// MetaProteinReport creates the TSV Protein report
func (evi Evidence) MetaProteinReport(workspace, brand string, channels int, hasDecoys, hasRazor, uniqueOnly, hasLabels bool) {

	var header string
	output := fmt.Sprintf("%s%sprotein.tsv", workspace, string(filepath.Separator))

	// create result file
	file, e := os.Create(output)
	if e != nil {
		msg.WriteFile(errors.New("cannot create protein report"), "error")
	}
	defer file.Close()

	// building the printing set tat may or not contain decoys
	var printSet ProteinEvidenceList
	for _, i := range evi.Proteins {
		if !hasDecoys {
			if !i.IsDecoy {
				printSet = append(printSet, i)
			}
		} else {
			printSet = append(printSet, i)
		}
	}

	header = "Group\tSubGroup\tProtein\tProtein ID\tEntry Name\tGene\tLength\tPercent Coverage\tOrganism\tProtein Description\tProtein Existence\tProtein Probability\tTop Peptide Probability\tTotal Peptides\tUnique Peptides\tRazor Peptides\tTotal Spectral Count\tUnique Spectral Count\tRazor Spectral Count\tTotal Intensity\tUnique Intensity\tRazor Intensity\tRazor Assigned Modifications\tRazor Observed Modifications\tIndistinguishable Proteins"

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
			if len(i.UniqueLabels.Channel1.CustomName) >= 1 {
				c1 = i.UniqueLabels.Channel1.CustomName
				c2 = i.UniqueLabels.Channel2.CustomName
				c3 = i.UniqueLabels.Channel3.CustomName
				c4 = i.UniqueLabels.Channel4.CustomName
				c5 = i.UniqueLabels.Channel5.CustomName
				c6 = i.UniqueLabels.Channel6.CustomName
				c7 = i.UniqueLabels.Channel7.CustomName
				c8 = i.UniqueLabels.Channel8.CustomName
				c9 = i.UniqueLabels.Channel9.CustomName
				c10 = i.UniqueLabels.Channel10.CustomName
				c11 = i.UniqueLabels.Channel11.CustomName
				c12 = i.UniqueLabels.Channel12.CustomName
				c13 = i.UniqueLabels.Channel13.CustomName
				c14 = i.UniqueLabels.Channel14.CustomName
				c15 = i.UniqueLabels.Channel15.CustomName
				c16 = i.UniqueLabels.Channel16.CustomName
				c17 = i.UniqueLabels.Channel17.CustomName
				c18 = i.UniqueLabels.Channel18.CustomName
				break
			}
		}

		header = strings.Replace(header, "Channel "+printSet[10].UniqueLabels.Channel1.Name, c1, -1)
		header = strings.Replace(header, "Channel "+printSet[10].UniqueLabels.Channel2.Name, c2, -1)
		header = strings.Replace(header, "Channel "+printSet[10].UniqueLabels.Channel3.Name, c3, -1)
		header = strings.Replace(header, "Channel "+printSet[10].UniqueLabels.Channel4.Name, c4, -1)
		header = strings.Replace(header, "Channel "+printSet[10].UniqueLabels.Channel5.Name, c5, -1)
		header = strings.Replace(header, "Channel "+printSet[10].UniqueLabels.Channel6.Name, c6, -1)
		header = strings.Replace(header, "Channel "+printSet[10].UniqueLabels.Channel7.Name, c7, -1)
		header = strings.Replace(header, "Channel "+printSet[10].UniqueLabels.Channel8.Name, c8, -1)
		header = strings.Replace(header, "Channel "+printSet[10].UniqueLabels.Channel9.Name, c9, -1)
		header = strings.Replace(header, "Channel "+printSet[10].UniqueLabels.Channel10.Name, c10, -1)
		header = strings.Replace(header, "Channel "+printSet[10].UniqueLabels.Channel11.Name, c11, -1)
		header = strings.Replace(header, "Channel "+printSet[10].UniqueLabels.Channel12.Name, c12, -1)
		header = strings.Replace(header, "Channel "+printSet[10].UniqueLabels.Channel13.Name, c13, -1)
		header = strings.Replace(header, "Channel "+printSet[10].UniqueLabels.Channel14.Name, c14, -1)
		header = strings.Replace(header, "Channel "+printSet[10].UniqueLabels.Channel15.Name, c15, -1)
		header = strings.Replace(header, "Channel "+printSet[10].UniqueLabels.Channel16.Name, c16, -1)
		header = strings.Replace(header, "Channel "+printSet[10].UniqueLabels.Channel17.Name, c17, -1)
		header = strings.Replace(header, "Channel "+printSet[10].UniqueLabels.Channel18.Name, c18, -1)
	}

	_, e = io.WriteString(file, header)
	if e != nil {
		msg.WriteToFile(e, "fatal")
	}

	for _, i := range printSet {

		var ip []string
		for k := range i.IndiProtein {
			ip = append(ip, k)
		}

		assL, obs := getModsList(i.Modifications.Index)

		// var uniqIons int
		// for _, j := range i.TotalPeptideIons {
		// 	if j.IsUnique {
		// 		uniqIons++
		// 	}
		// }

		// var urazorIons int
		// for _, j := range i.TotalPeptideIons {
		// 	if j.IsURazor {
		// 		urazorIons++
		// 	}
		// }

		sort.Strings(assL)
		sort.Strings(obs)
		sort.Strings(ip)

		// change between Unique+Razor and Unique only based on parameter defined on labelquant
		var reportIntensities [18]float64
		if uniqueOnly || !hasRazor {
			reportIntensities[0] = i.UniqueLabels.Channel1.Intensity
			reportIntensities[1] = i.UniqueLabels.Channel2.Intensity
			reportIntensities[2] = i.UniqueLabels.Channel3.Intensity
			reportIntensities[3] = i.UniqueLabels.Channel4.Intensity
			reportIntensities[4] = i.UniqueLabels.Channel5.Intensity
			reportIntensities[5] = i.UniqueLabels.Channel6.Intensity
			reportIntensities[6] = i.UniqueLabels.Channel7.Intensity
			reportIntensities[7] = i.UniqueLabels.Channel8.Intensity
			reportIntensities[8] = i.UniqueLabels.Channel9.Intensity
			reportIntensities[9] = i.UniqueLabels.Channel10.Intensity
			reportIntensities[10] = i.UniqueLabels.Channel11.Intensity
			reportIntensities[11] = i.UniqueLabels.Channel12.Intensity
			reportIntensities[12] = i.UniqueLabels.Channel13.Intensity
			reportIntensities[13] = i.UniqueLabels.Channel14.Intensity
			reportIntensities[14] = i.UniqueLabels.Channel15.Intensity
			reportIntensities[15] = i.UniqueLabels.Channel16.Intensity
			reportIntensities[16] = i.UniqueLabels.Channel17.Intensity
			reportIntensities[17] = i.UniqueLabels.Channel18.Intensity

		} else {
			reportIntensities[0] = i.URazorLabels.Channel1.Intensity
			reportIntensities[1] = i.URazorLabels.Channel2.Intensity
			reportIntensities[2] = i.URazorLabels.Channel3.Intensity
			reportIntensities[3] = i.URazorLabels.Channel4.Intensity
			reportIntensities[4] = i.URazorLabels.Channel5.Intensity
			reportIntensities[5] = i.URazorLabels.Channel6.Intensity
			reportIntensities[6] = i.URazorLabels.Channel7.Intensity
			reportIntensities[7] = i.URazorLabels.Channel8.Intensity
			reportIntensities[8] = i.URazorLabels.Channel9.Intensity
			reportIntensities[9] = i.URazorLabels.Channel10.Intensity
			reportIntensities[10] = i.URazorLabels.Channel11.Intensity
			reportIntensities[11] = i.URazorLabels.Channel12.Intensity
			reportIntensities[12] = i.URazorLabels.Channel13.Intensity
			reportIntensities[13] = i.URazorLabels.Channel14.Intensity
			reportIntensities[14] = i.URazorLabels.Channel15.Intensity
			reportIntensities[15] = i.URazorLabels.Channel16.Intensity
			reportIntensities[16] = i.URazorLabels.Channel17.Intensity
			reportIntensities[17] = i.URazorLabels.Channel18.Intensity
		}

		// proteins with almost no evidences, and completely shared with decoys are eliminated from the analysis,
		// in most cases proteins with one small peptide shared with a decoy
		line := fmt.Sprintf("%d\t%s\t%s\t%s\t%s\t%s\t%d\t%.2f\t%s\t%s\t%s\t%.4f\t%.4f\t%d\t%d\t%d\t%d\t%d\t%d\t%6.f\t%6.f\t%6.f\t%s\t%s\t%s",
			i.ProteinGroup,     // Group
			i.ProteinSubGroup,  // SubGroup
			i.PartHeader,       // Protein
			i.ProteinID,        // Protein ID
			i.EntryName,        // Entry Name
			i.GeneNames,        // Genes
			i.Length,           // Length
			i.Coverage,         // Percent Coverage
			i.Organism,         // Organism
			i.Description,      // Description
			i.ProteinExistence, // Protein Existence
			i.Probability,      // Protein Probability
			i.TopPepProb,       // Top Peptide Probability
			//len(i.TotalPeptideIons),  // Total Peptide Ions
			//uniqIons,                 // Unique Peptide Ions
			//urazorIons,               // Razor Peptide Ions
			len(i.TotalPeptides),     // Total Peptides
			len(i.UniquePeptides),    // Unique Peptides
			len(i.URazorPeptides),    // Razor Peptides
			i.TotalSpC,               // Total Spectral Count
			i.UniqueSpC,              // Unique Spectral Count
			i.URazorSpC,              // Razor Spectral Count
			i.TotalIntensity,         // Total Intensity
			i.UniqueIntensity,        // Unique Intensity
			i.URazorIntensity,        // Razor Intensity
			strings.Join(assL, ", "), // Razor Assigned Modifications
			strings.Join(obs, ", "),  // Razor Observed Modifications
			strings.Join(ip, ", "),   // Indistinguishable Proteins
		)

		switch channels {
		case 4:
			line = fmt.Sprintf("%s\t%.4f\t%.4f\t%.4f\t%.4f",
				line,
				reportIntensities[0],
				reportIntensities[1],
				reportIntensities[2],
				reportIntensities[3],
			)
		case 6:
			line = fmt.Sprintf("%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
				line,
				reportIntensities[0],
				reportIntensities[1],
				reportIntensities[2],
				reportIntensities[3],
				reportIntensities[4],
				reportIntensities[5],
			)
		case 8:
			line = fmt.Sprintf("%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
				line,
				reportIntensities[0],
				reportIntensities[1],
				reportIntensities[2],
				reportIntensities[3],
				reportIntensities[4],
				reportIntensities[5],
				reportIntensities[6],
				reportIntensities[7],
			)
		case 10:
			line = fmt.Sprintf("%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
				line,
				reportIntensities[0],
				reportIntensities[1],
				reportIntensities[2],
				reportIntensities[3],
				reportIntensities[4],
				reportIntensities[5],
				reportIntensities[6],
				reportIntensities[7],
				reportIntensities[8],
				reportIntensities[9],
			)
		case 11:
			line = fmt.Sprintf("%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
				line,
				reportIntensities[0],
				reportIntensities[1],
				reportIntensities[2],
				reportIntensities[3],
				reportIntensities[4],
				reportIntensities[5],
				reportIntensities[6],
				reportIntensities[7],
				reportIntensities[8],
				reportIntensities[9],
				reportIntensities[10],
			)
		case 16:
			line = fmt.Sprintf("%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
				line,
				reportIntensities[0],
				reportIntensities[1],
				reportIntensities[2],
				reportIntensities[3],
				reportIntensities[4],
				reportIntensities[5],
				reportIntensities[6],
				reportIntensities[7],
				reportIntensities[8],
				reportIntensities[9],
				reportIntensities[10],
				reportIntensities[11],
				reportIntensities[12],
				reportIntensities[13],
				reportIntensities[14],
				reportIntensities[15],
			)
		case 18:
			line = fmt.Sprintf("%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
				line,
				reportIntensities[0],
				reportIntensities[1],
				reportIntensities[2],
				reportIntensities[3],
				reportIntensities[4],
				reportIntensities[5],
				reportIntensities[6],
				reportIntensities[7],
				reportIntensities[8],
				reportIntensities[9],
				reportIntensities[10],
				reportIntensities[11],
				reportIntensities[12],
				reportIntensities[13],
				reportIntensities[14],
				reportIntensities[15],
				reportIntensities[16],
				reportIntensities[17],
			)
		default:
			header += ""
		}

		line += "\n"

		_, e = io.WriteString(file, line)
		if e != nil {
			msg.WriteToFile(e, "fatal")
		}

	}
}

// ProteinFastaReport saves to disk a filtered FASTA file with FDR aproved proteins
func (evi *Evidence) ProteinFastaReport(workspace string, hasDecoys bool) {

	output := fmt.Sprintf("%s%sprotein.fas", workspace, string(filepath.Separator))

	file, e := os.Create(output)
	if e != nil {
		msg.WriteFile(e, "fatal")
	}
	defer file.Close()

	// building the printing set tat may or not contain decoys
	var printSet ProteinEvidenceList
	for _, i := range evi.Proteins {
		if !hasDecoys {
			if !i.IsDecoy {
				printSet = append(printSet, i)
			}
		} else {
			printSet = append(printSet, i)
		}
	}

	for _, i := range printSet {
		header := i.OriginalHeader
		line := ">" + header + "\n" + i.Sequence + "\n"
		_, e = io.WriteString(file, line)
		if e != nil {
			msg.WriteToFile(e, "fatal")
		}
	}
}
