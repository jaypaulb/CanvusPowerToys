package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
)

func main() {
	srcDir := "webui/src"
	publicDir := "webui/public"

	// Create public directory structure
	if err := os.RemoveAll(publicDir); err != nil {
		fmt.Printf("Error removing public directory: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(publicDir, 0755); err != nil {
		fmt.Printf("Error creating public directory: %v\n", err)
		os.Exit(1)
	}

	// Initialize minifier
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/javascript", js.Minify)

	// Process all files
	err := filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Get relative path from srcDir
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		// Determine output path
		destPath := filepath.Join(publicDir, relPath)

		// Create destination directory
		destDir := filepath.Dir(destPath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", destDir, err)
		}

		// Read source file
		srcData, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		// Determine MIME type from extension
		ext := strings.ToLower(filepath.Ext(path))
		var mimeType string
		switch ext {
		case ".css":
			mimeType = "text/css"
		case ".js":
			mimeType = "text/javascript"
		case ".html":
			mimeType = "text/html"
		default:
			// Copy file as-is for other types
			if err := os.WriteFile(destPath, srcData, 0644); err != nil {
				return fmt.Errorf("failed to write %s: %w", destPath, err)
			}
			fmt.Printf("Copied: %s -> %s\n", relPath, relPath)
			return nil
		}

		// Minify file
		minified, err := m.Bytes(mimeType, srcData)
		if err != nil {
			return fmt.Errorf("failed to minify %s: %w", path, err)
		}

		// Write minified file
		if err := os.WriteFile(destPath, minified, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", destPath, err)
		}

		// Calculate size reduction
		originalSize := len(srcData)
		minifiedSize := len(minified)
		reduction := float64(originalSize-minifiedSize) / float64(originalSize) * 100

		fmt.Printf("Minified: %s (%.1f%% reduction: %d -> %d bytes)\n", relPath, reduction, originalSize, minifiedSize)
		return nil
	})

	if err != nil {
		fmt.Printf("Error processing assets: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nAsset processing complete!")
}

