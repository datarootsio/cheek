package cmd

import (
	"fmt"

	cheek "github.com/datarootsio/cheek/pkg"
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
	Run: func(cmd *cobra.Command, args []string) {
		c := cheek.NewConfig()
		viper.Unmarshal(&c)
		l := cheek.NewLogger(pretty, logLevel)
		cheek.RunSchedule(l, c, args[0])
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().BoolVarP(&pretty, "pretty", "p", true, "Output pretty formatted logs to console.")
	runCmd.PersistentFlags().BoolVarP(&suppressLogs, "suppress-logs", "s", false, "Do not output logs to stdout, only to file.")
	runCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", fmt.Sprintf("Set log level, can be one of %v|%v|%v|%v|%v|%v|%v (only applies to cheek specific logging)", zl.LevelTraceValue, zl.LevelDebugValue, zl.LevelInfoValue, zl.LevelWarnValue, zl.LevelErrorValue, zl.LevelFatalValue, zl.LevelPanicValue))
}
