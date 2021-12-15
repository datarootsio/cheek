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
	Long: `Terminal User Interface
	
Nothing cheek the UI.`,
	Run: func(cmd *cobra.Command, args []string) {
		cheek.TUI(httpPort, yamlFile)
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
	uiCmd.Flags().StringVarP(&yamlFile, "schedule", "s", "", "Define the schedule file to display.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
