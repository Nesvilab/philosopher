// Package cmd Report top level command
package cmd

import (
	"os"

	"philosopher/lib/met"
	"philosopher/lib/msg"
	"philosopher/lib/rep"
	"philosopher/lib/sys"

	"github.com/spf13/cobra"
)

// reportCmd represents the report commands
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Multi-level reporting for both narrow-searches and open-searches",
	Run: func(cmd *cobra.Command, args []string) {

		m.FunctionInitCheckUp()

		msg.Executing("Report ", Version)

		rep.Run(m)

		// store parameters on meta data
		m.Serialize()

		// clean tmp
		met.CleanTemp(m.Temp)

		msg.Done()
		return
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "report" {

		m.Restore(sys.Meta())

		reportCmd.Flags().BoolVarP(&m.Report.Decoys, "decoys", "", false, "add decoy observations to reports")
		reportCmd.Flags().BoolVarP(&m.Report.MSstats, "msstats", "", false, "create an output compatible with MSstats")
		reportCmd.Flags().BoolVarP(&m.Report.MZID, "mzid", "", false, "create a mzID output")
	}

	RootCmd.AddCommand(reportCmd)
}
