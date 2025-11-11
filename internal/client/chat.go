// Package client provides chat functionality
package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/hcd233/aris-mem-api/internal/config"
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
		infoColor.Println("  Please run 'login' command first")
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

	// Prepare request body
	reqBody := map[string]interface{}{
		"message": message,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// Create request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", h.token))
	req.Header.Set("Content-Type", "application/json")
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
	color.New(color.FgCyan, color.Bold).Print("AI: ")

	// Stream and render SSE response
	return h.renderSSEStream(resp.Body)
}

// renderSSEStream parses and renders SSE events
func (h *ChatHandler) renderSSEStream(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)

	// For accumulating tool call arguments
	toolCallBuffer := make(map[int]*ToolCallInfo)
	lastWasToolCall := false

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines and non-data lines
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		// Extract JSON data
		data := strings.TrimPrefix(line, "data: ")

		var event SSEEvent
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			continue
		}

		// Handle different event types
		switch event.DataType {
		case "heartbeat":
			// Ignore heartbeat events
			continue

		case "message":
			var msg ChatMessage
			if err := json.Unmarshal([]byte(event.Data), &msg); err != nil {
				continue
			}

			// Skip tool role messages (tool output)
			if msg.Role == "tool" {
				continue
			}

			// Handle assistant messages
			if msg.Role == "assistant" {
				// Handle tool calls
				if len(msg.ToolCalls) > 0 {
					for _, tc := range msg.ToolCalls {
						if tc.ID != "" && tc.Function.Name != "" {
							// Start of a new tool call
							if !lastWasToolCall {
								fmt.Println()
								lastWasToolCall = true
							}
							toolCallBuffer[tc.Index] = &ToolCallInfo{
								ID:        tc.ID,
								Name:      tc.Function.Name,
								Arguments: tc.Function.Arguments,
							}
						} else if tc.Function.Arguments != "" {
							// Accumulate arguments
							if info, exists := toolCallBuffer[tc.Index]; exists {
								info.Arguments += tc.Function.Arguments
								toolCallBuffer[tc.Index] = info
							}
						}
					}

					// Check if tool call is complete (has finish_reason)
					if msg.ResponseMeta.FinishReason == "tool_calls" {
						h.renderToolCalls(toolCallBuffer)
						toolCallBuffer = make(map[int]*ToolCallInfo)
						lastWasToolCall = false
					}
				} else if msg.Content != "" {
					// Regular content message
					if lastWasToolCall {
						fmt.Println()
						lastWasToolCall = false
					}
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

// renderToolCalls renders tool calls in a formatted way
func (h *ChatHandler) renderToolCalls(toolCalls map[int]*ToolCallInfo) {
	for _, info := range toolCalls {
		if info == nil || info.Name == "" {
			continue
		}

		// Parse and format arguments
		var args map[string]interface{}
		if err := json.Unmarshal([]byte(info.Arguments), &args); err == nil {
			formatted, _ := json.MarshalIndent(args, "  ", "  ")

			fmt.Println()
			color.New(color.FgYellow, color.Bold).Printf("🔧 Tool Call: %s\n", info.Name)
			color.New(color.FgYellow).Println("  Parameters:")
			fmt.Printf("  %s\n", string(formatted))
		}
	}
}

// SSEEvent represents an SSE event
type SSEEvent struct {
	DataType string `json:"dataType"`
	Data     string `json:"data"`
}

// ChatMessage represents a chat message
type ChatMessage struct {
	Role         string       `json:"role"`
	Content      string       `json:"content"`
	ToolCalls    []ToolCall   `json:"tool_calls"`
	ResponseMeta ResponseMeta `json:"response_meta"`
}

// ToolCall represents a tool call
type ToolCall struct {
	Index    int          `json:"index"`
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

// FunctionCall represents a function call
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ResponseMeta contains response metadata
type ResponseMeta struct {
	FinishReason string `json:"finish_reason"`
	Usage        *Usage `json:"usage"`
}

// Usage contains token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ToolCallInfo accumulates tool call information
type ToolCallInfo struct {
	ID        string
	Name      string
	Arguments string
}
