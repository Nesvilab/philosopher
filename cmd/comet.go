package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/ext/comet"
	"github.com/prvst/philosopher/lib/raw"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/spf13/cobra"
)

// cometCmd represents the comet command
var cometCmd = &cobra.Command{
	Use:   "comet",
	Short: "Peptide spectrum matching with Comet",
	Run: func(cmd *cobra.Command, args []string) {

		// verify if the command is been executed on a workspace directory
		if len(m.UUID) < 1 && len(m.Home) < 1 {
			e := &err.Error{Type: err.WorkspaceNotFound, Class: err.FATA}
			logrus.Fatal(e.Error())
		}

		logrus.Info("Executing Comet")
		var cmt = comet.New()

		if len(m.Comet.Param) < 1 {
			logrus.Fatal("No parameter file found. Run 'comet --help' for more information")
		}

		if m.Comet.Print == false && len(args) < 1 {
			logrus.Fatal("Missing parameter file or data file for analysis")
		}

		// deploy the binaries
		cmt.Deploy(m.OS, m.Arch)

		if m.Comet.Print == true {
			logrus.Info("Printing parameter file")
			sys.CopyFile(cmt.DefaultParam, filepath.Base(cmt.DefaultParam))
			return
		}

		// collect and store the mz files
		m.Comet.RawFiles = args

		// convert the param file to binary and store it in meta
		var binFile []byte
		paramAbs, _ := filepath.Abs(m.Comet.Param)
		binFile, e := ioutil.ReadFile(paramAbs)
		if e != nil {
			logrus.Fatal(e)
		}
		m.Comet.ParamFile = binFile

		// the indexing will help later in case other commands are used for qunatification
		// it will provide easy and fast access to mz data
		logrus.Info("Indexing spectra: please wait, this can take a few minutes")
		raw.IndexMz(args)

		// run comet
		// e = cmt.Run(args, m.Comet.Param)
		// if e != nil {
		// 	//logrus.Fatal(e)
		// }

		// store paramters on meta data
		m.Serialize()

		logrus.Info("Done")
		return
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "comet" {

		m.Restore(sys.Meta())

		cometCmd.Flags().BoolVarP(&m.Comet.Print, "print", "", false, "print a comet.params file")
		cometCmd.Flags().StringVarP(&m.Comet.Param, "param", "", "comet.params.txt", "comet parameter file")
	}

	RootCmd.AddCommand(cometCmd)
}
