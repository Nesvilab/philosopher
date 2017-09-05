package cmd

import (
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/cmsl/err"
	"github.com/prvst/philosopher/lib/meta"
	"github.com/prvst/philosopher/lib/quan"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/spf13/cobra"
)

var lbl quan.Quantify

// labelquantCmd represents the labelquant command
var labelquantCmd = &cobra.Command{
	Use:   "labelquant",
	Short: "Isobaric Labeling-Based Relative Quantification ",
	//Long:  `Provides methods for labeled data quantification`,
	Run: func(cmd *cobra.Command, args []string) {

		var m meta.Data
		m.Restore(sys.Meta())
		if len(m.UUID) < 1 && len(m.Home) < 1 {
			e := &err.Error{Type: err.WorkspaceNotFound, Class: err.FATA}
			logrus.Fatal(e.Error())
		}

		if len(lbl.Format) < 1 || len(lbl.Dir) < 1 {
			logrus.Fatal("You need to provide the path to the mz files and the correct extension.")
		}

		// hardcoded tmt for now
		lbl.Brand = "tmt"

		if strings.EqualFold(strings.ToLower(lbl.Format), "mzml") {
			lbl.Format = "mzML"
		} else if strings.EqualFold(lbl.Format, "mzxml") {
			lbl.Format = "mzXML"
		} else {
			logrus.Fatal("Unknown file format")
		}

		if len(lbl.ChanNorm) > 0 && lbl.IntNorm == true {
			i, err := strconv.Atoi(lbl.ChanNorm)
			if i > 10 || i < 1 || err != nil {
				logrus.Fatal("Inexisting channel number:", lbl.ChanNorm)
			}
			logrus.Fatal("You can choose only one method of normalization")
		} else if len(lbl.ChanNorm) == 0 && lbl.IntNorm == false {
			logrus.Fatal("Missing normalization method, type 'philosopher labelquant --help' for more information")
		}

		err := lbl.RunLabeledQuantification()
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Info("Done")
		return
	},
}

func init() {
	//labelquantCmd.Flags().StringVarP(&lbl.Brand, "brand", "", "", "type of label (tmt or itraq)")
	labelquantCmd.Flags().StringVarP(&lbl.Plex, "plex", "", "", "number of channels")
	labelquantCmd.Flags().Float64VarP(&lbl.Tol, "tol", "", 10, "m/z tolerance in ppm")
	labelquantCmd.Flags().StringVarP(&lbl.Dir, "dir", "", "", "folder path containing the raw files")
	labelquantCmd.Flags().StringVarP(&lbl.Format, "ext", "", "", "spectra file extension (mzML, mzXML)")
	labelquantCmd.Flags().Float64VarP(&lbl.Purity, "purity", "", 0.5, "ion purity threshold")
	labelquantCmd.Flags().StringVarP(&lbl.ChanNorm, "normToChannel", "", "", "normalize intensities to a control channel (provide a channel number as control)")
	labelquantCmd.Flags().BoolVarP(&lbl.IntNorm, "normToIntensity", "", false, "normalize intensities to the total intensity from all channels")

	RootCmd.AddCommand(labelquantCmd)

}
