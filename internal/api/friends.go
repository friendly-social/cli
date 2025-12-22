package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type generateFriendTokenResponse struct {
	Token FriendToken `json:"token"`
}

type addFriendRequest struct {
	Token  FriendToken `json:"token"`
	UserId UserId      `json:"userId"`
}

type addFriendResponse struct {
	Type string `json:"type"`
}

type friendRequestRequest struct {
	UserId     UserId         `json:"userId"`
	AccessHash UserAccessHash `json:"userAccessHash"`
}

// GenerateFriendToken creates token for Authorization's user by which another users can add them.
func (c *Client) GenerateFriendToken(auth *Authorization) (FriendToken, error) {
	resp, err := c.do("POST", "/friends/generate", auth, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf("unauthorized")
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("generate token failed: status %d", resp.StatusCode)
	}

	var tokenResp generateFriendTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	return tokenResp.Token, nil
}

// AddFriend makes request to add user with provided Token and ID to Authorization's friends list.
func (c *Client) AddFriend(auth *Authorization, token FriendToken, userId UserId) error {
	req := addFriendRequest{
		Token:  token,
		UserId: userId,
	}

	resp, err := c.do("POST", "/friends/add", auth, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("unauthorized")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("add friend failed: status %d", resp.StatusCode)
	}

	var addResp addFriendResponse
	if err := json.NewDecoder(resp.Body).Decode(&addResp); err != nil {
		return err
	}

	if addResp.Type == "FriendTokenExpired" {
		return fmt.Errorf("friend token expired")
	}

	return nil
}

// SendFriendRequest sends friend request from Authorization to user with provided ID and AccessHash.
func (c *Client) SendFriendRequest(auth *Authorization, userId UserId, accessHash UserAccessHash) error {
	req := friendRequestRequest{
		UserId:     userId,
		AccessHash: accessHash,
	}

	resp, err := c.do("POST", "/friends/request", auth, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("unauthorized")
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("user not found")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("send request failed: status %d", resp.StatusCode)
	}

	return nil
}

// DeclineFriendRequest declines Authorization's request from user with provided ID and AccessHash.
func (c *Client) DeclineFriendRequest(auth *Authorization, userId UserId, accessHash UserAccessHash) error {
	req := friendRequestRequest{
		UserId:     userId,
		AccessHash: accessHash,
	}

	resp, err := c.do("POST", "/friends/decline", auth, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("unauthorized")
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("user not found")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("decline request failed: status %d", resp.StatusCode)
	}

	return nil
}
