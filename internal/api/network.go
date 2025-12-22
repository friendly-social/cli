package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetNetworkDetails returns NetworkDetails for provided Authorization.
func (c *Client) GetNetworkDetails(auth *Authorization) (*NetworkDetails, error) {
	resp, err := c.do("GET", "/network/details", auth, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get network failed: status %d", resp.StatusCode)
	}

	var network NetworkDetails
	if err := json.NewDecoder(resp.Body).Decode(&network); err != nil {
		return nil, err
	}

	return &network, nil
}
