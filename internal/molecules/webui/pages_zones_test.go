package webui

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	webuiatoms "github.com/jaypaulb/CanvusPowerToys/internal/atoms/webui"
)

// TestHandleCreateZones_NoCanvas tests that HandleCreateZones returns 503 when canvas is not available
func TestHandleCreateZones_NoCanvas(t *testing.T) {
	// Create minimal canvas service (no canvas ID)
	canvasTracker := webuiatoms.NewCanvasTracker()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	canvasService := &CanvasService{
		apiBaseURL:       "http://test.example.com",
		authToken:        "test-token",
		installationName: "Unknown",
		canvasTracker:    canvasTracker,
		ctx:              ctx,
		cancel:           cancel,
	}

	apiClient := webuiatoms.NewAPIClient("http://test.example.com", "test-token")
	handler := &PagesHandler{
		canvasService: canvasService,
		apiClient:     apiClient,
	}

	// Create request with valid body
	reqBody := `{"gridSize": 3, "gridPattern": "Z"}`
	req := httptest.NewRequest("POST", "/create-zones", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.HandleCreateZones(w, req)

	// Check response
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected 503 Service Unavailable, got %d", w.Code)
	}

	// Check response body is JSON error
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || success {
		t.Error("Expected success: false in error response")
	}

	if errorMsg, ok := response["error"].(string); !ok || errorMsg != "Canvas not available" {
		t.Errorf("Expected error message 'Canvas not available', got: %v", errorMsg)
	}
}
