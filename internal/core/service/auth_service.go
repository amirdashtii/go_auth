package service

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/amirdashtii/go_auth/config"
	"github.com/amirdashtii/go_auth/controller/dto"
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
	db    ports.AuthRepository
	redis ports.InMemoryRespositoryContracts
}

func NewAuthService() *AuthService {
	dbRepo, err := repository.NewPGRepository()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	db := dbRepo.DB()
	authRepo := repository.NewPGAuthRepository(db)
	redisRepo, err := repository.NewRedisRepository()
	if err != nil {
		log.Fatalf("Failed to initialize redis: %v", err)
	}

	return &AuthService{
		db:    authRepo,
		redis: redisRepo,
	}
}

func (s *AuthService) Register(req *dto.RegisterRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := entities.User{
		PhoneNumber: req.PhoneNumber,
		Password:    string(hashedPassword),
		Status:      entities.Active,
		Role:        entities.UserRole,
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.Create(&user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (s *AuthService) Login(loginReq *dto.LoginRequest) (*entities.TokenPair, error) {
	user, err := s.db.FindUserByPhoneNumber(&loginReq.PhoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user.Status == entities.Deleted {
		return nil, fmt.Errorf("user not found")
	}
	if user.Status == entities.Deactivated {
		return nil, fmt.Errorf("user is deactivated")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	if user.Status != entities.Active {
		return nil, fmt.Errorf("user account is not active: %w", err)
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
		return nil, fmt.Errorf("failed to find stored token: %w", err)
	}

	if storedToken != refreshToken {
		return nil, fmt.Errorf("refresh token does not match stored token: %w", err)
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
		"role":       user.Role,
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
		return nil, fmt.Errorf("invalid token claims")
	}

	tokenType, ok := claims["token_type"].(string)
	if !ok || tokenType != expectedType {
		return nil, fmt.Errorf("invalid token type, expected %s", expectedType)
	}

	userIDValue, ok := claims["user_id"]
	if !ok {
		return nil, fmt.Errorf("invalid user ID in token")
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
			return nil, fmt.Errorf("invalid user ID format in token")
		}
	default:
		return nil, fmt.Errorf("unexpected user ID format in token")
	}

	user, err := s.db.FindUserByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user.Status == entities.Deleted {
		return nil, fmt.Errorf("user not found")
	}
	if user.Status == entities.Deactivated {
		return nil, fmt.Errorf("user is deactivated")
	}

	if user.Status != entities.Active {
		return nil, fmt.Errorf("user account is not active")
	}

	return user, nil
}

func (s *AuthService) ValidateToken(userID, token string) error {

	storedToken, err := s.redis.FindToken(userID + ":access")
	if err != nil {
		return fmt.Errorf("failed to find token: %w", err)
	}

	if storedToken != token {
		return fmt.Errorf("token does not match")
	}

	return nil
}
