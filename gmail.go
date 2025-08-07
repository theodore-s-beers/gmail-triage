package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strings"
	"unicode"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// Wraps Gmail API service
type GmailService struct {
	service *gmail.Service
}

// Simplified email message data for our CLI
type EmailMessage struct {
	ID      string
	From    string
	Subject string
	Snippet string
}

type EmailAction int

const (
	ActionPass EmailAction = iota
	ActionMarkRead
	ActionTrash
	ActionSpam
)

func initGmailService(ctx context.Context) (*GmailService, error) {
	credentialsPath := "credentials.json"

	b, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailModifyScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Gmail client: %v", err)
	}

	return &GmailService{service: srv}, nil
}

func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"

	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}

	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link, then paste the "+
		"authorization code below: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		fmt.Printf("Unable to read authorization code: %v", err)
		os.Exit(1)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		fmt.Printf("Unable to retrieve token: %v", err)
		os.Exit(1)
	}

	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)

	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		fmt.Printf("Unable to cache OAuth token: %v", err)
		return
	}
	defer f.Close()

	json.NewEncoder(f).Encode(token)
}

func getUnreadMessages(service *GmailService, searchString string) ([]*EmailMessage, error) {
	user := "me"

	// Build query: start with unread, add search string if provided
	query := "is:unread"

	searchString = strings.TrimSpace(searchString)
	if searchString != "" {
		query = fmt.Sprintf("is:unread %s", searchString)
	}

	// Returns a maximum of 100 messages by default
	req := service.service.Users.Messages.List(user).Q(query)
	resp, err := req.Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve messages: %v", err)
	}

	var messages []*EmailMessage
	for _, m := range resp.Messages {
		msg, err := service.service.Users.Messages.Get(user, m.Id).Do()
		if err != nil {
			continue // Skip irretrievable messages
		}

		email := &EmailMessage{
			ID:      msg.Id,
			Snippet: cleanSnippet(msg.Snippet),
		}

		for _, header := range msg.Payload.Headers {
			switch header.Name {
			case "From":
				email.From = header.Value
			case "Subject":
				email.Subject = header.Value
			}
		}

		messages = append(messages, email)
	}

	return messages, nil
}

func (s *GmailService) PerformAction(messageID string, action EmailAction) error {
	user := "me"

	switch action {

	case ActionMarkRead:
		req := &gmail.ModifyMessageRequest{
			RemoveLabelIds: []string{"UNREAD"},
		}
		_, err := s.service.Users.Messages.Modify(user, messageID, req).Do()
		return err

	case ActionTrash:
		_, err := s.service.Users.Messages.Trash(user, messageID).Do()
		return err

	case ActionSpam:
		req := &gmail.ModifyMessageRequest{
			AddLabelIds: []string{"SPAM"},
		}
		_, err := s.service.Users.Messages.Modify(user, messageID, req).Do()
		return err

	case ActionPass:
		return nil

	default:
		return fmt.Errorf("unknown action: %d", action)

	}
}

var badChars = []rune{
	'\u034F', // COMBINING GRAPHEME JOINER
}

func cleanSnippet(snippet string) string {
	decoded := html.UnescapeString(snippet) // Unescape HTML entities

	cleaned := strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) && !slices.Contains(badChars, r) {
			if unicode.IsSpace(r) {
				return ' '
			}

			return r
		}

		return -1 // Remove non-printable characters
	}, decoded)

	re := regexp.MustCompile(`\s+`)
	reduced := re.ReplaceAllString(cleaned, " ")

	return strings.TrimSpace(reduced)
}
