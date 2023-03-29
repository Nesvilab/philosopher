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
	Date            []byte            `xml:"date,attr"`
	SummaryXML      []byte            `xml:"summary_xml,attr"`
	AnalysisSummary []AnalysisSummary `xml:"analysis_summary"`
	MsmsRunSummary  MsmsRunSummary    `xml:"msms_run_summary"`
}

// AnalysisSummary struct
type AnalysisSummary struct {
	Analysis              []byte                `xml:"analysis,attr"`
	PeptideprophetSummary PeptideprophetSummary `xml:"peptideprophet_summary"`
}

// PeptideprophetSummary struct
type PeptideprophetSummary struct {
	DistributionPoint []DistributionPoint `xml:"distribution_point"`
}

// DistributionPoint ...
type DistributionPoint struct {
	Fvalue         float64 `xml:"fvalue,attr"`
	Obs1Distr      float64 `xml:"obs_1_distr,attr"`
	Model1PosDistr float64 `xml:"model_1_pos_distr,attr"`
	Model1NegDistr float64 `xml:"model_1_neg_distr,attr"`
	Obs2Distr      float64 `xml:"obs_2_distr,attr"`
	Model2PosDistr float64 `xml:"model_2_pos_distr,attr"`
	Model2NegDistr float64 `xml:"model_2_neg_distr,attr"`
	Obs3Distr      float64 `xml:"obs_3_distr,attr"`
	Model3PosDistr float64 `xml:"model_3_pos_distr,attr"`
	Model3NegDistr float64 `xml:"model_3_neg_distr,attr"`
	Obs4Distr      float64 `xml:"obs_4_distr,attr"`
	Model4PosDistr float64 `xml:"model_4_pos_distr,attr"`
	Model4NegDistr float64 `xml:"model_4_neg_distr,attr"`
	Obs5Distr      float64 `xml:"obs_5_distr,attr"`
	Model5PosDistr float64 `xml:"model_5_pos_distr,attr"`
	Model5NegDistr float64 `xml:"model_5_neg_distr,attr"`
	Obs6Distr      float64 `xml:"obs_6_distr,attr"`
	Model6PosDistr float64 `xml:"model_6_pos_distr,attr"`
	Model6NegDistr float64 `xml:"model_6_neg_distr,attr"`
	Obs7Distr      float64 `xml:"obs_7_distr,attr"`
	Model7PosDistr float64 `xml:"model_7_pos_distr,attr"`
	Model7NegDistr float64 `xml:"model_7_neg_distr,attr"`
}

// MsmsRunSummary tag
type MsmsRunSummary struct {
	BaseName       []byte          `xml:"base_name,attr"`
	SearchEngine   []byte          `xml:"search_engine,attr"`
	MsmsRunRummary []byte          `xml:"msms_run_summary,attr"`
	MsManufacturer []byte          `xml:"msManufacturer,attr"`
	MsModel        []byte          `xml:"msModel,attr"`
	MsIonization   []byte          `xml:"msIonization,attr"`
	MsMassAnalyzer []byte          `xml:"msMassAnalyzer,attr"`
	MsDetector     []byte          `xml:"msDetector,attr"`
	RawDataType    []byte          `xml:"raw_data_type,attr"`
	RawData        []byte          `xml:"raw_data,attr"`
	SearchSummary  SearchSummary   `xml:"search_summary"`
	SpectrumQuery  []SpectrumQuery `xml:"spectrum_query"`
}

// Specificity tag
type Specificity struct {
	Xmlname xml.Name `xml:"specificity"`
	Cut     []byte   `xml:"cut,attr"`
	NoCut   []byte   `xml:"no_cut,attr"`
	Sense   []byte   `xml:"sense,attr"`
}

// SearchSummary tag
type SearchSummary struct {
	SearchID               uint16                  `xml:"search_id,attr"`
	BaseName               []byte                  `xml:"base_name,attr"`
	SearchEngine           []byte                  `xml:"search_engine,attr"`
	SearchEngineVersion    []byte                  `xml:"search_engine_version,attr"`
	SearchDatabase         SearchDatabase          `xml:"search_database"`
	AminoAcidModifications []AminoacidModification `xml:"aminoacid_modification"`
	TerminalModifications  []TerminalModification  `xml:"terminal_modification"`
	Parameter              []Parameter             `xml:"parameter"`
}

// SearchDatabase tag
type SearchDatabase struct {
	XMLName   xml.Name `xml:"search_database"`
	LocalPath []byte   `xml:"local_path,attr"`
	Type      []byte   `xml:"type,attr"`
}

// AminoacidModification tag
type AminoacidModification struct {
	AminoAcid []byte  `xml:"aminoacid,attr"`
	MassDiff  float64 `xml:"massdiff,attr"`
	Mass      float64 `xml:"mass,attr"`
	Variable  []byte  `xml:"variable,attr"`
}

