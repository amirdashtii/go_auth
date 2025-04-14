package service

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/amirdashtii/go_auth/infrastructure/repository"
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenExpiration = 24 * time.Hour
)

type AuthService struct {
	db               ports.UserRepository
	whitelistedTokens map[string]time.Time
	mutex            sync.RWMutex
}

func NewAuthService() *AuthService {
	db, err := repository.NewPGRepository()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	return &AuthService{
		db:               db,
		whitelistedTokens: make(map[string]time.Time),
	}
}

// Helper function to create JWT token
func (s *AuthService) createToken(user *entities.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"is_admin": user.IsAdmin,
		"exp":      time.Now().Add(tokenExpiration).Unix(),
	})

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("JWT_SECRET not set in environment")
	}

	return token.SignedString([]byte(jwtSecret))
}

// Helper function to add token to whitelist
func (s *AuthService) addToWhitelist(token string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.whitelistedTokens[token] = time.Now().Add(tokenExpiration)
}

// Helper function to remove token from whitelist
func (s *AuthService) removeFromWhitelist(token string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.whitelistedTokens, token)
}

// Helper function to validate token in whitelist
func (s *AuthService) isTokenWhitelisted(token string) (bool, time.Time) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	expTime, exists := s.whitelistedTokens[token]
	return exists, expTime
}

func (s *AuthService) Register(user *entities.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = string(hashedPassword)
	user.IsActive = true
	user.IsAdmin = false
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if err := s.db.Create(user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.db.FindByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}

	if !user.IsActive {
		return "", errors.New("user account is not active")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	tokenString, err := s.createToken(user)
	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}

	s.addToWhitelist(tokenString)
	return tokenString, nil
}

func (s *AuthService) Logout(token string) error {
	s.removeFromWhitelist(token)
	return nil
}

func (s *AuthService) RefreshToken(token string) (string, error) {
	user, err := s.ValidateToken(token)
	if err != nil {
		return "", err
	}

	newToken, err := s.createToken(user)
	if err != nil {
		return "", fmt.Errorf("failed to create new token: %w", err)
	}

	s.mutex.Lock()
	delete(s.whitelistedTokens, token)
	s.whitelistedTokens[newToken] = time.Now().Add(tokenExpiration)
	s.mutex.Unlock()

	return newToken, nil
}

func (s *AuthService) ValidateToken(token string) (*entities.User, error) {
	isWhitelisted, expTime := s.isTokenWhitelisted(token)
	if !isWhitelisted {
		return nil, errors.New("token is not whitelisted")
	}

	if time.Now().After(expTime) {
		s.removeFromWhitelist(token)
		return nil, errors.New("token is expired")
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user ID: %w", err)
	}

	user, err := s.db.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if !user.IsActive {
		return nil, errors.New("user account is not active")
	}

	return user, nil
}
