package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/prvst/philosopher/lib/gth"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Philosopher version",
	Run: func(cmd *cobra.Command, args []string) {

		logrus.WithFields(logrus.Fields{
			"version": Version,
			"build":   Build,
		}).Info("Current Philosopher build and version")

		gth.UpdateChecker(Version, Build)

		return
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "version" {
	}

	RootCmd.AddCommand(versionCmd)
}
