package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, resp any, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("error marshaling response: %s", err.Error()), http.StatusInternalServerError)
	}
}

func WriteError(w http.ResponseWriter, err string, statusCode int) {
	resp := map[string]string{"error": err}
	WriteResponse(w, resp, statusCode)
}
