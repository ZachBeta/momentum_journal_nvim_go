# Momentum Journal Implementation Plan (Bubble Tea Version)

This document outlines a step-by-step approach to building the Momentum Journal application using Charm's Bubble Tea framework for the terminal UI. The plan is broken down into small, iterative chunks that build upon each other.

## Overall Architecture

Momentum Journal will be a single Go executable that provides:
1. A terminal UI built with Bubble Tea, featuring a two-pane layout with vim-like navigation
2. Integrated LLM functionality using Ollama locally or OpenRouter APIs
3. CLI subcommands for journal management
4. File persistence with autosave functionality

## Implementation Phases

### Phase 1: Core Application Setup
### Phase 2: Terminal UI Implementation
### Phase 3: LLM Integration
### Phase 4: Polish and Final Features

## Detailed Implementation Steps

### Phase 1: Core Application Setup

#### Step 1.1: Project Structure and CLI Framework

**Tasks:**
1. Initialize project and create basic directory structure
2. Set up CLI framework with subcommands
3. Implement configuration system with YAML support
4. Create basic logging

**Implementation Steps:**
1. Initialize Go module and create project skeleton
2. Set up CLI using a library like cobra or urfave/cli
3. Implement the `momentum new` command
4. Create configuration handling (~/.config/momentum_journal/config.yaml)

**LLM Prompt:**
```
I'm building a journaling application called "Momentum Journal" that will use Charm's Bubble Tea framework for a terminal UI. Let's start with the project setup and CLI framework.

Create the initial Go project with:
1. A main.go file that sets up a CLI using cobra (or similar library) with:
   - A root command with version and description
   - A 'new' subcommand to start a new journal entry
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

#### Step 1.2: Journal File Management

**Tasks:**
1. Implement journal entry creation and persistence
2. Create file naming and organization system
3. Add metadata tracking and storage

**Implementation Steps:**
1. Create journal package with entry creation/management
2. Implement ISO8601 timestamped filename generation
3. Add journal metadata structure and persistence
4. Create autosave functionality

**LLM Prompt:**
```
Now let's implement the file management system for our Momentum Journal application.

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
- Implement periodic autosave functionality (every 30 seconds or similar interval)
- Have unit tests for core functionality

Please also create helper functions to count words in markdown text, taking into account markdown syntax.
```

### Phase 2: Terminal UI Implementation

#### Step 2.1: Basic Bubble Tea Application

**Tasks:**
1. Set up basic Bubble Tea application structure
2. Create initial model and view
3. Implement basic event handling

**Implementation Steps:**
1. Create Bubble Tea application skeleton
2. Implement the basic model structure
3. Create simple view rendering
4. Add keyboard event handling

**LLM Prompt:**
```
Let's implement the basic Bubble Tea application structure for our Momentum Journal.

Create a terminal UI package using Bubble Tea (github.com/charmbracelet/bubbletea) that:
1. Sets up a basic application structure with:
   - A root model that will contain our application state
   - Update and View methods for the model
   - Keyboard event handling
2. Initializes the terminal UI when 'momentum new' is executed
3. Cleanly exits when the user presses 'q' or Ctrl+C

The implementation should:
- Follow Bubble Tea best practices for model organization
- Set up appropriate terminal settings (alternate screen, etc.)
- Include proper error handling
- Be integrated with the CLI framework from step 1.1

At this stage, we're just creating the basic application structure - we'll add the specific UI components in later steps.
```

#### Step 2.2: Two-Pane Layout Implementation

**Tasks:**
1. Create two-pane layout with writing area and conversation area
2. Implement vim-like navigation between panes
3. Set up status bar for progress tracking

**Implementation Steps:**
1. Implement split view layout with left/right panes
2. Create status bar at bottom of screen
3. Add vim-like keybindings for pane navigation
4. Set up focus management between panes

**LLM Prompt:**
```
Now let's implement the two-pane layout for our Momentum Journal application.

Create a UI layout using Bubble Tea and Lipgloss (github.com/charmbracelet/lipgloss) that:
1. Divides the screen into three main sections:
   - Left pane: Writing area (larger portion, approximately 60-70% of width)
   - Right pane: Conversation area with the AI agent
   - Bottom bar: Status bar showing word count, progress toward 750 words, and other relevant info
2. Implements vim-like navigation:
   - Ctrl+W followed by h/l to switch between panes
   - Tab key as an alternative to switch focus
3. Shows proper borders and styling to distinguish the panes

The implementation should:
- Use flexible layouts that adapt to terminal size
- Follow a composition-based approach with submodels for each pane
- Include state to track which pane has focus
- Show visual indication of which pane is active
- Display progress toward 3-page goal (750 words) in the status bar

