# Momentum Journal Implementation Tutorial - Phase 1

This tutorial will guide you through implementing Phase 1 of the Momentum Journal application: **Core Application Setup**. By the end of this tutorial, you'll have a functioning CLI application with configuration management and journal file handling.

## Prerequisites

Before starting, ensure you have:

- Go 1.21+ installed
- A code editor of your choice
- Basic knowledge of Go programming
- Git installed (for version control)

## Project Overview

In Phase 1, we'll implement:
1. Project structure and CLI framework
2. Configuration system with YAML support
3. Journal file management
4. Basic logging

Let's get started!

## Step 1: Initialize the Project

First, let's create a new directory and initialize a Go module:

```bash
# Create project directory
mkdir -p momentum_journal
cd momentum_journal

# Initialize Go module
go mod init github.com/yourusername/momentum_journal

# Create basic directory structure
mkdir -p cmd/momentum
mkdir -p internal/config
mkdir -p internal/journal
mkdir -p internal/ui
mkdir -p internal/llm
```

## Step 2: Set Up Dependencies

We'll need several dependencies for our project:

```bash
# CLI framework
go get github.com/spf13/cobra

# Configuration
go get github.com/spf13/viper

# Logging
go get go.uber.org/zap

# YAML support
go get gopkg.in/yaml.v3

# UI dependencies (we'll use these later)
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/lipgloss
go get github.com/charmbracelet/bubbles
```

## Step 3: Create Configuration Package

Let's start by implementing the configuration package. Create a file at `internal/config/config.go`:

```go
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	// LLM provider settings
	LLM struct {
		Provider   string `yaml:"provider"`    // "ollama" or "openrouter"
		APIKey     string `yaml:"api_key"`     // API key for openrouter
		ModelName  string `yaml:"model_name"`  // Model to use e.g., "llama3" for Ollama
		Endpoint   string `yaml:"endpoint"`    // API endpoint
		MaxTokens  int    `yaml:"max_tokens"`  // Maximum tokens for response
		Temperature float64 `yaml:"temperature"` // Temperature for generation
	} `yaml:"llm"`

	// Journal settings
	Journal struct {
		StorageDir   string `yaml:"storage_dir"`    // Directory to store journal files
		WordCountGoal int    `yaml:"word_count_goal"` // Default 750 words (3 pages)
		AutosaveInterval int `yaml:"autosave_interval"` // Autosave interval in seconds
	} `yaml:"journal"`

	// UI settings
	UI struct {
		Theme string `yaml:"theme"` // UI theme (light/dark)
	} `yaml:"ui"`

	logger *zap.Logger
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	c := &Config{}
	
	// Default LLM settings
	c.LLM.Provider = "ollama"
	c.LLM.ModelName = "llama3"
	c.LLM.Endpoint = "http://localhost:11434/api/generate"
	c.LLM.MaxTokens = 2048
	c.LLM.Temperature = 0.7
	
	// Default journal settings
	c.Journal.StorageDir = filepath.Join(homeDir, "momentum_journal", "journals")
	c.Journal.WordCountGoal = 750
	c.Journal.AutosaveInterval = 30
	
	// Default UI settings
	c.UI.Theme = "dark"
	
	return c
}

// ConfigPath returns the path to the config file
func ConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "momentum_journal_config.yaml"
	}
	
	configDir := filepath.Join(homeDir, ".config", "momentum_journal")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return filepath.Join(homeDir, "momentum_journal_config.yaml")
	}
	
	return filepath.Join(configDir, "config.yaml")
}

// Load loads the configuration from file
func Load(logger *zap.Logger) (*Config, error) {
	config := DefaultConfig()
	config.logger = logger
	
	configPath := ConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		logger.Info("Config file not found, creating default config", zap.String("path", configPath))
		if err := config.Save(); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		return config, nil
	}
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	
	logger.Info("Loaded configuration", zap.String("path", configPath))
	return config, nil
}

// Save saves the configuration to file
func (c *Config) Save() error {
	configPath := ConfigPath()
	
	// Ensure directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	if c.logger != nil {
		c.logger.Info("Saved configuration", zap.String("path", configPath))
	}
	
	return nil
}
```

