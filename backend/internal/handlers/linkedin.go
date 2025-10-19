package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johnson7543/chatsheet-assessment/internal/config"
	"github.com/johnson7543/chatsheet-assessment/internal/database"
	"github.com/johnson7543/chatsheet-assessment/internal/models"
)

// UnipileConnectRequest represents the request to Unipile API for account connection
type UnipileConnectRequest struct {
	Provider    string `json:"provider"`
	AccessToken string `json:"access_token,omitempty"` // For cookie auth (li_at cookie)
	Username    string `json:"username,omitempty"`     // For credentials auth
	Password    string `json:"password,omitempty"`     // For credentials auth
}

// UnipileConnectResponse represents the response from Unipile API
type UnipileConnectResponse struct {
	// Success fields
	AccountID string `json:"account_id"`
	Provider  string `json:"provider"`
	Name      string `json:"name,omitempty"`
	Username  string `json:"username,omitempty"`
	Status    any    `json:"status,omitempty"` // Can be string or number

	// Error fields (multiple formats supported)
	Error       string `json:"error,omitempty"`       // Generic error
	Message     string `json:"message,omitempty"`     // Generic message
	Description string `json:"description,omitempty"` // Error description

	// Unipile specific error format
	Type   string `json:"type,omitempty"`   // e.g., "errors/invalid_credentials"
	Title  string `json:"title,omitempty"`  // e.g., "Invalid credentials"
	Detail string `json:"detail,omitempty"` // e.g., "The provided credentials are invalid."
}

