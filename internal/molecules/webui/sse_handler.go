package webui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	webuimolecules "github.com/jaypaulb/CanvusPowerToys/internal/molecules/webui"
)

// SSEHandler handles Server-Sent Events for canvas_id updates.
type SSEHandler struct {
	canvasService *webuimolecules.CanvasService
}

// NewSSEHandler creates a new SSE handler.
func NewSSEHandler(canvasService *webuimolecules.CanvasService) *SSEHandler {
	return &SSEHandler{
		canvasService: canvasService,
	}
}

// HandleSubscribe handles the SSE subscription endpoint.
func (h *SSEHandler) HandleSubscribe(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

	// Create a channel to track client disconnection
	clientGone := r.Context().Done()

	// Send initial canvas state
	h.sendCanvasUpdate(w, h.canvasService.GetCanvasID(), h.canvasService.GetCanvasName())

	// Create ticker to send periodic updates
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Track last sent canvas_id to avoid duplicate sends
	lastCanvasID := h.canvasService.GetCanvasID()

	for {
		select {
		case <-clientGone:
			// Client disconnected
			return
		case <-ticker.C:
			// Periodic update check
			currentCanvasID := h.canvasService.GetCanvasID()
			if currentCanvasID != lastCanvasID {
				h.sendCanvasUpdate(w, currentCanvasID, h.canvasService.GetCanvasName())
				lastCanvasID = currentCanvasID
			} else {
				// Send keepalive
				h.sendKeepalive(w)
			}
		}
	}
}

// sendCanvasUpdate sends a canvas update event.
func (h *SSEHandler) sendCanvasUpdate(w http.ResponseWriter, canvasID, canvasName string) {
	event := map[string]interface{}{
		"canvas_id":   canvasID,
		"canvas_name": canvasName,
		"timestamp":   time.Now().Unix(),
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		fmt.Printf("Error marshaling canvas event: %v\n", err)
		return
	}

	// Send SSE formatted event
	fmt.Fprintf(w, "event: canvas_update\n")
	fmt.Fprintf(w, "data: %s\n\n", string(eventJSON))

	// Flush to ensure data is sent immediately
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

// sendKeepalive sends a keepalive comment to maintain connection.
func (h *SSEHandler) sendKeepalive(w http.ResponseWriter) {
	fmt.Fprintf(w, ": keepalive\n\n")
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

