package backup

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	// MaxBackups is the maximum number of backups to keep
	MaxBackups = 5
	// BackupSuffix is the suffix for backup files
	BackupSuffix = ".bak"
)

// Manager handles file backups with rotation.
type Manager struct {
	backupDir string
}

// NewManager creates a new backup manager.
func NewManager(backupDir string) *Manager {
	return &Manager{
		backupDir: backupDir,
	}
}

// CreateBackup creates a backup of a file if it has changed.
// Returns the backup path if created, empty string if no backup needed.
func (m *Manager) CreateBackup(filePath string) (string, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// No file to backup
		return "", nil
	}

	// Read original file
	originalData, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Ensure backup directory exists
	if err := os.MkdirAll(m.backupDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Generate backup filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("%s_%s%s",
		filepath.Base(filePath),
		timestamp,
		BackupSuffix)
	backupPath := filepath.Join(m.backupDir, backupName)

	// Check if this exact content was already backed up
	if m.hasIdenticalBackup(originalData) {
		// No need to create duplicate backup
		return "", nil
	}

	// Create backup file
	if err := os.WriteFile(backupPath, originalData, 0644); err != nil {
		return "", fmt.Errorf("failed to write backup: %w", err)
	}

	// Rotate old backups
	if err := m.rotateBackups(filepath.Base(filePath)); err != nil {
		// Log error but don't fail the backup
		// The backup was created successfully
	}

	return backupPath, nil
}

// hasIdenticalBackup checks if identical content was already backed up.
func (m *Manager) hasIdenticalBackup(data []byte) bool {
	entries, err := os.ReadDir(m.backupDir)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		backupPath := filepath.Join(m.backupDir, entry.Name())
		backupData, err := os.ReadFile(backupPath)
		if err != nil {
			continue
		}

		if len(backupData) == len(data) {
			identical := true
			for i := range data {
				if backupData[i] != data[i] {
					identical = false
					break
				}
			}
			if identical {
				return true
			}
		}
	}

	return false
}

// rotateBackups removes old backups, keeping only the most recent MaxBackups.
func (m *Manager) rotateBackups(baseName string) error {
	entries, err := os.ReadDir(m.backupDir)
	if err != nil {
		return err
	}

	// Find all backups for this file
	var backups []os.DirEntry
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		// Check if this is a backup of the target file
		if len(name) > len(baseName) && name[:len(baseName)] == baseName {
			backups = append(backups, entry)
		}
	}

	// If we have more than MaxBackups, remove the oldest
	if len(backups) > MaxBackups {
		// Sort by modification time (oldest first)
		// Simple approach: remove the ones with oldest timestamps in filename
		// For simplicity, we'll remove the ones that sort first alphabetically
		// (which should be oldest if timestamp format is consistent)
		toRemove := len(backups) - MaxBackups
		for i := 0; i < toRemove; i++ {
			backupPath := filepath.Join(m.backupDir, backups[i].Name())
			if err := os.Remove(backupPath); err != nil {
				return fmt.Errorf("failed to remove old backup: %w", err)
			}
		}
	}

	return nil
}

// RestoreBackup restores a file from a backup.
func (m *Manager) RestoreBackup(backupPath, targetPath string) error {
	src, err := os.Open(backupPath)
	if err != nil {
		return fmt.Errorf("failed to open backup: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("failed to create target file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy backup: %w", err)
	}

	return nil
}
