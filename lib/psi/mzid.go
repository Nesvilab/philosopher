package psi

import "encoding/xml"

// MzIdentML is the upper-most hierarchy level of mzIdentML with sub-containers
// for example describing software, protocols and search results (spectrum
// identifications or protein detection results)
type MzIdentML struct {
	XMLName                    xml.Name                   `xml:"MzIdentML"`
	CreationDate               string                     `xml:"creationDate,attr,omitempty"`
	Name                       string                     `xml:"name,attr,omitempty"`
	ID                         string                     `xml:"id,attr,omitempty"`
	Version                    string                     `xml:"version,attr,omitempty"`
	CvList                     CvList                     `xml:"cvList"`
	AnalysisSoftwareList       AnalysisSoftwareList       `xml:"AnalysisSoftwareList"`
	Provider                   Provider                   `xml:"Provider"`
	AuditCollection            AuditCollection            `xml:"AuditCollection"`
	AnalysisSampleCollection   AnalysisSampleCollection   `xml:"AnalysisSampleCollection"`
	SequenceCollection         SequenceCollection         `xml:"SequenceCollection"`
	AnalysisCollection         AnalysisCollection         `xml:"AnalysisCollection"`
	AnalysisProtocolCollection AnalysisProtocolCollection `xml:"AnalysisProtocolCollection"`
	DataCollection             DataCollection             `xml:"DataCollection"`
	BibliographicReference     []BibliographicReference   `xml:"BibliographicReference"`
}

// AnalysisSoftwareList is the software packages used to perform the analyses
type AnalysisSoftwareList struct {
	XMLName          xml.Name           `xml:"AnalysisSoftwareList"`
	AnalysisSoftware []AnalysisSoftware `xml:"AnalysisSoftware"`
}

// AnalysisSoftware is the software used for performing the analysis
type AnalysisSoftware struct {
	XMLName        xml.Name       `xml:"AnalysisSoftware"`
	ID             string         `xml:"id,attr,omitempty"`
	Name           string         `xml:"name,attr,omitempty"`
	URI            string         `xml:"uri,attr,omitempty"`
	Version        string         `xml:"version,attr,omitempty"`
	ContactRole    ContactRole    `xml:"ContactRole"`
	SoftwareName   SoftwareName   `xml:"SoftwareName"`
	Customizations Customizations `xml:"Customizations"`
}

// ContactRole is the Contact that provided the document instance
type ContactRole struct {
	XMLName    xml.Name `xml:"ContactRole"`
	ContactRef string   `xml:"contact_ref,attr,omitempty"`
	Role       Role     `xml:"Role"`
}

// Role is single entry from an ontology or a controlled vocabulary
type Role struct {
	XMLName xml.Name `xml:"Role"`
	CVParam CVParam  `xml:"cvParam"`
}

// SoftwareName is the name of the analysis software package, sourced from a CV
// if available
type SoftwareName struct {
	XMLName   xml.Name  `xml:"SoftwareName"`
	CVParam   CVParam   `xml:"cvParam"`
	UserParam UserParam `xml:"userParam"`
}

// Customizations is Any customizations to the software, such as alternative
// scoring mechanisms implemented, should be documented here as free text
type Customizations struct {
	XMLName xml.Name `xml:"Customizations"`
	Value   string   `xml:",chardata"`
}

// Provider is the Provider of the mzIdentML record in terms of the contact and
// software
type Provider struct {
	XMLName             xml.Name    `xml:"Provider"`
	AnalysisSoftwareRef string      `xml:"analysisSoftware_ref,attr,omitempty"`
	ID                  string      `xml:"id,attr,omitempty"`
	Name                string      `xml:"name,attr,omitempty"`
	ContactRole         ContactRole `xml:"ContactRole"`
}

// AuditCollection is the complete set of Contacts (people and organisations)
// for this file
type AuditCollection struct {
	XMLName      xml.Name     `xml:"AuditCollection"`
	Person       Person       `xml:"Person"`
	Organization Organization `xml:"Organization"`
}

// Person is a person's name and contact details. Any additional information
// such as the address, contact email etc. should be supplied using CV
// parameters or user parameters
type Person struct {
	XMLName     xml.Name      `xml:"Person"`
	FirstName   string        `xml:"firstName,attr,omitempty"`
	ID          string        `xml:"id,attr,omitempty"`
	LastName    string        `xml:"lastName,attr,omitempty"`
	MidInitials string        `xml:"midInitials,attr,omitempty"`
	Name        string        `xml:"name,attr,omitempty"`
	CVParam     []CVParam     `xml:"cvParam"`
	UserParam   []UserParam   `xml:"userParam"`
	Affiliation []Affiliation `xml:"Affiliation"`
}

// Affiliation is the organization a person belongs to
type Affiliation struct {
	XMLName         xml.Name `xml:"Affiliation"`
	OrganizationRef string   `xml:"organization_ref,attr,omitempty"`
}

