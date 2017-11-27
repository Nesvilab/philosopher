package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher/lib/aba"
	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/meta"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/spf13/cobra"
)

var a aba.Abacus

// abacusCmd represents the abacus command
var abacusCmd = &cobra.Command{
	Use:   "abacus",
	Short: "Combined analysis of LC-MS/MS results",
	//Long:  "Abacus aggregates data from multiple experiments, adjusts spectral counts to accurately account for peptides shared across multiple proteins, and performs common normalization steps",
	Run: func(cmd *cobra.Command, args []string) {

		var m meta.Data
		m.Restore(sys.Meta())
		if len(m.UUID) < 1 && len(m.Home) < 1 {
			e := &err.Error{Type: err.WorkspaceNotFound, Class: err.FATA}
			logrus.Fatal(e.Error())
		}

		if len(args) < 2 {
			logrus.Fatal("The combined analysis needs at least 2 result files to work")
		}

		err := a.Run(args)
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Info("Done")
		return
	},
}

func init() {

	a = aba.New()

	abacusCmd.Flags().StringVarP(&a.Comb, "comb", "", "", "combined file")
	abacusCmd.Flags().BoolVarP(&a.Razor, "razor", "", false, "use razor peptides for protein FDR scoring")
	abacusCmd.Flags().StringVarP(&a.Tag, "tag", "", "rev_", "decoy tag")
	abacusCmd.Flags().Float64VarP(&a.ProtProb, "prtProb", "", 0.9, "minimun protein probability")
	abacusCmd.Flags().Float64VarP(&a.PepProb, "pepProb", "", 0.5, "minimun peptide probability")

	RootCmd.AddCommand(abacusCmd)
}
