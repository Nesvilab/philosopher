// Package cmd Workspace top level command
package cmd

import (
	"fmt"
	"os"

	"philosopher/lib/msg"
	"philosopher/lib/wrk"

	ga "github.com/jpillora/go-ogle-analytics"
	"github.com/spf13/cobra"
)

var a, b, c, i, n bool

// workspaceCmd represents the workspace command
var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Manage the experiment workspace for the analysis",
	Run: func(cmd *cobra.Command, args []string) {

		msg.Executing("Workspace ", Version)

		wrk.Run(Version, Build, b, c, i, n)

		msg.Done()
		return
	},
}

func init() {

	workspaceCmd.Flags().BoolVarP(&i, "init", "", false, "initialize the workspace")
	workspaceCmd.Flags().BoolVarP(&b, "backup", "", false, "create a backup of the experiment meta data")
	workspaceCmd.Flags().BoolVarP(&c, "clean", "", false, "remove the workspace and all meta data. Experimental file are kept intact")
	workspaceCmd.Flags().BoolVarP(&a, "analytics", "", true, "reports when a workspace is created for usage estimation")
	workspaceCmd.Flags().BoolVarP(&n, "nocheck", "", false, "do not check for new versions")

	if len(os.Args) > 1 && os.Args[1] == "workspace" && a == true {

		// do not change this! This is for metric colletion, no user data is gatter, the software just reports back
		// the number of people using it and the geo location, just like any other website does.
		client, err := ga.NewClient("UA-111428141-1")
		if err != nil {
			_ = err
		}

		v := fmt.Sprintf("Version:%s", Version)
		b := fmt.Sprintf("Build:%s", Build)
		err = client.Send(ga.NewEvent("Philosopher", v).Label(b))
		if err != nil {
			_ = err
		}

	}

	RootCmd.AddCommand(workspaceCmd)
}
