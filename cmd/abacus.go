package cmd

import (
	"os"

	"github.com/prvst/philosopher/lib/aba"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// abacusCmd represents the abacus command
var abacusCmd = &cobra.Command{
	Use:   "abacus",
	Short: "Combined analysis of LC-MS/MS results",
	Run: func(cmd *cobra.Command, args []string) {

		e := m.FunctionInitCheckUp()
		if e != nil {
			logrus.Fatal(e)
		}

		if len(args) < 2 {
			logrus.Fatal("The combined analysis needs at least 2 result files to work")
		}

		logrus.Info("Executing Abacus ", Version)
		err := aba.Run(m.Abacus, m.Temp, args)
		if err != nil {
			logrus.Fatal(err)
		}

		// store parameters on meta data
		m.Serialize()

		// clean tmp
		met.CleanTemp(m.Temp)

		logrus.Info("Done")
		return
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "abacus" {

		m.Restore(sys.Meta())

		abacusCmd.Flags().StringVarP(&m.Abacus.CombPro, "protein", "", "", "combined protein file")
		abacusCmd.Flags().StringVarP(&m.Abacus.CombPep, "peptide", "", "", "combined peptide file")
		abacusCmd.Flags().StringVarP(&m.Abacus.Tag, "tag", "", "rev_", "decoy tag")
		abacusCmd.Flags().Float64VarP(&m.Abacus.ProtProb, "prtProb", "", 0.9, "minimum protein probability")
		abacusCmd.Flags().Float64VarP(&m.Abacus.PepProb, "pepProb", "", 0.5, "minimum peptide probability")
		abacusCmd.Flags().BoolVarP(&m.Abacus.Razor, "razor", "", false, "use razor peptides for protein FDR scoring")
		abacusCmd.Flags().BoolVarP(&m.Abacus.Picked, "picked", "", false, "apply the picked FDR algorithm before the protein scoring")
		abacusCmd.Flags().BoolVarP(&m.Abacus.Unique, "uniqueonly", "", false, "report TMT quantification based on only unique peptides")
		abacusCmd.Flags().BoolVarP(&m.Abacus.Labels, "labels", "", false, "indicates whether the data sets includes TMT labels or not")
		abacusCmd.Flags().BoolVarP(&m.Abacus.Reprint, "reprint", "", false, "create abacus reports using the Reprint format")
	}

	RootCmd.AddCommand(abacusCmd)
}
