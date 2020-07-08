package spc

import (
	"encoding/xml"
)

// PepXML is the root tag
type PepXML struct {
	Name                 string
	MsmsPipelineAnalysis MsmsPipelineAnalysis
}

// MsmsPipelineAnalysis tag
type MsmsPipelineAnalysis struct {
	XMLName         xml.Name          `xml:"msms_pipeline_analysis"`
	Date            string            `xml:"date,attr"`
	SummaryXML      []byte            `xml:"summary_xml,attr"`
	AnalysisSummary []AnalysisSummary `xml:"analysis_summary"`
	MsmsRunSummary  MsmsRunSummary    `xml:"msms_run_summary"`
}

// AnalysisSummary struct
type AnalysisSummary struct {
	XMLName               xml.Name              `xml:"analysis_summary"`
	Analysis              string                `xml:"analysis,attr,omitempty"`
	Time                  string                `xml:"time,attr,omitempty"`
	PeptideprophetSummary PeptideprophetSummary `xml:"peptideprophet_summary"`
}

// PeptideprophetSummary struct
type PeptideprophetSummary struct {
	XMLName           xml.Name            `xml:"peptideprophet_summary"`
	Version           []byte              `xml:"version,attr,omitempty"`
	Options           []byte              `xml:"options,attr,omitempty"`
	MixtureModel      []MixtureModel      `xml:"mixture_model"`
	DistributionPoint []DistributionPoint `xml:"distribution_point"`
}

// DistributionPoint ...
type DistributionPoint struct {
	XMLName        xml.Name `xml:"distribution_point"`
	Fvalue         float64  `xml:"fvalue,attr"`
	Obs1Distr      float64  `xml:"obs_1_distr,attr,omitempty"`
	Model1PosDistr float64  `xml:"model_1_pos_distr,attr,omitempty"`
	Model1NegDistr float64  `xml:"model_1_neg_distr,attr,omitempty"`
	Obs2Distr      float64  `xml:"obs_2_distr,attr,omitempty"`
	Model2PosDistr float64  `xml:"model_2_pos_distr,attr,omitempty"`
	Model2NegDistr float64  `xml:"model_2_neg_distr,attr,omitempty"`
	Obs3Distr      float64  `xml:"obs_3_distr,attr,omitempty"`
	Model3PosDistr float64  `xml:"model_3_pos_distr,attr,omitempty"`
	Model3NegDistr float64  `xml:"model_3_neg_distr,attr,omitempty"`
	Obs4Distr      float64  `xml:"obs_4_distr,attr,omitempty"`
	Model4PosDistr float64  `xml:"model_4_pos_distr,attr,omitempty"`
	Model4NegDistr float64  `xml:"model_4_neg_distr,attr,omitempty"`
	Obs5Distr      float64  `xml:"obs_5_distr,attr,omitempty"`
	Model5PosDistr float64  `xml:"model_5_pos_distr,attr,omitempty"`
	Model5NegDistr float64  `xml:"model_5_neg_distr,attr,omitempty"`
	Obs6Distr      float64  `xml:"obs_6_distr,attr,omitempty"`
	Model6PosDistr float64  `xml:"model_6_pos_distr,attr,omitempty"`
	Model6NegDistr float64  `xml:"model_6_neg_distr,attr,omitempty"`
	Obs7Distr      float64  `xml:"obs_7_distr,attr,omitempty"`
	Model7PosDistr float64  `xml:"model_7_pos_distr,attr,omitempty"`
	Model7NegDistr float64  `xml:"model_7_neg_distr,attr,omitempty"`
}

// MixtureModel struct
type MixtureModel struct {
	XMLName            xml.Name       `xml:"mixture_model"`
	PrecursorIonCharge uint8          `xml:"precursor_ion_charge,attr,omitempty"`
	Comments           []byte         `xml:"comments,attr,omitempty"`
	PriorProbability   float64        `xml:"prior_probability,attr,omitempty"`
	EstTotCorrect      float64        `xml:"est_tot_correct,attr,omitempty"`
	TotNumSpectra      float64        `xml:"tot_num_spectra,attr,omitempty"`
	NumIterations      float64        `xml:"num_iterations,attr,omitempty"`
	Mixturemodel       []Mixturemodel `xml:"mixturemodel"`
}

// Mixturemodel struct
type Mixturemodel struct {
	XMLName      xml.Name `xml:"mixturemodel"`
	Name         []byte   `xml:"name,attr"`
	PosBandwidth float64  `xml:"pos_bandwidth,attr"`
	NegBandwidth float64  `xml:"neg_bandwidth,attr"`
	Point        []Point  `xml:"point"`
}

// Point struct
type Point struct {
	XMLName xml.Name `xml:"point"`
	Value   float64  `xml:"value,attr"`
	PosDens float64  `xml:"pos_dens,attr"`
	NegDens float64  `xml:"neg_dens,attr"`
}

// MixturemodelDistribution struct
type MixturemodelDistribution struct {
	XMLName xml.Name `xml:"mixturemodel_distribution"`
	Name    []byte   `xml:"name,attr"`
}

