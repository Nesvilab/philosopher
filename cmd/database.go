package cmd

import (
	"os"

	"github.com/prvst/philosopher/lib/dat"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Target-Decoy database formatting",
	//Long: `The database command alows the creation and formatting of a Target-Decoy database. It also
	//provides options for downloading a fresh snapshot from UniProt`,
	Run: func(cmd *cobra.Command, args []string) {

		// if len(m.UUID) < 1 && len(m.Home) < 1 {
		// 	e := &err.Error{Type: err.WorkspaceNotFound, Class: err.FATA}
		// 	logrus.Fatal(e.Error())
		// }

		m.FunctionInitCheckUp()

		logrus.Info("Executing Database ", Version)

		m, e := dat.Run(m)
		if e != nil {
			logrus.Fatal(e)
		}

		// store paramters on meta data
		m.Serialize()

		logrus.Info("Done")
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
	}

	RootCmd.AddCommand(databaseCmd)
}
