package screenxml

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// ScreenXML represents the screen.xml structure.
type ScreenXML struct {
	XMLName   xml.Name      `xml:"multihead"`
	Type      string        `xml:"type,attr,omitempty"`
	DPI       *XMLAttr      `xml:"dpi,omitempty"`
	DPMS      *XMLAttr       `xml:"dpms,omitempty"`
	HwColorCorrection *XMLAttr `xml:"hw-color-correction,omitempty"`
	Iconify   *XMLAttr       `xml:"iconify,omitempty"`
	Vsync     *XMLAttr       `xml:"vsync,omitempty"`
	LayerSize *XMLAttr       `xml:"layer-size,omitempty"`
	Windows   []WindowConfig `xml:"window"`
}

// XMLAttr represents an XML element with type attribute (e.g., <dpi type="">value</dpi>).
type XMLAttr struct {
	Type  string `xml:"type,attr,omitempty"`
	Value string `xml:",chardata"`
}

// WindowConfig represents a window element in screen.xml.
type WindowConfig struct {
	Type            string        `xml:"type,attr"`
	DirectRendering *XMLAttr      `xml:"direct-rendering,omitempty"`
	Frameless       *XMLAttr      `xml:"frameless,omitempty"`
	FsaaSamples     *XMLAttr      `xml:"fsaa-samples,omitempty"`
	Fullscreen      *XMLAttr      `xml:"fullscreen,omitempty"`
	Location        *XMLAttr      `xml:"location,omitempty"`
	Resizable       *XMLAttr      `xml:"resizable,omitempty"`
	ScreenNumber    *XMLAttr      `xml:"screennumber,omitempty"`
	Size            *XMLAttr      `xml:"size,omitempty"`
	Areas           []AreaConfig  `xml:"area"`
}

// AreaConfig represents an area element within a window.
type AreaConfig struct {
	Type            string        `xml:"type,attr"`
	ColorCorrection  *ColorCorrection `xml:"colorcorrection,omitempty"`
	GraphicsLocation *XMLAttr     `xml:"graphicslocation,omitempty"`
	GraphicsSize     *XMLAttr     `xml:"graphicssize,omitempty"`
	Keystone         *Keystone    `xml:"keystone,omitempty"`
	Location         *XMLAttr     `xml:"location,omitempty"`
	Method           *XMLAttr     `xml:"method,omitempty"`
	RgbCube          *RgbCube     `xml:"rgbcube,omitempty"`
	Seams            *XMLAttr      `xml:"seams,omitempty"`
	Size             *XMLAttr      `xml:"size,omitempty"`
}

// ColorCorrection represents color correction settings.
type ColorCorrection struct {
	Type       string    `xml:"type,attr,omitempty"`
	Brightness *XMLAttr  `xml:"brightness,omitempty"`
	Contrast   *XMLAttr  `xml:"contrast,omitempty"`
	Gamma      *XMLAttr  `xml:"gamma,omitempty"`
	Red        string   `xml:"red,omitempty"`
	Green      string   `xml:"green,omitempty"`
	Blue       string   `xml:"blue,omitempty"`
}

// Keystone represents keystone correction.
type Keystone struct {
	Type      string    `xml:"type,attr,omitempty"`
	Rotations *XMLAttr  `xml:"rotations,omitempty"`
	V1        *XMLAttr  `xml:"v1,omitempty"`
	V2        *XMLAttr  `xml:"v2,omitempty"`
	V3        *XMLAttr  `xml:"v3,omitempty"`
	V4        *XMLAttr  `xml:"v4,omitempty"`
}

// RgbCube represents RGB cube settings.
type RgbCube struct {
	Type      string    `xml:"type,attr,omitempty"`
	Dimension *XMLAttr  `xml:"dimension,omitempty"`
	Division  *XMLAttr  `xml:"division,omitempty"`
	RgbTable  string   `xml:"rgb-table,omitempty"`
}

// XMLGenerator generates screen.xml from grid configuration.
type XMLGenerator struct {
	grid            *GridWidget
	gpuAssignments  *GPUAssignment
	touchAreas      *TouchAreaHandler
	resolution      Resolution
	defaultRes      Resolution
	areaPerGPU      bool // If true, create area per window (per GPU); if false, area per GPU output (per screen)
}

