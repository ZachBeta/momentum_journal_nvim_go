# Momentum Journal

A minimalist tool for following "The Artist's Way" journaling practice. It provides a distraction-free writing environment with AI support to help maintain writing flow.

## Features

- Two-pane terminal interface:
  - Writing pane with vim-like navigation and editing
  - Conversation pane with AI agent to facilitate reflection
- Local or cloud-based LLM integration via Ollama or OpenRouter
- Markdown file storage with metadata tracking
- Progress tracking toward a 750-word goal

## Building & Running

### Prerequisites

- Go 1.20 or later
- [Ollama](https://ollama.ai/) (optional, for local LLM support)

### Installation

```bash
# Clone the repository
git clone https://github.com/ZachBeta/momentum_journal_nvim_go.git
cd momentum_journal_nvim_go

# Build the project
go build -o bin/momentum ./cmd/momentum/...

# Run directly without building
go run ./cmd/momentum/... [command]
```

### Usage

```bash
# Create a new journal entry and open the TUI
momentum new

# List existing journal entries
momentum list

# Show help
momentum --help
```

## Key Bindings

- **Writing Pane:**
  - `i` - Enter Insert mode
  - `Esc` - Return to Normal mode
  - Vim-like movement: `h`, `j`, `k`, `l`, etc.

- **Navigation:**
  - `Tab` - Switch between writing and conversation panes
  - `q` or `Ctrl+C` - Quit the application

## Project Status

This is a work in progress. Currently implementing Phase 2 (Terminal UI) of the [implementation plan](bubbletea_plan.md).
