// Package service contains the implementation of the authentication service
package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/amirdashtii/go_auth/config"
	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/service/mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// TestRegister tests the user registration functionality
func TestRegister(t *testing.T) {
	// Initialize mock repositories
	mockAuthRepo := new(mocks.AuthRepository)
	mockRedisRepo := new(mocks.InMemoryRespositoryContracts)

	// Create service instance with mock repositories
	service := &AuthService{
		db:    mockAuthRepo,
		redis: mockRedisRepo,
	}

	// Create registration request
	req := &dto.RegisterRequest{
		PhoneNumber: "09123456789",
		Password:    "password123",
	}

	// Set up mock expectations
	mockAuthRepo.On("Create", mock.Anything).Return(nil).Once()

	// Execute registration
	err := service.Register(req)

	// Verify results
	assert.NoError(t, err)
	mockAuthRepo.AssertExpectations(t)
}

// TestRegister_DuplicateUser tests registration with a duplicate phone number
func TestRegister_DuplicateUser(t *testing.T) {
	// Initialize mock repositories
	mockAuthRepo := new(mocks.AuthRepository)
	mockRedisRepo := new(mocks.InMemoryRespositoryContracts)

	// Create service instance with mock repositories
	service := &AuthService{
		db:    mockAuthRepo,
		redis: mockRedisRepo,
	}

	// Create registration request
	req := &dto.RegisterRequest{
		PhoneNumber: "09123456789",
		Password:    "password123",
	}

	// Set up mock expectations
	// Expect Create to be called once and return a duplicate user error
	mockAuthRepo.On("Create", mock.Anything).Return(fmt.Errorf("user with this phone number already exists")).Once()

	// Execute registration
	err := service.Register(req)

	// Verify results
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
	mockAuthRepo.AssertExpectations(t)
}

// TestLogin tests the user login functionality
func TestLogin(t *testing.T) {
	// Initialize mock repositories
	mockAuthRepo := new(mocks.AuthRepository)
	mockRedisRepo := new(mocks.InMemoryRespositoryContracts)

	// Create service instance with mock repositories
	service := &AuthService{
		db:    mockAuthRepo,
		redis: mockRedisRepo,
	}

	// Create test user
	userID := uuid.New()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &entities.User{
		ID:          userID,
		PhoneNumber: "09123456789",
		Password:    string(hashedPassword),
		Role:        entities.UserRole,
	}

	// Create login request
	req := &dto.LoginRequest{
		PhoneNumber: "09123456789",
		Password:    "password123",
	}

	// Set up mock expectations
	mockAuthRepo.On("FindUserByPhoneNumber", &req.PhoneNumber).Return(user, nil).Once()
	mockRedisRepo.On("AddToken", mock.Anything, mock.Anything, mock.Anything).Return(nil).Twice()

	// Execute login
	tokens, err := service.Login(req)

	// Verify results
	assert.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
	mockAuthRepo.AssertExpectations(t)
	mockRedisRepo.AssertExpectations(t)
}

// TestLogin_InvalidPassword tests login with an incorrect password
func TestLogin_InvalidPassword(t *testing.T) {
	// Initialize mock repositories
	mockAuthRepo := new(mocks.AuthRepository)
	mockRedisRepo := new(mocks.InMemoryRespositoryContracts)

	// Create service instance with mock repositories
	service := &AuthService{
		db:    mockAuthRepo,
		redis: mockRedisRepo,
	}

	// Create test user with correct password
	userID := uuid.New()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct_password"), bcrypt.DefaultCost)
	user := &entities.User{
		ID:          userID,
		PhoneNumber: "09123456789",
		Password:    string(hashedPassword),
		Status:      entities.Active,
		Role:        entities.UserRole,
	}

	// Create login request with wrong password
	loginReq := &dto.LoginRequest{
		PhoneNumber: "09123456789",
		Password:    "wrong_password",
	}

	// Set up mock expectations
	// Expect FindUserByPhoneNumber to be called once and return the test user
	mockAuthRepo.On("FindUserByPhoneNumber", &loginReq.PhoneNumber).Return(user, nil).Once()

	// Execute login
	_, err := service.Login(loginReq)

	// Verify results
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid password")
	mockAuthRepo.AssertExpectations(t)
}

