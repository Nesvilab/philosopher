package fil

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"

	"github.com/Nesvilab/philosopher/lib/cla"
	"github.com/Nesvilab/philosopher/lib/id"
	"github.com/Nesvilab/philosopher/lib/msg"
	"github.com/Nesvilab/philosopher/lib/raz"

	"github.com/sirupsen/logrus"
)

// PepXMLFDRFilter processes and calculates the FDR at the PSM, Ion or Peptide level
// targetProbOpt (optional): minimum Probability bound, used for sequential FDR filtering.
// If targetProbOpt is provided (e.g., 0.96), only entries with Probability >= targetProb are considered when determining the threshold.
func PepXMLFDRFilter(input map[string]id.PepIDListPtrs, targetFDR float64, level, decoyTag, debug string, targetProbOpt ...float64) (id.PepIDListPtrs, float64) {

	//var msg string
	var calcFDR float64
	var list id.PepIDListPtrs
	var minProb float64 = 10

	// move all entries from map to list for psm, peptide, and ion
	if strings.EqualFold(level, "PSM") {
		for _, i := range input {
			for _, j := range i {
				list = append(list, j)
			}
		}
	} else if strings.EqualFold(level, "Peptide") {
		// 0 index means the one with the highest score
		for _, i := range input {
			list = append(list, i[0])
		}
	} else if strings.EqualFold(level, "Ion") {
		// 0 index means the one with the highest score
		for _, i := range input {
			list = append(list, i[0])
		}
	}

	// Optional minimum probability bound (sequential behavior).
	minProbBound := math.Inf(-1)
	if len(targetProbOpt) > 0 {
		minProbBound = targetProbOpt[0]
	}

	if len(list) > 0 {
		// compute pepID qvalues using Probability and add them to the list, the resulting list is sorted by Probability in descending order
		computePepIDQvalue(list, level, decoyTag)

		// determine the minimal Probability that satisfying the FDR threshold
		limit := (len(list) - 1)

		// retrieve extract q-value based on level and determine minProb
		getQvalue := func(item *id.PeptideIdentification, level string) float64 {
			switch level {
			case "PSM":
				return item.Qvalue
			case "Peptide":
				return item.PeptideQvalue
			case "Ion":
				return item.IonQvalue
			default:
				return math.NaN() // fallback if unknown level
			}
		}

		for j := 0; j < limit; j++ {
			if list[j].Probability < minProbBound {
				break // early break as list is sorted by Probability descending order
			}

			qval := getQvalue(list[j], level)
			if qval <= targetFDR {
				minProb = list[j].Probability
				calcFDR = qval
			}
		}

	}

	// retain only qualified ones to list
	cleanlist := make(id.PepIDListPtrs, 0)
	decoys := 0
	targets := 0

	for i := range list {
		if list[i].Probability >= minProb {

			cleanlist = append(cleanlist, list[i])

			if cla.IsDecoyPSM(*list[i], decoyTag) {
				decoys++
			} else {
				targets++
			}
		}
	}

	// print basic info
	if minProbBound < 0 {
		msg := fmt.Sprintf("Converged to %.2f %% FDR with %d %ss", calcFDR*100, targets, level)
		logrus.WithFields(logrus.Fields{
			"decoy":     decoys,
			"total":     (targets + decoys),
			"threshold": minProb,
		}).Info(msg)
	} else {
		logrus.WithFields(logrus.Fields{
			"decoy":     decoys,
			"total":     (targets + decoys),
			"threshold": minProb,
		}).Info(fmt.Sprintf("%d %ss", targets, level))
	}

	return cleanlist, minProb
}

