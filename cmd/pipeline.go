package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/prvst/philosopher/lib/aba"
	"github.com/prvst/philosopher/lib/clu"
	"github.com/prvst/philosopher/lib/dat"
	"github.com/prvst/philosopher/lib/ext/comet"
	"github.com/prvst/philosopher/lib/ext/msfragger"
	"github.com/prvst/philosopher/lib/ext/peptideprophet"
	"github.com/prvst/philosopher/lib/ext/proteinprophet"
	"github.com/prvst/philosopher/lib/ext/ptmprophet"
	"github.com/prvst/philosopher/lib/fil"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/pip"
	"github.com/prvst/philosopher/lib/qua"
	"github.com/prvst/philosopher/lib/rep"
	"github.com/prvst/philosopher/lib/sla"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/prvst/philosopher/lib/wrk"
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

		// For each dataset initialize workspace
		for _, i := range args {

			logrus.Info("Initiating the pipeline on ", i)

			// getting inside de the dataset folder
			dsAbs, _ := filepath.Abs(i)
			os.Chdir(dsAbs)

			// Workspace
			wrk.Run(Version, Build, false, false, true, false)

			// reload the meta data
			m.Restore(sys.Meta())

			// Database
			if p.Commands.Database == "yes" {
				m.Database = p.Database
				dat.Run(m)
				m.Serialize()
			}

			if p.Commands.Comet == "yes" && p.Commands.MSFragger == "yes" {
				logrus.Fatal("You can only specify one search engine at a time")
			}

			// return to the top level directory
			os.Chdir(dir)
		}

		// run database search on all files from top folder
		if p.Commands.Comet == "yes" || p.Commands.MSFragger == "yes" {
			var mzFiles []string

			for _, i := range args {

				// getting inside de the dataset folder
				dsAbs, _ := filepath.Abs(i)
				os.Chdir(dsAbs)

				m.Comet = p.Comet
				gobExtC := fmt.Sprintf("*.%s", p.Comet.RawExtension)
				filesC, e := filepath.Glob(gobExtC)
				if e != nil {
					logrus.Fatal(e)
				}

				if len(filesC) > 0 {
					for _, j := range filesC {
						f, _ := filepath.Abs(j)
						mzFiles = append(mzFiles, f)
					}
				}

				m.MSFragger = p.MSFragger
				gobExtM := fmt.Sprintf("*.%s", p.MSFragger.RawExtension)
				filesM, e := filepath.Glob(gobExtM)
				if e != nil {
					logrus.Fatal(e)
				}

				if len(filesM) > 0 {
					for _, j := range filesM {
						f, _ := filepath.Abs(j)
						mzFiles = append(mzFiles, f)
					}
				}

				// return to the top level directory
				os.Chdir(dir)
			}

			// Comet
			if p.Commands.Comet == "yes" {
				comet.Run(m, mzFiles)
				m.Serialize()
			}

			// MSFragger
			if p.Commands.MSFragger == "yes" {
				msfragger.Run(m, mzFiles)
				m.Serialize()
			}
		}

		// For each dataset run the Prophets inside them
		for _, i := range args {

			// getting inside de the dataset folder
			dsAbs, _ := filepath.Abs(i)
			os.Chdir(dsAbs)

			// reload the meta data
			m.Restore(sys.Meta())

			// PeptideProphet
			if p.Commands.PeptideProphet == "yes" {
				logrus.Info("Executing PeptideProphet on ", i)
				m.PeptideProphet = p.PeptideProphet
				m.PeptideProphet.Output = "interact"
				m.PeptideProphet.Combine = true
				gobExt := fmt.Sprintf("*.%s", p.PeptideProphet.FileExtension)
				files, e := filepath.Glob(gobExt)
				if e != nil {
					logrus.Fatal(e.Error())
				}
				peptideprophet.Run(m, files)
				m.Serialize()
			}

			// PTMProphet
			if p.Commands.PTMProphet == "yes" {
				logrus.Info("Executing PTMProphet on ", i)
				m.PTMProphet = p.PTMProphet
				var files []string
				files = append(files, "interact.pep.xml")
				m.PTMProphet.InputFiles = files
				ptmprophet.Run(m, files)
				m.Serialize()
			}

			// ProteinProphet
			if p.Commands.ProteinProphet == "yes" {
				logrus.Info("Executing ProteinProphet on ", i)
				m.ProteinProphet = p.ProteinProphet
				m.ProteinProphet.Output = "interact"
				var files []string
				if p.Commands.PTMProphet == "yes" {
					files = append(files, "interact.mod.pep.xml")
				} else {
					files = append(files, "interact.pep.xml")
				}
				proteinprophet.Run(m, files)
				m.Serialize()
			}

			// return to the top level directory
			os.Chdir(dir)
		}

		// Abacus
		var combinedProtXML string
		if p.Commands.Abacus == "yes" && len(p.Filter.Pox) == 0 {
			logrus.Info("Creating combined protein inference")
			os.Chdir(dir)
			meta.Restore(sys.Meta())
			meta.ProteinProphet = p.ProteinProphet
			meta.ProteinProphet.Output = "combined"
			var files []string
			for _, j := range args {
				fqn := fmt.Sprintf("%s%sinteract.pep.xml", j, string(filepath.Separator))
				if p.Commands.PTMProphet == "yes" {
					fqn = fmt.Sprintf("%s%sinteract.mod.pep.xml", j, string(filepath.Separator))
				}
				fqn, _ = filepath.Abs(fqn)
				files = append(files, fqn)
			}

			os.Chdir(dir)

			proteinprophet.Run(meta, files)
			combinedProtXML = fmt.Sprintf("%s%scombined.prot.xml", meta.Temp, string(filepath.Separator))

			m.Filter.Pox = combinedProtXML

			// copy to work directory
			sys.CopyFile(combinedProtXML, filepath.Base(combinedProtXML))

			m.Serialize()
		}

		// for each data set, run the filter and quantify
		for _, i := range args {

			// getting inside  each dataset folder again
			dsAbs, _ := filepath.Abs(i)
			os.Chdir(dsAbs)

			// Filter
			if p.Commands.Filter == "yes" {
				logrus.Info("Executing filter on ", i)
				m.Filter = p.Filter

				if len(m.Filter.Pex) == 0 {
					m.Filter.Pex = "interact.pep.xml"
					if p.Commands.PTMProphet == "yes" {
						m.Filter.Pex = "interact.mod.pep.xml"
					}
					if p.Commands.ProteinProphet == "yes" {
						m.Filter.Pox = "interact.prot.xml"
					}
				}

				if len(m.Filter.Pox) == 0 {
					if p.Commands.Abacus == "yes" {
						m.Filter.Pox = combinedProtXML
					}
				}

				m, e := fil.Run(m)
				if e != nil {
					logrus.Fatal(e.Error())
				}

				m.Serialize()
			}

			// FreeQuant
			if p.Commands.FreeQuant == "yes" {
				logrus.Info("Executing label-free quantification on ", i)
				m.Quantify = p.Freequant
				m.Quantify.Dir = dsAbs
				m.Quantify.Format = "mzML"
				e := qua.RunLabelFreeQuantification(m.Quantify)
				if e != nil {
					logrus.Fatal(e.Error())
				}
				m.Serialize()
			}

			// LabelQuant
			if p.Commands.LabelQuant == "yes" {
				logrus.Info("Executing label-based quantification on ", i)
				m.Quantify = p.LabelQuant
				m.Quantify.Dir = dsAbs
				m.Quantify.Format = "mzML"
				m.Quantify.Brand = "tmt"
				var e error
				m.Quantify, e = qua.RunTMTQuantification(m.Quantify, m.Filter.Mapmods)
				if e != nil {
					logrus.Fatal(e)
				}
				m.Serialize()
			}

			// Report
			if p.Commands.Report == "yes" {
				logrus.Info("Executing report on ", i)
				m.Report = p.Report
				rep.Run(m)
				m.Serialize()
			}

			// Cluster
			if p.Commands.Cluster == "yes" {
				logrus.Info("Executing cluster on ", i)
				m.Cluster = p.Cluster
				clu.GenerateReport(m)
				m.Serialize()

			}

			// return to the top level directory
			os.Chdir(dir)
		}

		// Abacus
		if p.Commands.Abacus == "yes" {
			logrus.Info("Executing abacus")
			os.Chdir(dir)
			m.Abacus = p.Abacus
			err := aba.Run(m.Abacus, m.Temp, args)
			if err != nil {
				logrus.Fatal(err)
			}
		}

		for _, i := range args {

			// getting inside de the dataset folder
			localDir, _ := filepath.Abs(i)
			os.Chdir(localDir)

			// Backup
			if p.Backup == true {
				wrk.Run(Version, Build, true, false, false, true)
			}

			// Clean
			if p.Clean == true {
				wrk.Run(Version, Build, false, true, false, true)
			}

		}

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
