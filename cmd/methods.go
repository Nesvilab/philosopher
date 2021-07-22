package cmd

import (
	"os"
	"philosopher/lib/met"
	"philosopher/lib/msg"
	"philosopher/lib/sys"
	"philosopher/lib/wmm"

	"github.com/spf13/cobra"
)

// methodsCmd represents the methods command
var methodsCmd = &cobra.Command{
	Use:    "methods",
	Hidden: true,
	Short:  "A write-my-methods function",

	Run: func(cmd *cobra.Command, args []string) {

		m.FunctionInitCheckUp()

		msg.Executing("Methods ", Version)

		wmm.Run(m)

		//m.Serialize()

		// clean tmp
		met.CleanTemp(m.Temp)

		msg.Done()
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "methods" {

		m.Restore(sys.Meta())

	}

	RootCmd.AddCommand(methodsCmd)
}