// MsmsRunSummary tag
type MsmsRunSummary struct {
	XMLName        xml.Name        `xml:"msms_run_summary"`
	BaseName       []byte          `xml:"base_name,attr,omitempty"`
	SearchEngine   []byte          `xml:"search_engine,attr,omitempty"`
	MsmsRunRummary []byte          `xml:"msms_run_summary,attr,omitempty"`
	MsManufacturer []byte          `xml:"msManufacturer,attr,omitempty"`
	MsModel        []byte          `xml:"msModel,attr,omitempty"`
	MsIonization   []byte          `xml:"msIonization,attr,omitempty"`
	MsMassAnalyzer []byte          `xml:"msMassAnalyzer,attr,omitempty"`
	MsDetector     []byte          `xml:"msDetector,attr,omitempty"`
	RawDataType    []byte          `xml:"raw_data_type,attr,omitempty"`
	RawData        []byte          `xml:"raw_data,attr,omitempty"`
	SampleEnzyme   SampleEnzyme    `xml:"sample_enzyme"`
	SearchSummary  SearchSummary   `xml:"search_summary"`
	SpectrumQuery  []SpectrumQuery `xml:"spectrum_query"`
}

// SampleEnzyme tag
type SampleEnzyme struct {
	XMLName     xml.Name    `xml:"sample_enzyme"`
	Name        []byte      `xml:"name,attr,omitempty"`
	Specificity Specificity `xml:"specificity"`
}

// Specificity tag
type Specificity struct {
	Xmlname xml.Name `xml:"specificity"`
	Cut     []byte   `xml:"cut,attr,omitempty"`
	NoCut   []byte   `xml:"no_cut,attr,omitempty"`
	Sense   []byte   `xml:"sense,attr,omitempty"`
}

// SearchSummary tag
type SearchSummary struct {
	XMLName                   xml.Name                    `xml:"search_summary"`
	SearchID                  uint16                      `xml:"search_id,attr,omitempty"`
	BaseName                  []byte                      `xml:"base_name,attr,omitempty"`
	SearchEngine              []byte                      `xml:"search_engine,attr,omitempty"`
	SearchEngineVersion       []byte                      `xml:"search_engine_version,attr,omitempty"`
	SearchDatabase            SearchDatabase              `xml:"search_database"`
	EnzymaticSearchConstraint []EnzymaticSearchConstraint `xml:"enzymatic_search_constraint"`
	AminoAcidModifications    []AminoacidModification     `xml:"aminoacid_modification"`
	TerminalModifications     []TerminalModification      `xml:"terminal_modification"`
	Parameter                 []Parameter                 `xml:"parameter"`
}

// SearchDatabase tag
type SearchDatabase struct {
	XMLName   xml.Name `xml:"search_database"`
	LocalPath []byte   `xml:"local_path,attr,omitempty"`
	Type      []byte   `xml:"type,attr"`
}

// EnzymaticSearchConstraint tag
type EnzymaticSearchConstraint struct {
	XMLName                 xml.Name `xml:"enzymatic_search_constraint"`
	Enzyme                  []byte   `xml:"enzyme,attr,omitempty"`
	MaxNumInternalCleavages uint32   `xml:"max_num_internal_cleavages,attr,omitempty"`
	MinNumTermini           uint8    `xml:"min_number_termini,attr,omitempty"`
}

// AminoacidModification tag
type AminoacidModification struct {
	XMLName   xml.Name `xml:"aminoacid_modification"`
	AminoAcid []byte   `xml:"aminoacid,attr,omitempty"`
	MassDiff  float64  `xml:"massdiff,attr,omitempty"`
	Mass      float64  `xml:"mass,attr,omitempty"`
	Variable  []byte   `xml:"variable,attr,omitempty"`
}

// TerminalModification tag
type TerminalModification struct {
	XMLName         xml.Name `xml:"terminal_modification"`
	MassDiff        float64  `xml:"massdiff,attr,omitempty"`
	ProteinTerminus []byte   `xml:"protein_terminus,attr,omitempty"`
	Mass            float64  `xml:"mass,attr,omitempty"`
	Terminus        []byte   `xml:"terminus,attr,omitempty"`
	Variable        []byte   `xml:"variable,attr,omitempty"`
}

// SpectrumQuery tag
type SpectrumQuery struct {
	XMLName                          xml.Name     `xml:"spectrum_query"`
	Spectrum                         string       `xml:"spectrum,attr,omitempty"`
	SpectrumNativeID                 []byte       `xml:"spectrumNativeID,attr,omitempty"`
	StartScan                        int          `xml:"start_scan,attr,omitempty"`
	EndScan                          int          `xml:"end_scan,attr,omitempty"`
	UncalibratedPrecursorNeutralMass float64      `xml:"uncalibrated_precursor_neutral_mass,attr,omitempty"`
	PrecursorNeutralMass             float64      `xml:"precursor_neutral_mass,attr,omitempty"`
	AssumedCharge                    uint8        `xml:"assumed_charge,attr,omitempty"`
	Index                            uint32       `xml:"index,attr,omitempty"`
	RetentionTimeSec                 float64      `xml:"retention_time_sec,attr,omitempty"`
	PrecursorIntensity               float64      `xml:"precursor_intensity,attr,omitempty"`
	IonMobility                      float64      `xml:"ion_mobility,attr,omitempty"`
	SearchResult                     SearchResult `xml:"search_result"`
}

