# Momentum Journal Implementation Plan

This document outlines a step-by-step approach to building the Momentum Journal application, broken down into small, iterative chunks that build upon each other. Each section includes detailed prompts for a code-generation LLM to implement the steps.

## Overall Architecture

The Momentum Journal consists of two main components:
1. A Go backend service that handles core functionality via HTTP/REST API
2. A Neovim plugin written in Lua that provides the user interface

## Implementation Phases

### Phase 1: Foundation and Core Backend
### Phase 2: Neovim Plugin Integration
### Phase 3: LLM Integration and Agent Functionality
### Phase 4: Polish and Final Features

## Detailed Implementation Steps

### Phase 1: Foundation and Core Backend

#### Step 1.1: Project Setup and Configuration Management

**Tasks:**
1. Create basic project structure
2. Implement configuration system with YAML support
3. Set up logging

**Implementation Steps:**
1. Initialize Go module and directory structure
2. Create configuration structs and YAML parsing
3. Implement configuration file handling (~/.config/momentum_journal/config.yaml)
4. Add basic logging framework

**LLM Prompt:**
```
I'm building a journaling application called "Momentum Journal" with a Go backend and Neovim frontend. Let's start with the project setup and configuration management.

Create the initial Go project structure with:
1. A main.go file with basic setup
2. A configuration package that:
   - Defines a Config struct with fields for LLM provider, API keys, journal storage location, and word count target
   - Loads/saves configuration from ~/.config/momentum_journal/config.yaml
   - Creates default configuration if none exists
3. Basic logging setup

The configuration should be in YAML format and include these fields:
- LLM provider (Ollama/OpenRouter)
- API keys and endpoints
- Journal storage location
- Word count target for "3 pages" (default: 750 words)

Please provide all necessary code files with proper error handling and comments.
```

#### Step 1.2: File Management System

**Tasks:**
1. Implement journal entry file creation and management
2. Set up automatic saving functionality
3. Create journal metadata structure

**Implementation Steps:**
1. Create journal package with entry creation/reading functions
2. Implement timestamped filename generation (ISO8601)
3. Add metadata tracking for journal entries
4. Create autosave functionality

**LLM Prompt:**
```
Building on our Momentum Journal application, let's implement the file management system for journal entries.

Create a journal package that:
1. Defines a JournalEntry struct with metadata (creation time, last modified, word count, etc.)
2. Implements functions to:
   - Create a new journal entry with ISO8601 timestamped filename (e.g., 2023-10-05T08:30-morning-pages.md)
   - Read existing journal entries
   - Update/save journal entries (for autosave functionality)
   - Track progress toward the 3-page goal (750 words by default)
3. Handles the journal storage directory (creating if it doesn't exist)

The implementation should:
- Use the configuration from the previous step to determine storage location
- Support markdown files
- Include proper error handling and logging
- Have unit tests for core functionality

This code should integrate with the configuration system we created earlier.
```

#### Step 1.3: HTTP API Server Implementation

**Tasks:**
1. Set up HTTP server with basic endpoints
2. Implement session management
3. Create API handlers for journal operations

**Implementation Steps:**
1. Create HTTP server setup with router
2. Implement session creation and management
3. Add endpoints for session context and status
4. Create configuration API endpoints

**LLM Prompt:**
```
Now let's implement the HTTP API server for our Momentum Journal application.

Create an API server that provides these endpoints:
1. POST /session/new - Creates a new writing session
2. POST /session/{id}/context - Accepts writing content to be processed
3. GET /session/{id}/status - Returns session statistics (word count, etc.)
4. GET /config - Retrieves current configuration
5. POST /config - Updates configuration

The implementation should:
1. Use a router (like gorilla/mux or chi)
2. Include session management that tracks:
   - Session ID
   - Associated journal entry
   - Session statistics (start time, current word count, etc.)
3. Have proper error handling and return appropriate HTTP status codes
4. Include middleware for logging and basic authentication if needed
5. Integrate with the configuration and journal packages from previous steps

Also implement a simple health check endpoint (GET /health) that returns the service status.

Include unit tests for each endpoint focusing on the main functionality.
```

