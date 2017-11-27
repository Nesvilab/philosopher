package xml

import (
	"encoding/gob"
	"errors"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher/lib/pro"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/prvst/philosopher/lib/tmt"
)

// ProtXML struct
type ProtXML struct {
	FileName   string
	DecoyTag   string
	Groups     GroupList
	RunOptions string
}

// GroupIdentification tag
type GroupIdentification struct {
	GroupNumber uint32
	Probability float64
	Proteins    ProtIDList
}

// ProteinIdentification struct
type ProteinIdentification struct {
	GroupNumber              uint32
	GroupSiblingID           string
	ProteinName              string
	UniqueStrippedPeptides   []string
	Length                   int
	PercentCoverage          float32
	PctSpectrumIDs           float32
	GroupProbability         float64
	Probability              float64
	Confidence               float64
	TopPepProb               float64
	IndistinguishableProtein []string
	PeptideIons              []PeptideIonIdentification
	HasRazor                 bool
	Picked                   int
}

// PeptideIonIdentification struct
type PeptideIonIdentification struct {
	PeptideSequence         string
	ModifiedPeptide         string
	Charge                  uint8
	InitialProbability      float64
	Weight                  float64
	GroupWeight             float64
	CalcNeutralPepMass      float64
	SharedParentProteins    int
	Razor                   int
	IsNondegenerateEvidence bool
	IsUnique                bool
	Labels                  tmt.Labels
}

// GroupList represents a protein group list
type GroupList []GroupIdentification

// ProtIDList list represents a list of custom protein identifications
type ProtIDList []ProteinIdentification

// Len function for sortng
func (p ProtIDList) Len() int {
	return len(p)
}

// Less function for sorting
func (p ProtIDList) Less(i, j int) bool {
	return p[i].TopPepProb > p[j].TopPepProb
}

