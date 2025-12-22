package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is an entity for interacting with Friendly API.
type Client struct {
	url  string
	http *http.Client
}

// NewClient creates basic Client with provided URL.
func NewClient(endpoint string) *Client {
	return &Client{
		url: endpoint,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewLocalhostClient creates Client with localhost URL and provided port.
func NewLocalhostClient(port int) *Client {
	return NewClient(fmt.Sprintf("http://localhost:%d", port))
}

// NewMeetacyClient creates Client with Meetacy URL.
func NewMeetacyClient() *Client {
	return NewClient("https://meetacy.app/friendly")
}

// do makes HTTP request to given path using provided data and returns HTTP response or error if something went wrong.
func (c *Client) do(method, path string, auth *Authorization, body any) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}

		bodyReader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequest(method, c.url+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if auth != nil {
		req.Header.Set("X-User-Id", fmt.Sprintf("%d", auth.Id))
		req.Header.Set("X-Token", string(auth.Token))
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}
