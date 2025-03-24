package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	providerName     string
	consumerPatterns []string
	reportPath       string
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify provider implementation against consumer expectations",
	Long: `Validates that a provider's implementation meets the expectations of its consumers.
This command downloads consumer mocks from the repository and compares them field-by-field
against the provider's actual implementation.

The verification process checks for:
- Breaking changes (where the provider removes or modifies fields that consumers use)
- Expected status codes
- Schema compatibility

Example:
  contract-testing verify --provider order-service --consumers user-service,notification-service --report ./reports`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Verifying %s against consumer expectations\n", providerName)
		
		// Initialize contract repository
		repoPath, _ := cmd.Flags().GetString("repository")
		contractRepo := repository.NewContractRepository(repoPath)
		
		// Load provider schema
		providerSchema, err := contractRepo.GetProviderSchema(providerName)
		if err != nil {
			return fmt.Errorf("failed to load provider schema: %w", err)
		}
		
		// Initialize verifier
		matcher := verifier.NewMatcher(providerSchema)
		
		// Process each consumer's mocks
		for _, consumer := range consumerPatterns {
			fmt.Printf("Checking compatibility with %s\n", consumer)
			
			// Get consumer mocks
			mocks, err := contractRepo.GetConsumerMocks(consumer)
			if err != nil {
				return fmt.Errorf("failed to load consumer mocks for %s: %w", consumer, err)
			}
			
			// Verify against each mock
			for _, mock := range mocks {
				result, err := matcher.ValidateMock(mock)
				if err != nil {
					fmt.Printf("Error validating mock %s: %v\n", mock.Name, err)
					continue
				}
				
				if result.Compatible {
					fmt.Printf("✅ %s: Compatible\n", mock.Name)
				} else {
					fmt.Printf("❌ %s: Incompatible - %s\n", mock.Name, result.Reason)
				}
			}
		}
		
		// Generate report
		reporter := verifier.NewReporter(reportPath)
		if err := reporter.GenerateReport(providerName, consumerPatterns, matcher.GetResults()); err != nil {
			return fmt.Errorf("failed to generate report: %w", err)
		}
		
		fmt.Printf("Verification complete. Report available at %s/%s-report.html\n", reportPath, providerName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
	
	verifyCmd.Flags().StringVarP(&providerName, "provider", "p", "", "Name of the provider service to verify (required)")
	verifyCmd.Flags().StringSliceVarP(&consumerPatterns, "consumers", "c", []string{}, "Comma-separated list of consumer services to verify against (required)")
	verifyCmd.Flags().StringVarP(&reportPath, "report", "o", "./reports", "Output path for verification reports")
	
	verifyCmd.MarkFlagRequired("provider")
	verifyCmd.MarkFlagRequired("consumers")
}