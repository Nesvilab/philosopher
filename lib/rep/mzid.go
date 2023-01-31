package rep

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Nesvilab/philosopher/lib/id"

	"github.com/Nesvilab/philosopher/lib/dat"
	"github.com/Nesvilab/philosopher/lib/psi"
)

// MzIdentMLReport creates a MzIdentML structure to be encoded
func (e Evidence) MzIdentMLReport(version, database string) {

	var mzid psi.MzIdentML

	t := time.Now()
	//var idCounter = 0

	// collect source file names
	var sourceMap = make(map[string]uint8)
	var sources []string
	for _, i := range e.PSM {
		s := strings.Split(i.Spectrum, ".")
		sourceMap[s[0]]++
	}

	for i := range sourceMap {
		sources = append(sources, i)
	}

	sort.Strings(sources)

	// load the database
	var dtb dat.Base
	dtb.Restore()

	// spectra evidence reference map
	var specRef = make(map[id.SpectrumType]string)

	// peptide evidence reference map
	var pepRef = make(map[string]string)

	// protein evidence reference map
	var proRef = make(map[string]string)

	// Header
	//mzid.Name = "foo"
	mzid.ID = "Philosopher"
	mzid.Version = "1.2.0"
	mzid.CreationDate = t.Format(time.ANSIC)
	mzid.Xmlns = "http://psidev.info/psi/pi/mzIdentML/1.2"
	mzid.XmlnsXsi = "http://www.w3.org/2001/XMLSchema-instance"
	mzid.XsiSchemaLocation = "http://psidev.info/psi/pi/mzIdentML/1.2 http://www.psidev.info/files/mzIdentML1.2.0.xsd"

	// CVlist
	mzid.CvList.CV = append(mzid.CvList.CV, psi.CV{ID: "PSI-MS", URI: "https://raw.githubusercontent.com/HUPO-PSI/psi-ms-CV/master/psi-ms.obo", FullName: "PSI-MS"})
	mzid.CvList.CV = append(mzid.CvList.CV, psi.CV{ID: "UNIMOD", URI: "http://www.unimod.org/obo/unimod.obo", FullName: "UNIMOD"})
	mzid.CvList.CV = append(mzid.CvList.CV, psi.CV{ID: "UO", URI: "https://raw.githubusercontent.com/bio-ontology-research-group/unit-ontology/master/unit.obo", FullName: "UNIT-ONTOLOGY"})
	mzid.CvList.CV = append(mzid.CvList.CV, psi.CV{ID: "PRIDE", URI: "https://github.com/PRIDE-Utilities/pride-ontology/blob/master/pride_cv.obo", FullName: "PRIDE"})
	mzid.CvList.Count = len(mzid.CvList.CV)

	// AnalysisSoftwareList
	aa := &psi.AnalysisSoftware{
		ID:      "Philosopher",
		Name:    "Philosopher toolkit",
		URI:     "https://philosopher.nesvilab.org",
		Version: version,
		ContactRole: psi.ContactRole{
			ContactRef: "PS_DEV",
			Role: psi.Role{
				CVParam: psi.CVParam{
					CVRef:     "PSI-MS",
					Accession: "MS:1001267",
					Name:      "software vendor",
				},
			},
		},
		SoftwareName: psi.SoftwareName{
			CVParam: psi.CVParam{
				CVRef:     "PSI-MS",
				Accession: "XXXX",
				Name:      "Philosopher",
			},
		},
		Customizations: psi.Customizations{
			Value: "No customizations",
		},
	}
	mzid.AnalysisSoftwareList.AnalysisSoftware = append(mzid.AnalysisSoftwareList.AnalysisSoftware, *aa)

	//Provider
	provider := &psi.Provider{
		ID: "PROVIDER",
		ContactRole: psi.ContactRole{
			ContactRef: "Philosopher_Author_FVL",
			Role: psi.Role{
				CVParam: psi.CVParam{
					CVRef:     "PSI-MS",
					Accession: "MS:1001271",
					Name:      "researcher",
				},
			},
		},
	}
	mzid.Provider = *provider

	// AuditCollection

	auditCol := &psi.AuditCollection{
		Person: psi.Person{
			ID:        "Philosopher_Author_FVL",
			LastName:  "da Veiga Leprevost",
			FirstName: "Felipe",
			CVParam: []psi.CVParam{
				{
					Name:      "contact email",
					Value:     "felipevl@umich.edu",
					CVRef:     "PSI-MS",
					Accession: "MS:1000589",
				},
				{
					Name:      "contact URL",
					Value:     "http://nesvilab.org",
					CVRef:     "PSI-MS",
					Accession: "MS:1000588",
				},
			},
			Affiliation: []psi.Affiliation{
				{
					OrganizationRef: "University of Michigan",
				},
			},
		},
		Organization: psi.Organization{
			ID:   "Nesvilab",
			Name: "Proteomics and Integrative Bioinformatics Lab",
			CVParam: []psi.CVParam{
				{
					Name:      "contact name",
					Value:     "Alexey I. Nesvizhskii",
					CVRef:     "PSI-MS",
					Accession: "MS:1000586",
				},
				{
					Name:      "contact address",
					Value:     "1301 Catherinse St., Ann Arbor, MI",
					CVRef:     "PSI-MS",
					Accession: "MS:1000587",
				},
				{
					Name:      "contact URL",
					Value:     "http://nesvilab.org",
					CVRef:     "PSI-MS",
					Accession: "MS:1000588",
				},
				{
					Name:      "contact email",
					Value:     "nesvi@med.umich.edu",
					CVRef:     "PSI-MS",
					Accession: "MS:1000589",
				},
			},
		},
	}
	mzid.AuditCollection = *auditCol

	// SequenceCollection - DBSequence
	idCounter := 0
	var seqs []psi.DBSequence
	for _, i := range dtb.Records {

		idCounter++

		db := &psi.DBSequence{
			ID:                fmt.Sprintf("DB_%d", idCounter),
			Accession:         i.ID,
			SearchDatabaseRef: dtb.FileName,
			CVParam: []psi.CVParam{
				{
					CVRef:     "PSI-MS",
					Accession: "MS:1001088",
					Name:      "protein description",
					Value:     i.Description,
				},
				{
					CVRef:     "PSI-MS",
					Accession: "MS:1001344",
					Name:      "AA sequence",
				},
			},
			Seq: psi.Seq{
				Value: i.Sequence,
			},
		}

		proRef[i.ID] = fmt.Sprintf("DB_%d", idCounter)
		seqs = append(seqs, *db)
	}
	mzid.SequenceCollection.DBSequence = seqs

	// SequenceCollection - Peptide
	var peps []psi.Peptide
	for _, i := range e.Peptides {

		p := psi.Peptide{
			ID: i.Sequence,
			PeptideSequence: psi.PeptideSequence{
				Value: i.Sequence,
			},
		}

		for _, j := range i.Modifications.IndexSlice {
			if j.Name != "Unknown" {
				mod := psi.Modification{
					//AvgMassDelta:          j.AverageMass,
					//MonoIsotopicMassDelta: j.MonoIsotopicMass,
					Residues: j.AminoAcid,
					Location: strconv.Itoa(j.Position),
					CVParam: []psi.CVParam{
						{
							CVRef:     "UNIMOD",
							Accession: j.ID,
							Name:      j.Name,
						},
					},
				}

				if mod.Residues == "N-term" {
					mod.Residues = ""
				}

				p.Modification = append(p.Modification, mod)
			}
		}

		peps = append(peps, p)
	}
	mzid.SequenceCollection.Peptide = peps

	// SequenceCollection - PeptideEvidence
	var pevs []psi.PeptideEvidence
	idCounter = 0
	for _, i := range e.PSM {

		idCounter++

		evi := psi.PeptideEvidence{
			DBSequenceRef: i.ProteinID,
			ID:            fmt.Sprintf("PepEv_%d", idCounter),
			IsDecoy:       strconv.FormatBool(i.IsDecoy),
			PeptideRef:    i.Peptide,
			Pre:           string(i.PrevAA),
			Post:          string(i.NextAA),
		}

		pepRef[i.Peptide] = fmt.Sprintf("PepEv_%d", idCounter)
		pevs = append(pevs, evi)
	}
	mzid.SequenceCollection.PeptideEvidence = pevs

	// AnalysisCollection
	idCounter = 0
	ac := &psi.AnalysisCollection{}
	for _, i := range sources {

		idCounter++

		si := &psi.SpectrumIdentification{
			SpectrumIdentificationListRef:     fmt.Sprintf("SIL_%d", idCounter),
			ID:                                fmt.Sprintf("SpecIdent_%d", idCounter),
			SpectrumIdentificationProtocolRef: fmt.Sprintf("SearchProtocol_%d", idCounter),
			InputSpectra: []psi.InputSpectra{
				{
					SpectraDataRef: i,
				},
			},
			SearchDatabaseRef: []psi.SearchDatabaseRef{
				{
					SearchDatabaseRef: dtb.FileName,
				},
			},
		}

		ac.SpectrumIdentification = append(ac.SpectrumIdentification, *si)
	}

	ac.ProteinDetection = psi.ProteinDetection{
		ProteinDetectionProtocolRef: "Philosopher_protocol",
		ProteinDetectionListRef:     "Protein Groups",
		ID:                          "Phi_1",
	}
	mzid.AnalysisCollection = *ac

	// AnalysisProtocolCollection
	apc := &psi.AnalysisProtocolCollection{
		SpectrumIdentificationProtocol: []psi.SpectrumIdentificationProtocol{
			{
				AnalysisSoftwareRef: "DatabaseSearch_ID",
				ID:                  "Search_Protocol_1",
				SearchType: psi.SearchType{
					CVParam: psi.CVParam{
						CVRef:     "PSI-MS",
						Accession: "MS:1001083",
						Name:      "ms-ms search",
					},
				},
				AdditionalSearchParams: psi.AdditionalSearchParams{
					CVParam: []psi.CVParam{
						{
							CVRef:     "PSI-MS",
							Accession: "MS:1001211",
							Name:      "parent mass type mono",
						},
						{
							CVRef:     "PSI-MS",
							Accession: "MS:1001256",
							Name:      "fragment mass type mono",
						},
						{
							CVRef:     "PSI-MS",
							Accession: "MS:1002492",
							Name:      "consensus scoring",
						},
						{
							CVRef:     "PSI-MS",
							Accession: "MS:1002490",
							Name:      "peptide-level scoring",
						},
						{
							CVRef:     "PSI-MS",
							Accession: "MS:1002497",
							Name:      "group PSMs by sequence with modifications",
						},
						{
							CVRef:     "PSI-MS",
							Accession: "MS:1002491",
							Name:      "modification localization scoring",
						},
					},
					UserParam: []psi.UserParam{
						{
							Name:  "MSFragger",
							Value: e.Parameters.MSFragger,
						},
						{
							Name:  "database_name",
							Value: e.Parameters.DatabaseName,
						},
						{
							Name:  "precursor_mass_lower",
							Value: e.Parameters.PrecursorMassLower,
						},
						{
							Name:  "precursor_mass_upper",
							Value: e.Parameters.PrecursorMassUpper,
						},
						{
							Name:  "precursor_mass_units",
							Value: e.Parameters.PrecursorMassUnits,
						},
						{
							Name:  "precursor_true_tolerance",
							Value: e.Parameters.PrecursorTrueTolerance,
						},
						{
							Name:  "precursor_true_units",
							Value: e.Parameters.PrecursorTrueUnits,
						},
						{
							Name:  "fragment_mass_tolerance",
							Value: e.Parameters.FragmentMassTolerance,
						},
						{
							Name:  "fragment_mass_units",
							Value: e.Parameters.FragmentMassUnits,
						},
						{
							Name:  "calibrate_mass",
							Value: e.Parameters.CalibrateMass,
						},
						{
							Name:  "ms1_tolerance_mad",
							Value: e.Parameters.Ms1ToleranceMad,
						},
						{
							Name:  "ms2_tolerance_mad",
							Value: e.Parameters.Ms2ToleranceMad,
						},
						{
							Name:  "evaluate_mass_calibration",
							Value: e.Parameters.EvaluateMassCalibration,
						},
						{
							Name:  "isotope_error",
							Value: e.Parameters.IsotopeError,
						},
						{
							Name:  "mass_offsets",
							Value: e.Parameters.MassOffsets,
						},
						{
							Name:  "precursor_mass_mode",
							Value: e.Parameters.PrecursorMassMode,
						},
						{
							Name:  "shifted_ions",
							Value: e.Parameters.ShiftedIons,
						},
						{
							Name:  "shifted_ions_exclude_ranges",
							Value: e.Parameters.ShiftedIonsExcludeRanges,
						},
						{
							Name:  "fragment_ion_series",
							Value: e.Parameters.FragmentIonSeries,
						},
						{
							Name:  "search_enzyme_name",
							Value: e.Parameters.SearchEnzymeName,
						},
						{
							Name:  "search_enzyme_cutafter",
							Value: e.Parameters.SearchEnzymeCutafter,
						},
						{
							Name:  "search_enzyme_butnotafter",
							Value: e.Parameters.SearchEnzymeButnotafter,
						},
						{
							Name:  "num_enzyme_termini",
							Value: e.Parameters.NumEnzymeTermini,
						},
						{
							Name:  "allowed_missed_cleavage",
							Value: e.Parameters.AllowedMissedCleavage,
						},
						{
							Name:  "clip_nTerm_M",
							Value: e.Parameters.ClipNTermM,
						},
						{
							Name:  "allow_multiple_variable_mods_on_residue",
							Value: e.Parameters.AllowMultipleVariableModsOnResidue,
						},
						{
							Name:  "max_variable_mods_per_mod",
							Value: e.Parameters.MaxVariableModsPerMod,
						},
						{
							Name:  "max_variable_mods_combinations",
							Value: e.Parameters.MaxVariableModsCombinations,
						},
						{
							Name:  "output_format",
							Value: e.Parameters.OutputFormat,
						},
						{
							Name:  "output_report_topN",
							Value: e.Parameters.OutputReportTopN,
						},
						{
							Name:  "output_max_expect",
							Value: e.Parameters.OutputMaxExpect,
						},
						{
							Name:  "report_alternative_proteins",
							Value: e.Parameters.ReportAlternativeProteins,
						},
						{
							Name:  "override_charge",
							Value: e.Parameters.OverrideCharge,
						},
						{
							Name:  "precursor_charge",
							Value: e.Parameters.PrecursorCharge,
						},
						{
							Name:  "digest_min_length",
							Value: e.Parameters.DigestMinLength,
						},
						{
							Name:  "digest_max_length",
							Value: e.Parameters.DigestMaxLength,
						},
						{
							Name:  "digest_mass_range",
							Value: e.Parameters.DigestMassRange,
						},
						{
							Name:  "max_fragment_charge",
							Value: e.Parameters.MaxFragmentCharge,
						},
						{
							Name:  "track_zero_topN",
							Value: e.Parameters.TrackZeroTopN,
						},
						{
							Name:  "zero_bin_accept_expect",
							Value: e.Parameters.ZeroBinAcceptExpect,
						},
						{
							Name:  "zero_bin_mult_expect",
							Value: e.Parameters.ZeroBinMultExpect,
						},
						{
							Name:  "add_topN_complementary",
							Value: e.Parameters.AddTopNComplementary,
						},
						{
							Name:  "minimum_peaks",
							Value: e.Parameters.MinimumPeaks,
						},
						{
							Name:  "use_topN_peaks",
							Value: e.Parameters.UseTopNPeaks,
						},
						{
							Name:  "min_fragments_modelling",
							Value: e.Parameters.MinFragmentsModelling,
						},
						{
							Name:  "min_matched_fragments",
							Value: e.Parameters.MinMatchedFragments,
						},
						{
							Name:  "minimum_ratio",
							Value: e.Parameters.MinimumRatio,
						},
						{
							Name:  "clear_mz_range",
							Value: e.Parameters.ClearMzRange,
						},
						{
							Name:  "variable_mod_01",
							Value: e.Parameters.VariableMod01,
						},
						{
							Name:  "variable_mod_02",
							Value: e.Parameters.VariableMod02,
						},
						{
							Name:  "add_C_cysteine",
							Value: e.Parameters.Cysteine,
						},
						{
							Name:  "add_Cterm_peptide",
							Value: e.Parameters.CTermPeptide,
						},
						{
							Name:  "add_Cterm_protein",
							Value: e.Parameters.CTermProtein,
						},
						{
							Name:  "add_D_aspartic_acid",
							Value: e.Parameters.AsparticAcid,
						},
						{
							Name:  "add_E_glutamic_acid",
							Value: e.Parameters.GlutamicAcid,
						},
						{
							Name:  "add_F_phenylalanine",
							Value: e.Parameters.Phenylalanine,
						},
						{
							Name:  "add_G_glycine",
							Value: e.Parameters.Glycine,
						},
						{
							Name:  "add_H_histidine",
							Value: e.Parameters.Histidine,
						},
						{
							Name:  "add_I_isoleucine",
							Value: e.Parameters.Isoleucine,
						},
						{
							Name:  "add_K_lysine",
							Value: e.Parameters.Lysine,
						},
						{
							Name:  "add_L_leucine",
							Value: e.Parameters.Leucine,
						},
						{
							Name:  "add_M_methionine",
							Value: e.Parameters.Methionine,
						},
						{
							Name:  "add_N_asparagine",
							Value: e.Parameters.Asparagine,
						},
						{
							Name:  "add_Nterm_peptide",
							Value: e.Parameters.NTermPeptide,
						},
						{
							Name:  "add_Nterm_protein",
							Value: e.Parameters.NTermProtein,
						},
						{
							Name:  "add_P_proline",
							Value: e.Parameters.Proline,
						},
						{
							Name:  "add_Q_glutamine",
							Value: e.Parameters.GlutamicAcid,
						},
						{
							Name:  "add_R_arginine",
							Value: e.Parameters.Arginine,
						},
						{
							Name:  "add_S_serine",
							Value: e.Parameters.Serine,
						},
						{
							Name:  "add_T_threonine",
							Value: e.Parameters.Threonine,
						},
						{
							Name:  "add_V_valine",
							Value: e.Parameters.Valine,
						},
						{
							Name:  "add_W_tryptophan",
							Value: e.Parameters.Tryptophan,
						},
						{
							Name:  "add_Y_tyrosine",
							Value: e.Parameters.Tyrosine,
						},
					},
				},
			},
		},
	}

	mzid.AnalysisProtocolCollection = *apc

	// DataCollection
	dta := psi.DataCollection{}

	for _, i := range sources {
		sf := &psi.SourceFile{
			ID:       i,
			Location: i,
			Name:     i,
		}

		dta.Inputs.SourceFile = append(dta.Inputs.SourceFile, *sf)
	}

	// DataCollection - Input - SearchDatabase
	sdb := &psi.SearchDatabase{
		ID:                   database,
		NumDatabaseSequences: len(dtb.Records),
		Location:             database,
		FileFormat: psi.FileFormat{
			CVParam: psi.CVParam{
				CVRef:     "PSI-MS",
				Accession: "MS:1001348",
				Name:      "FASTA format",
			},
		},
		DatabaseName: psi.DatabaseName{
			CVParam: psi.CVParam{
				CVRef:     "PSI-MS",
				Accession: "MS:1001073",
				Name:      "database type amino acid",
			},
			UserParam: psi.UserParam{
				Name: database,
			},
		},
	}
	mzid.DataCollection.Inputs.SearchDatabase = append(mzid.DataCollection.Inputs.SearchDatabase, *sdb)

	// DataCollection - Input - SpectraData
	for _, i := range sources {
		sd := &psi.SpectraData{
			Location: "./",
			ID:       i,
			Name:     i,
			FileFormat: psi.FileFormat{
				CVParam: psi.CVParam{
					CVRef:     "PSI-MS",
					Accession: "MS:1000584",
					Name:      "mzML format",
				},
			},
		}
		mzid.DataCollection.Inputs.SpectraData = append(mzid.DataCollection.Inputs.SpectraData, *sd)
	}

	// DataCollection - AnalysisData
	ad := &psi.AnalysisData{
		SpectrumIdentificationList: []psi.SpectrumIdentificationList{
			{
				ID: "SIL_1",
				FragmentationTable: psi.FragmentationTable{
					Measure: []psi.Measure{
						{
							ID: "Measure_MZ",
							CVParam: []psi.CVParam{
								{
									CVRef:         "PSI-MS",
									Accession:     "MS:1001225",
									Name:          "product ion m/z",
									UnitCvRef:     "PSI-MS",
									UnitAccession: "MS:1000040",
									UnitName:      "m/z",
								},
							},
						},
						{
							ID: "Measure_Int",
							CVParam: []psi.CVParam{
								{
									CVRef:         "PSI-MS",
									Accession:     "MS:1001226",
									Name:          "product ion intensity",
									UnitCvRef:     "PSI-MS",
									UnitAccession: "MS:1000131",
									UnitName:      "number of detector counts",
								},
							},
						},
						{
							ID: "Measure_Error",
							CVParam: []psi.CVParam{
								{
									CVRef:         "PSI-MS",
									Accession:     "MS:1001227",
									Name:          "product ion m/z error",
									UnitCvRef:     "PSI-MS",
									UnitAccession: "MS:1000040",
									UnitName:      "m/z",
								},
							},
						},
					},
				},
			},
		},
	}

	// DataCollection - SpectrumIdentificationResult
	idCounter = 0
	for i := 0; i <= len(sources)-1; i++ {
		for _, j := range e.PSM {
			if strings.Contains(j.Spectrum, sources[i]) {

				idCounter++

				sir := &psi.SpectrumIdentificationResult{
					SpectraDataRef: sources[i],
					ID:             fmt.Sprintf("Spectrum_%d", idCounter),
					SpectrumID:     fmt.Sprintf("%d", j.Index),
					SpectrumIdentificationItem: []psi.SpectrumIdentificationItem{
						{
							PassThreshold:            "true",
							Rank:                     j.HitRank,
							PeptideRef:               j.Peptide,
							CalculatedMassToCharge:   j.CalcNeutralPepMass,
							ChargeState:              j.AssumedCharge,
							ExperimentalMassToCharge: j.PrecursorNeutralMass,
							ID:                       fmt.Sprintf("SII_%d", j.HitRank),
							PeptideEvidenceRef: []psi.PeptideEvidenceRef{
								{
									PeptideEvidenceRef: pepRef[j.Peptide],
								},
							},
							Fragmentation: psi.Fragmentation{
								IonType: []psi.IonType{
									{
										CVParam: []psi.CVParam{
											{},
										},
										UserParam: []psi.UserParam{
											{},
										},
									},
								},
							},
							CVParam: []psi.CVParam{
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1000796",
									Name:      "spectrum title",
									Value:     j.SpectrumFileName().Str(),
								},
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1001192",
									Name:      "Expect value",
									Value:     fmt.Sprintf("%f", j.Expectation),
								},
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1000882",
									Name:      "protein",
									Value:     j.ProteinID,
								},
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1000886",
									Name:      "protein name",
									Value:     j.ProteinDescription,
								},
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1000934",
									Name:      "gene name",
									Value:     j.GeneName,
								},
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1000888",
									Name:      "modified peptide sequence",
									Value:     j.ModifiedPeptide,
								},
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1000894",
									Name:      "retention time",
									Value:     fmt.Sprintf("%f", j.RetentionTime),
								},
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1001976",
									Name:      "delta M",
									Value:     fmt.Sprintf("%f", j.Massdiff),
								},
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1002357",
									Name:      "PSM-level probability",
									Value:     fmt.Sprintf("%f", j.Probability),
								},
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1002252",
									Name:      "Comet:xcorr",
									Value:     fmt.Sprintf("%f", j.Xcorr),
								},
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1002253",
									Name:      "Comet:deltacn",
									Value:     fmt.Sprintf("%f", j.DeltaCN),
								},
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1002254",
									Name:      "Comet:deltacnstar",
									Value:     fmt.Sprintf("%f", j.DeltaCNStar),
								},
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1002255",
									Name:      "Comet:spscore",
									Value:     fmt.Sprintf("%f", j.SPScore),
								},
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1002256",
									Name:      "Comet:sprank",
									Value:     fmt.Sprintf("%f", j.SPRank),
								},
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1001331",
									Name:      "X! Tandem:hyperscore",
									Value:     fmt.Sprintf("%f", j.Hyperscore),
								},
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1002217",
									Name:      "decoy peptide",
									Value:     fmt.Sprintf("%v", j.IsDecoy),
								},
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1001843",
									Name:      "MS1 feature maximum intensity",
									Value:     fmt.Sprintf("%f", j.Intensity),
								},
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1001363",
									Name:      "peptide unique to one protein",
									Value:     fmt.Sprintf("%v", j.IsUnique),
								},
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1003015",
									Name:      "razor peptide",
									Value:     fmt.Sprintf("%v", j.IsURazor),
								},
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1002616",
									Name:      "TMT reagent 126",
									Value:     fmt.Sprintf("%f", j.Labels.Channel1.Intensity),
								},
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1002763",
									Name:      "TMT reagent 127N",
									Value:     fmt.Sprintf("%f", j.Labels.Channel2.Intensity),
								},
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1002764",
									Name:      "TMT reagent 127C",
									Value:     fmt.Sprintf("%f", j.Labels.Channel3.Intensity),
								},
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1002765",
									Name:      "TMT reagent 128N",
									Value:     fmt.Sprintf("%f", j.Labels.Channel4.Intensity),
								},
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1002766",
									Name:      "TMT reagent 128C",
									Value:     fmt.Sprintf("%f", j.Labels.Channel5.Intensity),
								},
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1002767",
									Name:      "TMT reagent 129N",
									Value:     fmt.Sprintf("%f", j.Labels.Channel6.Intensity),
								},
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1002768",
									Name:      "TMT reagent 129C",
									Value:     fmt.Sprintf("%f", j.Labels.Channel7.Intensity),
								},
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1002769",
									Name:      "TMT reagent 130N",
									Value:     fmt.Sprintf("%f", j.Labels.Channel8.Intensity),
								},
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1002770",
									Name:      "TMT reagent 130C",
									Value:     fmt.Sprintf("%f", j.Labels.Channel9.Intensity),
								},
								{
									CVRef:     "PSI-MS",
									Accession: "MS:1002621",
									Name:      "TMT reagent 131",
									Value:     fmt.Sprintf("%f", j.Labels.Channel10.Intensity),
								},
							},
							UserParam: []psi.UserParam{
								{
									Name:  "entry name",
									Value: j.EntryName,
								},
								{
									Name:  "TMT reagent 126 Label",
									Value: j.Labels.Channel1.Name,
								},
								{
									Name:  "TMT reagent 17N Label",
									Value: j.Labels.Channel2.Name,
								},
								{
									Name:  "TMT reagent 127C Label",
									Value: j.Labels.Channel3.Name,
								},
								{
									Name:  "TMT reagent 128N Label",
									Value: j.Labels.Channel4.Name,
								},
								{
									Name:  "TMT reagent 128C Label",
									Value: j.Labels.Channel5.Name,
								},
								{
									Name:  "TMT reagent 129N Label",
									Value: j.Labels.Channel6.Name,
								},
								{
									Name:  "TMT reagent 129C Label",
									Value: j.Labels.Channel7.Name,
								},
								{
									Name:  "TMT reagent 130N Label",
									Value: j.Labels.Channel8.Name,
								},
								{
									Name:  "TMT reagent 130C Label",
									Value: j.Labels.Channel9.Name,
								},
								{
									Name:  "TMT reagent 131 Label",
									Value: j.Labels.Channel10.Name,
								},
							},
						},
					},
				}

				specRef[j.SpectrumFileName()] = fmt.Sprintf("Spectrum_%d", idCounter)
				ad.SpectrumIdentificationList[0].SpectrumIdentificationResult = append(ad.SpectrumIdentificationList[0].SpectrumIdentificationResult, *sir)
			}
		}
	}

	// DataCollection - ProteinDetectionList
	idCounter = 0
	if len(e.Proteins) > 0 {
		pdl := &psi.ProteinDetectionList{
			ID: "protein groups",
		}

		var groupsMap = make(map[int]uint8)
		var groups []int

		for _, i := range e.Proteins {
			groupsMap[int(i.ProteinGroup)] = 0
		}

		for i := range groupsMap {
			groups = append(groups, i)
		}

		sort.Ints(groups)

		for _, i := range groups {

			idCounter++

			pag := &psi.ProteinAmbiguityGroup{
				ID: fmt.Sprintf("%d", i),
			}

			for _, j := range e.Proteins {
				if int(j.ProteinGroup) == i {

					pdh := &psi.ProteinDetectionHypothesis{
						ID:                j.ProteinSubGroup,
						PassThreshold:     "true",
						DBSquenceRef:      proRef[j.ProteinID],
						PeptideHypothesis: []psi.PeptideHypothesis{},
						CVParam: []psi.CVParam{
							{
								CVRef:     "PSI-MS",
								Accession: "MS:1000796",
								Name:      "spectrum title",
								Value:     "",
							},
						},
						UserParam: []psi.UserParam{
							{
								Name:  "original protein header",
								Value: j.OriginalHeader,
							},
							{
								Name:  "partial header",
								Value: j.PartHeader,
							},
						},
					}

					for _, k := range j.TotalPeptideIons {
						peph := &psi.PeptideHypothesis{
							PeptideEvidenceRef: pepRef[k.Sequence],
						}
						_ = k

						for l := range k.Spectra {
							siir := psi.SpectrumIdentificationItemRef{
								SpectrumIdentificationItemRef: specRef[l],
							}

							peph.SpectrumIdentificationItemRef = append(peph.SpectrumIdentificationItemRef, siir)
						}

						pdh.PeptideHypothesis = append(pdh.PeptideHypothesis, *peph)
					}

					pag.ProteinDetectionHypothesis = append(pag.ProteinDetectionHypothesis, *pdh)
				}
			}

			pdl.ProteinAmbiguityGroup = append(pdl.ProteinAmbiguityGroup, *pag)
		}

		ad.ProteinDetectionList = *pdl
	}

	mzid.DataCollection.AnalysisData = *ad

	// Burn!
	mzid.Write()

}
