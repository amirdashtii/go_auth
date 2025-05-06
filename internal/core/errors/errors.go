package errors

import "fmt"

type ErrorType string

const (
	ValidationError ErrorType = "VALIDATION_ERROR"
	AuthenticationError ErrorType = "AUTHENTICATION_ERROR"
	AuthorizationError ErrorType = "AUTHORIZATION_ERROR"
	NotFoundError ErrorType = "NOT_FOUND_ERROR"
	InternalError ErrorType = "INTERNAL_ERROR"
)

type ErrorMessage struct {
	English string
	Persian string
}

type CustomError struct {
	Type    ErrorType
	Message ErrorMessage
	Err     error
}

func (e *CustomError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Type, e.Message.English, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message.English)
}

func (e *CustomError) ErrorPersian() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Type, e.Message.Persian, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message.Persian)
}

func (e *CustomError) Unwrap() error {
	return e.Err
}

func New(errorType ErrorType, messageEn, messageFa string, err error) *CustomError {
	return &CustomError{
		Type: errorType,
		Message: ErrorMessage{
			English: messageEn,
			Persian: messageFa,
		},
		Err: err,
	}
}

func IsValidationError(err error) bool {
	if customErr, ok := err.(*CustomError); ok {
		return customErr.Type == ValidationError
	}
	return false
}

func IsAuthenticationError(err error) bool {
	if customErr, ok := err.(*CustomError); ok {
		return customErr.Type == AuthenticationError
	}
	return false
}

func IsAuthorizationError(err error) bool {
	if customErr, ok := err.(*CustomError); ok {
		return customErr.Type == AuthorizationError
	}
	return false
}

func IsNotFoundError(err error) bool {
	if customErr, ok := err.(*CustomError); ok {
		return customErr.Type == NotFoundError
	}
	return false
}

func IsInternalError(err error) bool {
	if customErr, ok := err.(*CustomError); ok {
		return customErr.Type == InternalError
	}
	return false
}
