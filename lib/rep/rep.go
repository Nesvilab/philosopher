package rep

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/cmsl/bio"
	"github.com/prvst/cmsl/utils"
	"github.com/prvst/philosopher/lib/clas"
	"github.com/prvst/philosopher/lib/data"
	"github.com/prvst/philosopher/lib/meta"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/prvst/philosopher/lib/tmt"
	"github.com/prvst/philosopher/lib/uni"
	"github.com/prvst/philosopher/lib/xml"
)

// Evidence ...
type Evidence struct {
	meta.Data
	PSM           PSMEvidenceList
	Ions          IonEvidenceList
	Peptides      PeptideEvidenceList
	Proteins      ProteinEvidenceList
	Mods          Modifications
	Modifications ModificationEvidence
	Combined      CombinedEvidenceList
}

// Modifications ...
type Modifications struct {
	DefinedModMassDiff  map[float64]float64
	DefinedModAminoAcid map[float64]string
}

// PSMEvidence struct
type PSMEvidence struct {
	Index                 uint32
	Spectrum              string
	Scan                  int
	Peptide               string
	Protein               string
	ModifiedPeptide       string
	AlternativeProteins   []string
	ModPositions          []uint16
	AssignedModMasses     []float64
	AssignedMassDiffs     []float64
	AssignedModifications []string
	ObservedModifications []string
	AssumedCharge         uint8
	HitRank               uint8
	PrecursorNeutralMass  float64
	PrecursorExpMass      float64
	RetentionTime         float64
	CalcNeutralPepMass    float64
	RawMassdiff           float64
	Massdiff              float64
	Probability           float64
	Expectation           float64
	Xcorr                 float64
	DeltaCN               float64
	SpRank                float64
	DiscriminantValue     float64
	ModNtermMass          float64
	Intensity             float64
	Purity                float64
	Labels                tmt.Labels
}

// PSMEvidenceList ...
type PSMEvidenceList []PSMEvidence

func (a PSMEvidenceList) Len() int           { return len(a) }
func (a PSMEvidenceList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a PSMEvidenceList) Less(i, j int) bool { return a[i].Spectrum < a[j].Spectrum }

// IonEvidence groups all valid info about peptide ions for reports
type IonEvidence struct {
	Sequence                string
	ModifiedSequence        string
	AssignedModifications   map[string]uint8
	ObservedModifications   map[string]uint8
	RetentionTime           string
	Spectra                 map[string]int
	MappedProteins          map[string]uint8
	IndiMappedProteins      map[string]uint8
	ChargeState             uint8
	Spc                     int
	MZ                      float64
	PeptideMass             float64
	IsNondegenerateEvidence bool
	Weight                  float64
	GroupWeight             float64
	Intensity               float64
	Probability             float64
	Expectation             float64
	IsRazor                 bool
	Labels                  tmt.Labels
	SummedLabelIntensity    float64
	ModifiedObservations    int
	UnModifiedObservations  int
}

// IonEvidenceList ...
type IonEvidenceList []IonEvidence

func (a IonEvidenceList) Len() int           { return len(a) }
func (a IonEvidenceList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a IonEvidenceList) Less(i, j int) bool { return a[i].Sequence < a[j].Sequence }

// PeptideEvidence groups all valid info about peptide ions for reports
type PeptideEvidence struct {
	Sequence               string
	ChargeState            map[uint8]uint8
	Spc                    int
	Intensity              float64
	ModifiedObservations   int
	UnModifiedObservations int
}

// PeptideEvidenceList ...
type PeptideEvidenceList []PeptideEvidence

func (a PeptideEvidenceList) Len() int           { return len(a) }
func (a PeptideEvidenceList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a PeptideEvidenceList) Less(i, j int) bool { return a[i].Sequence < a[j].Sequence }

// ProteinEvidence ...
type ProteinEvidence struct {
	OriginalHeader               string
	ProteinName                  string
	ProteinGroup                 uint32
	ProteinSubGroup              string
	ProteinID                    string
	EntryName                    string
	Description                  string
	Organism                     string
	Length                       int
	Coverage                     float32
	GeneNames                    string
	ProteinExistence             string
	Sequence                     string
	IndiProtein                  map[string]uint8
	UniqueStrippedPeptides       int
	TotalNumRazorPeptides        int
	TotalNumPeptideIons          int
	NumURazorPeptideIons         int // Unique + razor
	TotalPeptideIons             map[string]IonEvidence
	UniquePeptideIons            map[string]IonEvidence
	URazorPeptideIons            map[string]IonEvidence // Unique + razor
	TotalSpC                     int
	UniqueSpC                    int
	RazorSpC                     int // Unique + razor
	TotalIntensity               float64
	UniqueIntensity              float64
	RazorIntensity               float64 // Unique + razor
	Probability                  float64
	TopPepProb                   float64
	IsDecoy                      bool
	IsContaminant                bool
	TotalLabels                  tmt.Labels
	UniqueLabels                 tmt.Labels
	RazorLabels                  tmt.Labels // Unique + razor
	URazorModifiedObservations   int
	URazorUnModifiedObservations int
}

// ProteinEvidenceList list
type ProteinEvidenceList []ProteinEvidence

func (a ProteinEvidenceList) Len() int           { return len(a) }
func (a ProteinEvidenceList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ProteinEvidenceList) Less(i, j int) bool { return a[i].ProteinGroup < a[j].ProteinGroup }

// CombinedEvidence represents all combined proteins detected
type CombinedEvidence struct {
	GroupNumber             uint32
	SiblingID               string
	ProteinName             string
	ProteinID               string
	EntryName               string
	GeneNames               string
	Length                  int
	Names                   []string
	UniqueStrippedPeptides  int
	TotalPeptideIonStrings  map[string]int
	UniquePeptideIonStrings map[string]int
	TotalPeptideIons        int
	UniquePeptideIons       int
	SharedPeptideIons       int
	TotalSpc                []int
	UniqueSpc               []int
	ProteinProbability      float64
	TopPepProb              float64
}

