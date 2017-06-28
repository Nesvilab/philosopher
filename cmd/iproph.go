package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/prvst/cmsl/err"
	"github.com/prvst/philosopher/lib/ext/interprophet"
	"github.com/prvst/philosopher/lib/meta"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/spf13/cobra"
)

var ipt interprophet.InterProphet

// iprophCmd represents the iproph command
var iprophCmd = &cobra.Command{
	Use:   "iprophet",
	Short: "MS/MS integrative analysis",
	//Long:  "Multi-level integrative analysis of shotgun proteomic data\niProphet v5.0",
	Run: func(cmd *cobra.Command, args []string) {

		var m meta.Data
		m.Restore(sys.Meta())
		if len(m.UUID) < 1 && len(m.Home) < 1 {
			e := &err.Error{Type: err.WorkspaceNotFound, Class: err.FATA}
			logrus.Fatal(e.Error())
		}

		var err error

		// prepare binaries
		err = ipt.Deploy()
		if err != nil {
			logrus.Fatal(err)
		}

		// run
		err = ipt.Run(args)
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Info("Done")
		return
	},
}

func init() {

	ipt = interprophet.New()

	iprophCmd.Flags().StringVarP(&ipt.Threads, "threads", "", "", "specify threads to use (default 1)")
	iprophCmd.Flags().StringVarP(&ipt.Decoy, "decoy", "", "", "specify the decoy tag")
	iprophCmd.Flags().StringVarP(&ipt.Cat, "cat", "", "", "specify file listing peptide categories")
	iprophCmd.Flags().StringVarP(&ipt.MinProb, "minProb", "", "", "specify minimum probability of results to report")
	iprophCmd.Flags().StringVarP(&ipt.Output, "output", "", "iproph.pep.xml", "specify output name")
	iprophCmd.Flags().BoolVarP(&ipt.Length, "length", "", false, "use Peptide Length model")
	iprophCmd.Flags().BoolVarP(&ipt.Nofpkm, "nofpkm", "", false, "do not use FPKM model")
	iprophCmd.Flags().BoolVarP(&ipt.Nonss, "nonss", "", false, "do not use NSS model")
	iprophCmd.Flags().BoolVarP(&ipt.Nonse, "nonse", "", false, "do not use NSE model")
	iprophCmd.Flags().BoolVarP(&ipt.Nonrs, "nonrs", "", false, "do not use NRS model")
	iprophCmd.Flags().BoolVarP(&ipt.Nonsm, "nonsm", "", false, "do not use NSM model")
	iprophCmd.Flags().BoolVarP(&ipt.Nonsp, "nonsp", "", false, "do not use NSP model")
	iprophCmd.Flags().BoolVarP(&ipt.Sharpnse, "sharpnse", "", false, "Use more discriminating model for NSE in SWATH mode")
	iprophCmd.Flags().BoolVarP(&ipt.Nonsi, "nonsi", "", false, "do not use NSI model")

	RootCmd.AddCommand(iprophCmd)
}
