package webui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	webuiatoms "github.com/jaypaulb/CanvusPowerToys/internal/atoms/webui"
)

// MacrosHandler handles macros management API endpoints.
// NOTE: Macros are server-side operations that operate on widgets via the widgets API.
// They are NOT a Canvus API resource - they're abstractions built on top of widgets/anchors.
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

// filterWidgetsInZone filters widgets that are within a zone, excluding anchors and connectors.
func (h *MacrosHandler) filterWidgetsInZone(widgets []webuiatoms.Widget, zoneBB *webuiatoms.ZoneBoundingBox, excludeZoneID string) []webuiatoms.Widget {
	var filtered []webuiatoms.Widget
	for _, w := range widgets {
		if w.ID == excludeZoneID {
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
			filtered = append(filtered, w)
		}
	}
	return filtered
}

// HandleMove handles POST /api/macros/move - Move widgets from source zone to target zone.
func (h *MacrosHandler) HandleMove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SourceZoneID string `json:"sourceZoneId"`
		TargetZoneID string `json:"targetZoneId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.SourceZoneID == "" || req.TargetZoneID == "" {
		http.Error(w, "sourceZoneId and targetZoneId are required", http.StatusBadRequest)
		return
	}

	canvasID := h.canvasService.GetCanvasID()
	if canvasID == "" {
		http.Error(w, "Canvas not available", http.StatusServiceUnavailable)
		return
	}

	// Get zone bounding boxes
	sourceBB, err := webuiatoms.GetZoneBoundingBox(h.apiClient, canvasID, req.SourceZoneID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get source zone: %v", err), http.StatusInternalServerError)
		return
	}

	targetBB, err := webuiatoms.GetZoneBoundingBox(h.apiClient, canvasID, req.TargetZoneID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get target zone: %v", err), http.StatusInternalServerError)
		return
	}

	// Get all widgets
	allWidgets, err := webuiatoms.GetAllWidgets(h.apiClient, canvasID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get widgets: %v", err), http.StatusInternalServerError)
		return
	}

	// Filter widgets in source zone
	toMove := h.filterWidgetsInZone(allWidgets, sourceBB, "")

	// Transform and update each widget
	movedCount := 0
	for _, widget := range toMove {
		cloned := widget
		webuiatoms.TransformWidgetLocationAndScale(&cloned, sourceBB, targetBB)

		// PATCH widget via widgets API
		endpoint := fmt.Sprintf("/api/v1/canvases/%s/widgets/%s", canvasID, widget.ID)
		payload := map[string]interface{}{
			"location": cloned.Location,
			"scale":    cloned.Scale,
		}

		// Retry up to 3 times
		var lastErr error
		for tries := 0; tries < 3; tries++ {
			_, err := h.apiClient.Patch(endpoint, payload)
			if err == nil {
				movedCount++
				break
			}
			lastErr = err
		}
		if lastErr != nil && movedCount == 0 {
			http.Error(w, fmt.Sprintf("Failed to move widgets: %v", lastErr), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("%d widgets moved", movedCount),
	})
}

// HandleCopy handles POST /api/macros/copy - Copy widgets from source zone to target zone.
func (h *MacrosHandler) HandleCopy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SourceZoneID string `json:"sourceZoneId"`
		TargetZoneID string `json:"targetZoneId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.SourceZoneID == "" || req.TargetZoneID == "" {
		http.Error(w, "sourceZoneId and targetZoneId are required", http.StatusBadRequest)
		return
	}

	canvasID := h.canvasService.GetCanvasID()
	if canvasID == "" {
		http.Error(w, "Canvas not available", http.StatusServiceUnavailable)
		return
	}

	// Get zone bounding boxes
	sourceBB, err := webuiatoms.GetZoneBoundingBox(h.apiClient, canvasID, req.SourceZoneID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get source zone: %v", err), http.StatusInternalServerError)
		return
	}

	targetBB, err := webuiatoms.GetZoneBoundingBox(h.apiClient, canvasID, req.TargetZoneID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get target zone: %v", err), http.StatusInternalServerError)
		return
	}

	// Get all widgets
	allWidgets, err := webuiatoms.GetAllWidgets(h.apiClient, canvasID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get widgets: %v", err), http.StatusInternalServerError)
		return
	}

	// Identify widgets in source zone (excluding anchors, including connectors only if both endpoints in zone)
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

	// Build widgets to copy (normal widgets + connectors with both endpoints in zone)
	var toCopy []webuiatoms.Widget
	for _, w := range allWidgets {
		wt := strings.ToLower(w.WidgetType)
		if wt == "anchor" {
			continue
		}
		if wt == "connector" {
			// For connectors, check if both endpoints are in zone (simplified - would need src/dst parsing)
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

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("%d widgets copied", copiedCount),
	})
}

