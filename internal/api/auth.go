package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type generateRequest struct {
	Nickname    Nickname        `json:"nickname"`
	Description UserDescription `json:"description"`
	Interests   []Interest      `json:"interests"`
	Avatar      *FileDescriptor `json:"avatar"`
}

type generateResponse struct {
	Id         UserId         `json:"id"`
	AccessHash UserAccessHash `json:"accessHash"`
	Token      Token          `json:"token"`
}

// Generate makes request for creating account using provided data and returns Authorization structure.
func (c *Client) Generate(nickname Nickname, description UserDescription, interests []Interest, avatar *FileDescriptor) (*Authorization, error) {
	req := generateRequest{
		Nickname:    nickname,
		Description: description,
		Interests:   interests,
		Avatar:      avatar,
	}

	resp, err := c.do("POST", "/auth/generate", nil, req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("generate failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var genResp generateResponse
	if err := json.Unmarshal(bodyBytes, &genResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w, body: %s", err, string(bodyBytes))
	}

	return &Authorization{
		Id:         genResp.Id,
		AccessHash: genResp.AccessHash,
		Token:      genResp.Token,
	}, nil
}
