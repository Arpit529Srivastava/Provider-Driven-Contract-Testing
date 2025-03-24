package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// MockRepository handles operations for consumer mocks
type MockRepository struct {
	basePath string
}

// MockData represents a recorded consumer interaction
type MockData struct {
	Name         string                 `json:"name"`
	Method       string                 `json:"method"`
	Path         string                 `json:"path"`
	Request      map[string]interface{} `json:"request"`
	Response     map[string]interface{} `json:"response"`
	Service      string                 `json:"service"`
	Expectations map[string]interface{} `json:"expectations,omitempty"`
}

// NewMockRepository creates a new mock repository instance
func NewMockRepository(basePath string) (*MockRepository, error) {
	if basePath == "" {
		return nil, fmt.Errorf("base path cannot be empty")
	}

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create mock repository directory: %w", err)
	}

	return &MockRepository{
		basePath: basePath,
	}, nil
}

// GetAllConsumerMocks retrieves all mocks for a specific provider
func (r *MockRepository) GetAllConsumerMocks(providerName string) ([]MockData, error) {
	var allMocks []MockData
	consumersPath := filepath.Join(r.basePath, "consumers")

	// Get all consumer directories
	consumerDirs, err := ioutil.ReadDir(consumersPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read consumers directory: %w", err)
	}

	// For each consumer, look for mocks related to the specified provider
	for _, dir := range consumerDirs {
		if !dir.IsDir() {
			continue
		}

		consumerName := dir.Name()
		mocksPath := filepath.Join(consumersPath, consumerName, "mocks")
		
		// Check if mocks directory exists
		if _, err := os.Stat(mocksPath); os.IsNotExist(err) {
			continue
		}

		// Read mock files
		mockFiles, err := ioutil.ReadDir(mocksPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read mocks directory for consumer %s: %w", consumerName, err)
		}

		for _, file := range mockFiles {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
				continue
			}

			// Read and parse mock file
			mockPath := filepath.Join(mocksPath, file.Name())
			mockData, err := readMockFile(mockPath, consumerName)
			if err != nil {
				return nil, fmt.Errorf("failed to read mock file %s: %w", mockPath, err)
			}

			// Only include mocks related to the specified provider
			// This could be determined by path, service name, or other criteria
			if isProviderMock(mockData, providerName) {
				allMocks = append(allMocks, mockData)
			}
		}
	}

	return allMocks, nil
}

// StoreMock saves a mock to the repository
func (r *MockRepository) StoreMock(consumerName, mockName string, mockData MockData) error {
	// Create directory path
	mockDir := filepath.Join(r.basePath, "consumers", consumerName, "mocks")
	if err := os.MkdirAll(mockDir, 0755); err != nil {
		return fmt.Errorf("failed to create mock directory: %w", err)
	}

	// Prepare file path
	fileName := fmt.Sprintf("%s.json", mockName)
	filePath := filepath.Join(mockDir, fileName)

	// Set service name
	mockData.Service = consumerName

	// Convert to JSON
	jsonData, err := json.MarshalIndent(mockData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal mock data: %w", err)
	}

	// Write to file
	if err := ioutil.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write mock file: %w", err)
	}

	return nil
}

// GetMockForEndpoint retrieves mocks for a specific endpoint
func (r *MockRepository) GetMockForEndpoint(providerName, path, method string) ([]MockData, error) {
	allMocks, err := r.GetAllConsumerMocks(providerName)
	if err != nil {
		return nil, err
	}

	var endpointMocks []MockData
	for _, mock := range allMocks {
		if mock.Path == path && strings.EqualFold(mock.Method, method) {
			endpointMocks = append(endpointMocks, mock)
		}
	}

	return endpointMocks, nil
}

// Helper function to read a mock file
func readMockFile(path, consumerName string) (MockData, error) {
	var mockData MockData

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return mockData, err
	}

	if err := json.Unmarshal(data, &mockData); err != nil {
		return mockData, err
	}

	// Ensure service name is set
	if mockData.Service == "" {
		mockData.Service = consumerName
	}

	return mockData, nil
}

// Helper function to determine if a mock is related to a specific provider
func isProviderMock(mock MockData, providerName string) bool {
	// This is a simplified example. In reality, you might have more complex logic
	// to determine if a mock relates to a specific provider.
	
	// For example, check if the path contains the provider name
	if strings.Contains(mock.Path, providerName) {
		return true
	}
	
	// Or check for specific provider domains/paths
	if providerName == "order-service" && strings.HasPrefix(mock.Path, "/orders") {
		return true
	}
	
	return false
}