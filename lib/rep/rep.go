package rep

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/prvst/philosopher/lib/bio"
	"github.com/prvst/philosopher/lib/cla"
	"github.com/prvst/philosopher/lib/dat"
	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/id"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/prvst/philosopher/lib/tmt"
	"github.com/prvst/philosopher/lib/uni"
	"github.com/prvst/philosopher/lib/uti"
	"github.com/sirupsen/logrus"
)

// Evidence ...
type Evidence struct {
	Decoys        bool
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
	IonForm               string
	Protein               string
	RazorProtein          string
	ProteinID             string
	GeneName              string
	ModifiedPeptide       string
	MappedProteins        map[string]int
	ModPositions          []string
	AssignedModMasses     []float64
	AssignedMassDiffs     []float64
	AssignedAminoAcid     []string
	AssignedModifications map[string]uint16
	ObservedModifications map[string]uint16
	AssumedCharge         uint8
	HitRank               uint8
	PrecursorNeutralMass  float64
	PrecursorExpMass      float64
	RetentionTime         float64
	CalcNeutralPepMass    float64
	RawMassdiff           float64
	Massdiff              float64
	LocalizedPTMSites     map[string]int
	LocalizedPTMMassDiff  map[string]string
	Probability           float64
	Expectation           float64
	Xcorr                 float64
	DeltaCN               float64
	DeltaCNStar           float64
	SPScore               float64
	SPRank                float64
	Hyperscore            float64
	Nextscore             float64
	DiscriminantValue     float64
	Intensity             float64
	Purity                float64
	IsDecoy               bool
	IsUnique              bool
	IsURazor              bool
	Labels                tmt.Labels
}

// PSMEvidenceList ...
type PSMEvidenceList []PSMEvidence

func (a PSMEvidenceList) Len() int           { return len(a) }
func (a PSMEvidenceList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a PSMEvidenceList) Less(i, j int) bool { return a[i].Spectrum < a[j].Spectrum }

// IonEvidence groups all valid info about peptide ions for reports
type IonEvidence struct {
	Sequence               string
	IonForm                string
	ModifiedSequence       string
	RetentionTime          string
	ChargeState            uint8
	ModifiedObservations   int
	UnModifiedObservations int
	AssignedModifications  map[string]uint16
	ObservedModifications  map[string]uint16
	Spectra                map[string]int
	MappedProteins         map[string]int
	MZ                     float64
	PeptideMass            float64
	PrecursorNeutralMass   float64
	Weight                 float64
	GroupWeight            float64
	Intensity              float64
	Probability            float64
	Expectation            float64
	SummedLabelIntensity   float64
	IsUnique               bool
	IsURazor               bool
	IsDecoy                bool
	Labels                 tmt.Labels
	PhosphoLabels          tmt.Labels
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
	Spectra                map[string]uint8
	MappedProteins         map[string]uint16
	Spc                    int
	Intensity              float64
	ModifiedObservations   int
	UnModifiedObservations int
	IsDecoy                bool
	Labels                 tmt.Labels
	PhosphoLabels          tmt.Labels
}

// PeptideEvidenceList ...
type PeptideEvidenceList []PeptideEvidence

func (a PeptideEvidenceList) Len() int           { return len(a) }
func (a PeptideEvidenceList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a PeptideEvidenceList) Less(i, j int) bool { return a[i].Sequence < a[j].Sequence }

// ProteinEvidence ...
type ProteinEvidence struct {
	OriginalHeader               string
	PartHeader                   string
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
	SupportingSpectra            map[string]int
	IndiProtein                  map[string]uint8
	UniqueStrippedPeptides       int
	TotalPeptideIons             map[string]IonEvidence
	TotalSpC                     int
	UniqueSpC                    int
	URazorSpC                    int // Unique + razor
	TotalIntensity               float64
	UniqueIntensity              float64
	URazorIntensity              float64 // Unique + razor
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
	URazorLabels                 tmt.Labels // Unique + razor
	PhosphoTotalLabels           tmt.Labels
	PhosphoUniqueLabels          tmt.Labels
	PhosphoURazorLabels          tmt.Labels // Unique + razor
}

// ProteinEvidenceList list
type ProteinEvidenceList []ProteinEvidence

func (a ProteinEvidenceList) Len() int           { return len(a) }
func (a ProteinEvidenceList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ProteinEvidenceList) Less(i, j int) bool { return a[i].ProteinGroup < a[j].ProteinGroup }

// CombinedEvidence represents all combined proteins detected
type CombinedEvidence struct {
	GroupNumber            uint32
	SiblingID              string
	ProteinName            string
	ProteinID              string
	IndiProtein            []string
	EntryName              string
	GeneNames              string
	Description            string
	Length                 int
	Names                  []string
	UniqueStrippedPeptides int
	SupportingSpectra      map[string]string
	ProteinProbability     float64
	TopPepProb             float64
	PeptideIons            []id.PeptideIonIdentification
	TotalSpc               map[string]int
	UniqueSpc              map[string]int
	UrazorSpc              map[string]int
	TotalIntensity         map[string]float64
	UniqueIntensity        map[string]float64
	UrazorIntensity        map[string]float64
	TotalLabels            map[string]tmt.Labels
	UniqueLabels           map[string]tmt.Labels
	URazorLabels           map[string]tmt.Labels // Unique + razor
}

// CombinedEvidenceList ...
type CombinedEvidenceList []CombinedEvidence

func (a CombinedEvidenceList) Len() int           { return len(a) }
func (a CombinedEvidenceList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CombinedEvidenceList) Less(i, j int) bool { return a[i].GroupNumber < a[j].GroupNumber }

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

	var self Evidence

	return self
}

// Run is the main entry poit for Report
func Run(m met.Data) met.Data {

	var repo = New()

	err := repo.RestoreGranular()
	if err != nil {
		logrus.Fatal(err.Error())
	}

	if len(repo.Proteins) > 10 {

		logrus.Info("Creating Protein FASTA report")
		repo.ProteinFastaReport()

		if repo.Proteins[0].TotalLabels.Channel1.Intensity > 0 || repo.Proteins[10].TotalLabels.Channel1.Intensity > 0 {
			logrus.Info("Creating Protein TMT report")
			repo.ProteinTMTReport(m.Quantify.LabelNames, m.Quantify.Unique)
		} else {
			logrus.Info("Creating Protein report")
			repo.ProteinReport()
		}

	}

	// verifying if there is any quantification on labels
	if len(m.Quantify.Plex) > 0 {

		annotfile := fmt.Sprintf(".%sannotation.txt", string(filepath.Separator))
		annotfile, _ = filepath.Abs(annotfile)

		labelNames, _ := getLabelNames(annotfile)
		logrus.Info("Creating TMT PSM report")

		if m.SearchEngine == "MSFragger" {
			repo.PSMTMTFraggerReport(labelNames, m.Filter.Tag, m.Filter.Razor)
		} else {
			repo.PSMTMTReport(labelNames, m.Filter.Tag, m.Filter.Razor)
		}

		logrus.Info("Creating TMT peptide report")
		repo.PeptideTMTReport(labelNames)

		logrus.Info("Creating TMT peptide Ion report")
		repo.PeptideIonTMTReport(labelNames)

	} else {

		logrus.Info("Creating PSM report")
		if m.SearchEngine == "MSFragger" {
			repo.PSMFraggerReport(m.Filter.Tag, m.Filter.Razor)
		} else {
			repo.PSMReport(m.Filter.Tag, m.Filter.Razor)
		}

		logrus.Info("Creating peptide report")
		repo.PeptideReport()

		logrus.Info("Creating peptide Ion report")
		repo.PeptideIonReport()

	}

	if len(repo.Modifications.MassBins) > 0 {
		logrus.Info("Creating modification reports")
		repo.ModificationReport()

		if m.PTMProphet.InputFiles != nil || len(m.PTMProphet.InputFiles) > 0 {
			logrus.Info("Creating PTM localization report")
			repo.PSMLocalizationReport(m.Filter.Tag, m.Filter.Razor)
		}

		if len(m.Quantify.Plex) > 0 {
			//if repo.Proteins[0].TotalLabels.Channel1.Intensity > 0 || repo.Proteins[10].TotalLabels.Channel1.Intensity > 0 {
			logrus.Info("Creating TMT phospho protein report")
			repo.PhosphoProteinTMTReport(m.Quantify.LabelNames, m.Quantify.Unique)
		}

		logrus.Info("Plotting mass distribution")
		repo.PlotMassHist()
	}

	return m
}

