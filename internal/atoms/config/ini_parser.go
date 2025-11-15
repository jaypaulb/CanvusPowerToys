package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

// INIParser handles INI file parsing and writing.
type INIParser struct{}

// NewINIParser creates a new INI parser.
func NewINIParser() *INIParser {
	return &INIParser{}
}

// Read reads an INI file and returns the configuration.
func (p *INIParser) Read(filePath string) (*ini.File, error) {
	cfg, err := ini.Load(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load INI file: %w", err)
	}
	return cfg, nil
}

// Write writes an INI configuration to a file.
func (p *INIParser) Write(cfg *ini.File, filePath string) error {
	if err := cfg.SaveTo(filePath); err != nil {
		return fmt.Errorf("failed to save INI file: %w", err)
	}
	return nil
}

// GetSection gets a section from the INI file.
func (p *INIParser) GetSection(cfg *ini.File, sectionName string) (*ini.Section, error) {
	section := cfg.Section(sectionName)
	if section == nil {
		return nil, fmt.Errorf("section '%s' not found", sectionName)
	}
	return section, nil
}

// GetValue gets a value from a section.
func (p *INIParser) GetValue(section *ini.Section, key string) string {
	return section.Key(key).String()
}

// SetValue sets a value in a section.
func (p *INIParser) SetValue(section *ini.Section, key, value string) {
	section.Key(key).SetValue(value)
}

// CreateFile creates a new INI file.
func (p *INIParser) CreateFile(filePath string) (*ini.File, error) {
	cfg := ini.Empty()

	// Create parent directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	return cfg, nil
}
