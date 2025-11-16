package webui

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// WorkspaceSubscriber handles TCP JSON streaming subscription to Canvus workspace endpoint.
// The MTCS sends one JSON block per line, with \n as keepalive.
type WorkspaceSubscriber struct {
	clientID   string
	apiBaseURL string
	authToken  string
	httpClient *http.Client
}

// CanvasEvent represents a canvas_id update event from the workspace subscription.
type CanvasEvent struct {
	CanvasID   string
	CanvasName string
	Timestamp  time.Time
}

// NewWorkspaceSubscriber creates a new workspace subscriber.
func NewWorkspaceSubscriber(clientID, apiBaseURL, authToken string) *WorkspaceSubscriber {
	return &WorkspaceSubscriber{
		clientID:   clientID,
		apiBaseURL: strings.TrimSuffix(apiBaseURL, "/"),
		authToken:  authToken,
		// No timeout - TCP JSON streaming connection stays open indefinitely
		httpClient: &http.Client{
			Timeout: 0, // No timeout for streaming connections
		},
	}
}

// Subscribe connects to the workspace TCP JSON streaming endpoint and streams canvas_id updates.
// MTCS sends one JSON block per line, with \n as keepalive.
// Returns a channel of CanvasEvent and an error channel.
func (ws *WorkspaceSubscriber) Subscribe(ctx context.Context) (<-chan CanvasEvent, <-chan error) {
	eventChan := make(chan CanvasEvent, 10)
	errChan := make(chan error, 1)

	go func() {
		defer close(eventChan)
		defer close(errChan)

		// Ensure URL includes /api/v1 if not already present
		baseURL := ws.apiBaseURL
		if !strings.Contains(baseURL, "/api/v1") && !strings.Contains(baseURL, "/api") {
			baseURL = strings.TrimSuffix(baseURL, "/") + "/api/v1"
		}
		url := fmt.Sprintf("%s/clients/%s/workspaces/0/?subscribe", baseURL, ws.clientID)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := ws.connectAndStream(ctx, url, eventChan); err != nil {
					errChan <- err
					// Wait before reconnecting
					select {
					case <-ctx.Done():
						return
					case <-time.After(5 * time.Second):
						// Retry connection
					}
				}
			}
		}
	}()

	return eventChan, errChan
}

// connectAndStream establishes TCP JSON streaming connection and streams events.
// MTCS sends one JSON block per line, with \n as keepalive (empty lines).
func (ws *WorkspaceSubscriber) connectAndStream(ctx context.Context, url string, eventChan chan<- CanvasEvent) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Use Private-Token header (same as API client)
	req.Header.Set("Private-Token", ws.authToken)
	// Note: Not SSE - this is TCP JSON streaming, so no SSE headers needed

	resp, err := ws.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read line by line - each line is a JSON object, empty lines are keepalive
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			line := strings.TrimSpace(scanner.Text())
			// Skip empty lines (keepalive - just \n)
			if line == "" {
				continue
			}

			// Each non-empty line is a JSON object
			event := ws.parseEvent(line)
			if event != nil {
				eventChan <- *event
			}
		}
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
}

// parseEvent parses a JSON line and extracts canvas_id.
// Each line from MTCS is a complete JSON object.
func (ws *WorkspaceSubscriber) parseEvent(jsonLine string) *CanvasEvent {
	var eventData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonLine), &eventData); err != nil {
		// Not valid JSON - skip this line
		fmt.Printf("[WorkspaceSubscriber] Failed to parse JSON line: %v, line: %s\n", err, jsonLine)
		return nil
	}

	canvasID, ok := eventData["canvas_id"].(string)
	if !ok || canvasID == "" {
		// No canvas_id in this event - might be other workspace data
		return nil
	}

	canvasName, _ := eventData["canvas_name"].(string)

	return &CanvasEvent{
		CanvasID:   canvasID,
		CanvasName: canvasName,
		Timestamp:  time.Now(),
	}
}

// extractCanvasID extracts canvas_id from a string (fallback method).
func extractCanvasID(data string) string {
	// This is a placeholder - actual implementation depends on API format
	// Will be updated when we have actual API documentation
	if idx := strings.Index(data, `"canvas_id":"`); idx != -1 {
		start := idx + len(`"canvas_id":"`)
		if end := strings.Index(data[start:], `"`); end != -1 {
			return data[start : start+end]
		}
	}
	return ""
}
