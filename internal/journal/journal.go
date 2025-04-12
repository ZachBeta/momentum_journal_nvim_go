package journal

import (
	"fmt"
	"os"
	"path/filepath"

	//	"regexp" // Removed as we are simplifying CountWords
	"strings"
	"time"

	"github.com/ZachBeta/momentum_journal_nvim_go/internal/config" // Adjusted import path
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

// CountWords counts the number of words in text using basic tokenization.
func CountWords(text string) int {
	// Split by whitespace and count non-empty words
	words := strings.Fields(text)
	return len(words)
}
