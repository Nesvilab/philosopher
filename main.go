package main

import (
	"philosopher/cmd"
	"runtime/debug"
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
