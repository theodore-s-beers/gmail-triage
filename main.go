package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gmail-triage [flags]",
	Short: "A simple CLI for triaging unread Gmail messages",
	Long:  `This tool allows quick triaging of unread Gmail messages from the command line`,
	Run:   runTriage,
}

var (
	maxAge       int
	maxResults   int
	searchString string
)

func init() {
	rootCmd.Flags().IntVarP(&maxAge, "max-age", "a", -1, "maximum age of messages in days")
	rootCmd.Flags().IntVarP(&maxResults, "max-results", "m", 50, "maximum number of messages to fetch")
	rootCmd.Flags().StringVarP(&searchString, "search", "s", "", "search string to filter unread messages")
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

	maxResults = max(1, min(maxResults, 100))
	fmt.Printf("\n- Will search for up to %d unread messages\n", maxResults)

	if 0 <= maxAge && maxAge <= 365 {
		if maxAge == 1 {
			fmt.Printf("- With a maximum age of %d day\n", maxAge)
		} else {
			fmt.Printf("- With a maximum age of %d days\n", maxAge)
		}
	}

	searchString = strings.TrimSpace(searchString)
	if searchString != "" {
		fmt.Printf("- Matching the keyword(s) '%s'\n", searchString)
	}

	messages, err := getUnreadMessages(service, maxAge, maxResults, searchString)
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
