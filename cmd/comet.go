package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher/lib/ext/comet"
	"github.com/prvst/philosopher/lib/sys"

	"github.com/spf13/cobra"
)

// cometCmd represents the comet command
var cometCmd = &cobra.Command{
	Use:   "comet",
	Short: "MS/MS database search",
	//Long:  "Peptide Spectrum Matching using the Comet algorithm\nComet release 2016.01.rev 2",
	Run: func(cmd *cobra.Command, args []string) {

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
			sys.CopyFile(cmt.DefaultParam, filepath.Base(cmt.DefaultParam))
			return
		}

		paramAbs, _ := filepath.Abs(m.Comet.Param)

		var binFile []byte
		binFile, err := ioutil.ReadFile(paramAbs)
		if err != nil {
			logrus.Fatal(err)
		}

		m.Comet.ParamFile = binFile

		// run comet
		e := cmt.Run(args, m.Comet.Param)
		if e != nil {
			fmt.Println(e.Error())
		}

		// store paramters on meta data
		m.Serialize()

		logrus.Info("Done")
		return
	},
}

func init() {

	if os.Args[1] == "comet" {

		m.Restore(sys.Meta())

		cometCmd.Flags().BoolVarP(&m.Comet.Print, "print", "", false, "print a comet.params file")
		cometCmd.Flags().StringVarP(&m.Comet.Param, "param", "", "comet.params.txt", "comet parameter file")

		RootCmd.AddCommand(cometCmd)
	}

}
