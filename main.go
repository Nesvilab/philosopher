package main

import "philosopher/cmd"

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

	cmd.Version = version
	cmd.Build = build

	cmd.Execute()

}

// TODO update PRevAA
// TODO update protein fdr before filt
