// +build !windows

package cmd

import (
	"os"

	"github.com/prvst/philosopher/lib/ext/idconvert"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// idconvertCmd represents the idconvert command
var idconvertCmd = &cobra.Command{
	Use:   "idconvert",
	Short: "Convert mass spec identification file formats",

	Run: func(cmd *cobra.Command, args []string) {

		logrus.Info("Executing Idconvert ", Version)

		m, e := idconvert.Run(m, args)
		if e != nil {
			logrus.Fatal(e.Error())
		}

		m.Serialize()

		logrus.Info("Done")
		return

	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "idconvert" {

		m.Restore(sys.Meta())

		idconvertCmd.Flags().StringVarP(&m.Idconvert.Format, "format", "", "", "pepXML, mzIdentML, text")

	}

	RootCmd.AddCommand(idconvertCmd)
}