// NewXMLGenerator creates a new XML generator.
func NewXMLGenerator(grid *GridWidget, gpuAssignments *GPUAssignment, touchAreas *TouchAreaHandler, resolution Resolution) *XMLGenerator {
	return &XMLGenerator{
		grid:           grid,
		gpuAssignments: gpuAssignments,
		touchAreas:     touchAreas,
		resolution:     resolution,
		defaultRes:     CommonResolutions[0], // 1920x1080
		areaPerGPU:     false, // Default: area per Screen (per GPU output)
	}
}

// SetAreaPerGPU sets whether to create area per GPU (window) or per Screen (GPU output).
// true = per GPU (one area per window matching all outputs on that GPU)
// false = per Screen (one area per GPU output, each output is one screen)
func (xg *XMLGenerator) SetAreaPerGPU(perGPU bool) {
	xg.areaPerGPU = perGPU
}

// Generate generates screen.xml content from the grid configuration.
func (xg *XMLGenerator) Generate() ([]byte, error) {
	// Calculate total layer size from all assigned cells
	layerSize := xg.calculateLayerSize()

	screenXML := ScreenXML{
		Type:      "",
		DPI:       &XMLAttr{Type: "", Value: "40.053"},
		DPMS:      &XMLAttr{Type: "", Value: "0 0 0"},
		HwColorCorrection: &XMLAttr{Type: "", Value: "0"},
		Iconify:   &XMLAttr{Type: "", Value: "0"},
		Vsync:     &XMLAttr{Type: "", Value: "0"}, // Default for Windows
		LayerSize: &XMLAttr{Type: "", Value: layerSize},
		Windows:   []WindowConfig{},
	}

	// Group cells by GPU output to create windows
	gpuGroups := xg.groupCellsByGPU()

	// Counters for unique window and area names
	windowCounter := 1
	areaCounter := 1

	if xg.areaPerGPU {
		// Area per GPU (window) mode: Create one area per window matching all outputs on that GPU
		// Group by GPU number to create windows
		windowsByGPU := xg.groupByGPU(gpuGroups)

		// Create windows for each GPU
		for gpuNum, outputs := range windowsByGPU {
			window := xg.createWindowForGPU(gpuNum, outputs, &windowCounter, &areaCounter)
			if window != nil {
				screenXML.Windows = append(screenXML.Windows, *window)
			}
		}
	} else {
		// Area per Screen (GPU output) mode: Create one area per GPU output (each output is one screen)
		// Create a window for each GPU (group by GPU number, not GPU output)
		windowsByGPU := xg.groupByGPU(gpuGroups)

		// Create windows for each GPU, with one area per output
		for gpuNum, outputs := range windowsByGPU {
			window := xg.createWindowForScreen(gpuNum, outputs, &windowCounter, &areaCounter)
			if window != nil {
				screenXML.Windows = append(screenXML.Windows, *window)
			}
		}
	}

	// Generate XML with proper encoding and formatting
	xmlData, err := xml.MarshalIndent(screenXML, "", "  ") // Use 2 spaces for indentation
	if err != nil {
		return nil, fmt.Errorf("failed to marshal XML: %w", err)
	}

	// Add XML header with UTF-8 encoding and DOCTYPE
	var buf []byte
	buf = append(buf, []byte(`<?xml version="1.0" encoding="UTF-8"?>`+"\n")...)
	buf = append(buf, []byte(`<!DOCTYPE mtdoc>`+"\n")...)
	buf = append(buf, xmlData...)

	return buf, nil
}

// calculateLayerSize calculates the total layer size from cells with layer > 0.
// Only considers cells that have a layer index > 0 (not layer 0 or empty).
func (xg *XMLGenerator) calculateLayerSize() string {
	// Find the maximum graphics coordinates from cells with layer > 0
	maxX, maxY := 0, 0

	for row := 0; row < GridRows; row++ {
		for col := 0; col < GridCols; col++ {
			cell := xg.grid.GetCell(row, col)
			if cell != nil && cell.GPUOutput != "" {
				// Check if cell has a layer index > 0
				hasLayer := false
				if cell.Index != "" {
					var layerIdx int
					if _, err := fmt.Sscanf(cell.Index, "%d", &layerIdx); err == nil && layerIdx > 0 {
						hasLayer = true
					}
				}

				// Only consider cells with layer > 0
				if hasLayer {
					// Calculate cell's graphics coordinates
					x := col * cell.Resolution.Width
					y := row * cell.Resolution.Height
					cellMaxX := x + cell.Resolution.Width
					cellMaxY := y + cell.Resolution.Height

					if cellMaxX > maxX {
						maxX = cellMaxX
					}
					if cellMaxY > maxY {
						maxY = cellMaxY
					}
				}
			}
		}
	}

	// If no cells with layer assigned, return "0 0" for auto-calculation
	if maxX == 0 && maxY == 0 {
		return "0 0"
	}

	return fmt.Sprintf("%d %d", maxX, maxY)
}

