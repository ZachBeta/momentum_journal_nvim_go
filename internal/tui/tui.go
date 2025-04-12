package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	// log "github.com/sirupsen/logrus" // TBD: Add logging if needed
)

// focusState determines which pane has keyboard focus.
// We might add more states later (e.g., for command input).
type focusState int

const (
	writingPane focusState = iota
	conversationPane
)

// --- Sub-model Placeholders --- //

// TBD: Implement writingModel fully in Step 2.3
/* REMOVE PLACEHOLDER
type writingModel struct {
	width  int
	height int
}

func newWritingModel() writingModel {
	return writingModel{}
}
func (m *writingModel) SetSize(w, h int) { m.width, m.height = w, h }
func (m writingModel) View() string     { return "Writing Pane Placeholder" }
*/

// TBD: Implement convoModel fully in Step 2.4
type convoModel struct {
	width  int
	height int
}

func newConvoModel() convoModel {
	return convoModel{}
}
func (m *convoModel) SetSize(w, h int) { m.width, m.height = w, h }
func (m convoModel) View() string      { return "Conversation Pane Placeholder" }

// TBD: Implement statusBarModel fully later
type statusBarModel struct {
	width int
}

func newStatusBarModel() statusBarModel {
	return statusBarModel{}
}
func (m *statusBarModel) SetSize(w int) { m.width = w }
func (m statusBarModel) View() string {
	// Simple placeholder status
	return lipgloss.NewStyle().
		// Background(lipgloss.Color("7")). // Example styling
		// Foreground(lipgloss.Color("0")).
		Width(m.width).
		Render("Status: Word Count 0/750 | [Tab] to switch panes | [q] to quit")
}

// --- Main Model --- //

// model represents the state of the TUI application.
type model struct {
	width          int
	height         int
	focusedPane    focusState
	writingModel   writingModel
	convoModel     convoModel
	statusBarModel statusBarModel

	// Styles (can be customized later)
	paneStyle     lipgloss.Style
	focusedStyle  lipgloss.Style
	statusBarSyle lipgloss.Style

	quitting bool
}

// InitialModel creates the starting state for the Bubble Tea application.
func InitialModel() model {
	// Define base styles
	// We can make these configurable later
	paneStyle := lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")) // Dimmed border

	focusedStyle := paneStyle.Copy().
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("205")) // Highlighted border

	statusBarSyle := lipgloss.NewStyle()
	// statusBarSyle := lipgloss.NewStyle().
	// 	Background(lipgloss.Color("7")).
	// 	Foreground(lipgloss.Color("0"))

	m := model{
		writingModel:   NewWritingModel(),
		convoModel:     newConvoModel(),
		statusBarModel: newStatusBarModel(),
		focusedPane:    writingPane, // Start focus in writing pane
		paneStyle:      paneStyle,
		focusedStyle:   focusedStyle,
		statusBarSyle:  statusBarSyle,
	}
	return m
}

// Init is the first command that runs when the Bubble Tea program starts.
func (m model) Init() tea.Cmd {
	// Initialize sub-models and gather their initial commands
	// For now, only writingModel might have an initial command (like Blink)
	return m.writingModel.Init()
}

// Update handles incoming messages and updates the model's state.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	// Handle window resize events.
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Recalculate sizes and update sub-models
		m.updateSizes()
		// TBD: We might need to return update commands from sub-models if they react to resize

	// Handle keyboard events.
	case tea.KeyMsg:
		switch msg.String() {
		// Quit the application.
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		// Switch focus between panes.
		case "tab":
			if m.focusedPane == writingPane {
				m.focusedPane = conversationPane
				m.writingModel.Blur() // Blur the writing pane
				// cmd = m.convoModel.Focus() // TBD: Focus convo pane when implemented
			} else {
				m.focusedPane = writingPane
				// m.convoModel.Blur() // TBD: Blur convo pane
				cmd = m.writingModel.Focus() // Focus the writing pane
			}
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...) // Consume the message and return focus command

		// TBD: Handle Ctrl+W + h/l for switching focus as an alternative

		default:
			// Delegate other key presses to the focused pane
			switch m.focusedPane {
			case writingPane:
				m.writingModel, cmd = m.writingModel.Update(msg)
				cmds = append(cmds, cmd)
			case conversationPane:
				// TBD: Delegate to convoModel when implemented
				// m.convoModel, cmd = m.convoModel.Update(msg)
				// cmds = append(cmds, cmd)
			}
		}

		// TBD: Handle custom messages (e.g., word count updates, LLM responses).
	}

	// Update sizes again in case a command changed something that affects layout
	// This might be redundant if sub-model updates don't affect layout
	// m.updateSizes()

	// TBD: Update non-focused sub-models if necessary (e.g., pass tick messages)

	return m, tea.Batch(cmds...)
}

// updateSizes calculates and sets the dimensions for the sub-models based on the main model's width and height.
func (m *model) updateSizes() {
	statusBarHeight := lipgloss.Height(m.statusBarModel.View()) // Calculate actual height
	mainHeight := m.height - statusBarHeight

	// Simple 65/35 split, adjust as needed
	writingWidth := int(float64(m.width) * 0.65)
	// Ensure minimum width or handle edge cases if necessary
	if writingWidth < 10 {
		writingWidth = 10
	}
	convoWidth := m.width - writingWidth
	if convoWidth < 10 {
		convoWidth = 10
		// Adjust writing width if convo width hits minimum
		writingWidth = m.width - convoWidth
	}

	m.writingModel.SetSize(writingWidth-m.paneStyle.GetHorizontalBorderSize(), mainHeight-m.paneStyle.GetVerticalBorderSize())
	m.convoModel.SetSize(convoWidth-m.paneStyle.GetHorizontalBorderSize(), mainHeight-m.paneStyle.GetVerticalBorderSize())
	m.statusBarModel.SetSize(m.width)
}

// View renders the UI based on the current model state.
func (m model) View() string {
	// If quitting, show a final message.
	if m.quitting {
		return "Saving and quitting Momentum Journal...\n"
	}

	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	// Get views from sub-models
	writingView := m.writingModel.View()
	convoView := m.convoModel.View()
	statusBarView := m.statusBarModel.View()

	// Apply focus styling
	var styledWritingView, styledConvoView string
	if m.focusedPane == writingPane {
		styledWritingView = m.focusedStyle.Render(writingView)
		styledConvoView = m.paneStyle.Render(convoView)
	} else {
		styledWritingView = m.paneStyle.Render(writingView)
		styledConvoView = m.focusedStyle.Render(convoView)
	}

	// Set dimensions on the styled views before joining
	// Use GetWidth/GetHeight to account for borders/padding set by the style
	styledWritingView = lipgloss.NewStyle().Width(m.writingModel.width + m.paneStyle.GetHorizontalFrameSize()).Height(m.writingModel.height + m.paneStyle.GetVerticalFrameSize()).Render(styledWritingView)
	styledConvoView = lipgloss.NewStyle().Width(m.convoModel.width + m.paneStyle.GetHorizontalFrameSize()).Height(m.convoModel.height + m.paneStyle.GetVerticalFrameSize()).Render(styledConvoView)

	// Join the panes horizontally
	mainPane := lipgloss.JoinHorizontal(
		lipgloss.Top, // Align panes to the top
		styledWritingView,
		styledConvoView,
	)

	// Join the main pane and status bar vertically
	fullView := lipgloss.JoinVertical(
		lipgloss.Left, // Align items to the left
		mainPane,
		m.statusBarSyle.Width(m.width).Render(statusBarView), // Ensure status bar takes full width
	)

	return fullView
}
