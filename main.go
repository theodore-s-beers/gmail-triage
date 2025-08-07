package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gmail-triage [flags]",
	Short: "A simple CLI for triaging unread Gmail messages",
	Long:  `This tool allows quick triaging of unread Gmail messages from the command line`,
	Run:   runTriage,
}

var searchString string

func init() {
	rootCmd.Flags().StringVarP(&searchString, "search", "s", "", "search string to filter messages (combined with unread)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runTriage(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	fmt.Println("⏳ Initializing Gmail service...")
	service, err := initGmailService(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize Gmail service: %v", err)
	}

	if searchString != "" {
		fmt.Printf("⏳ Fetching unread messages matching '%s'...\n", searchString)
	} else {
		fmt.Println("⏳ Fetching unread messages...")
	}

	messages, err := getUnreadMessages(service, searchString)
	if err != nil {
		log.Fatalf("Failed to get unread messages: %v", err)
	}

	if len(messages) == 0 {
		fmt.Println("🎉 No unread messages!")
		return
	}

	if err := startTriage(service, messages); err != nil {
		log.Fatalf("Error during triage: %v", err)
	}
}