// CombinedEvidenceList ...
type CombinedEvidenceList []CombinedEvidence

// ModificationEvidence represents the list of modifications and the mod bins
type ModificationEvidence struct {
	MassBins     []MassBin
	AssignedBins []AssignedBin
	// MapModNamesToAverageMasses   map[string]float64
	// MapModMassToAverageMasses    map[float64]float64
	// MapModNamesToPSMs            map[string]PSMEvidenceList
	// MapModMassesToPSMs           map[float64]PSMEvidenceList
	// MapModMassToModNames         map[float64][]string
	// MapModNameToModMass          map[string]float64
	// MapModNameToMonoisotopicMass map[string]float64
	// MapModNameToComposition      map[string]string
	// UnModObservations            int
}

// MassBin represents each bin from the mass distribution
type MassBin struct {
	LowerMass     float64
	HigherRight   float64
	MassCenter    float64
	AverageMass   float64
	CorrectedMass float64
	Modifications []string
	Elements      PSMEvidenceList
	Spectra       map[string]int
}

// AssignedBin contains the bin data after collapsing entries into apex peaks
type AssignedBin struct {
	Mass                float64
	MappedModifications []string
	Elements            PSMEvidenceList
}

// AssignedBins is a list of assignedbins
type AssignedBins []AssignedBin

func (a AssignedBins) Len() int           { return len(a) }
func (a AssignedBins) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a AssignedBins) Less(i, j int) bool { return a[i].Mass < a[j].Mass }

// New constructor
func New() Evidence {

	var o Evidence
	var m meta.Data
	m.Restore(sys.Meta())

	o.UUID = m.UUID
	o.Distro = m.Distro
	o.Home = m.Home
	o.MetaFile = m.MetaFile
	o.MetaDir = m.MetaDir
	o.DB = m.DB
	o.Temp = m.Temp
	o.TimeStamp = m.TimeStamp
	o.OS = m.OS
	o.Arch = m.Arch

	return o
}

// AssemblePSMReport ...
func (e *Evidence) AssemblePSMReport(pep xml.PepIDList, decoyTag string) error {

	var list PSMEvidenceList

	// instantiate the unimod database for mapping fixed modification
	u := uni.New()
	u.ProcessUniMOD()

	// get the list of modifications for each psm spectrum
	var uniqModPSM = make(map[string][]string)
	for _, i := range e.Modifications.AssignedBins {
		for _, j := range i.Elements {
			uniqModPSM[j.Spectrum] = append(uniqModPSM[j.Spectrum], i.MappedModifications...)
		}
	}

	for _, i := range pep {

		if !clas.IsDecoyPSM(i, decoyTag) {

			var p PSMEvidence

			p.Index = i.Index
			p.Spectrum = i.Spectrum
			p.Scan = i.Scan
			p.Peptide = i.Peptide
			p.Protein = i.Protein
			p.ModifiedPeptide = i.ModifiedPeptide
			p.AlternativeProteins = i.AlternativeProteins
			p.ModPositions = i.ModPositions
			p.AssignedModMasses = i.AssignedModMasses
			p.AssignedMassDiffs = i.AssignedMassDiffs
			p.AssumedCharge = i.AssumedCharge
			p.HitRank = i.HitRank
			p.PrecursorNeutralMass = i.PrecursorNeutralMass
			p.PrecursorExpMass = i.PrecursorExpMass
			p.RetentionTime = i.RetentionTime
			p.CalcNeutralPepMass = i.CalcNeutralPepMass
			p.RawMassdiff = i.RawMassDiff
			p.Massdiff = i.Massdiff
			p.Probability = i.Probability
			p.Expectation = i.Expectation
			p.Xcorr = i.Xcorr
			p.DeltaCN = i.DeltaCN
			p.SpRank = i.SpRank
			p.DiscriminantValue = i.DiscriminantValue
			p.ModNtermMass = i.ModNtermMass
			p.Intensity = i.Intensity

			// search for an explanation for fixed modifications
			for _, j := range u.Modifications {
				mod := utils.Round(j.MonoMass, 5, 2)
				for _, l := range p.AssignedMassDiffs {
					if l >= (mod-0.2) && l <= (mod+0.2) {
						if !strings.Contains(j.Description, "substitution") {
							fullname := fmt.Sprintf("%s (%s)", j.Title, j.Description)
							p.AssignedModifications = append(p.AssignedModifications, fullname)
						}
						break
					}
				}
			}

			// search for an explanation for open search observed modifications
			v, ok := uniqModPSM[i.Spectrum]
			if ok {
				p.ObservedModifications = append(p.ObservedModifications, v...)
			}

			list = append(list, p)
		}
	}

	sort.Sort(list)
	e.PSM = list

	return nil
}

