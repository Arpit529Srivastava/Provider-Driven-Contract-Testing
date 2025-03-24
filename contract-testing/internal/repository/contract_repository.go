package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ContractRepository struct {
	basePath string
}

func NewContractRepository(basePath string) *ContractRepository {
	return &ContractRepository{
		basePath: basePath,
	}
}

func (r *ContractRepository) SaveProviderSchema(providerName, schemaContent string) error {
	dirPath := filepath.Join(r.basePath, "providers", providerName)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	filePath := filepath.Join(dirPath, "openapi.yaml")
	if err := os.WriteFile(filePath, []byte(schemaContent), 0644); err != nil {
		return fmt.Errorf("failed to write schema: %w", err)
	}
	
	return nil
}

func (r *ContractRepository) GetProviderSchema(providerName string) (string, error) {
	filePath := filepath.Join(r.basePath, "providers", providerName, "openapi.yaml")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read schema: %w", err)
	}
	
	return string(data), nil
}

func (r *ContractRepository) GetConsumerMocks(providerName string) (map[string][]string, error) {
	consumersPath := filepath.Join(r.basePath, "consumers")
	
	consumersDir, err := os.ReadDir(consumersPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read consumers directory: %w", err)
	}
	
	result := make(map[string][]string)
	
	for _, consumerDir := range consumersDir {
		if !consumerDir.IsDir() {
			continue
		}
		
		consumerName := consumerDir.Name()
		mocksPath := filepath.Join(consumersPath, consumerName, "mocks")
		
		mockFiles, err := os.ReadDir(mocksPath)
		if err != nil {
			continue // Skip if no mocks directory
		}
		
		for _, mockFile := range mockFiles {
			if mockFile.IsDir() || !strings.HasSuffix(mockFile.Name(), ".json") {
				continue
			}
			
			mockPath := filepath.Join(mocksPath, mockFile.Name())
			mockData, err := os.ReadFile(mockPath)
			if err != nil {
				continue
			}
			
			var mock map[string]interface{}
			if err := json.Unmarshal(mockData, &mock); err != nil {
				continue
			}
			
			// Check if this mock is for the specified provider
			mockProvider, ok := mock["provider"].(string)
			if !ok || mockProvider != providerName {
				continue
			}
			
			result[consumerName] = append(result[consumerName], mockPath)
		}
	}
	
	return result, nil
}