// groupByGPU groups GPU outputs by GPU number.
func (xg *XMLGenerator) groupByGPU(gpuGroups map[string][]CellCoord) map[int]map[string][]CellCoord {
	windowsByGPU := make(map[int]map[string][]CellCoord)

	for gpuOutput, cells := range gpuGroups {
		parts := strings.Split(gpuOutput, ":")
		if len(parts) != 2 {
			continue
		}
		var gpuNum int
		fmt.Sscanf(parts[0], "%d", &gpuNum)

		if windowsByGPU[gpuNum] == nil {
			windowsByGPU[gpuNum] = make(map[string][]CellCoord)
		}
		windowsByGPU[gpuNum][gpuOutput] = cells
	}

	return windowsByGPU
}

// createWindowForGPU creates a WindowConfig with a single area matching the window size.
// This is used when "area per GPU" mode is enabled (one area per window, matching all outputs on that GPU).
func (xg *XMLGenerator) createWindowForGPU(gpuNum int, outputs map[string][]CellCoord, windowCounter *int, areaCounter *int) *WindowConfig {
	if len(outputs) == 0 {
		return nil
	}

	// Calculate window size from all outputs
	minX, minY, maxX, maxY := xg.calculateWindowBounds(outputs)
	windowWidth := maxX - minX
	windowHeight := maxY - minY

	windowName := fmt.Sprintf("window%d", *windowCounter)
	*windowCounter++

	window := &WindowConfig{
		Type:            windowName,
		DirectRendering: &XMLAttr{Type: "", Value: "1"},
		Frameless:       &XMLAttr{Type: "", Value: "1"},
		FsaaSamples:     &XMLAttr{Type: "", Value: "4"},
		Fullscreen:      &XMLAttr{Type: "", Value: "0"},
		Location:        &XMLAttr{Type: "", Value: fmt.Sprintf("%d %d", minX, minY)},
		Resizable:       &XMLAttr{Type: "", Value: "0"},
		ScreenNumber:    &XMLAttr{Type: "", Value: "-1"},
		Size:            &XMLAttr{Type: "", Value: fmt.Sprintf("%d %d", windowWidth, windowHeight)},
		Areas:           []AreaConfig{},
	}

	// Create a single area that matches the window size (all outputs on this GPU)
	// Collect all cells from all outputs
	allCells := []CellCoord{}
	for _, cells := range outputs {
		allCells = append(allCells, cells...)
	}

	// Create one area for the entire window
	area := xg.createAreaForWindow(allCells, minX, minY, windowWidth, windowHeight, areaCounter)
	if area != nil {
		window.Areas = append(window.Areas, *area)
		*areaCounter++
	}

	return window
}

// calculateWindowBounds calculates the bounding box for a window from its outputs.
func (xg *XMLGenerator) calculateWindowBounds(outputs map[string][]CellCoord) (minX, minY, maxX, maxY int) {
	first := true
	for _, cells := range outputs {
		for _, cell := range cells {
			cellState := xg.grid.GetCell(cell.Row, cell.Col)
			if cellState == nil {
				continue
			}
			x := cell.Col * cellState.Resolution.Width
			y := cell.Row * cellState.Resolution.Height
			cellMaxX := x + cellState.Resolution.Width
			cellMaxY := y + cellState.Resolution.Height

			if first {
				minX, minY, maxX, maxY = x, y, cellMaxX, cellMaxY
				first = false
			} else {
				if x < minX {
					minX = x
				}
				if y < minY {
					minY = y
				}
				if cellMaxX > maxX {
					maxX = cellMaxX
				}
				if cellMaxY > maxY {
					maxY = cellMaxY
				}
			}
		}
	}
	return
}

// groupCellsByGPU groups grid cells by their assigned GPU output.
func (xg *XMLGenerator) groupCellsByGPU() map[string][]CellCoord {
	groups := make(map[string][]CellCoord)

	for row := 0; row < GridRows; row++ {
		for col := 0; col < GridCols; col++ {
			cell := xg.grid.GetCell(row, col)
			if cell != nil && cell.GPUOutput != "" {
				groups[cell.GPUOutput] = append(groups[cell.GPUOutput], CellCoord{Row: row, Col: col})
			}
		}
	}

	return groups
}

