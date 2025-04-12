# Momentum Journal Specification

## Overview
Momentum Journal is a minimalist tool designed to help artists follow "The Artist's Way" journaling practice, specifically the Morning Pages exercise (3 pages of stream of consciousness writing). The tool combines the power of neovim for text editing with AI assistance to help maintain writing flow when the user pauses or needs inspiration.

## Core Architecture
- **Backend**: Written in Go, providing core functionality via HTTP/REST API
- **Frontend**: Neovim plugin written in Lua that communicates with the Go backend
- **LLM Integration**: Support for both local LLMs via Ollama and cloud providers via OpenRouter-compatible API

## Core Features

### Minimum Viable Product (MVP)
1. Stream of consciousness writing in a markdown buffer
2. AI assistance via keyboard shortcut when the user pauses
3. Progress tracking toward 3-page goal (750 words)
4. Autosave functionality
5. Simple CLI command to start a new session
6. Socratic questioning as default agent personality

### User Interaction Flow
1. User runs `momentum_journal` or similar command
2. Tool creates a new timestamped markdown file and opens neovim with split panes:
   - Left: Writing buffer (markdown)
   - Right: Agent conversation pane
3. User writes in the left pane (starts in insert mode)
4. When user needs assistance, they press `<leader>q`
5. Cursor moves to agent window, displays "thinking..." and diagnostics
6. Agent processes context and asks a question based on the writing
7. User can either continue the conversation or switch back to writing
8. Autosave occurs throughout the session

## Technical Specifications

### Backend (Go)

#### API Endpoints
- `POST /session/new` - Create new writing session
- `POST /session/{id}/context` - Send current writing context to be processed
- `GET /session/{id}/status` - Get session statistics (word count, etc.)
- `GET /config` - Retrieve configuration
- `POST /config` - Update configuration

#### Data Models
- **Session**: Metadata about the current writing session
- **Context**: Writing content and history sent to the LLM
- **Response**: LLM response with questions/prompts
- **Configuration**: System settings

#### File Structure
- Journal entries stored in flat structure
- Naming convention: ISO8601 timestamp (minute-level precision) + suffix
  - Example: `/journals/2023-10-05T08:30-morning-pages.md`

#### Configuration
- Stored in `~/.config/momentum_journal/config.yaml`
- Configurable options:
  - LLM provider (Ollama/OpenRouter)
  - API keys and endpoints
  - Journal storage location
  - Word count target for "3 pages" (default: 750 words)
  - Custom agent personality (path to markdown file)

### Frontend (Neovim Plugin)

#### Key Bindings
- `<leader>q` - Trigger agent assistance

#### UI Elements
- Status bar showing progress toward 3-page goal (percentage)
- Split window layout (writing pane and agent pane)
- Diagnostic information in agent pane

#### Communication
- HTTP requests to the Go backend service
- Streaming responses from LLM shown in real-time

### LLM Integration
- Default: Local Ollama instance
- Alternative: OpenRouter-compatible API
- System prompt defining Socratic questioning personality
- Context management:
  - Full session accessible to LLM
  - Priority on most recent content
  - Appropriate token limit management

## Agent Personality
- Default: Socratic questioner
  - Asks thought-provoking questions
  - Helps user explore ideas more deeply
  - Avoids leading or judgmental framing
- Custom personalities via markdown configuration file
- Easy to modify or create new personalities

## Error Handling
- Graceful degradation if LLM service is unavailable
- Local caching to prevent data loss
- Clear error messages in the agent pane
- Automatic reconnection attempts

## Testing Plan
1. **Unit Tests**:
   - Go backend components
   - Configuration management
   - File operations
   - LLM API interaction

2. **Integration Tests**:
   - Backend API functionality
   - Communication between neovim and backend
   - LLM response handling

3. **User Testing**:
   - Writing flow and experience
   - Agent interaction quality
   - Performance with different LLMs

## Development Roadmap
1. **Phase 1** - Core Functionality:
   - Basic writing interface
   - Simple agent interaction
   - Local file storage

2. **Phase 2** - Enhanced Features:
   - Improved agent personalities
   - Better progress tracking
   - Session analytics

3. **Phase 3** - Additional Frontends:
   - Web interface
   - Mobile support
   - Other editor integrations

## Implementation Notes
- Start with prototype focused on core writing and agent interaction
- Keep architecture modular to support future frontends
- Prioritize writing experience and minimal friction to start writing 