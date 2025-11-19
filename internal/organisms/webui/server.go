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
	logger        *webuiatoms.WebUILogger
}

// NewServer creates a new WebUI server instance.
func NewServer(fileService *services.FileService, apiBaseURL, authToken, port, uploadDir string) (*Server, error) {
	// Create WebUI logger
	logger, _ := webuiatoms.NewWebUILogger()
	if logger != nil {
		logger.Logf("[Server] Creating WebUI server with apiBaseURL: '%s', port: %s", apiBaseURL, port)
		// Cleanup old logs in background
		go webuiatoms.CleanupOldLogs()
	}

	// Create API client
	apiClient := webuiatoms.NewAPIClient(apiBaseURL, authToken)

	// Create canvas service
	canvasService, err := webuimolecules.NewCanvasService(fileService, apiBaseURL, authToken)
	if err != nil {
		if logger != nil {
			logger.Logf("[Server] ERROR: Failed to create canvas service: %v", err)
		}
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
		logger:        logger,
	}, nil
}

// Start starts the HTTP server and canvas tracking.
func (s *Server) Start() error {
	if s.logger != nil {
		s.logger.Log("[Server.Start] Starting WebUI server...")
	}

	// Start canvas service (resolves client_id and starts subscription)
	if err := s.canvasService.Start(); err != nil {
		if s.logger != nil {
			s.logger.Logf("[Server.Start] ERROR: Canvas service failed to start: %v", err)
		}
		return fmt.Errorf("failed to start canvas service: %w", err)
	}

	if s.logger != nil {
		s.logger.Log("[Server.Start] Canvas service started successfully")
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
		if s.logger != nil {
			s.logger.Logf("[Server.Start] HTTP server listening on %s", s.httpServer.Addr)
		}
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			if s.logger != nil {
				s.logger.Logf("[Server.Start] ERROR: HTTP server error: %v", err)
			}
		}
	}()

	return nil
}

// Stop stops the HTTP server and canvas tracking.
func (s *Server) Stop() error {
	if s.logger != nil {
		s.logger.Log("[Server.Stop] Stopping WebUI server...")
	}

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
				if s.logger != nil {
					s.logger.Log("[Server.Stop] Server shutdown: Some connections did not close within timeout, forcing close")
				}
			} else {
				if s.logger != nil {
					s.logger.Logf("[Server.Stop] ERROR: Server shutdown error: %v", err)
				}
			}
			// Force close if graceful shutdown failed
			s.httpServer.Close()
			return fmt.Errorf("failed to shutdown server gracefully: %w", err)
		}
		if s.logger != nil {
			s.logger.Log("[Server.Stop] Server shutdown: All connections closed gracefully")
		}
	}

	// Close logger file
	if s.logger != nil {
		_ = s.logger.Close()
	}

	return nil
}

// GetLogPath returns the path to the WebUI log file.
func (s *Server) GetLogPath() string {
	if s.logger != nil {
		return s.logger.GetLogPath()
	}
	return ""
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
