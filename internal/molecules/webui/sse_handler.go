package webui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SSEHandler handles Server-Sent Events for canvas_id updates.
type SSEHandler struct {
	canvasService *CanvasService
}

// NewSSEHandler creates a new SSE handler.
func NewSSEHandler(canvasService *CanvasService) *SSEHandler {
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

	// Create a channel to track client disconnection and server shutdown
	// r.Context() is cancelled when the server shuts down, so this will close SSE connections
	clientGone := r.Context().Done()

	// Send initial canvas state
	h.sendCanvasUpdate(w, h.canvasService.GetCanvasID(), h.canvasService.GetCanvasName())

	// Create ticker to send periodic updates
	ticker := time.NewTicker(5 * time.Second) // Check more frequently for canvas name updates
	defer ticker.Stop()

	// Track last sent canvas_id and canvas_name to detect changes
	lastCanvasID := h.canvasService.GetCanvasID()
	lastCanvasName := h.canvasService.GetCanvasName()

	for {
		select {
		case <-clientGone:
			// Client disconnected or server shutting down
			fmt.Printf("[SSEHandler] Connection closed (client disconnected or server shutdown)\n")
			return
		case <-ticker.C:
			// Check if context is cancelled before sending (server might be shutting down)
			select {
			case <-clientGone:
				fmt.Printf("[SSEHandler] Context cancelled during tick, closing connection\n")
				return
			default:
				// Context still active, proceed with update
			}

			// Periodic update check - send update if canvas_id OR canvas_name changed
			currentCanvasID := h.canvasService.GetCanvasID()
			currentCanvasName := h.canvasService.GetCanvasName()

			if currentCanvasID != lastCanvasID || currentCanvasName != lastCanvasName {
				h.sendCanvasUpdate(w, currentCanvasID, currentCanvasName)
				lastCanvasID = currentCanvasID
				lastCanvasName = currentCanvasName
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
		"client_name": h.canvasService.GetClientName(),
		"client_id":   h.canvasService.GetClientID(),
		"timestamp":   time.Now().Unix(),
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		fmt.Printf("Error marshaling canvas event: %v\n", err)
		return
	}

	// Send SSE formatted event - check for write errors
	if _, err := fmt.Fprintf(w, "event: canvas_update\n"); err != nil {
		fmt.Printf("Error writing SSE event header: %v\n", err)
		return
	}
	if _, err := fmt.Fprintf(w, "data: %s\n\n", string(eventJSON)); err != nil {
		fmt.Printf("Error writing SSE event data: %v\n", err)
		return
	}

	// Flush to ensure data is sent immediately
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

// sendKeepalive sends a keepalive comment to maintain connection.
func (h *SSEHandler) sendKeepalive(w http.ResponseWriter) {
	if _, err := fmt.Fprintf(w, ": keepalive\n\n"); err != nil {
		// Connection likely closed, will be detected by clientGone channel
		return
	}
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

