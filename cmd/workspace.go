// Package cmd Workspace top level command
package cmd

import (
	"github.com/Nesvilab/philosopher/lib/msg"
	"github.com/Nesvilab/philosopher/lib/wrk"

	"github.com/spf13/cobra"
)

var backup, clean, initialize, nocheck bool
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
	workspaceCmd.Flags().BoolVarP(&nocheck, "nocheck", "", false, "do not check for new versions")
	workspaceCmd.Flags().StringVarP(&temp, "temp", "", "", "define a custom temporary folder for Philosopher")
	RootCmd.AddCommand(workspaceCmd)
}