Focus on getting the layout and navigation working correctly, we'll add the specific functionality for each pane in subsequent steps.
```

#### Step 2.3: Writing Pane Implementation

**Tasks:**
1. Implement vim-like text editor for the writing pane
2. Add core vim movement and editing features
3. Implement markdown syntax highlighting
4. Add word counting and progress tracking

**Implementation Steps:**
1. Create text editor component with vim-like modes
2. Implement basic movement and editing commands
3. Add markdown syntax highlighting
4. Integrate word counting with status bar

**LLM Prompt:**
```
Let's implement the writing pane with vim-like text editing capabilities for our Momentum Journal.

Create a text editor component that:
1. Implements core vim functionality:
   - Normal and insert modes with visual indicator
   - Movement commands: h/j/k/l, w/b, 0/$, gg/G
   - Basic editing: i/a/o, x, dd, yy/p
   - Search with / and n/N
2. Supports markdown text:
   - Syntax highlighting for basic markdown elements
   - Proper line wrapping for comfortable writing
3. Tracks word count and updates the status bar:
   - Counts words while ignoring markdown syntax
   - Updates progress toward 750-word goal
4. Integrates with the autosave functionality from earlier steps

The implementation should:
- Use a suitable text handling library (or implement one)
- Provide visual feedback for the current mode
- Handle UTF-8 text properly
- Update word count in real-time as the user types
- Support proper cursor positioning even with wrapped lines

Focus on creating a smooth, distraction-free writing experience with just enough vim functionality to be productive without being overwhelming.
```

#### Step 2.4: Conversation Pane Implementation

**Tasks:**
1. Create conversation display area
2. Implement structured response view with options
3. Add free-form input area
4. Set up option selection mechanics

**Implementation Steps:**
1. Create conversation history display
2. Implement agent response area with 3 options
3. Add free-form input capability
4. Set up option selection and submission

**LLM Prompt:**
```
Now let's implement the conversation pane for our Momentum Journal, which will display the AI agent's responses and allow user interaction.

Create a conversation component that:
1. Displays conversation history at the top:
   - Shows past exchanges between user and agent
   - Supports scrolling if history is lengthy
2. Shows the agent's latest response with 3 structured options:
   - Highlights selectable options
   - Allows navigation between options with arrow keys
   - Shows visual feedback for the currently selected option
3. Provides a free-form input area at the bottom:
   - Allows the user to type custom responses
   - Submits on Enter
   - Supports basic editing (backspace, etc.)
4. Implements a "thinking" state when waiting for LLM response

The implementation should:
- Use appropriate styling to distinguish between agent and user messages
- Provide clear visual hierarchy for conversation history, options, and input
- Include keyboard shortcuts for quick option selection (1/2/3 keys)
- Handle transitions between option selection and free-form input
- Use a placeholder in the free-form input ("Type your response...")

Focus on creating an intuitive conversation interface that makes it easy to both select from provided options and enter custom responses as needed.
```

### Phase 3: LLM Integration

#### Step 3.1: LLM Service Implementation

**Tasks:**
1. Create LLM service interface
2. Implement Ollama integration
3. Add OpenRouter API support
4. Create provider selection logic

**Implementation Steps:**
1. Define LLM service interface
2. Implement Ollama client
3. Add OpenRouter client
4. Create factory for selecting provider

**LLM Prompt:**
```
Let's implement the LLM service for our Momentum Journal application.

Create an LLM package that:
1. Defines a common interface for LLM providers:
   - Method to send requests and receive responses
   - Support for streaming responses
   - Error handling and retry logic
2. Implements two concrete providers:
   - Ollama client for local LLM usage
   - OpenRouter client for cloud LLM access
3. Creates a factory function that selects the appropriate provider based on configuration
4. Handles connection setup, error recovery, and proper shutdown

The implementation should:
- Support the configuration options defined earlier
- Include timeout handling and graceful degradation
- Properly manage resources (connections, goroutines, etc.)
- Support context cancellation for request interruption
- Include comprehensive error handling and logging

Focus on creating a clean, maintainable interface that abstracts away the differences between providers while allowing for provider-specific optimizations.
```

#### Step 3.2: Context Management and Prompt Templates

**Tasks:**
1. Implement context management for LLM requests
2. Create system prompt templates
3. Add journal content formatting
4. Implement Socratic questioning personality

**Implementation Steps:**
1. Create context window management
2. Implement template system for prompts
3. Add journal content formatting with focus on recent content
4. Create default Socratic questioning template

**LLM Prompt:**
```
Now let's implement the context management and prompt template system for our LLM integration.

Create a context and prompt system that:
1. Manages the context window for LLM requests:
   - Formats journal content to fit within token limits
   - Prioritizes recent paragraphs when truncation is needed
   - Includes conversation history appropriately