### Phase 2: Neovim Plugin Integration

#### Step 2.1: Basic Neovim Plugin Structure

**Tasks:**
1. Create Neovim plugin skeleton
2. Set up basic configuration in Lua
3. Implement plugin initialization

**Implementation Steps:**
1. Create plugin directory structure
2. Implement basic Lua plugin modules
3. Add plugin configuration handling
4. Set up plugin commands

**LLM Prompt:**
```
Let's create the Neovim plugin structure for Momentum Journal. This plugin will communicate with our Go backend service.

Create a Neovim plugin with:
1. Standard plugin directory structure (plugin/, lua/, etc.)
2. A main plugin file that handles:
   - Plugin initialization
   - Default configuration
   - Exposing a command to start Momentum Journal
3. A configuration module that integrates with the Go backend's configuration

The plugin should:
- Expose a command called `:MomentumJournal` that starts a new journaling session
- Allow users to configure plugin options in their Neovim config
- Check for the Go backend service availability on startup

Use modern Neovim plugin practices (init.lua approach, not old-style vimscript) and follow Neovim plugin best practices. Include comments explaining key functionality.
```

#### Step 2.2: Neovim UI Implementation

**Tasks:**
1. Implement split window layout
2. Create buffer management for writing and agent panes
3. Set up status line integration

**Implementation Steps:**
1. Create functions to manage split windows
2. Implement buffer creation and configuration
3. Add status line integration for progress tracking
4. Set up writing buffer with markdown syntax

**LLM Prompt:**
```
Building on our Momentum Journal Neovim plugin, let's implement the UI components.

Create the UI implementation with:
1. Split window layout:
   - Left pane: Writing buffer (markdown)
   - Right pane: Agent conversation buffer
2. Buffer management:
   - Create and configure the markdown buffer for journal writing
   - Create and configure the agent conversation buffer
   - Implement functions to switch between buffers
3. Status line integration:
   - Show progress toward 3-page goal (percentage or word count)
   - Display basic session information (time elapsed, etc.)
4. Writing buffer features:
   - Markdown syntax highlighting
   - Auto-insert mode when starting

The implementation should:
- Use Neovim's API for window and buffer manipulation
- Handle buffer-specific settings (like wrap, spell check for the writing buffer)
- Provide smooth navigation between the writing and agent panes
- Ensure the UI is clean and minimalist to avoid distractions

Focus on creating a distraction-free writing environment.
```

#### Step 2.3: HTTP Client Implementation in Lua

**Tasks:**
1. Create HTTP client for communicating with Go backend
2. Implement functions for each API endpoint
3. Handle response parsing and error management

**Implementation Steps:**
1. Implement basic HTTP client functionality in Lua
2. Add functions for each API endpoint
3. Create response parsing and error handling
4. Implement progress tracking communication

**LLM Prompt:**
```
Now let's implement the HTTP client functionality in our Neovim plugin to communicate with the Go backend.

Create a Lua HTTP client module that:
1. Provides functions for each API endpoint:
   - create_session() - Calls POST /session/new
   - send_context(session_id, content) - Calls POST /session/{id}/context
   - get_session_status(session_id) - Calls GET /session/{id}/status
   - get_config() - Calls GET /config
   - update_config(config) - Calls POST /config
2. Handles:
   - Request formation and sending
   - Response parsing (JSON)
   - Error handling and reporting
   - Asynchronous requests where appropriate

The implementation should:
- Use a Lua HTTP library compatible with Neovim (like plenary.nvim's curl wrapper, or nvim-lua/plenary.nvim)
- Include timeouts and retry logic
- Display appropriate error messages to the user
- Return properly structured Lua tables from responses

Ensure the client integrates with the UI components from the previous step and properly updates the status line with progress information.
```

