package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher/lib/wrk"
	"github.com/spf13/cobra"
)

var i, b, c bool

// workspaceCmd represents the workspace command
var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Manage the experiment workspace for the analysis",

	Run: func(cmd *cobra.Command, args []string) {

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
	RootCmd.AddCommand(workspaceCmd)
}
