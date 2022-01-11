package cmd

import (
	"fmt"
	"os"

	cheek "github.com/datarootsio/cheek/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// triggerCmd represents the trigger command
var triggerCmd = &cobra.Command{
	Use:   "trigger",
	Short: "Trigger a specific job by name",
	Long: `Trigger a specific job by name

The name should be defined in your schedule specs. Usage:
'cheek trigger my_schedule.yaml my_job'
`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		c := cheek.NewConfig()
		if err := viper.Unmarshal(&c); err != nil {
			fmt.Println("cannot init configuration")
			os.Exit(1)
		}
		l := cheek.NewLogger(logLevel, cheek.PrettyStdout())
		cheek.RunJob(l, c, args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(triggerCmd)
}