// CellCoord represents a cell coordinate.
type CellCoord struct {
	Row int
	Col int
}

// createAreaForGPU creates an AreaConfig for a GPU output group.
func (xg *XMLGenerator) createAreaForGPU(gpuOutput string, cells []CellCoord, areaCounter *int) *AreaConfig {
	if len(cells) == 0 {
		return nil
	}

	// Calculate bounding box for the cells
	minRow, maxRow := cells[0].Row, cells[0].Row
	minCol, maxCol := cells[0].Col, cells[0].Col

	// Get resolution from first cell
	firstCell := xg.grid.GetCell(cells[0].Row, cells[0].Col)
	if firstCell == nil {
		return nil
	}
	cellWidth := firstCell.Resolution.Width
	cellHeight := firstCell.Resolution.Height

	for _, cell := range cells {
		if cell.Row < minRow {
			minRow = cell.Row
		}
		if cell.Row > maxRow {
			maxRow = cell.Row
		}
		if cell.Col < minCol {
			minCol = cell.Col
		}
		if cell.Col > maxCol {
			maxCol = cell.Col
		}
	}

	// Calculate graphics coordinates (relative to layer)
	x := minCol * cellWidth
	y := minRow * cellHeight
	width := (maxCol - minCol + 1) * cellWidth
	height := (maxRow - minRow + 1) * cellHeight

	// Parse GPU output to determine type
	// Format: "gpu#:output#" (e.g., "1:2" means GPU 0, output 1)
	parts := strings.Split(gpuOutput, ":")
	if len(parts) != 2 {
		return nil
	}

	var gpuNum, outputNum int
	fmt.Sscanf(parts[0], "%d", &gpuNum)
	fmt.Sscanf(parts[1], "%d", &outputNum)
	// Note: area type is not used in the XML structure - it's determined by the GPU output

	// Create default color correction
	colorCorr := &ColorCorrection{
		Type:       "",
		Brightness: &XMLAttr{Type: "", Value: "0 0 0"},
		Contrast:  &XMLAttr{Type: "", Value: "1 1 1"},
		Gamma:     &XMLAttr{Type: "", Value: "1 1 1"},
		Red:       "0 0 1 1 ",
		Green:     "0 0 1 1 ",
		Blue:      "0 0 1 1 ",
	}

	// Create default keystone
	keystone := &Keystone{
		Type:      "",
		Rotations: &XMLAttr{Type: "", Value: "0"},
		V1:        &XMLAttr{Type: "", Value: "0 0"},
		V2:        &XMLAttr{Type: "", Value: "1 0"},
		V3:        &XMLAttr{Type: "", Value: "1 1"},
		V4:        &XMLAttr{Type: "", Value: "0 1"},
	}

	// Create default RGB cube
	rgbCube := &RgbCube{
		Type:      "",
		Dimension: &XMLAttr{Type: "", Value: "32"},
		Division:  &XMLAttr{Type: "", Value: "0"},
		RgbTable:  "",
	}

	areaName := fmt.Sprintf("area%d", *areaCounter)

	area := &AreaConfig{
		Type:             areaName,
		ColorCorrection:  colorCorr,
		GraphicsLocation: &XMLAttr{Type: "", Value: fmt.Sprintf("%d %d", x, y)},
		GraphicsSize:     &XMLAttr{Type: "", Value: fmt.Sprintf("%d %d", width, height)},
		Keystone:         keystone,
		Location:         &XMLAttr{Type: "", Value: "0 0"},
		Method:           &XMLAttr{Type: "", Value: "1"},
		RgbCube:          rgbCube,
		Seams:            &XMLAttr{Type: "", Value: "0 0 0 0"},
		Size:             &XMLAttr{Type: "", Value: fmt.Sprintf("%d %d", width, height)},
	}

	return area
}

