package cmd

import (
	"fmt"

	butt "github.com/bart6114/butt/pkg"
	zl "github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var (
	pretty       bool
	surpressLogs bool
	logLevel     string
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run path/to/schedule.yaml",
	Short: "Schedule & run jobs",
	Long:  "Schedule & run jobs",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		butt.RunSchedule(args[0], pretty, httpPort, surpressLogs, logLevel)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVarP(&pretty, "pretty", "p", false, "Output pretty formatted logs to console.")
	runCmd.Flags().BoolVarP(&surpressLogs, "surpress-logs", "s", false, "Do not output logs to stdout, only to file.")
	runCmd.Flags().StringVarP(&logLevel, "log-level", "l", "debug", fmt.Sprintf("Set log level, can be one of %v|%v|%v|%v|%v|%v|%v (only applies to butt specific logging)", zl.LevelTraceValue, zl.LevelDebugValue, zl.LevelInfoValue, zl.LevelWarnValue, zl.LevelErrorValue, zl.LevelFatalValue, zl.LevelPanicValue))
}
