package pip

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"philosopher/lib/dat"
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

	"philosopher/lib/ext/comet"
	"philosopher/lib/ext/msfragger"
	"philosopher/lib/met"
	"philosopher/lib/sys"
	"philosopher/lib/wrk"

	"github.com/ryanskidmore/parallel"
	"github.com/sirupsen/logrus"
)

// Directives contains the instructions to run a pipeline
type Directives struct {
	SlackToken     string             `yaml:"Slack Token"`
	SlackChannel   string             `yaml:"Slack Channel"`
	SlackUserID    string             `yaml:"Slack User ID"`
	Steps          Steps              `yaml:"Steps"`
	DatabaseSearch DatabaseSearch     `yaml:"Database Search"`
	PeptideProphet met.PeptideProphet `yaml:"Peptide Validation"`
	PTMProphet     met.PTMProphet     `yaml:"PTM Localization"`
	ProteinProphet met.ProteinProphet `yaml:"Protein Inference"`
	Filter         met.Filter         `yaml:"FDR Filtering"`
	Freequant      met.Quantify       `yaml:"Label-Free Quantification"`
	LabelQuant     met.Quantify       `yaml:"Isobaric Quantification"`
	Report         met.Report         `yaml:"Individual Reports"`
	BioQuant       met.BioQuant       `yaml:"Bio Cluster Quantification"`
	Abacus         met.Abacus         `yaml:"Integrated Reports"`
	TMTIntegrator  met.TMTIntegrator  `yaml:"Integrated Isobaric Quantification"`
}

// Steps contains the high-level elements of the analysis to be executed
type Steps struct {
	DatabaseSearch           string `yaml:"Database Search"`
	PeptideValidation        string `yaml:"Peptide Validation"`
	PTMLocalization          string `yaml:"PTM Localization"`
	ProteinInference         string `yaml:"Protein Inference"`
	LabelFreeQuantification  string `yaml:"Label-Free Quantification"`
	IsobaricQuantification   string `yaml:"Isobaric Quantification"`
	BioClusterQuantification string `yaml:"Bio Cluster Quantification"`
	FDRFiltering             string `yaml:"FDR Filtering"`
	IndividualReports        string `yaml:"Individual Reports"`
	IntegratedReports        string `yaml:"Integrated Reports"`
	TMTIntegrator            string `yaml:"Integrated Isobaric Quantification"`
}

// DatabaseSearch keeps the options related to the search step
type DatabaseSearch struct {
	SearchEngine    string        `yaml:"search_engine"`
	ProteinDatabase string        `yaml:"protein_database"`
	DecoyTag        string        `yaml:"decoy_tag"`
	MSFragger       met.MSFragger `yaml:"msfragger"`
	Comet           met.Comet     `yaml:"comet"`
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

	// Top-level Workspace
	//wrk.Run(Version, Build, "", false, false, true, true)

	for i := range data {

		logrus.Info("Initiating the workspace on ", data[i])

		// getting inside de the dataset folder
		dsAbs, _ := filepath.Abs(data[i])
		os.Chdir(dsAbs)

		// Workspace
		wrk.Run(Version, Build, "", false, false, true, true)

		// reload the meta data to avoid being overwritten with black home
		meta.Restore(sys.Meta())

		met.CleanTemp(meta.Temp)

		// return to the top level directory
		os.Chdir(dir)
	}

	return meta
}

