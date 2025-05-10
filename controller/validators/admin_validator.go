package validators

import (
	"fmt"

	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/errors"
	"github.com/amirdashtii/go_auth/internal/core/ports"
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

	roleType := entities.ParseRoleType(role)
	return roleType.String() != "Unknown"
}

// validateStatus checks if the status is valid without hardcoding status types
func validateStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return true
	}

	statusType := entities.ParseStatusType(status)
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


func getAdminCustomErrorMessage(field string) error {
	switch field {
	case "Sort":
		return errors.ErrInvalidSortField
	case "Role":
		return errors.ErrInvalidRoleField
	case "Status":
		return errors.ErrInvalidStatusField
	case "Order":
		return errors.ErrInvalidOrderField
	case "PhoneNumber":
		return errors.ErrInvalidPhoneNumber
	case "FirstName":
		return errors.ErrInvalidFirstName
	case "LastName":
		return errors.ErrInvalidLastName
	case "Email":
		return errors.ErrInvalidEmail
	case "Password":
		return errors.ErrInvalidPassword
  
	default:
		return errors.New(errors.ValidationError, fmt.Sprintf("%s Field is invalid.", field), fmt.Sprintf("فیلد %s نامعتبر است.", field), nil)
	}
}

func ValidateGetUsersRequest(req *dto.AdminGetUsersRequest, logger ports.Logger) error {
	if err := adminValidate.Struct(req); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			field := validationErrs[0].Field()
			logger.Error("Validation error",
				ports.F("error", err),
				ports.F("field", field),
			)
			return getAdminCustomErrorMessage(field)
		}
		logger.Error("Validation error",
			ports.F("error", err),
		)
		return errors.ErrInvalidRequest
	}
	return nil
}

func ValidateUpdateUserRequest(req *dto.AdminUserUpdateRequest, logger ports.Logger) error {
	if err := adminValidate.Struct(req); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			field := validationErrs[0].Field()
			logger.Error("Validation error",
				ports.F("error", err),
				ports.F("field", field),
			)
			return getAdminCustomErrorMessage(field)
		}
		logger.Error("Validation error",
			ports.F("error", err),
		)
		return errors.ErrInvalidRequest
	}
	return nil
}

func ValidateChangeRoleRequest(req *dto.AdminUserUpdateRoleRequest, logger ports.Logger) error {
	if err := adminValidate.Struct(req); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			field := validationErrs[0].Field()
			logger.Error("Validation error",
				ports.F("error", err),
				ports.F("field", field),
			)
			return getAdminCustomErrorMessage(field)
		}
		logger.Error("Validation error",
			ports.F("error", err),
		)
		return errors.ErrInvalidRequest
	}

	roleType := entities.ParseRoleType(req.Role)
	if roleType == entities.SuperAdminRole {
		logger.Error("Validation error",
			ports.F("error", "cannot change role to super admin"),
		)
		return errors.New(errors.ValidationError, "cannot change role to super admin", "نمی‌توان نقش را به نقش مدیر کل به عنوان مدیر کل تغییر داد.", nil)
	}

	return nil
}

func ValidateChangeStatusRequest(req *dto.AdminUserUpdateStatusRequest, logger ports.Logger) error {
	if err := adminValidate.Struct(req); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			field := validationErrs[0].Field()
			logger.Error("Validation error",
				ports.F("error", err),
				ports.F("field", field),
			)
			return getAdminCustomErrorMessage(field)
		}
		logger.Error("Validation error",
			ports.F("error", err),
		)
		return errors.ErrInvalidRequest
	}

	statusType := entities.ParseStatusType(req.Status)
	if statusType.String() == "Unknown" {
		logger.Error("Validation error",
			ports.F("error", "invalid status type"),
		)
		return errors.New(errors.ValidationError, "invalid status type", "نوع وضعیت نامعتبر است.", nil)
	}

	return nil
}
