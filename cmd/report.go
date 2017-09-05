package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/prvst/cmsl/err"
	"github.com/prvst/philosopher/lib/meta"
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

		var m meta.Data
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

			logrus.Info("Creating Protein identification report")
			if repo.Proteins[0].TotalLabels.Channel1.Mean > 0 || repo.Proteins[10].TotalLabels.Channel1.Mean > 0 {
				repo.ProteinQuantReport()
			} else {
				repo.ProteinReport()
			}

		}

		if repo.PSM[0].Labels.Channel1.Intensity > 0 || repo.PSM[10].Labels.Channel1.Intensity > 0 || repo.PSM[100].Labels.Channel1.Intensity > 0 {
			logrus.Info("Creating labeled PSM report")
			repo.PSMQuantReport()
		} else {
			logrus.Info("Creating PSM report")
			repo.PSMReport()
		}

		logrus.Info("Creating peptide Ion report")
		repo.PeptideIonReport()

		logrus.Info("Creating peptide report")
		repo.PeptideReport()

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

	RootCmd.AddCommand(reportCmd)
}