// Swap function for sorting
func (p ProtIDList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Read ...
func (p *ProtXML) Read(f string) error {

	var xml pro.XML
	err := xml.Parse(f)
	if err != nil {
		return err
	}

	var ptg = xml.ProteinSummary.ProteinGroup
	var groups GroupList

	for _, i := range ptg {

		var gi GroupIdentification
		var proteinList ProtIDList

		gi.GroupNumber = i.GroupNumber
		gi.Probability = i.Probability

		for jindex, j := range i.Protein {

			// correcting group probabilities
			if jindex == 0 {
				if i.Probability == 1 && j.Probability == 0 {
					j.Probability = float64(i.Probability)
				}
			}

			var ptid ProteinIdentification

			ptid.GroupNumber = i.GroupNumber
			ptid.GroupProbability = i.Probability
			ptid.Probability = i.Probability
			ptid.ProteinName = string(j.ProteinName)
			ptid.Probability = j.Probability
			ptid.PercentCoverage = j.PercentCoverage
			ptid.PctSpectrumIDs = j.PctSpectrumIDs
			ptid.TopPepProb = j.Peptide[0].InitialProbability
			ptid.GroupSiblingID = string(j.GroupSiblingID)

			if strings.EqualFold(string(j.Parameter.Name), "prot_length") {
				ptid.Length = j.Parameter.Value
			}

			// collect indistinguishable proteins
			if len(j.IndistinguishableProtein) > 0 {
				for _, k := range j.IndistinguishableProtein {
					ptid.IndistinguishableProtein = append(ptid.IndistinguishableProtein, k.ProteinName)
				}
			}

			for _, k := range j.Peptide {

				var pepid PeptideIonIdentification

				pepid.PeptideSequence = string(k.PeptideSequence)
				pepid.ModifiedPeptide = string(k.ModificationInfo.ModifiedPeptide)
				pepid.Charge = k.Charge
				pepid.InitialProbability = k.InitialProbability
				pepid.Weight = k.Weight
				pepid.GroupWeight = k.GroupWeight
				pepid.CalcNeutralPepMass = k.CalcNeutralPepMass
				pepid.SharedParentProteins = len(k.PeptideParentProtein)
				pepid.Razor = -1

				if string(k.IsNondegenerateEvidence) == "y" || string(k.IsNondegenerateEvidence) == "Y" {
					pepid.IsNondegenerateEvidence = true
				}

				// get the number of shared ions
				if pepid.Weight < 1 {
					pepid.IsUnique = false
				} else {
					pepid.IsUnique = true
				}

				if strings.EqualFold(string(k.IsNondegenerateEvidence), "Y") || strings.EqualFold(string(k.IsNondegenerateEvidence), "y") {
					pepid.IsNondegenerateEvidence = true
				} else {
					pepid.IsNondegenerateEvidence = false
				}

				ptid.PeptideIons = append(ptid.PeptideIons, pepid)

				pepid = PeptideIonIdentification{}
			}

			peps := strings.Split(string(j.UniqueStrippedPeptides), "+")
			ptid.UniqueStrippedPeptides = peps
			proteinList = append(proteinList, ptid)
		}

		gi.Proteins = proteinList
		groups = append(groups, gi)
	}

	p.RunOptions = string(xml.ProteinSummary.ProteinSummaryHeader.ProgramDetails.ProteinProphetDetails.RunOptions)
	p.Groups = groups

	if len(groups) == 0 {
		return errors.New("No Protein groups detected, check your file and try again")
	}

	return nil
}

// PromoteProteinIDs promotes protein identifications where the reference protein
// is indistinguishable to other target proteins.
func (p *ProtXML) PromoteProteinIDs() {

	for i := range p.Groups {
		for j := range p.Groups[i].Proteins {

			var list []string
			var ref string

			if strings.Contains(string(p.Groups[i].Proteins[j].ProteinName), p.DecoyTag) {
				for k := range p.Groups[i].Proteins[j].IndistinguishableProtein {
					if !strings.Contains(string(p.Groups[i].Proteins[j].IndistinguishableProtein[k]), p.DecoyTag) {
						list = append(list, string(p.Groups[i].Proteins[j].IndistinguishableProtein[k]))
					}
				}
			}

			if len(list) > 1 {
				for i := range list {
					if strings.Contains(list[i], "sp|") {
						ref = list[i]
						break
					} else {
						ref = list[i]
					}
				}
				p.Groups[i].Proteins[j].ProteinName = ref
			}

		}
	}

	return
}

// Serialize converts the whle structure to a gob file
func (p *ProtXML) Serialize() error {

	var err error

	// create a file
	dataFile, err := os.Create(sys.ProtxmlBin())
	if err != nil {
		return err
	}

	dataEncoder := gob.NewEncoder(dataFile)
	goberr := dataEncoder.Encode(p)
	if goberr != nil {
		logrus.Fatal("Cannot save results, Bad format", goberr)
	}
	dataFile.Close()

	return nil
}

// Restore reads philosopher results files and restore the data sctructure
func (p *ProtXML) Restore() error {

	file, _ := os.Open(sys.ProtxmlBin())

	dec := gob.NewDecoder(file)
	err := dec.Decode(&p)
	if err != nil {
		return errors.New("Could not restore Philosopher result. Please check file path")
	}

	return nil
}

// Serialize converts the whle structure to a gob file
func (p *ProtIDList) Serialize() error {

	var err error

	// create a file
	dataFile, err := os.Create(sys.ProBin())
	if err != nil {
		return err
	}

	dataEncoder := gob.NewEncoder(dataFile)
	goberr := dataEncoder.Encode(p)
	if goberr != nil {
		logrus.Fatal("Cannot save results, Bad format", goberr)
	}
	dataFile.Close()

	return nil
}

// Restore reads philosopher results files and restore the data sctructure
func (p *ProtIDList) Restore() error {

	file, _ := os.Open(sys.ProBin())

	dec := gob.NewDecoder(file)
	err := dec.Decode(&p)
	if err != nil {
		return errors.New("Could not restore Philosopher result. Please check file path")
	}

	return nil
}
