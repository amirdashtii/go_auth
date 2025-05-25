package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/errors"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockAuthService struct {
	mock.Mock
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string, fields ...ports.Field) {
	m.Called(msg, fields)
}

func (m *MockLogger) Error(msg string, fields ...ports.Field) {
	m.Called(msg, fields)
}

func (m *MockLogger) Debug(msg string, fields ...ports.Field) {
	m.Called(msg, fields)
}

func (m *MockLogger) Warn(msg string, fields ...ports.Field) {
	m.Called(msg, fields)
}

func (m *MockLogger) Fatal(msg string, fields ...ports.Field) {
	m.Called(msg, fields)
}

func (m *MockLogger) With(fields ...ports.Field) ports.Logger {
	m.Called(fields)
	return m
}

func (m *MockLogger) WithContext(ctx context.Context) ports.Logger {
	m.Called(ctx)
	return m
}

func (m *MockAuthService) Register(ctx context.Context, req *dto.RegisterRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockAuthService) Login(ctx context.Context, req *dto.LoginRequest) (*entities.TokenPair, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.TokenPair), args.Error(1)
}

func (m *MockAuthService) Logout(ctx context.Context, userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, refreshToken string) (*entities.TokenPair, error) {
	args := m.Called(refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.TokenPair), args.Error(1)
}

func (m *MockAuthService) ValidateToken(ctx context.Context, token string, userID string) error {
	args := m.Called(token, userID)
	return args.Error(0)
}

func TestRegisterHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		mockSetup      func(*MockAuthService, *MockLogger)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful registration",
			requestBody: map[string]interface{}{
				"phone_number": "09123456789",
				"password":     "Test123!@#",
			},
			mockSetup: func(m *MockAuthService, l *MockLogger) {
				m.On("Register", mock.AnythingOfType("*dto.RegisterRequest")).Return(nil)
				l.On("Error", mock.Anything, mock.Anything).Return()
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"message": "User registered successfully",
			},
		},
		{
			name: "invalid request format",
			requestBody: map[string]interface{}{
				"phone_number": "invalid",
			},
			mockSetup: func(m *MockAuthService, l *MockLogger) {
				l.On("Error", mock.Anything, mock.Anything).Return()
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"Type": "VALIDATION_ERROR",
					"Message": map[string]interface{}{
						"English": "Invalid request",
						"Persian": "درخواست نامعتبر است",
					},
					"Err": nil,
				},
			},
		},
		{
			name: "registration service error",
			requestBody: map[string]interface{}{
				"phone_number": "09123456789",
				"password":     "Test123!@#",
			},
			mockSetup: func(m *MockAuthService, l *MockLogger) {
				m.On("Register", mock.AnythingOfType("*dto.RegisterRequest")).Return(errors.New(errors.InternalError, "Registration failed", "ثبت نام با خطا مواجه شد", nil))
				l.On("Error", mock.Anything, mock.Anything).Return()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"Type": "INTERNAL_ERROR",
					"Message": map[string]interface{}{
						"English": "Registration failed",
						"Persian": "ثبت نام با خطا مواجه شد",
					},
					"Err": nil,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockAuthService)
			mockLogger := new(MockLogger)
			tt.mockSetup(mockSvc, mockLogger)

			handler := &AuthHTTPHandler{
				svc:    mockSvc,
				logger: mockLogger,
			}
			router := gin.New()
			router.POST("/register", handler.RegisterHandler)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			require.Equal(t, tt.expectedBody, response)
		})
	}
}

func TestLoginHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		mockSetup      func(*MockAuthService, *MockLogger)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful login",
			requestBody: map[string]interface{}{
				"phone_number": "09123456789",
				"password":     "Test123!@#",
			},
			mockSetup: func(m *MockAuthService, l *MockLogger) {
				m.On("Login", mock.AnythingOfType("*dto.LoginRequest")).Return(&entities.TokenPair{
					AccessToken:  "access_token",
					RefreshToken: "refresh_token",
				}, nil)
				l.On("Error", mock.Anything, mock.Anything).Return()
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"tokens": map[string]interface{}{
					"access_token":  "access_token",
					"refresh_token": "refresh_token",
				},
			},
		},
		{
			name: "invalid request format",
			requestBody: map[string]interface{}{
				"phone_number": "invalid",
			},
			mockSetup: func(m *MockAuthService, l *MockLogger) {
				l.On("Error", mock.Anything, mock.Anything).Return()
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"Type": "VALIDATION_ERROR",
					"Message": map[string]interface{}{
						"English": "Invalid request",
						"Persian": "درخواست نامعتبر است",
					},
					"Err": nil,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockAuthService)
			mockLogger := new(MockLogger)
			tt.mockSetup(mockSvc, mockLogger)

			handler := &AuthHTTPHandler{
				svc:    mockSvc,
				logger: mockLogger,
			}
			router := gin.New()
			router.POST("/login", handler.LoginHandler)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			require.Equal(t, tt.expectedBody, response)
		})
	}
}

func TestLogoutHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		mockSetup      func(*MockAuthService, *MockLogger)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "successful logout",
			userID: "user123",
			mockSetup: func(m *MockAuthService, l *MockLogger) {
				m.On("Logout", "user123").Return(nil)
				l.On("Error", mock.Anything, mock.Anything).Return()
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Logged out successfully",
			},
		},
		{
			name:   "missing user ID",
			userID: "",
			mockSetup: func(m *MockAuthService, l *MockLogger) {
				l.On("Error", mock.Anything, mock.Anything).Return()
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"Type": "AUTHENTICATION_ERROR",
					"Message": map[string]interface{}{
						"English": "Authentication required",
						"Persian": "لطفاً ابتدا وارد حساب کاربری خود شوید",
					},
					"Err": nil,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockAuthService)
			mockLogger := new(MockLogger)
			tt.mockSetup(mockSvc, mockLogger)

			handler := &AuthHTTPHandler{
				svc:    mockSvc,
				logger: mockLogger,
			}
			router := gin.New()
			router.POST("/logout", func(c *gin.Context) {
				if tt.userID != "" {
					c.Set("user_id", tt.userID)
				}
				handler.LogoutHandler(c)
			})

			req := httptest.NewRequest(http.MethodPost, "/logout", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			require.Equal(t, tt.expectedBody, response)
		})
	}
}
