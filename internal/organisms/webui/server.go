package webui

import (
	"context"
	"fmt"
	"net/http"
	"time"

	webuiatoms "github.com/jaypaulb/CanvusPowerToys/internal/atoms/webui"
	webuimolecules "github.com/jaypaulb/CanvusPowerToys/internal/molecules/webui"
	"github.com/jaypaulb/CanvusPowerToys/internal/organisms/services"
)

// Server represents the main WebUI HTTP server.
type Server struct {
	httpServer    *http.Server
	canvasService *webuimolecules.CanvasService
	apiRoutes     *webuimolecules.APIRoutes
	port          string
	fileService   *services.FileService
	apiBaseURL    string
	authToken     string
	uploadDir     string
}

// NewServer creates a new WebUI server instance.
func NewServer(fileService *services.FileService, apiBaseURL, authToken, port, uploadDir string) (*Server, error) {
	// Create API client
	apiClient := webuiatoms.NewAPIClient(apiBaseURL, authToken)

	// Create canvas service
	canvasService, err := webuimolecules.NewCanvasService(fileService, apiBaseURL, authToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create canvas service: %w", err)
	}

	// Create API routes
	apiRoutes := webuimolecules.NewAPIRoutes(canvasService, apiClient, uploadDir)

	return &Server{
		canvasService: canvasService,
		apiRoutes:     apiRoutes,
		port:          port,
		fileService:   fileService,
		apiBaseURL:    apiBaseURL,
		authToken:     authToken,
		uploadDir:     uploadDir,
	}, nil
}

// Start starts the HTTP server and canvas tracking.
func (s *Server) Start() error {
	// Start canvas service (resolves client_id and starts subscription)
	if err := s.canvasService.Start(); err != nil {
		return fmt.Errorf("failed to start canvas service: %w", err)
	}

	// Create HTTP mux
	mux := http.NewServeMux()

	// Register API routes
	s.apiRoutes.RegisterRoutes(mux)

	// Register static file handlers
	staticHandler := webuimolecules.NewStaticHandler()
	staticHandler.ServeFiles(mux)

	// Create HTTP server
	s.httpServer = &http.Server{
		Addr:         ":" + s.port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("WebUI server error: %v\n", err)
		}
	}()

	return nil
}

// Stop stops the HTTP server and canvas tracking.
func (s *Server) Stop() error {
	// Stop canvas service
	s.canvasService.Stop()

	// Shutdown HTTP server
	if s.httpServer != nil {
		// Reduced timeout to 5 seconds since SSE handler now checks context every 1 second
		// This should be sufficient for graceful shutdown of all connections
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(ctx); err != nil {
			// Log error but continue - server will still stop
			if err == context.DeadlineExceeded {
				fmt.Printf("Server shutdown: Some connections did not close within timeout, forcing close\n")
			} else {
				fmt.Printf("Server shutdown error: %v\n", err)
			}
			// Force close if graceful shutdown failed
			s.httpServer.Close()
			return fmt.Errorf("failed to shutdown server gracefully: %w", err)
		}
		fmt.Printf("Server shutdown: All connections closed gracefully\n")
	}

	return nil
}

// GetPort returns the server port.
func (s *Server) GetPort() string {
	return s.port
}

// IsRunning returns whether the server is running.
func (s *Server) IsRunning() bool {
	return s.httpServer != nil
}

// GetCanvasService returns the canvas service instance.
func (s *Server) GetCanvasService() *webuimolecules.CanvasService {
	return s.canvasService
}
