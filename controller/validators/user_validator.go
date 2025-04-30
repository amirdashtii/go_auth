package validators





// func ChangePasswordValidation(sl validator.StructLevel) {
// 	req := sl.Current().Interface().(ChangePasswordRequest)

// 	// Old password validation
// 	if err := sl.Validator().Var(req.OldPassword, "required"); err != nil {
// 		sl.ReportError(req.OldPassword, "old_password", "OldPassword", "Current password is required", "")
// 	}

// 	// New password validation
// 	if err := sl.Validator().Var(req.NewPassword, "required,password"); err != nil {
// 		sl.ReportError(req.NewPassword, "new_password", "NewPassword", "New password must be at least 8 characters and contain uppercase, lowercase letters and numbers", "")
// 	}
// }

// type UpdateProfileRequest struct {
// 	Email string `json:"email" validate:"omitempty,email"`
// }

// func UpdateProfileValidator(v *validator.Validate) {
// 	v.RegisterStructValidation(UpdateProfileValidation, UpdateProfileRequest{})
// }

// func UpdateProfileValidation(sl validator.StructLevel) {
// 	req := sl.Current().Interface().(UpdateProfileRequest)

// 	// Email validation
// 	if req.Email != "" {
// 		if err := sl.Validator().Var(req.Email, "email"); err != nil {
// 			sl.ReportError(req.Email, "email", "Email", "Email must be in a valid format", "")
// 		}
// 	}
// }