// Organization are entities like companies, universities, government agencies.
// Any additional information such as the address, email etc. should be supplied
// either as CV parameters or as user parameters.
type Organization struct {
	XMLName   xml.Name    `xml:"Organization"`
	ID        string      `xml:"id,attr,omitempty"`
	Name      string      `xml:"name,attr,omitempty"`
	CVParam   []CVParam   `xml:"cvParam"`
	UserParam []UserParam `xml:"userParam"`
	Parent    Parent      `xml:"Parent"`
}

// Parent is the containing organization (the university or business which a lab
// belongs to, etc.)
type Parent struct {
	XMLName         xml.Name `xml:"Parent"`
	OrganizationRef string   `xml:"organization_ref,attr,omitempty"`
}

// AnalysisSampleCollection is the samples analysed can optionally be recorded
// using CV terms for descriptions. If a composite sample has been analysed, the
// subsample association can be used to build a hierarchical description
type AnalysisSampleCollection struct {
	XMLName xml.Name `xml:"AnalysisSampleCollection"`
	Sample  []Sample `xml:"Sample"`
}

// Sample is a description of the sample analysed by mass spectrometry using
// CVParams or UserParams. If a composite sample has been analysed, a parent
// sample should be defined, which references subsamples. This represents any
// kind of substance used in an experimental workflow, such as whole organisms,
// cells, DNA, solutions, compounds and experimental substances
// (gels, arrays etc.)
type Sample struct {
	XMLName     xml.Name      `xml:"Sample"`
	ID          string        `xml:"id,attr,omitempty"`
	Name        string        `xml:"name,attr,omitempty"`
	ContactRole []ContactRole `xml:"ContactRole"`
	SubSample   []SubSample   `xml:"SubSample"`
	CVParam     []CVParam     `xml:"cvParam"`
	UserParam   []UserParam   `xml:"userParam"`
}

// SubSample is the references to the individual component samples within a
// mixed parent sample
type SubSample struct {
	XMLName   xml.Name `xml:"SubSample"`
	SampleRef string   `xml:"sample_ref,attr,omitempty"`
}

// SequenceCollection is the collection of sequences (DBSequence or Peptide)
// identified and their relationship between each other (PeptideEvidence) to be
// referenced elsewhere in the results
type SequenceCollection struct {
	XMLName         xml.Name          `xml:"SequenceCollection"`
	DBSequence      []DBSequence      `xml:"DBSequence"`
	Peptide         []Peptide         `xml:"Peptide"`
	PeptideEvidence []PeptideEvidence `xml:"PeptideEvidence"`
}

// DBSequence is a database sequence from the specified SearchDatabase
// (nucleic acid or amino acid). If the sequence is nucleic acid, the source
// nucleic acid sequence should be given in the seq attribute rather than a
// translated sequence
type DBSequence struct {
	XMLName           xml.Name    `xml:"DBSequence"`
	Accession         string      `xml:"accession,attr,omitempty"`
	ID                string      `xml:"id,attr,omitempty"`
	Length            string      `xml:"length,attr,omitempty"`
	Name              string      `xml:"name,attr,omitempty"`
	SearchDatabaseRef string      `xml:"searchDatabase_ref,attr,omitempty"`
	Seq               Seq         `xml:"Seq"`
	CVParam           []CVParam   `xml:"cvParam"`
	UserParam         []UserParam `xml:"userParam"`
}

// Seq is the actual sequence of amino acids or nucleic acid
type Seq struct {
	XMLName xml.Name `xml:"Seq"`
	Value   string   `xml:",chardata"`
}

// Peptide is One (poly)peptide (a sequence with modifications). The combination
// of Peptide sequence and modifications MUST be unique in the file
type Peptide struct {
	XMLName                  xml.Name                   `xml:"Peptide"`
	ID                       string                     `xml:"id,attr,omitempty"`
	Name                     string                     `xml:"name,attr,omitempty"`
	PeptideSequence          PeptideSequence            `xml:"PeptideSequence"`
	Modification             []Modification             `xml:"Modification"`
	SubstitutionModification []SubstitutionModification `xml:"SubstitutionModification"`
	CVParam                  []CVParam                  `xml:"cvParam"`
	UserParam                []UserParam                `xml:"userParam"`
}

// PeptideSequence is the amino acid sequence of the (poly)peptide. If a
// substitution modification has been found, the original sequence should be
// reported
type PeptideSequence struct {
	XMLName xml.Name `xml:"PeptideSequence"`
	Value   []byte   `xml:",chardata"`
}

