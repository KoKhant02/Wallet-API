package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	MintedSuccess = "Minted successfully."
	MintFailed    = "Minting failed!"
)

// BadRequestError represents an error with a 400 status code (Bad Request)
type BadRequestError struct {
	Message string `json:"message"`
}

// Error implements the error interface for BadRequestError
func (e *BadRequestError) Error() string {
	return fmt.Sprintf("Bad Request: %s", e.Message)
}

// WriteHTTPResponse writes the BadRequestError as an HTTP response with status code 400
func (e *BadRequestError) WriteHTTPResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": e.Message})
}
