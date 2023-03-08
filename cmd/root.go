// Package cmd provides the top level methods that correspond to the available commands
package cmd

import (
	"runtime"

	"github.com/Nesvilab/philosopher/lib/met"
	"github.com/Nesvilab/philosopher/lib/msg"
	"github.com/Nesvilab/philosopher/lib/sys"

	colorable "github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var m met.Data

var (
	// Version code
	Version string
	// Build code
	Build string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "philosopher",
	Short: "philosopher: a proteomics data analysis toolkit",
	Long:  "Philosopher: A toolkit for Proteomics data analysis and post-processing filtering",
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if e := RootCmd.Execute(); e != nil {
		msg.Custom(e, "trace")
	}
}

func init() {

	cobra.OnInitialize(initConfig)
	fmt := new(logrus.TextFormatter)
	fmt.TimestampFormat = "15:04:05"
	fmt.FullTimestamp = true
	fmt.DisableColors = false

	if runtime.GOOS == sys.Windows() {
		fmt.ForceColors = true
		logrus.SetOutput(colorable.NewColorableStdout())
	}

	logrus.SetFormatter(fmt)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
