package schema

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Generator struct {
	providerName string
	baseURL      string
	client       *http.Client
}

func NewGenerator(providerName, baseURL string) *Generator {
	return &Generator{
		providerName: providerName,
		baseURL:      baseURL,
		client:       &http.Client{},
	}
}

func (g *Generator) GenerateSchema(outputPath string) error {
	// In a real implementation, this would introspect the API
	// For this example, we'll create a simple schema for an order service
	
	schema := map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]interface{}{
			"title":       fmt.Sprintf("%s API", g.providerName),
			"description": fmt.Sprintf("API contract for %s", g.providerName),
			"version":     "1.0.0",
		},
		"servers": []map[string]interface{}{
			{
				"url": g.baseURL,
			},
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
									"type": "object",
									"required": []string{
										"userId", "items",
									},
									"properties": map[string]interface{}{
										"userId": map[string]interface{}{
											"type": "string",
										},
										"items": map[string]interface{}{
											"type": "array",
											"items": map[string]interface{}{
												"type": "object",
												"required": []string{
													"productId", "quantity",
												},
												"properties": map[string]interface{}{
													"productId": map[string]interface{}{
														"type": "string",
													},
													"quantity": map[string]interface{}{
														"type": "integer",
														"minimum": 1,
													},
												},
											},
										},
									},
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
										"type": "object",
										"properties": map[string]interface{}{
											"orderId": map[string]interface{}{
												"type": "string",
											},
											"status": map[string]interface{}{
												"type": "string",
												"enum": []string{"pending"},
											},
											"createdAt": map[string]interface{}{
												"type": "string",
												"format": "date-time",
											},
										},
									},
								},
							},
						},
						"400": map[string]interface{}{
							"description": "Invalid request",
						},
					},
				},
			},
			"/orders/{orderId}": map[string]interface{}{
				"get": map[string]interface{}{
					"summary": "Get order by ID",
					"parameters": []map[string]interface{}{
						{
							"name": "orderId",
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
										"type": "object",
										"properties": map[string]interface{}{
											"orderId": map[string]interface{}{
												"type": "string",
											},
											"userId": map[string]interface{}{
												"type": "string",
											},
											"status": map[string]interface{}{
												"type": "string",
												"enum": []string{"pending", "processing", "shipped", "delivered", "cancelled"},
											},
											"items": map[string]interface{}{
												"type": "array",
												"items": map[string]interface{}{
													"type": "object",
													"properties": map[string]interface{}{
														"productId": map[string]interface{}{
															"type": "string",
														},
														"quantity": map[string]interface{}{
															"type": "integer",
														},
													},
												},
											},
											"createdAt": map[string]interface{}{
												"type": "string",
												"format": "date-time",
											},
										},
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
	}
	
	// Create directory if it doesn't exist
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	// Convert to YAML and write to file
	yamlData, err := yaml.Marshal(schema)
	if err != nil {
		return fmt.Errorf("failed to marshal schema: %w", err)
	}
	
	if err := os.WriteFile(outputPath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write schema to file: %w", err)
	}
	
	return nil
}