// HandleGroups handles GET /api/macros/groups - List widget groups (computed from widgets).
func (h *MacrosHandler) HandleGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Groups are computed from widgets - for now return empty array
	// Full implementation would group widgets by color, title, etc.
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode([]interface{}{})
}

// HandlePinned handles GET /api/macros/pinned - List pinned widgets.
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

	// Get all widgets and filter pinned ones
	allWidgets, err := webuiatoms.GetAllWidgets(h.apiClient, canvasID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get widgets: %v", err), http.StatusInternalServerError)
		return
	}

	var pinned []webuiatoms.Widget
	for _, w := range allWidgets {
		if w.Pinned {
			pinned = append(pinned, w)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(pinned)
}

// HandleUnpin handles POST /api/macros/unpin-all - Unpin all widgets in a zone.
func (h *MacrosHandler) HandleUnpin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ZoneID string `json:"zoneId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ZoneID == "" {
		http.Error(w, "zoneId is required", http.StatusBadRequest)
		return
	}

	canvasID := h.canvasService.GetCanvasID()
	if canvasID == "" {
		http.Error(w, "Canvas not available", http.StatusServiceUnavailable)
		return
	}

	// Get zone bounding box
	zoneBB, err := webuiatoms.GetZoneBoundingBox(h.apiClient, canvasID, req.ZoneID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get zone: %v", err), http.StatusInternalServerError)
		return
	}

	// Get all widgets
	allWidgets, err := webuiatoms.GetAllWidgets(h.apiClient, canvasID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get widgets: %v", err), http.StatusInternalServerError)
		return
	}

	// Filter widgets in zone
	inZone := h.filterWidgetsInZone(allWidgets, zoneBB, req.ZoneID)

	// Unpin all widgets in zone
	unpinnedCount := 0
	for _, widget := range inZone {
		endpoint := fmt.Sprintf("/api/v1/canvases/%s/widgets/%s", canvasID, widget.ID)
		payload := map[string]interface{}{"pinned": false}
		_, err := h.apiClient.Patch(endpoint, payload)
		if err == nil {
			unpinnedCount++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("%d widgets unpinned", unpinnedCount),
	})
}

// HandlePinAll handles POST /api/macros/pin-all - Pin all widgets in a zone.
func (h *MacrosHandler) HandlePinAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ZoneID string `json:"zoneId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ZoneID == "" {
		http.Error(w, "zoneId is required", http.StatusBadRequest)
		return
	}

	canvasID := h.canvasService.GetCanvasID()
	if canvasID == "" {
		http.Error(w, "Canvas not available", http.StatusServiceUnavailable)
		return
	}

	// Get zone bounding box
	zoneBB, err := webuiatoms.GetZoneBoundingBox(h.apiClient, canvasID, req.ZoneID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get zone: %v", err), http.StatusInternalServerError)
		return
	}

	// Get all widgets
	allWidgets, err := webuiatoms.GetAllWidgets(h.apiClient, canvasID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get widgets: %v", err), http.StatusInternalServerError)
		return
	}

	// Filter widgets in zone
	inZone := h.filterWidgetsInZone(allWidgets, zoneBB, req.ZoneID)

	// Pin all widgets in zone
	pinnedCount := 0
	for _, widget := range inZone {
		endpoint := fmt.Sprintf("/api/v1/canvases/%s/widgets/%s", canvasID, widget.ID)
		payload := map[string]interface{}{"pinned": true}
		_, err := h.apiClient.Patch(endpoint, payload)
		if err == nil {
			pinnedCount++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("%d widgets pinned", pinnedCount),
	})
}