func computePepIDQvalue(list id.PepIDListPtrs, level, decoyTag string) {

	var targets uint
	var decoys uint

	// Sort by Probability descending. A deterministic secondary key on ties
	// guarantees reproducible ordering: Go's sort.Slice is not stable and the
	// input list is built from map iteration, so without a tie-breaker the
	// within-tier order (and therefore the FDR bookkeeping below) can differ
	// between runs on identical input.
	sort.Slice(list, func(i, j int) bool {
		if list[i].Probability != list[j].Probability {
			return list[i].Probability > list[j].Probability
		}
		si := list[i].SpectrumFileName().Str()
		sj := list[j].SpectrumFileName().Str()
		if si != sj {
			return si < sj
		}
		if list[i].Peptide != list[j].Peptide {
			return list[i].Peptide < list[j].Peptide
		}
		if list[i].AssumedCharge != list[j].AssumedCharge {
			return list[i].AssumedCharge < list[j].AssumedCharge
		}
		return list[i].Protein < list[j].Protein
	})

	// Build the Probability -> FDR map using cumulative counts at the END of
	// each probability tier (i.e. counts over ALL entries with Probability
	// >= p). Recording the value at the last position of the tier, rather
	// than the first, makes the FDR invariant to the relative order of tied
	// entries and matches the standard "FDR at threshold p" definition.
	var probFDRMap = make(map[float64]float64)
	limit := len(list)

	for i := 0; i < limit; i++ {
		pepID := list[i]

		if cla.IsDecoyPSM(*pepID, decoyTag) {
			decoys++
		} else {
			targets++
		}

		if i == limit-1 || list[i+1].Probability != pepID.Probability {
			if targets > 0 {
				probFDRMap[pepID.Probability] = float64(decoys) / float64(targets)
			} else {
				probFDRMap[pepID.Probability] = 1
			}
		}
	}

	// iterative over Probability from lowest to highest, assigning each Qvalue as the smallest FDR
	minFDR := probFDRMap[list[limit-1].Probability]

	assignQvalue := func(i int, val float64) {
		switch strings.ToLower(level) {
		case "psm":
			list[i].Qvalue = val
		case "peptide":
			list[i].PeptideQvalue = val
		case "ion":
			list[i].IonQvalue = val
		}
	}

	assignQvalue(limit-1, minFDR)

	for i := limit - 2; i >= 0; i-- {
		fdr := probFDRMap[list[i].Probability]
		if fdr < minFDR {
			minFDR = fdr
		}
		assignQvalue(i, minFDR)
	}

}

// PickedFDR employs the picked FDR strategy
func PickedFDR(p id.ProtXML) id.ProtXML {

	// var appMap = make(map[string]int)
	var targetMap = make(map[string]float64)
	var decoyMap = make(map[string]float64)
	var recordMap = make(map[string]int)

	// collect all proteins from every group
	for _, i := range p.Groups {
		for _, j := range i.Proteins {
			if cla.IsDecoyProtein(j, p.DecoyTag) {
				decoyMap[string(j.ProteinName)] = j.TopPepProb
			} else {
				targetMap[string(j.ProteinName)] = j.TopPepProb
			}
		}
	}

	// check unique targets
	for k := range targetMap {
		iKey := fmt.Sprintf("%s%s", p.DecoyTag, k)
		_, ok := decoyMap[iKey]
		if !ok {
			recordMap[k] = 1
		}
	}

	// check unique decoys
	for k := range decoyMap {
		iKey := strings.Replace(k, p.DecoyTag, "", -1)
		_, ok := targetMap[iKey]
		if !ok {
			recordMap[k] = 1
		}
	}

	// check paired observations
	for k, v := range targetMap {

		iKey := fmt.Sprintf("%s%s", p.DecoyTag, k)
		vok, ok := decoyMap[iKey]
		if ok {
			if vok > v {
				recordMap[k] = 0
				recordMap[iKey] = 1
			} else if v > vok {
				recordMap[k] = 1
				recordMap[iKey] = 0
			} else {
				recordMap[k] = 1
				recordMap[iKey] = 1
			}
		}
	}

	// collect all proteins from every group
	for i := range p.Groups {
		for j := range p.Groups[i].Proteins {

			v, ok := recordMap[string(p.Groups[i].Proteins[j].ProteinName)]
			if ok {
				p.Groups[i].Proteins[j].Picked = v
			}
		}
	}

	return p
}

// RazorCandidateMap is a list of razor candidates
//type RazorCandidateMap map[string]RazorCandidate

