package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher-source/lib/meta"
	"github.com/prvst/philosopher-source/lib/rep"
	"github.com/prvst/philosopher-source/lib/sys"
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
			logrus.Fatal("Workspace not found. Run 'philosopher init' to create a workspace")
		}

		repo.Restore()

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

		logrus.Info("Creating PSM report")
		repo.PSMReport()

		logrus.Info("Creating peptide Ion report")
		repo.PeptideIonReport()

		logrus.Info("Creating peptide report")
		repo.PeptideReport()

		if len(repo.Modifications.AssignedBins) > 0 {
			logrus.Info("Creating modification reports")
			repo.ModifiedPSMReport()
			repo.ModifiedPeptideIonReport()
			repo.ModifiedPeptideReport()
			repo.ModificationReport()
		}

		logrus.Info("Done")
		return
	},
}

func init() {

	repo = rep.New()

	RootCmd.AddCommand(reportCmd)
}
