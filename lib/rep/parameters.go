package rep

import (
	"github.com/nesvilab/philosopher/lib/spc"
)

// AssembleSearchParameters organizes the aprameters defined by the search engine
func (e *Evidence) AssembleSearchParameters(params []spc.Parameter) {

	for _, i := range params {
		if i.Name == "MSFragger" {
			e.Parameters.MSFragger = i.Value
		} else if i.Name == "database_name" {
			e.Parameters.DatabaseName = i.Value
		} else if i.Name == "num_threads" {
			e.Parameters.NumThreads = i.Value
		} else if i.Name == "precursor_mass_lower" {
			e.Parameters.PrecursorMassLower = i.Value
		} else if i.Name == "precursor_mass_upper" {
			e.Parameters.PrecursorMassUpper = i.Value
		} else if i.Name == "precursor_mass_units" {
			e.Parameters.PrecursorMassUnits = i.Value
		} else if i.Name == "precursor_true_tolerance" {
			e.Parameters.PrecursorTrueTolerance = i.Value
		} else if i.Name == "precursor_true_units" {
			e.Parameters.PrecursorTrueUnits = i.Value
		} else if i.Name == "fragment_mass_tolerance" {
			e.Parameters.FragmentMassTolerance = i.Value
		} else if i.Name == "fragment_mass_units" {
			e.Parameters.FragmentMassUnits = i.Value
		} else if i.Name == "calibrate_mass" {
			e.Parameters.CalibrateMass = i.Value
		} else if i.Name == "ms1_tolerance_mad" {
			e.Parameters.Ms1ToleranceMad = i.Value
		} else if i.Name == "ms2_tolerance_mad" {
			e.Parameters.Ms2ToleranceMad = i.Value
		} else if i.Name == "evaluate_mass_calibration" {
			e.Parameters.EvaluateMassCalibration = i.Value
		} else if i.Name == "isotope_error" {
			e.Parameters.IsotopeError = i.Value
		} else if i.Name == "mass_offsets" {
			e.Parameters.MassOffsets = i.Value
		} else if i.Name == "precursor_mass_mode" {
			e.Parameters.PrecursorMassMode = i.Value
		} else if i.Name == "shifted_ions" {
			e.Parameters.ShiftedIons = i.Value
		} else if i.Name == "shifted_ions_exclude_ranges" {
			e.Parameters.ShiftedIonsExcludeRanges = i.Value
		} else if i.Name == "fragment_ion_series" {
			e.Parameters.FragmentIonSeries = i.Value
		} else if i.Name == "search_enzyme_name" {
			e.Parameters.SearchEnzymeName = i.Value
		} else if i.Name == "search_enzyme_cutafter" {
			e.Parameters.SearchEnzymeCutafter = i.Value
		} else if i.Name == "search_enzyme_butnotafter" {
			e.Parameters.SearchEnzymeButnotafter = i.Value
		} else if i.Name == "num_enzyme_termini" {
			e.Parameters.NumEnzymeTermini = i.Value
		} else if i.Name == "allowed_missed_cleavage" {
			e.Parameters.AllowedMissedCleavage = i.Value
		} else if i.Name == "clip_nTerm_M" {
			e.Parameters.ClipNTermM = i.Value
		} else if i.Name == "allow_multiple_variable_mods_on_residue" {
			e.Parameters.AllowMultipleVariableModsOnResidue = i.Value
		} else if i.Name == "max_variable_mods_per_mod" {
			e.Parameters.MaxVariableModsPerMod = i.Value
		} else if i.Name == "max_variable_mods_combinations" {
			e.Parameters.MaxVariableModsCombinations = i.Value
		} else if i.Name == "output_file_extension" {
			e.Parameters.OutputFileExtension = i.Value
		} else if i.Name == "output_format" {
			e.Parameters.OutputFormat = i.Value
		} else if i.Name == "output_report_topN" {
			e.Parameters.OutputReportTopN = i.Value
		} else if i.Name == "output_max_expect" {
			e.Parameters.OutputMaxExpect = i.Value
		} else if i.Name == "report_alternative_proteins" {
			e.Parameters.ReportAlternativeProteins = i.Value
		} else if i.Name == "override_charge" {
			e.Parameters.OverrideCharge = i.Value
		} else if i.Name == "precursor_charge" {
			e.Parameters.PrecursorCharge = i.Value
		} else if i.Name == "digest_min_length" {
			e.Parameters.DigestMinLength = i.Value
		} else if i.Name == "digest_max_length" {
			e.Parameters.DigestMaxLength = i.Value
		} else if i.Name == "digest_mass_range" {
			e.Parameters.DigestMassRange = i.Value
		} else if i.Name == "max_fragment_charge" {
			e.Parameters.MaxFragmentCharge = i.Value
		} else if i.Name == "track_zero_topN" {
			e.Parameters.TrackZeroTopN = i.Value
		} else if i.Name == "zero_bin_accept_expect" {
			e.Parameters.ZeroBinAcceptExpect = i.Value
		} else if i.Name == "zero_bin_mult_expect" {
			e.Parameters.ZeroBinMultExpect = i.Value
		} else if i.Name == "add_topN_complementary" {
			e.Parameters.AddTopNComplementary = i.Value
		} else if i.Name == "minimum_peaks" {
			e.Parameters.MinimumPeaks = i.Value
		} else if i.Name == "use_topN_peaks" {
			e.Parameters.UseTopNPeaks = i.Value
		} else if i.Name == "min_fragments_modelling" {
			e.Parameters.MinFragmentsModelling = i.Value
		} else if i.Name == "min_matched_fragments" {
			e.Parameters.MinMatchedFragments = i.Value
		} else if i.Name == "minimum_ratio" {
			e.Parameters.MinimumRatio = i.Value
		} else if i.Name == "clear_mz_range" {
			e.Parameters.ClearMzRange = i.Value
		} else if i.Name == "variable_mod_01" {
			e.Parameters.VariableMod01 = i.Value
		} else if i.Name == "variable_mod_02" {
			e.Parameters.VariableMod02 = i.Value
		} else if i.Name == "add_A_alanine" {
			e.Parameters.Alanine = i.Value
		} else if i.Name == "add_C_cysteine" {
			e.Parameters.Cysteine = i.Value
		} else if i.Name == "add_Cterm_peptide" {
			e.Parameters.CTermPeptide = i.Value
		} else if i.Name == "add_Cterm_protein" {
			e.Parameters.CTermProtein = i.Value
		} else if i.Name == "add_D_aspartic_acid" {
			e.Parameters.AsparticAcid = i.Value
		} else if i.Name == "add_E_glutamic_acid" {
			e.Parameters.GlutamicAcid = i.Value
		} else if i.Name == "add_F_phenylalanine" {
			e.Parameters.Phenylalanine = i.Value
		} else if i.Name == "add_G_glycine" {
			e.Parameters.Glycine = i.Value
		} else if i.Name == "add_H_histidine" {
			e.Parameters.Histidine = i.Value
		} else if i.Name == "add_I_isoleucine" {
			e.Parameters.Isoleucine = i.Value
		} else if i.Name == "add_K_lysine" {
			e.Parameters.Lysine = i.Value
		} else if i.Name == "add_L_leucine" {
			e.Parameters.Leucine = i.Value
		} else if i.Name == "add_M_methionine" {
			e.Parameters.Methionine = i.Value
		} else if i.Name == "add_N_asparagine" {
			e.Parameters.Asparagine = i.Value
		} else if i.Name == "add_Nterm_peptide" {
			e.Parameters.NTermPeptide = i.Value
		} else if i.Name == "add_Nterm_protein" {
			e.Parameters.NTermProtein = i.Value
		} else if i.Name == "add_P_proline" {
			e.Parameters.Proline = i.Value
		} else if i.Name == "add_Q_glutamine" {
			e.Parameters.GlutamicAcid = i.Value
		} else if i.Name == "add_R_arginine" {
			e.Parameters.Arginine = i.Value
		} else if i.Name == "add_S_serine" {
			e.Parameters.Serine = i.Value
		} else if i.Name == "add_T_threonine" {
			e.Parameters.Threonine = i.Value
		} else if i.Name == "add_V_valine" {
			e.Parameters.Valine = i.Value
		} else if i.Name == "add_W_tryptophan" {
			e.Parameters.Tryptophan = i.Value
		} else if i.Name == "add_Y_tyrosine" {
			e.Parameters.Tyrosine = i.Value
		}
	}

	return
}
