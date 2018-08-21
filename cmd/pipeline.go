package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/pip"
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
		var meta = met.New(dir)
		os.Mkdir(meta.Temp, 0755)
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

		if m.Pipeline.Parallel == true {
			pip.ParallelRun(m, p, dir, Version, Build, args)
		} else {
			m = pip.Run(m, p, dir, Version, Build, args)
		}

		// // For each dataset ...
		// for _, i := range args {
		//
		// 	logrus.Info("Executing the pipeline on ", i)
		//
		// 	// getting inside de the dataset folder
		// 	dsAbs, _ := filepath.Abs(i)
		// 	os.Chdir(dsAbs)
		//
		// 	// Workspace
		// 	wrk.Run(Version, Build, false, false, true)
		//
		// 	// reload the meta data
		// 	m.Restore(sys.Meta())
		//
		// 	// Database
		// 	if p.Commands.Database == "yes" {
		// 		m.Database = p.Database
		// 		dat.Run(m)
		//
		// 		m.Serialize()
		// 		met.CleanTemp(m.Temp)
		// 	}
		//
		// 	if p.Commands.Comet == "yes" && p.Commands.MSFragger == "yes" {
		// 		logrus.Fatal("You can only specify one search engine at a time")
		// 	}
		//
		// 	// Comet
		// 	if p.Commands.Comet == "yes" {
		// 		m.Comet = p.Comet
		// 		gobExt := fmt.Sprintf("*.%s", p.Comet.RawExtension)
		// 		files, e := filepath.Glob(gobExt)
		// 		if e != nil {
		// 			logrus.Fatal(e)
		// 		}
		// 		comet.Run(m, files)
		//
		// 		m.Serialize()
		// 		met.CleanTemp(m.Temp)
		// 	}
		//
		// 	// MSFragger
		// 	if p.Commands.MSFragger == "yes" {
		// 		m.MSFragger = p.MSFragger
		// 		gobExt := fmt.Sprintf("*.%s", p.MSFragger.RawExtension)
		// 		files, e := filepath.Glob(gobExt)
		// 		if e != nil {
		// 			logrus.Fatal(e)
		// 		}
		// 		fragger.Run(m, files)
		//
		// 		m.Serialize()
		// 		met.CleanTemp(m.Temp)
		// 	}
		//
		// 	// PeptideProphet
		// 	if p.Commands.PeptideProphet == "yes" {
		// 		logrus.Info("Executing PeptideProphet on ", i)
		// 		m.PeptideProphet = p.PeptideProphet
		// 		m.PeptideProphet.Output = "interact"
		// 		m.PeptideProphet.Combine = true
		// 		gobExt := fmt.Sprintf("*.%s", p.PeptideProphet.FileExtension)
		// 		files, e := filepath.Glob(gobExt)
		// 		if e != nil {
		// 			logrus.Fatal(e)
		// 		}
		// 		peptideprophet.Run(m, files)
		//
		// 		m.Serialize()
		// 		met.CleanTemp(m.Temp)
		// 	}
		//
		// 	if p.Commands.PTMProphet == "yes" {
		// 		logrus.Info("Executing PTMProphet on ", i)
		// 		m.PTMProphet = p.PTMProphet
		// 		var files []string
		// 		files = append(files, "interact.pep.xml")
		// 		m.PTMProphet.InputFiles = files
		// 		ptmprophet.Run(m, files)
		//
		// 		m.Serialize()
		// 		met.CleanTemp(m.Temp)
		// 	}
		//
		// 	// ProteinProphet
		// 	if p.Commands.ProteinProphet == "yes" {
		// 		logrus.Info("Executing ProteinProphet on ", i)
		// 		m.ProteinProphet = p.ProteinProphet
		// 		m.ProteinProphet.Output = "interact"
		// 		var files []string
		// 		if p.Commands.PTMProphet == "yes" {
		// 			files = append(files, "interact.mod.pep.xml")
		// 		} else {
		// 			files = append(files, "interact.pep.xml")
		// 		}
		// 		proteinprophet.Run(m, files)
		//
		// 		m.Serialize()
		// 		met.CleanTemp(m.Temp)
		// 	}
		//
		// 	// return to the top level directory
		// 	os.Chdir(dir)
		// }
		//
		// // Abacus
		// var combinedProtXML string
		// if p.Commands.Abacus == "yes" {
		// 	logrus.Info("Creating combined protein inference")
		// 	// return to the top level directory
		// 	os.Chdir(dir)
		// 	m.ProteinProphet = p.ProteinProphet
		// 	m.ProteinProphet.Output = "combined"
		// 	var files []string
		// 	for _, j := range args {
		// 		fqn := fmt.Sprintf("%s%sinteract.pep.xml", j, string(filepath.Separator))
		// 		if p.Commands.PTMProphet == "yes" {
		// 			fqn = fmt.Sprintf("%s%sinteract.mod.pep.xml", j, string(filepath.Separator))
		// 		}
		// 		fqn, _ = filepath.Abs(fqn)
		// 		files = append(files, fqn)
		// 	}
		// 	proteinprophet.Run(m, files)
		// 	combinedProtXML = fmt.Sprintf("%s%scombined.prot.xml", m.Temp, string(filepath.Separator))
		//
		// 	// copy to work directory
		// 	sys.CopyFile(combinedProtXML, filepath.Base(combinedProtXML))
		//
		// 	m.Serialize()
		// 	met.CleanTemp(m.Temp)
		// }
		//
		// for _, i := range args {
		//
		// 	// getting inside  each dataset folder again
		// 	dsAbs, _ := filepath.Abs(i)
		// 	os.Chdir(dsAbs)
		//
		// 	// Filter
		// 	if p.Commands.Filter == "yes" {
		// 		logrus.Info("Executing filter on ", i)
		// 		m.Filter = p.Filter
		// 		m.Filter.Pex = "interact.pep.xml"
		// 		if p.Commands.PTMProphet == "yes" {
		// 			m.Filter.Pex = "interact.mod.pep.xml"
		// 		}
		// 		if p.Commands.ProteinProphet == "yes" {
		// 			m.Filter.Pox = "interact.prot.xml"
		// 		}
		// 		if p.Commands.Abacus == "yes" {
		// 			m.Filter.Pox = combinedProtXML
		// 		}
		// 		m, e := fil.Run(m)
		// 		if e != nil {
		// 			logrus.Fatal(e.Error())
		// 		}
		//
		// 		m.Serialize()
		// 		met.CleanTemp(m.Temp)
		// 	}
		//
		// 	// getting inside de the dataset folder again
		// 	os.Chdir(dsAbs)
		//
		// 	// FreeQuant
		// 	if p.Commands.FreeQuant == "yes" {
		// 		logrus.Info("Executing label-free quantification on ", i)
		// 		m.Quantify = p.Freequant
		// 		m.Quantify.Dir = dsAbs
		// 		m.Quantify.Format = "mzML"
		// 		e := qua.RunLabelFreeQuantification(m.Quantify)
		// 		if e != nil {
		// 			logrus.Fatal(e.Error())
		// 		}
		//
		// 		m.Serialize()
		// 		met.CleanTemp(m.Temp)
		// 	}
		//
		// 	// LabelQuant
		// 	if p.Commands.LabelQuant == "yes" {
		// 		logrus.Info("Executing label-based quantification on ", i)
		// 		m.Quantify = p.LabelQuant
		// 		m.Quantify.Dir = dsAbs
		// 		m.Quantify.Format = "mzML"
		// 		m.Quantify.Brand = "tmt"
		// 		var e error
		// 		m.Quantify, e = qua.RunTMTQuantification(m.Quantify, m.Filter.Mapmods)
		// 		if e != nil {
		// 			logrus.Fatal(e)
		// 		}
		//
		// 		m.Serialize()
		// 		met.CleanTemp(m.Temp)
		// 	}
		//
		// 	// Report
		// 	if p.Commands.Report == "yes" {
		// 		logrus.Info("Executing report on ", i)
		// 		rep.Run(m)
		//
		// 		m.Serialize()
		// 		met.CleanTemp(m.Temp)
		// 	}
		//
		// 	// Cluster
		// 	if p.Commands.Cluster == "yes" {
		// 		logrus.Info("Executing cluster on ", i)
		// 		m.Cluster = p.Cluster
		// 		clu.GenerateReport(m)
		//
		// 		m.Serialize()
		// 		met.CleanTemp(m.Temp)
		// 	}
		//
		// 	// return to the top level directory
		// 	os.Chdir(dir)
		// }
		//
		// // Abacus
		// if p.Commands.Abacus == "yes" {
		// 	logrus.Info("Executing abacus")
		// 	// return to the top level directory
		// 	os.Chdir(dir)
		// 	m.Abacus = p.Abacus
		// 	err := aba.Run(m.Abacus, m.Temp, args)
		// 	if err != nil {
		// 		logrus.Fatal(err)
		// 	}
		// }
		//
		// for _, i := range args {
		//
		// 	// getting inside de the dataset folder
		// 	localDir, _ := filepath.Abs(i)
		// 	os.Chdir(localDir)
		//
		// 	// Backup
		// 	if p.Backup == true {
		// 		wrk.Run(Version, Build, true, false, false)
		// 	}
		//
		// 	// Clean
		// 	if p.Clean == true {
		// 		wrk.Run(Version, Build, false, true, false)
		// 	}
		//
		// }

		// if len(p.SlackToken) > 0 {
		// 	sla.Run("Philosopher", p.SlackToken, "Philosopher pipeline is done", p.SlackChannel)
		// }

		logrus.Info("Done")
		return
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "pipeline" {

		m.Restore(sys.Meta())

		pipelineCmd.Flags().BoolVarP(&m.Pipeline.Print, "print", "", false, "print the pipeline configuration file")
		pipelineCmd.Flags().BoolVarP(&m.Pipeline.Parallel, "parallel", "", false, "uses a palelle mode for processing")
		pipelineCmd.Flags().StringVarP(&m.Pipeline.Directives, "config", "", "", "configuration file for the pipeline execution")

	}

	RootCmd.AddCommand(pipelineCmd)
}
