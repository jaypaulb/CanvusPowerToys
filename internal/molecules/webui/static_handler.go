package webui

import (
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
)

// StaticHandler handles serving static frontend files.
type StaticHandler struct {
	fileSystem fs.FS
}

// NewStaticHandler creates a new static file handler.
func NewStaticHandler() *StaticHandler {
	// Get the embedded public directory
	// The embed is at webui/public, so we access it directly
	// embeddedAssets already points to the public directory
	return &StaticHandler{
		fileSystem: embeddedAssets,
	}
}

// ServeFiles registers static file routes with the given mux.
func (sh *StaticHandler) ServeFiles(mux *http.ServeMux) {
	// Handle root and static files
	// Note: This should be registered AFTER API routes so API routes take precedence
	// In Go's http.ServeMux, more specific routes (like /api/...) take precedence over /
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Don't handle API routes - let API routes handler take care of them
		// More specific routes registered first will match before this catch-all
		if strings.HasPrefix(r.URL.Path, "/api/") {
			// This shouldn't happen if API routes are registered correctly,
			// but just in case, return 404
			http.NotFound(w, r)
			return
		}

		if r.URL.Path == "/" || r.URL.Path == "" {
			// Serve main.html for root
			sh.serveFile(w, r, "pages/html/main.html")
			return
		}

		// Check if it's a page request (ends with .html)
		if strings.HasSuffix(r.URL.Path, ".html") {
			// Map common routes to HTML files
			htmlPath := sh.mapRouteToHTML(r.URL.Path)
			sh.serveFile(w, r, htmlPath)
			return
		}

		// For all other paths, remove leading slash and try to find file
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			sh.serveFile(w, r, "pages/html/main.html")
			return
		}

		// Try direct path first (since embed root is "public", files are at root level)
		// But wait - embed is "public", so root contains "public" directory
		// So we need to check if path exists directly, or with public/ prefix

		// Normalize path separators (embedded FS always uses forward slashes)
		path = strings.ReplaceAll(path, "\\", "/")

		// First try: path as-is (for requests like /atoms/css/badge.css)
		// Since embed root is "public", we need "public/atoms/css/badge.css"
		publicPath := "public/" + path
		if _, err := fs.Stat(sh.fileSystem, publicPath); err == nil {
			sh.serveFile(w, r, publicPath)
			return
		}

		// Second try: path without public/ (in case embed structure is different)
		if _, err := fs.Stat(sh.fileSystem, path); err == nil {
			sh.serveFile(w, r, path)
			return
		}

		// Try with common prefixes
		for _, prefix := range []string{"public/pages/html/", "public/pages/js/", "public/pages/css/", "public/molecules/js/", "public/molecules/css/", "public/atoms/css/", "public/atoms/js/", "public/css/", "public/templates/css/", "public/templates/html/"} {
			fullPath := prefix + path
			if _, err := fs.Stat(sh.fileSystem, fullPath); err == nil {
				sh.serveFile(w, r, fullPath)
				return
			}
		}

		// Try without public/ prefix
		for _, prefix := range []string{"pages/html/", "pages/js/", "pages/css/", "molecules/js/", "molecules/css/", "atoms/css/", "atoms/js/", "css/", "templates/css/", "templates/html/"} {
			fullPath := prefix + path
			if _, err := fs.Stat(sh.fileSystem, fullPath); err == nil {
				sh.serveFile(w, r, fullPath)
				return
			}
		}

		// Not found
		fmt.Printf("[StaticHandler] File not found for path: %s\n", r.URL.Path)
		http.NotFound(w, r)
	})
}

// mapRouteToHTML maps URL routes to HTML file paths.
func (sh *StaticHandler) mapRouteToHTML(path string) string {
	// Remove leading slash and .html extension
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, ".html")

	// Map common routes (embed root is "public", so files are at "public/pages/html/...")
	switch path {
	case "", "index", "main":
		return "pages/html/main.html"
	case "pages":
		return "pages/html/pages.html"
	case "macros":
		return "pages/html/macros.html"
	case "remote-upload", "upload":
		return "pages/html/remote-upload.html"
	case "rcu":
		return "pages/html/rcu.html"
	default:
		// Try pages/html/{path}.html
		return "pages/html/" + path + ".html"
	}
}

// serveFile serves a specific file from the embedded filesystem.
func (sh *StaticHandler) serveFile(w http.ResponseWriter, r *http.Request, filePath string) {
	// Since embed is "//go:embed public", the filesystem root IS "public"
	// So we need to add "public/" prefix if not already present
	if !strings.HasPrefix(filePath, "public/") {
		filePath = "public/" + filePath
	}

	// IMPORTANT: Embedded filesystems always use forward slashes, even on Windows
	// Convert any backslashes to forward slashes (filepath.Clean converts to OS-specific separators)
	filePath = strings.ReplaceAll(filePath, "\\", "/")
	filePath = strings.TrimPrefix(filePath, "./") // Remove any leading ./

	// Debug: Log what we're trying to access
	fmt.Printf("[StaticHandler] Attempting to serve file: %s (request: %s)\n", filePath, r.URL.Path)

	// Check if file exists first
	if _, err := fs.Stat(sh.fileSystem, filePath); err != nil {
		fmt.Printf("[StaticHandler] File not found: %s, error: %v\n", filePath, err)
		// Try to list what's in the filesystem root for debugging
		if entries, listErr := fs.ReadDir(sh.fileSystem, "."); listErr == nil {
			fmt.Printf("[StaticHandler] Filesystem root contents:\n")
			for _, entry := range entries {
				fmt.Printf("  - %s (dir: %v)\n", entry.Name(), entry.IsDir())
			}
		}
		http.NotFound(w, r)
		return
	}

	// Read file from embedded filesystem
	data, err := fs.ReadFile(sh.fileSystem, filePath)
	if err != nil {
		fmt.Printf("[StaticHandler] Error reading file %s: %v\n", filePath, err)
		http.NotFound(w, r)
		return
	}

	fmt.Printf("[StaticHandler] Successfully serving file: %s (size: %d bytes)\n", filePath, len(data))

	// Set content type based on file extension
	ext := filepath.Ext(filePath)
	switch ext {
	case ".html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case ".css":
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	case ".js":
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Write file content
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

