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
	"github.com/prvst/philosopher/lib/pip"
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
		wrk.Run(Version, Build, false, false, false, true)

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

		logrus.Info("Done")
		return
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "pipeline" {

		m.Restore(sys.Meta())

		pipelineCmd.Flags().StringVarP(&m.Pipeline.Directives, "config", "", "", "configuration file for the pipeline execution")
		pipelineCmd.Flags().StringVarP(&m.Pipeline.Dataset, "dataset", "", "", "dataset directory")

	}

	RootCmd.AddCommand(pipelineCmd)
}
