package cmd

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/jpillora/go-ogle-analytics"
	"github.com/prvst/philosopher/lib/gth"
	"github.com/prvst/philosopher/lib/wrk"
	"github.com/spf13/cobra"
)

var a, b, c, i bool

// workspaceCmd represents the workspace command
var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Manage the experiment workspace for the analysis",
	Run: func(cmd *cobra.Command, args []string) {

		gth.UpdateChecker()

		if (i == true && b == true && c == true) || (i == true && b == true) || (i == true && c == true) || (c == true && b == true) {
			logrus.Fatal("this command accepts only one parameter")
		}

		if i == true {

			logrus.Info("Creating workspace")
			e := wrk.Init()
			if e != nil {
				if e.Class == "warning" {
					logrus.Warn(e.Error())
				}
			}
			logrus.Info("Done")
			return

		} else if b == true {

			logrus.Info("Creating backup")
			e := wrk.Backup()
			if e != nil {
				logrus.Warn(e.Error())
			}
			logrus.Info("Done")
			return

		} else if c == true {

			logrus.Info("Removing workspace")
			e := wrk.Clean()
			if e != nil {
				logrus.Warn(e.Error())
			}
			logrus.Info("Done")
			return

		}

		return
	},
}

func init() {

	workspaceCmd.Flags().BoolVarP(&i, "init", "", false, "Initialize the workspace")
	workspaceCmd.Flags().BoolVarP(&b, "backup", "", false, "create a backup of the experiment meta data")
	workspaceCmd.Flags().BoolVarP(&c, "clean", "", false, "Remove the workspace and all meta data. Experimental file are kept intact")
	workspaceCmd.Flags().BoolVarP(&a, "analytics", "", true, "Reports when a workspace is created for usage estimation")

	if len(os.Args) > 1 && os.Args[1] == "workspace" && a == true {

		// do not change this! This is for metric colletion, no user data is gatter, the software just reports back
		// the number of people using it and the geo location, just like any other website does.
		client, err := ga.NewClient("UA-111428141-1")
		if err != nil {
			panic(err)
		}

		err = client.Send(ga.NewEvent("Philosopher", "Workspace"))
		if err != nil {
			panic(err)
		}

	}

	RootCmd.AddCommand(workspaceCmd)
}
