package webui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	webuiatoms "github.com/jaypaulb/CanvusPowerToys/internal/atoms/webui"
)

// validateZoneRequest validates a zone request and returns the canvas ID.
func (h *MacrosHandler) validateZoneRequest(w http.ResponseWriter, r *http.Request, method string) (string, bool) {
	if r.Method != method {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return "", false
	}

	canvasID := h.canvasService.GetCanvasID()
	if canvasID == "" {
		http.Error(w, "Canvas not available", http.StatusServiceUnavailable)
		return "", false
	}

	return canvasID, true
}

// parseZoneIDRequest parses a request body expecting a zoneId field.
func parseZoneIDRequest(r *http.Request) (string, error) {
	var req struct {
		ZoneID string `json:"zoneId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return "", fmt.Errorf("invalid request body: %w", err)
	}
	if req.ZoneID == "" {
		return "", fmt.Errorf("zoneId is required")
	}
	return req.ZoneID, nil
}

// parseZonePairRequest parses a request body expecting sourceZoneId and targetZoneId fields.
func parseZonePairRequest(r *http.Request) (sourceZoneID, targetZoneID string, err error) {
	var req struct {
		SourceZoneID string `json:"sourceZoneId"`
		TargetZoneID string `json:"targetZoneId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return "", "", fmt.Errorf("invalid request body: %w", err)
	}
	if req.SourceZoneID == "" || req.TargetZoneID == "" {
		return "", "", fmt.Errorf("sourceZoneId and targetZoneId are required")
	}
	return req.SourceZoneID, req.TargetZoneID, nil
}

// sendJSONResponse sends a JSON response with CORS headers.
func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// sendErrorResponse sends an error response.
func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	http.Error(w, message, statusCode)
}

// moveWidgets moves widgets from source zone to target zone.
func (h *MacrosHandler) moveWidgets(canvasID, sourceZoneID, targetZoneID string) (int, error) {
	ops := NewMacrosOperations(h.apiClient, h.canvasService)

	// Get zone bounding boxes
	sourceBB, err := webuiatoms.GetZoneBoundingBox(h.apiClient, canvasID, sourceZoneID)
	if err != nil {
		return 0, fmt.Errorf("failed to get source zone: %w", err)
	}

	targetBB, err := webuiatoms.GetZoneBoundingBox(h.apiClient, canvasID, targetZoneID)
	if err != nil {
		return 0, fmt.Errorf("failed to get target zone: %w", err)
	}

	// Get all widgets
	allWidgets, err := webuiatoms.GetAllWidgets(h.apiClient, canvasID)
	if err != nil {
		return 0, fmt.Errorf("failed to get widgets: %w", err)
	}

	// Filter widgets in source zone
	toMove := FilterWidgetsInZone(allWidgets, sourceBB, "")

	// Transform and update each widget
	var updates []WidgetUpdate
	for _, widget := range toMove {
		cloned := widget
		webuiatoms.TransformWidgetLocationAndScale(&cloned, sourceBB, targetBB)
		updates = append(updates, WidgetUpdate{
			WidgetID: widget.ID,
			Payload: map[string]interface{}{
				"location": cloned.Location,
				"scale":    cloned.Scale,
			},
		})
	}

	return ops.BatchUpdateWidgets(canvasID, updates), nil
}

// copyWidgets copies widgets from source zone to target zone.
func (h *MacrosHandler) copyWidgets(canvasID, sourceZoneID, targetZoneID string) (int, error) {
	// Get zone bounding boxes
	sourceBB, err := webuiatoms.GetZoneBoundingBox(h.apiClient, canvasID, sourceZoneID)
	if err != nil {
		return 0, fmt.Errorf("failed to get source zone: %w", err)
	}

	targetBB, err := webuiatoms.GetZoneBoundingBox(h.apiClient, canvasID, targetZoneID)
	if err != nil {
		return 0, fmt.Errorf("failed to get target zone: %w", err)
	}

	// Get all widgets
	allWidgets, err := webuiatoms.GetAllWidgets(h.apiClient, canvasID)
	if err != nil {
		return 0, fmt.Errorf("failed to get widgets: %w", err)
	}

	// Identify widgets in source zone
	inZoneSet := make(map[string]bool)
	for _, w := range allWidgets {
		wt := strings.ToLower(w.WidgetType)
		if wt == "anchor" {
			continue
		}
		if wt != "connector" && w.Location != nil && webuiatoms.WidgetIsInZone(&w, sourceBB) {
			inZoneSet[w.ID] = true
		}
	}

	// Build widgets to copy
	var toCopy []webuiatoms.Widget
	for _, w := range allWidgets {
		wt := strings.ToLower(w.WidgetType)
		if wt == "anchor" || wt == "connector" {
			continue
		}
		if w.Location != nil && inZoneSet[w.ID] {
			toCopy = append(toCopy, w)
		}
	}

	// Copy widgets (create new widgets with transformed locations)
	copiedCount := 0
	for _, widget := range toCopy {
		cloned := widget
		cloned.ID = "" // Clear ID for new widget
		webuiatoms.TransformWidgetLocationAndScale(&cloned, sourceBB, targetBB)

		// Create widget payload
		payload := map[string]interface{}{
			"widget_type": cloned.WidgetType,
			"location":    cloned.Location,
			"scale":       cloned.Scale,
		}
		if cloned.Size != nil {
			payload["size"] = cloned.Size
		}
		if cloned.Title != "" {
			payload["title"] = cloned.Title
		}

		endpoint := fmt.Sprintf("/api/v1/canvases/%s/widgets", canvasID)
		_, err := h.apiClient.Post(endpoint, payload)
		if err == nil {
			copiedCount++
		}
	}

	return copiedCount, nil
}

