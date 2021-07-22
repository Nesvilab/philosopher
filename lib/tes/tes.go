package tes

import (
	"os"
	"philosopher/lib/wrk"
)

// SetupTestEnv pre-sets environment directory and meta folder
func SetupTestEnv() {

	os.Chdir("../../test/wrksp/")

	if _, err := os.Stat(".meta"); err != nil {
		if os.IsNotExist(err) {
			wrk.Init("0000", "0000", "")
		}
	}

}

// ShutDowTestEnv pre-sets environment directory and meta folder
func ShutDowTestEnv() {

	wrk.Clean()

}
