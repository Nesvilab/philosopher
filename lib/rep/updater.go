package rep

import (
	"fmt"
	"strings"

	"philosopher/lib/dat"
	"philosopher/lib/id"
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

// UpdateNumberOfEnzymaticTermini collects the NTT from ProteinProphet
// and passes along to the final Protein structure.
func (evi *Evidence) UpdateNumberOfEnzymaticTermini() {

	// restore the original prot.xml output
	var p id.ProtIDList
	p.Restore()

	// collect the updated ntt for each peptide-protein pair
	var nttPeptidetoProptein = make(map[string]uint8)

	for _, i := range p {
		for _, j := range i.PeptideIons {
			if !strings.Contains(i.ProteinName, "rev_") {
				key := fmt.Sprintf("%s#%s", j.PeptideSequence, i.ProteinName)
				nttPeptidetoProptein[key] = j.NumberOfEnzymaticTermini
			}
		}
	}

	for i := range evi.PSM {

		key := fmt.Sprintf("%s#%s", evi.PSM[i].Peptide, evi.PSM[i].Protein)
		ntt, ok := nttPeptidetoProptein[key]
		if ok {
			evi.PSM[i].NumberOfEnzymaticTermini = int(ntt)
		}
	}

	return
}

// UpdateIonStatus pushes back to ion and psm evideces the uniqueness and razorness status of each peptide and ion
func (evi *Evidence) UpdateIonStatus(decoyTag string) {

	var uniqueMap = make(map[string]bool)
	var urazorMap = make(map[string]string)
	var uniqueSeqMap = make(map[string]string)

	var PSMtoDelete []int
	var PeptidetoDelete []int
	var IontoDelete []int

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

		// the decoy tag checking is a failsafe mechanism to avoid proteins
		// with real complex razor case decisions to pass downstream
		// wrong classifications. If by any chance the protein gets assigned to
		// a razor decoy, this mechanism avoids the replacement

		rp, rOK := urazorMap[evi.PSM[i].IonForm]
		if rOK {

			evi.PSM[i].IsURazor = true

			// we found cases where the peptide maps to both target and decoy but is
			// assigned as razor to the decoy. the IF statement below replaces the
			// decoy by the target but it was removed because in some cases the protein
			// does not pass the FDR filtering.

			evi.PSM[i].MappedProteins[evi.PSM[i].Protein] = 0
			delete(evi.PSM[i].MappedProteins, rp)
			evi.PSM[i].Protein = rp

			if strings.Contains(rp, decoyTag) {
				evi.PSM[i].IsDecoy = true
			}
		}

		_, uOK := uniqueMap[evi.PSM[i].IonForm]
		if uOK {
			evi.PSM[i].IsUnique = true
		}

		if !rOK && !uOK {
			PSMtoDelete = append(PSMtoDelete, i)
		} else {
			uniqueSeqMap[evi.PSM[i].Peptide] = evi.PSM[i].Protein
		}
	}

	// remove PSMs mapping to Proteins that did not passed the FDR filter
	for _, i := range PSMtoDelete {
		evi.PSM = RemovePSMByIndex(evi.PSM, i)
	}

	for i := range evi.Ions {

		rp, rOK := urazorMap[evi.Ions[i].IonForm]
		if rOK {

			evi.Ions[i].IsURazor = true

			evi.Ions[i].MappedProteins[evi.Ions[i].Protein] = 0
			delete(evi.Ions[i].MappedProteins, rp)
			evi.Ions[i].Protein = rp

			if strings.Contains(rp, decoyTag) {
				evi.Ions[i].IsDecoy = true
			}

		}
		_, uOK := uniqueMap[evi.Ions[i].IonForm]
		if uOK {
			evi.Ions[i].IsUnique = true
		} else {
			evi.Ions[i].IsUnique = false
		}

		if !rOK && !uOK {
			IontoDelete = append(IontoDelete, i)
		}
	}

	// remove PSMs mapping to Proteins that did not passed the FDR filter
	for _, i := range IontoDelete {
		evi.Ions = RemoveIonsByIndex(evi.Ions, i)
	}

	for i := range evi.Peptides {

		v, ok := uniqueSeqMap[evi.Peptides[i].Sequence]
		if ok {
			evi.Peptides[i].MappedProteins[evi.Peptides[i].Protein] = 0
			delete(evi.Peptides[i].MappedProteins, v)
			evi.Peptides[i].Protein = v
		}

		if strings.Contains(v, decoyTag) {
			evi.Peptides[i].IsDecoy = true
		}

		if !ok {
			PeptidetoDelete = append(PeptidetoDelete, i)
		}
	}

	// remove Peptides mapping to Proteins that did not passed the FDR filter
	for _, i := range PeptidetoDelete {
		evi.Peptides = RemovePeptidesByIndex(evi.Peptides, i)
	}

	return
}

