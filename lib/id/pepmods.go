package id

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/Nesvilab/philosopher/lib/bio"
	"github.com/Nesvilab/philosopher/lib/mod"
	"github.com/Nesvilab/philosopher/lib/spc"
	"github.com/sirupsen/logrus"
)

// Monoisotopic masses of what an unmodified terminus carries: the N-terminal hydrogen and the
// C-terminal hydroxyl. MSFragger's mod_nterm_mass and mod_cterm_mass include them.
const (
	nTermHydrogenMass = 1.00782503207
	cTermHydroxylMass = 17.00273965
)

// declaredMassTolerance is how far, in Da, a pepXML mass may sit from the mass a search summary line
// states and still be read as that modification. The summary states masses with four decimals, spectrum
// queries carry five, and PTMProphet rewrites both from its own constants; those disagreements stay
// below 0.001 Da. Distinct modifications on one residue can be as close as 0.0095 Da (phospho and
// sulfo), so the tolerance has to stay well below that.
const declaredMassTolerance = 0.002

// minModificationDelta is the smallest delta, in Da, that can be a modification or a mass offset. A
// derived delta below it is a mass that drifted from the unmodified residue or from a declared
// modification, not a new modification.
const minModificationDelta = 0.5

// tmtProNTermMass is the mod_nterm_mass of a TMTpro-labelled peptide: hydrogen plus 304.2071.
// legacyTMTProNTermWindow is how far from it the 2020 PTMProphet workaround still applies; the original
// rule accepted any N-terminal mass whose text contained "305".
const (
	tmtProNTermMass         = 305.2150
	legacyTMTProNTermWindow = 0.5
)

// declaredEntry is a search summary modification with the total mass its header line states.
type declaredEntry struct {
	mass float64
	m    mod.Modification
}

// declaredModifications is the search summary table indexed for the lookups made per PSM: the entries
// of each site sorted by the mass the header states, the summed fixed delta of each site, and every
// declared delta. Sites compare case-insensitively ("N-term", "n-term").
type declaredModifications struct {
	index  map[string]mod.Modification
	bySite map[string][]declaredEntry
	fixed  map[string]float64
	deltas []float64
}

// newDeclaredModifications indexes a search summary table once per pepXML file.
func newDeclaredModifications(mods mod.Modifications) declaredModifications {

	d := declaredModifications{
		index:  mods.Index,
		bySite: make(map[string][]declaredEntry),
		fixed:  make(map[string]float64),
	}

	for _, v := range mods.Index {
		site := siteKey(v.AminoAcid)
		if mass, ok := headerMass(v.Index); ok {
			d.bySite[site] = append(d.bySite[site], declaredEntry{mass: mass, m: v})
		} else {
			logrus.Debugf("modification %q states no mass in its key and is matched by exact key only", v.Index)
		}
		if !v.Variable {
			d.fixed[site] += v.MassDiff
		}
		d.deltas = append(d.deltas, v.MassDiff)
	}

	for _, entries := range d.bySite {
		sort.Slice(entries, func(i, j int) bool {
			if entries[i].mass != entries[j].mass {
				return entries[i].mass < entries[j].mass
			}
			return entries[i].m.Index < entries[j].m.Index
		})
	}
	sort.Float64s(d.deltas)

	return d
}

// siteKey normalises a residue or terminus name for lookups.
func siteKey(site string) string {
	return strings.ToUpper(site)
}

// headerMass is the total mass stated in a search summary key such as "C#160.0307" or "N-term#43.0184".
func headerMass(key string) (float64, bool) {
	i := strings.LastIndex(key, "#")
	if i < 0 {
		return 0, false
	}
	mass, err := strconv.ParseFloat(key[i+1:], 64)
	return mass, err == nil
}

