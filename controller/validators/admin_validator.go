package validators

import (
	"fmt"

	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/errors"
	"github.com/go-playground/validator/v10"
)

var adminValidate *validator.Validate

func init() {
	adminValidate = validator.New()
	adminValidate.RegisterValidation("role", validateRole)
	adminValidate.RegisterValidation("status", validateStatus)
	adminValidate.RegisterValidation("sort", validateSort)
	adminValidate.RegisterValidation("order", validateOrder)
}

// validateRole checks if the role is valid without hardcoding role types
func validateRole(fl validator.FieldLevel) bool {
	role := fl.Field().String()
	if role == "" {
		return true
	}

	// Get role type from entities
	roleType := entities.ParseRoleType(role)
	// Check if it's a known role (ParseRoleType returns UserRole for unknown roles)
	return roleType.String() != "Unknown"
}

// validateStatus checks if the status is valid without hardcoding status types
func validateStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return true
	}

	// Get status type from entities
	statusType := entities.ParseStatusType(status)
	// Check if it's a known status (ParseStatusType returns Active for unknown statuses)
	return statusType.String() != "Unknown"
}

func validateSort(fl validator.FieldLevel) bool {
	validSortFields := map[string]bool{
		"":           true,
		"created_at": true,
		"updated_at": true,
		"email":      true,
		"role":       true,
		"status":     true,
		"first_name": true,
		"last_name":  true,
	}

	sort := fl.Field().String()
	return validSortFields[sort]
}

func validateOrder(fl validator.FieldLevel) bool {
	validOrders := map[string]bool{
		"":     true,
		"asc":  true,
		"desc": true,
	}

	order := fl.Field().String()
	return validOrders[order]
}

func getAdminCustomPersianErrorMessage(field string) string {
	switch field {
	case "Sort":
		return "فیلد مرتب\u200cسازی نامعتبر است."
	case "Role":
		return "فیلد نقش نامعتبر است."
	case "Status":
		return "فیلد وضعیت نامعتبر است."
	case "Order":
		return "فیلد ترتیب نامعتبر است."
	case "PhoneNumber":
		return "شماره موبایل باید با 09 شروع شده و 11 رقم باشد."
	case "FirstName":
		return "فیلد نام کوچک نامعتبر است."
	case "LastName":
		return "فیلد نام خانوادگی نامعتبر است."
	case "Email":
		return "فیلد ایمیل نامعتبر است."
	default:
		return fmt.Sprintf("فیلد %s نامعتبر است.", field)
	}
}

func getAdminCustomErrorEnglishMessageEn(field string) string {
	switch field {
	case "Sort":
		return "Sort field is invalid."
	case "Role":
		return "Role field is invalid."
	case "Status":
		return "Status field is invalid."
	case "Order":
		return "Order field is invalid."
	case "PhoneNumber":
		return "Phone number must start with 09 and be 11 digits."
	case "FirstName":
		return "First name field is invalid."
	case "LastName":
		return "Last name field is invalid."
	case "Email":
		return "Email field is invalid."
	default:
		return fmt.Sprintf("%s Field is invalid.", field)
	}
}

func ValidateGetUsersRequest(req *dto.AdminGetUsersRequest) error {
	if err := adminValidate.Struct(req); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			field := validationErrs[0].Field()
			return errors.New(
				errors.ValidationError,
				getAdminCustomErrorEnglishMessageEn(field),
				getAdminCustomPersianErrorMessage(field),
				nil,
			)
		}
		return err
	}
	return nil
}

func ValidateUpdateUserRequest(req *dto.AdminUserUpdateRequest) error {
	if err := adminValidate.Struct(req); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			field := validationErrs[0].Field()
			return errors.New(
				errors.ValidationError,
				getAdminCustomErrorEnglishMessageEn(field),
				getAdminCustomPersianErrorMessage(field),
				nil,
			)
		}
		return err
	}
	return nil
}

func ValidateChangeRoleRequest(req *dto.AdminUserUpdateRoleRequest) error {
	if err := adminValidate.Struct(req); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			field := validationErrs[0].Field()
			return errors.New(
				errors.ValidationError,
				getAdminCustomErrorEnglishMessageEn(field),
				getAdminCustomPersianErrorMessage(field),
				nil,
			)
		}
		return err
	}

	roleType := entities.ParseRoleType(req.Role)
	if roleType == entities.SuperAdminRole {
		return errors.New(errors.ValidationError, "cannot change role to super admin", "", nil)
	}

	// Additional check to ensure it's a known role type
	// if !validateRole(req.Role)	 {
	// 	return errors.New(errors.ValidationError, "invalid role type","",nil)
	// }

	return nil
}

func ValidateChangeStatusRequest(req *dto.AdminUserUpdateStatusRequest) error {
	if err := adminValidate.Struct(req); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			field := validationErrs[0].Field()
			return errors.New(
				errors.ValidationError,
				getAdminCustomErrorEnglishMessageEn(field),
				getAdminCustomPersianErrorMessage(field),
				nil,
			)
		}
		return err
	}

	statusType := entities.ParseStatusType(req.Status)
	// Additional check to ensure it's a known status type
	if statusType.String() == "Unknown" {
		return errors.New(errors.ValidationError, "invalid status type", "", nil)
	}

	return nil
}
