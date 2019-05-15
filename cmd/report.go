package cmd

import (
	"os"

	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/rep"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// reportCmd represents the report commands
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Multi-level reporting for both narrow-searches and open-searches",
	Run: func(cmd *cobra.Command, args []string) {

		e := m.FunctionInitCheckUp()
		if e != nil {
			logrus.Fatal(e)
		}

		logrus.Info("Executing Report ", Version)

		rep.Run(m)

		// store parameters on meta data
		m.Serialize()

		// clean tmp
		met.CleanTemp(m.Temp)

		logrus.Info("Done")
		return
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "report" {

		m.Restore(sys.Meta())

		reportCmd.Flags().BoolVarP(&m.Report.Decoys, "decoys", "", false, "add decoy observations to reports")
		reportCmd.Flags().BoolVarP(&m.Report.MSstats, "msstats", "", false, "create an output compatible to MSstats")
		reportCmd.Flags().BoolVarP(&m.Report.MZID, "mzID", "", false, "create a mzID output")
	}

	RootCmd.AddCommand(reportCmd)
}
