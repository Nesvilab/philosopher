package pip

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/prvst/philosopher/lib/ext/tmtintegrator"

	"github.com/prvst/philosopher/lib/aba"
	"github.com/prvst/philosopher/lib/clu"
	"github.com/prvst/philosopher/lib/ext/peptideprophet"
	"github.com/prvst/philosopher/lib/ext/proteinprophet"
	"github.com/prvst/philosopher/lib/ext/ptmprophet"
	"github.com/prvst/philosopher/lib/fil"
	"github.com/prvst/philosopher/lib/qua"
	"github.com/prvst/philosopher/lib/rep"

	"github.com/prvst/philosopher/lib/dat"
	"github.com/prvst/philosopher/lib/ext/comet"
	"github.com/prvst/philosopher/lib/ext/msfragger"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/prvst/philosopher/lib/wrk"
	"github.com/sirupsen/logrus"
)

// Directives contains the instructions to run a pipeline
type Directives struct {
	Initialize     bool               `yaml:"initialize"`
	Clean          bool               `yaml:"clean"`
	Backup         bool               `yaml:"backup"`
	Analtics       bool               `yaml:"analytics"`
	SlackToken     string             `yaml:"slackToken"`
	SlackChannel   string             `yaml:"slackChannel"`
	Commands       Commands           `yaml:"commands"`
	Database       met.Database       `yaml:"database"`
	MSFragger      met.MSFragger      `yaml:"msfragger"`
	Comet          met.Comet          `yaml:"comet"`
	PeptideProphet met.PeptideProphet `yaml:"peptideprophet"`
	PTMProphet     met.PTMProphet     `yaml:"ptmprophet"`
	ProteinProphet met.ProteinProphet `yaml:"proteinprophet"`
	Filter         met.Filter         `yaml:"filter"`
	Freequant      met.Quantify       `yaml:"freequant"`
	LabelQuant     met.Quantify       `yaml:"labelquant"`
	Report         met.Report         `yaml:"report"`
	Cluster        met.Cluster        `yaml:"cluster"`
	Abacus         met.Abacus         `yaml:"abacus"`
	TMTIntegrator  met.TMTIntegrator  `yaml:"tmt-integrator"`
}

// Commands struct {
type Commands struct {
	Database       string `yaml:"database"`
	MSFragger      string `yaml:"msfragger"`
	Comet          string `yaml:"comet"`
	PeptideProphet string `yaml:"peptideprophet"`
	PTMProphet     string `yaml:"ptmprophet"`
	ProteinProphet string `yaml:"proteinprophet"`
	Filter         string `yaml:"filter"`
	FreeQuant      string `yaml:"freequant"`
	LabelQuant     string `yaml:"labelquant"`
	Report         string `yaml:"report"`
	Cluster        string `yaml:"cluster"`
	Abacus         string `yaml:"abacus"`
	TMTIntegrator  string `yaml:"tmt-integrator"`
}

// DeployParameterFile ...
func DeployParameterFile(temp string) (string, error) {

	file := temp + string(filepath.Separator) + "philosopher.yml"

	param, err := Asset("philosopher.yml")
	if err != nil {
		return file, errors.New("Cannot deploy pipeline configuration file")
	}

	err = ioutil.WriteFile(file, param, sys.FilePermission())
	if err != nil {
		return file, errors.New("Cannot write pipeline parameter file")
	}

	return file, nil
}

// InitializeWorkspaces moves inside each data folder and initializes the Workspace with a database
func InitializeWorkspaces(meta met.Data, p Directives, dir, Version, Build string, data []string) met.Data {

	for _, i := range data {

		logrus.Info("Initiating the workspace on ", i)

		// getting inside de the dataset folder
		dsAbs, _ := filepath.Abs(i)
		os.Chdir(dsAbs)

		// Workspace
		wrk.Run(Version, Build, false, false, true, false)

		// reload the meta data
		meta.Restore(sys.Meta())

		// Database
		if p.Commands.Database == "yes" {
			meta.Database = p.Database
			dat.Run(meta)
			meta.Serialize()
		}

		if p.Commands.Comet == "yes" && p.Commands.MSFragger == "yes" {
			logrus.Fatal("You can only specify one search engine at a time")
		}

		// return to the top level directory
		os.Chdir(dir)
	}

	return meta
}

