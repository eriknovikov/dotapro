package errs

import (
	"errors"
	"fmt"
	"net/http"
)

// Base error types
var (
	// UNIMPLEMENTED is returned for features that are not yet implemented
	UNIMPLEMENTED = errors.New("unimplemented")
	
	// NOT_FOUND is returned when a resource is not found
	NOT_FOUND = errors.New("not found")
)

// AppError represents an application error with HTTP status code
type AppError struct {
	Err        error
	StatusCode int
	Message    string
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return "unknown error"
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new AppError
func NewAppError(err error, statusCode int, message string) *AppError {
	return &AppError{
		Err:        err,
		StatusCode: statusCode,
		Message:    message,
	}
}

// Common error constructors
func NewBadRequestError(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusNotFound,
		Message:    message,
	}
}

func NewInternalServerError(err error, message string) *AppError {
	return &AppError{
		Err:        err,
		StatusCode: http.StatusInternalServerError,
		Message:    message,
	}
}

func NewGatewayTimeoutError(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusGatewayTimeout,
		Message:    message,
	}
}

// WrapError wraps an existing error with additional context
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// IsNotFound checks if an error is a NOT_FOUND error
func IsNotFound(err error) bool {
	return errors.Is(err, NOT_FOUND)
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// GetStatusCode returns the HTTP status code for an error
// Returns 500 (Internal Server Error) for unknown errors
func GetStatusCode(err error) int {
	if appErr, ok := err.(*AppError); ok {
		return appErr.StatusCode
	}
	return http.StatusInternalServerError
}