// Modification is a molecule modification specification. If n modifications
// have been found on a peptide, there should be n instances of Modification.
// If multiple modifications are provided as cvParams, it is assumed that the
// modification is ambiguous i.e. one modification or another. A cvParam MUST be
// provided with the identification of the modification sourced from a suitable
// CV e.g. UNIMOD. If the modification is not present in the CV (and this will
// be checked by the semantic validator within a given tolerance window), there
// is a â€œunknown modificationâ€ CV term that MUST be used instead. A neutral
// loss should be defined as an additional CVParam within Modification. If more
// complex information should be given about neutral losses (such as
// presence/absence on particular product ions), this can additionally be
// encoded within the FragmentationArray
type Modification struct {
	XMLName      xml.Name  `xml:"Modification"`
	AvgMassDelta float64   `xml:"avgMassDelta,attr,omitempty"`
	Location     int       `xml:"location,attr,omitempty"`
	Residues     string    `xml:"residues,attr,omitempty"`
	CVParam      []CVParam `xml:"cvParam"`
}

// SubstitutionModification is a modification where one residue is substituted
// by another (amino acid change)
type SubstitutionModification struct {
	XMLName               xml.Name `xml:"SubstitutionModification"`
	AvgMassDelta          float64  `xml:"avgMassDelta,attr,omitempty"`
	Location              int      `xml:"location,attr,omitempty"`
	MonoisotopicMassDelta float64  `xml:"monoisotopicMassDelta,attr,omitempty"`
	OriginalResidue       string   `xml:"originalResidue,attr,omitempty"`
	ReplacementResidue    string   `xml:"replacementResidue,attr,omitempty"`
}

// PeptideEvidence  links a specific Peptide element to a specific position in a
// DBSequence. There MUST only be one PeptideEvidence item per
// Peptide-to-DBSequence-position
type PeptideEvidence struct {
	XMLName             xml.Name    `xml:"PeptideEvidence"`
	DBSequenceRef       string      `xml:"dBSequence_ref,attr,omitempty"`
	End                 int         `xml:"end,attr,omitempty"`
	Frame               string      `xml:"frame,attr,omitempty"`
	ID                  string      `xml:"id,attr,omitempty"`
	IsDecoy             bool        `xml:"isDecoy,attr,omitempty"`
	Name                string      `xml:"name,attr,omitempty"`
	PeptideRef          string      `xml:"peptide_ref,attr,omitempty"`
	Post                string      `xml:"post,attr,omitempty"`
	Pre                 string      `xml:"pre,attr,omitempty"`
	Start               string      `xml:"start,attr,omitempty"`
	TranslationTableRef string      `xml:"translationTable_ref,attr,omitempty"`
	CVParam             []CVParam   `xml:"cvParam"`
	UserParam           []UserParam `xml:"userParam"`
}

// AnalysisCollection is the analyses performed to get the results, which map
// the input and output data sets. Analyses are for example:
// SpectrumIdentification (resulting in peptides) or ProteinDetection
// (assemble proteins from peptides)
type AnalysisCollection struct {
	XMLName                xml.Name                 `xml:"AnalysisCollection"`
	SpectrumIdentification []SpectrumIdentification `xml:"SpectrumIdentification"`
	ProteinDetection       ProteinDetection         `xml:"ProteinDetection"`
}

// SpectrumIdentification is an analysis which tries to identify peptides in
// input spectra, referencing the database searched, the input spectra,
// the output results and the protocol that is run
type SpectrumIdentification struct {
	XMLName                           xml.Name            `xml:"SpectrumIdentification"`
	ActivityDate                      string              `xml:"activityDate,attr,omitempty"`
	ID                                string              `xml:"id,attr,omitempty"`
	Name                              string              `xml:"name,attr,omitempty"`
	SpectrumIdentificationListRef     string              `xml:"spectrumIdentificationList_ref,attr,omitempty"`
	SpectrumIdentificationProtocolRef string              `xml:"spectrumIdentificationProtocol_ref,attr,omitempty"`
	InputSpectra                      []InputSpectra      `xml:"InputSpectra"`
	SearchDatabaseRef                 []SearchDatabaseRef `xml:"SearchDatabaseRef"`
}

// InputSpectra is one of the spectra data sets used
type InputSpectra struct {
	XMLName        xml.Name `xml:"InputSpectra"`
	SpectraDataRef string   `xml:"spectraData_ref,attr,omitempty"`
}

// SearchDatabaseRef is one of the search databases used
type SearchDatabaseRef struct {
	XMLName           xml.Name `xml:"SearchDatabaseRef"`
	SearchDatabaseRef string   `xml:"searchDatabase_ref,attr,omitempty"`
}

