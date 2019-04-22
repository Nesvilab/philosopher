package rep

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/prvst/philosopher/lib/dat"
	"github.com/prvst/philosopher/lib/id"
	"github.com/prvst/philosopher/lib/mod"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/sirupsen/logrus"
)

// AssembleProteinReport ...
func (e *Evidence) AssembleProteinReport(pro id.ProtIDList, decoyTag string) error {

	var list ProteinEvidenceList
	var protMods = make(map[string][]mod.Modification)
	var err error

	var evidenceIons = make(map[string]IonEvidence)
	for _, i := range e.Ions {
		evidenceIons[i.IonForm] = i

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
		rep.Length = i.Length
		rep.Coverage = i.PercentCoverage
		rep.UniqueStrippedPeptides = len(i.UniqueStrippedPeptides)
		rep.Probability = i.Probability
		rep.TopPepProb = i.TopPepProb

		if strings.Contains(i.ProteinName, decoyTag) {
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

				v.MappedProteins = make(map[string]int)

				ref := v
				ref.Weight = k.Weight
				ref.GroupWeight = k.GroupWeight

				ref.MappedProteins[i.ProteinName]++
				ref.MappedProteins = make(map[string]int)
				for _, l := range k.PeptideParentProtein {
					ref.MappedProteins[l] = 0
				}

				ref.Modifications = k.Modifications

				ref.IsUnique = k.IsUnique
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

				ref.MappedProteins[i.ProteinName]++
				ref.MappedProteins = make(map[string]int)
				for _, l := range k.PeptideParentProtein {
					ref.MappedProteins[l] = 0
				}

				ref.Modifications = k.Modifications

				ref.IsUnique = k.IsUnique
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
			}

		}

		list = append(list, rep)
	}

	var dtb dat.Base
	dtb.Restore()

	if len(dtb.Records) < 1 {
		return errors.New("Cant locate database data")
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
	e.Proteins = list

	return err
}

