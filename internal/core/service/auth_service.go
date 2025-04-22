package service

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/amirdashtii/go_auth/config"
	"github.com/amirdashtii/go_auth/infrastructure/repository"
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenExpiration  = 1 * time.Hour
	refreshTokenExpiration = 7 * 24 * time.Hour
)

type AuthService struct {
	db    ports.UserRepository
	redis ports.InMemoryRespositoryContracts
}

func NewAuthService() *AuthService {
	db, err := repository.NewPGRepository()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	redisRepo, err := repository.NewRedisRepository()
	if err != nil {
		log.Fatalf("Failed to initialize redis: %v", err)
	}

	return &AuthService{
		db:    db,
		redis: redisRepo,
	}
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

func (s *AuthService) Login(email, password string) (*entities.TokenPair, error) {
	user, err := s.db.FindByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	if !user.IsActive {
		return nil, errors.New("user account is not active")
	}

	tokenPair, err := s.createTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create token pair: %w", err)
	}

	return tokenPair, nil
}

func (s *AuthService) Logout(userID string) error {

	err := s.redis.RemoveToken(userID + ":access")
	if err != nil {
		return err
	}

	err = s.redis.RemoveToken(userID + ":refresh")
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*entities.TokenPair, error) {

	user, err := s.parseAndValidateToken(refreshToken, "refresh")
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	storedToken, err := s.redis.FindToken(user.ID.String() + ":refresh")
	if err != nil {
		return nil, errors.New("failed to find stored token")
	}

	if storedToken != refreshToken {
		return nil, errors.New("refresh token does not match stored token")
	}

	s.Logout(user.ID.String())

	tokenPair, err := s.createTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create new token pair: %w", err)
	}

	return tokenPair, nil
}

func (s *AuthService) createTokenPair(user *entities.User) (*entities.TokenPair, error) {
	accessToken, err := s.createToken(user, accessTokenExpiration, "access")
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	refreshToken, err := s.createToken(user, refreshTokenExpiration, "refresh")
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	err = s.redis.AddToken(user.ID.String()+":access", accessToken, accessTokenExpiration)
	if err != nil {
		return &entities.TokenPair{}, fmt.Errorf("failed to store access token in redis: %w", err)
	}

	err = s.redis.AddToken(user.ID.String()+":refresh", refreshToken, refreshTokenExpiration)
	if err != nil {
		return &entities.TokenPair{}, fmt.Errorf("failed to store refresh token in redis: %w", err)
	}

	return &entities.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) createToken(user *entities.User, expiration time.Duration, tokenType string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    user.ID,
		"is_admin":   user.IsAdmin,
		"token_type": tokenType,
		"exp":        time.Now().Add(expiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	jwtSecret := config.JWT.Secret
	return token.SignedString([]byte(jwtSecret))
}

func (s *AuthService) parseAndValidateToken(token string, expectedType string) (*entities.User, error) {

	config, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	jwtSecret := config.JWT.Secret

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	tokenType, ok := claims["token_type"].(string)
	if !ok || tokenType != expectedType {
		return nil, fmt.Errorf("invalid token type, expected %s", expectedType)
	}

	userIDValue, ok := claims["user_id"]
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}

	var userID uuid.UUID
	switch v := userIDValue.(type) {
	case string:
		var err error
		userID, err = uuid.Parse(v)
		if err != nil {
			return nil, fmt.Errorf("failed to parse user ID: %w", err)
		}
	case map[string]interface{}:
		if uuidStr, ok := v["String"].(string); ok {
			var err error
			userID, err = uuid.Parse(uuidStr)
			if err != nil {
				return nil, fmt.Errorf("failed to parse user ID: %w", err)
			}
		} else {
			return nil, errors.New("invalid user ID format in token")
		}
	default:
		return nil, errors.New("unexpected user ID format in token")
	}

	user, err := s.db.FindByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if !user.IsActive {
		return nil, errors.New("user account is not active")
	}

	return user, nil
}

func (s *AuthService) ValidateToken(userID, token string) error {

	storedToken, err := s.redis.FindToken(userID + ":access")
	if err != nil {
		return errors.New("failed to find token")
	}

	if storedToken != token {
		return errors.New("token does not match")
	}

	return nil
}
