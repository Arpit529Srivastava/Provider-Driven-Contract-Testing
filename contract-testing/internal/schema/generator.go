package schema

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Generator handles the generation of OpenAPI schemas from provider services
type Generator struct {
	ServiceName string
	ServiceURL  string
}

// NewGenerator creates a new schema generator instance
func NewGenerator(serviceName, serviceURL string) *Generator {
	return &Generator{
		ServiceName: serviceName,
		ServiceURL:  serviceURL,
	}
}

// GenerateSchema analyzes a provider service and generates an OpenAPI schema
func (g *Generator) GenerateSchema() (map[string]interface{}, error) {
	// In a real implementation, this would use reflection, annotations, or swagger tools
	// to automatically generate the OpenAPI schema. For this example, we'll simulate
	// the generation for an order service.
	
	// Try to fetch OpenAPI spec from the service directly if it exposes a swagger endpoint
	resp, err := http.Get(fmt.Sprintf("%s/swagger/doc.json", g.ServiceURL))
	if err == nil && resp.StatusCode == http.StatusOK {
		var schema map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&schema); err != nil {
			return nil, fmt.Errorf("failed to decode swagger response: %w", err)
		}
		defer resp.Body.Close()
		return schema, nil
	}
	
	// If automated discovery fails, create a sample schema for demonstration purposes
	schema := map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]interface{}{
			"title":       g.ServiceName,
			"description": fmt.Sprintf("API for %s", g.ServiceName),
			"version":     "1.0.0",
		},
		"paths": map[string]interface{}{
			"/orders": map[string]interface{}{
				"post": map[string]interface{}{
					"summary": "Create a new order",
					"requestBody": map[string]interface{}{
						"required": true,
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"$ref": "#/components/schemas/OrderRequest",
								},
							},
						},
					},
					"responses": map[string]interface{}{
						"201": map[string]interface{}{
							"description": "Order created successfully",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/OrderResponse",
									},
								},
							},
						},
						"400": map[string]interface{}{
							"description": "Invalid input",
						},
					},
				},
				"get": map[string]interface{}{
					"summary": "List all orders",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "List of orders",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "array",
										"items": map[string]interface{}{
											"$ref": "#/components/schemas/OrderResponse",
										},
									},
								},
							},
						},
					},
				},
			},
			"/orders/{id}": map[string]interface{}{
				"get": map[string]interface{}{
					"summary": "Get an order by ID",
					"parameters": []map[string]interface{}{
						{
							"name": "id",
							"in": "path",
							"required": true,
							"schema": map[string]interface{}{
								"type": "string",
							},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Order details",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/OrderResponse",
									},
								},
							},
						},
						"404": map[string]interface{}{
							"description": "Order not found",
						},
					},
				},
			},
		},
		"components": map[string]interface{}{
			"schemas": map[string]interface{}{
				"OrderRequest": map[string]interface{}{
					"type": "object",
					"required": []string{"items", "customerId"},
					"properties": map[string]interface{}{
						"items": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"$ref": "#/components/schemas/OrderItem",
							},
						},
						"customerId": map[string]interface{}{
							"type": "string",
						},
						"shippingAddress": map[string]interface{}{
							"$ref": "#/components/schemas/Address",
						},
					},
				},
				"OrderResponse": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"orderId": map[string]interface{}{
							"type": "string",
						},
						"items": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"$ref": "#/components/schemas/OrderItem",
							},
						},
						"customerId": map[string]interface{}{
							"type": "string",
						},
						"status": map[string]interface{}{
							"type": "string",
							"enum": []string{"pending", "processing", "shipped", "delivered"},
							"default": "pending",
						},
						"createdAt": map[string]interface{}{
							"type": "string",
							"format": "date-time",
						},
						"total": map[string]interface{}{
							"type": "number",
							"format": "float",
						},
						"shippingAddress": map[string]interface{}{
							"$ref": "#/components/schemas/Address",
						},
					},
				},
				"OrderItem": map[string]interface{}{
					"type": "object",
					"required": []string{"productId", "quantity"},
					"properties": map[string]interface{}{
						"productId": map[string]interface{}{
							"type": "string",
						},
						"quantity": map[string]interface{}{
							"type": "integer",
							"minimum": 1,
						},
						"unitPrice": map[string]interface{}{
							"type": "number",
							"format": "float",
						},
					},
				},
				"Address": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"street": map[string]interface{}{
							"type": "string",
						},
						"city": map[string]interface{}{
							"type": "string",
						},
						"state": map[string]interface{}{
							"type": "string",
						},
						"zipCode": map[string]interface{}{
							"type": "string",
						},
						"country": map[string]interface{}{
							"type": "string",
						},
					},
				},
			},
		},
	}
	
	return schema, nil
}

// SaveSchema saves the generated schema to the specified output path
func (g *Generator) SaveSchema(schema map[string]interface{}, outputPath string) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	// Convert schema to YAML
	yamlData, err := yaml.Marshal(schema)
	if err != nil {
		return fmt.Errorf("failed to marshal schema to YAML: %w", err)
	}
	
	// Write schema to file
	outputFile := filepath.Join(outputPath, "openapi.yaml")
	if err := ioutil.WriteFile(outputFile, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write schema to file: %w", err)
	}
	
	return nil
}