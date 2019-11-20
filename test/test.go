package test

import (
	"os"
	"philosopher/lib/wrk"
)

// SetupTestEnv pre-sets environment directory and meta folder
func SetupTestEnv() {

	os.Chdir("../../test/wrksp/")
	wrk.Init("0000", "0000")

	return
}

// ShutDowTestEnv pre-sets environment directory and meta folder
func ShutDowTestEnv() {

	wrk.Clean()

	return
}