// mapModsFromPepXML receives a pepXML struct with modifications and adds them to the given struct
func (p *PeptideIdentification) mapModsFromPepXML(m spc.ModificationInfo, declared declaredModifications) {

	p.ModifiedPeptide = string(m.ModifiedPeptide)
	index := make(map[string]mod.Modification)
	residues := strings.Split(p.Peptide, "")

	for _, entry := range m.ModAminoacidMass {
		if v, ok := declared.residueModification(residues, entry, m.ModAminoacidMass); ok {
			index[v.Index] = v
		}
	}

	if m.ModNTermMass != 0 {
		if v, ok := declared.terminalModification("N-term", m.ModNTermMass, nTermHydrogenMass); ok {
			index[v.Index] = v
		}
	}

	if m.ModCTermMass != 0 {
		if v, ok := declared.terminalModification("C-term", m.ModCTermMass, cTermHydroxylMass); ok {
			index[v.Index] = v
		}
	}

	// if isotopicCorr >= 0.036386 || isotopicCorr <= -0.036386 {
	addObservedMassDiff(index, p.Massdiff)
	if len(index) != 0 {
		p.Modifications = mod.Modifications{Index: index}.ToSlice()
	}
}

// residueModification resolves one mod_aminoacid_mass entry: the declared modification its mass stands
// for or, for a mass the search summary does not declare (a mass offset MSFragger localized on the
// residue), a modification derived from the mass. An entry that fits neither is dropped with a log.
func (d declaredModifications) residueModification(residues []string, entry spc.ModAminoacidMass, entries []spc.ModAminoacidMass) (mod.Modification, bool) {

	if entry.Position < 1 || entry.Position > len(residues) {
		logrus.Debugf("dropping modification mass %.4f at position %d of %s: outside the peptide", entry.Mass, entry.Position, strings.Join(residues, ""))
		return mod.Modification{}, false
	}
	residue := residues[entry.Position-1]

	v, ok := d.findDeclared(residue, entry.Mass)
	if !ok {
		v, ok = d.deriveResidueModification(residue, entry.Mass, entry.Position, entries)
	}
	if !ok {
		logrus.Debugf("dropping modification mass %.4f on %s%d of %s: not declared and not derivable", entry.Mass, residue, entry.Position, strings.Join(residues, ""))
		return mod.Modification{}, false
	}

	v.Index = fmt.Sprintf("%s#%d#%.4f", residue, entry.Position, entry.Mass)
	v.Position = entry.Position
	return v, true
}

// terminalModification resolves a mod_nterm_mass or mod_cterm_mass: the declared terminal modification
// its mass stands for, the TMTpro label when PTMProphet moved the mass further than the tolerance, or a
// modification derived from the mass. A terminus carrying only its own group yields nothing. A declared
// match keeps the declared index, so the three paths can never add two entries for one terminus.
func (d declaredModifications) terminalModification(terminus string, mass float64, groupMass float64) (mod.Modification, bool) {

	v, ok := d.findDeclared(terminus, mass)
	if !ok && terminus == "N-term" {
		v, ok = d.legacyTMTProNTerm(mass)
	}
	if !ok {
		if v, ok = d.derivedModification(terminus, mass, mass-groupMass-d.fixedMassDiff(terminus)); !ok {
			logrus.Debugf("dropping %s modification mass %.4f: only the terminal group", terminus, mass)
			return mod.Modification{}, false
		}
		if v.Index == "" {
			v.Index = fmt.Sprintf("%s#%.4f", terminus, mass)
		}
	}

	v.AminoAcid = terminus
	return v, true
}

// legacyTMTProNTerm keeps the 2020 workaround for PTMProphet output whose mod_nterm_mass drifted away
// from the declared TMTpro mass: an N-terminal mass near 305.2150 maps to the declared TMTpro entry.
func (d declaredModifications) legacyTMTProNTerm(mass float64) (mod.Modification, bool) {
	if math.Abs(mass-tmtProNTermMass) > legacyTMTProNTermWindow {
		return mod.Modification{}, false
	}
	return d.findDeclared("N-term", tmtProNTermMass)
}

// findDeclared returns the search summary modification on a site ("C", "N-term") that a pepXML mass
// stands for: the exact four-decimal key, else the entry whose stated mass is nearest within
// declaredMassTolerance, which covers masses that PTMProphet or another engine rounded differently.
func (d declaredModifications) findDeclared(site string, mass float64) (mod.Modification, bool) {
	if v, ok := d.index[fmt.Sprintf("%s#%.4f", site, mass)]; ok {
		return v, true
	}
	return d.nearest(site, mass, declaredMassTolerance)
}

