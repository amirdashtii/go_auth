package validators

import (
	"fmt"
	"regexp"

	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/amirdashtii/go_auth/internal/core/errors"
	"github.com/amirdashtii/go_auth/internal/core/ports"
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

func getUserCustomErrorMessage(field string) error {
	switch field {
	case "PhoneNumber":
		return errors.ErrInvalidPhoneNumber
	case "FirstName":
		return errors.ErrInvalidFirstName
	case "LastName":
		return errors.ErrInvalidLastName
	case "Email":
		return errors.ErrInvalidEmail
	case "OldPassword":
		return errors.ErrInvalidOldPassword
	case "NewPassword":
		return errors.ErrInvalidNewPassword
	default:
		return errors.New(errors.ValidationError, fmt.Sprintf("%s Field is invalid.", field), fmt.Sprintf("فیلد %s نامعتبر است.", field), nil)
	}
}

func ValidateUserUpdateRequest(req *dto.UserUpdateRequest, logger ports.Logger) error {
	if err := userValidate.Struct(req); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			field := validationErrs[0].Field()
			logger.Error("Validation error",
				ports.F("error", err),
				ports.F("field", field),
			)
			return getUserCustomErrorMessage(field)
		}
		logger.Error("Validation error",
			ports.F("error", err),
		)
	}
	return nil
}

func ValidateChangePasswordRequest(req *dto.ChangePasswordRequest, logger ports.Logger) error {
	if err := userValidate.Struct(req); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			field := validationErrs[0].Field()
			logger.Error("Validation error",
				ports.F("error", err),
				ports.F("field", field),
			)
			return getUserCustomErrorMessage(field)
		}
		logger.Error("Validation error",
			ports.F("error", err),
		)
	}

	if req.NewPassword == req.OldPassword {
		logger.Error("Validation error",
			ports.F("error", "new password must be different from current password"),
		)
		return errors.New(errors.ValidationError, "new password must be different from current password", "پسورد جدید باید با پسورد قدیمی فرق داشته باشد.", nil)
	}

	return nil
}
