// Package cmd Bioquant top level command
package cmd

import (
	"os"

	"philosopher/lib/met"
	"philosopher/lib/msg"
	"philosopher/lib/qua"
	"philosopher/lib/sys"

	"github.com/spf13/cobra"
)

// bioquantCmd represents the bioquant command
var bioquantCmd = &cobra.Command{
	Use:   "bioquant",
	Short: "Protein report based on protein functional groups",
	Run: func(cmd *cobra.Command, args []string) {

		m.FunctionInitCheckUp()

		msg.Executing("Cluster ", Version)

		// run clustering
		qua.RunBioQuantification(m)

		// store paramters on meta data
		m.Serialize()

		// clean tmp
		met.CleanTemp(m.Temp)

		msg.Done()
		return
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "bioquant" {

		m.Restore(sys.Meta())

		bioquantCmd.Flags().StringVarP(&m.BioQuant.UID, "id", "", "", "UniProt proteome ID")
		bioquantCmd.Flags().Float64VarP(&m.BioQuant.Level, "level", "", 0.9, "cluster identity level")
	}

	RootCmd.AddCommand(bioquantCmd)
}
