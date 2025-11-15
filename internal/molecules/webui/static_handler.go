package webui

import (
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
	// Serve static files from embedded filesystem
	fileServer := http.FileServer(http.FS(sh.fileSystem))

	// Handle root - serve main.html
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

		// For all other paths, serve from embedded filesystem
		// Remove leading slash and serve
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "pages/html/main.html"
		}

		// Check if file exists
		if _, err := fs.Stat(sh.fileSystem, path); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		// Try with common prefixes
		for _, prefix := range []string{"pages/html/", "pages/js/", "pages/css/", "molecules/js/", "molecules/css/", "atoms/css/", "css/", "templates/css/", "templates/html/"} {
			fullPath := prefix + path
			if _, err := fs.Stat(sh.fileSystem, fullPath); err == nil {
				r.URL.Path = "/" + fullPath
				fileServer.ServeHTTP(w, r)
				return
			}
		}

		// Default: serve from filesystem (will 404 if not found)
		fileServer.ServeHTTP(w, r)
	})
}

// mapRouteToHTML maps URL routes to HTML file paths.
func (sh *StaticHandler) mapRouteToHTML(path string) string {
	// Remove leading slash and .html extension
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, ".html")

	// Map common routes
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
	// Clean the path
	filePath = filepath.Clean(filePath)

	// Read file from embedded filesystem
	data, err := fs.ReadFile(sh.fileSystem, filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

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

