package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/prvst/cmsl/err"
	"github.com/prvst/philosopher/lib/ext/proteinprophet"
	"github.com/prvst/philosopher/lib/meta"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/spf13/cobra"
)

var pop proteinprophet.ProteinProphet

// proprophCmd represents the proproph command
var proprophCmd = &cobra.Command{
	Use:   "proteinprophet",
	Short: "Protein identification validation",
	//Long:  "Statistical validation of protein identification based on peptide assignment to MS/MS spectra\nProteinProphet v5.0",
	Run: func(cmd *cobra.Command, args []string) {

		var m meta.Data
		m.Restore(sys.Meta())
		if len(m.UUID) < 1 && len(m.Home) < 1 {
			e := &err.Error{Type: err.WorkspaceNotFound, Class: err.FATA}
			logrus.Fatal(e.Error())
		}

		if len(args) < 1 {
			logrus.Fatal("No input file provided")
		}

		// deploy the binaries
		err := pop.Deploy()
		if err != nil {
			logrus.Fatal(err)
		}

		// run ProteinProphet
		err = pop.Run(args)
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Info("Done")

		return
	},
}

func init() {

	pop = proteinprophet.New()

	proprophCmd.Flags().BoolVarP(&pop.Iprophet, "iprophet", "", false, "input is from iProphet")
	proprophCmd.Flags().BoolVarP(&pop.ExcludeZ, "excludezeros", "", false, "exclude zero prob entries")
	proprophCmd.Flags().BoolVarP(&pop.Noprotlen, "noprotlen", "", false, "do not report protein length")
	proprophCmd.Flags().BoolVarP(&pop.Protmw, "protmw", "", false, "get protein mol weights")
	proprophCmd.Flags().BoolVarP(&pop.Icat, "icat", "", false, "highlight peptide cysteines")
	proprophCmd.Flags().BoolVarP(&pop.Glyc, "glyc", "", false, "highlight peptide N-glycosylation motif")
	proprophCmd.Flags().BoolVarP(&pop.Fpkm, "fpkm", "", false, "model protein FPKM values")
	proprophCmd.Flags().BoolVarP(&pop.NonSP, "nonsp", "", false, "do not use NSP model")
	proprophCmd.Flags().IntVarP(&pop.Minindep, "minindep", "", 0, "minimum percentage of independent peptides required for a protein")
	proprophCmd.Flags().Float64VarP(&pop.Minprob, "minprob", "", 0.05, "peptideProphet probabilty threshold")
	proprophCmd.Flags().IntVarP(&pop.Maxppmdiff, "maxppmdiff", "", 20, "maximum peptide mass difference in PPM")
	proprophCmd.Flags().BoolVarP(&pop.Accuracy, "accuracy", "", false, "equivalent to --minprob 0")
	proprophCmd.Flags().BoolVarP(&pop.Normprotlen, "normprotlen", "", false, "normalize NSP using Protein Length")
	proprophCmd.Flags().BoolVarP(&pop.Nogroupwts, "nogroupwts", "", false, "check peptide's Protein weight against the threshold (default: check peptide's Protein Group weight against threshold)")
	proprophCmd.Flags().BoolVarP(&pop.Instances, "instances", "", false, "use Expected Number of Ion Instances to adjust the peptide probabilities prior to NSP adjustment")
	proprophCmd.Flags().BoolVarP(&pop.Delude, "delude", "", false, "do NOT use peptide degeneracy information when assessing proteins")
	proprophCmd.Flags().BoolVarP(&pop.Nooccam, "nooccam", "", false, "non-conservative maximum protein list")
	proprophCmd.Flags().BoolVarP(&pop.Softoccam, "softoccam", "", false, "peptide weights are apportioned equally among proteins within each Protein Group (less conservative protein count estimate)")
	proprophCmd.Flags().BoolVarP(&pop.Confem, "confem", "", false, "use the EM to compute probability given the confidence")
	proprophCmd.Flags().BoolVarP(&pop.Logprobs, "logprobs", "", false, "use the log of the probabilities in the Confidence calculations")
	proprophCmd.Flags().BoolVarP(&pop.Allpeps, "allpeps", "", false, "consider all possible peptides in the database in the confidence model")
	proprophCmd.Flags().IntVarP(&pop.Mufactor, "mufactor", "", 1, "fudge factor to scale MU calculation")
	proprophCmd.Flags().BoolVarP(&pop.Unmapped, "unmapped", "", false, "report results for UNMAPPED proteins")
	proprophCmd.Flags().StringVarP(&pop.Output, "output", "", "interact", "Output name")
	//proprophCmd.Flags().BoolVarP(&pop.Asap, "asap", "", false, "compute ASAP ratios for protein entries (ASAP must have been run previously on interact dataset)")
	//proprophCmd.Flags().BoolVarP(&pop.Refresh, "refresh", "", false, "import manual changes to AAP ratios (after initially using ASAP option)")
	//proprophCmd.Flags().BoolVarP(&pop.Asapprophet, "asapprophet", "", false, "*new and Improved* compute ASAP ratios for protein entries (ASAP must have been run previously on all input interact datasets with mz/XML raw data format)")
	//proprophCmd.Flags().BoolVarP(&pop.Excludemods, "excludemods", "", false, "Exclude modified peptides (aside from those identified with variable modifications or isotope error correction) to be used for protein inference")

	RootCmd.AddCommand(proprophCmd)
}
