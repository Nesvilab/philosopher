package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher/lib/dat"
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

		logrus.Info("Executing the pipeline")

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

		// Creating the Workspace
		for _, i := range args {
			localDir, _ := filepath.Abs(i)
			os.Chdir(localDir)
			//logrus.Info("Creating Workspace on ", i)
			wrk.Run(Version, Build, false, false, false, true)
		}

		// Annotating the database
		if p.Commands.Database == "yes" {
			for _, i := range args {
				localDir, _ := filepath.Abs(i)
				os.Chdir(localDir)
				//logrus.Info("Creating Database on ", i)
				m.Database = p.Database
				dat.Run(m)
			}
		}

		logrus.Info("Done")
		return
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "pipeline" {

		m.Restore(sys.Meta())

		pipelineCmd.Flags().StringVarP(&m.Pipeline.Directives, "config", "", "", "configuration file for the pipeline execution")
	}

	RootCmd.AddCommand(pipelineCmd)
}
