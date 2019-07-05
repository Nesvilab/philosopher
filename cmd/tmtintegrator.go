package cmd

import (
	"os"

	"github.com/prvst/philosopher/lib/ext/tmtintegrator"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// tmtintegratorCmd represents the tmtintegrator command
var tmtintegratorCmd = &cobra.Command{
	Use:   "tmtintegrator",
	Short: "integrates channel abundances from multiple TMT samples",
	Run: func(cmd *cobra.Command, args []string) {

		e := m.FunctionInitCheckUp()
		if e != nil {
			logrus.Fatal(e)
		}

		logrus.Info("Executing TMT-Integrator ", Version)

		m, ee := tmtintegrator.Run(m, args)
		if ee != nil {
			logrus.Warn(ee)
		}
		m.Serialize()

		// clean tmp
		met.CleanTemp(m.Temp)

		logrus.Info("Done")
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
