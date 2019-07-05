package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/pip"
	"github.com/prvst/philosopher/lib/sla"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

// pipelineCmd represents the pipeline command
var pipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Automatic execution of consecutive analysis steps",
	Run: func(cmd *cobra.Command, args []string) {

		logrus.Info("Initializing the Pipeline ", Version)

		// get current directory
		dir, e := os.Getwd()
		if e != nil {
			logrus.Info("check folder permissions")
		}

		// create a virtual meta instance
		meta := met.New(dir)
		// if e != nil {
		// 	logrus.Fatal(e.Error())
		// }

		os.Mkdir(meta.Temp, sys.FilePermission())
		if _, e = os.Stat(meta.Temp); os.IsNotExist(e) {
			logrus.Info("Can't find temporary directory; check folder permissions")
		}

		param, e := pip.DeployParameterFile(meta.Temp)
		if e != nil {
			logrus.Fatal(e.Error())
		}

		if m.Pipeline.Print == true {
			logrus.Info("Printing parameter file")
			sys.CopyFile(param, filepath.Base(param))
			return
		}

		file, _ := filepath.Abs(m.Pipeline.Directives)

		y, e := ioutil.ReadFile(file)
		if e != nil {
			log.Fatal(e)
		}

		var p pip.Directives
		e = yaml.Unmarshal(y, &p)
		if e != nil {
			logrus.Fatal(e)
		}

		if len(args) < 1 {
			logrus.Fatal("You need to provide at least one dataset for the analysis.")
		} else if p.Commands.Abacus == "true" && len(args) < 2 {
			logrus.Fatal("You need to provide at least two datasets for the abacus integrative analysis.")
		}

		// Workspace - Database
		meta = pip.InitializeWorkspaces(meta, p, dir, Version, Build, args)

		// Comet - MSFragger
		meta = pip.DatabaseSearch(meta, p, dir, args)

		// PeptideProphet - PTMProphet - ProteinProphet
		meta = pip.Prophets(meta, p, dir, args)

		// Abacus - combined protxml
		meta = pip.CombinedProteinList(meta, p, dir, args)

		// Filter - Quantification - Clustering - Report
		meta = pip.FilterQuantifyReport(meta, p, dir, args)

		// Abacus
		meta = pip.Abacus(meta, p, dir, args)

		// TMT-Integrator
		meta = pip.TMTIntegrator(meta, p, dir, args)

		// Backup and Clean
		//pip.BackupAndClean(meta, p, dir, Version, Build, args)

		if len(p.SlackToken) > 0 {
			sla.Run("Philosopher", p.SlackToken, "Philosopher pipeline is done", p.SlackChannel)
		}

		logrus.Info("Done")
		return
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "pipeline" {

		m.Restore(sys.Meta())

		pipelineCmd.Flags().BoolVarP(&m.Pipeline.Print, "print", "", false, "print the pipeline configuration file")
		pipelineCmd.Flags().StringVarP(&m.Pipeline.Directives, "config", "", "", "configuration file for the pipeline execution")

	}

	RootCmd.AddCommand(pipelineCmd)
}
