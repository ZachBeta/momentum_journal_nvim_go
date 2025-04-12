package main

import (
	"fmt"
	"os"

	"github.com/ZachBeta/momentum_journal_nvim_go/internal/config"  // Adjusted import path
	"github.com/ZachBeta/momentum_journal_nvim_go/internal/logging" // Adjusted import path
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	debug      bool
	configFile string
	logger     *zap.Logger
	cfg        *config.Config
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "momentum",
	Short: "Momentum Journal - A terminal-based journaling app with AI assistance",
	Long: `Momentum Journal is a minimalist tool for following "The Artist's Way" journaling practice.
It provides a distraction-free writing environment with AI support to help maintain writing flow.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error

		// Initialize logger
		logger, err = logging.NewLogger(debug)
		if err != nil {
			return fmt.Errorf("failed to initialize logger: %w", err)
		}

		// Load configuration
		cfg, err = config.Load(logger)
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Add persistent flags for the root command
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug mode")
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Config file (default is $HOME/.config/momentum_journal/config.yaml)")
}
