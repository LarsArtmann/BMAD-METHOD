package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "template-health-endpoint",
	Short: "TypeSpec-driven health endpoint template generator",
	Long: `A comprehensive template generator for health endpoint APIs using TypeSpec.

This tool generates production-ready health endpoint implementations with:
- TypeSpec-first API definitions
- Go server implementations with OpenTelemetry integration
- TypeScript client SDKs
- Kubernetes deployment configurations
- Progressive complexity tiers (Basic → Intermediate → Advanced → Enterprise)

Examples:
  # Generate a basic health endpoint service
  template-health-endpoint generate --tier basic --name my-service

  # Generate an advanced service with full observability
  template-health-endpoint generate --tier advanced --name my-service --features opentelemetry,cloudevents

  # Validate TypeSpec schemas
  template-health-endpoint validate --schemas template-health/schemas/

For more information, visit: https://github.com/LarsArtmann/BMAD-METHOD`,
	Version: "1.0.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.template-health-endpoint.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".template-health-endpoint" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".template-health-endpoint")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil && verbose {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