// AssemblePSMReport ...
func (e *Evidence) AssemblePSMReport(pep id.PepIDList, decoyTag string) error {

	var list PSMEvidenceList

	// collect database information
	var dtb dat.Base
	dtb.Restore()

	var genes = make(map[string]string)
	var ptid = make(map[string]string)
	for _, j := range dtb.Records {
		genes[j.PartHeader] = j.GeneNames
		ptid[j.PartHeader] = j.ID
	}

	for _, i := range pep {

		var p PSMEvidence

		p.Index = i.Index
		p.Spectrum = i.Spectrum
		p.Scan = i.Scan
		p.Peptide = i.Peptide
		p.IonForm = fmt.Sprintf("%s#%d#%.4f", i.Peptide, i.AssumedCharge, i.CalcNeutralPepMass)
		p.Protein = i.Protein
		p.ModifiedPeptide = i.ModifiedPeptide
		p.ModPositions = i.ModPositions
		p.AssignedModMasses = i.AssignedModMasses
		p.AssignedMassDiffs = i.AssignedMassDiffs
		p.AssignedAminoAcid = i.AssignedAminoAcid
		p.AssumedCharge = i.AssumedCharge
		p.HitRank = i.HitRank
		p.PrecursorNeutralMass = i.PrecursorNeutralMass
		p.PrecursorExpMass = i.PrecursorExpMass
		p.RetentionTime = i.RetentionTime
		p.CalcNeutralPepMass = i.CalcNeutralPepMass
		p.RawMassdiff = i.RawMassDiff
		p.Massdiff = i.Massdiff
		p.LocalizedPTMSites = i.LocalizedPTMSites
		p.LocalizedPTMMassDiff = i.LocalizedPTMMassDiff
		p.Probability = i.Probability
		p.Expectation = i.Expectation
		p.Xcorr = i.Xcorr
		p.DeltaCN = i.DeltaCN
		p.SPRank = i.SPRank
		p.Hyperscore = i.Hyperscore
		p.Nextscore = i.Nextscore
		p.DiscriminantValue = i.DiscriminantValue
		p.Intensity = i.Intensity
		p.AssignedModifications = make(map[string]uint16)
		p.ObservedModifications = make(map[string]uint16)
		p.MappedProteins = make(map[string]int)

		for _, j := range i.AlternativeProteins {
			p.MappedProteins[j]++
		}

		gn, ok := genes[i.Protein]
		if ok {
			p.GeneName = gn
		}

		id, ok := ptid[i.Protein]
		if ok {
			p.ProteinID = id
		}

		// is this bservation a decoy ?
		if cla.IsDecoyPSM(i, decoyTag) {
			p.IsDecoy = true
		}

		list = append(list, p)
	}

	sort.Sort(list)
	e.PSM = list

	return nil
}

