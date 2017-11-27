package cmd

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher/lib/clu"
	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/spf13/cobra"
)

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Protein report based on protein clusters",
	Run: func(cmd *cobra.Command, args []string) {

		if len(m.UUID) < 1 && len(m.Home) < 1 {
			e := &err.Error{Type: err.WorkspaceNotFound, Class: err.FATA}
			logrus.Fatal(e.Error())
		}

		// run clustering
		clu.GenerateReport(m.Cluster, m.Home, m.Temp)

		// store paramters on meta data
		m.Serialize()

		logrus.Info("Done")
		return
	},
}

func init() {

	if os.Args[1] == "cluster" {

		m.Restore(sys.Meta())

		clusterCmd.Flags().StringVarP(&m.Cluster.UID, "id", "", "", "UniProt proteome ID")
		clusterCmd.Flags().Float64VarP(&m.Cluster.Level, "level", "", 0.9, "cluster identity level")
	}

	RootCmd.AddCommand(clusterCmd)
}
