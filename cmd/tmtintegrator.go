package cmd

import (
	"os"

	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/ext/tmtintegrator"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/spf13/cobra"
)

// tmtintegratorCmd represents the tmtintegrator command
var tmtintegratorCmd = &cobra.Command{
	Use:   "tmtintegrator",
	Short: "integrates channel abundances from multiple TMT samples",
	Run: func(cmd *cobra.Command, args []string) {

		m.FunctionInitCheckUp()

		err.Executing("TMT-Integrator ", Version)

		m := tmtintegrator.Run(m, args)

		m.Serialize()

		// clean tmp
		met.CleanTemp(m.Temp)

		err.Done()
		return

	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "tmtintegrator" {

		m.Restore(sys.Meta())

		tmtintegratorCmd.Flags().StringVarP(&m.TMTIntegrator.JarPath, "path", "", "", "")
		tmtintegratorCmd.Flags().StringVarP(&m.TMTIntegrator.Param, "param", "", "", "")
		tmtintegratorCmd.Flags().IntVarP(&m.TMTIntegrator.Memmory, "memmory", "", 8, "")
	}

	RootCmd.AddCommand(tmtintegratorCmd)

}
