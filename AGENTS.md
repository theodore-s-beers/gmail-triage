# AGENTS.md - Development Guidelines

## Build/Test Commands

- `go mod tidy` - Update dependencies
- `go build` - Build the binary
- `./gmail-triage` - Run the CLI tool
- `go run .` - Build and run in one command
- `go test ./...` - Run all tests (if any exist)

## Code Style Guidelines

### Imports

- Standard library imports first, then third-party, then local packages
- Group imports with blank lines between groups
- Use full package paths (e.g., `golang.org/x/oauth2`)

### Naming Conventions

- Use camelCase for variables and functions (e.g., `searchString`,
  `getUnreadMessages`)
- Use PascalCase for exported types and functions (e.g., `EmailMessage`,
  `GmailService`)
- Also use PascalCase for constants and enum-like values when exported (e.g.,
  `ActionMarkRead`); otherwise camelCase

### Error Handling

- Return errors as the last return value
- Use `fmt.Errorf()` for wrapping errors with context
- Use `log.Fatalf()` for fatal errors in main functions
- Handle errors immediately after function calls

### Types and Structs

- Define custom types for clarity (e.g., `EmailAction int`)
- Use struct embedding when appropriate
- Keep struct fields organized and well-documented
