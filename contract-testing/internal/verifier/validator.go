package verifier

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Arpit529srivastava/internal/schema"
)

type ValidationResult struct {
	ProviderName    string                  `json:"providerName"`
	SchemaPath      string                  `json:"schemaPath"`
	Timestamp       time.Time               `json:"timestamp"`
	ConsumerResults map[string]ConsumerResult `json:"consumerResults"`
	OverallSuccess  bool                    `json:"overallSuccess"`
}

type ConsumerResult struct {
	ConsumerName string        `json:"consumerName"`
	MatchResults []MatchResult `json:"matchResults"`
	Success      bool          `json:"success"`
}

type Validator struct {
	schemaPath  string
	mocksDir    string
	providerURL string
}

func NewValidator(schemaPath, mocksDir, providerURL string) *Validator {
	return &Validator{
		schemaPath:  schemaPath,
		mocksDir:    mocksDir,
		providerURL: providerURL,
	}
}

func (v *Validator) Validate() (*ValidationResult, error) {
	// Parse the schema
	parser := schema.NewParser(v.schemaPath)
	schemaData, err := parser.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse schema: %w", err)
	}
	
	// Get provider name from schema path
	providerName := filepath.Base(filepath.Dir(v.schemaPath))
	
	// Initialize the matcher
	matcher := NewMatcher(schemaData)
	
	// Initialize the result
	result := &ValidationResult{
		ProviderName:    providerName,
		SchemaPath:      v.schemaPath,
		Timestamp:       time.Now(),
		ConsumerResults: make(map[string]ConsumerResult),
		OverallSuccess:  true,
	}
	
	// Walk through mocks directory
	err = filepath.Walk(v.mocksDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".json") {
			return nil
		}
		
		// Read the mock to check if it's for this provider
		data, err := os.ReadFile(path)
		if err != nil {
			return nil // Skip this file
		}
		
		var mock struct {
			Provider string `json:"provider"`
			Consumer string `json:"consumer"`
		}
		
		if err := json.Unmarshal(data, &mock); err != nil {
			return nil // Skip this file
		}
		
		if mock.Provider != providerName {
			return nil // Skip mocks for other providers
		}
		
		// Match the mock against the schema
		matchResult, err := matcher.MatchMock(path)
		if err != nil {
			// Log the error but continue with other mocks
			fmt.Printf("Error matching mock %s: %v\n", path, err)
			return nil
		}
		
		// Update consumer results
		consumerResult, exists := result.ConsumerResults[mock.Consumer]
		if !exists {
			consumerResult = ConsumerResult{
				ConsumerName: mock.Consumer,
				MatchResults: []MatchResult{},
				Success:      true,
			}
		}
		
		consumerResult.MatchResults = append(consumerResult.MatchResults, matchResult)
		
		// Update success flag
		if !matchResult.IsCompatible {
			consumerResult.Success = false
			result.OverallSuccess = false
		}
		
		result.ConsumerResults[mock.Consumer] = consumerResult
		
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to process mocks: %w", err)
	}
	
	return result, nil
}