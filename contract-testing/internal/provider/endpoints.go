package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Endpoint struct {
	Path    string
	Method  string
	Handler http.HandlerFunc
}

func GetOrderServiceEndpoints() []Endpoint {
	return []Endpoint{
		{
			Path:    "/orders",
			Method:  "POST",
			Handler: createOrderHandler,
		},
		{
			Path:    "/orders/{orderId}",
			Method:  "GET",
			Handler: getOrderHandler,
		},
	}
}

func createOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request format",
		})
		return
	}
	
	// Validate required fields
	userID, ok := request["userId"].(string)
	if !ok || userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "userId is required and must be a string",
		})
		return
	}
	
	items, ok := request["items"].([]interface{})
	if !ok || len(items) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "items is required and must be a non-empty array",
		})
		return
	}
	
	// Create a new order (in a real implementation, this would interact with a database)
	orderID := "ord_" + generateRandomID()
	
	response := map[string]interface{}{
		"orderId":   orderID,
		"status":    "pending",
		"createdAt": time.Now().Format(time.RFC3339),
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func getOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	// Extract order ID from URL
	regex := regexp.MustCompile(`/orders/([^/]+)`)
	matches := regex.FindStringSubmatch(r.URL.Path)
	
	if len(matches) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid order ID",
		})
		return
	}
	
	orderID := matches[1]
	
	// In a real implementation, this would fetch the order from a database
	// For this example, we'll just generate a mock response
	
	// Simulate not found for certain IDs
	if strings.HasPrefix(orderID, "notfound_") {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Order not found",
		})
		return
	}
	
	response := map[string]interface{}{
		"orderId": orderID,
		"userId": "user_123",
		"status": "pending",
		"items": []map[string]interface{}{
			{
				"productId": "prod_1",
				"quantity":  2,
			},
		},
		"createdAt": time.Now().Format(time.RFC3339),
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func generateRandomID() string {
	// In a real implementation, this would generate a unique ID
	return fmt.Sprintf("%d", time.Now().UnixNano())
}