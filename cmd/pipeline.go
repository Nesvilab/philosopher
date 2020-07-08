// Package cmd Pipeline top level command
package cmd

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"philosopher/lib/met"
	"philosopher/lib/msg"
	"philosopher/lib/pip"
	"philosopher/lib/sla"
	"philosopher/lib/sys"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

// pipelineCmd represents the pipeline command
var pipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Automatic execution of consecutive analysis steps",
	Run: func(cmd *cobra.Command, args []string) {

		msg.Executing("Pipeline ", Version)

		// get current directory
		dir, e := os.Getwd()
		if e != nil {
			logrus.Info("check folder permissions")
		}

		// create a virtual meta instance
		meta := met.New(dir)

		os.Mkdir(meta.Temp, sys.FilePermission())
		if _, e = os.Stat(meta.Temp); os.IsNotExist(e) {
			msg.Custom(errors.New("Can't find temporary directory; check folder permissions"), "info")
		}

		if m.Pipeline.Print == true {
			param := pip.DeployParameterFile(meta.Temp)
			msg.Custom(errors.New("Printing parameter file"), "info")
			sys.CopyFile(param, filepath.Base(param))
			return
		}

		file, _ := filepath.Abs(m.Pipeline.Directives)

		y, e := ioutil.ReadFile(file)
		if e != nil {
			msg.ReadFile(e, "fatal")
		}

		var p pip.Directives
		e = yaml.Unmarshal(y, &p)
		if e != nil {
			msg.ReadFile(e, "fatal")
		}

		if len(args) < 1 {
			msg.NoParametersFound(errors.New("You need to provide at least one dataset for the analysis"), "fatal")
		} else if p.Commands.Abacus == "true" && len(args) < 2 {
			msg.NoParametersFound(errors.New("You need to provide at least two datasets for the abacus integrative analysis"), "fatal")
		}

		// Workspace - Database
		if p.Commands.Workspace == "yes" {
			meta = pip.InitializeWorkspaces(meta, p, dir, Version, Build, args)
		}

		// Comet - MSFragger
		if p.Commands.Comet == "yes" && p.Commands.MSFragger == "yes" {
			msg.Custom(errors.New("You can only specify one search engine at a time"), "fatal")
		} else if p.Commands.Comet == "yes" || p.Commands.MSFragger == "yes" {
			meta = pip.DatabaseSearch(meta, p, dir, args)
		}

		// PeptideProphet
		if p.Commands.PeptideProphet == "yes" {
			meta = pip.PeptideProphet(meta, p, dir, args)
		}

		// PTMProphet
		if p.Commands.PTMProphet == "yes" {
			meta = pip.PTMProphet(meta, p, dir, args)
		}

		// ProteinProphet
		if p.Commands.ProteinProphet == "yes" {
			meta = pip.ProteinProphet(meta, p, dir, args)
		}

		// Abacus - combined pepxml
		meta = pip.CombinedPeptideList(meta, p, dir, args)

		// Abacus - combined protxml
		meta = pip.CombinedProteinList(meta, p, dir, args)

		// FreeQuant
		if p.Commands.FreeQuant == "yes" {
			//if _, err := os.Stat(sys.LFQBin()); os.IsNotExist(err) {
			meta = pip.FreeQuant(meta, p, dir, args)
			//}
		}

		// LabelQuant
		if p.Commands.LabelQuant == "yes" {
			//if _, err := os.Stat(sys.IsoBin()); os.IsNotExist(err) {
			meta = pip.LabelQuant(meta, p, dir, args)
			//}
		}

		// Filter - Report
		meta = pip.FilterAndReport(meta, p, dir, args)

		// BioQuant
		if p.Commands.BioQuant == "yes" {
			meta = pip.BioQuant(meta, p, dir, args)
		}

		// Abacus
		meta = pip.Abacus(meta, p, dir, args)

		// TMT-Integrator
		meta = pip.TMTIntegrator(meta, p, dir, args)

		// Backup and Clean
		//pip.BackupAndClean(meta, p, dir, Version, Build, args)

		if len(p.SlackToken) > 0 {
			sla.Run("Philosopher", p.SlackToken, "Philosopher pipeline is done", p.SlackChannel, p.SlackUserID)
		}

		met.CleanTemp(meta.Temp)

		msg.Done()
		return
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "pipeline" {

		//m.Restore(sys.Meta())

		pipelineCmd.Flags().BoolVarP(&m.Pipeline.Print, "print", "", false, "print the pipeline configuration file")
		pipelineCmd.Flags().StringVarP(&m.Pipeline.Directives, "config", "", "", "configuration file for the pipeline execution")

	}

	RootCmd.AddCommand(pipelineCmd)
}
