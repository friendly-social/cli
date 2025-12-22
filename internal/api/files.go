package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type uploadFileResponse struct {
	Id         FileId         `json:"id"`
	AccessHash FileAccessHash `json:"accessHash"`
}

// GetFileURL returns URL for corresponding descriptor.
func (c *Client) GetFileURL(descriptor *FileDescriptor) string {
	return fmt.Sprintf("%s/files/download/%d/%s", c.url, descriptor.Id, descriptor.AccessHash)
}

// UploadFile uploads file from disk to server and returns corresponding descriptor.
func (c *Client) UploadFile(filename string, reader io.Reader) (*FileDescriptor, error) {
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

	req, err := http.NewRequest("POST", c.url+"/files/upload", body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.http.Do(req)
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
