package validators

// import (
// 	"github.com/amirdashtii/go_auth/controller/dto"
// 	"github.com/amirdashtii/go_auth/internal/core/entities"
// 	"github.com/go-playground/validator/v10"
// )

// func ValidateUser(user *dto.RegisterRequest) error {
// 	validate := validator.New()
// 	RegisterValidator(validate)
// 	return validate.Struct(user)
// }

// func RegisterValidator(v *validator.Validate) {
// 	v.RegisterStructValidation(RegisterValidation, entities.User{})
// 	v.RegisterValidation("password", validatePassword)
// }

// func RegisterValidation(sl validator.StructLevel) {
// 	user := sl.Current().Interface().(entities.User)

// 	// Email validation
// 	if err := sl.Validator().Var(user.Email, "required,email"); err != nil {
// 		sl.ReportError(user.Email, "email", "Email", "Email is required and must be in a valid format", "")
// 	}

// 	// Password validation
// 	if err := sl.Validator().Var(user.Password, "required,password"); err != nil {
// 		sl.ReportError(user.Password, "password", "Password", "Password must be at least 8 characters and contain uppercase, lowercase letters and numbers", "")
// 	}
// }

// func ValidateLogin(req *dto.LoginRequest) error {
// 	validate := validator.New()
// 	LoginValidator(validate)
// 	return validate.Struct(req)
// }

// func LoginValidator(v *validator.Validate) {
// 	v.RegisterStructValidation(LoginValidation, entities.User{})
// 	v.RegisterValidation("password", validatePassword)
// }

// func LoginValidation(sl validator.StructLevel) {
// 	user := sl.Current().Interface().(entities.User)

// 	// Email validation
// 	if err := sl.Validator().Var(user.Email, "required,email"); err != nil {
// 		sl.ReportError(user.Email, "email", "Email", "Email is required and must be in a valid format", "")
// 	}

// 	// Password validation
// 	if err := sl.Validator().Var(user.Password, "required"); err != nil {
// 		sl.ReportError(user.Password, "password", "Password", "Password is required", "")
// 	}
// }
