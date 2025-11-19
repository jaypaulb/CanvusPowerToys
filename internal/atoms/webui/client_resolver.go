package webui

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
func (r *ClientResolver) ResolveClientID(apiBaseURL, authToken, installationName string) (string, error) {
	// Create HTTP client with TLS verification disabled for self-signed certs
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}

	// Query Canvus API for clients
	// Base URL should already include /api/v1 if provided, but we'll be explicit
	url := fmt.Sprintf("%s/api/v1/clients", apiBaseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Private-Token", authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("API error: %d", resp.StatusCode)
	}

	// Parse response
	var clients []struct {
		ID               string `json:"id"`
		InstallationName string `json:"installation_name"`
		Name             string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&clients); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Find matching client by installation_name
	for _, client := range clients {
		if client.InstallationName == installationName {
			return client.ID, nil
		}
		// Also try matching by name as fallback
		if client.Name == installationName {
			return client.ID, nil
		}
	}

	return "", fmt.Errorf("client not found: no client with installation_name '%s'", installationName)
}
