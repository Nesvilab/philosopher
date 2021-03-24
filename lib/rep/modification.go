package rep

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"

	"philosopher/lib/msg"
	"philosopher/lib/obo"

	"philosopher/lib/sys"
	"philosopher/lib/uti"
)

// int[] matched_unimod_idx = new int[delta_mass_from_psm.length];
// for (int i = 0; i < delta_mass_from_psm.length; ++i) {
// 	double gap = 9999999;
// 	double matched_mass;
// 	int matched_idx;
//   for (int j = 0; j < unimod.length; ++j) { // can be optimized with binary search
//     if (Math.abs(delta_mass_from_psm[i] - unimod[j]) < gap) {
//     	gap = Math.abs(delta_mass_from_psm[i] - unimod[j])
//     	matched_mass = unimod[j];
// 		matched_idx = j;
//     }
//   }
// ​
//   if (gap < ppm * 1e-6 * matched_mass) {
//   	matched_unimod_idx[i] = matched_idx;
//   } else {
//   	matched_unimod_idx[i] = null;
//   }
// }

// MapMods maps PSMs to modifications based on their mass shifts
func (evi *Evidence) MapMods() {

	var modMap = make(map[float64]obo.Term)
	var modList []float64
	var ppm = float64(20)

	o := obo.NewUniModOntology()

	for _, i := range evi.PSM {
		modMap[i.Massdiff] = obo.Term{}
	}

	for k := range modMap {
		modList = append(modList, k)
	}

	for i := 0; i <= len(modList)-1; i++ {

		var gap = float64(9999999)
		var obo obo.Term
		var mass float64

		for j := range o.Terms {

			if math.Abs(modList[i]-o.Terms[j].MonoIsotopicMass) < gap {
				gap = math.Abs(modList[i] - o.Terms[j].MonoIsotopicMass)
				obo = o.Terms[j]
				mass = modList[i]
			}
		}

		if gap < (1e-6 * ppm * mass) {
			modMap[mass] = obo
		} else {
			delete(modMap, mass)
		}
	}

	for i := range evi.PSM {
		// for fixed and variable modifications
		for k, v := range evi.PSM[i].Modifications.Index {

			obo, ok := modMap[v.MassDiff]
			if ok {
				updatedMod := v

				_, ok := obo.Sites[v.AminoAcid]
				if ok {
					updatedMod.Name = obo.Name
					updatedMod.Definition = obo.Definition
					updatedMod.ID = obo.ID
					updatedMod.MonoIsotopicMass = obo.MonoIsotopicMass
					updatedMod.IsobaricMods[obo.Name]++
					evi.PSM[i].Modifications.Index[k] = updatedMod
				}
				if updatedMod.Type == "Observed" {
					updatedMod.Name = obo.Name
					updatedMod.Definition = obo.Definition
					updatedMod.ID = obo.ID
					updatedMod.MonoIsotopicMass = obo.MonoIsotopicMass
					updatedMod.IsobaricMods[obo.Name] = obo.MonoIsotopicMass
					evi.PSM[i].Modifications.Index[k] = updatedMod
				}
				break
			}
		}
	}
}

// AssembleModificationReport cretaes the modifications lists
func (evi *Evidence) AssembleModificationReport() {

	var modEvi ModificationEvidence

	var massWindow = float64(0.5)
	var binsize = float64(0.1)
	var amplitude = float64(1000)

	var bins []MassBin

	nBins := (amplitude*(1/binsize) + 1) * 2
	for i := 0; i <= int(nBins); i++ {
		var b MassBin

		b.LowerMass = -(amplitude) - (massWindow * binsize) + (float64(i) * binsize)
		b.LowerMass = uti.Round(b.LowerMass, 5, 4)

		b.HigherRight = -(amplitude) + (massWindow * binsize) + (float64(i) * binsize)
		b.HigherRight = uti.Round(b.HigherRight, 5, 4)

		b.MassCenter = -(amplitude) + (float64(i) * binsize)
		b.MassCenter = uti.Round(b.MassCenter, 5, 4)

		bins = append(bins, b)
	}

	// calculate the total number of PSMs per cluster
	for i := range evi.PSM {

		// the checklist will not allow the same PSM to be added multiple times to the
		// same bin in case multiple identical mods are present in te sequence
		var assignChecklist = make(map[float64]uint8)
		var obsChecklist = make(map[float64]uint8)

		for j := range bins {

			// for assigned mods
			// 0 here means something that doest not map to the pepXML header
			// like multiple mods on n-term
			for _, l := range evi.PSM[i].Modifications.Index {

				if l.MassDiff > bins[j].LowerMass && l.MassDiff <= bins[j].HigherRight && l.MassDiff != 0 {
					_, ok := assignChecklist[l.MassDiff]
					if !ok {
						if l.Type == "Assigned" {
							bins[j].AssignedMods = append(bins[j].AssignedMods, evi.PSM[i])
							assignChecklist[l.MassDiff] = 0
						}
					}
				}
			}

			// for delta masses
			if evi.PSM[i].Massdiff > bins[j].LowerMass && evi.PSM[i].Massdiff <= bins[j].HigherRight {
				_, ok := obsChecklist[evi.PSM[i].Massdiff]
				if !ok {
					bins[j].ObservedMods = append(bins[j].ObservedMods, evi.PSM[i])
					obsChecklist[evi.PSM[i].Massdiff] = 0
				}
			}

		}
	}

	// calculate average mass for each cluster
	var zeroBinMassDeviation float64
	for i := range bins {
		pep := bins[i].ObservedMods
		total := 0.0
		for j := range pep {
			total += pep[j].Massdiff
		}
		if len(bins[i].ObservedMods) > 0 {
			bins[i].AverageMass = (float64(total) / float64(len(pep)))
		} else {
			bins[i].AverageMass = 0
		}
		if bins[i].MassCenter == 0 {
			zeroBinMassDeviation = bins[i].AverageMass
		}

		bins[i].AverageMass = uti.Round(bins[i].AverageMass, 5, 4)
	}

	// correcting mass values based on Bin 0 average mass
	for i := range bins {
		if len(bins[i].ObservedMods) > 0 {
			if bins[i].AverageMass > 0 {
				bins[i].CorrectedMass = (bins[i].AverageMass - zeroBinMassDeviation)
			} else {
				bins[i].CorrectedMass = (bins[i].AverageMass + zeroBinMassDeviation)
			}
		} else {
			bins[i].CorrectedMass = bins[i].MassCenter
		}
		bins[i].CorrectedMass = uti.Round(bins[i].CorrectedMass, 5, 4)
	}

	modEvi.MassBins = bins
	evi.Modifications = modEvi

}

