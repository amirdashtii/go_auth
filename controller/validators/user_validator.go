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
	userValidate.RegisterValidation("password", validateAuthPassword)
}

func validateUserPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}
	
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)
	
	return hasUpper && hasLower && hasNumber && hasSpecial
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