// RazorFilter classifies peptides as razor
func RazorFilter(p id.ProtXML, minPepLen ...int) id.ProtXML {
	pepLen := 7
	if len(minPepLen) > 0 {
		pepLen = minPepLen[0]
	}

	var r raz.RazorMap = make(map[string]raz.RazorCandidate)
	var rList []string

	// perform a test load of the razor assingment, if there's a file, then the assignment is skipped, and the current file used.
	r.Restore(true)
	if len(r) == 0 {

		// for each peptide sequence, collapse all parent protein peptides from ions originated from the same sequence
		for _, i := range p.Groups {
			for _, j := range i.Proteins {
				for _, k := range j.PeptideIons {

					v, ok := r[k.PeptideSequence]
					if !ok {

						var rc raz.RazorCandidate
						rc.Sequence = k.PeptideSequence
						rc.MappedProteins = make(map[string]raz.MappedProtein)

						rc.MappedProteins[j.ProteinName] = raz.MappedProtein{
							W:   k.Weight,
							GW:  k.GroupWeight,
							TNP: j.TotalNumberPeptides,
							SID: j.GroupSiblingID,
						}

						for _, i := range j.IndistinguishableProtein {
							rc.MappedProteins[i] = raz.MappedProtein{
								W:   -1,
								GW:  -1,
								TNP: -1,
								SID: "zzz",
							}
						}

						for _, i := range k.PeptideParentProtein {
							rc.MappedProteins[i] = raz.MappedProtein{
								W:   -1,
								GW:  -1,
								TNP: -1,
								SID: "zzz",
							}
						}

						r[k.PeptideSequence] = rc

					} else {
						var c = v

						// doing like this will allow proteins that map to shared peptidesto be considered
						c.MappedProteins[j.ProteinName] = raz.MappedProtein{
							W:   k.Weight,
							GW:  k.GroupWeight,
							TNP: j.TotalNumberPeptides,
							SID: j.GroupSiblingID,
						}
						r[k.PeptideSequence] = c

					}

				}
			}
		}

		// this will make the assignment more deterministic
		for k := range r {
			rList = append(rList, k)
		}
		sort.Strings(rList)

		var razorPair = make(map[string]string)

		// get the best protein candidate for each peptide sequence and make the razor pair
		for _, k := range rList {
			// 1st pass: mark all cases with weight > 0.5
			for pt, w := range r[k].MappedProteins {
				if w.W > 0.5 {
					razorPair[k] = pt
				}
			}
		}

		// 2nd pass: mark all cases with highest group weight in the list
		for _, k := range rList {

			_, ok := razorPair[k]
			if !ok {

				var topPT string
				var topCount int
				var topGW float64
				var topTNP int
				var topGWMap = make(map[float64]uint8)
				var topTNPMap = make(map[int]uint8)

				if len(r[k].MappedProteins) == 1 {

					for pt := range r[k].MappedProteins {
						razorPair[k] = pt
					}

				} else if len(r[k].MappedProteins) > 1 {

					for pt, tnp := range r[k].MappedProteins {
						if tnp.GW >= topGW {
							topGW = tnp.GW
							topPT = pt
							topGWMap[topGW]++
						}
					}

					var tie bool
					if topGWMap[topGW] >= 2 {
						tie = true
					}

					if !tie {
						razorPair[k] = topPT

					} else {

						var tnpList []string
						for pt := range r[k].MappedProteins {
							tnpList = append(tnpList, pt)
						}

						sort.Strings(tnpList)

						for _, pt := range tnpList {
							if r[k].MappedProteins[pt].TNP >= topTNP {
								topTNP = r[k].MappedProteins[pt].TNP
								topPT = pt
								topTNPMap[topTNP]++
							}
						}

						var tie bool
						if topTNPMap[topTNP] >= 2 {
							tie = true
						}

						if !tie {

							var mplist []string
							for pt := range r[k].MappedProteins {
								mplist = append(mplist, pt)
							}
							sort.Strings(mplist)

							for _, pt := range mplist {
								if r[k].MappedProteins[pt].TNP >= topCount {
									topCount = r[k].MappedProteins[pt].TNP
									topPT = pt
								}
							}

							razorPair[k] = topPT

						} else {

							var idList []string
							for protein, id := range r[k].MappedProteins {
								idList = append(idList, fmt.Sprintf("%s#%s", id.SID, protein))
							}

							sort.Strings(idList)

							id := strings.Split(idList[0], "#")
							razorPair[k] = id[1]
						}

					}
				}
			}
		}

		for _, k := range rList {
			pt, ok := razorPair[k]
			if ok {
				razor := r[k]
				razor.MappedProtein = pt
				r[k] = razor
			}
		}

		r.Serialize()
	}

	for i := range p.Groups {
		for j := range p.Groups[i].Proteins {
			for k := range p.Groups[i].Proteins[j].PeptideIons {
				v, ok := r[string(p.Groups[i].Proteins[j].PeptideIons[k].PeptideSequence)]
				if ok {
					if p.Groups[i].Proteins[j].ProteinName == v.MappedProtein {
						p.Groups[i].Proteins[j].PeptideIons[k].Razor = 1
						p.Groups[i].Proteins[j].HasRazor = true
					}
				}
			}
		}
	}

	// mark as razor all peptides in the reference map
	for i := range p.Groups {
		for j := range p.Groups[i].Proteins {
			var r float64
			for k := range p.Groups[i].Proteins[j].PeptideIons {

				if p.Groups[i].Proteins[j].PeptideIons[k].Razor == 1 || p.Groups[i].Proteins[j].PeptideIons[k].IsUnique {
					if (p.Groups[i].Proteins[j].PeptideIons[k].InitialProbability > r) && (p.Groups[i].Proteins[j].PeptideIons[k].PeptideLength >= pepLen) {
						r = p.Groups[i].Proteins[j].PeptideIons[k].InitialProbability
					}
				}
			}
			p.Groups[i].Proteins[j].TopPepProb = r
		}
	}

	return p
}

