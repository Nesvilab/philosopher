package pip

import (
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/prvst/philosopher/lib/met"
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
