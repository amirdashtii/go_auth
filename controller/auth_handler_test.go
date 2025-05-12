package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockAuthService struct {
	mock.Mock
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
		mockSetup      func(*MockAuthService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful registration",
			requestBody: map[string]interface{}{
				"phone_number": "09123456789",
				"password":     "Test123!@#",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("Register", mock.AnythingOfType("*dto.RegisterRequest")).Return(nil)
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
			mockSetup:      func(m *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error":   "Invalid request format",
				"details": "Key: 'RegisterRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag",
			},
		},
		{
			name: "registration service error",
			requestBody: map[string]interface{}{
				"phone_number": "09123456789",
				"password":     "Test123!@#",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("Register", mock.AnythingOfType("*dto.RegisterRequest")).Return(errors.New("registration failed"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error":   "Registration failed",
				"details": "registration failed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockAuthService)
			tt.mockSetup(mockSvc)

			handler := &AuthHTTPHandler{svc: mockSvc}
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
		mockSetup      func(*MockAuthService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful login",
			requestBody: map[string]interface{}{
				"phone_number": "09123456789",
				"password":     "Test123!@#",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("Login", mock.AnythingOfType("*dto.LoginRequest")).Return(&entities.TokenPair{
					AccessToken:  "access_token",
					RefreshToken: "refresh_token",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Login successful",
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
			mockSetup:      func(m *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error":   "Invalid request format",
				"details": "Key: 'LoginRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag",
			},
		},
		{
			name: "login service error",
			requestBody: map[string]interface{}{
				"phone_number": "09123456789",
				"password":     "Test123!@#",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("Login", mock.AnythingOfType("*dto.LoginRequest")).Return(nil, errors.New("login failed"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error":   "Login failed",
				"details": "login failed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockAuthService)
			tt.mockSetup(mockSvc)

			handler := &AuthHTTPHandler{svc: mockSvc}
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
		userID         interface{}
		mockSetup      func(*MockAuthService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "successful logout",
			userID: "user123",
			mockSetup: func(m *MockAuthService) {
				m.On("Logout", "user123").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Logout successful",
			},
		},
		{
			name:           "missing user ID",
			userID:         nil,
			mockSetup:      func(m *MockAuthService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "User ID not found",
			},
		},
		{
			name:           "invalid user ID type",
			userID:         123,
			mockSetup:      func(m *MockAuthService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Invalid user ID type",
			},
		},
		{
			name:   "logout service error",
			userID: "user123",
			mockSetup: func(m *MockAuthService) {
				m.On("Logout", "user123").Return(errors.New("logout failed"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error":   "Logout failed",
				"details": "logout failed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockAuthService)
			tt.mockSetup(mockSvc)

			handler := &AuthHTTPHandler{svc: mockSvc}
			router := gin.New()
			router.POST("/logout", func(c *gin.Context) {
				if tt.userID != nil {
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

func TestRefreshTokenHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		mockSetup      func(*MockAuthService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful token refresh",
			requestBody: map[string]interface{}{
				"refresh_token": "valid_refresh_token",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("RefreshToken", "valid_refresh_token").Return(&entities.TokenPair{
					AccessToken:  "new_access_token",
					RefreshToken: "new_refresh_token",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Token refreshed successfully",
				"tokens": map[string]interface{}{
					"access_token":  "new_access_token",
					"refresh_token": "new_refresh_token",
				},
			},
		},
		{
			name: "invalid request format",
			requestBody: map[string]interface{}{
				"refresh_token": "",
			},
			mockSetup:      func(m *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error":   "Invalid request format",
				"details": "Key: 'RefreshTokenRequest.RefreshToken' Error:Field validation for 'RefreshToken' failed on the 'required' tag",
			},
		},
		{
			name: "refresh token service error",
			requestBody: map[string]interface{}{
				"refresh_token": "invalid_token",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("RefreshToken", "invalid_token").Return(nil, errors.New("token refresh failed"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error":   "Token refresh failed",
				"details": "token refresh failed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockAuthService)
			tt.mockSetup(mockSvc)

			handler := &AuthHTTPHandler{svc: mockSvc}
			router := gin.New()
			router.POST("/refresh-token", handler.RefreshTokenHandler)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/refresh-token", bytes.NewBuffer(body))
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
