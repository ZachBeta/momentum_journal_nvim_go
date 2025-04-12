package tui

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// writingMode determines the behavior of keybindings in the writing pane.
type writingMode int

const (
	modeNormal writingMode = iota
	modeInsert
)

// writingModel holds the state for the text editing pane.
type writingModel struct {
	textarea textarea.Model
	mode     writingMode
	width    int
	height   int
	// TBD: Yank buffer, word count, etc.
}

// NewWritingModel creates a new instance of the writing pane model.
func NewWritingModel() writingModel {
	ta := textarea.New()
	ta.Placeholder = "Start your morning pages..."
	ta.ShowLineNumbers = true // Let's enable line numbers
	ta.Focus()                // Start focused so the cursor initially shows

	// Configure styles (optional, can be customized later)
	// ta.FocusedStyle.CursorLine = lipgloss.NewStyle().Background(lipgloss.Color("62"))
	// ta.BlurredStyle.CursorLine = lipgloss.NewStyle()

	m := writingModel{
		textarea: ta,
		mode:     modeInsert, // Start in Insert mode for immediate typing
	}
	// Initially blur it, the main model will focus it based on state
	m.textarea.Blur()
	return m
}

// SetSize updates the dimensions of the writing pane.
func (m *writingModel) SetSize(w, h int) {
	m.width = w
	m.height = h // We might need to adjust height for the mode indicator
	indicatorHeight := lipgloss.Height(m.renderModeIndicator())
	m.textarea.SetWidth(w)
	m.textarea.SetHeight(h - indicatorHeight)
}

// Init initializes the writing model, returning an initial command.
func (m writingModel) Init() tea.Cmd {
	// If starting in Insert mode, start blinking the cursor.
	if m.mode == modeInsert {
		return textarea.Blink
	}
	return nil
}

// Update handles messages for the writing pane.
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
				return m, nil     // Consume Esc
			default:
				// Default textarea behavior for input
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
			case "a": // TBD: Insert after cursor
			case "o": // TBD: Insert new line below
			case "h", "j", "k", "l", "up", "down", "left", "right": // Basic movement
				// Pass movement keys to the textarea in normal mode too
				m.textarea, cmd = m.textarea.Update(msg)
				cmds = append(cmds, cmd)
			case "w", "b": // TBD: Word movement
			case "g": // TBD: Handle gg
			case "G": // TBD: Go to end
			case "x": // TBD: Delete character
			case "d": // TBD: Handle dd
			case "y": // TBD: Handle yy
			case "p": // TBD: Paste
			default:
				// Pass other keys (like PageUp/PageDown) for default textarea behavior
				m.textarea, cmd = m.textarea.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

// View renders the writing pane UI.
func (m writingModel) View() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		m.renderModeIndicator(),
		m.textarea.View(),
	)
}

// renderModeIndicator returns the visual indicator for the current mode.
func (m writingModel) renderModeIndicator() string {
	indicator := "[NORMAL]"
	if m.mode == modeInsert {
		indicator = "[INSERT]"
	}
	// TBD: Style the indicator (e.g., different colors)
	return lipgloss.NewStyle().Padding(0, 1).Render(indicator)
}

// Focus sets the writing pane to be focused.
func (m *writingModel) Focus() tea.Cmd {
	if m.mode == modeInsert {
		return m.textarea.Focus()
	}
	return nil
}

// Blur removes focus from the writing pane.
func (m *writingModel) Blur() {
	m.textarea.Blur()
}

// WordCount calculates the number of words in the textarea.
func (m writingModel) WordCount() int {
	// TBD: Implement more accurate word count logic (from Phase 1 helper)
	// For now, a simple split by space approximation
	return len(m.textarea.Value()) // Replace with actual count
}
