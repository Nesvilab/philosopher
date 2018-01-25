package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher/lib/dat"
	"github.com/prvst/philosopher/lib/ext/comet"
	"github.com/prvst/philosopher/lib/ext/peptideprophet"
	"github.com/prvst/philosopher/lib/ext/proteinprophet"
	"github.com/prvst/philosopher/lib/fil"
	"github.com/prvst/philosopher/lib/pip"
	"github.com/prvst/philosopher/lib/qua"
	"github.com/prvst/philosopher/lib/rep"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/prvst/philosopher/lib/wrk"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

// pipelineCmd represents the pipeline command
var pipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {

		param, e := pip.DeployParameterFile(m.Temp)
		if e != nil {
			logrus.Fatal(e.Error())
		}

		if m.Pipeline.Print == true {
			logrus.Info("Printing parameter file")
			sys.CopyFile(param, filepath.Base(param))
			return
		}

		logrus.Info("Executing the pipeline on ", m.Pipeline.Dataset)

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

		// getting inside de the dataset folder
		localDir, _ := filepath.Abs(m.Pipeline.Dataset)
		os.Chdir(localDir)

		// Workspace
		wrk.Run(Version, Build, false, false, true)

		// reload the meta data
		m.Restore(sys.Meta())

		// Database
		if p.Commands.Database == "yes" {
			m.Database = p.Database
			dat.Run(m)
		}

		// Comet
		if p.Commands.Comet == "yes" {
			m.Comet = p.Comet

			gobExt := fmt.Sprintf("*.%s", p.Comet.RawExtension)
			files, e := filepath.Glob(gobExt)
			if e != nil {
				logrus.Fatal(e)
			}

			comet.Run(m, files)
		}

		// PeptideProphet
		if p.Commands.PeptideProphet == "yes" {
			logrus.Info("Executing PeptideProphet")

			m.PeptideProphet = p.PeptideProphet
			m.PeptideProphet.Output = "interact"
			m.PeptideProphet.Combine = true

			gobExt := fmt.Sprintf("*.%s", p.PeptideProphet.FileExtension)
			files, e := filepath.Glob(gobExt)
			if e != nil {
				logrus.Fatal(e)
			}

			peptideprophet.Run(m, files)
		}

		// ProteinProphet
		if p.Commands.ProteinProphet == "yes" {
			logrus.Info("Executing ProteinProphet")

			m.ProteinProphet = p.ProteinProphet
			m.ProteinProphet.Output = "interact"

			var files []string
			files = append(files, "interact.pep.xml")

			proteinprophet.Run(m, files)
		}

		// Filter
		if p.Commands.Filter == "yes" {
			logrus.Info("Executing filter")

			m.Filter = p.Filter
			m.Filter.Pex = "interact.pep.xml"

			if p.Commands.ProteinProphet == "yes" {
				m.Filter.Pox = "interact.prot.xml"
			}

			e := fil.Run(m.Filter)
			if e != nil {
				logrus.Fatal(e.Error())
			}
		}

		// FreeQuant
		if p.Commands.FreeQuant == "yes" {
			logrus.Info("Executing label-free quantification")

			m.Quantify = p.Freequant
			m.Quantify.Dir = localDir
			m.Quantify.Format = "mzML"

			// run label-free quantification
			e := qua.RunLabelFreeQuantification(m.Quantify)
			if e != nil {
				logrus.Fatal(e.Error())
			}
		}

		// LabelQuant
		if p.Commands.LabelQuant == "yes" {
			logrus.Info("Executing label-based quantification")

			m.Quantify = p.LabelQuant
			m.Quantify.Dir = localDir
			m.Quantify.Format = "mzML"
			m.Quantify.Brand = "tmt"

			err := qua.RunTMTQuantification(m.Quantify)
			if err != nil {
				logrus.Fatal(err)
			}
		}

		if p.Commands.Report == "yes" {
			logrus.Info("Executing report")

			rep.Run(m)
		}

		m.Serialize()

		// Backup
		if p.Backup == true {
			wrk.Run(Version, Build, true, false, false)
		}

		// Clean
		if p.Clean == true {
			wrk.Run(Version, Build, false, true, false)
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
		pipelineCmd.Flags().StringVarP(&m.Pipeline.Dataset, "dataset", "", "", "dataset directory")

	}

	RootCmd.AddCommand(pipelineCmd)
}
