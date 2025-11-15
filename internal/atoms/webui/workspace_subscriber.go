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

// WorkspaceSubscriber handles SSE subscription to Canvus workspace endpoint.
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
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Subscribe connects to the workspace SSE endpoint and streams canvas_id updates.
// Returns a channel of CanvasEvent and an error channel.
func (ws *WorkspaceSubscriber) Subscribe(ctx context.Context) (<-chan CanvasEvent, <-chan error) {
	eventChan := make(chan CanvasEvent, 10)
	errChan := make(chan error, 1)

	go func() {
		defer close(eventChan)
		defer close(errChan)

		url := fmt.Sprintf("%s/clients/%s/workspaces/0/?subscribe", ws.apiBaseURL, ws.clientID)

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

// connectAndStream establishes SSE connection and streams events.
func (ws *WorkspaceSubscriber) connectAndStream(ctx context.Context, url string, eventChan chan<- CanvasEvent) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ws.authToken))
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := ws.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			line := scanner.Text()
			if strings.HasPrefix(line, "data: ") {
				data := strings.TrimPrefix(line, "data: ")
				event := ws.parseEvent(data)
				if event != nil {
					eventChan <- *event
				}
			}
		}
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
}

// parseEvent parses an SSE data line and extracts canvas_id.
func (ws *WorkspaceSubscriber) parseEvent(data string) *CanvasEvent {
	var eventData map[string]interface{}
	if err := json.Unmarshal([]byte(data), &eventData); err != nil {
		// Not JSON, try to extract canvas_id directly
		if strings.Contains(data, "canvas_id") {
			// Simple extraction - will be improved with actual API format
			return &CanvasEvent{
				CanvasID:  extractCanvasID(data),
				Timestamp: time.Now(),
			}
		}
		return nil
	}

	canvasID, ok := eventData["canvas_id"].(string)
	if !ok {
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
