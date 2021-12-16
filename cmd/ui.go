package cmd

import (
	cheek "github.com/datarootsio/cheek/pkg"
	"github.com/spf13/cobra"
)

var yamlFile string

// uiCmd represents the ui command
var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Terminal User Interface",
	Long: `cheek's UI

By default the UI will communicate with the scheduler to receive the latest state of the schedule specs. This requires the cheek scheduler to be up and running.

Alternatively, the '-schedule' flag allows you to provide a path to the specs YAML. These will be used as backup if the scheduler is not reachable.
`,
	Run: func(cmd *cobra.Command, args []string) {
		cheek.TUI(httpPort, yamlFile)
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
	uiCmd.Flags().StringVarP(&yamlFile, "schedule", "s", "", "Define the schedule file to use if cheek scheduler is not runinn.")
}