#### Step 2.4: Key Bindings and Commands

**Tasks:**
1. Implement `<leader>q` key binding
2. Create commands for journal operations
3. Set up buffer-local keymaps

**Implementation Steps:**
1. Add key binding for agent assistance
2. Implement buffer-local keymaps for navigation
3. Create user commands for common operations
4. Add autocommands for session behavior

**LLM Prompt:**
```
Let's implement the key bindings and commands for our Momentum Journal Neovim plugin.

Create the following functionality:
1. Key bindings:
   - `<leader>q` - Trigger agent assistance (send content to LLM and move cursor to agent window)
   - Buffer-local navigation keys to move between writing and agent panes
   - Normal mode keys for common operations

2. Commands:
   - `:MomentumStart` - Start a new journal session
   - `:MomentumSave` - Explicitly save current journal
   - `:MomentumAsk` - Same as `<leader>q` but as a command

3. Autocommands:
   - Autosave the journal buffer periodically
   - Update word count and progress in status line as user types
   - Handle buffer events appropriately

The implementation should:
- Use Neovim's API for defining keymaps and commands
- Make key bindings intuitive and consistent with Vim philosophy
- Include user feedback for operations (e.g., messages when saving)
- Support customization of key bindings via plugin configuration

Ensure all commands and keymaps are properly documented in code.
```

### Phase 3: LLM Integration and Agent Functionality

#### Step 3.1: Ollama Integration

**Tasks:**
1. Implement Ollama API client in Go
2. Create context formatting for LLM
3. Add response handling

**Implementation Steps:**
1. Create Ollama API client
2. Implement context formatting and prompt construction
3. Add streaming response handling
4. Create fallback mechanism

**LLM Prompt:**
```
Let's implement the Ollama LLM integration for our Momentum Journal backend.

Create an LLM package that:
1. Implements a client for the Ollama API:
   - Connection configuration
   - Request/response handling
   - Streaming support
2. Creates properly formatted prompts:
   - System prompt for Socratic questioning personality
   - Including journal content with proper context
   - Managing token limits
3. Processes responses to extract questions/assistance
4. Handles errors and provides fallback mechanisms

The implementation should:
- Support the configuration options from our config system
- Include proper error handling and logging
- Be tested with sample journal contents
- Support streaming responses

Also implement a prompt template system that allows for different agent personalities to be easily swapped in (starting with the Socratic questioner).
```

#### Step 3.2: OpenRouter API Integration

**Tasks:**
1. Implement OpenRouter API client
2. Ensure compatibility with Ollama interface
3. Add provider selection logic

**Implementation Steps:**
1. Create OpenRouter API client
2. Ensure compatible interface with Ollama client
3. Implement provider selection based on configuration
4. Add proper error handling and fallbacks

**LLM Prompt:**
```
Building on our LLM integration, let's add support for OpenRouter API as an alternative to Ollama.

Extend the LLM package to:
1. Implement an OpenRouter API client:
   - Connection configuration with API keys
   - Request/response handling compatible with our existing structure
   - Streaming support
2. Create a provider selection mechanism based on configuration
3. Ensure consistent interface between Ollama and OpenRouter clients
4. Add proper error handling and fallbacks

The implementation should:
- Share common code between providers where possible
- Use a factory pattern or similar to select the appropriate provider
- Support configuration changes without restarting the service
- Include tests for the OpenRouter client

Also update the configuration handling to support provider-specific options (like different model selections for each provider).
```

#### Step 3.3: Context Management and Token Handling

**Tasks:**
1. Implement context window management
2. Add token counting and limiting
3. Create prioritization for recent content

**Implementation Steps:**
1. Implement token counting functionality
2. Create context window management
3. Add prioritization for recent content
4. Implement session history handling

