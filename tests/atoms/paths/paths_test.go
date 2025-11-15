package paths_test

import (
	"os"
	"testing"

	"github.com/jaypaulb/CanvusPowerToys/internal/atoms/paths"
)

func TestGetAppDataPath(t *testing.T) {
	path, err := paths.GetAppDataPath()
	if err != nil {
		t.Fatalf("GetAppDataPath() error = %v", err)
	}
	if path == "" {
		t.Error("GetAppDataPath() returned empty path")
	}
}

func TestGetProgramDataPath(t *testing.T) {
	path, err := paths.GetProgramDataPath()
	if err != nil {
		t.Fatalf("GetProgramDataPath() error = %v", err)
	}
	if path == "" {
		t.Error("GetProgramDataPath() returned empty path")
	}
}

func TestGetLocalAppDataPath(t *testing.T) {
	path, err := paths.GetLocalAppDataPath()
	if err != nil {
		t.Fatalf("GetLocalAppDataPath() error = %v", err)
	}
	if path == "" {
		t.Error("GetLocalAppDataPath() returned empty path")
	}
}

func TestGetCanvusUserConfigPath(t *testing.T) {
	path, err := paths.GetCanvusUserConfigPath()
	if err != nil {
		t.Fatalf("GetCanvusUserConfigPath() error = %v", err)
	}
	if path == "" {
		t.Error("GetCanvusUserConfigPath() returned empty path")
	}
}

func TestGetCanvusSystemConfigPath(t *testing.T) {
	path, err := paths.GetCanvusSystemConfigPath()
	if err != nil {
		t.Fatalf("GetCanvusSystemConfigPath() error = %v", err)
	}
	if path == "" {
		t.Error("GetCanvusSystemConfigPath() returned empty path")
	}
}

func TestGetCanvusLogsPath(t *testing.T) {
	path, err := paths.GetCanvusLogsPath()
	if err != nil {
		t.Fatalf("GetCanvusLogsPath() error = %v", err)
	}
	if path == "" {
		t.Error("GetCanvusLogsPath() returned empty path")
	}
}

func TestJoinPath(t *testing.T) {
	result := paths.JoinPath("a", "b", "c")
	expected := "a/b/c"
	if result != expected && result != "a\\b\\c" {
		t.Errorf("JoinPath() = %v, want %v or a\\b\\c", result, expected)
	}
}

func TestFileExists(t *testing.T) {
	// Test with non-existent file
	if paths.FileExists("/nonexistent/file/path") {
		t.Error("FileExists() returned true for non-existent file")
	}

	// Test with existing file (current test file)
	if !paths.FileExists("paths_test.go") {
		t.Error("FileExists() returned false for existing file")
	}
}

func TestIsDir(t *testing.T) {
	// Test with non-existent path
	if paths.IsDir("/nonexistent/dir") {
		t.Error("IsDir() returned true for non-existent path")
	}

	// Test with current directory
	wd, _ := os.Getwd()
	if !paths.IsDir(wd) {
		t.Error("IsDir() returned false for existing directory")
	}
}
