package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/gin-gonic/gin"
)

// sensitiveFields contains the list of fields that should be masked in logs
var sensitiveFields = []string{
	"password",
	"token",
	"secret",
	"key",
	"authorization",
	"cookie",
	"session",
}

// responseWriter is a custom response writer that captures the response body
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// LoggerMiddleware creates a new logger middleware
func LoggerMiddleware(logger ports.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context() // Get the context
		l := logger.WithContext(ctx) // Create a new logger with the context

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Read request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Log request
		requestFields := []ports.Field{
			ports.F("client_ip", c.ClientIP()),
			ports.F("method", c.Request.Method),
			ports.F("path", path),
		}

		if raw != "" {
			requestFields = append(requestFields, ports.F("query", raw))
		}

		if len(requestBody) > 0 {
			maskedBody := maskSensitiveData(requestBody)
			requestFields = append(requestFields, ports.F("request_body", string(maskedBody)))
		}

		l.Info("Incoming Request", requestFields...) // Use l instead of logger

		// Create custom response writer
		responseBody := &bytes.Buffer{}
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:          responseBody,
		}
		c.Writer = writer

		// Process request
		c.Next()

		// Log response
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		responseFields := []ports.Field{
			ports.F("status", statusCode),
			ports.F("latency", latency),
			ports.F("path", path),
		}

		if responseBody.Len() > 0 {
			maskedResponse := maskSensitiveData(responseBody.Bytes())
			responseFields = append(responseFields, ports.F("response_body", string(maskedResponse)))
		}

		if errorMessage != "" {
			responseFields = append(responseFields, ports.F("error", errorMessage))
		}

		// Log response based on status code
		switch {
		case statusCode >= 500:
			l.Error("Server Response", responseFields...) // Use l instead of logger
		case statusCode >= 400:
			l.Warn("Client Response", responseFields...) // Use l instead of logger
		default:
			l.Info("Server Response", responseFields...) // Use l instead of logger
		}
	}
}

// maskSensitiveData masks sensitive fields in the request/response body
func maskSensitiveData(body []byte) []byte {
	// Try to parse as JSON
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		// If not JSON, return as is
		return body
	}

	// Mask sensitive fields
	maskSensitiveFields(data)

	// Convert back to JSON
	maskedBody, err := json.Marshal(data)
	if err != nil {
		return body
	}

	return maskedBody
}

// maskSensitiveFields recursively masks sensitive fields in the data
func maskSensitiveFields(data map[string]interface{}) {
	for key, value := range data {
		// Check if the key is sensitive
		if isSensitiveField(key) {
			data[key] = "********"
			continue
		}

		// If the value is a map, recursively mask its fields
		if nestedMap, ok := value.(map[string]interface{}); ok {
			maskSensitiveFields(nestedMap)
		}

		// If the value is a slice of maps, recursively mask each map's fields
		if nestedSlice, ok := value.([]interface{}); ok {
			for _, item := range nestedSlice {
				if nestedMap, ok := item.(map[string]interface{}); ok {
					maskSensitiveFields(nestedMap)
				}
			}
		}
	}
}

// isSensitiveField checks if a field name is in the sensitive fields list
func isSensitiveField(field string) bool {
	field = strings.ToLower(field)
	for _, sensitive := range sensitiveFields {
		if strings.Contains(field, sensitive) {
			return true
		}
	}
	return false
} 