// DatabaseSearch executes the search engines if requested
func DatabaseSearch(meta met.Data, p Directives, dir string, data []string) met.Data {

	if p.Commands.Comet == "yes" || p.Commands.MSFragger == "yes" {

		logrus.Info("Running the Database Search on all data")

		// reload the meta data
		meta.Restore(sys.Meta())

		var mzFiles []string

		for _, i := range data {

			// getting inside de the dataset folder
			dsAbs, _ := filepath.Abs(i)
			os.Chdir(dsAbs)

			meta.Comet = p.Comet
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

			meta.MSFragger = p.MSFragger
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

			// reload the meta data
			meta.Restore(sys.Meta())
		}

		// Comet
		if p.Commands.Comet == "yes" {
			comet.Run(meta, mzFiles)
			meta.Serialize()
		}

		// MSFragger
		if p.Commands.MSFragger == "yes" {
			msfragger.Run(meta, mzFiles)
			meta.Serialize()
		}
	}

	return meta
}

// Prophets execute the TPP Prophets
func Prophets(meta met.Data, p Directives, dir string, data []string) met.Data {

	if p.Commands.PeptideProphet == "yes" || p.Commands.ProteinProphet == "yes" || p.Commands.PTMProphet == "yes" {
		for _, i := range data {

			logrus.Info("Running the validation and inference on ", i)

			// getting inside de the dataset folder
			dsAbs, _ := filepath.Abs(i)
			os.Chdir(dsAbs)

			// reload the meta data
			meta.Restore(sys.Meta())

			// PeptideProphet
			if p.Commands.PeptideProphet == "yes" {
				logrus.Info("Executing PeptideProphet on ", i)
				meta.PeptideProphet = p.PeptideProphet
				meta.PeptideProphet.Output = "interact"
				meta.PeptideProphet.Combine = true
				gobExt := fmt.Sprintf("*.%s", p.PeptideProphet.FileExtension)
				files, e := filepath.Glob(gobExt)
				if e != nil {
					logrus.Fatal(e.Error())
				}
				peptideprophet.Run(meta, files)
				meta.Serialize()
			}

			// PTMProphet
			if p.Commands.PTMProphet == "yes" {
				logrus.Info("Executing PTMProphet on ", i)
				meta.PTMProphet = p.PTMProphet
				var files []string
				files = append(files, "interact.pep.xml")
				meta.PTMProphet.InputFiles = files
				ptmprophet.Run(meta, files)
				meta.Serialize()
			}

			// ProteinProphet
			if p.Commands.ProteinProphet == "yes" {
				logrus.Info("Executing ProteinProphet on ", i)
				meta.ProteinProphet = p.ProteinProphet
				meta.ProteinProphet.Output = "interact"
				var files []string
				if p.Commands.PTMProphet == "yes" {
					files = append(files, "interact.mod.pep.xml")
				} else {
					files = append(files, "interact.pep.xml")
				}
				proteinprophet.Run(meta, files)
				meta.Serialize()
			}

			// return to the top level directory
			os.Chdir(dir)

			// reload the meta data
			meta.Restore(sys.Meta())
		}
	}

	return meta
}

// CombinedProteinList executes ProteinProphet command creating the combined ProtXML
func CombinedProteinList(meta met.Data, p Directives, dir string, data []string) met.Data {

	var combinedProtXML string

	if p.Commands.Abacus == "yes" && len(p.Filter.Pox) == 0 {

		logrus.Info("Creating combined protein inference")

		// return to the top level directory
		os.Chdir(dir)

		// reload the meta data
		meta.Restore(sys.Meta())

		meta.ProteinProphet = p.ProteinProphet
		meta.ProteinProphet.Output = "combined"

		var files []string

		for _, j := range data {
			fqn := fmt.Sprintf("%s%sinteract.pep.xml", j, string(filepath.Separator))
			if p.Commands.PTMProphet == "yes" {
				fqn = fmt.Sprintf("%s%sinteract.mod.pep.xml", j, string(filepath.Separator))
			}
			fqn, _ = filepath.Abs(fqn)
			files = append(files, fqn)
		}

		// return to the top level directory
		os.Chdir(dir)

		// reload the meta data
		meta.Restore(sys.Meta())

		proteinprophet.Run(meta, files)
		combinedProtXML = fmt.Sprintf("%s%scombined.prot.xml", meta.Temp, string(filepath.Separator))

		meta.Filter.Pox = combinedProtXML

		// copy to work directory
		sys.CopyFile(combinedProtXML, filepath.Base(combinedProtXML))

		meta.Serialize()
	}

	return meta
}

