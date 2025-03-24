package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/y/contract-testing/internal/schema
)

var (
	serviceName string
	serviceURL  string
	outputPath  string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate OpenAPI schema from a provider service",
	Long: `Analyzes a provider service and generates an OpenAPI schema that represents
the contract. This schema includes all endpoints, data structures, request parameters,
and response formats that the provider exposes.

Example:
  contract-testing generate --service order-service --url http://localhost:8080 --output ./contracts/providers/order-service`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Generating schema for %s at %s\n", serviceName, serviceURL)
		
		generator := schema.NewGenerator(serviceName, serviceURL)
		schema, err := generator.GenerateSchema()
		if err != nil {
			return fmt.Errorf("failed to generate schema: %w", err)
		}
		
		if err := generator.SaveSchema(schema, outputPath); err != nil {
			return fmt.Errorf("failed to save schema: %w", err)
		}
		
		fmt.Printf("Successfully generated schema for %s at %s/openapi.yaml\n", serviceName, outputPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	
	generateCmd.Flags().StringVarP(&serviceName, "service", "s", "", "Name of the provider service (required)")
	generateCmd.Flags().StringVarP(&serviceURL, "url", "u", "", "URL of the provider service (required)")
	generateCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output path for the generated schema (required)")
	
	generateCmd.MarkFlagRequired("service")
	generateCmd.MarkFlagRequired("url")
	generateCmd.MarkFlagRequired("output")
}