// PSMReport report all psms from study that passed the FDR filter
func (e *Evidence) PSMReport() {

	output := fmt.Sprintf("%s%spsm.tsv", e.Temp, string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create report file:", err)
	}
	defer file.Close()

	_, err = io.WriteString(file, "Spectrum\tPeptide\tCharge\tRetention\tCalculated M/Z\tObserved M/Z\tOriginal Delta Mass\tAdjusted Delta Mass\tExperimental Mass\tPeptide Mass\tPeptideProphet Probability\tExpectation\tSpecified Modifications\tIdentified Modifications\n")
	if err != nil {
		logrus.Fatal("Cannot print PSM to file")
	}

	for _, i := range e.PSM {

		line := fmt.Sprintf("%s\t%s\t%d\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%s\t%s\n",
			i.Spectrum,
			i.Peptide,
			i.AssumedCharge,
			i.RetentionTime,
			((i.CalcNeutralPepMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge)),
			((i.PrecursorNeutralMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge)),
			i.RawMassdiff,
			i.Massdiff,
			i.PrecursorNeutralMass,
			i.CalcNeutralPepMass,
			i.Probability,
			i.Expectation,
			strings.Join(i.AssignedModifications, ", "),
			strings.Join(i.ObservedModifications, ", "),
		)
		_, err = io.WriteString(file, line)
		if err != nil {
			logrus.Fatal("Cannot print PSM to file")
		}
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}

// AssembleIonReport reports consist on ion reporting
func (e *Evidence) AssembleIonReport(ion xml.PepIDList, decoyTag string) error {

	var list IonEvidenceList
	var psmPtMap = make(map[string][]string)
	var err error

	// collapse all psm to protein based on Peptide-level identifications
	for _, i := range e.PSM {
		psmPtMap[i.Spectrum] = append(psmPtMap[i.Spectrum], i.Protein)
		psmPtMap[i.Spectrum] = append(psmPtMap[i.Spectrum], i.AlternativeProteins...)
	}

	for _, i := range ion {
		if !clas.IsDecoyPSM(i, decoyTag) {

			var pr IonEvidence

			pr.Spectra = make(map[string]int)
			pr.MappedProteins = make(map[string]uint8)
			pr.IndiMappedProteins = make(map[string]uint8)
			pr.ObservedModifications = make(map[string]uint8)
			pr.AssignedModifications = make(map[string]uint8)

			pr.Spectra[i.Spectrum]++

			pr.Sequence = i.Peptide
			pr.ModifiedSequence = i.ModifiedPeptide
			pr.MZ = utils.Round(((i.CalcNeutralPepMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge)), 5, 4)
			pr.ChargeState = i.AssumedCharge
			pr.PeptideMass = i.CalcNeutralPepMass
			pr.Probability = i.Probability
			pr.Expectation = i.Expectation

			// get he list of indi proteins from pepXML data
			v, ok := psmPtMap[i.Spectrum]
			if ok {
				for _, i := range v {
					pr.MappedProteins[i] = 0
					pr.IndiMappedProteins[i]++
				}
			}

			list = append(list, pr)
		}
	}

	sort.Sort(list)
	e.Ions = list

	return err
}

// UpdateIonAssignedAndObservedMods collects all Assigned and Observed modifications from
// individual PSM and assign them to ions
func (e *Evidence) UpdateIonAssignedAndObservedMods() {

	//var list IonEvidenceList

	for i := range e.Ions {
		var ion string
		if len(e.Ions[i].ModifiedSequence) > 0 {
			ion = fmt.Sprintf("%s#%d", e.Ions[i].ModifiedSequence, e.Ions[i].ChargeState)
		} else {
			ion = fmt.Sprintf("%s#%d", e.Ions[i].Sequence, e.Ions[i].ChargeState)
		}

		for _, j := range e.PSM {
			var psmIon string
			if len(j.ModifiedPeptide) > 0 {
				psmIon = fmt.Sprintf("%s#%d", j.ModifiedPeptide, j.AssumedCharge)
			} else {
				psmIon = fmt.Sprintf("%s#%d", j.Peptide, j.AssumedCharge)
			}

			if ion == psmIon {
				for _, k := range j.AssignedModifications {
					e.Ions[i].AssignedModifications[k]++
				}
				for _, k := range j.ObservedModifications {
					e.Ions[i].ObservedModifications[k]++
				}

				break
			}

		}

		//list = append(list, e.Ions[i])
	}

	//e.Ions = list

	return
}

// UpdateIonModCount counts how many times each ion is observed modified and not modified
func (e *Evidence) UpdateIonModCount() {

	// recreate the ion list from the main report object
	var AllIons = make(map[string]int)
	var ModIons = make(map[string]int)
	var UnModIons = make(map[string]int)

	for _, i := range e.Ions {
		var ion string
		if len(i.ModifiedSequence) > 0 {
			ion = fmt.Sprintf("%s#%d", i.ModifiedSequence, i.ChargeState)
		} else {
			ion = fmt.Sprintf("%s#%d", i.Sequence, i.ChargeState)
		}
		AllIons[ion] = 0
		ModIons[ion] = 0
		UnModIons[ion] = 0
	}

	// range over PSMs looking for modified and not modified evidences
	// if they exist on the ions map, get the numbers
	for _, i := range e.PSM {
		var psmIon string
		if len(i.ModifiedPeptide) > 0 {
			psmIon = fmt.Sprintf("%s#%d", i.ModifiedPeptide, i.AssumedCharge)
		} else {
			psmIon = fmt.Sprintf("%s#%d", i.Peptide, i.AssumedCharge)
		}

		// check the map
		_, ok := AllIons[psmIon]
		if ok {

			if i.Massdiff >= -0.99 && i.Massdiff <= 0.99 {
				UnModIons[psmIon]++
			} else {
				ModIons[psmIon]++
			}

		}
	}

	for i := range e.Ions {
		var ion string
		if len(e.Ions[i].ModifiedSequence) > 0 {
			ion = fmt.Sprintf("%s#%d", e.Ions[i].ModifiedSequence, e.Ions[i].ChargeState)
		} else {
			ion = fmt.Sprintf("%s#%d", e.Ions[i].Sequence, e.Ions[i].ChargeState)
		}

		v1, ok1 := UnModIons[ion]
		if ok1 {
			e.Ions[i].UnModifiedObservations = v1
		}

		v2, ok2 := ModIons[ion]
		if ok2 {
			e.Ions[i].ModifiedObservations = v2
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

// PeptideIonReport reports consist on ion reporting
func (e *Evidence) PeptideIonReport() {

	output := fmt.Sprintf("%s%sion.tsv", e.Temp, string(filepath.Separator))

	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Could not create peptide output file")
	}
	defer file.Close()

	_, err = io.WriteString(file, "Peptide Sequence\tM/Z\tCharge\tExperimental Mass\tProbability\tExpectation\tSpectral Count\tUnmodified Occurrences\tModified Occurrences\tIntensity\tAssigned Modifications\tObserved Modifications\tIntensity\tMapped Proteins\tProtein IDs\n")
	if err != nil {
		logrus.Fatal("Cannot create peptide ion report header")
	}

	// peptides with no mapped poteins are related to contaminants
	// and reverse sequences. They are dificult to clean because
	// in some cases they are shared between a match decoy and a target,
	// so they stay on the lists but cannot be mapped back to the
	// original proteins. These cases should be rare to find.
	for _, i := range e.Ions {

		var pts []string
		//var ipts []string

		if len(i.MappedProteins) > 0 {

			if len(e.Proteins) > 1 {

				for k := range i.MappedProteins {
					pts = append(pts, k)
				}

				// for k := range i.IndiMappedProteins {
				// 	ipts = append(ipts, k)
				// }

				var amods []string
				for j := range i.AssignedModifications {
					amods = append(amods, j)
				}

				var omods []string
				for j := range i.ObservedModifications {
					omods = append(omods, j)
				}

				line := fmt.Sprintf("%s\t%.4f\t%d\t%.4f\t%.4f\t%.4f\t%d\t%d\t%d\t%.4f\t%s\t%s\t%.4f\t%d\t%s\n",
					i.Sequence,
					i.MZ,
					i.ChargeState,
					i.PeptideMass,
					i.Probability,
					i.Expectation,
					i.Spc,
					i.UnModifiedObservations,
					i.ModifiedObservations,
					i.Intensity,
					strings.Join(amods, ", "),
					strings.Join(omods, ", "),
					i.Intensity,
					len(i.MappedProteins),
					strings.Join(pts, ", "),
				)
				_, err = io.WriteString(file, line)
				if err != nil {
					logrus.Fatal("Cannot print PSM to file")
				}
			}
		}
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}

// AssemblePeptideReport reports consist on ion reporting
func (e *Evidence) AssemblePeptideReport(pep xml.PepIDList, decoyTag string) error {

	var list PeptideEvidenceList
	var pepSeqMap = make(map[string]uint8)
	var pepCSMap = make(map[string][]uint8)
	var pepSpc = make(map[string]int)
	var pepInt = make(map[string]float64)
	var err error

	for _, i := range pep {
		if !clas.IsDecoyPSM(i, decoyTag) {
			pepSeqMap[i.Peptide] = 0
			pepSpc[i.Peptide] = 0
			pepInt[i.Peptide] = 0
		}
	}

	// TODO review this method, Intensity quant is not working
	for _, i := range e.PSM {
		_, ok := pepSeqMap[i.Peptide]
		if ok {
			pepCSMap[i.Peptide] = append(pepCSMap[i.Peptide], i.AssumedCharge)
			pepSpc[i.Peptide]++
			if i.Intensity > pepInt[i.Peptide] {
				pepInt[i.Peptide] = i.Intensity
			}
		}
	}

	for k := range pepSeqMap {

		var pep PeptideEvidence
		pep.ChargeState = make(map[uint8]uint8)
		pep.Sequence = k

		for _, i := range pepCSMap[k] {
			pep.ChargeState[i] = 0
		}
		pep.Spc = pepSpc[k]
		pep.Intensity = pepInt[k]

		list = append(list, pep)
	}

	sort.Sort(list)
	e.Peptides = list

	return err
}

// PeptideReport reports consist on ion reporting
func (e *Evidence) PeptideReport() {

	output := fmt.Sprintf("%s%speptide.tsv", e.Temp, string(filepath.Separator))

	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Could not create peptide output file")
	}
	defer file.Close()

	_, err = io.WriteString(file, "Peptide\tCharges\tSpectral Count\tUnmodified Occurrences\tModified Occurrences\n")
	if err != nil {
		logrus.Fatal("Cannot create peptide report header")
	}

	for _, i := range e.Peptides {

		var cs []string
		for j := range i.ChargeState {
			cs = append(cs, strconv.Itoa(int(j)))
		}

		line := fmt.Sprintf("%s\t%s\t%d\t%d\t%d\n",
			i.Sequence,
			strings.Join(cs, ", "),
			i.Spc,
			i.UnModifiedObservations,
			i.ModifiedObservations,
			//i.Intensity,
		)
		_, err = io.WriteString(file, line)
		if err != nil {
			logrus.Fatal("Cannot print PSM to file")
		}
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}

// ProteinFastaReport saves to disk a filtered FASTA file with FDR aproved proteins
func (e *Evidence) ProteinFastaReport() error {

	output := fmt.Sprintf("%s%sproteins.fas", e.Temp, string(filepath.Separator))

	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Could not create output file")
	}
	defer file.Close()

	for _, i := range e.Proteins {
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

// AssembleProteinReport ...
func (e *Evidence) AssembleProteinReport(pro xml.ProtIDList, decoyTag string) error {

	var list ProteinEvidenceList
	var err error

	for _, i := range pro {
		if !strings.HasPrefix(i.ProteinName, decoyTag) {

			var rep ProteinEvidence

			rep.TotalPeptideIons = make(map[string]IonEvidence)
			rep.UniquePeptideIons = make(map[string]IonEvidence)
			rep.URazorPeptideIons = make(map[string]IonEvidence)
			rep.IndiProtein = make(map[string]uint8)

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

			for j := range e.Ions {

				_, ok := e.Ions[j].MappedProteins[i.ProteinName]
				if ok {

					var ion string
					if len(e.Ions[j].ModifiedSequence) > 0 {
						ion = fmt.Sprintf("%s#%d", e.Ions[j].ModifiedSequence, e.Ions[j].ChargeState)
					} else {
						ion = fmt.Sprintf("%s#%d", e.Ions[j].Sequence, e.Ions[j].ChargeState)
					}

					rep.TotalNumPeptideIons++
					rep.TotalPeptideIons[ion] = e.Ions[j]

					for _, k := range i.PeptideIons {

						var ption string
						if len(k.ModifiedPeptide) > 0 {
							ption = fmt.Sprintf("%s#%d", k.ModifiedPeptide, k.Charge)
						} else {
							ption = fmt.Sprintf("%s#%d", k.PeptideSequence, k.Charge)
						}

						if ion == ption {

							if k.Razor == 1 {
								e.Ions[j].IsRazor = true
							}

							//if ion == ption {

							if k.IsUnique == true {
								rep.UniquePeptideIons[ion] = e.Ions[j]
								rep.URazorPeptideIons[ion] = e.Ions[j]
								rep.NumURazorPeptideIons++
								rep.TotalNumRazorPeptides++

								rep.URazorUnModifiedObservations += e.Ions[j].UnModifiedObservations
								rep.URazorModifiedObservations += e.Ions[j].ModifiedObservations
							}

							if k.Razor == 1 {
								rep.URazorUnModifiedObservations += e.Ions[j].UnModifiedObservations
								rep.URazorModifiedObservations += e.Ions[j].ModifiedObservations

								rep.URazorPeptideIons[ion] = e.Ions[j]
								rep.TotalNumRazorPeptides++
							}

						}

					}

				}
			}

			list = append(list, rep)
		}
	}

	var dtb data.Base
	dtb.Restore()

	if len(dtb.Records) < 1 {
		return errors.New("Cant locate database data")
	}

	for i := range list {
		for _, j := range dtb.Records {
			// fix the name sand headers and pull database information into proteinreport
			if strings.Contains(j.OriginalHeader, list[i].ProteinName) {
				if (j.IsDecoy == true && list[i].IsDecoy == true) || (j.IsDecoy == false && list[i].IsDecoy == false) {
					list[i].OriginalHeader = j.OriginalHeader
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
func (e *Evidence) ProteinReport() {

	// create result file
	output := fmt.Sprintf("%s%sreport.tsv", e.Temp, string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create report file:", err)
	}
	defer file.Close()

	line := fmt.Sprintf("Group\tSubGroup\tProtein ID\tEntry Name\tLength\tPercent Coverage\tOrganism\tDescription\tProtein Existence\tGenes\tProtein Probability\tTop Peptide Probability\tStripped Peptides\tTotal Peptide Ions\tUnique Peptide Ions\tTotal Spectral Count\tUnique Spectral Count\tRazor Spectral Count\tRazor Unmodified Occurrences\tRazor Modified Occurrences\tTotal Intensity\tUnique Intensity\tRazor Intensity\tIndistinguishable Proteins\n")

	n, err := io.WriteString(file, line)
	if err != nil {
		logrus.Fatal(n, err)
	}

	for _, i := range e.Proteins {

		var ip []string
		for k := range i.IndiProtein {
			ip = append(ip, k)
		}

		// proteins with almost no evidences, and completely shared with decoys are eliminated from the analysis,
		// in most cases proteins with one small peptide shared with a decoy
		//if len(i.TotalPeptideIons) > 0 {

		line = fmt.Sprintf("%d\t%s\t%s\t%s\t%d\t%.2f\t%s\t%s\t%s\t%s\t%.4f\t%.4f\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%6.f\t%6.f\t%6.f\t%s\t",
			i.ProteinGroup,                 // Group
			i.ProteinSubGroup,              // SubGroup
			i.ProteinID,                    // Protein ID
			i.EntryName,                    // Entry Name
			i.Length,                       // Length
			i.Coverage,                     // Percent Coverage
			i.Organism,                     // Organism
			i.Description,                  // Description
			i.ProteinExistence,             // Protein Existence
			i.GeneNames,                    // Genes
			i.Probability,                  // Protein Probability
			i.TopPepProb,                   // Top Peptide Probability
			i.UniqueStrippedPeptides,       // Stripped Peptides
			i.TotalNumPeptideIons,          // Total Peptide Ions
			i.NumURazorPeptideIons,         // Unique Peptide Ions
			i.TotalSpC,                     // Total Spectral Count
			i.UniqueSpC,                    // Unique Spectral Count
			i.RazorSpC,                     // Razor Spectral Count
			i.URazorUnModifiedObservations, // Unmodified Occurrences
			i.URazorModifiedObservations,   // Modified Occurrences
			i.TotalIntensity,               // Total Intensity
			i.UniqueIntensity,              // Unique Intensity
			i.RazorIntensity,               // Razor Intensity
			strings.Join(ip, ", "))         // Indistinguishable Proteins

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

// ProteinQuantReport ...
func (e *Evidence) ProteinQuantReport() {

	// create result file
	output := fmt.Sprintf("%s%sreport.tsv", e.Temp, string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create report file:", err)
	}
	defer file.Close()

	//Total Channel 1\tTotal Channel 2\tTotal Channel 3\t Total Channel 4\t Total Channel 5\t Total Channel 6\t Total Channel 7\tTotal Channel 8\tTotal Channel 9\tTotal Channel 10\tUnique Channel 1\tUnique Channel 2\tUnique Channel 3\tUnique Channel 4\tUnique Channel 5\tUnique Channel 6\tUnique Channel 7\tUnique Channel 8\tUnique Channel 9\tUnique Channel 10\n")
	line := fmt.Sprintf("Group\tSubGroup\tProtein ID\tEntry Name\tLength\tPercent Coverage\tDescription\tProtein Existence\tGenes\tProtein Probability\tTop Peptide Probability\tUnique Stripped Peptides\tRazor Peptides\tTotal Peptide Ions\tUnique Peptide Ions\tTotal Spectral Count\tUnique Spectral Count\tTotal Intensity\tUnique Intensity\tTotal Raw Channel 1\tTotal Raw Channel 2\tTotal Raw Channel 3\t Total Raw Channel 4\t Total Raw Channel 5\t Total Raw Channel 6\t Total Raw Channel 7\tTotal Raw Channel 8\tTotal Raw Channel 9\tTotal Raw Channel 10\tUnique Raw Channel 1\tUnique Raw Channel 2\tUnique Raw Channel 3\tUnique Raw Channel 4\tUnique Raw Channel 5\tUnique Raw Channel 6\tUnique Raw Channel 7\tUnique Raw Channel 8\tUnique Raw Channel 9\tUnique Raw Channel 10\tRazor Raw Channel 1\tRazor Raw Channel 2\tRazor Raw Channel 3\tRazor Raw Channel 4\tRazor Raw Channel 5\tRazor Raw Channel 6\tRazor Raw Channel 7\tRazor Raw Channel 8\tRazor Raw Channel 9\tRazor Raw Channel 10\tTotal Channel 1\tTotal Channel 2\tTotal Channel 3\t Total Channel 4\t Total Channel 5\t Total Channel 6\t Total Channel 7\tTotal Channel 8\tTotal Channel 9\tTotal Channel 10\tUnique Channel 1\tUnique Channel 2\tUnique Channel 3\tUnique Channel 4\tUnique Channel 5\tUnique Channel 6\tUnique Channel 7\tUnique Channel 8\tUnique Channel 9\tUnique Channel 10\tRazor Channel 1\tRazor Channel 2\tRazor Channel 3\tRazor Channel 4\tRazor Channel 5\tRazor Channel 6\tRazor Channel 7\tRazor Channel 8\tRazor Channel 9\tRazor Channel 10\tIndistinguishable Proteins\n")

	n, err := io.WriteString(file, line)
	if err != nil {
		logrus.Fatal(n, err)
	}

	for _, i := range e.Proteins {

		var ip []string
		for k := range i.IndiProtein {
			ip = append(ip, k)
		}

		//%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t
		if len(i.TotalPeptideIons) > 0 {
			line = fmt.Sprintf("%d\t%s\t%s\t%s\t%d\t%.2f\t%s\t%s\t%s\t%.4f\t%.4f\t%d\t%d\t%d\t%d\t%d\t%d\t%6.f\t%6.f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%s\t",
				i.ProteinGroup,
				i.ProteinSubGroup,
				i.ProteinID,
				i.EntryName,
				i.Length,
				i.Coverage,
				i.Description,
				i.ProteinExistence,
				i.GeneNames,
				i.Probability,
				i.TopPepProb,
				i.UniqueStrippedPeptides,
				i.TotalNumRazorPeptides,
				i.TotalNumPeptideIons,
				i.NumURazorPeptideIons,
				i.TotalSpC,
				i.UniqueSpC,
				i.TotalIntensity,
				i.UniqueIntensity,
				i.TotalLabels.Channel1.NormIntensity,
				i.TotalLabels.Channel2.NormIntensity,
				i.TotalLabels.Channel3.NormIntensity,
				i.TotalLabels.Channel4.NormIntensity,
				i.TotalLabels.Channel5.NormIntensity,
				i.TotalLabels.Channel6.NormIntensity,
				i.TotalLabels.Channel7.NormIntensity,
				i.TotalLabels.Channel8.NormIntensity,
				i.TotalLabels.Channel9.NormIntensity,
				i.TotalLabels.Channel10.NormIntensity,
				i.UniqueLabels.Channel1.NormIntensity,
				i.UniqueLabels.Channel2.NormIntensity,
				i.UniqueLabels.Channel3.NormIntensity,
				i.UniqueLabels.Channel4.NormIntensity,
				i.UniqueLabels.Channel5.NormIntensity,
				i.UniqueLabels.Channel6.NormIntensity,
				i.UniqueLabels.Channel7.NormIntensity,
				i.UniqueLabels.Channel8.NormIntensity,
				i.UniqueLabels.Channel9.NormIntensity,
				i.UniqueLabels.Channel10.NormIntensity,
				i.RazorLabels.Channel1.NormIntensity,
				i.RazorLabels.Channel2.NormIntensity,
				i.RazorLabels.Channel3.NormIntensity,
				i.RazorLabels.Channel4.NormIntensity,
				i.RazorLabels.Channel5.NormIntensity,
				i.RazorLabels.Channel6.NormIntensity,
				i.RazorLabels.Channel7.NormIntensity,
				i.RazorLabels.Channel8.NormIntensity,
				i.RazorLabels.Channel9.NormIntensity,
				i.RazorLabels.Channel10.NormIntensity,
				i.TotalLabels.Channel1.RatioIntensity,
				i.TotalLabels.Channel2.RatioIntensity,
				i.TotalLabels.Channel3.RatioIntensity,
				i.TotalLabels.Channel4.RatioIntensity,
				i.TotalLabels.Channel5.RatioIntensity,
				i.TotalLabels.Channel6.RatioIntensity,
				i.TotalLabels.Channel7.RatioIntensity,
				i.TotalLabels.Channel8.RatioIntensity,
				i.TotalLabels.Channel9.RatioIntensity,
				i.TotalLabels.Channel10.RatioIntensity,
				i.UniqueLabels.Channel1.RatioIntensity,
				i.UniqueLabels.Channel2.RatioIntensity,
				i.UniqueLabels.Channel3.RatioIntensity,
				i.UniqueLabels.Channel4.RatioIntensity,
				i.UniqueLabels.Channel5.RatioIntensity,
				i.UniqueLabels.Channel6.RatioIntensity,
				i.UniqueLabels.Channel7.RatioIntensity,
				i.UniqueLabels.Channel8.RatioIntensity,
				i.UniqueLabels.Channel9.RatioIntensity,
				i.UniqueLabels.Channel10.RatioIntensity,
				i.RazorLabels.Channel1.RatioIntensity,
				i.RazorLabels.Channel2.RatioIntensity,
				i.RazorLabels.Channel3.RatioIntensity,
				i.RazorLabels.Channel4.RatioIntensity,
				i.RazorLabels.Channel5.RatioIntensity,
				i.RazorLabels.Channel6.RatioIntensity,
				i.RazorLabels.Channel7.RatioIntensity,
				i.RazorLabels.Channel8.RatioIntensity,
				i.RazorLabels.Channel9.RatioIntensity,
				i.RazorLabels.Channel10.RatioIntensity,
				strings.Join(ip, ", "))

			line += "\n"
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

// AssembleModificationReport cretaes the modifications lists
func (e *Evidence) AssembleModificationReport() error {

	var modEvi ModificationEvidence

	var massWindow = float64(0.5)
	var binsize = float64(0.1)
	var amplitude = float64(500)

	var bins []MassBin

	nBins := (amplitude*(1/binsize) + 1) * 2
	for i := 0; i <= int(nBins); i++ {
		var b MassBin
		b.Spectra = make(map[string]int)

		b.LowerMass = -(amplitude) - (massWindow * binsize) + (float64(i) * binsize)
		b.LowerMass = utils.Round(b.LowerMass, 5, 4)

		b.HigherRight = -(amplitude) + (massWindow * binsize) + (float64(i) * binsize)
		b.HigherRight = utils.Round(b.HigherRight, 5, 4)

		b.MassCenter = -(amplitude) + (float64(i) * binsize)
		b.MassCenter = utils.Round(b.MassCenter, 5, 4)

		bins = append(bins, b)
	}

	// // lets consider fixed modifications from the search as massdiff values.
	// for i := range e.PSM {
	// 	if len(e.PSM[i].FixedMassDiffs) > 0 {
	// 		fmt.Println(e.PSM[i].FixedMassDiffs)
	// 		os.Exit(1)
	// 	}
	// }

	// calculate the total number of PSMs per cluster
	var counter int
	for i := range e.PSM {

		for j := range bins {
			if e.PSM[i].Massdiff > bins[j].LowerMass && e.PSM[i].Massdiff <= bins[j].HigherRight {
				bins[j].Elements = append(bins[j].Elements, e.PSM[i])
				bins[j].Spectra[e.PSM[i].Spectrum]++
				counter++
				break
			}
		}
	}

	// calculate average mass for each cluster
	var zeroBinMassDeviation float64
	for i := range bins {
		pep := bins[i].Elements
		total := 0.0
		for j := range pep {
			total += pep[j].Massdiff
		}
		if len(bins[i].Elements) > 0 {
			bins[i].AverageMass = (float64(total) / float64(len(pep)))
		} else {
			bins[i].AverageMass = 0
		}
		if bins[i].MassCenter == 0 {
			zeroBinMassDeviation = bins[i].AverageMass
		}

		bins[i].AverageMass = utils.Round(bins[i].AverageMass, 5, 4)
	}

	// correcting mass values based on Bin 0 average mass
	for i := range bins {
		if len(bins[i].Elements) > 0 {
			if bins[i].AverageMass > 0 {
				bins[i].CorrectedMass = (bins[i].AverageMass - zeroBinMassDeviation)
			} else {
				bins[i].CorrectedMass = (bins[i].AverageMass + zeroBinMassDeviation)
			}
		} else {
			bins[i].CorrectedMass = bins[i].MassCenter
		}
		bins[i].CorrectedMass = utils.Round(bins[i].CorrectedMass, 5, 4)
	}
	e.Modifications.MassBins = bins

	// inspect mass binning
	// for _, i := range bins {
	// 	fmt.Println(i.LowerMass, "\t", i.HigherRight, "\t", i.CorrectedMass, "\t", len(i.Elements))
	// }

	// This block applies the grouping logic to find apex n
	// var binSelection []float64
	// var elsSelector []int
	// var elsMap = make(map[float64]PSMEvidenceList)
	// var apexMass float64
	// var apexElements int
	// var selectedBins = make(map[float64]PSMEvidenceList)
	//
	// for _, i := range bins {
	//
	// 	if len(i.Elements) == 0 && len(binSelection) >= 3 {
	// 		//apexBinTolerance := utils.Round(float64(len(i.Elements))/float64(4), 5, 0)
	// 		apexBinTolerance := (apexElements / 4)
	// 		for j := 0; j <= len(binSelection)-1; j++ {
	// 			if elsSelector[j] >= apexBinTolerance {
	// 				selectedBins[apexMass] = append(selectedBins[apexMass], elsMap[binSelection[j]]...)
	// 			}
	// 		}
	// 		apexElements = 0
	// 		apexMass = 0
	// 		binSelection = nil
	// 		elsSelector = nil
	// 		elsMap = make(map[float64]PSMEvidenceList)
	// 	}
	//
	// 	if len(i.Elements) > 0 {
	// 		binSelection = append(binSelection, i.CorrectedMass)
	// 		elsSelector = append(elsSelector, len(i.Elements))
	// 		elsMap[i.CorrectedMass] = i.Elements
	// 		if len(i.Elements) > apexElements {
	// 			apexElements = len(i.Elements)
	// 			apexMass = i.CorrectedMass
	// 		}
	// 	}
	//
	// }

	// starting unimod parsing
	u := uni.New()
	u.ProcessUniMOD()

	// assign each modification to a certain bin based on a mass window
	var abins AssignedBins
	for _, i := range bins {

		if len(i.Elements) > 0 {

			var ab AssignedBin
			ab.Mass = utils.Round(i.CorrectedMass, 5, 2)
			ab.Elements = i.Elements

			for _, j := range u.Modifications {
				mod := utils.Round(j.MonoMass, 5, 2)

				if ab.Mass >= (mod-0.1) && ab.Mass <= (mod+0.1) {
					fullName := fmt.Sprintf("%s (%s)", j.Title, j.Description)
					ab.MappedModifications = append(ab.MappedModifications, fullName)
				}

			}

			// set all bins without unimod mappings to unknown
			if len(ab.MappedModifications) == 0 {
				ab.MappedModifications = append(ab.MappedModifications, "Unknown")
			}

			// reset the 0 bin to no modifications
			if ab.Mass == 0 {
				ab.MappedModifications = nil
			}

			abins = append(abins, ab)

		}
	}

	// var abins AssignedBins
	// for k, v := range selectedBins {
	//
	// 	var ab AssignedBin
	// 	ab.Mass = utils.Round(k, 5, 2)
	// 	ab.Elements = v
	//
	// 	for _, j := range u.Modifications {
	// 		mod := utils.Round(j.MonoMass, 5, 2)
	//
	// 		if ab.Mass >= (mod-0.1) && ab.Mass <= (mod+0.1) {
	// 			fullName := fmt.Sprintf("%s (%s)", j.Title, j.Description)
	// 			ab.MappedModifications = append(ab.MappedModifications, fullName)
	// 		}
	//
	// 	}
	//
	// 	// set all bins without unimod mappings to unknown
	// 	if len(ab.MappedModifications) == 0 {
	// 		ab.MappedModifications = append(ab.MappedModifications, "Unknown")
	// 	}
	//
	// 	// reset the 0 bin to no modifications
	// 	if ab.Mass == 0 {
	// 		ab.MappedModifications = nil
	// 	}
	//
	// 	abins = append(abins, ab)
	// }

	// inspect assigned mass binning
	// for _, i := range abins {
	// 	fmt.Println(i.Mass, "\t", len(i.Elements))
	// }

	e.Modifications = modEvi
	e.Modifications.MassBins = bins

	sort.Sort(abins)
	e.Modifications.AssignedBins = abins

	return nil
}

// ModificationReport ...
func (e *Evidence) ModificationReport() {

	// create result file
	output := fmt.Sprintf("%s%smodifications.tsv", e.Temp, string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create report file:", err)
	}
	defer file.Close()

	line := fmt.Sprintf("Mass Bin\tNumber of PSMs\tModification\n")

	n, err := io.WriteString(file, line)
	if err != nil {
		logrus.Fatal(n, err)
	}

	for _, i := range e.Modifications.AssignedBins {

		line = fmt.Sprintf("%.4f\t%d\t",
			i.Mass,          // mass bins
			len(i.Elements), // number of psms
		)

		// line = fmt.Sprintf("%.4f\t%d\t%s\t",
		// 	i.Mass,          // mass bins
		// 	len(i.Elements), // number of psms
		// 	strings.Join(i.MappedModifications, ", "))

		line += "\n"
		n, err := io.WriteString(file, line)
		if err != nil {
			logrus.Fatal(n, err)
		}

	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}

// PlotMassHist plots the delta mass histogram
func (e *Evidence) PlotMassHist() error {

	outfile := fmt.Sprintf("%s%sdelta-mass.html", e.Temp, string(filepath.Separator))

	file, err := os.Create(outfile)
	if err != nil {
		return errors.New("Could not create output for delta mass binning")
	}
	defer file.Close()

	var xvar []string
	var yvar []string

	for _, i := range e.Modifications.MassBins {
		xel := fmt.Sprintf("'%.2f',", i.MassCenter)
		xvar = append(xvar, xel)
		yel := fmt.Sprintf("'%d',", len(i.Elements))
		yvar = append(yvar, yel)
	}

	xline := fmt.Sprintf("	  x: %s,", xvar)
	yline := fmt.Sprintf("	  y: %s,", yvar)

	_, err = io.WriteString(file, "<head>\n")
	_, err = io.WriteString(file, "  <script src=\"https://cdn.plot.ly/plotly-latest.min.js\"></script>\n")
	_, err = io.WriteString(file, "</head>\n")
	_, err = io.WriteString(file, "<body>\n")
	_, err = io.WriteString(file, "<div id=\"myDiv\" style=\"width: 1024px; height: 768px;\"></div>\n")
	_, err = io.WriteString(file, "<script>\n")
	_, err = io.WriteString(file, "	var data = [{\n")
	_, err = io.WriteString(file, xline)
	_, err = io.WriteString(file, yline)
	_, err = io.WriteString(file, "	  type: 'bar'\n")
	_, err = io.WriteString(file, "	}];\n")
	_, err = io.WriteString(file, "	Plotly.newPlot('myDiv', data);\n")
	_, err = io.WriteString(file, "</script>\n")
	_, err = io.WriteString(file, "</body>")

	if err != nil {
		logrus.Warning("There was an error trying to plot the mass distribution")
	}

	// copy to work directory
	sys.CopyFile(outfile, filepath.Base(outfile))

	return nil
}

// Serialize converts the whle structure to a gob file
func (e *Evidence) Serialize() error {

	var err error

	// create a file
	dataFile, err := os.Create(sys.EvBin())
	if err != nil {
		return err
	}

	dataEncoder := gob.NewEncoder(dataFile)
	goberr := dataEncoder.Encode(e)
	if goberr != nil {
		msg := fmt.Sprintf("Cannot save results: %s", goberr)
		return errors.New(msg)
	}
	dataFile.Close()

	return nil
}

// Restore reads philosopher results files and restore the data sctructure
func (e *Evidence) Restore() error {

	file, _ := os.Open(sys.EvBin())

	dec := gob.NewDecoder(file)
	err := dec.Decode(&e)
	if err != nil {
		return errors.New("Could not restore Philosopher result. Please check file path")
	}

	return nil
}

// RestoreWithPath reads philosopher results files and restore the data sctructure
func (e *Evidence) RestoreWithPath(p string) error {

	var path string

	if strings.Contains(p, string(filepath.Separator)) {
		path = fmt.Sprintf("%s%s", p, sys.EvBin())
	} else {
		path = fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvBin())
	}

	file, _ := os.Open(path)

	dec := gob.NewDecoder(file)
	err := dec.Decode(&e)
	if err != nil {
		return errors.New("Could not restore Philosopher result. Please check file path")
	}

	return nil
}