// ProteinDetection is an Analysis which assembles a set of peptides
// (e.g. from a spectra search analysis) to proteins
type ProteinDetection struct {
	XMLName                      xml.Name                       `xml:"ProteinDetection"`
	ActivityDate                 string                         `xml:"activityDate,attr,omitempty"`
	ID                           string                         `xml:"id,attr,omitempty"`
	Name                         string                         `xml:"name,attr,omitempty"`
	ProteinDetectionListRef      string                         `xml:"proteinDetectionList_ref,attr,omitempty"`
	ProteinDetectionProtocolRef  string                         `xml:"proteinDetectionProtocol_ref,attr,omitempty"`
	InputSpectrumIdentifications []InputSpectrumIdentifications `xml:"InputSpectrumIdentifications"`
}

// InputSpectrumIdentifications is the lists of spectrum identifications that are input to the protein detection process
type InputSpectrumIdentifications struct {
	XMLName                       xml.Name `xml:"InputSpectrumIdentifications"`
	SpectrumIdentificationListRef string   `xml:"spectrumIdentificationList_ref,attr,omitempty"`
}

// AnalysisProtocolCollection is the collection of protocols which include the
// parameters and settings of the performed analyses
type AnalysisProtocolCollection struct {
	XMLName                        xml.Name                       `xml:"AnalysisProtocolCollection"`
	SpectrumIdentificationProtocol SpectrumIdentificationProtocol `xml:"SpectrumIdentificationProtocol"`
	ProteinDetectionProtocol       ProteinDetectionProtocol       `xml:"ProteinDetectionProtocol"`
}

// SpectrumIdentificationProtocol is the parameters and settings of a
// SpectrumIdentification analysis
type SpectrumIdentificationProtocol struct {
	XMLName                xml.Name               `xml:"SpectrumIdentificationProtocol"`
	AnalysisSoftwareRef    string                 `xml:"analysisSoftware_ref,attr,omitempty"`
	ID                     string                 `xml:"id,attr,omitempty"`
	Name                   string                 `xml:"name,attr,omitempty"`
	SearchType             SearchType             `xml:"SearchType"`
	AdditionalSearchParams AdditionalSearchParams `xml:"AdditionalSearchParams"`
	ModificationParams     ModificationParams     `xml:"ModificationParams"`
	Enzymes                Enzymes                `xml:"Enzymes"`
	MassTable              []MassTable            `xml:"MassTable"`
	FragmentTolerance      FragmentTolerance      `xml:"FragmentTolerance"`
	ParentTolerance        ParentTolerance        `xml:"ParentTolerance"`
	Threshold              Threshold              `xml:"Threshold"`
	DatabaseFilters        DatabaseFilters        `xml:"DatabaseFilters"`
	DatabaseTranslation    DatabaseTranslation    `xml:"DatabaseTranslation"`
}

// ProteinDetectionProtocol is the parameters and settings of a
// ProteinDetection process
type ProteinDetectionProtocol struct {
	XMLName             xml.Name       `xml:"ProteinDetectionProtocol"`
	AnalysisSoftwareRef string         `xml:"analysisSoftware_ref,attr,omitempty"`
	ID                  string         `xml:"id,attr,omitempty"`
	Name                string         `xml:"name,attr,omitempty"`
	AnalysisParams      AnalysisParams `xml:"AnalysisParams"`
	Threshold           Threshold      `xml:"Threshold"`
}

// AnalysisParams is the parameters and settings for the protein detection given
// as CV terms
type AnalysisParams struct {
	XMLName   xml.Name    `xml:"AnalysisParams"`
	CVParam   CVParam     `xml:"cvParam"`
	UserParam []UserParam `xml:"userParam"`
}

// SearchType is the type of search performed e.g. PMF, Tag searches, MS-MS
type SearchType struct {
	XMLName   xml.Name    `xml:"SearchType"`
	CVParam   CVParam     `xml:"cvParam"`
	UserParam []UserParam `xml:"userParam"`
}

// AdditionalSearchParams is the search parameters other than the modifications
// searched
type AdditionalSearchParams struct {
	XMLName   xml.Name    `xml:"AdditionalSearchParams"`
	CVParam   CVParam     `xml:"cvParam"`
	UserParam []UserParam `xml:"userParam"`
}

// ModificationParams is the specification of static/variable modifications
// (e.g. Oxidation of Methionine) that are to be considered in the spectra
// search
type ModificationParams struct {
	XMLName            xml.Name             `xml:"ModificationParams"`
	SearchModification []SearchModification `xml:"SearchModification"`
}

// SearchModification of a search modification as parameter for a spectra
// search. Contains the name of the modification, the mass, the specificity and
// whether it is a static modification
type SearchModification struct {
	XMLName          xml.Name           `xml:"SearchModification"`
	FixedMod         string             `xml:"fixedMod,attr,omitempty"`
	MassDelta        float64            `xml:"massDelta,attr,omitempty"`
	Residues         string             `xml:"residues,attr,omitempty"`
	SpecificityRules []SpecificityRules `xml:"SpecificityRules"`
	CVParam          []CVParam          `xml:"cvParam"`
}

