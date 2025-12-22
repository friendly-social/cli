package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetFeedQueue returns FeedQueue for provided Authorization.
func (c *Client) GetFeedQueue(auth *Authorization) (*FeedQueue, error) {
	resp, err := c.do("GET", "/feed/queue", auth, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get feed failed: status %d", resp.StatusCode)
	}

	var feed FeedQueue
	if err := json.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, err
	}

	return &feed, nil
}