// ProtXMLFilter filters the protein list under a specific fdr
func ProtXMLFilter(p id.ProtXML, targetFDR, pepProb, protProb float64, isPicked, isRazor bool, decoyTag string) id.ProtIDList {

	//var proteinIDs ProtIDList
	var list id.ProtIDList
	var targets float64
	var decoys float64
	var calcFDR float64
	var minProb float64 = 10

	// collect all proteins from every group into a list
	for i := range p.Groups {
		for j := range p.Groups[i].Proteins {

			if isRazor {

				if isPicked {
					if p.Groups[i].Proteins[j].Picked == 1 && p.Groups[i].Proteins[j].HasRazor {
						list = append(list, p.Groups[i].Proteins[j])
					}
				} else {
					if p.Groups[i].Proteins[j].HasRazor {
						list = append(list, p.Groups[i].Proteins[j])
					}
				}

			} else {

				if isPicked {
					if p.Groups[i].Proteins[j].TopPepProb >= protProb && p.Groups[i].Proteins[j].Picked == 1 {
						list = append(list, p.Groups[i].Proteins[j])
					}

				} else {
					if p.Groups[i].Proteins[j].TopPepProb >= pepProb && p.Groups[i].Proteins[j].TopPepProb >= protProb {
						list = append(list, p.Groups[i].Proteins[j])
					}
				}
			}
		}
	}

	// compute top protein qvalues using TopPepProb and add them to the list
	FDRMap := computeProteinQvalue(list, p.DecoyTag)

	// determine the minimal TopPepProb corresponding to the target FDR
	var topPepProb []float64
	for k := range FDRMap {
		topPepProb = append(topPepProb, k)
	}

	sort.Sort(sort.Reverse(sort.Float64Slice(topPepProb)))

	var currentTopPepProb = 10.0
	var currentFDR = 0.0
	var probArray []float64
	var probIndex int
	var probList = make(map[float64]uint8)

	for i := range topPepProb {

		// for inspections
		// fmt.Println(topPepProb[i], "\t", FDRMap[topPepProb[i]], "\t", uti.ToFixed(FDRMap[topPepProb[i]], 4), "\t", uti.Round(FDRMap[topPepProb[i]]*100, 5, 2))

		probArray = append(probArray, topPepProb[i])

		if FDRMap[topPepProb[i]] <= targetFDR {

			probList[topPepProb[i]] = 0
			minProb = topPepProb[i]
			calcFDR = FDRMap[topPepProb[i]]

			if topPepProb[i] < currentTopPepProb {
				currentTopPepProb = topPepProb[i]
				probIndex = i
			}

			if FDRMap[topPepProb[i]] > currentFDR {
				currentFDR = FDRMap[topPepProb[i]]
			}
		}

	}

	if currentTopPepProb == 10 {
		msg.Custom(errors.New("the protein FDR filter didn't reach the desired threshold, try a higher threshold using the --prot parameter"), "fatal")
	}

	// if the current protein block threshold is below the threshold, we look for the next one
	if currentFDR < targetFDR && probArray[len(probArray)-1] != currentTopPepProb {

		// test if the next protein block surpasses the threshold by 10%
		thresh := currentFDR + (0.01 * currentFDR)

		if FDRMap[probArray[probIndex+1]] <= thresh {

			probList[probArray[probIndex+1]] = 0

			minProb = probArray[probIndex+1]
			calcFDR = FDRMap[probArray[probIndex+1]]

		}
	}

	// for inspections
	// fmt.Println("index ", probIndex, probArray[probIndex], "minProb ", minProb, "calcFDR ", calcFDR)

	// retain only the qualified protein identifications
	var finalList id.ProtIDList
	decoys = 0
	targets = 0

	for i := range list {
		if list[i].TopPepProb >= currentTopPepProb {

			finalList = append(finalList, list[i])

			if cla.IsDecoyProtein(list[i], p.DecoyTag) {
				decoys++
			} else {
				targets++
			}
		}
	}

	msg := fmt.Sprintf("Converged to %.2f %% FDR with %0.f Proteins", (calcFDR * 100), targets)
	logrus.WithFields(logrus.Fields{
		"decoy":     decoys,
		"total":     (targets + decoys),
		"threshold": minProb,
	}).Info(msg)

	return finalList
}

