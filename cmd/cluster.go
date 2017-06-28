package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/prvst/cmsl/err"
	"github.com/prvst/philosopher/lib/clus"
	"github.com/prvst/philosopher/lib/meta"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/spf13/cobra"
)

var cls clus.Cluster

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Protein report based on protein clusters",
	//Long:  `Proteins are clustered based on sequence identity levels, and peptides are mapped to clusters, providing MS/MS evidence on a functional level.`,
	Run: func(cmd *cobra.Command, args []string) {

		var m meta.Data
		m.Restore(sys.Meta())
		if len(m.UUID) < 1 && len(m.Home) < 1 {
			e := &err.Error{Type: err.WorkspaceNotFound, Class: err.FATA}
			logrus.Fatal(e.Error())
		}

		cls.GenerateReport()

		logrus.Info("Done")

		return
	},
}

func init() {

	cls = clus.New()

	clusterCmd.Flags().StringVarP(&cls.UID, "id", "", "", "UniProt proteome ID")
	clusterCmd.Flags().Float64VarP(&cls.Level, "level", "", 0.9, "cluster identity level")

	RootCmd.AddCommand(clusterCmd)
}
