package webui

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jaypaulb/CanvusPowerToys/internal/atoms/paths"
)

// WebUILogger handles logging for WebUI operations to both console and file.
// This allows users to debug WebUI issues when running from GUI without console.
type WebUILogger struct {
	logFile *os.File
	logPath string
}

// NewWebUILogger creates a new WebUI logger instance.
// Returns logger and error if log file cannot be created.
func NewWebUILogger() (*WebUILogger, error) {
	logDir, err := getWebUILogDirectory()
	if err != nil {
		// Fallback: if we can't get log directory, still return a logger that prints to console
		return &WebUILogger{logPath: ""}, nil
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFileName := fmt.Sprintf("webui_%s.log", timestamp)
	logPath := filepath.Join(logDir, logFileName)

	// Create log file
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		// Log directory might not exist yet, return logger without file
		return &WebUILogger{logPath: ""}, nil
	}

	return &WebUILogger{
		logFile: file,
		logPath: logPath,
	}, nil
}

// Log writes a message to both console and log file.
func (wl *WebUILogger) Log(message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	formattedMessage := fmt.Sprintf("[%s] %s\n", timestamp, message)

	// Print to stdout
	fmt.Print(formattedMessage)

	// Write to file if available
	if wl.logFile != nil {
		_, _ = wl.logFile.WriteString(formattedMessage)
		// Don't close file - keep it open for multiple writes
	}
}

// Logf writes a formatted message to both console and log file.
func (wl *WebUILogger) Logf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	wl.Log(message)
}

// GetLogPath returns the path to the current log file.
func (wl *WebUILogger) GetLogPath() string {
	return wl.logPath
}

// Close closes the log file.
func (wl *WebUILogger) Close() error {
	if wl.logFile != nil {
		return wl.logFile.Close()
	}
	return nil
}

// getWebUILogDirectory returns the directory for WebUI logs.
// Creates the directory if it doesn't exist.
func getWebUILogDirectory() (string, error) {
	// Use same log directory as main app
	logDir, err := paths.GetCanvusLogsPath()
	if err != nil {
		return "", err
	}

	// Create webui subdirectory
	webUILogDir := filepath.Join(logDir, "webui")
	if err := os.MkdirAll(webUILogDir, 0755); err != nil {
		return "", err
	}

	return webUILogDir, nil
}

// CleanupOldLogs removes WebUI log files older than 7 days.
// This prevents log directory from growing too large.
func CleanupOldLogs() {
	logDir, err := getWebUILogDirectory()
	if err != nil {
		return
	}

	entries, err := os.ReadDir(logDir)
	if err != nil {
		return
	}

	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		// Remove files older than 7 days
		if info.ModTime().Before(sevenDaysAgo) {
			logPath := filepath.Join(logDir, entry.Name())
			_ = os.Remove(logPath)
		}
	}
}
