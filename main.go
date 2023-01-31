package main

import (
	"runtime/debug"

	"github.com/Nesvilab/philosopher/cmd"
)

var (
	// Version code
	Version string
	// Build code
	Build string

	version = "dev"
	build   = "build"
	//commit  = "none"
	//date    = "unknown"
)

func main() {
	debug.SetGCPercent(20)
	cmd.Version = version
	cmd.Build = build

	cmd.Execute()

}