// nearest returns the declared entry on a site whose stated mass is closest to a pepXML mass and
// closer than the bound. Ties go to the smaller key so the result does not depend on table order.
func (d declaredModifications) nearest(site string, mass float64, bound float64) (mod.Modification, bool) {

	entries := d.bySite[siteKey(site)]
	i := sort.Search(len(entries), func(k int) bool { return entries[k].mass >= mass })

	var best mod.Modification
	found := false
	gap := bound
	for _, k := range []int{i - 1, i} {
		if k < 0 || k >= len(entries) {
			continue
		}
		g := math.Abs(entries[k].mass - mass)
		if g < gap || (g == gap && found && entries[k].m.Index < best.Index) {
			best, gap, found = entries[k].m, g, true
		}
	}
	return best, found
}

// fixedMassDiff is the total delta of the fixed modifications declared for a residue or terminus.
func (d declaredModifications) fixedMassDiff(site string) float64 {
	return d.fixed[siteKey(site)]
}

// deriveResidueModification builds the modification for a mod_aminoacid_mass entry the search summary
// does not declare. Its mass is residue + fixed modifications + delta. When the same position also
// carries a declared entry, e.g. the fixed carbamidomethyl under a mass offset, that entry's mass is
// the more precise base; otherwise the base is the residue plus its declared fixed modifications.
func (d declaredModifications) deriveResidueModification(residue string, mass float64, position int, entries []spc.ModAminoacidMass) (mod.Modification, bool) {

	base, ok := d.declaredMassAtPosition(residue, position, mass, entries)
	if !ok {
		residueMass, known := bio.MonoIsotopeMassByCode(residue)
		if !known {
			return mod.Modification{}, false
		}
		base = residueMass + d.fixedMassDiff(residue)
	}
	return d.derivedModification(residue, mass, mass-base)
}

// derivedModification builds the modification for a site mass the search summary does not declare,
// from the delta between that mass and its base. A delta too small to be a modification is drift: the
// mass is read as the declared modification nearest to it within minModificationDelta, or as nothing.
// A real delta the summary declares on another site takes the declared value, so one mass offset
// reads the same on every residue it lands on.
func (d declaredModifications) derivedModification(site string, mass float64, delta float64) (mod.Modification, bool) {

	if math.Abs(delta) < minModificationDelta {
		return d.nearest(site, mass, minModificationDelta)
	}

	return mod.Modification{
		Type:      mod.Assigned,
		Variable:  true,
		AminoAcid: site,
		MassDiff:  d.declaredDeltaOrRounded(delta),
	}, true
}

// declaredDeltaOrRounded returns the declared delta nearest to a derived one within tolerance, or the
// derived delta rounded to four decimals.
func (d declaredModifications) declaredDeltaOrRounded(delta float64) float64 {

	i := sort.SearchFloat64s(d.deltas, delta)

	best := roundTo4Decimals(delta)
	gap := declaredMassTolerance
	for _, k := range []int{i - 1, i} {
		if k < 0 || k >= len(d.deltas) {
			continue
		}
		if g := math.Abs(d.deltas[k] - delta); g < gap {
			best, gap = d.deltas[k], g
		}
	}
	return best
}

// declaredMassAtPosition is the mass of another entry at the same position that the search summary
// declares, if there is one.
func (d declaredModifications) declaredMassAtPosition(residue string, position int, mass float64, entries []spc.ModAminoacidMass) (float64, bool) {
	for _, e := range entries {
		if e.Position != position || e.Mass == mass {
			continue
		}
		if _, ok := d.findDeclared(residue, e.Mass); ok {
			return e.Mass, true
		}
	}
	return 0, false
}

// roundTo4Decimals rounds a delta derived from five-decimal pepXML masses to the four decimals the
// search summary uses. The summary itself is loaded through uti.ToFixed, which truncates, so a delta
// that no header line declares can differ from a declared spelling in the last decimal.
func roundTo4Decimals(x float64) float64 {
	return math.Round(x*1e4) / 1e4
}

// addObservedMassDiff records the PSM's observed delta mass unless a modification already has that key.
func addObservedMassDiff(index map[string]mod.Modification, massDiff float64) {
	key := fmt.Sprintf("%.4f", massDiff)
	if _, ok := index[key]; !ok {
		index[key] = mod.Modification{
			Index:    key,
			Name:     "Unknown",
			Type:     mod.Observed,
			MassDiff: massDiff,
		}
	}
}
