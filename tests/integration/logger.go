package integration

import (
	"io"

	zerologger "github.com/amirdashtii/go_auth/infrastructure/logger"
	"github.com/amirdashtii/go_auth/internal/core/ports"
)

// NewTestLogger creates a new logger for testing
func NewTestLogger(serviceName string) ports.Logger {
	loggerConfig := ports.LoggerConfig{
		Level:       "debug",
		Environment: "test",
		ServiceName: serviceName,
		Output:      io.Discard, // Discard logs during tests
	}
	return zerologger.NewZerologLogger(loggerConfig)
}
