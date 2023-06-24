package cmd

import (
	"fmt"

	cheek "github.com/datarootsio/cheek/pkg"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var (
	httpPort string
	homeDir  string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cheek",
	Short: "Cheek",
	Long:  `cheek: the pico sized declarative job scheduler`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVar(&httpPort, "port", "8081", "port on which to open the http server for core to ui communication")
	rootCmd.PersistentFlags().StringVar(&homeDir, "homedir", cheek.CheekPath(), fmt.Sprintf("directory in which to save cheek's core & job logs, defaults to '%s'", cheek.CheekPath()))
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("cheek")
	viper.AutomaticEnv()

	// setting both default value AND bind value
	// because tests will often be called without flags
	// being set
	viper.SetDefault("port", "8081")
	if err := viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port")); err != nil {
		fmt.Printf("error binding pflag %s", err)
	}

	if err := viper.BindPFlag("suppressLogs", runCmd.PersistentFlags().Lookup("suppress-logs")); err != nil {
		fmt.Printf("error binding pflag %s", err)
	}

	if err := viper.BindPFlag("logLevel", runCmd.PersistentFlags().Lookup("log-level")); err != nil {
		fmt.Printf("error binding pflag %s", err)
	}

	if err := viper.BindPFlag("pretty", runCmd.PersistentFlags().Lookup("pretty")); err != nil {
		fmt.Printf("error binding pflag %s", err)
	}

	if err := viper.BindPFlag("homedir", rootCmd.PersistentFlags().Lookup("homedir")); err != nil {
		fmt.Printf("error binding pflag %s", err)
	}
}