// PSMReport report all psms from study that passed the FDR filter
func (e *Evidence) PSMReport(decoyTag string, hasRazor bool) {

	output := fmt.Sprintf("%s%spsm.tsv", sys.MetaDir(), string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create report file:", err)
	}
	defer file.Close()

	_, err = io.WriteString(file, "Spectrum\tPeptide\tModified Peptide\tCharge\tRetention\tCalculated M/Z\tObserved M/Z\tOriginal Delta Mass\tAdjusted Delta Mass\tExperimental Mass\tPeptide Mass\tXCorr\tDeltaCN\tDeltaCNStar\tSPScore\tSPRank\tExpectation\tHyperscore\tNextscore\tPeptideProphet Probability\tIntensity\tAssigned Modifications\tObserved Modifications\tNumber of Phospho Sites\tPhospho Site Localization\tIs Unique\tGene\tProtein\tMapped Proteins\n")
	if err != nil {
		logrus.Fatal("Cannot print PSM to file")
	}

	// building the printing set tat may or not contain decoys
	var printSet PSMEvidenceList
	for _, i := range e.PSM {
		if hasRazor == true {

			if i.IsURazor == true {
				if e.Decoys == false {
					if i.IsDecoy == false && len(i.Protein) > 0 && !strings.Contains(i.Protein, decoyTag) {
						printSet = append(printSet, i)
					}
				} else {
					printSet = append(printSet, i)
				}
			}

		} else {

			if e.Decoys == false {
				if i.IsDecoy == false && len(i.Protein) > 0 && !strings.Contains(i.Protein, decoyTag) {
					printSet = append(printSet, i)
				}
			} else {
				printSet = append(printSet, i)
			}

		}
	}

	for _, i := range printSet {

		var assL []string

		for j := 0; j <= len(i.ModPositions)-1; j++ {
			if i.AssignedMassDiffs[j] != 0 && i.AssignedAminoAcid[j] == "n" {
				loc := fmt.Sprintf("%s(%.4f)", i.ModPositions[j], i.AssignedMassDiffs[j])
				assL = append(assL, loc)
			}
		}

		for j := 0; j <= len(i.ModPositions)-1; j++ {
			if i.AssignedMassDiffs[j] != 0 && i.AssignedAminoAcid[j] != "n" && i.AssignedAminoAcid[j] != "c" {
				loc := fmt.Sprintf("%s%s(%.4f)", i.ModPositions[j], i.AssignedAminoAcid[j], i.AssignedMassDiffs[j])
				assL = append(assL, loc)
			}
		}

		for j := 0; j <= len(i.ModPositions)-1; j++ {
			if i.AssignedMassDiffs[j] != 0 && i.AssignedAminoAcid[j] == "c" {
				loc := fmt.Sprintf("%s(%.4f)", i.ModPositions[j], i.AssignedMassDiffs[j])
				assL = append(assL, loc)
			}
		}

		var mappedProteins []string
		for j := range i.MappedProteins {
			if j != i.Protein {
				mappedProteins = append(mappedProteins, j)
			}
		}

		var obs []string
		for j := range i.ObservedModifications {
			obs = append(obs, j)
		}

		//TODO FIX MODS
		line := fmt.Sprintf("%s\t%s\t%s\t%d\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%e\t%.4f\t%.4f\t%.4f\t%.4f\t%s\t%s\t%d\t%s\t%t\t%s\t%s\t%s\n",
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
			i.Intensity,
			strings.Join(assL, ", "),
			strings.Join(obs, ", "),
			i.LocalizedPTMSites["PTMProphet_STY79.9663"],
			i.LocalizedPTMMassDiff["PTMProphet_STY79.9663"],
			i.IsUnique,
			i.GeneName,
			i.Protein,
			strings.Join(mappedProteins, ", "),
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

// PSMTMTReport report all psms with TMT labels from study that passed the FDR filter
func (e *Evidence) PSMTMTReport(labels map[string]string, decoyTag string, hasRazor bool) {

	output := fmt.Sprintf("%s%spsm.tsv", sys.MetaDir(), string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create report file:", err)
	}
	defer file.Close()

	header := "Spectrum\tPeptide\tModified Peptide\tCharge\tRetention\tCalculated M/Z\tObserved M/Z\tOriginal Delta Mass\tAdjusted Delta Mass\tExperimental Mass\tPeptide Mass\tXCorr\tDeltaCN\tDeltaCNStar\tSPScore\tSPRank\tExpectation\tHyperscore\tNextscore\tPeptideProphet Probability\tIntensity\tIs Unique\tIs Used\tAssigned Modifications\tObserved Modifications\tNumber of Phospho Sites\tPhospho Site Localization\tGene\tProtein\tMapped Proteins\tPurity\t126 Abundance\t127N Abundance\t127C Abundance\t128N Abundance\t128C Abundance\t129N Abundance\t129C Abundance\t130N Abundance\t130C Abundance\t131N Abundance\t131C Abundance\n"

	if len(labels) > 0 {
		for k, v := range labels {
			header = strings.Replace(header, k, v, -1)
		}
	}

	_, err = io.WriteString(file, header)
	if err != nil {
		logrus.Fatal("Cannot print PSM to file")
	}

	// building the printing set tat may or not contain decoys
	var printSet PSMEvidenceList
	for _, i := range e.PSM {
		if hasRazor == true {

			if i.IsURazor == true {
				if e.Decoys == false {
					if i.IsDecoy == false && len(i.Protein) > 0 && !strings.Contains(i.Protein, decoyTag) {
						printSet = append(printSet, i)
					}
				} else {
					printSet = append(printSet, i)
				}
			}

		} else {

			if e.Decoys == false {
				if i.IsDecoy == false && len(i.Protein) > 0 && !strings.Contains(i.Protein, decoyTag) {
					printSet = append(printSet, i)
				}
			} else {
				printSet = append(printSet, i)
			}

		}
	}

	// ////////////////////////////////////////
	// ///// fix - NEEDS TO BE REMOVED
	// var dtb dat.Base
	// dtb.Restore()
	// var dbMap = make(map[string]string)
	// for _, j := range dtb.Records {
	// 	dbMap[j.PartHeader] = j.GeneNames
	// }
	// ///// fix - NEEDS TO BE REMOVED
	// ////////////////////////////////////////

	for _, i := range printSet {

		var assL []string

		for j := 0; j <= len(i.ModPositions)-1; j++ {
			if i.AssignedMassDiffs[j] != 0 && i.AssignedAminoAcid[j] == "n" {
				loc := fmt.Sprintf("%s(%.4f)", i.ModPositions[j], i.AssignedMassDiffs[j])
				assL = append(assL, loc)
			}
		}

		for j := 0; j <= len(i.ModPositions)-1; j++ {
			if i.AssignedMassDiffs[j] != 0 && i.AssignedAminoAcid[j] != "n" && i.AssignedAminoAcid[j] != "c" {
				loc := fmt.Sprintf("%s%s(%.4f)", i.ModPositions[j], i.AssignedAminoAcid[j], i.AssignedMassDiffs[j])
				assL = append(assL, loc)
			}
		}

		for j := 0; j <= len(i.ModPositions)-1; j++ {
			if i.AssignedMassDiffs[j] != 0 && i.AssignedAminoAcid[j] == "c" {
				loc := fmt.Sprintf("%s(%.4f)", i.ModPositions[j], i.AssignedMassDiffs[j])
				assL = append(assL, loc)
			}
		}

		var mappedProteins []string
		for j := range i.MappedProteins {
			if j != i.Protein {
				mappedProteins = append(mappedProteins, j)
			}
		}

		var obs []string
		for j := range i.ObservedModifications {
			obs = append(obs, j)
		}

		///// fix - NEEDS TO BE REMOVED
		//geneName := dbMap[i.Protein]
		///////////////

		//TODO FIX MDOS
		line := fmt.Sprintf("%s\t%s\t%s\t%d\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%e\t%.4f\t%.4f\t%.4f\t%.4f\t%t\t%t\t%s\t%s\t%d\t%s\t%s\t%s\t%s\t%.2f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\n",
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
			i.Intensity,
			i.IsUnique,
			i.Labels.IsUsed,
			strings.Join(assL, ", "),
			strings.Join(obs, ", "),
			i.LocalizedPTMSites["PTMProphet_STY79.9663"],
			i.LocalizedPTMMassDiff["PTMProphet_STY79.9663"],
			i.GeneName,
			i.Protein,
			strings.Join(mappedProteins, ", "),
			i.Purity,
			i.Labels.Channel1.Intensity,
			i.Labels.Channel2.Intensity,
			i.Labels.Channel3.Intensity,
			i.Labels.Channel4.Intensity,
			i.Labels.Channel5.Intensity,
			i.Labels.Channel6.Intensity,
			i.Labels.Channel7.Intensity,
			i.Labels.Channel8.Intensity,
			i.Labels.Channel9.Intensity,
			i.Labels.Channel10.Intensity,
			i.Labels.Channel11.Intensity,
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

// PSMFraggerReport report all psms from study that passed the FDR filter
func (e *Evidence) PSMFraggerReport(decoyTag string, hasRazor bool) {

	output := fmt.Sprintf("%s%spsm.tsv", sys.MetaDir(), string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create report file:", err)
	}
	defer file.Close()

	_, err = io.WriteString(file, "Spectrum\tPeptide\tModified Peptide\tCharge\tRetention\tCalculated M/Z\tObserved M/Z\tOriginal Delta Mass\tAdjusted Delta Mass\tExperimental Mass\tPeptide Mass\tExpectation\tHyperscore\tNextscore\tPeptideProphet Probability\tIntensity\tAssigned Modifications\tObserved Modifications\tNumber of Phospho Sites\tPhospho Site Localization\tIs Unique\tGene\tProtein\tMapped Proteins\n")
	if err != nil {
		logrus.Fatal("Cannot print PSM to file")
	}

	// building the printing set tat may or not contain decoys
	var printSet PSMEvidenceList
	for _, i := range e.PSM {
		if hasRazor == true {

			if i.IsURazor == true {
				if e.Decoys == false {
					if i.IsDecoy == false && len(i.Protein) > 0 && !strings.Contains(i.Protein, decoyTag) {
						printSet = append(printSet, i)
					}
				} else {
					printSet = append(printSet, i)
				}
			}

		} else {

			if e.Decoys == false {
				if i.IsDecoy == false && len(i.Protein) > 0 && !strings.Contains(i.Protein, decoyTag) {
					printSet = append(printSet, i)
				}
			} else {
				printSet = append(printSet, i)
			}

		}
	}

	for _, i := range printSet {

		var assL []string

		for j := 0; j <= len(i.ModPositions)-1; j++ {
			if i.AssignedMassDiffs[j] != 0 && i.AssignedAminoAcid[j] == "n" {
				loc := fmt.Sprintf("%s(%.4f)", i.ModPositions[j], i.AssignedMassDiffs[j])
				assL = append(assL, loc)
			}
		}

		for j := 0; j <= len(i.ModPositions)-1; j++ {
			if i.AssignedMassDiffs[j] != 0 && i.AssignedAminoAcid[j] != "n" && i.AssignedAminoAcid[j] != "c" {
				loc := fmt.Sprintf("%s%s(%.4f)", i.ModPositions[j], i.AssignedAminoAcid[j], i.AssignedMassDiffs[j])
				assL = append(assL, loc)
			}
		}

		for j := 0; j <= len(i.ModPositions)-1; j++ {
			if i.AssignedMassDiffs[j] != 0 && i.AssignedAminoAcid[j] == "c" {
				loc := fmt.Sprintf("%s(%.4f)", i.ModPositions[j], i.AssignedMassDiffs[j])
				assL = append(assL, loc)
			}
		}

		var mappedProteins []string
		for j := range i.MappedProteins {
			if j != i.Protein {
				mappedProteins = append(mappedProteins, j)
			}
		}

		var obs []string
		for j := range i.ObservedModifications {
			obs = append(obs, j)
		}

		//TODO FIX MODS
		line := fmt.Sprintf("%s\t%s\t%s\t%d\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%e\t%.4f\t%.4f\t%.4f\t%.4f\t%s\t%s\t%d\t%s\t%t\t%s\t%s\t%s\n",
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
			i.Expectation,
			i.Hyperscore,
			i.Nextscore,
			i.Probability,
			i.Intensity,
			strings.Join(assL, ", "),
			strings.Join(obs, ", "),
			i.LocalizedPTMSites["PTMProphet_STY79.9663"],
			i.LocalizedPTMMassDiff["PTMProphet_STY79.9663"],
			i.IsUnique,
			i.GeneName,
			i.Protein,
			strings.Join(mappedProteins, ", "),
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

// PSMTMTFraggerReport report all psms with TMT labels from study that passed the FDR filter
func (e *Evidence) PSMTMTFraggerReport(labels map[string]string, decoyTag string, hasRazor bool) {

	output := fmt.Sprintf("%s%spsm.tsv", sys.MetaDir(), string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create report file:", err)
	}
	defer file.Close()

	header := "Spectrum\tPeptide\tModified Peptide\tCharge\tRetention\tCalculated M/Z\tObserved M/Z\tOriginal Delta Mass\tAdjusted Delta Mass\tExperimental Mass\tPeptide Mass\tExpectation\tHyperscore\tNextscore\tPeptideProphet Probability\tIntensity\tIs Unique\tIs Used\tAssigned Modifications\tObserved Modifications\tNumber of Phospho Sites\tPhospho Site Localization\tGene\tProtein\tMapped Proteins\tPurity\t126 Abundance\t127N Abundance\t127C Abundance\t128N Abundance\t128C Abundance\t129N Abundance\t129C Abundance\t130N Abundance\t130C Abundance\t131N Abundance\t131C Abundance\n"

	if len(labels) > 0 {
		for k, v := range labels {
			header = strings.Replace(header, k, v, -1)
		}
	}

	_, err = io.WriteString(file, header)
	if err != nil {
		logrus.Fatal("Cannot print PSM to file")
	}

	// building the printing set tat may or not contain decoys
	var printSet PSMEvidenceList
	for _, i := range e.PSM {
		if hasRazor == true {

			if i.IsURazor == true {
				if e.Decoys == false {
					if i.IsDecoy == false && len(i.Protein) > 0 && !strings.Contains(i.Protein, decoyTag) {
						printSet = append(printSet, i)
					}
				} else {
					printSet = append(printSet, i)
				}
			}

		} else {

			if e.Decoys == false {
				if i.IsDecoy == false && len(i.Protein) > 0 && !strings.Contains(i.Protein, decoyTag) {
					printSet = append(printSet, i)
				}
			} else {
				printSet = append(printSet, i)
			}

		}
	}

	// ////////////////////////////////////////
	// ///// fix - NEEDS TO BE REMOVED
	// var dtb dat.Base
	// dtb.Restore()
	// var dbMap = make(map[string]string)
	// for _, j := range dtb.Records {
	// 	dbMap[j.PartHeader] = j.GeneNames
	// }
	// ///// fix - NEEDS TO BE REMOVED
	// ////////////////////////////////////////

	for _, i := range printSet {

		var assL []string

		for j := 0; j <= len(i.ModPositions)-1; j++ {
			if i.AssignedMassDiffs[j] != 0 && i.AssignedAminoAcid[j] == "n" {
				loc := fmt.Sprintf("%s(%.4f)", i.ModPositions[j], i.AssignedMassDiffs[j])
				assL = append(assL, loc)
			}
		}

		for j := 0; j <= len(i.ModPositions)-1; j++ {
			if i.AssignedMassDiffs[j] != 0 && i.AssignedAminoAcid[j] != "n" && i.AssignedAminoAcid[j] != "c" {
				loc := fmt.Sprintf("%s%s(%.4f)", i.ModPositions[j], i.AssignedAminoAcid[j], i.AssignedMassDiffs[j])
				assL = append(assL, loc)
			}
		}

		for j := 0; j <= len(i.ModPositions)-1; j++ {
			if i.AssignedMassDiffs[j] != 0 && i.AssignedAminoAcid[j] == "c" {
				loc := fmt.Sprintf("%s(%.4f)", i.ModPositions[j], i.AssignedMassDiffs[j])
				assL = append(assL, loc)
			}
		}

		var mappedProteins []string
		for j := range i.MappedProteins {
			if j != i.Protein {
				mappedProteins = append(mappedProteins, j)
			}
		}

		var obs []string
		for j := range i.ObservedModifications {
			obs = append(obs, j)
		}

		///// fix - NEEDS TO BE REMOVED
		//geneName := dbMap[i.Protein]
		///////////////

		// TODO FIZ MODS
		line := fmt.Sprintf("%s\t%s\t%s\t%d\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%e\t%.4f\t%.4f\t%.4f\t%.4f\t%t\t%t\t%s\t%s\t%d\t%s\t%s\t%s\t%s\t%.2f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\n",
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
			i.Expectation,
			i.Hyperscore,
			i.Nextscore,
			i.Probability,
			i.Intensity,
			i.IsUnique,
			i.Labels.IsUsed,
			strings.Join(assL, ", "),
			strings.Join(obs, ", "),
			i.LocalizedPTMSites["PTMProphet_STY79.9663"],
			i.LocalizedPTMMassDiff["PTMProphet_STY79.9663"],
			i.GeneName,
			i.Protein,
			strings.Join(mappedProteins, ", "),
			i.Purity,
			i.Labels.Channel1.Intensity,
			i.Labels.Channel2.Intensity,
			i.Labels.Channel3.Intensity,
			i.Labels.Channel4.Intensity,
			i.Labels.Channel5.Intensity,
			i.Labels.Channel6.Intensity,
			i.Labels.Channel7.Intensity,
			i.Labels.Channel8.Intensity,
			i.Labels.Channel9.Intensity,
			i.Labels.Channel10.Intensity,
			i.Labels.Channel11.Intensity,
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

// PSMLocalizationReport report ptm localization based on PTMProphet outputs
func (e *Evidence) PSMLocalizationReport(decoyTag string, hasRazor bool) {

	output := fmt.Sprintf("%s%slocalization.tsv", sys.MetaDir(), string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create report file:", err)
	}
	defer file.Close()

	_, err = io.WriteString(file, "Spectrum\tPeptide\tModified Peptide\tCharge\tRetention\tModification\tNumber of Sites\tObserved Mass Localization\n")
	if err != nil {
		logrus.Fatal("Cannot print PSM to file")
	}

	// building the printing set tat may or not contain decoys
	var printSet PSMEvidenceList
	for _, i := range e.PSM {
		if hasRazor == true {

			if i.IsURazor == true {
				if e.Decoys == false {
					if i.IsDecoy == false && len(i.Protein) > 0 && !strings.Contains(i.Protein, decoyTag) {
						printSet = append(printSet, i)
					}
				} else {
					printSet = append(printSet, i)
				}
			}

		} else {

			if e.Decoys == false {
				if i.IsDecoy == false && len(i.Protein) > 0 && !strings.Contains(i.Protein, decoyTag) {
					printSet = append(printSet, i)
				}
			} else {
				printSet = append(printSet, i)
			}

		}
	}

	for _, i := range printSet {
		for j := range i.LocalizedPTMMassDiff {
			line := fmt.Sprintf("%s\t%s\t%s\t%d\t%.4f\t%s\t%d\t%s\n",
				i.Spectrum,
				i.Peptide,
				i.ModifiedPeptide,
				i.AssumedCharge,
				i.RetentionTime,
				j,
				i.LocalizedPTMSites[j],
				i.LocalizedPTMMassDiff[j],
			)
			_, err = io.WriteString(file, line)
			if err != nil {
				logrus.Fatal("Cannot print PSM to file")
			}
		}
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}

// AssembleIonReport reports consist on ion reporting
func (e *Evidence) AssembleIonReport(ion id.PepIDList, decoyTag string) error {

	var list IonEvidenceList
	var psmPtMap = make(map[string][]string)
	var psmIonMap = make(map[string][]string)
	var assignedModMap = make(map[string][]string)
	var observedModMap = make(map[string][]string)
	var err error

	// collapse all psm to protein based on Peptide-level identifications
	for _, i := range e.PSM {

		psmIonMap[i.IonForm] = append(psmIonMap[i.IonForm], i.Spectrum)

		psmPtMap[i.Spectrum] = append(psmPtMap[i.Spectrum], i.Protein)
		//psmPtMap[i.Spectrum] = append(psmPtMap[i.Spectrum], i.AlternativeProteins...)

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
		var pr IonEvidence

		pr.IonForm = fmt.Sprintf("%s#%d#%.4f", i.Peptide, i.AssumedCharge, i.CalcNeutralPepMass)

		pr.Spectra = make(map[string]int)
		pr.MappedProteins = make(map[string]int)
		pr.ObservedModifications = make(map[string]uint16)
		pr.AssignedModifications = make(map[string]uint16)

		v, ok := psmIonMap[pr.IonForm]
		if ok {
			for _, j := range v {
				pr.Spectra[j]++
			}
		}

		pr.Sequence = i.Peptide
		pr.ModifiedSequence = i.ModifiedPeptide
		pr.MZ = uti.Round(((i.CalcNeutralPepMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge)), 5, 4)
		pr.ChargeState = i.AssumedCharge
		pr.PeptideMass = i.CalcNeutralPepMass
		pr.PrecursorNeutralMass = i.PrecursorNeutralMass
		pr.Probability = i.Probability
		pr.Expectation = i.Expectation

		pr.MappedProteins[i.Protein] = 0

		// get he list of indi proteins from pepXML data
		v, ok = psmPtMap[i.Spectrum]
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
		} else {
			pr.UnModifiedObservations++
		}

		// is this bservation a decoy ?
		if cla.IsDecoyPSM(i, decoyTag) {
			pr.IsDecoy = true
		}

		list = append(list, pr)
	}

	sort.Sort(list)
	e.Ions = list

	return err
}

// UpdateIonStatus pushes back to ion and psm evideces the uniqueness and razorness status of each peptide and ion
func (e *Evidence) UpdateIonStatus() {

	var uniqueMap = make(map[string]bool)
	var urazorMap = make(map[string]string)
	var ptMap = make(map[string]string)
	//var ptIndiMap = make(map[string][]string)

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
			//e.PSM[i].AlternativeProteins = ptIndiMap[v]
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

	for i := range e.Ions {

		v1, ok1 := UnModIons[e.Ions[i].IonForm]
		if ok1 {
			e.Ions[i].UnModifiedObservations = v1
		}

		v2, ok2 := ModIons[e.Ions[i].IonForm]
		if ok2 {
			e.Ions[i].ModifiedObservations = v2
		}

	}

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

	for i := range e.PSM {
		e.PSM[i].GeneName = dbMap[e.PSM[i].Protein]
	}

	return
}

// PeptideMap struct
type PeptideMap struct {
	Sequence     string
	RazorProtein string
	Proteins     map[string]int
}

// UpdateMappedProteins updates the list of mapped proteins on the data structures
func (e *Evidence) UpdateMappedProteins() {

	var list []PeptideMap
	var checkup = make(map[string]int)

	for _, i := range e.Proteins {
		for _, v := range i.TotalPeptideIons {

			_, ok := checkup[v.Sequence]
			if !ok {
				var pm PeptideMap

				pm.Sequence = v.Sequence
				pm.Proteins = v.MappedProteins
				pm.Proteins[i.PartHeader] = 0

				if v.IsURazor == true {
					pm.RazorProtein = i.PartHeader
				}

				list = append(list, pm)
				checkup[v.Sequence] = 0
			}

		}
	}

	for i := range e.PSM {
		for _, j := range list {
			if e.PSM[i].Peptide == j.Sequence {

				for k := range j.Proteins {
					e.PSM[i].MappedProteins[k]++
				}

				//e.PSM[i].MappedProteins = j.Proteins
				if len(e.PSM[i].RazorProtein) < 1 {
					e.PSM[i].RazorProtein = j.RazorProtein
					e.PSM[i].IsURazor = true
				}
				break
			}
		}

	}

	return
}

// UpdateIonAssignedAndObservedMods collects all Assigned and Observed modifications from
// individual PSM and assign them to ions
func (e *Evidence) UpdateIonAssignedAndObservedMods() {

	for i := range e.Ions {

		for _, j := range e.PSM {

			if e.Ions[i].IonForm == j.IonForm {
				for k := range j.AssignedModifications {
					e.Ions[i].AssignedModifications[k]++
				}
				for k := range j.ObservedModifications {
					e.Ions[i].ObservedModifications[k]++
				}

				break
			}
		}
	}

	return
}

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

// UpdateRecoveredPSMs brings the updated protein and mapped protein info to the recovered PSMs
// func (e *Evidence) UpdateRecoveredPSMs() {
//
// 	var seqProtein = make(map[string]string)
// 	var seqMappedProteins = make(map[string][]string)
//
// 	for _, i := range e.PSM {
// 		if i.Probability >= 0.9 {
// 			seqProtein[i.Peptide] = i.Protein
// 			for j := range i.MappedProteins {
// 				seqMappedProteins[i.Peptide] = append(seqMappedProteins[i.Peptide], j)
// 			}
// 		}
// 	}
//
// 	for i := range e.PSM {
// 		v, ok := seqProtein[e.PSM[i].Peptide]
// 		if ok && e.PSM[i].Probability < 0.9 {
//
// 			if e.PSM[i].Protein != v {
// 				e.PSM[i].Protein = v
// 				for _, j := range seqMappedProteins[e.PSM[i].Peptide] {
// 					e.PSM[i].MappedProteins[j] = 0
// 				}
// 			}
//
// 		}
// 	}
//
// 	return
// }

// PeptideIonReport reports consist on ion reporting
func (e *Evidence) PeptideIonReport() {

	output := fmt.Sprintf("%s%sion.tsv", sys.MetaDir(), string(filepath.Separator))

	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Could not create peptide output file")
	}
	defer file.Close()

	_, err = io.WriteString(file, "Peptide Sequence\tModified Sequence\tM/Z\tCharge\tExperimental Mass\tProbability\tExpectation\tSpectral Count\tUnmodified Observations\tModified Observations\tIntensity\tAssigned Modifications\tObserved Modifications\tIntensity\tMapped Proteins\tProtein IDs\n")
	if err != nil {
		logrus.Fatal("Cannot create peptide ion report header")
	}

	// building the printing set tat may or not contain decoys
	var printSet IonEvidenceList
	for _, i := range e.Ions {
		if e.Decoys == false {
			if i.IsDecoy == false {
				printSet = append(printSet, i)
			}
		} else {
			printSet = append(printSet, i)
		}
	}

	// peptides with no mapped poteins are related to contaminants
	// and reverse sequences. They are dificult to clean because
	// in some cases they are shared between a match decoy and a target,
	// so they stay on the lists but cannot be mapped back to the
	// original proteins. These cases should be rare to find.
	for _, i := range printSet {

		var pts []string
		//var ipts []string

		if len(i.MappedProteins) > 0 {

			//if len(e.Proteins) > 1 {

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

			line := fmt.Sprintf("%s\t%s\t%.4f\t%d\t%.4f\t%.4f\t%.4f\t%d\t%d\t%d\t%.4f\t%s\t%s\t%.4f\t%d\t%s\n",
				i.Sequence,
				i.ModifiedSequence,
				i.MZ,
				i.ChargeState,
				i.PeptideMass,
				i.Probability,
				i.Expectation,
				len(i.Spectra),
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
			//}
		}
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}

// PeptideIonTMTReport reports the ion table with TMT quantification
func (e *Evidence) PeptideIonTMTReport(labels map[string]string) {

	output := fmt.Sprintf("%s%sion.tsv", sys.MetaDir(), string(filepath.Separator))

	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Could not create peptide output file")
	}
	defer file.Close()

	header := "Peptide Sequence\tModified Sequence\tM/Z\tCharge\tExperimental Mass\tProbability\tExpectation\tSpectral Count\tUnmodified Observations\tModified Observations\tIntensity\tAssigned Modifications\tObserved Modifications\tIntensity\tMapped Proteins\tProtein IDs\t126 Abundance\t127N Abundance\t127C Abundance\t128N Abundance\t128C Abundance\t129N Abundance\t129C Abundance\t130N Abundance\t130C Abundance\t131N Abundance\t131C Abundance\n"

	if len(labels) > 0 {
		for k, v := range labels {
			header = strings.Replace(header, k, v, -1)
		}
	}

	_, err = io.WriteString(file, header)
	if err != nil {
		logrus.Fatal("Cannot create peptide ion report header")
	}

	// building the printing set tat may or not contain decoys
	var printSet IonEvidenceList
	for _, i := range e.Ions {
		if e.Decoys == false {
			if i.IsDecoy == false {
				printSet = append(printSet, i)
			}
		} else {
			printSet = append(printSet, i)
		}
	}

	// peptides with no mapped poteins are related to contaminants
	// and reverse sequences. They are dificult to clean because
	// in some cases they are shared between a match decoy and a target,
	// so they stay on the lists but cannot be mapped back to the
	// original proteins. These cases should be rare to find.
	for _, i := range printSet {

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

				line := fmt.Sprintf("%s\t%s\t%.4f\t%d\t%.4f\t%.4f\t%.4f\t%d\t%d\t%d\t%.4f\t%s\t%s\t%.4f\t%d\t%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\n",
					i.Sequence,
					i.ModifiedSequence,
					i.MZ,
					i.ChargeState,
					i.PeptideMass,
					i.Probability,
					i.Expectation,
					len(i.Spectra),
					i.UnModifiedObservations,
					i.ModifiedObservations,
					i.Intensity,
					strings.Join(amods, ", "),
					strings.Join(omods, ", "),
					i.Intensity,
					len(i.MappedProteins),
					strings.Join(pts, ", "),
					i.Labels.Channel1.Intensity,
					i.Labels.Channel2.Intensity,
					i.Labels.Channel3.Intensity,
					i.Labels.Channel4.Intensity,
					i.Labels.Channel5.Intensity,
					i.Labels.Channel6.Intensity,
					i.Labels.Channel7.Intensity,
					i.Labels.Channel8.Intensity,
					i.Labels.Channel9.Intensity,
					i.Labels.Channel10.Intensity,
					i.Labels.Channel11.Intensity,
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
func (e *Evidence) AssemblePeptideReport(pep id.PepIDList, decoyTag string) error {

	var list PeptideEvidenceList
	var pepSeqMap = make(map[string]bool) //is this a decoy
	var pepCSMap = make(map[string][]uint8)
	var pepInt = make(map[string]float64)
	var spectra = make(map[string][]string)
	var err error

	for _, i := range pep {
		if !cla.IsDecoyPSM(i, decoyTag) {
			pepSeqMap[i.Peptide] = false
		} else {
			pepSeqMap[i.Peptide] = true
		}
		pepInt[i.Peptide] = 0
	}

	for _, i := range e.PSM {
		_, ok := pepSeqMap[i.Peptide]
		if ok {

			pepCSMap[i.Peptide] = append(pepCSMap[i.Peptide], i.AssumedCharge)
			spectra[i.Peptide] = append(spectra[i.Peptide], i.Spectrum)

			if i.Intensity > pepInt[i.Peptide] {
				pepInt[i.Peptide] = i.Intensity
			}
		}
	}

	for k, v := range pepSeqMap {

		var pep PeptideEvidence
		pep.Spectra = make(map[string]uint8)
		pep.ChargeState = make(map[uint8]uint8)
		pep.MappedProteins = make(map[string]uint16)
		pep.Sequence = k

		for _, i := range spectra[k] {
			pep.Spectra[i] = 0
		}

		for _, i := range pepCSMap[k] {
			pep.ChargeState[i] = 0
		}
		pep.Spc = len(spectra[k])
		pep.Intensity = pepInt[k]

		// is this a decoy ?
		pep.IsDecoy = v

		list = append(list, pep)
	}

	sort.Sort(list)
	e.Peptides = list

	return err
}

// PeptideReport reports consist on ion reporting
func (e *Evidence) PeptideReport() {

	output := fmt.Sprintf("%s%speptide.tsv", sys.MetaDir(), string(filepath.Separator))

	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Could not create peptide output file")
	}
	defer file.Close()

	_, err = io.WriteString(file, "Peptide\tCharges\tSpectral Count\tUnmodified Observations\tModified Observations\n")
	if err != nil {
		logrus.Fatal("Cannot create peptide report header")
	}

	// building the printing set tat may or not contain decoys
	var printSet PeptideEvidenceList
	for _, i := range e.Peptides {
		if e.Decoys == false {
			if i.IsDecoy == false {
				printSet = append(printSet, i)
			}
		} else {
			printSet = append(printSet, i)
		}
	}

	for _, i := range printSet {

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

// PeptideTMTReport reports consist on ion reporting
func (e *Evidence) PeptideTMTReport(labels map[string]string) {

	output := fmt.Sprintf("%s%speptide.tsv", sys.MetaDir(), string(filepath.Separator))

	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Could not create peptide output file")
	}
	defer file.Close()

	header := "Peptide\tCharges\tSpectral Count\tUnmodified Observations\tModified Observations\t126 Abundance\t127N Abundance\t127C Abundance\t128N Abundance\t128C Abundance\t129N Abundance\t129C Abundance\t130N Abundance\t130C Abundance\t131N Abundance\t131C Abundance\n"

	if len(labels) > 0 {
		for k, v := range labels {
			header = strings.Replace(header, k, v, -1)
		}
	}

	_, err = io.WriteString(file, header)
	if err != nil {
		logrus.Fatal("Cannot create peptide report header")
	}

	// building the printing set tat may or not contain decoys
	var printSet PeptideEvidenceList
	for _, i := range e.Peptides {
		if e.Decoys == false {
			if i.IsDecoy == false {
				printSet = append(printSet, i)
			}
		} else {
			printSet = append(printSet, i)
		}
	}

	for _, i := range printSet {

		var cs []string
		for j := range i.ChargeState {
			cs = append(cs, strconv.Itoa(int(j)))
		}
		sort.Strings(cs)

		line := fmt.Sprintf("%s\t%s\t%d\t%d\t%d\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\n",
			i.Sequence,
			strings.Join(cs, ", "),
			i.Spc,
			i.UnModifiedObservations,
			i.ModifiedObservations,
			i.Labels.Channel1.Intensity,
			i.Labels.Channel2.Intensity,
			i.Labels.Channel3.Intensity,
			i.Labels.Channel4.Intensity,
			i.Labels.Channel5.Intensity,
			i.Labels.Channel6.Intensity,
			i.Labels.Channel7.Intensity,
			i.Labels.Channel8.Intensity,
			i.Labels.Channel9.Intensity,
			i.Labels.Channel10.Intensity,
			i.Labels.Channel11.Intensity,
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

// AssembleProteinReport ...
func (e *Evidence) AssembleProteinReport(pro id.ProtIDList, decoyTag string) error {

	var list ProteinEvidenceList
	var err error

	var evidenceIons = make(map[string]IonEvidence)
	for _, i := range e.Ions {
		evidenceIons[i.IonForm] = i
	}

	for _, i := range pro {

		var rep ProteinEvidence

		rep.SupportingSpectra = make(map[string]int)
		rep.TotalPeptideIons = make(map[string]IonEvidence)
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

				ref.IsUnique = k.IsUnique
				if k.Razor == 1 {
					ref.IsURazor = true
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

				ref.IsUnique = k.IsUnique
				if k.Razor == 1 {
					ref.IsURazor = true
				}

				rep.TotalPeptideIons[ion] = ref

			}

			if k.IsUnique == true {
				rep.URazorUnModifiedObservations += evidenceIons[ion].UnModifiedObservations
				rep.URazorModifiedObservations += evidenceIons[ion].ModifiedObservations

				for key, value := range evidenceIons[ion].AssignedModifications {
					rep.URazorAssignedModifications[key] += value
				}

				for key, value := range evidenceIons[ion].ObservedModifications {
					rep.URazorObservedModifications[key] += value
				}
			}

			if k.Razor == 1 {
				rep.URazorUnModifiedObservations += evidenceIons[ion].UnModifiedObservations
				rep.URazorModifiedObservations += evidenceIons[ion].ModifiedObservations

				for key, value := range evidenceIons[ion].AssignedModifications {
					rep.URazorAssignedModifications[key] += value
				}

				for key, value := range evidenceIons[ion].ObservedModifications {
					rep.URazorObservedModifications[key] += value
				}
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
func (e *Evidence) ProteinReport() {

	// create result file
	output := fmt.Sprintf("%s%sreport.tsv", sys.MetaDir(), string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create report file:", err)
	}
	defer file.Close()

	line := fmt.Sprintf("Group\tSubGroup\tProtein ID\tEntry Name\tLength\tPercent Coverage\tOrganism\tDescription\tProtein Existence\tGenes\tProtein Probability\tTop Peptide Probability\tStripped Peptides\tTotal Peptide Ions\tUnique Peptide Ions\tRazor Peptide Ions\tTotal Spectral Count\tUnique Spectral Count\tRazor Spectral Count\tTotal Intensity\tUnique Intensity\tRazor Intensity\tRazor Assigned Modifications\tRazor Observed Modifications\tIndistinguishable Proteins\n")

	n, err := io.WriteString(file, line)
	if err != nil {
		logrus.Fatal(n, err)
	}

	// building the printing set tat may or not contain decoys
	var printSet ProteinEvidenceList
	for _, i := range e.Proteins {
		if len(i.TotalPeptideIons) > 0 && i.TotalSpC > 0 {
			if e.Decoys == false {
				if i.IsDecoy == false {
					printSet = append(printSet, i)
				}
			} else {
				printSet = append(printSet, i)
			}
		}
	}

	for _, i := range printSet {

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

		// proteins with almost no evidences, and completely shared with decoys are eliminated from the analysis,
		// in most cases proteins with one small peptide shared with a decoy
		//if len(i.TotalPeptideIons) > 0 {

		line = fmt.Sprintf("%d\t%s\t%s\t%s\t%d\t%.2f\t%s\t%s\t%s\t%s\t%.4f\t%.4f\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%6.f\t%6.f\t%6.f\t%s\t%s\t%s\t",
			i.ProteinGroup,            // Group
			i.ProteinSubGroup,         // SubGroup
			i.ProteinID,               // Protein ID
			i.EntryName,               // Entry Name
			i.Length,                  // Length
			i.Coverage,                // Percent Coverage
			i.Organism,                // Organism
			i.Description,             // Description
			i.ProteinExistence,        // Protein Existence
			i.GeneNames,               // Genes
			i.Probability,             // Protein Probability
			i.TopPepProb,              // Top Peptide Probability
			i.UniqueStrippedPeptides,  // Stripped Peptides
			len(i.TotalPeptideIons),   // Total Peptide Ions
			uniqIons,                  // Unique Peptide Ions
			urazorIons,                // Razor Peptide Ions
			i.TotalSpC,                // Total Spectral Count
			i.UniqueSpC,               // Unique Spectral Count
			i.URazorSpC,               // Razor Spectral Count
			i.TotalIntensity,          // Total Intensity
			i.UniqueIntensity,         // Unique Intensity
			i.URazorIntensity,         // Razor Intensity
			strings.Join(amods, ", "), // Razor Assigned Modifications
			strings.Join(omods, ", "), // Razor Observed Modifications
			strings.Join(ip, ", "),    // Indistinguishable Proteins
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
func (e *Evidence) ProteinTMTReport(labels map[string]string, uniqueOnly bool) {

	// create result file
	output := fmt.Sprintf("%s%sreport.tsv", sys.MetaDir(), string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create report file:", err)
	}
	defer file.Close()

	line := fmt.Sprintf("Group\tSubGroup\tProtein ID\tEntry Name\tLength\tPercent Coverage\tDescription\tProtein Existence\tGenes\tProtein Probability\tTop Peptide Probability\tUnique Stripped Peptides\tTotal Peptide Ions\tUnique Peptide Ions\tRazor Peptides Ions\tTotal Spectral Count\tUnique Spectral Count\tRazor Spectral Count\tTotal Intensity\tUnique Intensity\tRazor Intensity\t126 Abundance\t127N Abundance\t127C Abundance\t128N Abundance\t128C Abundance\t129N Abundance\t129C Abundance\t130N Abundance\t130C Abundance\t131N Abundance\t131C Abundance\tIndistinguishableProteins\n")

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
		if e.Decoys == false {
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
			line = fmt.Sprintf("%d\t%s\t%s\t%s\t%d\t%.2f\t%s\t%s\t%s\t%.4f\t%.4f\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%6.f\t%6.f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%s\n",
				i.ProteinGroup,           // Group
				i.ProteinSubGroup,        // SubGroup
				i.ProteinID,              // Protein ID
				i.EntryName,              // EntryName
				i.Length,                 // Length
				i.Coverage,               // Percent Coverage
				i.Description,            // Description
				i.ProteinExistence,       // Protein Existence
				i.GeneNames,              // Genes
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
				strings.Join(ip, ", "))

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
func (e *Evidence) PhosphoProteinTMTReport(labels map[string]string, uniqueOnly bool) {

	// create result file
	output := fmt.Sprintf("%s%sphospho_report.tsv", sys.MetaDir(), string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create report file:", err)
	}
	defer file.Close()

	line := fmt.Sprintf("Group\tSubGroup\tProtein ID\tEntry Name\tLength\tPercent Coverage\tDescription\tProtein Existence\tGenes\tProtein Probability\tTop Peptide Probability\tUnique Stripped Peptides\tTotal Peptide Ions\tUnique Peptide Ions\tRazor Peptides Ions\tTotal Spectral Count\tUnique Spectral Count\tRazor Spectral Count\tTotal Intensity\tUnique Intensity\tRazor Intensity\t126 Abundance\t127N Abundance\t127C Abundance\t128N Abundance\t128C Abundance\t129N Abundance\t129C Abundance\t130N Abundance\t130C Abundance\t131N Abundancet\t131C Abundance\tIndistinguishableProteins\n")

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
		if e.Decoys == false {
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
			line = fmt.Sprintf("%d\t%s\t%s\t%s\t%d\t%.2f\t%s\t%s\t%s\t%.4f\t%.4f\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%6.f\t%6.f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%s\n",
				i.ProteinGroup,           // Group
				i.ProteinSubGroup,        // SubGroup
				i.ProteinID,              // Protein ID
				i.EntryName,              // EntryName
				i.Length,                 // Length
				i.Coverage,               // Percent Coverage
				i.Description,            // Description
				i.ProteinExistence,       // Protein Existence
				i.GeneNames,              // Genes
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
				strings.Join(ip, ", "))

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
func (e *Evidence) ProteinFastaReport() error {

	output := fmt.Sprintf("%s%sproteins.fas", sys.MetaDir(), string(filepath.Separator))

	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Could not create output file")
	}
	defer file.Close()

	// building the printing set tat may or not contain decoys
	var printSet ProteinEvidenceList
	for _, i := range e.Proteins {
		if e.Decoys == false {
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
		b.LowerMass = uti.Round(b.LowerMass, 5, 4)

		b.HigherRight = -(amplitude) + (massWindow * binsize) + (float64(i) * binsize)
		b.HigherRight = uti.Round(b.HigherRight, 5, 4)

		b.MassCenter = -(amplitude) + (float64(i) * binsize)
		b.MassCenter = uti.Round(b.MassCenter, 5, 4)

		bins = append(bins, b)
	}

	// calculate the total number of PSMs per cluster
	for i := range e.PSM {

		// the checklist will not allow the same PSM to be added multiple times to the
		// same bin in case multiple identical mods are present in te sequence
		var assignChecklist = make(map[float64]uint8)
		var obsChecklist = make(map[float64]uint8)

		for j := range bins {

			// for assigned mods
			// 0 here means something that doest not map to the pepXML header
			// like multiple mods on n-term
			for _, l := range e.PSM[i].AssignedMassDiffs {

				if l > bins[j].LowerMass && l <= bins[j].HigherRight && l != 0 {
					_, ok := assignChecklist[l]
					if !ok {
						bins[j].AssignedMods = append(bins[j].AssignedMods, e.PSM[i])
						assignChecklist[l] = 0
					}
				}
			}

			// for delta masses
			if e.PSM[i].Massdiff > bins[j].LowerMass && e.PSM[i].Massdiff <= bins[j].HigherRight {
				_, ok := obsChecklist[e.PSM[i].Massdiff]
				if !ok {
					bins[j].ObservedMods = append(bins[j].ObservedMods, e.PSM[i])
					obsChecklist[e.PSM[i].Massdiff] = 0
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

	//e.Modifications = modEvi
	//e.Modifications.MassBins = bins

	modEvi.MassBins = bins
	e.Modifications = modEvi

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
	output := fmt.Sprintf("%s%smodifications.tsv", sys.MetaDir(), string(filepath.Separator))

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

	outfile := fmt.Sprintf("%s%sdelta-mass.html", sys.MetaDir(), string(filepath.Separator))

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

// addCustomNames adds to the label structures user-defined names to be used on the TMT labels
func getLabelNames(annot string) (map[string]string, *err.Error) {

	var labels = make(map[string]string)

	file, e := os.Open(annot)
	if e != nil {
		return labels, &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: e.Error()}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		names := strings.Split(scanner.Text(), " ")
		labels[names[0]] = names[1]
	}

	if e = scanner.Err(); e != nil {
		return labels, &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: e.Error()}
	}

	return labels, nil
}
