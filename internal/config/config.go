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
		Provider    string  `yaml:"provider"`    // "ollama" or "openrouter"
		APIKey      string  `yaml:"api_key"`     // API key for openrouter
		ModelName   string  `yaml:"model_name"`  // Model to use e.g., "llama3" for Ollama
		Endpoint    string  `yaml:"endpoint"`    // API endpoint
		MaxTokens   int     `yaml:"max_tokens"`  // Maximum tokens for response
		Temperature float64 `yaml:"temperature"` // Temperature for generation
	} `yaml:"llm"`

	// Journal settings
	Journal struct {
		StorageDir       string `yaml:"storage_dir"`       // Directory to store journal files
		WordCountGoal    int    `yaml:"word_count_goal"`   // Default 750 words (3 pages)
		AutosaveInterval int    `yaml:"autosave_interval"` // Autosave interval in seconds
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
