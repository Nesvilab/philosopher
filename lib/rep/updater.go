package rep

import (
	"strings"

	"github.com/nesvilab/philosopher/lib/dat"
)

// PeptideMap struct
type PeptideMap struct {
	Sequence  string
	IonForm   string
	Protein   string
	ProteinID string
	Gene      string
	Proteins  map[string]int
}

// UpdateMappedProteins (DEPRECATED) updates the list of mapped proteins on the data structures
func (evi *Evidence) UpdateMappedProteins(decoyTag string) {

	var list = make(map[string]PeptideMap)
	var proteinMap = make(map[string]int8)

	// The PSM exclusion list was implemented on July 19 because e noticed that the psm.tsv
	// and protein tsv had a different number of unique protein IDs. The PSM tables had spectra
	// mapping to decoys and/or other proteins that do not exist in the final protein table. This
	// is most likely an effect of the backtracking with the promotion fo sequences based on the
	// alternative lists. Since these PSMs are mapping to proteins that do not enter the final
	// protein list, we decided to remove them and make both lists compatible in quantity and quality.
	var psmExclusion = make(map[string]uint8)

	for _, i := range evi.Proteins {
		for _, v := range i.TotalPeptideIons {

			_, ok := list[v.Sequence]
			if !ok {
				var pm PeptideMap

				pm.Sequence = v.Sequence
				pm.IonForm = v.IonForm
				pm.Proteins = v.MappedProteins
				pm.Protein = i.PartHeader
				pm.ProteinID = i.ProteinID
				pm.Gene = i.GeneNames
				pm.Proteins[i.PartHeader] = 0

				list[pm.Sequence] = pm
				proteinMap[i.PartHeader] = 0
			}
		}
	}

	// PSMs
	for i := range evi.PSM {
		v, ok := list[evi.PSM[i].Peptide]
		if ok {
			for k := range v.Proteins {
				evi.PSM[i].MappedProteins[k]++
			}
			if !strings.HasPrefix(v.Protein, decoyTag) && !strings.HasPrefix(evi.PSM[i].Protein, decoyTag) {
				evi.PSM[i].Protein = v.Protein
				evi.PSM[i].ProteinID = v.ProteinID
				evi.PSM[i].GeneName = v.Gene
			} else if strings.HasPrefix(v.Protein, decoyTag) && evi.PSM[i].IsDecoy {
				evi.PSM[i].Protein = v.Protein
				evi.PSM[i].ProteinID = v.ProteinID
				evi.PSM[i].GeneName = v.Gene
			}
		}
		_, ok = proteinMap[evi.PSM[i].Protein]
		if !ok {
			psmExclusion[evi.PSM[i].Spectrum] = 0
		}
	}

	var psm PSMEvidenceList
	for _, i := range evi.PSM {
		_, ok := psmExclusion[i.Spectrum]
		if !ok {
			psm = append(psm, i)
		}
	}
	evi.PSM = psm

	// Peptides
	for i := range evi.Peptides {
		v, ok := list[evi.Peptides[i].Sequence]
		if ok {
			for k := range v.Proteins {
				evi.Peptides[i].MappedProteins[k]++
			}
			if !strings.HasPrefix(v.Protein, decoyTag) && !strings.HasPrefix(evi.Peptides[i].Protein, decoyTag) {
				evi.Peptides[i].Protein = v.Protein
				evi.Peptides[i].ProteinID = v.ProteinID
				evi.Peptides[i].GeneName = v.Gene
			} else if strings.HasPrefix(v.Protein, decoyTag) && evi.Peptides[i].IsDecoy {
				evi.Peptides[i].Protein = v.Protein
				evi.Peptides[i].ProteinID = v.ProteinID
				evi.Peptides[i].GeneName = v.Gene
			}
		}
	}

	var pep PeptideEvidenceList
	for _, i := range evi.Peptides {
		pep = append(pep, i)
	}
	evi.Peptides = pep

	// Ions
	for i := range evi.Ions {
		v, ok := list[evi.Ions[i].Sequence]
		if ok {
			for k := range v.Proteins {
				evi.Ions[i].MappedProteins[k]++
			}
			if !strings.HasPrefix(v.Protein, decoyTag) && !strings.HasPrefix(evi.Ions[i].Protein, decoyTag) {
				evi.Ions[i].Protein = v.Protein
				evi.Ions[i].ProteinID = v.ProteinID
				evi.Ions[i].GeneName = v.Gene
			} else if strings.HasPrefix(v.Protein, decoyTag) && evi.Ions[i].IsDecoy {
				evi.Ions[i].Protein = v.Protein
				evi.Ions[i].ProteinID = v.ProteinID
				evi.Ions[i].GeneName = v.Gene
			}
		}
	}

	var ion IonEvidenceList
	for _, i := range evi.Ions {
		ion = append(ion, i)
	}
	evi.Ions = ion

	return
}