// HandleAutoGrid handles POST /api/macros/auto-grid - Organize widgets in a grid within a zone.
func (h *MacrosHandler) HandleAutoGrid(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ZoneID string `json:"zoneId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ZoneID == "" {
		http.Error(w, "zoneId is required", http.StatusBadRequest)
		return
	}

	canvasID := h.canvasService.GetCanvasID()
	if canvasID == "" {
		http.Error(w, "Canvas not available", http.StatusServiceUnavailable)
		return
	}

	// Get zone bounding box
	zoneBB, err := webuiatoms.GetZoneBoundingBox(h.apiClient, canvasID, req.ZoneID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get zone: %v", err), http.StatusInternalServerError)
		return
	}

	// Get all widgets
	allWidgets, err := webuiatoms.GetAllWidgets(h.apiClient, canvasID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get widgets: %v", err), http.StatusInternalServerError)
		return
	}

	// Filter widgets in zone
	inZone := h.filterWidgetsInZone(allWidgets, zoneBB, req.ZoneID)

	if len(inZone) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "No widgets found to auto-grid",
		})
		return
	}

	// Determine optimal grid size
	bestRows, bestCols := h.calculateOptimalGrid(len(inZone), zoneBB)
	cellWidth, cellHeight := h.calculateCellDimensions(zoneBB, bestRows, bestCols)

	// Position widgets in grid
	griddedCount := 0
	for i, widget := range inZone {
		row := i / bestCols
		col := i % bestCols
		x := zoneBB.X + 100 + float64(col)*(cellWidth+100)
		y := zoneBB.Y + 100 + float64(row)*(cellHeight+100)

		endpoint := fmt.Sprintf("/api/v1/canvases/%s/widgets/%s", canvasID, widget.ID)
		payload := map[string]interface{}{
			"location": map[string]float64{"x": x, "y": y},
		}
		_, err := h.apiClient.Patch(endpoint, payload)
		if err == nil {
			griddedCount++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("%d widgets organized in grid", griddedCount),
	})
}

// calculateOptimalGrid calculates the optimal grid dimensions for n widgets in a zone.
func (h *MacrosHandler) calculateOptimalGrid(n int, zoneBB *webuiatoms.ZoneBoundingBox) (rows, cols int) {
	aspectRatio := zoneBB.Width / zoneBB.Height

	bestRows := 1
	bestCols := n
	minEmptySpace := 1e10

	for r := 1; r <= n; r++ {
		c := (n + r - 1) / r // Ceiling division
		gridAspectRatio := float64(c) / float64(r)
		emptySpace := abs(aspectRatio - gridAspectRatio)

		if emptySpace < minEmptySpace {
			minEmptySpace = emptySpace
			bestRows = r
			bestCols = c
		}
	}

	return bestRows, bestCols
}

// calculateCellDimensions calculates cell width and height for a grid.
func (h *MacrosHandler) calculateCellDimensions(zoneBB *webuiatoms.ZoneBoundingBox, rows, cols int) (width, height float64) {
	buffer := 100.0
	width = (zoneBB.Width - buffer*float64(cols+1)) / float64(cols)
	height = (zoneBB.Height - buffer*float64(rows+1)) / float64(rows)
	return width, height
}

// HandleGroupColor handles POST /api/macros/group-color - Group widgets by color.
func (h *MacrosHandler) HandleGroupColor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ZoneID string `json:"zoneId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ZoneID == "" {
		http.Error(w, "zoneId is required", http.StatusBadRequest)
		return
	}

	canvasID := h.canvasService.GetCanvasID()
	if canvasID == "" {
		http.Error(w, "Canvas not available", http.StatusServiceUnavailable)
		return
	}

	// Get zone bounding box
	zoneBB, err := webuiatoms.GetZoneBoundingBox(h.apiClient, canvasID, req.ZoneID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get zone: %v", err), http.StatusInternalServerError)
		return
	}

	// Get all widgets
	allWidgets, err := webuiatoms.GetAllWidgets(h.apiClient, canvasID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get widgets: %v", err), http.StatusInternalServerError)
		return
	}

	// Filter and group widgets by color
	colorGroups := make(map[string][]webuiatoms.Widget)
	for _, w := range allWidgets {
		if w.ID == req.ZoneID {
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
			color := w.Color
			if color == "" {
				color = "default"
			}
			colorGroups[color] = append(colorGroups[color], w)
		}
	}

	// Position widgets by color groups
	groupedCount := h.positionWidgetGroups(colorGroups, zoneBB, canvasID)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("%d widgets grouped by color", groupedCount),
	})
}

