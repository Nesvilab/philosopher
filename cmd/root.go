package cmd

import (
	"fmt"
	"runtime"

	"github.com/Sirupsen/logrus"
	colorable "github.com/mattn/go-colorable"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// constants
const (
	Version = "1.0"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "philosopher",
	Short: "philospher: a proteomics data analysis pipeline",
	Long:  "Philosopher: A tool for Proteomics data analysis and post-processing filtering" + "\nversion: " + Version,
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
		//logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
		logrus.SetOutput(colorable.NewColorableStdout())
	}

	logrus.SetFormatter(fmt)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".philosopher")
	viper.AddConfigPath("$HOME")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
