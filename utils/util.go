package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// show all errors in json not plain text
type ErrorResponse struct {
	Error string `json:"error"`
}

// sendJSONError, show error and status code in json response
func SendJSONError(w http.ResponseWriter, error string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: error})
}

// Check QUERY Parameters
func ValidateQueryParam(r *http.Request, param string, defaultValue int) (int, error) {
	// Query parametresini al
	valueStr := r.URL.Query().Get(param)
	if valueStr == "" {
		return defaultValue, nil
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, err
	}
	return value, nil
}
