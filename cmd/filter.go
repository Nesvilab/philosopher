// Package cmd Filter top level command
package cmd

import (
	"errors"
	"os"

	"github.com/prvst/philosopher/lib/fil"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/msg"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/spf13/cobra"
)

// filterCmd represents the filter command
var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Statistical filtering, validation and False Discovery Rates assessment",
	Run: func(cmd *cobra.Command, args []string) {

		m.FunctionInitCheckUp()

		msg.Executing("Filter ", Version)

		// clean, clean clean
		os.RemoveAll(sys.EvBin())
		os.RemoveAll(sys.EvIonBin())
		os.RemoveAll(sys.EvModificationsBin())
		os.RemoveAll(sys.EvModificationsEvBin())
		os.RemoveAll(sys.EvPSMBin())
		os.RemoveAll(sys.EvPeptideBin())
		os.RemoveAll(sys.EvProteinBin())
		os.RemoveAll(sys.PsmBin())
		os.RemoveAll(sys.IonBin())
		os.RemoveAll(sys.PepBin())
		os.RemoveAll(sys.PepxmlBin())
		os.RemoveAll(sys.ProBin())
		os.RemoveAll(sys.ProtxmlBin())

		// check file existence
		if len(m.Filter.Pex) < 1 {
			msg.InputNotFound(errors.New("You must provide a pepXML file or a folder with one or more files, Run 'philosopher filter --help' for more information"), "fatal")
		}

		if len(m.Filter.Pox) == 0 && m.Filter.Razor == true {
			msg.Custom(errors.New("Razor option will be ignored because there is no protein inference data"), "warning")
			m.Filter.Razor = false
		}

		m := fil.Run(m)

		m.Serialize()

		// clean tmp
		met.CleanTemp(m.Temp)

		msg.Done()
		return
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "filter" {

		m.Restore(sys.Meta())

		filterCmd.Flags().StringVarP(&m.Filter.Pex, "pepxml", "", "", "pepXML file or directory containing a set of pepXML files")
		filterCmd.Flags().StringVarP(&m.Filter.Pox, "protxml", "", "", "protXML file path")
		filterCmd.Flags().StringVarP(&m.Filter.Tag, "tag", "", "rev_", "decoy tag")
		filterCmd.Flags().Float64VarP(&m.Filter.IonFDR, "ion", "", 0.01, "peptide ion FDR level")
		filterCmd.Flags().Float64VarP(&m.Filter.PepFDR, "pep", "", 0.01, "peptide FDR level")
		filterCmd.Flags().Float64VarP(&m.Filter.PsmFDR, "psm", "", 0.01, "psm FDR level")
		filterCmd.Flags().Float64VarP(&m.Filter.PtFDR, "prot", "", 0.01, "protein FDR level")
		filterCmd.Flags().Float64VarP(&m.Filter.PepProb, "pepProb", "", 0.7, "top peptide probability treshold for the FDR filtering")
		filterCmd.Flags().Float64VarP(&m.Filter.ProtProb, "protProb", "", 0.5, "protein probability treshold for the FDR filtering (not used with the razor algorithm)")
		filterCmd.Flags().Float64VarP(&m.Filter.Weight, "weight", "", 1, "threshold for defining peptide uniqueness")
		filterCmd.Flags().BoolVarP(&m.Filter.Seq, "sequential", "", false, "alternative algorithm that estimates FDR using both filtered PSM and Protein lists")
		filterCmd.Flags().BoolVarP(&m.Filter.Cap, "cappedsequential", "", false, "alternative algorithm that estimates FDR using both filtered PSM and Protein lists using a threshold cap from first pass")
		filterCmd.Flags().BoolVarP(&m.Filter.Model, "models", "", false, "print model distribution")
		filterCmd.Flags().BoolVarP(&m.Filter.Razor, "razor", "", false, "use razor peptides for protein FDR scoring")
		filterCmd.Flags().BoolVarP(&m.Filter.Picked, "picked", "", false, "apply the picked FDR algorithm before the protein scoring")
		filterCmd.Flags().BoolVarP(&m.Filter.Mapmods, "mapmods", "", false, "map modifications aquired by an open search")
		filterCmd.Flags().BoolVarP(&m.Filter.Fo, "fo", "", false, "")
		filterCmd.Flags().MarkHidden("fo")
	}

	RootCmd.AddCommand(filterCmd)
}
