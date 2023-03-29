package spc

import "encoding/xml"

// ProtXML is the root tag
type ProtXML struct {
	Name           string
	ProteinSummary ProteinSummary
}

// ProteinSummary tag is the root level
type ProteinSummary struct {
	ProteinSummaryHeader ProteinSummaryHeader `xml:"protein_summary_header"`
	ProteinGroup         []ProteinGroup       `xml:"protein_group"`
}

// ProteinSummaryHeader tag
type ProteinSummaryHeader struct {
	MinPeptideProbability       float32        `xml:"min_peptide_probability,attr"`
	MinPeptideWeight            float32        `xml:"min_peptide_weight,attr"`
	NumPredictedCorrectProteins float32        `xml:"num_predicted_correct_prots,attr"`
	TotalNumberSpectrumIDs      float32        `xml:"total_no_spectrum_ids,attr"`
	NumInput1Spectra            uint32         `xml:"num_input_1_spectra,attr"`
	NumInput2Spectra            uint32         `xml:"num_input_2_spectra,attr"`
	NumInput3Spectra            uint32         `xml:"num_input_3_spectra,attr"`
	NumInput4Spectra            uint32         `xml:"num_input_4_spectra,attr"`
	NumInput5Spectra            uint32         `xml:"num_input_5_spectra,attr"`
	ProgramDetails              ProgramDetails `xml:"program_details"`
}

// ProgramDetails tag
type ProgramDetails struct {
	Analysis              []byte                `xml:"analysis,attr"`
	Time                  []byte                `xml:"time,attr"`
	Version               []byte                `xml:"version,attr"`
	ProteinProphetDetails ProteinProphetDetails `xml:"proteinprophet_details"`
}

// ProteinProphetDetails tag
type ProteinProphetDetails struct {
	XMLName               xml.Name `xml:"proteinprophet_details"`
	OccamFlag             []byte   `xml:"occam_flag,attr"`
	GroupsFlag            []byte   `xml:"groups_flag,attr"`
	DegenFlag             []byte   `xml:"degen_flag,attr"`
	NSPFlag               []byte   `xml:"nsp_flag,attr"`
	FPKMFlag              []byte   `xml:"fpkm_flag,attr"`
	InitialPeptideWtIters []byte   `xml:"initial_peptide_wt_iters,attr"`
	NspDistributionIters  []byte   `xml:"nsp_distribution_iters,attr"`
	FinalPeptideWtIters   []byte   `xml:"final_peptide_wt_iters,attr"`
	RunOptions            []byte   `xml:"run_options,attr"`
}

// ProteinGroup tag
type ProteinGroup struct {
	GroupNumber uint32    `xml:"group_number,attr"`
	Probability float64   `xml:"probability,attr"`
	Protein     []Protein `xml:"protein"`
}

// Protein tag
type Protein struct {
	ProteinName                     []byte                     `xml:"protein_name,attr"`
	UniqueStrippedPeptides          []byte                     `xml:"unique_stripped_peptides,attr"`
	GroupSiblingID                  []byte                     `xml:"group_sibling_id,attr"`
	NumberIndistinguishableProteins int16                      `xml:"n_indistinguishable_proteins,attr"`
	TotalNumberPeptides             int                        `xml:"total_number_peptides,attr"`
	TotalNumberIndPeptides          int                        `xml:"total_number_distinct_peptides,attr"`
	PercentCoverage                 float32                    `xml:"percent_coverage,attr"`
	PctSpectrumIDs                  float32                    `xml:"pct_spectrum_ids,attr"`
	Probability                     float64                    `xml:"probability,attr"`
	Parameter                       Parameter                  `xml:"parameter"`
	Annotation                      Annotation                 `xml:"annotation"`
	IndistinguishableProtein        []IndistinguishableProtein `xml:"indistinguishable_protein"`
	Peptide                         []Peptide                  `xml:"peptide"`
}

// IndistinguishableProtein tag
type IndistinguishableProtein struct {
	ProteinName string `xml:"protein_name,attr"`
}

// Peptide tag
type Peptide struct {
	PeptideSequence         []byte                 `xml:"peptide_sequence,attr"`
	Charge                  uint8                  `xml:"charge,attr"`
	InitialProbability      float64                `xml:"initial_probability,attr"`
	Weight                  float64                `xml:"weight,attr"`
	GroupWeight             float64                `xml:"group_weight,attr"`
	IsNondegenerateEvidence []byte                 `xml:"is_nondegenerate_evidence,attr"`
	NEnzymaticTermini       uint8                  `xml:"n_enzymatic_termini,attr"`
	CalcNeutralPepMass      float64                `xml:"calc_neutral_pep_mass,attr"`
	ModificationInfo        ModificationInfo       `xml:"modification_info"`
	PeptideParentProtein    []PeptideParentProtein `xml:"peptide_parent_protein"`
}

// PeptideParentProtein tag
type PeptideParentProtein struct {
	ProteinName []byte `xml:"protein_name,attr"`
}

// IndistinguishablePeptide tag
type IndistinguishablePeptide struct {
	XMLName            xml.Name `xml:"indistinguishable_peptide"`
	PeptideSequence    []byte   `xml:"peptide_sequence,attr"`
	Charge             uint8    `xml:"charge,attr"`
	CalcNeutralPepMass float32  `xml:"calc_neutral_pep_mass,attr"`
}
