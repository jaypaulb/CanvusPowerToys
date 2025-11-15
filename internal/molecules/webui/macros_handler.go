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

// ZoneBoundingBox represents a zone's bounding box from an anchor.
type ZoneBoundingBox struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Scale  float64 `json:"scale"`
}

// Widget represents a widget from the Canvus API.
type Widget struct {
	ID         string                 `json:"id"`
	WidgetType string                 `json:"widget_type"`
	Location   *Location              `json:"location,omitempty"`
	Size       *Size                  `json:"size,omitempty"`
	Scale      float64                `json:"scale,omitempty"`
	Pinned     bool                   `json:"pinned,omitempty"`
	Title      string                 `json:"title,omitempty"`
	Color      string                 `json:"color,omitempty"`
	Data       map[string]interface{} `json:"-"`
}

// Location represents widget location.
type Location struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Size represents widget size.
type Size struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// getZoneBoundingBox gets the bounding box for a zone (anchor).
func (h *MacrosHandler) getZoneBoundingBox(zoneID string) (*ZoneBoundingBox, error) {
	canvasID := h.canvasService.GetCanvasID()
	if canvasID == "" {
		return nil, fmt.Errorf("canvas not available")
	}

	endpoint := fmt.Sprintf("/api/v1/canvases/%s/anchors/%s", canvasID, zoneID)
	data, err := h.apiClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get anchor: %w", err)
	}

	var anchor struct {
		Location *Location `json:"location"`
		Size     *Size     `json:"size"`
		Scale    float64   `json:"scale"`
	}
	if err := json.Unmarshal(data, &anchor); err != nil {
		return nil, fmt.Errorf("failed to parse anchor: %w", err)
	}

	if anchor.Location == nil || anchor.Size == nil {
		return nil, fmt.Errorf("invalid anchor data for zone ID: %s", zoneID)
	}

	return &ZoneBoundingBox{
		X:      anchor.Location.X,
		Y:      anchor.Location.Y,
		Width:  anchor.Size.Width,
		Height: anchor.Size.Height,
		Scale:  anchor.Scale,
	}, nil
}

// getAllWidgets fetches all widgets from the canvas.
func (h *MacrosHandler) getAllWidgets() ([]Widget, error) {
	canvasID := h.canvasService.GetCanvasID()
	if canvasID == "" {
		return nil, fmt.Errorf("canvas not available")
	}

	endpoint := fmt.Sprintf("/api/v1/canvases/%s/widgets", canvasID)
	data, err := h.apiClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get widgets: %w", err)
	}

	var widgets []Widget
	if err := json.Unmarshal(data, &widgets); err != nil {
		return nil, fmt.Errorf("failed to parse widgets: %w", err)
	}

	return widgets, nil
}

// widgetIsInZone checks if a widget is within a zone bounding box.
func widgetIsInZone(widget *Widget, zoneBB *ZoneBoundingBox) bool {
	if widget.Location == nil {
		return false
	}
	wx := widget.Location.X
	wy := widget.Location.Y
	withinX := wx >= zoneBB.X+2 && wx <= zoneBB.X+zoneBB.Width-2
	withinY := wy >= zoneBB.Y+2 && wy <= zoneBB.Y+zoneBB.Height-2
	return withinX && withinY
}

// transformWidgetLocationAndScale transforms widget location and scale from source to target zone.
func transformWidgetLocationAndScale(widget *Widget, sourceBB, targetBB *ZoneBoundingBox) {
	if widget.Location == nil {
		return
	}
	scaleFactor := targetBB.Width / sourceBB.Width
	deltaX := widget.Location.X - sourceBB.X
	deltaY := widget.Location.Y - sourceBB.Y
	widget.Location.X = targetBB.X + deltaX*scaleFactor
	widget.Location.Y = targetBB.Y + deltaY*scaleFactor
	oldScale := widget.Scale
	if oldScale == 0 {
		oldScale = 1
	}
	widget.Scale = oldScale * scaleFactor
}

