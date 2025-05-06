package validators

import (
    "fmt"
    "regexp"

    "github.com/amirdashtii/go_auth/controller/dto"
    "github.com/amirdashtii/go_auth/internal/core/errors"
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

func getCustomPersianErrorMessage(field string) string {
    switch field {
    case "PhoneNumber":
        return "شماره موبایل باید با 09 شروع شده و 11 رقم باشد."
    case "Password":
        return "رمز عبور باید حداقل ۸ کاراکتر و شامل حروف بزرگ، کوچک و عدد باشد."
    case "RefreshToken":
        return "توکن بروزرسانی نامعتبر است."
    default:
        return fmt.Sprintf("فیلد %s نامعتبر است.", field)
    }
}

func getCustomErrorEnglishMessageEn(field string) string {
    switch field {
    case "PhoneNumber":
        return "Phone number must start with 09 and be 11 digits."
    case "Password":
        return "Password must be at least 8 characters and include uppercase, lowercase, and a number."
    case "RefreshToken":
        return "Refresh token is invalid."
    default:
        return fmt.Sprintf("Field %s is invalid.", field)
    }
}

func ValidateRegisterRequest(req *dto.RegisterRequest) error {
    if err := authValidate.Struct(req); err != nil {
        if validationErrs, ok := err.(validator.ValidationErrors); ok {
            field := validationErrs[0].Field()
            return errors.New(
                errors.ValidationError,
                getCustomErrorEnglishMessageEn(field),
                getCustomPersianErrorMessage(field),
                nil,
            )
        }
        return err
    }
    return nil
}

func ValidateLoginRequest(req *dto.LoginRequest) error {
    if err := authValidate.Struct(req); err != nil {
        if validationErrs, ok := err.(validator.ValidationErrors); ok {
            field := validationErrs[0].Field()
            return errors.New(
                errors.ValidationError,
                getCustomErrorEnglishMessageEn(field),
                getCustomPersianErrorMessage(field),
                nil,
            )
        }
        return err
    }
    return nil
}

func ValidateRefreshTokenRequest(req *dto.RefreshTokenRequest) error {
    if err := authValidate.Struct(req); err != nil {
        if validationErrs, ok := err.(validator.ValidationErrors); ok {
            field := validationErrs[0].Field()
            return errors.New(
                errors.ValidationError,
                getCustomErrorEnglishMessageEn(field),
                getCustomPersianErrorMessage(field),
                nil,
            )
        }
        return err
    }
    return nil
}