# Momentum Journal - Phase 2 Tutorial: Building the Terminal UI

This tutorial guides you through implementing the terminal user interface (TUI) for the Momentum Journal application using Charm's Bubble Tea and Lipgloss libraries. This corresponds to Phase 2 of the [main implementation plan](bubbletea_plan.md).

**Goal:** Create a functional two-pane terminal UI with a writing area, a conversation area, a status bar, and basic vim-like navigation and editing capabilities.

**Prerequisites:** Completion of Phase 1 (Core Application Setup), including the CLI framework and basic journal file management.

## Step 2.1: Basic Bubble Tea Application Structure

**Goal:** Set up the foundational Bubble Tea application structure that will host our UI components.

**Tasks:**
1.  **Create UI Package:** Set up a dedicated package (e.g., `tui`) for the Bubble Tea code.
2.  **Root Model:** Define the main `model` struct that will hold the application's state (e.g., current view, focused pane, sub-models for panes).
3.  **Implement `Init`, `Update`, `View`:**
    *   `Init`: Perform initial setup, potentially loading data or setting initial state. Return an initial `Cmd`.
    *   `Update`: Handle incoming messages (`tea.Msg`) like keyboard events (`tea.KeyMsg`), window size changes (`tea.WindowSizeMsg`), and custom messages. Update the model's state and return the updated model and any new `Cmd`. Implement basic exit logic (e.g., on `Ctrl+C` or `q`).
    *   `View`: Render the UI based on the current model state. Return a `string`.
4.  **Integrate with CLI:** Modify the `momentum new` command (from Phase 1) to initialize and run the Bubble Tea program (`tea.NewProgram(initialModel).Run()`).
5.  **Terminal Setup:** Ensure the Bubble Tea program uses the alternate screen buffer and handles terminal resizing.

**Verification:** Running `momentum new` should launch a blank Bubble Tea application that exits cleanly when `Ctrl+C` is pressed.

**Code Snippets/Guidance:**

```go
package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	// ... other imports
)

type model struct {
	// TBD: State fields like dimensions, focused pane, sub-models
	quitting bool
}

func InitialModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	// TBD: Initial commands (e.g., tick for autosave)
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	// TBD: Handle other messages (WindowSizeMsg, custom messages)
	}
	// TBD: Delegate updates to sub-models
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "Quitting...
"
	}
	// TBD: Render layout using Lipgloss
	return "Hello, Bubble Tea!
" // Placeholder
}

// --- In main.go or relevant command file ---
// ...
// cmd := tui.InitialModel()
// p := tea.NewProgram(cmd, tea.WithAltScreen())
// if err := p.Run(); err != nil {
//     log.Fatal(err)
// }
// ...
```

## Step 2.2: Two-Pane Layout Implementation

**Goal:** Divide the terminal screen into a writing pane, a conversation pane, and a status bar using Lipgloss, and implement basic navigation between the panes.

**Tasks:**
1.  **Import Lipgloss:** Add `github.com/charmbracelet/lipgloss` to your imports.
2.  **Define Layout Structure:** In the root `model`, add fields to store dimensions (`width`, `height`) and potentially Lipgloss styles.
3.  **Sub-models:** Plan for sub-models for the `writingPane` and `conversationPane`. Add placeholders in the root model.
4.  **Status Bar:** Define a simple status bar component/model.
5.  **Handle `WindowSizeMsg`:** Update the model's `width` and `height` when the terminal is resized. Propagate size changes to sub-components if necessary.
6.  **Implement `View` Logic:** Use Lipgloss functions (`lipgloss.JoinVertical`, `lipgloss.JoinHorizontal`, `lipgloss.Place`, styling methods) to arrange the panes and status bar. The layout should be flexible and adapt to the available `width` and `height`. Aim for ~60-70% width for the writing pane.
7.  **Focus Management:** Add state to the root model to track the currently focused pane (e.g., `focusedPane int`).
8.  **Navigation:** Implement keybindings in the root `Update` function (e.g., `Ctrl+W` + `h/l`, or `Tab`) to change the `focusedPane` state.
9.  **Visual Feedback:** Use Lipgloss styling (e.g., border colors) in the `View` method to indicate which pane is currently active/focused.
10. **Status Bar Content:** Display basic information like word count (initially 0) and the target word count in the status bar.

