// Package cmd Freeequant top level command
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

// quantifyCmd represents the quantify command
var freequant = &cobra.Command{
	Use:   "freequant",
	Short: "Label-free Quantification ",
	Run: func(cmd *cobra.Command, args []string) {

		m.FunctionInitCheckUp()

		m.Quantify.Format = "mzML"
		if len(m.Quantify.Dir) < 1 {
			msg.InputNotFound(errors.New("You need to provide the path to the mz files and the correct extension"), "fatal")
		}

		msg.Executing("Label-free quantification ", Version)

		if strings.EqualFold(m.Quantify.Format, "mzml") {
			m.Quantify.Format = "mzML"
		} else if strings.EqualFold(m.Quantify.Format, "mzxml") {
			msg.InputNotFound(errors.New("Only the mzML format is supported"), "fatal")
			m.Quantify.Format = "mzXML"
		} else {
			msg.InputNotFound(errors.New("Unknown file format"), "fatal")
		}

		//forcing the larger time window to be the same as the smaller one
		m.Quantify.RTWin = m.Quantify.PTWin

		// run label-free quantification
		qua.RunLabelFreeQuantification(m.Quantify)

		// store parameters on meta data
		m.Serialize()

		// clean tmp
		met.CleanTemp(m.Temp)

		msg.Done()
		return
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "freequant" {

		m.Restore(sys.Meta())

		freequant.Flags().StringVarP(&m.Quantify.Dir, "dir", "", "", "folder path containing the raw files")
		freequant.Flags().Float64VarP(&m.Quantify.Tol, "tol", "", 10, "m/z tolerance in ppm")
		freequant.Flags().Float64VarP(&m.Quantify.PTWin, "ptw", "", 0.4, "specify the time windows for the peak (minute)")
		//freequant.Flags().BoolVarP(&m.Quantify.Isolated, "isolated", "", true, "use the isolated ion instead of the selected ion for quantification")
	}

	RootCmd.AddCommand(freequant)
}
