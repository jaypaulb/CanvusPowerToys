package screenxml

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

// ScreenXML represents the screen.xml structure.
type ScreenXML struct {
	XMLName xml.Name `xml:"screen"`
	MultiHead MultiHeadConfig `xml:"MultiHead"`
}

// MultiHeadConfig represents MultiHead configuration.
type MultiHeadConfig struct {
	Vsync    string `xml:"vsync,attr,omitempty"`
	LayerSize string `xml:"layer-size,attr,omitempty"`
	Areas    []Area `xml:"area"`
}

// Area represents a display area in screen.xml.
type Area struct {
	XMLName xml.Name `xml:"area"`
	Type    string   `xml:"type,attr"`
	Graphics Graphics `xml:"graphics"`
	Window   Window  `xml:"window,omitempty"`
}

// Graphics represents graphics coordinates for an area.
type Graphics struct {
	X      string `xml:"x,attr"`
	Y      string `xml:"y,attr"`
	Width  string `xml:"width,attr"`
	Height string `xml:"height,attr"`
}

// Window represents window attributes.
type Window struct {
	Location string `xml:"location,attr,omitempty"`
}

// XMLGenerator generates screen.xml from grid configuration.
type XMLGenerator struct {
	grid            *GridWidget
	gpuAssignments  *GPUAssignment
	touchAreas      *TouchAreaHandler
	resolution      Resolution
	defaultRes      Resolution
}

// NewXMLGenerator creates a new XML generator.
func NewXMLGenerator(grid *GridWidget, gpuAssignments *GPUAssignment, touchAreas *TouchAreaHandler, resolution Resolution) *XMLGenerator {
	return &XMLGenerator{
		grid:           grid,
		gpuAssignments: gpuAssignments,
		touchAreas:     touchAreas,
		resolution:     resolution,
		defaultRes:     CommonResolutions[0], // 1920x1080
	}
}

// Generate generates screen.xml content from the grid configuration.
func (xg *XMLGenerator) Generate() ([]byte, error) {
	screenXML := ScreenXML{
		MultiHead: MultiHeadConfig{
			Vsync:     "0", // Default for Windows
			LayerSize: "0 0", // Auto-calculate
			Areas:     []Area{},
		},
	}

	// Group cells by GPU output
	gpuGroups := xg.groupCellsByGPU()

	// Create areas for each GPU output
	for gpuOutput, cells := range gpuGroups {
		area := xg.createAreaForGPU(gpuOutput, cells)
		if area != nil {
			screenXML.MultiHead.Areas = append(screenXML.MultiHead.Areas, *area)
		}
	}

	// Generate XML
	var buf []byte
	// We'll use encoding/xml directly for proper formatting
	xmlData, err := xml.MarshalIndent(screenXML, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal XML: %w", err)
	}

	// Add XML header
	buf = append(buf, []byte(xml.Header)...)
	buf = append(buf, xmlData...)

	return buf, nil
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

// createAreaForGPU creates an Area XML element for a GPU output group.
func (xg *XMLGenerator) createAreaForGPU(gpuOutput string, cells []CellCoord) *Area {
	if len(cells) == 0 {
		return nil
	}

	// Calculate bounding box for the cells
	minRow, maxRow := cells[0].Row, cells[0].Row
	minCol, maxCol := cells[0].Col, cells[0].Col

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

	// Calculate graphics coordinates
	// Each cell represents one display output at default resolution
	cellWidth := xg.defaultRes.Width
	cellHeight := xg.defaultRes.Height

	x := minCol * cellWidth
	y := minRow * cellHeight
	width := (maxCol - minCol + 1) * cellWidth
	height := (maxRow - minRow + 1) * cellHeight

	// Parse GPU output to determine type
	// Format: "gpu#.output#"
	parts := strings.Split(gpuOutput, ".")
	if len(parts) != 2 {
		return nil
	}

	areaType := fmt.Sprintf("gpu%s.output%s", parts[0], parts[1])

	area := &Area{
		Type: areaType,
		Graphics: Graphics{
			X:      strconv.Itoa(x),
			Y:      strconv.Itoa(y),
			Width:  strconv.Itoa(width),
			Height: strconv.Itoa(height),
		},
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
	if len(screen.MultiHead.Areas) == 0 {
		return fmt.Errorf("no areas defined in screen.xml")
	}

	for i, area := range screen.MultiHead.Areas {
		if area.Type == "" {
			return fmt.Errorf("area %d missing type attribute", i)
		}
		if area.Graphics.Width == "" || area.Graphics.Height == "" {
			return fmt.Errorf("area %d missing graphics dimensions", i)
		}
	}

	return nil
}

