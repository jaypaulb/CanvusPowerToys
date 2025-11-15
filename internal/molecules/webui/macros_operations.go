package webui

import (
	"fmt"
	"strings"

	webuiatoms "github.com/jaypaulb/CanvusPowerToys/internal/atoms/webui"
)

// MacrosOperations provides reusable operations for macros functionality.
type MacrosOperations struct {
	apiClient     *webuiatoms.APIClient
	canvasService *CanvasService
}

// NewMacrosOperations creates a new macros operations helper.
func NewMacrosOperations(apiClient *webuiatoms.APIClient, canvasService *CanvasService) *MacrosOperations {
	return &MacrosOperations{
		apiClient:     apiClient,
		canvasService: canvasService,
	}
}

// GetZoneAndWidgets gets zone bounding box and all widgets in that zone.
func (mo *MacrosOperations) GetZoneAndWidgets(zoneID string) (*webuiatoms.ZoneBoundingBox, []webuiatoms.Widget, error) {
	canvasID := mo.canvasService.GetCanvasID()
	if canvasID == "" {
		return nil, nil, fmt.Errorf("canvas not available")
	}

	zoneBB, err := webuiatoms.GetZoneBoundingBox(mo.apiClient, canvasID, zoneID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get zone: %w", err)
	}

	allWidgets, err := webuiatoms.GetAllWidgets(mo.apiClient, canvasID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get widgets: %w", err)
	}

	return zoneBB, allWidgets, nil
}

// FilterWidgetsInZone filters widgets that are within a zone, excluding anchors and connectors.
func FilterWidgetsInZone(widgets []webuiatoms.Widget, zoneBB *webuiatoms.ZoneBoundingBox, excludeZoneID string) []webuiatoms.Widget {
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

// UpdateWidgetWithRetry updates a widget with retry logic (up to 3 attempts).
func (mo *MacrosOperations) UpdateWidgetWithRetry(canvasID, widgetID string, payload map[string]interface{}) error {
	endpoint := fmt.Sprintf("/api/v1/canvases/%s/widgets/%s", canvasID, widgetID)
	var lastErr error
	for tries := 0; tries < 3; tries++ {
		_, err := mo.apiClient.Patch(endpoint, payload)
		if err == nil {
			return nil
		}
		lastErr = err
	}
	return fmt.Errorf("failed after 3 attempts: %w", lastErr)
}

// BatchUpdateWidgets updates multiple widgets and returns count of successful updates.
func (mo *MacrosOperations) BatchUpdateWidgets(canvasID string, updates []WidgetUpdate) int {
	successCount := 0
	for _, update := range updates {
		if err := mo.UpdateWidgetWithRetry(canvasID, update.WidgetID, update.Payload); err == nil {
			successCount++
		}
	}
	return successCount
}

// WidgetUpdate represents a widget update operation.
type WidgetUpdate struct {
	WidgetID string
	Payload  map[string]interface{}
}

// CalculateOptimalGrid calculates the optimal grid dimensions for n widgets in a zone.
func CalculateOptimalGrid(n int, zoneBB *webuiatoms.ZoneBoundingBox) (rows, cols int) {
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

// CalculateCellDimensions calculates cell width and height for a grid.
func CalculateCellDimensions(zoneBB *webuiatoms.ZoneBoundingBox, rows, cols int) (width, height float64) {
	buffer := 100.0
	width = (zoneBB.Width - buffer*float64(cols+1)) / float64(cols)
	height = (zoneBB.Height - buffer*float64(rows+1)) / float64(rows)
	return width, height
}

// PositionWidgetGroups positions widget groups horizontally with vertical stacking within groups.
func (mo *MacrosOperations) PositionWidgetGroups(groups map[string][]webuiatoms.Widget, zoneBB *webuiatoms.ZoneBoundingBox, canvasID string) int {
	groupedCount := 0
	xOffset := zoneBB.X + 100
	for _, widgets := range groups {
		yOffset := zoneBB.Y + 100
		for _, widget := range widgets {
			endpoint := fmt.Sprintf("/api/v1/canvases/%s/widgets/%s", canvasID, widget.ID)
			payload := map[string]interface{}{
				"location": map[string]float64{"x": xOffset, "y": yOffset},
			}
			_, err := mo.apiClient.Patch(endpoint, payload)
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

