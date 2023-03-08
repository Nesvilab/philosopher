// Package cmd TMT-Integrator top level command
package cmd

import (
	"os"

	"github.com/Nesvilab/philosopher/lib/ext/tmtintegrator"
	"github.com/Nesvilab/philosopher/lib/met"
	"github.com/Nesvilab/philosopher/lib/msg"
	"github.com/Nesvilab/philosopher/lib/sys"

	"github.com/spf13/cobra"
)

// tmtintegratorCmd represents the tmtintegrator command
var tmtintegratorCmd = &cobra.Command{
	Use:   "tmtintegrator",
	Short: "integrates channel abundances from multiple TMT samples",
	Run: func(cmd *cobra.Command, args []string) {

		m.FunctionInitCheckUp()

		msg.Executing("TMT-Integrator ", Version)

		m := tmtintegrator.Run(m, args)

		m.Serialize()

		// clean tmp
		met.CleanTemp(m.Temp)

		msg.Done()
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "tmtintegrator" {

		m.Restore(sys.Meta())

		tmtintegratorCmd.Flags().StringVarP(&m.TMTIntegrator.JarPath, "path", "", "", "")
		tmtintegratorCmd.Flags().StringVarP(&m.TMTIntegrator.Param, "param", "", "", "")
		tmtintegratorCmd.Flags().IntVarP(&m.TMTIntegrator.Memory, "memory", "", 8, "")
	}

	RootCmd.AddCommand(tmtintegratorCmd)

}
