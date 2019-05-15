package rep

import (
	"github.com/prvst/philosopher/lib/dat"
)

// PeptideMap struct
type PeptideMap struct {
	Sequence     string
	Protein      string
	RazorProtein string
	Proteins     map[string]int
}

// UpdateMappedProteins updates the list of mapped proteins on the data structures
func (e *Evidence) UpdateMappedProteins() {

	var list = make(map[string]PeptideMap)
	var checkup = make(map[string]int)

	for _, i := range e.Proteins {
		for _, v := range i.TotalPeptideIons {

			_, ok := checkup[v.Sequence]
			if !ok {
				var pm PeptideMap

				pm.Sequence = v.Sequence
				pm.Proteins = v.MappedProteins
				pm.Proteins[i.PartHeader] = 0
				pm.RazorProtein = i.PartHeader

				list[pm.Sequence] = pm
				checkup[v.Sequence] = 0
			}
		}
	}

	for i := range e.PSM {
		v, ok := list[e.PSM[i].Peptide]
		if ok {

			for k := range v.Proteins {
				e.PSM[i].MappedProteins[k]++
			}

			if len(e.PSM[i].RazorProtein) < 1 {
				e.PSM[i].RazorProtein = v.RazorProtein
				e.PSM[i].IsURazor = true
			}
		}
	}

	for i := range e.Peptides {
		v, ok := list[e.Peptides[i].Sequence]
		if ok {

			for k := range v.Proteins {
				e.Peptides[i].MappedProteins[k]++
			}

			e.Peptides[i].Protein = v.RazorProtein
		}
	}

	for i := range e.Ions {
		v, ok := list[e.Ions[i].Sequence]
		if ok {

			for k := range v.Proteins {
				e.Ions[i].MappedProteins[k]++
			}

			e.Ions[i].Protein = v.RazorProtein
		}
	}

	return
}

// func (e *Evidence) UpdateMappedProteins() {

// 	var list []PeptideMap
// 	var checkup = make(map[string]int)

// 	log.Println("check-1")

// 	for _, i := range e.Proteins {
// 		for _, v := range i.TotalPeptideIons {

// 			_, ok := checkup[v.Sequence]
// 			if !ok {
// 				var pm PeptideMap

// 				pm.Sequence = v.Sequence
// 				pm.Proteins = v.MappedProteins
// 				pm.Proteins[i.PartHeader] = 0

// 				//if v.IsURazor == true {
// 				pm.RazorProtein = i.PartHeader
// 				//}

// 				list = append(list, pm)
// 				checkup[v.Sequence] = 0
// 			}

// 		}
// 	}

// 	log.Println("check-2")

// 	for i := range e.PSM {
// 		for _, j := range list {
// 			if e.PSM[i].Peptide == j.Sequence {

// 				for k := range j.Proteins {
// 					e.PSM[i].MappedProteins[k]++
// 				}

// 				if len(e.PSM[i].RazorProtein) < 1 {
// 					e.PSM[i].RazorProtein = j.RazorProtein
// 					e.PSM[i].IsURazor = true
// 				}
// 				break
// 			}
// 		}

// 	}

// 	log.Println("check-3")

// 	for i := range e.Peptides {
// 		for _, j := range list {
// 			if e.Peptides[i].Sequence == j.Sequence {

// 				for k := range j.Proteins {
// 					e.Peptides[i].MappedProteins[k]++
// 				}

// 				e.Peptides[i].Protein = j.RazorProtein

// 				break
// 			}
// 		}
// 	}

// 	log.Println("check-4")

// 	for i := range e.Ions {
// 		for _, j := range list {
// 			if e.Ions[i].Sequence == j.Sequence {

// 				for k := range j.Proteins {
// 					e.Ions[i].MappedProteins[k]++
// 				}

// 				e.Ions[i].Protein = j.RazorProtein

// 				break
// 			}
// 		}
// 	}

// 	log.Println("check-5")

// 	return
// }

