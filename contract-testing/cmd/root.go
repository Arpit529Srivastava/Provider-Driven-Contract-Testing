package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "contract-testing",
	Short: "Provider-driven contract testing tool",
	Long: `A comprehensive contract testing tool that implements provider-driven 
contract testing to ensure API compatibility between services.

This tool supports:
- OpenAPI schema generation from provider services
- Contract repository management
- Verification against consumer expectations
- Detailed compatibility reports`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags can be defined here
	rootCmd.PersistentFlags().StringP("repository", "r", "./contracts", "Path to contract repository")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
}