## Step 4: Create Logging Setup

Now, let's create a simple logging package in `internal/logging/logging.go`:

```go
package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new configured logger
func NewLogger(debug bool) (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	
	// Set the log level based on debug flag
	if debug {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	
	// Configure console output
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	
	// Create the logger
	return config.Build()
}

// FileLogger creates a logger that also writes to a file
func FileLogger(logPath string, debug bool) (*zap.Logger, error) {
	// Ensure log directory exists
	logDir := logPath
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}
	
	// Create a configuration
	config := zap.NewProductionConfig()
	
	// Set the log level based on debug flag
	if debug {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	
	// Configure outputs
	config.OutputPaths = []string{"stdout", logPath}
	config.ErrorOutputPaths = []string{"stderr", logPath}
	
	// Create the logger
	return config.Build()
}
```

## Step 5: Create Journal Package

Next, let's implement the journal package for file management in `internal/journal/journal.go`:

```go
package journal

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/yourusername/momentum_journal/internal/config"
	"go.uber.org/zap"
)

// JournalEntry represents a single journal entry
type JournalEntry struct {
	FilePath    string    `json:"file_path"`
	FileName    string    `json:"file_name"`
	CreatedAt   time.Time `json:"created_at"`
	ModifiedAt  time.Time `json:"modified_at"`
	WordCount   int       `json:"word_count"`
	Content     string    `json:"content"`
	IsCompleted bool      `json:"is_completed"` // True if the entry meets the word count goal
}

// Manager handles journal operations
type Manager struct {
	config *config.Config
	logger *zap.Logger
}

// NewManager creates a new journal manager
func NewManager(cfg *config.Config, logger *zap.Logger) (*Manager, error) {
	manager := &Manager{
		config: cfg,
		logger: logger,
	}
	
	// Ensure journal directory exists
	if err := os.MkdirAll(cfg.Journal.StorageDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create journal directory: %w", err)
	}
	
	return manager, nil
}

// CreateEntry creates a new journal entry
func (m *Manager) CreateEntry() (*JournalEntry, error) {
	now := time.Now()
	fileName := fmt.Sprintf("%s-morning-pages.md", now.Format("2006-01-02T15:04"))
	filePath := filepath.Join(m.config.Journal.StorageDir, fileName)
	
	entry := &JournalEntry{
		FilePath:   filePath,
		FileName:   fileName,
		CreatedAt:  now,
		ModifiedAt: now,
		WordCount:  0,
		Content:    "",
	}
	
	// Create initial file with metadata
	if err := m.SaveEntry(entry); err != nil {
		return nil, fmt.Errorf("failed to create journal entry: %w", err)
	}
	
	m.logger.Info("Created new journal entry", 
		zap.String("file", fileName),
		zap.Time("created_at", now))
	
	return entry, nil
}

// SaveEntry saves a journal entry to disk
func (m *Manager) SaveEntry(entry *JournalEntry) error {
	// Update modified time
	entry.ModifiedAt = time.Now()
	
	// Update word count
	entry.WordCount = CountWords(entry.Content)
	
	// Check if completed
	entry.IsCompleted = entry.WordCount >= m.config.Journal.WordCountGoal
	
	// Create directory if it doesn't exist
	dir := filepath.Dir(entry.FilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	// Write the file
	err := os.WriteFile(entry.FilePath, []byte(entry.Content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write journal entry: %w", err)
	}
	
	m.logger.Debug("Saved journal entry", 
		zap.String("file", entry.FileName),
		zap.Int("word_count", entry.WordCount),
		zap.Time("modified_at", entry.ModifiedAt))
	
	return nil
}

// ReadEntry reads a journal entry from disk
func (m *Manager) ReadEntry(filePath string) (*JournalEntry, error) {
	// Read file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}
	
	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	// Create entry
	entry := &JournalEntry{
		FilePath:   filePath,
		FileName:   filepath.Base(filePath),
		CreatedAt:  fileInfo.ModTime(), // This is an approximation, as we don't store creation time
		ModifiedAt: fileInfo.ModTime(),
		Content:    string(content),
		WordCount:  CountWords(string(content)),
	}
	
	// Check if completed
	entry.IsCompleted = entry.WordCount >= m.config.Journal.WordCountGoal
	
	return entry, nil
}

// ListEntries lists all journal entries
func (m *Manager) ListEntries() ([]*JournalEntry, error) {
	entries := []*JournalEntry{}
	
	// Read directory
	files, err := os.ReadDir(m.config.Journal.StorageDir)
	if err != nil {
		if os.IsNotExist(err) {
			return entries, nil
		}
		return nil, fmt.Errorf("failed to read journal directory: %w", err)
	}
	
	// Process each file
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		
		// Check if it's a markdown file
		if !strings.HasSuffix(file.Name(), ".md") {
			continue
		}
		
		// Read the entry
		entry, err := m.ReadEntry(filepath.Join(m.config.Journal.StorageDir, file.Name()))
		if err != nil {
			m.logger.Warn("Failed to read journal entry", 
				zap.String("file", file.Name()),
				zap.Error(err))
			continue
		}
		
		entries = append(entries, entry)
	}
	
	return entries, nil
}

// CountWords counts the number of words in text, ignoring markdown syntax
func CountWords(text string) int {
	// Remove markdown syntax (basic implementation, can be improved)
	// Remove headers, links, images, code blocks, etc.
	markdownSyntax := []string{
		`#.*`,                   // Headers
		`\[.*?\]\(.*?\)`,        // Links
		`!\[.*?\]\(.*?\)`,       // Images
		"```[\\s\\S]*?```",      // Code blocks
		"`.*?`",                 // Inline code
		`\*\*.*?\*\*`,           // Bold
		`\*.*?\*`,               // Italic
		`__.*?__`,               // Bold
		`_.*?_`,                 // Italic
		`~~.*?~~`,               // Strikethrough
		`>\s.*`,                 // Blockquotes
		`- .*`,                  // List items
		`\d+\. .*`,              // Numbered list items
	}
	
	content := text
	for _, pattern := range markdownSyntax {
		re := regexp.MustCompile(pattern)
		content = re.ReplaceAllString(content, "")
	}
	
	// Split by whitespace and count non-empty words
	words := strings.Fields(content)
	return len(words)
}
```

## Step 6: Implement the CLI with Cobra

Now, let's create our CLI application using Cobra. First, let's set up the root command in `cmd/momentum/root.go`:

```go
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/momentum_journal/internal/config"
	"github.com/yourusername/momentum_journal/internal/logging"
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
```

## Step 7: Create the "New" Command

Now, let's create the `new` command in `cmd/momentum/new.go`:

```go
package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/momentum_journal/internal/journal"
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
```

## Step 8: Create the List Command

Let's add a simple command to list existing journal entries in `cmd/momentum/list.go`:

```go
package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourusername/momentum_journal/internal/journal"
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
```

## Step 9: Create Main Entry Point

Finally, let's create the `main.go` file in `cmd/momentum/main.go`:

```go
package main

func main() {
	Execute()
}
```

## Step 10: Testing

To test our implementation so far:

1. Build the application:

```bash
go build -o bin/momentum cmd/momentum/*.go
```

2. Run the application:

```bash
# Show help
./momentum --help

# Create a new journal entry
./bin/momentum new

# List journal entries
./bin/momentum list
```

The application should:
- Create a configuration file at `~/.config/momentum_journal/config.yaml` if it doesn't exist
- Create journal entries in the configured storage directory
- List existing journal entries

## What's Next?

In Phase 2, we'll build the Bubble Tea UI with:
- Two-pane layout
- Vim-like text editor
- Conversation interface with the AI agent

You've successfully completed Phase 1 of the Momentum Journal implementation!

## Troubleshooting

### Common Issues:

1. **Go module errors**: Make sure your module path matches your GitHub username.
2. **Permission errors**: Check if you have write permissions for the configuration and journal directories.
3. **Missing dependencies**: Ensure you've installed all required Go packages.

### Getting Help:

If you encounter any issues, you can:
- Check the application logs for more detailed error information
- Run the application with `--debug` flag for verbose logging
- Review the code for any syntax errors or typos 