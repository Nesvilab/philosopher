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

	"github.com/nesvilab/philosopher/lib/dat"
	"github.com/nesvilab/philosopher/lib/id"
	"github.com/nesvilab/philosopher/lib/mod"
	"github.com/nesvilab/philosopher/lib/msg"
	"github.com/nesvilab/philosopher/lib/sys"
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

				//v.MappedProteins = make(map[string]int)

				ref := v
				ref.Weight = k.Weight
				ref.GroupWeight = k.GroupWeight

				//ref.MappedProteins = make(map[string]int)
				//ref.MappedProteins[i.ProteinName]++
				for _, l := range k.PeptideParentProtein {
					ref.MappedProteins[l] = 0
				}

				_, ok := ref.MappedProteins[i.ProteinName]
				if ok {
					delete(ref.MappedProteins, i.ProteinName)
				}

				ref.Modifications = k.Modifications

				// if k.IsNondegenerateEvidence == true {
				// 	ref.IsUnique = true
				// } else {
				// 	ref.IsUnique = false
				// }

				if len(ref.MappedProteins) == 0 {
					ref.IsUnique = true
				} else {
					ref.IsUnique = false
				}

				if ref.Weight >= weight {
					ref.IsUnique = true
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

				ref.Sequence = k.PeptideSequence
				ref.IonForm = ion
				ref.ModifiedSequence = k.ModifiedPeptide
				ref.ChargeState = k.Charge
				ref.Probability = k.InitialProbability
				ref.Weight = k.Weight
				ref.GroupWeight = k.GroupWeight
				ref.Labels = k.Labels

				ref.MappedProteins = make(map[string]int)
				//ref.MappedProteins[i.ProteinName]++
				for _, l := range k.PeptideParentProtein {
					ref.MappedProteins[l] = 0
				}

				_, ok := ref.MappedProteins[i.ProteinName]
				if ok {
					delete(ref.MappedProteins, i.ProteinName)
				}

				ref.Modifications = k.Modifications

				//ref.IsUnique = k.IsUnique
				// if k.Razor == 1 {
				// 	ref.IsURazor = true
				// }
				if len(ref.MappedProteins) == 0 {
					ref.IsUnique = true
				} else {
					ref.IsUnique = false
				}

				if ref.Weight >= weight {
					ref.IsUnique = true
				}

				// if k.IsNondegenerateEvidence == true {
				// 	ref.IsUnique = true
				// } else {
				// 	ref.IsUnique = false
				// }

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

		list = append(list, rep)
	}

	var dtb dat.Base
	dtb.Restore()

	if len(dtb.Records) < 1 {
		msg.DatabaseNotFound(errors.New(""), "fatal")
	}

	// fix the name sand headers and pull database information into proteinreport
	for i := range list {
		for _, j := range dtb.Records {
			if strings.Contains(j.OriginalHeader, list[i].ProteinName) {
				if (j.IsDecoy == true && list[i].IsDecoy == true) || (j.IsDecoy == false && list[i].IsDecoy == false) {
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

					break
				}
			}
		}
	}

	sort.Sort(list)
	evi.Proteins = list

	return
}

// MetaProteinReport creates the TSV Protein report
func (evi Evidence) MetaProteinReport(labels map[string]string, brand string, channels int, hasDecoys, uniqueOnly bool) {

	var header string
	output := fmt.Sprintf("%s%sprotein.tsv", sys.MetaDir(), string(filepath.Separator))

	// create result file
	file, e := os.Create(output)
	if e != nil {
		msg.WriteFile(errors.New("Cannot create protein report"), "error")
	}
	defer file.Close()

	// building the printing set tat may or not contain decoys
	var printSet ProteinEvidenceList
	for _, i := range evi.Proteins {
		if hasDecoys == false {
			if i.IsDecoy == false {
				printSet = append(printSet, i)
			}
		} else {
			printSet = append(printSet, i)
		}
	}

	header = fmt.Sprintf("Group\tSubGroup\tProtein\tProtein ID\tEntry Name\tGene\tLength\tPercent Coverage\tOrganism\tProtein Description\tProtein Existence\tProtein Probability\tTop Peptide Probability\tStripped Peptides\tTotal Peptide Ions\tUnique Peptide Ions\tRazor Peptide Ions\tTotal Spectral Count\tUnique Spectral Count\tRazor Spectral Count\tTotal Intensity\tUnique Intensity\tRazor Intensity\tRazor Assigned Modifications\tRazor Observed Modifications\tIndistinguishable Proteins")

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
		msg.WriteToFile(e, "fatal")
	}

	for _, i := range printSet {

		var ip []string
		for k := range i.IndiProtein {
			ip = append(ip, k)
		}

		assL, obs := getModsList(i.Modifications.Index)

		var uniqIons int
		for _, j := range i.TotalPeptideIons {
			if j.IsUnique == true {
				uniqIons++
			}
		}

		var urazorIons int
		for _, j := range i.TotalPeptideIons {
			if j.IsURazor == true {
				urazorIons++
			}
		}

		sort.Strings(assL)
		sort.Strings(obs)
		sort.Strings(ip)

		// change between Unique+Razor and Unique only based on paramter defined on labelquant
		var reportIntensities [16]float64
		if uniqueOnly == true {
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

		}

		// proteins with almost no evidences, and completely shared with decoys are eliminated from the analysis,
		// in most cases proteins with one small peptide shared with a decoy
		line := fmt.Sprintf("%d\t%s\t%s\t%s\t%s\t%s\t%d\t%.2f\t%s\t%s\t%s\t%.4f\t%.4f\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%6.f\t%6.f\t%6.f\t%s\t%s\t%s",
			i.ProteinGroup,           // Group
			i.ProteinSubGroup,        // SubGroup
			i.PartHeader,             // Protein
			i.ProteinID,              // Protein ID
			i.EntryName,              // Entry Name
			i.GeneNames,              // Genes
			i.Length,                 // Length
			i.Coverage,               // Percent Coverage
			i.Organism,               // Organism
			i.Description,            // Description
			i.ProteinExistence,       // Protein Existence
			i.Probability,            // Protein Probability
			i.TopPepProb,             // Top Peptide Probability
			i.UniqueStrippedPeptides, // Stripped Peptides
			len(i.TotalPeptideIons),  // Total Peptide Ions
			uniqIons,                 // Unique Peptide Ions
			urazorIons,               // Razor Peptide Ions
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

		if brand == "tmt" {
			switch channels {
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
			default:
				header += ""
			}
		}

		line += "\n"

		_, e = io.WriteString(file, line)
		if e != nil {
			msg.WriteToFile(e, "fatal")
		}

	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}

// ProteinFastaReport saves to disk a filtered FASTA file with FDR aproved proteins
func (evi *Evidence) ProteinFastaReport(hasDecoys bool) {

	output := fmt.Sprintf("%s%sprotein.fas", sys.MetaDir(), string(filepath.Separator))

	file, e := os.Create(output)
	if e != nil {
		msg.WriteFile(e, "fatal")
	}
	defer file.Close()

	// building the printing set tat may or not contain decoys
	var printSet ProteinEvidenceList
	for _, i := range evi.Proteins {
		if hasDecoys == false {
			if i.IsDecoy == false {
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

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}
