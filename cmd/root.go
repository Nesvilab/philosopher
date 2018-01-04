package cmd

import (
	"runtime"

	"github.com/Sirupsen/logrus"
	colorable "github.com/mattn/go-colorable"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/spf13/cobra"
)

var m met.Data

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "philosopher",
	Short: "philospher: a proteomics data analysis pipeline",
	Long:  "Philosopher: A tool for Proteomics data analysis and post-processing filtering",
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		logrus.Fatal(err)
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
