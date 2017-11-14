package cmd

import (
	"os"
	"runtime"

	"github.com/Sirupsen/logrus"
	colorable "github.com/mattn/go-colorable"
	"github.com/prvst/cmsl/err"
	"github.com/prvst/philosopher/lib/meta"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/spf13/cobra"
)

// constants
const (
	Version = "2.0"
)

var m meta.Data

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "philosopher",
	Short: "philospher: a proteomics data analysis toolkit",
	Long:  "Philosopher: A toolkit for Proteomics data analysis and post-processing filtering" + "\nversion: " + Version,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	// verify if the command is been executed on a workspace directory
	if os.Args[1] != "workspace" {
		if len(m.UUID) < 1 && len(m.Home) < 1 {
			e := &err.Error{Type: err.WorkspaceNotFound, Class: err.FATA}
			logrus.Fatal(e.Error())
		}
	}

	if err := RootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

// init sets up the configuration for the Cobra package and the Logrus logger
func init() {

	cobra.OnInitialize(initConfig)

	fmt := new(logrus.TextFormatter)
	fmt.TimestampFormat = "01:02:03"
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
	// empty
}
