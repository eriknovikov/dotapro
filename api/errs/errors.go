package errs

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrUnimplemented = errors.New("unimplemented")
	ErrNotFound      = errors.New("not found")
)

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

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(err error, statusCode int, message string) *AppError {
	return &AppError{
		Err:        err,
		StatusCode: statusCode,
		Message:    message,
	}
}

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
	return errors.Is(err, ErrNotFound)
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