// SpecificityRules is the specificity rules of the searched modification
// including for example the probability of a modification's presence or peptide
// or protein termini. Standard fixed or variable status should be provided by
// the attribute fixedMod
type SpecificityRules struct {
	XMLName   xml.Name    `xml:"SpecificityRules"`
	CVParam   []CVParam   `xml:"cvParam"`
	UserParam []UserParam `xml:"userParam"`
}

// Enzymes is the list of enzymes used in experiment
type Enzymes struct {
	XMLName     xml.Name `xml:"Enzymes"`
	Independent bool     `xml:"independent,attr,omitempty"`
	Enzyme      []Enzyme `xml:"Enzyme"`
}

// Enzyme is the details of an individual cleavage enzyme should be provided by
// giving a regular expression or a CV term if a "standard" enzyme cleavage has
// been performed
type Enzyme struct {
	XMLName         xml.Name   `xml:"Enzyme"`
	CTermGain       string     `xml:"cTermGain,attr,omitempty"`
	ID              string     `xml:"id,attr,omitempty"`
	MinDistance     int        `xml:"minDistance,attr,omitempty"`
	MissedCleavages int        `xml:"missedCleavages,attr,omitempty"`
	NTermGain       string     `xml:"nTermGain,attr,omitempty"`
	Name            string     `xml:"name,attr,omitempty"`
	SemiSpecific    bool       `xml:"semiSpecific,attr,omitempty"`
	SiteRegexp      SiteRegexp `xml:"SiteRegexp"`
	EnzymeName      EnzymeName `xml:"EnzymeName"`
}

// SiteRegexp is the Regular expression for specifying the enzyme cleavage site
type SiteRegexp struct {
	XMLName xml.Name `xml:"SiteRegexp"`
	Value   []byte   `xml:",chardata"`
}

// EnzymeName is the name of the enzyme from a CV
type EnzymeName struct {
	XMLName   xml.Name    `xml:"EnzymeName"`
	CVParam   []CVParam   `xml:"cvParam"`
	UserParam []UserParam `xml:"userParam"`
}

// MassTable is the masses of residues used in the search
type MassTable struct {
	XMLName          xml.Name           `xml:"MassTable"`
	ID               string             `xml:"id,attr,omitempty"`
	MSLevel          []int              `xml:"msLevel,attr,omitempty"`
	Name             string             `xml:"Name,attr,omitempty"`
	Residue          []Residue          `xml:"Residue"`
	AmbiguousResidue []AmbiguousResidue `xml:"AmbiguousResidue"`
	CVParam          []CVParam          `xml:"cvParam"`
	UserParam        []UserParam        `xml:"userParam"`
}

// Residue is the specification of a single residue within the mass table
type Residue struct {
	XMLName xml.Name `xml:"Residue"`
	Code    string   `xml:"code,attr,omitempty"`
	Mass    float64  `xml:"mass,attr,omitempty"`
}

// AmbiguousResidue is the specification of a single residue within the mass
// table
type AmbiguousResidue struct {
	XMLName   xml.Name    `xml:"AmbiguousResidue"`
	Code      string      `xml:"code,attr,omitempty"`
	CVParam   []CVParam   `xml:"cvParam"`
	UserParam []UserParam `xml:"userParam"`
}

// FragmentTolerance is the tolerance of the search given as a plus and minus
// value with units
type FragmentTolerance struct {
	XMLName xml.Name  `xml:"FragmentTolerance"`
	CVParam []CVParam `xml:"cvParam"`
}

// ParentTolerance is the tolerance of the search given as a plus and minus
// value with units
type ParentTolerance struct {
	XMLName xml.Name  `xml:"ParentTolerance"`
	CVParam []CVParam `xml:"cvParam"`
}

// Threshold is applied to determine that a result is significant. If multiple
// terms are used it is assumed that all conditions are satisfied by the passing
// results. Also applied to determine that a result is significant. If multiple
// terms are used it is assumed that all conditions are satisfied by the passing
// results
type Threshold struct {
	XMLName   xml.Name    `xml:"Threshold"`
	CVParam   []CVParam   `xml:"cvParam"`
	UserParam []UserParam `xml:"userParam"`
}

// DatabaseFilters is the specification of filters applied to the database
// searched
type DatabaseFilters struct {
	XMLName xml.Name `xml:"DatabaseFilters"`
	Filter  []Filter `xml:"Filter"`
}

// Filter applied to the search database. The filter MUST include at least one
// of Include and Exclude. If both are used, it is assumed that inclusion is
// performed first.
type Filter struct {
	XMLName    xml.Name   `xml:"Filter"`
	FilterType FilterType `xml:"FilterType"`
	Include    Include    `xml:"Include"`
	Exclude    Exclude    `xml:"Exclude"`
}

