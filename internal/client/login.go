// Package client provides CLI login functionality
package client

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/samber/lo"
)

const (
	maxRetries = 3
)

var (
	// Color definitions
	successColor = color.New(color.FgGreen, color.Bold)
	errorColor   = color.New(color.FgRed, color.Bold)
	infoColor    = color.New(color.FgCyan)
	warningColor = color.New(color.FgYellow)
	promptColor  = color.New(color.FgBlue, color.Bold)
	headerColor  = color.New(color.FgMagenta, color.Bold)
	commandColor = color.New(color.FgWhite, color.Bold)
	commandPromptColor = color.New(color.FgWhite, color.Bold)
	toolCallColor = color.New(color.FgYellow, color.Bold)
)

// LoginHandler 登录处理器
type LoginHandler struct {
	apiClient *APIClient
	scanner   *bufio.Scanner
}

// NewLoginHandler 创建登录处理器
//
//	return *LoginHandler
//	author centonhuang
//	update 2025-11-11 20:30:00
func NewLoginHandler() *LoginHandler {
	return &LoginHandler{
		apiClient: NewAPIClient(),
		scanner:   bufio.NewScanner(os.Stdin),
	}
}

// Execute executes the login flow
func (h *LoginHandler) Execute() error {
	// Step 1: Check local token
	tokenData, err := LoadToken()
	if err == nil && tokenData != nil {
		infoColor.Println("\n🔍 Checking existing credentials...")

		userResp, err := h.apiClient.GetCurrentUser(tokenData.AccessToken)
		if err == nil && userResp != nil && userResp.User != nil {
			// Token is valid, ask if user wants to re-login
			fmt.Println()
			successColor.Printf("✓ Already logged in as: ")
			fmt.Printf("%s (%s)\n", userResp.User.Name, userResp.User.Email)
			promptColor.Print("\n→ Re-login? (y/n): ")

			if !h.readYesNo() {
				infoColor.Println("✓ Keeping current session")
				return nil
			}
		} else {
			warningColor.Println("\n⚠ Token expired or invalid, please re-login")
		}
	} else {
		headerColor.Println("╔══════════════════════════════════════════════╗")
		headerColor.Println("║      Welcome to Aris Memory API Client       ║")
		headerColor.Println("╚══════════════════════════════════════════════╝")
	}

	// Step 2: Select login platform
	provider, err := h.selectProvider()
	if err != nil {
		return err
	}

	// Step 3: Get OAuth2 authorization URL
	loginResp, err := h.apiClient.OAuth2Login(provider)
	if err != nil {
		return fmt.Errorf("OAuth2 login request failed: %v", err)
	}

	fmt.Println()
	headerColor.Println("╔══════════════════════════════════════════════╗")
	headerColor.Println("║            Authorization Required            ║")
	headerColor.Println("╚══════════════════════════════════════════════╝")
	fmt.Println()
	infoColor.Println("📋 Please open the following URL in your browser:")
	fmt.Println()
	fmt.Println("  " + loginResp.RedirectURL + "  ")
	fmt.Println()
	warningColor.Println("⚠ After authorization, you will be redirected to a URL with parameters")
	fmt.Println("  Example: http://example.com/callback?state=xxx&code=yyy")
	fmt.Println()

	// Step 4: Get user input for state and code
	code, state, err := h.readCodeAndState()
	if err != nil {
		return err
	}

	// Step 5: Call OAuth2 callback
	infoColor.Println("\n🔄 Verifying authorization code...")
	callbackResp, err := h.apiClient.OAuth2Callback(provider, code, state)
	if err != nil {
		return fmt.Errorf("OAuth2 callback failed: %v", err)
	}

	// Step 6: Save token locally
	newTokenData := &TokenData{
		AccessToken:  callbackResp.AccessToken,
		RefreshToken: callbackResp.RefreshToken,
	}

	if err := SaveToken(newTokenData); err != nil {
		return fmt.Errorf("failed to save token: %v", err)
	}

	// Get user info
	fmt.Println()
	userResp, err := h.apiClient.GetCurrentUser(newTokenData.AccessToken)
	if err != nil {
		successColor.Println("✓ Login successful!")
	} else if userResp != nil && userResp.User != nil {
		successColor.Println("✓ Login successful!")
		fmt.Printf("  User: %s (%s)\n", userResp.User.Name, userResp.User.Email)
	}

	tokenFilePath := lo.Must1(getTokenFilePath())
	infoColor.Printf("  Token saved to: %s\n", tokenFilePath)

	return nil
}

// selectProvider prompts user to select login platform
func (h *LoginHandler) selectProvider() (string, error) {
	validProviders := map[string]bool{
		"github": true,
		"google": true,
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		promptColor.Print("\n→ Select login platform (github/google): ")

		if !h.scanner.Scan() {
			return "", fmt.Errorf("failed to read input")
		}

		provider := strings.TrimSpace(strings.ToLower(h.scanner.Text()))

		if validProviders[provider] {
			successColor.Printf("✓ Selected: %s\n", provider)
			return provider, nil
		}

		errorColor.Printf("✗ Invalid platform: %s\n", provider)
		warningColor.Println("  Please enter `github` or `google`")

		if attempt == maxRetries {
			return "", fmt.Errorf("maximum retries exceeded")
		}
	}

	return "", fmt.Errorf("failed to select platform")
}

// readCodeAndState reads authorization code and state from user
func (h *LoginHandler) readCodeAndState() (string, string, error) {
	for attempt := 1; attempt <= maxRetries; attempt++ {
		promptColor.Print("\n→ Enter `state` parameter: ")
		if !h.scanner.Scan() {
			return "", "", fmt.Errorf("failed to read input")
		}
		state := strings.TrimSpace(h.scanner.Text())

		promptColor.Print("→ Enter `code` parameter: ")
		if !h.scanner.Scan() {
			return "", "", fmt.Errorf("failed to read input")
		}
		code := strings.TrimSpace(h.scanner.Text())

		if state == "" || code == "" {
			errorColor.Println("✗ State and code cannot be empty")
			warningColor.Println("  Please try again")

			if attempt == maxRetries {
				return "", "", fmt.Errorf("maximum retries exceeded")
			}
			continue
		}

		successColor.Println("✓ Parameters received")
		return code, state, nil
	}

	return "", "", fmt.Errorf("failed to read code and state")
}

// readYesNo reads yes/no choice from user
func (h *LoginHandler) readYesNo() bool {
	if !h.scanner.Scan() {
		return false
	}

	answer := strings.TrimSpace(strings.ToLower(h.scanner.Text()))
	return answer == "y" || answer == "yes"
}
