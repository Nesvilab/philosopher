package msfragger

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/met"
	"github.com/sirupsen/logrus"
)

// MSFragger represents the tool configuration
type MSFragger struct {
	DefaultBin   string
	DefaultParam string
}

// New constructor
func New(temp string) MSFragger {

	var self MSFragger

	self.DefaultBin = ""
	self.DefaultParam = ""

	return self
}

// Run is the Fragger main entry point
func Run(m met.Data, args []string) (met.Data, *err.Error) {

	var frg = New(m.Temp)

	// collect and store the mz files
	m.MSFragger.RawFiles = args

	if len(m.MSFragger.Param) > 1 {
		// convert the param file to binary and store it in meta
		var binFile []byte
		paramAbs, _ := filepath.Abs(m.MSFragger.Param)
		binFile, e := ioutil.ReadFile(paramAbs)
		if e != nil {
			logrus.Fatal(e)
		}
		m.MSFragger.ParamFile = binFile
	}

	// run comet
	e := frg.Execute(m.MSFragger, args)
	if e != nil {
		//logrus.Fatal(e)
	}

	return m, nil
}

// Execute is the main fucntion to execute MSFragger
func (c *MSFragger) Execute(params met.MSFragger, cmdArgs []string) *err.Error {

	cmd := appendParams(params)

	for _, i := range cmdArgs {
		file, _ := filepath.Abs(i)
		cmd.Args = append(cmd.Args, file)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	e := cmd.Start()
	if e != nil {
		return nil
	}

	_ = cmd.Wait()

	return nil
}

func appendParams(params met.MSFragger) *exec.Cmd {

	mem := fmt.Sprintf("-Xmx%dG", params.Memmory)
	dbPath, _ := filepath.Abs(params.DatabaseName)
	jarPath, _ := filepath.Abs(params.JarPath)

	if len(params.VariableMod01) < 1 {
		params.VariableMod01 = ""
	}

	if len(params.VariableMod02) < 1 {
		params.VariableMod02 = ""
	}

	if len(params.VariableMod03) < 1 {
		params.VariableMod03 = ""
	}

	if len(params.VariableMod04) < 1 {
		params.VariableMod04 = ""
	}

	if len(params.VariableMod05) < 1 {
		params.VariableMod05 = ""
	}

	if len(params.VariableMod06) < 1 {
		params.VariableMod06 = ""
	}

	if len(params.VariableMod07) < 1 {
		params.VariableMod07 = ""
	}

	args := exec.Command("java",
		"-jar",
		mem,
		jarPath,
		"--database_name",
		dbPath,
		"--num_threads",
		fmt.Sprintf("%d", params.Threads),
		"--precursor_mass_lower",
		fmt.Sprintf("%d", params.PrecursorMassLower),
		"--precursor_mass_upper",
		fmt.Sprintf("%d", params.PrecursorMassUpper),
		"--precursor_mass_units",
		fmt.Sprintf("%d", params.PrecursorMassUnits),
		"--precursor_true_tolerance",
		fmt.Sprintf("%d", params.PrecursorTrueTolerance),
		"--fragment_mass_tolerance",
		fmt.Sprintf("%f", params.FragmentMassTolerance),
		"--fragment_mass_units",
		fmt.Sprintf("%d", params.FragmentMassUnits),
		"--calibrate_mass",
		fmt.Sprintf("%d", params.CalibrateMass),
		"--isotope_error",
		fmt.Sprintf("%s", params.IsotopeError),
		"--mass_offsets",
		fmt.Sprintf("%d", params.MassOffsets),
		"--localize_delta_mass",
		fmt.Sprintf("%d", params.LocalizeDeltaMass),
		"--precursor_mass_mode",
		fmt.Sprintf("%s", params.PrecursorMassMode),
		//"--shifted_ions_exclude_ranges",
		//fmt.Sprintf("%s", params.ShiftedIonsExcludeRanges),
		"--fragment_ion_series",
		fmt.Sprintf("%s", params.FragmentIonSeries),
		"--search_enzyme_name",
		fmt.Sprintf("%s", params.SearchEnzymeName),
		"--search_enzyme_cutafter",
		fmt.Sprintf("%s", params.SearchEnzymeCutafter),
		"--search_enzyme_butnotafter",
		fmt.Sprintf("%s", params.SearchEnzymeButNotAfter),
		"--num_enzyme_termini",
		fmt.Sprintf("%d", params.NumEnzymeTermini),
		"--allowed_missed_cleavage",
		fmt.Sprintf("%d", params.AllowedMissedCleavage),
		"--clip_nTerm_M",
		fmt.Sprintf("%d", params.ClipNTermM),
		"--variable_mod_01",
		fmt.Sprintf("%s", params.VariableMod01),
		"--variable_mod_02",
		fmt.Sprintf("%s", params.VariableMod02),
		"--variable_mod_03",
		fmt.Sprintf("%s", params.VariableMod03),
		"--variable_mod_04",
		fmt.Sprintf("%s", params.VariableMod04),
		"--allow_multiple_variable_mods_on_residue",
		fmt.Sprintf("%d", params.AllowMultipleVariableModsOnResidue),
		"--max_variable_mods_per_mod",
		fmt.Sprintf("%d", params.MaxVariableModsPerMod),
		"--max_variable_mods_combinations",
		fmt.Sprintf("%d", params.MaxVariableModsCombinations),
		"--output_file_extension",
		fmt.Sprintf("%s", params.OutputFileExtension),
		"--output_format",
		fmt.Sprintf("%s", params.OutputFormat),
		"--output_report_topN",
		fmt.Sprintf("%d", params.OutputReportTopN),
		"--output_max_expect",
		fmt.Sprintf("%d", params.OutputMaxExpect),
		"--report_alternative_proteins",
		fmt.Sprintf("%d", params.ReportAlternativeProteins),
		"--precursor_charge",
		fmt.Sprintf("%s", params.PrecursorCharge),
		"--override_charge",
		fmt.Sprintf("%d", params.OverrideCharge),
		"--digest_min_length",
		fmt.Sprintf("%d", params.DigestMinLength),
		"--digest_max_length",
		fmt.Sprintf("%d", params.DigestMaxLength),
		"--digest_mass_range",
		fmt.Sprintf("%s", params.DigestMassRange),
		"--max_fragment_charge",
		fmt.Sprintf("%d", params.MaxFragmentCharge),
		"--track_zero_topN",
		fmt.Sprintf("%d", params.TrackZeroTopN),
		"--zero_bin_accept_expect",
		fmt.Sprintf("%d", params.ZeroBinAcceptExpect),
		"--zero_bin_mult_expect",
		fmt.Sprintf("%d", params.ZeroBinMultExpect),
		"--add_topN_complementary",
		fmt.Sprintf("%d", params.AddTopNComplementary),
		"--minimum_peaks",
		fmt.Sprintf("%d", params.MinimumPeaks),
		"--use_topN_peaks",
		fmt.Sprintf("%d", params.UseTopNPeaks),
		"--min_fragments_modelling",
		fmt.Sprintf("%d", params.MinFragmentsModelling),
		"--min_matched_fragments",
		fmt.Sprintf("%d", params.MinMatchedFragments),
		"--minimum_ratio",
		fmt.Sprintf("%f", params.MinimumRatio),
		"--clear_mz_range",
		fmt.Sprintf("%s", params.ClearMzRange),
		"--add_C_cysteine",
		fmt.Sprintf("%f", params.AddCysteine),
		"--add_Cterm_peptide",
		fmt.Sprintf("%f", params.AddCtermPeptide),
		"--add_Cterm_protein",
		fmt.Sprintf("%f", params.AddCtermProtein),
		"--add_D_aspartic_acid",
		fmt.Sprintf("%f", params.AddAsparticAcid),
		"--add_E_glutamic_acid",
		fmt.Sprintf("%f", params.AddGlutamicAcid),
		"--add_F_phenylalanine",
		fmt.Sprintf("%f", params.AddPhenylAlnine),
		"--add_G_glycine",
		fmt.Sprintf("%f", params.AddGlycine),
		"--add_H_histidine",
		fmt.Sprintf("%f", params.AddHistidine),
		"--add_I_isoleucine",
		fmt.Sprintf("%f", params.AddIsoleucine),
		"--add_K_lysine",
		fmt.Sprintf("%f", params.AddLysine),
		"--add_L_leucine",
		fmt.Sprintf("%f", params.AddLeucine),
		"--add_M_methionine",
		fmt.Sprintf("%f", params.AddMethionine),
		"--add_N_asparagine",
		fmt.Sprintf("%f", params.AddAsparagine),
		"--add_Nterm_peptide",
		fmt.Sprintf("%f", params.AddNTermPeptide),
		"--add_Nterm_protein",
		fmt.Sprintf("%f", params.AddNtermProteine),
		"--add_P_proline",
		fmt.Sprintf("%f", params.AddProline),
		"--add_Q_glutamine",
		fmt.Sprintf("%f", params.AddGlutamine),
		"--add_R_arginine",
		fmt.Sprintf("%f", params.AddArginine),
		"--add_S_serine",
		fmt.Sprintf("%f", params.AddSerine),
		"--add_T_threonine",
		fmt.Sprintf("%f", params.AddThreonine),
		"--add_V_valine",
		fmt.Sprintf("%f", params.AddValine),
		"--add_W_tryptophan",
		fmt.Sprintf("%f", params.AddTryptophan),
		"--add_Y_tyrosine",
		fmt.Sprintf("%f", params.AddTyrosine),
	)

	return args
}
