package cmd

import (
	"fmt"

	"github.com/Arpit529srivastava/internal/verifier"
	"github.com/spf13/cobra"
)

var (
	schemaPath    string
	mocksDir      string
	providerURL   string
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify provider contracts against consumer mocks",
	Long:  `Validates the provider's implementation against consumer expectations by using the mocks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		validator := verifier.NewValidator(schemaPath, mocksDir, providerURL)
		results, err := validator.Validate()
		if err != nil {
			return err
		}
		
		reporter := verifier.NewReporter(results)
		summary := reporter.GenerateSummary()
		
		fmt.Println(summary)
		return nil
	},
}

func init() {
	verifyCmd.Flags().StringVarP(&schemaPath, "schema", "s", "", "Path to the provider schema (required)")
	verifyCmd.Flags().StringVarP(&mocksDir, "mocks", "m", "", "Directory containing consumer mocks (required)")
	verifyCmd.Flags().StringVarP(&providerURL, "url", "u", "", "Base URL of the provider service (required)")
	
	verifyCmd.MarkFlagRequired("schema")
	verifyCmd.MarkFlagRequired("mocks")
	verifyCmd.MarkFlagRequired("url")
}