// UpdateIonStatus pushes back to ion and psm evideces the uniqueness and razorness status of each peptide and ion
func (evi *Evidence) UpdateIonStatus(decoyTag string) {

	var uniqueMap = make(map[string]bool)
	var urazorMap = make(map[string]string)

	var uniqueSeqMap = make(map[string]string)

	for _, i := range evi.Proteins {

		for _, j := range i.TotalPeptideIons {
			if j.IsUnique == true {
				uniqueMap[j.IonForm] = true
			}
		}

		for _, j := range i.TotalPeptideIons {
			if j.IsURazor == true {
				urazorMap[j.IonForm] = i.PartHeader
			}
		}
	}

	for i := range evi.PSM {

		if len(evi.PSM[i].MappedProteins) == 0 {
			evi.PSM[i].IsUnique = true
		}

		_, uOK := uniqueMap[evi.PSM[i].IonForm]
		if uOK {
			evi.PSM[i].IsUnique = true
		}

		// the decoy tag checking is a failsafe mechanism to avoid proteins
		// with real complex razor case decisions to pass dowsntream
		// wrong classifications. If by any chance the protein gets assigned to
		// a razor decoy, this mchanism avoids the replacement
		rp, rOK := urazorMap[evi.PSM[i].IonForm]
		if rOK {
			evi.PSM[i].IsURazor = true
			if !strings.Contains(rp, decoyTag) {
				evi.PSM[i].MappedProteins[evi.PSM[i].Protein] = 0
				delete(evi.PSM[i].MappedProteins, rp)
				evi.PSM[i].Protein = rp
			}
		}

		uniqueSeqMap[evi.PSM[i].Peptide] = evi.PSM[i].Protein
	}

	for i := range evi.Ions {

		_, uOK := uniqueMap[evi.Ions[i].IonForm]
		if uOK {
			evi.Ions[i].IsUnique = true
		}

		rp, rOK := urazorMap[evi.Ions[i].IonForm]
		if rOK {
			evi.Ions[i].IsURazor = true
			if !strings.Contains(rp, decoyTag) {
				evi.Ions[i].MappedProteins[evi.Ions[i].Protein] = 0
				delete(evi.Ions[i].MappedProteins, rp)
				evi.Ions[i].Protein = rp
			}
		}
	}

	for i := range evi.Peptides {

		v, ok := uniqueSeqMap[evi.Peptides[i].Sequence]
		if ok {
			//if !strings.Contains(rp, decoyTag) {
			evi.Peptides[i].MappedProteins[evi.Peptides[i].Protein] = 0
			delete(evi.Peptides[i].MappedProteins, v)
			evi.Peptides[i].Protein = v
			//}
		}
	}

	return
}

// UpdateIonStatus pushes back to ion and psm evideces the uniqueness and razorness status of each peptide and ion
// func (evi *Evidence) UpdateIonStatus(decoyTag string) {

// 	var uniqueMap = make(map[string]bool)
// 	var urazorMap = make(map[string]string)
// 	//var ptMap = make(map[string]string)

// 	for _, i := range evi.Proteins {

// 		for _, j := range i.TotalPeptideIons {
// 			if j.IsUnique == true {
// 				uniqueMap[j.IonForm] = true
// 			}
// 		}

// 		for _, j := range i.TotalPeptideIons {
// 			if j.IsURazor == true {
// 				urazorMap[j.IonForm] = i.PartHeader
// 			}
// 		}
// 	}

// 	for i := range evi.PSM {

// 		if len(evi.PSM[i].MappedProteins) == 0 {
// 			evi.PSM[i].IsUnique = true
// 		}

// 		_, uOK := uniqueMap[evi.PSM[i].IonForm]
// 		if uOK {
// 			evi.PSM[i].IsUnique = true
// 		}

