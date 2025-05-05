package validators

import (
	"fmt"
	"strings"

	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/amirdashtii/go_auth/internal/core/entities"
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

func ValidateGetUsersRequest(req *dto.AdminGetUsersRequest) error {
	if err := adminValidate.Struct(req); err != nil {
		return err
	}
	return nil
}

func ValidateUpdateUserRequest(req *dto.AdminUserUpdateRequest) error {
	if err := adminValidate.Struct(req); err != nil {
		return err
	}
	return nil
}

func ValidateChangeRoleRequest(req *dto.AdminUserUpdateRoleRequest) error {
	if err := adminValidate.Struct(req); err != nil {
		return err
	}

	roleType := entities.ParseRoleType(req.Role)
	if roleType == entities.SuperAdminRole {
		return fmt.Errorf("cannot change role to super admin")
	}

	// Additional check to ensure it's a known role type
	if roleType.String() == "Unknown" {
		return fmt.Errorf("invalid role type")
	}

	return nil
}

func ValidateChangeStatusRequest(req *dto.AdminUserUpdateStatusRequest) error {
	if err := adminValidate.Struct(req); err != nil {
		return err
	}

	statusType := entities.ParseStatusType(req.Status)
	// Additional check to ensure it's a known status type
	if statusType.String() == "Unknown" {
		return fmt.Errorf("invalid status type")
	}

	return nil
}

func HandleAdminValidationError(err error) map[string]interface{} {
	validationErrors := make(map[string]interface{})
	
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			field := strings.ToLower(e.Field())
			validationErrors[field] = fmt.Sprintf("Invalid %s: %s", field, e.Tag())
		}
	}
	
	return map[string]interface{}{
		"error":   "Validation failed",
		"details": validationErrors,
	}
}