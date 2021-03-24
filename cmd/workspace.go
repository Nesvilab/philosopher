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

var analytics, backup, clean, initialize, nocheck bool
var temp string

// workspaceCmd represents the workspace command
var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Manage the experiment workspace for the analysis",
	Run: func(cmd *cobra.Command, args []string) {

		msg.Executing("Workspace ", Version)

		wrk.Run(Version, Build, temp, backup, clean, initialize, nocheck)

		msg.Done()
	},
}

func init() {

	workspaceCmd.Flags().BoolVarP(&initialize, "init", "", false, "initialize the workspace")
	workspaceCmd.Flags().BoolVarP(&backup, "backup", "", false, "create a backup of the experiment meta data")
	workspaceCmd.Flags().BoolVarP(&clean, "clean", "", false, "remove the workspace and all meta data. Experimental file are kept intact")
	workspaceCmd.Flags().BoolVarP(&analytics, "analytics", "", true, "reports when a workspace is created for usage estimation")
	workspaceCmd.Flags().BoolVarP(&nocheck, "nocheck", "", false, "do not check for new versions")
	workspaceCmd.Flags().StringVarP(&temp, "temp", "", "", "define a custom temporary folder for Philosopher")

	if len(os.Args) > 1 && os.Args[1] == "workspace" && analytics {

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