func computeProteinQvalue(list id.ProtIDList, decoyTag string) map[float64]float64 {

	var targets float64
	var decoys float64

	// Sort by TopPepProb descending with a deterministic secondary key on
	// ties (ProteinName). Without the tie-breaker, Go's non-stable sort and
	// upstream map iteration would randomize the order of equally-scored
	// entries, causing non-reproducible FDR values below.
	sort.Slice(list, func(i, j int) bool {
		if list[i].TopPepProb != list[j].TopPepProb {
			return list[i].TopPepProb > list[j].TopPepProb
		}
		return list[i].ProteinName < list[j].ProteinName
	})

	// Build the TopPepProb -> FDR map using cumulative counts at the END of
	// each probability tier (i.e. counts over ALL entries with TopPepProb
	// >= p). Recording the value at the last position of the tier, rather
	// than the first, makes the FDR invariant to the relative order of tied
	// entries.
	var probFDRMap = make(map[float64]float64)

	// compute FDR and put it to probFDRMap
	limit := len(list)
	for i := 0; i < limit; i++ {
		protID := list[i]

		if cla.IsDecoyProtein(list[i], decoyTag) {
			decoys++
		} else {
			targets++
		}

		if i == limit-1 || list[i+1].TopPepProb != protID.TopPepProb {
			if targets > 0 {
				probFDRMap[protID.TopPepProb] = decoys / targets
			} else {
				probFDRMap[protID.TopPepProb] = 1
			}
		}
	}

	// create a TopPepProb-to-Qvalue map, FDRMap (key: TopPepProb, value: Qvalue)
	minFDR := probFDRMap[list[limit-1].TopPepProb]
	list[limit-1].Qvalue = minFDR

	var FDRMap = make(map[float64]float64)
	FDRMap[list[limit-1].TopPepProb] = probFDRMap[list[limit-1].TopPepProb]

	for i := limit - 2; i >= 0; i-- {
		if probFDRMap[list[i].TopPepProb] < minFDR {
			minFDR = probFDRMap[list[i].TopPepProb]
		}
		list[i].Qvalue = minFDR

		_, ok := FDRMap[list[i].TopPepProb]
		if !ok {
			FDRMap[list[i].TopPepProb] = list[i].Qvalue
		}

	}

	return FDRMap
}

