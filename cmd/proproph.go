package cmd

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/ext/proteinprophet"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/spf13/cobra"
)

// proprophCmd represents the proproph command
var proprophCmd = &cobra.Command{
	Use:   "proteinprophet",
	Short: "Protein identification validation",
	//Long:  "Statistical validation of protein identification based on peptide assignment to MS/MS spectra\nProteinProphet v5.0",
	Run: func(cmd *cobra.Command, args []string) {

		if len(m.UUID) < 1 && len(m.Home) < 1 {
			e := &err.Error{Type: err.WorkspaceNotFound, Class: err.FATA}
			logrus.Fatal(e.Error())
		}

		logrus.Info("Executing ProteinProphet")
		var pop = proteinprophet.New()

		if len(args) < 1 {
			logrus.Fatal("No input file provided")
		}

		// deploy the binaries
		e := pop.Deploy(m.OS, m.Distro)
		if e != nil {
			logrus.Fatal(e.Message)
		}

		// run ProteinProphet
		xml, e := pop.Run(m.ProteinProphet, m.Home, m.Temp, args)
		if e != nil {
			logrus.Fatal(e.Message)
		}

		_ = xml
		// e = evi.NewInference()
		// if e != nil {
		// 	logrus.Fatal(e.Message)
		// }
		// evi.IndexProteinInference(xml)

		m.ProteinProphet.InputFiles = args
		m.Serialize()

		logrus.Info("Done")
		return
	},
}

func init() {

	if os.Args[1] == "proteinprophet" {

		m.Restore(sys.Meta())

		proprophCmd.Flags().BoolVarP(&m.ProteinProphet.Iprophet, "iprophet", "", false, "input is from iProphet")
		proprophCmd.Flags().BoolVarP(&m.ProteinProphet.ExcludeZ, "excludezeros", "", false, "exclude zero prob entries")
		proprophCmd.Flags().BoolVarP(&m.ProteinProphet.Noprotlen, "noprotlen", "", false, "do not report protein length")
		proprophCmd.Flags().BoolVarP(&m.ProteinProphet.Protmw, "protmw", "", false, "get protein mol weights")
		proprophCmd.Flags().BoolVarP(&m.ProteinProphet.Icat, "icat", "", false, "highlight peptide cysteines")
		proprophCmd.Flags().BoolVarP(&m.ProteinProphet.Glyc, "glyc", "", false, "highlight peptide N-glycosylation motif")
		proprophCmd.Flags().BoolVarP(&m.ProteinProphet.Fpkm, "fpkm", "", false, "model protein FPKM values")
		proprophCmd.Flags().BoolVarP(&m.ProteinProphet.NonSP, "nonsp", "", false, "do not use NSP model")
		proprophCmd.Flags().IntVarP(&m.ProteinProphet.Minindep, "minindep", "", 0, "minimum percentage of independent peptides required for a protein")
		proprophCmd.Flags().Float64VarP(&m.ProteinProphet.Minprob, "minprob", "", 0.05, "peptideProphet probabilty threshold")
		proprophCmd.Flags().IntVarP(&m.ProteinProphet.Maxppmdiff, "maxppmdiff", "", 20, "maximum peptide mass difference in PPM")
		proprophCmd.Flags().BoolVarP(&m.ProteinProphet.Accuracy, "accuracy", "", false, "equivalent to --minprob 0")
		proprophCmd.Flags().BoolVarP(&m.ProteinProphet.Normprotlen, "normprotlen", "", false, "normalize NSP using Protein Length")
		proprophCmd.Flags().BoolVarP(&m.ProteinProphet.Nogroupwts, "nogroupwts", "", false, "check peptide's Protein weight against the threshold (default: check peptide's Protein Group weight against threshold)")
		proprophCmd.Flags().BoolVarP(&m.ProteinProphet.Instances, "instances", "", false, "use Expected Number of Ion Instances to adjust the peptide probabilities prior to NSP adjustment")
		proprophCmd.Flags().BoolVarP(&m.ProteinProphet.Delude, "delude", "", false, "do NOT use peptide degeneracy information when assessing proteins")
		proprophCmd.Flags().BoolVarP(&m.ProteinProphet.Nooccam, "nooccam", "", false, "non-conservative maximum protein list")
		proprophCmd.Flags().BoolVarP(&m.ProteinProphet.Softoccam, "softoccam", "", false, "peptide weights are apportioned equally among proteins within each Protein Group (less conservative protein count estimate)")
		proprophCmd.Flags().BoolVarP(&m.ProteinProphet.Confem, "confem", "", false, "use the EM to compute probability given the confidence")
		proprophCmd.Flags().BoolVarP(&m.ProteinProphet.Logprobs, "logprobs", "", false, "use the log of the probabilities in the Confidence calculations")
		proprophCmd.Flags().BoolVarP(&m.ProteinProphet.Allpeps, "allpeps", "", false, "consider all possible peptides in the database in the confidence model")
		proprophCmd.Flags().IntVarP(&m.ProteinProphet.Mufactor, "mufactor", "", 1, "fudge factor to scale MU calculation")
		proprophCmd.Flags().BoolVarP(&m.ProteinProphet.Unmapped, "unmapped", "", false, "report results for UNMAPPED proteins")
		proprophCmd.Flags().StringVarP(&m.ProteinProphet.Output, "output", "", "interact", "Output name")
	}

	RootCmd.AddCommand(proprophCmd)
}
