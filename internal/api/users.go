package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetSelfDetails returns UserDetails for provided Authorization data.
func (c *Client) GetSelfDetails(auth *Authorization) (*UserDetails, error) {
	resp, err := c.do("GET", "/users/details", auth, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get details failed: status %d", resp.StatusCode)
	}

	var details UserDetails
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return nil, err
	}

	return &details, nil
}

// GetUserDetails returns UserDetails for provided user's ID and AccessHash from provided Authorization's perspective.
func (c *Client) GetUserDetails(auth *Authorization, userId UserId, accessHash UserAccessHash) (*UserDetails, error) {
	path := fmt.Sprintf("/users/details/%d/%s", userId, accessHash)
	resp, err := c.do("GET", path, auth, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get user details failed: status %d", resp.StatusCode)
	}

	var details UserDetails
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return nil, err
	}

	return &details, nil
}
