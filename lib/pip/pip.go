package pip

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/prvst/philosopher/lib/aba"
	"github.com/prvst/philosopher/lib/clu"
	"github.com/prvst/philosopher/lib/dat"
	"github.com/prvst/philosopher/lib/ext/comet"
	"github.com/prvst/philosopher/lib/ext/fragger"
	"github.com/prvst/philosopher/lib/ext/peptideprophet"
	"github.com/prvst/philosopher/lib/ext/proteinprophet"
	"github.com/prvst/philosopher/lib/ext/ptmprophet"
	"github.com/prvst/philosopher/lib/fil"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/qua"
	"github.com/prvst/philosopher/lib/rep"
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
}

func Run(m met.Data, p Directives, dir, Version, Build string, args []string) met.Data {

	// For each dataset ...
	for _, i := range args {

		logrus.Info("Executing the pipeline on ", i)

		// getting inside de the dataset folder
		dsAbs, _ := filepath.Abs(i)
		os.Chdir(dsAbs)

		// Workspace
		wrk.Run(Version, Build, false, false, true)

		// reload the meta data
		m.Restore(sys.Meta())

		// Database
		if p.Commands.Database == "yes" {
			m.Database = p.Database
			dat.Run(m)

			m.Serialize()
			met.CleanTemp(m.Temp)
		}

		if p.Commands.Comet == "yes" && p.Commands.MSFragger == "yes" {
			logrus.Fatal("You can only specify one search engine at a time")
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

			m.Serialize()
			met.CleanTemp(m.Temp)
		}

		// MSFragger
		if p.Commands.MSFragger == "yes" {
			m.MSFragger = p.MSFragger
			gobExt := fmt.Sprintf("*.%s", p.MSFragger.RawExtension)
			files, e := filepath.Glob(gobExt)
			if e != nil {
				logrus.Fatal(e)
			}
			fragger.Run(m, files)

			m.Serialize()
			met.CleanTemp(m.Temp)
		}

		// PeptideProphet
		if p.Commands.PeptideProphet == "yes" {
			logrus.Info("Executing PeptideProphet on ", i)
			m.PeptideProphet = p.PeptideProphet
			m.PeptideProphet.Output = "interact"
			m.PeptideProphet.Combine = true
			gobExt := fmt.Sprintf("*.%s", p.PeptideProphet.FileExtension)
			files, e := filepath.Glob(gobExt)
			if e != nil {
				logrus.Fatal(e)
			}
			peptideprophet.Run(m, files)

			m.Serialize()
			met.CleanTemp(m.Temp)
		}

		if p.Commands.PTMProphet == "yes" {
			logrus.Info("Executing PTMProphet on ", i)
			m.PTMProphet = p.PTMProphet
			var files []string
			files = append(files, "interact.pep.xml")
			m.PTMProphet.InputFiles = files
			ptmprophet.Run(m, files)

			m.Serialize()
			met.CleanTemp(m.Temp)
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
			met.CleanTemp(m.Temp)
		}

		// return to the top level directory
		os.Chdir(dir)
	}

	// Abacus
	var combinedProtXML string
	if p.Commands.Abacus == "yes" {
		logrus.Info("Creating combined protein inference")
		// return to the top level directory
		os.Chdir(dir)
		m.ProteinProphet = p.ProteinProphet
		m.ProteinProphet.Output = "combined"
		var files []string
		for _, j := range args {
			fqn := fmt.Sprintf("%s%sinteract.pep.xml", j, string(filepath.Separator))
			if p.Commands.PTMProphet == "yes" {
				fqn = fmt.Sprintf("%s%sinteract.mod.pep.xml", j, string(filepath.Separator))
			}
			fqn, _ = filepath.Abs(fqn)
			files = append(files, fqn)
		}
		proteinprophet.Run(m, files)
		combinedProtXML = fmt.Sprintf("%s%scombined.prot.xml", m.Temp, string(filepath.Separator))

		// copy to work directory
		sys.CopyFile(combinedProtXML, filepath.Base(combinedProtXML))

		m.Serialize()
		met.CleanTemp(m.Temp)
	}

	for _, i := range args {

		// getting inside  each dataset folder again
		dsAbs, _ := filepath.Abs(i)
		os.Chdir(dsAbs)

		// Filter
		if p.Commands.Filter == "yes" {
			logrus.Info("Executing filter on ", i)
			m.Filter = p.Filter
			m.Filter.Pex = "interact.pep.xml"
			if p.Commands.PTMProphet == "yes" {
				m.Filter.Pex = "interact.mod.pep.xml"
			}
			if p.Commands.ProteinProphet == "yes" {
				m.Filter.Pox = "interact.prot.xml"
			}
			if p.Commands.Abacus == "yes" {
				m.Filter.Pox = combinedProtXML
			}
			m, e := fil.Run(m)
			if e != nil {
				logrus.Fatal(e.Error())
			}

			m.Serialize()
			met.CleanTemp(m.Temp)
		}

		// getting inside de the dataset folder again
		os.Chdir(dsAbs)

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
			met.CleanTemp(m.Temp)
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
			met.CleanTemp(m.Temp)
		}

		// Report
		if p.Commands.Report == "yes" {
			logrus.Info("Executing report on ", i)
			rep.Run(m)

			m.Serialize()
			met.CleanTemp(m.Temp)
		}

		// Cluster
		if p.Commands.Cluster == "yes" {
			logrus.Info("Executing cluster on ", i)
			m.Cluster = p.Cluster
			clu.GenerateReport(m)

			m.Serialize()
			met.CleanTemp(m.Temp)
		}

		// return to the top level directory
		os.Chdir(dir)
	}

	// Abacus
	if p.Commands.Abacus == "yes" {
		logrus.Info("Executing abacus")
		// return to the top level directory
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
			wrk.Run(Version, Build, true, false, false)
		}

		// Clean
		if p.Clean == true {
			wrk.Run(Version, Build, false, true, false)
		}

	}

	return m
}

