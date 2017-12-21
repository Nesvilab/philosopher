package cmd

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/rep"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/spf13/cobra"
)

// reportCmd represents the report commands
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Multi-level reporting for both narrow-searches and open-searches",
	Run: func(cmd *cobra.Command, args []string) {

		if len(m.UUID) < 1 && len(m.Home) < 1 {
			e := &err.Error{Type: err.WorkspaceNotFound, Class: err.FATA}
			logrus.Fatal(e.Error())
		}

		logrus.Info("Executing report")
		var repo = rep.New()

		err := repo.RestoreGranular()
		if err != nil {
			logrus.Fatal(err.Error())
		}

		if len(repo.Proteins) > 0 {

			logrus.Info("Creating Protein FASTA report")
			repo.ProteinFastaReport()

			if repo.Proteins[0].TotalLabels.Channel1.Intensity > 0 || repo.Proteins[10].TotalLabels.Channel1.Intensity > 0 {
				logrus.Info("Creating Protein TMT report")
				repo.ProteinTMTReport(m.Quantify.Unique)
			} else {
				logrus.Info("Creating Protein report")
				repo.ProteinReport()
			}

		}

		// verifying if there is any quantification on labels
		if len(m.Quantify.Plex) > 0 {

			logrus.Info("Creating TMT PSM report")
			repo.PSMTMTReport(m.Filter.Tag)

			logrus.Info("Creating TMT peptide report")
			repo.PeptideTMTReport()

			logrus.Info("Creating TMT peptide Ion report")
			repo.PeptideIonTMTReport()

		} else {

			logrus.Info("Creating PSM report")
			repo.PSMReport(m.Filter.Tag)

			logrus.Info("Creating peptide report")
			repo.PeptideReport()

			logrus.Info("Creating peptide Ion report")
			repo.PeptideIonReport()

		}

		if len(repo.Modifications.MassBins) > 0 {
			logrus.Info("Creating modification reports")
			repo.ModificationReport()

			logrus.Info("Plotting mass distribution")
			repo.PlotMassHist()
		}

		// store parameters on meta data
		m.Serialize()

		logrus.Info("Done")
		return
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "report" {

		m.Restore(sys.Meta())

		reportCmd.Flags().BoolVarP(&m.Report.Decoys, "decoys", "", false, "add decoy observations to reports")
	}

	RootCmd.AddCommand(reportCmd)
}
