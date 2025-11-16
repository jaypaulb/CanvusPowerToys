package webui

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

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
	clientName          string // Actual client name from server
	installationName    string
	overrideClientName  string // Manual override for client name to monitor
	hasReceivedEvents   bool   // Track if we've received events (proves subscription works)
	lastEventTime       time.Time
	subscriptionStartTime time.Time // Track when subscription started
}

// NewCanvasService creates a new canvas service.
func NewCanvasService(fileService *services.FileService, apiBaseURL, authToken string) (*CanvasService, error) {
	fmt.Printf("[CanvasService] NewCanvasService called with apiBaseURL: '%s', authToken length: %d\n", apiBaseURL, len(authToken))
	clientResolver := webuiatoms.NewClientResolver(fileService)
	canvasTracker := webuiatoms.NewCanvasTracker()

	// Get installation name
	installationName, err := clientResolver.GetInstallationName()
	if err != nil {
		fmt.Printf("[CanvasService] ERROR: Failed to get installation name: %v\n", err)
		return nil, fmt.Errorf("failed to get installation name: %w", err)
	}
	fmt.Printf("[CanvasService] Installation name: '%s'\n", installationName)

	ctx, cancel := context.WithCancel(context.Background())

	cs := &CanvasService{
		clientResolver:      clientResolver,
		canvasTracker:      canvasTracker,
		ctx:                 ctx,
		cancel:              cancel,
		apiBaseURL:          apiBaseURL,
		authToken:           authToken,
		installationName:    installationName,
		workspaceSubscriber: nil, // Will be created after client_id resolution
	}
	fmt.Printf("[CanvasService] Created CanvasService with apiBaseURL: '%s'\n", cs.apiBaseURL)
	return cs, nil
}

