package cmd

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher/lib/met"
	"github.com/sanity-io/litter"
	"github.com/spf13/cobra"
	"github.com/vmihailenco/msgpack"
)

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:    "inspect",
	Short:  "Inspect meta data",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {

		var d met.Data

		file, _ := os.Open(os.Args[2])

		dec := msgpack.NewDecoder(file)
		err := dec.Decode(&d)
		if err != nil {
			logrus.Fatal("Could not restore meta data:", err)
		}

		litter.Dump(d)

		return
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "inspect" {
		RootCmd.AddCommand(inspectCmd)
	}

}
