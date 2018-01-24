package pip

import "github.com/prvst/philosopher/lib/met"

// Directives contains the instructions to run a pipeline
type Directives struct {
	Initialize     bool               `yaml:"initialize"`
	Clean          bool               `yaml:"clean"`
	Backup         bool               `yaml:"backup"`
	Analtics       bool               `yaml:"analytics"`
	Commands       Commands           `yaml:"commands"`
	Database       met.Database       `yaml:"database"`
	Comet          met.Comet          `yaml:"comet"`
	PeptideProphet met.PeptideProphet `yaml:"peptideprophet"`
	ProteinProphet met.ProteinProphet `yaml:"proteinrophet"`
}

// Commands struct {
type Commands struct {
	Database       string `yaml:"database"`
	Comet          string `yaml:"comet"`
	PeptideProphet string `yaml:"peptideprophet"`
	ProteinProphet string `yaml:"proteinprophet"`
}
