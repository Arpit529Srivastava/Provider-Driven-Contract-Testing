package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "contract-testing",
	Short: "A provider-driven contract testing tool",
	Long: `A tool for provider-driven contract testing that allows 
services to validate their API contracts against consumer expectations.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(verifyCmd)
	rootCmd.AddCommand(reportCmd)
}