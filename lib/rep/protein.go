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

	"philosopher/lib/dat"
	"philosopher/lib/id"
	"philosopher/lib/mod"
	"philosopher/lib/msg"
)

// AssembleProteinReport creates the post processed protein strcuture
func (evi *Evidence) AssembleProteinReport(pro id.ProtIDList, weight float64, decoyTag string) {

	var protMods = make(map[id.IonFormType][]mod.Modification)
	var evidenceIons = make(map[id.IonFormType]*IonEvidence)

	for idx, i := range evi.Ions {
		evidenceIons[i.IonForm()] = &evi.Ions[idx]
	}

	for _, i := range evi.PSM {
		for _, j := range i.Modifications.IndexSlice {
			protMods[i.IonForm()] = append(protMods[i.IonForm()], j)
		}
	}
	evi.Proteins = make(ProteinEvidenceList, len(pro))
	for idx, i := range pro {

		rep := &evi.Proteins[idx]

		rep.ProteinName = i.ProteinName
		rep.Description = i.Description
		rep.OriginalHeader = i.OriginalHeader
		rep.SupportingSpectra = make(map[id.SpectrumType]int)
		rep.TotalPeptideIons = make(map[id.IonFormType]IonEvidence)
		rep.IndiProtein = make(map[string]struct{})
		repModificationsIndex := make(map[string]mod.Modification)
		rep.ProteinGroup = i.GroupNumber
		rep.ProteinSubGroup = i.GroupSiblingID
		rep.Length = i.Length
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
			rep.IndiProtein[i.IndistinguishableProtein[j]] = struct{}{}
		}

		for _, k := range i.PeptideIons {

			//ion := fmt.Sprintf("%s#%d#%.4f", k.PeptideSequence, k.Charge, k.CalcNeutralPepMass)
			ion := k.IonForm()

			if v, ok := evidenceIons[ion]; ok {

				for spec := range v.Spectra {
					rep.SupportingSpectra[spec]++
				}

				ref := *v
				ref.Weight = k.Weight
				ref.GroupWeight = k.GroupWeight

				for _, l := range k.PeptideParentProtein {
					ref.MappedProteins[l] = 0
				}
				delete(ref.MappedProteins, i.ProteinName)

				ref.Modifications = k.Modifications.ToSlice()

				if len(ref.MappedProteins) == 0 && ref.Weight >= weight {
					ref.IsUnique = true
				} else {
					ref.IsUnique = false
				}

				if k.Razor == 1 {
					ref.IsURazor = true
				}
				refModifications := ref.Modifications.ToMap()
				if mods, ok := protMods[ion]; ok {
					for _, j := range mods {
						_, okMod := refModifications.Index[j.Index]
						if !okMod && k.IsUnique {
							refModifications.Index[j.Index] = j
							repModificationsIndex[j.Index] = j
						}

						if !okMod && k.Razor == 1 {
							refModifications.Index[j.Index] = j
							repModificationsIndex[j.Index] = j
						}
					}
				}
				ref.Modifications = refModifications.ToSlice()
				rep.TotalPeptideIons[ion] = ref

			} else {

				var ref IonEvidence
				ref.MappedProteins = make(map[string]int)
				ref.Spectra = make(map[id.SpectrumType]int)

				ref.Protein = i.ProteinName

				ref.Sequence = k.PeptideSequence
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

				ref.Modifications = k.Modifications.ToSlice()

				if len(ref.MappedProteins) == 0 && ref.Weight >= weight {
					ref.IsUnique = true
				} else {
					ref.IsUnique = false
				}
				refModifications := ref.Modifications.ToMap()
				if mods, ok := protMods[ion]; ok {
					for _, j := range mods {
						_, okMod := refModifications.Index[j.Index]
						if !okMod && k.IsUnique {
							refModifications.Index[j.Index] = j
							repModificationsIndex[j.Index] = j
						}

						if !okMod && k.Razor == 1 {
							refModifications.Index[j.Index] = j
							repModificationsIndex[j.Index] = j
						}
					}
				}
				ref.Modifications = refModifications.ToSlice()
				rep.TotalPeptideIons[ion] = ref
			}

		}
		if len(repModificationsIndex) != 0 {
			rep.Modifications = mod.Modifications{Index: repModificationsIndex}.ToSlice()
		}
	}

	var dtb dat.Base
	dtb.Restore()

	if len(dtb.Records) < 1 {
		msg.DatabaseNotFound(errors.New(""), "fatal")
	}

	// fix the name sand headers and pull database information into protein report
	for i := range evi.Proteins {

		pe := &evi.Proteins[i]

		for _, j := range dtb.Records {

			desc := strings.Replace(pe.Description, "|", " ", -1)
			ensName := j.PartHeader + " " + desc

			if (j.OriginalHeader == pe.OriginalHeader) || (ensName == pe.OriginalHeader) {

				if (j.IsDecoy && pe.IsDecoy) || (!j.IsDecoy && !pe.IsDecoy) {

					pe.OriginalHeader = j.OriginalHeader
					pe.PartHeader = j.PartHeader
					pe.ProteinID = j.ID
					pe.EntryName = j.EntryName
					pe.ProteinExistence = j.ProteinExistence
					pe.GeneNames = j.GeneNames
					pe.Sequence = j.Sequence
					pe.ProteinName = j.ProteinName
					pe.Organism = j.Organism

					// some simple headers might not have a full partheader, so we force them to be
					// the same as the EntryName
					if len(pe.PartHeader) == 0 {
						pe.PartHeader = pe.ProteinName
						pe.EntryName = pe.ProteinName
					}

					// uniprot entries have the description on ProteinName
					if len(j.Description) < 1 {
						pe.Description = j.ProteinName
					} else {
						pe.Description = j.Description
					}

					// for Ensemble entries without name
					if len(pe.ProteinName) == 0 {
						pe.ProteinName = pe.PartHeader
					}

					if pe.Length == 0 {
						pe.Length = len(j.Sequence)
					}

					// updating the protein ions
					for _, k := range pe.TotalPeptideIons {
						k.Protein = j.PartHeader
						k.ProteinID = j.ID
						k.GeneName = j.GeneNames
					}

					break
				}
			}
		}
	}

	sort.Sort(evi.Proteins)

}