// AnnotateDatabase annotates the database on the first ds, and copy the bin data to the other folders
func AnnotateDatabase(meta met.Data, p Directives, dir string, data []string) met.Data {

	var source string

	// getting inside de the dataset folder
	dsAbs, _ := filepath.Abs(data[0])
	os.Chdir(dsAbs)

	// reload the meta data to avoid being overwritten with black home
	meta.Restore(sys.Meta())

	meta.Database.Annot = p.DatabaseSearch.ProteinDatabase
	meta.Database.Tag = p.DatabaseSearch.DecoyTag
	dat.Run(meta)
	meta.Serialize()
	source = fmt.Sprintf("%s%s.meta%sdb.bin", dsAbs, string(filepath.Separator), string(filepath.Separator))

	// return to the top level directory
	os.Chdir(dir)

	for i := 1; i < len(data); i++ {
		//source := fmt.Sprintf("%s%s.meta%sdb.bin", data[0], string(filepath.Separator), string(filepath.Separator))
		destination := fmt.Sprintf("%s%s.meta%sdb.bin", data[i], string(filepath.Separator), string(filepath.Separator))

		// Read all content of src to data
		data, e := ioutil.ReadFile(source)
		if e != nil {
			log.Fatal(e)
		}

		e = ioutil.WriteFile(destination, data, 0644)
		if e != nil {
			log.Fatal(e)
		}

		//meta.Serialize()
	}

	met.CleanTemp(meta.Temp)

	return meta
}

// DBSearch executes the search engines if requested
func DBSearch(meta met.Data, p Directives, dir string, data []string) met.Data {

	logrus.Info("Running the Database Search")

	// if meta.Pipeline.Verbose == true {
	// 	p.DatabaseSearch.MSFragger.ToCmdString()
	// }

	// reload the meta data
	meta.Restore(sys.Meta())

	var mzFiles []string

	if p.DatabaseSearch.SearchEngine == "comet" {

		for _, i := range data {

			// getting inside de the dataset folder
			dsAbs, _ := filepath.Abs(i)
			os.Chdir(dsAbs)

			meta.Comet = p.DatabaseSearch.Comet
			gobExtC := fmt.Sprintf("*.%s", p.DatabaseSearch.Comet.RawExtension)
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

			// return to the top level directory
			os.Chdir(dir)
		}

		comet.Run(meta, mzFiles)
		meta.SearchEngine = "Comet"
	}

	if p.DatabaseSearch.SearchEngine == "msfragger" {

		for _, i := range data {

			// getting inside de the dataset folder
			dsAbs, _ := filepath.Abs(i)
			os.Chdir(dsAbs)

			meta.MSFragger = p.DatabaseSearch.MSFragger
			meta.MSFragger.DatabaseName = p.DatabaseSearch.ProteinDatabase
			meta.MSFragger.DecoyPrefix = p.DatabaseSearch.DecoyTag

			gobExtM := fmt.Sprintf("*.%s", p.DatabaseSearch.MSFragger.Extension)
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
		}

		// MSFragger
		if p.DatabaseSearch.SearchEngine == "msfragger" {
			msfragger.Run(meta, mzFiles)
			meta.SearchEngine = "MSFragger"
		}
	}

	met.CleanTemp(meta.Temp)

	return meta
}

