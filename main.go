package main

import (
	"github.com/prvst/philosopher/cmd"
)

var (
	// Version code
	Version string
	// Build code
	Build string
)

func main() {

	cmd.Version = Version
	cmd.Build = Build

	cmd.Execute()
}