**Verification:** Running `momentum new` should show a styled two-pane layout with a status bar. Pressing the navigation keys should visually change the focus between the panes. The layout should adjust reasonably when the terminal is resized.

**Code Snippets/Guidance:**

```go
// --- In tui/model.go ---
import (
	"github.com/charmbracelet/lipgloss"
	// ...
)

type focusState int

const (
	writingPane focusState = iota
	conversationPane
)

type model struct {
	width          int
	height         int
	focusedPane    focusState
	writingModel   writingModel   // Placeholder struct
	convoModel     convoModel     // Placeholder struct
	statusBarModel statusBarModel // Placeholder struct
	// Lipgloss styles
	quitting bool
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// TBD: Update sub-model sizes

	case tea.KeyMsg:
		switch msg.String() {
		// ... (quit keys)
		case "tab", "ctrl+w": // Simplified example, Ctrl+W usually needs sequence
			if m.focusedPane == writingPane {
				m.focusedPane = conversationPane
			} else {
				m.focusedPane = writingPane
			}
		}
	}
	// TBD: Delegate updates based on focus
	// cmd := m.updateFocusedPane(msg)
	return m, nil // , cmd
}

func (m model) View() string {
	// ... (quitting view)

	// TBD: Get views from sub-models
	writingView := "Writing Pane
(Focused: " + fmt.Sprint(m.focusedPane == writingPane) + ")"
	convoView := "Conversation Pane
(Focused: " + fmt.Sprint(m.focusedPane == conversationPane) + ")"
	statusBarView := "Status: Word Count 0/750"

	// Calculate pane dimensions
	statusBarHeight := lipgloss.Height(statusBarView)
	mainHeight := m.height - statusBarHeight
	writingWidth := int(float64(m.width) * 0.65) // ~65%
	convoWidth := m.width - writingWidth

	// Apply styling based on focus (example)
	// style := lipgloss.NewStyle().Border(lipgloss.NormalBorder())
	// focusedStyle := lipgloss.NewStyle().Border(lipgloss.ThickBorder()).BorderForeground(lipgloss.Color("205"))

	// TBD: Apply styles to writingView, convoView based on m.focusedPane

	mainView := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(writingWidth).Height(mainHeight).Render(writingView),
		lipgloss.NewStyle().Width(convoWidth).Height(mainHeight).Render(convoView),
	)

	return lipgloss.JoinVertical(lipgloss.Left, mainView, statusBarView)
}

// TBD: Define writingModel, convoModel, statusBarModel structs and their methods
```

## Step 2.3: Writing Pane Implementation

**Goal:** Implement the core text editing functionality within the left pane using vim-like keybindings.

**Tasks:**
1.  **Create `writingModel`:** Define the struct for the writing pane, likely using `github.com/charmbracelet/bubbles/textarea` or a custom implementation for more control.
2.  **Vim Modes:** Implement state management for Normal and Insert modes.
3.  **Text Area:** Integrate `textarea` or build functionality to display and edit multi-line text. Handle scrolling and line wrapping.
4.  **Keybindings (Normal Mode):** Implement basic movement (`h/j/k/l`, `w/b`, `0/$`, `gg/G`) and editing commands (`i/a/o` to enter Insert mode, `x`, `dd`, `yy/p`).
5.  **Keybindings (Insert Mode):** Allow text input, handle `Esc` to return to Normal mode.
6.  **Delegate Updates:** In the root model's `Update`, pass relevant messages (especially `tea.KeyMsg`) to the `writingModel`'s `Update` method *only when it's focused*.
7.  **Get View:** The `writingModel`'s `View()` method should return the rendered text area content as a string.
8.  **Word Count:** Implement logic (or use a helper from Phase 1) to count words in the text area's content, ignoring markdown. Update the status bar model (potentially via a message or shared state).
9.  **Autosave Integration:** Connect the text area content changes to the autosave mechanism (from Phase 1), perhaps triggering saves periodically or on specific events.
10. **(Optional) Syntax Highlighting:** Explore libraries or methods for basic markdown syntax highlighting within the text area.

