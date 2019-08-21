// Package cmd Database top level command
package cmd

import (
	"os"

	"github.com/prvst/philosopher/lib/dat"
	"github.com/prvst/philosopher/lib/msg"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/spf13/cobra"
)

var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Target-Decoy database formatting",
	Run: func(cmd *cobra.Command, args []string) {

		m.FunctionInitCheckUp()

		msg.Executing("Database ", Version)

		m = dat.Run(m)

		// store paramters on meta data
		m.Serialize()

		msg.Done()
		return
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "database" {

		m.Restore(sys.Meta())

		databaseCmd.Flags().StringVarP(&m.Database.ID, "id", "", "", "UniProt proteome ID")
		databaseCmd.Flags().StringVarP(&m.Database.Annot, "annotate", "", "", "process a ready-to-use database")
		databaseCmd.Flags().StringVarP(&m.Database.Enz, "enzyme", "", "trypsin", "enzyme for digestion (trypsin, lys_c, lys_n, glu_c, chymotrypsin)")
		databaseCmd.Flags().StringVarP(&m.Database.Tag, "prefix", "", "rev_", "define a decoy prefix")
		databaseCmd.Flags().StringVarP(&m.Database.Add, "add", "", "", "add custom sequences (UniProt FASTA format only)")
		databaseCmd.Flags().StringVarP(&m.Database.Custom, "custom", "", "", "use a pre formatted custom database")
		databaseCmd.Flags().BoolVarP(&m.Database.Crap, "contam", "", false, "add common contaminants")
		databaseCmd.Flags().BoolVarP(&m.Database.Rev, "reviewed", "", false, "use only reviwed sequences from Swiss-Prot")
		databaseCmd.Flags().BoolVarP(&m.Database.Iso, "isoform", "", false, "add isoform sequences")
		databaseCmd.Flags().BoolVarP(&m.Database.NoD, "nodecoys", "", false, "don't add decoys to the database")
	}

	RootCmd.AddCommand(databaseCmd)
}
