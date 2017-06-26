package cmd

import (
	"github.com/Sirupsen/logrus"
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
			logrus.Fatal("Workspace not found. Run 'philosopher init' to create a workspace")
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
	RootCmd.AddCommand(ptmprophetCmd)

}