**Verification:** When the writing pane is focused, you should be able to navigate using vim keys, switch between Normal and Insert modes, type text, and perform basic edits (`dd`, `p`). The word count in the status bar should update as you type.

**Code Snippets/Guidance:**

```go
// --- In tui/writing_pane.go ---
package tui

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	// ...
)

type writingMode int

const (
	modeNormal writingMode = iota
	modeInsert
)

type writingModel struct {
	textarea textarea.Model
	mode     writingMode
	width    int
	height   int
	// TBD: Yank buffer, etc.
}

func NewWritingModel() writingModel {
	ta := textarea.New()
	ta.Placeholder = "Start your morning pages..."
	// TBD: Configure textarea (styles, prompts)
	return writingModel{textarea: ta, mode: modeNormal}
}

func (m *writingModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.textarea.SetWidth(w)
	m.textarea.SetHeight(h) // Adjust as needed for borders/lines
}

func (m writingModel) Init() tea.Cmd {
	return textarea.Blink // Start cursor blinking
}

func (m writingModel) Update(msg tea.Msg) (writingModel, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.mode == modeInsert {
			switch msg.Type {
			case tea.KeyEsc:
				m.mode = modeNormal
				m.textarea.Blur() // Show static cursor in normal mode
			default:
				m.textarea, cmd = m.textarea.Update(msg)
				cmds = append(cmds, cmd)
				// TBD: Trigger word count update message
			}
		} else { // modeNormal
			switch msg.String() {
			case "i":
				m.mode = modeInsert
				m.textarea.Focus()
				cmds = append(cmds, textarea.Blink)
			case "a": // TBD
			case "o": // TBD
			case "h", "j", "k", "l": // TBD: Implement movement
			case "w", "b": // TBD
			case "g": // TBD: Handle gg
			case "G": // TBD
			case "x": // TBD: Delete character
			case "d": // TBD: Handle dd
			case "y": // TBD: Handle yy
			case "p": // TBD: Paste
			// TBD: Add other normal mode commands
			default:
				// Pass other keys for potential default textarea behavior (like scrolling)
				m.textarea, cmd = m.textarea.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m writingModel) View() string {
	// TBD: Add mode indicator
	modeIndicator := "[NORMAL]"
	if m.mode == modeInsert {
		modeIndicator = "[INSERT]"
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		modeIndicator,
		m.textarea.View(),
	)
}

func (m writingModel) WordCount() int {
	// TBD: Implement word count logic on m.textarea.Value()
	return 0
}
```

## Step 2.4: Conversation Pane Implementation

**Goal:** Build the right pane to display conversation history, show the AI agent's responses (with options), and allow user input.

