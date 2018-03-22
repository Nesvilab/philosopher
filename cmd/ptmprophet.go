package cmd

import (
	"os"

	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/ext/ptmprophet"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ptmprophetCmd represents the ptmprophet command
var ptmprophetCmd = &cobra.Command{
	Use:   "ptmprophet",
	Short: "PTM site localization",
	Run: func(cmd *cobra.Command, args []string) {

		if len(m.UUID) < 1 && len(m.Home) < 1 {
			e := &err.Error{Type: err.WorkspaceNotFound, Class: err.FATA}
			logrus.Fatal(e.Error())
		}

		logrus.Info("Executing PTMProphet ", Version)

		m.PTMProphet.InputFiles = args

		ptmprophet.Run(m, args)
		m.Serialize()

		logrus.Info("Done")
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "ptmprophet" {

		m.Restore(sys.Meta())

		ptmprophetCmd.Flags().StringVarP(&m.PTMProphet.Output, "output", "", "", "output prefix file name")
		ptmprophetCmd.Flags().StringVarP(&m.PTMProphet.Mods, "mods", "", "", "specify modifications. <amino acids, n, or c>:<mass_shift>:<neut_loss1>:...:<neut_lossN>,<amino acids, n, or c>:<mass_shift>:<neut_loss1>:...:<neut_lossN> (overrides the modifications from the interact.pep.xml file)")
		ptmprophetCmd.Flags().IntVarP(&m.PTMProphet.EM, "em", "", 1, "Set EM models to 0 (no EM), 1 (Intensity EM Model Applied) or 2 (Intensity and Matched Peaks EM Models Applied)")
		ptmprophetCmd.Flags().BoolVarP(&m.PTMProphet.KeepOld, "keepold", "", false, "retain old PTMProphet results in the pepXML file")
		ptmprophetCmd.Flags().BoolVarP(&m.PTMProphet.Verbose, "verbose", "", false, "produce Warnings to help troubleshoot potential PTM shuffling or mass difference issues")
		ptmprophetCmd.Flags().Float64VarP(&m.PTMProphet.MzTol, "mztol", "", 0.1, "use specified +/- MS2 mz tolerance on site specific ions")
		ptmprophetCmd.Flags().Float64VarP(&m.PTMProphet.PPMTol, "ppmtol", "", 1, "use specified +/- MS1 ppm tolerance on peptides which may have a slight offset depending on search parameters")
		ptmprophetCmd.Flags().Float64VarP(&m.PTMProphet.MinProb, "minprob", "", 0, "use specified minimum probability to evaluate peptides")
		ptmprophetCmd.Flags().BoolVarP(&m.PTMProphet.MassDiffMode, "massdiffmode", "", false, "use the Mass Difference and localize")
	}

	RootCmd.AddCommand(ptmprophetCmd)

}
