package cmd

import (
	"fmt"

	"github.com/Arpit529srivastava/internal/verifier"
	"github.com/spf13/cobra"
)

var (
	resultsPath string
	reportFormat string
	reportOutput string
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate detailed reports from verification results",
	Long:  `Creates comprehensive reports based on the verification results, highlighting compatibility issues.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		validationResult := &verifier.ValidationResult{}
		reporter := verifier.NewReporter(validationResult)
		err := reporter.GenerateReport(resultsPath, reportFormat, reportOutput)
		if err != nil {
			return err
		}
		
		fmt.Printf("Report generated successfully at: %s\n", reportOutput)
		return nil
	},
}

func init() {
	reportCmd.Flags().StringVarP(&resultsPath, "results", "r", "", "Path to verification results (required)")
	reportCmd.Flags().StringVarP(&reportFormat, "format", "f", "html", "Report format (html, json, markdown)")
	reportCmd.Flags().StringVarP(&reportOutput, "output", "o", "report.html", "Output path for the report")
	
	reportCmd.MarkFlagRequired("results")
}