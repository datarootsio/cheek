package cmd

import (
	"fmt"
	"os"

	cheek "github.com/bart6114/cheek/pkg"
	zl "github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	pretty       bool
	suppressLogs bool
	logLevel     string
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run path/to/schedule.yaml",
	Short: "Schedule & run jobs",
	Long:  "Schedule & run jobs",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c := cheek.NewConfig()
		if err := viper.Unmarshal(&c); err != nil {
			fmt.Println("cannot init configuration, error unmarshalling config: ", err)
			os.Exit(1)
		}
		if err := c.Init(); err != nil {
			fmt.Println("cannot init configuration", err)
			os.Exit(1)
		}

		l := cheek.NewLogger(logLevel, c.DB, cheek.PrettyStdout())
		return cheek.RunSchedule(l, c, args[0])
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().BoolVarP(&pretty, "pretty", "p", true, "Output pretty formatted logs to console.")
	runCmd.PersistentFlags().BoolVarP(&suppressLogs, "suppress-logs", "s", false, "Do not output logs to stdout, only to file.")
	runCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", fmt.Sprintf("Set log level, can be one of %v|%v|%v|%v|%v|%v|%v (only applies to cheek specific logging)", zl.LevelTraceValue, zl.LevelDebugValue, zl.LevelInfoValue, zl.LevelWarnValue, zl.LevelErrorValue, zl.LevelFatalValue, zl.LevelPanicValue))
}