2. Implements a template system for system prompts:
   - Defines a structure for agent personalities
   - Supports variable substitution in templates
   - Loads templates from configuration or embedded defaults
3. Creates the default Socratic questioner personality:
   - Asks thought-provoking questions based on journal content
   - Helps the user explore ideas more deeply
   - Provides 3 different question options for the user to choose from
   - Responds thoughtfully to free-form user inputs
4. Generates effective prompts for the LLM:
   - Creates clear system instructions for the desired behavior
   - Formats user content and conversation history properly
   - Includes appropriate control parameters (temperature, etc.)

The implementation should:
- Support different LLM models with varying token limits
- Be efficient to avoid unnecessary processing
- Create prompts that consistently generate the desired 3-option format
- Include fallback handling for unexpected LLM responses

Focus on creating prompts that reliably generate helpful, thought-provoking questions presented in a consistent 3-option format plus appropriate responses to user inputs.
```

#### Step 3.3: Agent Integration with UI

**Tasks:**
1. Connect the LLM service to the UI
2. Implement request/response flow
3. Add option parsing and formatting
4. Create conversation history management

**Implementation Steps:**
1. Connect UI events to LLM service calls
2. Implement streaming response handling in UI
3. Add option parsing and formatting
4. Create conversation history management

**LLM Prompt:**
```
Let's integrate the LLM service with our terminal UI to create the full agent experience.

Implement the integration between UI and LLM service:
1. Connect user actions to LLM requests:
   - Trigger initial agent response when a new session starts
   - Send journal content when user requests help (keyboard shortcut)
   - Send user selections or free-form responses to continue the conversation
2. Handle streaming responses in the UI:
   - Show "thinking" indicator while waiting
   - Display tokens as they arrive for a smoother experience
   - Parse completed responses to extract the 3 options
3. Format and display agent responses:
   - Extract and highlight the 3 question/suggestion options
   - Format free-form responses appropriately
   - Handle errors and unexpected formats gracefully
4. Manage conversation history:
   - Store and display past exchanges
   - Include relevant history in context for new requests
   - Implement scrolling for viewing longer histories

The implementation should:
- Run LLM requests in goroutines to keep UI responsive
- Handle cancellation if user switches context
- Include proper error handling and recovery
- Maintain conversation state across UI updates

Focus on creating a seamless experience where the agent feels responsive and helpful without interrupting the writing flow.
```

### Phase 4: Polish and Final Features

#### Step 4.1: Progress Tracking and Session Management

**Tasks:**
1. Improve word counting and progress tracking
2. Add session statistics
3. Implement session persistence
4. Create session review capabilities

**Implementation Steps:**
1. Enhance word counting with markdown awareness
2. Add detailed session statistics
3. Implement session state persistence
4. Create session review functionality

**LLM Prompt:**
```
Let's enhance the progress tracking and session management capabilities of our Momentum Journal.

Implement improved session features:
1. Enhanced progress tracking:
   - More accurate word counting that handles markdown syntax
   - Visual progress indicator in status bar (percentage and/or progress bar)
   - Celebration/notification when reaching the 750-word goal
2. Session statistics:
   - Writing speed (words per minute)
   - Session duration
   - Pause/active time tracking
   - Save statistics with journal entries for later reference
3. Session persistence:
   - Save and restore session state (including conversation)
   - Allow resuming an interrupted session
   - Auto-checkpoint functionality to prevent data loss
4. Session review:
   - Add 'momentum list' command to view past sessions
   - Implement basic search/filter functionality
   - Create session summary view

The implementation should:
- Update statistics in real-time without impacting performance
- Store metadata alongside journal entries
- Use an efficient format for session persistence
- Include appropriate visualizations for progress and statistics

Focus on creating features that encourage regular writing practice and provide meaningful feedback on progress.
```

#### Step 4.2: Error Handling and Reliability

**Tasks:**
1. Improve error handling throughout the application
2. Add reconnection and recovery mechanisms
3. Implement graceful degradation
4. Create user-friendly error messages

**Implementation Steps:**
1. Enhance error handling in critical paths
2. Add reconnection logic for LLM services
3. Implement graceful degradation modes
4. Create helpful user error messages

**LLM Prompt:**
```
Let's improve the error handling and reliability of our Momentum Journal application.

Enhance the system with:
1. Comprehensive error handling:
   - Categorize errors (user, system, network, LLM-specific)
   - Implement appropriate recovery strategies for each category
   - Prevent cascading failures when components fail
2. Recovery mechanisms:
   - Automatic reconnection to LLM services
   - Local caching to prevent data loss
   - Graceful degradation when services are unavailable
3. User feedback:
   - Clear, actionable error messages
   - Status indicators for system health
   - Guidance for resolving common issues
