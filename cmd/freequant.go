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
	//Long:  `Provides methods for MS1 Peak Intensity calculation based on XIC`,
	Run: func(cmd *cobra.Command, args []string) {

		// if len(m.UUID) < 1 && len(m.Home) < 1 {
		// 	e := &err.Error{Type: err.WorkspaceNotFound, Class: err.FATA}
		// 	logrus.Fatal(e.Error())
		// }

		m.FunctionInitCheckUp()

		// if len(m.Quantify.Format) < 1 || len(m.Quantify.Dir) < 1 {
		// 	logrus.Fatal("You need to provide the path to the mz files and the correct extension.")
		// }

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
		e := qua.RunLabelFreeQuantification(m.Quantify)
		if e != nil {
			logrus.Fatal(e.Error())
		}

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

		freequant.Flags().Float64VarP(&m.Quantify.Tol, "tol", "", 10, "m/z tolerance in ppm")
		freequant.Flags().StringVarP(&m.Quantify.Dir, "dir", "", "", "folder path containing the raw files")
		//freequant.Flags().StringVarP(&m.Quantify.Format, "ext", "", "", "spectra file extension (mzML, mzXML)")
		//freequant.Flags().Float64VarP(&m.Quantify.RTWin, "rtw", "", 3, "specify the retention time window for xic (minute)")
		freequant.Flags().Float64VarP(&m.Quantify.PTWin, "ptw", "", 0.4, "specify the time windows for the peak (minute)")
	}

	RootCmd.AddCommand(freequant)
}
