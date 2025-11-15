package webui

import (
	"context"
	"fmt"

	webuiatoms "github.com/jaypaulb/CanvusPowerToys/internal/atoms/webui"
	"github.com/jaypaulb/CanvusPowerToys/internal/organisms/services"
)

// CanvasService manages canvas tracking by composing client resolver and workspace subscriber.
type CanvasService struct {
	clientResolver      *webuiatoms.ClientResolver
	workspaceSubscriber *webuiatoms.WorkspaceSubscriber
	canvasTracker       *webuiatoms.CanvasTracker
	ctx                 context.Context
	cancel              context.CancelFunc
	apiBaseURL          string
	authToken           string
	clientID            string
	installationName    string
}

// NewCanvasService creates a new canvas service.
func NewCanvasService(fileService *services.FileService, apiBaseURL, authToken string) (*CanvasService, error) {
	clientResolver := webuiatoms.NewClientResolver(fileService)
	canvasTracker := webuiatoms.NewCanvasTracker()

	// Get installation name
	installationName, err := clientResolver.GetInstallationName()
	if err != nil {
		return nil, fmt.Errorf("failed to get installation name: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &CanvasService{
		clientResolver:      clientResolver,
		canvasTracker:      canvasTracker,
		ctx:                 ctx,
		cancel:              cancel,
		apiBaseURL:          apiBaseURL,
		authToken:           authToken,
		installationName:    installationName,
		workspaceSubscriber: nil, // Will be created after client_id resolution
	}, nil
}

// Start initializes client_id resolution and starts workspace subscription.
func (cs *CanvasService) Start() error {
	// Resolve client_id from installation_name
	clientID, err := cs.clientResolver.ResolveClientID(cs.apiBaseURL, cs.authToken, cs.installationName)
	if err != nil {
		return fmt.Errorf("failed to resolve client_id: %w", err)
	}

	cs.clientID = clientID

	// Create workspace subscriber
	cs.workspaceSubscriber = webuiatoms.NewWorkspaceSubscriber(
		cs.clientID,
		cs.apiBaseURL,
		cs.authToken,
	)

	// Start subscription
	eventChan, errChan := cs.workspaceSubscriber.Subscribe(cs.ctx)

	// Process events in background
	go cs.processEvents(eventChan, errChan)

	return nil
}

// Stop stops the canvas service and cancels subscriptions.
func (cs *CanvasService) Stop() {
	if cs.cancel != nil {
		cs.cancel()
	}
}

// processEvents processes canvas events from the workspace subscription.
func (cs *CanvasService) processEvents(eventChan <-chan webuiatoms.CanvasEvent, errChan <-chan error) {
	for {
		select {
		case <-cs.ctx.Done():
			return
		case event, ok := <-eventChan:
			if !ok {
				return
			}
			// Update canvas tracker with new canvas_id
			cs.canvasTracker.UpdateCanvas(event.CanvasID, event.CanvasName)
		case err, ok := <-errChan:
			if !ok {
				return
			}
			// Log error (will be handled by error handling system)
			fmt.Printf("Canvas service error: %v\n", err)
			// Reconnection is handled by workspace_subscriber
		}
	}
}

// GetCanvasID returns the current canvas_id.
func (cs *CanvasService) GetCanvasID() string {
	return cs.canvasTracker.GetCanvasID()
}

// GetCanvasName returns the current canvas_name.
func (cs *CanvasService) GetCanvasName() string {
	return cs.canvasTracker.GetCanvasName()
}

// GetInstallationName returns the installation name.
func (cs *CanvasService) GetInstallationName() string {
	return cs.installationName
}

// GetClientID returns the resolved client_id.
func (cs *CanvasService) GetClientID() string {
	return cs.clientID
}

// IsConnected returns whether the service is connected and tracking.
func (cs *CanvasService) IsConnected() bool {
	return cs.clientID != "" && cs.workspaceSubscriber != nil
}

