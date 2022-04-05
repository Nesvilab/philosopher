package rep

import (
	"strings"

	"philosopher/lib/dat"
	"philosopher/lib/id"
	"philosopher/lib/uti"
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
func (evi *Evidence) UpdateNumberOfEnzymaticTermini(decoyTag string) {

	// restore the original prot.xml output
	var p id.ProtIDList
	p.Restore()

	// collect the updated ntt for each peptide-protein pair
	//var nttPeptidetoProptein = make(map[string]uint8)
	type k struct {
		a string
		b string
	}
	var nttPeptidetoProptein = make(map[k]uint8)

	for _, i := range p {
		for _, j := range i.PeptideIons {
			if !strings.Contains(i.ProteinName, decoyTag) {
				nttPeptidetoProptein[k{j.PeptideSequence, i.ProteinName}] = j.NumberOfEnzymaticTermini
			}
		}
	}

	for i := range evi.PSM {
		if ntt, ok := nttPeptidetoProptein[k{evi.PSM[i].Peptide, evi.PSM[i].Protein}]; ok {
			evi.PSM[i].NumberOfEnzymaticTermini = ntt
		}
	}
}

// UpdateIonStatus pushes back to ion and psm evideces the uniqueness and razorness status of each peptide and ion
func (evi *Evidence) UpdateIonStatus(decoyTag string) {

	var uniqueIons = make(map[id.IonFormType]bool)
	var razorIons = make(map[id.IonFormType]string)
	var uniquePeptides = make(map[string]string)
	var razorPeptides = make(map[string]string)

	for _, i := range evi.Proteins {
		for _, j := range i.TotalPeptideIons {

			if j.IsUnique {
				uniqueIons[j.IonForm()] = true
				uniquePeptides[j.Sequence] = i.PartHeader
			}

			if j.IsURazor {
				razorIons[j.IonForm()] = i.PartHeader
				razorPeptides[j.Sequence] = i.PartHeader
			}
		}
	}

	for i := range evi.PSM {
		// the decoy tag checking is a failsafe mechanism to avoid proteins
		// with real complex razor case decisions to pass downstream
		// wrong classifications. If by any chance the protein gets assigned to
		// a razor decoy, this mechanism avoids the replacement

		rp, rOK := razorIons[evi.PSM[i].IonForm()]
		if rOK {

			evi.PSM[i].IsURazor = true

			// we found cases where the peptide maps to both target and decoy but is
			// assigned as razor to the decoy. the IF statement below replaces the
			// decoy by the target but it was removed because in some cases the protein
			// does not pass the FDR filtering.

			evi.PSM[i].MappedProteins[evi.PSM[i].Protein] = 0
			delete(evi.PSM[i].MappedProteins, rp)
			evi.PSM[i].Protein = rp

			// if strings.Contains(rp, decoyTag) {
			// 	evi.PSM[i].IsDecoy = true
			// } else {
			// 	evi.PSM[i].IsDecoy = false
			// }
		}

		if !evi.PSM[i].IsURazor {
			sp, sOK := razorPeptides[evi.PSM[i].Peptide]
			if sOK {

				evi.PSM[i].IsURazor = true

				// we found cases where the peptide maps to both target and decoy but is
				// assigned as razor to the decoy. the IF statement below replaces the
				// decoy by the target but it was removed because in some cases the protein
				// does not pass the FDR filtering.

				evi.PSM[i].MappedProteins[evi.PSM[i].Protein] = 0
				delete(evi.PSM[i].MappedProteins, sp)
				evi.PSM[i].Protein = sp

				if strings.Contains(sp, decoyTag) {
					evi.PSM[i].IsDecoy = true
				}
			}

			_, uOK := uniqueIons[evi.PSM[i].IonForm()]
			if uOK {
				evi.PSM[i].IsUnique = true
			}

			uniquePeptides[evi.PSM[i].Peptide] = evi.PSM[i].Protein
		}
	}

	for i := range evi.Ions {
		rp, rOK := razorIons[evi.Ions[i].IonForm()]
		if rOK {

			evi.Ions[i].IsURazor = true

			evi.Ions[i].MappedProteins[evi.Ions[i].Protein] = 0
			delete(evi.Ions[i].MappedProteins, rp)
			evi.Ions[i].Protein = rp

			if strings.Contains(rp, decoyTag) {
				evi.Ions[i].IsDecoy = true
			}

		}
		_, uOK := uniqueIons[evi.Ions[i].IonForm()]
		if uOK {
			evi.Ions[i].IsUnique = true
		} else {
			evi.Ions[i].IsUnique = false
		}
	}

	for i := range evi.Peptides {
		rp, rOK := razorPeptides[evi.Peptides[i].Sequence]
		if rOK {

			evi.Peptides[i].IsURazor = true

			evi.Peptides[i].MappedProteins[evi.Peptides[i].Protein] = 0
			delete(evi.Peptides[i].MappedProteins, rp)
			evi.Peptides[i].Protein = rp

			if strings.Contains(rp, decoyTag) {
				evi.Peptides[i].IsDecoy = true
			}

		}
		_, uOK := uniquePeptides[evi.Peptides[i].Sequence]
		if uOK {
			evi.Peptides[i].IsUnique = true
		} else {
			evi.Peptides[i].IsUnique = false
		}
	}
}

// UpdateIonModCount counts how many times each ion is observed modified and not modified
func (evi *Evidence) UpdateIonModCount() {

	// recreate the ion list from the main report object
	var AllIons = make(map[id.IonFormType]struct{})
	var ModIons = make(map[id.IonFormType]int)
	var UnModIons = make(map[id.IonFormType]int)

	for _, i := range evi.Ions {
		AllIons[i.IonForm()] = struct{}{}
		ModIons[i.IonForm()] = 0
		UnModIons[i.IonForm()] = 0
	}

	// range over PSMs looking for modified and not modified evidences
	// if they exist on the ions map, get the numbers
	for _, i := range evi.PSM {

		// check the map
		_, ok := AllIons[i.IonForm()]
		if ok {

			if i.Massdiff >= -0.99 && i.Massdiff <= 0.99 {
				UnModIons[i.IonForm()]++
			} else {
				ModIons[i.IonForm()]++
			}

		}
	}
}

// SyncPSMToProteins forces the synchronization between the filtered proteins, and the remaining structures.
func (evi *Evidence) SyncPSMToProteins(decoy string) {
	var totalSpc = make(map[string][]id.SpectrumType, len(evi.PSM))
	var uniqueSpc = make(map[string][]id.SpectrumType, len(evi.PSM))
	var razorSpc = make(map[string][]id.SpectrumType, len(evi.PSM))

	var totalPeptides = make(map[string][]string, len(evi.PSM))
	var uniquePeptides = make(map[string][]string, len(evi.PSM))
	var razorPeptides = make(map[string][]string, len(evi.PSM))

	var proteinIndex = make(map[string]struct{})

	for _, i := range evi.Proteins {
		proteinIndex[i.PartHeader] = struct{}{}
	}

	// for _, i := range evi.PSM {
	// 	_, ok := proteinIndex[i.Protein]
	// 	if ok {
	// 		newPSM = append(newPSM, i)
	// 	}
	// }
	// evi.PSM = newPSM

	// for _, i := range evi.Ions {
	// 	_, ok := proteinIndex[i.Protein]
	// 	if ok {
	// 		newIons = append(newIons, i)
	// 	}
	// }
	// evi.Ions = newIons

	// for _, i := range evi.Peptides {
	// 	_, ok := proteinIndex[i.Protein]
	// 	if ok {
	// 		newPeptides = append(newPeptides, i)
	// 	}
	// }
	// evi.Peptides = newPeptides

	for _, i := range evi.PSM {
		//if !i.IsDecoy {

		// Total
		totalSpc[i.Protein] = append(totalSpc[i.Protein], i.SpectrumFileName())
		totalPeptides[i.Protein] = append(totalPeptides[i.Protein], i.Peptide)
		for j := range i.MappedProteins {
			totalSpc[j] = append(totalSpc[j], i.SpectrumFileName())
			totalPeptides[j] = append(totalPeptides[j], i.Peptide)
		}

		if i.IsUnique {
			uniqueSpc[i.Protein] = append(uniqueSpc[i.Protein], i.SpectrumFileName())
			uniquePeptides[i.Protein] = append(uniquePeptides[i.Protein], i.Peptide)
		}

		if i.IsURazor {
			razorSpc[i.Protein] = append(razorSpc[i.Protein], i.SpectrumFileName())
			razorPeptides[i.Protein] = append(razorPeptides[i.Protein], i.Peptide)
		}
		//}
	}

	for k, v := range totalPeptides {
		totalPeptides[k] = uti.RemoveDuplicateStrings(v)
	}

	for k, v := range uniquePeptides {
		uniquePeptides[k] = uti.RemoveDuplicateStrings(v)
	}

	for k, v := range razorPeptides {
		razorPeptides[k] = uti.RemoveDuplicateStrings(v)
	}

	for i := range evi.Proteins {

		evi.Proteins[i].SupportingSpectra = make(map[id.SpectrumType]int)
		evi.Proteins[i].TotalSpC = 0
		evi.Proteins[i].UniqueSpC = 0
		evi.Proteins[i].URazorSpC = 0
		evi.Proteins[i].TotalPeptides = make(map[string]int)
		evi.Proteins[i].UniquePeptides = make(map[string]int)
		evi.Proteins[i].URazorPeptides = make(map[string]int)

		if v, ok := totalSpc[evi.Proteins[i].PartHeader]; ok {
			evi.Proteins[i].TotalSpC += len(v)
			for _, j := range v {
				evi.Proteins[i].SupportingSpectra[j]++
			}
		}

		if v, ok := totalPeptides[evi.Proteins[i].PartHeader]; ok {
			for _, j := range v {
				evi.Proteins[i].TotalPeptides[j]++
			}
		}

		if v, ok := uniqueSpc[evi.Proteins[i].PartHeader]; ok {
			evi.Proteins[i].UniqueSpC += len(v)
		}

		if v, ok := uniquePeptides[evi.Proteins[i].PartHeader]; ok {
			for _, j := range v {
				evi.Proteins[i].UniquePeptides[j]++
			}
		}

		if v, ok := razorSpc[evi.Proteins[i].PartHeader]; ok {
			evi.Proteins[i].URazorSpC += len(v)
		}

		if v, ok := razorPeptides[evi.Proteins[i].PartHeader]; ok {
			for _, j := range v {
				evi.Proteins[i].URazorPeptides[j]++
			}
		}
	}

	{
		proteinIndex = make(map[string]struct{}, len(evi.Proteins))
		newProteins := make([]int, 0, len(evi.Proteins))
		for idx, i := range evi.Proteins {
			if len(i.SupportingSpectra) > 0 {
				proteinIndex[i.PartHeader] = struct{}{}
				newProteins = append(newProteins, idx)
			}
		}
		for idx, i := range newProteins {
			evi.Proteins[idx] = evi.Proteins[i]
		}
		evi.Proteins = evi.Proteins[:len(newProteins)]
	}
	{
		newPSM := make([]int, 0, len(evi.PSM))
		for idx, i := range evi.PSM {
			if _, ok := proteinIndex[i.Protein]; ok {
				newPSM = append(newPSM, idx)
			}
		}
		for idx, i := range newPSM {
			evi.PSM[idx] = evi.PSM[i]
		}
		evi.PSM = evi.PSM[:len(newPSM)]
	}
	{
		newIons := make([]int, 0, len(evi.Ions))
		for idx, i := range evi.Ions {
			if _, ok := proteinIndex[i.Protein]; ok {
				newIons = append(newIons, idx)
			}
		}
		for idx, i := range newIons {
			evi.Ions[idx] = evi.Ions[i]
		}
		evi.Ions = evi.Ions[:len(newIons)]
	}
	{
		newPeptides := make([]int, 0, len(evi.Peptides))
		for idx, i := range evi.Peptides {
			if _, ok := proteinIndex[i.Protein]; ok {
				newPeptides = append(newPeptides, idx)
			}
		}
		for idx, i := range newPeptides {
			evi.Peptides[idx] = evi.Peptides[i]
		}
		evi.Peptides = evi.Peptides[:len(newPeptides)]
	}
}

// SyncPSMToPeptides forces the synchronization between the filtered peptides, and the remaining structures.
func (evi Evidence) SyncPSMToPeptides(decoy string) Evidence {

	var spectra = make(map[string][]id.SpectrumType)

	for _, i := range evi.PSM {
		if !i.IsDecoy {
			spectra[i.Peptide] = append(spectra[i.Peptide], i.SpectrumFileName())
		}
	}

	for i := range evi.Peptides {

		evi.Peptides[i].Spc = 0
		evi.Peptides[i].Spectra = make(map[id.SpectrumType]uint8)

		if v, ok := spectra[evi.Peptides[i].Sequence]; ok {

			//evi.Peptides[i].IsDecoy = false

			for _, j := range v {
				evi.Peptides[i].Spectra[j]++
			}
			evi.Peptides[i].Spc = len(v)

		}
	}

	return evi
}

// SyncPSMToPeptideIons forces the synchronization between the filtered ions, and the remaining structures.
func (evi Evidence) SyncPSMToPeptideIons(decoy string) Evidence {

	var spectra = make(map[id.IonFormType][]id.SpectrumType)

	for _, i := range evi.PSM {
		if !i.IsDecoy {
			spectra[i.IonForm()] = append(spectra[i.IonForm()], i.SpectrumFileName())
		}
	}

	for i := range evi.Ions {

		evi.Ions[i].Spectra = make(map[id.SpectrumType]int)

		v, ok := spectra[evi.Ions[i].IonForm()]
		if ok {

			//evi.Ions[i].IsDecoy = false

			for _, j := range v {
				evi.Ions[i].Spectra[j]++
			}
		}
	}

	return evi
}

// UpdateLayerswithDatabase will fix the protein and gene assignments based on the database data
func (evi *Evidence) UpdateLayerswithDatabase(decoyTag string) {
	type liteRecord struct {
		ID          string
		EntryName   string
		GeneNames   string
		Description string
		Sequence    string
	}
	var recordMap = make(map[string]liteRecord)

	{
		var dtb dat.Base
		dtb.Restore()
		for _, j := range dtb.Records {
			recordMap[j.PartHeader] = liteRecord{j.ID, j.EntryName, j.GeneNames, strings.TrimSpace(j.Description), j.Sequence}
		}
	}

	type prevNextAA struct {
		prev byte
		next byte
	}
	var pepPrevNextAA = make(map[string]prevNextAA)

	replacerIL := strings.NewReplacer("L", "I")
	for i := range evi.PSM {

		rec := recordMap[evi.PSM[i].Protein]
		evi.PSM[i].ProteinID = rec.ID
		evi.PSM[i].EntryName = rec.EntryName
		evi.PSM[i].GeneName = rec.GeneNames
		evi.PSM[i].ProteinDescription = rec.Description

		// update mapped genes
		for k := range evi.PSM[i].MappedProteins {
			if !strings.Contains(k, decoyTag) {
				evi.PSM[i].MappedGenes[recordMap[k].GeneNames] = struct{}{}
			}
		}

		// map the peptide to the protein
		mstart := strings.Index(replacerIL.Replace(rec.Sequence), replacerIL.Replace(evi.PSM[i].Peptide))
		mend := mstart + len(evi.PSM[i].Peptide)

		if mstart != -1 {
			evi.PSM[i].ProteinStart = mstart + 1
			evi.PSM[i].ProteinEnd = mend

			if (mstart) <= 0 {
				evi.PSM[i].PrevAA = '-'
			} else {
				evi.PSM[i].PrevAA = rec.Sequence[mstart-1]
			}

			if (mend + 1) >= len(rec.Sequence) {
				evi.PSM[i].NextAA = '-'
			} else {
				evi.PSM[i].NextAA = rec.Sequence[mend]
			}

		}

		pepPrevNextAA[evi.PSM[i].Peptide] = prevNextAA{evi.PSM[i].PrevAA, evi.PSM[i].NextAA}

	}

	for i := range evi.Ions {

		id := evi.Ions[i].Protein
		if evi.Ions[i].IsDecoy {
			id = strings.Replace(id, decoyTag, "", 1)
		}
		rec := recordMap[id]

		evi.Ions[i].ProteinID = rec.ID
		evi.Ions[i].EntryName = rec.EntryName
		evi.Ions[i].GeneName = rec.GeneNames
		evi.Ions[i].ProteinDescription = rec.Description

		// update mapped genes
		for k := range evi.Ions[i].MappedProteins {
			if !strings.Contains(k, decoyTag) {
				evi.Ions[i].MappedGenes[recordMap[k].GeneNames] = struct{}{}
			}
		}
		pnAA := pepPrevNextAA[evi.Ions[i].Sequence]
		evi.Ions[i].PrevAA = pnAA.prev
		evi.Ions[i].NextAA = pnAA.next
	}

	for i := range evi.Peptides {

		id := evi.Peptides[i].Protein
		if evi.Peptides[i].IsDecoy {
			id = strings.Replace(id, decoyTag, "", 1)
		}
		rec := recordMap[id]
		evi.Peptides[i].ProteinID = rec.ID
		evi.Peptides[i].EntryName = rec.EntryName
		evi.Peptides[i].GeneName = rec.GeneNames
		evi.Peptides[i].ProteinDescription = rec.Description

		// update mapped genes
		for k := range evi.Peptides[i].MappedProteins {
			if !strings.Contains(k, decoyTag) {
				evi.Peptides[i].MappedGenes[recordMap[k].GeneNames] = struct{}{}
			}
		}
		pnAA := pepPrevNextAA[evi.Peptides[i].Sequence]
		evi.Peptides[i].PrevAA = pnAA.prev
		evi.Peptides[i].NextAA = pnAA.next
	}
}

// UpdateSupportingSpectra pushes back from PSM to Protein the new supporting spectra from razor results
func (evi *Evidence) UpdateSupportingSpectra() {

	var ptSupSpec = make(map[string][]id.SpectrumType)
	var uniqueSpec = make(map[id.IonFormType][]id.SpectrumType)
	var razorSpec = make(map[id.IonFormType][]id.SpectrumType)

	var totalPeptides = make(map[string][]string)
	var uniquePeptides = make(map[string][]string)
	var razorPeptides = make(map[string][]string)

	for _, i := range evi.PSM {

		_, ok := ptSupSpec[i.Protein]
		if !ok {
			ptSupSpec[i.Protein] = append(ptSupSpec[i.Protein], i.SpectrumFileName())
		}

		if i.IsUnique {
			_, ok := uniqueSpec[i.IonForm()]
			if !ok {
				uniqueSpec[i.IonForm()] = append(uniqueSpec[i.IonForm()], i.SpectrumFileName())
			}
		}

		if i.IsURazor {
			_, ok := razorSpec[i.IonForm()]
			if !ok {
				razorSpec[i.IonForm()] = append(razorSpec[i.IonForm()], i.SpectrumFileName())
			}
		}
	}

	for _, i := range evi.Peptides {

		totalPeptides[i.Protein] = append(totalPeptides[i.Protein], i.Sequence)
		for j := range i.MappedProteins {
			totalPeptides[j] = append(totalPeptides[j], i.Sequence)
		}

		if i.IsUnique {
			uniquePeptides[i.Protein] = append(uniquePeptides[i.Protein], i.Sequence)
		}

		if i.IsURazor {
			razorPeptides[i.Protein] = append(razorPeptides[i.Protein], i.Sequence)
		}
	}

	for k, v := range totalPeptides {
		totalPeptides[k] = uti.RemoveDuplicateStrings(v)
	}

	for k, v := range uniquePeptides {
		uniquePeptides[k] = uti.RemoveDuplicateStrings(v)
	}

	for k, v := range razorPeptides {
		razorPeptides[k] = uti.RemoveDuplicateStrings(v)
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

			Up, UOK := uniqueSpec[evi.Proteins[i].TotalPeptideIons[k].IonForm()]
			if UOK && evi.Proteins[i].TotalPeptideIons[k].IsUnique {
				for _, l := range Up {
					evi.Proteins[i].TotalPeptideIons[k].Spectra[l] = 0
				}
			}

			Rp, ROK := razorSpec[evi.Proteins[i].TotalPeptideIons[k].IonForm()]
			if ROK && evi.Proteins[i].TotalPeptideIons[k].IsURazor {
				for _, l := range Rp {
					evi.Proteins[i].TotalPeptideIons[k].Spectra[l] = 0
				}
			}

		}

		vTP, okTP := totalPeptides[evi.Proteins[i].PartHeader]
		if okTP {
			for _, j := range vTP {
				evi.Proteins[i].TotalPeptides[j]++
			}
		}

		vuP, okuP := uniquePeptides[evi.Proteins[i].PartHeader]
		if okuP {
			for _, j := range vuP {
				evi.Proteins[i].UniquePeptides[j]++
			}
		}

		vRP, okRP := razorPeptides[evi.Proteins[i].PartHeader]
		if okRP {
			for _, j := range vRP {
				evi.Proteins[i].URazorPeptides[j]++
			}
		}

	}
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
}