// PeptideProphet executes PeptideProphet in Parallel mode
func PeptideProphet(meta met.Data, p Directives, dir string, data []string) met.Data {

	if !p.PeptideProphet.Concurrent {
		for _, i := range data {

			logrus.Info("Running the validation and inference on ", i)

			// getting inside de the dataset folder
			dsAbs, _ := filepath.Abs(i)
			os.Chdir(dsAbs)

			// // reload the meta data
			meta.Restore(sys.Meta())

			// PeptideProphet
			if p.Steps.PeptideValidation == "yes" {
				logrus.Info("Executing PeptideProphet on ", i)
				meta.PeptideProphet = p.PeptideProphet
				meta.PeptideProphet.Database = p.DatabaseSearch.ProteinDatabase
				meta.PeptideProphet.Decoy = p.DatabaseSearch.DecoyTag
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

			// return to the top level directory
			os.Chdir(dir)
		}
	} else {

		pll := parallel.New()
		worker, err := pll.NewWorker("worker", &parallel.WorkerConfig{Parallelism: len(data)})

		if err != nil {
			log.Fatalf("FATAL: Failed to create new worker: %v", err)
		}

		worker.SetExecution(func(wh *parallel.WorkerHelper, args interface{}) {

			ds := fmt.Sprintf("%v", args)

			logrus.Info("Running the validation and inference on ", ds)

			meta.Restore(sys.Meta())
			meta.Database.Annot = p.DatabaseSearch.ProteinDatabase
			meta.Database.Tag = p.DatabaseSearch.DecoyTag

			// getting inside de the dataset folder
			dsAbs, _ := filepath.Abs(ds)
			absMeta := fmt.Sprintf("%s%s%s", dsAbs, string(filepath.Separator), sys.Meta())

			// reload the meta data
			meta.Restore(absMeta)

			// PeptideProphet
			logrus.Info("Executing PeptideProphet on ", ds)
			meta.PeptideProphet = p.PeptideProphet
			meta.PeptideProphet.Database = p.DatabaseSearch.ProteinDatabase
			meta.PeptideProphet.Decoy = p.DatabaseSearch.DecoyTag
			meta.PeptideProphet.Output = "interact"
			meta.PeptideProphet.Combine = true

			gobExt := fmt.Sprintf("%s%s*.%s", dsAbs, string(filepath.Separator), p.PeptideProphet.FileExtension)

			files, e := filepath.Glob(gobExt)
			if e != nil {
				msg.Custom(e, "fatal")
			}

			peptideprophet.Run(meta, files)

			wh.Done()
		})

		for _, ds := range data {
			worker.Start(interface{}(ds))
		}
		worker.Wait()
	}

	os.Chdir(dir)

	return meta
}

// PTMProphet execute the TPP PTMProphet
func PTMProphet(meta met.Data, p Directives, dir string, data []string) met.Data {

	for _, i := range data {

		logrus.Info("Running the validation and inference on ", i)

		// getting inside de the dataset folder
		dsAbs, _ := filepath.Abs(i)
		os.Chdir(dsAbs)

		// reload the meta data
		meta.Restore(sys.Meta())

		// PTMProphet
		if p.Steps.PTMLocalization == "yes" {
			logrus.Info("Executing PTMProphet on ", i)
			meta.PTMProphet = p.PTMProphet
			var files []string
			files = append(files, "interact.pep.xml")
			meta.PTMProphet.InputFiles = files
			meta.PTMProphet.KeepOld = true
			ptmprophet.Run(meta, files)
			meta.Serialize()
		}

		// return to the top level directory
		os.Chdir(dir)
	}

	return meta
}

// ProteinProphet execute the TPP ProteinProphet
func ProteinProphet(meta met.Data, p Directives, dir string, data []string) met.Data {

	for _, i := range data {

		logrus.Info("Running the validation and inference on ", i)

		// getting inside de the dataset folder
		dsAbs, _ := filepath.Abs(i)
		os.Chdir(dsAbs)

		// reload the meta data
		meta.Restore(sys.Meta())

		// ProteinProphet
		if p.Steps.ProteinInference == "yes" {
			logrus.Info("Executing ProteinProphet on ", i)
			meta.ProteinProphet = p.ProteinProphet
			meta.ProteinProphet.Output = "interact"
			var files []string
			if p.Steps.PTMLocalization == "yes" {
				files = append(files, "interact.mod.pep.xml")
			} else {
				files = append(files, "interact.pep.xml")
			}
			proteinprophet.Run(meta, files)
			meta.Serialize()
			met.CleanTemp(meta.Temp)
		}

		// return to the top level directory
		os.Chdir(dir)
	}

	return meta
}

// CombinedPeptideList executes iProphet command creating the combined PepXML
func CombinedPeptideList(meta met.Data, p Directives, dir string, data []string) met.Data {

	var combinedPepXML string

	if p.Steps.IntegratedReports == "yes" && p.Abacus.Peptide && len(p.Filter.Pex) == 0 {

		logrus.Info("Integrating peptide validation")

		// return to the top level directory
		os.Chdir(dir)

		// reload the meta data
		meta.Restore(sys.Meta())

		var files []string

		for _, j := range data {
			fqn := fmt.Sprintf("%s%sinteract.pep.xml", j, string(filepath.Separator))
			fqn, _ = filepath.Abs(fqn)
			files = append(files, fqn)
		}

		meta.Home = dir
		meta.InterProphet.Output = "combined"
		meta.InterProphet.Nonsp = true
		meta.InterProphet.InputFiles = files
		meta.InterProphet.Decoy = "rev_"
		meta.InterProphet.Threads = 6
		meta.InterProphet.MinProb = p.Abacus.PepProb

		// run
		meta = interprophet.Run(meta, files)

		combinedPepXML = fmt.Sprintf("%s%scombined.pep.xml", meta.Temp, string(filepath.Separator))

		// copy to work directory
		sys.CopyFile(combinedPepXML, filepath.Base(combinedPepXML))
	}

	return meta
}

// CombinedProteinList executes ProteinProphet command creating the combined ProtXML
func CombinedProteinList(meta met.Data, p Directives, dir string, data []string) met.Data {

	var combinedProtXML string

	logrus.Info("Creating combined protein inference")

	if p.Steps.IntegratedReports == "yes" && p.Abacus.Protein && len(p.Filter.Pox) == 0 {

		// return to the top level directory
		os.Chdir(dir)

		// reload the meta data
		meta.Restore(sys.Meta())

		meta.Home = dir
		meta.ProteinProphet = p.ProteinProphet
		meta.ProteinProphet.Output = "combined"
		meta.ProteinProphet.Minprob = p.Abacus.ProtProb

		var files []string

		for _, j := range data {
			fqn := fmt.Sprintf("%s%sinteract.pep.xml", j, string(filepath.Separator))
			if p.Steps.PTMLocalization == "yes" {
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
	}

	// when running the pipeline, we want to run the inference only once
	// and them copy the meta data to all data folders.
	os.Chdir(dir)

	protXML := fil.ReadProtXMLInput("combined.prot.xml", p.DatabaseSearch.DecoyTag, p.Filter.Weight)
	proBin := fil.ProcessProteinIdentifications(protXML, p.Filter.PtFDR, p.Filter.PepFDR, p.Filter.ProtProb, p.Abacus.Picked, p.Abacus.Razor, p.Filter.Fo, true, p.DatabaseSearch.DecoyTag)

	for _, i := range data {
		dest := fmt.Sprintf("%s%s.meta%spro.bin", i, string(filepath.Separator), string(filepath.Separator))
		sys.CopyFile(proBin, dest)
	}

	e := os.RemoveAll(path.Dir(proBin))
	if e != nil {
		log.Fatal(e)
	}

	return meta
}

// FreeQuant executes the LFQ method
func FreeQuant(meta met.Data, p Directives, dir string, data []string) met.Data {

	for _, i := range data {

		// getting inside  each dataset folder again
		dsAbs, _ := filepath.Abs(i)
		os.Chdir(dsAbs)

		// reload the meta data
		meta.Restore(sys.Meta())

		if _, err := os.Stat(sys.LFQBin()); err == nil {
			return meta
		}

		logrus.Info("Executing label-free quantification on ", i)

		meta.Quantify = p.Freequant
		meta.Quantify.Dir = dsAbs
		meta.Quantify.Format = "mzML"
		meta.Quantify.Pex = fmt.Sprintf("%s%sinteract.pep.xml", dsAbs, string(filepath.Separator))
		meta.Quantify.Tag = "rev_"

		qua.RunLabelFreeQuantification(meta.Quantify)

		meta.Serialize()

		// return to the top level directory
		os.Chdir(dir)

	}

	return meta
}

// LabelQuant executes the isobaric-tag quantification method
func LabelQuant(meta met.Data, p Directives, dir string, data []string) met.Data {

	for _, i := range data {

		// getting inside  each dataset folder again
		dsAbs, _ := filepath.Abs(i)
		os.Chdir(dsAbs)

		//annotation, _ := filepath.Glob("annotation*.txt")
		//fullAnnotation, _ := filepath.Abs(annotation[0])

		fullAnnotation := "annotation.txt"

		// reload the meta data
		meta.Restore(sys.Meta())

		if _, err := os.Stat(sys.IsoBin()); err == nil {
			return meta
		}

		logrus.Info("Executing label-based quantification on ", i)

		meta.Quantify = p.LabelQuant
		meta.Quantify.Dir = dsAbs
		meta.Quantify.Format = "mzML"
		meta.Quantify.Annot = fullAnnotation
		meta.Quantify.Brand = p.LabelQuant.Brand
		meta.Quantify.Pex = fmt.Sprintf("%s%sinteract.pep.xml", dsAbs, string(filepath.Separator))
		meta.Quantify.Tag = "rev_"

		meta.Quantify = qua.RunIsobaricLabelQuantification(meta.Quantify, meta.Filter.Mapmods)

		meta.Serialize()

		// return to the top level directory
		os.Chdir(dir)
	}

	return meta
}

// BioQuant executes the bioquant quantification method
func BioQuant(meta met.Data, p Directives, dir string, data []string) met.Data {

	for _, i := range data {

		// getting inside  each dataset folder again
		dsAbs, _ := filepath.Abs(i)
		os.Chdir(dsAbs)

		// reload the meta data
		meta.Restore(sys.Meta())

		logrus.Info("Executing cluster on ", i)

		meta.BioQuant = p.BioQuant

		//clu.GenerateReport(meta)
		qua.RunBioQuantification(meta)

		meta.Serialize()

		// return to the top level directory
		os.Chdir(dir)
	}

	// return to the top level directory
	os.Chdir(dir)
	return meta
}

// Filter executes the Filter, Quantify and Report commands in tandem
func Filter(meta met.Data, p Directives, dir string, data []string) met.Data {

	// this is the virtual home directory where the pipeline is being executed.
	//vHome := meta.Home

	for _, i := range data {

		// getting inside  each dataset folder again
		dsAbs, _ := filepath.Abs(i)
		os.Chdir(dsAbs)

		// reload the meta data
		meta.Restore(sys.Meta())

		// Filter
		if p.Steps.FDRFiltering == "yes" {

			logrus.Info("Executing filter on ", i)
			meta.Filter = p.Filter
			meta.Filter.Tag = p.DatabaseSearch.DecoyTag

			if len(p.Filter.Pex) == 0 {
				meta.Filter.Pex = "interact.pep.xml"
				if p.Steps.PTMLocalization == "yes" {
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

			if p.Steps.IntegratedReports == "yes" && len(p.Filter.Pox) == 0 {
				meta.Filter.Pox = "combined"
			} else if p.Steps.IntegratedReports == "no" && !p.Abacus.Protein && len(p.Filter.Pox) == 0 {
				meta.Filter.Pox = ""
				meta.Filter.Razor = false
				meta.Filter.TwoD = false
				meta.Filter.Seq = false
			}

			meta := fil.Run(meta)

			meta.Serialize()
		}

		// return to the top level directory
		os.Chdir(dir)
	}

	return meta
}

// Report executes the Report commands
func Report(meta met.Data, p Directives, dir string, data []string) met.Data {

	for _, i := range data {

		// getting inside  each dataset folder again
		dsAbs, _ := filepath.Abs(i)
		os.Chdir(dsAbs)

		// reload the meta data
		meta.Restore(sys.Meta())

		// Report
		if p.Steps.IndividualReports == "yes" {

			logrus.Info("Executing report on ", i)

			meta.Report = p.Report

			rep.Run(meta)
			meta.Serialize()
		}

		// return to the top level directory
		os.Chdir(dir)
	}

	return meta
}

// Abacus loads all data and creates the combined protein report
func Abacus(meta met.Data, p Directives, dir string, data []string) met.Data {

	if p.Steps.IntegratedReports == "yes" {

		logrus.Info("Executing abacus")

		// return to the top level directory
		os.Chdir(dir)

		meta.Abacus = p.Abacus
		meta.Abacus.Tag = p.DatabaseSearch.DecoyTag
		meta.Abacus.Picked = p.Filter.Picked
		meta.Abacus.Razor = p.Filter.Razor

		if len(p.LabelQuant.Annot) > 0 {
			meta.Abacus.Labels = true
		}

		aba.Run(meta, data)
	}

	return meta
}

// TMTIntegrator executes TMT-I on all PSM results
func TMTIntegrator(meta met.Data, p Directives, dir string, data []string) met.Data {

	if p.Steps.TMTIntegrator == "yes" {

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
