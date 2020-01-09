// Package cmd MSFragger top level command
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"philosopher/lib/ext/msfragger"
	"philosopher/lib/met"
	"philosopher/lib/msg"
	"philosopher/lib/sys"
)

// msfraggerCmd represents the msfragger command
var msfraggerCmd = &cobra.Command{
	Use:   "msfragger",
	Short: "Ultra fast and comprehensive peptide identification in mass spectrometry–based proteomics",
	Run: func(cmd *cobra.Command, args []string) {

		m.FunctionInitCheckUp()

		msg.Executing("MSFragger ", Version)

		m := msfragger.Run(m, args)

		m.Serialize()

		// clean tmp
		met.CleanTemp(m.Temp)

		msg.Done()
		return

	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "msfragger" {

		m.Restore(sys.Meta())

		msfraggerCmd.Flags().StringVarP(&m.MSFragger.JarPath, "path", "", "", "")
		msfraggerCmd.Flags().StringVarP(&m.MSFragger.Param, "param", "", "", "")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.Memory, "memory", "", 8, "")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.Threads, "num_threads", "", 0, "CPU to set num threads; else specify num threads directly (max 64)")
		//msfraggerCmd.Flags().StringVarP(&m.MSFragger.RawExtension, "raw", "", "", "")
		msfraggerCmd.Flags().StringVarP(&m.MSFragger.DatabaseName, "database_name", "", "", "")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.PrecursorMassLower, "precursor_mass_lower", "", -20, "")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.PrecursorMassUpper, "precursor_mass_upper", "", 20, "")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.PrecursorMassUnits, "precursor_mass_units", "", 1, "0=Daltons, 1=ppm")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.PrecursorTrueTolerance, "precursor_true_tolerance", "", 20, "")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.PrecursorTrueUnits, "precursor_true_units", "", 1, "0=Daltons, 1=ppm")
		msfraggerCmd.Flags().Float64VarP(&m.MSFragger.FragmentMassTolerance, "fragment_mass_tolerance", "", 20, "")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.FragmentMassUnits, "fragment_mass_units", "", 1, "")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.CalibrateMass, "calibrate_mass", "", 0, "0=Off, 1=On, 2=On and find optimal parameters")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.Deisotope, "deisotope", "", 0, "0=off, -1/0/1/2/3 (standard C13 error)")
		msfraggerCmd.Flags().StringVarP(&m.MSFragger.IsotopeError, "isotope_error", "-1/0/1/2/3", "0", "0=off, -1/0/1/2/3 (standard C13 error)")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.MassOffsets, "mass_offsets", "", 0, "allow for additional precursor mass window shifts")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.LocalizeDeltaMass, "localize_delta_mass", "", 0, "")
		msfraggerCmd.Flags().StringVarP(&m.MSFragger.PrecursorMassMode, "precursor_mass_mode", "", "selected", "")
		//msfraggerCmd.Flags().StringVarP(&m.MSFragger.DeltaMassExcludeRanges, "delta_mass_exclude_ranges", "", "(-1.5,3.5)", "")
		msfraggerCmd.Flags().StringVarP(&m.MSFragger.FragmentIonSeries, "fragment_ion_series", "", "b,y", "")
		msfraggerCmd.Flags().StringVarP(&m.MSFragger.SearchEnzymeName, "search_enzyme_name", "", "Trypsin", "")
		msfraggerCmd.Flags().StringVarP(&m.MSFragger.SearchEnzymeCutafter, "search_enzyme_cutafter", "", "KR", "")
		msfraggerCmd.Flags().StringVarP(&m.MSFragger.SearchEnzymeButNotAfter, "search_enzyme_butnotafter", "", "P", "")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.NumEnzymeTermini, "num_enzyme_termini", "", 2, "2 for enzymatic, 1 for semi-enzymatic, 0 for nonspecific digestion")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.AllowedMissedCleavage, "allowed_missed_cleavage", "", 1, "maximum value is 5")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.ClipNTermM, "clip_nTerm_M", "", 1, "")
		msfraggerCmd.Flags().StringVarP(&m.MSFragger.VariableMod01, "variable_mod_01", "", "15.99490 M", "")
		msfraggerCmd.Flags().StringVarP(&m.MSFragger.VariableMod02, "variable_mod_02", "", "42.01060 [^", "")
		msfraggerCmd.Flags().StringVarP(&m.MSFragger.VariableMod03, "variable_mod_03", "", "", "")
		msfraggerCmd.Flags().StringVarP(&m.MSFragger.VariableMod04, "variable_mod_04", "", "", "")
		msfraggerCmd.Flags().StringVarP(&m.MSFragger.VariableMod05, "variable_mod_05", "", "", "")
		msfraggerCmd.Flags().StringVarP(&m.MSFragger.VariableMod06, "variable_mod_06", "", "", "")
		msfraggerCmd.Flags().StringVarP(&m.MSFragger.VariableMod07, "variable_mod_07", "", "", "")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.AllowMultipleVariableModsOnResidue, "allow_multiple_variable_mods_on_residue", "", 1, "static mods are not considered")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.MaxVariableModsPerPeptide, "max_variable_mods_per_peptide", "", 3, "maximun of 5")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.MaxVariableModsCombinations, "max_variable_mods_combinations", "", 5000, "maximum of 65534, limits number of modified peptides generated from sequence")
		msfraggerCmd.Flags().StringVarP(&m.MSFragger.OutputFileExtension, "output_file_extension", "", "pepXML", "")
		msfraggerCmd.Flags().StringVarP(&m.MSFragger.OutputFormat, "output_format", "", "pepXML", "")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.OutputReportTopN, "output_report_topN", "", 1, "")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.OutputMaxExpect, "output_max_expect", "", 50, "")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.ReportAlternativeProteins, "report_alternative_proteins", "", 0, "0=no, 1=yes")
		msfraggerCmd.Flags().StringVarP(&m.MSFragger.PrecursorCharge, "precursor_charge", "", "1 4", "precursor charge range to analyze; does not override any existing charge; 0 as 1st entry ignores parameter")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.OverrideCharge, "override_charge", "", 0, "")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.DigestMinLength, "digest_min_length", "", 7, "")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.DigestMaxLength, "digest_max_length", "", 50, "")
		msfraggerCmd.Flags().StringVarP(&m.MSFragger.DigestMassRange, "digest_mass_range", "", "500.0 5000.0", "")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.MaxFragmentCharge, "max_fragment_charge", "", 2, "")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.TrackZeroTopN, "track_zero_topN", "", 0, "in addition to topN results, keep track of top results in zero bin")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.ZeroBinAcceptExpect, "zero_bin_accept_expect", "", 0, "boost top zero bin entry to top if it has expect under 0.01 - set to 0 to disable")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.ZeroBinMultExpect, "zero_bin_mult_expect", "", 1, "disabled if above passes - multiply expect of zero bin for ordering purposes (does not affect reported expect)")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.AddTopNComplementary, "add_topN_complementary", "", 0, "")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.MinimumPeaks, "minimum_peaks", "", 15, "required minimum number of peaks in spectrum to search (default 10)")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.UseTopNPeaks, "use_topN_peaks", "", 100, "")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.MinFragmentsModelling, "min_fragments_modelling", "", 2, "")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.MinMatchedFragments, "min_matched_fragments", "", 4, "")
		msfraggerCmd.Flags().Float64VarP(&m.MSFragger.MinimumRatio, "minimum_ratio", "", 0.01, "")
		msfraggerCmd.Flags().StringVarP(&m.MSFragger.ClearMzRange, "clear_mz_range", "", "0.0 0.0", "")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.RemovePrecursorPeak, "remove_precursor_peak", "", 0, "remove precursor peaks from tandem mass spectra. 0=not remove; 1=remove the peak with precursor charge;")
		msfraggerCmd.Flags().StringVarP(&m.MSFragger.RemovePrecursorRange, "remove_precursor_range", "", "-1.5,1.5", "m/z range in removing precursor peaks. Unit: Da.")
		msfraggerCmd.Flags().IntVarP(&m.MSFragger.IntensityTransform, "intensity_transform", "", 0, "transform peaks intensities with sqrt root. 0 = not transform; 1 = transform using sqrt root.")
		msfraggerCmd.Flags().Float64VarP(&m.MSFragger.AddCtermPeptide, "add_Cterm_peptide", "", 0.000000, "")
		msfraggerCmd.Flags().Float64VarP(&m.MSFragger.AddCtermProtein, "add_Cterm_protein", "", 0.000000, "")
		msfraggerCmd.Flags().Float64VarP(&m.MSFragger.AddNTermPeptide, "add_Nterm_peptide", "", 0.000000, "")
		msfraggerCmd.Flags().Float64VarP(&m.MSFragger.AddNtermProteine, "add_Nterm_protein", "", 0.000000, "")
		msfraggerCmd.Flags().Float64VarP(&m.MSFragger.AddAlanine, "add_A_alanine", "", 0.000000, "")
		msfraggerCmd.Flags().Float64VarP(&m.MSFragger.AddCysteine, "add_C_cysteine", "", 0.000000, "")
		msfraggerCmd.Flags().Float64VarP(&m.MSFragger.AddAsparticAcid, "add_D_aspartic_acid", "", 0.000000, "")
		msfraggerCmd.Flags().Float64VarP(&m.MSFragger.AddGlutamicAcid, "add_E_glutamic_acid", "", 0.000000, "")
		msfraggerCmd.Flags().Float64VarP(&m.MSFragger.AddPhenylAlnine, "add_F_phenylalanine", "", 0.000000, "")
		msfraggerCmd.Flags().Float64VarP(&m.MSFragger.AddGlycine, "add_G_glycine", "", 0.000000, "")
		msfraggerCmd.Flags().Float64VarP(&m.MSFragger.AddHistidine, "add_H_histidine", "", 0.000000, "")
		msfraggerCmd.Flags().Float64VarP(&m.MSFragger.AddIsoleucine, "add_I_isoleucine", "", 0.000000, "")
		msfraggerCmd.Flags().Float64VarP(&m.MSFragger.AddLysine, "add_K_lysine", "", 0.000000, "")
		msfraggerCmd.Flags().Float64VarP(&m.MSFragger.AddLeucine, "add_L_leucine", "", 0.000000, "")
		msfraggerCmd.Flags().Float64VarP(&m.MSFragger.AddMethionine, "add_M_methionine", "", 0.000000, "")
		msfraggerCmd.Flags().Float64VarP(&m.MSFragger.AddAsparagine, "add_N_asparagine", "", 0.000000, "")
		msfraggerCmd.Flags().Float64VarP(&m.MSFragger.AddProline, "add_P_proline", "", 0.000000, "")
		msfraggerCmd.Flags().Float64VarP(&m.MSFragger.AddGlutamine, "add_Q_glutamine", "", 0.000000, "")
		msfraggerCmd.Flags().Float64VarP(&m.MSFragger.AddArginine, "add_R_arginine", "", 0.000000, "")
		msfraggerCmd.Flags().Float64VarP(&m.MSFragger.AddSerine, "add_S_serine", "", 0.000000, "")
		msfraggerCmd.Flags().Float64VarP(&m.MSFragger.AddThreonine, "add_T_threonine", "", 0.000000, "")
		msfraggerCmd.Flags().Float64VarP(&m.MSFragger.AddValine, "add_V_valine", "", 0.000000, "")
		msfraggerCmd.Flags().Float64VarP(&m.MSFragger.AddTryptophan, "add_W_tryptophan", "", 0.000000, "")
		msfraggerCmd.Flags().Float64VarP(&m.MSFragger.AddTyrosine, "add_Y_tyrosine", "", 0.000000, "")

	}

	RootCmd.AddCommand(msfraggerCmd)

}
