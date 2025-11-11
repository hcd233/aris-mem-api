package cmd

import (
	"fmt"
	"os"

	"github.com/hcd233/aris-mem-api/internal/client"

	"github.com/spf13/cobra"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Client Command Group",
	Long:  `Client command group for authentication and interaction`,
}

var loginClientCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the service",
	Long:  `Authenticate using OAuth2 and obtain access tokens`,
	Run: func(_ *cobra.Command, _ []string) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("\n❌ An error occurred during login")
				os.Exit(1)
			}
		}()

		handler := client.NewLoginHandler()
		if err := handler.Execute(); err != nil {
			fmt.Printf("\n❌ Login failed: %v\n", err)
			os.Exit(1)
		}

		os.Exit(0)
	},
}

var chatClientCmd = &cobra.Command{
	Use:   "chat",
	Short: "Chat with AI assistant",
	Long:  `Interactive chat with Aris Mem AI`,
	Run: func(_ *cobra.Command, _ []string) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("\n❌ An error occurred during chat")
				os.Exit(1)
			}
		}()

		handler := client.NewChatHandler()
		if err := handler.Execute(); err != nil {
			os.Exit(1)
		}

		os.Exit(0)
	},
}

func init() {
	clientCmd.AddCommand(loginClientCmd)
	clientCmd.AddCommand(chatClientCmd)
	rootCmd.AddCommand(clientCmd)
}