// FilterType is the type of filter e.g. database taxonomy filter, pi filter,
// mw filter
type FilterType struct {
	XMLName   xml.Name    `xml:"FilterType"`
	CVParam   []CVParam   `xml:"cvParam"`
	UserParam []UserParam `xml:"userParam"`
}

// Include is all sequences fulfilling the specifed criteria are included
type Include struct {
	XMLName   xml.Name    `xml:"Include"`
	CVParam   []CVParam   `xml:"cvParam"`
	UserParam []UserParam `xml:"userParam"`
}

// Exclude all sequences fulfilling the specifed criteria are excluded
type Exclude struct {
	XMLName   xml.Name    `xml:"Exclude"`
	CVParam   []CVParam   `xml:"cvParam"`
	UserParam []UserParam `xml:"userParam"`
}

// DatabaseTranslation is a specification of how a nucleic acid sequence
// database was translated for searching
type DatabaseTranslation struct {
	XMLName          xml.Name           `xml:"DatabaseTranslation"`
	Frames           string             `xml:"frames,attr,omitempty"`
	TranslationTable []TranslationTable `xml:"TranslationTable"`
}

// TranslationTable is the table used to translate codons into nucleic acids
// e.g. by reference to the NCBI translation table
type TranslationTable struct {
	XMLName xml.Name  `xml:"TranslationTable"`
	CVParam []CVParam `xml:"cvParam"`
}

// DataCollection is the collection of input and output data sets of the
// analyses
type DataCollection struct {
	XMLName      xml.Name     `xml:"DataCollection"`
	Inputs       Inputs       `xml:"Inputs"`
	AnalysisData AnalysisData `xml:"AnalysisData"`
}

// Inputs is the inputs to the analyses including the databases searched, the
// spectral data and the source file converted to mzIdentML
type Inputs struct {
	XMLName        xml.Name         `xml:"Inputs"`
	SourceFile     []SourceFile     `xml:"SourceFile"`
	SearchDatabase []SearchDatabase `xml:"SearchDatabase"`
	SpectraData    []SpectraData    `xml:"SpectraData"`
}

// SearchDatabase is a database for searching mass spectra. Examples include a
// set of amino acid sequence entries, nucleotide databases (e.g. 6 frame
// translated) or annotated spectra libraries
type SearchDatabase struct {
	XMLName                     xml.Name                    `xml:"SearchDatabase"`
	ID                          string                      `xml:"id,attr,omitempty"`
	Location                    string                      `xml:"location,attr,omitempty"`
	Name                        string                      `xml:"name,attr,omitempty"`
	NumDatabaseSequences        string                      `xml:"numDatabaseSequences,attr,omitempty"`
	NumResidues                 string                      `xml:"numResidues,attr,omitempty"`
	ReleaseDate                 string                      `xml:"releaseDate,attr,omitempty"`
	Version                     string                      `xml:"version,attr,omitempty"`
	ExternalFormatDocumentation ExternalFormatDocumentation `xml:"ExternalFormatDocumentation"`
	FileFormat                  FileFormat                  `xml:"FileFormat"`
	DatabaseName                DatabaseName                `xml:"DatabaseName"`
	CVParam                     []CVParam                   `xml:"cvParam"`
}

// ExternalFormatDocumentation is a URI to access documentation and tools to
// interpret the external format of the ExternalData instance. For example,
// XML Schema or static libraries (APIs) to access binary formats
type ExternalFormatDocumentation struct {
	XMLName xml.Name `xml:"ExternalFormatDocumentation"`
	Value   []byte   `xml:",chardata"`
}

// FileFormat is the format of the ExternalData file, for example "tiff" for
// image files.
type FileFormat struct {
	XMLName xml.Name `xml:"FileFormat"`
	CVParam CVParam  `xml:"cvParam"`
}

// DatabaseName is the database name may be given as a cvParam if it maps
// exactly to one of the release databases listed in the CV, otherwise a
// userParam should be used
type DatabaseName struct {
	XMLName   xml.Name    `xml:"DatabaseName"`
	CVParam   []CVParam   `xml:"cvParam"`
	UserParam []UserParam `xml:"userParam"`
}

// SpectraData should be used
type SpectraData struct {
	XMLName                     xml.Name                    `xml:"SpectraData"`
	ID                          string                      `xml:"id,attr,omitempty"`
	Location                    string                      `xml:"location,attr,omitempty"`
	Name                        string                      `xml:"name,attr,omitempty"`
	ExternalFormatDocumentation ExternalFormatDocumentation `xml:"ExternalFormatDocumentation"`
	FileFormat                  FileFormat                  `xml:"FileFormat"`
	SpectrumIDFormat            SpectrumIDFormat            `xml:"SpectrumIDFormat"`
}

// SpectrumIDFormat is the format of the spectrum identifier within the source
// file
type SpectrumIDFormat struct {
	XMLName xml.Name `xml:"SpectrumIDFormat"`
	CVParam CVParam  `xml:"cvParam"`
}

