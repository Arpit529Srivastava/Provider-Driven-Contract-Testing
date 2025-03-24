package cmd

import (
	"fmt"

	"github.com/Arpit529srivastava/internal/schema"
	"github.com/spf13/cobra"
)

var (
	providerName string
	baseURL      string
	outputPath   string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate OpenAPI schema from provider service",
	Long:  `Analyzes the provider service API and generates an OpenAPI schema that represents the contract.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		generator := schema.NewGenerator(providerName, baseURL)
		err := generator.GenerateSchema(outputPath)
		if err != nil {
			return err
		}
		fmt.Printf("Schema generated successfully at: %s\n", outputPath)
		return nil
	},
}

func init() {
	generateCmd.Flags().StringVarP(&providerName, "provider", "p", "", "Name of the provider service (required)")
	generateCmd.Flags().StringVarP(&baseURL, "url", "u", "", "Base URL of the provider service (required)")
	generateCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output path for the generated schema (required)")
	
	generateCmd.MarkFlagRequired("provider")
	generateCmd.MarkFlagRequired("url")
	generateCmd.MarkFlagRequired("output")
}