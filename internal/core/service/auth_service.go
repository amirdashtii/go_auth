package service

import (
	"context"
	"os"
	"time"

	"github.com/amirdashtii/go_auth/config"
	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/amirdashtii/go_auth/infrastructure/logger"
	"github.com/amirdashtii/go_auth/infrastructure/repository"
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/errors"
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
	db        ports.AuthRepository
	redis     ports.InMemoryRespositoryContracts
	logger    ports.Logger
	jwtSecret string
}

func NewAuthService() *AuthService {
	dbRepo, err := repository.GetPGRepository()
	if err != nil {
		panic(errors.ErrDatabaseInit)
	}
	db := dbRepo.DB()

	config, err := config.LoadConfig()
	if err != nil {
		panic(errors.ErrLoadConfig)
	}

	// Initialize logger with both file and console output
	loggerConfig := ports.LoggerConfig{
		Level:       "info",
		Environment: "development",
		ServiceName: "go_auth",
		Output:      os.Stdout,
	}
	appLogger := logger.NewZerologLogger(loggerConfig)

	authRepo := repository.NewPGAuthRepository(db, appLogger)
	redisRepo, err := repository.GetRedisRepository(appLogger)
	if err != nil {
		panic(errors.ErrRedisInit)
	}

	return &AuthService{
		db:        authRepo,
		redis:     redisRepo,
		logger:    appLogger,
		jwtSecret: config.JWT.Secret,
	}
}

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) error {
	if ctx.Err() != nil {
		s.logger.Error("Context cancelled while registering user",
			ports.F("error", ctx.Err()),
			ports.F("phone_number", req.PhoneNumber),
		)
		return errors.ErrContextCancelled
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Error hashing password",
			ports.F("error", err),
		)
		return errors.ErrCreateUser
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}
	
	user := &entities.User{
		ID:          uuid.New(),
		PhoneNumber: req.PhoneNumber,
		Password:    string(hashedPassword),
		Status:      entities.Active,
		Role:        entities.UserRole,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, loginReq *dto.LoginRequest) (*entities.TokenPair, error) {
	if ctx.Err() != nil {
		s.logger.Error("Context cancelled while logging in user",
			ports.F("error", ctx.Err()),
			ports.F("phone_number", loginReq.PhoneNumber),
		)
		return nil, errors.ErrContextCancelled
	}
	
	user, err := s.db.FindUserByPhoneNumber(ctx, &loginReq.PhoneNumber)
	if err != nil {
		return nil, err
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		s.logger.Error("Invalid password",
			ports.F("error", err),
			ports.F("user_id", user.ID),
		)
		return nil, errors.ErrInvalidCredentials
	}

	// Check user status
	if user.Status == entities.Deleted {
		s.logger.Error("User is deleted",
			ports.F("error", err),
			ports.F("user_id", user.ID),
		)
		return nil, errors.ErrInvalidCredentials
	}
	if user.Status == entities.Deactivated {
		s.logger.Error("User is deactivated",
			ports.F("user_id", user.ID),
		)
		return nil, errors.ErrAccountDeactivated
	}

	// Generate tokens
	tokens, err := s.createTokenPair(ctx, user)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *AuthService) Logout(ctx context.Context, userID string) error {
	if ctx.Err() != nil {
		s.logger.Error("Context cancelled while logging out user",
			ports.F("error", ctx.Err()),
			ports.F("user_id", userID),
		)
		return errors.ErrContextCancelled
	}

	err := s.redis.RemoveToken(ctx, userID+":access")
	if err != nil {
		return err
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	err = s.redis.RemoveToken(ctx, userID+":refresh")
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*entities.TokenPair, error) {
	if ctx.Err() != nil {
		s.logger.Error("Context cancelled while refreshing token",
			ports.F("error", ctx.Err()),
			ports.F("refresh_token", refreshToken),
		)
		return nil, errors.ErrContextCancelled
	}

	user, err := s.parseAndValidateToken(ctx, refreshToken, "refresh")
	if err != nil {
		return nil, err
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	storedToken, err := s.redis.FindToken(ctx, user.ID.String()+":refresh")
	if err != nil {
		return nil, err
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if storedToken != refreshToken {
		s.logger.Error("Invalid refresh token",
			ports.F("user_id", user.ID),
		)
		return nil, errors.ErrInvalidToken
	}

	err = s.Logout(ctx, user.ID.String())
	if err != nil {
		return nil, err
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	tokenPair, err := s.createTokenPair(ctx, user)
	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}

func (s *AuthService) createTokenPair(ctx context.Context, user *entities.User) (*entities.TokenPair, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	accessToken, err := s.createToken(ctx, user, accessTokenExpiration, "access")
	if err != nil {
		return nil, err
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	refreshToken, err := s.createToken(ctx, user, refreshTokenExpiration, "refresh")
	if err != nil {
		return nil, err
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	err = s.redis.AddToken(ctx, user.ID.String()+":access", accessToken, accessTokenExpiration)
	if err != nil {
		return nil, err
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	err = s.redis.AddToken(ctx, user.ID.String()+":refresh", refreshToken, refreshTokenExpiration)
	if err != nil {
		return nil, err
	}

	return &entities.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) createToken(ctx context.Context, user *entities.User, expiration time.Duration, tokenType string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	claims := jwt.MapClaims{
		"user_id":    user.ID,
		"role":       user.Role,
		"token_type": tokenType,
		"exp":        time.Now().Add(expiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		s.logger.Error("Error creating token",
			ports.F("error", err),
			ports.F("user_id", user.ID),
		)
		return "", errors.ErrTokenCreation
	}

	return tokenString, nil
}

func (s *AuthService) parseAndValidateToken(ctx context.Context, token string, expectedType string) (*entities.User, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		s.logger.Error("Error parsing token",
			ports.F("error", err),
			ports.F("token", token),
		)
		return nil, errors.ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		s.logger.Error("Invalid token claims",
			ports.F("token", token),
		)
		return nil, errors.ErrInvalidToken
	}

	tokenType, ok := claims["token_type"].(string)
	if !ok || tokenType != expectedType {
		s.logger.Error("Invalid token type",
			ports.F("token", token),
		)
		return nil, errors.ErrInvalidToken
	}

	userIDValue, ok := claims["user_id"]
	if !ok {
		s.logger.Error("Invalid token claims",
			ports.F("token", token),
		)
		return nil, errors.ErrInvalidToken
	}

	var userID uuid.UUID
	switch v := userIDValue.(type) {
	case string:
		userID, err = uuid.Parse(v)
		if err != nil {
			s.logger.Error("Invalid user ID",
				ports.F("error", err),
				ports.F("token", token),
			)
			return nil, errors.ErrInvalidToken
		}
	case map[string]interface{}:
		if uuidStr, ok := v["String"].(string); ok {
			userID, err = uuid.Parse(uuidStr)
			if err != nil {
				s.logger.Error("Invalid user ID",
					ports.F("error", err),
					ports.F("token", token),
				)
				return nil, errors.ErrInvalidToken
			}
		} else {
			s.logger.Error("Invalid user ID",
				ports.F("token", token),
			)
			return nil, errors.ErrInvalidToken
		}
	default:
		s.logger.Error("Invalid user ID",
			ports.F("token", token),
		)
		return nil, errors.ErrInvalidToken
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	user, err := s.db.FindUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if user.Status == entities.Deleted {
		s.logger.Error("User is deleted",
			ports.F("user_id", userID),
		)
		return nil, errors.ErrInvalidCredentials
	}
	if user.Status == entities.Deactivated {
		s.logger.Error("User is deactivated",
			ports.F("user_id", userID),
		)
		return nil, errors.ErrAccountDeactivated
	}

	return user, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, userID, token string) error {
	if ctx.Err() != nil {
		s.logger.Error("Context cancelled while validating token",
			ports.F("error", ctx.Err()),
			ports.F("user_id", userID),
			ports.F("token", token),
		)
		return errors.ErrContextCancelled
	}

	storedToken, err := s.redis.FindToken(ctx, userID+":access")
	if err != nil {
		return err
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	if storedToken != token {
		s.logger.Error("Invalid access token",
			ports.F("user_id", userID),
		)
		return errors.ErrInvalidToken
	}

	return nil
}