// HandleGroupTitle handles POST /api/macros/group-title - Group widgets by title.
func (h *MacrosHandler) HandleGroupTitle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ZoneID string `json:"zoneId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ZoneID == "" {
		http.Error(w, "zoneId is required", http.StatusBadRequest)
		return
	}

	canvasID := h.canvasService.GetCanvasID()
	if canvasID == "" {
		http.Error(w, "Canvas not available", http.StatusServiceUnavailable)
		return
	}

	// Get zone bounding box
	zoneBB, err := webuiatoms.GetZoneBoundingBox(h.apiClient, canvasID, req.ZoneID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get zone: %v", err), http.StatusInternalServerError)
		return
	}

	// Get all widgets
	allWidgets, err := webuiatoms.GetAllWidgets(h.apiClient, canvasID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get widgets: %v", err), http.StatusInternalServerError)
		return
	}

	// Filter widgets in zone
	inZone := h.filterWidgetsInZone(allWidgets, zoneBB, req.ZoneID)

	if len(inZone) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "No widgets found to group by title",
		})
		return
	}

	// Sort by title
	sort.Slice(inZone, func(i, j int) bool {
		return strings.ToLower(inZone[i].Title) < strings.ToLower(inZone[j].Title)
	})

	// Group by title
	titleGroups := make(map[string][]webuiatoms.Widget)
	for _, w := range inZone {
		title := w.Title
		if title == "" {
			title = "untitled"
		}
		titleGroups[title] = append(titleGroups[title], w)
	}

	// Position widgets by title groups
	groupedCount := h.positionWidgetGroups(titleGroups, zoneBB, canvasID)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("%d widgets grouped by title", groupedCount),
	})
}

// positionWidgetGroups positions widget groups horizontally with vertical stacking within groups.
func (h *MacrosHandler) positionWidgetGroups(groups map[string][]webuiatoms.Widget, zoneBB *webuiatoms.ZoneBoundingBox, canvasID string) int {
	groupedCount := 0
	xOffset := zoneBB.X + 100
	for _, widgets := range groups {
		yOffset := zoneBB.Y + 100
		for _, widget := range widgets {
			endpoint := fmt.Sprintf("/api/v1/canvases/%s/widgets/%s", canvasID, widget.ID)
			payload := map[string]interface{}{
				"location": map[string]float64{"x": xOffset, "y": yOffset},
			}
			_, err := h.apiClient.Patch(endpoint, payload)
			if err == nil {
				groupedCount++
			}
			yOffset += 200
		}
		xOffset += 300
	}
	return groupedCount
}

// abs returns absolute value of a float64.
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
