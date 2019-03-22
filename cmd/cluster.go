package cmd

import (
	"os"

	"github.com/prvst/philosopher/lib/clu"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Protein report based on protein clusters",
	Run: func(cmd *cobra.Command, args []string) {

		m.FunctionInitCheckUp()

		logrus.Info("Executing Cluster ", Version)
		// run clustering
		e := clu.GenerateReport(m)
		if e != nil {
			logrus.Fatal(e.Error())
		}

		// store paramters on meta data
		m.Serialize()

		// clean tmp
		met.CleanTemp(m.Temp)

		logrus.Info("Done")
		return
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "cluster" {

		m.Restore(sys.Meta())

		clusterCmd.Flags().StringVarP(&m.Cluster.UID, "id", "", "", "UniProt proteome ID")
		clusterCmd.Flags().Float64VarP(&m.Cluster.Level, "level", "", 0.9, "cluster identity level")
	}

	RootCmd.AddCommand(clusterCmd)
}