// func (evi *Evidence) UpdateIonStatus(decoyTag string) {

// 	var uniqueMap = make(map[string]bool)
// 	var urazorMap = make(map[string]string)
// 	var uniqueSeqMap = make(map[string]string)

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

// 		// the decoy tag checking is a failsafe mechanism to avoid proteins
// 		// with real complex razor case decisions to pass downstream
// 		// wrong classifications. If by any chance the protein gets assigned to
// 		// a razor decoy, this mechanism avoids the replacement

// 		rp, rOK := urazorMap[evi.PSM[i].IonForm]
// 		if rOK {

// 			evi.PSM[i].IsURazor = true

// 			// we found cases where the peptide maps to both target and decoy but is
// 			// assigned as razor to the decoy. the IF statement below replaces the
// 			// decoy by the target but it was removed because in some cases the protein
// 			// does not pass the FDR filtering.

// 			evi.PSM[i].MappedProteins[evi.PSM[i].Protein] = 0
// 			delete(evi.PSM[i].MappedProteins, rp)
// 			evi.PSM[i].Protein = rp

// 			if strings.Contains(rp, decoyTag) {
// 				evi.PSM[i].IsDecoy = true
// 			}
// 		}

// 		_, uOK := uniqueMap[evi.PSM[i].IonForm]
// 		if uOK {
// 			evi.PSM[i].IsUnique = true
// 		}

// 		uniqueSeqMap[evi.PSM[i].Peptide] = evi.PSM[i].Protein
// 	}

// 	var uniques1 = make(map[string]uint8)
// 	for _, x := range evi.PSM {
// 		if !x.IsDecoy {
// 			uniques1[x.Protein] = 0
// 		}
// 	}
// 	fmt.Println(len(uniques1))

// 	for i := range evi.Ions {

// 		rp, rOK := urazorMap[evi.Ions[i].IonForm]
// 		if rOK {

// 			evi.Ions[i].IsURazor = true

// 			evi.Ions[i].MappedProteins[evi.Ions[i].Protein] = 0
// 			delete(evi.Ions[i].MappedProteins, rp)
// 			evi.Ions[i].Protein = rp

// 			if strings.Contains(rp, decoyTag) {
// 				evi.Ions[i].IsDecoy = true
// 			}

// 		}
// 		_, uOK := uniqueMap[evi.Ions[i].IonForm]
// 		if uOK {
// 			evi.Ions[i].IsUnique = true
// 		} else {
// 			evi.Ions[i].IsUnique = false
// 		}
// 	}

// 	for i := range evi.Peptides {

// 		v, ok := uniqueSeqMap[evi.Peptides[i].Sequence]
// 		if ok {
// 			evi.Peptides[i].MappedProteins[evi.Peptides[i].Protein] = 0
// 			delete(evi.Peptides[i].MappedProteins, v)
// 			evi.Peptides[i].Protein = v
// 		}

// 		if strings.Contains(v, decoyTag) {
// 			evi.Peptides[i].IsDecoy = true
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

		// update mapped genes
		for k := range evi.PSM[i].MappedProteins {
			if !strings.Contains(k, decoyTag) {
				evi.PSM[i].MappedGenes[geneMap[k]] = 0
			}
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

		// update mapped genes
		for k := range evi.Ions[i].MappedProteins {
			if !strings.Contains(k, decoyTag) {
				evi.Ions[i].MappedGenes[geneMap[k]] = 0
			}
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

		// update mapped genes
		for k := range evi.Peptides[i].MappedProteins {
			if !strings.Contains(k, decoyTag) {
				evi.Peptides[i].MappedGenes[geneMap[k]] = 0
			}
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