// FilterQuantifyReport executes the Filter, Quantify and Report commands in tandem
func FilterQuantifyReport(meta met.Data, p Directives, dir string, data []string) met.Data {

	for _, i := range data {

		// getting inside  each dataset folder again
		dsAbs, _ := filepath.Abs(i)
		os.Chdir(dsAbs)

		// reload the meta data
		meta.Restore(sys.Meta())

		// Filter
		if p.Commands.Filter == "yes" {

			logrus.Info("Executing filter on ", i)
			meta.Filter = p.Filter

			if len(meta.Filter.Pex) == 0 {
				meta.Filter.Pex = "interact.pep.xml"
			}

			if len(meta.Filter.Pox) == 0 {
				meta.Filter.Pox = "interact.prot.xml"
			}

			if len(meta.Filter.Pox) == 0 && p.Commands.Abacus == "yes" {
				meta.Filter.Pox = fmt.Sprintf("%s%scombined.prot.xml", meta.Temp, string(filepath.Separator))
			}

			meta, e := fil.Run(meta)
			if e != nil {
				logrus.Fatal(e.Error())
			}

			meta.Serialize()
		}

		// FreeQuant
		if p.Commands.FreeQuant == "yes" {

			logrus.Info("Executing label-free quantification on ", i)

			meta.Quantify = p.Freequant
			meta.Quantify.Dir = dsAbs
			meta.Quantify.Format = "mzML"

			e := qua.RunLabelFreeQuantification(meta.Quantify)
			if e != nil {
				logrus.Fatal(e.Error())
			}
			meta.Serialize()
		}

		// LabelQuant
		if p.Commands.LabelQuant == "yes" {

			logrus.Info("Executing label-based quantification on ", i)

			meta.Quantify = p.LabelQuant
			meta.Quantify.Dir = dsAbs
			meta.Quantify.Format = "mzML"
			meta.Quantify.Brand = "tmt"

			var e error
			meta.Quantify, e = qua.RunTMTQuantification(meta.Quantify, meta.Filter.Mapmods)
			if e != nil {
				logrus.Fatal(e)
			}
			meta.Serialize()
		}

		// Report
		if p.Commands.Report == "yes" {

			logrus.Info("Executing report on ", i)

			meta.Report = p.Report

			rep.Run(meta)
			meta.Serialize()
		}

		// Cluster
		if p.Commands.Cluster == "yes" {

			logrus.Info("Executing cluster on ", i)

			meta.Cluster = p.Cluster

			clu.GenerateReport(meta)
			meta.Serialize()

		}

		// return to the top level directory
		os.Chdir(dir)

		// reload the meta data
		meta.Restore(sys.Meta())
	}

	return meta
}

// Abacus loads all data and creates the combined protein report
func Abacus(meta met.Data, p Directives, dir string, data []string) met.Data {

	if p.Commands.Abacus == "yes" {

		logrus.Info("Executing abacus")

		// return to the top level directory
		os.Chdir(dir)

		// reload the meta data
		meta.Restore(sys.Meta())

		meta.Abacus = p.Abacus
		err := aba.Run(meta.Abacus, meta.Temp, data)
		if err != nil {
			logrus.Fatal(err)
		}
	}

	return meta
}

// BackupAndClean stores the results in a zip file and removes all meta data
func BackupAndClean(meta met.Data, p Directives, dir, Version, Build string, data []string) {

	logrus.Info("Savig results and cleaning the workspaces")

	for _, i := range data {

		// getting inside de the dataset folder
		localDir, _ := filepath.Abs(i)
		os.Chdir(localDir)

		// reload the meta data
		meta.Restore(sys.Meta())

		// Backup
		if p.Backup == true {
			wrk.Run(Version, Build, true, false, false, true)
		}

		// Clean
		if p.Clean == true {
			wrk.Run(Version, Build, false, true, false, true)
		}

	}

	return
}

// TMTIntegrator executes TMT-I on all PSM results
func TMTIntegrator(meta met.Data, p Directives, dir string, data []string) met.Data {

	if p.Commands.TMTIntegrator == "yes" {

		logrus.Info("Running TMT-Integrator")

		// reload the meta data
		meta.Restore(sys.Meta())

		var psms []string

		for _, i := range data {

			// getting inside de the dataset folder
			dsAbs, _ := filepath.Abs(i)
			os.Chdir(dsAbs)

			// reload the meta data
			meta.Restore(sys.Meta())

			meta.TMTIntegrator = p.TMTIntegrator

			psms = append(psms, fmt.Sprintf("%s%spsm.tsv", dsAbs, string(filepath.Separator)))

			_, err := tmtintegrator.Run(meta, psms)
			if err != nil {
				logrus.Fatal(err)
			}

		}

	}

	return meta
}
