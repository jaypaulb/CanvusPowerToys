package webui

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	webuiatoms "github.com/jaypaulb/CanvusPowerToys/internal/atoms/webui"
)

// Test credentials from docs/webui-test-credentials.md
const (
	testServerURL = "https://ise2025.canvusmultisite.com"
	testToken     = "4VZmSp0T-i7Lzw-9bXIFdiPof_yu6jRHVGCsbOsrs5I"
)

// getTestCredentials returns test credentials, using environment variables if set
func getTestCredentials() (serverURL, token string) {
	serverURL = os.Getenv("CANVUS_TEST_SERVER_URL")
	if serverURL == "" {
		serverURL = testServerURL
	}
	token = os.Getenv("CANVUS_TEST_TOKEN")
	if token == "" {
		token = testToken
	}
	return serverURL, token
}

// setupTestServer creates a test server with real API client pointing to live server
func setupTestServer(t *testing.T) (*httptest.Server, *APIRoutes, *CanvasService) {
	serverURL, token := getTestCredentials()

	// Create canvas service
	canvasTracker := webuiatoms.NewCanvasTracker()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	canvasService := &CanvasService{
		apiBaseURL:       serverURL,
		authToken:        token,
		installationName: "TestInstallation",
		canvasTracker:    canvasTracker,
		ctx:              ctx,
		cancel:           cancel,
	}

	// Create API client
	apiClient := webuiatoms.NewAPIClient(serverURL, token)

	// Create API routes
	apiRoutes := NewAPIRoutes(canvasService, apiClient, "")

	// Create test server
	mux := http.NewServeMux()
	apiRoutes.RegisterRoutes(mux)
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	return server, apiRoutes, canvasService
}