// sequentialFDRControl estimates FDR levels by applying a second filter where all
// proteins from the protein filtered list are matched against filtered PSMs
func sequentialFDRControl(pep id.PepIDList, pro id.ProtIDList, psm, psmT, peptide, peptideT, ion, ionT float64, decoyTag string) {

	extPep := extractPSMfromPepXML("sequential", pep, pro)

	// organize enties by score (probability or expectation)
	sort.Sort(extPep)

	uniqPsms := GetUniquePSMs(extPep)
	uniqPeps := GetUniquePeptides(extPep)
	uniqIons := getUniquePeptideIons(extPep)

	logrus.WithFields(logrus.Fields{
		"psms":     len(uniqPsms),
		"peptides": len(uniqPeps),
		"ions":     len(uniqIons),
	}).Info("Applying sequential FDR estimation")

	wg := sync.WaitGroup{}
	wg.Add(3)
	filteredPSM, _ := PepXMLFDRFilter(uniqPsms, psm, "PSM", decoyTag, "", 0)
	go func() { defer wg.Done(); filteredPSM.Serialize("psm") }()
	filteredPeptides, _ := PepXMLFDRFilter(uniqPeps, peptide, "Peptide", decoyTag, "", peptideT)
	go func() { defer wg.Done(); filteredPeptides.Serialize("pep") }()
	filteredIons, _ := PepXMLFDRFilter(uniqIons, ion, "Ion", decoyTag, "", ionT)
	go func() { defer wg.Done(); filteredIons.Serialize("ion") }()
	wg.Wait()

}

// twoDFDRFilter estimates FDR levels by applying a second filter by regenerating
// a protein list with decoys from protXML and pepXML.
func twoDFDRFilter(pep id.PepIDList, pro id.ProtIDList, psm, peptide, ion float64, decoyTag string) {

	// filter protein list at given FDR level and regenerate protein list by adding pairing decoys
	//logrus.Info("Creating mirror image from filtered protein list")
	mirrorProteinList := mirrorProteinList(pro, decoyTag)

	// get new protein list profile
	//logrus.Info(protxml.ProteinProfileWithList(mirrorProteinList, pa.Tag, pa.Con))
	t, d := proteinProfileWithList(mirrorProteinList, decoyTag)
	logrus.WithFields(logrus.Fields{
		"target": t,
		"decoy":  d,
	}).Info("2D FDR estimation: Protein mirror image")

	// get PSM from the original pepXML using protein REGENERATED protein list, using protein names
	extPep := extractPSMfromPepXML("2d", pep, mirrorProteinList)

	// organize enties by score (probability or expectation)
	sort.Sort(extPep)

	uniqPsms := GetUniquePSMs(extPep)
	uniqPeps := GetUniquePeptides(extPep)
	uniqIons := getUniquePeptideIons(extPep)

	logrus.WithFields(logrus.Fields{
		"psms":     len(uniqPsms),
		"peptides": len(uniqPeps),
		"ions":     len(uniqIons),
	}).Info("Second filtering results")

	filteredPSM, _ := PepXMLFDRFilter(uniqPsms, psm, "PSM", decoyTag, "")
	filteredPeptides, _ := PepXMLFDRFilter(uniqPeps, peptide, "Peptide", decoyTag, "")
	filteredIons, _ := PepXMLFDRFilter(uniqIons, ion, "Ion", decoyTag, "")

	filteredPSM.Serialize("psm")
	filteredPeptides.Serialize("pep")
	filteredIons.Serialize("ion")
}

// mirrorProteinList takes a filtered list and regenerate the correspondedn decoys
func mirrorProteinList(p id.ProtIDList, decoyTag string) id.ProtIDList {

	var targets = make(map[string]uint8)
	var decoys = make(map[string]uint8)

	// get filtered list
	var list id.ProtIDList
	for _, i := range p {
		if !cla.IsDecoyProtein(i, decoyTag) {
			list = append(list, i)
		}
	}

	// get the list of identified taget proteins
	for _, i := range p {
		if cla.IsDecoy(i.ProteinName, decoyTag) {
			decoys[i.ProteinName] = 0
		} else {
			targets[i.ProteinName] = 0
		}
	}

	// collect all original protein ids in case we need to put them on mirror list
	var refMap = make(map[string]id.ProteinIdentification)
	for _, i := range p {
		refMap[i.ProteinName] = i
	}

	// add decoys correspondent to the given targets.
	// first check if the opposite list doesn't have an entry already.
	// if not, search for the mirror entry on the original list, if found
	// move it to the mirror list, otherwise add fake entry.
	for _, k := range list {
		decoy := decoyTag + k.ProteinName
		v, ok := refMap[decoy]
		if ok {
			list = append(list, v)
		} else {
			var pt id.ProteinIdentification
			pt.ProteinName = decoy
			list = append(list, pt)
		}
	}

	return list
}