// 		// the decoy tag checking is a failsafe mechanism to avoid proteins
// 		// with real complex razor case decisions to pass dowsntream
// 		// wrong classifications. If by any chance the protein gets assigned to
// 		// a razor decoy, this mchanism avoids the replacement
// 		rp, rOK := urazorMap[evi.PSM[i].IonForm]
// 		if rOK {
// 			evi.PSM[i].IsURazor = true
// 			if !strings.Contains(rp, decoyTag) {
// 				evi.PSM[i].MappedProteins[evi.PSM[i].Protein] = 0
// 				delete(evi.PSM[i].MappedProteins, rp)
// 				evi.PSM[i].Protein = rp
// 			}
// 		}

// 		// v, ok := ptMap[evi.PSM[i].IonForm]
// 		// if ok {
// 		// 	evi.PSM[i].Protein = v
// 		// }
// 	}

// 	for i := range evi.Ions {

// 		_, uOK := uniqueMap[evi.Ions[i].IonForm]
// 		if uOK {
// 			evi.Ions[i].IsUnique = true
// 		}

// 		// _, rOK := urazorMap[evi.Ions[i].IonForm]
// 		// if rOK {
// 		// 	evi.Ions[i].IsURazor = true
// 		// }

// 		// a razor decoy, this mchanism avoids the replacement
// 		rp, rOK := urazorMap[evi.Ions[i].IonForm]
// 		if rOK {
// 			evi.Ions[i].IsURazor = true
// 			if !strings.Contains(rp, decoyTag) {
// 				evi.Ions[i].MappedProteins[evi.Ions[i].Protein] = 0
// 				delete(evi.Ions[i].MappedProteins, rp)
// 				evi.Ions[i].Protein = rp
// 			}
// 		}
// 	}

// 	return
// }

// UpdateIonModCount counts how many times each ion is observed modified and not modified
func (evi *Evidence) UpdateIonModCount() {

	// recreate the ion list from the main report object
	var AllIons = make(map[string]int)
	var ModIons = make(map[string]int)
	var UnModIons = make(map[string]int)

	for _, i := range evi.Ions {
		AllIons[i.IonForm] = 0
		ModIons[i.IonForm] = 0
		UnModIons[i.IonForm] = 0
	}

	// range over PSMs looking for modified and not modified evidences
	// if they exist on the ions map, get the numbers
	for _, i := range evi.PSM {

		// check the map
		_, ok := AllIons[i.IonForm]
		if ok {

			if i.Massdiff >= -0.99 && i.Massdiff <= 0.99 {
				UnModIons[i.IonForm]++
			} else {
				ModIons[i.IonForm]++
			}

		}
	}

	return
}

// UpdateLayerswithDatabase will fix the protein and gene assignments based on the database data
func (evi *Evidence) UpdateLayerswithDatabase(decoyTag string) {

	var dtb dat.Base
	dtb.Restore()

	var proteinIDMap = make(map[string]string)
	var entryNameMap = make(map[string]string)
	var geneMap = make(map[string]string)
	var descriptionMap = make(map[string]string)

	for _, j := range dtb.Records {
		if j.IsDecoy == false {
			proteinIDMap[j.PartHeader] = j.ID
			entryNameMap[j.PartHeader] = j.EntryName
			geneMap[j.PartHeader] = j.GeneNames
			descriptionMap[j.PartHeader] = j.Description
		}
	}

	for i := range evi.PSM {

		id := evi.PSM[i].Protein
		if evi.PSM[i].IsDecoy {
			id = strings.Replace(id, decoyTag, "", 1)
		}

		evi.PSM[i].ProteinID = proteinIDMap[id]
		evi.PSM[i].EntryName = entryNameMap[id]
		evi.PSM[i].GeneName = geneMap[id]
		evi.PSM[i].ProteinDescription = descriptionMap[id]

		for k := range evi.PSM[i].MappedProteins {
			evi.PSM[i].MappedGenes[geneMap[k]] = 0
		}
	}

	for i := range evi.Ions {

		id := evi.Ions[i].Protein
		if evi.Ions[i].IsDecoy {
			id = strings.Replace(id, decoyTag, "", 1)
		}

		evi.Ions[i].ProteinID = proteinIDMap[id]
		evi.Ions[i].EntryName = entryNameMap[id]
		evi.Ions[i].GeneName = geneMap[id]
		evi.Ions[i].ProteinDescription = descriptionMap[id]

		for k := range evi.Ions[i].MappedProteins {
			evi.Ions[i].MappedGenes[geneMap[k]] = 0
		}
	}

	for i := range evi.Peptides {

		id := evi.Peptides[i].Protein
		if evi.Peptides[i].IsDecoy {
			id = strings.Replace(id, decoyTag, "", 1)
		}

		evi.Peptides[i].ProteinID = proteinIDMap[id]
		evi.Peptides[i].EntryName = entryNameMap[id]
		evi.Peptides[i].GeneName = geneMap[id]
		evi.Peptides[i].ProteinDescription = descriptionMap[id]

		for k := range evi.Peptides[i].MappedProteins {
			evi.Peptides[i].MappedGenes[geneMap[k]] = 0
		}
	}

	return
}

