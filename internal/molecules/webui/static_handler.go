package webui

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
)

//go:embed ../../webui/public/*
var publicFiles embed.FS

// StaticHandler handles serving embedded static files.
type StaticHandler struct{}

// NewStaticHandler creates a new static file handler.
func NewStaticHandler() *StaticHandler {
	return &StaticHandler{}
}

// ServeFiles serves embedded static files.
func (h *StaticHandler) ServeFiles(mux *http.ServeMux) {
	// Get subdirectory from embed
	fsys, err := fs.Sub(publicFiles, "webui/public")
	if err != nil {
		// If embed fails, return (files will be served from disk in dev)
		return
	}

	fileServer := http.FileServer(http.FS(fsys))

	// Serve static files
	mux.Handle("/css/", http.StripPrefix("/css/", fileServer))
	mux.Handle("/js/", http.StripPrefix("/js/", fileServer))
	mux.Handle("/atoms/", http.StripPrefix("/atoms/", fileServer))
	mux.Handle("/molecules/", http.StripPrefix("/molecules/", fileServer))
	mux.Handle("/templates/", http.StripPrefix("/templates/", fileServer))
	mux.Handle("/pages/", http.StripPrefix("/pages/", fileServer))

	// Serve HTML pages
	mux.HandleFunc("/", h.handleIndex)
	mux.HandleFunc("/main.html", h.handlePage("main.html"))
	mux.HandleFunc("/pages.html", h.handlePage("pages.html"))
	mux.HandleFunc("/macros.html", h.handlePage("macros.html"))
	mux.HandleFunc("/remote-upload.html", h.handlePage("remote-upload.html"))
	mux.HandleFunc("/rcu.html", h.handlePage("rcu.html"))
}

// handleIndex serves the main page.
func (h *StaticHandler) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	h.handlePage("main.html")(w, r)
}

// handlePage returns a handler for a specific HTML page.
func (h *StaticHandler) handlePage(pageName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fsys, err := fs.Sub(publicFiles, "webui/public")
		if err != nil {
			http.Error(w, "Files not available", http.StatusInternalServerError)
			return
		}

		file, err := fsys.Open("pages/html/" + pageName)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer file.Close()

		// Read file content
		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(content)
	}
}

