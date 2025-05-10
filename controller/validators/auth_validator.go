package validators

import (
	"fmt"
	"regexp"

	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/amirdashtii/go_auth/internal/core/errors"
	"github.com/amirdashtii/go_auth/internal/core/ports"
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

func getAuthCustomErrorMessage(field string) error {
    switch field {
    case "PhoneNumber":
        return errors.ErrInvalidPhoneNumber
    case "Password":
        return errors.ErrInvalidPassword
    case "RefreshToken":
        return errors.ErrInvalidRefreshToken
    default:
        return errors.New(errors.ValidationError, fmt.Sprintf("Field %s is invalid.", field), fmt.Sprintf("فیلد %s نامعتبر است.", field), nil)
    }
}


func ValidateRegisterRequest(req *dto.RegisterRequest, logger ports.Logger) error {
    if err := authValidate.Struct(req); err != nil {
        if validationErrs, ok := err.(validator.ValidationErrors); ok {
            field := validationErrs[0].Field()
            logger.Error("Validation error",
				ports.F("error", err),
				ports.F("field", field),
			)   
            return getAuthCustomErrorMessage(field)
        }
		logger.Error("Validation error",
			ports.F("error", err),
		)
    }
    return nil
}

func ValidateLoginRequest(req *dto.LoginRequest, logger ports.Logger) error {
    if err := authValidate.Struct(req); err != nil {
        if validationErrs, ok := err.(validator.ValidationErrors); ok {
            field := validationErrs[0].Field()
            logger.Error("Validation error",
				ports.F("error", err),
				ports.F("field", field),
			)
            return getAuthCustomErrorMessage(field)
        }
		logger.Error("Validation error",
			ports.F("error", err),
		)
    }
    return nil
}

func ValidateRefreshTokenRequest(req *dto.RefreshTokenRequest, logger ports.Logger) error {
    if err := authValidate.Struct(req); err != nil {
        if validationErrs, ok := err.(validator.ValidationErrors); ok {
            field := validationErrs[0].Field()
            logger.Error("Validation error",
				ports.F("error", err),
				ports.F("field", field),
			)
                return getAuthCustomErrorMessage(field)
        }
		logger.Error("Validation error",
			ports.F("error", err),
		)
    }
    return nil
}