package verifier

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"
)

type Reporter struct {
	results *ValidationResult
}

func NewReporter(results *ValidationResult) *Reporter {
	return &Reporter{
		results: results,
	}
}

func (r *Reporter) GenerateSummary() string {
	if r.results == nil {
		return "No validation results available"
	}
	
	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf("Validation Results for Provider: %s\n", r.results.ProviderName))
	sb.WriteString(fmt.Sprintf("Schema: %s\n", r.results.SchemaPath))
	sb.WriteString(fmt.Sprintf("Timestamp: %s\n\n", r.results.Timestamp.Format("2006-01-02 15:04:05")))
	
	if r.results.OverallSuccess {
		sb.WriteString("✅ Overall: All consumer contracts are compatible\n\n")
	} else {
		sb.WriteString("❌ Overall: Some consumer contracts are incompatible\n\n")
	}
	
	// Display results for each consumer
	sb.WriteString("Consumer Results:\n")
	
	for consumer, result := range r.results.ConsumerResults {
		if result.Success {
			sb.WriteString(fmt.Sprintf("  ✅ %s: All expectations met\n", consumer))
		} else {
			sb.WriteString(fmt.Sprintf("  ❌ %s: Incompatibilities found\n", consumer))
			
			// List issues for each mock
			for _, matchResult := range result.MatchResults {
				if !matchResult.IsCompatible {
					sb.WriteString(fmt.Sprintf("    - Mock: %s\n", matchResult.Mock.Description))
					
					for _, issue := range matchResult.Issues {
						sb.WriteString(fmt.Sprintf("      • %s: %s\n", issue.Path, issue.Description))
					}
				}
			}
		}
	}
	
	return sb.String()
}

func (r *Reporter) GenerateReport(resultsPath, format, outputPath string) error {
	// Load results if not provided
	if r.results == nil && resultsPath != "" {
		data, err := os.ReadFile(resultsPath)
		if err != nil {
			return fmt.Errorf("failed to read results file: %w", err)
		}
		
		var results ValidationResult
		if err := json.Unmarshal(data, &results); err != nil {
			return fmt.Errorf("failed to parse results: %w", err)
		}
		
		r.results = &results
	}
	
	if r.results == nil {
		return fmt.Errorf("no validation results available")
	}
	
	switch format {
	case "json":
		return r.generateJSONReport(outputPath)
	case "html":
		return r.generateHTMLReport(outputPath)
	case "markdown":
		return r.generateMarkdownReport(outputPath)
	default:
		return fmt.Errorf("unsupported report format: %s", format)
	}
}

func (r *Reporter) generateJSONReport(outputPath string) error {
	data, err := json.MarshalIndent(r.results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}
	
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write report: %w", err)
	}
	
	return nil
}

func (r *Reporter) generateHTMLReport(outputPath string) error {
	// This is simplified - in a real implementation, you'd have a proper HTML template
	htmlTemplate := `<!DOCTYPE html>
<html>
<head>
    <title>Contract Testing Report - {{.ProviderName}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        h1, h2, h3 { color: #333; }
        .success { color: green; }
        .failure { color: red; }
        .issue { margin-left: 20px; }
        .summary { margin: 20px 0; padding: 10px; background-color: #f8f8f8; }
    </style>
</head>
<body>
    <h1>Contract Testing Report</h1>
    
    <div class="summary">
        <h2>Summary</h2>
        <p><strong>Provider:</strong> {{.ProviderName}}</p>
        <p><strong>Schema:</strong> {{.SchemaPath}}</p>
        <p><strong>Timestamp:</strong> {{.Timestamp}}</p>
        {{if .OverallSuccess}}
            <p class="success"><strong>Overall Status:</strong> All consumer contracts are compatible</p>
        {{else}}
            <p class="failure"><strong>Overall Status:</strong> Some consumer contracts are incompatible</p>
        {{end}}
    </div>
    
    <h2>Consumer Results</h2>
    {{range $consumer, $result := .ConsumerResults}}
        <h3>{{$consumer}}</h3>
        {{if $result.Success}}
            <p class="success">✅ All expectations met</p>
        {{else}}
            <p class="failure">❌ Incompatibilities found</p>
            {{range $mock := $result.MatchResults}}
                {{if not $mock.IsCompatible}}
                    <div class="issue">
                        <h4>{{$mock.Mock.Description}}</h4>
                        <ul>
                            {{range $issue := $mock.Issues}}
                                <li>
                                    <strong>{{$issue.Path}}:</strong> {{$issue.Description}}
                                    ({{$issue.Severity}})
                                </li>
                            {{end}}
                        </ul>
                    </div>
                {{end}}
            {{en{{end}}
            {{end}}
        {{end}}
    {{end}}
</body>
</html>`

	tmpl, err := template.New("report").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}
	
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create report file: %w", err)
	}
	defer file.Close()
	
	if err := tmpl.Execute(file, r.results); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	
	return nil
}

func (r *Reporter) generateMarkdownReport(outputPath string) error {
	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf("# Contract Testing Report - %s\n\n", r.results.ProviderName))
	
	sb.WriteString("## Summary\n\n")
	sb.WriteString(fmt.Sprintf("- **Provider:** %s\n", r.results.ProviderName))
	sb.WriteString(fmt.Sprintf("- **Schema:** %s\n", r.results.SchemaPath))
	sb.WriteString(fmt.Sprintf("- **Timestamp:** %s\n", r.results.Timestamp.Format("2006-01-02 15:04:05")))
	
	if r.results.OverallSuccess {
		sb.WriteString("\n**Overall Status:** ✅ All consumer contracts are compatible\n\n")
	} else {
		sb.WriteString("\n**Overall Status:** ❌ Some consumer contracts are incompatible\n\n")
	}
	
	sb.WriteString("## Consumer Results\n\n")
	
	for consumer, result := range r.results.ConsumerResults {
		sb.WriteString(fmt.Sprintf("### %s\n\n", consumer))
		
		if result.Success {
			sb.WriteString("✅ All expectations met\n\n")
		} else {
			sb.WriteString("❌ Incompatibilities found\n\n")
			
			for _, matchResult := range result.MatchResults {
				if !matchResult.IsCompatible {
					sb.WriteString(fmt.Sprintf("#### %s\n\n", matchResult.Mock.Description))
					
					for _, issue := range matchResult.Issues {
						sb.WriteString(fmt.Sprintf("- **%s:** %s (%s)\n", issue.Path, issue.Description, issue.Severity))
					}
					
					sb.WriteString("\n")
				}
			}
		}
	}
	
	if err := os.WriteFile(outputPath, []byte(sb.String()), 0644); err != nil {
		return fmt.Errorf("failed to write report: %w", err)
	}
	
	return nil
}