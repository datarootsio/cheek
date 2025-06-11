package cmd

import (
	cheek "github.com/bart6114/cheek/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// triggerCmd represents the trigger command
var triggerCmd = &cobra.Command{
	Use:   "trigger {schedule.yaml} {job_name}",
	Short: "Trigger a specific job by name",
	Long: `Trigger a specific job by name

The name should be defined in your schedule specs. Usage:
'cheek trigger my_schedule.yaml my_job'
`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c := cheek.NewConfig()
		if err := viper.Unmarshal(&c); err != nil {
			return err
		}
		if err := c.Init(); err != nil {
			return err
		}

		l := cheek.NewLogger(logLevel, c.DB, cheek.PrettyStdout())
		_, err := cheek.RunJob(l, c, args[0], args[1])
		return err
	},
}

func init() {
	rootCmd.AddCommand(triggerCmd)
}
