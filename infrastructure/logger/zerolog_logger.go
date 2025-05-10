package logger

import (
	"context"
	"time"

	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/rs/zerolog"
)

type zerologLogger struct {
	logger zerolog.Logger
	config ports.LoggerConfig
}

// NewZerologLogger creates a new zerolog logger
func NewZerologLogger(config ports.LoggerConfig) ports.Logger {
	// Set up zerolog
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(getLogLevel(config.Level))

	// Create logger
	logger := zerolog.New(config.Output).
		With().
		Timestamp().
		Str("service", config.ServiceName).
		Str("env", config.Environment).
		Logger()

	return &zerologLogger{
		logger: logger,
		config: config,
	}
}

func (l *zerologLogger) Debug(msg string, fields ...ports.Field) {
	l.logger.Debug().Fields(convertFields(fields)).Msg(msg)
}

func (l *zerologLogger) Info(msg string, fields ...ports.Field) {
	l.logger.Info().Fields(convertFields(fields)).Msg(msg)
}

func (l *zerologLogger) Warn(msg string, fields ...ports.Field) {
	l.logger.Warn().Fields(convertFields(fields)).Msg(msg)
}

func (l *zerologLogger) Error(msg string, fields ...ports.Field) {
	l.logger.Error().Fields(convertFields(fields)).Msg(msg)
}

func (l *zerologLogger) Fatal(msg string, fields ...ports.Field) {
	l.logger.Fatal().Fields(convertFields(fields)).Msg(msg)
}

func (l *zerologLogger) With(fields ...ports.Field) ports.Logger {
	return &zerologLogger{
		logger: l.logger.With().Fields(convertFields(fields)).Logger(),
		config: l.config,
	}
}

func (l *zerologLogger) WithContext(ctx context.Context) ports.Logger {
	return &zerologLogger{
		logger: l.logger.With().Fields(convertFields(getContextFields(ctx))).Logger(),
		config: l.config,
	}
}

// Helper functions

func getLogLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

func convertFields(fields []ports.Field) map[string]interface{} {
	result := make(map[string]interface{})
	for _, field := range fields {
		result[field.Key] = field.Value
	}
	return result
}

func getContextFields(ctx context.Context) []ports.Field {

	return nil
} 