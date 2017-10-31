package cmd

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/cmsl/err"
	"github.com/prvst/philosopher/lib/fil"
	"github.com/prvst/philosopher/lib/meta"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/spf13/cobra"
)

var fp fil.Filter

// filterCmd represents the filter command
var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Statistical filtering, validation and False Discovery Rates assessment",
	//Long:  `Custom algorithms for multi-level False Discovery Rates scoring and evaluation`,
	Run: func(cmd *cobra.Command, args []string) {

		var m meta.Data
		m.Restore(sys.Meta())
		if len(m.UUID) < 1 && len(m.Home) < 1 {
			e := &err.Error{Type: err.WorkspaceNotFound, Class: err.FATA}
			logrus.Fatal(e.Error())
		}

		/// clean, clean clean
		os.RemoveAll(sys.EvBin())
		os.RemoveAll(sys.EvIonBin())
		os.RemoveAll(sys.EvModificationsBin())
		os.RemoveAll(sys.EvModificationsEvBin())
		os.RemoveAll(sys.EvPSMBin())
		os.RemoveAll(sys.EvPeptideBin())
		os.RemoveAll(sys.EvProteinBin())
		os.RemoveAll(sys.PsmBin())
		os.RemoveAll(sys.IonBin())
		os.RemoveAll(sys.PepBin())
		os.RemoveAll(sys.PepxmlBin())
		os.RemoveAll(sys.ProBin())
		os.RemoveAll(sys.ProtxmlBin())

		// check file existence
		if len(fp.Pex) < 1 {
			logrus.Fatal("You must provide a pepXML file or a folder with one or more files, Run 'philosopher filter --help' for more information")
		}

		//stat.Run(fp, psmFDR, pepFDR, ionFDR, prtFDR, pepProb, prtProb)
		e := fp.Run(fp.Psmfdr, fp.Pepfdr, fp.Ionfdr, fp.Ptfdr, fp.PepProb, fp.ProtProb, fp.Picked, fp.Razor, fp.Mapmods)
		if e != nil {
			logrus.Fatal(e.Error())
		}

		// m.Experimental.DecoyTag = fp.Tag
		// m.Experimental.ConTag = fp.Con
		// m.Experimental.PsmFDR = psmFDR
		// m.Experimental.PepFDR = pepFDR
		// m.Experimental.IonFDR = ionFDR
		// m.Experimental.PrtFDR = prtFDR
		// m.Experimental.PepProb = pepProb
		// m.Experimental.PrtProb = prtProb

		m.Serialize()

		logrus.Info("Done")
		return
	},
}

func init() {

	fp = fil.New()

	filterCmd.Flags().StringVarP(&fp.Pex, "pepxml", "", "", "pepXML file or directory containing a set of pepXML files")
	filterCmd.Flags().StringVarP(&fp.Pox, "protxml", "", "", "protXML file path")
	filterCmd.Flags().StringVarP(&fp.Tag, "tag", "", "rev_", "decoy tag")
	filterCmd.Flags().Float64VarP(&fp.Ionfdr, "ion", "", 0.01, "peptide ion FDR level")
	filterCmd.Flags().Float64VarP(&fp.Pepfdr, "pep", "", 0.01, "peptide FDR level")
	filterCmd.Flags().Float64VarP(&fp.Psmfdr, "psm", "", 0.01, "psm FDR level")
	filterCmd.Flags().Float64VarP(&fp.Ptfdr, "prot", "", 0.01, "protein FDR level")
	filterCmd.Flags().Float64VarP(&fp.PepProb, "pepProb", "", 0.7, "top peptide probability treshold for the FDR filtering")
	filterCmd.Flags().Float64VarP(&fp.ProtProb, "protProb", "", 0.5, "protein probability treshold for the FDR filtering (not used with the razor algorithm)")
	filterCmd.Flags().BoolVarP(&fp.Seq, "sequential", "", false, "alternative algorithm that estimates FDR using both filtered PSM and Protein lists")
	filterCmd.Flags().BoolVarP(&fp.Model, "models", "", false, "print model distribution")
	filterCmd.Flags().BoolVarP(&fp.Razor, "razor", "", false, "use razor peptides for protein FDR scoring")
	filterCmd.Flags().BoolVarP(&fp.Picked, "picked", "", false, "apply the picked FDR algorithm before the protein scoring")
	filterCmd.Flags().BoolVarP(&fp.Mapmods, "mapmods", "", false, "map modifications aquired by an open search")

	RootCmd.AddCommand(filterCmd)
}