**LLM Prompt:**
```
Now let's implement the context management and token handling for our LLM integration.

Create a context management system that:
1. Manages the context window for LLM requests:
   - Counts tokens in content (implement or use a tokenizer)
   - Ensures requests stay within model token limits
   - Prioritizes recent content when truncation is needed
2. Implements a sliding window approach for journal content:
   - Full session available to LLM
   - Emphasis on most recent paragraphs
   - Proper handling of conversation history
3. Optimizes context usage:
   - Summarizes older content if needed
   - Tracks important themes/topics for continuity

The implementation should:
- Support different models with varying token limits
- Be efficient to avoid unnecessary processing
- Include unit tests for token counting and window management
- Integrate with our existing LLM clients

Also implement session history tracking to maintain conversation continuity between requests.
```

#### Step 3.4: Agent Personality System

**Tasks:**
1. Create system prompt templates
2. Implement personality configuration
3. Add default Socratic questioning personality

**Implementation Steps:**
1. Create system prompt template structure
2. Implement personality configuration loading
3. Add default Socratic questioning template
4. Create personality switching mechanism

**LLM Prompt:**
```
Let's implement the agent personality system for our Momentum Journal.

Create a personality system that:
1. Defines a structure for agent personalities:
   - System prompt templates
   - Behavior parameters
   - Question/response styles
2. Implements the default Socratic questioner personality:
   - Asks thought-provoking questions
   - Helps explore ideas more deeply
   - Avoids leading or judgmental framing
3. Supports loading custom personalities from markdown files:
   - Parser for personality definition files
   - Validation of personality parameters
   - Easy switching between personalities

The implementation should:
- Use a template system for prompts
- Support variable substitution in templates
- Include a library of helper functions for common personality behaviors
- Have proper documentation for creating custom personalities

Ensure the personality system integrates with our LLM package and respects the configuration settings.
```

### Phase 4: Polish and Final Features

#### Step 4.1: Progress Tracking and Statistics

**Tasks:**
1. Implement word count tracking
2. Add progress calculation toward 3-page goal
3. Create session statistics

**Implementation Steps:**
1. Implement accurate word counting
2. Create progress calculation functionality
3. Add session statistics tracking
4. Implement real-time updates

**LLM Prompt:**
```
Let's implement the progress tracking and statistics features for Momentum Journal.

Create a statistics system that:
1. Tracks writing progress:
   - Accurate word count (handling markdown syntax properly)
   - Progress toward 3-page goal (750 words by default)
   - Writing speed and session duration
2. Calculates and updates statistics in real-time:
   - Updates as the user types
   - Provides formatted output for status displays
   - Tracks historical data across sessions
3. Exposes statistics through the API:
   - Current session statistics
   - Historical trends (optional)

The implementation should:
- Be efficient to avoid performance impact during writing
- Update the Neovim status line with progress information
- Store statistics with journal entries for later reference
- Include unit tests for accuracy

Also implement a progress bar or percentage display for the 3-page goal in the Neovim status line.
```

#### Step 4.2: Error Handling and Reliability

**Tasks:**
1. Improve error handling throughout the system
2. Add reconnection logic for LLM services
3. Implement data recovery mechanisms

**Implementation Steps:**
1. Enhance error handling in critical paths
2. Add reconnection logic for LLM providers
3. Implement journal data recovery mechanisms
4. Create user-friendly error messages

**LLM Prompt:**
```
Let's improve the error handling and reliability of our Momentum Journal application.

Enhance the system with:
1. Comprehensive error handling:
   - Graceful degradation when LLM services are unavailable
   - Recovery mechanisms for network interruptions
   - User-friendly error messages in the Neovim UI
2. LLM service reliability:
   - Automatic reconnection attempts
   - Fallback providers if primary is unavailable
   - Queuing of requests during outages
3. Data protection:
   - Local caching to prevent data loss
   - Automatic recovery of unsaved changes
   - Backup mechanisms for journal entries

The implementation should:
- Focus on maintaining the writing flow despite technical issues
- Include appropriate logging for debugging
- Recover gracefully from various failure scenarios
- Provide clear feedback to the user when issues occur

Update both the Go backend and Neovim plugin to handle errors consistently and maintain a good user experience even during failures.
```