// ModificationReport ...
func (evi *Evidence) ModificationReport() {

	// create result file
	output := fmt.Sprintf("%s%smodifications.tsv", sys.MetaDir(), string(filepath.Separator))

	// create result file
	file, e := os.Create(output)
	if e != nil {
		msg.WriteFile(errors.New("could not create report files"), "error")
	}
	defer file.Close()

	line := "Mass Bin\tPSMs with Assigned Modifications\tPSMs with Observed Modifications\n"

	_, e = io.WriteString(file, line)
	if e != nil {
		msg.WriteToFile(e, "fatal")
	}

	for _, i := range evi.Modifications.MassBins {

		line = fmt.Sprintf("%.4f\t%d\t%d",
			i.CorrectedMass,
			len(i.AssignedMods),
			len(i.ObservedMods),
		)

		line += "\n"
		_, e = io.WriteString(file, line)
		if e != nil {
			msg.WriteToFile(e, "fatal")
		}

	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))
}

// PlotMassHist plots the delta mass histogram
func (evi *Evidence) PlotMassHist() {

	outfile := fmt.Sprintf("%s%sdelta-mass.html", sys.MetaDir(), string(filepath.Separator))

	file, e := os.Create(outfile)
	if e != nil {
		msg.WriteFile(errors.New("could not create output for delta mass binning"), "fatal")
	}
	defer file.Close()

	var xvar []string
	var y1var []string
	var y2var []string

	for _, i := range evi.Modifications.MassBins {
		if i.MassCenter >= -501 && i.MassCenter <= 501 {
			xel := fmt.Sprintf("'%.2f',", i.MassCenter)
			xvar = append(xvar, xel)
			y1el := fmt.Sprintf("'%d',", len(i.AssignedMods))
			y1var = append(y1var, y1el)
			y2el := fmt.Sprintf("'%d',", len(i.ObservedMods))
			y2var = append(y2var, y2el)
		}
	}

	xAxis := fmt.Sprintf("	  x: %s,", xvar)
	AssAxis := fmt.Sprintf("	  y: %s,", y1var)
	ObsAxis := fmt.Sprintf("	  y: %s,", y2var)

	io.WriteString(file, "<head>\n")
	io.WriteString(file, "  <script src=\"https://cdn.plot.ly/plotly-latest.min.js\"></script>\n")
	io.WriteString(file, "</head>\n")
	io.WriteString(file, "<body>\n")
	io.WriteString(file, "<div id=\"myDiv\" style=\"width: 1024px; height: 768px;\"></div>\n")
	io.WriteString(file, "<script>\n")
	io.WriteString(file, "var trace1 = {")
	io.WriteString(file, xAxis)
	io.WriteString(file, ObsAxis)
	io.WriteString(file, "name: 'Observed',")
	io.WriteString(file, "type: 'bar',")
	io.WriteString(file, "};")
	io.WriteString(file, "var trace2 = {")
	io.WriteString(file, xAxis)
	io.WriteString(file, AssAxis)
	io.WriteString(file, "name: 'Assigned',")
	io.WriteString(file, "type: 'bar',")
	io.WriteString(file, "};")
	io.WriteString(file, "var data = [trace1, trace2];\n")
	io.WriteString(file, "var layout = {barmode: 'stack', title: 'Distribution of Mass Modifications', xaxis: {title: 'mass bins'}, yaxis: {title: '# PSMs'}};\n")
	io.WriteString(file, "Plotly.newPlot('myDiv', data, layout);\n")
	io.WriteString(file, "</script>\n")
	io.WriteString(file, "</body>")

	if e != nil {
		msg.Custom(errors.New("there was an error trying to plot the mass distribution"), "fatal")
	}

	// copy to work directory
	sys.CopyFile(outfile, filepath.Base(outfile))

}