// SearchResult tag
type SearchResult struct {
	XMLName   xml.Name    `xml:"search_result"`
	SearchHit []SearchHit `xml:"search_hit"`
}

// SearchHit tag
type SearchHit struct {
	XMLName             xml.Name             `xml:"search_hit"`
	HitRank             uint8                `xml:"hit_rank,attr"`
	Peptide             string               `xml:"peptide,attr"`
	PrevAA              string               `xml:"peptide_prev_aa,attr"`
	NextAA              string               `xml:"peptide_next_aa,attr"`
	Protein             string               `xml:"protein,attr"`
	ProteinDescr        string               `xml:"protein_descr,attr"`
	TotalProteins       uint16               `xml:"num_tot_proteins,attr"`
	MatchedIons         uint16               `xml:"num_matched_ions,attr"`
	TotalIons           uint16               `xml:"tot_num_ions,attr"`
	CalcNeutralPepMass  float64              `xml:"calc_neutral_pep_mass,attr"`
	Massdiff            float64              `xml:"massdiff,attr"`
	TotalTerm           uint8                `xml:"num_tol_term,attr"`
	MissedCleavages     uint8                `xml:"num_missed_cleavages,attr"`
	MatchedPeptides     uint32               `xml:"num_matched_peptides,attr"`
	IsRejected          uint8                `xml:"is_rejected,attr"`
	Score               []SearchScore        `xml:"search_score"`
	ModificationInfo    ModificationInfo     `xml:"modification_info"`
	AnalysisResult      []AnalysisResult     `xml:"analysis_result"`
	AlternativeProteins []AlternativeProtein `xml:"alternative_protein"`
}

// AlternativeProtein tag
type AlternativeProtein struct {
	XMLName     xml.Name `xml:"alternative_protein"`
	Protein     string   `xml:"protein,attr"`
	Description string   `xml:"protein_descr,attr,omitempty"`
	NumTolTerm  int8     `xml:"num_tol_tem,attr,omitempty"`
	PepPrevAA   []byte   `xml:"peptide_prev_aa,attr,omitempty"`
	PepNextAA   []byte   `xml:"peptide_next_aa,attr,omitempty"`
}

// AnalysisResult tag
type AnalysisResult struct {
	XMLName              xml.Name             `xml:"analysis_result"`
	Analysis             string               `xml:"analysis,attr"`
	PeptideProphetResult PeptideProphetResult `xml:"peptideprophet_result"`
	InterProphetResult   InterProphetResult   `xml:"interprophet_result"`
	PTMProphetResult     []PTMProphetResult   `xml:"ptmprophet_result"`
	SearchScoreSummary   SearchScoreSummary   `xml:"search_score_summary,omitempty"`
}

// PeptideProphetResult tag
type PeptideProphetResult struct {
	XMLName            xml.Name           `xml:"peptideprophet_result"`
	Probability        float64            `xml:"probability,attr,omitempty"`
	AllNttProb         []byte             `xml:"all_ntt_prob,attr,omitempty"`
	SearchScoreSummary SearchScoreSummary `xml:"search_score_summary,omitempty"`
}

// InterProphetResult tag
type InterProphetResult struct {
	XMLName     xml.Name `xml:"interprophet_result"`
	Probability float64  `xml:"probability,attr"`
	AllNttProb  []byte   `xml:"all_ntt_prob,attr"`
}

// PTMProphetResult tag
type PTMProphetResult struct {
	XMLName                 xml.Name                  `xml:"ptmprophet_result"`
	Prior                   float64                   `xml:"prior,attr"`
	PTM                     []byte                    `xml:"ptm,attr"`
	PTMPeptide              []byte                    `xml:"ptm_peptide,attr"`
	ModAminoAcidProbability []ModAminoAcidProbability `xml:"mod_aminoacid_probability"`
}

// ModAminoAcidProbability tag
type ModAminoAcidProbability struct {
	XMLName     xml.Name `xml:"mod_aminoacid_probability"`
	Position    int      `xml:"position,attr"`
	Probability float32  `xml:"probability,attr"`
}

// SearchScoreSummary tag
type SearchScoreSummary struct {
	XMLName   xml.Name    `xml:"search_score_summary"`
	Parameter []Parameter `xml:"parameter"`
}

// SearchScore tag
type SearchScore struct {
	XMLName xml.Name `xml:"search_score"`
	Name    []byte   `xml:"name,attr"`
	Value   string   `xml:"value,attr"`
}

// ProphetModel struct
type ProphetModel struct {
	Charge uint8
	Points map[string]uint8
}
