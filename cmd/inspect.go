package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/rep"
	"github.com/spf13/cobra"
	"github.com/vmihailenco/msgpack"
)

var object string
var key string

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:    "inspect",
	Short:  "Inspect meta data",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {

		//file, _ := os.Open(os.Args[2])

		if object == "meta" {

			var o met.Data

			target := fmt.Sprintf(".meta%smeta.bin", string(filepath.Separator))
			file, _ := os.Open(target)

			dec := msgpack.NewDecoder(file)
			err := dec.Decode(&o)
			if err != nil {
				logrus.Fatal("Could not restore meta data:", err)
			}
			spew.Dump(o)

		} else if object == "psm" {

			target := fmt.Sprintf(".meta%sev.psm.bin", string(filepath.Separator))
			file, _ := os.Open(target)

			var o rep.PSMEvidenceList
			dec := msgpack.NewDecoder(file)
			err := dec.Decode(&o)
			if err != nil {
				logrus.Fatal("Could not restore meta data:", err)
			}
			spew.Dump(o)

		} else if object == "protein" {

			target := fmt.Sprintf(".meta%sev.pro.bin", string(filepath.Separator))
			file, _ := os.Open(target)

			var o rep.ProteinEvidenceList
			dec := msgpack.NewDecoder(file)
			err := dec.Decode(&o)
			if err != nil {
				logrus.Fatal("Could not restore meta data:", err)
			}

			if len(key) > 0 {

				for _, i := range o {
					if i.ProteinID == key {
						spew.Dump(i)
					}
				}

			} else {
				spew.Dump(o)
			}

		}

		return
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "inspect" {
		inspectCmd.Flags().StringVarP(&object, "object", "", "meta", "object to inspect")
		inspectCmd.Flags().StringVarP(&key, "key", "", "", "individual ID to inspect")

		RootCmd.AddCommand(inspectCmd)
	}

}
