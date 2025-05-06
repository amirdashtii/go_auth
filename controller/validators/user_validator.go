package validators

import (
	"fmt"
	"regexp"

	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/amirdashtii/go_auth/internal/core/errors"
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

func getUserCustomPersianErrorMessage(field string) string {
	switch field {
	case "PhoneNumber":
		return "شماره موبایل باید با 09 شروع شده و 11 رقم باشد."
	case "FirstName":
		return "فیلد نام کوچک نامعتبر است."
	case "LastName":
		return "فیلد نام خانوادگی نامعتبر است."
	case "Email":
		return "فیلد ایمیل نامعتبر است."
	case "OldPassword":
		return "پسورد قدیمی باید حداقل 8 کاراکتر و شامل حروف بزرگ، کوچک و یک عدد باشد."
	case "NewPassword":
		return "پسورد جدید باید حداقل 8 کاراکتر و شامل حروف بزرگ، کوچک و یک عدد باشد."
	default:
		return fmt.Sprintf("فیلد %s نامعتبر است.", field)
	}
}

func getUserCustomErrorEnglishMessageEn(field string) string {
	switch field {
	case "PhoneNumber":
		return "Phone number must start with 09 and be 11 digits."
	case "FirstName":
		return "First name field is invalid."
	case "LastName":
		return "Last name field is invalid."
	case "Email":
		return "Email field is invalid."
	case "OldPassword":
		return "Old password must be at least 8 characters and include uppercase, lowercase, and a number."
	case "NewPassword":
		return "New password must be at least 8 characters and include uppercase, lowercase, and a number."
	default:
		return fmt.Sprintf("%s Field is invalid.", field)
	}
}

func ValidateUserUpdateRequest(req *dto.UserUpdateRequest) error {
	if err := userValidate.Struct(req); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			field := validationErrs[0].Field()
			return errors.New(
				errors.ValidationError,
				getUserCustomErrorEnglishMessageEn(field),
				getUserCustomPersianErrorMessage(field),
				nil,
			)
		}
		return err
	}
	return nil
}

func ValidateChangePasswordRequest(req *dto.ChangePasswordRequest) error {
	if err := userValidate.Struct(req); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			field := validationErrs[0].Field()
			return errors.New(
				errors.ValidationError,
				getUserCustomErrorEnglishMessageEn(field),
				getUserCustomPersianErrorMessage(field),
				nil,
			)
		}
		return err
	}

	if req.NewPassword == req.OldPassword {
		return errors.New(errors.ValidationError, "new password must be different from current password", "پسورد جدید باید با پسورد قدیمی فرق داشته باشد.", nil)
	}

	return nil
}
