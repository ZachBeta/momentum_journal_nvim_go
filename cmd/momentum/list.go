package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ZachBeta/momentum_journal_nvim_go/internal/journal" // Adjusted import path
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List journal entries",
	Long:  `List all journal entries with basic information.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create journal manager
		journalManager, err := journal.NewManager(cfg, logger)
		if err != nil {
			return fmt.Errorf("failed to create journal manager: %w", err)
		}

		// List entries
		entries, err := journalManager.ListEntries()
		if err != nil {
			return fmt.Errorf("failed to list journal entries: %w", err)
		}

		if len(entries) == 0 {
			fmt.Println("No journal entries found.")
			return nil
		}

		// Create a tab writer for nice formatting
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "DATE\tTIME\tWORDS\tCOMPLETE\tFILE")
		fmt.Fprintln(w, "----\t----\t-----\t--------\t----")

		for _, entry := range entries {
			fmt.Fprintf(w, "%s\t%s\t%d\t%v\t%s\n",
				entry.CreatedAt.Format("2006-01-02"),
				entry.CreatedAt.Format("15:04"),
				entry.WordCount,
				entry.IsCompleted,
				entry.FileName)
		}

		w.Flush()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