// TestLogin_DeactivatedUser tests login for a deactivated user
func TestLogin_DeactivatedUser(t *testing.T) {
	// Initialize mock repositories
	mockAuthRepo := new(mocks.AuthRepository)
	mockRedisRepo := new(mocks.InMemoryRespositoryContracts)

	// Create service instance with mock repositories
	service := &AuthService{
		db:    mockAuthRepo,
		redis: mockRedisRepo,
	}

	// Create test user with deactivated status
	userID := uuid.New()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &entities.User{
		ID:          userID,
		PhoneNumber: "09123456789",
		Password:    string(hashedPassword),
		Status:      entities.Deactivated,
		Role:        entities.UserRole,
	}

	// Create login request
	loginReq := &dto.LoginRequest{
		PhoneNumber: "09123456789",
		Password:    "password123",
	}

	// Set up mock expectations
	// Expect FindUserByPhoneNumber to be called once and return the deactivated user
	mockAuthRepo.On("FindUserByPhoneNumber", &loginReq.PhoneNumber).Return(user, nil).Once()

	// Execute login
	_, err := service.Login(loginReq)

	// Verify results
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user is deactivated")
	mockAuthRepo.AssertExpectations(t)
}

// TestLogin_DeletedUser tests login for a deleted user
func TestLogin_DeletedUser(t *testing.T) {
	// Initialize mock repositories
	mockAuthRepo := new(mocks.AuthRepository)
	mockRedisRepo := new(mocks.InMemoryRespositoryContracts)

	// Create service instance with mock repositories
	service := &AuthService{
		db:    mockAuthRepo,
		redis: mockRedisRepo,
	}

	// Create test user with deleted status
	userID := uuid.New()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &entities.User{
		ID:          userID,
		PhoneNumber: "09123456789",
		Password:    string(hashedPassword),
		Status:      entities.Deleted,
		Role:        entities.UserRole,
	}

	// Create login request
	loginReq := &dto.LoginRequest{
		PhoneNumber: "09123456789",
		Password:    "password123",
	}

	// Set up mock expectations
	// Expect FindUserByPhoneNumber to be called once and return the deleted user
	mockAuthRepo.On("FindUserByPhoneNumber", &loginReq.PhoneNumber).Return(user, nil).Once()

	// Execute login
	_, err := service.Login(loginReq)

	// Verify results
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	mockAuthRepo.AssertExpectations(t)
}

// TestLogin_RedisError tests login when Redis operations fail
func TestLogin_RedisError(t *testing.T) {
	// Initialize mock repositories
	mockAuthRepo := new(mocks.AuthRepository)
	mockRedisRepo := new(mocks.InMemoryRespositoryContracts)

	// Create service instance with mock repositories
	service := &AuthService{
		db:    mockAuthRepo,
		redis: mockRedisRepo,
	}

	// Create test user
	userID := uuid.New()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &entities.User{
		ID:          userID,
		PhoneNumber: "09123456789",
		Password:    string(hashedPassword),
		Status:      entities.Active,
		Role:        entities.UserRole,
	}

	// Create login request
	loginReq := &dto.LoginRequest{
		PhoneNumber: "09123456789",
		Password:    "password123",
	}

	// Set up mock expectations
	// Expect FindUserByPhoneNumber to be called once and return the test user
	mockAuthRepo.On("FindUserByPhoneNumber", &loginReq.PhoneNumber).Return(user, nil).Once()
	// Expect AddToken to be called once and return a Redis error
	mockRedisRepo.On("AddToken", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("redis error")).Once()

	// Execute login
	_, err := service.Login(loginReq)

	// Verify results
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to store access token in redis")
	mockAuthRepo.AssertExpectations(t)
	mockRedisRepo.AssertExpectations(t)
}

