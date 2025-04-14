package validators

import (
	"regexp"

	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/go-playground/validator/v10"
)

func HandleValidationError(err error) map[string]interface{} {
	if validationErr, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)
		for _, e := range validationErr {
			errors[e.Field()] = e.Tag()
		}
		return map[string]interface{}{
			"error": "Validation failed",
			"details": errors,
		}
	}
	return map[string]interface{}{
		"error": "Validation failed",
		"details": err.Error(),
	}
}


func ValidateUser(user entities.User) error {
	validate := validator.New()
	RegisterValidator(validate)
	return validate.Struct(user)
}

func RegisterValidator(v *validator.Validate) {
	v.RegisterStructValidation(RegisterValidation, entities.User{})
	v.RegisterValidation("password", validatePassword)
}

func RegisterValidation(sl validator.StructLevel) {
	user := sl.Current().Interface().(entities.User)
	
	// Email validation
	if err := sl.Validator().Var(user.Email, "required,email"); err != nil {
		sl.ReportError(user.Email, "email", "Email", "Email is required and must be in a valid format", "")
	}

	// Password validation
	if err := sl.Validator().Var(user.Password, "required,password"); err != nil {
		sl.ReportError(user.Password, "password", "Password", "Password must be at least 8 characters and contain uppercase, lowercase letters and numbers", "")
	}
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	
	// Minimum 8 characters
	if len(password) < 8 {
		return false
	}

	// Contains uppercase letters
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	// Contains lowercase letters
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	// Contains numbers
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	return hasUpper && hasLower && hasNumber
}

func LoginValidator(v *validator.Validate) {
	v.RegisterStructValidation(LoginValidation, entities.User{})
}

func LoginValidation(sl validator.StructLevel) {
	user := sl.Current().Interface().(entities.User)
	
	// Email validation
	if err := sl.Validator().Var(user.Email, "required,email"); err != nil {
		sl.ReportError(user.Email, "email", "Email", "Email is required and must be in a valid format", "")
	}

	// Password validation
	if err := sl.Validator().Var(user.Password, "required"); err != nil {
		sl.ReportError(user.Password, "password", "Password", "Password is required", "")
	}
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,password"`
}

func ChangePasswordValidator(v *validator.Validate) {
	v.RegisterStructValidation(ChangePasswordValidation, ChangePasswordRequest{})
}

func ChangePasswordValidation(sl validator.StructLevel) {
	req := sl.Current().Interface().(ChangePasswordRequest)
	
	// Old password validation
	if err := sl.Validator().Var(req.OldPassword, "required"); err != nil {
		sl.ReportError(req.OldPassword, "old_password", "OldPassword", "Current password is required", "")
	}

	// New password validation
	if err := sl.Validator().Var(req.NewPassword, "required,password"); err != nil {
		sl.ReportError(req.NewPassword, "new_password", "NewPassword", "New password must be at least 8 characters and contain uppercase, lowercase letters and numbers", "")
	}
}

type UpdateProfileRequest struct {
	Email string `json:"email" validate:"omitempty,email"`
}

func UpdateProfileValidator(v *validator.Validate) {
	v.RegisterStructValidation(UpdateProfileValidation, UpdateProfileRequest{})
}

func UpdateProfileValidation(sl validator.StructLevel) {
	req := sl.Current().Interface().(UpdateProfileRequest)
	
	// Email validation
	if req.Email != "" {
		if err := sl.Validator().Var(req.Email, "email"); err != nil {
			sl.ReportError(req.Email, "email", "Email", "Email must be in a valid format", "")
		}
	}
}


