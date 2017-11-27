package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/rep"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/spf13/cobra"
)

var repo rep.Evidence

// reportCmd represents the report commands
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Multi-level reporting for both narrow-searches and open-searches",
	//Long:  `Creates peptide-level and protein-level reportsbased on the experimental results.`,
	Run: func(cmd *cobra.Command, args []string) {

		var m met.Data
		m.Restore(sys.Meta())
		if len(m.UUID) < 1 && len(m.Home) < 1 {
			e := &err.Error{Type: err.WorkspaceNotFound, Class: err.FATA}
			logrus.Fatal(e.Error())
		}

		//repo.Restore()
		err := repo.RestoreGranular()
		if err != nil {
			logrus.Fatal(err.Error())
		}

		if len(repo.Proteins) > 0 {

			logrus.Info("Creating Protein FASTA report")
			repo.ProteinFastaReport()

			if repo.Proteins[0].TotalLabels.Channel1.Intensity > 0 || repo.Proteins[10].TotalLabels.Channel1.Intensity > 0 {
				logrus.Info("Creating Protein TMT report")
				repo.ProteinTMTReport()
			} else {
				logrus.Info("Creating Protein report")
				repo.ProteinReport()
			}

		}

		// verifying if there is any quantification on labels
		var lblMarker float64
		for i := 0; i <= 1000; i++ {
			lblMarker += repo.PSM[i].Labels.Channel1.Intensity
		}

		if lblMarker > 0 {

			logrus.Info("Creating TMT PSM report")
			repo.PSMQTMTReport()

			logrus.Info("Creating TMT peptide report")
			repo.PeptideTMTReport()

			logrus.Info("Creating TMT peptide Ion report")
			repo.PeptideIonTMTReport()

		} else {

			logrus.Info("Creating PSM report")
			repo.PSMReport()

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

		logrus.Info("Done")
		return
	},
}

func init() {

	repo = rep.New()

	reportCmd.Flags().BoolVarP(&repo.Decoys, "decoys", "", false, "add decoy observations to reports")

	RootCmd.AddCommand(reportCmd)
}
