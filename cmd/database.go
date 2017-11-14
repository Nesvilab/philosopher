package cmd

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher/lib/data"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/spf13/cobra"
)

var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Target-Decoy database formatting",
	//Long: `The database command alows the creation and formatting of a Target-Decoy database. It also
	//provides options for downloading a fresh snapshot from UniProt`,
	Run: func(cmd *cobra.Command, args []string) {

		// store paramters on meta data
		m.Serialize()

		var db data.Base

		if len(m.Database.Annot) > 0 {

			logrus.Info("Processing database")

			err := db.ProcessDB(m.Database.Annot, m.Database.Tag)
			if err != nil {
				logrus.Fatal(err)
			}

			err = db.Serialize()
			if err != nil {
				logrus.Fatal(err)
			}

			logrus.Info("Done")
			return

		}

		if len(m.Database.ID) < 1 && len(m.Database.Custom) < 1 {
			logrus.Fatal("You need to provide a taxon ID or a custom FASTA file")
		}

		if m.Database.Crap == false {
			logrus.Warning("Contaminants are not going to be added to database")
		}

		if len(m.Database.Custom) < 1 {

			logrus.Info("Fetching database")
			db.Fetch(m.Database.ID, m.Temp, m.Database.Iso, m.Database.Rev)

		} else {
			db.UniProtDB = m.Database.Custom
		}

		logrus.Info("Processing decoys")
		db.Create(m.Temp, m.Database.Add, m.Database.Enz, m.Database.Tag, m.Database.Crap)

		logrus.Info("Creating file")
		db.Save(m.Home, m.Temp, m.Database.Tag)

		err := db.Serialize()
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Info("Done")
		return
	},
}

func init() {

	if os.Args[1] == "database" {

		m.Restore(sys.Meta())

		databaseCmd.Flags().StringVarP(&m.Database.ID, "id", "", "", "UniProt proteome ID")
		databaseCmd.Flags().StringVarP(&m.Database.Annot, "annotate", "", "", "process a ready-to-use database")
		databaseCmd.Flags().StringVarP(&m.Database.Enz, "enzyme", "", "trypsin", "enzyme for digestion (trypsin, lys_c, lys_n, chymotrypsin)")
		databaseCmd.Flags().StringVarP(&m.Database.Tag, "prefix", "", "rev_", "decoy prefix to be added")
		databaseCmd.Flags().StringVarP(&m.Database.Add, "add", "", "", "add custom sequences (UniProt FASTA format only)")
		databaseCmd.Flags().StringVarP(&m.Database.Custom, "custom", "", "", "use a pre formatted custom database")
		databaseCmd.Flags().BoolVarP(&m.Database.Crap, "contam", "", false, "add common contaminants")
		databaseCmd.Flags().BoolVarP(&m.Database.Rev, "reviewed", "", false, "use only reviwed sequences from Swiss-Prot")
		databaseCmd.Flags().BoolVarP(&m.Database.Iso, "isoform", "", false, "add isoform sequences")

		RootCmd.AddCommand(databaseCmd)
	}

}
