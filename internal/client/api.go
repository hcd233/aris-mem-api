// Package client provides API client functionality
package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/hcd233/aris-mem-api/internal/config"
	"github.com/hcd233/aris-mem-api/internal/protocol/dto"
)

// APIClient represents the HTTP API client
type APIClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewAPIClient creates a new API client instance
func NewAPIClient() *APIClient {
	baseURL := config.ServerEndpoint

	if baseURL == "" {
		color.Red("✗ Server endpoint is not configured")
		color.Yellow("  Please set SERVER_ENDPOINT environment variable")
		os.Exit(1)
	}

	return &APIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetCurrentUser retrieves current user information
func (c *APIClient) GetCurrentUser(accessToken string) (*dto.GetCurUserRsp, error) {
	url := fmt.Sprintf("%s/api/v1/user/current", c.baseURL)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized: token expired or invalid")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %d, body: %s", resp.StatusCode, string(body))
	}

	var result dto.GetCurUserRsp
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	// Check if there's an error in the response body
	if result.Error != nil {
		// Handle unauthorized error (token expired)
		if result.Error.Code == 10001 {
			return nil, fmt.Errorf("unauthorized: token expired or invalid")
		}
		return nil, fmt.Errorf("API error: %s", result.Error.Error())
	}

	return &result, nil
}

// OAuth2Login initiates OAuth2 login flow
func (c *APIClient) OAuth2Login(platform string) (*dto.LoginResp, error) {
	url := fmt.Sprintf("%s/api/v1/oauth2/login?platform=%s", c.baseURL, platform)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %d, body: %s", resp.StatusCode, string(body))
	}

	var result dto.LoginResp
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	// Check if there's an error in the response body
	if result.Error != nil {
		return nil, fmt.Errorf("API error: %s", result.Error.Error())
	}

	return &result, nil
}

// OAuth2Callback handles OAuth2 callback with authorization code
func (c *APIClient) OAuth2Callback(platform, code, state string) (*dto.CallbackRsp, error) {
	url := fmt.Sprintf("%s/api/v1/oauth2/callback?platform=%s&code=%s&state=%s", c.baseURL, platform, code, state)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %d, body: %s", resp.StatusCode, string(body))
	}

	var result dto.CallbackRsp
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	// Check if there's an error in the response body
	if result.Error != nil {
		return nil, fmt.Errorf("API error: %s", result.Error.Error())
	}

	return &result, nil
}
