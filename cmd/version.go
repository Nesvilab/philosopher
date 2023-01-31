// Package cmd Version top level command
package cmd

import (
	"github.com/Nesvilab/philosopher/lib/gth"

	"github.com/sirupsen/logrus"
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
	},
}

func init() {

	RootCmd.AddCommand(versionCmd)
}
