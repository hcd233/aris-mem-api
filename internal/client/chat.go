// Package client provides chat functionality
package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/adk"
	"github.com/hcd233/aris-mem-api/internal/config"
	"github.com/hcd233/aris-mem-api/internal/dto"
	"github.com/samber/lo"
)

// ChatHandler handles chat interactions
type ChatHandler struct {
	apiClient *APIClient
	scanner   *bufio.Scanner
	token     string
}

// NewChatHandler creates a new chat handler
func NewChatHandler() *ChatHandler {
	return &ChatHandler{
		apiClient: NewAPIClient(),
		scanner:   bufio.NewScanner(os.Stdin),
	}
}

// Execute starts the chat session
func (h *ChatHandler) Execute() error {
	// Step 1: Check login status
	tokenData, err := LoadToken()
	if err != nil || tokenData == nil {
		errorColor.Println("\n✗ Not logged in")
		infoColor.Println("  Please run `client login` command first")
		return fmt.Errorf("not logged in")
	}

	h.token = tokenData.AccessToken

	// Step 2: Get and display user info
	infoColor.Println("\n🔍 Verifying credentials...")
	userResp, err := h.apiClient.GetCurrentUser(h.token)
	if err != nil {
		errorColor.Println("✗ Failed to verify credentials")
		warningColor.Println("  Please login again")
		return err
	}

	if userResp == nil || userResp.User == nil {
		errorColor.Println("✗ Invalid user data")
		return fmt.Errorf("invalid user data")
	}

	// Display welcome message
	fmt.Println()
	headerColor.Println("╔══════════════════════════════════════════════╗")
	headerColor.Println("║          AI Assistant Chat Session           ║")
	headerColor.Println("╚══════════════════════════════════════════════╝")
	fmt.Println()
	successColor.Printf("✓ Logged in as: ")
	fmt.Printf("%s (%s)\n", userResp.User.Name, userResp.User.Email)
	fmt.Println()
	infoColor.Println("💡 Type your message and press Enter to chat")
	infoColor.Println("   Type 'exit' or 'quit' to end the session")

	fmt.Println()

	// Step 3: Enter chat loop
	return h.chatLoop()
}

// chatLoop handles the main chat interaction loop
func (h *ChatHandler) chatLoop() error {
	for {
		// Display prompt
		promptColor.Print("You: ")

		// Read user input
		if !h.scanner.Scan() {
			break
		}

		input := strings.TrimSpace(h.scanner.Text())

		// Ignore empty input
		if input == "" {
			continue
		}

		// Check for exit commands
		if input == "exit" || input == "quit" {
			infoColor.Println("\n👋 Goodbye!")
			break
		}

		// Send message and display response
		if err := h.sendMessage(input); err != nil {
			errorColor.Printf("✗ Error: %v\n", err)
			fmt.Println()
		}
	}

	return nil
}

// sendMessage sends a message to the AI and streams the response
func (h *ChatHandler) sendMessage(message string) error {
	baseURL := config.ServerEndpoint
	url := fmt.Sprintf("%s/api/v1/agent/chat", baseURL)

	// Prepare multipart form body
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err := writer.WriteField("content", message); err != nil {
		return err
	}

	contentType := writer.FormDataContentType()

	if err := writer.Close(); err != nil {
		return err
	}

	// Create request
	req, err := http.NewRequest(http.MethodPost, url, &body)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", h.token))
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "text/event-stream")

	// Send request
	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Display AI label
	fmt.Println()
	commandColor.Print("AI: ")

	// Stream and render SSE response
	return h.renderSSEStream(resp.Body)
}

// renderSSEStream parses and renders SSE events
func (h *ChatHandler) renderSSEStream(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines and non-data lines
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		// Extract JSON data
		data := strings.TrimPrefix(line, "data: ")

		var event dto.SSEResponse
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			continue
		}

		// Handle different event types
		switch event.DataType {
		case "heartbeat":
			// Ignore heartbeat events
			continue

		case "message":
			var msg adk.Message
			if err := sonic.Unmarshal([]byte(lo.Must1(sonic.Marshal(event.Data))), &msg); err != nil {
				continue
			}

			// Skip tool role messages (tool output)
			if msg.Role == "tool" {
				toolCallColor.Println("🔧 Tool Call Done!")
				fmt.Println()
				continue
			}

			// Handle assistant messages
			if msg.Role == "assistant" {
				// Handle tool calls
				if len(msg.ToolCalls) > 0 {
					for _, tc := range msg.ToolCalls {
						if tc.Function.Name != "" {
							toolCallColor.Printf("\n\n🔧 Tool Call: %s\n", tc.Function.Name)
						}
					}
				} else if msg.Content != "" {
					// Regular content message
					fmt.Print(msg.Content)
				}

				// Handle completion
				if msg.ResponseMeta.FinishReason == "stop" {
					fmt.Println()
				}
			}
		}
	}

	fmt.Println()
	return scanner.Err()
}
