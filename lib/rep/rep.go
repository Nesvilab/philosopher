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
	"github.com/prvst/cmsl/err"
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
	Index                     uint32
	Spectrum                  string
	Scan                      int
	Peptide                   string
	Protein                   string
	ModifiedPeptide           string
	AlternativeProteins       []string
	AlternativeTargetProteins []string
	ModPositions              []uint16
	AssignedModMasses         []float64
	AssignedMassDiffs         []float64
	AssignedModifications     map[string]uint16
	ObservedModifications     map[string]uint16
	AssumedCharge             uint8
	HitRank                   uint8
	PrecursorNeutralMass      float64
	PrecursorExpMass          float64
	RetentionTime             float64
	CalcNeutralPepMass        float64
	RawMassdiff               float64
	Massdiff                  float64
	LocalizedMassDiff         string
	Probability               float64
	Expectation               float64
	Xcorr                     float64
	DeltaCN                   float64
	DeltaCNStar               float64
	SPScore                   float64
	SPRank                    float64
	Hyperscore                float64
	Nextscore                 float64
	DiscriminantValue         float64
	ModNtermMass              float64
	Intensity                 float64
	Purity                    float64
	Labels                    tmt.Labels
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
	AssignedModifications   map[string]uint16
	ObservedModifications   map[string]uint16
	RetentionTime           string
	Spectra                 map[string]int
	MappedProteins          map[string]uint8
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
	URazorModifiedObservations   int
	URazorUnModifiedObservations int
	URazorAssignedModifications  map[string]uint16
	URazorObservedModifications  map[string]uint16
	TotalLabels                  tmt.Labels
	UniqueLabels                 tmt.Labels
	RazorLabels                  tmt.Labels // Unique + razor
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
	MassBins []MassBin
}

