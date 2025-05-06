package service

import (
	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/amirdashtii/go_auth/infrastructure/repository"
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/errors"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db ports.UserRepository
}

func NewUserService() *UserService {
	dbRepo, err := repository.NewPGRepository()
	if err != nil {
		panic(err)
	}
	db := dbRepo.DB()
	userRepo := repository.NewPGUserRepository(db)
	return &UserService{
		db: userRepo,
	}
}

func (s *UserService) GetProfile(userID *uuid.UUID) (*dto.UserProfileResponse, error) {
	user, err := s.db.FindUserByID(userID)
	if err != nil {
		return nil, errors.New(errors.NotFoundError, "user not found", "کاربر یافت نشد", err)
	}
	if user.Status == entities.Deleted {
		return nil, errors.New(errors.NotFoundError, "user not found", "کاربر یافت نشد", nil)
	}
	if user.Status == entities.Deactivated {
		return nil, errors.New(errors.AuthenticationError, "user is deactivated", "حساب کاربری غیرفعال است", nil)
	}

	resp := dto.UserProfileResponse{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
	return &resp, nil
}

func (s *UserService) UpdateProfile(userID *uuid.UUID, req *dto.UserUpdateRequest) error {
	user := &entities.User{
		ID:          *userID,
		PhoneNumber: req.PhoneNumber,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
	}

	if err := s.db.Update(user); err != nil {
		return errors.New(errors.InternalError, "failed to update user profile", "خطا در به\u200cروزرسانی پروفایل کاربر", err)
	}
	return nil
}

func (s *UserService) ChangePassword(userID *uuid.UUID, changePasswordReq *dto.ChangePasswordRequest) error {
	currentUser, err := s.db.FindUserByID(userID)
	if err != nil {
		return errors.New(errors.NotFoundError, "user not found", "کاربر یافت نشد", err)
	}

	if currentUser.Status == entities.Deleted {
		return errors.New(errors.NotFoundError, "user not found", "کاربر یافت نشد", nil)
	}
	if currentUser.Status == entities.Deactivated {
		return errors.New(errors.AuthenticationError, "user is deactivated", "حساب کاربری غیرفعال است", nil)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(changePasswordReq.OldPassword)); err != nil {
		return errors.New(errors.AuthenticationError, "current password is incorrect", "رمز عبور فعلی اشتباه است", err)
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(changePasswordReq.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New(errors.InternalError, "failed to hash password", "خطا در رمزنگاری رمز عبور", err)
	}

	user := &entities.User{}
	user.ID = *userID
	user.Password = string(hashedNewPassword)
	if err := s.db.Update(user); err != nil {
		return errors.New(errors.InternalError, "failed to update password", "خطا در به\u200cروزرسانی رمز عبور", err)
	}
	return nil
}

func (s *UserService) DeleteProfile(userID *uuid.UUID) error {
	if err := s.db.Delete(userID); err != nil {
		return errors.New(errors.InternalError, "failed to delete user", "خطا در حذف کاربر", err)
	}
	return nil
}
