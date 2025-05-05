package validators

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/go-playground/validator/v10"
)

var authValidate *validator.Validate

func init() {
	authValidate = validator.New()
	authValidate.RegisterValidation("phone", ValidatePhoneNumber)
	authValidate.RegisterValidation("password", ValidateAuthPassword)
}

func ValidatePhoneNumber(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	pattern := `^09[0-9]{9}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

func ValidateAuthPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}
	
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	
	return hasUpper && hasLower && hasNumber
}

func ValidateRegisterRequest(req *dto.RegisterRequest) error {
	if err := authValidate.Struct(req); err != nil {
		return err
	}
	return nil
}

func ValidateLoginRequest(req *dto.LoginRequest) error {
	if err := authValidate.Struct(req); err != nil {
		return err
	}
	return nil
}

func ValidateRefreshTokenRequest(req *dto.RefreshTokenRequest) error {
	if err := authValidate.Struct(req); err != nil {
		return err
	}
	return nil
}

func HandleValidationError(err error) map[string]interface{} {
	validationErrors := make(map[string]interface{})
	
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			field := strings.ToLower(e.Field())
			validationErrors[field] = fmt.Sprintf("Invalid %s: %s", field, e.Tag())
		}
	}
	
	return map[string]interface{}{
		"error":   "Validation failed",
		"details": validationErrors,
	}
}