// Package friendly provides a Go SDK for the Friendly API
package friendly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

// Core types - these serialize directly as their wrapped type
type UserId int64

type UserAccessHash string

type Token string

type FriendToken string

type FileId int64

type FileAccessHash string

type Nickname string

type UserDescription string

type Interest string

type FileDescriptor struct {
	Id         FileId         `json:"id"`
	AccessHash FileAccessHash `json:"accessHash"`
}

type Authorization struct {
	Id         UserId         `json:"id"`
	AccessHash UserAccessHash `json:"accessHash"`
	Token      Token          `json:"token"`
}

type UserDetails struct {
	Id          UserId          `json:"id"`
	AccessHash  UserAccessHash  `json:"accessHash"`
	Nickname    Nickname        `json:"nickname"`
	Description UserDescription `json:"description"`
	Interests   []Interest      `json:"interests"`
	Avatar      *FileDescriptor `json:"avatar"`
}

type NetworkDetails struct {
	Friends []UserDetails `json:"friends"`
}

type FeedEntry struct {
	IsExtendedNetwork bool          `json:"isExtendedNetwork"`
	CommonFriends     []UserDetails `json:"commonFriends"`
	Details           UserDetails   `json:"details"`
}

type FeedQueue struct {
	Entries []FeedEntry `json:"entries"`
}

// Validation functions
func NewNickname(s string) (Nickname, error) {
	if len(s) > 256 {
		return "", fmt.Errorf("nickname too long: %d > 256", len(s))
	}
	return Nickname(s), nil
}

func MustNickname(s string) Nickname {
	n, err := NewNickname(s)
	if err != nil {
		panic(err)
	}
	return n
}

func NewUserDescription(s string) (UserDescription, error) {
	if len(s) > 1024 {
		return "", fmt.Errorf("description too long: %d > 1024", len(s))
	}
	return UserDescription(s), nil
}

func MustUserDescription(s string) UserDescription {
	d, err := NewUserDescription(s)
	if err != nil {
		panic(err)
	}
	return d
}

func NewInterest(s string) (Interest, error) {
	if len(s) > 64 {
		return "", fmt.Errorf("interest too long: %d > 64", len(s))
	}
	return Interest(s), nil
}

func MustInterest(s string) Interest {
	i, err := NewInterest(s)
	if err != nil {
		panic(err)
	}
	return i
}

func NewToken(s string) (Token, error) {
	if len(s) != 256 {
		return "", fmt.Errorf("token must be 256 characters, got %d", len(s))
	}
	return Token(s), nil
}

func NewUserAccessHash(s string) (UserAccessHash, error) {
	if len(s) != 256 {
		return "", fmt.Errorf("access hash must be 256 characters, got %d", len(s))
	}
	return UserAccessHash(s), nil
}

func NewFileAccessHash(s string) (FileAccessHash, error) {
	if len(s) != 256 {
		return "", fmt.Errorf("file access hash must be 256 characters, got %d", len(s))
	}
	return FileAccessHash(s), nil
}

func NewFriendToken(s string) (FriendToken, error) {
	if len(s) != 256 {
		return "", fmt.Errorf("friend token must be 256 characters, got %d", len(s))
	}
	return FriendToken(s), nil
}

// Client
type Client struct {
	endpoint   string
	httpClient *http.Client
}

func NewClient(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func NewLocalhostClient(port int) *Client {
	return NewClient(fmt.Sprintf("http://localhost:%d", port))
}

func NewMeetacyClient() *Client {
	return NewClient("https://meetacy.app/friendly")
}

// HTTP helpers
func (c *Client) doRequest(method, path string, auth *Authorization, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequest(method, c.endpoint+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if auth != nil {
		req.Header.Set("X-User-Id", fmt.Sprintf("%d", auth.Id))
		req.Header.Set("X-Token", string(auth.Token))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

// Auth API
type generateRequest struct {
	Nickname    Nickname        `json:"nickname"`
	Description UserDescription `json:"description"`
	Interests   []Interest      `json:"interests"`
	Avatar      *FileDescriptor `json:"avatar"`
}

type generateResponse struct {
	Token      Token          `json:"token"`
	Id         UserId         `json:"id"`
	AccessHash UserAccessHash `json:"accessHash"`
}

func (c *Client) Generate(nickname Nickname, description UserDescription, interests []Interest, avatar *FileDescriptor) (*Authorization, error) {
	req := generateRequest{
		Nickname:    nickname,
		Description: description,
		Interests:   interests,
		Avatar:      avatar,
	}

	resp, err := c.doRequest("POST", "/auth/generate", nil, req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read body for error details
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

// Users API
func (c *Client) GetSelfDetails(auth *Authorization) (*UserDetails, error) {
	resp, err := c.doRequest("GET", "/users/details", auth, nil)
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

func (c *Client) GetUserDetails(auth *Authorization, userId UserId, accessHash UserAccessHash) (*UserDetails, error) {
	path := fmt.Sprintf("/users/details/%d/%s", userId, accessHash)
	resp, err := c.doRequest("GET", path, auth, nil)
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

// Friends API
type generateFriendTokenResponse struct {
	Token FriendToken `json:"token"`
}

func (c *Client) GenerateFriendToken(auth *Authorization) (FriendToken, error) {
	resp, err := c.doRequest("POST", "/friends/generate", auth, nil)
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

type addFriendRequest struct {
	Token  FriendToken `json:"token"`
	UserId UserId      `json:"userId"`
}

type addFriendResponse struct {
	Type string `json:"type"`
}

func (c *Client) AddFriend(auth *Authorization, token FriendToken, userId UserId) error {
	req := addFriendRequest{
		Token:  token,
		UserId: userId,
	}

	resp, err := c.doRequest("POST", "/friends/add", auth, req)
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

type friendRequestRequest struct {
	UserId     UserId         `json:"userId"`
	AccessHash UserAccessHash `json:"userAccessHash"`
}

func (c *Client) SendFriendRequest(auth *Authorization, userId UserId, accessHash UserAccessHash) error {
	req := friendRequestRequest{
		UserId:     userId,
		AccessHash: accessHash,
	}

	resp, err := c.doRequest("POST", "/friends/request", auth, req)
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

func (c *Client) DeclineFriendRequest(auth *Authorization, userId UserId, accessHash UserAccessHash) error {
	req := friendRequestRequest{
		UserId:     userId,
		AccessHash: accessHash,
	}

	resp, err := c.doRequest("POST", "/friends/decline", auth, req)
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

// Network API
func (c *Client) GetNetworkDetails(auth *Authorization) (*NetworkDetails, error) {
	resp, err := c.doRequest("GET", "/network/details", auth, nil)
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

// Feed API
func (c *Client) GetFeedQueue(auth *Authorization) (*FeedQueue, error) {
	resp, err := c.doRequest("GET", "/feed/queue", auth, nil)
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

// Files API
type uploadFileResponse struct {
	Id         FileId         `json:"id"`
	AccessHash FileAccessHash `json:"accessHash"`
}

func (c *Client) UploadFile(filename string, contentType string, reader io.Reader) (*FileDescriptor, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(part, reader); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.endpoint+"/files/upload", body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upload failed: status %d", resp.StatusCode)
	}

	var uploadResp uploadFileResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return nil, err
	}

	return &FileDescriptor{
		Id:         uploadResp.Id,
		AccessHash: uploadResp.AccessHash,
	}, nil
}

func (c *Client) GetFileURL(descriptor *FileDescriptor) string {
	return fmt.Sprintf("%s/files/download/%d/%s",
		c.endpoint, descriptor.Id, descriptor.AccessHash)
}
