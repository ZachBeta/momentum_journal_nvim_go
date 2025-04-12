package logging

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new configured logger
func NewLogger(debug bool) (*zap.Logger, error) {
	// Use development config for more console-friendly output
	config := zap.NewDevelopmentConfig()

	// Set the log level based on debug flag
	if debug {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// Keep console output configuration (colors, time format)
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Keep colored levels
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// OutputPaths and ErrorOutputPaths are typically stderr for DevelopmentConfig
	// Let's stick with the DevelopmentConfig defaults for now.
	// config.OutputPaths = []string{"stdout"}
	// config.ErrorOutputPaths = []string{"stderr"}

	// Create the logger
	return config.Build()
}

// FileLogger creates a logger that also writes to a file
func FileLogger(logPath string, debug bool) (*zap.Logger, error) {
	// Ensure log directory exists
	logDir := filepath.Dir(logPath) // Corrected: use filepath.Dir to get the directory
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	// Create a configuration - let's use Development for file too for consistency?
	// Or keep Production? Let's keep Production for file logging for now.
	config := zap.NewProductionConfig()

	// Set the log level based on debug flag
	if debug {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// Configure outputs - writing JSON to file, nothing to stdout/stderr from here
	config.OutputPaths = []string{logPath}
	config.ErrorOutputPaths = []string{logPath}

	// Create the logger
	return config.Build()
}
