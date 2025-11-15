package webui

import (
	"encoding/json"
	"net/http"

	webuimolecules "github.com/jaypaulb/CanvusPowerToys/internal/molecules/webui"
)

// APIRoutes handles registration of API routes for the WebUI server.
type APIRoutes struct {
	canvasService *webuimolecules.CanvasService
	sseHandler    *SSEHandler
}

// NewAPIRoutes creates a new API routes handler.
func NewAPIRoutes(canvasService *webuimolecules.CanvasService) *APIRoutes {
	sseHandler := NewSSEHandler(canvasService)

	return &APIRoutes{
		canvasService: canvasService,
		sseHandler:    sseHandler,
	}
}

// RegisterRoutes registers all API routes with the given mux.
func (ar *APIRoutes) RegisterRoutes(mux *http.ServeMux) {
	// SSE endpoint for canvas_id updates
	mux.HandleFunc("/api/subscribe-workspace", ar.sseHandler.HandleSubscribe)

	// Canvas info endpoint (current canvas_id and canvas_name)
	mux.HandleFunc("/api/canvas/info", ar.handleCanvasInfo)

	// Installation info endpoint
	mux.HandleFunc("/api/installation/info", ar.handleInstallationInfo)

	// Health check endpoint
	mux.HandleFunc("/api/health", ar.handleHealth)
}

// handleCanvasInfo returns current canvas information.
func (ar *APIRoutes) handleCanvasInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	response := map[string]interface{}{
		"canvas_id":   ar.canvasService.GetCanvasID(),
		"canvas_name": ar.canvasService.GetCanvasName(),
		"client_id":   ar.canvasService.GetClientID(),
		"connected":   ar.canvasService.IsConnected(),
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// handleInstallationInfo returns installation information.
func (ar *APIRoutes) handleInstallationInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	response := map[string]interface{}{
		"installation_name": ar.canvasService.GetInstallationName(),
		"client_id":         ar.canvasService.GetClientID(),
		"connected":         ar.canvasService.IsConnected(),
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// handleHealth returns server health status.
func (ar *APIRoutes) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	response := map[string]interface{}{
		"status":    "ok",
		"connected": ar.canvasService.IsConnected(),
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

