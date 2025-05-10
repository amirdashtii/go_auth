package logger

import (
	"os"
	"path/filepath"
	"time"

	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/rs/zerolog"
)

// NewFileLogger creates a new logger that writes to a file
func NewFileLogger(config ports.LoggerConfig) (ports.Logger, error) {
	// Create logs directory if it doesn't exist
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return nil, err
	}

	// Create log file with date in name
	currentTime := time.Now().Format("2006-01-02")
	logFileName := filepath.Join(logsDir, "app-"+currentTime+".log")
	
	// Open log file
	logFile, err := os.OpenFile(
		logFileName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return nil, err
	}

	// Create multi-writer to write to both file and stdout
	multi := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout}, logFile)

	// Create logger
	logger := zerolog.New(multi).
		With().
		Timestamp().
		Str("service", config.ServiceName).
		Str("env", config.Environment).
		Logger()

	return &zerologLogger{
		logger: logger,
		config: config,
	}, nil
} 