// UpdateSupportingSpectra pushes back from PSM to Protein the new supporting spectra from razor results
func (evi *Evidence) UpdateSupportingSpectra() {

	var ptSupSpec = make(map[string][]string)
	var uniqueSpec = make(map[string][]string)
	var razorSpec = make(map[string][]string)

	for _, i := range evi.PSM {

		_, ok := ptSupSpec[i.Protein]
		if !ok {
			ptSupSpec[i.Protein] = append(ptSupSpec[i.Protein], i.Spectrum)
		} else {
			ptSupSpec[i.Protein] = append(ptSupSpec[i.Protein], i.Spectrum)
		}

		if i.IsUnique == true {
			_, ok := uniqueSpec[i.IonForm]
			if !ok {
				uniqueSpec[i.IonForm] = append(uniqueSpec[i.IonForm], i.Spectrum)
			} else {
				uniqueSpec[i.IonForm] = append(uniqueSpec[i.IonForm], i.Spectrum)
			}
		}

		if i.IsURazor == true {
			_, ok := razorSpec[i.IonForm]
			if !ok {
				razorSpec[i.IonForm] = append(razorSpec[i.IonForm], i.Spectrum)
			} else {
				razorSpec[i.IonForm] = append(razorSpec[i.IonForm], i.Spectrum)
			}
		}

	}

	for i := range evi.Proteins {
		for j := range evi.Proteins[i].TotalPeptideIons {

			if len(evi.Proteins[i].TotalPeptideIons[j].Spectra) == 0 {
				delete(evi.Proteins[i].TotalPeptideIons, j)
			}
		}
	}

	for i := range evi.Proteins {

		v, ok := ptSupSpec[evi.Proteins[i].PartHeader]
		if ok {
			for _, j := range v {
				evi.Proteins[i].SupportingSpectra[j] = 0
			}
		}

		for k := range evi.Proteins[i].TotalPeptideIons {

			Up, UOK := uniqueSpec[evi.Proteins[i].TotalPeptideIons[k].IonForm]
			if UOK && evi.Proteins[i].TotalPeptideIons[k].IsUnique == true {
				for _, l := range Up {
					evi.Proteins[i].TotalPeptideIons[k].Spectra[l] = 0
				}
			}

			Rp, ROK := razorSpec[evi.Proteins[i].TotalPeptideIons[k].IonForm]
			if ROK && evi.Proteins[i].TotalPeptideIons[k].IsURazor == true {
				for _, l := range Rp {
					evi.Proteins[i].TotalPeptideIons[k].Spectra[l] = 0
				}
			}

		}

	}

	return
}

// UpdatePeptideModCount counts how many times each peptide is observed modified and not modified
func (evi *Evidence) UpdatePeptideModCount() {

	// recreate the ion list from the main report object
	var all = make(map[string]int)
	var mod = make(map[string]int)
	var unmod = make(map[string]int)

	for _, i := range evi.Peptides {
		all[i.Sequence] = 0
		mod[i.Sequence] = 0
		unmod[i.Sequence] = 0
	}

	// range over PSMs looking for modified and not modified evidences
	// if they exist on the ions map, get the numbers
	for _, i := range evi.PSM {

		_, ok := all[i.Peptide]
		if ok {

			if i.Massdiff >= -0.99 && i.Massdiff <= 0.99 {
				unmod[i.Peptide]++
			} else {
				mod[i.Peptide]++
			}

		}
	}

	for i := range evi.Peptides {

		v1, ok1 := unmod[evi.Peptides[i].Sequence]
		if ok1 {
			evi.Peptides[i].UnModifiedObservations = v1
		}

		v2, ok2 := mod[evi.Peptides[i].Sequence]
		if ok2 {
			evi.Peptides[i].ModifiedObservations = v2
		}

	}

	return
}