// pinWidgetsInZone pins or unpins all widgets in a zone.
func (h *MacrosHandler) pinWidgetsInZone(canvasID, zoneID string, pinned bool) (int, error) {
	ops := NewMacrosOperations(h.apiClient, h.canvasService)

	zoneBB, allWidgets, err := ops.GetZoneAndWidgets(zoneID)
	if err != nil {
		return 0, err
	}

	// Filter widgets in zone
	inZone := FilterWidgetsInZone(allWidgets, zoneBB, zoneID)

	// Update widgets
	var updates []WidgetUpdate
	for _, widget := range inZone {
		updates = append(updates, WidgetUpdate{
			WidgetID: widget.ID,
			Payload:  map[string]interface{}{"pinned": pinned},
		})
	}

	return ops.BatchUpdateWidgets(canvasID, updates), nil
}

// organizeWidgetsInGrid organizes widgets in a grid within a zone.
func (h *MacrosHandler) organizeWidgetsInGrid(canvasID, zoneID string) (int, error) {
	ops := NewMacrosOperations(h.apiClient, h.canvasService)

	zoneBB, allWidgets, err := ops.GetZoneAndWidgets(zoneID)
	if err != nil {
		return 0, err
	}

	// Filter widgets in zone
	inZone := FilterWidgetsInZone(allWidgets, zoneBB, zoneID)

	if len(inZone) == 0 {
		return 0, nil
	}

	// Determine optimal grid size
	bestRows, bestCols := CalculateOptimalGrid(len(inZone), zoneBB)
	cellWidth, cellHeight := CalculateCellDimensions(zoneBB, bestRows, bestCols)

	// Position widgets in grid
	var updates []WidgetUpdate
	buffer := 100.0
	for i, widget := range inZone {
		row := i / bestCols
		col := i % bestCols
		x := zoneBB.X + buffer + float64(col)*(cellWidth+buffer)
		y := zoneBB.Y + buffer + float64(row)*(cellHeight+buffer)

		updates = append(updates, WidgetUpdate{
			WidgetID: widget.ID,
			Payload: map[string]interface{}{
				"location": map[string]float64{"x": x, "y": y},
			},
		})
	}

	return ops.BatchUpdateWidgets(canvasID, updates), nil
}

// groupWidgetsByAttribute groups widgets by an attribute (color or title) and positions them.
func (h *MacrosHandler) groupWidgetsByAttribute(canvasID, zoneID string, getAttribute func(webuiatoms.Widget) string) (int, error) {
	ops := NewMacrosOperations(h.apiClient, h.canvasService)

	zoneBB, allWidgets, err := ops.GetZoneAndWidgets(zoneID)
	if err != nil {
		return 0, err
	}

	// Filter and group widgets
	groups := make(map[string][]webuiatoms.Widget)
	for _, w := range allWidgets {
		if w.ID == zoneID {
			continue
		}
		if w.Location == nil {
			continue
		}
		wt := strings.ToLower(w.WidgetType)
		if wt == "connector" || wt == "anchor" {
			continue
		}
		if webuiatoms.WidgetIsInZone(&w, zoneBB) {
			attr := getAttribute(w)
			if attr == "" {
				attr = "default"
			}
			groups[attr] = append(groups[attr], w)
		}
	}

	return ops.PositionWidgetGroups(groups, zoneBB, canvasID), nil
}