func ParallelRun(m met.Data, p Directives, dir, Version, Build string, args []string) error {

	var metArray []met.Data

	// For each dataset ...
	for _, i := range args {
		logrus.Info("Executing the pipeline on ", i)

		// getting inside de the dataset folder
		os.Chdir(dir)
		dsAbs, _ := filepath.Abs(i)
		os.Chdir(dsAbs)

		// Workspace
		wrk.Run(Version, Build, false, false, true)

		// reload the meta data
		m.Restore(sys.Meta())
		metArray = append(metArray, m)
	}

	fmt.Println(metArray)

	// // Database
	// if p.Commands.Database == "yes" {
	// 	m.Database = p.Database
	// 	dat.Run(m)
	//
	// 	m.Serialize()
	// 	met.CleanTemp(m.Temp)
	// }

	return nil
}

// func ParallelRun(m met.Data, p Directives, dir, Version, Build string, args []string) met.Data {
//
// 	var wg sync.WaitGroup
// 	wg.Add(len(args))
// 	go func(args []string) {
// 		for _, arg := range args {
// 			// getting inside de the dataset folder
// 			os.Chdir(dir)
// 			dsAbs, _ := filepath.Abs(arg)
// 			os.Chdir(dsAbs)
//
// 			// Workspace
// 			wrk.Run(Version, Build, false, false, true)
//
// 			// reload the meta data
// 			m.Restore(sys.Meta())
//
// 			// Database
// 			if p.Commands.Database == "yes" {
// 				m.Database = p.Database
// 				dat.Run(m)
// 				m.Serialize()
// 			}
//
// 			_ = m
// 			wg.Done()
// 		}
// 	}(args)
// 	wg.Wait()
//
// 	// returning to the pipeline root directory
// 	os.Chdir(dir)
//
// 	// PeptideProphet
// 	if p.Commands.PeptideProphet == "yes" {
//
// 		// For each dataset ...
// 		var wg sync.WaitGroup
// 		for _, i := range args {
// 			wg.Add(1)
//
// 			os.Chdir(dir)
//
// 			// getting inside de the dataset folder
// 			dsAbs, _ := filepath.Abs(i)
// 			os.Chdir(dsAbs)
// 			fmt.Println(dsAbs)
//
// 			// reload the meta data
// 			m.Restore(sys.Meta())
// 			fmt.Println(m.ProjectName)
//
// 			logrus.Info("Executing PeptideProphet on ", i)
// 			m.PeptideProphet = p.PeptideProphet
// 			m.PeptideProphet.Output = "interact"
// 			m.PeptideProphet.Combine = true
// 			gobExt := fmt.Sprintf("*.%s", p.PeptideProphet.FileExtension)
// 			files, e := filepath.Glob(gobExt)
// 			if e != nil {
// 				logrus.Fatal(e)
// 			}
//
// 			go func(m met.Data, files []string) {
// 				defer wg.Done()
// 				//fmt.Println(m.ProjectName)
// 				//peptideprophet.Run(m, files)
// 				m.Serialize()
// 				met.CleanTemp(m.Temp)
// 			}(m, files)
//
// 		}
// 		wg.Wait()
//
// 	}
//
// 	fmt.Println("finished all")
//
// 	return m
// }

