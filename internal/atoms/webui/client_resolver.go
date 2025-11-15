package webui

import (
	"fmt"

	"github.com/jaypaulb/CanvusPowerToys/internal/atoms/config"
	"github.com/jaypaulb/CanvusPowerToys/internal/organisms/services"
)

// ClientResolver handles resolution of client_id from installation_name.
type ClientResolver struct {
	fileService *services.FileService
	iniParser   *config.INIParser
}

// NewClientResolver creates a new client resolver.
func NewClientResolver(fileService *services.FileService) *ClientResolver {
	return &ClientResolver{
		fileService: fileService,
		iniParser:   config.NewINIParser(),
	}
}

// GetInstallationName reads installation_name from mt-canvus.ini or falls back to device name.
func (r *ClientResolver) GetInstallationName() (string, error) {
	iniPath := r.fileService.DetectMtCanvusIni()
	if iniPath == "" {
		// No INI file found, use device name
		return GetDeviceName()
	}

	iniFile, err := r.iniParser.Read(iniPath)
	if err != nil {
		// If we can't read INI, fall back to device name
		return GetDeviceName()
	}

	// Try to get installation_name from [canvas] section
	canvasSection := iniFile.Section("canvas")
	if canvasSection == nil {
		// No [canvas] section, use device name
		return GetDeviceName()
	}

	installationName := canvasSection.Key("installation_name").String()
	if installationName == "" {
		// Not set in INI, use device name
		return GetDeviceName()
	}

	return installationName, nil
}

// ResolveClientID queries the Canvus API to find client_id matching installation_name.
// This will be implemented when we have the API client ready.
// For now, this is a placeholder that returns an error.
func (r *ClientResolver) ResolveClientID(apiBaseURL, authToken, installationName string) (string, error) {
	// TODO: Implement API call to GET /api/v1/clients
	// TODO: Match installation_name in response
	// TODO: Return matching client_id
	// TODO: Handle errors (not found, API errors, etc.)
	return "", fmt.Errorf("not implemented: requires Canvus API client")
}