// ProteinReport ...
func (e *Evidence) ProteinReport(hasDecoys bool) {

	// create result file
	output := fmt.Sprintf("%s%sprotein.tsv", sys.MetaDir(), string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create protein report:", err)
	}
	defer file.Close()

	line := fmt.Sprintf("Group\tSubGroup\tProtein\tProtein ID\tEntry Name\tGene\tLength\tPercent Coverage\tOrganism\tProtein Description\tProtein Existence\tProtein Probability\tTop Peptide Probability\tStripped Peptides\tTotal Peptide Ions\tUnique Peptide Ions\tRazor Peptide Ions\tTotal Spectral Count\tUnique Spectral Count\tRazor Spectral Count\tTotal Intensity\tUnique Intensity\tRazor Intensity\tRazor Assigned Modifications\tRazor Observed Modifications\tIndistinguishable Proteins\n")

	n, err := io.WriteString(file, line)
	if err != nil {
		logrus.Fatal(n, err)
	}

	// building the printing set tat may or not contain decoys
	var printSet ProteinEvidenceList
	for _, i := range e.Proteins {
		if hasDecoys == false {
			if i.IsDecoy == false {
				printSet = append(printSet, i)
			}
		} else {
			printSet = append(printSet, i)
		}
	}

	for _, i := range printSet {

		var ip []string
		for k := range i.IndiProtein {
			ip = append(ip, k)
		}

		assL, obs := getModsList(i.Modifications.Index)

		var uniqIons int
		for _, j := range i.TotalPeptideIons {
			//if j.IsNondegenerateEvidence == true {
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

		// proteins with almost no evidences, and completely shared with decoys are eliminated from the analysis,
		// in most cases proteins with one small peptide shared with a decoy
		//if len(i.TotalPeptideIons) > 0 {

		line = fmt.Sprintf("%d\t%s\t%s\t%s\t%s\t%s\t%d\t%.2f\t%s\t%s\t%s\t%.4f\t%.4f\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%6.f\t%6.f\t%6.f\t%s\t%s\t%s\t",
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

		line += "\n"
		n, err := io.WriteString(file, line)
		if err != nil {
			logrus.Fatal(n, err)
		}
		//}

	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}

// ProteinTMTReport ...
func (e *Evidence) ProteinTMTReport(labels map[string]string, uniqueOnly, hasDecoys bool) {

	// create result file
	output := fmt.Sprintf("%s%sprotein.tsv", sys.MetaDir(), string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create report file:", err)
	}
	defer file.Close()

	line := fmt.Sprintf("Group\tSubGroup\tProtein\tProtein ID\tEntry Name\tGene\tLength\tPercent Coverage\tOrganism\tDescription\tProtein Existence\tProtein Probability\tTop Peptide Probability\tUnique Stripped Peptides\tTotal Peptide Ions\tUnique Peptide Ions\tRazor Peptides Ions\tTotal Spectral Count\tUnique Spectral Count\tRazor Spectral Count\tTotal Intensity\tUnique Intensity\tRazor Intensity\t126 Abundance\t127N Abundance\t127C Abundance\t128N Abundance\t128C Abundance\t129N Abundance\t129C Abundance\t130N Abundance\t130C Abundance\t131N Abundance\t131C Abundance\tRazor Assigned Modifications\tRazor Observed Modifications\tIndistinguishable Proteins\n")

	if len(labels) > 0 {
		for k, v := range labels {
			line = strings.Replace(line, k, v, -1)
		}
	}

	n, err := io.WriteString(file, line)
	if err != nil {
		logrus.Fatal(n, err)
	}

	// building the printing set tat may or not contain decoys
	var printSet ProteinEvidenceList
	for _, i := range e.Proteins {
		if hasDecoys == false {
			if i.IsDecoy == false {
				printSet = append(printSet, i)
			}
		} else {
			printSet = append(printSet, i)
		}
	}

	for _, i := range printSet {

		var ip []string
		for k := range i.IndiProtein {
			ip = append(ip, k)
		}

		var assL []string
		var obs []string

		for _, j := range i.Modifications.Index {
			if j.Type == "Assigned" && j.Variable == "Y" {
				assL = append(assL, fmt.Sprintf("%s%s:%s", j.Position, j.AminoAcid, j.Name))
			} else if j.Type == "Observed" {
				obs = append(obs, fmt.Sprintf("%.4f:%s", j.MassDiff, j.Name))
			}
		}

		var uniqIons int
		for _, j := range i.TotalPeptideIons {
			//if j.IsNondegenerateEvidence == true {
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

		// change between Unique+Razor and Unique only based on paramter defined on labelquant
		var reportIntensities [11]float64
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
		}

		if len(i.TotalPeptideIons) > 0 {
			line = fmt.Sprintf("%d\t%s\t%s\t%s\t%s\t%s\t%d\t%.2f\t%s\t%s\t%s\t%.4f\t%.4f\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%6.f\t%6.f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%s\t%s\t%s\n",
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
				i.TopPepProb,             // Top peptide Probability
				i.UniqueStrippedPeptides, // Unique Stripped Peptides
				len(i.TotalPeptideIons),  // Total peptide Ions
				uniqIons,                 // Unique Peptide Ions
				urazorIons,               // Unique+Razor peptide Ions
				i.TotalSpC,               // Total Spectral Count
				i.UniqueSpC,              // Unique Spectral Count
				i.URazorSpC,              // Unique+Razor Spectral Count
				i.TotalIntensity,         // Total Intensity
				i.UniqueIntensity,        // Unique Intensity
				i.URazorIntensity,        // Razor Intensity
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
				strings.Join(assL, ", "), // Razor Assigned Modifications
				strings.Join(obs, ", "),  // Razor Observed Modifications
				strings.Join(ip, ", "),
			) // Indistinguishable Proteins

			//			line += "\n"
			n, err := io.WriteString(file, line)
			if err != nil {
				logrus.Fatal(n, err)
			}
		}
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}

// PhosphoProteinTMTReport ...
func (e *Evidence) PhosphoProteinTMTReport(labels map[string]string, uniqueOnly, hasDecoys bool) {

	// create result file
	output := fmt.Sprintf("%s%sphosphoprotein.tsv", sys.MetaDir(), string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create report file:", err)
	}
	defer file.Close()

	line := fmt.Sprintf("Group\tSubGroup\tProtein\tProtein ID\tEntry Name\tGene\tLength\tPercent Coverage\tOrganism\tDescription\tProtein Existence\tProtein Probability\tTop Peptide Probability\tUnique Stripped Peptides\tTotal Peptide Ions\tUnique Peptide Ions\tRazor Peptides Ions\tTotal Spectral Count\tUnique Spectral Count\tRazor Spectral Count\tTotal Intensity\tUnique Intensity\tRazor Intensity\t126 Abundance\t127N Abundance\t127C Abundance\t128N Abundance\t128C Abundance\t129N Abundance\t129C Abundance\t130N Abundance\t130C Abundance\t131N Abundancet\t131C Abundance\tRazor Assigned Modifications\tRazor Observed Modifications\tIndistinguishableProteins\n")

	if len(labels) > 0 {
		for k, v := range labels {
			line = strings.Replace(line, k, v, -1)
		}
	}

	n, err := io.WriteString(file, line)
	if err != nil {
		logrus.Fatal(n, err)
	}

	// building the printing set tat may or not contain decoys
	var printSet ProteinEvidenceList
	for _, i := range e.Proteins {
		if hasDecoys == false {
			if i.IsDecoy == false {
				printSet = append(printSet, i)
			}
		} else {
			printSet = append(printSet, i)
		}
	}

	for _, i := range printSet {

		var ip []string
		for k := range i.IndiProtein {
			ip = append(ip, k)
		}

		assL, obs := getModsList(i.Modifications.Index)

		var uniqIons int
		for _, j := range i.TotalPeptideIons {
			//if j.IsNondegenerateEvidence == true {
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

		// change between Unique+Razor and Unique only based on paramter defined on labelquant
		var reportIntensities [11]float64
		if uniqueOnly == true {
			reportIntensities[0] = i.PhosphoUniqueLabels.Channel1.Intensity
			reportIntensities[1] = i.PhosphoUniqueLabels.Channel2.Intensity
			reportIntensities[2] = i.PhosphoUniqueLabels.Channel3.Intensity
			reportIntensities[3] = i.PhosphoUniqueLabels.Channel4.Intensity
			reportIntensities[4] = i.PhosphoUniqueLabels.Channel5.Intensity
			reportIntensities[5] = i.PhosphoUniqueLabels.Channel6.Intensity
			reportIntensities[6] = i.PhosphoUniqueLabels.Channel7.Intensity
			reportIntensities[7] = i.PhosphoUniqueLabels.Channel8.Intensity
			reportIntensities[8] = i.PhosphoUniqueLabels.Channel9.Intensity
			reportIntensities[9] = i.PhosphoUniqueLabels.Channel10.Intensity
			reportIntensities[10] = i.PhosphoUniqueLabels.Channel11.Intensity
		} else {
			reportIntensities[0] = i.PhosphoURazorLabels.Channel1.Intensity
			reportIntensities[1] = i.PhosphoURazorLabels.Channel2.Intensity
			reportIntensities[2] = i.PhosphoURazorLabels.Channel3.Intensity
			reportIntensities[3] = i.PhosphoURazorLabels.Channel4.Intensity
			reportIntensities[4] = i.PhosphoURazorLabels.Channel5.Intensity
			reportIntensities[5] = i.PhosphoURazorLabels.Channel6.Intensity
			reportIntensities[6] = i.PhosphoURazorLabels.Channel7.Intensity
			reportIntensities[7] = i.PhosphoURazorLabels.Channel8.Intensity
			reportIntensities[8] = i.PhosphoURazorLabels.Channel9.Intensity
			reportIntensities[9] = i.PhosphoURazorLabels.Channel10.Intensity
			reportIntensities[10] = i.PhosphoURazorLabels.Channel11.Intensity
		}

		if len(i.TotalPeptideIons) > 0 {
			line = fmt.Sprintf("%d\t%s\t%s\t%s\t%s\t%s\t%d\t%.2f\t%s\t%s\t%s\t%.4f\t%.4f\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%6.f\t%6.f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%s\t%s\t%s\n",
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
				i.TopPepProb,             // Top peptide Probability
				i.UniqueStrippedPeptides, // Unique Stripped Peptides
				len(i.TotalPeptideIons),  // Total peptide Ions
				uniqIons,                 // Unique Peptide Ions
				urazorIons,               // Unique+Razor peptide Ions
				i.TotalSpC,               // Total Spectral Count
				i.UniqueSpC,              // Unique Spectral Count
				i.URazorSpC,              // Unique+Razor Spectral Count
				i.TotalIntensity,         // Total Intensity
				i.UniqueIntensity,        // Unique Intensity
				i.URazorIntensity,        // Razor Intensity
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
				strings.Join(assL, ", "), // Razor Assigned Modifications
				strings.Join(obs, ", "),  // Razor Observed Modifications
				strings.Join(ip, ", "),
			)

			//			line += "\n"
			n, err := io.WriteString(file, line)
			if err != nil {
				logrus.Fatal(n, err)
			}
		}
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}

// ProteinFastaReport saves to disk a filtered FASTA file with FDR aproved proteins
func (e *Evidence) ProteinFastaReport(hasDecoys bool) error {

	output := fmt.Sprintf("%s%sproteins.fas", sys.MetaDir(), string(filepath.Separator))

	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Could not create output file")
	}
	defer file.Close()

	// building the printing set tat may or not contain decoys
	var printSet ProteinEvidenceList
	for _, i := range e.Proteins {
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
		_, err = io.WriteString(file, line)
		if err != nil {
			return errors.New("Cannot print PSM to file")
		}
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return nil
}
