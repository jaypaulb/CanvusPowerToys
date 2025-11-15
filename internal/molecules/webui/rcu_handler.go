package webui

import (
	"encoding/json"
	"fmt"
	"net/http"

	webuiatoms "github.com/jaypaulb/CanvusPowerToys/internal/atoms/webui"
)

// RCUHandler handles RCU configuration API endpoints.
type RCUHandler struct {
	apiClient     *webuiatoms.APIClient
	canvasService *CanvasService
}

// NewRCUHandler creates a new RCU handler.
func NewRCUHandler(apiClient *webuiatoms.APIClient, canvasService *CanvasService) *RCUHandler {
	return &RCUHandler{
		apiClient:     apiClient,
		canvasService: canvasService,
	}
}

// HandleConfig handles GET/POST /api/rcu/config - Get/Set RCU configuration.
func (h *RCUHandler) HandleConfig(w http.ResponseWriter, r *http.Request) {
	canvasID := h.canvasService.GetCanvasID()
	if canvasID == "" {
		http.Error(w, "Canvas not available", http.StatusServiceUnavailable)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Get RCU config from Canvus API
		endpoint := fmt.Sprintf("/api/v1/canvases/%s/rcu/config", canvasID)
		data, err := h.apiClient.Get(endpoint)
		if err != nil {
			// Return default config if not found
			defaultConfig := map[string]interface{}{
				"enabled": false,
				"port":    8080,
				"timeout": 30,
			}
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(defaultConfig)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write(data)

	case http.MethodPost:
		var req struct {
			Enabled bool `json:"enabled"`
			Port    int  `json:"port"`
			Timeout int  `json:"timeout"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Update RCU config via Canvus API
		endpoint := fmt.Sprintf("/api/v1/canvases/%s/rcu/config", canvasID)
		data, err := h.apiClient.Post(endpoint, req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to update config: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write(data)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleStatus handles GET /api/rcu/status - Get RCU status.
func (h *RCUHandler) HandleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	canvasID := h.canvasService.GetCanvasID()
	if canvasID == "" {
		http.Error(w, "Canvas not available", http.StatusServiceUnavailable)
		return
	}

	// Get RCU status from Canvus API
	endpoint := fmt.Sprintf("/api/v1/canvases/%s/rcu/status", canvasID)
	data, err := h.apiClient.Get(endpoint)
	if err != nil {
		// Return default status if not found
		defaultStatus := map[string]interface{}{
			"connected":   false,
			"last_update": nil,
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(defaultStatus)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// HandleTest handles POST /api/rcu/test - Test RCU connection.
func (h *RCUHandler) HandleTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	canvasID := h.canvasService.GetCanvasID()
	if canvasID == "" {
		http.Error(w, "Canvas not available", http.StatusServiceUnavailable)
		return
	}

	// Test RCU connection via Canvus API
	endpoint := fmt.Sprintf("/api/v1/canvases/%s/rcu/test", canvasID)
	data, err := h.apiClient.Post(endpoint, nil)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

