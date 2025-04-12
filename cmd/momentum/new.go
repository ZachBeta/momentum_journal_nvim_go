package main

import (
	"fmt"
	"log" // Use standard log for fatal errors from Bubble Tea

	"github.com/ZachBeta/momentum_journal_nvim_go/internal/journal" // Adjusted import path
	"github.com/ZachBeta/momentum_journal_nvim_go/internal/tui"     // Import the new TUI package
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Start a new journal entry",
	Long: `Create a new journal entry and open the Momentum Journal interface.
This command starts a new writing session with the specified settings.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create journal manager
		journalManager, err := journal.NewManager(cfg, logger)
		if err != nil {
			// Log error using zap before returning
			logger.Error("Failed to create journal manager", zap.Error(err))
			return fmt.Errorf("failed to create journal manager: %w", err)
		}

		// Create new entry
		// TBD: We might want to pass the entry or its path to the TUI model later
		_, err = journalManager.CreateEntry()
		if err != nil {
			// Log error using zap before returning
			logger.Error("Failed to create journal entry", zap.Error(err))
			return fmt.Errorf("failed to create journal entry: %w", err)
		}

		// Initialize the TUI model
		tuiModel := tui.InitialModel()

		// Create and run the Bubble Tea program
		// Using tea.WithAltScreen() provides a dedicated screen for the TUI
		// Using tea.WithMouseCellMotion() enables mouse support (optional but often useful)
		p := tea.NewProgram(tuiModel, tea.WithAltScreen()) //, tea.WithMouseCellMotion())

		logger.Info("Starting Momentum Journal TUI...")

		// Run the program. This blocks until the program exits.
		if _, err := p.Run(); err != nil {
			// Log the error from Bubble Tea using standard log or zap
			logger.Error("Error running Bubble Tea program", zap.Error(err))
			// Use standard log for fatal errors that terminate the app immediately after TUI fails
			log.Fatalf("Alas, there's been an error: %v", err)
			// The return below might not be reached if log.Fatalf exits, but good practice.
			return fmt.Errorf("error running TUI: %w", err)
		}

		logger.Info("Momentum Journal TUI finished.")
		// TBD: Perform any cleanup after the TUI exits if needed.
		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
