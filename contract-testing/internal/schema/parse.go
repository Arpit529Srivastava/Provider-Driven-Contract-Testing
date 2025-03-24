package schema

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Parser struct {
	schemaPath string
}

func NewParser(schemaPath string) *Parser {
	return &Parser{
		schemaPath: schemaPath,
	}
}

func (p *Parser) Parse() (map[string]interface{}, error) {
	data, err := os.ReadFile(p.schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}
	
	var schema map[string]interface{}
	if err := yaml.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("failed to parse schema: %w", err)
	}
	
	return schema, nil
}

func (p *Parser) GetEndpoints() ([]string, error) {
	schema, err := p.Parse()
	if err != nil {
		return nil, err
	}
	
	paths, ok := schema["paths"].(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid schema format: paths not found or not a map")
	}
	
	var endpoints []string
	for path := range paths {
		endpoints = append(endpoints, fmt.Sprintf("%v", path))
	}
	
	return endpoints, nil
}