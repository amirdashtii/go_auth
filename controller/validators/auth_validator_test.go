package validators

import (
	"context"
	"os"
	"testing"

	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type mockLogger struct{}

func (m *mockLogger) Info(msg string, fields ...ports.Field)  {}
func (m *mockLogger) Error(msg string, fields ...ports.Field) {}
func (m *mockLogger) Debug(msg string, fields ...ports.Field) {}
func (m *mockLogger) Warn(msg string, fields ...ports.Field)  {}
func (m *mockLogger) Fatal(msg string, fields ...ports.Field) { os.Exit(1) }
func (m *mockLogger) With(fields ...ports.Field) ports.Logger { return m }
func (m *mockLogger) WithContext(ctx context.Context) ports.Logger { return m }

var testLogger = &mockLogger{}

func TestValidatePhoneNumber(t *testing.T) {
	tests := []struct {
		name     string
		phone    string
		expected bool
	}{
		{
			name:     "valid phone number",
			phone:    "09123456789",
			expected: true,
		},
		{
			name:     "invalid phone number - wrong prefix",
			phone:    "08123456789",
			expected: false,
		},
		{
			name:     "invalid phone number - too short",
			phone:    "0912345678",
			expected: false,
		},
		{
			name:     "invalid phone number - too long",
			phone:    "091234567890",
			expected: false,
		},
		{
			name:     "invalid phone number - contains letters",
			phone:    "0912345678a",
			expected: false,
		},
	}

	v := validator.New()
	v.RegisterValidation("phone", ValidatePhoneNumber)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Var(tt.phone, "phone")
			if tt.expected {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestValidateAuthPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{
			name:     "valid password",
			password: "Test1234",
			expected: true,
		},
		{
			name:     "invalid password - too short",
			password: "Test123",
			expected: false,
		},
		{
			name:     "invalid password - no uppercase",
			password: "test1234",
			expected: false,
		},
		{
			name:     "invalid password - no lowercase",
			password: "TEST1234",
			expected: false,
		},
		{
			name:     "invalid password - no numbers",
			password: "TestTest",
			expected: false,
		},
	}

	v := validator.New()
	v.RegisterValidation("password", ValidateAuthPassword)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Var(tt.password, "password")
			if tt.expected {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestValidateRegisterRequest(t *testing.T) {
	tests := []struct {
		name    string
		request *dto.RegisterRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: &dto.RegisterRequest{
				PhoneNumber: "09123456789",
				Password:    "Test1234",
			},
			wantErr: false,
		},
		{
			name: "invalid phone",
			request: &dto.RegisterRequest{
				PhoneNumber: "08123456789",
				Password:    "Test1234",
			},
			wantErr: true,
		},
		{
			name: "invalid password",
			request: &dto.RegisterRequest{
				PhoneNumber: "09123456789",
				Password:    "test",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRegisterRequest(tt.request, testLogger)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateLoginRequest(t *testing.T) {
	tests := []struct {
		name    string
		request *dto.LoginRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: &dto.LoginRequest{
				PhoneNumber: "09123456789",
				Password:    "Test1234",
			},
			wantErr: false,
		},
		{
			name: "invalid phone",
			request: &dto.LoginRequest{
				PhoneNumber: "08123456789",
				Password:    "Test1234",
			},
			wantErr: true,
		},
		{
			name: "invalid password",
			request: &dto.LoginRequest{
				PhoneNumber: "09123456789",
				Password:    "test",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLoginRequest(tt.request, testLogger)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateRefreshTokenRequest(t *testing.T) {
	tests := []struct {
		name    string
		request *dto.RefreshTokenRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: &dto.RefreshTokenRequest{
				RefreshToken: "valid.refresh.token",
			},
			wantErr: false,
		},
		{
			name: "empty refresh token",
			request: &dto.RefreshTokenRequest{
				RefreshToken: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRefreshTokenRequest(tt.request, testLogger)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}