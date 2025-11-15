package webui

import (
	"encoding/json"
	"fmt"
)

// ZoneBoundingBox represents a zone's bounding box from an anchor.
type ZoneBoundingBox struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Scale  float64 `json:"scale"`
}

// GetZoneBoundingBox gets the bounding box for a zone (anchor) from the Canvus API.
func GetZoneBoundingBox(apiClient *APIClient, canvasID, zoneID string) (*ZoneBoundingBox, error) {
	if canvasID == "" {
		return nil, fmt.Errorf("canvas not available")
	}

	endpoint := fmt.Sprintf("/api/v1/canvases/%s/anchors/%s", canvasID, zoneID)
	data, err := apiClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get anchor: %w", err)
	}

	var anchor struct {
		Location *struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
		} `json:"location"`
		Size *struct {
			Width  float64 `json:"width"`
			Height float64 `json:"height"`
		} `json:"size"`
		Scale float64 `json:"scale"`
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

// WidgetIsInZone checks if a widget is within a zone bounding box.
func WidgetIsInZone(widget *Widget, zoneBB *ZoneBoundingBox) bool {
	if widget.Location == nil {
		return false
	}
	wx := widget.Location.X
	wy := widget.Location.Y
	withinX := wx >= zoneBB.X+2 && wx <= zoneBB.X+zoneBB.Width-2
	withinY := wy >= zoneBB.Y+2 && wy <= zoneBB.Y+zoneBB.Height-2
	return withinX && withinY
}

