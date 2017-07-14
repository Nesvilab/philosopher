package cmd

import (
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/cmsl/err"
	"github.com/prvst/philosopher/lib/meta"
	"github.com/prvst/philosopher/lib/quan"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/spf13/cobra"
)

var qnt quan.Quantify

// quantifyCmd represents the quantify command
var freequant = &cobra.Command{
	Use:   "freequant",
	Short: "Label-free Quantification ",
	//Long:  `Provides methods for MS1 Peak Intensity calculation based on XIC`,
	Run: func(cmd *cobra.Command, args []string) {

		var m meta.Data
		m.Restore(sys.Meta())
		if len(m.UUID) < 1 && len(m.Home) < 1 {
			e := &err.Error{Type: err.WorkspaceNotFound, Class: err.FATA}
			logrus.Fatal(e.Error())
		}

		if len(qnt.Format) < 1 || len(qnt.Dir) < 1 {
			logrus.Fatal("You need to provide the path to the mz files and the correct extension.")
		}

		if strings.EqualFold(qnt.Format, "mzml") {
			qnt.Format = "mzML"
		} else if strings.EqualFold(qnt.Format, "mzxml") {
			qnt.Format = "mzXML"
		} else {
			logrus.Fatal("Unknown file format")
		}

		qnt.RunLabelFreeQuantification()

		logrus.Info("Done")
		return
	},
}

func init() {

	qnt = quan.New()

	freequant.Flags().Float64VarP(&qnt.Tol, "tol", "", 10, "m/z tolerance in ppm")
	freequant.Flags().StringVarP(&qnt.Dir, "dir", "", "", "folder path containing the raw files")
	freequant.Flags().StringVarP(&qnt.Format, "ext", "", "", "spectra file extension (mzML, mzXML)")
	freequant.Flags().Float64VarP(&qnt.RTWin, "rtw", "", 3, "specify the retention time window for xic (minute)")
	freequant.Flags().Float64VarP(&qnt.PTWin, "ptw", "", 0.2, "specify the time windows for the peak (minute)")

	RootCmd.AddCommand(freequant)
}
