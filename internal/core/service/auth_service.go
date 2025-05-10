package service

import (
	// "database/sql"

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
	db    ports.AuthRepository
	redis ports.InMemoryRespositoryContracts
	logger ports.Logger
}

func NewAuthService() *AuthService {
	dbRepo, err := repository.NewPGRepository()
	if err != nil {
		panic(errors.ErrDatabaseInit)
	}
	db := dbRepo.DB()

	// Create log file
	logFile, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	// Initialize logger with both file and console output
	loggerConfig := ports.LoggerConfig{
		Level:       "info",
		Environment: "development",
		ServiceName: "go_auth",
		Output:      logFile,
	}
	appLogger := logger.NewZerologLogger(loggerConfig)

	authRepo := repository.NewPGAuthRepository(db, appLogger)
	redisRepo, err := repository.NewRedisRepository(appLogger)
	if err != nil {
		panic(errors.ErrRedisInit)
	}

	return &AuthService{
		db:    authRepo,
		redis: redisRepo,
		logger: appLogger,
	}
}

func (s *AuthService) Register(req *dto.RegisterRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Error generating password hash",
			ports.F("error", err),
			ports.F("phone_number", req.PhoneNumber),
		)	
		return errors.ErrCreateUser
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
		return err
	}

	return nil
}

func (s *AuthService) Login(loginReq *dto.LoginRequest) (*entities.TokenPair, error) {
	user, err := s.db.FindUserByPhoneNumber(&loginReq.PhoneNumber)
	if err != nil {
		return nil, err
	}
	if user.Status == entities.Deleted {
		s.logger.Error("User is deleted",
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

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		s.logger.Error("Password is incorrect",
			ports.F("user_id", user.ID),
		)
		return nil, errors.ErrLogin
	}

	tokenPair, err := s.createTokenPair(user)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	storedToken, err := s.redis.FindToken(user.ID.String() + ":refresh")
	if err != nil {
		return nil, err
	}

	if storedToken != refreshToken {
		s.logger.Error("Invalid refresh token",
			ports.F("user_id", user.ID),
		)
		return nil, errors.ErrInvalidToken
	}

	err = s.Logout(user.ID.String())
	if err != nil {
		return nil, err
	}

	tokenPair, err := s.createTokenPair(user)
	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}

func (s *AuthService) createTokenPair(user *entities.User) (*entities.TokenPair, error) {
	accessToken, err := s.createToken(user, accessTokenExpiration, "access")
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.createToken(user, refreshTokenExpiration, "refresh")
	if err != nil {
		return nil, err
	}

	err = s.redis.AddToken(user.ID.String()+":access", accessToken, accessTokenExpiration)
	if err != nil {
		return nil, err
	}

	err = s.redis.AddToken(user.ID.String()+":refresh", refreshToken, refreshTokenExpiration)
	if err != nil {
		return nil, err
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
		return "", err
	}

	jwtSecret := config.JWT.Secret
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		s.logger.Error("Error creating token",
			ports.F("error", err),
			ports.F("user_id", user.ID),
		)
		return "", errors.ErrTokenCreation
	}

	return tokenString, nil
}

func (s *AuthService) parseAndValidateToken(token string, expectedType string) (*entities.User, error) {
	config, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	jwtSecret := config.JWT.Secret

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
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

	user, err := s.db.FindUserByID(userID)
	if err != nil {
		return nil, err
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

func (s *AuthService) ValidateToken(userID, token string) error {
	storedToken, err := s.redis.FindToken(userID + ":access")
	if err != nil {
		return err
	}

	if storedToken != token {
		s.logger.Error("Invalid access token",
			ports.F("user_id", userID),
		)
		return errors.ErrInvalidToken
	}

	return nil
}
