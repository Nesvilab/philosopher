package pip

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"philosopher/lib/msg"

	"philosopher/lib/ext/interprophet"
	"philosopher/lib/ext/tmtintegrator"

	"philosopher/lib/aba"
	"philosopher/lib/ext/peptideprophet"
	"philosopher/lib/ext/proteinprophet"
	"philosopher/lib/ext/ptmprophet"
	"philosopher/lib/fil"
	"philosopher/lib/qua"
	"philosopher/lib/rep"

	"github.com/sirupsen/logrus"
	"philosopher/lib/dat"
	"philosopher/lib/ext/comet"
	"philosopher/lib/ext/msfragger"
	"philosopher/lib/met"
	"philosopher/lib/sys"
	"philosopher/lib/wrk"
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
	BioQuant       met.BioQuant       `yaml:"bioquant"`
	Abacus         met.Abacus         `yaml:"abacus"`
	TMTIntegrator  met.TMTIntegrator  `yaml:"tmtintegrator"`
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
	BioQuant       string `yaml:"bioquant"`
	Abacus         string `yaml:"abacus"`
	TMTIntegrator  string `yaml:"tmtintegrator"`
}

// DeployParameterFile deploys the pipeline yaml config file
func DeployParameterFile(temp string) string {

	file := temp + string(filepath.Separator) + "philosopher.yml"

	param, e := Asset("philosopher.yml")
	if e != nil {
		msg.DeployAsset(errors.New("pipeline configuration file"), "fatal")
	}

	e = ioutil.WriteFile(file, param, sys.FilePermission())
	if e != nil {
		msg.DeployAsset(errors.New("pipeline configuration file"), "fatal")
	}

	return file
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
			msg.Custom(errors.New("You can only specify one search engine at a time"), "fatal")
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
				msg.Custom(e, "fatal")
			}

			if len(filesC) > 0 {
				for _, j := range filesC {
					f, _ := filepath.Abs(j)
					mzFiles = append(mzFiles, f)
				}
			}

			meta.MSFragger = p.MSFragger
			meta.MSFragger.DatabaseName = p.Database.Annot

			gobExtM := fmt.Sprintf("*.%s", p.MSFragger.RawExtension)
			filesM, e := filepath.Glob(gobExtM)
			if e != nil {
				msg.Custom(e, "fatal")
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
			//meta.Restore(sys.Meta())
		}

		// Comet
		if p.Commands.Comet == "yes" {
			comet.Run(meta, mzFiles)
			meta.SearchEngine = "Comet"
			//meta.Serialize()
		}

		// MSFragger
		if p.Commands.MSFragger == "yes" {
			msfragger.Run(meta, mzFiles)
			meta.SearchEngine = "MSFragger"
			//meta.Serialize()
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
				meta.PeptideProphet.Database = p.Database.Annot
				meta.PeptideProphet.Decoy = p.Database.Tag
				meta.PeptideProphet.Output = "interact"
				meta.PeptideProphet.Combine = true
				gobExt := fmt.Sprintf("*.%s", p.PeptideProphet.FileExtension)
				files, e := filepath.Glob(gobExt)
				if e != nil {
					msg.Custom(e, "fatal")
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
			//meta.Restore(sys.Meta())
		}
	}

	return meta
}

// ParallelProphets execute the TPP Prophets
func ParallelProphets(meta met.Data, p Directives, dir string, data []string) met.Data {

	var wg sync.WaitGroup
	wg.Add(len(data))

	if p.Commands.PeptideProphet == "yes" || p.Commands.ProteinProphet == "yes" || p.Commands.PTMProphet == "yes" {

		if p.Commands.PeptideProphet == "yes" {
			for _, ds := range data {

				db := p.Database.Annot

				go func(ds, db string) {
					defer wg.Done()

					logrus.Info("Running the validation and inference on ", ds)

					// getting inside de the dataset folder
					dsAbs, _ := filepath.Abs(ds)
					absMeta := fmt.Sprintf("%s%s%s", dsAbs, string(filepath.Separator), sys.Meta())

					// reload the meta data
					meta.Restore(absMeta)

					// PeptideProphet
					logrus.Info("Executing PeptideProphet on ", ds)
					meta.PeptideProphet = p.PeptideProphet
					meta.PeptideProphet.Database = p.Database.Annot
					meta.PeptideProphet.Decoy = p.Database.Tag
					meta.PeptideProphet.Output = "interact"
					meta.PeptideProphet.Combine = true

					gobExt := fmt.Sprintf("%s%s*.%s", dsAbs, string(filepath.Separator), p.PeptideProphet.FileExtension)

					files, e := filepath.Glob(gobExt)
					if e != nil {
						msg.Custom(e, "fatal")
					}

					peptideprophet.Run(meta, files)

					// give a chance to the execution to untangle the output
					time.Sleep(time.Second * 1)

					//meta.Serialize()

				}(ds, db)
			}

			wg.Wait()
		}

		// 	// PTMProphet
		// 	if p.Commands.PTMProphet == "yes" {
		// 		logrus.Info("Executing PTMProphet on ", i)
		// 		meta.PTMProphet = p.PTMProphet
		// 		var files []string
		// 		files = append(files, "interact.pep.xml")
		// 		meta.PTMProphet.InputFiles = files
		// 		ptmprophet.Run(meta, files)
		// 		meta.Serialize()
		// 	}

		// 	// ProteinProphet
		// 	if p.Commands.ProteinProphet == "yes" {
		// 		logrus.Info("Executing ProteinProphet on ", i)
		// 		meta.ProteinProphet = p.ProteinProphet
		// 		meta.ProteinProphet.Output = "interact"
		// 		var files []string
		// 		if p.Commands.PTMProphet == "yes" {
		// 			files = append(files, "interact.mod.pep.xml")
		// 		} else {
		// 			files = append(files, "interact.pep.xml")
		// 		}
		// 		proteinprophet.Run(meta, files)
		// 		meta.Serialize()
		// 	}

		// return to the top level directory
		//os.Chdir(dir)

		// reload the meta data
		//meta.Restore(sys.Meta())
		// /}
	}

	os.Chdir(dir)

	return meta
}

