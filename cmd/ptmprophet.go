package cmd

import (
	"os"

	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/ext/ptmprophet"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/spf13/cobra"
)

// ptmprophetCmd represents the ptmprophet command
var ptmprophetCmd = &cobra.Command{
	Use:   "ptmprophet",
	Short: "PTM site localization",
	Run: func(cmd *cobra.Command, args []string) {

		m.FunctionInitCheckUp()

		err.Executing("PTMProphet ", Version)

		m.PTMProphet.InputFiles = args
		ptmprophet.Run(m, args)
		m.Serialize()

		// clean tmp
		met.CleanTemp(m.Temp)

		err.Done()
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "ptmprophet" {

		m.Restore(sys.Meta())

		ptmprophetCmd.Flags().StringVarP(&m.PTMProphet.Output, "output", "", "", "output prefix file name")
		ptmprophetCmd.Flags().StringVarP(&m.PTMProphet.Mods, "mods", "", "", "<amino acids, n, or c>:<mass_shift>:<neut_loss1>:...:<neut_lossN>,<amino acids, n, or c>:<mass_shift>:<neut_loss1>:...:<neut_lossN> (overrides the modifications from the interact.pep.xml file)")
		ptmprophetCmd.Flags().StringVarP(&m.PTMProphet.NIons, "nions", "", "", "use specified N-term ions, separate multiple ions by commas (default: a,b for CID, c for ETD)")
		ptmprophetCmd.Flags().StringVarP(&m.PTMProphet.CIons, "cions", "", "", "use specified C-term ions, separate multiple ions by commas (default: y for CID, z for ETD)")
		ptmprophetCmd.Flags().IntVarP(&m.PTMProphet.EM, "em", "", 2, "set EM models to 0 (no EM), 1 (Intensity EM Model Applied) or 2 (Intensity and Matched Peaks EM Models Applied)")
		ptmprophetCmd.Flags().IntVarP(&m.PTMProphet.FragPPMTol, "fragppmtol", "", 15, "when computing PSM-specific mass_offset and mass_tolerance, use specified default +/- MS2 mz tolerance on fragment ions")
		ptmprophetCmd.Flags().IntVarP(&m.PTMProphet.MaxThreads, "maxthreads", "", 1, "use specified number of threads for processing")
		ptmprophetCmd.Flags().IntVarP(&m.PTMProphet.MaxFragZ, "maxfragz", "", 0, "limit maximum fragment charge (default: 0=precursor charge, negative values subtract from precursor charge)")
		ptmprophetCmd.Flags().IntVarP(&m.PTMProphet.Mino, "mino", "", 0, "use specified number of pseudo-counts when computing Oscore")
		ptmprophetCmd.Flags().IntVarP(&m.PTMProphet.MassOffset, "massoffset", "", 0, "adjust the massdiff by offset <number>")
		ptmprophetCmd.Flags().BoolVarP(&m.PTMProphet.Static, "static", "", false, "use static fragppmtol for all PSMs instead of dynamically estimates offsets and tolerances")
		ptmprophetCmd.Flags().BoolVarP(&m.PTMProphet.KeepOld, "keepold", "", false, "retain old PTMProphet results in the pepXML file")
		ptmprophetCmd.Flags().BoolVarP(&m.PTMProphet.Verbose, "verbose", "", false, "produce Warnings to help troubleshoot potential PTM shuffling or mass difference issues")
		ptmprophetCmd.Flags().BoolVarP(&m.PTMProphet.Lability, "lability", "", false, "compute Lability of PTMs")
		ptmprophetCmd.Flags().BoolVarP(&m.PTMProphet.Direct, "direct", "", false, "use only direct evidence for evaluating PTM site probabilities")
		ptmprophetCmd.Flags().BoolVarP(&m.PTMProphet.Ifrags, "ifrags", "", false, "use internal fragments for localization")
		ptmprophetCmd.Flags().BoolVarP(&m.PTMProphet.Autodirect, "autodirect", "", false, "use direct evidence when the lability is high, use in combination with LABILITY")
		ptmprophetCmd.Flags().BoolVarP(&m.PTMProphet.NoMinoFactor, "nominofactor", "", false, "disable MINO factor correction when MINO= is set greater than 0 (default: apply MINO factor correction)")
		ptmprophetCmd.Flags().Float64VarP(&m.PTMProphet.PPMTol, "ppmtol", "", 1, "use specified +/- MS1 ppm tolerance on peptides which may have a slight offset depending on search parameters")
		ptmprophetCmd.Flags().Float64VarP(&m.PTMProphet.MinProb, "minprob", "", 0.9, "use specified minimum probability to evaluate peptides")
		ptmprophetCmd.Flags().BoolVarP(&m.PTMProphet.MassDiffMode, "massdiffmode", "", false, "use the Mass Difference and localize")
	}

	RootCmd.AddCommand(ptmprophetCmd)

}