// TerminalModification tag
type TerminalModification struct {
	MassDiff float64 `xml:"massdiff,attr"`
	Mass     float64 `xml:"mass,attr"`
	Terminus []byte  `xml:"terminus,attr"`
	Variable []byte  `xml:"variable,attr"`
}

// SpectrumQuery tag
type SpectrumQuery struct {
	CompensationVoltage              string       `xml:"compensation_voltage,attr"`
	Spectrum                         []byte       `xml:"spectrum,attr"`
	StartScan                        int          `xml:"start_scan,attr"`
	EndScan                          int          `xml:"end_scan,attr"`
	AssumedCharge                    uint8        `xml:"assumed_charge,attr"`
	Index                            uint32       `xml:"index,attr"`
	RetentionTimeSec                 float64      `xml:"retention_time_sec,attr"`
	IonMobility                      float64      `xml:"ion_mobility,attr"`
	UncalibratedPrecursorNeutralMass float64      `xml:"uncalibrated_precursor_neutral_mass,attr"`
	PrecursorNeutralMass             float64      `xml:"precursor_neutral_mass,attr"`
	SearchResult                     SearchResult `xml:"search_result"`
}

// SearchResult tag
type SearchResult struct {
	SearchHit []SearchHit `xml:"search_hit"`
}

// SearchHit tag
type SearchHit struct {
	HitRank             uint8                `xml:"hit_rank,attr"`
	Peptide             []byte               `xml:"peptide,attr"`
	PrevAA              []byte               `xml:"peptide_prev_aa,attr"`
	NextAA              []byte               `xml:"peptide_next_aa,attr"`
	Protein             []byte               `xml:"protein,attr"`
	Class               []byte               `xml:"fdr_group,attr"`
	TotalTerm           uint8                `xml:"num_tol_term,attr"`
	MissedCleavages     uint8                `xml:"num_missed_cleavages,attr"`
	IsRejected          uint8                `xml:"is_rejected,attr"`
	TotalProteins       uint32               `xml:"num_tot_proteins,attr"`
	MatchedIons         uint16               `xml:"num_matched_ions,attr"`
	TotalIons           uint16               `xml:"tot_num_ions,attr"`
	MatchedPeptides     uint32               `xml:"num_matched_peptides,attr"`
	CalcNeutralPepMass  float64              `xml:"calc_neutral_pep_mass,attr"`
	Massdiff            float64              `xml:"massdiff,attr"`
	ModificationInfo    ModificationInfo     `xml:"modification_info"`
	Score               []SearchScore        `xml:"search_score"`
	AnalysisResult      []AnalysisResult     `xml:"analysis_result"`
	AlternativeProteins []AlternativeProtein `xml:"alternative_protein"`
	PTMResult           PTMResult            `xml:"ptm_result"`
}

// AlternativeProtein tag
type AlternativeProtein struct {
	Protein   []byte `xml:"protein,attr"`
	PepPrevAA []byte `xml:"peptide_prev_aa,attr"`
	PepNextAA []byte `xml:"peptide_next_aa,attr"`
}

// AnalysisResult tag
type AnalysisResult struct {
	Analysis             []byte               `xml:"analysis,attr"`
	PeptideProphetResult PeptideProphetResult `xml:"peptideprophet_result"`
	InterProphetResult   InterProphetResult   `xml:"interprophet_result"`
	PTMProphetResult     []PTMProphetResult   `xml:"ptmprophet_result"`
}

// PeptideProphetResult tag
type PeptideProphetResult struct {
	Probability float64 `xml:"probability,attr"`
}

// InterProphetResult tag
type InterProphetResult struct {
	Probability float64 `xml:"probability,attr"`
}

// PTMProphetResult tag
type PTMProphetResult struct {
	PTM                     []byte                    `xml:"ptm,attr"`
	PTMPeptide              []byte                    `xml:"ptm_peptide,attr"`
	ModAminoAcidProbability []ModAminoAcidProbability `xml:"mod_aminoacid_probability"`
}

// ModAminoAcidProbability tag
type ModAminoAcidProbability struct {
	Position    int     `xml:"position,attr"`
	Probability float32 `xml:"probability,attr"`
}

// SearchScoreSummary tag
type SearchScoreSummary struct {
	XMLName   xml.Name    `xml:"search_score_summary"`
	Parameter []Parameter `xml:"parameter"`
}

// SearchScore tag
type SearchScore struct {
	Name  []byte `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// ProphetModel tag
type ProphetModel struct {
	Charge uint8
	Points map[string]uint8
}

// PTMResult tag
type PTMResult struct {
	BestScoreWithPTM    string `xml:"best_score_with_ptm,attr"`
	ScoreWithoutPTM     string `xml:"score_without_ptm,attr"`
	LocalizationPeptide string `xml:"localization_peptide,attr"`
}
