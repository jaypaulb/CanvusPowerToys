package services

import (
	"os"
	"path/filepath"

	"github.com/jaypaulb/CanvusPowerToys/internal/atoms/paths"
)

// FileService handles file detection and path operations.
type FileService struct {
	userConfigPath   string
	systemConfigPath string
}

// NewFileService creates a new file service instance.
func NewFileService() (*FileService, error) {
	userPath, err := paths.GetCanvusUserConfigPath()
	if err != nil {
		return nil, err
	}

	systemPath, err := paths.GetCanvusSystemConfigPath()
	if err != nil {
		return nil, err
	}

	return &FileService{
		userConfigPath:   userPath,
		systemConfigPath: systemPath,
	}, nil
}

// DetectMtCanvusIni attempts to find mt-canvus.ini in standard locations.
// Returns the path if found, empty string if not found.
func (fs *FileService) DetectMtCanvusIni() string {
	// Check user config location first
	userIni := filepath.Join(fs.userConfigPath, "mt-canvus.ini")
	if paths.FileExists(userIni) {
		return userIni
	}

	// Check system config location
	systemIni := filepath.Join(fs.systemConfigPath, "mt-canvus.ini")
	if paths.FileExists(systemIni) {
		return systemIni
	}

	return ""
}

// DetectScreenXml attempts to find screen.xml in standard locations.
// Returns the path if found, empty string if not found.
func (fs *FileService) DetectScreenXml() string {
	// Check user config location first
	userXml := filepath.Join(fs.userConfigPath, "screen.xml")
	if paths.FileExists(userXml) {
		return userXml
	}

	// Check system config location
	systemXml := filepath.Join(fs.systemConfigPath, "screen.xml")
	if paths.FileExists(systemXml) {
		return systemXml
	}

	return ""
}

// DetectMenuYml attempts to find menu.yml in standard locations.
// Returns the path if found, empty string if not found.
func (fs *FileService) DetectMenuYml() string {
	// Check user config location first
	userYml := filepath.Join(fs.userConfigPath, "menu.yml")
	if paths.FileExists(userYml) {
		return userYml
	}

	// Check system config location
	systemYml := filepath.Join(fs.systemConfigPath, "menu.yml")
	if paths.FileExists(systemYml) {
		return systemYml
	}

	return ""
}

// GetUserConfigPath returns the user config directory path.
func (fs *FileService) GetUserConfigPath() string {
	return fs.userConfigPath
}

// GetSystemConfigPath returns the system config directory path.
func (fs *FileService) GetSystemConfigPath() string {
	return fs.systemConfigPath
}

// EnsureDirectory creates a directory if it doesn't exist.
func (fs *FileService) EnsureDirectory(path string) error {
	return os.MkdirAll(path, 0755)
}
