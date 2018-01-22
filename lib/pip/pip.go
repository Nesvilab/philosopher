package pip

import "github.com/prvst/philosopher/lib/met"

// Directives contains the instructions to run a pipeline
type Directives struct {
	Initialize bool         `yaml:"initialize"`
	Clean      bool         `yaml:"clean"`
	Backup     bool         `yaml:"backup"`
	Analtics   bool         `yaml:"analytics"`
	Commands   Commands     `yaml:"commands"`
	Database   met.Database `yaml:"database"`
}

// Commands struct {
type Commands struct {
	Database string `yaml:"database"`
}
