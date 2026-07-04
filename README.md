# Budget CLI

A terminal-based TUI application for managing personal personal finance budgets.

## Project Overview

`budget-cli` is built with Bubble Tea and Lip Gloss to provide a responsive text user interface. The application currently includes:

- A startup loading screen
- A home screen with a centered banner
- A menu screen with keyboard navigation
- Shared terminal window sizing across submodels via a single `WindowSize` pointer

## Key Libraries

- `github.com/charmbracelet/bubbletea` - TUI framework for managing Bubble Tea models, messages, and program lifecycle
- `github.com/charmbracelet/bubbles` - reusable Bubble Tea components such as spinner widgets
- `github.com/charmbracelet/lipgloss` - terminal styling and layout

## Current Architecture

- `main.go` starts the app with `models.NewRootModel()` and runs it in an alternate screen
- `models/root.go` manages the application state machine and routes messages between views
- `models/loading.go` implements the loading spinner screen
- `models/home.go` implements the home screen and transitions into the menu
- `models/menu.go` implements the menu navigation UI
- `styles/base_styles.go` contains shared Lip Gloss styles
- `tests/` contains black-box tests exercising the public model behavior

## Running the Project

```bash
cd budget_cli
go run .
```

## Running Tests

```bash
go test ./...
```

## Next Steps

- Add persistent budget storage using a file or embedded database
- Implement budget creation and editing flows
- Add transaction tracking and monthly summaries
- Improve menu layout and key bindings for better usability
- Add more tests for navigation and rendering behavior
