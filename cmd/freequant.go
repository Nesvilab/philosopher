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

// quantifyCmd represents the quantify command
var freequant = &cobra.Command{
	Use:   "freequant",
	Short: "Label-free Quantification ",
	Run: func(cmd *cobra.Command, args []string) {

		m.FunctionInitCheckUp()

		m.Quantify.Format = "mzML"
		if len(m.Quantify.Dir) < 1 {
			logrus.Fatal("You need to provide the path to the mz files and the correct extension.")
		}

		logrus.Info("Executing label-free quantification ", Version)

		if strings.EqualFold(m.Quantify.Format, "mzml") {
			m.Quantify.Format = "mzML"
		} else if strings.EqualFold(m.Quantify.Format, "mzxml") {
			logrus.Fatal("Only the mzML format is supported")
			m.Quantify.Format = "mzXML"
		} else {
			logrus.Fatal("Unknown file format")
		}

		//forcing the larger time window to be the same as the smaller one
		m.Quantify.RTWin = m.Quantify.PTWin

		// run label-free quantification
		qua.RunLabelFreeQuantification(m.Quantify)

		// store paramters on meta data
		m.Serialize()

		// clean tmp
		met.CleanTemp(m.Temp)

		logrus.Info("Done")
		return
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "freequant" {

		m.Restore(sys.Meta())

		freequant.Flags().StringVarP(&m.Quantify.Dir, "dir", "", "", "folder path containing the raw files")
		freequant.Flags().BoolVarP(&m.Quantify.Isolated, "isolated", "", false, "use the isolated ion instead of the selected ion for quantification")
		freequant.Flags().Float64VarP(&m.Quantify.Tol, "tol", "", 10, "m/z tolerance in ppm")
		freequant.Flags().Float64VarP(&m.Quantify.PTWin, "ptw", "", 0.4, "specify the time windows for the peak (minute)")
	}

	RootCmd.AddCommand(freequant)
}
