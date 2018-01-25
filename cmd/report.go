package cmd

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/rep"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/spf13/cobra"
)

// reportCmd represents the report commands
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Multi-level reporting for both narrow-searches and open-searches",
	Run: func(cmd *cobra.Command, args []string) {

		if len(m.UUID) < 1 && len(m.Home) < 1 {
			e := &err.Error{Type: err.WorkspaceNotFound, Class: err.FATA}
			logrus.Fatal(e.Error())
		}

		logrus.Info("Executing report")

		rep.Run(m)

		// store parameters on meta data
		m.Serialize()

		logrus.Info("Done")
		return
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "report" {

		m.Restore(sys.Meta())

		reportCmd.Flags().BoolVarP(&m.Report.Decoys, "decoys", "", false, "add decoy observations to reports")
	}

	RootCmd.AddCommand(reportCmd)
}
