package utils

import (
	"dotapro-lambda-api/errs"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

// WriteResponse writes a JSON response with the given status code.
// The Content-Type header is set to application/json.
func WriteResponse(w http.ResponseWriter, resp any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("error marshaling response: %s", err.Error()), http.StatusInternalServerError)
	}
}

// WriteError writes an error response with the given status code.
// The error is returned as a JSON object with an "error" key.
func WriteError(w http.ResponseWriter, err string, statusCode int) {
	// Only log error responses (non-2xx status codes)
	if statusCode >= 400 {
		log.Error().
			Str("error", err).
			Int("status_code", statusCode).
			Msg("API error response")
	}
	resp := map[string]string{"error": err}
	WriteResponse(w, resp, statusCode)
}

// WriteAppError writes an AppError response with the appropriate status code.
func WriteAppError(w http.ResponseWriter, appErr *errs.AppError) {
	WriteError(w, appErr.Error(), appErr.StatusCode)
}

// HandleError handles any error and writes the appropriate response.
// It handles AppError, context errors, and generic errors.
func HandleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	// Check if it's an AppError
	if appErr, ok := err.(*errs.AppError); ok {
		WriteAppError(w, appErr)
		return
	}

	// Handle context errors
	if err == errs.NOT_FOUND {
		WriteError(w, err.Error(), http.StatusNotFound)
		return
	}

	// Default to internal server error
	WriteError(w, "internal server error", http.StatusInternalServerError)
}