// MassBin represents each bin from the mass distribution
type MassBin struct {
	LowerMass     float64
	HigherRight   float64
	MassCenter    float64
	AverageMass   float64
	CorrectedMass float64
	Modifications []string
	AssignedMods  PSMEvidenceList
	ObservedMods  PSMEvidenceList
}

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
			p.AlternativeTargetProteins = i.AlternativeTargetProteins
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
			p.LocalizedMassDiff = i.LocalizedMassDiff
			p.Probability = i.Probability
			p.Expectation = i.Expectation
			p.Xcorr = i.Xcorr
			p.DeltaCN = i.DeltaCN
			p.SPRank = i.SPRank
			p.Hyperscore = i.Hyperscore
			p.Nextscore = i.Nextscore
			p.DiscriminantValue = i.DiscriminantValue
			p.ModNtermMass = i.ModNtermMass
			p.Intensity = i.Intensity
			p.AssignedModifications = make(map[string]uint16)
			p.ObservedModifications = make(map[string]uint16)

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

	_, err = io.WriteString(file, "Spectrum\tPeptide\tModified Peptide with Assigned Modifications\tCharge\tRetention\tCalculated M/Z\tObserved M/Z\tOriginal Delta Mass\tAdjusted Delta Mass\tExperimental Mass\tPeptide Mass\tXCorr\tDeltaCN\tDeltaCNStar\tSPScore\tSPRank\tExpectation\tHyperscore\tNextscore\tPeptideProphet Probability\tAssigned Modifications\tOberved Modifications\tDelta Mass Localization\tMapped Proteins\tProtein\tAlternative Proteins\n")
	if err != nil {
		logrus.Fatal("Cannot print PSM to file")
	}

	for _, i := range e.PSM {

		var ass []string
		for j := range i.AssignedModifications {
			ass = append(ass, j)
		}

		var obs []string
		for j := range i.ObservedModifications {
			obs = append(obs, j)
		}

		line := fmt.Sprintf("%s\t%s\t%s\t%d\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%e\t%.4f\t%.4f\t%.4f\t%s\t%s\t%s\t%d\t%s\t%s\n",
			i.Spectrum,
			i.Peptide,
			i.ModifiedPeptide,
			i.AssumedCharge,
			i.RetentionTime,
			((i.CalcNeutralPepMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge)),
			((i.PrecursorNeutralMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge)),
			i.RawMassdiff,
			i.Massdiff,
			i.PrecursorNeutralMass,
			i.CalcNeutralPepMass,
			i.Xcorr,
			i.DeltaCN,
			i.DeltaCNStar,
			i.SPScore,
			i.SPRank,
			i.Expectation,
			i.Hyperscore,
			i.Nextscore,
			i.Probability,
			strings.Join(ass, ", "),
			strings.Join(obs, ", "),
			i.LocalizedMassDiff,
			len(i.AlternativeTargetProteins)+1,
			i.Protein,
			strings.Join(i.AlternativeTargetProteins, ", "),
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
	var assignedModMap = make(map[string][]string)
	var observedModMap = make(map[string][]string)
	var err error

	// collapse all psm to protein based on Peptide-level identifications
	for _, i := range e.PSM {
		psmPtMap[i.Spectrum] = append(psmPtMap[i.Spectrum], i.Protein)
		psmPtMap[i.Spectrum] = append(psmPtMap[i.Spectrum], i.AlternativeProteins...)

		// get the list of all assigned modifications
		if len(i.AssignedModifications) > 0 {
			for k := range i.AssignedModifications {
				assignedModMap[i.Spectrum] = append(assignedModMap[i.Spectrum], k)
			}
		}

		// get the list of all observed modifications
		if len(i.ObservedModifications) > 0 {
			for k := range i.ObservedModifications {
				observedModMap[i.Spectrum] = append(observedModMap[i.Spectrum], k)
			}
		}
	}

	for _, i := range ion {
		if !clas.IsDecoyPSM(i, decoyTag) {

			var pr IonEvidence

			pr.Spectra = make(map[string]int)
			pr.MappedProteins = make(map[string]uint8)
			pr.ObservedModifications = make(map[string]uint16)
			pr.AssignedModifications = make(map[string]uint16)

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
				for _, j := range v {
					pr.MappedProteins[j] = 0
				}
			}

			va, oka := assignedModMap[i.Spectrum]
			if oka {
				for _, j := range va {
					pr.AssignedModifications[j] = 0
				}
			}

			vo, oko := observedModMap[i.Spectrum]
			if oko {
				for _, j := range vo {
					pr.ObservedModifications[j] = 0
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
				for k := range j.AssignedModifications {
					e.Ions[i].AssignedModifications[k]++
				}
				for k := range j.ObservedModifications {
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

// PeptideIonReport reports consist on ion reporting
func (e *Evidence) PeptideIonReport() {

	output := fmt.Sprintf("%s%sion.tsv", e.Temp, string(filepath.Separator))

	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Could not create peptide output file")
	}
	defer file.Close()

	_, err = io.WriteString(file, "Peptide Sequence\tM/Z\tCharge\tExperimental Mass\tProbability\tExpectation\tSpectral Count\tUnmodified Observations\tModified Observations\tIntensity\tAssigned Modifications\tObserved Modifications\tIntensity\tMapped Proteins\tProtein IDs\n")
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

	_, err = io.WriteString(file, "Peptide\tCharges\tSpectral Count\tUnmodified Observations\tModified Observations\n")
	if err != nil {
		logrus.Fatal("Cannot create peptide report header")
	}

	for _, i := range e.Peptides {

		var cs []string
		for j := range i.ChargeState {
			cs = append(cs, strconv.Itoa(int(j)))
		}
		sort.Strings(cs)

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
			rep.URazorAssignedModifications = make(map[string]uint16)
			rep.URazorObservedModifications = make(map[string]uint16)

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

								for key, value := range e.Ions[j].AssignedModifications {
									rep.URazorAssignedModifications[key] += value
								}

								for key, value := range e.Ions[j].ObservedModifications {
									rep.URazorObservedModifications[key] += value
								}

							}

							if k.Razor == 1 {
								rep.URazorUnModifiedObservations += e.Ions[j].UnModifiedObservations
								rep.URazorModifiedObservations += e.Ions[j].ModifiedObservations

								rep.URazorPeptideIons[ion] = e.Ions[j]
								rep.TotalNumRazorPeptides++

								for key, value := range e.Ions[j].AssignedModifications {
									rep.URazorAssignedModifications[key] += value
								}

								for key, value := range e.Ions[j].ObservedModifications {
									rep.URazorObservedModifications[key] += value
								}

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

	line := fmt.Sprintf("Group\tSubGroup\tProtein ID\tEntry Name\tLength\tPercent Coverage\tOrganism\tDescription\tProtein Existence\tGenes\tProtein Probability\tTop Peptide Probability\tStripped Peptides\tTotal Peptide Ions\tUnique Peptide Ions\tTotal Spectral Count\tUnique Spectral Count\tRazor Spectral Count\tRazor Unmodified Observations\tRazor Modified Observations\tTotal Intensity\tUnique Intensity\tRazor Intensity\tRazor Assigned Modifications\tRazor Observed Modifications\tIndistinguishable Proteins\n")

	n, err := io.WriteString(file, line)
	if err != nil {
		logrus.Fatal(n, err)
	}

	for _, i := range e.Proteins {

		var ip []string
		for k := range i.IndiProtein {
			ip = append(ip, k)
		}

		var amods []string
		if len(i.URazorAssignedModifications) > 0 {
			for j := range i.URazorAssignedModifications {
				amods = append(amods, j)
			}
		}

		var omods []string
		if len(i.URazorObservedModifications) > 0 {
			for j := range i.URazorObservedModifications {
				omods = append(omods, j)
			}
		}

		// proteins with almost no evidences, and completely shared with decoys are eliminated from the analysis,
		// in most cases proteins with one small peptide shared with a decoy
		//if len(i.TotalPeptideIons) > 0 {

		line = fmt.Sprintf("%d\t%s\t%s\t%s\t%d\t%.2f\t%s\t%s\t%s\t%s\t%.4f\t%.4f\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%6.f\t%6.f\t%6.f\t%s\t%s\t%s\t",
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
			strings.Join(amods, ", "),      // Razor Assigned Modifications
			strings.Join(omods, ", "),      // Razor Observed Modifications
			strings.Join(ip, ", "),         // Indistinguishable Proteins
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

		b.LowerMass = -(amplitude) - (massWindow * binsize) + (float64(i) * binsize)
		b.LowerMass = utils.Round(b.LowerMass, 5, 4)

		b.HigherRight = -(amplitude) + (massWindow * binsize) + (float64(i) * binsize)
		b.HigherRight = utils.Round(b.HigherRight, 5, 4)

		b.MassCenter = -(amplitude) + (float64(i) * binsize)
		b.MassCenter = utils.Round(b.MassCenter, 5, 4)

		bins = append(bins, b)
	}

	// calculate the total number of PSMs per cluster
	var counter int
	for i := range e.PSM {

		for j := range bins {

			// for assigned mods
			for l := range e.PSM[i].AssignedMassDiffs {
				if e.PSM[i].AssignedMassDiffs[l] > bins[j].LowerMass && e.PSM[i].AssignedMassDiffs[l] <= bins[j].HigherRight {
					bins[j].AssignedMods = append(bins[j].AssignedMods, e.PSM[i])
					counter++
					break
				}
			}

			// for delta masses
			if e.PSM[i].Massdiff > bins[j].LowerMass && e.PSM[i].Massdiff <= bins[j].HigherRight {
				bins[j].ObservedMods = append(bins[j].ObservedMods, e.PSM[i])
				counter++
				break
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

		bins[i].AverageMass = utils.Round(bins[i].AverageMass, 5, 4)
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
		bins[i].CorrectedMass = utils.Round(bins[i].CorrectedMass, 5, 4)
	}

	e.Modifications.MassBins = bins
	e.Modifications = modEvi
	e.Modifications.MassBins = bins

	return nil
}

// MapMassDiffToUniMod maps PSMs to modifications based on their mass shifts
func (e *Evidence) MapMassDiffToUniMod() *err.Error {

	// 10 ppm
	var tolerance = 0.01

	u := uni.New()
	u.ProcessUniMOD()

	for _, i := range u.Modifications {

		for j := range e.PSM {

			// for fixed and variable modifications
			for k := range e.PSM[j].AssignedMassDiffs {
				if e.PSM[j].AssignedMassDiffs[k] >= (i.MonoMass-tolerance) && e.PSM[j].AssignedMassDiffs[k] <= (i.MonoMass+tolerance) {
					if !strings.Contains(i.Description, "substitution") {
						fullname := fmt.Sprintf("%.4f:%s (%s)", i.MonoMass, i.Title, i.Description)
						e.PSM[j].AssignedModifications[fullname] = 0
					}
				}
			}

			// for delta masses
			if e.PSM[j].Massdiff >= (i.MonoMass-tolerance) && e.PSM[j].Massdiff <= (i.MonoMass+tolerance) {
				fullName := fmt.Sprintf("%.4f:%s (%s)", i.MonoMass, i.Title, i.Description)
				_, ok := e.PSM[j].AssignedModifications[fullName]
				if !ok {
					e.PSM[j].ObservedModifications[fullName] = 0
				}
			}

		}
	}

	for j := range e.PSM {
		if e.PSM[j].Massdiff != 0 && len(e.PSM[j].ObservedModifications) == 0 {
			e.PSM[j].ObservedModifications["Unknown"] = 0
		}
	}

	return nil
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

	line := fmt.Sprintf("Mass Bin\tPSMs with Assigned Modifications\tPSMs with Observed Modifications\n")

	n, err := io.WriteString(file, line)
	if err != nil {
		logrus.Fatal(n, err)
	}

	for _, i := range e.Modifications.MassBins {

		line = fmt.Sprintf("%.4f\t%d\t%d",
			i.CorrectedMass,
			len(i.AssignedMods),
			len(i.ObservedMods),
		)

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
	var y1var []string
	var y2var []string

	for _, i := range e.Modifications.MassBins {
		xel := fmt.Sprintf("'%.2f',", i.MassCenter)
		xvar = append(xvar, xel)
		y1el := fmt.Sprintf("'%d',", len(i.AssignedMods))
		y1var = append(y1var, y1el)
		y2el := fmt.Sprintf("'%d',", len(i.ObservedMods))
		y2var = append(y2var, y2el)
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

	if err != nil {
		logrus.Warning("There was an error trying to plot the mass distribution")
	}

	// copy to work directory
	sys.CopyFile(outfile, filepath.Base(outfile))

	return nil
}

// // PlotMassHist plots the delta mass histogram
// func (e *Evidence) PlotMassHist() error {
//
// 	outfile := fmt.Sprintf("%s%sdelta-mass.html", e.Temp, string(filepath.Separator))
//
// 	file, err := os.Create(outfile)
// 	if err != nil {
// 		return errors.New("Could not create output for delta mass binning")
// 	}
// 	defer file.Close()
//
// 	var xvar []string
// 	var yvar []string
//
// 	for _, i := range e.Modifications.MassBins {
// 		xel := fmt.Sprintf("'%.2f',", i.MassCenter)
// 		xvar = append(xvar, xel)
// 		yel := fmt.Sprintf("'%d',", len(i.ObservedMods))
// 		yvar = append(yvar, yel)
// 	}
//
// 	xline := fmt.Sprintf("	  x: %s,", xvar)
// 	yline := fmt.Sprintf("	  y: %s,", yvar)
//
// 	io.WriteString(file, "<head>\n")
// 	io.WriteString(file, "  <script src=\"https://cdn.plot.ly/plotly-latest.min.js\"></script>\n")
// 	io.WriteString(file, "</head>\n")
// 	io.WriteString(file, "<body>\n")
// 	io.WriteString(file, "<div id=\"myDiv\" style=\"width: 1024px; height: 768px;\"></div>\n")
// 	io.WriteString(file, "<script>\n")
// 	io.WriteString(file, "	var data = [{\n")
// 	io.WriteString(file, xline)
// 	io.WriteString(file, yline)
// 	io.WriteString(file, "	  type: 'bar'\n")
// 	io.WriteString(file, "	}];\n")
// 	io.WriteString(file, "	Plotly.newPlot('myDiv', data);\n")
// 	io.WriteString(file, "</script>\n")
// 	io.WriteString(file, "</body>")
//
// 	if err != nil {
// 		logrus.Warning("There was an error trying to plot the mass distribution")
// 	}
//
// 	// copy to work directory
// 	sys.CopyFile(outfile, filepath.Base(outfile))
//
// 	return nil
// }

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