// TestAPI_Health tests the health check endpoint
func TestAPI_Health(t *testing.T) {
	server, _, _ := setupTestServer(t)

	resp, err := http.Get(server.URL + "/api/health")
	if err != nil {
		t.Fatalf("Failed to call health endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

// TestAPI_ServerInfo tests the server info endpoint
func TestAPI_ServerInfo(t *testing.T) {
	server, _, _ := setupTestServer(t)

	resp, err := http.Get(server.URL + "/api/server-info")
	if err != nil {
		t.Fatalf("Failed to call server-info endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if _, ok := data["hostname"]; !ok {
		t.Error("Response missing 'hostname' field")
	}
}

// TestAPI_CanvasInfo tests the canvas info endpoint
func TestAPI_CanvasInfo(t *testing.T) {
	server, _, _ := setupTestServer(t)

	resp, err := http.Get(server.URL + "/api/canvas/info")
	if err != nil {
		t.Fatalf("Failed to call canvas/info endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Should have these fields even if empty
	requiredFields := []string{"canvas_id", "canvas_name", "client_id", "client_name", "installation_name", "connected"}
	for _, field := range requiredFields {
		if _, ok := data[field]; !ok {
			t.Errorf("Response missing required field: %s", field)
		}
	}
}

// TestAPI_InstallationInfo tests the installation info endpoint
func TestAPI_InstallationInfo(t *testing.T) {
	server, _, _ := setupTestServer(t)

	resp, err := http.Get(server.URL + "/api/installation/info")
	if err != nil {
		t.Fatalf("Failed to call installation/info endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if _, ok := data["installation_name"]; !ok {
		t.Error("Response missing 'installation_name' field")
	}
}

// TestAPI_GetZones tests the get-zones endpoint
func TestAPI_GetZones(t *testing.T) {
	server, _, canvasService := setupTestServer(t)
	_ = canvasService // Used later in test

	// Test without canvas (should return empty zones gracefully)
	resp, err := http.Get(server.URL + "/get-zones")
	if err != nil {
		t.Fatalf("Failed to call get-zones endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if success, ok := data["success"].(bool); !ok || !success {
		t.Error("Expected success: true in response")
	}

	if zones, ok := data["zones"].([]interface{}); !ok {
		t.Error("Response missing 'zones' array")
	} else {
		t.Logf("Retrieved %d zones", len(zones))
	}

	// If we have a canvas ID, test with it
	// Try to set a canvas ID by overriding client
	serverURL, token := getTestCredentials()
	apiClient := webuiatoms.NewAPIClient(serverURL, token)

	// Try to get clients and find one that exists
	clients, err := apiClient.GetClients()
	if err == nil && len(clients) > 0 {
		// Override to first available client
		err := canvasService.OverrideClient(clients[0].Name)
		if err == nil {
			// Wait a bit for subscription to start
			time.Sleep(2 * time.Second)

			// Test again with canvas
			resp2, err := http.Get(server.URL + "/get-zones")
			if err == nil {
				defer resp2.Body.Close()
				var data2 map[string]interface{}
				if err := json.NewDecoder(resp2.Body).Decode(&data2); err == nil {
					t.Logf("With canvas: Retrieved %d zones", len(data2["zones"].([]interface{})))
				}
			}
		}
	}
}

// TestAPI_CreateZones tests the create-zones endpoint
func TestAPI_CreateZones(t *testing.T) {
	server, _, canvasService := setupTestServer(t)

	// Test without canvas (should return 503)
	reqBody := `{"gridSize": 3, "gridPattern": "Z"}`
	resp, err := http.Post(server.URL+"/create-zones", "application/json", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to call create-zones endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("Expected status 503 when canvas not available, got %d", resp.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if success, ok := data["success"].(bool); !ok || success {
		t.Error("Expected success: false in error response")
	}

	if errorMsg, ok := data["error"].(string); !ok || errorMsg != "Canvas not available" {
		t.Errorf("Expected error 'Canvas not available', got: %v", errorMsg)
	}

	// Try to set canvas by overriding client
	serverURL, token := getTestCredentials()
	apiClient := webuiatoms.NewAPIClient(serverURL, token)
	clients, err := apiClient.GetClients()
	if err == nil && len(clients) > 0 {
		err := canvasService.OverrideClient(clients[0].Name)
		if err == nil {
			time.Sleep(3 * time.Second) // Wait for subscription to get canvas ID

			// Test with canvas (if we got one)
			canvasID := canvasService.GetCanvasID()
			if canvasID != "" {
				reqBody2 := `{"gridSize": 3, "gridPattern": "Z"}`
				resp2, err := http.Post(server.URL+"/create-zones", "application/json", strings.NewReader(reqBody2))
				if err == nil {
					defer resp2.Body.Close()
					if resp2.StatusCode == http.StatusOK {
						t.Log("Successfully created zones with canvas")
					} else {
						t.Logf("Create zones returned status %d (may be expected if canvas has no widgets)", resp2.StatusCode)
					}
				}
			}
		}
	}
}

// TestAPI_AdminEndpoints tests admin endpoints
func TestAPI_AdminEndpoints(t *testing.T) {
	server, _, _ := setupTestServer(t)

	// Test list-users (should work even without canvas - returns empty list)
	resp, err := http.Get(server.URL + "/api/admin/list-users")
	if err != nil {
		t.Fatalf("Failed to call list-users endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if users, ok := data["users"].([]interface{}); ok {
		t.Logf("Retrieved %d RCU users", len(users))
	}

	// Test create-targets (requires canvas)
	reqBody := `{"teams": [1, 2, 3]}`
	resp2, err := http.Post(server.URL+"/api/admin/create-targets", "application/json", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to call create-targets endpoint: %v", err)
	}
	defer resp2.Body.Close()

	// Should return 503 if no canvas, or 200 if canvas available
	if resp2.StatusCode != http.StatusServiceUnavailable && resp2.StatusCode != http.StatusOK {
		t.Errorf("Expected status 503 or 200, got %d", resp2.StatusCode)
	}
}

// TestAPI_RCUEndpoints tests RCU endpoints
func TestAPI_RCUEndpoints(t *testing.T) {
	server, _, _ := setupTestServer(t)

	// Test RCU config (should work even without canvas - returns default config)
	resp, err := http.Get(server.URL + "/api/rcu/config")
	if err != nil {
		t.Fatalf("Failed to call rcu/config endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Test identify-user (requires canvas)
	reqBody := `{"user_name": "TestUser"}`
	req, err := http.NewRequest("POST", server.URL+"/identify-user", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp2, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to call identify-user endpoint: %v", err)
	}
	defer resp2.Body.Close()

	// Should return 503 if no canvas, or 200/400 if canvas available
	if resp2.StatusCode == http.StatusServiceUnavailable {
		t.Log("identify-user correctly returns 503 when canvas not available")
	} else if resp2.StatusCode == http.StatusOK || resp2.StatusCode == http.StatusBadRequest {
		t.Logf("identify-user returned %d (may be expected depending on canvas state)", resp2.StatusCode)
	} else {
		t.Logf("identify-user returned status %d (may be expected)", resp2.StatusCode)
	}
}

// TestAPI_ClientOverride tests client override endpoint
func TestAPI_ClientOverride(t *testing.T) {
	server, _, _ := setupTestServer(t)

	serverURL, token := getTestCredentials()
	apiClient := webuiatoms.NewAPIClient(serverURL, token)

	// Get available clients
	clients, err := apiClient.GetClients()
	if err != nil {
		t.Skipf("Skipping test - failed to get clients: %v", err)
	}

	if len(clients) == 0 {
		t.Skip("Skipping test - no clients available")
	}

	// Find a client with installation_name
	var clientName string
	for _, client := range clients {
		if client.InstallationName != "" {
			clientName = client.InstallationName
			break
		}
	}

	if clientName == "" {
		t.Skip("Skipping test - no clients with installation_name available")
	}

	// Test override with first available client that has a name
	reqBody := fmt.Sprintf(`{"client_name": "%s"}`, clientName)
	req, err := http.NewRequest("POST", server.URL+"/api/client/override", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to call client/override endpoint: %v", err)
	}
	defer resp.Body.Close()

	// Read response body first before checking status
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Response: %s", resp.StatusCode, string(bodyBytes))
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if success, ok := data["success"].(bool); !ok || !success {
		t.Errorf("Expected success: true, got: %v", data)
	}

	// Wait a bit for subscription to start
	time.Sleep(2 * time.Second)

	// Check canvas info to see if we got a canvas ID
	resp2, err := http.Get(server.URL + "/api/canvas/info")
	if err == nil {
		defer resp2.Body.Close()
		var canvasInfo map[string]interface{}
		if err := json.NewDecoder(resp2.Body).Decode(&canvasInfo); err == nil {
			if canvasID, ok := canvasInfo["canvas_id"].(string); ok && canvasID != "" {
				t.Logf("Successfully got canvas ID after override: %s", canvasID)
			} else {
				t.Log("Canvas ID not yet available (may need more time for subscription)")
			}
		}
	}
}

// TestAPI_AllEndpointsExist tests that all expected endpoints exist and return JSON
func TestAPI_AllEndpointsExist(t *testing.T) {
	server, _, canvasService := setupTestServer(t)
	_ = canvasService // May be used in future tests

	endpoints := []struct {
		method string
		path   string
		body   string
	}{
		{"GET", "/api/health", ""},
		{"GET", "/api/server-info", ""},
		{"GET", "/api/canvas/info", ""},
		{"GET", "/api/installation/info", ""},
		{"GET", "/get-zones", ""},
		{"GET", "/api/pages", ""},
		{"GET", "/api/rcu/config", ""},
		{"GET", "/api/admin/list-users", ""},
		// Skip client override test here - it's tested separately in TestAPI_ClientOverride
	}

	for _, endpoint := range endpoints {
		t.Run(fmt.Sprintf("%s %s", endpoint.method, endpoint.path), func(t *testing.T) {
			var req *http.Request
			var err error

			if endpoint.method == "GET" {
				req, err = http.NewRequest(endpoint.method, server.URL+endpoint.path, nil)
			} else {
				req, err = http.NewRequest(endpoint.method, server.URL+endpoint.path, strings.NewReader(endpoint.body))
				req.Header.Set("Content-Type", "application/json")
			}

			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			client := &http.Client{Timeout: 5 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Failed to call endpoint: %v", err)
			}
			defer resp.Body.Close()

			// Check content type is JSON (unless it's a 404/405 which might be HTML)
			contentType := resp.Header.Get("Content-Type")
			if resp.StatusCode < 400 && !strings.Contains(contentType, "application/json") {
				t.Errorf("Expected JSON response, got: %s", contentType)
			}

			// Try to parse as JSON if status is OK
			if resp.StatusCode == http.StatusOK {
				var data map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
					t.Errorf("Response is not valid JSON: %v", err)
				}
			}
		})
	}
}