// createWindowForScreen creates a WindowConfig with one area per GPU output (per screen).
// This is used when "area per Screen" mode is enabled (each GPU output = one screen).
func (xg *XMLGenerator) createWindowForScreen(gpuNum int, outputs map[string][]CellCoord, windowCounter *int, areaCounter *int) *WindowConfig {
	if len(outputs) == 0 {
		return nil
	}

	// Calculate window size from all outputs
	minX, minY, maxX, maxY := xg.calculateWindowBounds(outputs)
	windowWidth := maxX - minX
	windowHeight := maxY - minY

	windowName := fmt.Sprintf("window%d", *windowCounter)
	*windowCounter++

	window := &WindowConfig{
		Type:            windowName,
		DirectRendering: &XMLAttr{Type: "", Value: "1"},
		Frameless:       &XMLAttr{Type: "", Value: "1"},
		FsaaSamples:     &XMLAttr{Type: "", Value: "4"},
		Fullscreen:      &XMLAttr{Type: "", Value: "0"},
		Location:        &XMLAttr{Type: "", Value: fmt.Sprintf("%d %d", minX, minY)},
		Resizable:       &XMLAttr{Type: "", Value: "0"},
		ScreenNumber:    &XMLAttr{Type: "", Value: "-1"},
		Size:            &XMLAttr{Type: "", Value: fmt.Sprintf("%d %d", windowWidth, windowHeight)},
		Areas:           []AreaConfig{},
	}

	// Create one area for each GPU output (each output = one screen)
	for gpuOutput, cells := range outputs {
		area := xg.createAreaForGPU(gpuOutput, cells, areaCounter)
		if area != nil {
			window.Areas = append(window.Areas, *area)
			*areaCounter++
		}
	}

	return window
}

// createAreaForWindow creates a single AreaConfig that matches the window size.
func (xg *XMLGenerator) createAreaForWindow(cells []CellCoord, x, y, width, height int, areaCounter *int) *AreaConfig {
	if len(cells) == 0 {
		return nil
	}

	// Create default color correction
	colorCorr := &ColorCorrection{
		Type:       "",
		Brightness: &XMLAttr{Type: "", Value: "0 0 0"},
		Contrast:  &XMLAttr{Type: "", Value: "1 1 1"},
		Gamma:     &XMLAttr{Type: "", Value: "1 1 1"},
		Red:       "0 0 1 1 ",
		Green:     "0 0 1 1 ",
		Blue:      "0 0 1 1 ",
	}

	// Create default keystone
	keystone := &Keystone{
		Type:      "",
		Rotations: &XMLAttr{Type: "", Value: "0"},
		V1:        &XMLAttr{Type: "", Value: "0 0"},
		V2:        &XMLAttr{Type: "", Value: "1 0"},
		V3:        &XMLAttr{Type: "", Value: "1 1"},
		V4:        &XMLAttr{Type: "", Value: "0 1"},
	}

	// Create default RGB cube
	rgbCube := &RgbCube{
		Type:      "",
		Dimension: &XMLAttr{Type: "", Value: "32"},
		Division:  &XMLAttr{Type: "", Value: "0"},
		RgbTable:  "",
	}

	area := &AreaConfig{
		Type:             "area",
		ColorCorrection:  colorCorr,
		GraphicsLocation: &XMLAttr{Type: "", Value: fmt.Sprintf("%d %d", x, y)},
		GraphicsSize:     &XMLAttr{Type: "", Value: fmt.Sprintf("%d %d", width, height)},
		Keystone:         keystone,
		Location:         &XMLAttr{Type: "", Value: "0 0"},
		Method:           &XMLAttr{Type: "", Value: "1"},
		RgbCube:          rgbCube,
		Seams:            &XMLAttr{Type: "", Value: "0 0 0 0"},
		Size:             &XMLAttr{Type: "", Value: fmt.Sprintf("%d %d", width, height)},
	}

	return area
}

// Validate validates the generated XML structure.
func (xg *XMLGenerator) Validate(xmlData []byte) error {
	var screen ScreenXML
	if err := xml.Unmarshal(xmlData, &screen); err != nil {
		return fmt.Errorf("invalid XML structure: %w", err)
	}

	// Basic validation
	if len(screen.Windows) == 0 {
		return fmt.Errorf("no windows defined in screen.xml")
	}

	for i, window := range screen.Windows {
		if len(window.Areas) == 0 {
			return fmt.Errorf("window %d has no areas", i)
		}
		for j, area := range window.Areas {
			if area.Type == "" {
				return fmt.Errorf("window %d, area %d missing type attribute", i, j)
			}
			if area.GraphicsSize == nil || area.GraphicsSize.Value == "" {
				return fmt.Errorf("window %d, area %d missing graphics size", i, j)
			}
		}
	}

	return nil
}

