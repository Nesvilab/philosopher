package cmd

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/ext/ptmprophet"
	"github.com/prvst/philosopher/lib/sys"
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

		logrus.Info("Executing PTMProphet")
		var ptm = ptmprophet.New()

		// deploy the binaries
		e := ptm.Deploy(m.OS, m.Distro)
		if e != nil {
			fmt.Println(e.Message)
		}

		// run
		xml, e := ptm.Run(m.PTMProphet, args)
		if e != nil {
			fmt.Println(e.Message)
		}

		_ = xml
		//evi.IndexIdentification(xml, m.PeptideProphet.Decoy)

		m.PTMProphet.InputFiles = args

		m.Serialize()
		logrus.Info("Done")

	},
}

func init() {

	if os.Args[1] == "ptmprophet" {

		m.Restore(sys.Meta())

		ptmprophetCmd.Flags().StringVarP(&m.PTMProphet.Output, "output", "", "", "output prefix file name")
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
