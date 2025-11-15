package webui

import (
	"encoding/json"
	"net/http"
	"strings"

	webuiatoms "github.com/jaypaulb/CanvusPowerToys/internal/atoms/webui"
)

// APIRoutes handles registration of API routes for the WebUI server.
type APIRoutes struct {
	canvasService *CanvasService
	sseHandler    *SSEHandler
	apiClient     *webuiatoms.APIClient
	pagesHandler  *PagesHandler
	macrosHandler *MacrosHandler
	uploadHandler *UploadHandler
	rcuHandler    *RCUHandler
}

// NewAPIRoutes creates a new API routes handler.
func NewAPIRoutes(canvasService *CanvasService, apiClient *webuiatoms.APIClient, uploadDir string) *APIRoutes {
	sseHandler := NewSSEHandler(canvasService)
	pagesHandler := NewPagesHandler(apiClient, canvasService)
	macrosHandler := NewMacrosHandler(apiClient, canvasService)
	uploadHandler := NewUploadHandler(apiClient, canvasService, uploadDir)
	rcuHandler := NewRCUHandler(apiClient, canvasService)

	return &APIRoutes{
		canvasService: canvasService,
		sseHandler:    sseHandler,
		apiClient:     apiClient,
		pagesHandler:  pagesHandler,
		macrosHandler: macrosHandler,
		uploadHandler: uploadHandler,
		rcuHandler:    rcuHandler,
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

	// Pages endpoints
	mux.HandleFunc("/api/pages", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			ar.pagesHandler.HandleList(w, r)
		} else if r.Method == http.MethodPost {
			ar.pagesHandler.HandleCreate(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Macros endpoints
	mux.HandleFunc("/api/macros/groups", ar.macrosHandler.HandleGroups)
	mux.HandleFunc("/api/macros/pinned", ar.macrosHandler.HandlePinned)
	mux.HandleFunc("/api/macros/move", ar.macrosHandler.HandleMove)
	mux.HandleFunc("/api/macros/copy", ar.macrosHandler.HandleCopy)
	mux.HandleFunc("/api/macros/pin-all", ar.macrosHandler.HandlePinAll)
	mux.HandleFunc("/api/macros/unpin-all", ar.macrosHandler.HandleUnpin)
	mux.HandleFunc("/api/macros/auto-grid", ar.macrosHandler.HandleAutoGrid)
	mux.HandleFunc("/api/macros/group-color", ar.macrosHandler.HandleGroupColor)
	mux.HandleFunc("/api/macros/group-title", ar.macrosHandler.HandleGroupTitle)

	// Remote upload endpoints
	mux.HandleFunc("/api/remote-upload", ar.uploadHandler.HandleUpload)
	mux.HandleFunc("/api/remote-upload/history", ar.uploadHandler.HandleHistory)

	// RCU endpoints
	mux.HandleFunc("/api/rcu/config", ar.rcuHandler.HandleConfig)
	mux.HandleFunc("/api/rcu/status", ar.rcuHandler.HandleStatus)
	mux.HandleFunc("/api/rcu/test", ar.rcuHandler.HandleTest)
}

// contains checks if a string contains a substring.
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
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