**Tasks:**
1.  **Create `convoModel`:** Define the struct for the conversation pane. Consider using `github.com/charmbracelet/bubbles/viewport` for the history and maybe `textarea` for the input.
2.  **Conversation History:** Store past user messages and agent responses (e.g., in a slice of structs).
3.  **History Display:** Use the `viewport` to display the formatted conversation history, making it scrollable.
4.  **Agent Response Area:** Design how to display the latest agent message, potentially highlighting 3 selectable options.
5.  **Input Area:** Add a `textarea` (single-line or limited height) at the bottom for free-form user input.
6.  **Option Selection:** Implement state and keybindings (arrow keys, numbers) to navigate and select one of the agent's options when available.
7.  **Input Handling:** When the input area is focused (within the conversation pane), allow typing and handle `Enter` to submit the message.
8.  **Delegate Updates:** Similar to the writing pane, the root model should delegate updates to the `convoModel` when it's focused. The `convoModel` itself needs to manage focus between history/options and the input area.
9.  **Get View:** The `convoModel`'s `View()` method should render the history, current agent response/options, and input area.
10. **"Thinking" State:** Add a visual indicator (e.g., "Agent is thinking...") when waiting for an LLM response (this will be triggered by messages in Phase 3).

**Verification:** When the conversation pane is focused, you should be able to scroll through placeholder history, potentially select placeholder options (if implemented visually), and type text into the input area.

**Code Snippets/Guidance:**

```go
// --- In tui/convo_pane.go ---
package tui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	// ...
)

// TBD: Define message struct (User/Agent, Content, Options)

type convoModel struct {
	viewport    viewport.Model
	inputArea   textarea.Model
	messages    []string // Placeholder for structured messages
	width       int
	height      int
	isThinking  bool
	selectedOpt int // Index of selected option, -1 if none
	// TBD: State for focus (history/options vs input)
}

func NewConvoModel() convoModel {
	vp := viewport.New(0, 0) // Size set later
	ta := textarea.New()
	ta.Placeholder = "Type your response or select an option..."
	ta.SetHeight(1) // Single line input
	// TBD: Configure viewport and textarea
	return convoModel{
		viewport:  vp,
		inputArea: ta,
		messages:  []string{"Agent: Welcome!", "User: Hello."}, // Example history
	}
}

func (m *convoModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	inputHeight := m.inputArea.Height() + 2 // + borders/padding
	m.viewport.Width = w
	m.viewport.Height = h - inputHeight
	m.inputArea.SetWidth(w - 2) // Account for borders/padding
	m.viewport.SetContent(m.renderHistory()) // Update content on resize
}

func (m convoModel) Init() tea.Cmd {
	return nil
}

func (m convoModel) Update(msg tea.Msg) (convoModel, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	// TBD: Handle focus switching between history/options and input area

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// TBD: Handle keys based on focus (viewport scrolling, option selection, input typing)
		// Example: If input is focused
		// if m.inputFocused {
		// 	if msg.Type == tea.KeyEnter {
		// 		// Handle submission
		// 		userInput := m.inputArea.Value()
		// 		m.messages = append(m.messages, "User: "+userInput)
		//      m.inputArea.Reset()
		// 		// TBD: Send message to LLM (Phase 3)
		// 		m.isThinking = true
		// 	} else {
		// 		m.inputArea, cmd = m.inputArea.Update(msg)
		// 		cmds = append(cmds, cmd)
		// 	}
		// } else { // History/Options focused
		// 	// Handle viewport scrolling, option selection (Up/Down/Enter, 1/2/3)
		// 	m.viewport, cmd = m.viewport.Update(msg)
		// 	cmds = append(cmds, cmd)
		// }
	}

	m.viewport.SetContent(m.renderHistory()) // Update content after potential message add

	return m, tea.Batch(cmds...)
}

func (m convoModel) renderHistory() string {
	// TBD: Format m.messages nicely
	// TBD: Render current agent options if available and highlight selectedOpt
	history := lipgloss.JoinVertical(lipgloss.Left, m.messages...)
	if m.isThinking {
		history += "
Agent is thinking..."
	}
	return history
}

func (m convoModel) View() string {
	// TBD: Style input area, potentially indicating focus
	return lipgloss.JoinVertical(lipgloss.Left,
		m.viewport.View(),
		m.inputArea.View(),
	)
}

```

---

Completing these steps will result in a visually structured TUI ready for the LLM integration in Phase 3. Remember to test frequently and commit your changes incrementally. 