#### Step 4.3: Documentation and Testing

**Tasks:**
1. Create comprehensive documentation
2. Improve test coverage
3. Add usage examples and guides

**Implementation Steps:**
1. Create user documentation
2. Write developer documentation
3. Improve test coverage across the codebase
4. Add usage examples and guides

**LLM Prompt:**
```
Now let's focus on documentation and testing for our Momentum Journal application.

Create comprehensive documentation and tests:
1. User documentation:
   - Installation guide
   - Configuration options
   - Usage instructions with examples
   - Troubleshooting section
2. Developer documentation:
   - Architecture overview
   - API documentation
   - Extending the system (custom personalities, etc.)
   - Contribution guidelines
3. Testing improvements:
   - Increase unit test coverage
   - Add integration tests
   - Create end-to-end test scenarios
   - Performance testing

The documentation should:
- Be clear and accessible for both users and developers
- Include examples for common tasks
- Provide diagrams for architecture and data flow
- Be formatted in markdown for easy maintenance

For testing, focus on critical paths first and ensure that the core functionality has thorough test coverage.
```

#### Step 4.4: Integration and Final Polish

**Tasks:**
1. Ensure all components work together seamlessly
2. Add final UI polish and improvements
3. Perform end-to-end testing
4. Package for distribution

**Implementation Steps:**
1. Perform end-to-end integration testing
2. Refine UI elements and interaction
3. Create installation packaging
4. Add final polish and improvements

**LLM Prompt:**
```
Let's complete our Momentum Journal implementation with final integration and polish.

Implement the following:
1. End-to-end integration:
   - Ensure all components work together seamlessly
   - Fix any integration issues discovered
   - Optimize performance of critical paths
2. UI refinements:
   - Add visual polish to the Neovim interface
   - Improve status line information and formatting
   - Enhance the agent conversation presentation
3. Installation and distribution:
   - Create installation scripts
   - Package the Go backend for easy installation
   - Make the Neovim plugin compatible with popular plugin managers
4. Final testing:
   - Perform user acceptance testing
   - Check for any remaining bugs or issues
   - Validate against the original requirements

The implementation should:
- Focus on creating a smooth, distraction-free writing experience
- Prioritize reliability and performance
- Make installation and setup straightforward
- Ensure all features from the specification are working correctly

Also create a quick-start guide and ensure the default configuration provides a good out-of-the-box experience.
```

## Project Timeline and Milestones

### Milestone 1: Core Backend (Steps 1.1-1.3)
- Functioning configuration system
- Journal file management
- Basic HTTP API server

### Milestone 2: Neovim Integration (Steps 2.1-2.4)
- Working Neovim plugin
- Split window UI
- Communication with backend

### Milestone 3: LLM Integration (Steps 3.1-3.4)
- Ollama and OpenRouter support
- Context management
- Agent personality system

### Milestone 4: Final Product (Steps 4.1-4.4)
- Progress tracking
- Improved reliability
- Documentation and testing
- Final polish and distribution

## Implementation Best Practices

1. **Incremental Development**
   - Implement and test one component at a time
   - Ensure each step builds cleanly on the previous ones
   - Commit code frequently with descriptive messages

2. **Error Handling**
   - Implement robust error handling from the start
   - Use appropriate error types and messages
   - Ensure errors are logged and reported appropriately

3. **Testing**
   - Write tests alongside implementation
   - Focus on core functionality first
   - Include both unit and integration tests

4. **Documentation**
   - Document code as it's written
   - Create user documentation for each feature
   - Include examples and common use cases

5. **User Experience**
   - Prioritize the writing experience
   - Minimize distractions
   - Make interactions intuitive and consistent 