// TestLogout tests the user logout functionality
func TestLogout(t *testing.T) {
	// Initialize mock repositories
	mockAuthRepo := new(mocks.AuthRepository)
	mockRedisRepo := new(mocks.InMemoryRespositoryContracts)

	// Create service instance with mock repositories
	service := &AuthService{
		db:    mockAuthRepo,
		redis: mockRedisRepo,
	}

	// Create test user ID
	userID := uuid.New()

	// Set up mock expectations
	// Expect RemoveToken to be called twice - once for access token and once for refresh token
	mockRedisRepo.On("RemoveToken", userID.String()+":access").Return(nil).Once()
	mockRedisRepo.On("RemoveToken", userID.String()+":refresh").Return(nil).Once()

	// Execute logout
	err := service.Logout(userID.String())

	// Verify results
	assert.NoError(t, err)
	mockRedisRepo.AssertExpectations(t)
}

// TestLogout_RedisError tests logout when Redis operations fail
func TestLogout_RedisError(t *testing.T) {
	// Initialize mock repositories
	mockAuthRepo := new(mocks.AuthRepository)
	mockRedisRepo := new(mocks.InMemoryRespositoryContracts)

	// Create service instance with mock repositories
	service := &AuthService{
		db:    mockAuthRepo,
		redis: mockRedisRepo,
	}

	// Create test user ID
	userID := uuid.New()

	// Set up mock expectations
	// Expect RemoveToken to be called once for access token and return a Redis error
	mockRedisRepo.On("RemoveToken", userID.String()+":access").Return(fmt.Errorf("redis error")).Once()

	// Execute logout
	err := service.Logout(userID.String())

	// Verify results
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "redis error")
	mockRedisRepo.AssertExpectations(t)
}

// TestRefreshToken tests the token refresh functionality
func TestRefreshToken(t *testing.T) {
	// Initialize mock repositories
	mockAuthRepo := new(mocks.AuthRepository)
	mockRedisRepo := new(mocks.InMemoryRespositoryContracts)

	// Create service instance with mock repositories
	service := &AuthService{
		db:    mockAuthRepo,
		redis: mockRedisRepo,
	}

	// Create test user
	userID := uuid.New()
	user := &entities.User{
		ID:          userID,
		PhoneNumber: "09123456789",
		Role:        entities.UserRole,
	}

	// Create refresh token
	cfg, _ := config.LoadConfig()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    userID.String(),
		"role":       user.Role,
		"token_type": "refresh",
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	})
	refreshToken, _ := token.SignedString([]byte(cfg.JWT.Secret))

	// Set up mock expectations
	mockRedisRepo.On("FindToken", userID.String()+":refresh").Return(refreshToken, nil).Once()
	mockAuthRepo.On("FindUserByID", userID).Return(user, nil).Once()
	mockRedisRepo.On("RemoveToken", userID.String()+":access").Return(nil).Once()
	mockRedisRepo.On("RemoveToken", userID.String()+":refresh").Return(nil).Once()
	mockRedisRepo.On("AddToken", mock.Anything, mock.Anything, mock.Anything).Return(nil).Twice()

	// Execute refresh token
	tokens, err := service.RefreshToken(refreshToken)

	// Verify results
	assert.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
	mockAuthRepo.AssertExpectations(t)
	mockRedisRepo.AssertExpectations(t)
}

