package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/johnson7543/chatsheet-assessment/internal/config"
)

// UnipileService handles interactions with the Unipile API
type UnipileService struct {
	apiKey string
	apiURL string
	client *http.Client
}

// NewUnipileService creates a new Unipile service
func NewUnipileService() *UnipileService {
	return &UnipileService{
		apiKey: config.App.UnipileAPIKey,
		apiURL: config.App.Unipile.APIURL,
		client: &http.Client{},
	}
}

// ConnectRequest represents a request to connect an account
type ConnectRequest struct {
	Provider string                 `json:"provider"`
	Cookie   string                 `json:"cookie,omitempty"`
	Username string                 `json:"username,omitempty"`
	Password string                 `json:"password,omitempty"`
	Type     string                 `json:"type,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ConnectResponse represents a response from Unipile
type ConnectResponse struct {
	AccountID   string `json:"account_id"`
	Provider    string `json:"provider"`
	Name        string `json:"name,omitempty"`
	Username    string `json:"username,omitempty"`
	Status      string `json:"status,omitempty"`
	Message     string `json:"message,omitempty"`
	Error       string `json:"error,omitempty"`
	Description string `json:"description,omitempty"`
}

// ConnectAccount connects an account via Unipile API
func (s *UnipileService) ConnectAccount(req ConnectRequest) (*ConnectResponse, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("Unipile API key is not configured")
	}

	// Prepare request body
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/accounts", s.apiURL)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-KEY", s.apiKey)

	// Make the request
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call Unipile API: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Parse response
	var unipileResp ConnectResponse
	if err := json.Unmarshal(body, &unipileResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// Check for errors
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		errorMsg := unipileResp.Error
		if errorMsg == "" {
			errorMsg = unipileResp.Message
		}
		if errorMsg == "" {
			errorMsg = "Unknown error from Unipile API"
		}
		return nil, fmt.Errorf("Unipile API error: %s", errorMsg)
	}

	// Validate response
	if unipileResp.AccountID == "" {
		return nil, fmt.Errorf("invalid response from Unipile API: missing account_id")
	}

	return &unipileResp, nil
}

// ConnectWithCookie connects a LinkedIn account using cookie authentication
func (s *UnipileService) ConnectWithCookie(cookie string) (accountID, accountName string, err error) {
	req := ConnectRequest{
		Provider: "linkedin",
		Cookie:   cookie,
		Type:     "cookie",
	}

	resp, err := s.ConnectAccount(req)
	if err != nil {
		return "", "", err
	}

	name := resp.Name
	if name == "" {
		name = resp.Username
	}

	return resp.AccountID, name, nil
}

// ConnectWithCredentials connects a LinkedIn account using username/password
func (s *UnipileService) ConnectWithCredentials(username, password string) (accountID, accountName string, err error) {
	req := ConnectRequest{
		Provider: "linkedin",
		Username: username,
		Password: password,
		Type:     "credentials",
	}

	resp, err := s.ConnectAccount(req)
	if err != nil {
		return "", "", err
	}

	name := resp.Name
	if name == "" {
		name = resp.Username
	}

	return resp.AccountID, name, nil
}
