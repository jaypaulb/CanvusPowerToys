package webui

import (
	"encoding/json"
	"fmt"
	"net/http"

	webuiatoms "github.com/jaypaulb/CanvusPowerToys/internal/atoms/webui"
)

// PagesHandler handles pages management API endpoints.
type PagesHandler struct {
	apiClient *webuiatoms.APIClient
	canvasService *CanvasService
}

// NewPagesHandler creates a new pages handler.
func NewPagesHandler(apiClient *webuiatoms.APIClient, canvasService *CanvasService) *PagesHandler {
	return &PagesHandler{
		apiClient:     apiClient,
		canvasService: canvasService,
	}
}

// HandleList handles GET /api/pages - List all pages.
func (h *PagesHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	canvasID := h.canvasService.GetCanvasID()
	if canvasID == "" {
		http.Error(w, "Canvas not available", http.StatusServiceUnavailable)
		return
	}

	// Fetch anchors (pages) from Canvus API
	// Note: "Pages" in WebUI are "Anchors" in Canvus API
	endpoint := fmt.Sprintf("/api/v1/canvases/%s/anchors", canvasID)
	data, err := h.apiClient.Get(endpoint)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch pages: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// HandleCreate handles POST /api/pages - Create a new anchor (page).
// Note: Anchors are created as widgets with widget_type="Anchor" in the Canvus API.
func (h *PagesHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Page name is required", http.StatusBadRequest)
		return
	}

	canvasID := h.canvasService.GetCanvasID()
	if canvasID == "" {
		http.Error(w, "Canvas not available", http.StatusServiceUnavailable)
		return
	}

	// Create anchor (page) as a widget via Canvus API
	// Anchors are widgets with widget_type="Anchor"
	payload := map[string]interface{}{
		"widget_type": "Anchor",
		"anchor_name": req.Name,
	}
	endpoint := fmt.Sprintf("/api/v1/canvases/%s/widgets", canvasID)
	data, err := h.apiClient.Post(endpoint, payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create page: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}


