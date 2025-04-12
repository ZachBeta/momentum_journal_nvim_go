package main

import (
	"fmt"

	"github.com/ZachBeta/momentum_journal_nvim_go/internal/journal" // Adjusted import path
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
			return fmt.Errorf("failed to create journal manager: %w", err)
		}

		// Create new entry
		entry, err := journalManager.CreateEntry()
		if err != nil {
			return fmt.Errorf("failed to create journal entry: %w", err)
		}

		logger.Info("Created new journal entry",
			zap.String("file", entry.FileName),
			zap.String("path", entry.FilePath))

		// In a later step, we'll start the Bubble Tea UI here
		fmt.Printf("Created new journal entry: %s\n", entry.FilePath)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