// TestRefreshToken_ExpiredToken tests refresh with an expired token
func TestRefreshToken_ExpiredToken(t *testing.T) {
	mockAuthRepo := new(mocks.AuthRepository)
	mockRedisRepo := new(mocks.InMemoryRespositoryContracts)

	service := &AuthService{
		db:    mockAuthRepo,
		redis: mockRedisRepo,
	}

	userID := uuid.New()
	user := &entities.User{
		Role: entities.UserRole,
	}

	config, _ := config.LoadConfig()
	claims := jwt.MapClaims{
		"user_id":    userID,
		"role":       user.Role,
		"token_type": "refresh",
		"exp":        time.Now().Add(-1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, _ := token.SignedString([]byte(config.JWT.Secret))

	_, err := service.RefreshToken(refreshToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid refresh token")
}

// TestRefreshToken_RedisError tests refresh when Redis operations fail
func TestRefreshToken_RedisError(t *testing.T) {
	// Initialize mock repositories
	mockAuthRepo := new(mocks.AuthRepository)
	mockRedisRepo := new(mocks.InMemoryRespositoryContracts)

	// Create service instance with mock repositories
	service := &AuthService{
		db:    mockAuthRepo,
		redis: mockRedisRepo,
	}

	// Create test user
	userID := uuid.New()
	user := &entities.User{
		ID:          userID,
		PhoneNumber: "09123456789",
		Status:      entities.Active,
		Role:        entities.UserRole,
	}

	// Create valid refresh token
	config, _ := config.LoadConfig()
	claims := jwt.MapClaims{
		"user_id":    userID,
		"role":       user.Role,
		"token_type": "refresh",
		"exp":        time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, _ := token.SignedString([]byte(config.JWT.Secret))

	// Set up mock expectations
	// Expect FindUserByID to be called once and return the test user
	mockAuthRepo.On("FindUserByID", userID).Return(user, nil).Once()
	// Expect FindToken to be called once and return the refresh token
	mockRedisRepo.On("FindToken", userID.String()+":refresh").Return(refreshToken, nil).Once()
	// Expect RemoveToken to be called twice - once for access token and once for refresh token
	mockRedisRepo.On("RemoveToken", userID.String()+":access").Return(nil).Once()
	mockRedisRepo.On("RemoveToken", userID.String()+":refresh").Return(nil).Once()
	// Expect AddToken to be called once and return a Redis error
	mockRedisRepo.On("AddToken", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("redis error")).Once()

	// Execute refresh token
	_, err := service.RefreshToken(refreshToken)

	// Verify results
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to store access token in redis")
	mockAuthRepo.AssertExpectations(t)
	mockRedisRepo.AssertExpectations(t)
}

// TestRefreshToken_InvalidToken tests refresh with an invalid token
func TestRefreshToken_InvalidToken(t *testing.T) {
	// Initialize mock repositories
	mockAuthRepo := new(mocks.AuthRepository)
	mockRedisRepo := new(mocks.InMemoryRespositoryContracts)

	// Create service instance with mock repositories
	service := &AuthService{
		db:    mockAuthRepo,
		redis: mockRedisRepo,
	}

	// Create invalid refresh token
	refreshToken := "invalid_refresh_token"

	// Execute refresh token
	_, err := service.RefreshToken(refreshToken)

	// Verify results
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid refresh token")
}

// TestRefreshToken_UserNotFound tests refresh for a non-existent user
func TestRefreshToken_UserNotFound(t *testing.T) {
	// Initialize mock repositories
	mockAuthRepo := new(mocks.AuthRepository)
	mockRedisRepo := new(mocks.InMemoryRespositoryContracts)

	// Create service instance with mock repositories
	service := &AuthService{
		db:    mockAuthRepo,
		redis: mockRedisRepo,
	}

	// Create test user ID
	userID := uuid.New()

	// Create valid refresh token
	config, _ := config.LoadConfig()
	claims := jwt.MapClaims{
		"user_id":    userID,
		"role":       entities.UserRole,
		"token_type": "refresh",
		"exp":        time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, _ := token.SignedString([]byte(config.JWT.Secret))

	// Set up mock expectations
	// Expect FindUserByID to be called and return user not found error
	mockAuthRepo.On("FindUserByID", userID).Return(nil, fmt.Errorf("user not found")).Once()

	// Execute refresh token
	_, err := service.RefreshToken(refreshToken)

	// Verify results
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	mockAuthRepo.AssertExpectations(t)
	mockRedisRepo.AssertExpectations(t)
}

// TestRefreshToken_TokenMismatch tests refresh when the token doesn't match
func TestRefreshToken_TokenMismatch(t *testing.T) {
	// Initialize mock repositories
	mockAuthRepo := new(mocks.AuthRepository)
	mockRedisRepo := new(mocks.InMemoryRespositoryContracts)

	// Create service instance with mock repositories
	service := &AuthService{
		db:    mockAuthRepo,
		redis: mockRedisRepo,
	}

	// Create test user
	userID := uuid.New()
	user := &entities.User{
		ID:          userID,
		PhoneNumber: "09123456789",
		Status:      entities.Active,
		Role:        entities.UserRole,
	}

	// Create valid refresh token
	config, _ := config.LoadConfig()
	claims := jwt.MapClaims{
		"user_id":    userID,
		"role":       user.Role,
		"token_type": "refresh",
		"exp":        time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, _ := token.SignedString([]byte(config.JWT.Secret))

	// Create a different token to be returned by FindToken
	storedToken := "different_refresh_token"

	// Set up mock expectations
	// Expect FindUserByID to be called once and return the test user
	mockAuthRepo.On("FindUserByID", userID).Return(user, nil).Once()
	// Expect FindToken to be called once and return a different token
	mockRedisRepo.On("FindToken", userID.String()+":refresh").Return(storedToken, nil).Once()

	// Execute refresh token
	_, err := service.RefreshToken(refreshToken)

	// Verify results
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "refresh token does not match stored token")
	mockAuthRepo.AssertExpectations(t)
	mockRedisRepo.AssertExpectations(t)
}

// TestValidateToken tests the token validation functionality
func TestValidateToken(t *testing.T) {
	// Initialize mock repositories
	mockAuthRepo := new(mocks.AuthRepository)
	mockRedisRepo := new(mocks.InMemoryRespositoryContracts)

	// Create service instance with mock repositories
	service := &AuthService{
		db:    mockAuthRepo,
		redis: mockRedisRepo,
	}

	// Create test user
	userID := uuid.New()

	// Create access token
	cfg, _ := config.LoadConfig()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"type":    "access",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	accessToken, _ := token.SignedString([]byte(cfg.JWT.Secret))


	// Set up mock expectations
	mockRedisRepo.On("FindToken", userID.String()+":access").Return(accessToken, nil).Once()

	// Execute validate token
	err := service.ValidateToken(userID.String(), accessToken)

	// Verify results
	assert.NoError(t, err)
	mockAuthRepo.AssertExpectations(t)
	mockRedisRepo.AssertExpectations(t)
}

// TestValidateToken_InvalidToken tests token validation with an invalid token
func TestValidateToken_InvalidToken(t *testing.T) {
	// Initialize mock repositories
	mockAuthRepo := new(mocks.AuthRepository)
	mockRedisRepo := new(mocks.InMemoryRespositoryContracts)

	// Create service instance with mock repositories
	service := &AuthService{
		db:    mockAuthRepo,
		redis: mockRedisRepo,
	}

	// Create test user
	userID := uuid.New()

	// Set up mock expectations
	mockRedisRepo.On("FindToken", userID.String()+":access").Return("", fmt.Errorf("token not found")).Once()

	// Execute validate token with invalid token
	err := service.ValidateToken(userID.String(), "invalid_token")

	// Verify results
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find token")
	mockRedisRepo.AssertExpectations(t)
}
// ... existing code ...

// TestValidateToken_TokenNotFound tests token validation when token is not found in Redis
func TestValidateToken_TokenNotFound(t *testing.T) {
	// Initialize mock repositories
	mockAuthRepo := new(mocks.AuthRepository)
	mockRedisRepo := new(mocks.InMemoryRespositoryContracts)

	// Create service instance with mock repositories
	service := &AuthService{
		db:    mockAuthRepo,
		redis: mockRedisRepo,
	}

	// Create test user
	userID := uuid.New()

	// Create access token
	cfg, _ := config.LoadConfig()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"type":    "access",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	accessToken, _ := token.SignedString([]byte(cfg.JWT.Secret))

	// Set up mock expectations
	mockRedisRepo.On("FindToken", userID.String()+":access").Return("", fmt.Errorf("token not found")).Once()
	
	// Execute validate token
	err := service.ValidateToken(userID.String(), accessToken)

	// Verify results
	assert.Error(t, err)
	mockAuthRepo.AssertExpectations(t)
	mockRedisRepo.AssertExpectations(t)
}


// TestNewAuthService tests the creation of a new auth service
func TestNewAuthService(t *testing.T) {
	// Initialize mock repositories
	mockAuthRepo := new(mocks.AuthRepository)
	mockRedisRepo := new(mocks.InMemoryRespositoryContracts)

	// Create service instance with mock repositories
	service := &AuthService{
		db:    mockAuthRepo,
		redis: mockRedisRepo,
	}

	// Verify service instance
	assert.NotNil(t, service)
	assert.Equal(t, mockAuthRepo, service.db)
	assert.Equal(t, mockRedisRepo, service.redis)
}

// TestParseAndValidateToken_ExpiredToken tests token parsing with expired token
func TestParseAndValidateToken_ExpiredToken(t *testing.T) {
	// Create service instance with mock repositories
	service := &AuthService{}

	// Create test user
	userID := uuid.New()

	// Create expired token
	cfg, _ := config.LoadConfig()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"type":    "access",
		"exp":     time.Now().Add(-time.Hour).Unix(),
	})
	expiredToken, _ := token.SignedString([]byte(cfg.JWT.Secret))

	// Execute validate token
	_, err := service.parseAndValidateToken(expiredToken, "access")

	// Verify results
	assert.Error(t, err)

}

// TestParseAndValidateToken_InvalidSignature tests token parsing with invalid signature
func TestParseAndValidateToken_InvalidSignature(t *testing.T) {
	// Create service instance with mock repositories
	service := &AuthService{}

	// Create test user
	userID := uuid.New()

	// Create token with invalid signature
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"type":    "access",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})

	// Create invalid token with wrong secret
	invalidToken, _ := token.SignedString([]byte("wrong_secret"))	

	// Execute validate token
	_, err := service.parseAndValidateToken(invalidToken, "access")

	// Verify results
	assert.Error(t, err)

}

// TestParseAndValidateToken_MissingClaims tests token parsing with missing required claims
func TestParseAndValidateToken_MissingClaims(t *testing.T) {


	// Create service instance with mock repositories
	service := &AuthService{}

	// Create token with missing claims
	cfg, _ := config.LoadConfig()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": uuid.New().String(),
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	invalidToken, _ := token.SignedString([]byte(cfg.JWT.Secret))

	// Execute validate token
	_, err := service.parseAndValidateToken(invalidToken, "access")

	// Verify results
	assert.Error(t, err)
}

// TestParseAndValidateToken_MissingUserID tests token parsing when user ID is missing
func TestParseAndValidateToken_MissingUserID(t *testing.T) {
	// Create service instance with mock repositories
	service := &AuthService{}

	// Create token without user ID
	cfg, _ := config.LoadConfig()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"type": "access",
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	invalidToken, _ := token.SignedString([]byte(cfg.JWT.Secret))

	// Execute validate token
	_, err := service.parseAndValidateToken(invalidToken, "access")

	// Verify results
	assert.Error(t, err)
}

// TestParseAndValidateToken_InvalidUserIDFormat tests token parsing with invalid user ID format
func TestParseAndValidateToken_InvalidUserIDFormat(t *testing.T) {
	// Create service instance with mock repositories
	service := &AuthService{}

	// Create token with invalid user ID format
	cfg, _ := config.LoadConfig()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": 123,
		"type":    "access",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	invalidToken, _ := token.SignedString([]byte(cfg.JWT.Secret))

	// Execute validate token
	_, err := service.parseAndValidateToken(invalidToken, "access")

	// Verify results
	assert.Error(t, err)
}

// TestParseAndValidateToken_InvalidUserIDString tests token parsing with invalid user ID string
func TestParseAndValidateToken_InvalidUserIDString(t *testing.T) {
	// Create service instance with mock repositories
	service := &AuthService{}

	// Create token with invalid user ID string
	cfg, _ := config.LoadConfig()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "not-a-uuid",
		"type":    "access",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	invalidToken, _ := token.SignedString([]byte(cfg.JWT.Secret))

	// Execute validate token
	_, err := service.parseAndValidateToken(invalidToken, "access")

	// Verify results
	assert.Error(t, err)
}