// Package cmd Inspect top level command
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"philosopher/lib/dat"
	"philosopher/lib/fil"
	"philosopher/lib/met"
	"philosopher/lib/msg"
	"philosopher/lib/qua"
	"philosopher/lib/rep"
	"philosopher/lib/sys"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
	"github.com/vmihailenco/msgpack/v5"
)

var object string
var key string

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:    "inspect",
	Short:  "Inspect meta data",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {

		switch object {
		case "meta":
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
		case "psm":
			var o rep.PSMEvidenceList

			target := fmt.Sprintf(".meta%spsm.bin", string(filepath.Separator))

			sys.Restore(&o, target, false)

			spew.Dump(o)

		case "peptide":
			var o rep.PeptideEvidenceList

			target := fmt.Sprintf(".meta%spep.bin", string(filepath.Separator))
			file, _ := os.Open(target)

			dec := msgpack.NewDecoder(file)
			e := dec.Decode(&o)
			if e != nil {
				msg.DecodeMsgPck(e, "fatal")
			}
			spew.Dump(o)
		case "db":
			var o dat.Base

			target := fmt.Sprintf(".meta%sdb.bin", string(filepath.Separator))
			file, _ := os.Open(target)

			dec := msgpack.NewDecoder(file)
			e := dec.Decode(&o)
			if e != nil {
				msg.DecodeMsgPck(e, "fatal")
			}
			spew.Dump(o.Records)
		case "lfq":
			var o qua.LFQ

			target := fmt.Sprintf(".meta%slfq.bin", string(filepath.Separator))
			file, _ := os.Open(target)

			dec := msgpack.NewDecoder(file)
			e := dec.Decode(&o)
			if e != nil {
				msg.DecodeMsgPck(e, "fatal")
			}
			spew.Dump(o.Intensities)
		case "razor":
			var o fil.RazorMap

			target := fmt.Sprintf(".meta%srazor.bin", string(filepath.Separator))
			file, _ := os.Open(target)

			dec := msgpack.NewDecoder(file)
			e := dec.Decode(&o)
			if e != nil {
				msg.DecodeMsgPck(e, "fatal")
			}
			spew.Dump(o)
		case "protein":
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
		default:
			msg.Custom(errors.New("the option is not available"), "fatal")
		}

	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "inspect" {
		inspectCmd.Flags().StringVarP(&object, "object", "", "meta", "object to inspect")
		inspectCmd.Flags().StringVarP(&key, "key", "", "", "individual ID to inspect")

		RootCmd.AddCommand(inspectCmd)
	}

}
