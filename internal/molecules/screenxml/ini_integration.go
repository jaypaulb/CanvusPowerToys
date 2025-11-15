package screenxml

import (
	"fmt"
	"strings"

	"gopkg.in/ini.v1"

	"github.com/jaypaulb/CanvusPowerToys/internal/atoms/config"
)

// INIIntegration handles integration with mt-canvus.ini.
type INIIntegration struct {
	iniParser *config.INIParser
}

// NewINIIntegration creates a new INI integration handler.
func NewINIIntegration() *INIIntegration {
	return &INIIntegration{
		iniParser: config.NewINIParser(),
	}
}

// DetectVideoOutputs detects video outputs (areas not in layout).
func (ii *INIIntegration) DetectVideoOutputs(grid *GridWidget) []string {
	var videoOutputs []string

	for row := 0; row < GridRows; row++ {
		for col := 0; col < GridCols; col++ {
			cell := grid.GetCell(row, col)
			if cell != nil && cell.GPUOutput != "" && !cell.IsLayoutArea {
				// This is a video output (has GPU but not in layout area)
				videoOutputs = append(videoOutputs, cell.GPUOutput)
			}
		}
	}

	return videoOutputs
}

// GenerateVideoOutputConfig generates video-output configuration for mt-canvus.ini.
func (ii *INIIntegration) GenerateVideoOutputConfig(videoOutputs []string) string {
	if len(videoOutputs) == 0 {
		return ""
	}

	// Format: video-output=gpu0.output1,gpu0.output2,...
	return strings.Join(videoOutputs, ",")
}

// UpdateMtCanvusIni updates mt-canvus.ini with video-output configuration.
func (ii *INIIntegration) UpdateMtCanvusIni(filePath string, videoOutputs []string) error {
	// Read existing INI file
	cfg, err := ii.iniParser.Read(filePath)
	if err != nil {
		return fmt.Errorf("failed to read mt-canvus.ini: %w", err)
	}

	// Get or create [MultiHead] section
	section, err := cfg.GetSection("MultiHead")
	if err != nil {
		section, err = cfg.NewSection("MultiHead")
		if err != nil {
			return fmt.Errorf("failed to create MultiHead section: %w", err)
		}
	}

	// Set video-output
	videoOutputValue := ii.GenerateVideoOutputConfig(videoOutputs)
	if videoOutputValue != "" {
		ii.iniParser.SetValue(section, "video-output", videoOutputValue)
	}

	// Write back to file
	if err := ii.iniParser.Write(cfg, filePath); err != nil {
		return fmt.Errorf("failed to write mt-canvus.ini: %w", err)
	}

	return nil
}

// ShouldUpdateIni checks if video-output should be updated in mt-canvus.ini.
func (ii *INIIntegration) ShouldUpdateIni(grid *GridWidget) bool {
	videoOutputs := ii.DetectVideoOutputs(grid)
	return len(videoOutputs) > 0
}