4. Diagnostic capabilities:
   - Improved logging with appropriate levels
   - Troubleshooting mode for debugging
   - Self-checks for common configuration issues

The implementation should:
- Focus on maintaining the writing flow despite technical issues
- Handle edge cases like network interruptions, service outages
- Provide appropriate fallbacks when primary functions fail
- Balance information detail with usability in error presentations

Focus on creating a resilient application that handles failures gracefully while keeping the user informed about what's happening.
```

#### Step 4.3: Documentation and Configuration

**Tasks:**
1. Create comprehensive documentation
2. Enhance configuration options
3. Add customization capabilities
4. Implement help and guidance features

**Implementation Steps:**
1. Create user documentation
2. Add more configuration options
3. Implement customization capabilities
4. Create integrated help features

**LLM Prompt:**
```
Now let's focus on documentation and configuration for our Momentum Journal application.

Implement documentation and configuration enhancements:
1. User documentation:
   - Create README.md with installation and getting started
   - Add detailed usage guide with examples
   - Include configuration reference
   - Provide troubleshooting section
2. Enhanced configuration:
   - Add more customization options (colors, layout, keybindings)
   - Support configuration reload without restart
   - Implement validation and helpful error messages for misconfigurations
   - Add CLI flags for overriding config options
3. Customization capabilities:
   - Allow custom prompt templates for different agent personalities
   - Support theme customization (colors, borders, etc.)
   - Implement custom keybinding support
4. Integrated help:
   - Add 'momentum help' command with detailed guidance
   - Implement context-sensitive help within the application
   - Create keyboard shortcut to view available commands
   - Add tips for effective journaling practice

The implementation should:
- Make documentation accessible and comprehensive
- Support configuration through both file and CLI options
- Use sensible defaults while allowing customization
- Include examples for common customizations

Focus on making the application easy to configure while providing good documentation for users to get the most out of it.
```

#### Step 4.4: Polish and Distribution

**Tasks:**
1. Perform final UI polish and refinements
2. Optimize performance and resource usage
3. Create installation package
4. Add final quality-of-life features

**Implementation Steps:**
1. Polish UI elements and interactions
2. Optimize performance in key areas
3. Create installation scripts and packages
4. Add finishing touches and refinements

**LLM Prompt:**
```
Let's complete our Momentum Journal implementation with final polish and distribution preparation.

Implement these final enhancements:
1. UI polish:
   - Refine colors, borders, and spacing for optimal aesthetics
   - Add subtle animations for transitions (if appropriate)
   - Ensure consistent styling across all components
   - Optimize layout for different terminal sizes
2. Performance optimizations:
   - Profile and optimize hot paths
   - Reduce memory usage where possible
   - Optimize startup time
   - Ensure efficient handling of large journal entries
3. Distribution preparation:
   - Create installation script
   - Add Homebrew formula template
   - Prepare for distribution via GitHub releases
   - Include appropriate license and contribution guidelines
4. Final features:
   - Add keyboard shortcut reference
   - Implement session summary after writing
   - Add export functionality to other formats
   - Create a welcome experience for first-time users

The implementation should:
- Focus on creating a polished, professional feel
- Ensure the application runs well on various environments
- Make installation and updates straightforward
- Address any remaining rough edges or inconsistencies

Focus on making a memorable, delightful application that users will want to return to daily for their writing practice.
```

## Project Timeline and Milestones

### Milestone 1: Core Application (Steps 1.1-2.2)
- Functioning CLI and configuration
- Journal file management
- Basic two-pane Bubble Tea UI

### Milestone 2: Interactive Components (Steps 2.3-2.4)
- Working vim-like text editor
- Conversation pane with options and input

### Milestone 3: LLM Integration (Steps 3.1-3.3)
- Ollama and OpenRouter support
- Context management
- Agent integration with UI

### Milestone 4: Final Product (Steps 4.1-4.4)
- Enhanced progress tracking
- Improved reliability
- Documentation and configuration
- Polish and distribution

## Implementation Best Practices

1. **Incremental Development**
   - Implement one component at a time
   - Test each component thoroughly before moving on
   - Commit code frequently with descriptive messages

2. **Code Organization**
   - Use clear package boundaries
   - Follow Go idioms and best practices
   - Implement clean interfaces between components

3. **UI Considerations**
   - Ensure responsive UI even during LLM operations
   - Test with different terminal sizes and types
   - Consider accessibility where possible

4. **Testing**
   - Write unit tests for core functionality
   - Add integration tests for component interactions
   - Include end-to-end tests for critical flows

5. **User Experience**
   - Focus on creating a distraction-free writing environment
   - Ensure keyboard shortcuts are intuitive
   - Make error messages helpful and actionable 