// CombinedPeptideList executes iProphet command creating the combined PepXML
func CombinedPeptideList(meta met.Data, p Directives, dir string, data []string) met.Data {

	var combinedPepXML string

	if p.Commands.Abacus == "yes" && p.Abacus.Peptide == true && len(p.Filter.Pex) == 0 {

		logrus.Info("Integrating peptide validation")

		// return to the top level directory
		os.Chdir(dir)

		// reload the meta data
		meta.Restore(sys.Meta())

		var files []string

		for _, j := range data {
			fqn := fmt.Sprintf("%s%sinteract.pep.xml", j, string(filepath.Separator))
			// if p.Commands.PTMProphet == "yes" {
			// 	fqn = fmt.Sprintf("%s%sinteract.mod.pep.xml", j, string(filepath.Separator))
			// }
			fqn, _ = filepath.Abs(fqn)
			files = append(files, fqn)
		}

		meta.InterProphet.Output = "combined"
		meta.InterProphet.Nonsp = true
		meta.InterProphet.InputFiles = files
		meta.InterProphet.Decoy = "rev_"
		meta.InterProphet.Threads = 6

		// run
		meta = interprophet.Run(meta, files)

		combinedPepXML = fmt.Sprintf("%s%scombined.pep.xml", meta.Temp, string(filepath.Separator))

		// copy to work directory
		sys.CopyFile(combinedPepXML, filepath.Base(combinedPepXML))

		//meta.Serialize()
	}

	return meta
}

// CombinedProteinList executes ProteinProphet command creating the combined ProtXML
func CombinedProteinList(meta met.Data, p Directives, dir string, data []string) met.Data {

	var combinedProtXML string

	if p.Commands.Abacus == "yes" && p.Abacus.Protein == true && len(p.Filter.Pox) == 0 {

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

		proteinprophet.Run(meta, files)
		combinedProtXML = fmt.Sprintf("%s%scombined.prot.xml", meta.Temp, string(filepath.Separator))

		meta.Filter.Pox = combinedProtXML

		// copy to work directory
		sys.CopyFile(combinedProtXML, filepath.Base(combinedProtXML))

		//meta.Serialize()
	}

	return meta
}

// FilterQuantifyReport executes the Filter, Quantify and Report commands in tandem
func FilterQuantifyReport(meta met.Data, p Directives, dir string, data []string) met.Data {

	// this is the virtual home directory where the pipeline is being executed.
	vHome := meta.Home

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
			meta.Filter.Tag = p.Database.Tag

			if len(p.Filter.Pex) == 0 {
				meta.Filter.Pex = "interact.pep.xml"
				if p.Commands.PTMProphet == "yes" {
					meta.Filter.Pex = "interact.mod.pep.xml"
				}
			} else {
				meta.Filter.Pex = p.Filter.Pex
			}

			if len(p.Filter.Pox) == 0 {
				meta.Filter.Pox = "interact.prot.xml"
			} else {
				meta.Filter.Pox = p.Filter.Pox
			}

			if p.Commands.Abacus == "yes" && p.Abacus.Protein == true && len(p.Filter.Pox) == 0 {
				meta.Filter.Pox = fmt.Sprintf("%s%scombined.prot.xml", vHome, string(filepath.Separator))
			}

			meta := fil.Run(meta)

			meta.Serialize()
		}

		// FreeQuant
		if p.Commands.FreeQuant == "yes" {

			logrus.Info("Executing label-free quantification on ", i)

			meta.Quantify = p.Freequant
			meta.Quantify.Dir = dsAbs
			meta.Quantify.Format = "mzML"

			qua.RunLabelFreeQuantification(meta.Quantify)

			meta.Serialize()
		}

		// LabelQuant
		if p.Commands.LabelQuant == "yes" {

			logrus.Info("Executing label-based quantification on ", i)

			meta.Quantify = p.LabelQuant
			meta.Quantify.Dir = dsAbs
			meta.Quantify.Format = "mzML"
			meta.Quantify.Brand = "tmt"

			meta.Quantify = qua.RunIsobaricLabelQuantification(meta.Quantify, meta.Filter.Mapmods)

			meta.Serialize()
		}

		// Report
		if p.Commands.Report == "yes" {

			logrus.Info("Executing report on ", i)

			meta.Report = p.Report

			rep.Run(meta)
			meta.Serialize()
		}

		// BioQuant
		if p.Commands.BioQuant == "yes" {

			logrus.Info("Executing cluster on ", i)

			meta.BioQuant = p.BioQuant

			//clu.GenerateReport(meta)
			qua.RunBioQuantification(meta)
			meta.Serialize()

		}

		// return to the top level directory
		os.Chdir(dir)

		// reload the meta data
		//meta.Restore(sys.Meta())
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
		meta.Abacus.Tag = p.Database.Tag
		meta.Abacus.Picked = p.Filter.Picked
		meta.Abacus.Razor = p.Filter.Razor

		if len(p.LabelQuant.Annot) > 0 {
			meta.Abacus.Labels = true
		}

		aba.Run(meta, data)
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
			meta.TMTIntegrator = p.TMTIntegrator
			psms = append(psms, fmt.Sprintf("%s%spsm.tsv", i, string(filepath.Separator)))
		}

		tmtintegrator.Run(meta, psms)
	}

	return meta
}
