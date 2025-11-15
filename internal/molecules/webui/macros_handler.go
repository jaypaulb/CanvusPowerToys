package webui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	webuiatoms "github.com/jaypaulb/CanvusPowerToys/internal/atoms/webui"
)

// MacrosHandler handles macros management API endpoints.
// NOTE: The Canvus API documentation does not include macros endpoints.
// Macros may be implemented as widgets or may require a different approach.
// This handler uses assumed endpoints that need verification against actual API.
type MacrosHandler struct {
	apiClient     *webuiatoms.APIClient
	canvasService *CanvasService
}

// NewMacrosHandler creates a new macros handler.
func NewMacrosHandler(apiClient *webuiatoms.APIClient, canvasService *CanvasService) *MacrosHandler {
	return &MacrosHandler{
		apiClient:     apiClient,
		canvasService: canvasService,
	}
}

// HandleGroups handles GET /api/macros/groups - List macro groups.
func (h *MacrosHandler) HandleGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	canvasID := h.canvasService.GetCanvasID()
	if canvasID == "" {
		http.Error(w, "Canvas not available", http.StatusServiceUnavailable)
		return
	}

	// Fetch macro groups from Canvus API
	endpoint := fmt.Sprintf("/api/v1/canvases/%s/macros/groups", canvasID)
	data, err := h.apiClient.Get(endpoint)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch groups: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// HandlePinned handles GET /api/macros/pinned - List pinned macros.
func (h *MacrosHandler) HandlePinned(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	canvasID := h.canvasService.GetCanvasID()
	if canvasID == "" {
		http.Error(w, "Canvas not available", http.StatusServiceUnavailable)
		return
	}

	// Fetch pinned macros from Canvus API
	endpoint := fmt.Sprintf("/api/v1/canvases/%s/macros?pinned=true", canvasID)
	data, err := h.apiClient.Get(endpoint)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch pinned macros: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// HandleMove handles POST /api/macros/{id}/move - Move macro to different group.
func (h *MacrosHandler) HandleMove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract macro ID from path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 || pathParts[0] != "api" || pathParts[1] != "macros" {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	macroID := pathParts[2]
	if pathParts[3] != "move" {
		http.Error(w, "Invalid endpoint", http.StatusBadRequest)
		return
	}

	var req struct {
		GroupID string `json:"group_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	canvasID := h.canvasService.GetCanvasID()
	if canvasID == "" {
		http.Error(w, "Canvas not available", http.StatusServiceUnavailable)
		return
	}

	// Move macro via Canvus API
	endpoint := fmt.Sprintf("/api/v1/canvases/%s/macros/%s/move", canvasID, macroID)
	data, err := h.apiClient.Post(endpoint, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to move macro: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// HandleCopy handles POST /api/macros/{id}/copy - Copy macro to different group.
func (h *MacrosHandler) HandleCopy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract macro ID from path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 || pathParts[0] != "api" || pathParts[1] != "macros" {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	macroID := pathParts[2]
	if pathParts[3] != "copy" {
		http.Error(w, "Invalid endpoint", http.StatusBadRequest)
		return
	}

	var req struct {
		GroupID string `json:"group_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	canvasID := h.canvasService.GetCanvasID()
	if canvasID == "" {
		http.Error(w, "Canvas not available", http.StatusServiceUnavailable)
		return
	}

	// Copy macro via Canvus API
	endpoint := fmt.Sprintf("/api/v1/canvases/%s/macros/%s/copy", canvasID, macroID)
	data, err := h.apiClient.Post(endpoint, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to copy macro: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// HandleUnpin handles POST /api/macros/{id}/unpin - Unpin a macro.
func (h *MacrosHandler) HandleUnpin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract macro ID from path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 || pathParts[0] != "api" || pathParts[1] != "macros" {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	macroID := pathParts[2]
	if pathParts[3] != "unpin" {
		http.Error(w, "Invalid endpoint", http.StatusBadRequest)
		return
	}

	canvasID := h.canvasService.GetCanvasID()
	if canvasID == "" {
		http.Error(w, "Canvas not available", http.StatusServiceUnavailable)
		return
	}

	// Unpin macro via Canvus API
	endpoint := fmt.Sprintf("/api/v1/canvases/%s/macros/%s/unpin", canvasID, macroID)
	_, err := h.apiClient.Post(endpoint, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unpin macro: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