// func ParallelRun(m met.Data, args []string, dir, Version, Build string) met.Data {
//
// 	file, _ := filepath.Abs(m.Pipeline.Directives)
//
// 	y, e := ioutil.ReadFile(file)
// 	if e != nil {
// 		log.Fatal(e)
// 	}
//
// 	var p Directives
// 	e = yaml.Unmarshal(y, &p)
// 	if e != nil {
// 		logrus.Fatal(e)
// 	}
//
// 	if len(args) < 1 {
// 		logrus.Fatal("You need to provide at least one dataset for the analysis.")
// 	} else if p.Commands.Abacus == "true" && len(args) < 2 {
// 		logrus.Fatal("You need to provide at least two datasets for the abacus integrative analysis.")
// 	}
//
// 	//var wg sync.WaitGroup
// 	//wg.Add(len(args))
// 	// go func(args []string) {
// 	//
// 	// 	for _, arg := range args {
// 	//
// 	// 		fmt.Println(arg)
// 	//
// 	// 		// getting inside de the dataset folder
// 	// 		os.Chdir(dir)
// 	// 		dsAbs, _ := filepath.Abs(arg)
// 	// 		os.Chdir(dsAbs)
// 	//
// 	// 		fmt.Println(dsAbs)
// 	//
// 	// 		// Workspace
// 	// 		wrk.Run(Version, Build, false, false, true)
// 	//
// 	// 		// reload the meta data
// 	// 		m.Restore(sys.Meta())
// 	//
// 	// 		// Database
// 	// 		if p.Commands.Database == "yes" {
// 	// 			m.Database = p.Database
// 	// 			dat.Run(m)
// 	// 			m.Serialize()
// 	// 		}
// 	//
// 	// 		wg.Done()
// 	// 	}
// 	//
// 	// }(args)
// 	// wg.Wait()
//
// 	// 		// getting inside de the dataset folder
// 	os.Chdir(dir)
//
// 	var wgPep sync.WaitGroup
// 	wgPep.Add(len(args))
// 	for _, arg := range args {
//
// 		// getting inside de the dataset folder
// 		os.Chdir(dir)
// 		pwd, err := os.Getwd()
// 		if err != nil {
// 			fmt.Println(err)
// 			os.Exit(1)
// 		}
// 		fmt.Println(pwd)
//
// 		go func(arg string) {
//
// 			os.Chdir(m.Home)
// 			dsAbs, _ := filepath.Abs(arg)
// 			os.Chdir(dsAbs)
// 			fmt.Println(dsAbs)
// 			// reload the meta data
// 			m.Restore(sys.Meta())
//
// 			// pwd, err := os.Getwd()
// 			// if err != nil {
// 			// 	fmt.Println(err)
// 			// 	os.Exit(1)
// 			// }
// 			// fmt.Println(pwd)
//
// 			logrus.Info("Executing PeptideProphet on ", arg)
// 			// m.PeptideProphet = p.PeptideProphet
// 			// m.PeptideProphet.Output = "interact"
// 			// m.PeptideProphet.Combine = true
// 			// gobExt := fmt.Sprintf("*.%s", p.PeptideProphet.FileExtension)
// 			// files, e := filepath.Glob(gobExt)
// 			// if e != nil {
// 			// 	logrus.Fatal(e)
// 			// }
// 			// peptideprophet.Run(m, files)
//
// 			m.Serialize()
// 			met.CleanTemp(m.Temp)
//
// 			wgPep.Done()
// 		}(arg)
// 	}
// 	wgPep.Wait()
//
// 	return m
// }

// DeployParameterFile ...
func DeployParameterFile(temp string) (string, error) {

	file := temp + string(filepath.Separator) + "philosopher.yaml"

	param, err := Asset("philosopher.yaml")
	if err != nil {
		return file, errors.New("Cannot deploy Comet parameter file")
	}

	err = ioutil.WriteFile(file, param, 0644)
	if err != nil {
		return file, errors.New("Cannot deploy pipeline parameter file")
	}

	return file, nil
}