// ConnectLinkedInWithCookie handles LinkedIn connection using cookie authentication
func ConnectLinkedInWithCookie(c *gin.Context) {
	var req models.LinkedInCookieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	userID := c.GetUint("user_id")

	// Call Unipile API with cookie authentication
	unipileReq := UnipileConnectRequest{
		Provider:    "LINKEDIN",
		AccessToken: req.Cookie, // li_at cookie value
	}

	accountID, accountName, err := callUnipileAPI(unipileReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Save to database
	linkedAccount := models.LinkedAccount{
		UserID:      userID,
		Provider:    "linkedin",
		AccountID:   accountID,
		AccountName: accountName,
	}

	if err := database.DB.Create(&linkedAccount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to save account"})
		return
	}

	response := models.LinkedInConnectResponse{
		Message:   "LinkedIn account connected successfully",
		AccountID: accountID,
		Account:   linkedAccount,
	}

	c.JSON(http.StatusOK, response)
}

// ConnectLinkedInWithCredentials handles LinkedIn connection using username/password
func ConnectLinkedInWithCredentials(c *gin.Context) {
	log.Println("=== LinkedIn Connect with Credentials START ===")

	var req models.LinkedInCredentialsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("ERROR: Failed to bind JSON request: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	log.Printf("User ID: %d", userID)
	log.Printf("Username provided: %s", req.Username)
	log.Printf("Password length: %d", len(req.Password))

	// Call Unipile API with credentials
	unipileReq := UnipileConnectRequest{
		Provider: "LINKEDIN",
		Username: req.Username,
		Password: req.Password,
	}

	log.Println("Calling Unipile API...")
	accountID, accountName, err := callUnipileAPI(unipileReq)
	if err != nil {
		log.Printf("ERROR: Unipile API call failed: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	log.Printf("SUCCESS: Received account_id: %s, account_name: %s", accountID, accountName)

	// Save to database
	log.Println("Saving to database...")
	linkedAccount := models.LinkedAccount{
		UserID:      userID,
		Provider:    "linkedin",
		AccountID:   accountID,
		AccountName: accountName,
	}

	if err := database.DB.Create(&linkedAccount).Error; err != nil {
		log.Printf("ERROR: Failed to save to database: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to save account"})
		return
	}
	log.Printf("Database saved successfully, ID: %d", linkedAccount.ID)

	response := models.LinkedInConnectResponse{
		Message:   "LinkedIn account connected successfully",
		AccountID: accountID,
		Account:   linkedAccount,
	}

	log.Println("=== LinkedIn Connect with Credentials END (SUCCESS) ===")
	c.JSON(http.StatusOK, response)
}

// callUnipileAPI makes the actual API call to Unipile
func callUnipileAPI(req UnipileConnectRequest) (accountID, accountName string, err error) {
	log.Println("--- callUnipileAPI START ---")

	// Check if API key is configured
	if config.App.UnipileAPIKey == "" {
		log.Println("ERROR: Unipile API key is not configured")
		return "", "", fmt.Errorf("config error: Unipile API key is not configured")
	}
	log.Printf("Unipile API URL: %s", config.App.Unipile.APIURL)

	// Prepare request body
	jsonData, err := json.Marshal(req)
	if err != nil {
		log.Printf("ERROR: Failed to marshal request: %v", err)
		return "", "", fmt.Errorf("failed to marshal request: %v", err)
	}
	log.Printf("Request payload: %s", string(jsonData))

	// Create HTTP request to Unipile API
	url := fmt.Sprintf("%s/accounts", config.App.Unipile.APIURL)
	log.Printf("Making POST request to: %s", url)

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("ERROR: Failed to create HTTP request: %v", err)
		return "", "", fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-KEY", config.App.UnipileAPIKey)
	log.Printf("Request headers set: Content-Type=application/json, X-API-KEY=***%s", config.App.UnipileAPIKey[len(config.App.UnipileAPIKey)-4:])

	// Make the request
	log.Println("Sending request to Unipile...")
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		log.Printf("ERROR: HTTP request failed: %v", err)
		return "", "", fmt.Errorf("failed to call Unipile API: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("Response Status Code: %d", resp.StatusCode)
	log.Printf("Response Headers: %v", resp.Header)

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: Failed to read response body: %v", err)
		return "", "", fmt.Errorf("failed to read response: %v", err)
	}
	log.Printf("Response Body: %s", string(body))

	// Parse response
	var unipileResp UnipileConnectResponse
	if err := json.Unmarshal(body, &unipileResp); err != nil {
		log.Printf("ERROR: Failed to parse JSON response: %v", err)
		log.Printf("Raw response body: %s", string(body))
		return "", "", fmt.Errorf("failed to parse response: %v", err)
	}
	log.Printf("Parsed response: AccountID=%s, Provider=%s, Name=%s, Title=%s, Detail=%s, Error=%s",
		unipileResp.AccountID, unipileResp.Provider, unipileResp.Name, unipileResp.Title, unipileResp.Detail, unipileResp.Error)

	// Check for errors in response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		// Build comprehensive error message from various possible error fields
		errorMsg := getUnipileErrorMessage(unipileResp)
		log.Printf("ERROR: Unipile API returned error status: %d, message: %s", resp.StatusCode, errorMsg)
		return "", "", fmt.Errorf("%s", errorMsg)
	}

	// Validate response
	if unipileResp.AccountID == "" {
		log.Println("ERROR: Response missing account_id field")
		return "", "", fmt.Errorf("invalid response from Unipile API: missing account_id")
	}

	// Return the account ID and name
	name := unipileResp.Name
	if name == "" {
		name = unipileResp.Username
	}

	log.Printf("SUCCESS: Returning account_id=%s, name=%s", unipileResp.AccountID, name)
	log.Println("--- callUnipileAPI END ---")
	return unipileResp.AccountID, name, nil
}

// getUnipileErrorMessage extracts the best error message from Unipile response
func getUnipileErrorMessage(resp UnipileConnectResponse) string {
	// Priority order for error messages:
	// 1. Detail (most descriptive)
	// 2. Title (error title)
	// 3. Error (generic error field)
	// 4. Description (alternative description)
	// 5. Message (generic message)
	// 6. Type (error type)

	if resp.Detail != "" {
		// If we have both title and detail, combine them
		if resp.Title != "" {
			return fmt.Sprintf("%s: %s", resp.Title, resp.Detail)
		}
		return resp.Detail
	}

	if resp.Title != "" {
		return resp.Title
	}

	if resp.Error != "" {
		if resp.Description != "" {
			return fmt.Sprintf("%s: %s", resp.Error, resp.Description)
		}
		return resp.Error
	}

	if resp.Description != "" {
		return resp.Description
	}

	if resp.Message != "" {
		return resp.Message
	}

	if resp.Type != "" {
		return fmt.Sprintf("Error type: %s", resp.Type)
	}

	return "Unknown error from Unipile API"
}
