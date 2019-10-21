// Package cmd Custer top level command
package cmd

import (
	"os"

	"github.com/nesvilab/philosopher/lib/clu"
	"github.com/nesvilab/philosopher/lib/met"
	"github.com/nesvilab/philosopher/lib/msg"
	"github.com/nesvilab/philosopher/lib/sys"
	"github.com/spf13/cobra"
)

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Protein report based on protein clusters",
	Run: func(cmd *cobra.Command, args []string) {

		m.FunctionInitCheckUp()

		msg.Executing("Cluster ", Version)

		// run clustering
		clu.GenerateReport(m)

		// store paramters on meta data
		m.Serialize()

		// clean tmp
		met.CleanTemp(m.Temp)

		msg.Done()
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