// ProteinReport creates the TSV Protein report
func (eviProteins ProteinEvidenceList) ProteinReport(workspace, brand, decoyTag string, channels int, hasDecoys, hasRazor, uniqueOnly, hasLabels, hasPrefix, removeContam bool) {

	var header string
	var output string

	if hasPrefix {
		output = fmt.Sprintf("%s%s%s_protein.tsv", workspace, string(filepath.Separator), path.Base(workspace))
	} else {
		output = fmt.Sprintf("%s%sprotein.tsv", workspace, string(filepath.Separator))
	}

	// create result file
	file, e := os.Create(output)
	bw := bufio.NewWriter(file)
	if e != nil {
		msg.WriteFile(errors.New("cannot create protein report"), "error")
	}
	defer file.Close()
	defer bw.Flush()
	// building the printing set tat may or not contain decoys
	var printSet []*ProteinEvidence
	for idx, i := range eviProteins {

		if removeContam && (strings.HasPrefix(i.OriginalHeader, "contam_") || strings.HasPrefix(i.OriginalHeader, "Cont_")) {
			continue
		}

		if !hasDecoys {
			if !i.IsDecoy {
				printSet = append(printSet, &eviProteins[idx])
			}
		} else {
			printSet = append(printSet, &eviProteins[idx])
		}
	}

	header = "Protein\tProtein ID\tEntry Name\tGene\tLength\tOrganism\tProtein Description\tProtein Existence\tCoverage\tProtein Probability\tTop Peptide Probability\tTotal Peptides\tUnique Peptides\tRazor Peptides\tTotal Spectral Count\tUnique Spectral Count\tRazor Spectral Count\tTotal Intensity\tUnique Intensity\tRazor Intensity\tRazor Assigned Modifications\tRazor Observed Modifications\tIndistinguishable Proteins"

	var headerIndex int
	for i := range printSet {
		if printSet[i].UniqueLabels != nil && len(printSet[i].UniqueLabels.Channel1.CustomName) > 0 {
			headerIndex = i
			break
		}
	}

	if brand == "tmt" {
		switch channels {
		case 6:
			header = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s",
				header,
				printSet[headerIndex].URazorLabels.Channel1.CustomName,
				printSet[headerIndex].URazorLabels.Channel2.CustomName,
				printSet[headerIndex].URazorLabels.Channel5.CustomName,
				printSet[headerIndex].URazorLabels.Channel6.CustomName,
				printSet[headerIndex].URazorLabels.Channel9.CustomName,
				printSet[headerIndex].URazorLabels.Channel10.CustomName,
			)
		case 10:
			header = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
				header,
				printSet[headerIndex].URazorLabels.Channel1.CustomName,
				printSet[headerIndex].URazorLabels.Channel2.CustomName,
				printSet[headerIndex].URazorLabels.Channel3.CustomName,
				printSet[headerIndex].URazorLabels.Channel4.CustomName,
				printSet[headerIndex].URazorLabels.Channel5.CustomName,
				printSet[headerIndex].URazorLabels.Channel6.CustomName,
				printSet[headerIndex].URazorLabels.Channel7.CustomName,
				printSet[headerIndex].URazorLabels.Channel8.CustomName,
				printSet[headerIndex].URazorLabels.Channel9.CustomName,
				printSet[headerIndex].URazorLabels.Channel10.CustomName,
			)
		case 11:
			header = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
				header,
				printSet[headerIndex].URazorLabels.Channel1.CustomName,
				printSet[headerIndex].URazorLabels.Channel2.CustomName,
				printSet[headerIndex].URazorLabels.Channel3.CustomName,
				printSet[headerIndex].URazorLabels.Channel4.CustomName,
				printSet[headerIndex].URazorLabels.Channel5.CustomName,
				printSet[headerIndex].URazorLabels.Channel6.CustomName,
				printSet[headerIndex].URazorLabels.Channel7.CustomName,
				printSet[headerIndex].URazorLabels.Channel8.CustomName,
				printSet[headerIndex].URazorLabels.Channel9.CustomName,
				printSet[headerIndex].URazorLabels.Channel10.CustomName,
				printSet[headerIndex].URazorLabels.Channel11.CustomName,
			)
		case 16:
			header = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
				header,
				printSet[headerIndex].URazorLabels.Channel1.CustomName,
				printSet[headerIndex].URazorLabels.Channel2.CustomName,
				printSet[headerIndex].URazorLabels.Channel3.CustomName,
				printSet[headerIndex].URazorLabels.Channel4.CustomName,
				printSet[headerIndex].URazorLabels.Channel5.CustomName,
				printSet[headerIndex].URazorLabels.Channel6.CustomName,
				printSet[headerIndex].URazorLabels.Channel7.CustomName,
				printSet[headerIndex].URazorLabels.Channel8.CustomName,
				printSet[headerIndex].URazorLabels.Channel9.CustomName,
				printSet[headerIndex].URazorLabels.Channel10.CustomName,
				printSet[headerIndex].URazorLabels.Channel11.CustomName,
				printSet[headerIndex].URazorLabels.Channel12.CustomName,
				printSet[headerIndex].URazorLabels.Channel13.CustomName,
				printSet[headerIndex].URazorLabels.Channel14.CustomName,
				printSet[headerIndex].URazorLabels.Channel15.CustomName,
				printSet[headerIndex].URazorLabels.Channel16.CustomName,
			)
		case 18:
			header = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
				header,
				printSet[headerIndex].URazorLabels.Channel1.CustomName,
				printSet[headerIndex].URazorLabels.Channel2.CustomName,
				printSet[headerIndex].URazorLabels.Channel3.CustomName,
				printSet[headerIndex].URazorLabels.Channel4.CustomName,
				printSet[headerIndex].URazorLabels.Channel5.CustomName,
				printSet[headerIndex].URazorLabels.Channel6.CustomName,
				printSet[headerIndex].URazorLabels.Channel7.CustomName,
				printSet[headerIndex].URazorLabels.Channel8.CustomName,
				printSet[headerIndex].URazorLabels.Channel9.CustomName,
				printSet[headerIndex].URazorLabels.Channel10.CustomName,
				printSet[headerIndex].URazorLabels.Channel11.CustomName,
				printSet[headerIndex].URazorLabels.Channel12.CustomName,
				printSet[headerIndex].URazorLabels.Channel13.CustomName,
				printSet[headerIndex].URazorLabels.Channel14.CustomName,
				printSet[headerIndex].URazorLabels.Channel15.CustomName,
				printSet[headerIndex].URazorLabels.Channel16.CustomName,
				printSet[headerIndex].URazorLabels.Channel17.CustomName,
				printSet[headerIndex].URazorLabels.Channel18.CustomName,
			)
		default:
			header += ""
		}
	} else if brand == "itraq" {
		switch channels {
		case 4:
			header = fmt.Sprintf("%s\t%s\t%s\t%s\t%s",
				header,
				printSet[headerIndex].URazorLabels.Channel1.CustomName,
				printSet[headerIndex].URazorLabels.Channel2.CustomName,
				printSet[headerIndex].URazorLabels.Channel3.CustomName,
				printSet[headerIndex].URazorLabels.Channel4.CustomName,
			)
		case 8:
			header = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
				header,
				printSet[headerIndex].URazorLabels.Channel1.CustomName,
				printSet[headerIndex].URazorLabels.Channel2.CustomName,
				printSet[headerIndex].URazorLabels.Channel3.CustomName,
				printSet[headerIndex].URazorLabels.Channel4.CustomName,
				printSet[headerIndex].URazorLabels.Channel5.CustomName,
				printSet[headerIndex].URazorLabels.Channel6.CustomName,
				printSet[headerIndex].URazorLabels.Channel7.CustomName,
				printSet[headerIndex].URazorLabels.Channel8.CustomName,
			)
		default:
			header += ""
		}
	} else if brand == "xtag" {
		header = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
			header,
			printSet[headerIndex].URazorLabels.Channel1.CustomName,
			printSet[headerIndex].URazorLabels.Channel2.CustomName,
			printSet[headerIndex].URazorLabels.Channel3.CustomName,
			printSet[headerIndex].URazorLabels.Channel4.CustomName,
			printSet[headerIndex].URazorLabels.Channel5.CustomName,
			printSet[headerIndex].URazorLabels.Channel6.CustomName,
			printSet[headerIndex].URazorLabels.Channel7.CustomName,
			printSet[headerIndex].URazorLabels.Channel8.CustomName,
			printSet[headerIndex].URazorLabels.Channel9.CustomName,
			printSet[headerIndex].URazorLabels.Channel10.CustomName,
			printSet[headerIndex].URazorLabels.Channel11.CustomName,
			printSet[headerIndex].URazorLabels.Channel12.CustomName,
			printSet[headerIndex].URazorLabels.Channel13.CustomName,
			printSet[headerIndex].URazorLabels.Channel14.CustomName,
			printSet[headerIndex].URazorLabels.Channel15.CustomName,
			printSet[headerIndex].URazorLabels.Channel16.CustomName,
			printSet[headerIndex].URazorLabels.Channel17.CustomName,
			printSet[headerIndex].URazorLabels.Channel18.CustomName,
		)
	}

	header += "\n"

	_, e = io.WriteString(bw, header)
	if e != nil {
		msg.WriteToFile(e, "fatal")
	}

	for _, i := range printSet {

		var ip []string
		for k := range i.IndiProtein {
			ip = append(ip, k)
		}

		assL, obs := getModsList(i.Modifications.ToMap().Index)

		sort.Strings(assL)
		sort.Strings(obs)
		sort.Strings(ip)

		// change between Unique+Razor and Unique only based on parameter defined on labelquant
		var reportIntensities [18]float64
		if uniqueOnly || !hasRazor {
			if i.UniqueLabels != nil {
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
			}
		} else {
			if i.URazorLabels != nil {
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
		}

		// append decoy tags on the gene and proteinID names
		if i.IsDecoy {
			i.ProteinID = decoyTag + i.ProteinID
			i.GeneNames = decoyTag + i.GeneNames
			i.EntryName = decoyTag + i.EntryName
		}

		// proteins with almost no evidences, and completely shared with decoys are eliminated from the an	alysis,
		// in most cases proteins with one small peptide shared with a decoy
		line := fmt.Sprintf("%s\t%s\t%s\t%s\t%d\t%s\t%s\t%s\t%.2f\t%.4f\t%.4f\t%d\t%d\t%d\t%d\t%d\t%d\t%6.f\t%6.f\t%6.f\t%s\t%s\t%s",
			i.PartHeader,             // Protein
			i.ProteinID,              // Protein ID
			i.EntryName,              // Entry Name
			i.GeneNames,              // Genes
			i.Length,                 // Length
			i.Organism,               // Organism
			i.Description,            // Description
			i.ProteinExistence,       // Protein Existence
			i.Coverage,               // Coverage
			i.Probability,            // Protein Probability
			i.TopPepProb,             // Top Peptide Probability
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

		if brand == "tmt" || brand == "itraq" {
			switch channels {
			case 2:
				line = fmt.Sprintf("%s\t%.4f\t%.4f",
					line,
					reportIntensities[0],
					reportIntensities[1],
				)
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
					reportIntensities[4],
					reportIntensities[5],
					reportIntensities[8],
					reportIntensities[9],
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
		} else if brand == "xtag" {
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
		}

		line += "\n"

		_, e = io.WriteString(bw, line)
		if e != nil {
			msg.WriteToFile(e, "fatal")
		}

	}
}

// ProteinFastaReport saves to disk a filtered FASTA file with FDR aproved proteins
func (eviProteins ProteinEvidenceList) ProteinFastaReport(workspace string, hasDecoys bool) {

	output := fmt.Sprintf("%s%sprotein.fas", workspace, string(filepath.Separator))

	file, e := os.Create(output)
	bw := bufio.NewWriter(file)
	if e != nil {
		msg.WriteFile(e, "fatal")
	}
	defer file.Close()
	defer bw.Flush()
	// building the printing set tat may or not contain decoys
	var printSet []*ProteinEvidence
	for idx, i := range eviProteins {
		if !hasDecoys {
			if !i.IsDecoy {
				printSet = append(printSet, &eviProteins[idx])
			}
		} else {
			printSet = append(printSet, &eviProteins[idx])
		}
	}

	for _, i := range printSet {
		header := i.OriginalHeader
		line := ">" + header + "\n" + i.Sequence + "\n"
		_, e = io.WriteString(bw, line)
		if e != nil {
			msg.WriteToFile(e, "fatal")
		}
	}
}
