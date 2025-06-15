package cmd

import (
	"fmt"

	cheek "github.com/bart6114/cheek/pkg"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "The version of cheek",
	Long:  "The version of cheek",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cheek.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
