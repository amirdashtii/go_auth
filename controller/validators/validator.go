package validators

// import (
// 	"regexp"

// 	"github.com/go-playground/validator/v10"
// )

// func HandleValidationError(err error) map[string]interface{} {
// 	if validationErr, ok := err.(validator.ValidationErrors); ok {
// 		errors := make(map[string]string)
// 		for _, e := range validationErr {
// 			errors[e.Field()] = e.Tag()
// 		}
// 		return map[string]interface{}{
// 			"error":   "Validation failed",
// 			"details": errors,
// 		}
// 	}
// 	return map[string]interface{}{
// 		"error":   "Validation failed",
// 		"details": err.Error(),
// 	}
// }

// func validatePassword(fl validator.FieldLevel) bool {
// 	password := fl.Field().String()

// 	// Minimum 8 characters
// 	if len(password) < 8 {
// 		return false
// 	}

// 	// Contains uppercase letters
// 	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
// 	// Contains lowercase letters
// 	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
// 	// Contains numbers
// 	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

// 	return hasUpper && hasLower && hasNumber
// }
