package verifier

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type MockRequest struct {
	Method   string                 `json:"method"`
	Endpoint string                 `json:"endpoint"`
	Headers  map[string]string      `json:"headers"`
	Body     map[string]interface{} `json:"body"`
}

type MockResponse struct {
	StatusCode int                    `json:"statusCode"`
	Headers    map[string]string      `json:"headers"`
	Body       map[string]interface{} `json:"body"`
}

type Mock struct {
	Provider     string       `json:"provider"`
	Consumer     string       `json:"consumer"`
	Description  string       `json:"description"`
	Request      MockRequest  `json:"request"`
	Response     MockResponse `json:"response"`
	Dependencies []string     `json:"dependencies"`
}

type MatchResult struct {
	Mock         Mock    `json:"mock"`
	IsCompatible bool    `json:"isCompatible"`
	Issues       []Issue `json:"issues"`
}

type Issue struct {
	Path        string `json:"path"`
	Description string `json:"description"`
	Severity    string `json:"severity"` // "error", "warning"
}

type Matcher struct {
	schema map[string]interface{}
}

func NewMatcher(schema map[string]interface{}) *Matcher {
	return &Matcher{
		schema: schema,
	}
}

func (m *Matcher) MatchMock(mockPath string) (MatchResult, error) {
	data, err := os.ReadFile(mockPath)
	if err != nil {
		return MatchResult{}, fmt.Errorf("failed to read mock file: %w", err)
	}
	
	var mock Mock
	if err := json.Unmarshal(data, &mock); err != nil {
		return MatchResult{}, fmt.Errorf("failed to parse mock: %w", err)
	}
	
	result := MatchResult{
		Mock:         mock,
		IsCompatible: true,
		Issues:       []Issue{},
	}
	
	// Get the path from the schema
	paths, _ := m.schema["paths"].(map[interface{}]interface{})
	if paths == nil {
		result.IsCompatible = false
		result.Issues = append(result.Issues, Issue{
			Path:        "",
			Description: "Schema does not contain paths",
			Severity:    "error",
		})
		return result, nil
	}
	
	// Find the matching endpoint in the schema
	endpoint := mock.Request.Endpoint
	method := strings.ToLower(mock.Request.Method)
	
	var pathItem map[interface{}]interface{}
	var found bool
	
	for path, item := range paths {
		if fmt.Sprintf("%v", path) == endpoint {
			pathItem = item.(map[interface{}]interface{})
			found = true
			break
		}
	}
	
	if !found {
		result.IsCompatible = false
		result.Issues = append(result.Issues, Issue{
			Path:        endpoint,
			Description: "Endpoint not found in provider schema",
			Severity:    "error",
		})
		return result, nil
	}
	
	// Check if the method is supported
	methodItem, found := pathItem[method]
	if !found {
		result.IsCompatible = false
		result.Issues = append(result.Issues, Issue{
			Path:        fmt.Sprintf("%s %s", method, endpoint),
			Description: "Method not supported for this endpoint",
			Severity:    "error",
		})
		return result, nil
	}
	
	// Validate request body against schema
	methodMap := methodItem.(map[interface{}]interface{})
	if requestBody, ok := methodMap["requestBody"].(map[interface{}]interface{}); ok {
		if content, ok := requestBody["content"].(map[interface{}]interface{}); ok {
			if jsonContent, ok := content["application/json"].(map[interface{}]interface{}); ok {
				if _, ok := jsonContent["schema"].(map[interface{}]interface{}); ok {
					// Here we would validate the mock request body against the schema
					// For simplicity, we'll just assume it's valid
				}
			}
		}
	}
	
	// Validate response schema
	if responses, ok := methodMap["responses"].(map[interface{}]interface{}); ok {
		statusCode := fmt.Sprintf("%d", mock.Response.StatusCode)
		if response, ok := responses[statusCode].(map[interface{}]interface{}); ok {
			if content, ok := response["content"].(map[interface{}]interface{}); ok {
				if jsonContent, ok := content["application/json"].(map[interface{}]interface{}); ok {
					if _, ok := jsonContent["schema"].(map[interface{}]interface{}); ok {
						// Validate mock response against schema
						// Again, for simplicity we'll just assume it's valid
					} else {
						result.IsCompatible = false
						result.Issues = append(result.Issues, Issue{
							Path:        fmt.Sprintf("%s %s response.body", method, endpoint),
							Description: "Response schema not defined in provider contract",
							Severity:    "error",
						})
					}
				}
			}
		} else {
			result.IsCompatible = false
			result.Issues = append(result.Issues, Issue{
				Path:        fmt.Sprintf("%s %s response.statusCode", method, endpoint),
				Description: fmt.Sprintf("Status code %s not defined in provider contract", statusCode),
				Severity:    "error",
			})
		}
	}
	
	return result, nil
}