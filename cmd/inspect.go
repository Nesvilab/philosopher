package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/prvst/philosopher/lib/msg"

	"github.com/davecgh/go-spew/spew"
	"github.com/prvst/philosopher/lib/dat"
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

		if object == "meta" {

			var o met.Data

			target := fmt.Sprintf(".meta%smeta.bin", string(filepath.Separator))
			file, _ := os.Open(target)

			dec := msgpack.NewDecoder(file)
			e := dec.Decode(&o)
			if e != nil {
				msg.DecodeMsgPck(e, "fatal")
			}

			if key == "session" {
				fmt.Println(o.UUID)
			} else {
				spew.Dump(o)
			}

		} else if object == "parameters" {

			var o rep.SearchParametersEvidence

			target := fmt.Sprintf(".meta%sev.param.bin", string(filepath.Separator))
			file, _ := os.Open(target)

			dec := msgpack.NewDecoder(file)
			e := dec.Decode(&o)
			if e != nil {
				msg.DecodeMsgPck(e, "fatal")
			}
			spew.Dump(o)

		} else if object == "psm" {

			var o rep.PSMEvidenceList

			target := fmt.Sprintf(".meta%sev.psm.bin", string(filepath.Separator))
			file, _ := os.Open(target)

			dec := msgpack.NewDecoder(file)
			e := dec.Decode(&o)
			if e != nil {
				msg.DecodeMsgPck(e, "fatal")
			}
			spew.Dump(o)

		} else if object == "db" {

			var o dat.Base

			target := fmt.Sprintf(".meta%sdb.bin", string(filepath.Separator))
			file, _ := os.Open(target)

			dec := msgpack.NewDecoder(file)
			e := dec.Decode(&o)
			if e != nil {
				msg.DecodeMsgPck(e, "fatal")
			}
			spew.Dump(o.Records)

		} else if object == "protein" {

			var o rep.ProteinEvidenceList

			target := fmt.Sprintf(".meta%spro.bin", string(filepath.Separator))
			file, _ := os.Open(target)

			dec := msgpack.NewDecoder(file)
			e := dec.Decode(&o)
			if e != nil {
				msg.DecodeMsgPck(e, "fatal")
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
