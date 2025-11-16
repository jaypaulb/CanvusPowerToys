package webui

import (
	"testing"
)

// TestCanvasService_MinimalMode tests that a minimal canvas service can be created
// and GetCanvasID() returns empty string without panicking
func TestCanvasService_MinimalMode(t *testing.T) {
	// Create minimal canvas service (as done in manager.go when auto-detection fails)
	cs := &CanvasService{
		apiBaseURL:       "http://test.example.com",
		authToken:        "test-token",
		installationName: "Unknown",
	}

	// GetCanvasID should return empty string, not panic
	canvasID := cs.GetCanvasID()
	if canvasID != "" {
		t.Errorf("Expected empty canvas ID for minimal service, got: %s", canvasID)
	}

	// GetCanvasName should return empty string
	canvasName := cs.GetCanvasName()
	if canvasName != "" {
		t.Errorf("Expected empty canvas name for minimal service, got: %s", canvasName)
	}

	// IsConnected should return false
	if cs.IsConnected() {
		t.Error("Expected minimal service to not be connected")
	}
}

// TestCanvasService_GetCanvasID_WithNilTracker tests that GetCanvasID handles nil canvasTracker gracefully
func TestCanvasService_GetCanvasID_WithNilTracker(t *testing.T) {
	cs := &CanvasService{
		apiBaseURL:       "http://test.example.com",
		authToken:        "test-token",
		installationName: "Unknown",
		canvasTracker:    nil, // Explicitly nil
	}

	// Should not panic, should return empty string
	canvasID := cs.GetCanvasID()
	if canvasID != "" {
		t.Errorf("Expected empty canvas ID when tracker is nil, got: %s", canvasID)
	}
}

