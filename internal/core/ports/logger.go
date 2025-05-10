package ports

import (
	"context"
	"io"
)

// Logger interface defines the logging contract
type Logger interface {
	// Debug logs a debug message
	Debug(msg string, fields ...Field)
	// Info logs an info message
	Info(msg string, fields ...Field)
	// Warn logs a warning message
	Warn(msg string, fields ...Field)
	// Error logs an error message
	Error(msg string, fields ...Field)
	// Fatal logs a fatal message and then calls os.Exit(1)
	Fatal(msg string, fields ...Field)
	// With returns a new logger with the given fields added to its context
	With(fields ...Field) Logger
	// WithContext returns a new logger with the given context
	WithContext(ctx context.Context) Logger
}

// Field represents a key-value pair for structured logging
type Field struct {
	Key   string
	Value interface{}
}

// F creates a new Field
func F(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// LoggerConfig holds the configuration for the logger
type LoggerConfig struct {
	// Output is the writer where logs will be written
	Output io.Writer
	// Level is the minimum log level that will be logged
	Level string
	// Environment is the current environment (development, production, etc.)
	Environment string
	// ServiceName is the name of the service
	ServiceName string
} 