// Start initializes client_id resolution and starts workspace subscription.
// Returns error if resolution fails, but service can still be used for manual override.
func (cs *CanvasService) Start() error {
	// If clientResolver is nil, service was created in minimal mode (auto-detection failed)
	// User will need to manually override via WebUI
	if cs.clientResolver == nil {
		fmt.Printf("[CanvasService] Running in minimal mode - client override required via WebUI\n")
		return nil // Not an error - just needs manual override
	}

	// Resolve client_id from installation_name
	clientID, err := cs.clientResolver.ResolveClientID(cs.apiBaseURL, cs.authToken, cs.installationName)
	if err != nil {
		return fmt.Errorf("failed to resolve client_id: %w", err)
	}

	cs.clientID = clientID

	// Fetch client name from server to verify it exists
	cs.fetchClientName()

	// Create workspace subscriber
	cs.workspaceSubscriber = webuiatoms.NewWorkspaceSubscriber(
		cs.clientID,
		cs.apiBaseURL,
		cs.authToken,
	)

	// Start subscription
	cs.subscriptionStartTime = time.Now()
	eventChan, errChan := cs.workspaceSubscriber.Subscribe(cs.ctx)

	// Process events in background
	go cs.processEvents(eventChan, errChan)

	// Also start polling fallback - fetch canvas_id directly from workspace API
	// This ensures we get canvas_id even if SSE subscription has issues
	go cs.pollWorkspaceCanvasID()

	// Fetch initial canvas name if we have a canvas_id but no canvas_name
	go func() {
		// Wait a moment for initial events to arrive
		time.Sleep(2 * time.Second)
		canvasID := cs.canvasTracker.GetCanvasID()
		canvasName := cs.canvasTracker.GetCanvasName()
		if canvasID != "" && canvasName == "" {
			// Fetch canvas name from API
			fetchedName, err := cs.fetchCanvasName(canvasID)
			if err == nil && fetchedName != "" {
				cs.canvasTracker.UpdateCanvas(canvasID, fetchedName)
			}
		}
	}()

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
			// Mark that we've received events (proves subscription is working)
			cs.hasReceivedEvents = true
			cs.lastEventTime = time.Now()

			// If canvas_name is empty, fetch it from the API
			canvasName := event.CanvasName
			if canvasName == "" && event.CanvasID != "" {
				// Fetch canvas name from API
				fetchedName, err := cs.fetchCanvasName(event.CanvasID)
				if err == nil && fetchedName != "" {
					canvasName = fetchedName
				}
			}
			// Update canvas tracker with new canvas_id and canvas_name
			cs.canvasTracker.UpdateCanvas(event.CanvasID, canvasName)
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

// fetchCanvasName fetches the canvas name from the Canvus API using canvas_id.
// Endpoint: GET /api/v1/canvases/{canvasID}
func (cs *CanvasService) fetchCanvasName(canvasID string) (string, error) {
	if canvasID == "" {
		return "", fmt.Errorf("canvas_id is empty")
	}

	fmt.Printf("[CanvasService] fetchCanvasName: Fetching canvas name for canvasID: %s\n", canvasID)
	apiClient := webuiatoms.NewAPIClient(cs.apiBaseURL, cs.authToken)
	endpoint := fmt.Sprintf("/api/v1/canvases/%s", canvasID)
	fmt.Printf("[CanvasService] fetchCanvasName: Calling endpoint: %s\n", endpoint)

	data, err := apiClient.Get(endpoint)
	if err != nil {
		fmt.Printf("[CanvasService] fetchCanvasName: ERROR - Failed to fetch canvas: %v\n", err)
		return "", fmt.Errorf("failed to fetch canvas: %w", err)
	}

	var canvas map[string]interface{}
	if err := json.Unmarshal(data, &canvas); err != nil {
		fmt.Printf("[CanvasService] fetchCanvasName: ERROR - Failed to parse canvas JSON: %v\n", err)
		fmt.Printf("[CanvasService] fetchCanvasName: Response body: %s\n", string(data))
		return "", fmt.Errorf("failed to parse canvas: %w", err)
	}

	if name, ok := canvas["name"].(string); ok && name != "" {
		fmt.Printf("[CanvasService] fetchCanvasName: Success - Found canvas name: '%s'\n", name)
		return name, nil
	}

	fmt.Printf("[CanvasService] fetchCanvasName: ERROR - Canvas name not found in response. Response keys: %v\n", getMapKeys(canvas))
	return "", fmt.Errorf("canvas name not found in response")
}

// getMapKeys returns the keys of a map as a slice (for debugging)
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// GetCanvasID returns the current canvas_id.
func (cs *CanvasService) GetCanvasID() string {
	if cs.canvasTracker == nil {
		return ""
	}
	return cs.canvasTracker.GetCanvasID()
}

// GetCanvasName returns the current canvas_name.
// If canvas_name is empty but we have a canvas_id, attempts to fetch it from the API.
func (cs *CanvasService) GetCanvasName() string {
	if cs.canvasTracker == nil {
		return ""
	}

	canvasName := cs.canvasTracker.GetCanvasName()
	canvasID := cs.canvasTracker.GetCanvasID()

	// If we have a canvas_id but no canvas_name, try to fetch it
	if canvasID != "" && canvasName == "" {
		fmt.Printf("[CanvasService] GetCanvasName: Canvas ID exists but name is empty, fetching...\n")
		// Fetch asynchronously to avoid blocking
		go func() {
			fetchedName, err := cs.fetchCanvasName(canvasID)
			if err == nil && fetchedName != "" {
				cs.canvasTracker.UpdateCanvas(canvasID, fetchedName)
				fmt.Printf("[CanvasService] GetCanvasName: Updated canvas name to: '%s'\n", fetchedName)
			}
		}()
	}

	return canvasName
}

// GetInstallationName returns the installation name.
func (cs *CanvasService) GetInstallationName() string {
	return cs.installationName
}

// GetClientID returns the resolved client_id.
func (cs *CanvasService) GetClientID() string {
	return cs.clientID
}

// GetClientName returns the actual client name from the server (if available).
func (cs *CanvasService) GetClientName() string {
	return cs.clientName
}

// IsConnected returns whether the service is connected and tracking.
// This checks if we have a clientID and workspace subscriber.
// Having a clientID means we've successfully resolved a client, which is sufficient for connection.
// Note: clientName may be empty if the client exists but has no name, or if fetchClientName hasn't run yet.
func (cs *CanvasService) IsConnected() bool {
	// Must have a clientID and workspace subscriber to be considered connected
	if cs.clientID == "" || cs.workspaceSubscriber == nil {
		return false
	}
	// If we've received events, we're definitely connected
	if cs.hasReceivedEvents {
		return true
	}
	// If no events received yet, allow grace period of 30 seconds after subscription start
	// (allows for initial connection time)
	if !cs.subscriptionStartTime.IsZero() {
		return time.Since(cs.subscriptionStartTime) < 30*time.Second
	}
	// If we have a clientID and subscriber, consider connected (even if no events yet)
	// This handles cases where events are slow or subscription is still establishing
	return true
}

// fetchClientName fetches the client name from the server to verify it exists.
// This is informational only - doesn't affect connection status.
// A client may exist but have an empty name, which is valid.
func (cs *CanvasService) fetchClientName() {
	if cs.clientID == "" {
		return
	}

	apiClient := webuiatoms.NewAPIClient(cs.apiBaseURL, cs.authToken)
	clients, err := apiClient.GetClients()
	if err != nil {
		fmt.Printf("[CanvasService] Failed to fetch clients to verify client name: %v\n", err)
		// Don't clear clientName on error - keep existing value if any
		return
	}

	// Find client by ID
	for _, client := range clients {
		if client.ID == cs.clientID {
			cs.clientName = client.InstallationName
			return
		}
	}

	// Client not found by ID - log warning but don't clear clientName
	// The clientID might be valid but the API call might have failed or client list might be stale
	fmt.Printf("[CanvasService] WARNING: Client ID %s not found in clients list. This may be temporary. Available clients: %v\n", cs.clientID, clients)
	// Don't clear clientName - keep it if we had it before, as the clientID is still valid
}

// pollWorkspaceCanvasID polls the workspace API to get canvas_id directly (fallback to subscription).
// This ensures we get canvas_id even if SSE subscription has issues.
// Only polls if we don't have canvas_id yet, then stops.
func (cs *CanvasService) pollWorkspaceCanvasID() {
	if cs.clientID == "" {
		return
	}

	// Wait a bit before first poll to give SSE subscription a chance
	time.Sleep(2 * time.Second)

	apiClient := webuiatoms.NewAPIClient(cs.apiBaseURL, cs.authToken)

	// Poll up to 6 times (30 seconds total) or until we get canvas_id
	maxAttempts := 6
	for attempt := 0; attempt < maxAttempts; attempt++ {
		select {
		case <-cs.ctx.Done():
			return
		default:
			// Check if we already have canvas_id (from SSE subscription)
			if cs.canvasTracker.GetCanvasID() != "" {
				fmt.Printf("[CanvasService] Polling stopped - canvas_id obtained from subscription\n")
				return
			}

			// Fetch workspace 0 data directly from API
			endpoint := fmt.Sprintf("/api/v1/clients/%s/workspaces/0", cs.clientID)
			fmt.Printf("[CanvasService] Polling workspace (attempt %d/%d)...\n", attempt+1, maxAttempts)
			body, err := apiClient.Get(endpoint)
			if err != nil {
				fmt.Printf("[CanvasService] Polling workspace failed: %v\n", err)
				time.Sleep(5 * time.Second)
				continue
			}

			// Parse workspace data to extract canvas_id
			var workspaceData map[string]interface{}
			if err := json.Unmarshal(body, &workspaceData); err != nil {
				fmt.Printf("[CanvasService] Failed to parse workspace data: %v\n", err)
				time.Sleep(5 * time.Second)
				continue
			}

			canvasID, ok := workspaceData["canvas_id"].(string)
			if !ok || canvasID == "" {
				fmt.Printf("[CanvasService] No canvas_id in workspace data yet\n")
				time.Sleep(5 * time.Second)
				continue
			}

			canvasName, _ := workspaceData["canvas_name"].(string)
			cs.canvasTracker.UpdateCanvas(canvasID, canvasName)
			fmt.Printf("[CanvasService] Polling fallback: Found canvas_id=%s, canvas_name=%s\n", canvasID, canvasName)
			return // Success, stop polling
		}
	}
	fmt.Printf("[CanvasService] Polling stopped after %d attempts - canvas_id not found\n", maxAttempts)
}

// OverrideClient manually sets a client name to monitor instead of using installation name.
func (cs *CanvasService) OverrideClient(clientName string) error {
	fmt.Printf("[CanvasService] OverrideClient called with clientName: '%s'\n", clientName)
	if clientName == "" {
		// Clear override - use installation name again
		cs.overrideClientName = ""
		// Restart with installation name (if available)
		if cs.installationName != "" && cs.installationName != "Unknown" {
			return cs.restartWithClientName(cs.installationName)
		}
		return fmt.Errorf("cannot clear override: no installation name available")
	}

	cs.overrideClientName = clientName
	// Restart subscription with new client name
	return cs.restartWithClientName(clientName)
}

// restartWithClientName restarts the workspace subscription with a specific client name.
func (cs *CanvasService) restartWithClientName(clientName string) error {
	fmt.Printf("[CanvasService] restartWithClientName called with clientName: '%s'\n", clientName)
	fmt.Printf("[CanvasService] API Base URL: %s\n", cs.apiBaseURL)

	// Stop current subscription
	if cs.workspaceSubscriber != nil {
		fmt.Printf("[CanvasService] Stopping current subscription\n")
		cs.Stop()
	}

	// Create new context
	cs.ctx, cs.cancel = context.WithCancel(context.Background())

	// Always use direct API lookup for manual override (matches by client name, not installation_name)
	// This ensures we can override to any client by name, regardless of installation_name
	apiClient := webuiatoms.NewAPIClient(cs.apiBaseURL, cs.authToken)
	fmt.Printf("[CanvasService] Created API client, fetching clients list from: %s/api/v1/clients\n", cs.apiBaseURL)

	// Get clients list and find matching client
	clients, err := apiClient.GetClients()
	if err != nil {
		fmt.Printf("[CanvasService] ERROR: Failed to get clients list: %v\n", err)
		return fmt.Errorf("failed to get clients list: %w", err)
	}
	fmt.Printf("[CanvasService] Successfully fetched %d clients from API\n", len(clients))

	// Log all available clients for debugging
	fmt.Printf("[CanvasService] Available clients from API:\n")
	for i, client := range clients {
		fmt.Printf("[CanvasService]   [%d] ID: %s, InstallationName: '%s'\n", i, client.ID, client.InstallationName)
	}

	// Find client by installation_name (case-insensitive match)
	var clientID string
	var foundClientName string
	clientNameLower := strings.ToLower(clientName)
	fmt.Printf("[CanvasService] Searching for client with installation_name (case-insensitive): '%s'\n", clientName)
	for _, client := range clients {
		installationNameLower := strings.ToLower(client.InstallationName)
		fmt.Printf("[CanvasService]   Comparing '%s' with installation_name '%s' (lower: '%s')\n",
			clientName, client.InstallationName, installationNameLower)

		if installationNameLower == clientNameLower {
			clientID = client.ID
			foundClientName = client.InstallationName
			fmt.Printf("[CanvasService] MATCH FOUND! Client ID: %s, InstallationName: '%s'\n", clientID, foundClientName)
			break
		}
	}

	if clientID == "" {
		// Provide helpful error message with available client installation_names
		availableNames := make([]string, 0, len(clients))
		for _, client := range clients {
			availableNames = append(availableNames, client.InstallationName)
		}
		fmt.Printf("[CanvasService] ERROR: No client found with installation_name '%s'. Available clients: %v\n", clientName, availableNames)
		if len(availableNames) > 0 {
			return fmt.Errorf("client not found: no client with installation_name '%s'. Available clients: %v", clientName, availableNames)
		}
		return fmt.Errorf("client not found: no client with installation_name '%s'", clientName)
	}

	cs.clientID = clientID
	fmt.Printf("[CanvasService] Set clientID to: %s\n", cs.clientID)
	// Update clientName to the actual name from server (in case we matched by installation_name)
	if foundClientName != "" {
		cs.clientName = foundClientName
		fmt.Printf("[CanvasService] Set clientName to: '%s'\n", cs.clientName)
	}

	// Reset event tracking
	cs.hasReceivedEvents = false
	cs.lastEventTime = time.Time{}
	fmt.Printf("[CanvasService] Reset event tracking\n")

	// Fetch client name from server to verify it exists
	fmt.Printf("[CanvasService] Fetching client name from server...\n")
	cs.fetchClientName()
	fmt.Printf("[CanvasService] Client name after fetch: '%s'\n", cs.clientName)

	// Create new workspace subscriber
	fmt.Printf("[CanvasService] Creating workspace subscriber for clientID: %s\n", cs.clientID)
	cs.workspaceSubscriber = webuiatoms.NewWorkspaceSubscriber(
		cs.clientID,
		cs.apiBaseURL,
		cs.authToken,
	)

	// Start subscription
	cs.subscriptionStartTime = time.Now()
	fmt.Printf("[CanvasService] Starting workspace subscription...\n")
	eventChan, errChan := cs.workspaceSubscriber.Subscribe(cs.ctx)

	// Process events in background
	fmt.Printf("[CanvasService] Starting event processing goroutine\n")
	go cs.processEvents(eventChan, errChan)

	// Also start polling fallback - fetch canvas_id directly from workspace API
	// This ensures we get canvas_id even if SSE subscription has issues
	fmt.Printf("[CanvasService] Starting polling fallback goroutine\n")
	go cs.pollWorkspaceCanvasID()

	fmt.Printf("[CanvasService] restartWithClientName completed successfully\n")

	// Fetch initial canvas name if we have a canvas_id but no canvas_name
	go func() {
		// Wait a moment for initial events to arrive
		time.Sleep(2 * time.Second)
		canvasID := cs.canvasTracker.GetCanvasID()
		canvasName := cs.canvasTracker.GetCanvasName()
		if canvasID != "" && canvasName == "" {
			// Fetch canvas name from API
			fetchedName, err := cs.fetchCanvasName(canvasID)
			if err == nil && fetchedName != "" {
				cs.canvasTracker.UpdateCanvas(canvasID, fetchedName)
			}
		}
	}()

	return nil
}

