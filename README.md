# Gmail Triage CLI

This is a command-line tool meant for quick triage of unread Gmail messages.

## Setup

### Google Cloud setup

1. Go to the [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the Gmail API:
   - Go to "APIs & Services" > "Library"
   - Search for "Gmail API" and enable it
4. Generate credentials:
   - Go to "APIs & Services" > "Credentials"
   - Click "Create Credentials" > "OAuth client ID"
   - Choose "Desktop application" and enter a name for your client
   - Download the JSON file; save as `credentials.json` in the project root
5. Add yourself as a user
   - Go to "APIs & Services" > "OAuth consent screen"
   - Complete whatever steps Google requires (e.g., "Branding")
   - Add your own Gmail address under "Audience" > "Test users"

_Parts of this process can be confusing; but you only need to do it once._

### Build and run

```sh
go mod tidy
go build
./gmail-triage
```

### First run

On first run, the app will:

1. Generate a link for Gmail authorization
2. Save your OAuth token to `token.json`
3. Start showing your unread messages

## Usage

The CLI will show each unread message with:

- Sender
- Subject line
- Date
- Preview snippet

For each message, choose:

- `r`: Mark as read
- `t`: Move to trash
- `s`: Mark as spam
- `p`: Pass (skip, do nothing)
- `q`: Quit the program

## Security notes

- Keep `credentials.json` and `token.json` private
- Ensure these files are in `.gitignore` if using version control
- The OAuth token will refresh automatically when needed

## TODO

1. **Test OAuth flow** – Make sure authentication works end-to-end
2. **Error handling** – Add better error handling throughout
3. **Rate limiting** – Add delays between API calls if needed
4. **Batch operations** – Consider batching actions for better performance
5. **Configuration** – Add config file for preferences (max messages, etc.)
6. **Enhanced display** – Maybe add colors or better formatting
7. **Undo functionality** – Cache recent actions in case of mistakes

## API limits

The Gmail API allows:

- 250 quota units per user per second
- Most operations cost 1–5 units
- Should be fine for personal use

The current implementation fetches, by default, up to 50 unread messages at a
time.
