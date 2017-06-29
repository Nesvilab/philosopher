package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/prvst/cmsl/err"
	"github.com/prvst/philosopher/lib/ext/ptmprophet"
	"github.com/prvst/philosopher/lib/meta"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/spf13/cobra"
)

var ptm ptmprophet.PTMProphet

// ptmprophetCmd represents the ptmprophet command
var ptmprophetCmd = &cobra.Command{
	Use:   "ptmprophet",
	Short: "PTM site localisation",
	Run: func(cmd *cobra.Command, args []string) {

		var m meta.Data
		m.Restore(sys.Meta())
		if len(m.UUID) < 1 && len(m.Home) < 1 {
			e := &err.Error{Type: err.WorkspaceNotFound, Class: err.FATA}
			logrus.Fatal(e.Error())
		}

		// deploy the binaries
		err := ptm.Deploy()
		if err != nil {
			logrus.Fatal(err)
		}

		// run
		err = ptm.Run(args)
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Info("Done")

	},
}

func init() {

	ptm = ptmprophet.New()

	ptmprophetCmd.Flags().StringVarP(&ptm.Output, "output", "", "", "output prefix file name")
	//ptmprophetCmd.Flags().BoolVarP(&ptm.NoUpdate, "noupdate", "", false, "don't update modification_info tags in pepXML")
	ptmprophetCmd.Flags().Int8VarP(&ptm.EM, "em", "", 1, "Set EM models to 0 (no EM), 1 (Intensity EM Model Applied) or 2 (Intensity and Matched Peaks EM Models Applied)")
	ptmprophetCmd.Flags().BoolVarP(&ptm.KeepOld, "keepold", "", false, "retain old PTMProphet results in the pepXML file")
	ptmprophetCmd.Flags().BoolVarP(&ptm.Verbose, "verbose", "", false, "produce Warnings to help troubleshoot potential PTM shuffling or mass difference issues")
	ptmprophetCmd.Flags().Float64VarP(&ptm.MzTol, "mztol", "", 0.1, "use specified +/- MS2 mz tolerance on site specific ions")
	ptmprophetCmd.Flags().Float64VarP(&ptm.PPMTol, "ppmtol", "", 1, "use specified +/- MS1 ppm tolerance on peptides which may have a slight offset depending on search parameters")
	ptmprophetCmd.Flags().Float64VarP(&ptm.MinProb, "minprob", "", 0, "use specified minimum probability to evaluate peptides")
	ptmprophetCmd.Flags().BoolVarP(&ptm.MassDiffMode, "massdiffmode", "", false, "use the Mass Difference and localize")

	RootCmd.AddCommand(ptmprophetCmd)

}
