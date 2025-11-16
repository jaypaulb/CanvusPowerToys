package webui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// APIClient handles authenticated requests to the Canvus Server API.
type APIClient struct {
	baseURL   string
	authToken string
	httpClient *http.Client
}

// NewAPIClient creates a new API client for Canvus Server.
func NewAPIClient(baseURL, authToken string) *APIClient {
	fmt.Printf("[APIClient] NewAPIClient created with baseURL: '%s', authToken length: %d\n", baseURL, len(authToken))
	return &APIClient{
		baseURL:   baseURL,
		authToken: authToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Get performs a GET request to the Canvus API.
func (c *APIClient) Get(endpoint string) ([]byte, error) {
	url := c.baseURL + endpoint
	fmt.Printf("[APIClient] GET %s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("[APIClient] ERROR: Failed to create request: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Private-Token", c.authToken)
	req.Header.Set("Content-Type", "application/json")
	fmt.Printf("[APIClient] Request headers: Private-Token=%s (length: %d)\n",
		func() string {
			if len(c.authToken) > 10 {
				return c.authToken[:10] + "..."
			}
			return c.authToken
		}(), len(c.authToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Printf("[APIClient] ERROR: Request failed: %v\n", err)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("[APIClient] Response status: %d %s\n", resp.StatusCode, resp.Status)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[APIClient] ERROR: Failed to read response: %v\n", err)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Printf("[APIClient] ERROR: API returned error status %d: %s\n", resp.StatusCode, string(body))
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	fmt.Printf("[APIClient] Success: Received %d bytes\n", len(body))
	return body, nil
}

// Post performs a POST request to the Canvus API.
func (c *APIClient) Post(endpoint string, data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	url := c.baseURL + endpoint
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Private-Token", c.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Put performs a PUT request to the Canvus API.
func (c *APIClient) Put(endpoint string, data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	url := c.baseURL + endpoint
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Private-Token", c.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Patch performs a PATCH request to the Canvus API.
func (c *APIClient) Patch(endpoint string, data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	url := c.baseURL + endpoint
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Private-Token", c.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Delete performs a DELETE request to the Canvus API.
func (c *APIClient) Delete(endpoint string) error {
	url := c.baseURL + endpoint
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Private-Token", c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	return nil
}

// PostMultipart performs a multipart/form-data POST request to the Canvus API.
// This is used for file uploads where the API expects:
// - json: JSON metadata as a form field
// - data: File binary data as a form field
func (c *APIClient) PostMultipart(endpoint string, jsonData map[string]interface{}, fileData io.Reader, fileName string) ([]byte, error) {
	url := c.baseURL + endpoint
	fmt.Printf("[APIClient] PostMultipart: %s\n", url)

	// Write JSON part
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json data: %w", err)
	}

	// Build multipart form manually
	boundary := "----WebKitFormBoundary" + fmt.Sprintf("%d", time.Now().UnixNano())

	var buf bytes.Buffer
	// JSON part
	buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	buf.WriteString("Content-Disposition: form-data; name=\"json\"\r\n")
	buf.WriteString("Content-Type: application/json\r\n\r\n")
	buf.Write(jsonBytes)
	buf.WriteString("\r\n")

	// File part
	buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	buf.WriteString(fmt.Sprintf("Content-Disposition: form-data; name=\"data\"; filename=\"%s\"\r\n", fileName))
	buf.WriteString("Content-Type: application/octet-stream\r\n\r\n")

	// Copy file data
	if _, err := io.Copy(&buf, fileData); err != nil {
		return nil, fmt.Errorf("failed to copy file data: %w", err)
	}

	buf.WriteString("\r\n")
	buf.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Private-Token", c.authToken)
	req.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", boundary))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Printf("[APIClient] ERROR: PostMultipart request failed: %v\n", err)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[APIClient] ERROR: PostMultipart failed to read response: %v\n", err)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	fmt.Printf("[APIClient] PostMultipart response status: %d %s\n", resp.StatusCode, resp.Status)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Printf("[APIClient] ERROR: PostMultipart API returned error status %d: %s\n", resp.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	fmt.Printf("[APIClient] PostMultipart success: Received %d bytes\n", len(bodyBytes))
	return bodyBytes, nil
}

// Client represents a Canvus client.
type Client struct {
	ID               string `json:"id"`
	InstallationName string `json:"installation_name"`
}

// GetClients fetches the list of clients from the Canvus API.
func (c *APIClient) GetClients() ([]Client, error) {
	url := c.baseURL + "/api/v1/clients"
	fmt.Printf("[APIClient] GetClients: Calling %s\n", url)
	body, err := c.Get("/api/v1/clients")
	if err != nil {
		fmt.Printf("[APIClient] ERROR: GetClients failed: %v\n", err)
		return nil, fmt.Errorf("failed to get clients: %w", err)
	}

	fmt.Printf("[APIClient] GetClients: Received %d bytes response\n", len(body))
	var clients []Client
	if err := json.Unmarshal(body, &clients); err != nil {
		fmt.Printf("[APIClient] ERROR: Failed to unmarshal clients: %v\n", err)
		fmt.Printf("[APIClient] Response body: %s\n", string(body))
		return nil, fmt.Errorf("failed to unmarshal clients: %w", err)
	}

	fmt.Printf("[APIClient] GetClients: Successfully parsed %d clients\n", len(clients))
	return clients, nil
}


