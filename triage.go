package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func startTriage(service *GmailService, messages []*EmailMessage) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\n📧 Found %d unread messages\n", len(messages))
	fmt.Println("Commands: (r)ead, (t)rash, (s)pam, (p)ass, (q)uit")
	fmt.Println(strings.Repeat("-", 60))

	for i, msg := range messages {
		if err := displayMessage(i+1, len(messages), msg); err != nil {
			return err
		}

		action, quit, err := promptForAction(reader)
		if err != nil {
			return err
		}

		if quit {
			fmt.Println("Goodbye!")
			break
		}

		if action != ActionPass {
			if err := service.PerformAction(msg.ID, action); err != nil {
				fmt.Printf("❌ Error performing action: %v\n", err)
				continue
			}

			fmt.Printf("✅ %s\n", getActionDescription(action))
		}

		fmt.Println()
	}

	return nil
}

func displayMessage(current, total int, msg *EmailMessage) error {
	fmt.Printf("\n[%d/%d] ", current, total)

	from := msg.From
	if len(from) > 50 {
		from = from[:47] + "..."
	}

	subject := msg.Subject
	if subject == "" {
		subject = "(no subject)"
	}
	if len(subject) > 60 {
		subject = subject[:57] + "..."
	}

	fmt.Printf("From: %s\n", from)
	fmt.Printf("Subject: %s\n", subject)

	if msg.Snippet != "" {
		snippet := msg.Snippet
		if len(snippet) > 80 {
			snippet = snippet[:77] + "..."
		}
		fmt.Printf("Preview: %s\n", snippet)
	}

	return nil
}

func promptForAction(reader *bufio.Reader) (EmailAction, bool, error) {
	for {
		fmt.Print("Action [r/d/s/p/q]: ")

		input, err := reader.ReadString('\n')
		if err != nil {
			return ActionPass, false, err
		}

		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "r", "read":
			return ActionMarkRead, false, nil
		case "t", "trash":
			return ActionTrash, false, nil
		case "s", "spam":
			return ActionSpam, false, nil
		case "p", "pass":
			return ActionPass, false, nil
		case "q", "quit":
			return ActionPass, true, nil
		default:
			fmt.Println("Invalid option. Use: (r)ead, (t)rash, (s)pam, (p)ass, (q)uit")
			continue
		}
	}
}

func getActionDescription(action EmailAction) string {
	switch action {
	case ActionMarkRead:
		return "Marked as read"
	case ActionTrash:
		return "Moved to trash"
	case ActionSpam:
		return "Marked as spam"
	case ActionPass:
		return "Passed"
	default:
		return "Unknown action"
	}
}
