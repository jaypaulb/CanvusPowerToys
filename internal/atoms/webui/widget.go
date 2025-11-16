package webui

import (
	"encoding/json"
	"fmt"
	"strings"
)

// WidgetLocation represents widget location.
type WidgetLocation struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// WidgetSize represents widget size.
type WidgetSize struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// Widget represents a widget from the Canvus API.
type Widget struct {
	ID         string                 `json:"id"`
	WidgetType string                 `json:"widget_type"`
	Location   *WidgetLocation        `json:"location,omitempty"`
	Size       *WidgetSize            `json:"size,omitempty"`
	Scale      float64                `json:"scale,omitempty"`
	Pinned     bool                   `json:"pinned,omitempty"`
	Title      string                 `json:"title,omitempty"`
	Color      string                 `json:"color,omitempty"`
	Data       map[string]interface{} `json:"-"`
}

// GetAllWidgets fetches all widgets from the canvas.
func GetAllWidgets(apiClient *APIClient, canvasID string) ([]Widget, error) {
	if canvasID == "" {
		return nil, fmt.Errorf("canvas not available")
	}

	endpoint := fmt.Sprintf("/api/v1/canvases/%s/widgets", canvasID)
	fmt.Printf("[GetAllWidgets] Fetching widgets from %s\n", endpoint)
	data, err := apiClient.Get(endpoint)
	if err != nil {
		fmt.Printf("[GetAllWidgets] ERROR: Failed to get widgets: %v\n", err)
		return nil, fmt.Errorf("failed to get widgets: %w", err)
	}

	var widgets []Widget
	if err := json.Unmarshal(data, &widgets); err != nil {
		fmt.Printf("[GetAllWidgets] ERROR: Failed to parse widgets JSON: %v\n", err)
		fmt.Printf("[GetAllWidgets] Raw response length: %d bytes\n", len(data))
		return nil, fmt.Errorf("failed to parse widgets: %w", err)
	}

	fmt.Printf("[GetAllWidgets] Retrieved %d widgets\n", len(widgets))

	// Log widget types breakdown
	typeCounts := make(map[string]int)
	for _, w := range widgets {
		wt := w.WidgetType
		if wt == "" {
			wt = "unknown"
		}
		typeCounts[wt]++
	}
	fmt.Printf("[GetAllWidgets] Widget types: %v\n", typeCounts)

	return widgets, nil
}

// TransformWidgetLocationAndScale transforms widget location and scale from source to target zone.
func TransformWidgetLocationAndScale(widget *Widget, sourceBB, targetBB *ZoneBoundingBox) {
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

// GetWidgetPatchEndpoint returns the API endpoint for patching a widget based on widget_type.
func GetWidgetPatchEndpoint(widgetType string) string {
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