// AnalysisData are sets generated by the analyses, including peptide and
// protein lists
type AnalysisData struct {
	XMLName                    xml.Name                     `xml:"AnalysisData"`
	SpectrumIdentificationList []SpectrumIdentificationList `xml:"SpectrumIdentificationList"`
	ProteinDetectionList       ProteinDetectionList         `xml:"ProteinDetectionList"`
}

// SpectrumIdentificationList is the set of all search results from
// SpectrumIdentification
type SpectrumIdentificationList struct {
	XMLName                      xml.Name                       `xml:"SpectrumIdentificationList"`
	ID                           string                         `xml:"id,attr,omitempty"`
	Name                         string                         `xml:"name,attr,omitempty"`
	NumSequencesSearched         float64                        `xml:"numSequencesSearched,attr,omitempty"`
	FragmentationTable           FragmentationTable             `xml:"FragmentationTable"`
	SpectrumIdentificationResult []SpectrumIdentificationResult `xml:"SpectrumIdentificationResult"`
	CVParam                      []CVParam                      `xml:"cvParam"`
	UserParam                    []UserParam                    `xml:"userParam"`
}

// FragmentationTable is the the types of measures that will be reported in
// generic arrays for each SpectrumIdentificationItem e.g. product ion m/z,
// product ion intensity, product ion m/z error
type FragmentationTable struct {
	XMLName xml.Name  `xml:"FragmentationTable"`
	Measure []Measure `xml:"Measure"`
}

// Measure references to CV terms defining the measures about product ions to be
// reported in SpectrumIdentificationItem
type Measure struct {
	XMLName xml.Name  `xml:"Measure"`
	ID      string    `xml:"id,attr,omitempty"`
	Name    string    `xml:"name,attr,omitempty"`
	CVParam []CVParam `xml:"cvParam"`
}

// SpectrumIdentificationResult is All identifications made from searching one
// spectrum. For PMF data, all peptide identifications will be listed underneath
// as SpectrumIdentificationItems. For MS/MS data, there will be ranked
// SpectrumIdentificationItems corresponding to possible different peptide IDs
type SpectrumIdentificationResult struct {
	XMLName                    xml.Name                     `xml:"SpectrumIdentificationResult"`
	ID                         string                       `xml:"id,attr,omitempty"`
	Name                       string                       `xml:"name,attr,omitempty"`
	SpectraDataRef             string                       `xml:"spectraData_ref,attr,omitempty"`
	SpectrumID                 string                       `xml:"spectrumID,attr,omitempty"`
	SpectrumIdentificationItem []SpectrumIdentificationItem `xml:"SpectrumIdentificationItem"`
}

// SpectrumIdentificationItem is an identification of a single (poly)peptide,
// resulting from querying an input spectra, along with the set of confidence
// values for that identification. PeptideEvidence elements should be given
// for all mappings of the corresponding Peptide sequence within protein
// sequences
type SpectrumIdentificationItem struct {
	XMLName                  xml.Name             `xml:"SpectrumIdentificationItem"`
	CalculatedMassToCharge   float64              `xml:"calculatedMassToCharge,attr,omitempty"`
	CalculatedPI             float64              `xml:"calculatedPI,attr,omitempty"`
	ChargeState              int                  `xml:"chargeState,attr,omitempty"`
	ExperimentalMassToCharge float64              `xml:"experimentalMassToCharge,attr,omitempty"`
	ID                       string               `xml:"id,attr,omitempty"`
	MassTableRef             string               `xml:"massTable_ref,attr,omitempty"`
	Name                     string               `xml:"name,attr,omitempty"`
	PassThreshold            bool                 `xml:"passThreshold,attr,omitempty"`
	PeptideRef               string               `xml:"peptide_ref,attr,omitempty"`
	Rank                     int                  `xml:"rank,attr,omitempty"`
	SampleRef                string               `xml:"sample_ref,attr,omitempty"`
	PeptideEvidenceRef       []PeptideEvidenceRef `xml:"PeptideEvidenceRef"`
	Fragmentation            Fragmentation        `xml:"Fragmentation"`
	CVParam                  []CVParam            `xml:"cvParam"`
}

// PeptideEvidenceRef reference to the PeptideEvidence element identified.
// If a specific sequence can be assigned to multiple proteins and or positions
// in a protein all possible PeptideEvidence elements should be referenced here
type PeptideEvidenceRef struct {
	XMLName            xml.Name `xml:"PeptideEvidenceRef"`
	PeptideEvidenceRef string   `xml:"peptideEvidence_ref,attr,omitempty"`
}

// Fragmentation is the product ions identified in this result
type Fragmentation struct {
	XMLName xml.Name  `xml:"Fragmentation"`
	IonType []IonType `xml:"IonType"`
}

