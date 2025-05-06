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
		return errors.New(errors.InternalError, "failed to hash password", "خطا در رمزنگاری رمز عبور", nil)
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
		return errors.New(errors.InternalError, "failed to create user", "خطا در ایجاد کاربر", nil)
	}

	return nil
}

func (s *AuthService) Login(loginReq *dto.LoginRequest) (*entities.TokenPair, error) {
	user, err := s.db.FindUserByPhoneNumber(&loginReq.PhoneNumber)
	// TODO
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(errors.NotFoundError, "user not found", "کاربر یافت نشد", nil)
		}
		return nil, errors.New(errors.InternalError, "failed to find user", "خطا در جستجوی کاربر", nil)
	}
	if user.Status == entities.Deleted {
		return nil, errors.New(errors.NotFoundError, "user not found", "کاربر یافت نشد", nil)
	}
	if user.Status == entities.Deactivated {
		return nil, errors.New(errors.AuthenticationError, "user is deactivated", "حساب کاربری غیرفعال است", nil)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		return nil, errors.New(errors.AuthenticationError, "invalid password", "رمز عبور نامعتبر است", nil)
	}

	if user.Status != entities.Active {
		return nil, errors.New(errors.AuthenticationError, "user account is not active", "حساب کاربری فعال نیست", nil)
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
		return errors.New(errors.InternalError, "failed to remove access token", "خطا در حذف توکن دسترسی", nil)
	}

	err = s.redis.RemoveToken(userID + ":refresh")
	if err != nil {
		return errors.New(errors.InternalError, "failed to remove refresh token", "خطا در حذف توکن بروزرسانی", nil)
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
		return nil, errors.New(errors.InternalError, "failed to find stored token", "خطا در یافتن توکن ذخیره شده", nil)
	}

	if storedToken != refreshToken {
		return nil, errors.New(errors.AuthenticationError, "refresh token does not match stored token", "توکن بروزرسانی با توکن ذخیره شده مطابقت ندارد", nil)
	}

	err=s.Logout(user.ID.String())
	if err != nil{
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
		return &entities.TokenPair{}, err
	}

	err = s.redis.AddToken(user.ID.String()+":refresh", refreshToken, refreshTokenExpiration)
	if err != nil {
		return &entities.TokenPair{}, err
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
	return token.SignedString([]byte(jwtSecret))
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
		return nil, errors.New(errors.AuthenticationError, "failed to parse token", "خطا در تجزیه توکن", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New(errors.AuthenticationError, "invalid token claims", "اطلاعات توکن نامعتبر است", nil)
	}

	tokenType, ok := claims["token_type"].(string)
	if !ok || tokenType != expectedType {
		return nil, errors.New(errors.AuthenticationError, fmt.Sprintf("invalid token type, expected %s", expectedType), "نوع توکن نامعتبر است", nil)
	}

	userIDValue, ok := claims["user_id"]
	if !ok {
		return nil, errors.New(errors.AuthenticationError, "invalid user ID in token", "شناسه کاربر در توکن نامعتبر است", nil)
	}

	var userID uuid.UUID
	switch v := userIDValue.(type) {
	case string:
		userID, err = uuid.Parse(v)
		if err != nil {
			return nil, errors.New(errors.AuthenticationError, "failed to parse user ID", "خطا در تجزیه شناسه کاربر", nil)
		}
	case map[string]interface{}:
		if uuidStr, ok := v["String"].(string); ok {
			userID, err = uuid.Parse(uuidStr)
			if err != nil {
				return nil, errors.New(errors.AuthenticationError, "failed to parse user ID", "خطا در تجزیه شناسه کاربر", nil)
			}
		} else {
			return nil, errors.New(errors.AuthenticationError, "invalid user ID format in token", "فرمت شناسه کاربر در توکن نامعتبر است", nil)
		}
	default:
		return nil, errors.New(errors.AuthenticationError, "unexpected user ID format in token", "فرمت غیرمنتظره شناسه کاربر در توکن", nil)
	}

	user, err := s.db.FindUserByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(errors.NotFoundError, "user not found", "کاربر یافت نشد", nil)
		}
		return nil, errors.New(errors.InternalError, "failed to find user", "خطا در یافتن کاربر", nil)
	}
	if user.Status == entities.Deleted {
		return nil, errors.New(errors.NotFoundError, "user not found", "کاربر یافت نشد", nil)
	}
	if user.Status == entities.Deactivated {
		return nil, errors.New(errors.AuthenticationError, "user is deactivated", "حساب کاربری غیرفعال است", nil)
	}

	if user.Status != entities.Active {
		return nil, errors.New(errors.AuthenticationError, "user account is not active", "حساب کاربری فعال نیست", nil)
	}

	return user, nil
}

func (s *AuthService) ValidateToken(userID, token string) error {
	storedToken, err := s.redis.FindToken(userID + ":access")
	if err != nil {
		return err
	}

	if storedToken != token {
		return errors.New(errors.AuthenticationError, "token does not match", "توکن مطابقت ندارد", nil)
	}

	return nil
}
