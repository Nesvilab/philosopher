// Package cmd Pipeline top level command
package cmd

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Nesvilab/philosopher/lib/met"
	"github.com/Nesvilab/philosopher/lib/msg"
	"github.com/Nesvilab/philosopher/lib/pip"
	"github.com/Nesvilab/philosopher/lib/sla"
	"github.com/Nesvilab/philosopher/lib/sys"

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
			msg.Custom(errors.New("can't find temporary directory; check folder permissions"), "info")
		}

		if m.Pipeline.Print {
			param := pip.DeployParameterFile(meta.Temp)
			msg.Custom(errors.New("printing parameter file"), "info")
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
			msg.NoParametersFound(errors.New("you need to provide at least one dataset for the analysis"), "fatal")
		} else if p.Steps.IntegratedReports == "true" && len(args) < 2 {
			msg.NoParametersFound(errors.New("you need to provide at least two datasets for the abacus integrative analysis"), "fatal")
		}

		// Workspace - Database
		meta = pip.InitializeWorkspaces(meta, p, dir, Version, Build, args)

		meta = pip.AnnotateDatabase(meta, p, dir, args)

		// if m.Pipeline.Verbose == true {
		// 	meta.Pipeline.Verbose = true
		// }

		// Comet - MSFragger
		if p.Steps.DatabaseSearch == "yes" {
			meta = pip.DBSearch(meta, p, dir, args)
		}

		// PeptideProphet
		if p.Steps.PeptideValidation == "yes" {
			meta = pip.PeptideProphet(meta, p, dir, args)
		}

		// PTMProphet
		if p.Steps.PTMLocalization == "yes" {
			meta = pip.PTMProphet(meta, p, dir, args)
		}

		// ProteinProphet
		if p.Steps.ProteinInference == "yes" {
			meta = pip.ProteinProphet(meta, p, dir, args)
		}

		if p.Steps.IntegratedReports == "yes" {
			// Abacus - combined pepxml
			meta = pip.CombinedPeptideList(meta, p, dir, args)

			// Abacus - combined protxml
			meta = pip.CombinedProteinList(meta, p, dir, args)
		}

		// Filter
		if p.Steps.FDRFiltering == "yes" {
			meta = pip.Filter(meta, p, dir, args)
		}

		// FreeQuant
		if p.Steps.LabelFreeQuantification == "yes" {
			meta = pip.FreeQuant(meta, p, dir, args)
		}

		// LabelQuant
		if p.Steps.IsobaricQuantification == "yes" {
			meta = pip.LabelQuant(meta, p, dir, args)
		}

		// Report
		if p.Steps.IndividualReports == "yes" {
			meta = pip.Report(meta, p, dir, args)
		}

		// BioQuant
		if p.Steps.BioClusterQuantification == "yes" {
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
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "pipeline" {

		pipelineCmd.Flags().BoolVarP(&m.Pipeline.Print, "print", "", false, "print the pipeline configuration file")
		pipelineCmd.Flags().BoolVarP(&m.Pipeline.Verbose, "verbose", "", false, "show the parameters for each command that is executed")
		pipelineCmd.Flags().StringVarP(&m.Pipeline.Directives, "config", "", "", "configuration file for the pipeline execution")
		pipelineCmd.Flags().MarkHidden("verbose")

	}

	RootCmd.AddCommand(pipelineCmd)
}