// IonType defines the index of fragmentation ions being reported, importing a
// CV term for the Type of ion e.g. b ion. Example: if b3 b7 b8 and b10 have
// been identified, the index attribute will contain 3 7 8 10, and the
// corresponding values will be reported in parallel arrays below
type IonType struct {
	XMLName       xml.Name        `xml:"IonType"`
	Charge        int             `xml:"charge,attr,omitempty"`
	Index         []string        `xml:"index,attr,omitempty"`
	FragmentArray []FragmentArray `xml:"FragmentArray"`
	CVParam       []CVParam       `xml:"cvParam"`
	UserParam     []UserParam     `xml:"userParam"`
}

// FragmentArray is an array of values for a given type of measure and for a
// particular ion type, in parallel to the index of ions identified
type FragmentArray struct {
	XMLName    xml.Name `xml:"FragmentArray"`
	MeasureRef string   `xml:"measure_ref,attr,omitempty"`
	Values     []string `xml:"values,attr,omitempty"`
}

// ProteinDetectionList is the protein list resulting from a protein detection
// process
type ProteinDetectionList struct {
	XMLName               xml.Name                `xml:"ProteinDetectionList"`
	ID                    string                  `xml:"id,attr,omitempty"`
	Name                  string                  `xml:"name,attr,omitempty"`
	ProteinAmbiguityGroup []ProteinAmbiguityGroup `xml:"ProteinAmbiguityGroup"`
	CVParam               []CVParam               `xml:"cvParam"`
	UserParam             []UserParam             `xml:"userParam"`
}

// ProteinAmbiguityGroup is a set of logically related results from a protein
// detection, for example to represent conflicting assignments of peptides to
// proteins
type ProteinAmbiguityGroup struct {
	XMLName                    xml.Name                     `xml:"ProteinAmbiguityGroup"`
	ID                         string                       `xml:"id,attr,omitempty"`
	Name                       string                       `xml:"name,attr,omitempty"`
	ProteinDetectionHypothesis []ProteinDetectionHypothesis `xml:"ProteinDetectionHypothesis"`
	CVParam                    []CVParam                    `xml:"cvParam"`
	UserParam                  []UserParam                  `xml:"userParam"`
}

// ProteinDetectionHypothesis is a single result of the ProteinDetection
// analysis (i.e. a protein)
type ProteinDetectionHypothesis struct {
	XMLName           xml.Name            `xml:"ProteinDetectionHypothesis"`
	DBSquenceRef      string              `xml:"dBSequence_ref,attr,omitempty"`
	ID                string              `xml:"id,attr,omitempty"`
	Name              string              `xml:"name,attr,omitempty"`
	PassThreshold     bool                `xml:"passThreasold,attr,omitempty"`
	PeptideHypothesis []PeptideHypothesis `xml:"PeptideHypothesis"`
	CVParam           []CVParam           `xml:"cvParam"`
	UserParam         []UserParam         `xml:"userParam"`
}

// PeptideHypothesis is the evidence on which this ProteinHypothesis is based by
// reference to a PeptideEvidence element
type PeptideHypothesis struct {
	XMLName                       xml.Name                        `xml:"PeptideHypothesis"`
	PeptideEvidenceRef            string                          `xml:"peptideEvidence_ref,attr,omitempty"`
	SpectrumIdentificationItemRef []SpectrumIdentificationItemRef `xml:"SpectrumIdentificationItemRef,attr,omitempty"`
}

// SpectrumIdentificationItemRef Reference(s) to the SpectrumIdentificationItem
// element(s) that support the given PeptideEvidence element. Using these
// references it is possible to indicate which spectra were actually accepted as
// evidence for this peptide identification in the given protein
type SpectrumIdentificationItemRef struct {
	XMLName                       xml.Name `xml:"SpectrumIdentificationItemRef"`
	SpectrumIdentificationItemRef string   `xml:"spectrumIdentificationItem_ref,attr,omitempty"`
}

// BibliographicReference is any bibliographic references associated with the
// file
type BibliographicReference struct {
	XMLName     xml.Name `xml:"BibliographicReference"`
	Authors     string   `xml:"authors,attr,omitempty"`
	DOI         string   `xml:"doi,attr,omitempty"`
	Editor      string   `xml:"editor,attr,omitempty"`
	ID          string   `xml:"id,attr,omitempty"`
	Issue       string   `xml:"issue,attr,omitempty"`
	Name        string   `xml:"name,attr,omitempty"`
	Pages       string   `xml:"pages,attr,omitempty"`
	Publication string   `xml:"publication,attr,omitempty"`
	Publisher   string   `xml:"publisher,attr,omitempty"`
	Title       string   `xml:"title,attr,omitempty"`
	Volume      string   `xml:"volume,attr,omitempty"`
	Year        string   `xml:"year,attr,omitempty"`
}
