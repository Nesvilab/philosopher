package cmd

import (
	"os"
	"strings"

	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/qua"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/sirupsen/logrus"
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
			logrus.Fatal("You need to provide the path to the mz files and the correct extension.")
		}

		if len(m.Quantify.Plex) < 1 {
			logrus.Fatal("You need to especify the experiment Plex")
		}

		// hardcoded tmt for now
		logrus.Info("Executing label-based quantification ", Version)
		m.Quantify.Brand = "tmt"

		if strings.EqualFold(strings.ToLower(m.Quantify.Format), "mzml") {
			m.Quantify.Format = "mzML"
		} else if strings.EqualFold(m.Quantify.Format, "mzxml") {
			logrus.Fatal("Only the mzML format is supported")
			m.Quantify.Format = "mzXML"
		} else {
			logrus.Fatal("Unknown file format")
		}

		m.Quantify = qua.RunTMTQuantification(m.Quantify, m.Filter.Mapmods)

		// store paramters on meta data
		m.Serialize()

		// clean tmp
		met.CleanTemp(m.Temp)

		logrus.Info("Done")
		return
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "labelquant" {

		m.Restore(sys.Meta())

		labelquantCmd.Flags().StringVarP(&m.Quantify.Annot, "annot", "", "", "annotation file with custom names for the TMT channels")
		labelquantCmd.Flags().StringVarP(&m.Quantify.Plex, "plex", "", "", "number of channels")
		labelquantCmd.Flags().StringVarP(&m.Quantify.Dir, "dir", "", "", "folder path containing the raw files")
		labelquantCmd.Flags().Float64VarP(&m.Quantify.Tol, "tol", "", 20, "m/z tolerance in ppm")
		labelquantCmd.Flags().IntVarP(&m.Quantify.Level, "level", "", 2, "ms level for the quantification")
		labelquantCmd.Flags().Float64VarP(&m.Quantify.Purity, "purity", "", 0.5, "ion purity threshold")
		labelquantCmd.Flags().Float64VarP(&m.Quantify.MinProb, "minprob", "", 0.7, "only use PSMs with a minimun probability score")
		labelquantCmd.Flags().Float64VarP(&m.Quantify.RemoveLow, "removelow", "", 0.0, "ignore the lower % of PSMs based on their summed abundances. 0 Means no removal, entry value must be decimal")
		labelquantCmd.Flags().BoolVarP(&m.Quantify.Unique, "uniqueonly", "", false, "report quantification based on only unique peptides")
		labelquantCmd.Flags().BoolVarP(&m.Quantify.BestPSM, "bestpsm", "", false, "select the best PSMs for protein quantification")

	}

	RootCmd.AddCommand(labelquantCmd)
}