// UpdateIonModCount counts how many times each ion is observed modified and not modified
func (e *Evidence) UpdateIonModCount() {

	// recreate the ion list from the main report object
	var AllIons = make(map[string]int)
	var ModIons = make(map[string]int)
	var UnModIons = make(map[string]int)

	for _, i := range e.Ions {
		AllIons[i.IonForm] = 0
		ModIons[i.IonForm] = 0
		UnModIons[i.IonForm] = 0
	}

	// range over PSMs looking for modified and not modified evidences
	// if they exist on the ions map, get the numbers
	for _, i := range e.PSM {

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

	// for i := range e.Ions {

	// 	v1, ok1 := UnModIons[e.Ions[i].IonForm]
	// 	if ok1 {
	// 		e.Ions[i].UnModifiedObservations = v1
	// 	}

	// 	v2, ok2 := ModIons[e.Ions[i].IonForm]
	// 	if ok2 {
	// 		e.Ions[i].ModifiedObservations = v2
	// 	}

	// }

	return
}

// UpdateProteinStatus assignes the razor protein to THE protein column and removes it form the alt Protein
// it basically swiches places with the current protein assignment from pepXML
func (e *Evidence) UpdateProteinStatus() {

	for i := range e.PSM {

		if e.PSM[i].IsURazor == true && e.PSM[i].Protein != e.PSM[i].RazorProtein {

			var altProteins []string

			// push the selected protein to the top fo the list
			altProteins = append(altProteins, e.PSM[i].Protein)

			// replace the selected protein by the razor one
			e.PSM[i].Protein = e.PSM[i].RazorProtein

		} else if e.PSM[i].IsURazor == false && e.PSM[i].Protein != e.PSM[i].RazorProtein {
			e.PSM[i].RazorProtein = e.PSM[i].Protein
		}
	}

	return
}

// UpdateGeneNames will fix the gene name assignment after razor assingment
func (e *Evidence) UpdateGeneNames() {

	var dtb dat.Base
	dtb.Restore()

	var dbMap = make(map[string]string)
	for _, j := range dtb.Records {
		dbMap[j.PartHeader] = j.GeneNames
	}

	var descMap = make(map[string]string)
	for _, j := range dtb.Records {
		descMap[j.PartHeader] = j.ProteinName
	}

	var idMap = make(map[string]string)
	for _, j := range dtb.Records {
		idMap[j.PartHeader] = j.ID
	}

	var entryMap = make(map[string]string)
	for _, j := range dtb.Records {
		entryMap[j.PartHeader] = j.EntryName
	}

	for i := range e.PSM {
		e.PSM[i].GeneName = dbMap[e.PSM[i].Protein]
		e.PSM[i].ProteinDescription = descMap[e.PSM[i].Protein]
		e.PSM[i].ProteinID = idMap[e.PSM[i].Protein]
		e.PSM[i].EntryName = entryMap[e.PSM[i].Protein]
	}

	for i := range e.Ions {
		e.Ions[i].GeneName = dbMap[e.Ions[i].Protein]
		e.Ions[i].ProteinDescription = descMap[e.Ions[i].Protein]
		e.Ions[i].ProteinID = idMap[e.Ions[i].Protein]
		e.Ions[i].EntryName = entryMap[e.Ions[i].Protein]
	}

	for i := range e.Peptides {
		e.Peptides[i].GeneName = dbMap[e.Peptides[i].Protein]
		e.Peptides[i].ProteinDescription = descMap[e.Peptides[i].Protein]
		e.Peptides[i].ProteinID = idMap[e.Peptides[i].Protein]
		e.Peptides[i].EntryName = entryMap[e.Peptides[i].Protein]
	}

	return
}

// UpdateIonStatus pushes back to ion and psm evideces the uniqueness and razorness status of each peptide and ion
func (e *Evidence) UpdateIonStatus() {

	var uniqueMap = make(map[string]bool)
	var urazorMap = make(map[string]string)
	var ptMap = make(map[string]string)

	for _, i := range e.Proteins {

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

	for i := range e.PSM {

		if len(e.PSM[i].MappedProteins) == 0 {
			e.PSM[i].IsUnique = true
		}

		_, uOK := uniqueMap[e.PSM[i].IonForm]
		if uOK {
			e.PSM[i].IsUnique = true
		}

		rp, rOK := urazorMap[e.PSM[i].IonForm]
		if rOK {
			e.PSM[i].IsURazor = true
			e.PSM[i].RazorProtein = rp
		}

		v, ok := ptMap[e.PSM[i].IonForm]
		if ok {
			e.PSM[i].Protein = v
		}
	}

	for i := range e.Ions {

		_, uOK := uniqueMap[e.Ions[i].IonForm]
		if uOK {
			e.Ions[i].IsUnique = true
		}

		_, rOK := urazorMap[e.Ions[i].IonForm]
		if rOK {
			e.Ions[i].IsURazor = true
		}

	}

	return
}

// UpdateIonAssignedAndObservedMods collects all Assigned and Observed modifications from
// individual PSM and assign them to ions
//func (e *Evidence) UpdateIonAssignedAndObservedMods() {

// var aMap = make(map[string][]string)
// var aPepMap = make(map[string][]string)

// // collect the assigned modifications from the PSM data
// for _, i := range e.PSM {
// 	for j := range i.AssignedModifications {
// 		aMap[i.IonForm] = append(aMap[i.IonForm], j)
// 		aPepMap[i.Peptide] = append(aPepMap[i.Peptide], j)
// 	}
// }

// // forward it to the Ion data
// for i := range e.Ions {
// 	v, ok := aMap[e.Ions[i].IonForm]
// 	if ok {
// 		for _, j := range v {
// 			e.Ions[i].AssignedModifications[j] = 0
// 		}
// 	}
// }

// // forward it to the peptide data
// for i := range e.Peptides {
// 	v, ok := aPepMap[e.Peptides[i].Sequence]
// 	if ok {
// 		for _, j := range v {
// 			e.Peptides[i].AssignedModifications[j] = 0
// 		}
// 	}
// }

//return
//}

// UpdateSupportingSpectra pushes back from SM to Protein the new supporting spectra from razor results
func (e *Evidence) UpdateSupportingSpectra() {

	var ptSupSpec = make(map[string][]string)
	var uniqueSpec = make(map[string][]string)
	var razorSpec = make(map[string][]string)

	for _, i := range e.PSM {

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

	for i := range e.Proteins {
		for j := range e.Proteins[i].TotalPeptideIons {

			for k := range e.Proteins[i].TotalPeptideIons[j].Spectra {
				delete(e.Proteins[i].TotalPeptideIons[k].Spectra, k)
			}

		}
	}

	for i := range e.Proteins {

		v, ok := ptSupSpec[e.Proteins[i].PartHeader]
		if ok {
			for _, j := range v {
				e.Proteins[i].SupportingSpectra[j] = 0
			}
		}

		for k := range e.Proteins[i].TotalPeptideIons {

			Up, UOK := uniqueSpec[e.Proteins[i].TotalPeptideIons[k].IonForm]
			if UOK && e.Proteins[i].TotalPeptideIons[k].IsUnique == true {
				for _, l := range Up {
					e.Proteins[i].TotalPeptideIons[k].Spectra[l] = 0
				}
			}

			Rp, ROK := razorSpec[e.Proteins[i].TotalPeptideIons[k].IonForm]
			if ROK && e.Proteins[i].TotalPeptideIons[k].IsURazor == true {
				for _, l := range Rp {
					e.Proteins[i].TotalPeptideIons[k].Spectra[l] = 0
				}
			}

		}

	}

	return
}

// UpdatePeptideModCount counts how many times each peptide is observed modified and not modified
func (e *Evidence) UpdatePeptideModCount() {

	// recreate the ion list from the main report object
	var all = make(map[string]int)
	var mod = make(map[string]int)
	var unmod = make(map[string]int)

	for _, i := range e.Peptides {
		all[i.Sequence] = 0
		mod[i.Sequence] = 0
		unmod[i.Sequence] = 0
	}

	// range over PSMs looking for modified and not modified evidences
	// if they exist on the ions map, get the numbers
	for _, i := range e.PSM {

		_, ok := all[i.Peptide]
		if ok {

			if i.Massdiff >= -0.99 && i.Massdiff <= 0.99 {
				unmod[i.Peptide]++
			} else {
				mod[i.Peptide]++
			}

		}
	}

	for i := range e.Peptides {

		v1, ok1 := unmod[e.Peptides[i].Sequence]
		if ok1 {
			e.Peptides[i].UnModifiedObservations = v1
		}

		v2, ok2 := mod[e.Peptides[i].Sequence]
		if ok2 {
			e.Peptides[i].ModifiedObservations = v2
		}

	}

	return
}
