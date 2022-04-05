// Package cmd Labelquant top level command
package cmd

import (
	"errors"
	"os"
	"strings"

	"philosopher/lib/met"
	"philosopher/lib/msg"
	"philosopher/lib/qua"
	"philosopher/lib/sys"

	"github.com/spf13/cobra"
)

// labelquantCmd represents the labelquant command
var labelquantCmd = &cobra.Command{
	Use:   "labelquant",
	Short: "Isobaric Labeling-Based Relative Quantification ",
	Run: func(cmd *cobra.Command, args []string) {

		m.FunctionInitCheckUp()

		m.Quantify.Format = "mzML"

		if len(m.Quantify.Format) < 1 || len(m.Quantify.Dir) < 1 {
			msg.InputNotFound(errors.New("you need to provide the path to the mz files and the correct extension"), "fatal")
		}

		if len(m.Quantify.Plex) < 1 {
			msg.InputNotFound(errors.New("you need to specify the experiment Plex"), "fatal")
		}

		msg.Executing("Isobaric-label quantification ", Version)

		if strings.EqualFold(strings.ToLower(m.Quantify.Format), "mzml") {
			m.Quantify.Format = "mzML"
		} else if strings.EqualFold(m.Quantify.Format, "mzxml") {
			msg.InputNotFound(errors.New("only the mzML format is supported"), "fatal")
			m.Quantify.Format = "mzXML"
		} else {
			msg.InputNotFound(errors.New("unknown file format"), "fatal")
		}

		m.Quantify = qua.RunIsobaricLabelQuantification(m.Quantify, m.Filter.Mapmods)

		// store parameters on meta data
		m.Serialize()

		// clean tmp
		met.CleanTemp(m.Temp)

		msg.Done()
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "labelquant" {

		m.Restore(sys.Meta())

		labelquantCmd.Flags().StringVarP(&m.Quantify.Annot, "annot", "", "", "annotation file with custom names for the TMT channels")
		labelquantCmd.Flags().StringVarP(&m.Quantify.Plex, "plex", "", "", "number of reporter ion channels")
		labelquantCmd.Flags().StringVarP(&m.Quantify.Dir, "dir", "", "", "folder path containing the raw files")
		labelquantCmd.Flags().StringVarP(&m.Quantify.Brand, "brand", "", "", "isobaric labeling brand (tmt, itraq)")
		labelquantCmd.Flags().Float64VarP(&m.Quantify.Tol, "tol", "", 20, "m/z tolerance in ppm")
		labelquantCmd.Flags().IntVarP(&m.Quantify.Level, "level", "", 2, "ms level for the quantification")
		labelquantCmd.Flags().Float64VarP(&m.Quantify.Purity, "purity", "", 0.5, "ion purity threshold")
		labelquantCmd.Flags().Float64VarP(&m.Quantify.MinProb, "minprob", "", 0.7, "only use PSMs with the specified minimum probability score")
		labelquantCmd.Flags().Float64VarP(&m.Quantify.RemoveLow, "removelow", "", 0.0, "ignore the lower % of PSMs based on their summed abundances. 0 means no removal, entry value must be a decimal")
		labelquantCmd.Flags().BoolVarP(&m.Quantify.Unique, "uniqueonly", "", false, "report quantification based only on unique peptides")
		labelquantCmd.Flags().BoolVarP(&m.Quantify.BestPSM, "bestpsm", "", false, "select the best PSMs for protein quantification")
		//labelquantCmd.Flags().BoolVarP(&m.Quantify.Raw, "raw", "", false, "read raw files instead of converted XML")

	}

	RootCmd.AddCommand(labelquantCmd)
}
