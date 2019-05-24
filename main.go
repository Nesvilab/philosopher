package main

import (
	"github.com/prvst/philosopher/cmd"
)

var (
	// Version code
	Version string
	// Build code
	Build string

	version = "dev"
	build   = "build"
	commit  = "none"
	date    = "unknown"
)

func main() {

	cmd.Version = version
	cmd.Build = build

	cmd.Execute()
}