// getWidgetPatchEndpoint returns the API endpoint for patching a widget based on widget_type.
func getWidgetPatchEndpoint(widgetType string) string {
	wt := strings.ToLower(widgetType)
	switch wt {
	case "note":
		return "/notes"
	case "browser":
		return "/browsers"
	case "image":
		return "/images"
	case "pdf":
		return "/pdfs"
	case "video":
		return "/videos"
	case "connector":
		return "/connectors"
	case "anchor":
		return "/anchors"
	default:
		return "/notes" // Default fallback
	}
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
	sourceBB, err := h.getZoneBoundingBox(req.SourceZoneID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get source zone: %v", err), http.StatusInternalServerError)
		return
	}

	targetBB, err := h.getZoneBoundingBox(req.TargetZoneID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get target zone: %v", err), http.StatusInternalServerError)
		return
	}

	// Get all widgets
	allWidgets, err := h.getAllWidgets()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get widgets: %v", err), http.StatusInternalServerError)
		return
	}

	// Filter widgets in source zone (exclude anchors and connectors)
	var toMove []Widget
	for _, w := range allWidgets {
		wt := strings.ToLower(w.WidgetType)
		if wt == "connector" || wt == "anchor" {
			continue
		}
		if w.Location == nil {
			continue
		}
		if widgetIsInZone(&w, sourceBB) {
			toMove = append(toMove, w)
		}
	}

	// Transform and update each widget
	movedCount := 0
	for _, widget := range toMove {
		cloned := widget
		transformWidgetLocationAndScale(&cloned, sourceBB, targetBB)

		// PATCH widget via widgets API (simpler - use widgets endpoint directly)
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
			if tries < 2 {
				// Wait 500ms before retry (would need time.Sleep, but keeping it simple for now)
			}
		}
		if lastErr != nil && movedCount == 0 {
			// If first widget fails, return error
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
	sourceBB, err := h.getZoneBoundingBox(req.SourceZoneID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get source zone: %v", err), http.StatusInternalServerError)
		return
	}

	targetBB, err := h.getZoneBoundingBox(req.TargetZoneID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get target zone: %v", err), http.StatusInternalServerError)
		return
	}

	// Get all widgets
	allWidgets, err := h.getAllWidgets()
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
		if wt != "connector" && w.Location != nil && widgetIsInZone(&w, sourceBB) {
			inZoneSet[w.ID] = true
		}
	}

	// Build widgets to copy (normal widgets + connectors with both endpoints in zone)
	var toCopy []Widget
	for _, w := range allWidgets {
		wt := strings.ToLower(w.WidgetType)
		if wt == "anchor" {
			continue
		}
		if wt == "connector" {
			// For connectors, check if both endpoints are in zone (simplified - would need src/dst parsing)
			// For now, skip connectors in copy
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
		transformWidgetLocationAndScale(&cloned, sourceBB, targetBB)

		// Create widget payload (simplified - would need full widget structure)
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
	allWidgets, err := h.getAllWidgets()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get widgets: %v", err), http.StatusInternalServerError)
		return
	}

	var pinned []Widget
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
	zoneBB, err := h.getZoneBoundingBox(req.ZoneID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get zone: %v", err), http.StatusInternalServerError)
		return
	}

	// Get all widgets
	allWidgets, err := h.getAllWidgets()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get widgets: %v", err), http.StatusInternalServerError)
		return
	}

	// Filter widgets in zone (exclude anchors and connectors)
	var inZone []Widget
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
		if widgetIsInZone(&w, zoneBB) {
			inZone = append(inZone, w)
		}
	}

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
	zoneBB, err := h.getZoneBoundingBox(req.ZoneID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get zone: %v", err), http.StatusInternalServerError)
		return
	}

	// Get all widgets
	allWidgets, err := h.getAllWidgets()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get widgets: %v", err), http.StatusInternalServerError)
		return
	}

	// Filter widgets in zone (exclude anchors and connectors)
	var inZone []Widget
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
		if widgetIsInZone(&w, zoneBB) {
			inZone = append(inZone, w)
		}
	}

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
	zoneBB, err := h.getZoneBoundingBox(req.ZoneID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get zone: %v", err), http.StatusInternalServerError)
		return
	}

	// Get all widgets
	allWidgets, err := h.getAllWidgets()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get widgets: %v", err), http.StatusInternalServerError)
		return
	}

	// Filter widgets in zone
	var inZone []Widget
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
		if widgetIsInZone(&w, zoneBB) {
			inZone = append(inZone, w)
		}
	}

	if len(inZone) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "No widgets found to auto-grid",
		})
		return
	}

	// Determine grid size (simplified - calculate optimal rows/cols)
	n := len(inZone)
	aspectRatio := zoneBB.Width / zoneBB.Height
	buffer := 100.0

	bestRows := 1
	bestCols := n
	minEmptySpace := 1e10

	for rows := 1; rows <= n; rows++ {
		cols := (n + rows - 1) / rows // Ceiling division
		// Calculate cell dimensions for aspect ratio comparison
		_ = (zoneBB.Width - buffer*float64(cols+1)) / float64(cols)  // cellWidth (for calculation)
		_ = (zoneBB.Height - buffer*float64(rows+1)) / float64(rows) // cellHeight (for calculation)
		gridAspectRatio := float64(cols) / float64(rows)
		emptySpace := abs(aspectRatio - gridAspectRatio)

		if emptySpace < minEmptySpace {
			minEmptySpace = emptySpace
			bestRows = rows
			bestCols = cols
		}
	}

	// Calculate cell dimensions
	cellWidth := (zoneBB.Width - buffer*float64(bestCols+1)) / float64(bestCols)
	cellHeight := (zoneBB.Height - buffer*float64(bestRows+1)) / float64(bestRows)

	// Position widgets in grid
	griddedCount := 0
	for i, widget := range inZone {
		row := i / bestCols
		col := i % bestCols
		x := zoneBB.X + buffer + float64(col)*(cellWidth+buffer)
		y := zoneBB.Y + buffer + float64(row)*(cellHeight+buffer)

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
	zoneBB, err := h.getZoneBoundingBox(req.ZoneID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get zone: %v", err), http.StatusInternalServerError)
		return
	}

	// Get all widgets
	allWidgets, err := h.getAllWidgets()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get widgets: %v", err), http.StatusInternalServerError)
		return
	}

	// Filter and group widgets by color
	colorGroups := make(map[string][]Widget)
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
		if widgetIsInZone(&w, zoneBB) {
			color := w.Color
			if color == "" {
				color = "default"
			}
			colorGroups[color] = append(colorGroups[color], w)
		}
	}

	// Position widgets by color groups (simplified - arrange horizontally)
	groupedCount := 0
	xOffset := zoneBB.X + 100
	for color, widgets := range colorGroups {
		_ = color // Color group identifier
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
			yOffset += 200 // Stack vertically
		}
		xOffset += 300 // Next color group to the right
	}

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
	zoneBB, err := h.getZoneBoundingBox(req.ZoneID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get zone: %v", err), http.StatusInternalServerError)
		return
	}

	// Get all widgets
	allWidgets, err := h.getAllWidgets()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get widgets: %v", err), http.StatusInternalServerError)
		return
	}

	// Filter widgets in zone
	var inZone []Widget
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
		if widgetIsInZone(&w, zoneBB) {
			inZone = append(inZone, w)
		}
	}

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
	titleGroups := make(map[string][]Widget)
	for _, w := range inZone {
		title := w.Title
		if title == "" {
			title = "untitled"
		}
		titleGroups[title] = append(titleGroups[title], w)
	}

	// Position widgets by title groups
	groupedCount := 0
	xOffset := zoneBB.X + 100
	for _, widgets := range titleGroups {
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

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("%d widgets grouped by title", groupedCount),
	})
}

// abs returns absolute value of a float64.
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
