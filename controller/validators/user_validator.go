package validators

import (
	"fmt"
	"regexp"

	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/go-playground/validator/v10"
)

var userValidate *validator.Validate

func init() {
	userValidate = validator.New()
	userValidate.RegisterValidation("password", validatePassword)
	userValidate.RegisterValidation("phone", validatePhone)
	userValidate.RegisterValidation("name", validateName)
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return isValidPassword(password)
}

func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if phone == "" {
		return true
	}
	// Iranian phone number format: +98XXXXXXXXXX
	phoneRegex := regexp.MustCompile(`^\+98[0-9]{10}$`)
	return phoneRegex.MatchString(phone)
}

func validateName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	if name == "" {
		return true
	}
	return len(name) >= 2 && len(name) <= 50
}

func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	
	return hasUpper && hasLower && hasNumber
}

func ValidateUserUpdateRequest(req *dto.UserUpdateRequest) error {
	if err := userValidate.Struct(req); err != nil {
		return err
	}
	return nil
}

func ValidateChangePasswordRequest(req *dto.ChangePasswordRequest) error {
	if err := userValidate.Struct(req); err != nil {
		return err
	}

	if req.NewPassword == req.OldPassword {
		return fmt.Errorf("new password must be different from current password")
	}

	return nil
}
