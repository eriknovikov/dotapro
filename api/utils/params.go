package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// ParseInt64Param parses an int64 query parameter with the given key.
// Returns nil if the parameter is not present.
// Returns an error if the parameter value is invalid.
func ParseInt64Param(params url.Values, key string) (*int64, error) {
	value := params.Get(key)
	if value == "" {
		return nil, nil
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid %s: %w", key, err)
	}

	return &parsed, nil
}

// ParseIntParam parses an int query parameter with the given key.
// Returns 0 if the parameter is not present.
// Returns an error if the parameter value is invalid.
func ParseIntParam(params url.Values, key string) (int, error) {
	value := params.Get(key)
	if value == "" {
		return 0, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}

	return parsed, nil
}

// ParseRequiredInt64Param parses a required int64 query parameter with the given key.
// Returns an error if the parameter is not present or invalid.
func ParseRequiredInt64Param(params url.Values, key string) (int64, error) {
	value := params.Get(key)
	if value == "" {
		return 0, fmt.Errorf("missing required parameter: %s", key)
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}

	return parsed, nil
}

// ParseStringParam parses a string query parameter with the given key.
// Returns empty string if the parameter is not present.
func ParseStringParam(params url.Values, key string) string {
	return params.Get(key)
}

// ParseRequiredStringParam parses a required string query parameter with the given key.
// Returns an error if the parameter is not present.
func ParseRequiredStringParam(params url.Values, key string) (string, error) {
	value := params.Get(key)
	if value == "" {
		return "", fmt.Errorf("missing required parameter: %s", key)
	}
	return value, nil
}

// WriteParamError writes a parameter parsing error to the response.
func WriteParamError(w http.ResponseWriter, paramName string, err error) {
	http.Error(w, fmt.Sprintf("invalid %s: %v", paramName, err), http.StatusBadRequest)
}
