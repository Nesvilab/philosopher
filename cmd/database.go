package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher/lib/dat"
	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/spf13/cobra"
)

var dtb dat.Base

var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Target-Decoy database formatting",
	//Long: `The database command alows the creation and formatting of a Target-Decoy database. It also
	//provides options for downloading a fresh snapshot from UniProt`,
	Run: func(cmd *cobra.Command, args []string) {

		var m met.Data
		m.Restore(sys.Meta())
		if len(m.UUID) < 1 && len(m.Home) < 1 {
			e := &err.Error{Type: err.WorkspaceNotFound, Class: err.FATA}
			logrus.Fatal(e.Error())
		}

		if len(dtb.Annot) > 0 {

			logrus.Info("Processing database")

			var u dat.Base
			err := u.ProcessDB(dtb.Annot, dtb.Tag)
			if err != nil {
				logrus.Fatal(err)
			}

			err = u.Serialize()
			if err != nil {
				logrus.Fatal(err)
			}

			logrus.Info("Done")
			return

		}

		if len(dtb.ID) < 1 && len(dtb.Custom) < 1 {
			logrus.Fatal("You need to provide a taxon ID or a custom FASTA file")
		}

		if dtb.Crap == false {
			logrus.Warning("Contaminants are not going to be added to database")
		}

		if len(dtb.Custom) < 1 {
			logrus.Info("Fetching database")
			dtb.Fetch()
		} else {
			dtb.UniProtDB = dtb.Custom
		}

		logrus.Info("Processing decoys")
		dtb.Create()

		logrus.Info("Creating file")
		dtb.Save()

		err := dtb.Serialize()
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Info("Done")
		return
	},
}

func init() {

	dtb = dat.New()

	databaseCmd.Flags().StringVarP(&dtb.ID, "id", "", "", "UniProt proteome ID")
	databaseCmd.Flags().StringVarP(&dtb.Annot, "annotate", "", "", "process a ready-to-use database")
	databaseCmd.Flags().StringVarP(&dtb.Enz, "enzyme", "", "trypsin", "enzyme for digestion (trypsin, lys_c, lys_n, chymotrypsin)")
	databaseCmd.Flags().StringVarP(&dtb.Tag, "prefix", "", "rev_", "decoy prefix to be added")
	databaseCmd.Flags().StringVarP(&dtb.Add, "add", "", "", "add custom sequences (UniProt FASTA format only)")
	databaseCmd.Flags().StringVarP(&dtb.Custom, "custom", "", "", "use a pre formatted custom database")
	databaseCmd.Flags().BoolVarP(&dtb.Crap, "contam", "", false, "add common contaminants")
	databaseCmd.Flags().BoolVarP(&dtb.Rev, "reviewed", "", false, "use only reviwed sequences from Swiss-Prot")
	databaseCmd.Flags().BoolVarP(&dtb.Iso, "isoform", "", false, "add isoform sequences")

	RootCmd.AddCommand(databaseCmd)
}
