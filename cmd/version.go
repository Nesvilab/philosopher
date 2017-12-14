// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Philosopher version",
	Run: func(cmd *cobra.Command, args []string) {
		t := time.Now()
		fmt.Printf("Version: %d%02d%02d\n", t.Year(), t.Month(), t.Day())
		fmt.Printf("Build: %d%02d%02d.%02d%02d\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute())
		return
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "version" {
	}

	RootCmd.AddCommand(versionCmd)
}
